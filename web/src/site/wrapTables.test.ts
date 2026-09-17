import { expect, it } from 'vitest'
import { wrapTables } from './wrapTables'

it('每個表格包進可橫向捲動的容器，重複呼叫不重複包', () => {
  const root = document.createElement('div')
  root.innerHTML = '<p>x</p><table id="a"><tr><td><table id="b"><tr><td>1</td></tr></table></td></tr></table>'
  wrapTables(root)
  wrapTables(root)
  for (const id of ['a', 'b']) {
    const table = root.querySelector(`#${id}`)!
    expect(table.parentElement?.className).toBe('table-scroll')
    expect(table.parentElement?.parentElement?.className).not.toBe('table-scroll')
  }
  expect(root.querySelectorAll('.table-scroll')).toHaveLength(2)
})
