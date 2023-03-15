package nhttp

import (
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

// RouterOpt represents available option to init a router
type RouterOpt struct {
	RootRouter *mux.Router
	BasePath   string
	Logger     nlog.Logger
}

type Router struct {
	*mux.Router
	logger      nlog.Logger
	middlewares map[string]Middleware
}

func (r *Router) NewHandler(fn HandlerFunc) Handler {
	return Handler{Logger: r.logger, Fn: fn}
}

func (r *Router) Handle(path string, fn HandlerFunc) *mux.Route {
	// Create handler
	h := Handler{Logger: r.logger, Fn: fn}
	return r.NewRoute().Path(path).Handler(h)
}

func (r *Router) HandleWithMiddleware(path string, middlewareName string, fn HandlerFunc) *mux.Route {
	// Create handler
	h := Handler{Logger: r.logger, Fn: fn}

	// Get middleware
	m, ok := r.middlewares[middlewareName]
	if ok {
		// If middleware found, add handler into middleware
		h = m(h)
	}

	// Register handler
	return r.NewRoute().Path(path).Handler(h)
}

func (r *Router) RegisterMiddleware(name string, m Middleware) {
	r.middlewares[name] = m
}

// NewCORSRouter return new router that add headers for CORS
func NewCORSRouter(r *Router) http.Handler {
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Authorization", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions})
	return handlers.CORS(originsOk, headersOk, methodsOk)(r)
}

func NewApiRouter(opt RouterOpt) *Router {
	// Create Router
	r := Router{
		Router:      opt.RootRouter.PathPrefix(opt.BasePath).Subrouter(),
		logger:      opt.Logger,
		middlewares: make(map[string]Middleware),
	}

	// Set standard error handler
	r.MethodNotAllowedHandler = r.NewHandler(HandleMethodNotAllowed)
	r.NotFoundHandler = r.NewHandler(HandleNotFound)

	// Return routes
	return &r
}

func HandleNotFound(_ *http.Request) (*Success, error) {
	return nil, ErrNotFound
}

func HandleMethodNotAllowed(_ *http.Request) (*Success, error) {
	return nil, ErrMethodNotAllowed
}
