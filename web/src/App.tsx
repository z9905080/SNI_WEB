import { lazy, Suspense } from 'react'
import { Route, Routes } from 'react-router'
import SiteApp from '@/site/SiteApp'

const AdminApp = lazy(() => import('@/admin/AdminApp'))

export default function App() {
  return (
    <Routes>
      {/* 後台 chunk 有 1MB 以上，載入慢時 fallback 若為 null 會是一片白，看起來像壞掉 */}
      <Route
        path="/admin/*"
        element={
          <Suspense
            fallback={
              <p role="status" className="flex min-h-screen items-center justify-center text-muted-foreground">
                載入中…
              </p>
            }
          >
            <AdminApp />
          </Suspense>
        }
      />
      <Route path="*" element={<SiteApp />} />
    </Routes>
  )
}
