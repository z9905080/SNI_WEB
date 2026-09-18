package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/z9905080/SNI_WEB/backend/internal/auth"
	"github.com/z9905080/SNI_WEB/backend/internal/db"
	"github.com/z9905080/SNI_WEB/backend/internal/db/dbgen"
)

func userCmd(ctx context.Context, sub string, args []string, stdin io.Reader, stdout, stderr io.Writer, getenv func(string) string) error {
	if sub != "create" && sub != "passwd" {
		return errUsage
	}
	flags := flag.NewFlagSet("user "+sub, flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	account := flags.String("account", "", "")
	name := flags.String("name", "", "")
	fromStdin := flags.Bool("password-stdin", false, "")
	if err := flags.Parse(args); err != nil || *account == "" || (sub == "create" && *name == "") {
		return errUsage
	}

	password, err := readPassword(stdin, stderr, *fromStdin)
	if err != nil {
		return err
	}
	dsn := getenv("DATABASE_URL")
	if dsn == "" {
		return errors.New("缺少環境變數：DATABASE_URL")
	}
	conn, err := db.Open(ctx, dsn)
	if err != nil {
		return err
	}
	defer conn.Close()
	q := dbgen.New(conn)

	if sub == "create" {
		id, err := auth.CreateUser(ctx, q, *account, *name, password)
		if err != nil {
			return err
		}
		fmt.Fprintf(stdout, "已建立帳號 %s（id=%d）\n", *account, id)
		return nil
	}
	if err := auth.SetPassword(ctx, q, *account, password); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "已更新 %s 的密碼\n", *account)
	return nil
}

// readPassword 從 stdin 讀一行，或在終端機上不回顯地輸入兩次。
func readPassword(stdin io.Reader, prompt io.Writer, fromStdin bool) (string, error) {
	if fromStdin {
		line, err := bufio.NewReader(stdin).ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return "", err
		}
		return strings.TrimRight(line, "\r\n"), nil
	}
	f, ok := stdin.(*os.File)
	if !ok || !term.IsTerminal(int(f.Fd())) {
		return "", errors.New("非互動環境請使用 --password-stdin")
	}
	ask := func(label string) (string, error) {
		fmt.Fprint(prompt, label)
		b, err := term.ReadPassword(int(f.Fd()))
		fmt.Fprintln(prompt)
		return string(b), err
	}
	p1, err := ask("密碼：")
	if err != nil {
		return "", err
	}
	p2, err := ask("再輸入一次：")
	if err != nil {
		return "", err
	}
	if p1 != p2 {
		return "", errors.New("兩次輸入的密碼不一致")
	}
	return p1, nil
}
