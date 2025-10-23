import { LCDClient } from "@osmonauts/lcd";
import { AffiliateInfoRequest, AffiliateInfoResponseSDKType, ReferredByRequest, ReferredByResponseSDKType, AllAffiliateTiersRequest, AllAffiliateTiersResponseSDKType, AffiliateWhitelistRequest, AffiliateWhitelistResponseSDKType, AffiliateParametersRequest, AffiliateParametersResponseSDKType, AffiliateOverridesRequest, AffiliateOverridesResponseSDKType } from "./query";
export class LCDQueryClient {
  req: LCDClient;

  constructor({
    requestClient
  }: {
    requestClient: LCDClient;
  }) {
    this.req = requestClient;
    this.affiliateInfo = this.affiliateInfo.bind(this);
    this.referredBy = this.referredBy.bind(this);
    this.allAffiliateTiers = this.allAffiliateTiers.bind(this);
    this.affiliateWhitelist = this.affiliateWhitelist.bind(this);
    this.affiliateParameters = this.affiliateParameters.bind(this);
    this.affiliateOverrides = this.affiliateOverrides.bind(this);
  }
  /* Query AffiliateInfo returns the affiliate info for a given address. */


  async affiliateInfo(params: AffiliateInfoRequest): Promise<AffiliateInfoResponseSDKType> {
    const endpoint = `dydxprotocol/affiliates/affiliate_info/${params.address}`;
    return await this.req.get<AffiliateInfoResponseSDKType>(endpoint);
  }
  /* Query ReferredBy returns the affiliate that referred a given address. */


  async referredBy(params: ReferredByRequest): Promise<ReferredByResponseSDKType> {
    const endpoint = `dydxprotocol/affiliates/referred_by/${params.address}`;
    return await this.req.get<ReferredByResponseSDKType>(endpoint);
  }
  /* Query AllAffiliateTiers returns all affiliate tiers. */


  async allAffiliateTiers(_params: AllAffiliateTiersRequest = {}): Promise<AllAffiliateTiersResponseSDKType> {
    const endpoint = `dydxprotocol/affiliates/all_affiliate_tiers`;
    return await this.req.get<AllAffiliateTiersResponseSDKType>(endpoint);
  }
  /* Query AffiliateWhitelist returns the affiliate whitelist. */


  async affiliateWhitelist(_params: AffiliateWhitelistRequest = {}): Promise<AffiliateWhitelistResponseSDKType> {
    const endpoint = `dydxprotocol/affiliates/affiliate_whitelist`;
    return await this.req.get<AffiliateWhitelistResponseSDKType>(endpoint);
  }
  /* Query AffiliateParameters returns the affiliate parameters. */


  async affiliateParameters(_params: AffiliateParametersRequest = {}): Promise<AffiliateParametersResponseSDKType> {
    const endpoint = `dydxprotocol/affiliates/affiliate_parameters`;
    return await this.req.get<AffiliateParametersResponseSDKType>(endpoint);
  }
  /* Query AffiliateOverrides returns the affiliate overrides. */


  async affiliateOverrides(_params: AffiliateOverridesRequest = {}): Promise<AffiliateOverridesResponseSDKType> {
    const endpoint = `dydxprotocol/affiliates/affiliate_overrides`;
    return await this.req.get<AffiliateOverridesResponseSDKType>(endpoint);
  }

}