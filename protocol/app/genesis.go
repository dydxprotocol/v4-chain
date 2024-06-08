package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/basic_manager"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	consumerTypes "github.com/ethos-works/ethos/ethos-chain/x/ccv/consumer/types"
	"github.com/ethos-works/ethos/ethos-chain/x/ccv/types"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

// The genesis state of the blockchain is represented here as a map of raw json
// messages key'd by a identifier string.
// The identifier is used to determine which module genesis information belongs
// to so it may be appropriately routed during init chain.
// Within this application default genesis information is retrieved from
// the ModuleBasicManager which populates json from each BasicModule
// object provided to it during init.
type GenesisState map[string]json.RawMessage

// Map of supported versions for consumer genesis transformation
type IcsVersion string

const (
	v2_x   IcsVersion = "v2.x"
	v3_0_x IcsVersion = "v3.0.x"
	v3_1_x IcsVersion = "v3.1.x"
	v3_2_x IcsVersion = "v3.2.x"
	v3_3_x IcsVersion = "v3.3.x"
	v4_x_x IcsVersion = "v4.x"
)

var TransformationVersions map[string]IcsVersion = map[string]IcsVersion{
	"v2.x":   v2_x,
	"v3.0.x": v3_0_x,
	"v3.1.x": v3_1_x,
	"v3.2.x": v3_2_x,
	"v3.3.x": v3_3_x,
	"v4.x":   v4_x_x,
}

// Transformation of consumer genesis content as it is exported from a provider version v1,2,3
// to a format readable by current consumer implementation.
func transformToNew(jsonRaw []byte, ctx client.Context) (json.RawMessage, error) {
	// v1,2,3 uses deprecated fields of GenesisState type
	oldConsumerGenesis := consumerTypes.GenesisState{}
	err := ctx.Codec.UnmarshalJSON(jsonRaw, &oldConsumerGenesis)
	if err != nil {
		return nil, fmt.Errorf("reading consumer genesis data failed: %s", err)
	}

	initialValSet := oldConsumerGenesis.InitialValSet
	// transformation from >= v3.3.x
	if len(initialValSet) == 0 {
		initialValSet = oldConsumerGenesis.Provider.InitialValSet
	}

	clientState := oldConsumerGenesis.ProviderClientState
	if clientState == nil {
		clientState = oldConsumerGenesis.Provider.ClientState
	}

	consensusState := oldConsumerGenesis.ProviderConsensusState
	if consensusState == nil {
		consensusState = oldConsumerGenesis.Provider.ConsensusState
	}

	// Use DefaultRetryDelayPeriod if not set
	if oldConsumerGenesis.Params.RetryDelayPeriod == 0 {
		oldConsumerGenesis.Params.RetryDelayPeriod = types.DefaultRetryDelayPeriod
	}

	// Versions before v3.3.x of provider genesis data fills up deprecated fields
	// ProviderClientState, ConsensusState and InitialValSet in type GenesisState
	newGenesis := types.ConsumerGenesisState{
		Params: oldConsumerGenesis.Params,
		Provider: types.ProviderInfo{
			ClientState:    clientState,
			ConsensusState: consensusState,
			InitialValSet:  initialValSet,
		},
		NewChain: oldConsumerGenesis.NewChain,
	}

	newJson, err := ctx.Codec.MarshalJSON(&newGenesis)
	if err != nil {
		return nil, fmt.Errorf("failed marshalling data to new type: %s", err)
	}
	return newJson, nil
}

// Transformation of consumer genesis content as it is exported by current provider version
// to a format supported by consumer version v3.3.x
func transformToV33(jsonRaw []byte, ctx client.Context) ([]byte, error) {
	// v1,2,3 uses deprecated fields of GenesisState type
	srcConGen := consumerTypes.GenesisState{}
	err := ctx.Codec.UnmarshalJSON(jsonRaw, &srcConGen)
	if err != nil {
		return nil, fmt.Errorf("reading consumer genesis data failed: %s", err)
	}

	// Remove retry_delay_period from 'params'
	params, err := ctx.Codec.MarshalJSON(&srcConGen.Params)
	if err != nil {
		return nil, err
	}
	tmp := map[string]json.RawMessage{}
	if err := json.Unmarshal(params, &tmp); err != nil {
		return nil, fmt.Errorf("unmarshalling 'params' failed: %v", err)
	}
	_, exists := tmp["retry_delay_period"]
	if exists {
		delete(tmp, "retry_delay_period")
	}
	params, err = json.Marshal(tmp)
	if err != nil {
		return nil, err
	}

	// Marshal GenesisState and patch 'params' value
	result, err := ctx.Codec.MarshalJSON(&srcConGen)
	if err != nil {
		return nil, err
	}
	genState := map[string]json.RawMessage{}
	if err := json.Unmarshal(result, &genState); err != nil {
		return nil, fmt.Errorf("unmarshalling 'GenesisState' failed: %v", err)
	}
	genState["params"] = params

	result, err = json.Marshal(genState)
	if err != nil {
		return nil, fmt.Errorf("marshalling transformation result failed: %v", err)
	}
	return result, nil
}

// Transformation of consumer genesis content as it is exported from current provider version
// to a format readable by consumer implementation of version v2.x
// Use removePreHashKey to remove prehash_key_before_comparison from result.
func transformToV2(jsonRaw []byte, ctx client.Context, removePreHashKey bool) (json.RawMessage, error) {
	// populate deprecated fields of GenesisState used by version v2.x
	srcConGen := consumerTypes.GenesisState{}
	err := ctx.Codec.UnmarshalJSON(jsonRaw, &srcConGen)
	if err != nil {
		return nil, fmt.Errorf("reading consumer genesis data failed: %s", err)
	}

	// remove retry_delay_period from 'params' if present (introduced in v4.x)
	params, err := ctx.Codec.MarshalJSON(&srcConGen.Params)
	if err != nil {
		return nil, err
	}
	paramsMap := map[string]json.RawMessage{}
	if err := json.Unmarshal(params, &paramsMap); err != nil {
		return nil, fmt.Errorf("unmarshalling 'params' failed: %v", err)
	}
	_, exists := paramsMap["retry_delay_period"]
	if exists {
		delete(paramsMap, "retry_delay_period")
	}
	params, err = json.Marshal(paramsMap)
	if err != nil {
		return nil, err
	}

	// marshal GenesisState and patch 'params' value
	result, err := ctx.Codec.MarshalJSON(&srcConGen)
	if err != nil {
		return nil, err
	}
	genState := map[string]json.RawMessage{}
	if err := json.Unmarshal(result, &genState); err != nil {
		return nil, fmt.Errorf("unmarshalling 'GenesisState' failed: %v", err)
	}
	genState["params"] = params

	provider, err := ctx.Codec.MarshalJSON(&srcConGen.Provider)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling 'Provider' failed: %v", err)
	}
	providerMap := map[string]json.RawMessage{}
	if err := json.Unmarshal(provider, &providerMap); err != nil {
		return nil, fmt.Errorf("unmarshalling 'provider' failed: %v", err)
	}

	// patch .initial_val_set form .provider.initial_val_set if needed
	if len(srcConGen.Provider.InitialValSet) > 0 {
		valSet, exists := providerMap["initial_val_set"]
		if !exists {
			return nil, fmt.Errorf("'initial_val_set' not found in provider data")
		}
		_, exists = genState["initial_val_set"]
		if exists {
			genState["initial_val_set"] = valSet
		}
	}

	// patch .provider_consensus_state from provider.consensus_state if needed
	if srcConGen.Provider.ConsensusState != nil {
		valSet, exists := providerMap["consensus_state"]
		if !exists {
			return nil, fmt.Errorf("'consensus_state' not found in provider data")
		}
		_, exists = genState["provider_consensus_state"]
		if exists {
			genState["provider_consensus_state"] = valSet
		}
	}

	// patch .provider_client_state from provider.client_state if needed
	if srcConGen.Provider.ClientState != nil {
		clientState, exists := providerMap["client_state"]
		if !exists {
			return nil, fmt.Errorf("'client_state' not found in provider data")
		}
		_, exists = genState["provider_client_state"]
		if exists {
			genState["provider_client_state"] = clientState
		}
	}

	// delete .provider entry (introduced in v3.3.x)
	delete(genState, "provider")

	// Marshall final result
	result, err = json.Marshal(genState)
	if err != nil {
		return nil, fmt.Errorf("marshalling transformation result failed: %v", err)
	}

	if removePreHashKey {
		// remove all `prehash_key_before_comparison` entries not supported in v2.x (see ics23)
		re := regexp.MustCompile(`,\s*"prehash_key_before_comparison"\s*:\s*(false|true)`)
		result = re.ReplaceAll(result, []byte{})
	}
	return result, nil
}

// transformGenesis transforms ccv consumer genesis data to the specified target version
// Returns the transformed data or an error in case the transformation failed or the format is not supported by current implementation
func transformGenesis(ctx client.Context, targetVersion IcsVersion, jsonRaw []byte) (json.RawMessage, error) {
	var newConsumerGenesis json.RawMessage = nil
	var err error

	switch targetVersion {
	// v2.x, v3.0-v3.2 share same consumer genesis type
	case v2_x:
		newConsumerGenesis, err = transformToV2(jsonRaw, ctx, true)
	case v3_0_x, v3_1_x, v3_2_x:
		// same as v2 replacement without need of `prehash_key_before_comparison` removal
		newConsumerGenesis, err = transformToV2(jsonRaw, ctx, false)
	case v3_3_x:
		newConsumerGenesis, err = transformToV33(jsonRaw, ctx)
	case v4_x_x:
		newConsumerGenesis, err = transformToNew(jsonRaw, ctx)
	default:
		err = fmt.Errorf("unsupported target version '%s'. Run %s --help",
			targetVersion, version.AppName)
	}

	if err != nil {
		return nil, fmt.Errorf("transformation failed: %v", err)
	}
	return newConsumerGenesis, err
}

// Transform a consumer genesis json file exported from a given ccv provider version
// to a consumer genesis json format supported by current ccv consumer version or v2.x
// This allows user to patch consumer genesis of
//   - current implementation from exports of provider of < v3.3.x
//   - v2.x from exports of provider >= v3.2.x
//
// Result will be written to defined output.
func TransformConsumerGenesis(cmd *cobra.Command, args []string) error {
	sourceFile := args[0]
	jsonRaw, err := os.ReadFile(filepath.Clean(sourceFile))
	if err != nil {
		return err
	}

	clientCtx := client.GetClientContextFromCmd(cmd)
	version, err := cmd.Flags().GetString("to")
	if err != nil {
		return fmt.Errorf("error getting targetVersion %v", err)
	}
	targetVersion, exists := TransformationVersions[version]
	if !exists {
		return fmt.Errorf("unsupported target version '%s'", version)
	}

	// try to transform data to target format
	newConsumerGenesis, err := transformGenesis(clientCtx, targetVersion, jsonRaw)
	if err != nil {
		return err
	}

	bz, err := newConsumerGenesis.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed exporting new consumer genesis to JSON: %s", err)
	}

	sortedBz, err := sdk.SortJSON(bz)
	if err != nil {
		return fmt.Errorf("failed sorting transformed consumer genesis JSON: %s", err)
	}

	cmd.Println(string(sortedBz))
	return nil
}

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState(cdc codec.JSONCodec) GenesisState {
	return basic_manager.ModuleBasics.DefaultGenesis(cdc)
}

// GetConsumerGenesisTransformCmd transforms Consumer Genesis JSON content exported from a
// provider version v1,v2 or v3 to a JSON format supported by this consumer version.
func GetConsumerGenesisTransformCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transform [-to version] genesis-file",
		Short: "Transform CCV consumer genesis data exported to a specific target format",
		Long: strings.TrimSpace(
			fmt.Sprintf(`
Transform the consumer genesis data exported from a provider version v1,v2, v3, v4 to a specified consumer target version.
The result is printed to STDOUT.

Note: Content to be transformed is not the consumer genesis file itself but the exported content from provider chain which is used to patch the consumer genesis file!

Example:
$ %s transform /path/to/ccv_consumer_genesis.json
$ %s --to v2.x transform /path/to/ccv_consumer_genesis.json
`, version.AppName, version.AppName),
		),
		Args: cobra.RangeArgs(1, 2),
		RunE: TransformConsumerGenesis,
	}
	cmd.Flags().String("to", string(v4_x_x),
		fmt.Sprintf("target version for consumer genesis. Supported versions %s",
			maps.Keys(TransformationVersions)))
	return cmd
}
