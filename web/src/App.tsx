import { lazy, Suspense } from 'react'
import { Route, Routes } from 'react-router'

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
      <Route path="*" element={<p className="p-8">前台建置中</p>} />
    </Routes>
  )
}
