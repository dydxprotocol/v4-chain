import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgSetMarketMapperRevenueShare, MsgSetMarketMapperRevenueShareResponse } from "./tx";
/** Msg defines the Msg service. */

export interface Msg {
  /** SetMarketMapperRevenueShare creates a new revenue share for a market mapper. */
  setMarketMapperRevenueShare(request: MsgSetMarketMapperRevenueShare): Promise<MsgSetMarketMapperRevenueShareResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.setMarketMapperRevenueShare = this.setMarketMapperRevenueShare.bind(this);
  }

  setMarketMapperRevenueShare(request: MsgSetMarketMapperRevenueShare): Promise<MsgSetMarketMapperRevenueShareResponse> {
    const data = MsgSetMarketMapperRevenueShare.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.revshare.Msg", "SetMarketMapperRevenueShare", data);
    return promise.then(data => MsgSetMarketMapperRevenueShareResponse.decode(new _m0.Reader(data)));
  }

}