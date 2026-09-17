import { Navigate, Route, Routes } from 'react-router'
import { Toaster } from '@/components/ui/sonner'
import { RequireAuth, useRedirectOn401 } from './auth'
import GroupsPage from './groups/GroupsPage'
import LoginPage from './LoginPage'

const Todo = ({ name }: { name: string }) => <p>{name}（建置中）</p>

export default function AdminApp() {
  useRedirectOn401()
  return (
    <>
      <Routes>
        <Route path="login" element={<LoginPage />} />
        <Route element={<RequireAuth />}>
          <Route index element={<Navigate to="groups" replace />} />
          <Route path="groups" element={<GroupsPage />} />
          <Route path="pages/new" element={<Todo name="新增頁面" />} />
          <Route path="pages/:id" element={<Todo name="編輯頁面" />} />
          <Route path="carousels" element={<Todo name="輪播圖" />} />
          <Route path="marquees" element={<Todo name="跑馬燈" />} />
          <Route path="images" element={<Todo name="圖庫" />} />
          <Route path="settings" element={<Todo name="網站設定" />} />
          <Route path="*" element={<Navigate to="groups" replace />} />
        </Route>
      </Routes>
      <Toaster richColors position="top-center" />
    </>
  )
}
