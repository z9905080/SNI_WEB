import { useEffect } from 'react'
import { Route, Routes, useLocation, useParams } from 'react-router'
import { ApiError } from '@/shared/api'
import { publicConfig } from '@/shared/config'
import type { SiteData } from '@/shared/types'
import { initAnalytics, trackPageView } from './analytics'
import Footer from './Footer'
import Header from './Header'
import NotFound from './NotFound'
import PageView from './PageView'
import { useHome, usePage, useSite } from './queries'
import { useDocumentTitle } from './useDocumentTitle'

export default function SiteApp() {
  const { data: site } = useSite()
  const { pathname } = useLocation()

  useEffect(() => initAnalytics(publicConfig.gaMeasurementId), [])
  useEffect(() => {
    window.scrollTo(0, 0)
    trackPageView(pathname)
  }, [pathname])

  return (
    <div className="flex min-h-screen flex-col">
      <Header title={site?.web_title ?? ''} subtitle={site?.web_sub_title ?? ''} menu={site?.menu ?? []} />
      <main className="flex-1">
        <Routes>
          <Route index element={<HomeRoute site={site} />} />
          <Route path="page/:id" element={<PageRoute site={site} />} />
          <Route path="*" element={<NotFound />} />
        </Routes>
      </main>
      <Footer title={site?.web_title ?? ''} facebookUrl={site?.facebook_url ?? ''} counterScriptUrl={publicConfig.counterScriptUrl} />
    </div>
  )
}

function HomeRoute({ site }: { site?: SiteData }) {
  const query = useHome()
  useDocumentTitle(site && (site.web_sub_title ? `${site.web_title}｜${site.web_sub_title}` : site.web_title))
  return <PageView query={query} greeting={site?.web_sub_title} />
}

function PageRoute({ site }: { site?: SiteData }) {
  const { id = '' } = useParams()
  const query = usePage(id)
  const page = query.data?.page
  const groupName = site?.menu.find((g) => g.id === page?.group_id)?.name
  useDocumentTitle(page && site ? `${page.name}｜${site.web_title}` : undefined)
  if (query.error instanceof ApiError && (query.error.status === 404 || query.error.status === 400)) {
    return <NotFound />
  }
  return <PageView query={query} groupName={groupName ?? ''} />
}
