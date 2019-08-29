package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	http1 "net/http"
	"os"
	endpoint "pikabu-control/control/pkg/endpoint"
	"time"

	http "github.com/go-kit/kit/transport/http"
	handlers "github.com/gorilla/handlers"
	mux "github.com/gorilla/mux"
)

// makeApiHandler creates the handler logic
func makeApiHandler(m *mux.Router, endpoints endpoint.Endpoints, options []http.ServerOption) {
	m.Methods("POST").Path("/api/").Handler(
		handlers.CORS(
			handlers.AllowedMethods([]string{"POST"}),
			handlers.AllowedOrigins([]string{"*"}),
		)(http.NewServer(endpoints.ApiEndpoint, decodeApiRequest, encodeApiResponse, options...)))
	m.Methods("POST").Path("/api").Handler(
		handlers.CORS(
			handlers.AllowedMethods([]string{"POST"}),
			handlers.AllowedOrigins([]string{"*"}),
		)(http.NewServer(endpoints.ApiEndpoint, decodeApiRequest, encodeApiResponse, options...)))
	m.Methods("POST").Path("/api/{category}/{service}").Handler(
		handlers.CORS(
			handlers.AllowedMethods([]string{"POST"}),
			handlers.AllowedOrigins([]string{"*"}),
		)(http.NewServer(endpoints.ApiEndpoint, decodeApiRequest, encodeApiResponse, options...)))

	m.Methods("OPTIONS").Path("/api/").Handler(
		handlers.CORS(
			handlers.AllowedMethods([]string{"OPTIONS", "POST"}),
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedHeaders([]string{"content-type"}),
		)(http.NewServer(endpoints.ApiEndpoint, decodeApiRequest, encodeApiResponse, options...)))
	m.Methods("OPTIONS").Path("/api").Handler(
		handlers.CORS(
			handlers.AllowedMethods([]string{"OPTIONS", "POST"}),
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedHeaders([]string{"content-type"}),
		)(http.NewServer(endpoints.ApiEndpoint, decodeApiRequest, encodeApiResponse, options...)))
	m.Methods("OPTIONS").Path("/api/{category}/{service}").Handler(
		handlers.CORS(
			handlers.AllowedMethods([]string{"OPTIONS", "POST"}),
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedHeaders([]string{"content-type"}),
		)(http.NewServer(endpoints.ApiEndpoint, decodeApiRequest, encodeApiResponse, options...)))}

// decodeApiRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeApiRequest(_ context.Context, r *http1.Request) (interface{}, error) {
	vars := mux.Vars(r)
	req := endpoint.ApiRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if len(req.Req.Category) == 0 && len(vars["category"]) > 0 {
		req.Req.Category = vars["category"]
	}
	if len(req.Req.Service) == 0 && len(vars["service"]) > 0 {
		req.Req.Service = vars["service"]
	}
	return req, err
}

// encodeApiResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeApiResponse(ctx context.Context, w http1.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}
func TestFunc(w http1.ResponseWriter, r *http1.Request) {
	filePath := os.Getenv("PIKABU_FRONTEND_HOME") + r.URL.Path
	//fmt.Fprintf(os.Stderr, "DEBUG: %v\n", filePath)
	if _, err := os.Stat(filePath); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		fmt.Fprintf(w, "404 Not Found Page !! : %v%v", r.Host, r.URL.Path)
		return
	}

	http1.ServeFile(w, r, filePath)
}
func WebSocketFunc(w http1.ResponseWriter, r *http1.Request) {
	service.ServeWebSocket(w, r)
}
// makeRootHandler creates the handler logic
func makeRootHandler(m *mux.Router, endpoints endpoint.Endpoints, options []http.ServerOption) {
	m.NotFoundHandler = http1.HandlerFunc(TestFunc)
	m.HandleFunc("/ws", WebSocketFunc)
	//m.Methods("GET").Path("/ws").HandlerFunc(WebSocketFunc)
	m.Methods("GET").Path("/").Handler(http1.FileServer(http1.Dir(os.Getenv("PIKABU_FRONTEND_HOME"))))
}

// decodeRootRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeRootRequest(_ context.Context, r *http1.Request) (interface{}, error) {
	req := endpoint.RootRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// encodeRootResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeRootResponse(ctx context.Context, w http1.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}

// makeFileHandler creates the handler logic
func makeFileHandler(m *mux.Router, endpoints endpoint.Endpoints, options []http.ServerOption) {
	m.Methods("GET").Path("/files/{access_token}/{sep1}/{id}/{sep2}/{sep3}/{date}/{file}").Handler(handlers.CORS(handlers.AllowedMethods([]string{"GET"}), handlers.AllowedOrigins([]string{"*"}))(http.NewServer(endpoints.FileEndpoint, decodeFileRequest, encodeFileResponse, options...)))
	m.Methods("POST").Path("/files/{access_token}/{sep1}/{id}/{sep2}/{sep3}/{date}/{file}").Handler(handlers.CORS(handlers.AllowedMethods([]string{"POST"}), handlers.AllowedOrigins([]string{"*"}))(http.NewServer(endpoints.FileEndpoint, decodeFileRequest, encodeFileResponse, options...)))
}

// decodeFileRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeFileRequest(_ context.Context, r *http1.Request) (interface{}, error) {
	vars := mux.Vars(r)

	f := service.FileVars{
		Method:      r.Method,
		Id:          vars["id"],
		AccessToken: vars["access_token"],
		File:        vars["file"],
	}
	f.Date, _ = time.Parse("20060102150405", vars["date"])
	f.Separators = append(f.Separators, vars["sep1"])
	f.Separators = append(f.Separators, vars["sep2"])
	f.Separators = append(f.Separators, vars["sep3"])

	f.R = r
	req := endpoint.FileRequest{
		Req: f,
	}
	return req, err
}

// encodeFileResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeFileResponse(ctx context.Context, w http1.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		w.WriteHeader(http1.StatusUnauthorized)
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	res := response.(endpoint.FileResponse)
	if res.Res.Method == "GET" {
		http1.ServeFile(w, res.Res.R, res.Res.StaticFilePath)
	} else if res.Res.Method == "POST" {
		if err = service.FilePost(res.Res.R, res.Res.StaticFilePath); err != nil {
			fmt.Fprintf(os.Stderr, "file post error: %v\n", err)
			return
		}
	}

	return
}
func ErrorEncoder(_ context.Context, err error, w http1.ResponseWriter) {
	w.WriteHeader(err2code(err))
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}
func ErrorDecoder(r *http1.Response) error {
	var w errorWrapper
	if err := json.NewDecoder(r.Body).Decode(&w); err != nil {
		return err
	}
	return errors.New(w.Error)
}

// This is used to set the http status, see an example here :
// https://github.com/go-kit/kit/blob/master/examples/addsvc/pkg/addtransport/http.go#L133
func err2code(err error) int {
	return http1.StatusInternalServerError
}

type errorWrapper struct {
	Error string `json:"error"`
}
