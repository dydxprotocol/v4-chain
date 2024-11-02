import { ClobPair, ClobPairSDKType } from "./clob_pair";
import { LiquidationsConfig, LiquidationsConfigSDKType } from "./liquidations_config";
import { BlockRateLimitConfiguration, BlockRateLimitConfigurationSDKType } from "./block_rate_limit_config";
import { EquityTierLimitConfiguration, EquityTierLimitConfigurationSDKType } from "./equity_tier_limit_config";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** GenesisState defines the clob module's genesis state. */

export interface GenesisState {
  clobPairs: ClobPair[];
  liquidationsConfig?: LiquidationsConfig;
  blockRateLimitConfig?: BlockRateLimitConfiguration;
  equityTierLimitConfig?: EquityTierLimitConfiguration;
}
/** GenesisState defines the clob module's genesis state. */

export interface GenesisStateSDKType {
  clob_pairs: ClobPairSDKType[];
  liquidations_config?: LiquidationsConfigSDKType;
  block_rate_limit_config?: BlockRateLimitConfigurationSDKType;
  equity_tier_limit_config?: EquityTierLimitConfigurationSDKType;
}

function createBaseGenesisState(): GenesisState {
  return {
    clobPairs: [],
    liquidationsConfig: undefined,
    blockRateLimitConfig: undefined,
    equityTierLimitConfig: undefined
  };
}

export const GenesisState = {
  encode(message: GenesisState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.clobPairs) {
      ClobPair.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    if (message.liquidationsConfig !== undefined) {
      LiquidationsConfig.encode(message.liquidationsConfig, writer.uint32(18).fork()).ldelim();
    }

    if (message.blockRateLimitConfig !== undefined) {
      BlockRateLimitConfiguration.encode(message.blockRateLimitConfig, writer.uint32(26).fork()).ldelim();
    }

    if (message.equityTierLimitConfig !== undefined) {
      EquityTierLimitConfiguration.encode(message.equityTierLimitConfig, writer.uint32(34).fork()).ldelim();
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
          message.clobPairs.push(ClobPair.decode(reader, reader.uint32()));
          break;

        case 2:
          message.liquidationsConfig = LiquidationsConfig.decode(reader, reader.uint32());
          break;

        case 3:
          message.blockRateLimitConfig = BlockRateLimitConfiguration.decode(reader, reader.uint32());
          break;

        case 4:
          message.equityTierLimitConfig = EquityTierLimitConfiguration.decode(reader, reader.uint32());
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
    message.clobPairs = object.clobPairs?.map(e => ClobPair.fromPartial(e)) || [];
    message.liquidationsConfig = object.liquidationsConfig !== undefined && object.liquidationsConfig !== null ? LiquidationsConfig.fromPartial(object.liquidationsConfig) : undefined;
    message.blockRateLimitConfig = object.blockRateLimitConfig !== undefined && object.blockRateLimitConfig !== null ? BlockRateLimitConfiguration.fromPartial(object.blockRateLimitConfig) : undefined;
    message.equityTierLimitConfig = object.equityTierLimitConfig !== undefined && object.equityTierLimitConfig !== null ? EquityTierLimitConfiguration.fromPartial(object.equityTierLimitConfig) : undefined;
    return message;
  }

};