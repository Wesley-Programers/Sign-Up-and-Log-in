package middleware

import (
	"time"
	"net/http"
	
	"ShieldAuth-API/internal/security"
)

func RateLimitMiddleware(limiter *security.RedisLimiter, keyPrefix string, maxAttempts int, duration time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			userID, ok := r.Context().Value(UserIDKey).(string)
			if !ok || userID == "" {
				http.Error(w, "Error", http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			key := keyPrefix + userID
			
			allowed, err := limiter.CheckLimit(ctx, key, maxAttempts, duration)
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
