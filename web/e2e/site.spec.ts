import { expect, test } from '@playwright/test'

test('首頁顯示輪播、跑馬燈與內容', async ({ page }) => {
  await page.goto('/')
  await expect(page).toHaveTitle('生長之家｜一句簡單的感謝')
  await expect(page.getByRole('region', { name: '輪播圖' }).locator('img').first()).toBeVisible()
  await expect(page.getByRole('region', { name: '最新公告' }).getByText('合掌感謝！').first()).toBeVisible()
  await expect(page.getByText('●官方網站重新改版，歡迎瀏覽！')).toBeVisible()
  await expect(page.locator('.legacy-content .table-scroll table')).toHaveCount(1)
  await expect(page.locator('.legacy-content iframe')).toHaveAttribute('src', /youtube\.com\/embed/)
})

test('伺服器注入 SEO meta', async ({ request }) => {
  const html = await (await request.get('/page/3')).text()
  expect(html).toContain('<title>練成會｜生長之家</title>')
  expect(html).toContain('property="og:image" content="http://localhost:8080/php/picture/sample-1.jpg"')
})

test('舊 hash 網址轉到新網址', async ({ page }) => {
  await page.goto('/#/3')
  await expect(page).toHaveURL('/page/3')
  await expect(page.getByRole('heading', { name: '練成會' })).toBeVisible()
  await page.goto('/#/')
  await expect(page).toHaveURL('/')
})

test('內文舊連結以前端路由切換，不重新載入', async ({ page }) => {
  await page.goto('/')
  await page.evaluate(() => Object.assign(window, { __marker: 1 }))
  await page.locator('.legacy-content').getByRole('link', { name: '練成會' }).click()
  await expect(page).toHaveURL('/page/3')
  expect(await page.evaluate(() => (window as unknown as { __marker?: number }).__marker)).toBe(1)
})

test('桌機下拉選單', async ({ page }) => {
  await page.goto('/')
  const nav = page.getByRole('navigation', { name: '主選單' })
  await nav.getByRole('button', { name: '主要行事活動' }).hover()
  await nav.getByRole('link', { name: '誌友會' }).click()
  await expect(page).toHaveURL('/page/4')
  await expect(page).toHaveTitle('誌友會｜生長之家')
})

test('不存在的頁面回 404', async ({ page }) => {
  const res = await page.goto('/page/9999')
  expect(res?.status()).toBe(404)
  await expect(page.getByRole('heading', { name: '找不到這個頁面' })).toBeVisible()
})

test.describe('手機', () => {
  test.use({ viewport: { width: 375, height: 800 } })

  test('沒有橫向捲動，抽屜選單可用', async ({ page }) => {
    for (const path of ['/', '/page/3']) {
      await page.goto(path)
      await expect(page.locator('.legacy-content')).toBeVisible()
      const overflow = await page.evaluate(() => document.documentElement.scrollWidth - window.innerWidth)
      expect(overflow, path).toBeLessThanOrEqual(0)
    }

    await page.getByRole('button', { name: '開啟選單' }).click()
    const drawer = page.getByRole('dialog', { name: '選單' })
    await drawer.getByRole('button', { name: '關於我們' }).click()
    await drawer.getByRole('link', { name: '生長之家簡介' }).click()
    await expect(page).toHaveURL('/page/2')
    await expect(drawer).toBeHidden()
  })
})
