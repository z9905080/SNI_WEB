package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/z9905080/SNI_WEB/backend/internal/db/dbgen"
	"github.com/z9905080/SNI_WEB/backend/internal/testutil"
)

func setup(t *testing.T) (*Sessions, *dbgen.Queries, *time.Time) {
	t.Helper()
	conn, _ := testutil.MySQL(t)
	q := dbgen.New(conn)
	now := time.Date(2026, 9, 16, 12, 0, 0, 0, time.UTC)
	s := NewSessions(q, func() time.Time { return now })
	if _, err := CreateUser(context.Background(), q, "admin", "管理者", "password1"); err != nil {
		t.Fatal(err)
	}
	return s, q, &now
}

func TestLoginAndValidate(t *testing.T) {
	ctx := context.Background()
	s, _, now := setup(t)

	if _, _, _, err := s.Login(ctx, "admin", "wrong-pass"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("錯誤密碼：%v", err)
	}
	if _, _, _, err := s.Login(ctx, "nobody", "password1"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("不存在的帳號：%v", err)
	}

	u, raw, exp, err := s.Login(ctx, "admin", "password1")
	if err != nil {
		t.Fatal(err)
	}
	if u.Account != "admin" || u.Name != "管理者" || u.ID == 0 {
		t.Fatalf("user %+v", u)
	}
	if !exp.Equal(now.Add(time.Hour)) {
		t.Fatalf("expires %v", exp)
	}

	got, renewed, err := s.Validate(ctx, raw)
	if err != nil || got != u || !renewed.IsZero() {
		t.Fatalf("剛登入不需延長：%+v %v %v", got, renewed, err)
	}

	*now = now.Add(31 * time.Minute) // 剩 29 分鐘
	_, renewed, err = s.Validate(ctx, raw)
	if err != nil || !renewed.Equal(now.Add(time.Hour)) {
		t.Fatalf("應延長：%v %v", renewed, err)
	}

	*now = renewed // 恰好到期
	if _, _, err := s.Validate(ctx, raw); !errors.Is(err, ErrInvalidSession) {
		t.Fatalf("過期應拒絕：%v", err)
	}

	if _, _, err := s.Validate(ctx, "garbage"); !errors.Is(err, ErrInvalidSession) {
		t.Fatalf("got %v", err)
	}
}

func TestNewLoginReplacesOldSessionAndLogout(t *testing.T) {
	ctx := context.Background()
	s, _, _ := setup(t)
	_, first, _, _ := s.Login(ctx, "admin", "password1")
	u, second, _, _ := s.Login(ctx, "admin", "password1")
	if _, _, err := s.Validate(ctx, first); !errors.Is(err, ErrInvalidSession) {
		t.Fatalf("舊 session 應失效：%v", err)
	}
	if _, _, err := s.Validate(ctx, second); err != nil {
		t.Fatal(err)
	}
	if err := s.Logout(ctx, u.ID); err != nil {
		t.Fatal(err)
	}
	if _, _, err := s.Validate(ctx, second); !errors.Is(err, ErrInvalidSession) {
		t.Fatalf("登出後應失效：%v", err)
	}
}

func TestTokenStoredHashed(t *testing.T) {
	ctx := context.Background()
	s, q, _ := setup(t)
	u, raw, _, _ := s.Login(ctx, "admin", "password1")
	row, err := q.GetSessionByToken(ctx, HashToken(raw))
	if err != nil || row.UserID != u.ID {
		t.Fatalf("資料庫應存雜湊：%+v %v", row, err)
	}
	if _, err := q.GetSessionByToken(ctx, raw); err == nil {
		t.Fatal("資料庫不應存原始 token")
	}
}

func TestCreateUserAndSetPassword(t *testing.T) {
	ctx := context.Background()
	s, q, _ := setup(t)
	if _, err := CreateUser(ctx, q, "admin", "x", "password2"); !errors.Is(err, ErrAccountExists) {
		t.Fatalf("got %v", err)
	}
	if _, err := CreateUser(ctx, q, "short", "x", "123"); !errors.Is(err, ErrPasswordTooShort) {
		t.Fatalf("got %v", err)
	}
	if err := SetPassword(ctx, q, "admin", "newpassword"); err != nil {
		t.Fatal(err)
	}
	if _, _, _, err := s.Login(ctx, "admin", "newpassword"); err != nil {
		t.Fatalf("新密碼應可登入：%v", err)
	}
	if err := SetPassword(ctx, q, "ghost", "newpassword"); !errors.Is(err, ErrAccountNotFound) {
		t.Fatalf("got %v", err)
	}
}

func TestUserContext(t *testing.T) {
	if _, ok := UserFrom(context.Background()); ok {
		t.Fatal("空 context 不應有 user")
	}
	ctx := WithUser(context.Background(), User{ID: 1, Account: "a"})
	if u, ok := UserFrom(ctx); !ok || u.ID != 1 {
		t.Fatalf("%+v %v", u, ok)
	}
}
