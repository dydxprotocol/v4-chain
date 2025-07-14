import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgSetMarketMapperRevenueShare, MsgSetMarketMapperRevenueShareResponse, MsgSetMarketMapperRevShareDetailsForMarket, MsgSetMarketMapperRevShareDetailsForMarketResponse, MsgUpdateUnconditionalRevShareConfig, MsgUpdateUnconditionalRevShareConfigResponse, MsgSetOrderRouterRevShare, MsgSetOrderRouterRevShareResponse } from "./tx";
/** Msg defines the Msg service. */

export interface Msg {
  /**
   * SetMarketMapperRevenueShare sets the revenue share for a market
   * mapper.
   */
  setMarketMapperRevenueShare(request: MsgSetMarketMapperRevenueShare): Promise<MsgSetMarketMapperRevenueShareResponse>;
  /**
   * SetMarketMapperRevenueShareDetails sets the revenue share details for a
   * market mapper.
   */

  setMarketMapperRevShareDetailsForMarket(request: MsgSetMarketMapperRevShareDetailsForMarket): Promise<MsgSetMarketMapperRevShareDetailsForMarketResponse>;
  /** UpdateUnconditionalRevShareConfig sets the unconditional revshare config */

  updateUnconditionalRevShareConfig(request: MsgUpdateUnconditionalRevShareConfig): Promise<MsgUpdateUnconditionalRevShareConfigResponse>;
  /** SetOrderRouterRevShare sets the revenue share for an order router. */

  setOrderRouterRevShare(request: MsgSetOrderRouterRevShare): Promise<MsgSetOrderRouterRevShareResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.setMarketMapperRevenueShare = this.setMarketMapperRevenueShare.bind(this);
    this.setMarketMapperRevShareDetailsForMarket = this.setMarketMapperRevShareDetailsForMarket.bind(this);
    this.updateUnconditionalRevShareConfig = this.updateUnconditionalRevShareConfig.bind(this);
    this.setOrderRouterRevShare = this.setOrderRouterRevShare.bind(this);
  }

  setMarketMapperRevenueShare(request: MsgSetMarketMapperRevenueShare): Promise<MsgSetMarketMapperRevenueShareResponse> {
    const data = MsgSetMarketMapperRevenueShare.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.revshare.Msg", "SetMarketMapperRevenueShare", data);
    return promise.then(data => MsgSetMarketMapperRevenueShareResponse.decode(new _m0.Reader(data)));
  }

  setMarketMapperRevShareDetailsForMarket(request: MsgSetMarketMapperRevShareDetailsForMarket): Promise<MsgSetMarketMapperRevShareDetailsForMarketResponse> {
    const data = MsgSetMarketMapperRevShareDetailsForMarket.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.revshare.Msg", "SetMarketMapperRevShareDetailsForMarket", data);
    return promise.then(data => MsgSetMarketMapperRevShareDetailsForMarketResponse.decode(new _m0.Reader(data)));
  }

  updateUnconditionalRevShareConfig(request: MsgUpdateUnconditionalRevShareConfig): Promise<MsgUpdateUnconditionalRevShareConfigResponse> {
    const data = MsgUpdateUnconditionalRevShareConfig.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.revshare.Msg", "UpdateUnconditionalRevShareConfig", data);
    return promise.then(data => MsgUpdateUnconditionalRevShareConfigResponse.decode(new _m0.Reader(data)));
  }

  setOrderRouterRevShare(request: MsgSetOrderRouterRevShare): Promise<MsgSetOrderRouterRevShareResponse> {
    const data = MsgSetOrderRouterRevShare.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.revshare.Msg", "SetOrderRouterRevShare", data);
    return promise.then(data => MsgSetOrderRouterRevShareResponse.decode(new _m0.Reader(data)));
  }

}