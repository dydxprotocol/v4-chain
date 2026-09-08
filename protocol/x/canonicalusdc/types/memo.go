package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type FundingMemo struct {
	Version        uint32 `json:"version"`
	Action         string `json:"action"`
	Controller     string `json:"controller"`
	NobleRecipient string `json:"noble_recipient"`
}

type canonicalMemoEnvelope struct {
	CanonicalUsdc *FundingMemo `json:"canonical_usdc"`
}

func ParseFundingMemo(memo string) (FundingMemo, bool, error) {
	if memo == "" {
		return FundingMemo{}, false, nil
	}
	if len(memo) > MaxMemoLength {
		return FundingMemo{}, false, ErrInvalidFundingMemo
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(memo), &fields); err != nil {
		return FundingMemo{}, false, nil
	}
	if _, ok := fields["canonical_usdc"]; !ok {
		return FundingMemo{}, false, nil
	}
	decoder := json.NewDecoder(bytes.NewBufferString(memo))
	decoder.DisallowUnknownFields()
	var envelope canonicalMemoEnvelope
	if err := decoder.Decode(&envelope); err != nil {
		return FundingMemo{}, true, fmt.Errorf("%w: %v", ErrInvalidFundingMemo, err)
	}
	if err := ensureJSONEOF(decoder); err != nil {
		return FundingMemo{}, true, err
	}
	if envelope.CanonicalUsdc == nil {
		return FundingMemo{}, true, ErrInvalidFundingMemo
	}
	return *envelope.CanonicalUsdc, true, nil
}

func ensureJSONEOF(decoder *json.Decoder) error {
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return fmt.Errorf("%w: multiple JSON values", ErrInvalidFundingMemo)
		}
		return fmt.Errorf("%w: %v", ErrInvalidFundingMemo, err)
	}
	return nil
}
