import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgAddAuthenticator, MsgAddAuthenticatorResponse, MsgRemoveAuthenticator, MsgRemoveAuthenticatorResponse } from "./tx";
/** Msg defines the Msg service. */

export interface Msg {
  /** AddAuthenticator adds an authenticator to an account. */
  addAuthenticator(request: MsgAddAuthenticator): Promise<MsgAddAuthenticatorResponse>;
  /** RemoveAuthenticator removes an authenticator from an account. */

  removeAuthenticator(request: MsgRemoveAuthenticator): Promise<MsgRemoveAuthenticatorResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.addAuthenticator = this.addAuthenticator.bind(this);
    this.removeAuthenticator = this.removeAuthenticator.bind(this);
  }

  addAuthenticator(request: MsgAddAuthenticator): Promise<MsgAddAuthenticatorResponse> {
    const data = MsgAddAuthenticator.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.accountplus.Msg", "AddAuthenticator", data);
    return promise.then(data => MsgAddAuthenticatorResponse.decode(new _m0.Reader(data)));
  }

  removeAuthenticator(request: MsgRemoveAuthenticator): Promise<MsgRemoveAuthenticatorResponse> {
    const data = MsgRemoveAuthenticator.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.accountplus.Msg", "RemoveAuthenticator", data);
    return promise.then(data => MsgRemoveAuthenticatorResponse.decode(new _m0.Reader(data)));
  }

}