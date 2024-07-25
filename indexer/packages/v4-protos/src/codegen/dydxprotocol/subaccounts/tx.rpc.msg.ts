import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgClaimYieldForSubaccount, MsgClaimYieldForSubaccountResponse } from "./tx";
/** Msg defines the Msg service. */

export interface Msg {
  /** ClaimYieldForSubaccount claims the yield for the provided subaccount and persists it to state. */
  claimYieldForSubaccount(request: MsgClaimYieldForSubaccount): Promise<MsgClaimYieldForSubaccountResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.claimYieldForSubaccount = this.claimYieldForSubaccount.bind(this);
  }

  claimYieldForSubaccount(request: MsgClaimYieldForSubaccount): Promise<MsgClaimYieldForSubaccountResponse> {
    const data = MsgClaimYieldForSubaccount.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.subaccounts.Msg", "ClaimYieldForSubaccount", data);
    return promise.then(data => MsgClaimYieldForSubaccountResponse.decode(new _m0.Reader(data)));
  }

}