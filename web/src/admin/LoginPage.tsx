import { type FormEvent } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useNavigate, useSearchParams } from 'react-router'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useDocumentTitle } from '@/site/useDocumentTitle'
import { adminApi, adminKeys } from './api'
import { safeRedirect } from './auth'

export default function LoginPage() {
  useDocumentTitle('後台登入')
  const [params] = useSearchParams()
  const navigate = useNavigate()
  const qc = useQueryClient()
  const login = useMutation({
    mutationFn: (v: { account: string; password: string }) => adminApi.login(v.account, v.password),
    onSuccess: (user) => {
      qc.setQueryData(adminKeys.me, user)
      navigate(safeRedirect(params.get('redirect')), { replace: true })
    },
  })

  const onSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    const form = new FormData(e.currentTarget)
    login.mutate({ account: String(form.get('account')), password: String(form.get('password')) })
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-muted p-4">
      <Card className="w-full max-w-sm">
        <CardHeader className="items-center text-center">
          <img src="/logo.png" alt="" className="mx-auto h-14" />
          <CardTitle>生長之家後台</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={onSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="account">帳號</Label>
              <Input id="account" name="account" autoComplete="username" required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="password">密碼</Label>
              <Input id="password" name="password" type="password" autoComplete="current-password" required />
            </div>
            {login.isError && (
              <p role="alert" className="text-sm text-destructive">
                {login.error.message}
              </p>
            )}
            <Button type="submit" className="w-full" disabled={login.isPending}>
              {login.isPending ? '登入中…' : '登入'}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
