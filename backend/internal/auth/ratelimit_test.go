package auth

import (
	"testing"
	"time"
)

func TestLimiter(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	l := NewLimiter(3, 15*time.Minute, func() time.Time { return now })

	for i := range 3 {
		if l.Blocked("a") {
			t.Fatalf("第 %d 次前不應被擋", i+1)
		}
		l.Add("a")
		now = now.Add(time.Minute)
	}
	if !l.Blocked("a") {
		t.Fatal("達上限應被擋")
	}
	if l.Blocked("b") {
		t.Fatal("不同 key 不受影響")
	}

	// 第一次紀錄在 00:00，15 分鐘後過期
	now = time.Date(2026, 1, 1, 0, 15, 0, 0, time.UTC)
	if l.Blocked("a") {
		t.Fatal("最舊的一筆過期後應解除")
	}

	l.Add("a")
	if !l.Blocked("a") {
		t.Fatal("再加一筆又達上限")
	}
	l.Reset("a")
	if l.Blocked("a") {
		t.Fatal("Reset 後應解除")
	}
}
