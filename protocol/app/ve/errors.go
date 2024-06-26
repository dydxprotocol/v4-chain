package ve

import (
	"fmt"
)

// PreBlockError is an error that is returned when the pre-block simulation fails.
type PreBlockError struct {
	Err error
}

func (e PreBlockError) Error() string {
	return fmt.Sprintf("finalize block error: %s", e.Err.Error())
}

func (e PreBlockError) Label() string {
	return "PreBlockError"
}

// ErrPanic is an error that is returned when a panic occurs in the ABCI handler.
type ErrPanic struct {
	Err error
}

func (e ErrPanic) Error() string {
	return fmt.Sprintf("panic: %s", e.Err.Error())
}

func (e ErrPanic) Label() string {
	return "Panic"
}

// OracleClientError is an error that is returned when the oracle client's response is invalid.
type OracleClientError struct {
	Err error
}

func (e OracleClientError) Error() string {
	return fmt.Sprintf("oracle client error: %s", e.Err.Error())
}

func (e OracleClientError) Label() string {
	return "OracleClientError"
}

// TransformPricesError is an error that is returned when there is a failure in attempting to transform the prices returned
// from the oracle server to the format expected by the validator set.
type TransformPricesError struct {
	Err error
}

func (e TransformPricesError) Error() string {
	return fmt.Sprintf("prices transform error: %s", e.Err.Error())
}

func (e TransformPricesError) Label() string {
	return "TransformPricesError"
}

// ValidateVoteExtensionError is an error that is returned when there is a failure in validating a vote extension.
type ValidateVoteExtensionError struct {
	Err error
}

func (e ValidateVoteExtensionError) Error() string {
	return fmt.Sprintf("validate vote extension error: %s", e.Err.Error())
}

func (e ValidateVoteExtensionError) Label() string {
	return "ValidateVoteExtensionError"
}
