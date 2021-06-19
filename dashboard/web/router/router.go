package router

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/siacentral/sia-host-dashboard/web"
)

var (
	//ErrNotRunning returned from shutdown if the server was not started
	ErrNotRunning = errors.New("server is not running")
)

//NewRouter creates a new router
func NewRouter(endpoints []APIEndpoint, opts APIOptions) (router *APIRouter) {
	router = &APIRouter{
		middleware: []MiddlewareFunc{rateLimitMiddleware},
		endpoints:  endpoints,
		options:    opts,
	}

	return
}

//AddMiddleware appends a new middleware function to the router
func (router *APIRouter) AddMiddleware(middleware MiddlewareFunc) {
	router.middleware = append(router.middleware, middleware)
}

//ListenAndServe starts the router listening for connections
func (router *APIRouter) ListenAndServe() error {
	var handler http.Handler
	r := mux.NewRouter().StrictSlash(true)

	apiRouter := r.PathPrefix("/api").Subrouter()

	r.Path("/favicon.ico").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		_, _ = w.Write([]byte{})
	})

	for _, endpoint := range router.endpoints {
		apiRouter.Methods(endpoint.Method).Path(endpoint.Pattern).
			Name(endpoint.Name).Handler(router.attachMiddleware(endpoint))
	}

	r.PathPrefix("/").Handler(http.FileServer(web.Assets))

	if router.options.CORS.Enabled {
		handler = handlers.CORS(handlers.AllowedHeaders(router.options.CORS.Headers),
			handlers.AllowedOrigins(router.options.CORS.Origins),
			handlers.AllowedMethods(router.options.CORS.Methods))(r)
	} else {
		handler = r
	}

	handler = http.TimeoutHandler(handler, time.Minute*1, "timeout")

	l, err := net.Listen("tcp", router.options.ListenAddress)
	if err != nil {
		return err
	}

	router.server = &http.Server{
		Handler:      handler,
		Addr:         router.options.ListenAddress,
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}

	if err := router.server.Serve(l); err != nil && errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

//Shutdown attempts to gracefully shutdown the running server
func (router *APIRouter) Shutdown(ctx context.Context) error {
	if router.server == nil {
		return ErrNotRunning
	}

	return router.server.Shutdown(ctx)
}

func chainMiddleware(router *APIRouter, endpoint APIEndpoint, middleware ...MiddlewareFunc) APIHandlerFunc {
	if len(middleware) == 0 {
		return endpoint.Handler
	}

	return middleware[0](router, endpoint, chainMiddleware(router, endpoint, middleware[1:]...))
}

func (router *APIRouter) attachMiddleware(endpoint APIEndpoint) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		var remoteAddr string

		if forwardHeader := r.Header.Get("X-Forwarded-For"); len(forwardHeader) > 0 {
			remoteAddr = forwardHeader
		} else {
			remoteAddr = r.RemoteAddr
		}

		request := &APIRequest{
			Request:   r,
			IPAddress: remoteAddr,
			Timestamp: time.Now(),
		}

		defer func() { logRequest(endpoint, request, err) }()

		chainMiddleware(router, endpoint, append(router.middleware, endpoint.Middleware...)...)(w, request)
	})
}

func logRequest(endpoint APIEndpoint, request *APIRequest, err error) {
	if err != nil {
		log.Println(err)
	}

	log.Printf(
		"%s %s %s %s %s",
		request.IPAddress,
		endpoint.Name,
		endpoint.Method,
		request.Request.URL.Path,
		time.Since(request.Timestamp),
	)
}

//HandleError returns an error response with the specified message and http status
func HandleError(message string, status int, w http.ResponseWriter, r *APIRequest) {
	var resp APIResponse

	resp.Message = message
	resp.Type = "error"

	SendJSONResponse(resp, status, w, r)
}

//SendJSONResponse SendJSONResponse
func SendJSONResponse(response interface{}, status int, w http.ResponseWriter, r *APIRequest) {
	data, err := json.Marshal(response)

	if err != nil {
		w.WriteHeader(500)
		log.Println("Send Response:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(data)

	if err != nil {
		log.Println("Write data:", err)
	}
}
