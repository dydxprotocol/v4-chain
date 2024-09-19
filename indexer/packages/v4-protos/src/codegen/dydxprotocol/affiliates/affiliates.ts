import * as _m0 from "protobufjs/minimal";
import { DeepPartial, Long } from "../../helpers";
/** AffiliateTiers defines the affiliate tiers. */

export interface AffiliateTiers {
  /** All affiliate tiers */
  tiers: AffiliateTiers_Tier[];
}
/** AffiliateTiers defines the affiliate tiers. */

export interface AffiliateTiersSDKType {
  /** All affiliate tiers */
  tiers: AffiliateTiers_TierSDKType[];
}
/** Tier defines an affiliate tier. */

export interface AffiliateTiers_Tier {
  /** Required all-time referred volume in quote quantums. */
  reqReferredVolumeQuoteQuantums: Long;
  /** Required currently staked native tokens (in whole coins). */

  reqStakedWholeCoins: number;
  /** Taker fee share in parts-per-million. */

  takerFeeSharePpm: number;
}
/** Tier defines an affiliate tier. */

export interface AffiliateTiers_TierSDKType {
  /** Required all-time referred volume in quote quantums. */
  req_referred_volume_quote_quantums: Long;
  /** Required currently staked native tokens (in whole coins). */

  req_staked_whole_coins: number;
  /** Taker fee share in parts-per-million. */

  taker_fee_share_ppm: number;
}
/**
 * AffiliateWhitelist specifies the whitelisted affiliates.
 * If an address is in the whitelist, then the affiliate fee share in
 * this object will override fee share from the regular affiliate tiers above.
 */

export interface AffiliateWhitelist {
  /** All affiliate whitelist tiers. */
  tiers: AffiliateWhitelist_Tier[];
}
/**
 * AffiliateWhitelist specifies the whitelisted affiliates.
 * If an address is in the whitelist, then the affiliate fee share in
 * this object will override fee share from the regular affiliate tiers above.
 */

export interface AffiliateWhitelistSDKType {
  /** All affiliate whitelist tiers. */
  tiers: AffiliateWhitelist_TierSDKType[];
}
/** Tier defines an affiliate whitelist tier. */

export interface AffiliateWhitelist_Tier {
  /** List of unique whitelisted addresses. */
  addresses: string[];
  /** Taker fee share in parts-per-million. */

  takerFeeSharePpm: number;
}
/** Tier defines an affiliate whitelist tier. */

export interface AffiliateWhitelist_TierSDKType {
  /** List of unique whitelisted addresses. */
  addresses: string[];
  /** Taker fee share in parts-per-million. */

  taker_fee_share_ppm: number;
}

function createBaseAffiliateTiers(): AffiliateTiers {
  return {
    tiers: []
  };
}

export const AffiliateTiers = {
  encode(message: AffiliateTiers, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.tiers) {
      AffiliateTiers_Tier.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AffiliateTiers {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAffiliateTiers();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.tiers.push(AffiliateTiers_Tier.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<AffiliateTiers>): AffiliateTiers {
    const message = createBaseAffiliateTiers();
    message.tiers = object.tiers?.map(e => AffiliateTiers_Tier.fromPartial(e)) || [];
    return message;
  }

};

function createBaseAffiliateTiers_Tier(): AffiliateTiers_Tier {
  return {
    reqReferredVolumeQuoteQuantums: Long.UZERO,
    reqStakedWholeCoins: 0,
    takerFeeSharePpm: 0
  };
}

export const AffiliateTiers_Tier = {
  encode(message: AffiliateTiers_Tier, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (!message.reqReferredVolumeQuoteQuantums.isZero()) {
      writer.uint32(8).uint64(message.reqReferredVolumeQuoteQuantums);
    }

    if (message.reqStakedWholeCoins !== 0) {
      writer.uint32(16).uint32(message.reqStakedWholeCoins);
    }

    if (message.takerFeeSharePpm !== 0) {
      writer.uint32(24).uint32(message.takerFeeSharePpm);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AffiliateTiers_Tier {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAffiliateTiers_Tier();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.reqReferredVolumeQuoteQuantums = (reader.uint64() as Long);
          break;

        case 2:
          message.reqStakedWholeCoins = reader.uint32();
          break;

        case 3:
          message.takerFeeSharePpm = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<AffiliateTiers_Tier>): AffiliateTiers_Tier {
    const message = createBaseAffiliateTiers_Tier();
    message.reqReferredVolumeQuoteQuantums = object.reqReferredVolumeQuoteQuantums !== undefined && object.reqReferredVolumeQuoteQuantums !== null ? Long.fromValue(object.reqReferredVolumeQuoteQuantums) : Long.UZERO;
    message.reqStakedWholeCoins = object.reqStakedWholeCoins ?? 0;
    message.takerFeeSharePpm = object.takerFeeSharePpm ?? 0;
    return message;
  }

};

function createBaseAffiliateWhitelist(): AffiliateWhitelist {
  return {
    tiers: []
  };
}

export const AffiliateWhitelist = {
  encode(message: AffiliateWhitelist, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.tiers) {
      AffiliateWhitelist_Tier.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AffiliateWhitelist {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAffiliateWhitelist();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.tiers.push(AffiliateWhitelist_Tier.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<AffiliateWhitelist>): AffiliateWhitelist {
    const message = createBaseAffiliateWhitelist();
    message.tiers = object.tiers?.map(e => AffiliateWhitelist_Tier.fromPartial(e)) || [];
    return message;
  }

};

function createBaseAffiliateWhitelist_Tier(): AffiliateWhitelist_Tier {
  return {
    addresses: [],
    takerFeeSharePpm: 0
  };
}

export const AffiliateWhitelist_Tier = {
  encode(message: AffiliateWhitelist_Tier, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.addresses) {
      writer.uint32(10).string(v!);
    }

    if (message.takerFeeSharePpm !== 0) {
      writer.uint32(16).uint32(message.takerFeeSharePpm);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AffiliateWhitelist_Tier {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAffiliateWhitelist_Tier();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.addresses.push(reader.string());
          break;

        case 2:
          message.takerFeeSharePpm = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<AffiliateWhitelist_Tier>): AffiliateWhitelist_Tier {
    const message = createBaseAffiliateWhitelist_Tier();
    message.addresses = object.addresses?.map(e => e) || [];
    message.takerFeeSharePpm = object.takerFeeSharePpm ?? 0;
    return message;
  }

};