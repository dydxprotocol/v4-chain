import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgSetLimitParams, MsgSetLimitParamsResponse, MsgUpdateSDAIConversionRate, MsgUpdateSDAIConversionRateResponse } from "./tx";
/** Msg defines the Msg service. */

export interface Msg {
  /** SetLimitParams sets a `LimitParams` object in state. */
  setLimitParams(request: MsgSetLimitParams): Promise<MsgSetLimitParamsResponse>;
  /** UpdateSDAIConversionRate updates the sDAI conversion rate with a asociated block height */

  updateSDAIConversionRate(request: MsgUpdateSDAIConversionRate): Promise<MsgUpdateSDAIConversionRateResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.setLimitParams = this.setLimitParams.bind(this);
    this.updateSDAIConversionRate = this.updateSDAIConversionRate.bind(this);
  }

  setLimitParams(request: MsgSetLimitParams): Promise<MsgSetLimitParamsResponse> {
    const data = MsgSetLimitParams.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.ratelimit.Msg", "SetLimitParams", data);
    return promise.then(data => MsgSetLimitParamsResponse.decode(new _m0.Reader(data)));
  }

  updateSDAIConversionRate(request: MsgUpdateSDAIConversionRate): Promise<MsgUpdateSDAIConversionRateResponse> {
    const data = MsgUpdateSDAIConversionRate.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.ratelimit.Msg", "UpdateSDAIConversionRate", data);
    return promise.then(data => MsgUpdateSDAIConversionRateResponse.decode(new _m0.Reader(data)));
  }

}