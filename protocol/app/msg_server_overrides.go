package app

import (
	"fmt"

	gogogrpc "github.com/cosmos/gogoproto/grpc"
	"google.golang.org/grpc"
)

const ibcTransferMsgServiceName = "ibc.applications.transfer.v1.Msg"

// msgServerOverrides replaces selected message service implementations as
// modules register them. Each service is still registered exactly once with
// the application's MsgServiceRouter.
type msgServerOverrides struct {
	delegate  gogogrpc.Server
	overrides map[string]any
}

var _ gogogrpc.Server = msgServerOverrides{}

func newMsgServerOverrides(
	delegate gogogrpc.Server,
	overrides map[string]any,
) gogogrpc.Server {
	if delegate == nil {
		panic("message server delegate must not be nil")
	}

	copiedOverrides := make(map[string]any, len(overrides))
	for serviceName, implementation := range overrides {
		if serviceName == "" {
			panic("message server override name must not be empty")
		}
		if implementation == nil {
			panic(fmt.Sprintf("message server override for %s must not be nil", serviceName))
		}
		copiedOverrides[serviceName] = implementation
	}

	return msgServerOverrides{
		delegate:  delegate,
		overrides: copiedOverrides,
	}
}

func (s msgServerOverrides) RegisterService(
	descriptor *grpc.ServiceDesc,
	implementation any,
) {
	if override, ok := s.overrides[descriptor.ServiceName]; ok {
		implementation = override
	}
	s.delegate.RegisterService(descriptor, implementation)
}
