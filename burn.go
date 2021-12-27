package burn

import (
	"log"
	"os"
)

type Burn struct {
	inject.Injector
	handlers []Handler
	action   Handler
	logger   *log.Logger
}

func New() *Burn {
	m := &Burn{Injector: inject.New(), action: func() {}, logger: log.New(os.Stdout, "[burn] ", 0)}
	m.Map(m.logger)
	m.Map(defaultReturnHandler())
	return m
}

func (m *Burn) Handlers(handlers ...Handler) {
	m.handlers = make([]Handler, 0)
	for _, handler := range handlers {
		m.Use(handler)
	}
}

type Handler interface{}

type ClassicBurn struct {
	*Burn
	Router
}

type Context interface {
	inject.Injector
	Next()
	Written() bool
}

type context struct {
	inject.Injector
	handlers []Handler
	action   Handler
	rw       ResponseWriter
	index    int
}
