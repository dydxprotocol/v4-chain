import { LCDClient } from "@osmonauts/lcd";
import { GetAuthenticatorRequest, GetAuthenticatorResponseSDKType, GetAuthenticatorsRequest, GetAuthenticatorsResponseSDKType } from "./query";
export class LCDQueryClient {
  req: LCDClient;

  constructor({
    requestClient
  }: {
    requestClient: LCDClient;
  }) {
    this.req = requestClient;
    this.getAuthenticator = this.getAuthenticator.bind(this);
    this.getAuthenticators = this.getAuthenticators.bind(this);
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

}