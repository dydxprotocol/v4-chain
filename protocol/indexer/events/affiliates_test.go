package events

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRegisterAffiliateEventV1_Success(t *testing.T) {
	registerAffiliateEvent := NewRegisterAffiliateEventV1(
		"referee",
		"affiliate",
	)
	expectedRegisterAffiliateEventProto := &RegisterAffiliateEventV1{
		Referee:   "referee",
		Affiliate: "affiliate",
	}
	require.Equal(t, expectedRegisterAffiliateEventProto, registerAffiliateEvent)
}
