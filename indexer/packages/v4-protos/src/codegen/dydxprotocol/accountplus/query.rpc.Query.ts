import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { GetAuthenticatorRequest, GetAuthenticatorResponse, GetAuthenticatorsRequest, GetAuthenticatorsResponse } from "./query";
/** Query defines the gRPC querier service. */

export interface Query {
  /** Queries a single authenticator by account and authenticator ID. */
  getAuthenticator(request: GetAuthenticatorRequest): Promise<GetAuthenticatorResponse>;
  /** Queries all authenticators for a given account. */

  getAuthenticators(request: GetAuthenticatorsRequest): Promise<GetAuthenticatorsResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.getAuthenticator = this.getAuthenticator.bind(this);
    this.getAuthenticators = this.getAuthenticators.bind(this);
  }

  getAuthenticator(request: GetAuthenticatorRequest): Promise<GetAuthenticatorResponse> {
    const data = GetAuthenticatorRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.accountplus.Query", "GetAuthenticator", data);
    return promise.then(data => GetAuthenticatorResponse.decode(new _m0.Reader(data)));
  }

  getAuthenticators(request: GetAuthenticatorsRequest): Promise<GetAuthenticatorsResponse> {
    const data = GetAuthenticatorsRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.accountplus.Query", "GetAuthenticators", data);
    return promise.then(data => GetAuthenticatorsResponse.decode(new _m0.Reader(data)));
  }

}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    getAuthenticator(request: GetAuthenticatorRequest): Promise<GetAuthenticatorResponse> {
      return queryService.getAuthenticator(request);
    },

    getAuthenticators(request: GetAuthenticatorsRequest): Promise<GetAuthenticatorsResponse> {
      return queryService.getAuthenticators(request);
    }

  };
};