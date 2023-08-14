package msgsender_test

import (
	"github.com/Shopify/sarama"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/msgsender"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMessage_AddHeader(t *testing.T) {
	tests := map[string]struct {
		// Input
		message msgsender.Message
		header  msgsender.MessageHeader

		// Expectation
		expectedMessage msgsender.Message
	}{
		"Adds header to message with nil headers": {
			message: msgsender.Message{
				Key:     []byte{0x0},
				Value:   []byte{0x1},
				Headers: nil,
			},
			header: msgsender.MessageHeader{
				Key:   []byte{0x2},
				Value: []byte{0x3},
			},
			expectedMessage: msgsender.Message{
				Key:   []byte{0x0},
				Value: []byte{0x1},
				Headers: []sarama.RecordHeader{
					{
						Key:   []byte{0x2},
						Value: []byte{0x3},
					},
				},
			},
		},
		"Adds header to to message with empty headers": {
			message: msgsender.Message{
				Key:     []byte{0x0},
				Value:   []byte{0x1},
				Headers: make([]sarama.RecordHeader, 0),
			},
			header: msgsender.MessageHeader{
				Key:   []byte{0x4},
				Value: []byte{0x5},
			},
			expectedMessage: msgsender.Message{
				Key:   []byte{0x0},
				Value: []byte{0x1},
				Headers: []sarama.RecordHeader{
					{
						Key:   []byte{0x4},
						Value: []byte{0x5},
					},
				},
			},
		},
		"Adds header to to message with existing headers": {
			message: msgsender.Message{
				Key:   []byte{0x0},
				Value: []byte{0x1},
				Headers: []sarama.RecordHeader{
					{
						Key:   []byte{0x6},
						Value: []byte{0x7},
					},
				},
			},
			header: msgsender.MessageHeader{
				Key:   []byte{0x8},
				Value: []byte{0x9},
			},
			expectedMessage: msgsender.Message{
				Key:   []byte{0x0},
				Value: []byte{0x1},
				Headers: []sarama.RecordHeader{
					{
						Key:   []byte{0x6},
						Value: []byte{0x7},
					},
					{
						Key:   []byte{0x8},
						Value: []byte{0x9},
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.expectedMessage, tc.message.AddHeader(tc.header))
		})
	}
}
