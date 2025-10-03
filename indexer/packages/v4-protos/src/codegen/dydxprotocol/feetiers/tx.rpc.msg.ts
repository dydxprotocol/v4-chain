import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgUpdatePerpetualFeeParams, MsgUpdatePerpetualFeeParamsResponse, MsgSetFeeHolidayParams, MsgSetFeeHolidayParamsResponse } from "./tx";
/** Msg defines the Msg service. */

export interface Msg {
  /** UpdatePerpetualFeeParams updates the PerpetualFeeParams in state. */
  updatePerpetualFeeParams(request: MsgUpdatePerpetualFeeParams): Promise<MsgUpdatePerpetualFeeParamsResponse>;
  /** SetFeeHolidayParams sets the no fee holiday period for each CLOB pair. */

  setFeeHolidayParams(request: MsgSetFeeHolidayParams): Promise<MsgSetFeeHolidayParamsResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.updatePerpetualFeeParams = this.updatePerpetualFeeParams.bind(this);
    this.setFeeHolidayParams = this.setFeeHolidayParams.bind(this);
  }

  updatePerpetualFeeParams(request: MsgUpdatePerpetualFeeParams): Promise<MsgUpdatePerpetualFeeParamsResponse> {
    const data = MsgUpdatePerpetualFeeParams.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.feetiers.Msg", "UpdatePerpetualFeeParams", data);
    return promise.then(data => MsgUpdatePerpetualFeeParamsResponse.decode(new _m0.Reader(data)));
  }

  setFeeHolidayParams(request: MsgSetFeeHolidayParams): Promise<MsgSetFeeHolidayParamsResponse> {
    const data = MsgSetFeeHolidayParams.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.feetiers.Msg", "SetFeeHolidayParams", data);
    return promise.then(data => MsgSetFeeHolidayParamsResponse.decode(new _m0.Reader(data)));
  }

}