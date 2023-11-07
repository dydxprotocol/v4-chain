package msgs

import (
	"sort"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TestBlockWithMsgs struct {
	Block uint32
	Msgs  []TestSdkMsg
}

type TestSdkMsg struct {
	Msg              sdk.Msg
	ExpectedIsOk     bool
	ExpectedRespCode uint32
}

// SampleMsg is a struct containing a sample msg and its name.
type SampleMsg struct {
	Name string
	Msg  sdk.Msg
}

// GetMsgNameWithModuleVersion returns the name of the msg type along with module and version, given its type url.
func GetMsgNameWithModuleVersion(typeUrl string) string {
	tokens := strings.Split(typeUrl, ".")
	if !IsValidMsgFormat(tokens) {
		panic("invalid type url: " + typeUrl)
	}

	lastThreeTokens := tokens[len(tokens)-3:]
	result := strings.Join(lastThreeTokens, ".")
	return result
}

// IsValidMsgFormat returns true if the given tokens are of the form: "<module>.<version>.Msg<MsgName>"
func IsValidMsgFormat(tokens []string) bool {
	tokenLen := len(tokens)
	if tokenLen < 3 {
		return false
	}

	if tokens[tokenLen-2] == "" || tokens[tokenLen-3] == "" {
		return false
	}

	return strings.HasPrefix(tokens[tokenLen-1], "Msg")
}

// GetSampleMsgs returns a list of sample msgs for each non-nil map value in the input.
func GetNonNilSampleMsgs(typeUrlToSampleMsgMap map[string]sdk.Msg) []SampleMsg {
	sampleMsgs := make([]SampleMsg, 0)
	for name, sample := range typeUrlToSampleMsgMap {
		if sample != nil {
			shortName := GetMsgNameWithModuleVersion(name)
			sampleMsgs = append(sampleMsgs, SampleMsg{shortName, sample})
		}
	}
	sort.Slice(sampleMsgs, func(i, j int) bool {
		return sampleMsgs[i].Name < sampleMsgs[j].Name
	})
	return sampleMsgs
}

// CopyMap returns a copy of the given map.
func CopyMap(m map[string]sdk.Msg) map[string]sdk.Msg {
	result := make(map[string]sdk.Msg)
	for k, v := range m {
		result[k] = v
	}
	return result
}
