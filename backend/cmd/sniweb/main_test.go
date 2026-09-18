package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/z9905080/SNI_WEB/backend/internal/auth"
	"github.com/z9905080/SNI_WEB/backend/internal/db/dbgen"
	"github.com/z9905080/SNI_WEB/backend/internal/testutil"
)

func runCmd(t *testing.T, env map[string]string, stdin string, args ...string) (int, string, string) {
	t.Helper()
	var out, errOut bytes.Buffer
	code := run(context.Background(), args, strings.NewReader(stdin), &out, &errOut, func(k string) string { return env[k] })
	return code, out.String(), errOut.String()
}

func TestRunUsage(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"沒有參數", nil},
		{"未知指令", []string{"nope"}},
		{"user 沒有子指令", []string{"user"}},
		{"user 未知子指令", []string{"user", "delete", "--account", "a"}},
		{"create 缺 name", []string{"user", "create", "--account", "a", "--password-stdin"}},
		{"passwd 缺 account", []string{"user", "passwd", "--password-stdin"}},
		{"未知 flag", []string{"user", "passwd", "--account", "a", "--bogus"}},
		{"migrate 多餘參數", []string{"migrate", "up"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, errOut := runCmd(t, nil, "", tt.args...)
			if code != 2 || !strings.Contains(errOut, "用法") {
				t.Fatalf("code=%d stderr=%q", code, errOut)
			}
		})
	}
}

func TestServeRequiresConfig(t *testing.T) {
	code, _, errOut := runCmd(t, nil, "", "serve")
	if code != 1 || !strings.Contains(errOut, "DATABASE_URL") {
		t.Fatalf("code=%d stderr=%q", code, errOut)
	}
}

// migrate 只該要求 DATABASE_URL，不必湊齊 serve 用的其他設定
func TestMigrateRequiresOnlyDatabaseURL(t *testing.T) {
	code, _, errOut := runCmd(t, nil, "", "migrate")
	if code != 1 || !strings.Contains(errOut, "DATABASE_URL") {
		t.Fatalf("code=%d stderr=%q", code, errOut)
	}
	if strings.Contains(errOut, "PUBLIC_BASE_URL") {
		t.Fatalf("migrate 不該要求 PUBLIC_BASE_URL：%q", errOut)
	}
}

func TestUserCommands(t *testing.T) {
	conn, dsn := testutil.MySQL(t)
	env := map[string]string{"DATABASE_URL": dsn}
	sessions := auth.NewSessions(dbgen.New(conn), time.Now)
	ctx := context.Background()

	code, out, errOut := runCmd(t, env, "password1\n", "user", "create", "--account", "bob", "--name", "鮑伯", "--password-stdin")
	if code != 0 || !strings.Contains(out, "bob") {
		t.Fatalf("create: code=%d out=%q err=%q", code, out, errOut)
	}
	if _, _, _, err := sessions.Login(ctx, "bob", "password1"); err != nil {
		t.Fatalf("新帳號應可登入：%v", err)
	}

	if code, _, errOut := runCmd(t, env, "password1\n", "user", "create", "--account", "bob", "--name", "x", "--password-stdin"); code != 1 || !strings.Contains(errOut, "帳號已存在") {
		t.Fatalf("重複帳號：code=%d err=%q", code, errOut)
	}
	if code, _, errOut := runCmd(t, env, "short\n", "user", "create", "--account", "amy", "--name", "x", "--password-stdin"); code != 1 || !strings.Contains(errOut, "8 個字元") {
		t.Fatalf("密碼太短：code=%d err=%q", code, errOut)
	}

	if code, _, errOut := runCmd(t, env, "newpassword", "user", "passwd", "--account", "bob", "--password-stdin"); code != 0 {
		t.Fatalf("passwd: code=%d err=%q", code, errOut)
	}
	if _, _, _, err := sessions.Login(ctx, "bob", "newpassword"); err != nil {
		t.Fatalf("新密碼應可登入：%v", err)
	}
	if code, _, errOut := runCmd(t, env, "newpassword\n", "user", "passwd", "--account", "ghost", "--password-stdin"); code != 1 || !strings.Contains(errOut, "帳號不存在") {
		t.Fatalf("不存在的帳號：code=%d err=%q", code, errOut)
	}

	// 非終端機又沒有 --password-stdin
	if code, _, errOut := runCmd(t, env, "x\n", "user", "passwd", "--account", "bob"); code != 1 || !strings.Contains(errOut, "--password-stdin") {
		t.Fatalf("code=%d err=%q", code, errOut)
	}
}
