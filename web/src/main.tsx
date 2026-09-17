import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { createBrowserRouter, RouterProvider } from 'react-router'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { retryUnlessClientError } from '@/shared/api'
import { legacyHashTarget } from '@/site/legacyLinks'
import App from './App'
import './index.css'

// 舊網址 /#/5 → /page/5；在 Router 啟動前改寫，避免多一次導覽
const legacyTarget = legacyHashTarget(window.location.hash)
if (legacyTarget) window.history.replaceState(null, '', legacyTarget)

const queryClient = new QueryClient({
  defaultOptions: { queries: { staleTime: 60_000, retry: retryUnlessClientError } },
})

// data router 才能使用 useBlocker（編輯頁未存檔提示）
const router = createBrowserRouter([{ path: '*', element: <App /> }])

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>
  </StrictMode>,
)
