import { MsgDelayMessage } from "./tx";
export const AminoConverter = {
  "/dydxprotocol.delaymsg.MsgDelayMessage": {
    aminoType: "/dydxprotocol.delaymsg.MsgDelayMessage",
    toAmino: MsgDelayMessage.toAmino,
    fromAmino: MsgDelayMessage.fromAmino
  }
};