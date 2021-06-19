package router

import (
	"net/http"
	"time"
)

type (
	//APIOptions options when initializing the API
	APIOptions struct {
		ListenAddress string        `json:"listen_address"`
		CORS          CORSOptions   `json:"cors"`
		RateInterval  time.Duration `json:"rate_interval"`
		RateLimit     uint64        `json:"rate_limit"`
	}

	//CORSOptions cors options for the API
	CORSOptions struct {
		Enabled bool     `json:"enabled"`
		Headers []string `json:"headers"`
		Origins []string `json:"origins"`
		Methods []string `json:"methods"`
	}

	//APIRequest a request made to the API
	APIRequest struct {
		*http.Request
		AccessKey    string
		AccessSecret string
		AuthToken    string
		IPAddress    string
		Timestamp    time.Time
	}

	//MiddlewareFunc a middleware function for the router
	MiddlewareFunc func(*APIRouter, APIEndpoint, APIHandlerFunc) APIHandlerFunc

	//APIHandlerFunc a handler for an API endpoint
	APIHandlerFunc func(http.ResponseWriter, *APIRequest)

	//APIEndpoint an endpoint of the API to retrieve or set data
	APIEndpoint struct {
		Name        string
		Method      string
		Pattern     string
		Secure      bool
		Permissions []string
		Middleware  []MiddlewareFunc
		Handler     APIHandlerFunc
	}

	//APIResponse APIResponse
	APIResponse struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	}

	//AuthToken AuthToken
	AuthToken struct {
		UserID         string    `json:"user_id"`
		IPAddress      string    `json:"ip_addr"`
		Permissions    []string  `json:"permissions"`
		ExpirationDate time.Time `json:"expiration_date"`
	}

	//APIRouter an API router
	APIRouter struct {
		options    APIOptions
		middleware []MiddlewareFunc
		endpoints  []APIEndpoint
		server     *http.Server
	}
)
