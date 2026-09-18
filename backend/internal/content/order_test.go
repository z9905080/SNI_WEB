package content

import (
	"slices"
	"testing"
)

func TestParseIDList(t *testing.T) {
	cases := map[string][]int64{
		"":            nil,
		"   ":         nil,
		"null":        nil,
		"not json":    nil,
		`["1","2"]`:   nil,
		"[1.5]":       nil,
		"[]":          {},
		"[5,20,38]":   {5, 20, 38},
		" [3, 1, 2] ": {3, 1, 2},
	}
	for in, want := range cases {
		if got := ParseIDList(in); !slices.Equal(got, want) {
			t.Errorf("ParseIDList(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestFormatIDList(t *testing.T) {
	if got := FormatIDList(nil); got != "[]" {
		t.Errorf("nil → %q", got)
	}
	if got := FormatIDList([]int64{3, 1}); got != "[3,1]" {
		t.Errorf("got %q", got)
	}
}

func TestSortByOrder(t *testing.T) {
	type item struct{ ID int64 }
	id := func(i item) int64 { return i.ID }
	ids := func(items []item) []int64 {
		var out []int64
		for _, i := range items {
			out = append(out, i.ID)
		}
		return out
	}

	items := []item{{7}, {2}, {5}, {9}, {1}}
	// 陣列中的依序；9、1 不在陣列中 → 依 id 遞增接在後面；99 已不存在 → 忽略；重複的 5 只算一次
	got := ids(SortByOrder(items, []int64{5, 99, 7, 5, 2}, id))
	if want := []int64{5, 7, 2, 1, 9}; !slices.Equal(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}

	got = ids(SortByOrder(items, nil, id))
	if want := []int64{1, 2, 5, 7, 9}; !slices.Equal(got, want) {
		t.Fatalf("空排序應依 id 遞增：got %v", got)
	}

	if out := SortByOrder([]item{}, []int64{1}, id); out == nil || len(out) != 0 {
		t.Fatalf("空輸入應回傳非 nil 空切片：%#v", out)
	}
}

func TestSortByOrderDoesNotMutateInput(t *testing.T) {
	in := []int64{3, 1, 2}
	SortByOrder(in, nil, func(v int64) int64 { return v })
	if !slices.Equal(in, []int64{3, 1, 2}) {
		t.Fatalf("輸入被修改：%v", in)
	}
}

func TestSameIDSet(t *testing.T) {
	cases := []struct {
		a, b []int64
		want bool
	}{
		{[]int64{1, 2, 3}, []int64{3, 1, 2}, true},
		{nil, []int64{}, true},
		{[]int64{1, 2}, []int64{1, 2, 3}, false},
		{[]int64{1, 1}, []int64{1, 2}, false},
		{[]int64{1, 4}, []int64{1, 2}, false},
	}
	for _, c := range cases {
		if got := SameIDSet(c.a, c.b); got != c.want {
			t.Errorf("SameIDSet(%v,%v)=%v", c.a, c.b, got)
		}
	}
}

func TestAppendAndRemoveID(t *testing.T) {
	if got := AppendID("", 4); got != "[4]" {
		t.Errorf("AppendID 空字串 → %q", got)
	}
	if got := AppendID("[1,4,2]", 4); got != "[1,2,4]" {
		t.Errorf("AppendID 既有 id 應移到尾端 → %q", got)
	}
	if got := AppendID("garbage", 3); got != "[3]" {
		t.Errorf("AppendID 無效 JSON → %q", got)
	}
	if got := RemoveID("[1,4,2,4]", 4); got != "[1,2]" {
		t.Errorf("RemoveID → %q", got)
	}
	if got := RemoveID("", 4); got != "[]" {
		t.Errorf("RemoveID 空字串 → %q", got)
	}
}
