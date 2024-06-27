import { MsgDepositToSubaccount, MsgWithdrawFromSubaccount, MsgSendFromModuleToAccount } from "./transfer";
import { Rpc } from "../../helpers";
import { MsgCreateTransfer, MsgCreateTransferResponse, MsgDepositToSubaccountResponse, MsgWithdrawFromSubaccountResponse, MsgSendFromModuleToAccountResponse } from "./tx";
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
    /**
     * SendFromModuleToAccount initiates a new transfer from a module to an
     * `x/bank` account (should only be executed by governance).
     */
    sendFromModuleToAccount(request: MsgSendFromModuleToAccount): Promise<MsgSendFromModuleToAccountResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    createTransfer(request: MsgCreateTransfer): Promise<MsgCreateTransferResponse>;
    depositToSubaccount(request: MsgDepositToSubaccount): Promise<MsgDepositToSubaccountResponse>;
    withdrawFromSubaccount(request: MsgWithdrawFromSubaccount): Promise<MsgWithdrawFromSubaccountResponse>;
    sendFromModuleToAccount(request: MsgSendFromModuleToAccount): Promise<MsgSendFromModuleToAccountResponse>;
}
