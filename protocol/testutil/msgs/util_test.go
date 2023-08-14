package msgs_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/msgs"

	"github.com/stretchr/testify/require"
)

const (
	exampleMsgTypeUrl = "/cosmos.bank.v1beta1.MsgSend"
)

func TestGetMsgNameWithModuleVersion(t *testing.T) {
	tests := map[string]struct {
		input          string
		expectedResult string
		expectedPanic  string
	}{
		"Invalid: empty input": {
			input:         "",
			expectedPanic: "invalid type url: ",
		},
		"Invalid: only 1 tokens": {
			input:         "token1.",
			expectedPanic: "invalid type url: doesNotMatter.",
		},
		"Invalid: only 2 tokens": {
			input:         "token1.token2",
			expectedPanic: "invalid type url: token1.token2.",
		},
		"Invalid: empty last token": {
			input:         "token1.token2.",
			expectedPanic: "invalid type url: token1.token2.",
		},
		"Invalid: non-msg last token": {
			input:         "token1.token2.NotMsgPrefix",
			expectedPanic: "invalid type url: token1.token2.NotMsgPrefix",
		},
		"Valid token": {
			input:          exampleMsgTypeUrl,
			expectedResult: "bank.v1beta1.MsgSend",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.expectedPanic != "" {
				require.PanicsWithValue(t, "invalid type url: "+tc.input, func() {
					msgs.GetMsgNameWithModuleVersion(tc.input)
				})
			} else {
				require.Equal(t, tc.expectedResult, msgs.GetMsgNameWithModuleVersion(tc.input))
			}
		})
	}
}

func TestGetNonNilSampleMsgs(t *testing.T) {
	tests := map[string]struct {
		input          map[string]sdk.Msg
		expectedResult []msgs.SampleMsg
	}{
		"Empty input": {
			input:          map[string]sdk.Msg{},
			expectedResult: []msgs.SampleMsg{},
		},
		"All nil value input": {
			input: map[string]sdk.Msg{
				"token1":          nil,
				exampleMsgTypeUrl: nil,
			},
			expectedResult: []msgs.SampleMsg{},
		},
		"Valid input": {
			input: map[string]sdk.Msg{
				"token1":                               nil,
				"/cosmos.bank.v1beta1.MsgSendResponse": nil,
				exampleMsgTypeUrl:                      &testdata.TestMsg{Signers: []string{"meh"}},
			},
			expectedResult: []msgs.SampleMsg{
				{
					Name: "bank.v1beta1.MsgSend",
					Msg:  &testdata.TestMsg{Signers: []string{"meh"}},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := msgs.GetNonNilSampleMsgs(tc.input)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}
