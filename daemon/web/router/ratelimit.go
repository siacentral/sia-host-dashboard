package router

import (
	"net/http"
	"sync"
	"time"
)

type (
	limitter struct {
		Requests  uint64
		ResetTime time.Time
	}
)

var (
	rateUsers = make(map[string]limitter)
	mu        = sync.Mutex{}
)

func rateLimitMiddleware(router *APIRouter, endpoint APIEndpoint, handler APIHandlerFunc) APIHandlerFunc {
	return APIHandlerFunc(func(w http.ResponseWriter, r *APIRequest) {
		mu.Lock()

		defer mu.Unlock()

		limitter := rateUsers[r.IPAddress]
		current := time.Now()

		if limitter.ResetTime.Before(current) {
			limitter.ResetTime = current.Add(router.options.RateInterval)
			limitter.Requests = 0
		}

		limitter.Requests++

		if limitter.Requests > router.options.RateLimit {
			HandleError("too many requests", 429, w, r)
			return
		}

		rateUsers[r.IPAddress] = limitter

		handler(w, r)
	})
}
