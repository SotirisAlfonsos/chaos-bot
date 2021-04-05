package v1

import (
	"context"

	"github.com/SotirisAlfonsos/chaos-bot/common/cpu"
	v1 "github.com/SotirisAlfonsos/chaos-bot/proto/grpc/v1"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"go.opencensus.io/trace"
)

// CPUManager is the rpc for cpu failures
type CPUManager struct {
	CPU    *cpu.CPU
	Logger log.Logger
	*v1.UnimplementedCPUServer
}

// Start a new cpu injection based on the percentage specified
func (cm *CPUManager) Start(ctx context.Context, req *v1.CPURequest) (*v1.StatusResponse, error) {
	ctx, span := trace.StartSpan(ctx, "v1.api.cpu.Start")
	defer span.End()

	resp := make(chan response, 1)

	go func() {
		resp <- cm.startCPU(req)
	}()

	select {
	case <-ctx.Done():
		<-resp
		_ = level.Warn(cm.Logger).Log("msg", "Context error encountered", "err", ctx.Err())
		return prepareResponse("", ctx.Err())
	case r := <-resp:
		return prepareResponse(r.message, r.err)
	}
}

// Recover the cpu injection
func (cm *CPUManager) Recover(ctx context.Context, _ *v1.CPURequest) (*v1.StatusResponse, error) {
	ctx, span := trace.StartSpan(ctx, "v1.api.cpu.Stop")
	defer span.End()

	resp := make(chan response, 1)

	go func() {
		resp <- cm.recoverCPU()
	}()

	select {
	case <-ctx.Done():
		<-resp
		_ = level.Warn(cm.Logger).Log("msg", "Context error encountered", "err", ctx.Err())
		return prepareResponse("", ctx.Err())
	case r := <-resp:
		return prepareResponse(r.message, r.err)
	}
}

func (cm *CPUManager) startCPU(req *v1.CPURequest) response {
	message, err := cm.CPU.Start(int(req.Percentage))

	return response{
		message: message,
		err:     err,
	}
}

func (cm *CPUManager) recoverCPU() response {
	message, err := cm.CPU.Recover()

	return response{
		message: message,
		err:     err,
	}
}
