import { PerpetualParams, PerpetualParamsAmino, PerpetualParamsSDKType, LiquidityTier, LiquidityTierAmino, LiquidityTierSDKType } from "./perpetual";
import { Params, ParamsAmino, ParamsSDKType } from "./params";
import { BinaryReader, BinaryWriter } from "../../binary";
/** MsgCreatePerpetual is a message used by x/gov to create a new perpetual. */
export interface MsgCreatePerpetual {
  /** The address that controls the module. */
  authority: string;
  /** `params` defines parameters for the new perpetual market. */
  params: PerpetualParams;
}
export interface MsgCreatePerpetualProtoMsg {
  typeUrl: "/dydxprotocol.perpetuals.MsgCreatePerpetual";
  value: Uint8Array;
}
/** MsgCreatePerpetual is a message used by x/gov to create a new perpetual. */
export interface MsgCreatePerpetualAmino {
  /** The address that controls the module. */
  authority?: string;
  /** `params` defines parameters for the new perpetual market. */
  params?: PerpetualParamsAmino;
}
export interface MsgCreatePerpetualAminoMsg {
  type: "/dydxprotocol.perpetuals.MsgCreatePerpetual";
  value: MsgCreatePerpetualAmino;
}
/** MsgCreatePerpetual is a message used by x/gov to create a new perpetual. */
export interface MsgCreatePerpetualSDKType {
  authority: string;
  params: PerpetualParamsSDKType;
}
/**
 * MsgCreatePerpetualResponse defines the CreatePerpetual
 * response type.
 */
export interface MsgCreatePerpetualResponse {}
export interface MsgCreatePerpetualResponseProtoMsg {
  typeUrl: "/dydxprotocol.perpetuals.MsgCreatePerpetualResponse";
  value: Uint8Array;
}
/**
 * MsgCreatePerpetualResponse defines the CreatePerpetual
 * response type.
 */
export interface MsgCreatePerpetualResponseAmino {}
export interface MsgCreatePerpetualResponseAminoMsg {
  type: "/dydxprotocol.perpetuals.MsgCreatePerpetualResponse";
  value: MsgCreatePerpetualResponseAmino;
}
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
  liquidityTier: LiquidityTier;
}
export interface MsgSetLiquidityTierProtoMsg {
  typeUrl: "/dydxprotocol.perpetuals.MsgSetLiquidityTier";
  value: Uint8Array;
}
/**
 * MsgSetLiquidityTier is a message used by x/gov to create or update a
 * liquidity tier.
 */
export interface MsgSetLiquidityTierAmino {
  /** The address that controls the module. */
  authority?: string;
  /** The liquidity tier to create or update. */
  liquidity_tier?: LiquidityTierAmino;
}
export interface MsgSetLiquidityTierAminoMsg {
  type: "/dydxprotocol.perpetuals.MsgSetLiquidityTier";
  value: MsgSetLiquidityTierAmino;
}
/**
 * MsgSetLiquidityTier is a message used by x/gov to create or update a
 * liquidity tier.
 */
export interface MsgSetLiquidityTierSDKType {
  authority: string;
  liquidity_tier: LiquidityTierSDKType;
}
/** MsgSetLiquidityTierResponse defines the SetLiquidityTier response type. */
export interface MsgSetLiquidityTierResponse {}
export interface MsgSetLiquidityTierResponseProtoMsg {
  typeUrl: "/dydxprotocol.perpetuals.MsgSetLiquidityTierResponse";
  value: Uint8Array;
}
/** MsgSetLiquidityTierResponse defines the SetLiquidityTier response type. */
export interface MsgSetLiquidityTierResponseAmino {}
export interface MsgSetLiquidityTierResponseAminoMsg {
  type: "/dydxprotocol.perpetuals.MsgSetLiquidityTierResponse";
  value: MsgSetLiquidityTierResponseAmino;
}
/** MsgSetLiquidityTierResponse defines the SetLiquidityTier response type. */
export interface MsgSetLiquidityTierResponseSDKType {}
/**
 * MsgUpdatePerpetualParams is a message used by x/gov to update the parameters
 * of a perpetual.
 */
export interface MsgUpdatePerpetualParams {
  authority: string;
  /** The perpetual to update. Each field must be set. */
  perpetualParams: PerpetualParams;
}
export interface MsgUpdatePerpetualParamsProtoMsg {
  typeUrl: "/dydxprotocol.perpetuals.MsgUpdatePerpetualParams";
  value: Uint8Array;
}
/**
 * MsgUpdatePerpetualParams is a message used by x/gov to update the parameters
 * of a perpetual.
 */
export interface MsgUpdatePerpetualParamsAmino {
  authority?: string;
  /** The perpetual to update. Each field must be set. */
  perpetual_params?: PerpetualParamsAmino;
}
export interface MsgUpdatePerpetualParamsAminoMsg {
  type: "/dydxprotocol.perpetuals.MsgUpdatePerpetualParams";
  value: MsgUpdatePerpetualParamsAmino;
}
/**
 * MsgUpdatePerpetualParams is a message used by x/gov to update the parameters
 * of a perpetual.
 */
export interface MsgUpdatePerpetualParamsSDKType {
  authority: string;
  perpetual_params: PerpetualParamsSDKType;
}
/**
 * MsgUpdatePerpetualParamsResponse defines the UpdatePerpetualParams
 * response type.
 */
export interface MsgUpdatePerpetualParamsResponse {}
export interface MsgUpdatePerpetualParamsResponseProtoMsg {
  typeUrl: "/dydxprotocol.perpetuals.MsgUpdatePerpetualParamsResponse";
  value: Uint8Array;
}
/**
 * MsgUpdatePerpetualParamsResponse defines the UpdatePerpetualParams
 * response type.
 */
export interface MsgUpdatePerpetualParamsResponseAmino {}
export interface MsgUpdatePerpetualParamsResponseAminoMsg {
  type: "/dydxprotocol.perpetuals.MsgUpdatePerpetualParamsResponse";
  value: MsgUpdatePerpetualParamsResponseAmino;
}
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
export interface FundingPremiumProtoMsg {
  typeUrl: "/dydxprotocol.perpetuals.FundingPremium";
  value: Uint8Array;
}
/**
 * FundingPremium represents a funding premium value for a perpetual
 * market. Can be used to represent a premium vote or a premium sample.
 */
export interface FundingPremiumAmino {
  /** The id of the perpetual market. */
  perpetual_id?: number;
  /** The sampled premium rate. In parts-per-million. */
  premium_ppm?: number;
}
export interface FundingPremiumAminoMsg {
  type: "/dydxprotocol.perpetuals.FundingPremium";
  value: FundingPremiumAmino;
}
/**
 * FundingPremium represents a funding premium value for a perpetual
 * market. Can be used to represent a premium vote or a premium sample.
 */
export interface FundingPremiumSDKType {
  perpetual_id: number;
  premium_ppm: number;
}
/** MsgAddPremiumVotes is a request type for the AddPremiumVotes method. */
export interface MsgAddPremiumVotes {
  votes: FundingPremium[];
}
export interface MsgAddPremiumVotesProtoMsg {
  typeUrl: "/dydxprotocol.perpetuals.MsgAddPremiumVotes";
  value: Uint8Array;
}
/** MsgAddPremiumVotes is a request type for the AddPremiumVotes method. */
export interface MsgAddPremiumVotesAmino {
  votes?: FundingPremiumAmino[];
}
export interface MsgAddPremiumVotesAminoMsg {
  type: "/dydxprotocol.perpetuals.MsgAddPremiumVotes";
  value: MsgAddPremiumVotesAmino;
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
export interface MsgAddPremiumVotesResponseProtoMsg {
  typeUrl: "/dydxprotocol.perpetuals.MsgAddPremiumVotesResponse";
  value: Uint8Array;
}
/**
 * MsgAddPremiumVotesResponse defines the AddPremiumVotes
 * response type.
 */
export interface MsgAddPremiumVotesResponseAmino {}
export interface MsgAddPremiumVotesResponseAminoMsg {
  type: "/dydxprotocol.perpetuals.MsgAddPremiumVotesResponse";
  value: MsgAddPremiumVotesResponseAmino;
}
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
  params: Params;
}
export interface MsgUpdateParamsProtoMsg {
  typeUrl: "/dydxprotocol.perpetuals.MsgUpdateParams";
  value: Uint8Array;
}
/**
 * MsgUpdateParams is a message used by x/gov to update the parameters of the
 * perpetuals module.
 */
export interface MsgUpdateParamsAmino {
  authority?: string;
  /** The parameters to update. Each field must be set. */
  params?: ParamsAmino;
}
export interface MsgUpdateParamsAminoMsg {
  type: "/dydxprotocol.perpetuals.MsgUpdateParams";
  value: MsgUpdateParamsAmino;
}
/**
 * MsgUpdateParams is a message used by x/gov to update the parameters of the
 * perpetuals module.
 */
export interface MsgUpdateParamsSDKType {
  authority: string;
  params: ParamsSDKType;
}
/** MsgUpdateParamsResponse defines the UpdateParams response type. */
export interface MsgUpdateParamsResponse {}
export interface MsgUpdateParamsResponseProtoMsg {
  typeUrl: "/dydxprotocol.perpetuals.MsgUpdateParamsResponse";
  value: Uint8Array;
}
/** MsgUpdateParamsResponse defines the UpdateParams response type. */
export interface MsgUpdateParamsResponseAmino {}
export interface MsgUpdateParamsResponseAminoMsg {
  type: "/dydxprotocol.perpetuals.MsgUpdateParamsResponse";
  value: MsgUpdateParamsResponseAmino;
}
/** MsgUpdateParamsResponse defines the UpdateParams response type. */
export interface MsgUpdateParamsResponseSDKType {}
function createBaseMsgCreatePerpetual(): MsgCreatePerpetual {
  return {
    authority: "",
    params: PerpetualParams.fromPartial({})
  };
}
export const MsgCreatePerpetual = {
  typeUrl: "/dydxprotocol.perpetuals.MsgCreatePerpetual",
  encode(message: MsgCreatePerpetual, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }
    if (message.params !== undefined) {
      PerpetualParams.encode(message.params, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgCreatePerpetual {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<MsgCreatePerpetual>): MsgCreatePerpetual {
    const message = createBaseMsgCreatePerpetual();
    message.authority = object.authority ?? "";
    message.params = object.params !== undefined && object.params !== null ? PerpetualParams.fromPartial(object.params) : undefined;
    return message;
  },
  fromAmino(object: MsgCreatePerpetualAmino): MsgCreatePerpetual {
    const message = createBaseMsgCreatePerpetual();
    if (object.authority !== undefined && object.authority !== null) {
      message.authority = object.authority;
    }
    if (object.params !== undefined && object.params !== null) {
      message.params = PerpetualParams.fromAmino(object.params);
    }
    return message;
  },
  toAmino(message: MsgCreatePerpetual): MsgCreatePerpetualAmino {
    const obj: any = {};
    obj.authority = message.authority;
    obj.params = message.params ? PerpetualParams.toAmino(message.params) : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgCreatePerpetualAminoMsg): MsgCreatePerpetual {
    return MsgCreatePerpetual.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgCreatePerpetualProtoMsg): MsgCreatePerpetual {
    return MsgCreatePerpetual.decode(message.value);
  },
  toProto(message: MsgCreatePerpetual): Uint8Array {
    return MsgCreatePerpetual.encode(message).finish();
  },
  toProtoMsg(message: MsgCreatePerpetual): MsgCreatePerpetualProtoMsg {
    return {
      typeUrl: "/dydxprotocol.perpetuals.MsgCreatePerpetual",
      value: MsgCreatePerpetual.encode(message).finish()
    };
  }
};
function createBaseMsgCreatePerpetualResponse(): MsgCreatePerpetualResponse {
  return {};
}
export const MsgCreatePerpetualResponse = {
  typeUrl: "/dydxprotocol.perpetuals.MsgCreatePerpetualResponse",
  encode(_: MsgCreatePerpetualResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgCreatePerpetualResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(_: Partial<MsgCreatePerpetualResponse>): MsgCreatePerpetualResponse {
    const message = createBaseMsgCreatePerpetualResponse();
    return message;
  },
  fromAmino(_: MsgCreatePerpetualResponseAmino): MsgCreatePerpetualResponse {
    const message = createBaseMsgCreatePerpetualResponse();
    return message;
  },
  toAmino(_: MsgCreatePerpetualResponse): MsgCreatePerpetualResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgCreatePerpetualResponseAminoMsg): MsgCreatePerpetualResponse {
    return MsgCreatePerpetualResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgCreatePerpetualResponseProtoMsg): MsgCreatePerpetualResponse {
    return MsgCreatePerpetualResponse.decode(message.value);
  },
  toProto(message: MsgCreatePerpetualResponse): Uint8Array {
    return MsgCreatePerpetualResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgCreatePerpetualResponse): MsgCreatePerpetualResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.perpetuals.MsgCreatePerpetualResponse",
      value: MsgCreatePerpetualResponse.encode(message).finish()
    };
  }
};
function createBaseMsgSetLiquidityTier(): MsgSetLiquidityTier {
  return {
    authority: "",
    liquidityTier: LiquidityTier.fromPartial({})
  };
}
export const MsgSetLiquidityTier = {
  typeUrl: "/dydxprotocol.perpetuals.MsgSetLiquidityTier",
  encode(message: MsgSetLiquidityTier, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }
    if (message.liquidityTier !== undefined) {
      LiquidityTier.encode(message.liquidityTier, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgSetLiquidityTier {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<MsgSetLiquidityTier>): MsgSetLiquidityTier {
    const message = createBaseMsgSetLiquidityTier();
    message.authority = object.authority ?? "";
    message.liquidityTier = object.liquidityTier !== undefined && object.liquidityTier !== null ? LiquidityTier.fromPartial(object.liquidityTier) : undefined;
    return message;
  },
  fromAmino(object: MsgSetLiquidityTierAmino): MsgSetLiquidityTier {
    const message = createBaseMsgSetLiquidityTier();
    if (object.authority !== undefined && object.authority !== null) {
      message.authority = object.authority;
    }
    if (object.liquidity_tier !== undefined && object.liquidity_tier !== null) {
      message.liquidityTier = LiquidityTier.fromAmino(object.liquidity_tier);
    }
    return message;
  },
  toAmino(message: MsgSetLiquidityTier): MsgSetLiquidityTierAmino {
    const obj: any = {};
    obj.authority = message.authority;
    obj.liquidity_tier = message.liquidityTier ? LiquidityTier.toAmino(message.liquidityTier) : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgSetLiquidityTierAminoMsg): MsgSetLiquidityTier {
    return MsgSetLiquidityTier.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgSetLiquidityTierProtoMsg): MsgSetLiquidityTier {
    return MsgSetLiquidityTier.decode(message.value);
  },
  toProto(message: MsgSetLiquidityTier): Uint8Array {
    return MsgSetLiquidityTier.encode(message).finish();
  },
  toProtoMsg(message: MsgSetLiquidityTier): MsgSetLiquidityTierProtoMsg {
    return {
      typeUrl: "/dydxprotocol.perpetuals.MsgSetLiquidityTier",
      value: MsgSetLiquidityTier.encode(message).finish()
    };
  }
};
function createBaseMsgSetLiquidityTierResponse(): MsgSetLiquidityTierResponse {
  return {};
}
export const MsgSetLiquidityTierResponse = {
  typeUrl: "/dydxprotocol.perpetuals.MsgSetLiquidityTierResponse",
  encode(_: MsgSetLiquidityTierResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgSetLiquidityTierResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(_: Partial<MsgSetLiquidityTierResponse>): MsgSetLiquidityTierResponse {
    const message = createBaseMsgSetLiquidityTierResponse();
    return message;
  },
  fromAmino(_: MsgSetLiquidityTierResponseAmino): MsgSetLiquidityTierResponse {
    const message = createBaseMsgSetLiquidityTierResponse();
    return message;
  },
  toAmino(_: MsgSetLiquidityTierResponse): MsgSetLiquidityTierResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgSetLiquidityTierResponseAminoMsg): MsgSetLiquidityTierResponse {
    return MsgSetLiquidityTierResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgSetLiquidityTierResponseProtoMsg): MsgSetLiquidityTierResponse {
    return MsgSetLiquidityTierResponse.decode(message.value);
  },
  toProto(message: MsgSetLiquidityTierResponse): Uint8Array {
    return MsgSetLiquidityTierResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgSetLiquidityTierResponse): MsgSetLiquidityTierResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.perpetuals.MsgSetLiquidityTierResponse",
      value: MsgSetLiquidityTierResponse.encode(message).finish()
    };
  }
};
function createBaseMsgUpdatePerpetualParams(): MsgUpdatePerpetualParams {
  return {
    authority: "",
    perpetualParams: PerpetualParams.fromPartial({})
  };
}
export const MsgUpdatePerpetualParams = {
  typeUrl: "/dydxprotocol.perpetuals.MsgUpdatePerpetualParams",
  encode(message: MsgUpdatePerpetualParams, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }
    if (message.perpetualParams !== undefined) {
      PerpetualParams.encode(message.perpetualParams, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdatePerpetualParams {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<MsgUpdatePerpetualParams>): MsgUpdatePerpetualParams {
    const message = createBaseMsgUpdatePerpetualParams();
    message.authority = object.authority ?? "";
    message.perpetualParams = object.perpetualParams !== undefined && object.perpetualParams !== null ? PerpetualParams.fromPartial(object.perpetualParams) : undefined;
    return message;
  },
  fromAmino(object: MsgUpdatePerpetualParamsAmino): MsgUpdatePerpetualParams {
    const message = createBaseMsgUpdatePerpetualParams();
    if (object.authority !== undefined && object.authority !== null) {
      message.authority = object.authority;
    }
    if (object.perpetual_params !== undefined && object.perpetual_params !== null) {
      message.perpetualParams = PerpetualParams.fromAmino(object.perpetual_params);
    }
    return message;
  },
  toAmino(message: MsgUpdatePerpetualParams): MsgUpdatePerpetualParamsAmino {
    const obj: any = {};
    obj.authority = message.authority;
    obj.perpetual_params = message.perpetualParams ? PerpetualParams.toAmino(message.perpetualParams) : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgUpdatePerpetualParamsAminoMsg): MsgUpdatePerpetualParams {
    return MsgUpdatePerpetualParams.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdatePerpetualParamsProtoMsg): MsgUpdatePerpetualParams {
    return MsgUpdatePerpetualParams.decode(message.value);
  },
  toProto(message: MsgUpdatePerpetualParams): Uint8Array {
    return MsgUpdatePerpetualParams.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdatePerpetualParams): MsgUpdatePerpetualParamsProtoMsg {
    return {
      typeUrl: "/dydxprotocol.perpetuals.MsgUpdatePerpetualParams",
      value: MsgUpdatePerpetualParams.encode(message).finish()
    };
  }
};
function createBaseMsgUpdatePerpetualParamsResponse(): MsgUpdatePerpetualParamsResponse {
  return {};
}
export const MsgUpdatePerpetualParamsResponse = {
  typeUrl: "/dydxprotocol.perpetuals.MsgUpdatePerpetualParamsResponse",
  encode(_: MsgUpdatePerpetualParamsResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdatePerpetualParamsResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(_: Partial<MsgUpdatePerpetualParamsResponse>): MsgUpdatePerpetualParamsResponse {
    const message = createBaseMsgUpdatePerpetualParamsResponse();
    return message;
  },
  fromAmino(_: MsgUpdatePerpetualParamsResponseAmino): MsgUpdatePerpetualParamsResponse {
    const message = createBaseMsgUpdatePerpetualParamsResponse();
    return message;
  },
  toAmino(_: MsgUpdatePerpetualParamsResponse): MsgUpdatePerpetualParamsResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgUpdatePerpetualParamsResponseAminoMsg): MsgUpdatePerpetualParamsResponse {
    return MsgUpdatePerpetualParamsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdatePerpetualParamsResponseProtoMsg): MsgUpdatePerpetualParamsResponse {
    return MsgUpdatePerpetualParamsResponse.decode(message.value);
  },
  toProto(message: MsgUpdatePerpetualParamsResponse): Uint8Array {
    return MsgUpdatePerpetualParamsResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdatePerpetualParamsResponse): MsgUpdatePerpetualParamsResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.perpetuals.MsgUpdatePerpetualParamsResponse",
      value: MsgUpdatePerpetualParamsResponse.encode(message).finish()
    };
  }
};
function createBaseFundingPremium(): FundingPremium {
  return {
    perpetualId: 0,
    premiumPpm: 0
  };
}
export const FundingPremium = {
  typeUrl: "/dydxprotocol.perpetuals.FundingPremium",
  encode(message: FundingPremium, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.perpetualId !== 0) {
      writer.uint32(8).uint32(message.perpetualId);
    }
    if (message.premiumPpm !== 0) {
      writer.uint32(16).int32(message.premiumPpm);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): FundingPremium {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<FundingPremium>): FundingPremium {
    const message = createBaseFundingPremium();
    message.perpetualId = object.perpetualId ?? 0;
    message.premiumPpm = object.premiumPpm ?? 0;
    return message;
  },
  fromAmino(object: FundingPremiumAmino): FundingPremium {
    const message = createBaseFundingPremium();
    if (object.perpetual_id !== undefined && object.perpetual_id !== null) {
      message.perpetualId = object.perpetual_id;
    }
    if (object.premium_ppm !== undefined && object.premium_ppm !== null) {
      message.premiumPpm = object.premium_ppm;
    }
    return message;
  },
  toAmino(message: FundingPremium): FundingPremiumAmino {
    const obj: any = {};
    obj.perpetual_id = message.perpetualId;
    obj.premium_ppm = message.premiumPpm;
    return obj;
  },
  fromAminoMsg(object: FundingPremiumAminoMsg): FundingPremium {
    return FundingPremium.fromAmino(object.value);
  },
  fromProtoMsg(message: FundingPremiumProtoMsg): FundingPremium {
    return FundingPremium.decode(message.value);
  },
  toProto(message: FundingPremium): Uint8Array {
    return FundingPremium.encode(message).finish();
  },
  toProtoMsg(message: FundingPremium): FundingPremiumProtoMsg {
    return {
      typeUrl: "/dydxprotocol.perpetuals.FundingPremium",
      value: FundingPremium.encode(message).finish()
    };
  }
};
function createBaseMsgAddPremiumVotes(): MsgAddPremiumVotes {
  return {
    votes: []
  };
}
export const MsgAddPremiumVotes = {
  typeUrl: "/dydxprotocol.perpetuals.MsgAddPremiumVotes",
  encode(message: MsgAddPremiumVotes, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.votes) {
      FundingPremium.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgAddPremiumVotes {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<MsgAddPremiumVotes>): MsgAddPremiumVotes {
    const message = createBaseMsgAddPremiumVotes();
    message.votes = object.votes?.map(e => FundingPremium.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: MsgAddPremiumVotesAmino): MsgAddPremiumVotes {
    const message = createBaseMsgAddPremiumVotes();
    message.votes = object.votes?.map(e => FundingPremium.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: MsgAddPremiumVotes): MsgAddPremiumVotesAmino {
    const obj: any = {};
    if (message.votes) {
      obj.votes = message.votes.map(e => e ? FundingPremium.toAmino(e) : undefined);
    } else {
      obj.votes = [];
    }
    return obj;
  },
  fromAminoMsg(object: MsgAddPremiumVotesAminoMsg): MsgAddPremiumVotes {
    return MsgAddPremiumVotes.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgAddPremiumVotesProtoMsg): MsgAddPremiumVotes {
    return MsgAddPremiumVotes.decode(message.value);
  },
  toProto(message: MsgAddPremiumVotes): Uint8Array {
    return MsgAddPremiumVotes.encode(message).finish();
  },
  toProtoMsg(message: MsgAddPremiumVotes): MsgAddPremiumVotesProtoMsg {
    return {
      typeUrl: "/dydxprotocol.perpetuals.MsgAddPremiumVotes",
      value: MsgAddPremiumVotes.encode(message).finish()
    };
  }
};
function createBaseMsgAddPremiumVotesResponse(): MsgAddPremiumVotesResponse {
  return {};
}
export const MsgAddPremiumVotesResponse = {
  typeUrl: "/dydxprotocol.perpetuals.MsgAddPremiumVotesResponse",
  encode(_: MsgAddPremiumVotesResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgAddPremiumVotesResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(_: Partial<MsgAddPremiumVotesResponse>): MsgAddPremiumVotesResponse {
    const message = createBaseMsgAddPremiumVotesResponse();
    return message;
  },
  fromAmino(_: MsgAddPremiumVotesResponseAmino): MsgAddPremiumVotesResponse {
    const message = createBaseMsgAddPremiumVotesResponse();
    return message;
  },
  toAmino(_: MsgAddPremiumVotesResponse): MsgAddPremiumVotesResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgAddPremiumVotesResponseAminoMsg): MsgAddPremiumVotesResponse {
    return MsgAddPremiumVotesResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgAddPremiumVotesResponseProtoMsg): MsgAddPremiumVotesResponse {
    return MsgAddPremiumVotesResponse.decode(message.value);
  },
  toProto(message: MsgAddPremiumVotesResponse): Uint8Array {
    return MsgAddPremiumVotesResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgAddPremiumVotesResponse): MsgAddPremiumVotesResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.perpetuals.MsgAddPremiumVotesResponse",
      value: MsgAddPremiumVotesResponse.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateParams(): MsgUpdateParams {
  return {
    authority: "",
    params: Params.fromPartial({})
  };
}
export const MsgUpdateParams = {
  typeUrl: "/dydxprotocol.perpetuals.MsgUpdateParams",
  encode(message: MsgUpdateParams, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateParams {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<MsgUpdateParams>): MsgUpdateParams {
    const message = createBaseMsgUpdateParams();
    message.authority = object.authority ?? "";
    message.params = object.params !== undefined && object.params !== null ? Params.fromPartial(object.params) : undefined;
    return message;
  },
  fromAmino(object: MsgUpdateParamsAmino): MsgUpdateParams {
    const message = createBaseMsgUpdateParams();
    if (object.authority !== undefined && object.authority !== null) {
      message.authority = object.authority;
    }
    if (object.params !== undefined && object.params !== null) {
      message.params = Params.fromAmino(object.params);
    }
    return message;
  },
  toAmino(message: MsgUpdateParams): MsgUpdateParamsAmino {
    const obj: any = {};
    obj.authority = message.authority;
    obj.params = message.params ? Params.toAmino(message.params) : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgUpdateParamsAminoMsg): MsgUpdateParams {
    return MsgUpdateParams.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateParamsProtoMsg): MsgUpdateParams {
    return MsgUpdateParams.decode(message.value);
  },
  toProto(message: MsgUpdateParams): Uint8Array {
    return MsgUpdateParams.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateParams): MsgUpdateParamsProtoMsg {
    return {
      typeUrl: "/dydxprotocol.perpetuals.MsgUpdateParams",
      value: MsgUpdateParams.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateParamsResponse(): MsgUpdateParamsResponse {
  return {};
}
export const MsgUpdateParamsResponse = {
  typeUrl: "/dydxprotocol.perpetuals.MsgUpdateParamsResponse",
  encode(_: MsgUpdateParamsResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateParamsResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(_: Partial<MsgUpdateParamsResponse>): MsgUpdateParamsResponse {
    const message = createBaseMsgUpdateParamsResponse();
    return message;
  },
  fromAmino(_: MsgUpdateParamsResponseAmino): MsgUpdateParamsResponse {
    const message = createBaseMsgUpdateParamsResponse();
    return message;
  },
  toAmino(_: MsgUpdateParamsResponse): MsgUpdateParamsResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgUpdateParamsResponseAminoMsg): MsgUpdateParamsResponse {
    return MsgUpdateParamsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateParamsResponseProtoMsg): MsgUpdateParamsResponse {
    return MsgUpdateParamsResponse.decode(message.value);
  },
  toProto(message: MsgUpdateParamsResponse): Uint8Array {
    return MsgUpdateParamsResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateParamsResponse): MsgUpdateParamsResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.perpetuals.MsgUpdateParamsResponse",
      value: MsgUpdateParamsResponse.encode(message).finish()
    };
  }
};