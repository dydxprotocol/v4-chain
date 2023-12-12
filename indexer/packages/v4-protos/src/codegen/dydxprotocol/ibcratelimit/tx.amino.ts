import { MsgSetLimitParams, MsgDeleteLimitParams } from "./tx";
export const AminoConverter = {
  "/dydxprotocol.ibcratelimit.MsgSetLimitParams": {
    aminoType: "/dydxprotocol.ibcratelimit.MsgSetLimitParams",
    toAmino: MsgSetLimitParams.toAmino,
    fromAmino: MsgSetLimitParams.fromAmino
  },
  "/dydxprotocol.ibcratelimit.MsgDeleteLimitParams": {
    aminoType: "/dydxprotocol.ibcratelimit.MsgDeleteLimitParams",
    toAmino: MsgDeleteLimitParams.toAmino,
    fromAmino: MsgDeleteLimitParams.fromAmino
  }
};