import { LCDClient } from "@osmonauts/lcd";
import { QueryParamsRequest, QueryParamsResponseSDKType, GetAuthenticatorRequest, GetAuthenticatorResponseSDKType, GetAuthenticatorsRequest, GetAuthenticatorsResponseSDKType, AccountStateRequest, AccountStateResponseSDKType } from "./query";
export class LCDQueryClient {
  req: LCDClient;

  constructor({
    requestClient
  }: {
    requestClient: LCDClient;
  }) {
    this.req = requestClient;
    this.params = this.params.bind(this);
    this.getAuthenticator = this.getAuthenticator.bind(this);
    this.getAuthenticators = this.getAuthenticators.bind(this);
    this.accountState = this.accountState.bind(this);
  }
  /* Parameters queries the parameters of the module. */


  async params(_params: QueryParamsRequest = {}): Promise<QueryParamsResponseSDKType> {
    const endpoint = `dydxprotocol/accountplus/params`;
    return await this.req.get<QueryParamsResponseSDKType>(endpoint);
  }
  /* Queries a single authenticator by account and authenticator ID. */


  async getAuthenticator(params: GetAuthenticatorRequest): Promise<GetAuthenticatorResponseSDKType> {
    const endpoint = `dydxprotocol/accountplus/authenticator/${params.account}/${params.authenticatorId}`;
    return await this.req.get<GetAuthenticatorResponseSDKType>(endpoint);
  }
  /* Queries all authenticators for a given account. */


  async getAuthenticators(params: GetAuthenticatorsRequest): Promise<GetAuthenticatorsResponseSDKType> {
    const endpoint = `dydxprotocol/accountplus/authenticators/${params.account}`;
    return await this.req.get<GetAuthenticatorsResponseSDKType>(endpoint);
  }
  /* Queries for an account state (timestamp nonce). */


  async accountState(params: AccountStateRequest): Promise<AccountStateResponseSDKType> {
    const endpoint = `dydxprotocol/accountplus/account_state/${params.address}`;
    return await this.req.get<AccountStateResponseSDKType>(endpoint);
  }

}