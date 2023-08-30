import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgAddPremiumVotes, MsgAddPremiumVotesResponse, MsgCreatePerpetual, MsgCreatePerpetualResponse } from "./tx";
/** Msg defines the Msg service. */

export interface Msg {
  /**
   * AddPremiumVotes add new samples of the funding premiums to the
   * application.
   */
  addPremiumVotes(request: MsgAddPremiumVotes): Promise<MsgAddPremiumVotesResponse>;
  /** CreatePerpetual creates a new perpetual object. */

  createPerpetual(request: MsgCreatePerpetual): Promise<MsgCreatePerpetualResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.addPremiumVotes = this.addPremiumVotes.bind(this);
    this.createPerpetual = this.createPerpetual.bind(this);
  }

  addPremiumVotes(request: MsgAddPremiumVotes): Promise<MsgAddPremiumVotesResponse> {
    const data = MsgAddPremiumVotes.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.perpetuals.Msg", "AddPremiumVotes", data);
    return promise.then(data => MsgAddPremiumVotesResponse.decode(new _m0.Reader(data)));
  }

  createPerpetual(request: MsgCreatePerpetual): Promise<MsgCreatePerpetualResponse> {
    const data = MsgCreatePerpetual.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.perpetuals.Msg", "CreatePerpetual", data);
    return promise.then(data => MsgCreatePerpetualResponse.decode(new _m0.Reader(data)));
  }

}