package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

var visitors = make(map[string]time.Time)
var mu sync.Mutex

// RateLimiter limits requests per IP
func RateLimiter(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		mu.Lock()
		last, exists := visitors[ip]
		if exists && time.Since(last) < time.Second {
			mu.Unlock()
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		visitors[ip] = time.Now()
		mu.Unlock()

		next(w, r)
	}
}
