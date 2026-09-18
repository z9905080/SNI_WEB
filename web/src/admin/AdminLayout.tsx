import { useMutation, useQueryClient } from '@tanstack/react-query'
import { NavLink, Outlet, useNavigate } from 'react-router'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import type { User } from '@/shared/types'
import { adminApi, adminKeys } from './api'

const links = [
  { to: '/admin/groups', label: '頁籤與頁面' },
  { to: '/admin/carousels', label: '輪播圖' },
  { to: '/admin/marquees', label: '跑馬燈' },
  { to: '/admin/images', label: '圖庫' },
  { to: '/admin/settings', label: '網站設定' },
]

export default function AdminLayout({ user }: { user: User }) {
  const qc = useQueryClient()
  const navigate = useNavigate()
  const logout = useMutation({
    mutationFn: adminApi.logout,
    onSuccess: () => {
      qc.removeQueries({ queryKey: adminKeys.all })
      navigate('/admin/login', { replace: true })
    },
    onError: (e) => toast.error(e.message),
  })

  return (
    <div className="min-h-screen bg-muted/40 md:flex">
      <aside className="border-b bg-background md:sticky md:top-0 md:h-screen md:w-56 md:shrink-0 md:border-b-0 md:border-r">
        <div className="flex items-center gap-2 p-4">
          <img src="/logo.png" alt="" className="h-8" />
          <span className="font-bold">後台管理</span>
        </div>
        <nav aria-label="後台選單" className="flex gap-1 overflow-x-auto px-2 pb-2 md:flex-col md:pb-0">
          {links.map((l) => (
            <NavLink
              key={l.to}
              to={l.to}
              className={({ isActive }) =>
                `whitespace-nowrap rounded-md px-3 py-2 text-sm ${isActive ? 'bg-primary text-primary-foreground' : 'hover:bg-muted'}`
              }
            >
              {l.label}
            </NavLink>
          ))}
        </nav>
        <div className="hidden space-y-2 border-t p-4 text-sm md:absolute md:bottom-0 md:block md:w-full">
          <p className="text-muted-foreground">{user.name}</p>
          <div className="flex gap-2">
            <Button variant="outline" size="sm" asChild>
              <a href="/" target="_blank" rel="noreferrer">查看網站</a>
            </Button>
            <Button variant="ghost" size="sm" onClick={() => logout.mutate()} disabled={logout.isPending}>
              登出
            </Button>
          </div>
        </div>
      </aside>
      <main className="min-w-0 flex-1 p-4 md:p-8">
        <div className="mb-4 flex items-center justify-end gap-2 text-sm md:hidden">
          <span className="text-muted-foreground">{user.name}</span>
          <Button variant="ghost" size="sm" onClick={() => logout.mutate()}>
            登出
          </Button>
        </div>
        <Outlet />
      </main>
    </div>
  )
}
