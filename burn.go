package burn

import (
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/codegangsta/inject"
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

func (m *Burn) Action(handler Handler) {
	validateHandler(handler)
	m.action = handler
}

func (m *Burn) Logger(logger *log.Logger) {
	m.logger = logger
	m.Map(m.logger)
}

func (m *Burn) Use(handler Handler) {
	validateHandler(handler)

	m.handlers = append(m.handlers, handler)
}

func (m *Burn) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	m.createContext(res, req).run()
}

func (m *Burn) RunOnAddr(addr string) {
	logger := m.Injector.Get(reflect.TypeOf(m.logger)).Interface().(*log.Logger)
	logger.Printf("listening on %s (%s)\n", addr, Env)
	logger.Fatalln(http.ListenAndServe(addr, m))
}

func (m *Burn) Run() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	host := os.Getenv("HOST")

	m.RunOnAddr(host + ":" + port)
}

func (m *Burn) createContext(res http.ResponseWriter, req *http.Request) *context {
	c := &context{inject.New(), m.handlers, m.action, NewResponseWriter(res), 0}
	c.SetParent(m)
	c.MapTo(c, (*Context)(nil))
	c.MapTo(c.rw, (*http.ResponseWriter)(nil))
	c.Map(req)
	return c
}

type ClassicBurn struct {
	*Burn
	Router
}

func Classic() *ClassicBurn {
	r := NewRouter()
	m := New()
	m.Use(Logger())
	m.Use(Recovery())
	m.Use(Static("public"))
	m.MapTo(r, (*Routes)(nil))
	m.Action(r.Handle)
	return &ClassicBurn{m, r}
}

type Handler interface{}

func validateHandler(handler Handler) {
	if reflect.TypeOf(handler).Kind() != reflect.Func {
		panic("burn handler must be a callable func")
	}
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

func (c *context) handler() Handler {
	if c.index < len(c.handlers) {
		return c.handlers[c.index]
	}
	if c.index == len(c.handlers) {
		return c.action
	}
	panic("invalid index for context handler")
}

func (c *context) Next() {
	c.index += 1
	c.run()
}

func (c *context) Written() bool {
	return c.rw.Written()
}

func (c *context) run() {
	for c.index <= len(c.handlers) {
		_, err := c.Invoke(c.handler())
		if err != nil {
			panic(err)
		}
		c.index += 1

		if c.Written() {
			return
		}
	}
}
