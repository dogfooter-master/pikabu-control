// THIS FILE IS AUTO GENERATED BY GK-CLI DO NOT EDIT!!
package http

import (
	http "github.com/go-kit/kit/transport/http"
	mux "github.com/gorilla/mux"
	http1 "net/http"
	endpoint "pikabu-control/control/pkg/endpoint"
)

// NewHTTPHandler returns a handler that makes a set of endpoints available on
// predefined paths.
func NewHTTPHandler(endpoints endpoint.Endpoints, options map[string][]http.ServerOption) http1.Handler {
	m := mux.NewRouter()
	makeApiHandler(m, endpoints, options["Api"])
	makeRootHandler(m, endpoints, options["Root"])
	makeFileHandler(m, endpoints, options["File"])
	return m
}
