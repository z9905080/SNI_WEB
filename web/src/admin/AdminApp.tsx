import { Navigate, Route, Routes } from 'react-router'
import { Toaster } from '@/components/ui/sonner'
import { RequireAuth, useRedirectOn401 } from './auth'
import CarouselsPage from './carousels/CarouselsPage'
import PageEditorPage from './editor/PageEditorPage'
import GroupsPage from './groups/GroupsPage'
import ImagesPage from './images/ImagesPage'
import LoginPage from './LoginPage'
import MarqueesPage from './marquees/MarqueesPage'
import SettingsPage from './settings/SettingsPage'

export default function AdminApp() {
  useRedirectOn401()
  return (
    <>
      <Routes>
        <Route path="login" element={<LoginPage />} />
        <Route element={<RequireAuth />}>
          <Route index element={<Navigate to="groups" replace />} />
          <Route path="groups" element={<GroupsPage />} />
          <Route path="pages/new" element={<PageEditorPage />} />
          <Route path="pages/:id" element={<PageEditorPage />} />
          <Route path="carousels" element={<CarouselsPage />} />
          <Route path="marquees" element={<MarqueesPage />} />
          <Route path="images" element={<ImagesPage />} />
          <Route path="settings" element={<SettingsPage />} />
          <Route path="*" element={<Navigate to="groups" replace />} />
        </Route>
      </Routes>
      <Toaster richColors position="top-center" />
    </>
  )
}
