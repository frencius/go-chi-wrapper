package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

// IServer is an interface that has list of basic http server and middleware functions
type IServer interface {
	Listen(port int) *http.Server
	Close(server *http.Server) error

	Get(pattern string, f http.HandlerFunc)
	Post(pattern string, f http.HandlerFunc)
	Put(pattern string, f http.HandlerFunc)
	Patch(pattern string, f http.HandlerFunc)
	Options(pattern string, f http.HandlerFunc)
	Connect(pattern string, f http.HandlerFunc)
	Head(pattern string, f http.HandlerFunc)
	Delete(pattern string, f http.HandlerFunc)
	Trace(pattern string, f http.HandlerFunc)

	WithValue(key string, value interface{})
	AllowCORS()
	Swagger(pattern string, path string)
}

// Server is an HTTP Server instance contains Chi router
type Server struct {
	ctx    context.Context
	router *chi.Mux
}

// New creates Server instance that contains Chi router
// ctx is passed and store in Chi instance
func New(ctx context.Context) *Server {
	s := new(Server)

	s.ctx = ctx
	s.router = chi.NewRouter()

	return s
}

// Listen returns HTTP server
// The ListenAndServe function run in goroutine
func (h *Server) Listen(port int) *http.Server {
	srv := &http.Server{Addr: fmt.Sprintf(":%v", port), Handler: h.router}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Print("HTTP SERVER ERROR", err)
		}
	}()

	log.Printf(fmt.Sprintf("HTTP SERVER STARTED ON PORT %v", port))

	return srv
}

// Close shutdowns HTTP server
func (h *Server) Close(server *http.Server) error {
	ctxShutdown, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer func() {
		cancel()
	}()

	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Print("API SERVER FAILED TO SHUTDOWN", err)

		return err
	}

	log.Print("API SERVER SUCCESSFULLY SHUTDOWN")

	return nil
}

// Route is used to create subroute
func (h *Server) Route(pattern string, fn func(r IServer)) {
	subRouter := New(h.ctx)
	if fn != nil {
		fn(subRouter)
	}

	h.router.Mount(pattern, subRouter)
}

// ServeHTTP has same definition with (mx *Mux) ServeHTTP()
// see ServeHTTP in chi library
func (h *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

// Get adds the route `pattern` that matches a GET http method to
// execute the `handlerFn` http.HandlerFunc.
func (h *Server) Get(pattern string, f http.HandlerFunc) {
	h.router.Get(pattern, f)
}

// Post adds the route `pattern` that matches a POST http method to
// execute the `handlerFn` http.HandlerFunc.
func (h *Server) Post(pattern string, f http.HandlerFunc) {
	h.router.Post(pattern, f)
}

// Put adds the route `pattern` that matches a PUT http method to
// execute the `handlerFn` http.HandlerFunc.
func (h *Server) Put(pattern string, f http.HandlerFunc) {
	h.router.Put(pattern, f)
}

// Patch adds the route `pattern` that matches a PATCH http method to
// execute the `handlerFn` http.HandlerFunc.
func (h *Server) Patch(pattern string, f http.HandlerFunc) {
	h.router.Patch(pattern, f)
}

// Options adds the route `pattern` that matches a OPTIONS http method to
// execute the `handlerFn` http.HandlerFunc.
func (h *Server) Options(pattern string, f http.HandlerFunc) {
	h.router.Options(pattern, f)
}

// Connect adds the route `pattern` that matches a CONNECT http method to
// execute the `handlerFn` http.HandlerFunc.
func (h *Server) Connect(pattern string, f http.HandlerFunc) {
	h.router.Connect(pattern, f)
}

// Head adds the route `pattern` that matches a HEAD http method to
// execute the `handlerFn` http.HandlerFunc.
func (h *Server) Head(pattern string, f http.HandlerFunc) {
	h.router.Head(pattern, f)
}

// Delete adds the route `pattern` that matches a DELETE http method to
// execute the `handlerFn` http.HandlerFunc.
func (h *Server) Delete(pattern string, f http.HandlerFunc) {
	h.router.Delete(pattern, f)
}

// Trace adds the route `pattern` that matches a TRACE http method to
// execute the `handlerFn` http.HandlerFunc.
func (h *Server) Trace(pattern string, f http.HandlerFunc) {
	h.router.Trace(pattern, f)
}

// WithValue is a short-hand middleware to set a key/value on the request context
// this middleware utilized `github.com/go-chi/chi/middleware`
func (h *Server) WithValue(key string, value interface{}) {
	h.router.Use(
		middleware.WithValue(key, value),
	)
}

// AllowCORS is a middleware to allow HTTP request from JS (web browser)
func (h *Server) AllowCORS() {
	corsHTTP := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "OPTIONS", "CONNECT", "HEAD", "TRACE"},
		AllowedHeaders: []string{"Accept", "Accept-Encoding", "Authorization", "Content-Length", "Content-Type", "X-CSRF-Token"},
	})

	h.router.Use(corsHTTP.Handler)
}

// Swagger is a middleware for generating swagger documentation by giving path of swagger yaml or json file
func (h *Server) Swagger(pattern string, path string) {
	Path = path
	h.router.Get(
		pattern,
		httpSwagger.Handler(httpSwagger.URL("doc.json")))
}
