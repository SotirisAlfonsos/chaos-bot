package v1

import (
	"context"

	"github.com/SotirisAlfonsos/chaos-bot/common/server"
	v1 "github.com/SotirisAlfonsos/chaos-bot/proto/grpc/v1"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"go.opencensus.io/trace"
)

type ServerManager struct {
	Server server.Server
	Logger log.Logger
	*v1.UnimplementedServerServer
}

func (sm *ServerManager) Stop(ctx context.Context, _ *v1.ServerRequest) (*v1.StatusResponse, error) {
	ctx, span := trace.StartSpan(ctx, "v1.api.server.Stop")
	defer span.End()

	resp := make(chan response, 1)

	go func() {
		resp <- sm.stop()
	}()

	select {
	case <-ctx.Done():
		<-resp
		_ = level.Warn(sm.Logger).Log("msg", "Context error encountered", "err", ctx.Err())
		return prepareResponse("", ctx.Err())
	case r := <-resp:
		return prepareResponse(r.message, r.err)
	}
}

func (sm *ServerManager) stop() response {
	message, err := sm.Server.StopUnix()

	return response{
		message: message,
		err:     err,
	}
}
