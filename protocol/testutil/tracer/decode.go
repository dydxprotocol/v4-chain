package tracer

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const ReadOperation = "read"
const WriteOperation = "write"
const DeleteOperation = "delete"

type TraceOperation struct {
	Operation string                 `json:"operation"`
	Key       string                 `json:"key"`
	Value     string                 `json:"value"`
	Metadata  map[string]interface{} `json:"metadata"`
}

type TraceDecoder struct {
	bytes.Buffer
}

func (td *TraceDecoder) GetOperations() (operations []TraceOperation) {
	// Trim the string of trailing newlines.
	opString := strings.TrimSuffix(td.String(), "\n")
	if len(opString) == 0 {
		return operations
	}

	// Split the string on newlines and parse the JSON.
	operationStrings := strings.Split(opString, "\n")
	for _, operationString := range operationStrings {
		var operation = TraceOperation{}
		err := json.Unmarshal([]byte(operationString), &operation)
		if err != nil {
			panic(err)
		}

		operations = append(operations, operation)
	}

	return operations
}

func (td *TraceDecoder) GetWriteOperations() (operations []TraceOperation) {
	for _, operation := range td.GetOperations() {
		if operation.Operation == WriteOperation || operation.Operation == DeleteOperation {
			operations = append(operations, operation)
		}
	}

	return operations
}

func (td *TraceDecoder) RequireKeyPrefixWrittenInSequence(
	t testing.TB,
	keys []string,
) {
	writeOperations := td.GetWriteOperations()
	decodedWriteOperations := make([]string, 0, len(writeOperations))
	for _, operation := range writeOperations {
		s, err := base64.StdEncoding.DecodeString(operation.Key)
		require.NoError(t, err)

		decodedWriteOperations = append(decodedWriteOperations, string(s))
	}

	require.Len(
		t,
		decodedWriteOperations,
		len(keys),
		fmt.Sprintf(
			"Different number of write operations (%+v) performed than expected (%+v)",
			len(decodedWriteOperations),
			len(keys),
		),
	)

	for i, dwo := range decodedWriteOperations {
		require.True(
			t,
			strings.HasPrefix(
				dwo,
				keys[i],
			),
			fmt.Sprintf(
				"Keys were not written in sequence.\nExpected:\n%s\nFound:\n%s\nWrite index: %d",
				keys[i],
				dwo,
				i,
			),
		)
	}
}

func (td *TraceDecoder) RequireReadWriteInSequence(
	t testing.TB,
	expectedOperations []TraceOperation,
) {
	operations := td.GetOperations()
	require.Len(
		t,
		operations,
		len(expectedOperations),
		"Different number of operations performed than expected",
	)

	for i, op := range operations {
		require.Equal(t, expectedOperations[i].Operation, op.Operation)

		key, err := base64.StdEncoding.DecodeString(op.Key)
		require.NoError(t, err)
		require.Equal(t, expectedOperations[i].Key, string(key))

		value, err := base64.StdEncoding.DecodeString(op.Value)
		require.NoError(t, err)
		require.Equal(t, expectedOperations[i].Value, string(value))
	}
}

// RequireKeyPrefixesWritten checks if the list of key prefixes were
// written to. This should be used when the order of writes do not have
// a deterministic ordering, i.e when state is branched and written back to.
// TODO(CLOB-851) update this function to assert ordering of key prefixes within stores.
func (td *TraceDecoder) RequireKeyPrefixesWritten(
	t testing.TB,
	keys []string,
) {
	writeOperations := td.GetWriteOperations()
	decodedWriteOperations := make([]string, 0, len(writeOperations))

	for _, operation := range writeOperations {
		s, err := base64.StdEncoding.DecodeString(operation.Key)
		require.NoError(t, err)

		decodedWriteOperations = append(decodedWriteOperations, string(s))
	}

	for _, key := range keys {
		foundMatchingWriteOperation := false
		for _, dwo := range decodedWriteOperations {
			if strings.HasPrefix(
				dwo,
				key,
			) {
				foundMatchingWriteOperation = true
				break
			}
		}
		require.Truef(
			t,
			foundMatchingWriteOperation,
			"Did not find a matching write operation for key %+v",
			key,
		)
	}
}
