import { MsgDepositToSubaccount, MsgWithdrawFromSubaccount, MsgSendFromModuleToAccount } from "./transfer";
import { GeneratedType, Registry } from "@cosmjs/proto-signing";
import { MsgCreateTransfer } from "./tx";
export const registry: ReadonlyArray<[string, GeneratedType]> = [["/dydxprotocol.sending.MsgCreateTransfer", MsgCreateTransfer], ["/dydxprotocol.sending.MsgDepositToSubaccount", MsgDepositToSubaccount], ["/dydxprotocol.sending.MsgWithdrawFromSubaccount", MsgWithdrawFromSubaccount], ["/dydxprotocol.sending.MsgSendFromModuleToAccount", MsgSendFromModuleToAccount]];
export const load = (protoRegistry: Registry) => {
  registry.forEach(([typeUrl, mod]) => {
    protoRegistry.register(typeUrl, mod);
  });
};
export const MessageComposer = {
  encoded: {
    createTransfer(value: MsgCreateTransfer) {
      return {
        typeUrl: "/dydxprotocol.sending.MsgCreateTransfer",
        value: MsgCreateTransfer.encode(value).finish()
      };
    },
    depositToSubaccount(value: MsgDepositToSubaccount) {
      return {
        typeUrl: "/dydxprotocol.sending.MsgDepositToSubaccount",
        value: MsgDepositToSubaccount.encode(value).finish()
      };
    },
    withdrawFromSubaccount(value: MsgWithdrawFromSubaccount) {
      return {
        typeUrl: "/dydxprotocol.sending.MsgWithdrawFromSubaccount",
        value: MsgWithdrawFromSubaccount.encode(value).finish()
      };
    },
    sendFromModuleToAccount(value: MsgSendFromModuleToAccount) {
      return {
        typeUrl: "/dydxprotocol.sending.MsgSendFromModuleToAccount",
        value: MsgSendFromModuleToAccount.encode(value).finish()
      };
    }
  },
  withTypeUrl: {
    createTransfer(value: MsgCreateTransfer) {
      return {
        typeUrl: "/dydxprotocol.sending.MsgCreateTransfer",
        value
      };
    },
    depositToSubaccount(value: MsgDepositToSubaccount) {
      return {
        typeUrl: "/dydxprotocol.sending.MsgDepositToSubaccount",
        value
      };
    },
    withdrawFromSubaccount(value: MsgWithdrawFromSubaccount) {
      return {
        typeUrl: "/dydxprotocol.sending.MsgWithdrawFromSubaccount",
        value
      };
    },
    sendFromModuleToAccount(value: MsgSendFromModuleToAccount) {
      return {
        typeUrl: "/dydxprotocol.sending.MsgSendFromModuleToAccount",
        value
      };
    }
  },
  fromPartial: {
    createTransfer(value: MsgCreateTransfer) {
      return {
        typeUrl: "/dydxprotocol.sending.MsgCreateTransfer",
        value: MsgCreateTransfer.fromPartial(value)
      };
    },
    depositToSubaccount(value: MsgDepositToSubaccount) {
      return {
        typeUrl: "/dydxprotocol.sending.MsgDepositToSubaccount",
        value: MsgDepositToSubaccount.fromPartial(value)
      };
    },
    withdrawFromSubaccount(value: MsgWithdrawFromSubaccount) {
      return {
        typeUrl: "/dydxprotocol.sending.MsgWithdrawFromSubaccount",
        value: MsgWithdrawFromSubaccount.fromPartial(value)
      };
    },
    sendFromModuleToAccount(value: MsgSendFromModuleToAccount) {
      return {
        typeUrl: "/dydxprotocol.sending.MsgSendFromModuleToAccount",
        value: MsgSendFromModuleToAccount.fromPartial(value)
      };
    }
  }
};