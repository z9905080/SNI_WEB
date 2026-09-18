import path from 'node:path'
import { expect, test, type Page } from '@playwright/test'

const origin = { Origin: 'http://localhost:8080' }

async function login(page: Page, password = 'admin1234') {
  await page.goto('/admin/groups')
  await expect(page).toHaveURL(/\/admin\/login\?redirect=%2Fadmin%2Fgroups/)
  await page.getByLabel('帳號').fill('admin')
  await page.getByLabel('密碼').fill(password)
  const response = page.waitForResponse('**/api/v1/admin/auth/login')
  await page.getByRole('button', { name: '登入' }).click()
  await response
}

test('密碼錯誤時顯示訊息', async ({ page }) => {
  await login(page, 'wrong-password')
  await expect(page.getByRole('alert')).toHaveText('帳號或密碼錯誤')
})

test('新增頁面並在前台看到', async ({ page }) => {
  await login(page)
  await expect(page).toHaveURL('/admin/groups')
  const name = `E2E 頁面 ${Date.now()}`

  await page.getByRole('list', { name: '頁籤' }).getByRole('button', { name: /^關於我們/ }).click()
  await page.getByRole('link', { name: '新增頁面' }).click()
  await page.getByLabel('頁面名稱').fill(name)
  await expect(page.getByLabel('所屬頁籤')).toHaveValue('3')

  await page.getByLabel('頁面內容').click()
  await page.keyboard.type('這是 E2E 建立的內容，')
  await page.getByRole('button', { name: '粗體' }).click()
  await page.keyboard.type('粗體字')
  await page.getByRole('button', { name: '儲存', exact: true }).click()

  await expect(page.getByText('已儲存')).toBeVisible()
  await expect(page).toHaveURL(/\/admin\/pages\/\d+$/)
  const id = page.url().split('/').pop()

  await page.goto(`/page/${id}`)
  await expect(page.getByRole('heading', { name })).toBeVisible()
  await expect(page.locator('.legacy-content strong')).toHaveText('粗體字')
})

test('無法保留格式的頁面以原始碼模式開啟並可直接儲存', async ({ page }) => {
  await login(page)
  const res = await page.request.post('/api/v1/admin/pages', {
    headers: origin,
    data: { name: `E2E 原始碼 ${Date.now()}`, group_id: 3, html: '<div class="box"><p>框內文字</p></div>' },
  })
  expect(res.ok()).toBeTruthy()
  const { page: created } = await res.json()

  await page.goto(`/admin/pages/${created.id}`)
  await expect(page.getByRole('status')).toContainText('HTML 原始碼模式')
  await expect(page.getByRole('button', { name: 'HTML 原始碼' })).toHaveAttribute('aria-pressed', 'true')

  await page.locator('.cm-content').click()
  await page.keyboard.press('ControlOrMeta+End')
  await page.keyboard.insertText('<p>新增段落</p>')
  await page.getByRole('button', { name: '儲存', exact: true }).click()
  await expect(page.getByText('已儲存')).toBeVisible()

  const saved = await (await page.request.get(`/api/v1/pages/${created.id}`)).json()
  expect(saved.page.html).toContain('<div class="box">')
  expect(saved.page.html).toContain('<p>新增段落</p>')

  await page.getByRole('button', { name: '視覺編輯' }).click()
  await expect(page.getByRole('alertdialog')).toContainText('無法在視覺編輯器中保留')
  await page.getByRole('button', { name: '取消' }).click()
  await expect(page.getByRole('button', { name: 'HTML 原始碼' })).toHaveAttribute('aria-pressed', 'true')
})

test('未儲存時離開會提示', async ({ page }) => {
  await login(page)
  await page.goto('/admin/pages/4')
  await page.getByLabel('頁面名稱').fill('誌友會（修改中）')
  await page.getByRole('navigation', { name: '後台選單' }).getByRole('link', { name: '輪播圖' }).click()
  await expect(page.getByRole('alertdialog')).toContainText('還有變更沒有儲存')
  await page.getByRole('button', { name: '取消' }).click()
  await expect(page).toHaveURL('/admin/pages/4')
})

test('以鍵盤排序頁籤並反映到前台', async ({ page }) => {
  await login(page)
  const items = () => page.getByRole('list', { name: '頁籤' }).getByRole('listitem')
  await expect(items()).toHaveText([/首頁/, /關於我們/, /主要行事活動/])

  const dragging = items().filter({ hasText: '關於我們' })
  const handle = dragging.getByRole('button', { name: '拖曳排序' })
  // dnd-kit 的鍵盤感應器在無障礙量測完成前偶爾不會立即回應方向鍵，於headless瀏覽器中屬已知的時序問題；
  // 重試整段「拿起、移動、放下」直到反映到畫面上，不放寬最終斷言。
  await expect(async () => {
    await handle.focus()
    await page.keyboard.press('Space')
    await expect(dragging).toHaveClass(/opacity-80/, { timeout: 500 })
    await page.keyboard.press('ArrowDown')
    await page.keyboard.press('Space')
    await expect(items()).toHaveText([/首頁/, /主要行事活動/, /關於我們/], { timeout: 500 })
  }).toPass({ timeout: 20_000 })

  await expect
    .poll(async () => (await (await page.request.get('/api/v1/site')).json()).menu.map((g: { name: string }) => g.name))
    .toEqual(['首頁', '主要行事活動', '關於我們'])
  await page.reload()
  await expect(items()).toHaveText([/首頁/, /主要行事活動/, /關於我們/])
})

test('上傳圖片後出現在圖庫', async ({ page }) => {
  await login(page)
  await page.getByRole('navigation', { name: '後台選單' }).getByRole('link', { name: '圖庫' }).click()
  await page.getByLabel('選擇圖片').setInputFiles(path.resolve(import.meta.dirname, '../public/logo.png'))
  await expect(page.getByText('已完成')).toBeVisible()
  await expect(page.getByRole('img', { name: /^\d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2}.*\.png$/ }).first()).toBeVisible()
})

test('輪播圖編輯視窗中開啟圖片選擇器並選圖後儲存', async ({ page }) => {
  await login(page)
  await page.goto('/admin/carousels')
  await page.getByRole('button', { name: '新增輪播圖' }).click()

  const dialog = page.getByRole('dialog', { name: '新增輪播圖' })
  await dialog.getByRole('button', { name: '選擇圖片' }).click()

  const picker = page.getByRole('dialog', { name: '選擇圖片' })
  await expect(picker).toBeVisible()
  await picker.getByRole('button', { name: 'sample-1.jpg' }).click()

  await expect(picker).toBeHidden()
  await expect(dialog).toBeVisible()
  await expect(dialog.getByRole('img', { name: '已選擇的圖片' })).toHaveAttribute('src', '/php/picture/sample-1.jpg')

  const url = `https://example.com/e2e-${Date.now()}`
  await dialog.getByLabel('連結網址').fill(url)
  await dialog.getByRole('button', { name: '儲存', exact: true }).click()

  await expect(page.getByText('已儲存輪播圖')).toBeVisible()
  await expect(page.getByText(url)).toBeVisible()
})

test('登出後需要重新登入', async ({ page }) => {
  await login(page)
  await page.getByRole('button', { name: '登出' }).first().click()
  await expect(page).toHaveURL('/admin/login')
  await page.goto('/admin/settings')
  await expect(page).toHaveURL(/\/admin\/login\?redirect=/)
})

test('前台入口 bundle 不含後台程式', async ({ request }) => {
  const html = await (await request.get('/')).text()
  const entry = /<script type="module" crossorigin src="([^"]+)"/.exec(html)?.[1]
  expect(entry).toBeTruthy()
  const js = await (await request.get(entry!)).text()
  for (const marker of ['ProseMirror', 'DndDescribedBy', 'cm-content']) {
    expect(js, marker).not.toContain(marker)
  }
})
