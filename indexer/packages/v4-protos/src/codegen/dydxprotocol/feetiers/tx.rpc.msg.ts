import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgUpdatePerpetualFeeParams, MsgUpdatePerpetualFeeParamsResponse, MsgSetFeeDiscountCampaignParams, MsgSetFeeDiscountCampaignParamsResponse } from "./tx";
/** Msg defines the Msg service. */

export interface Msg {
  /** UpdatePerpetualFeeParams updates the PerpetualFeeParams in state. */
  updatePerpetualFeeParams(request: MsgUpdatePerpetualFeeParams): Promise<MsgUpdatePerpetualFeeParamsResponse>;
  /**
   * SetFeeDiscountCampaignParams sets or updates fee discount campaigns for
   * specific CLOB pairs.
   */

  setFeeDiscountCampaignParams(request: MsgSetFeeDiscountCampaignParams): Promise<MsgSetFeeDiscountCampaignParamsResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.updatePerpetualFeeParams = this.updatePerpetualFeeParams.bind(this);
    this.setFeeDiscountCampaignParams = this.setFeeDiscountCampaignParams.bind(this);
  }

  updatePerpetualFeeParams(request: MsgUpdatePerpetualFeeParams): Promise<MsgUpdatePerpetualFeeParamsResponse> {
    const data = MsgUpdatePerpetualFeeParams.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.feetiers.Msg", "UpdatePerpetualFeeParams", data);
    return promise.then(data => MsgUpdatePerpetualFeeParamsResponse.decode(new _m0.Reader(data)));
  }

  setFeeDiscountCampaignParams(request: MsgSetFeeDiscountCampaignParams): Promise<MsgSetFeeDiscountCampaignParamsResponse> {
    const data = MsgSetFeeDiscountCampaignParams.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.feetiers.Msg", "SetFeeDiscountCampaignParams", data);
    return promise.then(data => MsgSetFeeDiscountCampaignParamsResponse.decode(new _m0.Reader(data)));
  }

}