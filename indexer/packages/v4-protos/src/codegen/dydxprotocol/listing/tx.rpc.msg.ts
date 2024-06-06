import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgEnablePermissionlessMarketListing, MsgEnablePermissionlessMarketListingResponse } from "./tx";
/** Msg defines the Msg service. */

export interface Msg {
  /**
   * EnablePermissionlessMarketListing enables/disables permissionless market
   * listing
   */
  enablePermissionlessMarketListing(request: MsgEnablePermissionlessMarketListing): Promise<MsgEnablePermissionlessMarketListingResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.enablePermissionlessMarketListing = this.enablePermissionlessMarketListing.bind(this);
  }

  enablePermissionlessMarketListing(request: MsgEnablePermissionlessMarketListing): Promise<MsgEnablePermissionlessMarketListingResponse> {
    const data = MsgEnablePermissionlessMarketListing.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.listing.Msg", "EnablePermissionlessMarketListing", data);
    return promise.then(data => MsgEnablePermissionlessMarketListingResponse.decode(new _m0.Reader(data)));
  }

}