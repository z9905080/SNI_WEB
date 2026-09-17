// shadcn/ui 目前版本的元件（web/src/components/ui/*）直接從 "cn" 套件匯入 cn()；
// 這裡重新匯出到慣例路徑 @/lib/utils，供之後手寫的元件沿用相同慣例。
export { cn } from 'cn'
