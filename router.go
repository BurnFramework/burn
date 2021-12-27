package burn

import "net/http"

type Params map[string]string

type Router interface {
	Routes

	Group(string, func(Router), ...Handler)
	Get(string, ...Handler) Route
	Patch(string, ...Handler) Route
	Post(string, ...Handler) Route
	Put(string, ...Handler) Route
	Delete(string, ...Handler) Route
	Options(string, ...Handler) Route
	Head(string, ...Handler) Route
	Any(string, ...Handler) Route
	AddRoute(string, string, ...Handler) Route
	NotFound(...Handler)
	Handle(http.ResponseWriter, *http.Request, Context)
}

type router struct {
	routes   []*route
	notFound []*notFound
}

func NewRouter() {

}
