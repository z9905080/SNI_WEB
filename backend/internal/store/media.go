package store

import (
	"context"

	"github.com/z9905080/SNI_WEB/backend/internal/db/dbgen"
)

func (s *Store) Carousel(ctx context.Context, id int64) (Carousel, error) {
	r, err := s.q.GetCarousel(ctx, id)
	if err != nil {
		return Carousel{}, notFound(err)
	}
	return Carousel{ID: r.ID, Image: r.Image, URL: r.Url}, nil
}

func (s *Store) CreateCarousel(ctx context.Context, image, url string) (Carousel, error) {
	id, err := s.q.CreateCarousel(ctx, dbgen.CreateCarouselParams{Image: image, Url: url})
	if err != nil {
		return Carousel{}, err
	}
	return Carousel{ID: id, Image: image, URL: url}, nil
}

func (s *Store) UpdateCarousel(ctx context.Context, c Carousel) error {
	return affected(s.q.UpdateCarousel(ctx, dbgen.UpdateCarouselParams{Image: c.Image, Url: c.URL, ID: c.ID}))
}

func (s *Store) DeleteCarousel(ctx context.Context, id int64) error {
	return affected(s.q.DeleteCarousel(ctx, id))
}

func (s *Store) Marquee(ctx context.Context, id int64) (Marquee, error) {
	r, err := s.q.GetMarquee(ctx, id)
	if err != nil {
		return Marquee{}, notFound(err)
	}
	return Marquee{ID: r.ID, Text: r.Text, Color: r.Color}, nil
}

func (s *Store) CreateMarquee(ctx context.Context, text, color string) (Marquee, error) {
	id, err := s.q.CreateMarquee(ctx, dbgen.CreateMarqueeParams{Text: text, Color: color})
	if err != nil {
		return Marquee{}, err
	}
	return Marquee{ID: id, Text: text, Color: color}, nil
}

func (s *Store) UpdateMarquee(ctx context.Context, m Marquee) error {
	return affected(s.q.UpdateMarquee(ctx, dbgen.UpdateMarqueeParams{Text: m.Text, Color: m.Color, ID: m.ID}))
}

func (s *Store) DeleteMarquee(ctx context.Context, id int64) error {
	return affected(s.q.DeleteMarquee(ctx, id))
}

func affected(n int64, err error) error {
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

type ImageUsages struct {
	Pages     []PageRef  `json:"pages"`
	Carousels []Carousel `json:"carousels"`
}

// ImageUsages 找出內容或輪播圖中引用該檔名的項目（以 "picture/檔名" 比對，避免誤判純文字）。
func (s *Store) ImageUsages(ctx context.Context, name string) (ImageUsages, error) {
	needle := "picture/" + name
	pages, err := s.q.FindPagesUsingText(ctx, needle)
	if err != nil {
		return ImageUsages{}, err
	}
	carousels, err := s.q.FindCarouselsUsingText(ctx, needle)
	if err != nil {
		return ImageUsages{}, err
	}
	u := ImageUsages{Pages: make([]PageRef, 0, len(pages)), Carousels: make([]Carousel, 0, len(carousels))}
	for _, p := range pages {
		u.Pages = append(u.Pages, PageRef{ID: p.ID, Name: p.PageName})
	}
	for _, c := range carousels {
		u.Carousels = append(u.Carousels, Carousel{ID: c.ID, Image: c.Image, URL: c.Url})
	}
	return u, nil
}
