import { Rpc } from "../../helpers";
import { BinaryReader } from "../../binary";
import { MsgSetLimitParams, MsgSetLimitParamsResponse, MsgDeleteLimitParams, MsgDeleteLimitParamsResponse } from "./tx";
/** Msg defines the Msg service. */
export interface Msg {
  /** SetLimitParams sets a `LimitParams` object in state. */
  setLimitParams(request: MsgSetLimitParams): Promise<MsgSetLimitParamsResponse>;
  /** DeleteLimitParams removes a `LimitParams` object from state. */
  deleteLimitParams(request: MsgDeleteLimitParams): Promise<MsgDeleteLimitParamsResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.setLimitParams = this.setLimitParams.bind(this);
    this.deleteLimitParams = this.deleteLimitParams.bind(this);
  }
  setLimitParams(request: MsgSetLimitParams): Promise<MsgSetLimitParamsResponse> {
    const data = MsgSetLimitParams.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.ibcratelimit.Msg", "SetLimitParams", data);
    return promise.then(data => MsgSetLimitParamsResponse.decode(new BinaryReader(data)));
  }
  deleteLimitParams(request: MsgDeleteLimitParams): Promise<MsgDeleteLimitParamsResponse> {
    const data = MsgDeleteLimitParams.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.ibcratelimit.Msg", "DeleteLimitParams", data);
    return promise.then(data => MsgDeleteLimitParamsResponse.decode(new BinaryReader(data)));
  }
}