import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgAddAuthenticator, MsgAddAuthenticatorResponse, MsgRemoveAuthenticator, MsgRemoveAuthenticatorResponse, MsgSetActiveState, MsgSetActiveStateResponse } from "./tx";
/** Msg defines the Msg service. */

export interface Msg {
  /** AddAuthenticator adds an authenticator to an account. */
  addAuthenticator(request: MsgAddAuthenticator): Promise<MsgAddAuthenticatorResponse>;
  /** RemoveAuthenticator removes an authenticator from an account. */

  removeAuthenticator(request: MsgRemoveAuthenticator): Promise<MsgRemoveAuthenticatorResponse>;
  /**
   * SetActiveState sets the active state of the authenticator.
   * Primarily used for circuit breaking.
   */

  setActiveState(request: MsgSetActiveState): Promise<MsgSetActiveStateResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.addAuthenticator = this.addAuthenticator.bind(this);
    this.removeAuthenticator = this.removeAuthenticator.bind(this);
    this.setActiveState = this.setActiveState.bind(this);
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

  setActiveState(request: MsgSetActiveState): Promise<MsgSetActiveStateResponse> {
    const data = MsgSetActiveState.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.accountplus.Msg", "SetActiveState", data);
    return promise.then(data => MsgSetActiveStateResponse.decode(new _m0.Reader(data)));
  }

}