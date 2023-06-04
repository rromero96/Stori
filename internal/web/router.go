/* package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"github.com/go-chi/chi/v5"
)

// Handler is a type that handles a http request within our framework.
type Handler func(w http.ResponseWriter, r *http.Request) error

// Middleware is a function designed to run some code before and/or after
// another Handler. It is designed to remove boilerplate or other concerns not
// direct to any given Handler.
type Middleware func(http.HandlerFunc) http.HandlerFunc

// wrapMiddleware creates a new handler by wrapping mw around a final
// handler. The middlewares' Handlers will be executed by requests in the order
// they are provided.
func wrapMiddleware(handler http.HandlerFunc, mw []Middleware) http.HandlerFunc {
	// Loop backwards through the middleware invoking each one. Replace the
	// handler with the new wrapped handler. Looping backwards ensures that the
	// first middleware of the slice is the first to be executed by requests.
	for i := len(mw) - 1; i >= 0; i-- {
		h := mw[i]
		if h != nil {
			handler = h(handler)
		}
	}

	return handler
}

// ErrorEncoder is used to define centralized error handler for your application.
type ErrorEncoder func(ctx context.Context, err error, w http.ResponseWriter)

// DefaultErrorEncoder writes the error to the ResponseWriter, by default a
// content type of text/plain, a body of the plain text of the error, and a
// status code of 500. If the error implements Headerer, the provided headers
// will be applied to the response. If the error implements json.Marshaler, and
// the marshaling succeeds, a content type of application/json and the JSON
// encoded form of the error will be used. If the error implements StatusCoder,
// the provided StatusCode will be used instead of 500.
func DefaultErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	contentType, body := "text/plain; charset=utf-8", []byte(err.Error())
	if m, ok := err.(json.Marshaler); ok {
		if jsonBody, marshalErr := m.MarshalJSON(); marshalErr == nil {
			contentType, body = "application/json; charset=utf-8", jsonBody
		}
	}

	w.Header().Set("Content-Type", contentType)
	if h, ok := err.(Headerer); ok {
		for k, values := range h.Headers() {
			for _, v := range values {
				w.Header().Add(k, v)
			}
		}
	}

	code := http.StatusInternalServerError
	if sc, ok := err.(StatusCoder); ok {
		code = sc.StatusCode()
	}

	w.WriteHeader(code)
	w.Write(body)
}

// StatusCoder is checked by DefaultErrorEncoder. If an error value implements
// StatusCoder, the StatusCode will be used when encoding the error. By default,
// StatusInternalServerError (500) is used.
type StatusCoder interface {
	StatusCode() int
}

// Headerer is checked by DefaultErrorEncoder. If an error value implements
// Headerer, the provided headers will be applied to the response writer, after
// the Content-Type is set.
type Headerer interface {
	Headers() http.Header
}

// ErrorHandler receives a transport error to be processed for diagnostic purposes.
type ErrorHandler func(ctx context.Context, err error)

// DefaultNotFoundHandler handler for routing paths that could not be found.
var DefaultNotFoundHandler = func(w http.ResponseWriter, r *http.Request) {
	err := NewErrorf(http.StatusNotFound, "resource %s not found", r.URL.Path)
	_ = EncodeJSON(w, err, http.StatusNotFound)
}

// Router is a http.Handler which can be used to dispatch requests to different
// handler functions via configurable routes.
type Router struct {
	mux        *chi.Mux
	mw         []Middleware
	errEncoder ErrorEncoder
	errHandler ErrorHandler
}

// New instantiates a `Router`.
func New() *Router {
	mux := chi.NewRouter()
	mux.NotFound(DefaultNotFoundHandler)

	return &Router{
		mux:        mux,
		errEncoder: DefaultErrorEncoder,
	}
}

// Use appends a middleware handler to the middleware stack.
func (r *Router) Use(middlewares ...Middleware) {
	r.mw = append(r.mw, middlewares...)
}

// NotFound sets a custom http.HandlerFunc for routing paths that could
// not be found. The default 404 handler is `http.NotFound`.
func (r *Router) NotFound(fn http.HandlerFunc) {
	r.mux.NotFound(fn)
}

// ErrorEncoder sets the given fn as ErrorEncoder.
func (r *Router) ErrorEncoder(fn ErrorEncoder) {
	r.errEncoder = fn
}

// ErrorHandler sets the given fn as ErrorHandler.
func (r *Router) ErrorHandler(fn ErrorHandler) {
	r.errHandler = fn
}

// Group creates a new RouteGroup with the given p prefix and middlewares which are
// chained after this Router's middlewares.
func (r *Router) Group(p string, mw ...Middleware) *RouteGroup {
	return &RouteGroup{
		router: r,
		path:   p,
		mw:     mw,
	}
}

func (r *Router) handlerAdapter(handler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		err := handler(w, req)
		if err == nil {
			return
		}

		r.errHandler(req.Context(), err)
		r.errEncoder(req.Context(), err, w)
	}
}

// Method adds the route `pattern` that matches `method` http method to
// execute the `handler` http.Handler wrapped by `mw`.
func (r *Router) Method(method, pattern string, handler Handler, mw ...Middleware) {
	r.mux.Method(method, pattern, r.wrapHandler(r.handlerAdapter(handler), mw...))
}

// Any adds the route `pattern` that matches any http method to execute the `handler` http.Handler wrapped by `mw`.
func (r *Router) Any(pattern string, handler Handler, mw ...Middleware) {
	r.mux.Handle(pattern, r.wrapHandler(r.handlerAdapter(handler), mw...))
}

func (r *Router) wrapHandler(handler http.HandlerFunc, mw ...Middleware) http.HandlerFunc {
	// First wrap handler specific middleware around this handler.
	handler = wrapMiddleware(handler, mw)

	// Add the application's general middleware to the handler chain.
	handler = wrapMiddleware(handler, r.mw)

	return func(w http.ResponseWriter, r *http.Request) {
		routeCtx := chi.RouteContext(r.Context())
		routeParams := routeCtx.URLParams

		params := make(URIParams, len(routeParams.Keys))
		for i := range routeParams.Keys {
			params[routeParams.Keys[i]] = routeParams.Values[i]
		}

		r = r.WithContext(WithParams(r.Context(), params))
		handler(w, r)
	}
}

// Get is a shortcut for r.Method(http.MethodGet, pattern, handle, mw).
func (r *Router) Get(pattern string, handler Handler, mw ...Middleware) {
	r.Method(http.MethodGet, pattern, handler, mw...)
}

// Head is a shortcut for r.Method(http.MethodHead, pattern, handle, mw).
func (r *Router) Head(pattern string, handler Handler, mw ...Middleware) {
	r.Method(http.MethodHead, pattern, handler, mw...)
}

// Options is a shortcut for r.Method(http.MethodOptions, pattern, handle, mw).
func (r *Router) Options(pattern string, handler Handler, mw ...Middleware) {
	r.Method(http.MethodOptions, pattern, handler, mw...)
}

// Post is a shortcut for r.Method(http.MethodPost, pattern, handle, mw).
func (r *Router) Post(pattern string, handler Handler, mw ...Middleware) {
	r.Method(http.MethodPost, pattern, handler, mw...)
}

// Put is a shortcut for r.Method(http.MethodPut, pattern, handle, mw).
func (r *Router) Put(pattern string, handler Handler, mw ...Middleware) {
	r.Method(http.MethodPut, pattern, handler, mw...)
}

// Patch is a shortcut for r.Method(http.MethodPatch, pattern, handle, mw).
func (r *Router) Patch(pattern string, handler Handler, mw ...Middleware) {
	r.Method(http.MethodPatch, pattern, handler, mw...)
}

// Delete is a shortcut for r.Method(http.MethodDelete, pattern, handle, mw).
func (r *Router) Delete(pattern string, handler Handler, mw ...Middleware) {
	r.Method(http.MethodDelete, pattern, handler, mw...)
}

// Trace is a shortcut for r.Method(http.MethodTrace, pattern, handle, mw).
func (r *Router) Trace(pattern string, handler Handler, mw ...Middleware) {
	r.Method(http.MethodTrace, pattern, handler, mw...)
}

// ServeHTTP conforms to the http.Handler interface.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// Route describes the details of a routing handler.
type Route struct {
	Method      string
	Route       string
	Handler     http.Handler
	Middlewares []func(http.Handler) http.Handler
}

// Routes returns the routing tree in an easily traversable structure.
func (r *Router) Routes() ([]Route, error) {
	var routes []Route
	walkFunc := func(method string, route string, handler http.Handler, mw ...func(http.Handler) http.Handler) error {
		routes = append(routes, Route{
			Method:      method,
			Route:       route,
			Handler:     handler,
			Middlewares: mw,
		})
		return nil
	}

	if err := chi.Walk(r.mux, walkFunc); err != nil {
		return nil, fmt.Errorf("generating routes: %v", err)
	}

	return routes, nil
}

// RouteGroup represents a group of routes that share the same path prefix and middlewares.
type RouteGroup struct {
	router *Router
	path   string
	mw     []Middleware
}

// Group creates a new RouteGroup with the given path relative to the existing RouteGroup path
// and middlewares which are chained after this RouteGroup's middlewares.
func (g *RouteGroup) Group(p string, mw ...Middleware) *RouteGroup {
	return g.router.Group(path.Join(g.path, p), g.appendMiddlewares(mw)...)
}

// Method adds the route `pattern` that matches `method` http method to
// execute the `handler` http.Handler wrapped by `mw`.
func (g *RouteGroup) Method(method, pattern string, handler Handler, mw ...Middleware) {
	g.router.Method(method, path.Join(g.path, pattern), handler, g.appendMiddlewares(mw)...)
}

// Any adds the route `pattern` that matches any http method to execute the `handler` http.Handler wrapped by `mw`.
func (g *RouteGroup) Any(pattern string, handler Handler, mw ...Middleware) {
	g.router.Any(path.Join(g.path, pattern), handler, g.appendMiddlewares(mw)...)
}

func (g *RouteGroup) appendMiddlewares(mw []Middleware) []Middleware {
	var m []Middleware
	m = append(m, g.mw...)
	m = append(m, mw...)
	return m
}

// Get is a shortcut for g.Method(http.MethodGet, pattern, handle, mw).
func (g *RouteGroup) Get(pattern string, handler Handler, mw ...Middleware) {
	g.Method(http.MethodGet, pattern, handler, mw...)
}

// Head is a shortcut for g.Method(http.MethodHead, pattern, handle, mw).
func (g *RouteGroup) Head(pattern string, handler Handler, mw ...Middleware) {
	g.Method(http.MethodHead, pattern, handler, mw...)
}

// Options is a shortcut for g.Method(http.MethodOptions, pattern, handle, mw).
func (g *RouteGroup) Options(pattern string, handler Handler, mw ...Middleware) {
	g.Method(http.MethodOptions, pattern, handler, mw...)
}

// Post is a shortcut for g.Method(http.MethodPost, pattern, handle, mw).
func (g *RouteGroup) Post(pattern string, handler Handler, mw ...Middleware) {
	g.Method(http.MethodPost, pattern, handler, mw...)
}

// Put is a shortcut for g.Method(http.MethodPut, pattern, handle, mw).
func (g *RouteGroup) Put(pattern string, handler Handler, mw ...Middleware) {
	g.Method(http.MethodPut, pattern, handler, mw...)
}

// Patch is a shortcut for g.Method(http.MethodPatch, pattern, handle, mw).
func (g *RouteGroup) Patch(pattern string, handler Handler, mw ...Middleware) {
	g.Method(http.MethodPatch, pattern, handler, mw...)
}

// Delete is a shortcut for g.Method(http.MethodDelete, pattern, handle, mw).
func (g *RouteGroup) Delete(pattern string, handler Handler, mw ...Middleware) {
	g.Method(http.MethodDelete, pattern, handler, mw...)
}

// Trace is a shortcut for g.Method(http.MethodTrace, pattern, handle, mw).
func (g *RouteGroup) Trace(pattern string, handler Handler, mw ...Middleware) {
	g.Method(http.MethodTrace, pattern, handler, mw...)
}
*/