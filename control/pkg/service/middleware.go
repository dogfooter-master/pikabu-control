package service

import (
	"context"
	"pikabu-control/control/pkg/grpc/pb"

	log "github.com/go-kit/kit/log"
)

type Middleware func(ControlService) ControlService

type loggingMiddleware struct {
	logger log.Logger
	next   ControlService
}

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next ControlService) ControlService {
		return &loggingMiddleware{logger, next}
	}

}

func (l loggingMiddleware) Api(ctx context.Context, req Payload) (res Payload, err error) {
	defer func() {
		l.logger.Log("method", "Api", "req", req, "res", res, "err", err)
	}()
	return l.next.Api(ctx, req)
}
func (l loggingMiddleware) Root(ctx context.Context, req *pb.RootRequest) (res *pb.RootReply, err error) {
	defer func() {
		l.logger.Log("method", "Root", "req", req, "res", res, "err", err)
	}()
	return l.next.Root(ctx, req)
}
func (l loggingMiddleware) File(ctx context.Context, req FileVars) (res FileVars, err error) {
	defer func() {
		l.logger.Log("method", "File", "req", req, "res", res, "err", err)
	}()
	return l.next.File(ctx, req)
}
