// Package content 放置與資料庫、HTTP 無關的內容處理規則。
package content

import (
	"encoding/json"
	"slices"
	"strings"
)

// ParseIDList 解析 page_sort / page_group_sort 的 JSON 整數陣列；空字串或無效內容視為空陣列。
func ParseIDList(s string) []int64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	var ids []int64
	if err := json.Unmarshal([]byte(s), &ids); err != nil {
		return nil
	}
	return ids
}

func FormatIDList(ids []int64) string {
	if ids == nil {
		ids = []int64{}
	}
	b, _ := json.Marshal(ids)
	return string(b)
}

// SortByOrder 依 order 排列 items；不在 order 中的項目依 id 遞增接在最後，order 中不存在的 id 忽略。
func SortByOrder[T any](items []T, order []int64, id func(T) int64) []T {
	byID := make(map[int64]T, len(items))
	for _, it := range items {
		byID[id(it)] = it
	}
	out := make([]T, 0, len(items))
	used := make(map[int64]bool, len(items))
	for _, oid := range order {
		if it, ok := byID[oid]; ok && !used[oid] {
			out = append(out, it)
			used[oid] = true
		}
	}
	rest := make([]T, 0, len(items)-len(out))
	for _, it := range items {
		if !used[id(it)] {
			rest = append(rest, it)
		}
	}
	slices.SortStableFunc(rest, func(a, b T) int {
		switch x, y := id(a), id(b); {
		case x < y:
			return -1
		case x > y:
			return 1
		}
		return 0
	})
	return append(out, rest...)
}

// SameIDSet 判斷兩個陣列是否為相同集合（含重複次數）。
func SameIDSet(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	x, y := slices.Clone(a), slices.Clone(b)
	slices.Sort(x)
	slices.Sort(y)
	return slices.Equal(x, y)
}

func AppendID(list string, id int64) string {
	ids := slices.DeleteFunc(ParseIDList(list), func(v int64) bool { return v == id })
	return FormatIDList(append(ids, id))
}

func RemoveID(list string, id int64) string {
	return FormatIDList(slices.DeleteFunc(ParseIDList(list), func(v int64) bool { return v == id }))
}
