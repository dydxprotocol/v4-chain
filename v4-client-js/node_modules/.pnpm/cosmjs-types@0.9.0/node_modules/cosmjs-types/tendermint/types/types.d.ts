import { Proof } from "../crypto/proof";
import { Consensus } from "../version/types";
import { Timestamp } from "../../google/protobuf/timestamp";
import { ValidatorSet } from "./validator";
import { BinaryReader, BinaryWriter } from "../../binary";
export declare const protobufPackage = "tendermint.types";
/** BlockIdFlag indicates which BlcokID the signature is for */
export declare enum BlockIDFlag {
    BLOCK_ID_FLAG_UNKNOWN = 0,
    BLOCK_ID_FLAG_ABSENT = 1,
    BLOCK_ID_FLAG_COMMIT = 2,
    BLOCK_ID_FLAG_NIL = 3,
    UNRECOGNIZED = -1
}
export declare function blockIDFlagFromJSON(object: any): BlockIDFlag;
export declare function blockIDFlagToJSON(object: BlockIDFlag): string;
/** SignedMsgType is a type of signed message in the consensus. */
export declare enum SignedMsgType {
    SIGNED_MSG_TYPE_UNKNOWN = 0,
    /** SIGNED_MSG_TYPE_PREVOTE - Votes */
    SIGNED_MSG_TYPE_PREVOTE = 1,
    SIGNED_MSG_TYPE_PRECOMMIT = 2,
    /** SIGNED_MSG_TYPE_PROPOSAL - Proposals */
    SIGNED_MSG_TYPE_PROPOSAL = 32,
    UNRECOGNIZED = -1
}
export declare function signedMsgTypeFromJSON(object: any): SignedMsgType;
export declare function signedMsgTypeToJSON(object: SignedMsgType): string;
/** PartsetHeader */
export interface PartSetHeader {
    total: number;
    hash: Uint8Array;
}
export interface Part {
    index: number;
    bytes: Uint8Array;
    proof: Proof;
}
/** BlockID */
export interface BlockID {
    hash: Uint8Array;
    partSetHeader: PartSetHeader;
}
/** Header defines the structure of a block header. */
export interface Header {
    /** basic block info */
    version: Consensus;
    chainId: string;
    height: bigint;
    time: Timestamp;
    /** prev block info */
    lastBlockId: BlockID;
    /** hashes of block data */
    lastCommitHash: Uint8Array;
    dataHash: Uint8Array;
    /** hashes from the app output from the prev block */
    validatorsHash: Uint8Array;
    /** validators for the next block */
    nextValidatorsHash: Uint8Array;
    /** consensus params for current block */
    consensusHash: Uint8Array;
    /** state after txs from the previous block */
    appHash: Uint8Array;
    lastResultsHash: Uint8Array;
    /** consensus info */
    evidenceHash: Uint8Array;
    /** original proposer of the block */
    proposerAddress: Uint8Array;
}
/** Data contains the set of transactions included in the block */
export interface Data {
    /**
     * Txs that will be applied by state @ block.Height+1.
     * NOTE: not all txs here are valid.  We're just agreeing on the order first.
     * This means that block.AppHash does not include these txs.
     */
    txs: Uint8Array[];
}
/**
 * Vote represents a prevote, precommit, or commit vote from validators for
 * consensus.
 */
export interface Vote {
    type: SignedMsgType;
    height: bigint;
    round: number;
    blockId: BlockID;
    timestamp: Timestamp;
    validatorAddress: Uint8Array;
    validatorIndex: number;
    signature: Uint8Array;
}
/** Commit contains the evidence that a block was committed by a set of validators. */
export interface Commit {
    height: bigint;
    round: number;
    blockId: BlockID;
    signatures: CommitSig[];
}
/** CommitSig is a part of the Vote included in a Commit. */
export interface CommitSig {
    blockIdFlag: BlockIDFlag;
    validatorAddress: Uint8Array;
    timestamp: Timestamp;
    signature: Uint8Array;
}
export interface Proposal {
    type: SignedMsgType;
    height: bigint;
    round: number;
    polRound: number;
    blockId: BlockID;
    timestamp: Timestamp;
    signature: Uint8Array;
}
export interface SignedHeader {
    header?: Header;
    commit?: Commit;
}
export interface LightBlock {
    signedHeader?: SignedHeader;
    validatorSet?: ValidatorSet;
}
export interface BlockMeta {
    blockId: BlockID;
    blockSize: bigint;
    header: Header;
    numTxs: bigint;
}
/** TxProof represents a Merkle proof of the presence of a transaction in the Merkle tree. */
export interface TxProof {
    rootHash: Uint8Array;
    data: Uint8Array;
    proof?: Proof;
}
export declare const PartSetHeader: {
    typeUrl: string;
    encode(message: PartSetHeader, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): PartSetHeader;
    fromJSON(object: any): PartSetHeader;
    toJSON(message: PartSetHeader): unknown;
    fromPartial<I extends {
        total?: number | undefined;
        hash?: Uint8Array | undefined;
    } & {
        total?: number | undefined;
        hash?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof PartSetHeader>, never>>(object: I): PartSetHeader;
};
export declare const Part: {
    typeUrl: string;
    encode(message: Part, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Part;
    fromJSON(object: any): Part;
    toJSON(message: Part): unknown;
    fromPartial<I extends {
        index?: number | undefined;
        bytes?: Uint8Array | undefined;
        proof?: {
            total?: bigint | undefined;
            index?: bigint | undefined;
            leafHash?: Uint8Array | undefined;
            aunts?: Uint8Array[] | undefined;
        } | undefined;
    } & {
        index?: number | undefined;
        bytes?: Uint8Array | undefined;
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
    } & Record<Exclude<keyof I, keyof Part>, never>>(object: I): Part;
};
export declare const BlockID: {
    typeUrl: string;
    encode(message: BlockID, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): BlockID;
    fromJSON(object: any): BlockID;
    toJSON(message: BlockID): unknown;
    fromPartial<I extends {
        hash?: Uint8Array | undefined;
        partSetHeader?: {
            total?: number | undefined;
            hash?: Uint8Array | undefined;
        } | undefined;
    } & {
        hash?: Uint8Array | undefined;
        partSetHeader?: ({
            total?: number | undefined;
            hash?: Uint8Array | undefined;
        } & {
            total?: number | undefined;
            hash?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["partSetHeader"], keyof PartSetHeader>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof BlockID>, never>>(object: I): BlockID;
};
export declare const Header: {
    typeUrl: string;
    encode(message: Header, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Header;
    fromJSON(object: any): Header;
    toJSON(message: Header): unknown;
    fromPartial<I extends {
        version?: {
            block?: bigint | undefined;
            app?: bigint | undefined;
        } | undefined;
        chainId?: string | undefined;
        height?: bigint | undefined;
        time?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
        lastBlockId?: {
            hash?: Uint8Array | undefined;
            partSetHeader?: {
                total?: number | undefined;
                hash?: Uint8Array | undefined;
            } | undefined;
        } | undefined;
        lastCommitHash?: Uint8Array | undefined;
        dataHash?: Uint8Array | undefined;
        validatorsHash?: Uint8Array | undefined;
        nextValidatorsHash?: Uint8Array | undefined;
        consensusHash?: Uint8Array | undefined;
        appHash?: Uint8Array | undefined;
        lastResultsHash?: Uint8Array | undefined;
        evidenceHash?: Uint8Array | undefined;
        proposerAddress?: Uint8Array | undefined;
    } & {
        version?: ({
            block?: bigint | undefined;
            app?: bigint | undefined;
        } & {
            block?: bigint | undefined;
            app?: bigint | undefined;
        } & Record<Exclude<keyof I["version"], keyof Consensus>, never>) | undefined;
        chainId?: string | undefined;
        height?: bigint | undefined;
        time?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["time"], keyof Timestamp>, never>) | undefined;
        lastBlockId?: ({
            hash?: Uint8Array | undefined;
            partSetHeader?: {
                total?: number | undefined;
                hash?: Uint8Array | undefined;
            } | undefined;
        } & {
            hash?: Uint8Array | undefined;
            partSetHeader?: ({
                total?: number | undefined;
                hash?: Uint8Array | undefined;
            } & {
                total?: number | undefined;
                hash?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["lastBlockId"]["partSetHeader"], keyof PartSetHeader>, never>) | undefined;
        } & Record<Exclude<keyof I["lastBlockId"], keyof BlockID>, never>) | undefined;
        lastCommitHash?: Uint8Array | undefined;
        dataHash?: Uint8Array | undefined;
        validatorsHash?: Uint8Array | undefined;
        nextValidatorsHash?: Uint8Array | undefined;
        consensusHash?: Uint8Array | undefined;
        appHash?: Uint8Array | undefined;
        lastResultsHash?: Uint8Array | undefined;
        evidenceHash?: Uint8Array | undefined;
        proposerAddress?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof Header>, never>>(object: I): Header;
};
export declare const Data: {
    typeUrl: string;
    encode(message: Data, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Data;
    fromJSON(object: any): Data;
    toJSON(message: Data): unknown;
    fromPartial<I extends {
        txs?: Uint8Array[] | undefined;
    } & {
        txs?: (Uint8Array[] & Uint8Array[] & Record<Exclude<keyof I["txs"], keyof Uint8Array[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "txs">, never>>(object: I): Data;
};
export declare const Vote: {
    typeUrl: string;
    encode(message: Vote, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Vote;
    fromJSON(object: any): Vote;
    toJSON(message: Vote): unknown;
    fromPartial<I extends {
        type?: SignedMsgType | undefined;
        height?: bigint | undefined;
        round?: number | undefined;
        blockId?: {
            hash?: Uint8Array | undefined;
            partSetHeader?: {
                total?: number | undefined;
                hash?: Uint8Array | undefined;
            } | undefined;
        } | undefined;
        timestamp?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
        validatorAddress?: Uint8Array | undefined;
        validatorIndex?: number | undefined;
        signature?: Uint8Array | undefined;
    } & {
        type?: SignedMsgType | undefined;
        height?: bigint | undefined;
        round?: number | undefined;
        blockId?: ({
            hash?: Uint8Array | undefined;
            partSetHeader?: {
                total?: number | undefined;
                hash?: Uint8Array | undefined;
            } | undefined;
        } & {
            hash?: Uint8Array | undefined;
            partSetHeader?: ({
                total?: number | undefined;
                hash?: Uint8Array | undefined;
            } & {
                total?: number | undefined;
                hash?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["blockId"]["partSetHeader"], keyof PartSetHeader>, never>) | undefined;
        } & Record<Exclude<keyof I["blockId"], keyof BlockID>, never>) | undefined;
        timestamp?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["timestamp"], keyof Timestamp>, never>) | undefined;
        validatorAddress?: Uint8Array | undefined;
        validatorIndex?: number | undefined;
        signature?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof Vote>, never>>(object: I): Vote;
};
export declare const Commit: {
    typeUrl: string;
    encode(message: Commit, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Commit;
    fromJSON(object: any): Commit;
    toJSON(message: Commit): unknown;
    fromPartial<I extends {
        height?: bigint | undefined;
        round?: number | undefined;
        blockId?: {
            hash?: Uint8Array | undefined;
            partSetHeader?: {
                total?: number | undefined;
                hash?: Uint8Array | undefined;
            } | undefined;
        } | undefined;
        signatures?: {
            blockIdFlag?: BlockIDFlag | undefined;
            validatorAddress?: Uint8Array | undefined;
            timestamp?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            signature?: Uint8Array | undefined;
        }[] | undefined;
    } & {
        height?: bigint | undefined;
        round?: number | undefined;
        blockId?: ({
            hash?: Uint8Array | undefined;
            partSetHeader?: {
                total?: number | undefined;
                hash?: Uint8Array | undefined;
            } | undefined;
        } & {
            hash?: Uint8Array | undefined;
            partSetHeader?: ({
                total?: number | undefined;
                hash?: Uint8Array | undefined;
            } & {
                total?: number | undefined;
                hash?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["blockId"]["partSetHeader"], keyof PartSetHeader>, never>) | undefined;
        } & Record<Exclude<keyof I["blockId"], keyof BlockID>, never>) | undefined;
        signatures?: ({
            blockIdFlag?: BlockIDFlag | undefined;
            validatorAddress?: Uint8Array | undefined;
            timestamp?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            signature?: Uint8Array | undefined;
        }[] & ({
            blockIdFlag?: BlockIDFlag | undefined;
            validatorAddress?: Uint8Array | undefined;
            timestamp?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            signature?: Uint8Array | undefined;
        } & {
            blockIdFlag?: BlockIDFlag | undefined;
            validatorAddress?: Uint8Array | undefined;
            timestamp?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["signatures"][number]["timestamp"], keyof Timestamp>, never>) | undefined;
            signature?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["signatures"][number], keyof CommitSig>, never>)[] & Record<Exclude<keyof I["signatures"], keyof {
            blockIdFlag?: BlockIDFlag | undefined;
            validatorAddress?: Uint8Array | undefined;
            timestamp?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            signature?: Uint8Array | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof Commit>, never>>(object: I): Commit;
};
export declare const CommitSig: {
    typeUrl: string;
    encode(message: CommitSig, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): CommitSig;
    fromJSON(object: any): CommitSig;
    toJSON(message: CommitSig): unknown;
    fromPartial<I extends {
        blockIdFlag?: BlockIDFlag | undefined;
        validatorAddress?: Uint8Array | undefined;
        timestamp?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
        signature?: Uint8Array | undefined;
    } & {
        blockIdFlag?: BlockIDFlag | undefined;
        validatorAddress?: Uint8Array | undefined;
        timestamp?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["timestamp"], keyof Timestamp>, never>) | undefined;
        signature?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof CommitSig>, never>>(object: I): CommitSig;
};
export declare const Proposal: {
    typeUrl: string;
    encode(message: Proposal, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Proposal;
    fromJSON(object: any): Proposal;
    toJSON(message: Proposal): unknown;
    fromPartial<I extends {
        type?: SignedMsgType | undefined;
        height?: bigint | undefined;
        round?: number | undefined;
        polRound?: number | undefined;
        blockId?: {
            hash?: Uint8Array | undefined;
            partSetHeader?: {
                total?: number | undefined;
                hash?: Uint8Array | undefined;
            } | undefined;
        } | undefined;
        timestamp?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
        signature?: Uint8Array | undefined;
    } & {
        type?: SignedMsgType | undefined;
        height?: bigint | undefined;
        round?: number | undefined;
        polRound?: number | undefined;
        blockId?: ({
            hash?: Uint8Array | undefined;
            partSetHeader?: {
                total?: number | undefined;
                hash?: Uint8Array | undefined;
            } | undefined;
        } & {
            hash?: Uint8Array | undefined;
            partSetHeader?: ({
                total?: number | undefined;
                hash?: Uint8Array | undefined;
            } & {
                total?: number | undefined;
                hash?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["blockId"]["partSetHeader"], keyof PartSetHeader>, never>) | undefined;
        } & Record<Exclude<keyof I["blockId"], keyof BlockID>, never>) | undefined;
        timestamp?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["timestamp"], keyof Timestamp>, never>) | undefined;
        signature?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof Proposal>, never>>(object: I): Proposal;
};
export declare const SignedHeader: {
    typeUrl: string;
    encode(message: SignedHeader, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): SignedHeader;
    fromJSON(object: any): SignedHeader;
    toJSON(message: SignedHeader): unknown;
    fromPartial<I extends {
        header?: {
            version?: {
                block?: bigint | undefined;
                app?: bigint | undefined;
            } | undefined;
            chainId?: string | undefined;
            height?: bigint | undefined;
            time?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            lastBlockId?: {
                hash?: Uint8Array | undefined;
                partSetHeader?: {
                    total?: number | undefined;
                    hash?: Uint8Array | undefined;
                } | undefined;
            } | undefined;
            lastCommitHash?: Uint8Array | undefined;
            dataHash?: Uint8Array | undefined;
            validatorsHash?: Uint8Array | undefined;
            nextValidatorsHash?: Uint8Array | undefined;
            consensusHash?: Uint8Array | undefined;
            appHash?: Uint8Array | undefined;
            lastResultsHash?: Uint8Array | undefined;
            evidenceHash?: Uint8Array | undefined;
            proposerAddress?: Uint8Array | undefined;
        } | undefined;
        commit?: {
            height?: bigint | undefined;
            round?: number | undefined;
            blockId?: {
                hash?: Uint8Array | undefined;
                partSetHeader?: {
                    total?: number | undefined;
                    hash?: Uint8Array | undefined;
                } | undefined;
            } | undefined;
            signatures?: {
                blockIdFlag?: BlockIDFlag | undefined;
                validatorAddress?: Uint8Array | undefined;
                timestamp?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                signature?: Uint8Array | undefined;
            }[] | undefined;
        } | undefined;
    } & {
        header?: ({
            version?: {
                block?: bigint | undefined;
                app?: bigint | undefined;
            } | undefined;
            chainId?: string | undefined;
            height?: bigint | undefined;
            time?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            lastBlockId?: {
                hash?: Uint8Array | undefined;
                partSetHeader?: {
                    total?: number | undefined;
                    hash?: Uint8Array | undefined;
                } | undefined;
            } | undefined;
            lastCommitHash?: Uint8Array | undefined;
            dataHash?: Uint8Array | undefined;
            validatorsHash?: Uint8Array | undefined;
            nextValidatorsHash?: Uint8Array | undefined;
            consensusHash?: Uint8Array | undefined;
            appHash?: Uint8Array | undefined;
            lastResultsHash?: Uint8Array | undefined;
            evidenceHash?: Uint8Array | undefined;
            proposerAddress?: Uint8Array | undefined;
        } & {
            version?: ({
                block?: bigint | undefined;
                app?: bigint | undefined;
            } & {
                block?: bigint | undefined;
                app?: bigint | undefined;
            } & Record<Exclude<keyof I["header"]["version"], keyof Consensus>, never>) | undefined;
            chainId?: string | undefined;
            height?: bigint | undefined;
            time?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["header"]["time"], keyof Timestamp>, never>) | undefined;
            lastBlockId?: ({
                hash?: Uint8Array | undefined;
                partSetHeader?: {
                    total?: number | undefined;
                    hash?: Uint8Array | undefined;
                } | undefined;
            } & {
                hash?: Uint8Array | undefined;
                partSetHeader?: ({
                    total?: number | undefined;
                    hash?: Uint8Array | undefined;
                } & {
                    total?: number | undefined;
                    hash?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["header"]["lastBlockId"]["partSetHeader"], keyof PartSetHeader>, never>) | undefined;
            } & Record<Exclude<keyof I["header"]["lastBlockId"], keyof BlockID>, never>) | undefined;
            lastCommitHash?: Uint8Array | undefined;
            dataHash?: Uint8Array | undefined;
            validatorsHash?: Uint8Array | undefined;
            nextValidatorsHash?: Uint8Array | undefined;
            consensusHash?: Uint8Array | undefined;
            appHash?: Uint8Array | undefined;
            lastResultsHash?: Uint8Array | undefined;
            evidenceHash?: Uint8Array | undefined;
            proposerAddress?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["header"], keyof Header>, never>) | undefined;
        commit?: ({
            height?: bigint | undefined;
            round?: number | undefined;
            blockId?: {
                hash?: Uint8Array | undefined;
                partSetHeader?: {
                    total?: number | undefined;
                    hash?: Uint8Array | undefined;
                } | undefined;
            } | undefined;
            signatures?: {
                blockIdFlag?: BlockIDFlag | undefined;
                validatorAddress?: Uint8Array | undefined;
                timestamp?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                signature?: Uint8Array | undefined;
            }[] | undefined;
        } & {
            height?: bigint | undefined;
            round?: number | undefined;
            blockId?: ({
                hash?: Uint8Array | undefined;
                partSetHeader?: {
                    total?: number | undefined;
                    hash?: Uint8Array | undefined;
                } | undefined;
            } & {
                hash?: Uint8Array | undefined;
                partSetHeader?: ({
                    total?: number | undefined;
                    hash?: Uint8Array | undefined;
                } & {
                    total?: number | undefined;
                    hash?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["commit"]["blockId"]["partSetHeader"], keyof PartSetHeader>, never>) | undefined;
            } & Record<Exclude<keyof I["commit"]["blockId"], keyof BlockID>, never>) | undefined;
            signatures?: ({
                blockIdFlag?: BlockIDFlag | undefined;
                validatorAddress?: Uint8Array | undefined;
                timestamp?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                signature?: Uint8Array | undefined;
            }[] & ({
                blockIdFlag?: BlockIDFlag | undefined;
                validatorAddress?: Uint8Array | undefined;
                timestamp?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                signature?: Uint8Array | undefined;
            } & {
                blockIdFlag?: BlockIDFlag | undefined;
                validatorAddress?: Uint8Array | undefined;
                timestamp?: ({
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & Record<Exclude<keyof I["commit"]["signatures"][number]["timestamp"], keyof Timestamp>, never>) | undefined;
                signature?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["commit"]["signatures"][number], keyof CommitSig>, never>)[] & Record<Exclude<keyof I["commit"]["signatures"], keyof {
                blockIdFlag?: BlockIDFlag | undefined;
                validatorAddress?: Uint8Array | undefined;
                timestamp?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                signature?: Uint8Array | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["commit"], keyof Commit>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof SignedHeader>, never>>(object: I): SignedHeader;
};
export declare const LightBlock: {
    typeUrl: string;
    encode(message: LightBlock, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): LightBlock;
    fromJSON(object: any): LightBlock;
    toJSON(message: LightBlock): unknown;
    fromPartial<I extends {
        signedHeader?: {
            header?: {
                version?: {
                    block?: bigint | undefined;
                    app?: bigint | undefined;
                } | undefined;
                chainId?: string | undefined;
                height?: bigint | undefined;
                time?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                lastBlockId?: {
                    hash?: Uint8Array | undefined;
                    partSetHeader?: {
                        total?: number | undefined;
                        hash?: Uint8Array | undefined;
                    } | undefined;
                } | undefined;
                lastCommitHash?: Uint8Array | undefined;
                dataHash?: Uint8Array | undefined;
                validatorsHash?: Uint8Array | undefined;
                nextValidatorsHash?: Uint8Array | undefined;
                consensusHash?: Uint8Array | undefined;
                appHash?: Uint8Array | undefined;
                lastResultsHash?: Uint8Array | undefined;
                evidenceHash?: Uint8Array | undefined;
                proposerAddress?: Uint8Array | undefined;
            } | undefined;
            commit?: {
                height?: bigint | undefined;
                round?: number | undefined;
                blockId?: {
                    hash?: Uint8Array | undefined;
                    partSetHeader?: {
                        total?: number | undefined;
                        hash?: Uint8Array | undefined;
                    } | undefined;
                } | undefined;
                signatures?: {
                    blockIdFlag?: BlockIDFlag | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    timestamp?: {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } | undefined;
                    signature?: Uint8Array | undefined;
                }[] | undefined;
            } | undefined;
        } | undefined;
        validatorSet?: {
            validators?: {
                address?: Uint8Array | undefined;
                pubKey?: {
                    ed25519?: Uint8Array | undefined;
                    secp256k1?: Uint8Array | undefined;
                } | undefined;
                votingPower?: bigint | undefined;
                proposerPriority?: bigint | undefined;
            }[] | undefined;
            proposer?: {
                address?: Uint8Array | undefined;
                pubKey?: {
                    ed25519?: Uint8Array | undefined;
                    secp256k1?: Uint8Array | undefined;
                } | undefined;
                votingPower?: bigint | undefined;
                proposerPriority?: bigint | undefined;
            } | undefined;
            totalVotingPower?: bigint | undefined;
        } | undefined;
    } & {
        signedHeader?: ({
            header?: {
                version?: {
                    block?: bigint | undefined;
                    app?: bigint | undefined;
                } | undefined;
                chainId?: string | undefined;
                height?: bigint | undefined;
                time?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                lastBlockId?: {
                    hash?: Uint8Array | undefined;
                    partSetHeader?: {
                        total?: number | undefined;
                        hash?: Uint8Array | undefined;
                    } | undefined;
                } | undefined;
                lastCommitHash?: Uint8Array | undefined;
                dataHash?: Uint8Array | undefined;
                validatorsHash?: Uint8Array | undefined;
                nextValidatorsHash?: Uint8Array | undefined;
                consensusHash?: Uint8Array | undefined;
                appHash?: Uint8Array | undefined;
                lastResultsHash?: Uint8Array | undefined;
                evidenceHash?: Uint8Array | undefined;
                proposerAddress?: Uint8Array | undefined;
            } | undefined;
            commit?: {
                height?: bigint | undefined;
                round?: number | undefined;
                blockId?: {
                    hash?: Uint8Array | undefined;
                    partSetHeader?: {
                        total?: number | undefined;
                        hash?: Uint8Array | undefined;
                    } | undefined;
                } | undefined;
                signatures?: {
                    blockIdFlag?: BlockIDFlag | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    timestamp?: {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } | undefined;
                    signature?: Uint8Array | undefined;
                }[] | undefined;
            } | undefined;
        } & {
            header?: ({
                version?: {
                    block?: bigint | undefined;
                    app?: bigint | undefined;
                } | undefined;
                chainId?: string | undefined;
                height?: bigint | undefined;
                time?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                lastBlockId?: {
                    hash?: Uint8Array | undefined;
                    partSetHeader?: {
                        total?: number | undefined;
                        hash?: Uint8Array | undefined;
                    } | undefined;
                } | undefined;
                lastCommitHash?: Uint8Array | undefined;
                dataHash?: Uint8Array | undefined;
                validatorsHash?: Uint8Array | undefined;
                nextValidatorsHash?: Uint8Array | undefined;
                consensusHash?: Uint8Array | undefined;
                appHash?: Uint8Array | undefined;
                lastResultsHash?: Uint8Array | undefined;
                evidenceHash?: Uint8Array | undefined;
                proposerAddress?: Uint8Array | undefined;
            } & {
                version?: ({
                    block?: bigint | undefined;
                    app?: bigint | undefined;
                } & {
                    block?: bigint | undefined;
                    app?: bigint | undefined;
                } & Record<Exclude<keyof I["signedHeader"]["header"]["version"], keyof Consensus>, never>) | undefined;
                chainId?: string | undefined;
                height?: bigint | undefined;
                time?: ({
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & Record<Exclude<keyof I["signedHeader"]["header"]["time"], keyof Timestamp>, never>) | undefined;
                lastBlockId?: ({
                    hash?: Uint8Array | undefined;
                    partSetHeader?: {
                        total?: number | undefined;
                        hash?: Uint8Array | undefined;
                    } | undefined;
                } & {
                    hash?: Uint8Array | undefined;
                    partSetHeader?: ({
                        total?: number | undefined;
                        hash?: Uint8Array | undefined;
                    } & {
                        total?: number | undefined;
                        hash?: Uint8Array | undefined;
                    } & Record<Exclude<keyof I["signedHeader"]["header"]["lastBlockId"]["partSetHeader"], keyof PartSetHeader>, never>) | undefined;
                } & Record<Exclude<keyof I["signedHeader"]["header"]["lastBlockId"], keyof BlockID>, never>) | undefined;
                lastCommitHash?: Uint8Array | undefined;
                dataHash?: Uint8Array | undefined;
                validatorsHash?: Uint8Array | undefined;
                nextValidatorsHash?: Uint8Array | undefined;
                consensusHash?: Uint8Array | undefined;
                appHash?: Uint8Array | undefined;
                lastResultsHash?: Uint8Array | undefined;
                evidenceHash?: Uint8Array | undefined;
                proposerAddress?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["signedHeader"]["header"], keyof Header>, never>) | undefined;
            commit?: ({
                height?: bigint | undefined;
                round?: number | undefined;
                blockId?: {
                    hash?: Uint8Array | undefined;
                    partSetHeader?: {
                        total?: number | undefined;
                        hash?: Uint8Array | undefined;
                    } | undefined;
                } | undefined;
                signatures?: {
                    blockIdFlag?: BlockIDFlag | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    timestamp?: {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } | undefined;
                    signature?: Uint8Array | undefined;
                }[] | undefined;
            } & {
                height?: bigint | undefined;
                round?: number | undefined;
                blockId?: ({
                    hash?: Uint8Array | undefined;
                    partSetHeader?: {
                        total?: number | undefined;
                        hash?: Uint8Array | undefined;
                    } | undefined;
                } & {
                    hash?: Uint8Array | undefined;
                    partSetHeader?: ({
                        total?: number | undefined;
                        hash?: Uint8Array | undefined;
                    } & {
                        total?: number | undefined;
                        hash?: Uint8Array | undefined;
                    } & Record<Exclude<keyof I["signedHeader"]["commit"]["blockId"]["partSetHeader"], keyof PartSetHeader>, never>) | undefined;
                } & Record<Exclude<keyof I["signedHeader"]["commit"]["blockId"], keyof BlockID>, never>) | undefined;
                signatures?: ({
                    blockIdFlag?: BlockIDFlag | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    timestamp?: {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } | undefined;
                    signature?: Uint8Array | undefined;
                }[] & ({
                    blockIdFlag?: BlockIDFlag | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    timestamp?: {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } | undefined;
                    signature?: Uint8Array | undefined;
                } & {
                    blockIdFlag?: BlockIDFlag | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    timestamp?: ({
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } & {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } & Record<Exclude<keyof I["signedHeader"]["commit"]["signatures"][number]["timestamp"], keyof Timestamp>, never>) | undefined;
                    signature?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["signedHeader"]["commit"]["signatures"][number], keyof CommitSig>, never>)[] & Record<Exclude<keyof I["signedHeader"]["commit"]["signatures"], keyof {
                    blockIdFlag?: BlockIDFlag | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    timestamp?: {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } | undefined;
                    signature?: Uint8Array | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["signedHeader"]["commit"], keyof Commit>, never>) | undefined;
        } & Record<Exclude<keyof I["signedHeader"], keyof SignedHeader>, never>) | undefined;
        validatorSet?: ({
            validators?: {
                address?: Uint8Array | undefined;
                pubKey?: {
                    ed25519?: Uint8Array | undefined;
                    secp256k1?: Uint8Array | undefined;
                } | undefined;
                votingPower?: bigint | undefined;
                proposerPriority?: bigint | undefined;
            }[] | undefined;
            proposer?: {
                address?: Uint8Array | undefined;
                pubKey?: {
                    ed25519?: Uint8Array | undefined;
                    secp256k1?: Uint8Array | undefined;
                } | undefined;
                votingPower?: bigint | undefined;
                proposerPriority?: bigint | undefined;
            } | undefined;
            totalVotingPower?: bigint | undefined;
        } & {
            validators?: ({
                address?: Uint8Array | undefined;
                pubKey?: {
                    ed25519?: Uint8Array | undefined;
                    secp256k1?: Uint8Array | undefined;
                } | undefined;
                votingPower?: bigint | undefined;
                proposerPriority?: bigint | undefined;
            }[] & ({
                address?: Uint8Array | undefined;
                pubKey?: {
                    ed25519?: Uint8Array | undefined;
                    secp256k1?: Uint8Array | undefined;
                } | undefined;
                votingPower?: bigint | undefined;
                proposerPriority?: bigint | undefined;
            } & {
                address?: Uint8Array | undefined;
                pubKey?: ({
                    ed25519?: Uint8Array | undefined;
                    secp256k1?: Uint8Array | undefined;
                } & {
                    ed25519?: Uint8Array | undefined;
                    secp256k1?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["validatorSet"]["validators"][number]["pubKey"], keyof import("../crypto/keys").PublicKey>, never>) | undefined;
                votingPower?: bigint | undefined;
                proposerPriority?: bigint | undefined;
            } & Record<Exclude<keyof I["validatorSet"]["validators"][number], keyof import("./validator").Validator>, never>)[] & Record<Exclude<keyof I["validatorSet"]["validators"], keyof {
                address?: Uint8Array | undefined;
                pubKey?: {
                    ed25519?: Uint8Array | undefined;
                    secp256k1?: Uint8Array | undefined;
                } | undefined;
                votingPower?: bigint | undefined;
                proposerPriority?: bigint | undefined;
            }[]>, never>) | undefined;
            proposer?: ({
                address?: Uint8Array | undefined;
                pubKey?: {
                    ed25519?: Uint8Array | undefined;
                    secp256k1?: Uint8Array | undefined;
                } | undefined;
                votingPower?: bigint | undefined;
                proposerPriority?: bigint | undefined;
            } & {
                address?: Uint8Array | undefined;
                pubKey?: ({
                    ed25519?: Uint8Array | undefined;
                    secp256k1?: Uint8Array | undefined;
                } & {
                    ed25519?: Uint8Array | undefined;
                    secp256k1?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["validatorSet"]["proposer"]["pubKey"], keyof import("../crypto/keys").PublicKey>, never>) | undefined;
                votingPower?: bigint | undefined;
                proposerPriority?: bigint | undefined;
            } & Record<Exclude<keyof I["validatorSet"]["proposer"], keyof import("./validator").Validator>, never>) | undefined;
            totalVotingPower?: bigint | undefined;
        } & Record<Exclude<keyof I["validatorSet"], keyof ValidatorSet>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof LightBlock>, never>>(object: I): LightBlock;
};
export declare const BlockMeta: {
    typeUrl: string;
    encode(message: BlockMeta, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): BlockMeta;
    fromJSON(object: any): BlockMeta;
    toJSON(message: BlockMeta): unknown;
    fromPartial<I extends {
        blockId?: {
            hash?: Uint8Array | undefined;
            partSetHeader?: {
                total?: number | undefined;
                hash?: Uint8Array | undefined;
            } | undefined;
        } | undefined;
        blockSize?: bigint | undefined;
        header?: {
            version?: {
                block?: bigint | undefined;
                app?: bigint | undefined;
            } | undefined;
            chainId?: string | undefined;
            height?: bigint | undefined;
            time?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            lastBlockId?: {
                hash?: Uint8Array | undefined;
                partSetHeader?: {
                    total?: number | undefined;
                    hash?: Uint8Array | undefined;
                } | undefined;
            } | undefined;
            lastCommitHash?: Uint8Array | undefined;
            dataHash?: Uint8Array | undefined;
            validatorsHash?: Uint8Array | undefined;
            nextValidatorsHash?: Uint8Array | undefined;
            consensusHash?: Uint8Array | undefined;
            appHash?: Uint8Array | undefined;
            lastResultsHash?: Uint8Array | undefined;
            evidenceHash?: Uint8Array | undefined;
            proposerAddress?: Uint8Array | undefined;
        } | undefined;
        numTxs?: bigint | undefined;
    } & {
        blockId?: ({
            hash?: Uint8Array | undefined;
            partSetHeader?: {
                total?: number | undefined;
                hash?: Uint8Array | undefined;
            } | undefined;
        } & {
            hash?: Uint8Array | undefined;
            partSetHeader?: ({
                total?: number | undefined;
                hash?: Uint8Array | undefined;
            } & {
                total?: number | undefined;
                hash?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["blockId"]["partSetHeader"], keyof PartSetHeader>, never>) | undefined;
        } & Record<Exclude<keyof I["blockId"], keyof BlockID>, never>) | undefined;
        blockSize?: bigint | undefined;
        header?: ({
            version?: {
                block?: bigint | undefined;
                app?: bigint | undefined;
            } | undefined;
            chainId?: string | undefined;
            height?: bigint | undefined;
            time?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            lastBlockId?: {
                hash?: Uint8Array | undefined;
                partSetHeader?: {
                    total?: number | undefined;
                    hash?: Uint8Array | undefined;
                } | undefined;
            } | undefined;
            lastCommitHash?: Uint8Array | undefined;
            dataHash?: Uint8Array | undefined;
            validatorsHash?: Uint8Array | undefined;
            nextValidatorsHash?: Uint8Array | undefined;
            consensusHash?: Uint8Array | undefined;
            appHash?: Uint8Array | undefined;
            lastResultsHash?: Uint8Array | undefined;
            evidenceHash?: Uint8Array | undefined;
            proposerAddress?: Uint8Array | undefined;
        } & {
            version?: ({
                block?: bigint | undefined;
                app?: bigint | undefined;
            } & {
                block?: bigint | undefined;
                app?: bigint | undefined;
            } & Record<Exclude<keyof I["header"]["version"], keyof Consensus>, never>) | undefined;
            chainId?: string | undefined;
            height?: bigint | undefined;
            time?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["header"]["time"], keyof Timestamp>, never>) | undefined;
            lastBlockId?: ({
                hash?: Uint8Array | undefined;
                partSetHeader?: {
                    total?: number | undefined;
                    hash?: Uint8Array | undefined;
                } | undefined;
            } & {
                hash?: Uint8Array | undefined;
                partSetHeader?: ({
                    total?: number | undefined;
                    hash?: Uint8Array | undefined;
                } & {
                    total?: number | undefined;
                    hash?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["header"]["lastBlockId"]["partSetHeader"], keyof PartSetHeader>, never>) | undefined;
            } & Record<Exclude<keyof I["header"]["lastBlockId"], keyof BlockID>, never>) | undefined;
            lastCommitHash?: Uint8Array | undefined;
            dataHash?: Uint8Array | undefined;
            validatorsHash?: Uint8Array | undefined;
            nextValidatorsHash?: Uint8Array | undefined;
            consensusHash?: Uint8Array | undefined;
            appHash?: Uint8Array | undefined;
            lastResultsHash?: Uint8Array | undefined;
            evidenceHash?: Uint8Array | undefined;
            proposerAddress?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["header"], keyof Header>, never>) | undefined;
        numTxs?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof BlockMeta>, never>>(object: I): BlockMeta;
};
export declare const TxProof: {
    typeUrl: string;
    encode(message: TxProof, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): TxProof;
    fromJSON(object: any): TxProof;
    toJSON(message: TxProof): unknown;
    fromPartial<I extends {
        rootHash?: Uint8Array | undefined;
        data?: Uint8Array | undefined;
        proof?: {
            total?: bigint | undefined;
            index?: bigint | undefined;
            leafHash?: Uint8Array | undefined;
            aunts?: Uint8Array[] | undefined;
        } | undefined;
    } & {
        rootHash?: Uint8Array | undefined;
        data?: Uint8Array | undefined;
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
    } & Record<Exclude<keyof I, keyof TxProof>, never>>(object: I): TxProof;
};
