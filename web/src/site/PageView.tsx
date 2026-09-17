import type { UseQueryResult } from '@tanstack/react-query'
import type { PageData } from '@/shared/types'

type Props = { query: UseQueryResult<PageData>; greeting?: string; groupName?: string }

export default function PageView({ query, groupName }: Props) {
  if (!query.data) return null
  const { page } = query.data
  return <article className="mx-auto max-w-[42rem] px-4 py-8">{groupName && page && <h1>{page.name}</h1>}</article>
}
