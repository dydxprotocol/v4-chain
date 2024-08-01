import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgDepositToVault, MsgDepositToVaultResponse, MsgUpdateParams, MsgUpdateParamsResponse } from "./tx";
/** Msg defines the Msg service. */

export interface Msg {
  /** DepositToVault deposits funds into a vault. */
  depositToVault(request: MsgDepositToVault): Promise<MsgDepositToVaultResponse>;
  /** UpdateParams updates the Params in state. */

  updateParams(request: MsgUpdateParams): Promise<MsgUpdateParamsResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.depositToVault = this.depositToVault.bind(this);
    this.updateParams = this.updateParams.bind(this);
  }

  depositToVault(request: MsgDepositToVault): Promise<MsgDepositToVaultResponse> {
    const data = MsgDepositToVault.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.vault.Msg", "DepositToVault", data);
    return promise.then(data => MsgDepositToVaultResponse.decode(new _m0.Reader(data)));
  }

  updateParams(request: MsgUpdateParams): Promise<MsgUpdateParamsResponse> {
    const data = MsgUpdateParams.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.vault.Msg", "UpdateParams", data);
    return promise.then(data => MsgUpdateParamsResponse.decode(new _m0.Reader(data)));
  }

}