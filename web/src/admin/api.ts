import { api, ApiError } from '@/shared/api'
import type { Carousel, Group, ImageItem, ImageUsages, Marquee, Page, Settings, User } from '@/shared/types'

export type PageInput = { name: string; group_id: number; html: string }
export type CarouselInput = { image: string; url: string }
export type MarqueeInput = { text: string; color: string }

export const adminKeys = {
  all: ['admin'] as const,
  me: ['admin', 'me'] as const,
  groups: ['admin', 'groups'] as const,
  page: (id: number) => ['admin', 'page', id] as const,
  carousels: ['admin', 'carousels'] as const,
  marquees: ['admin', 'marquees'] as const,
  settings: ['admin', 'settings'] as const,
  images: ['admin', 'images'] as const,
}

const A = '/admin'
const post = (body: unknown) => ({ method: 'POST', body })
const patch = (body: unknown) => ({ method: 'PATCH', body })
const put = (body: unknown) => ({ method: 'PUT', body })
const del = { method: 'DELETE' }

export const adminApi = {
  me: () => api<{ user: User }>(`${A}/auth/me`).then((r) => r.user),
  login: (account: string, password: string) =>
    api<{ user: User }>(`${A}/auth/login`, post({ account, password })).then((r) => r.user),
  logout: () => api<void>(`${A}/auth/logout`, { method: 'POST' }),

  groups: () => api<{ groups: Group[] }>(`${A}/groups`).then((r) => r.groups),
  createGroup: (name: string) => api<{ group: Group }>(`${A}/groups`, post({ name })).then((r) => r.group),
  renameGroup: (id: number, name: string) => api<unknown>(`${A}/groups/${id}`, patch({ name })).then(() => {}),
  deleteGroup: (id: number) => api<void>(`${A}/groups/${id}`, del),
  setGroupOrder: (ids: number[]) => api<void>(`${A}/groups/order`, put({ ids })),
  setPageOrder: (groupId: number, ids: number[]) => api<void>(`${A}/groups/${groupId}/pages/order`, put({ ids })),

  page: (id: number) => api<{ page: Page }>(`${A}/pages/${id}`).then((r) => r.page),
  createPage: (p: PageInput) => api<{ page: Page }>(`${A}/pages`, post(p)).then((r) => r.page),
  updatePage: (id: number, p: Partial<PageInput>) => api<{ page: Page }>(`${A}/pages/${id}`, patch(p)).then((r) => r.page),
  deletePage: (id: number) => api<void>(`${A}/pages/${id}`, del),

  carousels: () => api<{ carousels: Carousel[] }>(`${A}/carousels`).then((r) => r.carousels),
  createCarousel: (c: CarouselInput) => api<{ carousel: Carousel }>(`${A}/carousels`, post(c)).then((r) => r.carousel),
  updateCarousel: (id: number, c: CarouselInput) =>
    api<{ carousel: Carousel }>(`${A}/carousels/${id}`, patch(c)).then((r) => r.carousel),
  deleteCarousel: (id: number) => api<void>(`${A}/carousels/${id}`, del),

  marquees: () => api<{ marquees: Marquee[] }>(`${A}/marquees`).then((r) => r.marquees),
  createMarquee: (m: MarqueeInput) => api<{ marquee: Marquee }>(`${A}/marquees`, post(m)).then((r) => r.marquee),
  updateMarquee: (id: number, m: MarqueeInput) =>
    api<{ marquee: Marquee }>(`${A}/marquees/${id}`, patch(m)).then((r) => r.marquee),
  deleteMarquee: (id: number) => api<void>(`${A}/marquees/${id}`, del),

  settings: () => api<{ settings: Settings }>(`${A}/settings`).then((r) => r.settings),
  saveSettings: (s: Settings) => api<{ settings: Settings }>(`${A}/settings`, put(s)).then((r) => r.settings),

  images: () => api<{ images: ImageItem[] }>(`${A}/images`).then((r) => r.images),
  imageUsages: (name: string) =>
    api<{ usages: ImageUsages }>(`${A}/images/${encodeURIComponent(name)}/usages`).then((r) => r.usages),
  deleteImage: (name: string) => api<void>(`${A}/images/${encodeURIComponent(name)}`, del),
}

// fetch 無法回報上傳進度，因此上傳改用 XMLHttpRequest
export function uploadImage(file: File, onProgress: (ratio: number) => void): Promise<ImageItem> {
  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest()
    const body = new FormData()
    body.append('files', file)
    xhr.open('POST', `/api/v1${A}/images`)
    xhr.upload.onprogress = (e) => {
      if (e.lengthComputable) onProgress(e.loaded / e.total)
    }
    xhr.onload = () => {
      let data: { images?: ImageItem[]; error?: { code: string; message: string; request_id?: string } } | null = null
      try {
        data = JSON.parse(xhr.responseText)
      } catch {
        // 非 JSON 回應，下面統一處理
      }
      if (xhr.status >= 200 && xhr.status < 300 && data?.images?.[0]) {
        onProgress(1)
        resolve(data.images[0])
        return
      }
      const err = data?.error
      reject(new ApiError(xhr.status, err?.code ?? 'unknown', err?.message ?? `上傳失敗（${xhr.status}）`, err?.request_id))
    }
    xhr.onerror = () => reject(new ApiError(0, 'network', '無法連線到伺服器，請檢查網路後再試'))
    xhr.send(body)
  })
}
