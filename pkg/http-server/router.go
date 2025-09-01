package httpserver

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/soulcodex/deus-cargo-tracker/pkg/errutil"
)

var (
	ErrStartingHTTPServer = errutil.NewError("http server start errored")
	ErrShutdownHTTPServer = errutil.NewError("http server shutdown errored")
)

type Router struct {
	muxRouter *mux.Router
	server    *http.Server
	addr      string
}

type Middleware func(next http.Handler) http.Handler

func New(options ...RouterConfigFunc) Router {
	opts := NewRouterConfig(options...)

	readTimeout := time.Duration(opts.ReadTimeout) * time.Second
	writeTimeout := time.Duration(opts.WriteTimeout) * time.Second

	m := muxRouter()
	for _, middleware := range opts.Middlewares {
		m.Use(mux.MiddlewareFunc(middleware))
	}

	h := &http.Server{
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
	}

	return Router{m, h, net.JoinHostPort(opts.Host, strconv.Itoa(opts.Port))}
}

func muxRouter() *mux.Router {
	return mux.NewRouter().StrictSlash(true)
}

func (router *Router) Use(mws ...Middleware) {
	for _, mw := range mws {
		router.muxRouter.Use(mux.MiddlewareFunc(mw))
	}
}

func (router *Router) ListenAndServe() error {
	router.server.Addr = router.addr
	router.server.Handler = router.muxRouter

	// Use default options
	err := router.server.ListenAndServe()
	if err == nil {
		return nil
	}

	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return ErrStartingHTTPServer.Wrap(err)
}

func (router *Router) Shutdown(ctx context.Context) error {
	if err := router.server.Shutdown(ctx); err != nil {
		return ErrShutdownHTTPServer.Wrap(err)
	}

	return nil
}

func (router *Router) GetMuxRouter() *mux.Router {
	return router.muxRouter
}

func (router *Router) AddMiddleware(middlewares ...Middleware) {
	for _, middleware := range middlewares {
		router.muxRouter.Use(mux.MiddlewareFunc(middleware))
	}
}

func (router *Router) handleMultipleMethods(
	methods []string,
	path string,
	handler http.Handler,
	middlewares ...Middleware,
) {
	subRouter := router.muxRouter.NewRoute().Subrouter()
	for _, middleware := range middlewares {
		subRouter.Use(mux.MiddlewareFunc(middleware))
	}

	subRouter.HandleFunc(path, handler.ServeHTTP).Methods(methods...)
}

func (router *Router) Handle(methods []string, path string, handler http.Handler, middlewares ...Middleware) {
	router.handleMultipleMethods(methods, path, handler, middlewares...)
}

func (router *Router) Get(path string, handler http.Handler, middlewares ...Middleware) {
	router.Handle([]string{http.MethodGet}, path, handler, middlewares...)
}

func (router *Router) Post(path string, handler http.Handler, middlewares ...Middleware) {
	router.Handle([]string{http.MethodPost}, path, handler, middlewares...)
}

func (router *Router) Put(path string, handler http.Handler, middlewares ...Middleware) {
	router.Handle([]string{http.MethodPut}, path, handler, middlewares...)
}

func (router *Router) Patch(path string, handler http.Handler, middlewares ...Middleware) {
	router.Handle([]string{http.MethodPatch}, path, handler, middlewares...)
}

func (router *Router) Delete(path string, handler http.Handler, middlewares ...Middleware) {
	router.Handle([]string{http.MethodDelete}, path, handler, middlewares...)
}

func (router *Router) Head(path string, handler http.Handler, middlewares ...Middleware) {
	router.Handle([]string{http.MethodHead}, path, handler, middlewares...)
}

func (router *Router) Options(path string, handler http.Handler, middlewares ...Middleware) {
	router.Handle([]string{http.MethodOptions}, path, handler, middlewares...)
}

func (router *Router) Route(methods []string, path string, handler http.Handler, middlewares ...Middleware) {
	router.handleMultipleMethods(methods, path, handler, middlewares...)
}
