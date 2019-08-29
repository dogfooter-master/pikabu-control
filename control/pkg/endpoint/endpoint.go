package endpoint

import (
	"context"
	"pikabu-control/control/pkg/grpc/pb"
	service "pikabu-control/control/pkg/service"

	endpoint "github.com/go-kit/kit/endpoint"
)

// ApiRequest collects the request parameters for the Api method.
type ApiRequest struct {
	Req service.Payload `json:"data"`
}

// ApiResponse collects the response parameters for the Api method.
type ApiResponse struct {
	Res service.Payload `json:"data"`
	Err error           `json:"err"`
}

// MakeApiEndpoint returns an endpoint that invokes Api on the service.
func MakeApiEndpoint(s service.ControlService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ApiRequest)
		res, err := s.Api(ctx, req.Req)
		return ApiResponse{
			Err: err,
			Res: res,
		}, nil
	}
}

// Failed implements Failer.
func (r ApiResponse) Failed() error {
	return r.Err
}

// RootRequest collects the request parameters for the Root method.
type RootRequest struct {
	Req *pb.RootRequest `json:"req"`
}

// RootResponse collects the response parameters for the Root method.
type RootResponse struct {
	Res *pb.RootReply `json:"res"`
	Err error         `json:"err"`
}

// MakeRootEndpoint returns an endpoint that invokes Root on the service.
func MakeRootEndpoint(s service.ControlService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(RootRequest)
		res, err := s.Root(ctx, req.Req)
		return RootResponse{
			Err: err,
			Res: res,
		}, nil
	}
}

// Failed implements Failer.
func (r RootResponse) Failed() error {
	return r.Err
}

// FileRequest collects the request parameters for the File method.
type FileRequest struct {
	Req service.FileVars `json:"req"`
}

// FileResponse collects the response parameters for the File method.
type FileResponse struct {
	Res service.FileVars `json:"res"`
	Err error            `json:"err"`
}

// MakeFileEndpoint returns an endpoint that invokes File on the service.
func MakeFileEndpoint(s service.ControlService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(FileRequest)
		res, err := s.File(ctx, req.Req)
		return FileResponse{
			Err: err,
			Res: res,
		}, nil
	}
}

// Failed implements Failer.
func (r FileResponse) Failed() error {
	return r.Err
}

// Failure is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so they've
// failed, and if so encode them using a separate write path based on the error.
type Failure interface {
	Failed() error
}

// Api implements Service. Primarily useful in a client.
func (e Endpoints) Api(ctx context.Context, req service.Payload) (res service.Payload, err error) {
	request := ApiRequest{Req: req}
	response, err := e.ApiEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(ApiResponse).Res, response.(ApiResponse).Err
}

// Root implements Service. Primarily useful in a client.
func (e Endpoints) Root(ctx context.Context, req *pb.RootRequest) (res *pb.RootReply, err error) {
	request := RootRequest{Req: req}
	response, err := e.RootEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(RootResponse).Res, response.(RootResponse).Err
}

// File implements Service. Primarily useful in a client.
func (e Endpoints) File(ctx context.Context, req service.FileVars) (res service.FileVars, err error) {
	request := FileRequest{Req: req}
	response, err := e.FileEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(FileResponse).Res, response.(FileResponse).Err
}
