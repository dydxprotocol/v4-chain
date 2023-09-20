package msgs_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/app/msgs"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/stretchr/testify/require"
)

func TestAllTypeMessages(t *testing.T) {
	allTypes := make([]string, 0)

	// Use reflect.
	app := testapp.DefaultTestApp(nil)
	interfaceRegistry := app.InterfaceRegistry()
	typeUrl := reflect.ValueOf(interfaceRegistry).Elem().FieldByName("typeURLMap")
	if !typeUrl.IsValid() || typeUrl.Kind() != reflect.Map {
		require.Fail(t, "typeURLMap is not a map")
	}

	// Grab all typeURLs.
	for _, key := range typeUrl.MapKeys() {
		keyStr := key.String()
		allTypes = append(allTypes, keyStr)
	}

	// Sorting is needed since MapKeys() returns a random order.
	sort.Strings(allTypes)

	// Assert.
	require.Equal(t, allTypes, lib.GetSortedKeys[sort.StringSlice](msgs.AllTypeMessages))
}

func TestAllTypeMessages_SumOfDistinctLists(t *testing.T) {
	// The following should fail if there's a duplicate message in any of the lists.
	expectedAllTypeMsgs := lib.MergeAllMapsMustHaveDistinctKeys(
		msgs.AppInjectedMsgSamples,
		msgs.InternalMsgSamplesAll,
		msgs.NestedMsgSamples,
		msgs.UnsupportedMsgSamples,
		msgs.NormalMsgs,
	)
	require.Equal(
		t,
		lib.GetSortedKeys[sort.StringSlice](expectedAllTypeMsgs),
		lib.GetSortedKeys[sort.StringSlice](msgs.AllTypeMessages),
	)
}

func TestAllTypeMessages_EachMsgBelongsToSingleListOnly(t *testing.T) {
	// Each message must be in exactly one of the following lists:
	// 	a. app-injected msg list
	// 	b. internal-only msg list
	// 	c. nested msg list
	//  d. unsupported msg list
	// 	e. normal msg list
	for k := range msgs.AllTypeMessages {
		listCount := 0

		if _, isAppInjectedMsg := msgs.AppInjectedMsgSamples[k]; isAppInjectedMsg {
			listCount++
		}
		if _, isInternalOnlyMsg := msgs.InternalMsgSamplesAll[k]; isInternalOnlyMsg {
			listCount++
		}
		if _, isNestedMsg := msgs.NestedMsgSamples[k]; isNestedMsg {
			listCount++
		}
		if _, isUnsupportedMsg := msgs.UnsupportedMsgSamples[k]; isUnsupportedMsg {
			listCount++
		}
		if _, isNormalMsg := msgs.NormalMsgs[k]; isNormalMsg {
			listCount++
		}
		require.Equal(t, 1, listCount, "msg %s belongs to %d lists", k, listCount)
	}
}

func TestAllTypeMessages_NoUnregisteredMsg(t *testing.T) {
	for k := range msgs.AllTypeMessages {
		_, isUnregisteredMsg := msgs.UnregisteredMsgs[k]
		require.False(t, isUnregisteredMsg, "msg %s is unregistered", k)
	}
}

func TestDisallowMsgs(t *testing.T) {
	expectedDisallowMsgs := lib.MergeAllMapsMustHaveDistinctKeys(
		msgs.AppInjectedMsgSamples,
		msgs.InternalMsgSamplesAll,
		msgs.NestedMsgSamples,
		msgs.UnsupportedMsgSamples,
	)
	require.Equal(t, expectedDisallowMsgs, msgs.DisallowMsgs)
}

func TestAllowMsgs(t *testing.T) {
	require.Equal(t, msgs.NormalMsgs, msgs.AllowMsgs)
}

func TestAllTypeMessages_SumOfAllowDisallow_MinusUnregistered(t *testing.T) {
	expectedAllTypeMsgs := lib.MergeAllMapsMustHaveDistinctKeys(
		msgs.DisallowMsgs,
		msgs.AllowMsgs,
	)

	for k := range msgs.UnregisteredMsgs {
		_, exists := msgs.AllTypeMessages[k]
		require.False(t, exists, "msg %s is unregistered", k)
	}
	require.Equal(
		t,
		lib.GetSortedKeys[sort.StringSlice](expectedAllTypeMsgs),
		lib.GetSortedKeys[sort.StringSlice](msgs.AllTypeMessages),
	)
}
