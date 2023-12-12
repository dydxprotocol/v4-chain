import { MsgProposedOperations, MsgPlaceOrder, MsgCancelOrder, MsgCreateClobPair, MsgUpdateClobPair, MsgUpdateEquityTierLimitConfiguration, MsgUpdateBlockRateLimitConfiguration, MsgUpdateLiquidationsConfig } from "./tx";
export const AminoConverter = {
  "/dydxprotocol.clob.MsgProposedOperations": {
    aminoType: "/dydxprotocol.clob.MsgProposedOperations",
    toAmino: MsgProposedOperations.toAmino,
    fromAmino: MsgProposedOperations.fromAmino
  },
  "/dydxprotocol.clob.MsgPlaceOrder": {
    aminoType: "/dydxprotocol.clob.MsgPlaceOrder",
    toAmino: MsgPlaceOrder.toAmino,
    fromAmino: MsgPlaceOrder.fromAmino
  },
  "/dydxprotocol.clob.MsgCancelOrder": {
    aminoType: "/dydxprotocol.clob.MsgCancelOrder",
    toAmino: MsgCancelOrder.toAmino,
    fromAmino: MsgCancelOrder.fromAmino
  },
  "/dydxprotocol.clob.MsgCreateClobPair": {
    aminoType: "/dydxprotocol.clob.MsgCreateClobPair",
    toAmino: MsgCreateClobPair.toAmino,
    fromAmino: MsgCreateClobPair.fromAmino
  },
  "/dydxprotocol.clob.MsgUpdateClobPair": {
    aminoType: "/dydxprotocol.clob.MsgUpdateClobPair",
    toAmino: MsgUpdateClobPair.toAmino,
    fromAmino: MsgUpdateClobPair.fromAmino
  },
  "/dydxprotocol.clob.MsgUpdateEquityTierLimitConfiguration": {
    aminoType: "/dydxprotocol.clob.MsgUpdateEquityTierLimitConfiguration",
    toAmino: MsgUpdateEquityTierLimitConfiguration.toAmino,
    fromAmino: MsgUpdateEquityTierLimitConfiguration.fromAmino
  },
  "/dydxprotocol.clob.MsgUpdateBlockRateLimitConfiguration": {
    aminoType: "/dydxprotocol.clob.MsgUpdateBlockRateLimitConfiguration",
    toAmino: MsgUpdateBlockRateLimitConfiguration.toAmino,
    fromAmino: MsgUpdateBlockRateLimitConfiguration.fromAmino
  },
  "/dydxprotocol.clob.MsgUpdateLiquidationsConfig": {
    aminoType: "/dydxprotocol.clob.MsgUpdateLiquidationsConfig",
    toAmino: MsgUpdateLiquidationsConfig.toAmino,
    fromAmino: MsgUpdateLiquidationsConfig.fromAmino
  }
};