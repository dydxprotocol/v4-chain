import { MsgSetVestEntry, MsgDeleteVestEntry } from "./tx";
export const AminoConverter = {
  "/dydxprotocol.vest.MsgSetVestEntry": {
    aminoType: "/dydxprotocol.vest.MsgSetVestEntry",
    toAmino: MsgSetVestEntry.toAmino,
    fromAmino: MsgSetVestEntry.fromAmino
  },
  "/dydxprotocol.vest.MsgDeleteVestEntry": {
    aminoType: "/dydxprotocol.vest.MsgDeleteVestEntry",
    toAmino: MsgDeleteVestEntry.toAmino,
    fromAmino: MsgDeleteVestEntry.fromAmino
  }
};