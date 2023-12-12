import { MsgDepositToSubaccount, MsgWithdrawFromSubaccount, MsgSendFromModuleToAccount } from "./transfer";
import { MsgCreateTransfer } from "./tx";
export const AminoConverter = {
  "/dydxprotocol.sending.MsgCreateTransfer": {
    aminoType: "/dydxprotocol.sending.MsgCreateTransfer",
    toAmino: MsgCreateTransfer.toAmino,
    fromAmino: MsgCreateTransfer.fromAmino
  },
  "/dydxprotocol.sending.MsgDepositToSubaccount": {
    aminoType: "/dydxprotocol.sending.MsgDepositToSubaccount",
    toAmino: MsgDepositToSubaccount.toAmino,
    fromAmino: MsgDepositToSubaccount.fromAmino
  },
  "/dydxprotocol.sending.MsgWithdrawFromSubaccount": {
    aminoType: "/dydxprotocol.sending.MsgWithdrawFromSubaccount",
    toAmino: MsgWithdrawFromSubaccount.toAmino,
    fromAmino: MsgWithdrawFromSubaccount.fromAmino
  },
  "/dydxprotocol.sending.MsgSendFromModuleToAccount": {
    aminoType: "/dydxprotocol.sending.MsgSendFromModuleToAccount",
    toAmino: MsgSendFromModuleToAccount.toAmino,
    fromAmino: MsgSendFromModuleToAccount.fromAmino
  }
};