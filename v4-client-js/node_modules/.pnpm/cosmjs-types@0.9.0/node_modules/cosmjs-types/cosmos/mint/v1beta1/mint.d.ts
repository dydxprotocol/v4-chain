import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmos.mint.v1beta1";
/** Minter represents the minting state. */
export interface Minter {
    /** current annual inflation rate */
    inflation: string;
    /** current annual expected provisions */
    annualProvisions: string;
}
/** Params defines the parameters for the x/mint module. */
export interface Params {
    /** type of coin to mint */
    mintDenom: string;
    /** maximum annual change in inflation rate */
    inflationRateChange: string;
    /** maximum inflation rate */
    inflationMax: string;
    /** minimum inflation rate */
    inflationMin: string;
    /** goal of percent bonded atoms */
    goalBonded: string;
    /** expected blocks per year */
    blocksPerYear: bigint;
}
export declare const Minter: {
    typeUrl: string;
    encode(message: Minter, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Minter;
    fromJSON(object: any): Minter;
    toJSON(message: Minter): unknown;
    fromPartial<I extends {
        inflation?: string | undefined;
        annualProvisions?: string | undefined;
    } & {
        inflation?: string | undefined;
        annualProvisions?: string | undefined;
    } & Record<Exclude<keyof I, keyof Minter>, never>>(object: I): Minter;
};
export declare const Params: {
    typeUrl: string;
    encode(message: Params, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Params;
    fromJSON(object: any): Params;
    toJSON(message: Params): unknown;
    fromPartial<I extends {
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
    } & Record<Exclude<keyof I, keyof Params>, never>>(object: I): Params;
};
