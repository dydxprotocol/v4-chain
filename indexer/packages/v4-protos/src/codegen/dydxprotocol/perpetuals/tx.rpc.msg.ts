import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgAddPremiumVotes, MsgAddPremiumVotesResponse } from "./tx";
/** Msg defines the Msg service. */

export interface Msg {
  /**
   * AddPremiumVotes add new samples of the funding premiums to the
   * application.
   */
  addPremiumVotes(request: MsgAddPremiumVotes): Promise<MsgAddPremiumVotesResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.addPremiumVotes = this.addPremiumVotes.bind(this);
  }

  addPremiumVotes(request: MsgAddPremiumVotes): Promise<MsgAddPremiumVotesResponse> {
    const data = MsgAddPremiumVotes.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.perpetuals.Msg", "AddPremiumVotes", data);
    return promise.then(data => MsgAddPremiumVotesResponse.decode(new _m0.Reader(data)));
  }

}