// inspired by https://blog.logrocket.com/rate-limiting-go-application/ and https://www.alexedwards.net/blog/how-to-rate-limit-http-requests

package middlewares

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/grid-stream-org/api/internal/custom_error"
	"golang.org/x/time/rate"
)

// client stores rate limiter and last seen timestamp
type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	clients sync.Map // Use sync.Map because it's threadsafe
)

// Cleanup expired clients every minute
func init() {
	go func() {
		for {
			time.Sleep(time.Minute)
			clients.Range(func(ip, value interface{}) bool {
				client := value.(*client)
				if time.Since(client.lastSeen) > 3*time.Minute {
					clients.Delete(ip)
				}
				return true
			})
		}
	}()
}

// Per-client rate limiter middleware
func PerClientRateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract client IP
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Get or create a rate limiter for the client
		value, exists := clients.Load(ip)
		if !exists {
			// TODO: may need to change limits (3 req/sec, burst of 6), this means that user will be able to make 3 req/sec after 6 req/sec burst
			limiter := rate.NewLimiter(3, 6)
			value = &client{limiter: limiter, lastSeen: time.Now()}
			clients.Store(ip, value)
		}

		cli := value.(*client)
		cli.lastSeen = time.Now()

		// Check if request is allowed
		if !cli.limiter.Allow() {
			message := custom_error.New(http.StatusTooManyRequests, "Too many requests", errors.New("429 Too many requests"))
			w.WriteHeader(http.StatusTooManyRequests)
			err = json.NewEncoder(w).Encode(&message)
			if err != nil {
				http.Error(w, "Internal Server error", http.StatusInternalServerError)
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}
