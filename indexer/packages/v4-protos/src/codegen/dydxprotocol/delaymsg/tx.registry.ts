import { GeneratedType, Registry } from "@cosmjs/proto-signing";
import { MsgDelayMessage } from "./tx";
export const registry: ReadonlyArray<[string, GeneratedType]> = [["/dydxprotocol.delaymsg.MsgDelayMessage", MsgDelayMessage]];
export const load = (protoRegistry: Registry) => {
  registry.forEach(([typeUrl, mod]) => {
    protoRegistry.register(typeUrl, mod);
  });
};
export const MessageComposer = {
  encoded: {
    delayMessage(value: MsgDelayMessage) {
      return {
        typeUrl: "/dydxprotocol.delaymsg.MsgDelayMessage",
        value: MsgDelayMessage.encode(value).finish()
      };
    }
  },
  withTypeUrl: {
    delayMessage(value: MsgDelayMessage) {
      return {
        typeUrl: "/dydxprotocol.delaymsg.MsgDelayMessage",
        value
      };
    }
  },
  fromPartial: {
    delayMessage(value: MsgDelayMessage) {
      return {
        typeUrl: "/dydxprotocol.delaymsg.MsgDelayMessage",
        value: MsgDelayMessage.fromPartial(value)
      };
    }
  }
};