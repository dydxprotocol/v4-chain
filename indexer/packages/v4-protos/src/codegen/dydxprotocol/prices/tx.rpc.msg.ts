import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgUpdateMarketPrices, MsgUpdateMarketPricesResponse, MsgCreateOracleMarket, MsgCreateOracleMarketResponse, MsgUpdateMarketParam, MsgUpdateMarketParamResponse } from "./tx";
/** Msg defines the Msg service. */

export interface Msg {
  /**
   * UpdateMarketPrices updates the oracle price of a market relative to
   * quoteCurrency.
   */
  updateMarketPrices(request: MsgUpdateMarketPrices): Promise<MsgUpdateMarketPricesResponse>;
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
    this.updateMarketPrices = this.updateMarketPrices.bind(this);
    this.createOracleMarket = this.createOracleMarket.bind(this);
    this.updateMarketParam = this.updateMarketParam.bind(this);
  }

  updateMarketPrices(request: MsgUpdateMarketPrices): Promise<MsgUpdateMarketPricesResponse> {
    const data = MsgUpdateMarketPrices.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.prices.Msg", "UpdateMarketPrices", data);
    return promise.then(data => MsgUpdateMarketPricesResponse.decode(new _m0.Reader(data)));
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