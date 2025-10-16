import { PerpetualFeeParams, PerpetualFeeParamsSDKType } from "./params";
import { StakingTier, StakingTierSDKType } from "./staking_tier";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** GenesisState defines the feetiers module's genesis state. */

export interface GenesisState {
  /** The parameters for perpetual fees. */
  params?: PerpetualFeeParams;
  /** The staking tiers. */

  stakingTiers: StakingTier[];
}
/** GenesisState defines the feetiers module's genesis state. */

export interface GenesisStateSDKType {
  /** The parameters for perpetual fees. */
  params?: PerpetualFeeParamsSDKType;
  /** The staking tiers. */

  staking_tiers: StakingTierSDKType[];
}

function createBaseGenesisState(): GenesisState {
  return {
    params: undefined,
    stakingTiers: []
  };
}

export const GenesisState = {
  encode(message: GenesisState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      PerpetualFeeParams.encode(message.params, writer.uint32(10).fork()).ldelim();
    }

    for (const v of message.stakingTiers) {
      StakingTier.encode(v!, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGenesisState();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.params = PerpetualFeeParams.decode(reader, reader.uint32());
          break;

        case 2:
          message.stakingTiers.push(StakingTier.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<GenesisState>): GenesisState {
    const message = createBaseGenesisState();
    message.params = object.params !== undefined && object.params !== null ? PerpetualFeeParams.fromPartial(object.params) : undefined;
    message.stakingTiers = object.stakingTiers?.map(e => StakingTier.fromPartial(e)) || [];
    return message;
  }

};