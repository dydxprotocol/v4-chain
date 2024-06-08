import { Coin } from "../../base/v1beta1/coin";
import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmos.crisis.v1beta1";
/** GenesisState defines the crisis module's genesis state. */
export interface GenesisState {
    /**
     * constant_fee is the fee used to verify the invariant in the crisis
     * module.
     */
    constantFee: Coin;
}
export declare const GenesisState: {
    typeUrl: string;
    encode(message: GenesisState, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
    fromPartial<I extends {
        constantFee?: {
            denom?: string | undefined;
            amount?: string | undefined;
        } | undefined;
    } & {
        constantFee?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["constantFee"], keyof Coin>, never>) | undefined;
    } & Record<Exclude<keyof I, "constantFee">, never>>(object: I): GenesisState;
};
