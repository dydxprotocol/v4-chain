import { Timestamp } from "../../../google/protobuf/timestamp";
import { BinaryReader, BinaryWriter } from "../../../binary";
import { toTimestamp, fromTimestamp } from "../../../helpers";
/** UpdateMarketPriceRequest is a request message updating market prices. */
export interface UpdateMarketPricesRequest {
  marketPriceUpdates: MarketPriceUpdate[];
}
export interface UpdateMarketPricesRequestProtoMsg {
  typeUrl: "/dydxprotocol.daemons.pricefeed.UpdateMarketPricesRequest";
  value: Uint8Array;
}
/** UpdateMarketPriceRequest is a request message updating market prices. */
export interface UpdateMarketPricesRequestAmino {
  market_price_updates?: MarketPriceUpdateAmino[];
}
export interface UpdateMarketPricesRequestAminoMsg {
  type: "/dydxprotocol.daemons.pricefeed.UpdateMarketPricesRequest";
  value: UpdateMarketPricesRequestAmino;
}
/** UpdateMarketPriceRequest is a request message updating market prices. */
export interface UpdateMarketPricesRequestSDKType {
  market_price_updates: MarketPriceUpdateSDKType[];
}
/** UpdateMarketPricesResponse is a response message for updating market prices. */
export interface UpdateMarketPricesResponse {}
export interface UpdateMarketPricesResponseProtoMsg {
  typeUrl: "/dydxprotocol.daemons.pricefeed.UpdateMarketPricesResponse";
  value: Uint8Array;
}
/** UpdateMarketPricesResponse is a response message for updating market prices. */
export interface UpdateMarketPricesResponseAmino {}
export interface UpdateMarketPricesResponseAminoMsg {
  type: "/dydxprotocol.daemons.pricefeed.UpdateMarketPricesResponse";
  value: UpdateMarketPricesResponseAmino;
}
/** UpdateMarketPricesResponse is a response message for updating market prices. */
export interface UpdateMarketPricesResponseSDKType {}
/** ExchangePrice represents a specific exchange's market price */
export interface ExchangePrice {
  exchangeId: string;
  price: bigint;
  lastUpdateTime?: Date;
}
export interface ExchangePriceProtoMsg {
  typeUrl: "/dydxprotocol.daemons.pricefeed.ExchangePrice";
  value: Uint8Array;
}
/** ExchangePrice represents a specific exchange's market price */
export interface ExchangePriceAmino {
  exchange_id?: string;
  price?: string;
  last_update_time?: string;
}
export interface ExchangePriceAminoMsg {
  type: "/dydxprotocol.daemons.pricefeed.ExchangePrice";
  value: ExchangePriceAmino;
}
/** ExchangePrice represents a specific exchange's market price */
export interface ExchangePriceSDKType {
  exchange_id: string;
  price: bigint;
  last_update_time?: Date;
}
/** MarketPriceUpdate represents an update to a single market */
export interface MarketPriceUpdate {
  marketId: number;
  exchangePrices: ExchangePrice[];
}
export interface MarketPriceUpdateProtoMsg {
  typeUrl: "/dydxprotocol.daemons.pricefeed.MarketPriceUpdate";
  value: Uint8Array;
}
/** MarketPriceUpdate represents an update to a single market */
export interface MarketPriceUpdateAmino {
  market_id?: number;
  exchange_prices?: ExchangePriceAmino[];
}
export interface MarketPriceUpdateAminoMsg {
  type: "/dydxprotocol.daemons.pricefeed.MarketPriceUpdate";
  value: MarketPriceUpdateAmino;
}
/** MarketPriceUpdate represents an update to a single market */
export interface MarketPriceUpdateSDKType {
  market_id: number;
  exchange_prices: ExchangePriceSDKType[];
}
function createBaseUpdateMarketPricesRequest(): UpdateMarketPricesRequest {
  return {
    marketPriceUpdates: []
  };
}
export const UpdateMarketPricesRequest = {
  typeUrl: "/dydxprotocol.daemons.pricefeed.UpdateMarketPricesRequest",
  encode(message: UpdateMarketPricesRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.marketPriceUpdates) {
      MarketPriceUpdate.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): UpdateMarketPricesRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateMarketPricesRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.marketPriceUpdates.push(MarketPriceUpdate.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<UpdateMarketPricesRequest>): UpdateMarketPricesRequest {
    const message = createBaseUpdateMarketPricesRequest();
    message.marketPriceUpdates = object.marketPriceUpdates?.map(e => MarketPriceUpdate.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: UpdateMarketPricesRequestAmino): UpdateMarketPricesRequest {
    const message = createBaseUpdateMarketPricesRequest();
    message.marketPriceUpdates = object.market_price_updates?.map(e => MarketPriceUpdate.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: UpdateMarketPricesRequest): UpdateMarketPricesRequestAmino {
    const obj: any = {};
    if (message.marketPriceUpdates) {
      obj.market_price_updates = message.marketPriceUpdates.map(e => e ? MarketPriceUpdate.toAmino(e) : undefined);
    } else {
      obj.market_price_updates = [];
    }
    return obj;
  },
  fromAminoMsg(object: UpdateMarketPricesRequestAminoMsg): UpdateMarketPricesRequest {
    return UpdateMarketPricesRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: UpdateMarketPricesRequestProtoMsg): UpdateMarketPricesRequest {
    return UpdateMarketPricesRequest.decode(message.value);
  },
  toProto(message: UpdateMarketPricesRequest): Uint8Array {
    return UpdateMarketPricesRequest.encode(message).finish();
  },
  toProtoMsg(message: UpdateMarketPricesRequest): UpdateMarketPricesRequestProtoMsg {
    return {
      typeUrl: "/dydxprotocol.daemons.pricefeed.UpdateMarketPricesRequest",
      value: UpdateMarketPricesRequest.encode(message).finish()
    };
  }
};
function createBaseUpdateMarketPricesResponse(): UpdateMarketPricesResponse {
  return {};
}
export const UpdateMarketPricesResponse = {
  typeUrl: "/dydxprotocol.daemons.pricefeed.UpdateMarketPricesResponse",
  encode(_: UpdateMarketPricesResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): UpdateMarketPricesResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateMarketPricesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(_: Partial<UpdateMarketPricesResponse>): UpdateMarketPricesResponse {
    const message = createBaseUpdateMarketPricesResponse();
    return message;
  },
  fromAmino(_: UpdateMarketPricesResponseAmino): UpdateMarketPricesResponse {
    const message = createBaseUpdateMarketPricesResponse();
    return message;
  },
  toAmino(_: UpdateMarketPricesResponse): UpdateMarketPricesResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: UpdateMarketPricesResponseAminoMsg): UpdateMarketPricesResponse {
    return UpdateMarketPricesResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: UpdateMarketPricesResponseProtoMsg): UpdateMarketPricesResponse {
    return UpdateMarketPricesResponse.decode(message.value);
  },
  toProto(message: UpdateMarketPricesResponse): Uint8Array {
    return UpdateMarketPricesResponse.encode(message).finish();
  },
  toProtoMsg(message: UpdateMarketPricesResponse): UpdateMarketPricesResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.daemons.pricefeed.UpdateMarketPricesResponse",
      value: UpdateMarketPricesResponse.encode(message).finish()
    };
  }
};
function createBaseExchangePrice(): ExchangePrice {
  return {
    exchangeId: "",
    price: BigInt(0),
    lastUpdateTime: undefined
  };
}
export const ExchangePrice = {
  typeUrl: "/dydxprotocol.daemons.pricefeed.ExchangePrice",
  encode(message: ExchangePrice, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.exchangeId !== "") {
      writer.uint32(10).string(message.exchangeId);
    }
    if (message.price !== BigInt(0)) {
      writer.uint32(16).uint64(message.price);
    }
    if (message.lastUpdateTime !== undefined) {
      Timestamp.encode(toTimestamp(message.lastUpdateTime), writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): ExchangePrice {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExchangePrice();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.exchangeId = reader.string();
          break;
        case 2:
          message.price = reader.uint64();
          break;
        case 3:
          message.lastUpdateTime = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<ExchangePrice>): ExchangePrice {
    const message = createBaseExchangePrice();
    message.exchangeId = object.exchangeId ?? "";
    message.price = object.price !== undefined && object.price !== null ? BigInt(object.price.toString()) : BigInt(0);
    message.lastUpdateTime = object.lastUpdateTime ?? undefined;
    return message;
  },
  fromAmino(object: ExchangePriceAmino): ExchangePrice {
    const message = createBaseExchangePrice();
    if (object.exchange_id !== undefined && object.exchange_id !== null) {
      message.exchangeId = object.exchange_id;
    }
    if (object.price !== undefined && object.price !== null) {
      message.price = BigInt(object.price);
    }
    if (object.last_update_time !== undefined && object.last_update_time !== null) {
      message.lastUpdateTime = fromTimestamp(Timestamp.fromAmino(object.last_update_time));
    }
    return message;
  },
  toAmino(message: ExchangePrice): ExchangePriceAmino {
    const obj: any = {};
    obj.exchange_id = message.exchangeId;
    obj.price = message.price ? message.price.toString() : undefined;
    obj.last_update_time = message.lastUpdateTime ? Timestamp.toAmino(toTimestamp(message.lastUpdateTime)) : undefined;
    return obj;
  },
  fromAminoMsg(object: ExchangePriceAminoMsg): ExchangePrice {
    return ExchangePrice.fromAmino(object.value);
  },
  fromProtoMsg(message: ExchangePriceProtoMsg): ExchangePrice {
    return ExchangePrice.decode(message.value);
  },
  toProto(message: ExchangePrice): Uint8Array {
    return ExchangePrice.encode(message).finish();
  },
  toProtoMsg(message: ExchangePrice): ExchangePriceProtoMsg {
    return {
      typeUrl: "/dydxprotocol.daemons.pricefeed.ExchangePrice",
      value: ExchangePrice.encode(message).finish()
    };
  }
};
function createBaseMarketPriceUpdate(): MarketPriceUpdate {
  return {
    marketId: 0,
    exchangePrices: []
  };
}
export const MarketPriceUpdate = {
  typeUrl: "/dydxprotocol.daemons.pricefeed.MarketPriceUpdate",
  encode(message: MarketPriceUpdate, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.marketId !== 0) {
      writer.uint32(8).uint32(message.marketId);
    }
    for (const v of message.exchangePrices) {
      ExchangePrice.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MarketPriceUpdate {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketPriceUpdate();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.marketId = reader.uint32();
          break;
        case 2:
          message.exchangePrices.push(ExchangePrice.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MarketPriceUpdate>): MarketPriceUpdate {
    const message = createBaseMarketPriceUpdate();
    message.marketId = object.marketId ?? 0;
    message.exchangePrices = object.exchangePrices?.map(e => ExchangePrice.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: MarketPriceUpdateAmino): MarketPriceUpdate {
    const message = createBaseMarketPriceUpdate();
    if (object.market_id !== undefined && object.market_id !== null) {
      message.marketId = object.market_id;
    }
    message.exchangePrices = object.exchange_prices?.map(e => ExchangePrice.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: MarketPriceUpdate): MarketPriceUpdateAmino {
    const obj: any = {};
    obj.market_id = message.marketId;
    if (message.exchangePrices) {
      obj.exchange_prices = message.exchangePrices.map(e => e ? ExchangePrice.toAmino(e) : undefined);
    } else {
      obj.exchange_prices = [];
    }
    return obj;
  },
  fromAminoMsg(object: MarketPriceUpdateAminoMsg): MarketPriceUpdate {
    return MarketPriceUpdate.fromAmino(object.value);
  },
  fromProtoMsg(message: MarketPriceUpdateProtoMsg): MarketPriceUpdate {
    return MarketPriceUpdate.decode(message.value);
  },
  toProto(message: MarketPriceUpdate): Uint8Array {
    return MarketPriceUpdate.encode(message).finish();
  },
  toProtoMsg(message: MarketPriceUpdate): MarketPriceUpdateProtoMsg {
    return {
      typeUrl: "/dydxprotocol.daemons.pricefeed.MarketPriceUpdate",
      value: MarketPriceUpdate.encode(message).finish()
    };
  }
};