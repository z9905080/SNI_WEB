import { useEffect, useRef } from 'react'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { useLocation, useNavigate } from 'react-router'
import { ApiError } from '@/shared/api'
import AdminLayout from './AdminLayout'
import { adminApi, adminKeys } from './api'

export const useMe = () => useQuery({ queryKey: adminKeys.me, queryFn: adminApi.me, staleTime: Infinity })

export function safeRedirect(value: string | null) {
  return value && value.startsWith('/admin/') ? value : '/admin/groups'
}

// 任何後台 query / mutation 收到 401 時，清掉後台快取並導向登入頁
export function useRedirectOn401() {
  const qc = useQueryClient()
  const navigate = useNavigate()
  const location = useLocation()
  const here = useRef(location)
  useEffect(() => {
    here.current = location
  })

  useEffect(() => {
    const handle = (error: unknown) => {
      if (!(error instanceof ApiError) || error.status !== 401) return
      const { pathname, search } = here.current
      if (pathname === '/admin/login') return
      qc.removeQueries({ queryKey: adminKeys.all })
      navigate(`/admin/login?redirect=${encodeURIComponent(pathname + search)}`, { replace: true })
    }
    const unsubQueries = qc.getQueryCache().subscribe((e) => {
      if (e.type === 'updated' && e.action.type === 'error') handle(e.action.error)
    })
    const unsubMutations = qc.getMutationCache().subscribe((e) => {
      if (e.type === 'updated' && e.action.type === 'error') handle(e.action.error)
    })
    return () => {
      unsubQueries()
      unsubMutations()
    }
  }, [qc, navigate])
}

export function RequireAuth() {
  const me = useMe()
  if (me.data) return <AdminLayout user={me.data} />
  if (me.isError && !(me.error instanceof ApiError && me.error.status === 401)) {
    return <p role="alert" className="p-8 text-destructive">{me.error.message}</p>
  }
  return <p className="p-8 text-muted-foreground">載入中…</p>
}
