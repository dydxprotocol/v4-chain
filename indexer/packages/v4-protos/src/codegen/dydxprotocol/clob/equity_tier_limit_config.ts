import { BinaryReader, BinaryWriter } from "../../binary";
import { bytesFromBase64, base64FromBytes } from "../../helpers";
/**
 * Defines the set of equity tiers to limit how many open orders
 * a subaccount is allowed to have.
 */
export interface EquityTierLimitConfiguration {
  /**
   * How many short term stateful orders are allowed per equity tier.
   * Specifying 0 values disables this limit.
   */
  shortTermOrderEquityTiers: EquityTierLimit[];
  /**
   * How many open stateful orders are allowed per equity tier.
   * Specifying 0 values disables this limit.
   */
  statefulOrderEquityTiers: EquityTierLimit[];
}
export interface EquityTierLimitConfigurationProtoMsg {
  typeUrl: "/dydxprotocol.clob.EquityTierLimitConfiguration";
  value: Uint8Array;
}
/**
 * Defines the set of equity tiers to limit how many open orders
 * a subaccount is allowed to have.
 */
export interface EquityTierLimitConfigurationAmino {
  /**
   * How many short term stateful orders are allowed per equity tier.
   * Specifying 0 values disables this limit.
   */
  short_term_order_equity_tiers?: EquityTierLimitAmino[];
  /**
   * How many open stateful orders are allowed per equity tier.
   * Specifying 0 values disables this limit.
   */
  stateful_order_equity_tiers?: EquityTierLimitAmino[];
}
export interface EquityTierLimitConfigurationAminoMsg {
  type: "/dydxprotocol.clob.EquityTierLimitConfiguration";
  value: EquityTierLimitConfigurationAmino;
}
/**
 * Defines the set of equity tiers to limit how many open orders
 * a subaccount is allowed to have.
 */
export interface EquityTierLimitConfigurationSDKType {
  short_term_order_equity_tiers: EquityTierLimitSDKType[];
  stateful_order_equity_tiers: EquityTierLimitSDKType[];
}
/** Defines an equity tier limit. */
export interface EquityTierLimit {
  /** The total net collateral in USDC quote quantums of equity required. */
  usdTncRequired: Uint8Array;
  /** What the limit is for `usd_tnc_required`. */
  limit: number;
}
export interface EquityTierLimitProtoMsg {
  typeUrl: "/dydxprotocol.clob.EquityTierLimit";
  value: Uint8Array;
}
/** Defines an equity tier limit. */
export interface EquityTierLimitAmino {
  /** The total net collateral in USDC quote quantums of equity required. */
  usd_tnc_required?: string;
  /** What the limit is for `usd_tnc_required`. */
  limit?: number;
}
export interface EquityTierLimitAminoMsg {
  type: "/dydxprotocol.clob.EquityTierLimit";
  value: EquityTierLimitAmino;
}
/** Defines an equity tier limit. */
export interface EquityTierLimitSDKType {
  usd_tnc_required: Uint8Array;
  limit: number;
}
function createBaseEquityTierLimitConfiguration(): EquityTierLimitConfiguration {
  return {
    shortTermOrderEquityTiers: [],
    statefulOrderEquityTiers: []
  };
}
export const EquityTierLimitConfiguration = {
  typeUrl: "/dydxprotocol.clob.EquityTierLimitConfiguration",
  encode(message: EquityTierLimitConfiguration, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.shortTermOrderEquityTiers) {
      EquityTierLimit.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.statefulOrderEquityTiers) {
      EquityTierLimit.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): EquityTierLimitConfiguration {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEquityTierLimitConfiguration();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.shortTermOrderEquityTiers.push(EquityTierLimit.decode(reader, reader.uint32()));
          break;
        case 2:
          message.statefulOrderEquityTiers.push(EquityTierLimit.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<EquityTierLimitConfiguration>): EquityTierLimitConfiguration {
    const message = createBaseEquityTierLimitConfiguration();
    message.shortTermOrderEquityTiers = object.shortTermOrderEquityTiers?.map(e => EquityTierLimit.fromPartial(e)) || [];
    message.statefulOrderEquityTiers = object.statefulOrderEquityTiers?.map(e => EquityTierLimit.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: EquityTierLimitConfigurationAmino): EquityTierLimitConfiguration {
    const message = createBaseEquityTierLimitConfiguration();
    message.shortTermOrderEquityTiers = object.short_term_order_equity_tiers?.map(e => EquityTierLimit.fromAmino(e)) || [];
    message.statefulOrderEquityTiers = object.stateful_order_equity_tiers?.map(e => EquityTierLimit.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: EquityTierLimitConfiguration): EquityTierLimitConfigurationAmino {
    const obj: any = {};
    if (message.shortTermOrderEquityTiers) {
      obj.short_term_order_equity_tiers = message.shortTermOrderEquityTiers.map(e => e ? EquityTierLimit.toAmino(e) : undefined);
    } else {
      obj.short_term_order_equity_tiers = [];
    }
    if (message.statefulOrderEquityTiers) {
      obj.stateful_order_equity_tiers = message.statefulOrderEquityTiers.map(e => e ? EquityTierLimit.toAmino(e) : undefined);
    } else {
      obj.stateful_order_equity_tiers = [];
    }
    return obj;
  },
  fromAminoMsg(object: EquityTierLimitConfigurationAminoMsg): EquityTierLimitConfiguration {
    return EquityTierLimitConfiguration.fromAmino(object.value);
  },
  fromProtoMsg(message: EquityTierLimitConfigurationProtoMsg): EquityTierLimitConfiguration {
    return EquityTierLimitConfiguration.decode(message.value);
  },
  toProto(message: EquityTierLimitConfiguration): Uint8Array {
    return EquityTierLimitConfiguration.encode(message).finish();
  },
  toProtoMsg(message: EquityTierLimitConfiguration): EquityTierLimitConfigurationProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.EquityTierLimitConfiguration",
      value: EquityTierLimitConfiguration.encode(message).finish()
    };
  }
};
function createBaseEquityTierLimit(): EquityTierLimit {
  return {
    usdTncRequired: new Uint8Array(),
    limit: 0
  };
}
export const EquityTierLimit = {
  typeUrl: "/dydxprotocol.clob.EquityTierLimit",
  encode(message: EquityTierLimit, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.usdTncRequired.length !== 0) {
      writer.uint32(10).bytes(message.usdTncRequired);
    }
    if (message.limit !== 0) {
      writer.uint32(16).uint32(message.limit);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): EquityTierLimit {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEquityTierLimit();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.usdTncRequired = reader.bytes();
          break;
        case 2:
          message.limit = reader.uint32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<EquityTierLimit>): EquityTierLimit {
    const message = createBaseEquityTierLimit();
    message.usdTncRequired = object.usdTncRequired ?? new Uint8Array();
    message.limit = object.limit ?? 0;
    return message;
  },
  fromAmino(object: EquityTierLimitAmino): EquityTierLimit {
    const message = createBaseEquityTierLimit();
    if (object.usd_tnc_required !== undefined && object.usd_tnc_required !== null) {
      message.usdTncRequired = bytesFromBase64(object.usd_tnc_required);
    }
    if (object.limit !== undefined && object.limit !== null) {
      message.limit = object.limit;
    }
    return message;
  },
  toAmino(message: EquityTierLimit): EquityTierLimitAmino {
    const obj: any = {};
    obj.usd_tnc_required = message.usdTncRequired ? base64FromBytes(message.usdTncRequired) : undefined;
    obj.limit = message.limit;
    return obj;
  },
  fromAminoMsg(object: EquityTierLimitAminoMsg): EquityTierLimit {
    return EquityTierLimit.fromAmino(object.value);
  },
  fromProtoMsg(message: EquityTierLimitProtoMsg): EquityTierLimit {
    return EquityTierLimit.decode(message.value);
  },
  toProto(message: EquityTierLimit): Uint8Array {
    return EquityTierLimit.encode(message).finish();
  },
  toProtoMsg(message: EquityTierLimit): EquityTierLimitProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.EquityTierLimit",
      value: EquityTierLimit.encode(message).finish()
    };
  }
};