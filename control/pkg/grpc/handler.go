package grpc

import (
	"context"
	"errors"
	grpc "github.com/go-kit/kit/transport/grpc"
	context1 "golang.org/x/net/context"
	endpoint "pikabu-control/control/pkg/endpoint"
	pb "pikabu-control/control/pkg/grpc/pb"
)

// makeApiHandler creates the handler logic
func makeApiHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.ApiEndpoint, decodeApiRequest, encodeApiResponse, options...)
}

// decodeApiResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain sum request.
// TODO implement the decoder
func decodeApiRequest(_ context.Context, r interface{}) (interface{}, error) {
	return nil, errors.New("'Control' Decoder is not impelemented")
}

// encodeApiResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
// TODO implement the encoder
func encodeApiResponse(_ context.Context, r interface{}) (interface{}, error) {
	return nil, errors.New("'Control' Encoder is not impelemented")
}
func (g *grpcServer) Api(ctx context1.Context, req *pb.ApiRequest) (*pb.ApiReply, error) {
	_, rep, err := g.api.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ApiReply), nil
}

// makeRootHandler creates the handler logic
func makeRootHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.RootEndpoint, decodeRootRequest, encodeRootResponse, options...)
}

// decodeRootResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain sum request.
// TODO implement the decoder
func decodeRootRequest(_ context.Context, r interface{}) (interface{}, error) {
	return nil, errors.New("'Control' Decoder is not impelemented")
}

// encodeRootResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
// TODO implement the encoder
func encodeRootResponse(_ context.Context, r interface{}) (interface{}, error) {
	return nil, errors.New("'Control' Encoder is not impelemented")
}
func (g *grpcServer) Root(ctx context1.Context, req *pb.RootRequest) (*pb.RootReply, error) {
	_, rep, err := g.root.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.RootReply), nil
}

// makeFileHandler creates the handler logic
func makeFileHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.FileEndpoint, decodeFileRequest, encodeFileResponse, options...)
}

// decodeFileResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain sum request.
// TODO implement the decoder
func decodeFileRequest(_ context.Context, r interface{}) (interface{}, error) {
	return nil, errors.New("'Control' Decoder is not impelemented")
}

// encodeFileResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
// TODO implement the encoder
func encodeFileResponse(_ context.Context, r interface{}) (interface{}, error) {
	return nil, errors.New("'Control' Encoder is not impelemented")
}
func (g *grpcServer) File(ctx context1.Context, req *pb.FileRequest) (*pb.FileReply, error) {
	_, rep, err := g.file.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.FileReply), nil
}
