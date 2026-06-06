package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/OmarAraby/go-ecommerce/internal/api/response"
)

type client struct {
	limiter  *rate.Limiter // token bucket limiter for this client
	lastSeen time.Time     // last time this client made a request on this api
}

// IPRateLimiter tracks a token-bucket limiter per client IP.
type IPRateLimiter struct {
	mu      sync.Mutex
	clients map[string]*client
	r       rate.Limit // tokens added per second
	b       int        // max burst size
}

func newIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	rl := &IPRateLimiter{
		clients: make(map[string]*client),
		r:       r,
		b:       b,
	}
	go rl.cleanup()
	return rl
}

func (rl *IPRateLimiter) get(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	c, ok := rl.clients[ip]
	if !ok {
		c = &client{limiter: rate.NewLimiter(rl.r, rl.b)}
		rl.clients[ip] = c
	}
	c.lastSeen = time.Now()
	return c.limiter
}

// cleanup removes entries that have not made a request in the last 3 minutes.
func (rl *IPRateLimiter) cleanup() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, c := range rl.clients {
			if time.Since(c.lastSeen) > 3*time.Minute {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// RateLimit returns a middleware that allows r requests/second with a burst of b per client IP.
// Requests that exceed the limit receive 429 Too Many Requests.
func RateLimit(r rate.Limit, b int) func(http.Handler) http.Handler {
	rl := newIPRateLimiter(r, b)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ip := clientIP(req)
			if !rl.get(ip).Allow() {
				w.Header().Set("Retry-After", "1")
				response.TooManyRequests(w)
				return
			}
			next.ServeHTTP(w, req)
		})
	}
}

// clientIP extracts the real client IP, respecting proxy headers.
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.TrimSpace(strings.Split(xff, ",")[0])
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
