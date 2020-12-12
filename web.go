package web

import (
	"context"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

// ctxKey represents the type of value for the context key.
type ctxKey int

// KeyValues is how request values are stored/retrieved.
const KeyValues ctxKey = 1

// Values represent state for each request.
type Values struct {
	Now        time.Time
	StatusCode int
}

// Handler is a version of an HTTP Handler that takes in a context
// and can return an error.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers.
type App struct {
	*mux.Router
	shutdown chan os.Signal
	mw       []Middleware
}

// NewApp creates an App value that handle a set of routes for the application.
func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {
	app := App{
		Router:   mux.NewRouter(),
		shutdown: shutdown,
		mw:       mw,
	}

	return &app
}

// Handle is our mechanism for mounting Handlers for a given HTTP verb and path
func (a *App) Handle(method string, path string, handler Handler, mw ...Middleware) {
	// First wrap handler specific middleware around this handler
	// then, the application's general middleware
	handler = wrapMiddleware(mw, handler)
	handler = wrapMiddleware(a.mw, handler)

	// The function to execute for each request.
	h := func(w http.ResponseWriter, r *http.Request) {

		// Set the context with the required values to
		// process the request.
		v := Values{
			Now: time.Now(),
		}
		ctx := context.WithValue(context.Background(), KeyValues, &v)

		// Call the wrapped handler functions.
		if err := handler(ctx, w, r); err != nil {
			a.SignalShutdown()
			return
		}
	}

	a.Router.HandleFunc(path, h).Methods(method)
}

// ServeHTTP implements the http.Handler interface.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.Router.ServeHTTP(w, r)
}

// SignalShutdown is used to gracefully shutdown the app when an
// integrity issue is identified.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

// Params returns the web call parameters from the request.
func Params(r *http.Request) map[string]string {
	return mux.Vars(r)
}
