import { GeneratedType, Registry } from "@cosmjs/proto-signing";
import { MsgUpdateDowntimeParams } from "./tx";
export const registry: ReadonlyArray<[string, GeneratedType]> = [["/dydxprotocol.blocktime.MsgUpdateDowntimeParams", MsgUpdateDowntimeParams]];
export const load = (protoRegistry: Registry) => {
  registry.forEach(([typeUrl, mod]) => {
    protoRegistry.register(typeUrl, mod);
  });
};
export const MessageComposer = {
  encoded: {
    updateDowntimeParams(value: MsgUpdateDowntimeParams) {
      return {
        typeUrl: "/dydxprotocol.blocktime.MsgUpdateDowntimeParams",
        value: MsgUpdateDowntimeParams.encode(value).finish()
      };
    }
  },
  withTypeUrl: {
    updateDowntimeParams(value: MsgUpdateDowntimeParams) {
      return {
        typeUrl: "/dydxprotocol.blocktime.MsgUpdateDowntimeParams",
        value
      };
    }
  },
  fromPartial: {
    updateDowntimeParams(value: MsgUpdateDowntimeParams) {
      return {
        typeUrl: "/dydxprotocol.blocktime.MsgUpdateDowntimeParams",
        value: MsgUpdateDowntimeParams.fromPartial(value)
      };
    }
  }
};