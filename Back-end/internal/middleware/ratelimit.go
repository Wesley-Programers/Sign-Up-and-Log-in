package middleware

import (
	"net/http"
	"time"

	"index/Back-end/internal/security"
)

func RateLimitMiddleware(limiter *security.RedisLimiter, maxAttempts int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			userID, ok := r.Context().Value(UserIDKey).(string)
			if !ok || userID == "" {
				http.Error(w, "Error", http.StatusAccepted)
				return
			}

			ctx := r.Context()
			key := "change-email:" + userID
			
			allowed, err := limiter.CheckLimit(ctx, key, 3, 24*time.Hour)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			if !allowed {
				http.Error(w, "Lot of attempts", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}