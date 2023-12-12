import { MarketParam, MarketParamSDKType } from "./market_param";
import { MarketPrice, MarketPriceSDKType } from "./market_price";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** GenesisState defines the prices module's genesis state. */

export interface GenesisState {
  marketParams: MarketParam[];
  marketPrices: MarketPrice[];
}
/** GenesisState defines the prices module's genesis state. */

export interface GenesisStateSDKType {
  market_params: MarketParamSDKType[];
  market_prices: MarketPriceSDKType[];
}

function createBaseGenesisState(): GenesisState {
  return {
    marketParams: [],
    marketPrices: []
  };
}

export const GenesisState = {
  encode(message: GenesisState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.marketParams) {
      MarketParam.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    for (const v of message.marketPrices) {
      MarketPrice.encode(v!, writer.uint32(18).fork()).ldelim();
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
          message.marketParams.push(MarketParam.decode(reader, reader.uint32()));
          break;

        case 2:
          message.marketPrices.push(MarketPrice.decode(reader, reader.uint32()));
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
    message.marketParams = object.marketParams?.map(e => MarketParam.fromPartial(e)) || [];
    message.marketPrices = object.marketPrices?.map(e => MarketPrice.fromPartial(e)) || [];
    return message;
  }

};