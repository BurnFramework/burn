package responsewriter

import (
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

func NewResponseWriter(rw http.ResponseWriter) ResponseWriter {
	newRw := responseWriter{rw, 0, 0, nil}
	if cn, ok := rw.(http.CloseNotifier); ok {
		return &closeNotifyResponseWriter{newRw, cn}
	}
	return &newRw
}

type responseWriter struct {
	http.ResponseWriter
	status     int
	size       int
	BeforeFunc []BeforeFunc
}

type closeNotifyResponseWriter struct {
	responseWriter
	closeNotifier http.CloseNotifier
}
