import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgSetMarketsHardCap, MsgSetMarketsHardCapResponse } from "./tx";
/** Msg defines the Msg service. */

export interface Msg {
  /** SetMarketsHardCap sets a hard cap on the number of markets listed */
  setMarketsHardCap(request: MsgSetMarketsHardCap): Promise<MsgSetMarketsHardCapResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.setMarketsHardCap = this.setMarketsHardCap.bind(this);
  }

  setMarketsHardCap(request: MsgSetMarketsHardCap): Promise<MsgSetMarketsHardCapResponse> {
    const data = MsgSetMarketsHardCap.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.listing.Msg", "SetMarketsHardCap", data);
    return promise.then(data => MsgSetMarketsHardCapResponse.decode(new _m0.Reader(data)));
  }

}