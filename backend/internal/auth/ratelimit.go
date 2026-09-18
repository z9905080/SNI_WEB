package auth

import (
	"sync"
	"time"
)

// Limiter 是記憶體中的滑動視窗計數器，僅適用單一實例。
type Limiter struct {
	mu     sync.Mutex
	limit  int
	window time.Duration
	now    func() time.Time
	hits   map[string][]time.Time
}

func NewLimiter(limit int, window time.Duration, now func() time.Time) *Limiter {
	return &Limiter{limit: limit, window: window, now: now, hits: map[string][]time.Time{}}
}

func (l *Limiter) Blocked(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.prune(key, l.now())) >= l.limit
}

func (l *Limiter) Add(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	l.hits[key] = append(l.prune(key, now), now)
	// ponytail: 超過 1 萬個 key 才整批清理；多實例部署時改用共用儲存（如 Redis）
	if len(l.hits) > 10000 {
		for k := range l.hits {
			l.prune(k, now)
		}
	}
}

func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.hits, key)
}

// prune 移除過期紀錄並回傳剩餘紀錄；呼叫端需持有鎖。
func (l *Limiter) prune(key string, now time.Time) []time.Time {
	hits := l.hits[key]
	cutoff := now.Add(-l.window)
	i := 0
	for i < len(hits) && !hits[i].After(cutoff) {
		i++
	}
	hits = hits[i:]
	if len(hits) == 0 {
		delete(l.hits, key)
		return nil
	}
	l.hits[key] = hits
	return hits
}
