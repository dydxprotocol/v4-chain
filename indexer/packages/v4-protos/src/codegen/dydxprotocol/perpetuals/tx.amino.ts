import { MsgAddPremiumVotes, MsgCreatePerpetual, MsgSetLiquidityTier, MsgUpdatePerpetualParams, MsgUpdateParams } from "./tx";
export const AminoConverter = {
  "/dydxprotocol.perpetuals.MsgAddPremiumVotes": {
    aminoType: "/dydxprotocol.perpetuals.MsgAddPremiumVotes",
    toAmino: MsgAddPremiumVotes.toAmino,
    fromAmino: MsgAddPremiumVotes.fromAmino
  },
  "/dydxprotocol.perpetuals.MsgCreatePerpetual": {
    aminoType: "/dydxprotocol.perpetuals.MsgCreatePerpetual",
    toAmino: MsgCreatePerpetual.toAmino,
    fromAmino: MsgCreatePerpetual.fromAmino
  },
  "/dydxprotocol.perpetuals.MsgSetLiquidityTier": {
    aminoType: "/dydxprotocol.perpetuals.MsgSetLiquidityTier",
    toAmino: MsgSetLiquidityTier.toAmino,
    fromAmino: MsgSetLiquidityTier.fromAmino
  },
  "/dydxprotocol.perpetuals.MsgUpdatePerpetualParams": {
    aminoType: "/dydxprotocol.perpetuals.MsgUpdatePerpetualParams",
    toAmino: MsgUpdatePerpetualParams.toAmino,
    fromAmino: MsgUpdatePerpetualParams.fromAmino
  },
  "/dydxprotocol.perpetuals.MsgUpdateParams": {
    aminoType: "/dydxprotocol.perpetuals.MsgUpdateParams",
    toAmino: MsgUpdateParams.toAmino,
    fromAmino: MsgUpdateParams.fromAmino
  }
};