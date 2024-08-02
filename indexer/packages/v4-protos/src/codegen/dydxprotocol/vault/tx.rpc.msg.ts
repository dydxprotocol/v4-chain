import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgDepositToVault, MsgDepositToVaultResponse, MsgUpdateDefaultQuotingParams, MsgUpdateDefaultQuotingParamsResponse } from "./tx";
/** Msg defines the Msg service. */

export interface Msg {
  /** DepositToVault deposits funds into a vault. */
  depositToVault(request: MsgDepositToVault): Promise<MsgDepositToVaultResponse>;
  /** UpdateDefaultQuotingParams updates the default quoting params in state. */

  updateDefaultQuotingParams(request: MsgUpdateDefaultQuotingParams): Promise<MsgUpdateDefaultQuotingParamsResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.depositToVault = this.depositToVault.bind(this);
    this.updateDefaultQuotingParams = this.updateDefaultQuotingParams.bind(this);
  }

  depositToVault(request: MsgDepositToVault): Promise<MsgDepositToVaultResponse> {
    const data = MsgDepositToVault.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.vault.Msg", "DepositToVault", data);
    return promise.then(data => MsgDepositToVaultResponse.decode(new _m0.Reader(data)));
  }

  updateDefaultQuotingParams(request: MsgUpdateDefaultQuotingParams): Promise<MsgUpdateDefaultQuotingParamsResponse> {
    const data = MsgUpdateDefaultQuotingParams.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.vault.Msg", "UpdateDefaultQuotingParams", data);
    return promise.then(data => MsgUpdateDefaultQuotingParamsResponse.decode(new _m0.Reader(data)));
  }

}