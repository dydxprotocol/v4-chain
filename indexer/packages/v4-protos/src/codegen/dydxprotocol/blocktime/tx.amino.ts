import { MsgUpdateDowntimeParams } from "./tx";
export const AminoConverter = {
  "/dydxprotocol.blocktime.MsgUpdateDowntimeParams": {
    aminoType: "/dydxprotocol.blocktime.MsgUpdateDowntimeParams",
    toAmino: MsgUpdateDowntimeParams.toAmino,
    fromAmino: MsgUpdateDowntimeParams.fromAmino
  }
};