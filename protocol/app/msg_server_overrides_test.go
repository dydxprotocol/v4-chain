package app

import (
	"context"
	"testing"

	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

type recordedService struct {
	descriptor     *grpc.ServiceDesc
	implementation any
}

type recordingMsgServer struct {
	services []recordedService
}

type transferMsgServer struct {
	name string
}

var _ ibctransfertypes.MsgServer = (*transferMsgServer)(nil)

func (*transferMsgServer) Transfer(
	context.Context,
	*ibctransfertypes.MsgTransfer,
) (*ibctransfertypes.MsgTransferResponse, error) {
	return nil, nil
}

func (*transferMsgServer) UpdateParams(
	context.Context,
	*ibctransfertypes.MsgUpdateParams,
) (*ibctransfertypes.MsgUpdateParamsResponse, error) {
	return nil, nil
}

func (s *recordingMsgServer) RegisterService(
	descriptor *grpc.ServiceDesc,
	implementation any,
) {
	s.services = append(s.services, recordedService{
		descriptor:     descriptor,
		implementation: implementation,
	})
}

func TestMsgServerOverrides(t *testing.T) {
	recorder := &recordingMsgServer{}
	transferImplementation := &transferMsgServer{name: "transfer"}
	originalTransferImplementation := &transferMsgServer{name: "original"}
	server := newMsgServerOverrides(recorder, map[string]any{
		ibcTransferMsgServiceName: transferImplementation,
	})

	ibctransfertypes.RegisterMsgServer(server, originalTransferImplementation)

	originalImplementation := &struct{ name string }{name: "original"}
	otherDescriptor := &grpc.ServiceDesc{ServiceName: "dydxprotocol.test.Msg"}
	server.RegisterService(otherDescriptor, originalImplementation)

	require.Len(t, recorder.services, 2)
	require.Equal(t, ibcTransferMsgServiceName, recorder.services[0].descriptor.ServiceName)
	require.Same(t, transferImplementation, recorder.services[0].implementation)
	require.Same(t, otherDescriptor, recorder.services[1].descriptor)
	require.Same(t, originalImplementation, recorder.services[1].implementation)
}

func TestMsgServerOverridesCopiesOverrides(t *testing.T) {
	recorder := &recordingMsgServer{}
	originalImplementation := &struct{ name string }{name: "original"}
	overrideImplementation := &struct{ name string }{name: "override"}
	overrides := map[string]any{
		ibcTransferMsgServiceName: overrideImplementation,
	}
	server := newMsgServerOverrides(recorder, overrides)
	delete(overrides, ibcTransferMsgServiceName)

	server.RegisterService(
		&grpc.ServiceDesc{ServiceName: ibcTransferMsgServiceName},
		originalImplementation,
	)

	require.Len(t, recorder.services, 1)
	require.Same(t, overrideImplementation, recorder.services[0].implementation)
}
