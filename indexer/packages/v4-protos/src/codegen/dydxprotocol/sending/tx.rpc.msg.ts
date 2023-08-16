import { MsgDepositToSubaccount, MsgWithdrawFromSubaccount } from "./transfer";
import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgCreateTransfer, MsgCreateTransferResponse, MsgDepositToSubaccountResponse, MsgWithdrawFromSubaccountResponse } from "./tx";
/** Msg defines the Msg service. */

export interface Msg {
  /** CreateTransfer initiates a new transfer between subaccounts. */
  createTransfer(request: MsgCreateTransfer): Promise<MsgCreateTransferResponse>;
  /**
   * DepositToSubaccount initiates a new transfer from an `x/bank` account
   * to an `x/subaccounts` subaccount.
   */

  depositToSubaccount(request: MsgDepositToSubaccount): Promise<MsgDepositToSubaccountResponse>;
  /**
   * WithdrawFromSubaccount initiates a new transfer from an `x/subaccounts`
   * subaccount to an `x/bank` account.
   */

  withdrawFromSubaccount(request: MsgWithdrawFromSubaccount): Promise<MsgWithdrawFromSubaccountResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.createTransfer = this.createTransfer.bind(this);
    this.depositToSubaccount = this.depositToSubaccount.bind(this);
    this.withdrawFromSubaccount = this.withdrawFromSubaccount.bind(this);
  }

  createTransfer(request: MsgCreateTransfer): Promise<MsgCreateTransferResponse> {
    const data = MsgCreateTransfer.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.sending.Msg", "CreateTransfer", data);
    return promise.then(data => MsgCreateTransferResponse.decode(new _m0.Reader(data)));
  }

  depositToSubaccount(request: MsgDepositToSubaccount): Promise<MsgDepositToSubaccountResponse> {
    const data = MsgDepositToSubaccount.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.sending.Msg", "DepositToSubaccount", data);
    return promise.then(data => MsgDepositToSubaccountResponse.decode(new _m0.Reader(data)));
  }

  withdrawFromSubaccount(request: MsgWithdrawFromSubaccount): Promise<MsgWithdrawFromSubaccountResponse> {
    const data = MsgWithdrawFromSubaccount.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.sending.Msg", "WithdrawFromSubaccount", data);
    return promise.then(data => MsgWithdrawFromSubaccountResponse.decode(new _m0.Reader(data)));
  }

}