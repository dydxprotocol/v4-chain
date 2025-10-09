import { PerpetualFeeParams, PerpetualFeeParamsSDKType } from "./params";
import { PerMarketFeeDiscountParams, PerMarketFeeDiscountParamsSDKType } from "./per_market_fee_discount";
import { StakingTier, StakingTierSDKType } from "./staking_tier";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** MsgUpdatePerpetualFeeParams is the Msg/UpdatePerpetualFeeParams request type. */

export interface MsgUpdatePerpetualFeeParams {
  authority: string;
  /** Defines the parameters to update. All parameters must be supplied. */

  params?: PerpetualFeeParams;
}
/** MsgUpdatePerpetualFeeParams is the Msg/UpdatePerpetualFeeParams request type. */

export interface MsgUpdatePerpetualFeeParamsSDKType {
  authority: string;
  /** Defines the parameters to update. All parameters must be supplied. */

  params?: PerpetualFeeParamsSDKType;
}
/**
 * MsgUpdatePerpetualFeeParamsResponse is the Msg/UpdatePerpetualFeeParams
 * response type.
 */

export interface MsgUpdatePerpetualFeeParamsResponse {}
/**
 * MsgUpdatePerpetualFeeParamsResponse is the Msg/UpdatePerpetualFeeParams
 * response type.
 */

export interface MsgUpdatePerpetualFeeParamsResponseSDKType {}
/**
 * MsgSetMarketFeeDiscountParams is the Msg/SetMarketFeeDiscountParams
 * request type.
 */

export interface MsgSetMarketFeeDiscountParams {
  /** authority is the address that controls the module */
  authority: string;
  /** The per-market fee discount parameters to create or update */

  params: PerMarketFeeDiscountParams[];
}
/**
 * MsgSetMarketFeeDiscountParams is the Msg/SetMarketFeeDiscountParams
 * request type.
 */

export interface MsgSetMarketFeeDiscountParamsSDKType {
  /** authority is the address that controls the module */
  authority: string;
  /** The per-market fee discount parameters to create or update */

  params: PerMarketFeeDiscountParamsSDKType[];
}
/**
 * MsgSetMarketFeeDiscountParamsResponse is the
 * Msg/SetMarketFeeDiscountParams response type.
 */

export interface MsgSetMarketFeeDiscountParamsResponse {}
/**
 * MsgSetMarketFeeDiscountParamsResponse is the
 * Msg/SetMarketFeeDiscountParams response type.
 */

export interface MsgSetMarketFeeDiscountParamsResponseSDKType {}
/** MsgSetStakingTiers is the Msg/SetStakingTiers request type. */

export interface MsgSetStakingTiers {
  authority: string;
  /** List of Staking tiers */

  stakingTiers: StakingTier[];
}
/** MsgSetStakingTiers is the Msg/SetStakingTiers request type. */

export interface MsgSetStakingTiersSDKType {
  authority: string;
  /** List of Staking tiers */

  staking_tiers: StakingTierSDKType[];
}
/** MsgSetStakingTiersResponse is the Msg/SetStakingTiers response type. */

export interface MsgSetStakingTiersResponse {}
/** MsgSetStakingTiersResponse is the Msg/SetStakingTiers response type. */

export interface MsgSetStakingTiersResponseSDKType {}

function createBaseMsgUpdatePerpetualFeeParams(): MsgUpdatePerpetualFeeParams {
  return {
    authority: "",
    params: undefined
  };
}

export const MsgUpdatePerpetualFeeParams = {
  encode(message: MsgUpdatePerpetualFeeParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.params !== undefined) {
      PerpetualFeeParams.encode(message.params, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdatePerpetualFeeParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdatePerpetualFeeParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;

        case 2:
          message.params = PerpetualFeeParams.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgUpdatePerpetualFeeParams>): MsgUpdatePerpetualFeeParams {
    const message = createBaseMsgUpdatePerpetualFeeParams();
    message.authority = object.authority ?? "";
    message.params = object.params !== undefined && object.params !== null ? PerpetualFeeParams.fromPartial(object.params) : undefined;
    return message;
  }

};

function createBaseMsgUpdatePerpetualFeeParamsResponse(): MsgUpdatePerpetualFeeParamsResponse {
  return {};
}

export const MsgUpdatePerpetualFeeParamsResponse = {
  encode(_: MsgUpdatePerpetualFeeParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdatePerpetualFeeParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdatePerpetualFeeParamsResponse();

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

  fromPartial(_: DeepPartial<MsgUpdatePerpetualFeeParamsResponse>): MsgUpdatePerpetualFeeParamsResponse {
    const message = createBaseMsgUpdatePerpetualFeeParamsResponse();
    return message;
  }

};

function createBaseMsgSetMarketFeeDiscountParams(): MsgSetMarketFeeDiscountParams {
  return {
    authority: "",
    params: []
  };
}

export const MsgSetMarketFeeDiscountParams = {
  encode(message: MsgSetMarketFeeDiscountParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    for (const v of message.params) {
      PerMarketFeeDiscountParams.encode(v!, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetMarketFeeDiscountParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetMarketFeeDiscountParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;

        case 2:
          message.params.push(PerMarketFeeDiscountParams.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgSetMarketFeeDiscountParams>): MsgSetMarketFeeDiscountParams {
    const message = createBaseMsgSetMarketFeeDiscountParams();
    message.authority = object.authority ?? "";
    message.params = object.params?.map(e => PerMarketFeeDiscountParams.fromPartial(e)) || [];
    return message;
  }

};

function createBaseMsgSetMarketFeeDiscountParamsResponse(): MsgSetMarketFeeDiscountParamsResponse {
  return {};
}

export const MsgSetMarketFeeDiscountParamsResponse = {
  encode(_: MsgSetMarketFeeDiscountParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetMarketFeeDiscountParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetMarketFeeDiscountParamsResponse();

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

  fromPartial(_: DeepPartial<MsgSetMarketFeeDiscountParamsResponse>): MsgSetMarketFeeDiscountParamsResponse {
    const message = createBaseMsgSetMarketFeeDiscountParamsResponse();
    return message;
  }

};

function createBaseMsgSetStakingTiers(): MsgSetStakingTiers {
  return {
    authority: "",
    stakingTiers: []
  };
}

export const MsgSetStakingTiers = {
  encode(message: MsgSetStakingTiers, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    for (const v of message.stakingTiers) {
      StakingTier.encode(v!, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetStakingTiers {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetStakingTiers();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
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

  fromPartial(object: DeepPartial<MsgSetStakingTiers>): MsgSetStakingTiers {
    const message = createBaseMsgSetStakingTiers();
    message.authority = object.authority ?? "";
    message.stakingTiers = object.stakingTiers?.map(e => StakingTier.fromPartial(e)) || [];
    return message;
  }

};

function createBaseMsgSetStakingTiersResponse(): MsgSetStakingTiersResponse {
  return {};
}

export const MsgSetStakingTiersResponse = {
  encode(_: MsgSetStakingTiersResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetStakingTiersResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetStakingTiersResponse();

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

  fromPartial(_: DeepPartial<MsgSetStakingTiersResponse>): MsgSetStakingTiersResponse {
    const message = createBaseMsgSetStakingTiersResponse();
    return message;
  }

};