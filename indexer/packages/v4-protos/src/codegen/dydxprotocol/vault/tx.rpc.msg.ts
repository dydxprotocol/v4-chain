import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgDepositToVault, MsgDepositToVaultResponse } from "./tx";
/** Msg defines the Msg service. */

export interface Msg {
  /** DepositToVault deposits funds into a vault. */
  depositToVault(request: MsgDepositToVault): Promise<MsgDepositToVaultResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.depositToVault = this.depositToVault.bind(this);
  }

  depositToVault(request: MsgDepositToVault): Promise<MsgDepositToVaultResponse> {
    const data = MsgDepositToVault.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.vault.Msg", "DepositToVault", data);
    return promise.then(data => MsgDepositToVaultResponse.decode(new _m0.Reader(data)));
  }

}