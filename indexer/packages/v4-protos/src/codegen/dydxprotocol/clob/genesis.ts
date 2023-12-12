import { ClobPair, ClobPairAmino, ClobPairSDKType } from "./clob_pair";
import { LiquidationsConfig, LiquidationsConfigAmino, LiquidationsConfigSDKType } from "./liquidations_config";
import { BlockRateLimitConfiguration, BlockRateLimitConfigurationAmino, BlockRateLimitConfigurationSDKType } from "./block_rate_limit_config";
import { EquityTierLimitConfiguration, EquityTierLimitConfigurationAmino, EquityTierLimitConfigurationSDKType } from "./equity_tier_limit_config";
import { BinaryReader, BinaryWriter } from "../../binary";
/** GenesisState defines the clob module's genesis state. */
export interface GenesisState {
  clobPairs: ClobPair[];
  liquidationsConfig: LiquidationsConfig;
  blockRateLimitConfig: BlockRateLimitConfiguration;
  equityTierLimitConfig: EquityTierLimitConfiguration;
}
export interface GenesisStateProtoMsg {
  typeUrl: "/dydxprotocol.clob.GenesisState";
  value: Uint8Array;
}
/** GenesisState defines the clob module's genesis state. */
export interface GenesisStateAmino {
  clob_pairs?: ClobPairAmino[];
  liquidations_config?: LiquidationsConfigAmino;
  block_rate_limit_config?: BlockRateLimitConfigurationAmino;
  equity_tier_limit_config?: EquityTierLimitConfigurationAmino;
}
export interface GenesisStateAminoMsg {
  type: "/dydxprotocol.clob.GenesisState";
  value: GenesisStateAmino;
}
/** GenesisState defines the clob module's genesis state. */
export interface GenesisStateSDKType {
  clob_pairs: ClobPairSDKType[];
  liquidations_config: LiquidationsConfigSDKType;
  block_rate_limit_config: BlockRateLimitConfigurationSDKType;
  equity_tier_limit_config: EquityTierLimitConfigurationSDKType;
}
function createBaseGenesisState(): GenesisState {
  return {
    clobPairs: [],
    liquidationsConfig: LiquidationsConfig.fromPartial({}),
    blockRateLimitConfig: BlockRateLimitConfiguration.fromPartial({}),
    equityTierLimitConfig: EquityTierLimitConfiguration.fromPartial({})
  };
}
export const GenesisState = {
  typeUrl: "/dydxprotocol.clob.GenesisState",
  encode(message: GenesisState, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
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
  decode(input: BinaryReader | Uint8Array, length?: number): GenesisState {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<GenesisState>): GenesisState {
    const message = createBaseGenesisState();
    message.clobPairs = object.clobPairs?.map(e => ClobPair.fromPartial(e)) || [];
    message.liquidationsConfig = object.liquidationsConfig !== undefined && object.liquidationsConfig !== null ? LiquidationsConfig.fromPartial(object.liquidationsConfig) : undefined;
    message.blockRateLimitConfig = object.blockRateLimitConfig !== undefined && object.blockRateLimitConfig !== null ? BlockRateLimitConfiguration.fromPartial(object.blockRateLimitConfig) : undefined;
    message.equityTierLimitConfig = object.equityTierLimitConfig !== undefined && object.equityTierLimitConfig !== null ? EquityTierLimitConfiguration.fromPartial(object.equityTierLimitConfig) : undefined;
    return message;
  },
  fromAmino(object: GenesisStateAmino): GenesisState {
    const message = createBaseGenesisState();
    message.clobPairs = object.clob_pairs?.map(e => ClobPair.fromAmino(e)) || [];
    if (object.liquidations_config !== undefined && object.liquidations_config !== null) {
      message.liquidationsConfig = LiquidationsConfig.fromAmino(object.liquidations_config);
    }
    if (object.block_rate_limit_config !== undefined && object.block_rate_limit_config !== null) {
      message.blockRateLimitConfig = BlockRateLimitConfiguration.fromAmino(object.block_rate_limit_config);
    }
    if (object.equity_tier_limit_config !== undefined && object.equity_tier_limit_config !== null) {
      message.equityTierLimitConfig = EquityTierLimitConfiguration.fromAmino(object.equity_tier_limit_config);
    }
    return message;
  },
  toAmino(message: GenesisState): GenesisStateAmino {
    const obj: any = {};
    if (message.clobPairs) {
      obj.clob_pairs = message.clobPairs.map(e => e ? ClobPair.toAmino(e) : undefined);
    } else {
      obj.clob_pairs = [];
    }
    obj.liquidations_config = message.liquidationsConfig ? LiquidationsConfig.toAmino(message.liquidationsConfig) : undefined;
    obj.block_rate_limit_config = message.blockRateLimitConfig ? BlockRateLimitConfiguration.toAmino(message.blockRateLimitConfig) : undefined;
    obj.equity_tier_limit_config = message.equityTierLimitConfig ? EquityTierLimitConfiguration.toAmino(message.equityTierLimitConfig) : undefined;
    return obj;
  },
  fromAminoMsg(object: GenesisStateAminoMsg): GenesisState {
    return GenesisState.fromAmino(object.value);
  },
  fromProtoMsg(message: GenesisStateProtoMsg): GenesisState {
    return GenesisState.decode(message.value);
  },
  toProto(message: GenesisState): Uint8Array {
    return GenesisState.encode(message).finish();
  },
  toProtoMsg(message: GenesisState): GenesisStateProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.GenesisState",
      value: GenesisState.encode(message).finish()
    };
  }
};