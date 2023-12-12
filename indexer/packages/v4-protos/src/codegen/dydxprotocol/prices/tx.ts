import { MarketParam, MarketParamAmino, MarketParamSDKType } from "./market_param";
import { BinaryReader, BinaryWriter } from "../../binary";
/**
 * MsgCreateOracleMarket is a message used by x/gov for creating a new oracle
 * market.
 */
export interface MsgCreateOracleMarket {
  /** The address that controls the module. */
  authority: string;
  /** `params` defines parameters for the new oracle market. */
  params: MarketParam;
}
export interface MsgCreateOracleMarketProtoMsg {
  typeUrl: "/dydxprotocol.prices.MsgCreateOracleMarket";
  value: Uint8Array;
}
/**
 * MsgCreateOracleMarket is a message used by x/gov for creating a new oracle
 * market.
 */
export interface MsgCreateOracleMarketAmino {
  /** The address that controls the module. */
  authority?: string;
  /** `params` defines parameters for the new oracle market. */
  params?: MarketParamAmino;
}
export interface MsgCreateOracleMarketAminoMsg {
  type: "/dydxprotocol.prices.MsgCreateOracleMarket";
  value: MsgCreateOracleMarketAmino;
}
/**
 * MsgCreateOracleMarket is a message used by x/gov for creating a new oracle
 * market.
 */
export interface MsgCreateOracleMarketSDKType {
  authority: string;
  params: MarketParamSDKType;
}
/** MsgCreateOracleMarketResponse defines the CreateOracleMarket response type. */
export interface MsgCreateOracleMarketResponse {}
export interface MsgCreateOracleMarketResponseProtoMsg {
  typeUrl: "/dydxprotocol.prices.MsgCreateOracleMarketResponse";
  value: Uint8Array;
}
/** MsgCreateOracleMarketResponse defines the CreateOracleMarket response type. */
export interface MsgCreateOracleMarketResponseAmino {}
export interface MsgCreateOracleMarketResponseAminoMsg {
  type: "/dydxprotocol.prices.MsgCreateOracleMarketResponse";
  value: MsgCreateOracleMarketResponseAmino;
}
/** MsgCreateOracleMarketResponse defines the CreateOracleMarket response type. */
export interface MsgCreateOracleMarketResponseSDKType {}
/** MsgUpdateMarketPrices is a request type for the UpdateMarketPrices method. */
export interface MsgUpdateMarketPrices {
  marketPriceUpdates: MsgUpdateMarketPrices_MarketPrice[];
}
export interface MsgUpdateMarketPricesProtoMsg {
  typeUrl: "/dydxprotocol.prices.MsgUpdateMarketPrices";
  value: Uint8Array;
}
/** MsgUpdateMarketPrices is a request type for the UpdateMarketPrices method. */
export interface MsgUpdateMarketPricesAmino {
  market_price_updates?: MsgUpdateMarketPrices_MarketPriceAmino[];
}
export interface MsgUpdateMarketPricesAminoMsg {
  type: "/dydxprotocol.prices.MsgUpdateMarketPrices";
  value: MsgUpdateMarketPricesAmino;
}
/** MsgUpdateMarketPrices is a request type for the UpdateMarketPrices method. */
export interface MsgUpdateMarketPricesSDKType {
  market_price_updates: MsgUpdateMarketPrices_MarketPriceSDKType[];
}
/** MarketPrice represents a price update for a single market */
export interface MsgUpdateMarketPrices_MarketPrice {
  /** The id of market to update */
  marketId: number;
  /** The updated price */
  price: bigint;
}
export interface MsgUpdateMarketPrices_MarketPriceProtoMsg {
  typeUrl: "/dydxprotocol.prices.MarketPrice";
  value: Uint8Array;
}
/** MarketPrice represents a price update for a single market */
export interface MsgUpdateMarketPrices_MarketPriceAmino {
  /** The id of market to update */
  market_id?: number;
  /** The updated price */
  price?: string;
}
export interface MsgUpdateMarketPrices_MarketPriceAminoMsg {
  type: "/dydxprotocol.prices.MarketPrice";
  value: MsgUpdateMarketPrices_MarketPriceAmino;
}
/** MarketPrice represents a price update for a single market */
export interface MsgUpdateMarketPrices_MarketPriceSDKType {
  market_id: number;
  price: bigint;
}
/**
 * MsgUpdateMarketPricesResponse defines the MsgUpdateMarketPrices response
 * type.
 */
export interface MsgUpdateMarketPricesResponse {}
export interface MsgUpdateMarketPricesResponseProtoMsg {
  typeUrl: "/dydxprotocol.prices.MsgUpdateMarketPricesResponse";
  value: Uint8Array;
}
/**
 * MsgUpdateMarketPricesResponse defines the MsgUpdateMarketPrices response
 * type.
 */
export interface MsgUpdateMarketPricesResponseAmino {}
export interface MsgUpdateMarketPricesResponseAminoMsg {
  type: "/dydxprotocol.prices.MsgUpdateMarketPricesResponse";
  value: MsgUpdateMarketPricesResponseAmino;
}
/**
 * MsgUpdateMarketPricesResponse defines the MsgUpdateMarketPrices response
 * type.
 */
export interface MsgUpdateMarketPricesResponseSDKType {}
/**
 * MsgUpdateMarketParam is a message used by x/gov for updating the parameters
 * of an oracle market.
 */
export interface MsgUpdateMarketParam {
  authority: string;
  /** The market param to update. Each field must be set. */
  marketParam: MarketParam;
}
export interface MsgUpdateMarketParamProtoMsg {
  typeUrl: "/dydxprotocol.prices.MsgUpdateMarketParam";
  value: Uint8Array;
}
/**
 * MsgUpdateMarketParam is a message used by x/gov for updating the parameters
 * of an oracle market.
 */
export interface MsgUpdateMarketParamAmino {
  authority?: string;
  /** The market param to update. Each field must be set. */
  market_param?: MarketParamAmino;
}
export interface MsgUpdateMarketParamAminoMsg {
  type: "/dydxprotocol.prices.MsgUpdateMarketParam";
  value: MsgUpdateMarketParamAmino;
}
/**
 * MsgUpdateMarketParam is a message used by x/gov for updating the parameters
 * of an oracle market.
 */
export interface MsgUpdateMarketParamSDKType {
  authority: string;
  market_param: MarketParamSDKType;
}
/** MsgUpdateMarketParamResponse defines the UpdateMarketParam response type. */
export interface MsgUpdateMarketParamResponse {}
export interface MsgUpdateMarketParamResponseProtoMsg {
  typeUrl: "/dydxprotocol.prices.MsgUpdateMarketParamResponse";
  value: Uint8Array;
}
/** MsgUpdateMarketParamResponse defines the UpdateMarketParam response type. */
export interface MsgUpdateMarketParamResponseAmino {}
export interface MsgUpdateMarketParamResponseAminoMsg {
  type: "/dydxprotocol.prices.MsgUpdateMarketParamResponse";
  value: MsgUpdateMarketParamResponseAmino;
}
/** MsgUpdateMarketParamResponse defines the UpdateMarketParam response type. */
export interface MsgUpdateMarketParamResponseSDKType {}
function createBaseMsgCreateOracleMarket(): MsgCreateOracleMarket {
  return {
    authority: "",
    params: MarketParam.fromPartial({})
  };
}
export const MsgCreateOracleMarket = {
  typeUrl: "/dydxprotocol.prices.MsgCreateOracleMarket",
  encode(message: MsgCreateOracleMarket, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }
    if (message.params !== undefined) {
      MarketParam.encode(message.params, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgCreateOracleMarket {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCreateOracleMarket();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;
        case 2:
          message.params = MarketParam.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgCreateOracleMarket>): MsgCreateOracleMarket {
    const message = createBaseMsgCreateOracleMarket();
    message.authority = object.authority ?? "";
    message.params = object.params !== undefined && object.params !== null ? MarketParam.fromPartial(object.params) : undefined;
    return message;
  },
  fromAmino(object: MsgCreateOracleMarketAmino): MsgCreateOracleMarket {
    const message = createBaseMsgCreateOracleMarket();
    if (object.authority !== undefined && object.authority !== null) {
      message.authority = object.authority;
    }
    if (object.params !== undefined && object.params !== null) {
      message.params = MarketParam.fromAmino(object.params);
    }
    return message;
  },
  toAmino(message: MsgCreateOracleMarket): MsgCreateOracleMarketAmino {
    const obj: any = {};
    obj.authority = message.authority;
    obj.params = message.params ? MarketParam.toAmino(message.params) : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgCreateOracleMarketAminoMsg): MsgCreateOracleMarket {
    return MsgCreateOracleMarket.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgCreateOracleMarketProtoMsg): MsgCreateOracleMarket {
    return MsgCreateOracleMarket.decode(message.value);
  },
  toProto(message: MsgCreateOracleMarket): Uint8Array {
    return MsgCreateOracleMarket.encode(message).finish();
  },
  toProtoMsg(message: MsgCreateOracleMarket): MsgCreateOracleMarketProtoMsg {
    return {
      typeUrl: "/dydxprotocol.prices.MsgCreateOracleMarket",
      value: MsgCreateOracleMarket.encode(message).finish()
    };
  }
};
function createBaseMsgCreateOracleMarketResponse(): MsgCreateOracleMarketResponse {
  return {};
}
export const MsgCreateOracleMarketResponse = {
  typeUrl: "/dydxprotocol.prices.MsgCreateOracleMarketResponse",
  encode(_: MsgCreateOracleMarketResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgCreateOracleMarketResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCreateOracleMarketResponse();
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
  fromPartial(_: Partial<MsgCreateOracleMarketResponse>): MsgCreateOracleMarketResponse {
    const message = createBaseMsgCreateOracleMarketResponse();
    return message;
  },
  fromAmino(_: MsgCreateOracleMarketResponseAmino): MsgCreateOracleMarketResponse {
    const message = createBaseMsgCreateOracleMarketResponse();
    return message;
  },
  toAmino(_: MsgCreateOracleMarketResponse): MsgCreateOracleMarketResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgCreateOracleMarketResponseAminoMsg): MsgCreateOracleMarketResponse {
    return MsgCreateOracleMarketResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgCreateOracleMarketResponseProtoMsg): MsgCreateOracleMarketResponse {
    return MsgCreateOracleMarketResponse.decode(message.value);
  },
  toProto(message: MsgCreateOracleMarketResponse): Uint8Array {
    return MsgCreateOracleMarketResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgCreateOracleMarketResponse): MsgCreateOracleMarketResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.prices.MsgCreateOracleMarketResponse",
      value: MsgCreateOracleMarketResponse.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateMarketPrices(): MsgUpdateMarketPrices {
  return {
    marketPriceUpdates: []
  };
}
export const MsgUpdateMarketPrices = {
  typeUrl: "/dydxprotocol.prices.MsgUpdateMarketPrices",
  encode(message: MsgUpdateMarketPrices, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.marketPriceUpdates) {
      MsgUpdateMarketPrices_MarketPrice.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateMarketPrices {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateMarketPrices();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.marketPriceUpdates.push(MsgUpdateMarketPrices_MarketPrice.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgUpdateMarketPrices>): MsgUpdateMarketPrices {
    const message = createBaseMsgUpdateMarketPrices();
    message.marketPriceUpdates = object.marketPriceUpdates?.map(e => MsgUpdateMarketPrices_MarketPrice.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: MsgUpdateMarketPricesAmino): MsgUpdateMarketPrices {
    const message = createBaseMsgUpdateMarketPrices();
    message.marketPriceUpdates = object.market_price_updates?.map(e => MsgUpdateMarketPrices_MarketPrice.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: MsgUpdateMarketPrices): MsgUpdateMarketPricesAmino {
    const obj: any = {};
    if (message.marketPriceUpdates) {
      obj.market_price_updates = message.marketPriceUpdates.map(e => e ? MsgUpdateMarketPrices_MarketPrice.toAmino(e) : undefined);
    } else {
      obj.market_price_updates = [];
    }
    return obj;
  },
  fromAminoMsg(object: MsgUpdateMarketPricesAminoMsg): MsgUpdateMarketPrices {
    return MsgUpdateMarketPrices.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateMarketPricesProtoMsg): MsgUpdateMarketPrices {
    return MsgUpdateMarketPrices.decode(message.value);
  },
  toProto(message: MsgUpdateMarketPrices): Uint8Array {
    return MsgUpdateMarketPrices.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateMarketPrices): MsgUpdateMarketPricesProtoMsg {
    return {
      typeUrl: "/dydxprotocol.prices.MsgUpdateMarketPrices",
      value: MsgUpdateMarketPrices.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateMarketPrices_MarketPrice(): MsgUpdateMarketPrices_MarketPrice {
  return {
    marketId: 0,
    price: BigInt(0)
  };
}
export const MsgUpdateMarketPrices_MarketPrice = {
  typeUrl: "/dydxprotocol.prices.MarketPrice",
  encode(message: MsgUpdateMarketPrices_MarketPrice, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.marketId !== 0) {
      writer.uint32(8).uint32(message.marketId);
    }
    if (message.price !== BigInt(0)) {
      writer.uint32(16).uint64(message.price);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateMarketPrices_MarketPrice {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateMarketPrices_MarketPrice();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.marketId = reader.uint32();
          break;
        case 2:
          message.price = reader.uint64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgUpdateMarketPrices_MarketPrice>): MsgUpdateMarketPrices_MarketPrice {
    const message = createBaseMsgUpdateMarketPrices_MarketPrice();
    message.marketId = object.marketId ?? 0;
    message.price = object.price !== undefined && object.price !== null ? BigInt(object.price.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: MsgUpdateMarketPrices_MarketPriceAmino): MsgUpdateMarketPrices_MarketPrice {
    const message = createBaseMsgUpdateMarketPrices_MarketPrice();
    if (object.market_id !== undefined && object.market_id !== null) {
      message.marketId = object.market_id;
    }
    if (object.price !== undefined && object.price !== null) {
      message.price = BigInt(object.price);
    }
    return message;
  },
  toAmino(message: MsgUpdateMarketPrices_MarketPrice): MsgUpdateMarketPrices_MarketPriceAmino {
    const obj: any = {};
    obj.market_id = message.marketId;
    obj.price = message.price ? message.price.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgUpdateMarketPrices_MarketPriceAminoMsg): MsgUpdateMarketPrices_MarketPrice {
    return MsgUpdateMarketPrices_MarketPrice.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateMarketPrices_MarketPriceProtoMsg): MsgUpdateMarketPrices_MarketPrice {
    return MsgUpdateMarketPrices_MarketPrice.decode(message.value);
  },
  toProto(message: MsgUpdateMarketPrices_MarketPrice): Uint8Array {
    return MsgUpdateMarketPrices_MarketPrice.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateMarketPrices_MarketPrice): MsgUpdateMarketPrices_MarketPriceProtoMsg {
    return {
      typeUrl: "/dydxprotocol.prices.MarketPrice",
      value: MsgUpdateMarketPrices_MarketPrice.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateMarketPricesResponse(): MsgUpdateMarketPricesResponse {
  return {};
}
export const MsgUpdateMarketPricesResponse = {
  typeUrl: "/dydxprotocol.prices.MsgUpdateMarketPricesResponse",
  encode(_: MsgUpdateMarketPricesResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateMarketPricesResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateMarketPricesResponse();
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
  fromPartial(_: Partial<MsgUpdateMarketPricesResponse>): MsgUpdateMarketPricesResponse {
    const message = createBaseMsgUpdateMarketPricesResponse();
    return message;
  },
  fromAmino(_: MsgUpdateMarketPricesResponseAmino): MsgUpdateMarketPricesResponse {
    const message = createBaseMsgUpdateMarketPricesResponse();
    return message;
  },
  toAmino(_: MsgUpdateMarketPricesResponse): MsgUpdateMarketPricesResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgUpdateMarketPricesResponseAminoMsg): MsgUpdateMarketPricesResponse {
    return MsgUpdateMarketPricesResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateMarketPricesResponseProtoMsg): MsgUpdateMarketPricesResponse {
    return MsgUpdateMarketPricesResponse.decode(message.value);
  },
  toProto(message: MsgUpdateMarketPricesResponse): Uint8Array {
    return MsgUpdateMarketPricesResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateMarketPricesResponse): MsgUpdateMarketPricesResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.prices.MsgUpdateMarketPricesResponse",
      value: MsgUpdateMarketPricesResponse.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateMarketParam(): MsgUpdateMarketParam {
  return {
    authority: "",
    marketParam: MarketParam.fromPartial({})
  };
}
export const MsgUpdateMarketParam = {
  typeUrl: "/dydxprotocol.prices.MsgUpdateMarketParam",
  encode(message: MsgUpdateMarketParam, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }
    if (message.marketParam !== undefined) {
      MarketParam.encode(message.marketParam, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateMarketParam {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateMarketParam();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;
        case 2:
          message.marketParam = MarketParam.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgUpdateMarketParam>): MsgUpdateMarketParam {
    const message = createBaseMsgUpdateMarketParam();
    message.authority = object.authority ?? "";
    message.marketParam = object.marketParam !== undefined && object.marketParam !== null ? MarketParam.fromPartial(object.marketParam) : undefined;
    return message;
  },
  fromAmino(object: MsgUpdateMarketParamAmino): MsgUpdateMarketParam {
    const message = createBaseMsgUpdateMarketParam();
    if (object.authority !== undefined && object.authority !== null) {
      message.authority = object.authority;
    }
    if (object.market_param !== undefined && object.market_param !== null) {
      message.marketParam = MarketParam.fromAmino(object.market_param);
    }
    return message;
  },
  toAmino(message: MsgUpdateMarketParam): MsgUpdateMarketParamAmino {
    const obj: any = {};
    obj.authority = message.authority;
    obj.market_param = message.marketParam ? MarketParam.toAmino(message.marketParam) : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgUpdateMarketParamAminoMsg): MsgUpdateMarketParam {
    return MsgUpdateMarketParam.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateMarketParamProtoMsg): MsgUpdateMarketParam {
    return MsgUpdateMarketParam.decode(message.value);
  },
  toProto(message: MsgUpdateMarketParam): Uint8Array {
    return MsgUpdateMarketParam.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateMarketParam): MsgUpdateMarketParamProtoMsg {
    return {
      typeUrl: "/dydxprotocol.prices.MsgUpdateMarketParam",
      value: MsgUpdateMarketParam.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateMarketParamResponse(): MsgUpdateMarketParamResponse {
  return {};
}
export const MsgUpdateMarketParamResponse = {
  typeUrl: "/dydxprotocol.prices.MsgUpdateMarketParamResponse",
  encode(_: MsgUpdateMarketParamResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateMarketParamResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateMarketParamResponse();
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
  fromPartial(_: Partial<MsgUpdateMarketParamResponse>): MsgUpdateMarketParamResponse {
    const message = createBaseMsgUpdateMarketParamResponse();
    return message;
  },
  fromAmino(_: MsgUpdateMarketParamResponseAmino): MsgUpdateMarketParamResponse {
    const message = createBaseMsgUpdateMarketParamResponse();
    return message;
  },
  toAmino(_: MsgUpdateMarketParamResponse): MsgUpdateMarketParamResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgUpdateMarketParamResponseAminoMsg): MsgUpdateMarketParamResponse {
    return MsgUpdateMarketParamResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateMarketParamResponseProtoMsg): MsgUpdateMarketParamResponse {
    return MsgUpdateMarketParamResponse.decode(message.value);
  },
  toProto(message: MsgUpdateMarketParamResponse): Uint8Array {
    return MsgUpdateMarketParamResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateMarketParamResponse): MsgUpdateMarketParamResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.prices.MsgUpdateMarketParamResponse",
      value: MsgUpdateMarketParamResponse.encode(message).finish()
    };
  }
};