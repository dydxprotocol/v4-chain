import { Rpc } from "../../helpers";
import { MsgAddPremiumVotes, MsgAddPremiumVotesResponse, MsgCreatePerpetual, MsgCreatePerpetualResponse, MsgSetLiquidityTier, MsgSetLiquidityTierResponse, MsgUpdatePerpetualParams, MsgUpdatePerpetualParamsResponse, MsgUpdateParams, MsgUpdateParamsResponse } from "./tx";
/** Msg defines the Msg service. */
export interface Msg {
    /**
     * AddPremiumVotes add new samples of the funding premiums to the
     * application.
     */
    addPremiumVotes(request: MsgAddPremiumVotes): Promise<MsgAddPremiumVotesResponse>;
    /** CreatePerpetual creates a new perpetual object. */
    createPerpetual(request: MsgCreatePerpetual): Promise<MsgCreatePerpetualResponse>;
    /**
     * SetLiquidityTier creates an liquidity tier if the ID doesn't exist, and
     * updates the existing liquidity tier otherwise.
     */
    setLiquidityTier(request: MsgSetLiquidityTier): Promise<MsgSetLiquidityTierResponse>;
    /** UpdatePerpetualParams updates the parameters of a perpetual market. */
    updatePerpetualParams(request: MsgUpdatePerpetualParams): Promise<MsgUpdatePerpetualParamsResponse>;
    /** UpdateParams updates the parameters of perpetuals module. */
    updateParams(request: MsgUpdateParams): Promise<MsgUpdateParamsResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    addPremiumVotes(request: MsgAddPremiumVotes): Promise<MsgAddPremiumVotesResponse>;
    createPerpetual(request: MsgCreatePerpetual): Promise<MsgCreatePerpetualResponse>;
    setLiquidityTier(request: MsgSetLiquidityTier): Promise<MsgSetLiquidityTierResponse>;
    updatePerpetualParams(request: MsgUpdatePerpetualParams): Promise<MsgUpdatePerpetualParamsResponse>;
    updateParams(request: MsgUpdateParams): Promise<MsgUpdateParamsResponse>;
}
