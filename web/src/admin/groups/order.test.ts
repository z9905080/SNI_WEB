import { expect, it } from 'vitest'
import { applyOrder } from './order'

it('依 id 陣列重新排列，忽略不存在的 id', () => {
  const items = [{ id: 1 }, { id: 2 }, { id: 3 }]
  expect(applyOrder(items, [3, 9, 1, 2], (i) => i.id)).toEqual([{ id: 3 }, { id: 1 }, { id: 2 }])
})
