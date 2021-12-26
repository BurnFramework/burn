package responsewriter

import (
  "net"
  "net/http"
)

type ResponseWriter interface {
    http.ResponseWriter
    http.Flusher
    http.Hijacker
    
    Status() int
    Size() int
    
    Writer() bool
    
    Before(BeforeFunc)
}

type BeforeFunc func(ResponseWriter)

