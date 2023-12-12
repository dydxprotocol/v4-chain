import { MsgUpdateParams } from "./tx";
export const AminoConverter = {
  "/dydxprotocol.stats.MsgUpdateParams": {
    aminoType: "/dydxprotocol.stats.MsgUpdateParams",
    toAmino: MsgUpdateParams.toAmino,
    fromAmino: MsgUpdateParams.fromAmino
  }
};