import { Rpc } from "../../helpers";
import { BinaryReader } from "../../binary";
import { MsgUpdateDowntimeParams, MsgUpdateDowntimeParamsResponse } from "./tx";
/** Msg defines the Msg service. */
export interface Msg {
  /** UpdateDowntimeParams updates the DowntimeParams in state. */
  updateDowntimeParams(request: MsgUpdateDowntimeParams): Promise<MsgUpdateDowntimeParamsResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.updateDowntimeParams = this.updateDowntimeParams.bind(this);
  }
  updateDowntimeParams(request: MsgUpdateDowntimeParams): Promise<MsgUpdateDowntimeParamsResponse> {
    const data = MsgUpdateDowntimeParams.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.blocktime.Msg", "UpdateDowntimeParams", data);
    return promise.then(data => MsgUpdateDowntimeParamsResponse.decode(new BinaryReader(data)));
  }
}