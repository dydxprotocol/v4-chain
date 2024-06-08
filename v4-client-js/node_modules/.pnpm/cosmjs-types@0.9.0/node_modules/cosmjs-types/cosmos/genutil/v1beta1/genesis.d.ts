import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmos.genutil.v1beta1";
/** GenesisState defines the raw genesis transaction in JSON. */
export interface GenesisState {
    /** gen_txs defines the genesis transactions. */
    genTxs: Uint8Array[];
}
export declare const GenesisState: {
    typeUrl: string;
    encode(message: GenesisState, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
    fromPartial<I extends {
        genTxs?: Uint8Array[] | undefined;
    } & {
        genTxs?: (Uint8Array[] & Uint8Array[] & Record<Exclude<keyof I["genTxs"], keyof Uint8Array[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "genTxs">, never>>(object: I): GenesisState;
};
