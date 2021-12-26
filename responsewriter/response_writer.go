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

func (rw *responseWriter) WriteHeader(s int) {
	rw.callBefore()
	rw.ResponseWriter.WriteHeader(s)
	rw.status = s
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.Written() {
		rw.WriteHeader(http.StatusOK)
	}
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func (rw *responseWriter) Status() int {
	return rw.status
}

type closeNotifyResponseWriter struct {
	responseWriter
	closeNotifier http.CloseNotifier
}
