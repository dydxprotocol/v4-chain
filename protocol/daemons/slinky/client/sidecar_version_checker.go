package client

import (
	"context"
	"fmt"
	"strings"

	"cosmossdk.io/log"
	"github.com/hashicorp/go-version"

	oracleclient "github.com/dydxprotocol/slinky/service/clients/oracle"
	"github.com/dydxprotocol/slinky/service/servers/oracle/types"
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
	slinky oracleclient.OracleClient
	logger log.Logger
}

func NewSidecarVersionChecker(slinky oracleclient.OracleClient, logger log.Logger) SidecarVersionChecker {
	return &SidecarVersionCheckerImpl{
		slinky: slinky,
		logger: logger,
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

func (s *SidecarVersionCheckerImpl) CheckSidecarVersion(ctx context.Context) error {
	// Retrieve sidecar version via gRPC
	slinkyResponse, err := s.slinky.Version(ctx, &types.QueryVersionRequest{})
	if err != nil {
		return err
	}

	versionStr := slinkyResponse.Version
	if idx := strings.LastIndex(versionStr, "/"); idx != -1 {
		versionStr = versionStr[idx+1:]
	}

	current, err := version.NewVersion(versionStr)
	if err != nil {
		return fmt.Errorf("failed to parse current version: %w", err)
	}

	minimum, err := version.NewVersion(MinSidecarVersion)
	if err != nil {
		return fmt.Errorf("failed to parse minimum version: %w", err)
	}

	// Compare versions
	if current.LessThan(minimum) {
		return fmt.Errorf("Sidecar version %s is less than minimum required version %s. "+
			"The node will shut down soon", current, minimum)
	}

	// Version is acceptable
	s.logger.Info("Sidecar version check passed", "version", current)
	return nil
}
