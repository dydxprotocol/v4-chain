import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgUpdatePerpetualFeeParams, MsgUpdatePerpetualFeeParamsResponse, MsgSetMarketFeeDiscountParams, MsgSetMarketFeeDiscountParamsResponse, MsgSetStakingTiers, MsgSetStakingTiersResponse } from "./tx";
/** Msg defines the Msg service. */

export interface Msg {
  /** UpdatePerpetualFeeParams updates the PerpetualFeeParams in state. */
  updatePerpetualFeeParams(request: MsgUpdatePerpetualFeeParams): Promise<MsgUpdatePerpetualFeeParamsResponse>;
  /**
   * SetMarketFeeDiscountParams sets or updates PerMarketFeeDiscountParams for
   * specific CLOB pairs.
   */

  setMarketFeeDiscountParams(request: MsgSetMarketFeeDiscountParams): Promise<MsgSetMarketFeeDiscountParamsResponse>;
  /** SetStakingTiers sets the staking tiers in state. */

  setStakingTiers(request: MsgSetStakingTiers): Promise<MsgSetStakingTiersResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.updatePerpetualFeeParams = this.updatePerpetualFeeParams.bind(this);
    this.setMarketFeeDiscountParams = this.setMarketFeeDiscountParams.bind(this);
    this.setStakingTiers = this.setStakingTiers.bind(this);
  }

  updatePerpetualFeeParams(request: MsgUpdatePerpetualFeeParams): Promise<MsgUpdatePerpetualFeeParamsResponse> {
    const data = MsgUpdatePerpetualFeeParams.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.feetiers.Msg", "UpdatePerpetualFeeParams", data);
    return promise.then(data => MsgUpdatePerpetualFeeParamsResponse.decode(new _m0.Reader(data)));
  }

  setMarketFeeDiscountParams(request: MsgSetMarketFeeDiscountParams): Promise<MsgSetMarketFeeDiscountParamsResponse> {
    const data = MsgSetMarketFeeDiscountParams.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.feetiers.Msg", "SetMarketFeeDiscountParams", data);
    return promise.then(data => MsgSetMarketFeeDiscountParamsResponse.decode(new _m0.Reader(data)));
  }

  setStakingTiers(request: MsgSetStakingTiers): Promise<MsgSetStakingTiersResponse> {
    const data = MsgSetStakingTiers.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.feetiers.Msg", "SetStakingTiers", data);
    return promise.then(data => MsgSetStakingTiersResponse.decode(new _m0.Reader(data)));
  }

}