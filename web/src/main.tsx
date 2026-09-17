import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter } from 'react-router'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import App from './App'
import './index.css'
import { legacyHashTarget } from '@/site/legacyLinks'

// 舊網址 /#/5 → /page/5；在 Router 啟動前改寫，避免多一次導覽
const legacyTarget = legacyHashTarget(window.location.hash)
if (legacyTarget) window.history.replaceState(null, '', legacyTarget)

const queryClient = new QueryClient({
  defaultOptions: { queries: { staleTime: 60_000, retry: 1 } },
})

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <App />
      </BrowserRouter>
    </QueryClientProvider>
  </StrictMode>,
)
