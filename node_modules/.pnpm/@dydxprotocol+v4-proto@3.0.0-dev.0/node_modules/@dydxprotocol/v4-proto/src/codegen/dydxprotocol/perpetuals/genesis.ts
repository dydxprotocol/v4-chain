import { Perpetual, PerpetualSDKType, LiquidityTier, LiquidityTierSDKType } from "./perpetual";
import { Params, ParamsSDKType } from "./params";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** GenesisState defines the perpetuals module's genesis state. */

export interface GenesisState {
  perpetuals: Perpetual[];
  liquidityTiers: LiquidityTier[];
  params?: Params;
}
/** GenesisState defines the perpetuals module's genesis state. */

export interface GenesisStateSDKType {
  perpetuals: PerpetualSDKType[];
  liquidity_tiers: LiquidityTierSDKType[];
  params?: ParamsSDKType;
}

function createBaseGenesisState(): GenesisState {
  return {
    perpetuals: [],
    liquidityTiers: [],
    params: undefined
  };
}

export const GenesisState = {
  encode(message: GenesisState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.perpetuals) {
      Perpetual.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    for (const v of message.liquidityTiers) {
      LiquidityTier.encode(v!, writer.uint32(18).fork()).ldelim();
    }

    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(26).fork()).ldelim();
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
          message.perpetuals.push(Perpetual.decode(reader, reader.uint32()));
          break;

        case 2:
          message.liquidityTiers.push(LiquidityTier.decode(reader, reader.uint32()));
          break;

        case 3:
          message.params = Params.decode(reader, reader.uint32());
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
    message.perpetuals = object.perpetuals?.map(e => Perpetual.fromPartial(e)) || [];
    message.liquidityTiers = object.liquidityTiers?.map(e => LiquidityTier.fromPartial(e)) || [];
    message.params = object.params !== undefined && object.params !== null ? Params.fromPartial(object.params) : undefined;
    return message;
  }

};