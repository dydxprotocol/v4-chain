import { Minter, Params } from "./mint";
import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmos.mint.v1beta1";
/** GenesisState defines the mint module's genesis state. */
export interface GenesisState {
    /** minter is a space for holding current inflation information. */
    minter: Minter;
    /** params defines all the parameters of the module. */
    params: Params;
}
export declare const GenesisState: {
    typeUrl: string;
    encode(message: GenesisState, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
    fromPartial<I extends {
        minter?: {
            inflation?: string | undefined;
            annualProvisions?: string | undefined;
        } | undefined;
        params?: {
            mintDenom?: string | undefined;
            inflationRateChange?: string | undefined;
            inflationMax?: string | undefined;
            inflationMin?: string | undefined;
            goalBonded?: string | undefined;
            blocksPerYear?: bigint | undefined;
        } | undefined;
    } & {
        minter?: ({
            inflation?: string | undefined;
            annualProvisions?: string | undefined;
        } & {
            inflation?: string | undefined;
            annualProvisions?: string | undefined;
        } & Record<Exclude<keyof I["minter"], keyof Minter>, never>) | undefined;
        params?: ({
            mintDenom?: string | undefined;
            inflationRateChange?: string | undefined;
            inflationMax?: string | undefined;
            inflationMin?: string | undefined;
            goalBonded?: string | undefined;
            blocksPerYear?: bigint | undefined;
        } & {
            mintDenom?: string | undefined;
            inflationRateChange?: string | undefined;
            inflationMax?: string | undefined;
            inflationMin?: string | undefined;
            goalBonded?: string | undefined;
            blocksPerYear?: bigint | undefined;
        } & Record<Exclude<keyof I["params"], keyof Params>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof GenesisState>, never>>(object: I): GenesisState;
};
