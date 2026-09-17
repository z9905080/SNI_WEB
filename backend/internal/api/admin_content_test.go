package api

import (
	"strings"
	"testing"

	"github.com/z9905080/SNI_WEB/backend/internal/store"
)

const (
	contentFixture = `INSERT INTO page_group (id, group_name, page_sort) VALUES (1, '首頁', ''), (2, '活動', '[11,10]')`
	contentPages   = `INSERT INTO page_content (id, page_group_id, page_name, html_context) VALUES (10, 2, 'A', '<p>a</p>'), (11, 2, 'B', '')`
)

type groupsBody struct {
	Groups []store.Group `json:"groups"`
}
type groupBody struct {
	Group store.Group `json:"group"`
}
type pageOnlyBody struct {
	Page store.Page `json:"page"`
}

func TestAdminRequiresLogin(t *testing.T) {
	h := newHarness(t, contentFixture)
	for _, c := range []struct{ method, path string }{
		{"GET", "/api/v1/admin/groups"},
		{"POST", "/api/v1/admin/groups"},
		{"PUT", "/api/v1/admin/groups/order"},
		{"GET", "/api/v1/admin/pages/10"},
		{"DELETE", "/api/v1/admin/pages/10"},
	} {
		expectError(t, h.do(c.method, c.path, map[string]any{}), 401, "unauthorized")
	}
}

func TestGroupEndpoints(t *testing.T) {
	h := newHarness(t, contentFixture, contentPages)
	h.login()

	list := decode[groupsBody](t, h.do("GET", "/api/v1/admin/groups", nil))
	if len(list.Groups) != 2 || list.Groups[1].Pages[0].ID != 11 {
		t.Fatalf("%+v", list)
	}

	expectError(t, h.do("POST", "/api/v1/admin/groups", map[string]string{"name": "  "}), 400, "validation")
	expectError(t, h.do("POST", "/api/v1/admin/groups", map[string]string{"name": strings.Repeat("字", 51)}), 400, "validation")

	resp := h.do("POST", "/api/v1/admin/groups", map[string]string{"name": " 新群組 "})
	expectStatus(t, resp, 201)
	g := decode[groupBody](t, resp).Group
	if g.Name != "新群組" || g.ID == 0 {
		t.Fatalf("%+v", g)
	}

	renamed := decode[groupBody](t, h.do("PATCH", "/api/v1/admin/groups/"+itoa(g.ID), map[string]string{"name": "改名"}))
	if renamed.Group.Name != "改名" {
		t.Fatalf("%+v", renamed)
	}
	expectError(t, h.do("PATCH", "/api/v1/admin/groups/999", map[string]string{"name": "x"}), 404, "not_found")

	expectError(t, h.do("PUT", "/api/v1/admin/groups/order", map[string]any{"ids": []int64{1}}), 400, "invalid_order")
	expectError(t, h.do("PUT", "/api/v1/admin/groups/order", map[string]any{"ids": []string{"1"}}), 400, "invalid_json")
	expectError(t, h.do("PUT", "/api/v1/admin/groups/order", map[string]any{}), 400, "validation")
	expectStatus(t, h.do("PUT", "/api/v1/admin/groups/order", map[string]any{"ids": []int64{g.ID, 2, 1}}), 204)

	list = decode[groupsBody](t, h.do("GET", "/api/v1/admin/groups", nil))
	if list.Groups[0].ID != g.ID {
		t.Fatalf("排序未生效：%+v", list.Groups)
	}

	expectStatus(t, h.do("PUT", "/api/v1/admin/groups/2/pages/order", map[string]any{"ids": []int64{10, 11}}), 204)
	expectError(t, h.do("PUT", "/api/v1/admin/groups/2/pages/order", map[string]any{"ids": []int64{10}}), 400, "invalid_order")

	expectError(t, h.do("DELETE", "/api/v1/admin/groups/1", nil), 409, "group_protected")
	expectError(t, h.do("DELETE", "/api/v1/admin/groups/2", nil), 409, "group_not_empty")
	expectStatus(t, h.do("DELETE", "/api/v1/admin/groups/"+itoa(g.ID), nil), 204)
	expectError(t, h.do("DELETE", "/api/v1/admin/groups/"+itoa(g.ID), nil), 404, "not_found")
}

func TestGroupNameUnsupportedCharacters(t *testing.T) {
	h := newHarness(t)
	h.login()
	expectError(t, h.do("POST", "/api/v1/admin/groups", map[string]string{"name": "感謝🙏"}), 400, "unsupported_characters")
}

func TestPageEndpoints(t *testing.T) {
	h := newHarness(t, contentFixture, contentPages)
	h.login()

	expectError(t, h.do("POST", "/api/v1/admin/pages", map[string]any{"name": "", "group_id": 2}), 400, "validation")
	expectError(t, h.do("POST", "/api/v1/admin/pages", map[string]any{"name": "x", "group_id": 0}), 400, "validation")
	expectError(t, h.do("POST", "/api/v1/admin/pages", map[string]any{"name": "x", "group_id": 999}), 400, "invalid_group")

	resp := h.do("POST", "/api/v1/admin/pages", map[string]any{"name": "新頁", "group_id": 2, "html": `<p>內容</p><script>x</script>`})
	expectStatus(t, resp, 201)
	p := decode[pageOnlyBody](t, resp).Page
	if p.ID == 0 || p.HTML != "<p>內容</p>" || p.GroupID != 2 {
		t.Fatalf("%+v", p)
	}

	got := decode[pageOnlyBody](t, h.do("GET", "/api/v1/admin/pages/"+itoa(p.ID), nil)).Page
	if got != p {
		t.Fatalf("%+v", got)
	}

	// 只改名稱，其他欄位保留
	upd := decode[pageOnlyBody](t, h.do("PATCH", "/api/v1/admin/pages/"+itoa(p.ID), map[string]any{"name": "改名"})).Page
	if upd.Name != "改名" || upd.HTML != "<p>內容</p>" || upd.GroupID != 2 {
		t.Fatalf("%+v", upd)
	}
	// 換群組
	upd = decode[pageOnlyBody](t, h.do("PATCH", "/api/v1/admin/pages/"+itoa(p.ID), map[string]any{"group_id": 1, "html": "<p>新</p>"})).Page
	if upd.GroupID != 1 || upd.HTML != "<p>新</p>" {
		t.Fatalf("%+v", upd)
	}
	expectError(t, h.do("PATCH", "/api/v1/admin/pages/"+itoa(p.ID), map[string]any{"name": " "}), 400, "validation")
	expectError(t, h.do("PATCH", "/api/v1/admin/pages/999", map[string]any{"name": "x"}), 404, "not_found")
	expectError(t, h.do("PATCH", "/api/v1/admin/pages/"+itoa(p.ID), map[string]any{"html": strings.Repeat("字", 30000)}), 400, "content_too_long")

	expectStatus(t, h.do("DELETE", "/api/v1/admin/pages/"+itoa(p.ID), nil), 204)
	expectError(t, h.do("GET", "/api/v1/admin/pages/"+itoa(p.ID), nil), 404, "not_found")
}
