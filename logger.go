package burn

import (
	"log"
	"net/http"
	"time"
)

func Logger() Handler {
	return func(res http.ResponseWriter, req *http.Requet, c Context, log *log.Logger) {
		start := time.Now()

		addr := req.Header.Get("X-Real-IP")
		if addr == "" {
			addr = req.Header.Get("X-Forwarded-For")
		}
	}

	log.Printf("Started %s %s for %s", req.Method, req.URL.Path, addr)
}
