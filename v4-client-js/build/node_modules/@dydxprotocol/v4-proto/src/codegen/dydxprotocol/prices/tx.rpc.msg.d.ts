import { Rpc } from "../../helpers";
import { MsgUpdateMarketPrices, MsgUpdateMarketPricesResponse, MsgCreateOracleMarket, MsgCreateOracleMarketResponse, MsgUpdateMarketParam, MsgUpdateMarketParamResponse } from "./tx";
/** Msg defines the Msg service. */
export interface Msg {
    /**
     * UpdateMarketPrices updates the oracle price of a market relative to
     * quoteCurrency.
     */
    updateMarketPrices(request: MsgUpdateMarketPrices): Promise<MsgUpdateMarketPricesResponse>;
    /** CreateOracleMarket creates a new oracle market. */
    createOracleMarket(request: MsgCreateOracleMarket): Promise<MsgCreateOracleMarketResponse>;
    /**
     * UpdateMarketParams allows governance to update the parameters of an
     * oracle market.
     */
    updateMarketParam(request: MsgUpdateMarketParam): Promise<MsgUpdateMarketParamResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    updateMarketPrices(request: MsgUpdateMarketPrices): Promise<MsgUpdateMarketPricesResponse>;
    createOracleMarket(request: MsgCreateOracleMarket): Promise<MsgCreateOracleMarketResponse>;
    updateMarketParam(request: MsgUpdateMarketParam): Promise<MsgUpdateMarketParamResponse>;
}
