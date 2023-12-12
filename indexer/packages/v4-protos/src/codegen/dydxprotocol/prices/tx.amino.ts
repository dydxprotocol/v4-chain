import { MsgUpdateMarketPrices, MsgCreateOracleMarket, MsgUpdateMarketParam } from "./tx";
export const AminoConverter = {
  "/dydxprotocol.prices.MsgUpdateMarketPrices": {
    aminoType: "/dydxprotocol.prices.MsgUpdateMarketPrices",
    toAmino: MsgUpdateMarketPrices.toAmino,
    fromAmino: MsgUpdateMarketPrices.fromAmino
  },
  "/dydxprotocol.prices.MsgCreateOracleMarket": {
    aminoType: "/dydxprotocol.prices.MsgCreateOracleMarket",
    toAmino: MsgCreateOracleMarket.toAmino,
    fromAmino: MsgCreateOracleMarket.fromAmino
  },
  "/dydxprotocol.prices.MsgUpdateMarketParam": {
    aminoType: "/dydxprotocol.prices.MsgUpdateMarketParam",
    toAmino: MsgUpdateMarketParam.toAmino,
    fromAmino: MsgUpdateMarketParam.fromAmino
  }
};