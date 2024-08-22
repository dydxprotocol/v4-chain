import { AffiliateTiers, AffiliateTiersSDKType } from "./affiliates";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/**
 * AffiliateInfoRequest is the request type for the Query/AffiliateInfo RPC
 * method.
 */

export interface AffiliateInfoRequest {
  address: string;
}
/**
 * AffiliateInfoRequest is the request type for the Query/AffiliateInfo RPC
 * method.
 */

export interface AffiliateInfoRequestSDKType {
  address: string;
}
/**
 * AffiliateInfoResponse is the response type for the Query/AffiliateInfo RPC
 * method.
 */

export interface AffiliateInfoResponse {
  /** The affiliate's tier. */
  tier: number;
  /** The affiliate's taker fee share in parts-per-million. */

  feeSharePpm: number;
  /** The affiliate's all-time referred volume in quote quantums. */

  referredVolume: Uint8Array;
  /** The affiliate's currently staked native tokens (in whole coins). */

  stakedAmount: Uint8Array;
}
/**
 * AffiliateInfoResponse is the response type for the Query/AffiliateInfo RPC
 * method.
 */

export interface AffiliateInfoResponseSDKType {
  /** The affiliate's tier. */
  tier: number;
  /** The affiliate's taker fee share in parts-per-million. */

  fee_share_ppm: number;
  /** The affiliate's all-time referred volume in quote quantums. */

  referred_volume: Uint8Array;
  /** The affiliate's currently staked native tokens (in whole coins). */

  staked_amount: Uint8Array;
}
/** ReferredByRequest is the request type for the Query/ReferredBy RPC method. */

export interface ReferredByRequest {
  /** The address to query. */
  address: string;
}
/** ReferredByRequest is the request type for the Query/ReferredBy RPC method. */

export interface ReferredByRequestSDKType {
  /** The address to query. */
  address: string;
}
/** ReferredByResponse is the response type for the Query/ReferredBy RPC method. */

export interface ReferredByResponse {
  /** The affiliate's address that referred the queried address. */
  affiliateAddress: string;
}
/** ReferredByResponse is the response type for the Query/ReferredBy RPC method. */

export interface ReferredByResponseSDKType {
  /** The affiliate's address that referred the queried address. */
  affiliate_address: string;
}
/**
 * AllAffiliateTiersRequest is the request type for the Query/AllAffiliateTiers
 * RPC method.
 */

export interface AllAffiliateTiersRequest {}
/**
 * AllAffiliateTiersRequest is the request type for the Query/AllAffiliateTiers
 * RPC method.
 */

export interface AllAffiliateTiersRequestSDKType {}
/**
 * AllAffiliateTiersResponse is the response type for the
 * Query/AllAffiliateTiers RPC method.
 */

export interface AllAffiliateTiersResponse {
  /** All affiliate tiers information. */
  tiers?: AffiliateTiers;
}
/**
 * AllAffiliateTiersResponse is the response type for the
 * Query/AllAffiliateTiers RPC method.
 */

export interface AllAffiliateTiersResponseSDKType {
  /** All affiliate tiers information. */
  tiers?: AffiliateTiersSDKType;
}

function createBaseAffiliateInfoRequest(): AffiliateInfoRequest {
  return {
    address: ""
  };
}

export const AffiliateInfoRequest = {
  encode(message: AffiliateInfoRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AffiliateInfoRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAffiliateInfoRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<AffiliateInfoRequest>): AffiliateInfoRequest {
    const message = createBaseAffiliateInfoRequest();
    message.address = object.address ?? "";
    return message;
  }

};

function createBaseAffiliateInfoResponse(): AffiliateInfoResponse {
  return {
    tier: 0,
    feeSharePpm: 0,
    referredVolume: new Uint8Array(),
    stakedAmount: new Uint8Array()
  };
}

export const AffiliateInfoResponse = {
  encode(message: AffiliateInfoResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.tier !== 0) {
      writer.uint32(8).uint32(message.tier);
    }

    if (message.feeSharePpm !== 0) {
      writer.uint32(16).uint32(message.feeSharePpm);
    }

    if (message.referredVolume.length !== 0) {
      writer.uint32(26).bytes(message.referredVolume);
    }

    if (message.stakedAmount.length !== 0) {
      writer.uint32(34).bytes(message.stakedAmount);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AffiliateInfoResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAffiliateInfoResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.tier = reader.uint32();
          break;

        case 2:
          message.feeSharePpm = reader.uint32();
          break;

        case 3:
          message.referredVolume = reader.bytes();
          break;

        case 4:
          message.stakedAmount = reader.bytes();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<AffiliateInfoResponse>): AffiliateInfoResponse {
    const message = createBaseAffiliateInfoResponse();
    message.tier = object.tier ?? 0;
    message.feeSharePpm = object.feeSharePpm ?? 0;
    message.referredVolume = object.referredVolume ?? new Uint8Array();
    message.stakedAmount = object.stakedAmount ?? new Uint8Array();
    return message;
  }

};

function createBaseReferredByRequest(): ReferredByRequest {
  return {
    address: ""
  };
}

export const ReferredByRequest = {
  encode(message: ReferredByRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReferredByRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReferredByRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<ReferredByRequest>): ReferredByRequest {
    const message = createBaseReferredByRequest();
    message.address = object.address ?? "";
    return message;
  }

};

function createBaseReferredByResponse(): ReferredByResponse {
  return {
    affiliateAddress: ""
  };
}

export const ReferredByResponse = {
  encode(message: ReferredByResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.affiliateAddress !== "") {
      writer.uint32(10).string(message.affiliateAddress);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReferredByResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReferredByResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.affiliateAddress = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<ReferredByResponse>): ReferredByResponse {
    const message = createBaseReferredByResponse();
    message.affiliateAddress = object.affiliateAddress ?? "";
    return message;
  }

};

function createBaseAllAffiliateTiersRequest(): AllAffiliateTiersRequest {
  return {};
}

export const AllAffiliateTiersRequest = {
  encode(_: AllAffiliateTiersRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AllAffiliateTiersRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAllAffiliateTiersRequest();

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

  fromPartial(_: DeepPartial<AllAffiliateTiersRequest>): AllAffiliateTiersRequest {
    const message = createBaseAllAffiliateTiersRequest();
    return message;
  }

};

function createBaseAllAffiliateTiersResponse(): AllAffiliateTiersResponse {
  return {
    tiers: undefined
  };
}

export const AllAffiliateTiersResponse = {
  encode(message: AllAffiliateTiersResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.tiers !== undefined) {
      AffiliateTiers.encode(message.tiers, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AllAffiliateTiersResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAllAffiliateTiersResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.tiers = AffiliateTiers.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<AllAffiliateTiersResponse>): AllAffiliateTiersResponse {
    const message = createBaseAllAffiliateTiersResponse();
    message.tiers = object.tiers !== undefined && object.tiers !== null ? AffiliateTiers.fromPartial(object.tiers) : undefined;
    return message;
  }

};