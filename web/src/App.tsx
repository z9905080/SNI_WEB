import { lazy, Suspense } from 'react'
import { Route, Routes } from 'react-router'
import SiteApp from '@/site/SiteApp'

const AdminApp = lazy(() => import('@/admin/AdminApp'))

export default function App() {
  return (
    <Routes>
      <Route
        path="/admin/*"
        element={
          <Suspense fallback={null}>
            <AdminApp />
          </Suspense>
        }
      />
      <Route path="*" element={<SiteApp />} />
    </Routes>
  )
}
