import { LCDClient } from "@osmonauts/lcd";
import { AffiliateInfoRequest, AffiliateInfoResponseSDKType, ReferredByRequest, ReferredByResponseSDKType, AllAffiliateTiersRequest, AllAffiliateTiersResponseSDKType, AffiliateWhitelistRequest, AffiliateWhitelistResponseSDKType, BrokerAffiliateRequest, BrokerAffiliateResponseSDKType } from "./query";
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
    this.brokerAffiliate = this.brokerAffiliate.bind(this);
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
  /* Query BrokerAffiliate returns the broker affiliate for a given broker ID. */


  async brokerAffiliate(params: BrokerAffiliateRequest): Promise<BrokerAffiliateResponseSDKType> {
    const endpoint = `dydxprotocol/affiliates/broker_affiliate/${params.brokerId}`;
    return await this.req.get<BrokerAffiliateResponseSDKType>(endpoint);
  }

}