export function wrapTables(root: HTMLElement) {
  for (const table of root.querySelectorAll('table')) {
    if (table.parentElement?.classList.contains('table-scroll')) continue
    const wrapper = document.createElement('div')
    wrapper.className = 'table-scroll'
    table.replaceWith(wrapper)
    wrapper.append(table)
  }
}
