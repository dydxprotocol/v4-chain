import * as _m0 from "protobufjs/minimal";
import { Long, DeepPartial } from "../../helpers";
/**
 * MarketMapperRevShareDetails specifies any details associated with the market
 * mapper revenue share
 */

export interface MarketMapperRevShareDetails {
  /** Unix timestamp recorded when the market revenue share expires */
  expirationTs: Long;
}
/**
 * MarketMapperRevShareDetails specifies any details associated with the market
 * mapper revenue share
 */

export interface MarketMapperRevShareDetailsSDKType {
  /** Unix timestamp recorded when the market revenue share expires */
  expiration_ts: Long;
}
/**
 * UnconditionalRevShareConfig stores recipients that
 * receive a share of net revenue unconditionally.
 */

export interface UnconditionalRevShareConfig {
  /** Configs for each recipient. */
  configs: UnconditionalRevShareConfig_RecipientConfig[];
}
/**
 * UnconditionalRevShareConfig stores recipients that
 * receive a share of net revenue unconditionally.
 */

export interface UnconditionalRevShareConfigSDKType {
  /** Configs for each recipient. */
  configs: UnconditionalRevShareConfig_RecipientConfigSDKType[];
}
/** Describes the config of a recipient */

export interface UnconditionalRevShareConfig_RecipientConfig {
  /** Address of the recepient. */
  address: string;
  /** Percentage of net revenue to share with recipient, in parts-per-million. */

  sharePpm: number;
}
/** Describes the config of a recipient */

export interface UnconditionalRevShareConfig_RecipientConfigSDKType {
  /** Address of the recepient. */
  address: string;
  /** Percentage of net revenue to share with recipient, in parts-per-million. */

  share_ppm: number;
}
/** Message to set the order router revenue share */

export interface OrderRouterRevShare {
  /** The address of the order router. */
  address: string;
  /** The share of the revenue to be paid to the order router. */

  sharePpm: number;
}
/** Message to set the order router revenue share */

export interface OrderRouterRevShareSDKType {
  /** The address of the order router. */
  address: string;
  /** The share of the revenue to be paid to the order router. */

  share_ppm: number;
}

function createBaseMarketMapperRevShareDetails(): MarketMapperRevShareDetails {
  return {
    expirationTs: Long.UZERO
  };
}

export const MarketMapperRevShareDetails = {
  encode(message: MarketMapperRevShareDetails, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (!message.expirationTs.isZero()) {
      writer.uint32(8).uint64(message.expirationTs);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketMapperRevShareDetails {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketMapperRevShareDetails();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.expirationTs = (reader.uint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MarketMapperRevShareDetails>): MarketMapperRevShareDetails {
    const message = createBaseMarketMapperRevShareDetails();
    message.expirationTs = object.expirationTs !== undefined && object.expirationTs !== null ? Long.fromValue(object.expirationTs) : Long.UZERO;
    return message;
  }

};

function createBaseUnconditionalRevShareConfig(): UnconditionalRevShareConfig {
  return {
    configs: []
  };
}

export const UnconditionalRevShareConfig = {
  encode(message: UnconditionalRevShareConfig, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.configs) {
      UnconditionalRevShareConfig_RecipientConfig.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UnconditionalRevShareConfig {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUnconditionalRevShareConfig();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.configs.push(UnconditionalRevShareConfig_RecipientConfig.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<UnconditionalRevShareConfig>): UnconditionalRevShareConfig {
    const message = createBaseUnconditionalRevShareConfig();
    message.configs = object.configs?.map(e => UnconditionalRevShareConfig_RecipientConfig.fromPartial(e)) || [];
    return message;
  }

};

function createBaseUnconditionalRevShareConfig_RecipientConfig(): UnconditionalRevShareConfig_RecipientConfig {
  return {
    address: "",
    sharePpm: 0
  };
}

export const UnconditionalRevShareConfig_RecipientConfig = {
  encode(message: UnconditionalRevShareConfig_RecipientConfig, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }

    if (message.sharePpm !== 0) {
      writer.uint32(16).uint32(message.sharePpm);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UnconditionalRevShareConfig_RecipientConfig {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUnconditionalRevShareConfig_RecipientConfig();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;

        case 2:
          message.sharePpm = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<UnconditionalRevShareConfig_RecipientConfig>): UnconditionalRevShareConfig_RecipientConfig {
    const message = createBaseUnconditionalRevShareConfig_RecipientConfig();
    message.address = object.address ?? "";
    message.sharePpm = object.sharePpm ?? 0;
    return message;
  }

};

function createBaseOrderRouterRevShare(): OrderRouterRevShare {
  return {
    address: "",
    sharePpm: 0
  };
}

export const OrderRouterRevShare = {
  encode(message: OrderRouterRevShare, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }

    if (message.sharePpm !== 0) {
      writer.uint32(16).uint32(message.sharePpm);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): OrderRouterRevShare {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOrderRouterRevShare();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;

        case 2:
          message.sharePpm = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<OrderRouterRevShare>): OrderRouterRevShare {
    const message = createBaseOrderRouterRevShare();
    message.address = object.address ?? "";
    message.sharePpm = object.sharePpm ?? 0;
    return message;
  }

};