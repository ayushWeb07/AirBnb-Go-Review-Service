package middlewares

import (
	"net/http"
	"time"

	"github.com/ayushWeb07/AirBnb-Go-Review-Service/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Review-Service/internal/utils"
	"github.com/go-chi/httprate"
)

func RateLimiter(serverConfig *config.ServerConfig) func(next http.Handler) http.Handler {
	return httprate.Limit(
		serverConfig.RequestsPerMinute,
		time.Minute,
		httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint),
		httprate.WithLimitHandler(func(resWriter http.ResponseWriter, req *http.Request) {
			utils.WriteJsonResponse(http.StatusTooManyRequests, resWriter, map[string]any{
				"success": false,
				"message": "Too many requests",
				"error":   "You have been rate-limited. Please, slow down and try again after sometime",
			})
		}),
	)
}
