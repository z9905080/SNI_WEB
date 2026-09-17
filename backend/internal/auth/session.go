package auth

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"

	"github.com/z9905080/SNI_WEB/backend/internal/db/dbgen"
)

const (
	SessionTTL   = time.Hour
	RefreshBelow = 30 * time.Minute
)

var (
	ErrInvalidCredentials = errors.New("auth: invalid credentials")
	ErrInvalidSession     = errors.New("auth: invalid session")
	ErrAccountExists      = errors.New("帳號已存在")
	ErrAccountNotFound    = errors.New("帳號不存在")
)

type User struct {
	ID      int64  `json:"id"`
	Account string `json:"account"`
	Name    string `json:"name"`
}

type Sessions struct {
	q   *dbgen.Queries
	now func() time.Time
}

func NewSessions(q *dbgen.Queries, now func() time.Time) *Sessions {
	return &Sessions{q: q, now: now}
}

var (
	dummyOnce sync.Once
	dummyHash string
)

// 帳號不存在時仍執行一次 bcrypt，避免從回應時間分辨帳號是否存在。
func burnPasswordCheck(pw string) {
	dummyOnce.Do(func() { dummyHash, _ = HashPassword("dummy-password") })
	CheckPassword(dummyHash, pw)
}

func (s *Sessions) expiry() time.Time {
	return s.now().UTC().Add(SessionTTL).Truncate(time.Second)
}

func (s *Sessions) Login(ctx context.Context, account, password string) (User, string, time.Time, error) {
	row, err := s.q.GetUserByAccount(ctx, account)
	if errors.Is(err, sql.ErrNoRows) {
		burnPasswordCheck(password)
		return User{}, "", time.Time{}, ErrInvalidCredentials
	}
	if err != nil {
		return User{}, "", time.Time{}, err
	}
	if !CheckPassword(row.Pwd, password) {
		return User{}, "", time.Time{}, ErrInvalidCredentials
	}
	raw, hash := NewToken()
	exp := s.expiry()
	if err := s.q.UpsertSession(ctx, dbgen.UpsertSessionParams{UserID: row.ID, Token: hash, ExpireTime: exp}); err != nil {
		return User{}, "", time.Time{}, err
	}
	return User{ID: row.ID, Account: row.Account, Name: row.UserName}, raw, exp, nil
}

func (s *Sessions) Validate(ctx context.Context, raw string) (User, time.Time, error) {
	if raw == "" {
		return User{}, time.Time{}, ErrInvalidSession
	}
	sess, err := s.q.GetSessionByToken(ctx, HashToken(raw))
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, time.Time{}, ErrInvalidSession
	}
	if err != nil {
		return User{}, time.Time{}, err
	}
	now := s.now().UTC()
	if !now.Before(sess.ExpireTime) {
		return User{}, time.Time{}, ErrInvalidSession
	}
	row, err := s.q.GetUserByID(ctx, sess.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, time.Time{}, ErrInvalidSession
	}
	if err != nil {
		return User{}, time.Time{}, err
	}
	var renewed time.Time
	if sess.ExpireTime.Sub(now) < RefreshBelow {
		renewed = s.expiry()
		if err := s.q.ExtendSession(ctx, dbgen.ExtendSessionParams{ExpireTime: renewed, UserID: sess.UserID}); err != nil {
			return User{}, time.Time{}, err
		}
	}
	return User{ID: row.ID, Account: row.Account, Name: row.UserName}, renewed, nil
}

func (s *Sessions) Logout(ctx context.Context, userID int64) error {
	return s.q.DeleteSession(ctx, userID)
}

func CreateUser(ctx context.Context, q *dbgen.Queries, account, name, password string) (int64, error) {
	if _, err := q.GetUserByAccount(ctx, account); err == nil {
		return 0, ErrAccountExists
	} else if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	hash, err := HashPassword(password)
	if err != nil {
		return 0, err
	}
	return q.CreateUser(ctx, dbgen.CreateUserParams{Account: account, Pwd: hash, UserName: name})
}

func SetPassword(ctx context.Context, q *dbgen.Queries, account, password string) error {
	hash, err := HashPassword(password)
	if err != nil {
		return err
	}
	n, err := q.UpdateUserPassword(ctx, dbgen.UpdateUserPasswordParams{Pwd: hash, Account: account})
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrAccountNotFound
	}
	return nil
}

type ctxKey struct{}

func WithUser(ctx context.Context, u User) context.Context {
	return context.WithValue(ctx, ctxKey{}, u)
}

func UserFrom(ctx context.Context) (User, bool) {
	u, ok := ctx.Value(ctxKey{}).(User)
	return u, ok
}
