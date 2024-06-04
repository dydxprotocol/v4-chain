import { BinaryReader, BinaryWriter } from "../../binary";
export declare const protobufPackage = "tendermint.crypto";
export interface Proof {
    total: bigint;
    index: bigint;
    leafHash: Uint8Array;
    aunts: Uint8Array[];
}
export interface ValueOp {
    /** Encoded in ProofOp.Key. */
    key: Uint8Array;
    /** To encode in ProofOp.Data */
    proof?: Proof;
}
export interface DominoOp {
    key: string;
    input: string;
    output: string;
}
/**
 * ProofOp defines an operation used for calculating Merkle root
 * The data could be arbitrary format, providing nessecary data
 * for example neighbouring node hash
 */
export interface ProofOp {
    type: string;
    key: Uint8Array;
    data: Uint8Array;
}
/** ProofOps is Merkle proof defined by the list of ProofOps */
export interface ProofOps {
    ops: ProofOp[];
}
export declare const Proof: {
    typeUrl: string;
    encode(message: Proof, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Proof;
    fromJSON(object: any): Proof;
    toJSON(message: Proof): unknown;
    fromPartial<I extends {
        total?: bigint | undefined;
        index?: bigint | undefined;
        leafHash?: Uint8Array | undefined;
        aunts?: Uint8Array[] | undefined;
    } & {
        total?: bigint | undefined;
        index?: bigint | undefined;
        leafHash?: Uint8Array | undefined;
        aunts?: (Uint8Array[] & Uint8Array[] & Record<Exclude<keyof I["aunts"], keyof Uint8Array[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof Proof>, never>>(object: I): Proof;
};
export declare const ValueOp: {
    typeUrl: string;
    encode(message: ValueOp, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ValueOp;
    fromJSON(object: any): ValueOp;
    toJSON(message: ValueOp): unknown;
    fromPartial<I extends {
        key?: Uint8Array | undefined;
        proof?: {
            total?: bigint | undefined;
            index?: bigint | undefined;
            leafHash?: Uint8Array | undefined;
            aunts?: Uint8Array[] | undefined;
        } | undefined;
    } & {
        key?: Uint8Array | undefined;
        proof?: ({
            total?: bigint | undefined;
            index?: bigint | undefined;
            leafHash?: Uint8Array | undefined;
            aunts?: Uint8Array[] | undefined;
        } & {
            total?: bigint | undefined;
            index?: bigint | undefined;
            leafHash?: Uint8Array | undefined;
            aunts?: (Uint8Array[] & Uint8Array[] & Record<Exclude<keyof I["proof"]["aunts"], keyof Uint8Array[]>, never>) | undefined;
        } & Record<Exclude<keyof I["proof"], keyof Proof>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof ValueOp>, never>>(object: I): ValueOp;
};
export declare const DominoOp: {
    typeUrl: string;
    encode(message: DominoOp, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): DominoOp;
    fromJSON(object: any): DominoOp;
    toJSON(message: DominoOp): unknown;
    fromPartial<I extends {
        key?: string | undefined;
        input?: string | undefined;
        output?: string | undefined;
    } & {
        key?: string | undefined;
        input?: string | undefined;
        output?: string | undefined;
    } & Record<Exclude<keyof I, keyof DominoOp>, never>>(object: I): DominoOp;
};
export declare const ProofOp: {
    typeUrl: string;
    encode(message: ProofOp, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ProofOp;
    fromJSON(object: any): ProofOp;
    toJSON(message: ProofOp): unknown;
    fromPartial<I extends {
        type?: string | undefined;
        key?: Uint8Array | undefined;
        data?: Uint8Array | undefined;
    } & {
        type?: string | undefined;
        key?: Uint8Array | undefined;
        data?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof ProofOp>, never>>(object: I): ProofOp;
};
export declare const ProofOps: {
    typeUrl: string;
    encode(message: ProofOps, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ProofOps;
    fromJSON(object: any): ProofOps;
    toJSON(message: ProofOps): unknown;
    fromPartial<I extends {
        ops?: {
            type?: string | undefined;
            key?: Uint8Array | undefined;
            data?: Uint8Array | undefined;
        }[] | undefined;
    } & {
        ops?: ({
            type?: string | undefined;
            key?: Uint8Array | undefined;
            data?: Uint8Array | undefined;
        }[] & ({
            type?: string | undefined;
            key?: Uint8Array | undefined;
            data?: Uint8Array | undefined;
        } & {
            type?: string | undefined;
            key?: Uint8Array | undefined;
            data?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["ops"][number], keyof ProofOp>, never>)[] & Record<Exclude<keyof I["ops"], keyof {
            type?: string | undefined;
            key?: Uint8Array | undefined;
            data?: Uint8Array | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "ops">, never>>(object: I): ProofOps;
};
