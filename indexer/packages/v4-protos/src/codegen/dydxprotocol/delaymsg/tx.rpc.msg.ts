import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgDelayMessage, MsgDelayMessageResponse } from "./tx";
/** Msg defines the Msg service. */

export interface Msg {
  /**
   * DelayMessage delays the execution of a message for a given number of
   * blocks.
   */
  delayMessage(request: MsgDelayMessage): Promise<MsgDelayMessageResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.delayMessage = this.delayMessage.bind(this);
  }

  delayMessage(request: MsgDelayMessage): Promise<MsgDelayMessageResponse> {
    const data = MsgDelayMessage.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.delaymsg.Msg", "DelayMessage", data);
    return promise.then(data => MsgDelayMessageResponse.decode(new _m0.Reader(data)));
  }

}