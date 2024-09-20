import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgDepositToMegavault, MsgDepositToMegavaultResponse, MsgWithdrawFromMegavault, MsgWithdrawFromMegavaultResponse, MsgUpdateDefaultQuotingParams, MsgUpdateDefaultQuotingParamsResponse, MsgUpdateOperatorParams, MsgUpdateOperatorParamsResponse, MsgSetVaultParams, MsgSetVaultParamsResponse, MsgUnlockShares, MsgUnlockSharesResponse, MsgAllocateToVault, MsgAllocateToVaultResponse, MsgRetrieveFromVault, MsgRetrieveFromVaultResponse } from "./tx";
/** Msg defines the Msg service. */

export interface Msg {
  /** DepositToMegavault deposits funds into megavault. */
  depositToMegavault(request: MsgDepositToMegavault): Promise<MsgDepositToMegavaultResponse>;
  /** WithdrawFromMegavault withdraws shares from megavault. */

  withdrawFromMegavault(request: MsgWithdrawFromMegavault): Promise<MsgWithdrawFromMegavaultResponse>;
  /** UpdateDefaultQuotingParams updates the default quoting params in state. */

  updateDefaultQuotingParams(request: MsgUpdateDefaultQuotingParams): Promise<MsgUpdateDefaultQuotingParamsResponse>;
  /** UpdateOperatorParams sets the parameters regarding megavault operator. */

  updateOperatorParams(request: MsgUpdateOperatorParams): Promise<MsgUpdateOperatorParamsResponse>;
  /** SetVaultParams sets the parameters of a specific vault. */

  setVaultParams(request: MsgSetVaultParams): Promise<MsgSetVaultParamsResponse>;
  /**
   * UnlockShares unlocks an owner's shares that are due to unlock by the block
   * height that this transaction is included in.
   */

  unlockShares(request: MsgUnlockShares): Promise<MsgUnlockSharesResponse>;
  /** AllocateToVault allocates funds from main vault to a vault. */

  allocateToVault(request: MsgAllocateToVault): Promise<MsgAllocateToVaultResponse>;
  /** RetrieveFromVault retrieves funds from a vault to main vault. */

  retrieveFromVault(request: MsgRetrieveFromVault): Promise<MsgRetrieveFromVaultResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.depositToMegavault = this.depositToMegavault.bind(this);
    this.withdrawFromMegavault = this.withdrawFromMegavault.bind(this);
    this.updateDefaultQuotingParams = this.updateDefaultQuotingParams.bind(this);
    this.updateOperatorParams = this.updateOperatorParams.bind(this);
    this.setVaultParams = this.setVaultParams.bind(this);
    this.unlockShares = this.unlockShares.bind(this);
    this.allocateToVault = this.allocateToVault.bind(this);
    this.retrieveFromVault = this.retrieveFromVault.bind(this);
  }

  depositToMegavault(request: MsgDepositToMegavault): Promise<MsgDepositToMegavaultResponse> {
    const data = MsgDepositToMegavault.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.vault.Msg", "DepositToMegavault", data);
    return promise.then(data => MsgDepositToMegavaultResponse.decode(new _m0.Reader(data)));
  }

  withdrawFromMegavault(request: MsgWithdrawFromMegavault): Promise<MsgWithdrawFromMegavaultResponse> {
    const data = MsgWithdrawFromMegavault.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.vault.Msg", "WithdrawFromMegavault", data);
    return promise.then(data => MsgWithdrawFromMegavaultResponse.decode(new _m0.Reader(data)));
  }

  updateDefaultQuotingParams(request: MsgUpdateDefaultQuotingParams): Promise<MsgUpdateDefaultQuotingParamsResponse> {
    const data = MsgUpdateDefaultQuotingParams.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.vault.Msg", "UpdateDefaultQuotingParams", data);
    return promise.then(data => MsgUpdateDefaultQuotingParamsResponse.decode(new _m0.Reader(data)));
  }

  updateOperatorParams(request: MsgUpdateOperatorParams): Promise<MsgUpdateOperatorParamsResponse> {
    const data = MsgUpdateOperatorParams.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.vault.Msg", "UpdateOperatorParams", data);
    return promise.then(data => MsgUpdateOperatorParamsResponse.decode(new _m0.Reader(data)));
  }

  setVaultParams(request: MsgSetVaultParams): Promise<MsgSetVaultParamsResponse> {
    const data = MsgSetVaultParams.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.vault.Msg", "SetVaultParams", data);
    return promise.then(data => MsgSetVaultParamsResponse.decode(new _m0.Reader(data)));
  }

  unlockShares(request: MsgUnlockShares): Promise<MsgUnlockSharesResponse> {
    const data = MsgUnlockShares.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.vault.Msg", "UnlockShares", data);
    return promise.then(data => MsgUnlockSharesResponse.decode(new _m0.Reader(data)));
  }

  allocateToVault(request: MsgAllocateToVault): Promise<MsgAllocateToVaultResponse> {
    const data = MsgAllocateToVault.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.vault.Msg", "AllocateToVault", data);
    return promise.then(data => MsgAllocateToVaultResponse.decode(new _m0.Reader(data)));
  }

  retrieveFromVault(request: MsgRetrieveFromVault): Promise<MsgRetrieveFromVaultResponse> {
    const data = MsgRetrieveFromVault.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.vault.Msg", "RetrieveFromVault", data);
    return promise.then(data => MsgRetrieveFromVaultResponse.decode(new _m0.Reader(data)));
  }

}