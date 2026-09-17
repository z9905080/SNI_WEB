package store

import (
	"context"
	"errors"
	"slices"

	"github.com/z9905080/SNI_WEB/backend/internal/content"
	"github.com/z9905080/SNI_WEB/backend/internal/db/dbgen"
)

// updateConfig 以列鎖讀取設定值、套用 fn 後寫回；設定列不存在時新增。
func updateConfig(ctx context.Context, q *dbgen.Queries, key string, fn func(string) string) error {
	row, err := q.GetConfigForUpdate(ctx, key)
	if errors.Is(notFound(err), ErrNotFound) {
		return q.InsertConfig(ctx, dbgen.InsertConfigParams{DataKey: key, DataValue: fn("")})
	}
	if err != nil {
		return err
	}
	return q.UpdateConfig(ctx, dbgen.UpdateConfigParams{DataValue: fn(row.DataValue), ID: row.ID})
}

func setConfig(ctx context.Context, q *dbgen.Queries, key, value string) error {
	return updateConfig(ctx, q, key, func(string) string { return value })
}

func (s *Store) CreateGroup(ctx context.Context, name string) (Group, error) {
	var id int64
	err := s.inTx(ctx, func(q *dbgen.Queries) error {
		var err error
		if id, err = q.CreateGroup(ctx, name); err != nil {
			return err
		}
		return updateConfig(ctx, q, KeyGroupSort, func(old string) string { return content.AppendID(old, id) })
	})
	if err != nil {
		return Group{}, err
	}
	return Group{ID: id, Name: name, Pages: []PageRef{}}, nil
}

func (s *Store) RenameGroup(ctx context.Context, id int64, name string) error {
	n, err := s.q.RenameGroup(ctx, dbgen.RenameGroupParams{GroupName: name, ID: id})
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) DeleteGroup(ctx context.Context, id int64) error {
	if id == HomeGroupID {
		return ErrGroupProtected
	}
	return s.inTx(ctx, func(q *dbgen.Queries) error {
		if _, err := q.GetGroupForUpdate(ctx, id); err != nil {
			return notFound(err)
		}
		n, err := q.CountPagesInGroup(ctx, id)
		if err != nil {
			return err
		}
		if n > 0 {
			return ErrGroupNotEmpty
		}
		if err := q.DeleteGroup(ctx, id); err != nil {
			return err
		}
		return updateConfig(ctx, q, KeyGroupSort, func(old string) string { return content.RemoveID(old, id) })
	})
}

func (s *Store) SetGroupOrder(ctx context.Context, ids []int64) error {
	return s.inTx(ctx, func(q *dbgen.Queries) error {
		groups, err := q.ListGroups(ctx)
		if err != nil {
			return err
		}
		existing := make([]int64, 0, len(groups))
		for _, g := range groups {
			existing = append(existing, g.ID)
		}
		if !content.SameIDSet(ids, existing) {
			return ErrOrderMismatch
		}
		return setConfig(ctx, q, KeyGroupSort, content.FormatIDList(ids))
	})
}

func (s *Store) SetPageOrder(ctx context.Context, groupID int64, ids []int64) error {
	return s.inTx(ctx, func(q *dbgen.Queries) error {
		if _, err := q.GetGroupForUpdate(ctx, groupID); err != nil {
			return notFound(err)
		}
		existing, err := q.ListPageIDsInGroup(ctx, groupID)
		if err != nil {
			return err
		}
		if !content.SameIDSet(ids, existing) {
			return ErrOrderMismatch
		}
		return q.UpdateGroupPageSort(ctx, dbgen.UpdateGroupPageSortParams{PageSort: content.FormatIDList(ids), ID: groupID})
	})
}

func cleanHTML(html string) (string, error) {
	out := content.SanitizeHTML(html)
	if len(out) > MaxHTMLBytes {
		return "", ErrContentTooLong
	}
	return out, nil
}

func (s *Store) CreatePage(ctx context.Context, p Page) (Page, error) {
	html, err := cleanHTML(p.HTML)
	if err != nil {
		return Page{}, err
	}
	p.HTML = html
	err = s.inTx(ctx, func(q *dbgen.Queries) error {
		g, err := q.GetGroupForUpdate(ctx, p.GroupID)
		if err != nil {
			if errors.Is(notFound(err), ErrNotFound) {
				return ErrInvalidGroup
			}
			return err
		}
		if p.ID, err = q.CreatePage(ctx, dbgen.CreatePageParams{PageGroupID: p.GroupID, PageName: p.Name, HtmlContext: p.HTML}); err != nil {
			return err
		}
		return q.UpdateGroupPageSort(ctx, dbgen.UpdateGroupPageSortParams{PageSort: content.AppendID(g.PageSort, p.ID), ID: g.ID})
	})
	if err != nil {
		return Page{}, err
	}
	return p, nil
}

func (s *Store) UpdatePage(ctx context.Context, p Page) (Page, error) {
	html, err := cleanHTML(p.HTML)
	if err != nil {
		return Page{}, err
	}
	p.HTML = html
	err = s.inTx(ctx, func(q *dbgen.Queries) error {
		old, err := q.GetPageForUpdate(ctx, p.ID)
		if err != nil {
			return notFound(err)
		}
		if old.PageGroupID != p.GroupID {
			if err := moveBetweenGroups(ctx, q, p.ID, old.PageGroupID, p.GroupID); err != nil {
				return err
			}
		}
		return q.UpdatePage(ctx, dbgen.UpdatePageParams{PageGroupID: p.GroupID, PageName: p.Name, HtmlContext: p.HTML, ID: p.ID})
	})
	if err != nil {
		return Page{}, err
	}
	return p, nil
}

// moveBetweenGroups 依 id 遞增順序鎖定兩個群組以避免死結，再更新雙方的 page_sort。
func moveBetweenGroups(ctx context.Context, q *dbgen.Queries, pageID, from, to int64) error {
	ids := []int64{from, to}
	slices.Sort(ids)
	locked := map[int64]dbgen.PageGroup{}
	for _, gid := range ids {
		g, err := q.GetGroupForUpdate(ctx, gid)
		if errors.Is(notFound(err), ErrNotFound) {
			continue
		}
		if err != nil {
			return err
		}
		locked[gid] = g
	}
	dst, ok := locked[to]
	if !ok {
		return ErrInvalidGroup
	}
	if src, ok := locked[from]; ok {
		if err := q.UpdateGroupPageSort(ctx, dbgen.UpdateGroupPageSortParams{PageSort: content.RemoveID(src.PageSort, pageID), ID: src.ID}); err != nil {
			return err
		}
	}
	return q.UpdateGroupPageSort(ctx, dbgen.UpdateGroupPageSortParams{PageSort: content.AppendID(dst.PageSort, pageID), ID: dst.ID})
}

func (s *Store) DeletePage(ctx context.Context, id int64) error {
	return s.inTx(ctx, func(q *dbgen.Queries) error {
		old, err := q.GetPageForUpdate(ctx, id)
		if err != nil {
			return notFound(err)
		}
		if err := q.DeletePage(ctx, id); err != nil {
			return err
		}
		g, err := q.GetGroupForUpdate(ctx, old.PageGroupID)
		if errors.Is(notFound(err), ErrNotFound) {
			return nil
		}
		if err != nil {
			return err
		}
		return q.UpdateGroupPageSort(ctx, dbgen.UpdateGroupPageSortParams{PageSort: content.RemoveID(g.PageSort, id), ID: g.ID})
	})
}

func (s *Store) UpdateSettings(ctx context.Context, v Settings) error {
	return s.inTx(ctx, func(q *dbgen.Queries) error {
		for key, val := range map[string]string{
			KeyWebTitle:    v.WebTitle,
			KeyWebSubTitle: v.WebSubTitle,
			KeyFacebookURL: v.FacebookURL,
		} {
			if err := setConfig(ctx, q, key, val); err != nil {
				return err
			}
		}
		return nil
	})
}
