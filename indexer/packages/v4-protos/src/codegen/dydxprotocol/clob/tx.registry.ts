import { GeneratedType, Registry } from "@cosmjs/proto-signing";
import { MsgProposedOperations, MsgPlaceOrder, MsgCancelOrder, MsgCreateClobPair, MsgUpdateClobPair, MsgUpdateEquityTierLimitConfiguration, MsgUpdateBlockRateLimitConfiguration, MsgUpdateLiquidationsConfig } from "./tx";
export const registry: ReadonlyArray<[string, GeneratedType]> = [["/dydxprotocol.clob.MsgProposedOperations", MsgProposedOperations], ["/dydxprotocol.clob.MsgPlaceOrder", MsgPlaceOrder], ["/dydxprotocol.clob.MsgCancelOrder", MsgCancelOrder], ["/dydxprotocol.clob.MsgCreateClobPair", MsgCreateClobPair], ["/dydxprotocol.clob.MsgUpdateClobPair", MsgUpdateClobPair], ["/dydxprotocol.clob.MsgUpdateEquityTierLimitConfiguration", MsgUpdateEquityTierLimitConfiguration], ["/dydxprotocol.clob.MsgUpdateBlockRateLimitConfiguration", MsgUpdateBlockRateLimitConfiguration], ["/dydxprotocol.clob.MsgUpdateLiquidationsConfig", MsgUpdateLiquidationsConfig]];
export const load = (protoRegistry: Registry) => {
  registry.forEach(([typeUrl, mod]) => {
    protoRegistry.register(typeUrl, mod);
  });
};
export const MessageComposer = {
  encoded: {
    proposedOperations(value: MsgProposedOperations) {
      return {
        typeUrl: "/dydxprotocol.clob.MsgProposedOperations",
        value: MsgProposedOperations.encode(value).finish()
      };
    },
    placeOrder(value: MsgPlaceOrder) {
      return {
        typeUrl: "/dydxprotocol.clob.MsgPlaceOrder",
        value: MsgPlaceOrder.encode(value).finish()
      };
    },
    cancelOrder(value: MsgCancelOrder) {
      return {
        typeUrl: "/dydxprotocol.clob.MsgCancelOrder",
        value: MsgCancelOrder.encode(value).finish()
      };
    },
    createClobPair(value: MsgCreateClobPair) {
      return {
        typeUrl: "/dydxprotocol.clob.MsgCreateClobPair",
        value: MsgCreateClobPair.encode(value).finish()
      };
    },
    updateClobPair(value: MsgUpdateClobPair) {
      return {
        typeUrl: "/dydxprotocol.clob.MsgUpdateClobPair",
        value: MsgUpdateClobPair.encode(value).finish()
      };
    },
    updateEquityTierLimitConfiguration(value: MsgUpdateEquityTierLimitConfiguration) {
      return {
        typeUrl: "/dydxprotocol.clob.MsgUpdateEquityTierLimitConfiguration",
        value: MsgUpdateEquityTierLimitConfiguration.encode(value).finish()
      };
    },
    updateBlockRateLimitConfiguration(value: MsgUpdateBlockRateLimitConfiguration) {
      return {
        typeUrl: "/dydxprotocol.clob.MsgUpdateBlockRateLimitConfiguration",
        value: MsgUpdateBlockRateLimitConfiguration.encode(value).finish()
      };
    },
    updateLiquidationsConfig(value: MsgUpdateLiquidationsConfig) {
      return {
        typeUrl: "/dydxprotocol.clob.MsgUpdateLiquidationsConfig",
        value: MsgUpdateLiquidationsConfig.encode(value).finish()
      };
    }
  },
  withTypeUrl: {
    proposedOperations(value: MsgProposedOperations) {
      return {
        typeUrl: "/dydxprotocol.clob.MsgProposedOperations",
        value
      };
    },
    placeOrder(value: MsgPlaceOrder) {
      return {
        typeUrl: "/dydxprotocol.clob.MsgPlaceOrder",
        value
      };
    },
    cancelOrder(value: MsgCancelOrder) {
      return {
        typeUrl: "/dydxprotocol.clob.MsgCancelOrder",
        value
      };
    },
    createClobPair(value: MsgCreateClobPair) {
      return {
        typeUrl: "/dydxprotocol.clob.MsgCreateClobPair",
        value
      };
    },
    updateClobPair(value: MsgUpdateClobPair) {
      return {
        typeUrl: "/dydxprotocol.clob.MsgUpdateClobPair",
        value
      };
    },
    updateEquityTierLimitConfiguration(value: MsgUpdateEquityTierLimitConfiguration) {
      return {
        typeUrl: "/dydxprotocol.clob.MsgUpdateEquityTierLimitConfiguration",
        value
      };
    },
    updateBlockRateLimitConfiguration(value: MsgUpdateBlockRateLimitConfiguration) {
      return {
        typeUrl: "/dydxprotocol.clob.MsgUpdateBlockRateLimitConfiguration",
        value
      };
    },
    updateLiquidationsConfig(value: MsgUpdateLiquidationsConfig) {
      return {
        typeUrl: "/dydxprotocol.clob.MsgUpdateLiquidationsConfig",
        value
      };
    }
  },
  fromPartial: {
    proposedOperations(value: MsgProposedOperations) {
      return {
        typeUrl: "/dydxprotocol.clob.MsgProposedOperations",
        value: MsgProposedOperations.fromPartial(value)
      };
    },
    placeOrder(value: MsgPlaceOrder) {
      return {
        typeUrl: "/dydxprotocol.clob.MsgPlaceOrder",
        value: MsgPlaceOrder.fromPartial(value)
      };
    },
    cancelOrder(value: MsgCancelOrder) {
      return {
        typeUrl: "/dydxprotocol.clob.MsgCancelOrder",
        value: MsgCancelOrder.fromPartial(value)
      };
    },
    createClobPair(value: MsgCreateClobPair) {
      return {
        typeUrl: "/dydxprotocol.clob.MsgCreateClobPair",
        value: MsgCreateClobPair.fromPartial(value)
      };
    },
    updateClobPair(value: MsgUpdateClobPair) {
      return {
        typeUrl: "/dydxprotocol.clob.MsgUpdateClobPair",
        value: MsgUpdateClobPair.fromPartial(value)
      };
    },
    updateEquityTierLimitConfiguration(value: MsgUpdateEquityTierLimitConfiguration) {
      return {
        typeUrl: "/dydxprotocol.clob.MsgUpdateEquityTierLimitConfiguration",
        value: MsgUpdateEquityTierLimitConfiguration.fromPartial(value)
      };
    },
    updateBlockRateLimitConfiguration(value: MsgUpdateBlockRateLimitConfiguration) {
      return {
        typeUrl: "/dydxprotocol.clob.MsgUpdateBlockRateLimitConfiguration",
        value: MsgUpdateBlockRateLimitConfiguration.fromPartial(value)
      };
    },
    updateLiquidationsConfig(value: MsgUpdateLiquidationsConfig) {
      return {
        typeUrl: "/dydxprotocol.clob.MsgUpdateLiquidationsConfig",
        value: MsgUpdateLiquidationsConfig.fromPartial(value)
      };
    }
  }
};