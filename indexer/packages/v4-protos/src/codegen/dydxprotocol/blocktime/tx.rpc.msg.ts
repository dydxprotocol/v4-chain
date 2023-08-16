import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgUpdateDowntimeParams, MsgUpdateDowntimeParamsResponse, MsgIsDelayedBlock, MsgIsDelayedBlockResponse } from "./tx";
/** Msg defines the Msg service. */

export interface Msg {
  /** UpdateDowntimeParams updates the DowntimeParams in state. */
  updateDowntimeParams(request: MsgUpdateDowntimeParams): Promise<MsgUpdateDowntimeParamsResponse>;
  /**
   * IsDelayedBlock indicates a significant difference between wall time and the
   * time of the proposed block.
   */

  isDelayedBlock(request: MsgIsDelayedBlock): Promise<MsgIsDelayedBlockResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.updateDowntimeParams = this.updateDowntimeParams.bind(this);
    this.isDelayedBlock = this.isDelayedBlock.bind(this);
  }

  updateDowntimeParams(request: MsgUpdateDowntimeParams): Promise<MsgUpdateDowntimeParamsResponse> {
    const data = MsgUpdateDowntimeParams.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.blocktime.Msg", "UpdateDowntimeParams", data);
    return promise.then(data => MsgUpdateDowntimeParamsResponse.decode(new _m0.Reader(data)));
  }

  isDelayedBlock(request: MsgIsDelayedBlock): Promise<MsgIsDelayedBlockResponse> {
    const data = MsgIsDelayedBlock.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.blocktime.Msg", "IsDelayedBlock", data);
    return promise.then(data => MsgIsDelayedBlockResponse.decode(new _m0.Reader(data)));
  }

}