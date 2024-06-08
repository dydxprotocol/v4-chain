import { Rpc } from "../../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgUnjail, MsgUnjailResponse, MsgUpdateParams, MsgUpdateParamsResponse } from "./tx";
/** Msg defines the slashing Msg service. */

export interface Msg {
  /**
   * Unjail defines a method for unjailing a jailed validator, thus returning
   * them into the bonded validator set, so they can begin receiving provisions
   * and rewards again.
   */
  unjail(request: MsgUnjail): Promise<MsgUnjailResponse>;
  /**
   * UpdateParams defines a governance operation for updating the x/slashing module
   * parameters. The authority defaults to the x/gov module account.
   * 
   * Since: cosmos-sdk 0.47
   */

  updateParams(request: MsgUpdateParams): Promise<MsgUpdateParamsResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.unjail = this.unjail.bind(this);
    this.updateParams = this.updateParams.bind(this);
  }

  unjail(request: MsgUnjail): Promise<MsgUnjailResponse> {
    const data = MsgUnjail.encode(request).finish();
    const promise = this.rpc.request("cosmos.slashing.v1beta1.Msg", "Unjail", data);
    return promise.then(data => MsgUnjailResponse.decode(new _m0.Reader(data)));
  }

  updateParams(request: MsgUpdateParams): Promise<MsgUpdateParamsResponse> {
    const data = MsgUpdateParams.encode(request).finish();
    const promise = this.rpc.request("cosmos.slashing.v1beta1.Msg", "UpdateParams", data);
    return promise.then(data => MsgUpdateParamsResponse.decode(new _m0.Reader(data)));
  }

}