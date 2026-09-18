package store

import (
	"errors"
	"testing"
)

func TestCarouselCRUD(t *testing.T) {
	s, _ := newTestStore(t)
	c, err := s.CreateCarousel(ctx, "/php/picture/a.jpg", "https://x")
	if err != nil || c.ID == 0 {
		t.Fatalf("%+v %v", c, err)
	}
	c.URL = ""
	if err := s.UpdateCarousel(ctx, c); err != nil {
		t.Fatal(err)
	}
	got, err := s.Carousel(ctx, c.ID)
	if err != nil || got != c {
		t.Fatalf("%+v %v", got, err)
	}
	if err := s.UpdateCarousel(ctx, Carousel{ID: 999, Image: "x"}); !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v", err)
	}
	if err := s.DeleteCarousel(ctx, c.ID); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteCarousel(ctx, c.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v", err)
	}
	if _, err := s.Carousel(ctx, c.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v", err)
	}
}

func TestMarqueeCRUD(t *testing.T) {
	s, _ := newTestStore(t)
	m, err := s.CreateMarquee(ctx, "合掌感謝！", "#1EFF00")
	if err != nil || m.ID == 0 {
		t.Fatalf("%+v %v", m, err)
	}
	m.Text = "改"
	if err := s.UpdateMarquee(ctx, m); err != nil {
		t.Fatal(err)
	}
	list, _ := s.Marquees(ctx)
	if len(list) != 1 || list[0] != m {
		t.Fatalf("%+v", list)
	}
	if err := s.DeleteMarquee(ctx, m.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Marquee(ctx, m.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v", err)
	}
}

func TestImageUsages(t *testing.T) {
	s, _ := newTestStore(t,
		`INSERT INTO page_content (id, page_group_id, page_name, html_context) VALUES
		 (1, 1, '用到', '<img src="http://www.seicho-no-ie.org.tw/php/picture/2019-12-25_08-27-33.jpg">'),
		 (2, 1, '相似檔名', '<img src="/php/picture/2019-12-25_08-27-33-abc123.jpg">'),
		 (3, 1, '無關', '<p>2019-12-25_08-27-33.jpg 只是文字</p>')`,
		`INSERT INTO carousel (id, image, url) VALUES (5, '/php/picture/2019-12-25_08-27-33.jpg', ''), (6, '/picture/2019-12-25_08-27-33.jpg', ''), (7, '/php/picture/other.jpg', '')`)
	u, err := s.ImageUsages(ctx, "2019-12-25_08-27-33.jpg")
	if err != nil {
		t.Fatal(err)
	}
	if len(u.Pages) != 1 || u.Pages[0] != (PageRef{ID: 1, Name: "用到"}) {
		t.Fatalf("pages %+v", u.Pages)
	}
	if len(u.Carousels) != 2 || u.Carousels[0].ID != 5 || u.Carousels[1].ID != 6 {
		t.Fatalf("carousels %+v", u.Carousels)
	}

	u, _ = s.ImageUsages(ctx, "unused.jpg")
	if u.Pages == nil || u.Carousels == nil {
		t.Fatalf("應為空切片：%#v", u)
	}
}
