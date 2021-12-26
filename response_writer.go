package burn

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

type ResponseWriter interface {
	http.ResponseWriter
	http.Flusher
	http.Hijacker
	Status() int
	Written() bool
	Size() int
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
	status      int
	size        int
	beforeFuncs []BeforeFunc
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

func (rw *responseWriter) Size() int {
	return rw.size
}

func (rw *responseWriter) Written() bool {
	return rw.status != 0
}

func (rw *responseWriter) Before(before BeforeFunc) {
	rw.beforeFuncs = append(rw.beforeFuncs, before)
}

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("the ResponseWriter doesn't support the Hijacker interface")
	}
	return hijacker.Hijack()
}

func (rw *responseWriter) callBefore() {
	for i := len(rw.beforeFuncs) - 1; i >= 0; i-- {
		rw.beforeFuncs[i](rw)
	}
}

func (rw *responseWriter) Flush() {
	flusher, ok := rw.ResponseWriter.(http.Flusher)
	if ok {
		flusher.Flush()
	}
}

type closeNotifyResponseWriter struct {
	responseWriter
	closeNotifier http.CloseNotifier
}

func (rw *closeNotifyResponseWriter) CloseNotify() <-chan bool {
	return rw.closeNotifier.CloseNotify()
}
