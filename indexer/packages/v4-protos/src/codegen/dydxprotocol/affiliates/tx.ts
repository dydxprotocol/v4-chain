import { AffiliateTiers, AffiliateTiersSDKType } from "./affiliates";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** Message to register a referee-affiliate relationship */

export interface MsgRegisterAffiliate {
  /** Address of the referee */
  referee: string;
  /** Address of the affiliate */

  affiliate: string;
}
/** Message to register a referee-affiliate relationship */

export interface MsgRegisterAffiliateSDKType {
  /** Address of the referee */
  referee: string;
  /** Address of the affiliate */

  affiliate: string;
}
/** Response to MsgRegisterAffiliate */

export interface MsgRegisterAffiliateResponse {}
/** Response to MsgRegisterAffiliate */

export interface MsgRegisterAffiliateResponseSDKType {}
/** Message to update affiliate tiers */

export interface MsgUpdateAffiliateTiers {
  /** Authority sending this message. Will be sent by gov */
  authority: string;
  /** Updated affiliate tiers information */

  tiers?: AffiliateTiers;
}
/** Message to update affiliate tiers */

export interface MsgUpdateAffiliateTiersSDKType {
  /** Authority sending this message. Will be sent by gov */
  authority: string;
  /** Updated affiliate tiers information */

  tiers?: AffiliateTiersSDKType;
}
/** Response to MsgUpdateAffiliateTiers */

export interface MsgUpdateAffiliateTiersResponse {}
/** Response to MsgUpdateAffiliateTiers */

export interface MsgUpdateAffiliateTiersResponseSDKType {}

function createBaseMsgRegisterAffiliate(): MsgRegisterAffiliate {
  return {
    referee: "",
    affiliate: ""
  };
}

export const MsgRegisterAffiliate = {
  encode(message: MsgRegisterAffiliate, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.referee !== "") {
      writer.uint32(10).string(message.referee);
    }

    if (message.affiliate !== "") {
      writer.uint32(18).string(message.affiliate);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgRegisterAffiliate {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRegisterAffiliate();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.referee = reader.string();
          break;

        case 2:
          message.affiliate = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgRegisterAffiliate>): MsgRegisterAffiliate {
    const message = createBaseMsgRegisterAffiliate();
    message.referee = object.referee ?? "";
    message.affiliate = object.affiliate ?? "";
    return message;
  }

};

function createBaseMsgRegisterAffiliateResponse(): MsgRegisterAffiliateResponse {
  return {};
}

export const MsgRegisterAffiliateResponse = {
  encode(_: MsgRegisterAffiliateResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgRegisterAffiliateResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRegisterAffiliateResponse();

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

  fromPartial(_: DeepPartial<MsgRegisterAffiliateResponse>): MsgRegisterAffiliateResponse {
    const message = createBaseMsgRegisterAffiliateResponse();
    return message;
  }

};

function createBaseMsgUpdateAffiliateTiers(): MsgUpdateAffiliateTiers {
  return {
    authority: "",
    tiers: undefined
  };
}

export const MsgUpdateAffiliateTiers = {
  encode(message: MsgUpdateAffiliateTiers, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.tiers !== undefined) {
      AffiliateTiers.encode(message.tiers, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateAffiliateTiers {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateAffiliateTiers();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;

        case 2:
          message.tiers = AffiliateTiers.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgUpdateAffiliateTiers>): MsgUpdateAffiliateTiers {
    const message = createBaseMsgUpdateAffiliateTiers();
    message.authority = object.authority ?? "";
    message.tiers = object.tiers !== undefined && object.tiers !== null ? AffiliateTiers.fromPartial(object.tiers) : undefined;
    return message;
  }

};

function createBaseMsgUpdateAffiliateTiersResponse(): MsgUpdateAffiliateTiersResponse {
  return {};
}

export const MsgUpdateAffiliateTiersResponse = {
  encode(_: MsgUpdateAffiliateTiersResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateAffiliateTiersResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateAffiliateTiersResponse();

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

  fromPartial(_: DeepPartial<MsgUpdateAffiliateTiersResponse>): MsgUpdateAffiliateTiersResponse {
    const message = createBaseMsgUpdateAffiliateTiersResponse();
    return message;
  }

};