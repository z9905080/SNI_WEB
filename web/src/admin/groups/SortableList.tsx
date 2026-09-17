import type { ReactNode } from 'react'
import {
  closestCenter,
  DndContext,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
  type DragEndEvent,
} from '@dnd-kit/core'
import { restrictToVerticalAxis } from '@dnd-kit/modifiers'
import {
  arrayMove,
  SortableContext,
  sortableKeyboardCoordinates,
  useSortable,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'

type Props<T> = {
  label: string
  items: T[]
  getId: (item: T) => number
  onReorder: (ids: number[]) => void
  renderItem: (item: T, handle: ReactNode) => ReactNode
}

const announcements = {
  onDragStart: () => '已拿起項目，使用方向鍵移動，空白鍵放下',
  onDragOver: () => '',
  onDragEnd: () => '已放下項目',
  onDragCancel: () => '已取消移動',
}

export default function SortableList<T>({ label, items, getId, onReorder, renderItem }: Props<T>) {
  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 4 } }),
    useSensor(KeyboardSensor, { coordinateGetter: sortableKeyboardCoordinates }),
  )
  const ids = items.map(getId)

  const onDragEnd = ({ active, over }: DragEndEvent) => {
    if (!over || active.id === over.id) return
    onReorder(arrayMove(ids, ids.indexOf(Number(active.id)), ids.indexOf(Number(over.id))))
  }

  return (
    <DndContext
      sensors={sensors}
      collisionDetection={closestCenter}
      modifiers={[restrictToVerticalAxis]}
      onDragEnd={onDragEnd}
      accessibility={{ announcements, screenReaderInstructions: { draggable: '按空白鍵拿起，方向鍵移動，再按空白鍵放下' } }}
    >
      <SortableContext items={ids} strategy={verticalListSortingStrategy}>
        <ul aria-label={label} className="space-y-2">
          {items.map((item) => (
            <SortableRow key={getId(item)} id={getId(item)}>
              {(handle) => renderItem(item, handle)}
            </SortableRow>
          ))}
        </ul>
      </SortableContext>
    </DndContext>
  )
}

function SortableRow({ id, children }: { id: number; children: (handle: ReactNode) => ReactNode }) {
  const { attributes, listeners, setNodeRef, setActivatorNodeRef, transform, transition, isDragging } = useSortable({ id })
  const handle = (
    <button
      type="button"
      ref={setActivatorNodeRef}
      {...attributes}
      {...listeners}
      aria-label="拖曳排序"
      className="cursor-grab touch-none rounded px-1 text-muted-foreground hover:bg-muted active:cursor-grabbing"
    >
      ⋮⋮
    </button>
  )
  return (
    <li
      ref={setNodeRef}
      style={{ transform: CSS.Transform.toString(transform), transition }}
      className={isDragging ? 'relative z-10 opacity-80 shadow-lg' : undefined}
    >
      {children(handle)}
    </li>
  )
}
