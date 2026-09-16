package store

import (
	"context"

	"github.com/z9905080/SNI_WEB/backend/internal/content"
	"github.com/z9905080/SNI_WEB/backend/internal/db/dbgen"
)

func configValue(cfgs []dbgen.WebConfig, key string) string {
	for _, c := range cfgs {
		if c.DataKey == key {
			return c.DataValue
		}
	}
	return ""
}

func (s *Store) Settings(ctx context.Context) (Settings, error) {
	cfgs, err := s.q.ListConfig(ctx)
	if err != nil {
		return Settings{}, err
	}
	return Settings{
		WebTitle:    configValue(cfgs, KeyWebTitle),
		WebSubTitle: configValue(cfgs, KeyWebSubTitle),
		FacebookURL: configValue(cfgs, KeyFacebookURL),
	}, nil
}

func (s *Store) Nav(ctx context.Context) ([]Group, error) {
	cfgs, err := s.q.ListConfig(ctx)
	if err != nil {
		return nil, err
	}
	groups, err := s.q.ListGroups(ctx)
	if err != nil {
		return nil, err
	}
	pages, err := s.q.ListPageSummaries(ctx)
	if err != nil {
		return nil, err
	}

	byGroup := map[int64][]dbgen.ListPageSummariesRow{}
	for _, p := range pages {
		byGroup[p.PageGroupID] = append(byGroup[p.PageGroupID], p)
	}

	groups = content.SortByOrder(groups, content.ParseIDList(configValue(cfgs, KeyGroupSort)),
		func(g dbgen.PageGroup) int64 { return g.ID })

	out := make([]Group, 0, len(groups))
	for _, g := range groups {
		sorted := content.SortByOrder(byGroup[g.ID], content.ParseIDList(g.PageSort),
			func(p dbgen.ListPageSummariesRow) int64 { return p.ID })
		refs := make([]PageRef, 0, len(sorted))
		for _, p := range sorted {
			refs = append(refs, PageRef{ID: p.ID, Name: p.PageName})
		}
		out = append(out, Group{ID: g.ID, Name: g.GroupName, Pages: refs})
	}
	return out, nil
}

func (s *Store) HomePage(ctx context.Context) (Page, bool, error) {
	g, err := s.q.GetGroup(ctx, HomeGroupID)
	if err != nil {
		if notFound(err) == ErrNotFound {
			return Page{}, false, nil
		}
		return Page{}, false, err
	}
	ids, err := s.q.ListPageIDsInGroup(ctx, HomeGroupID)
	if err != nil || len(ids) == 0 {
		return Page{}, false, err
	}
	ordered := content.SortByOrder(ids, content.ParseIDList(g.PageSort), func(id int64) int64 { return id })
	p, err := s.Page(ctx, ordered[0])
	if err != nil {
		return Page{}, false, err
	}
	return p, true, nil
}

func (s *Store) Page(ctx context.Context, id int64) (Page, error) {
	p, err := s.q.GetPage(ctx, id)
	if err != nil {
		return Page{}, notFound(err)
	}
	return toPage(p), nil
}

func toPage(p dbgen.PageContent) Page {
	return Page{ID: p.ID, GroupID: p.PageGroupID, Name: p.PageName, HTML: p.HtmlContext}
}

func (s *Store) Carousels(ctx context.Context) ([]Carousel, error) {
	rows, err := s.q.ListCarousels(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]Carousel, 0, len(rows))
	for _, r := range rows {
		out = append(out, Carousel{ID: r.ID, Image: r.Image, URL: r.Url})
	}
	return out, nil
}

func (s *Store) Marquees(ctx context.Context) ([]Marquee, error) {
	rows, err := s.q.ListMarquees(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]Marquee, 0, len(rows))
	for _, r := range rows {
		out = append(out, Marquee{ID: r.ID, Text: r.Text, Color: r.Color})
	}
	return out, nil
}
