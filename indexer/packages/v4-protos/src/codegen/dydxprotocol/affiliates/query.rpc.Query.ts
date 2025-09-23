import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { AffiliateInfoRequest, AffiliateInfoResponse, ReferredByRequest, ReferredByResponse, AllAffiliateTiersRequest, AllAffiliateTiersResponse, AffiliateWhitelistRequest, AffiliateWhitelistResponse, AffiliateParametersRequest, AffiliateParametersResponse, AffiliateOverridesRequest, AffiliateOverridesResponse } from "./query";
/** Query defines the gRPC querier service. */

export interface Query {
  /** Query AffiliateInfo returns the affiliate info for a given address. */
  affiliateInfo(request: AffiliateInfoRequest): Promise<AffiliateInfoResponse>;
  /** Query ReferredBy returns the affiliate that referred a given address. */

  referredBy(request: ReferredByRequest): Promise<ReferredByResponse>;
  /** Query AllAffiliateTiers returns all affiliate tiers. */

  allAffiliateTiers(request?: AllAffiliateTiersRequest): Promise<AllAffiliateTiersResponse>;
  /** Query AffiliateWhitelist returns the affiliate whitelist. */

  affiliateWhitelist(request?: AffiliateWhitelistRequest): Promise<AffiliateWhitelistResponse>;
  /** Query AffiliateParameters returns the affiliate parameters. */

  affiliateParameters(request?: AffiliateParametersRequest): Promise<AffiliateParametersResponse>;
  /** Query AffiliateOverrides returns the affiliate overrides. */

  affiliateOverrides(request?: AffiliateOverridesRequest): Promise<AffiliateOverridesResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.affiliateInfo = this.affiliateInfo.bind(this);
    this.referredBy = this.referredBy.bind(this);
    this.allAffiliateTiers = this.allAffiliateTiers.bind(this);
    this.affiliateWhitelist = this.affiliateWhitelist.bind(this);
    this.affiliateParameters = this.affiliateParameters.bind(this);
    this.affiliateOverrides = this.affiliateOverrides.bind(this);
  }

  affiliateInfo(request: AffiliateInfoRequest): Promise<AffiliateInfoResponse> {
    const data = AffiliateInfoRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.affiliates.Query", "AffiliateInfo", data);
    return promise.then(data => AffiliateInfoResponse.decode(new _m0.Reader(data)));
  }

  referredBy(request: ReferredByRequest): Promise<ReferredByResponse> {
    const data = ReferredByRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.affiliates.Query", "ReferredBy", data);
    return promise.then(data => ReferredByResponse.decode(new _m0.Reader(data)));
  }

  allAffiliateTiers(request: AllAffiliateTiersRequest = {}): Promise<AllAffiliateTiersResponse> {
    const data = AllAffiliateTiersRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.affiliates.Query", "AllAffiliateTiers", data);
    return promise.then(data => AllAffiliateTiersResponse.decode(new _m0.Reader(data)));
  }

  affiliateWhitelist(request: AffiliateWhitelistRequest = {}): Promise<AffiliateWhitelistResponse> {
    const data = AffiliateWhitelistRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.affiliates.Query", "AffiliateWhitelist", data);
    return promise.then(data => AffiliateWhitelistResponse.decode(new _m0.Reader(data)));
  }

  affiliateParameters(request: AffiliateParametersRequest = {}): Promise<AffiliateParametersResponse> {
    const data = AffiliateParametersRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.affiliates.Query", "AffiliateParameters", data);
    return promise.then(data => AffiliateParametersResponse.decode(new _m0.Reader(data)));
  }

  affiliateOverrides(request: AffiliateOverridesRequest = {}): Promise<AffiliateOverridesResponse> {
    const data = AffiliateOverridesRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.affiliates.Query", "AffiliateOverrides", data);
    return promise.then(data => AffiliateOverridesResponse.decode(new _m0.Reader(data)));
  }

}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    affiliateInfo(request: AffiliateInfoRequest): Promise<AffiliateInfoResponse> {
      return queryService.affiliateInfo(request);
    },

    referredBy(request: ReferredByRequest): Promise<ReferredByResponse> {
      return queryService.referredBy(request);
    },

    allAffiliateTiers(request?: AllAffiliateTiersRequest): Promise<AllAffiliateTiersResponse> {
      return queryService.allAffiliateTiers(request);
    },

    affiliateWhitelist(request?: AffiliateWhitelistRequest): Promise<AffiliateWhitelistResponse> {
      return queryService.affiliateWhitelist(request);
    },

    affiliateParameters(request?: AffiliateParametersRequest): Promise<AffiliateParametersResponse> {
      return queryService.affiliateParameters(request);
    },

    affiliateOverrides(request?: AffiliateOverridesRequest): Promise<AffiliateOverridesResponse> {
      return queryService.affiliateOverrides(request);
    }

  };
};