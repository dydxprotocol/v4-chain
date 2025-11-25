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
  /**
   * Required currently staked native tokens (in whole coins).
   * This is deprecated
   */

  /** @deprecated */

  reqStakedWholeCoins: number;
  /** Taker fee share in parts-per-million. */

  takerFeeSharePpm: number;
}
/** Tier defines an affiliate tier. */

export interface AffiliateTiers_TierSDKType {
  /** Required all-time referred volume in quote quantums. */
  req_referred_volume_quote_quantums: Long;
  /**
   * Required currently staked native tokens (in whole coins).
   * This is deprecated
   */

  /** @deprecated */

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
/** AffiliateParameters defines the parameters for the affiliate program. */

export interface AffiliateParameters {
  /**
   * Maximum attributable volume for a referred user in a 30d rolling window in
   * notional
   */
  maximum_30dAttributableVolumePerReferredUserQuoteQuantums: Long;
  /** Referred user automatically gets set to this fee tier */

  refereeMinimumFeeTierIdx: number;
  /**
   * Maximum affiliate revenue for a referred user in a 30d rolling window in
   * quote quantums
   */

  maximum_30dAffiliateRevenuePerReferredUserQuoteQuantums: Long;
}
/** AffiliateParameters defines the parameters for the affiliate program. */

export interface AffiliateParametersSDKType {
  /**
   * Maximum attributable volume for a referred user in a 30d rolling window in
   * notional
   */
  maximum_30d_attributable_volume_per_referred_user_quote_quantums: Long;
  /** Referred user automatically gets set to this fee tier */

  referee_minimum_fee_tier_idx: number;
  /**
   * Maximum affiliate revenue for a referred user in a 30d rolling window in
   * quote quantums
   */

  maximum_30d_affiliate_revenue_per_referred_user_quote_quantums: Long;
}
/** AffiliateOverrides defines the affiliate whitelist. */

export interface AffiliateOverrides {
  /**
   * List of unique whitelisted addresses.
   * These are automatically put at the maximum affiliate tier
   */
  addresses: string[];
}
/** AffiliateOverrides defines the affiliate whitelist. */

export interface AffiliateOverridesSDKType {
  /**
   * List of unique whitelisted addresses.
   * These are automatically put at the maximum affiliate tier
   */
  addresses: string[];
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

function createBaseAffiliateParameters(): AffiliateParameters {
  return {
    maximum_30dAttributableVolumePerReferredUserQuoteQuantums: Long.UZERO,
    refereeMinimumFeeTierIdx: 0,
    maximum_30dAffiliateRevenuePerReferredUserQuoteQuantums: Long.UZERO
  };
}

export const AffiliateParameters = {
  encode(message: AffiliateParameters, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (!message.maximum_30dAttributableVolumePerReferredUserQuoteQuantums.isZero()) {
      writer.uint32(8).uint64(message.maximum_30dAttributableVolumePerReferredUserQuoteQuantums);
    }

    if (message.refereeMinimumFeeTierIdx !== 0) {
      writer.uint32(16).uint32(message.refereeMinimumFeeTierIdx);
    }

    if (!message.maximum_30dAffiliateRevenuePerReferredUserQuoteQuantums.isZero()) {
      writer.uint32(24).uint64(message.maximum_30dAffiliateRevenuePerReferredUserQuoteQuantums);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AffiliateParameters {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAffiliateParameters();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.maximum_30dAttributableVolumePerReferredUserQuoteQuantums = (reader.uint64() as Long);
          break;

        case 2:
          message.refereeMinimumFeeTierIdx = reader.uint32();
          break;

        case 3:
          message.maximum_30dAffiliateRevenuePerReferredUserQuoteQuantums = (reader.uint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<AffiliateParameters>): AffiliateParameters {
    const message = createBaseAffiliateParameters();
    message.maximum_30dAttributableVolumePerReferredUserQuoteQuantums = object.maximum_30dAttributableVolumePerReferredUserQuoteQuantums !== undefined && object.maximum_30dAttributableVolumePerReferredUserQuoteQuantums !== null ? Long.fromValue(object.maximum_30dAttributableVolumePerReferredUserQuoteQuantums) : Long.UZERO;
    message.refereeMinimumFeeTierIdx = object.refereeMinimumFeeTierIdx ?? 0;
    message.maximum_30dAffiliateRevenuePerReferredUserQuoteQuantums = object.maximum_30dAffiliateRevenuePerReferredUserQuoteQuantums !== undefined && object.maximum_30dAffiliateRevenuePerReferredUserQuoteQuantums !== null ? Long.fromValue(object.maximum_30dAffiliateRevenuePerReferredUserQuoteQuantums) : Long.UZERO;
    return message;
  }

};

function createBaseAffiliateOverrides(): AffiliateOverrides {
  return {
    addresses: []
  };
}

export const AffiliateOverrides = {
  encode(message: AffiliateOverrides, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.addresses) {
      writer.uint32(10).string(v!);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AffiliateOverrides {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAffiliateOverrides();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.addresses.push(reader.string());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<AffiliateOverrides>): AffiliateOverrides {
    const message = createBaseAffiliateOverrides();
    message.addresses = object.addresses?.map(e => e) || [];
    return message;
  }

};