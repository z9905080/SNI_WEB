import { Link } from 'react-router'
import { useDocumentTitle } from './useDocumentTitle'

export default function NotFound() {
  useDocumentTitle('找不到頁面')
  return (
    <div className="mx-auto max-w-[42rem] px-4 py-24">
      <h1 className="font-serif text-3xl font-bold text-brand-dark">找不到這個頁面</h1>
      <p className="mt-4 text-lg text-slate-600">這個頁面可能已經移除，或網址有誤。可以從上方選單找到其他頁面。</p>
      <Link to="/" className="mt-8 inline-block rounded-md bg-brand px-5 py-2.5 text-white hover:bg-brand-dark">
        回到首頁
      </Link>
    </div>
  )
}
