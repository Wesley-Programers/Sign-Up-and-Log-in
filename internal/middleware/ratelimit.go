package middleware

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"ShieldAuth-API/internal/security"
)

func RateLimitMiddleware(limiter *security.RedisLimiter, keyPrefix string, maxAttempts int, duration time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			key := EmailKeyFunc(r)
			if key == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "identify provider missing for rate limiting"})
				return
			}
			
			allowed, err := limiter.CheckLimit(r.Context(), key, maxAttempts, duration)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			if !allowed {
				w.Header().Set("Retry-After", strconv.Itoa(int(duration.Seconds())))
				http.Error(w, "Too many requests. Please try again later.", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
