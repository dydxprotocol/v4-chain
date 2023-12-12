import { MsgAcknowledgeBridges, MsgCompleteBridge, MsgUpdateEventParams, MsgUpdateProposeParams, MsgUpdateSafetyParams } from "./tx";
export const AminoConverter = {
  "/dydxprotocol.bridge.MsgAcknowledgeBridges": {
    aminoType: "/dydxprotocol.bridge.MsgAcknowledgeBridges",
    toAmino: MsgAcknowledgeBridges.toAmino,
    fromAmino: MsgAcknowledgeBridges.fromAmino
  },
  "/dydxprotocol.bridge.MsgCompleteBridge": {
    aminoType: "/dydxprotocol.bridge.MsgCompleteBridge",
    toAmino: MsgCompleteBridge.toAmino,
    fromAmino: MsgCompleteBridge.fromAmino
  },
  "/dydxprotocol.bridge.MsgUpdateEventParams": {
    aminoType: "/dydxprotocol.bridge.MsgUpdateEventParams",
    toAmino: MsgUpdateEventParams.toAmino,
    fromAmino: MsgUpdateEventParams.fromAmino
  },
  "/dydxprotocol.bridge.MsgUpdateProposeParams": {
    aminoType: "/dydxprotocol.bridge.MsgUpdateProposeParams",
    toAmino: MsgUpdateProposeParams.toAmino,
    fromAmino: MsgUpdateProposeParams.fromAmino
  },
  "/dydxprotocol.bridge.MsgUpdateSafetyParams": {
    aminoType: "/dydxprotocol.bridge.MsgUpdateSafetyParams",
    toAmino: MsgUpdateSafetyParams.toAmino,
    fromAmino: MsgUpdateSafetyParams.fromAmino
  }
};