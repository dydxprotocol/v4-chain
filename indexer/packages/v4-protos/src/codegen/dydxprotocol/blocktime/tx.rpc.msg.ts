import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgUpdateDowntimeParams, MsgUpdateDowntimeParamsResponse, MsgUpdateSynchronyParams, MsgUpdateSynchronyParamsResponse } from "./tx";
/** Msg defines the Msg service. */

export interface Msg {
  /** UpdateDowntimeParams updates the DowntimeParams in state. */
  updateDowntimeParams(request: MsgUpdateDowntimeParams): Promise<MsgUpdateDowntimeParamsResponse>;
  /** UpdateSynchronyParams updates the SynchronyParams in state. */

  updateSynchronyParams(request: MsgUpdateSynchronyParams): Promise<MsgUpdateSynchronyParamsResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.updateDowntimeParams = this.updateDowntimeParams.bind(this);
    this.updateSynchronyParams = this.updateSynchronyParams.bind(this);
  }

  updateDowntimeParams(request: MsgUpdateDowntimeParams): Promise<MsgUpdateDowntimeParamsResponse> {
    const data = MsgUpdateDowntimeParams.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.blocktime.Msg", "UpdateDowntimeParams", data);
    return promise.then(data => MsgUpdateDowntimeParamsResponse.decode(new _m0.Reader(data)));
  }

  updateSynchronyParams(request: MsgUpdateSynchronyParams): Promise<MsgUpdateSynchronyParamsResponse> {
    const data = MsgUpdateSynchronyParams.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.blocktime.Msg", "UpdateSynchronyParams", data);
    return promise.then(data => MsgUpdateSynchronyParamsResponse.decode(new _m0.Reader(data)));
  }

}