import { GeneratedType, Registry } from "@cosmjs/proto-signing";
import { MsgUpdateMarketPrices, MsgCreateOracleMarket, MsgUpdateMarketParam } from "./tx";
export const registry: ReadonlyArray<[string, GeneratedType]> = [["/dydxprotocol.prices.MsgUpdateMarketPrices", MsgUpdateMarketPrices], ["/dydxprotocol.prices.MsgCreateOracleMarket", MsgCreateOracleMarket], ["/dydxprotocol.prices.MsgUpdateMarketParam", MsgUpdateMarketParam]];
export const load = (protoRegistry: Registry) => {
  registry.forEach(([typeUrl, mod]) => {
    protoRegistry.register(typeUrl, mod);
  });
};
export const MessageComposer = {
  encoded: {
    updateMarketPrices(value: MsgUpdateMarketPrices) {
      return {
        typeUrl: "/dydxprotocol.prices.MsgUpdateMarketPrices",
        value: MsgUpdateMarketPrices.encode(value).finish()
      };
    },
    createOracleMarket(value: MsgCreateOracleMarket) {
      return {
        typeUrl: "/dydxprotocol.prices.MsgCreateOracleMarket",
        value: MsgCreateOracleMarket.encode(value).finish()
      };
    },
    updateMarketParam(value: MsgUpdateMarketParam) {
      return {
        typeUrl: "/dydxprotocol.prices.MsgUpdateMarketParam",
        value: MsgUpdateMarketParam.encode(value).finish()
      };
    }
  },
  withTypeUrl: {
    updateMarketPrices(value: MsgUpdateMarketPrices) {
      return {
        typeUrl: "/dydxprotocol.prices.MsgUpdateMarketPrices",
        value
      };
    },
    createOracleMarket(value: MsgCreateOracleMarket) {
      return {
        typeUrl: "/dydxprotocol.prices.MsgCreateOracleMarket",
        value
      };
    },
    updateMarketParam(value: MsgUpdateMarketParam) {
      return {
        typeUrl: "/dydxprotocol.prices.MsgUpdateMarketParam",
        value
      };
    }
  },
  fromPartial: {
    updateMarketPrices(value: MsgUpdateMarketPrices) {
      return {
        typeUrl: "/dydxprotocol.prices.MsgUpdateMarketPrices",
        value: MsgUpdateMarketPrices.fromPartial(value)
      };
    },
    createOracleMarket(value: MsgCreateOracleMarket) {
      return {
        typeUrl: "/dydxprotocol.prices.MsgCreateOracleMarket",
        value: MsgCreateOracleMarket.fromPartial(value)
      };
    },
    updateMarketParam(value: MsgUpdateMarketParam) {
      return {
        typeUrl: "/dydxprotocol.prices.MsgUpdateMarketParam",
        value: MsgUpdateMarketParam.fromPartial(value)
      };
    }
  }
};