import { PerpetualParams, PerpetualParamsSDKType, LiquidityTier, LiquidityTierSDKType } from "./perpetual";
import { Params, ParamsSDKType } from "./params";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** MsgCreatePerpetual is a message used by x/gov to create a new perpetual. */

export interface MsgCreatePerpetual {
  /** The address that controls the module. */
  authority: string;
  /** `params` defines parameters for the new perpetual market. */

  params?: PerpetualParams;
}
/** MsgCreatePerpetual is a message used by x/gov to create a new perpetual. */

export interface MsgCreatePerpetualSDKType {
  /** The address that controls the module. */
  authority: string;
  /** `params` defines parameters for the new perpetual market. */

  params?: PerpetualParamsSDKType;
}
/**
 * MsgCreatePerpetualResponse defines the CreatePerpetual
 * response type.
 */

export interface MsgCreatePerpetualResponse {}
/**
 * MsgCreatePerpetualResponse defines the CreatePerpetual
 * response type.
 */

export interface MsgCreatePerpetualResponseSDKType {}
/**
 * MsgSetLiquidityTier is a message used by x/gov to create or update a
 * liquidity tier.
 */

export interface MsgSetLiquidityTier {
  /** The address that controls the module. */
  authority: string;
  /** The liquidity tier to create or update. */

  liquidityTier?: LiquidityTier;
}
/**
 * MsgSetLiquidityTier is a message used by x/gov to create or update a
 * liquidity tier.
 */

export interface MsgSetLiquidityTierSDKType {
  /** The address that controls the module. */
  authority: string;
  /** The liquidity tier to create or update. */

  liquidity_tier?: LiquidityTierSDKType;
}
/** MsgSetLiquidityTierResponse defines the SetLiquidityTier response type. */

export interface MsgSetLiquidityTierResponse {}
/** MsgSetLiquidityTierResponse defines the SetLiquidityTier response type. */

export interface MsgSetLiquidityTierResponseSDKType {}
/**
 * MsgUpdatePerpetualParams is a message used by x/gov to update the parameters
 * of a perpetual.
 */

export interface MsgUpdatePerpetualParams {
  authority: string;
  /** The perpetual to update. Each field must be set. */

  perpetualParams?: PerpetualParams;
}
/**
 * MsgUpdatePerpetualParams is a message used by x/gov to update the parameters
 * of a perpetual.
 */

export interface MsgUpdatePerpetualParamsSDKType {
  authority: string;
  /** The perpetual to update. Each field must be set. */

  perpetual_params?: PerpetualParamsSDKType;
}
/**
 * MsgUpdatePerpetualParamsResponse defines the UpdatePerpetualParams
 * response type.
 */

export interface MsgUpdatePerpetualParamsResponse {}
/**
 * MsgUpdatePerpetualParamsResponse defines the UpdatePerpetualParams
 * response type.
 */

export interface MsgUpdatePerpetualParamsResponseSDKType {}
/**
 * FundingPremium represents a funding premium value for a perpetual
 * market. Can be used to represent a premium vote or a premium sample.
 */

export interface FundingPremium {
  /** The id of the perpetual market. */
  perpetualId: number;
  /** The sampled premium rate. In parts-per-million. */

  premiumPpm: number;
}
/**
 * FundingPremium represents a funding premium value for a perpetual
 * market. Can be used to represent a premium vote or a premium sample.
 */

export interface FundingPremiumSDKType {
  /** The id of the perpetual market. */
  perpetual_id: number;
  /** The sampled premium rate. In parts-per-million. */

  premium_ppm: number;
}
/** MsgAddPremiumVotes is a request type for the AddPremiumVotes method. */

export interface MsgAddPremiumVotes {
  votes: FundingPremium[];
}
/** MsgAddPremiumVotes is a request type for the AddPremiumVotes method. */

export interface MsgAddPremiumVotesSDKType {
  votes: FundingPremiumSDKType[];
}
/**
 * MsgAddPremiumVotesResponse defines the AddPremiumVotes
 * response type.
 */

export interface MsgAddPremiumVotesResponse {}
/**
 * MsgAddPremiumVotesResponse defines the AddPremiumVotes
 * response type.
 */

export interface MsgAddPremiumVotesResponseSDKType {}
/**
 * MsgUpdateParams is a message used by x/gov to update the parameters of the
 * perpetuals module.
 */

export interface MsgUpdateParams {
  authority: string;
  /** The parameters to update. Each field must be set. */

  params?: Params;
}
/**
 * MsgUpdateParams is a message used by x/gov to update the parameters of the
 * perpetuals module.
 */

export interface MsgUpdateParamsSDKType {
  authority: string;
  /** The parameters to update. Each field must be set. */

  params?: ParamsSDKType;
}
/** MsgUpdateParamsResponse defines the UpdateParams response type. */

export interface MsgUpdateParamsResponse {}
/** MsgUpdateParamsResponse defines the UpdateParams response type. */

export interface MsgUpdateParamsResponseSDKType {}

function createBaseMsgCreatePerpetual(): MsgCreatePerpetual {
  return {
    authority: "",
    params: undefined
  };
}

export const MsgCreatePerpetual = {
  encode(message: MsgCreatePerpetual, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.params !== undefined) {
      PerpetualParams.encode(message.params, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgCreatePerpetual {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCreatePerpetual();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;

        case 2:
          message.params = PerpetualParams.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgCreatePerpetual>): MsgCreatePerpetual {
    const message = createBaseMsgCreatePerpetual();
    message.authority = object.authority ?? "";
    message.params = object.params !== undefined && object.params !== null ? PerpetualParams.fromPartial(object.params) : undefined;
    return message;
  }

};

function createBaseMsgCreatePerpetualResponse(): MsgCreatePerpetualResponse {
  return {};
}

export const MsgCreatePerpetualResponse = {
  encode(_: MsgCreatePerpetualResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgCreatePerpetualResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCreatePerpetualResponse();

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

  fromPartial(_: DeepPartial<MsgCreatePerpetualResponse>): MsgCreatePerpetualResponse {
    const message = createBaseMsgCreatePerpetualResponse();
    return message;
  }

};

function createBaseMsgSetLiquidityTier(): MsgSetLiquidityTier {
  return {
    authority: "",
    liquidityTier: undefined
  };
}

export const MsgSetLiquidityTier = {
  encode(message: MsgSetLiquidityTier, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.liquidityTier !== undefined) {
      LiquidityTier.encode(message.liquidityTier, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetLiquidityTier {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetLiquidityTier();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;

        case 2:
          message.liquidityTier = LiquidityTier.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgSetLiquidityTier>): MsgSetLiquidityTier {
    const message = createBaseMsgSetLiquidityTier();
    message.authority = object.authority ?? "";
    message.liquidityTier = object.liquidityTier !== undefined && object.liquidityTier !== null ? LiquidityTier.fromPartial(object.liquidityTier) : undefined;
    return message;
  }

};

function createBaseMsgSetLiquidityTierResponse(): MsgSetLiquidityTierResponse {
  return {};
}

export const MsgSetLiquidityTierResponse = {
  encode(_: MsgSetLiquidityTierResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetLiquidityTierResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetLiquidityTierResponse();

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

  fromPartial(_: DeepPartial<MsgSetLiquidityTierResponse>): MsgSetLiquidityTierResponse {
    const message = createBaseMsgSetLiquidityTierResponse();
    return message;
  }

};

function createBaseMsgUpdatePerpetualParams(): MsgUpdatePerpetualParams {
  return {
    authority: "",
    perpetualParams: undefined
  };
}

export const MsgUpdatePerpetualParams = {
  encode(message: MsgUpdatePerpetualParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.perpetualParams !== undefined) {
      PerpetualParams.encode(message.perpetualParams, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdatePerpetualParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdatePerpetualParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;

        case 2:
          message.perpetualParams = PerpetualParams.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgUpdatePerpetualParams>): MsgUpdatePerpetualParams {
    const message = createBaseMsgUpdatePerpetualParams();
    message.authority = object.authority ?? "";
    message.perpetualParams = object.perpetualParams !== undefined && object.perpetualParams !== null ? PerpetualParams.fromPartial(object.perpetualParams) : undefined;
    return message;
  }

};

function createBaseMsgUpdatePerpetualParamsResponse(): MsgUpdatePerpetualParamsResponse {
  return {};
}

export const MsgUpdatePerpetualParamsResponse = {
  encode(_: MsgUpdatePerpetualParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdatePerpetualParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdatePerpetualParamsResponse();

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

  fromPartial(_: DeepPartial<MsgUpdatePerpetualParamsResponse>): MsgUpdatePerpetualParamsResponse {
    const message = createBaseMsgUpdatePerpetualParamsResponse();
    return message;
  }

};

function createBaseFundingPremium(): FundingPremium {
  return {
    perpetualId: 0,
    premiumPpm: 0
  };
}

export const FundingPremium = {
  encode(message: FundingPremium, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.perpetualId !== 0) {
      writer.uint32(8).uint32(message.perpetualId);
    }

    if (message.premiumPpm !== 0) {
      writer.uint32(16).int32(message.premiumPpm);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FundingPremium {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFundingPremium();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.perpetualId = reader.uint32();
          break;

        case 2:
          message.premiumPpm = reader.int32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<FundingPremium>): FundingPremium {
    const message = createBaseFundingPremium();
    message.perpetualId = object.perpetualId ?? 0;
    message.premiumPpm = object.premiumPpm ?? 0;
    return message;
  }

};

function createBaseMsgAddPremiumVotes(): MsgAddPremiumVotes {
  return {
    votes: []
  };
}

export const MsgAddPremiumVotes = {
  encode(message: MsgAddPremiumVotes, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.votes) {
      FundingPremium.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgAddPremiumVotes {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgAddPremiumVotes();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.votes.push(FundingPremium.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgAddPremiumVotes>): MsgAddPremiumVotes {
    const message = createBaseMsgAddPremiumVotes();
    message.votes = object.votes?.map(e => FundingPremium.fromPartial(e)) || [];
    return message;
  }

};

function createBaseMsgAddPremiumVotesResponse(): MsgAddPremiumVotesResponse {
  return {};
}

export const MsgAddPremiumVotesResponse = {
  encode(_: MsgAddPremiumVotesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgAddPremiumVotesResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgAddPremiumVotesResponse();

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

  fromPartial(_: DeepPartial<MsgAddPremiumVotesResponse>): MsgAddPremiumVotesResponse {
    const message = createBaseMsgAddPremiumVotesResponse();
    return message;
  }

};

function createBaseMsgUpdateParams(): MsgUpdateParams {
  return {
    authority: "",
    params: undefined
  };
}

export const MsgUpdateParams = {
  encode(message: MsgUpdateParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;

        case 2:
          message.params = Params.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgUpdateParams>): MsgUpdateParams {
    const message = createBaseMsgUpdateParams();
    message.authority = object.authority ?? "";
    message.params = object.params !== undefined && object.params !== null ? Params.fromPartial(object.params) : undefined;
    return message;
  }

};

function createBaseMsgUpdateParamsResponse(): MsgUpdateParamsResponse {
  return {};
}

export const MsgUpdateParamsResponse = {
  encode(_: MsgUpdateParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateParamsResponse();

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

  fromPartial(_: DeepPartial<MsgUpdateParamsResponse>): MsgUpdateParamsResponse {
    const message = createBaseMsgUpdateParamsResponse();
    return message;
  }

};