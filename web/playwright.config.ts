import { defineConfig, devices } from '@playwright/test'

// 需先 make up；make e2e 會建置前端、重設資料庫並放入範例圖片
export default defineConfig({
  testDir: './e2e',
  workers: 1,
  use: { baseURL: 'http://localhost:8080', trace: 'retain-on-failure' },
  projects: [{ name: 'chromium', use: { ...devices['Desktop Chrome'] } }],
  webServer: {
    command: 'go run ./backend/cmd/sniweb serve',
    cwd: '..',
    url: 'http://localhost:8080/readyz',
    reuseExistingServer: false,
    timeout: 120_000,
    env: {
      PORT: '8080',
      DATABASE_URL: 'root:sniweb@tcp(127.0.0.1:3306)/sniweb',
      PUBLIC_BASE_URL: 'http://localhost:8080',
      COOKIE_SECURE: 'false',
      STORAGE_DRIVER: 'disk',
      STORAGE_DISK_DIR: './tmp/picture',
    },
  },
})
