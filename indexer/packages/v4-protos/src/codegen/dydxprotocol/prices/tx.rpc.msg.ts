import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgCreateOracleMarket, MsgCreateOracleMarketResponse, MsgUpdateMarketParam, MsgUpdateMarketParamResponse } from "./tx";
/** Msg defines the Msg service. */

export interface Msg {
  /** CreateOracleMarket creates a new oracle market. */
  createOracleMarket(request: MsgCreateOracleMarket): Promise<MsgCreateOracleMarketResponse>;
  /**
   * UpdateMarketParams allows governance to update the parameters of an
   * oracle market.
   */

  updateMarketParam(request: MsgUpdateMarketParam): Promise<MsgUpdateMarketParamResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.createOracleMarket = this.createOracleMarket.bind(this);
    this.updateMarketParam = this.updateMarketParam.bind(this);
  }

  createOracleMarket(request: MsgCreateOracleMarket): Promise<MsgCreateOracleMarketResponse> {
    const data = MsgCreateOracleMarket.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.prices.Msg", "CreateOracleMarket", data);
    return promise.then(data => MsgCreateOracleMarketResponse.decode(new _m0.Reader(data)));
  }

  updateMarketParam(request: MsgUpdateMarketParam): Promise<MsgUpdateMarketParamResponse> {
    const data = MsgUpdateMarketParam.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.prices.Msg", "UpdateMarketParam", data);
    return promise.then(data => MsgUpdateMarketParamResponse.decode(new _m0.Reader(data)));
  }

}