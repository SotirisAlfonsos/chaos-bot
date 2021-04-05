package v1

import (
	"context"

	"github.com/SotirisAlfonsos/chaos-bot/common/network"
	v1 "github.com/SotirisAlfonsos/chaos-bot/proto/grpc/v1"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"go.opencensus.io/trace"
)

// NetworkManager is the rpc for network failure injections
type NetworkManager struct {
	Network *network.Network
	Logger  log.Logger
	*v1.UnimplementedNetworkServer
}

// NewNetworkManager will create the rpc for network failures with a logger attached
func NewNetworkManager(logger log.Logger) *NetworkManager {
	return &NetworkManager{
		Network: network.New(logger),
		Logger:  logger,
	}
}

// Start a network injection according to the network request
func (nm *NetworkManager) Start(ctx context.Context, req *v1.NetworkRequest) (*v1.StatusResponse, error) {
	ctx, span := trace.StartSpan(ctx, "v1.api.network.Start")
	defer span.End()

	resp := make(chan response, 1)

	go func() {
		resp <- nm.startNetwork(req)
	}()

	select {
	case <-ctx.Done():
		<-resp
		_ = level.Warn(nm.Logger).Log("msg", "Context error encountered", "err", ctx.Err())
		return prepareResponse("", ctx.Err())
	case r := <-resp:
		return prepareResponse(r.message, r.err)
	}
}

// Recover a network injection
func (nm *NetworkManager) Recover(ctx context.Context, req *v1.NetworkRequest) (*v1.StatusResponse, error) {
	ctx, span := trace.StartSpan(ctx, "v1.api.network.Stop")
	defer span.End()

	resp := make(chan response, 1)

	go func() {
		resp <- nm.recoverNetwork(req.Device)
	}()

	select {
	case <-ctx.Done():
		<-resp
		_ = level.Warn(nm.Logger).Log("msg", "Context error encountered", "err", ctx.Err())
		return prepareResponse("", ctx.Err())
	case r := <-resp:
		return prepareResponse(r.message, r.err)
	}
}

func (nm *NetworkManager) startNetwork(req *v1.NetworkRequest) response {
	message, err := nm.Network.Start(req)

	return response{
		message: message,
		err:     err,
	}
}

func (nm *NetworkManager) recoverNetwork(device string) response {
	message, err := nm.Network.Recover(device)

	return response{
		message: message,
		err:     err,
	}
}
