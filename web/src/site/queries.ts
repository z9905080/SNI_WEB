import { useQuery } from '@tanstack/react-query'
import { api } from '@/shared/api'
import type { PageData, SiteData } from '@/shared/types'

export const useSite = () =>
  useQuery({ queryKey: ['site'], queryFn: ({ signal }) => api<SiteData>('/site', { signal }) })

export const useHome = () =>
  useQuery({ queryKey: ['home'], queryFn: ({ signal }) => api<PageData>('/home', { signal }) })

export const usePage = (id: string) =>
  useQuery({ queryKey: ['page', id], queryFn: ({ signal }) => api<PageData>(`/pages/${id}`, { signal }) })
