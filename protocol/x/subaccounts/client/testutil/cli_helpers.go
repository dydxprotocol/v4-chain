package testutil

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/cosmos/cosmos-sdk/testutil"
)

// MsgQuerySubaccountExec executes a query for the given subaccount id.
func MsgQuerySubaccountExec(
	owner string,
	number uint32,
) (testutil.BufferWriter, error) {

	queryCmd := exec.Command("bash", "-c", "docker exec interchain-security-instance interchain-security-cd query subaccounts show-subaccount "+owner+" "+fmt.Sprint(number)+" --node tcp://7.7.8.4:26658 -o json")
	var transferOut bytes.Buffer
	var stdTransferErr bytes.Buffer
	queryCmd.Stdout = &transferOut
	queryCmd.Stderr = &stdTransferErr
	err := queryCmd.Run()
	if err != nil {
		return nil, err
	}
	return &transferOut, nil
}
