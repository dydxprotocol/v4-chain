import { Duration } from "../../google/protobuf/duration";
import { BinaryReader, BinaryWriter } from "../../binary";
export declare const protobufPackage = "tendermint.types";
/**
 * ConsensusParams contains consensus critical parameters that determine the
 * validity of blocks.
 */
export interface ConsensusParams {
    block?: BlockParams;
    evidence?: EvidenceParams;
    validator?: ValidatorParams;
    version?: VersionParams;
}
/** BlockParams contains limits on the block size. */
export interface BlockParams {
    /**
     * Max block size, in bytes.
     * Note: must be greater than 0
     */
    maxBytes: bigint;
    /**
     * Max gas per block.
     * Note: must be greater or equal to -1
     */
    maxGas: bigint;
}
/** EvidenceParams determine how we handle evidence of malfeasance. */
export interface EvidenceParams {
    /**
     * Max age of evidence, in blocks.
     *
     * The basic formula for calculating this is: MaxAgeDuration / {average block
     * time}.
     */
    maxAgeNumBlocks: bigint;
    /**
     * Max age of evidence, in time.
     *
     * It should correspond with an app's "unbonding period" or other similar
     * mechanism for handling [Nothing-At-Stake
     * attacks](https://github.com/ethereum/wiki/wiki/Proof-of-Stake-FAQ#what-is-the-nothing-at-stake-problem-and-how-can-it-be-fixed).
     */
    maxAgeDuration: Duration;
    /**
     * This sets the maximum size of total evidence in bytes that can be committed in a single block.
     * and should fall comfortably under the max block bytes.
     * Default is 1048576 or 1MB
     */
    maxBytes: bigint;
}
/**
 * ValidatorParams restrict the public key types validators can use.
 * NOTE: uses ABCI pubkey naming, not Amino names.
 */
export interface ValidatorParams {
    pubKeyTypes: string[];
}
/** VersionParams contains the ABCI application version. */
export interface VersionParams {
    app: bigint;
}
/**
 * HashedParams is a subset of ConsensusParams.
 *
 * It is hashed into the Header.ConsensusHash.
 */
export interface HashedParams {
    blockMaxBytes: bigint;
    blockMaxGas: bigint;
}
export declare const ConsensusParams: {
    typeUrl: string;
    encode(message: ConsensusParams, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ConsensusParams;
    fromJSON(object: any): ConsensusParams;
    toJSON(message: ConsensusParams): unknown;
    fromPartial<I extends {
        block?: {
            maxBytes?: bigint | undefined;
            maxGas?: bigint | undefined;
        } | undefined;
        evidence?: {
            maxAgeNumBlocks?: bigint | undefined;
            maxAgeDuration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            maxBytes?: bigint | undefined;
        } | undefined;
        validator?: {
            pubKeyTypes?: string[] | undefined;
        } | undefined;
        version?: {
            app?: bigint | undefined;
        } | undefined;
    } & {
        block?: ({
            maxBytes?: bigint | undefined;
            maxGas?: bigint | undefined;
        } & {
            maxBytes?: bigint | undefined;
            maxGas?: bigint | undefined;
        } & Record<Exclude<keyof I["block"], keyof BlockParams>, never>) | undefined;
        evidence?: ({
            maxAgeNumBlocks?: bigint | undefined;
            maxAgeDuration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            maxBytes?: bigint | undefined;
        } & {
            maxAgeNumBlocks?: bigint | undefined;
            maxAgeDuration?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["evidence"]["maxAgeDuration"], keyof Duration>, never>) | undefined;
            maxBytes?: bigint | undefined;
        } & Record<Exclude<keyof I["evidence"], keyof EvidenceParams>, never>) | undefined;
        validator?: ({
            pubKeyTypes?: string[] | undefined;
        } & {
            pubKeyTypes?: (string[] & string[] & Record<Exclude<keyof I["validator"]["pubKeyTypes"], keyof string[]>, never>) | undefined;
        } & Record<Exclude<keyof I["validator"], "pubKeyTypes">, never>) | undefined;
        version?: ({
            app?: bigint | undefined;
        } & {
            app?: bigint | undefined;
        } & Record<Exclude<keyof I["version"], "app">, never>) | undefined;
    } & Record<Exclude<keyof I, keyof ConsensusParams>, never>>(object: I): ConsensusParams;
};
export declare const BlockParams: {
    typeUrl: string;
    encode(message: BlockParams, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): BlockParams;
    fromJSON(object: any): BlockParams;
    toJSON(message: BlockParams): unknown;
    fromPartial<I extends {
        maxBytes?: bigint | undefined;
        maxGas?: bigint | undefined;
    } & {
        maxBytes?: bigint | undefined;
        maxGas?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof BlockParams>, never>>(object: I): BlockParams;
};
export declare const EvidenceParams: {
    typeUrl: string;
    encode(message: EvidenceParams, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): EvidenceParams;
    fromJSON(object: any): EvidenceParams;
    toJSON(message: EvidenceParams): unknown;
    fromPartial<I extends {
        maxAgeNumBlocks?: bigint | undefined;
        maxAgeDuration?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
        maxBytes?: bigint | undefined;
    } & {
        maxAgeNumBlocks?: bigint | undefined;
        maxAgeDuration?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["maxAgeDuration"], keyof Duration>, never>) | undefined;
        maxBytes?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof EvidenceParams>, never>>(object: I): EvidenceParams;
};
export declare const ValidatorParams: {
    typeUrl: string;
    encode(message: ValidatorParams, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ValidatorParams;
    fromJSON(object: any): ValidatorParams;
    toJSON(message: ValidatorParams): unknown;
    fromPartial<I extends {
        pubKeyTypes?: string[] | undefined;
    } & {
        pubKeyTypes?: (string[] & string[] & Record<Exclude<keyof I["pubKeyTypes"], keyof string[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "pubKeyTypes">, never>>(object: I): ValidatorParams;
};
export declare const VersionParams: {
    typeUrl: string;
    encode(message: VersionParams, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): VersionParams;
    fromJSON(object: any): VersionParams;
    toJSON(message: VersionParams): unknown;
    fromPartial<I extends {
        app?: bigint | undefined;
    } & {
        app?: bigint | undefined;
    } & Record<Exclude<keyof I, "app">, never>>(object: I): VersionParams;
};
export declare const HashedParams: {
    typeUrl: string;
    encode(message: HashedParams, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): HashedParams;
    fromJSON(object: any): HashedParams;
    toJSON(message: HashedParams): unknown;
    fromPartial<I extends {
        blockMaxBytes?: bigint | undefined;
        blockMaxGas?: bigint | undefined;
    } & {
        blockMaxBytes?: bigint | undefined;
        blockMaxGas?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof HashedParams>, never>>(object: I): HashedParams;
};
