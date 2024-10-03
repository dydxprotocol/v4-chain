package client

import (
	"context"
	"fmt"

	"cosmossdk.io/log"
	"github.com/hashicorp/go-version"

	oracleclient "github.com/skip-mev/connect/v2/service/clients/oracle"
	"github.com/skip-mev/connect/v2/service/servers/oracle/types"
)

const (
	MinSidecarVersion = "v1.0.12"
)

// SidecarVersionChecker is a lightweight process run in a goroutine by the slinky client.
// Its purpose is to query the running sidecar version and check if it is at least a minimum
// acceptable version.
type SidecarVersionChecker interface {
	Start(ctx context.Context) error
	Stop()
	CheckSidecarVersion(context.Context) error
}

// SidecarVersionCheckerImpl implements the SidecarVersionChecker interface.
type SidecarVersionCheckerImpl struct {
	Logger log.Logger
	slinky oracleclient.OracleClient
}

func NewSidecarVersionChecker(logger log.Logger, slinky oracleclient.OracleClient) SidecarVersionChecker {
	return &SidecarVersionCheckerImpl{
		Logger: logger,
		slinky: slinky,
	}
}

// Start initializes the underlying connections of the SidecarVersionChecker.
func (s *SidecarVersionCheckerImpl) Start(
	ctx context.Context) error {
	return s.slinky.Start(ctx)
}

// Stop closes all existing connections.
func (s *SidecarVersionCheckerImpl) Stop() {
	_ = s.slinky.Stop()
}

func (p *SidecarVersionCheckerImpl) CheckSidecarVersion(ctx context.Context) error {
	// get prices from slinky sidecar via GRPC
	slinkyResponse, err := p.slinky.Version(ctx, &types.QueryVersionRequest{})
	if err != nil {
		return err
	}
	current, err := version.NewVersion(slinkyResponse.Version)
	if err != nil {
		return fmt.Errorf("failed to parse current version: %w", err)
	}

	minimum, err := version.NewVersion(MinSidecarVersion)
	if err != nil {
		return fmt.Errorf("failed to parse minimum version: %w", err)
	}

	// Compare versions
	if current.LessThan(minimum) {
		return fmt.Errorf("sidecar version %s is less than minimum required version %s", current, minimum)
	}

	// Version is acceptable
	p.Logger.Info("Sidecar version check passed", "version", current)
	return nil

}
