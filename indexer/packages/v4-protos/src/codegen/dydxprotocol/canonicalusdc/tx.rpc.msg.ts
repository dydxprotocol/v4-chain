import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgUpdateControls, MsgUpdateControlsResponse } from "./tx";
/** Msg defines governance actions. Ordinary users keep using
 ibc.applications.transfer.v1.MsgTransfer. */

export interface Msg {
  /** UpdateControls updates canonical routing controls. */
  updateControls(request: MsgUpdateControls): Promise<MsgUpdateControlsResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.updateControls = this.updateControls.bind(this);
  }

  updateControls(request: MsgUpdateControls): Promise<MsgUpdateControlsResponse> {
    const data = MsgUpdateControls.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.canonicalusdc.Msg", "UpdateControls", data);
    return promise.then(data => MsgUpdateControlsResponse.decode(new _m0.Reader(data)));
  }

}