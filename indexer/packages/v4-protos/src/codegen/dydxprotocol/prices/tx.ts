import { MarketParam, MarketParamSDKType } from "./market_param";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/**
 * MsgCreateOracleMarket is a message used by x/gov for creating a new oracle
 * market.
 */

export interface MsgCreateOracleMarket {
  /** The address that controls the module. */
  authority: string;
  /** `params` defines parameters for the new oracle market. */

  params?: MarketParam;
}
/**
 * MsgCreateOracleMarket is a message used by x/gov for creating a new oracle
 * market.
 */

export interface MsgCreateOracleMarketSDKType {
  /** The address that controls the module. */
  authority: string;
  /** `params` defines parameters for the new oracle market. */

  params?: MarketParamSDKType;
}
/** MsgCreateOracleMarketResponse defines the CreateOracleMarket response type. */

export interface MsgCreateOracleMarketResponse {}
/** MsgCreateOracleMarketResponse defines the CreateOracleMarket response type. */

export interface MsgCreateOracleMarketResponseSDKType {}
/**
 * MsgUpdateMarketParam is a message used by x/gov for updating the parameters
 * of an oracle market.
 */

export interface MsgUpdateMarketParam {
  authority: string;
  /** The market param to update. Each field must be set. */

  marketParam?: MarketParam;
}
/**
 * MsgUpdateMarketParam is a message used by x/gov for updating the parameters
 * of an oracle market.
 */

export interface MsgUpdateMarketParamSDKType {
  authority: string;
  /** The market param to update. Each field must be set. */

  market_param?: MarketParamSDKType;
}
/** MsgUpdateMarketParamResponse defines the UpdateMarketParam response type. */

export interface MsgUpdateMarketParamResponse {}
/** MsgUpdateMarketParamResponse defines the UpdateMarketParam response type. */

export interface MsgUpdateMarketParamResponseSDKType {}

function createBaseMsgCreateOracleMarket(): MsgCreateOracleMarket {
  return {
    authority: "",
    params: undefined
  };
}

export const MsgCreateOracleMarket = {
  encode(message: MsgCreateOracleMarket, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.params !== undefined) {
      MarketParam.encode(message.params, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgCreateOracleMarket {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
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

  fromPartial(object: DeepPartial<MsgCreateOracleMarket>): MsgCreateOracleMarket {
    const message = createBaseMsgCreateOracleMarket();
    message.authority = object.authority ?? "";
    message.params = object.params !== undefined && object.params !== null ? MarketParam.fromPartial(object.params) : undefined;
    return message;
  }

};

function createBaseMsgCreateOracleMarketResponse(): MsgCreateOracleMarketResponse {
  return {};
}

export const MsgCreateOracleMarketResponse = {
  encode(_: MsgCreateOracleMarketResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgCreateOracleMarketResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
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

  fromPartial(_: DeepPartial<MsgCreateOracleMarketResponse>): MsgCreateOracleMarketResponse {
    const message = createBaseMsgCreateOracleMarketResponse();
    return message;
  }

};

function createBaseMsgUpdateMarketParam(): MsgUpdateMarketParam {
  return {
    authority: "",
    marketParam: undefined
  };
}

export const MsgUpdateMarketParam = {
  encode(message: MsgUpdateMarketParam, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.marketParam !== undefined) {
      MarketParam.encode(message.marketParam, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateMarketParam {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
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

  fromPartial(object: DeepPartial<MsgUpdateMarketParam>): MsgUpdateMarketParam {
    const message = createBaseMsgUpdateMarketParam();
    message.authority = object.authority ?? "";
    message.marketParam = object.marketParam !== undefined && object.marketParam !== null ? MarketParam.fromPartial(object.marketParam) : undefined;
    return message;
  }

};

function createBaseMsgUpdateMarketParamResponse(): MsgUpdateMarketParamResponse {
  return {};
}

export const MsgUpdateMarketParamResponse = {
  encode(_: MsgUpdateMarketParamResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateMarketParamResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
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

  fromPartial(_: DeepPartial<MsgUpdateMarketParamResponse>): MsgUpdateMarketParamResponse {
    const message = createBaseMsgUpdateMarketParamResponse();
    return message;
  }

};