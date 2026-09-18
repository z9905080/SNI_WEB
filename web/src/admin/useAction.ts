import { useMutation, useQueryClient, type QueryKey } from '@tanstack/react-query'
import { toast } from 'sonner'

// 後台多數操作的共同模式：呼叫 API → 成功後重新載入清單 → 失敗時顯示錯誤
export function useAction(invalidate: QueryKey) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (fn: () => Promise<unknown>) => fn(),
    onSuccess: () => qc.invalidateQueries({ queryKey: invalidate }),
    onError: (e) => toast.error(e.message),
  })
}
