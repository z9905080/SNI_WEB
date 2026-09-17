export function applyOrder<T>(items: T[], ids: number[], getId: (item: T) => number): T[] {
  const byId = new Map(items.map((item) => [getId(item), item]))
  return ids.flatMap((id) => byId.get(id) ?? [])
}
