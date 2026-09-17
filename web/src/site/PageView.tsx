import type { UseQueryResult } from '@tanstack/react-query'
import type { PageData } from '@/shared/types'
import Banner from './Banner'
import LegacyContent from './LegacyContent'
import Marquee from './Marquee'

type Props = { query: UseQueryResult<PageData>; greeting?: string; groupName?: string }

export default function PageView({ query, greeting, groupName }: Props) {
  if (query.isPending) return <Skeleton />
  if (query.isError) {
    return (
      <div role="alert" className="mx-auto max-w-[42rem] px-4 py-24">
        <p className="text-lg text-ink">{query.error.message}</p>
        <button
          type="button"
          onClick={() => query.refetch()}
          className="mt-6 rounded-md bg-brand px-5 py-2.5 text-white hover:bg-brand-dark"
        >
          重新載入
        </button>
      </div>
    )
  }

  const { page, carousels, marquees } = query.data
  const isInner = groupName !== undefined
  return (
    <>
      <Banner items={carousels} />
      <Marquee items={marquees} />
      <article className="mx-auto max-w-[42rem] px-4 pb-8 pt-10 sm:pt-16">
        {greeting && (
          // 首頁唯一的強調：以明體大字排出網站副標題
          <p className="mb-12 font-serif text-4xl font-bold leading-snug tracking-wide text-brand-dark sm:text-6xl">
            {greeting}
          </p>
        )}
        {isInner && page && (
          <header className="mb-8">
            {groupName && <p className="text-sm text-slate-500">{groupName}</p>}
            <h1 className="mt-1 font-serif text-3xl font-bold leading-snug text-brand-dark sm:text-4xl">{page.name}</h1>
          </header>
        )}
        {page ? <LegacyContent html={page.html} /> : <p className="text-slate-500">這裡還沒有內容。</p>}
      </article>
    </>
  )
}

function Skeleton() {
  return (
    <div aria-label="載入中" aria-busy="true" className="animate-pulse">
      <div className="mx-auto aspect-[16/7] max-w-6xl bg-brand-soft sm:mt-6 sm:rounded-md" />
      <div className="mx-auto max-w-[42rem] space-y-4 px-4 py-12">
        <div className="h-8 w-1/3 rounded bg-slate-200" />
        <div className="h-4 rounded bg-slate-200" />
        <div className="h-4 w-5/6 rounded bg-slate-200" />
        <div className="h-4 w-2/3 rounded bg-slate-200" />
      </div>
    </div>
  )
}
