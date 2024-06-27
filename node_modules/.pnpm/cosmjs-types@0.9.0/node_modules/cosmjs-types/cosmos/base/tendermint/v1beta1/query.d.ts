import { PageRequest, PageResponse } from "../../query/v1beta1/pagination";
import { Any } from "../../../../google/protobuf/any";
import { BlockID } from "../../../../tendermint/types/types";
import { Block as Block1 } from "../../../../tendermint/types/block";
import { Block as Block2 } from "./types";
import { DefaultNodeInfo } from "../../../../tendermint/p2p/types";
import { BinaryReader, BinaryWriter } from "../../../../binary";
import { Rpc } from "../../../../helpers";
export declare const protobufPackage = "cosmos.base.tendermint.v1beta1";
/** GetValidatorSetByHeightRequest is the request type for the Query/GetValidatorSetByHeight RPC method. */
export interface GetValidatorSetByHeightRequest {
    height: bigint;
    /** pagination defines an pagination for the request. */
    pagination?: PageRequest;
}
/** GetValidatorSetByHeightResponse is the response type for the Query/GetValidatorSetByHeight RPC method. */
export interface GetValidatorSetByHeightResponse {
    blockHeight: bigint;
    validators: Validator[];
    /** pagination defines an pagination for the response. */
    pagination?: PageResponse;
}
/** GetLatestValidatorSetRequest is the request type for the Query/GetValidatorSetByHeight RPC method. */
export interface GetLatestValidatorSetRequest {
    /** pagination defines an pagination for the request. */
    pagination?: PageRequest;
}
/** GetLatestValidatorSetResponse is the response type for the Query/GetValidatorSetByHeight RPC method. */
export interface GetLatestValidatorSetResponse {
    blockHeight: bigint;
    validators: Validator[];
    /** pagination defines an pagination for the response. */
    pagination?: PageResponse;
}
/** Validator is the type for the validator-set. */
export interface Validator {
    address: string;
    pubKey?: Any;
    votingPower: bigint;
    proposerPriority: bigint;
}
/** GetBlockByHeightRequest is the request type for the Query/GetBlockByHeight RPC method. */
export interface GetBlockByHeightRequest {
    height: bigint;
}
/** GetBlockByHeightResponse is the response type for the Query/GetBlockByHeight RPC method. */
export interface GetBlockByHeightResponse {
    blockId?: BlockID;
    /** Deprecated: please use `sdk_block` instead */
    block?: Block1;
    /** Since: cosmos-sdk 0.47 */
    sdkBlock?: Block2;
}
/** GetLatestBlockRequest is the request type for the Query/GetLatestBlock RPC method. */
export interface GetLatestBlockRequest {
}
/** GetLatestBlockResponse is the response type for the Query/GetLatestBlock RPC method. */
export interface GetLatestBlockResponse {
    blockId?: BlockID;
    /** Deprecated: please use `sdk_block` instead */
    block?: Block1;
    /** Since: cosmos-sdk 0.47 */
    sdkBlock?: Block2;
}
/** GetSyncingRequest is the request type for the Query/GetSyncing RPC method. */
export interface GetSyncingRequest {
}
/** GetSyncingResponse is the response type for the Query/GetSyncing RPC method. */
export interface GetSyncingResponse {
    syncing: boolean;
}
/** GetNodeInfoRequest is the request type for the Query/GetNodeInfo RPC method. */
export interface GetNodeInfoRequest {
}
/** GetNodeInfoResponse is the response type for the Query/GetNodeInfo RPC method. */
export interface GetNodeInfoResponse {
    defaultNodeInfo?: DefaultNodeInfo;
    applicationVersion?: VersionInfo;
}
/** VersionInfo is the type for the GetNodeInfoResponse message. */
export interface VersionInfo {
    name: string;
    appName: string;
    version: string;
    gitCommit: string;
    buildTags: string;
    goVersion: string;
    buildDeps: Module[];
    /** Since: cosmos-sdk 0.43 */
    cosmosSdkVersion: string;
}
/** Module is the type for VersionInfo */
export interface Module {
    /** module path */
    path: string;
    /** module version */
    version: string;
    /** checksum */
    sum: string;
}
/** ABCIQueryRequest defines the request structure for the ABCIQuery gRPC query. */
export interface ABCIQueryRequest {
    data: Uint8Array;
    path: string;
    height: bigint;
    prove: boolean;
}
/**
 * ABCIQueryResponse defines the response structure for the ABCIQuery gRPC query.
 *
 * Note: This type is a duplicate of the ResponseQuery proto type defined in
 * Tendermint.
 */
export interface ABCIQueryResponse {
    code: number;
    /** nondeterministic */
    log: string;
    /** nondeterministic */
    info: string;
    index: bigint;
    key: Uint8Array;
    value: Uint8Array;
    proofOps?: ProofOps;
    height: bigint;
    codespace: string;
}
/**
 * ProofOp defines an operation used for calculating Merkle root. The data could
 * be arbitrary format, providing necessary data for example neighbouring node
 * hash.
 *
 * Note: This type is a duplicate of the ProofOp proto type defined in Tendermint.
 */
export interface ProofOp {
    type: string;
    key: Uint8Array;
    data: Uint8Array;
}
/**
 * ProofOps is Merkle proof defined by the list of ProofOps.
 *
 * Note: This type is a duplicate of the ProofOps proto type defined in Tendermint.
 */
export interface ProofOps {
    ops: ProofOp[];
}
export declare const GetValidatorSetByHeightRequest: {
    typeUrl: string;
    encode(message: GetValidatorSetByHeightRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GetValidatorSetByHeightRequest;
    fromJSON(object: any): GetValidatorSetByHeightRequest;
    toJSON(message: GetValidatorSetByHeightRequest): unknown;
    fromPartial<I extends {
        height?: bigint | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        height?: bigint | undefined;
        pagination?: ({
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } & {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageRequest>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof GetValidatorSetByHeightRequest>, never>>(object: I): GetValidatorSetByHeightRequest;
};
export declare const GetValidatorSetByHeightResponse: {
    typeUrl: string;
    encode(message: GetValidatorSetByHeightResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GetValidatorSetByHeightResponse;
    fromJSON(object: any): GetValidatorSetByHeightResponse;
    toJSON(message: GetValidatorSetByHeightResponse): unknown;
    fromPartial<I extends {
        blockHeight?: bigint | undefined;
        validators?: {
            address?: string | undefined;
            pubKey?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            votingPower?: bigint | undefined;
            proposerPriority?: bigint | undefined;
        }[] | undefined;
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
        blockHeight?: bigint | undefined;
        validators?: ({
            address?: string | undefined;
            pubKey?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            votingPower?: bigint | undefined;
            proposerPriority?: bigint | undefined;
        }[] & ({
            address?: string | undefined;
            pubKey?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            votingPower?: bigint | undefined;
            proposerPriority?: bigint | undefined;
        } & {
            address?: string | undefined;
            pubKey?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["validators"][number]["pubKey"], keyof Any>, never>) | undefined;
            votingPower?: bigint | undefined;
            proposerPriority?: bigint | undefined;
        } & Record<Exclude<keyof I["validators"][number], keyof Validator>, never>)[] & Record<Exclude<keyof I["validators"], keyof {
            address?: string | undefined;
            pubKey?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            votingPower?: bigint | undefined;
            proposerPriority?: bigint | undefined;
        }[]>, never>) | undefined;
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof GetValidatorSetByHeightResponse>, never>>(object: I): GetValidatorSetByHeightResponse;
};
export declare const GetLatestValidatorSetRequest: {
    typeUrl: string;
    encode(message: GetLatestValidatorSetRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GetLatestValidatorSetRequest;
    fromJSON(object: any): GetLatestValidatorSetRequest;
    toJSON(message: GetLatestValidatorSetRequest): unknown;
    fromPartial<I extends {
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        pagination?: ({
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } & {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageRequest>, never>) | undefined;
    } & Record<Exclude<keyof I, "pagination">, never>>(object: I): GetLatestValidatorSetRequest;
};
export declare const GetLatestValidatorSetResponse: {
    typeUrl: string;
    encode(message: GetLatestValidatorSetResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GetLatestValidatorSetResponse;
    fromJSON(object: any): GetLatestValidatorSetResponse;
    toJSON(message: GetLatestValidatorSetResponse): unknown;
    fromPartial<I extends {
        blockHeight?: bigint | undefined;
        validators?: {
            address?: string | undefined;
            pubKey?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            votingPower?: bigint | undefined;
            proposerPriority?: bigint | undefined;
        }[] | undefined;
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
        blockHeight?: bigint | undefined;
        validators?: ({
            address?: string | undefined;
            pubKey?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            votingPower?: bigint | undefined;
            proposerPriority?: bigint | undefined;
        }[] & ({
            address?: string | undefined;
            pubKey?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            votingPower?: bigint | undefined;
            proposerPriority?: bigint | undefined;
        } & {
            address?: string | undefined;
            pubKey?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["validators"][number]["pubKey"], keyof Any>, never>) | undefined;
            votingPower?: bigint | undefined;
            proposerPriority?: bigint | undefined;
        } & Record<Exclude<keyof I["validators"][number], keyof Validator>, never>)[] & Record<Exclude<keyof I["validators"], keyof {
            address?: string | undefined;
            pubKey?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            votingPower?: bigint | undefined;
            proposerPriority?: bigint | undefined;
        }[]>, never>) | undefined;
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof GetLatestValidatorSetResponse>, never>>(object: I): GetLatestValidatorSetResponse;
};
export declare const Validator: {
    typeUrl: string;
    encode(message: Validator, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Validator;
    fromJSON(object: any): Validator;
    toJSON(message: Validator): unknown;
    fromPartial<I extends {
        address?: string | undefined;
        pubKey?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
        votingPower?: bigint | undefined;
        proposerPriority?: bigint | undefined;
    } & {
        address?: string | undefined;
        pubKey?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["pubKey"], keyof Any>, never>) | undefined;
        votingPower?: bigint | undefined;
        proposerPriority?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof Validator>, never>>(object: I): Validator;
};
export declare const GetBlockByHeightRequest: {
    typeUrl: string;
    encode(message: GetBlockByHeightRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GetBlockByHeightRequest;
    fromJSON(object: any): GetBlockByHeightRequest;
    toJSON(message: GetBlockByHeightRequest): unknown;
    fromPartial<I extends {
        height?: bigint | undefined;
    } & {
        height?: bigint | undefined;
    } & Record<Exclude<keyof I, "height">, never>>(object: I): GetBlockByHeightRequest;
};
export declare const GetBlockByHeightResponse: {
    typeUrl: string;
    encode(message: GetBlockByHeightResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GetBlockByHeightResponse;
    fromJSON(object: any): GetBlockByHeightResponse;
    toJSON(message: GetBlockByHeightResponse): unknown;
    fromPartial<I extends {
        blockId?: {
            hash?: Uint8Array | undefined;
            partSetHeader?: {
                total?: number | undefined;
                hash?: Uint8Array | undefined;
            } | undefined;
        } | undefined;
        block?: {
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
            data?: {
                txs?: Uint8Array[] | undefined;
            } | undefined;
            evidence?: {
                evidence?: {
                    duplicateVoteEvidence?: {
                        voteA?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        voteB?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        totalVotingPower?: bigint | undefined;
                        validatorPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                    lightClientAttackEvidence?: {
                        conflictingBlock?: {
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                        } | undefined;
                        commonHeight?: bigint | undefined;
                        byzantineValidators?: {
                            address?: Uint8Array | undefined;
                            pubKey?: {
                                ed25519?: Uint8Array | undefined;
                                secp256k1?: Uint8Array | undefined;
                            } | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        }[] | undefined;
                        totalVotingPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                }[] | undefined;
            } | undefined;
            lastCommit?: {
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
                    blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    timestamp?: {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } | undefined;
                    signature?: Uint8Array | undefined;
                }[] | undefined;
            } | undefined;
        } | undefined;
        sdkBlock?: {
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
                proposerAddress?: string | undefined;
            } | undefined;
            data?: {
                txs?: Uint8Array[] | undefined;
            } | undefined;
            evidence?: {
                evidence?: {
                    duplicateVoteEvidence?: {
                        voteA?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        voteB?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        totalVotingPower?: bigint | undefined;
                        validatorPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                    lightClientAttackEvidence?: {
                        conflictingBlock?: {
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                        } | undefined;
                        commonHeight?: bigint | undefined;
                        byzantineValidators?: {
                            address?: Uint8Array | undefined;
                            pubKey?: {
                                ed25519?: Uint8Array | undefined;
                                secp256k1?: Uint8Array | undefined;
                            } | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        }[] | undefined;
                        totalVotingPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                }[] | undefined;
            } | undefined;
            lastCommit?: {
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
                    blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    timestamp?: {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } | undefined;
                    signature?: Uint8Array | undefined;
                }[] | undefined;
            } | undefined;
        } | undefined;
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
            } & Record<Exclude<keyof I["blockId"]["partSetHeader"], keyof import("../../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
        } & Record<Exclude<keyof I["blockId"], keyof BlockID>, never>) | undefined;
        block?: ({
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
            data?: {
                txs?: Uint8Array[] | undefined;
            } | undefined;
            evidence?: {
                evidence?: {
                    duplicateVoteEvidence?: {
                        voteA?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        voteB?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        totalVotingPower?: bigint | undefined;
                        validatorPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                    lightClientAttackEvidence?: {
                        conflictingBlock?: {
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                        } | undefined;
                        commonHeight?: bigint | undefined;
                        byzantineValidators?: {
                            address?: Uint8Array | undefined;
                            pubKey?: {
                                ed25519?: Uint8Array | undefined;
                                secp256k1?: Uint8Array | undefined;
                            } | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        }[] | undefined;
                        totalVotingPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                }[] | undefined;
            } | undefined;
            lastCommit?: {
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
                    blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                } & Record<Exclude<keyof I["block"]["header"]["version"], keyof import("../../../../tendermint/version/types").Consensus>, never>) | undefined;
                chainId?: string | undefined;
                height?: bigint | undefined;
                time?: ({
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & Record<Exclude<keyof I["block"]["header"]["time"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
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
                    } & Record<Exclude<keyof I["block"]["header"]["lastBlockId"]["partSetHeader"], keyof import("../../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                } & Record<Exclude<keyof I["block"]["header"]["lastBlockId"], keyof BlockID>, never>) | undefined;
                lastCommitHash?: Uint8Array | undefined;
                dataHash?: Uint8Array | undefined;
                validatorsHash?: Uint8Array | undefined;
                nextValidatorsHash?: Uint8Array | undefined;
                consensusHash?: Uint8Array | undefined;
                appHash?: Uint8Array | undefined;
                lastResultsHash?: Uint8Array | undefined;
                evidenceHash?: Uint8Array | undefined;
                proposerAddress?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["block"]["header"], keyof import("../../../../tendermint/types/types").Header>, never>) | undefined;
            data?: ({
                txs?: Uint8Array[] | undefined;
            } & {
                txs?: (Uint8Array[] & Uint8Array[] & Record<Exclude<keyof I["block"]["data"]["txs"], keyof Uint8Array[]>, never>) | undefined;
            } & Record<Exclude<keyof I["block"]["data"], "txs">, never>) | undefined;
            evidence?: ({
                evidence?: {
                    duplicateVoteEvidence?: {
                        voteA?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        voteB?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        totalVotingPower?: bigint | undefined;
                        validatorPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                    lightClientAttackEvidence?: {
                        conflictingBlock?: {
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                        } | undefined;
                        commonHeight?: bigint | undefined;
                        byzantineValidators?: {
                            address?: Uint8Array | undefined;
                            pubKey?: {
                                ed25519?: Uint8Array | undefined;
                                secp256k1?: Uint8Array | undefined;
                            } | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        }[] | undefined;
                        totalVotingPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                }[] | undefined;
            } & {
                evidence?: ({
                    duplicateVoteEvidence?: {
                        voteA?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        voteB?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        totalVotingPower?: bigint | undefined;
                        validatorPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                    lightClientAttackEvidence?: {
                        conflictingBlock?: {
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                        } | undefined;
                        commonHeight?: bigint | undefined;
                        byzantineValidators?: {
                            address?: Uint8Array | undefined;
                            pubKey?: {
                                ed25519?: Uint8Array | undefined;
                                secp256k1?: Uint8Array | undefined;
                            } | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        }[] | undefined;
                        totalVotingPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                }[] & ({
                    duplicateVoteEvidence?: {
                        voteA?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        voteB?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        totalVotingPower?: bigint | undefined;
                        validatorPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                    lightClientAttackEvidence?: {
                        conflictingBlock?: {
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                        } | undefined;
                        commonHeight?: bigint | undefined;
                        byzantineValidators?: {
                            address?: Uint8Array | undefined;
                            pubKey?: {
                                ed25519?: Uint8Array | undefined;
                                secp256k1?: Uint8Array | undefined;
                            } | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        }[] | undefined;
                        totalVotingPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                } & {
                    duplicateVoteEvidence?: ({
                        voteA?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        voteB?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        totalVotingPower?: bigint | undefined;
                        validatorPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } & {
                        voteA?: ({
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                                } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteA"]["blockId"]["partSetHeader"], keyof import("../../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                            } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteA"]["blockId"], keyof BlockID>, never>) | undefined;
                            timestamp?: ({
                                seconds?: bigint | undefined;
                                nanos?: number | undefined;
                            } & {
                                seconds?: bigint | undefined;
                                nanos?: number | undefined;
                            } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteA"]["timestamp"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                            validatorAddress?: Uint8Array | undefined;
                            validatorIndex?: number | undefined;
                            signature?: Uint8Array | undefined;
                        } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteA"], keyof import("../../../../tendermint/types/types").Vote>, never>) | undefined;
                        voteB?: ({
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                                } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteB"]["blockId"]["partSetHeader"], keyof import("../../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                            } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteB"]["blockId"], keyof BlockID>, never>) | undefined;
                            timestamp?: ({
                                seconds?: bigint | undefined;
                                nanos?: number | undefined;
                            } & {
                                seconds?: bigint | undefined;
                                nanos?: number | undefined;
                            } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteB"]["timestamp"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                            validatorAddress?: Uint8Array | undefined;
                            validatorIndex?: number | undefined;
                            signature?: Uint8Array | undefined;
                        } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteB"], keyof import("../../../../tendermint/types/types").Vote>, never>) | undefined;
                        totalVotingPower?: bigint | undefined;
                        validatorPower?: bigint | undefined;
                        timestamp?: ({
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } & {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["timestamp"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                    } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["duplicateVoteEvidence"], keyof import("../../../../tendermint/types/evidence").DuplicateVoteEvidence>, never>) | undefined;
                    lightClientAttackEvidence?: ({
                        conflictingBlock?: {
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                        } | undefined;
                        commonHeight?: bigint | undefined;
                        byzantineValidators?: {
                            address?: Uint8Array | undefined;
                            pubKey?: {
                                ed25519?: Uint8Array | undefined;
                                secp256k1?: Uint8Array | undefined;
                            } | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        }[] | undefined;
                        totalVotingPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } & {
                        conflictingBlock?: ({
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                                    } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"]["version"], keyof import("../../../../tendermint/version/types").Consensus>, never>) | undefined;
                                    chainId?: string | undefined;
                                    height?: bigint | undefined;
                                    time?: ({
                                        seconds?: bigint | undefined;
                                        nanos?: number | undefined;
                                    } & {
                                        seconds?: bigint | undefined;
                                        nanos?: number | undefined;
                                    } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"]["time"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
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
                                        } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"]["lastBlockId"]["partSetHeader"], keyof import("../../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                                    } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"]["lastBlockId"], keyof BlockID>, never>) | undefined;
                                    lastCommitHash?: Uint8Array | undefined;
                                    dataHash?: Uint8Array | undefined;
                                    validatorsHash?: Uint8Array | undefined;
                                    nextValidatorsHash?: Uint8Array | undefined;
                                    consensusHash?: Uint8Array | undefined;
                                    appHash?: Uint8Array | undefined;
                                    lastResultsHash?: Uint8Array | undefined;
                                    evidenceHash?: Uint8Array | undefined;
                                    proposerAddress?: Uint8Array | undefined;
                                } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"], keyof import("../../../../tendermint/types/types").Header>, never>) | undefined;
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                                        } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["blockId"]["partSetHeader"], keyof import("../../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                                    } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["blockId"], keyof BlockID>, never>) | undefined;
                                    signatures?: ({
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                                        validatorAddress?: Uint8Array | undefined;
                                        timestamp?: {
                                            seconds?: bigint | undefined;
                                            nanos?: number | undefined;
                                        } | undefined;
                                        signature?: Uint8Array | undefined;
                                    }[] & ({
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                                        validatorAddress?: Uint8Array | undefined;
                                        timestamp?: {
                                            seconds?: bigint | undefined;
                                            nanos?: number | undefined;
                                        } | undefined;
                                        signature?: Uint8Array | undefined;
                                    } & {
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                                        validatorAddress?: Uint8Array | undefined;
                                        timestamp?: ({
                                            seconds?: bigint | undefined;
                                            nanos?: number | undefined;
                                        } & {
                                            seconds?: bigint | undefined;
                                            nanos?: number | undefined;
                                        } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["signatures"][number]["timestamp"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                                        signature?: Uint8Array | undefined;
                                    } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["signatures"][number], keyof import("../../../../tendermint/types/types").CommitSig>, never>)[] & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["signatures"], keyof {
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                                        validatorAddress?: Uint8Array | undefined;
                                        timestamp?: {
                                            seconds?: bigint | undefined;
                                            nanos?: number | undefined;
                                        } | undefined;
                                        signature?: Uint8Array | undefined;
                                    }[]>, never>) | undefined;
                                } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"], keyof import("../../../../tendermint/types/types").Commit>, never>) | undefined;
                            } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"], keyof import("../../../../tendermint/types/types").SignedHeader>, never>) | undefined;
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
                                    } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["validators"][number]["pubKey"], keyof import("../../../../tendermint/crypto/keys").PublicKey>, never>) | undefined;
                                    votingPower?: bigint | undefined;
                                    proposerPriority?: bigint | undefined;
                                } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["validators"][number], keyof import("../../../../tendermint/types/validator").Validator>, never>)[] & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["validators"], keyof {
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
                                    } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["proposer"]["pubKey"], keyof import("../../../../tendermint/crypto/keys").PublicKey>, never>) | undefined;
                                    votingPower?: bigint | undefined;
                                    proposerPriority?: bigint | undefined;
                                } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["proposer"], keyof import("../../../../tendermint/types/validator").Validator>, never>) | undefined;
                                totalVotingPower?: bigint | undefined;
                            } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"], keyof import("../../../../tendermint/types/validator").ValidatorSet>, never>) | undefined;
                        } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"], keyof import("../../../../tendermint/types/types").LightBlock>, never>) | undefined;
                        commonHeight?: bigint | undefined;
                        byzantineValidators?: ({
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
                            } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["byzantineValidators"][number]["pubKey"], keyof import("../../../../tendermint/crypto/keys").PublicKey>, never>) | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["byzantineValidators"][number], keyof import("../../../../tendermint/types/validator").Validator>, never>)[] & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["byzantineValidators"], keyof {
                            address?: Uint8Array | undefined;
                            pubKey?: {
                                ed25519?: Uint8Array | undefined;
                                secp256k1?: Uint8Array | undefined;
                            } | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        }[]>, never>) | undefined;
                        totalVotingPower?: bigint | undefined;
                        timestamp?: ({
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } & {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["timestamp"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                    } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"], keyof import("../../../../tendermint/types/evidence").LightClientAttackEvidence>, never>) | undefined;
                } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number], keyof import("../../../../tendermint/types/evidence").Evidence>, never>)[] & Record<Exclude<keyof I["block"]["evidence"]["evidence"], keyof {
                    duplicateVoteEvidence?: {
                        voteA?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        voteB?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        totalVotingPower?: bigint | undefined;
                        validatorPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                    lightClientAttackEvidence?: {
                        conflictingBlock?: {
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                        } | undefined;
                        commonHeight?: bigint | undefined;
                        byzantineValidators?: {
                            address?: Uint8Array | undefined;
                            pubKey?: {
                                ed25519?: Uint8Array | undefined;
                                secp256k1?: Uint8Array | undefined;
                            } | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        }[] | undefined;
                        totalVotingPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["block"]["evidence"], "evidence">, never>) | undefined;
            lastCommit?: ({
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
                    blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                    } & Record<Exclude<keyof I["block"]["lastCommit"]["blockId"]["partSetHeader"], keyof import("../../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                } & Record<Exclude<keyof I["block"]["lastCommit"]["blockId"], keyof BlockID>, never>) | undefined;
                signatures?: ({
                    blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    timestamp?: {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } | undefined;
                    signature?: Uint8Array | undefined;
                }[] & ({
                    blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    timestamp?: {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } | undefined;
                    signature?: Uint8Array | undefined;
                } & {
                    blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    timestamp?: ({
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } & {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } & Record<Exclude<keyof I["block"]["lastCommit"]["signatures"][number]["timestamp"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                    signature?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["block"]["lastCommit"]["signatures"][number], keyof import("../../../../tendermint/types/types").CommitSig>, never>)[] & Record<Exclude<keyof I["block"]["lastCommit"]["signatures"], keyof {
                    blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    timestamp?: {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } | undefined;
                    signature?: Uint8Array | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["block"]["lastCommit"], keyof import("../../../../tendermint/types/types").Commit>, never>) | undefined;
        } & Record<Exclude<keyof I["block"], keyof Block1>, never>) | undefined;
        sdkBlock?: ({
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
                proposerAddress?: string | undefined;
            } | undefined;
            data?: {
                txs?: Uint8Array[] | undefined;
            } | undefined;
            evidence?: {
                evidence?: {
                    duplicateVoteEvidence?: {
                        voteA?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        voteB?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        totalVotingPower?: bigint | undefined;
                        validatorPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                    lightClientAttackEvidence?: {
                        conflictingBlock?: {
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                        } | undefined;
                        commonHeight?: bigint | undefined;
                        byzantineValidators?: {
                            address?: Uint8Array | undefined;
                            pubKey?: {
                                ed25519?: Uint8Array | undefined;
                                secp256k1?: Uint8Array | undefined;
                            } | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        }[] | undefined;
                        totalVotingPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                }[] | undefined;
            } | undefined;
            lastCommit?: {
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
                    blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                proposerAddress?: string | undefined;
            } & {
                version?: ({
                    block?: bigint | undefined;
                    app?: bigint | undefined;
                } & {
                    block?: bigint | undefined;
                    app?: bigint | undefined;
                } & Record<Exclude<keyof I["sdkBlock"]["header"]["version"], keyof import("../../../../tendermint/version/types").Consensus>, never>) | undefined;
                chainId?: string | undefined;
                height?: bigint | undefined;
                time?: ({
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & Record<Exclude<keyof I["sdkBlock"]["header"]["time"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
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
                    } & Record<Exclude<keyof I["sdkBlock"]["header"]["lastBlockId"]["partSetHeader"], keyof import("../../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                } & Record<Exclude<keyof I["sdkBlock"]["header"]["lastBlockId"], keyof BlockID>, never>) | undefined;
                lastCommitHash?: Uint8Array | undefined;
                dataHash?: Uint8Array | undefined;
                validatorsHash?: Uint8Array | undefined;
                nextValidatorsHash?: Uint8Array | undefined;
                consensusHash?: Uint8Array | undefined;
                appHash?: Uint8Array | undefined;
                lastResultsHash?: Uint8Array | undefined;
                evidenceHash?: Uint8Array | undefined;
                proposerAddress?: string | undefined;
            } & Record<Exclude<keyof I["sdkBlock"]["header"], keyof import("./types").Header>, never>) | undefined;
            data?: ({
                txs?: Uint8Array[] | undefined;
            } & {
                txs?: (Uint8Array[] & Uint8Array[] & Record<Exclude<keyof I["sdkBlock"]["data"]["txs"], keyof Uint8Array[]>, never>) | undefined;
            } & Record<Exclude<keyof I["sdkBlock"]["data"], "txs">, never>) | undefined;
            evidence?: ({
                evidence?: {
                    duplicateVoteEvidence?: {
                        voteA?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        voteB?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        totalVotingPower?: bigint | undefined;
                        validatorPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                    lightClientAttackEvidence?: {
                        conflictingBlock?: {
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                        } | undefined;
                        commonHeight?: bigint | undefined;
                        byzantineValidators?: {
                            address?: Uint8Array | undefined;
                            pubKey?: {
                                ed25519?: Uint8Array | undefined;
                                secp256k1?: Uint8Array | undefined;
                            } | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        }[] | undefined;
                        totalVotingPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                }[] | undefined;
            } & {
                evidence?: ({
                    duplicateVoteEvidence?: {
                        voteA?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        voteB?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        totalVotingPower?: bigint | undefined;
                        validatorPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                    lightClientAttackEvidence?: {
                        conflictingBlock?: {
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                        } | undefined;
                        commonHeight?: bigint | undefined;
                        byzantineValidators?: {
                            address?: Uint8Array | undefined;
                            pubKey?: {
                                ed25519?: Uint8Array | undefined;
                                secp256k1?: Uint8Array | undefined;
                            } | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        }[] | undefined;
                        totalVotingPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                }[] & ({
                    duplicateVoteEvidence?: {
                        voteA?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        voteB?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        totalVotingPower?: bigint | undefined;
                        validatorPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                    lightClientAttackEvidence?: {
                        conflictingBlock?: {
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                        } | undefined;
                        commonHeight?: bigint | undefined;
                        byzantineValidators?: {
                            address?: Uint8Array | undefined;
                            pubKey?: {
                                ed25519?: Uint8Array | undefined;
                                secp256k1?: Uint8Array | undefined;
                            } | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        }[] | undefined;
                        totalVotingPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                } & {
                    duplicateVoteEvidence?: ({
                        voteA?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        voteB?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        totalVotingPower?: bigint | undefined;
                        validatorPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } & {
                        voteA?: ({
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                                } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteA"]["blockId"]["partSetHeader"], keyof import("../../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                            } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteA"]["blockId"], keyof BlockID>, never>) | undefined;
                            timestamp?: ({
                                seconds?: bigint | undefined;
                                nanos?: number | undefined;
                            } & {
                                seconds?: bigint | undefined;
                                nanos?: number | undefined;
                            } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteA"]["timestamp"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                            validatorAddress?: Uint8Array | undefined;
                            validatorIndex?: number | undefined;
                            signature?: Uint8Array | undefined;
                        } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteA"], keyof import("../../../../tendermint/types/types").Vote>, never>) | undefined;
                        voteB?: ({
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                                } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteB"]["blockId"]["partSetHeader"], keyof import("../../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                            } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteB"]["blockId"], keyof BlockID>, never>) | undefined;
                            timestamp?: ({
                                seconds?: bigint | undefined;
                                nanos?: number | undefined;
                            } & {
                                seconds?: bigint | undefined;
                                nanos?: number | undefined;
                            } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteB"]["timestamp"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                            validatorAddress?: Uint8Array | undefined;
                            validatorIndex?: number | undefined;
                            signature?: Uint8Array | undefined;
                        } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteB"], keyof import("../../../../tendermint/types/types").Vote>, never>) | undefined;
                        totalVotingPower?: bigint | undefined;
                        validatorPower?: bigint | undefined;
                        timestamp?: ({
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } & {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["timestamp"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                    } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["duplicateVoteEvidence"], keyof import("../../../../tendermint/types/evidence").DuplicateVoteEvidence>, never>) | undefined;
                    lightClientAttackEvidence?: ({
                        conflictingBlock?: {
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                        } | undefined;
                        commonHeight?: bigint | undefined;
                        byzantineValidators?: {
                            address?: Uint8Array | undefined;
                            pubKey?: {
                                ed25519?: Uint8Array | undefined;
                                secp256k1?: Uint8Array | undefined;
                            } | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        }[] | undefined;
                        totalVotingPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } & {
                        conflictingBlock?: ({
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                                    } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"]["version"], keyof import("../../../../tendermint/version/types").Consensus>, never>) | undefined;
                                    chainId?: string | undefined;
                                    height?: bigint | undefined;
                                    time?: ({
                                        seconds?: bigint | undefined;
                                        nanos?: number | undefined;
                                    } & {
                                        seconds?: bigint | undefined;
                                        nanos?: number | undefined;
                                    } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"]["time"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
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
                                        } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"]["lastBlockId"]["partSetHeader"], keyof import("../../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                                    } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"]["lastBlockId"], keyof BlockID>, never>) | undefined;
                                    lastCommitHash?: Uint8Array | undefined;
                                    dataHash?: Uint8Array | undefined;
                                    validatorsHash?: Uint8Array | undefined;
                                    nextValidatorsHash?: Uint8Array | undefined;
                                    consensusHash?: Uint8Array | undefined;
                                    appHash?: Uint8Array | undefined;
                                    lastResultsHash?: Uint8Array | undefined;
                                    evidenceHash?: Uint8Array | undefined;
                                    proposerAddress?: Uint8Array | undefined;
                                } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"], keyof import("../../../../tendermint/types/types").Header>, never>) | undefined;
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                                        } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["blockId"]["partSetHeader"], keyof import("../../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                                    } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["blockId"], keyof BlockID>, never>) | undefined;
                                    signatures?: ({
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                                        validatorAddress?: Uint8Array | undefined;
                                        timestamp?: {
                                            seconds?: bigint | undefined;
                                            nanos?: number | undefined;
                                        } | undefined;
                                        signature?: Uint8Array | undefined;
                                    }[] & ({
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                                        validatorAddress?: Uint8Array | undefined;
                                        timestamp?: {
                                            seconds?: bigint | undefined;
                                            nanos?: number | undefined;
                                        } | undefined;
                                        signature?: Uint8Array | undefined;
                                    } & {
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                                        validatorAddress?: Uint8Array | undefined;
                                        timestamp?: ({
                                            seconds?: bigint | undefined;
                                            nanos?: number | undefined;
                                        } & {
                                            seconds?: bigint | undefined;
                                            nanos?: number | undefined;
                                        } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["signatures"][number]["timestamp"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                                        signature?: Uint8Array | undefined;
                                    } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["signatures"][number], keyof import("../../../../tendermint/types/types").CommitSig>, never>)[] & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["signatures"], keyof {
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                                        validatorAddress?: Uint8Array | undefined;
                                        timestamp?: {
                                            seconds?: bigint | undefined;
                                            nanos?: number | undefined;
                                        } | undefined;
                                        signature?: Uint8Array | undefined;
                                    }[]>, never>) | undefined;
                                } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"], keyof import("../../../../tendermint/types/types").Commit>, never>) | undefined;
                            } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"], keyof import("../../../../tendermint/types/types").SignedHeader>, never>) | undefined;
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
                                    } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["validators"][number]["pubKey"], keyof import("../../../../tendermint/crypto/keys").PublicKey>, never>) | undefined;
                                    votingPower?: bigint | undefined;
                                    proposerPriority?: bigint | undefined;
                                } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["validators"][number], keyof import("../../../../tendermint/types/validator").Validator>, never>)[] & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["validators"], keyof {
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
                                    } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["proposer"]["pubKey"], keyof import("../../../../tendermint/crypto/keys").PublicKey>, never>) | undefined;
                                    votingPower?: bigint | undefined;
                                    proposerPriority?: bigint | undefined;
                                } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["proposer"], keyof import("../../../../tendermint/types/validator").Validator>, never>) | undefined;
                                totalVotingPower?: bigint | undefined;
                            } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"], keyof import("../../../../tendermint/types/validator").ValidatorSet>, never>) | undefined;
                        } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"], keyof import("../../../../tendermint/types/types").LightBlock>, never>) | undefined;
                        commonHeight?: bigint | undefined;
                        byzantineValidators?: ({
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
                            } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["byzantineValidators"][number]["pubKey"], keyof import("../../../../tendermint/crypto/keys").PublicKey>, never>) | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["byzantineValidators"][number], keyof import("../../../../tendermint/types/validator").Validator>, never>)[] & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["byzantineValidators"], keyof {
                            address?: Uint8Array | undefined;
                            pubKey?: {
                                ed25519?: Uint8Array | undefined;
                                secp256k1?: Uint8Array | undefined;
                            } | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        }[]>, never>) | undefined;
                        totalVotingPower?: bigint | undefined;
                        timestamp?: ({
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } & {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["timestamp"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                    } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"], keyof import("../../../../tendermint/types/evidence").LightClientAttackEvidence>, never>) | undefined;
                } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number], keyof import("../../../../tendermint/types/evidence").Evidence>, never>)[] & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"], keyof {
                    duplicateVoteEvidence?: {
                        voteA?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        voteB?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        totalVotingPower?: bigint | undefined;
                        validatorPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                    lightClientAttackEvidence?: {
                        conflictingBlock?: {
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                        } | undefined;
                        commonHeight?: bigint | undefined;
                        byzantineValidators?: {
                            address?: Uint8Array | undefined;
                            pubKey?: {
                                ed25519?: Uint8Array | undefined;
                                secp256k1?: Uint8Array | undefined;
                            } | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        }[] | undefined;
                        totalVotingPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["sdkBlock"]["evidence"], "evidence">, never>) | undefined;
            lastCommit?: ({
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
                    blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                    } & Record<Exclude<keyof I["sdkBlock"]["lastCommit"]["blockId"]["partSetHeader"], keyof import("../../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                } & Record<Exclude<keyof I["sdkBlock"]["lastCommit"]["blockId"], keyof BlockID>, never>) | undefined;
                signatures?: ({
                    blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    timestamp?: {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } | undefined;
                    signature?: Uint8Array | undefined;
                }[] & ({
                    blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    timestamp?: {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } | undefined;
                    signature?: Uint8Array | undefined;
                } & {
                    blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    timestamp?: ({
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } & {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } & Record<Exclude<keyof I["sdkBlock"]["lastCommit"]["signatures"][number]["timestamp"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                    signature?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["sdkBlock"]["lastCommit"]["signatures"][number], keyof import("../../../../tendermint/types/types").CommitSig>, never>)[] & Record<Exclude<keyof I["sdkBlock"]["lastCommit"]["signatures"], keyof {
                    blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    timestamp?: {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } | undefined;
                    signature?: Uint8Array | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["sdkBlock"]["lastCommit"], keyof import("../../../../tendermint/types/types").Commit>, never>) | undefined;
        } & Record<Exclude<keyof I["sdkBlock"], keyof Block2>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof GetBlockByHeightResponse>, never>>(object: I): GetBlockByHeightResponse;
};
export declare const GetLatestBlockRequest: {
    typeUrl: string;
    encode(_: GetLatestBlockRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GetLatestBlockRequest;
    fromJSON(_: any): GetLatestBlockRequest;
    toJSON(_: GetLatestBlockRequest): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): GetLatestBlockRequest;
};
export declare const GetLatestBlockResponse: {
    typeUrl: string;
    encode(message: GetLatestBlockResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GetLatestBlockResponse;
    fromJSON(object: any): GetLatestBlockResponse;
    toJSON(message: GetLatestBlockResponse): unknown;
    fromPartial<I extends {
        blockId?: {
            hash?: Uint8Array | undefined;
            partSetHeader?: {
                total?: number | undefined;
                hash?: Uint8Array | undefined;
            } | undefined;
        } | undefined;
        block?: {
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
            data?: {
                txs?: Uint8Array[] | undefined;
            } | undefined;
            evidence?: {
                evidence?: {
                    duplicateVoteEvidence?: {
                        voteA?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        voteB?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        totalVotingPower?: bigint | undefined;
                        validatorPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                    lightClientAttackEvidence?: {
                        conflictingBlock?: {
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                        } | undefined;
                        commonHeight?: bigint | undefined;
                        byzantineValidators?: {
                            address?: Uint8Array | undefined;
                            pubKey?: {
                                ed25519?: Uint8Array | undefined;
                                secp256k1?: Uint8Array | undefined;
                            } | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        }[] | undefined;
                        totalVotingPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                }[] | undefined;
            } | undefined;
            lastCommit?: {
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
                    blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    timestamp?: {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } | undefined;
                    signature?: Uint8Array | undefined;
                }[] | undefined;
            } | undefined;
        } | undefined;
        sdkBlock?: {
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
                proposerAddress?: string | undefined;
            } | undefined;
            data?: {
                txs?: Uint8Array[] | undefined;
            } | undefined;
            evidence?: {
                evidence?: {
                    duplicateVoteEvidence?: {
                        voteA?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        voteB?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        totalVotingPower?: bigint | undefined;
                        validatorPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                    lightClientAttackEvidence?: {
                        conflictingBlock?: {
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                        } | undefined;
                        commonHeight?: bigint | undefined;
                        byzantineValidators?: {
                            address?: Uint8Array | undefined;
                            pubKey?: {
                                ed25519?: Uint8Array | undefined;
                                secp256k1?: Uint8Array | undefined;
                            } | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        }[] | undefined;
                        totalVotingPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                }[] | undefined;
            } | undefined;
            lastCommit?: {
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
                    blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    timestamp?: {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } | undefined;
                    signature?: Uint8Array | undefined;
                }[] | undefined;
            } | undefined;
        } | undefined;
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
            } & Record<Exclude<keyof I["blockId"]["partSetHeader"], keyof import("../../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
        } & Record<Exclude<keyof I["blockId"], keyof BlockID>, never>) | undefined;
        block?: ({
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
            data?: {
                txs?: Uint8Array[] | undefined;
            } | undefined;
            evidence?: {
                evidence?: {
                    duplicateVoteEvidence?: {
                        voteA?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        voteB?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        totalVotingPower?: bigint | undefined;
                        validatorPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                    lightClientAttackEvidence?: {
                        conflictingBlock?: {
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                        } | undefined;
                        commonHeight?: bigint | undefined;
                        byzantineValidators?: {
                            address?: Uint8Array | undefined;
                            pubKey?: {
                                ed25519?: Uint8Array | undefined;
                                secp256k1?: Uint8Array | undefined;
                            } | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        }[] | undefined;
                        totalVotingPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                }[] | undefined;
            } | undefined;
            lastCommit?: {
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
                    blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                } & Record<Exclude<keyof I["block"]["header"]["version"], keyof import("../../../../tendermint/version/types").Consensus>, never>) | undefined;
                chainId?: string | undefined;
                height?: bigint | undefined;
                time?: ({
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & Record<Exclude<keyof I["block"]["header"]["time"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
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
                    } & Record<Exclude<keyof I["block"]["header"]["lastBlockId"]["partSetHeader"], keyof import("../../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                } & Record<Exclude<keyof I["block"]["header"]["lastBlockId"], keyof BlockID>, never>) | undefined;
                lastCommitHash?: Uint8Array | undefined;
                dataHash?: Uint8Array | undefined;
                validatorsHash?: Uint8Array | undefined;
                nextValidatorsHash?: Uint8Array | undefined;
                consensusHash?: Uint8Array | undefined;
                appHash?: Uint8Array | undefined;
                lastResultsHash?: Uint8Array | undefined;
                evidenceHash?: Uint8Array | undefined;
                proposerAddress?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["block"]["header"], keyof import("../../../../tendermint/types/types").Header>, never>) | undefined;
            data?: ({
                txs?: Uint8Array[] | undefined;
            } & {
                txs?: (Uint8Array[] & Uint8Array[] & Record<Exclude<keyof I["block"]["data"]["txs"], keyof Uint8Array[]>, never>) | undefined;
            } & Record<Exclude<keyof I["block"]["data"], "txs">, never>) | undefined;
            evidence?: ({
                evidence?: {
                    duplicateVoteEvidence?: {
                        voteA?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        voteB?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        totalVotingPower?: bigint | undefined;
                        validatorPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                    lightClientAttackEvidence?: {
                        conflictingBlock?: {
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                        } | undefined;
                        commonHeight?: bigint | undefined;
                        byzantineValidators?: {
                            address?: Uint8Array | undefined;
                            pubKey?: {
                                ed25519?: Uint8Array | undefined;
                                secp256k1?: Uint8Array | undefined;
                            } | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        }[] | undefined;
                        totalVotingPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                }[] | undefined;
            } & {
                evidence?: ({
                    duplicateVoteEvidence?: {
                        voteA?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        voteB?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        totalVotingPower?: bigint | undefined;
                        validatorPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                    lightClientAttackEvidence?: {
                        conflictingBlock?: {
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                        } | undefined;
                        commonHeight?: bigint | undefined;
                        byzantineValidators?: {
                            address?: Uint8Array | undefined;
                            pubKey?: {
                                ed25519?: Uint8Array | undefined;
                                secp256k1?: Uint8Array | undefined;
                            } | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        }[] | undefined;
                        totalVotingPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                }[] & ({
                    duplicateVoteEvidence?: {
                        voteA?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        voteB?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        totalVotingPower?: bigint | undefined;
                        validatorPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                    lightClientAttackEvidence?: {
                        conflictingBlock?: {
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                        } | undefined;
                        commonHeight?: bigint | undefined;
                        byzantineValidators?: {
                            address?: Uint8Array | undefined;
                            pubKey?: {
                                ed25519?: Uint8Array | undefined;
                                secp256k1?: Uint8Array | undefined;
                            } | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        }[] | undefined;
                        totalVotingPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                } & {
                    duplicateVoteEvidence?: ({
                        voteA?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        voteB?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        totalVotingPower?: bigint | undefined;
                        validatorPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } & {
                        voteA?: ({
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                                } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteA"]["blockId"]["partSetHeader"], keyof import("../../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                            } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteA"]["blockId"], keyof BlockID>, never>) | undefined;
                            timestamp?: ({
                                seconds?: bigint | undefined;
                                nanos?: number | undefined;
                            } & {
                                seconds?: bigint | undefined;
                                nanos?: number | undefined;
                            } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteA"]["timestamp"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                            validatorAddress?: Uint8Array | undefined;
                            validatorIndex?: number | undefined;
                            signature?: Uint8Array | undefined;
                        } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteA"], keyof import("../../../../tendermint/types/types").Vote>, never>) | undefined;
                        voteB?: ({
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                                } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteB"]["blockId"]["partSetHeader"], keyof import("../../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                            } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteB"]["blockId"], keyof BlockID>, never>) | undefined;
                            timestamp?: ({
                                seconds?: bigint | undefined;
                                nanos?: number | undefined;
                            } & {
                                seconds?: bigint | undefined;
                                nanos?: number | undefined;
                            } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteB"]["timestamp"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                            validatorAddress?: Uint8Array | undefined;
                            validatorIndex?: number | undefined;
                            signature?: Uint8Array | undefined;
                        } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteB"], keyof import("../../../../tendermint/types/types").Vote>, never>) | undefined;
                        totalVotingPower?: bigint | undefined;
                        validatorPower?: bigint | undefined;
                        timestamp?: ({
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } & {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["timestamp"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                    } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["duplicateVoteEvidence"], keyof import("../../../../tendermint/types/evidence").DuplicateVoteEvidence>, never>) | undefined;
                    lightClientAttackEvidence?: ({
                        conflictingBlock?: {
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                        } | undefined;
                        commonHeight?: bigint | undefined;
                        byzantineValidators?: {
                            address?: Uint8Array | undefined;
                            pubKey?: {
                                ed25519?: Uint8Array | undefined;
                                secp256k1?: Uint8Array | undefined;
                            } | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        }[] | undefined;
                        totalVotingPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } & {
                        conflictingBlock?: ({
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                                    } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"]["version"], keyof import("../../../../tendermint/version/types").Consensus>, never>) | undefined;
                                    chainId?: string | undefined;
                                    height?: bigint | undefined;
                                    time?: ({
                                        seconds?: bigint | undefined;
                                        nanos?: number | undefined;
                                    } & {
                                        seconds?: bigint | undefined;
                                        nanos?: number | undefined;
                                    } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"]["time"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
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
                                        } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"]["lastBlockId"]["partSetHeader"], keyof import("../../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                                    } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"]["lastBlockId"], keyof BlockID>, never>) | undefined;
                                    lastCommitHash?: Uint8Array | undefined;
                                    dataHash?: Uint8Array | undefined;
                                    validatorsHash?: Uint8Array | undefined;
                                    nextValidatorsHash?: Uint8Array | undefined;
                                    consensusHash?: Uint8Array | undefined;
                                    appHash?: Uint8Array | undefined;
                                    lastResultsHash?: Uint8Array | undefined;
                                    evidenceHash?: Uint8Array | undefined;
                                    proposerAddress?: Uint8Array | undefined;
                                } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"], keyof import("../../../../tendermint/types/types").Header>, never>) | undefined;
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                                        } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["blockId"]["partSetHeader"], keyof import("../../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                                    } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["blockId"], keyof BlockID>, never>) | undefined;
                                    signatures?: ({
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                                        validatorAddress?: Uint8Array | undefined;
                                        timestamp?: {
                                            seconds?: bigint | undefined;
                                            nanos?: number | undefined;
                                        } | undefined;
                                        signature?: Uint8Array | undefined;
                                    }[] & ({
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                                        validatorAddress?: Uint8Array | undefined;
                                        timestamp?: {
                                            seconds?: bigint | undefined;
                                            nanos?: number | undefined;
                                        } | undefined;
                                        signature?: Uint8Array | undefined;
                                    } & {
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                                        validatorAddress?: Uint8Array | undefined;
                                        timestamp?: ({
                                            seconds?: bigint | undefined;
                                            nanos?: number | undefined;
                                        } & {
                                            seconds?: bigint | undefined;
                                            nanos?: number | undefined;
                                        } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["signatures"][number]["timestamp"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                                        signature?: Uint8Array | undefined;
                                    } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["signatures"][number], keyof import("../../../../tendermint/types/types").CommitSig>, never>)[] & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["signatures"], keyof {
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                                        validatorAddress?: Uint8Array | undefined;
                                        timestamp?: {
                                            seconds?: bigint | undefined;
                                            nanos?: number | undefined;
                                        } | undefined;
                                        signature?: Uint8Array | undefined;
                                    }[]>, never>) | undefined;
                                } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"], keyof import("../../../../tendermint/types/types").Commit>, never>) | undefined;
                            } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"], keyof import("../../../../tendermint/types/types").SignedHeader>, never>) | undefined;
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
                                    } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["validators"][number]["pubKey"], keyof import("../../../../tendermint/crypto/keys").PublicKey>, never>) | undefined;
                                    votingPower?: bigint | undefined;
                                    proposerPriority?: bigint | undefined;
                                } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["validators"][number], keyof import("../../../../tendermint/types/validator").Validator>, never>)[] & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["validators"], keyof {
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
                                    } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["proposer"]["pubKey"], keyof import("../../../../tendermint/crypto/keys").PublicKey>, never>) | undefined;
                                    votingPower?: bigint | undefined;
                                    proposerPriority?: bigint | undefined;
                                } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["proposer"], keyof import("../../../../tendermint/types/validator").Validator>, never>) | undefined;
                                totalVotingPower?: bigint | undefined;
                            } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"], keyof import("../../../../tendermint/types/validator").ValidatorSet>, never>) | undefined;
                        } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"], keyof import("../../../../tendermint/types/types").LightBlock>, never>) | undefined;
                        commonHeight?: bigint | undefined;
                        byzantineValidators?: ({
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
                            } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["byzantineValidators"][number]["pubKey"], keyof import("../../../../tendermint/crypto/keys").PublicKey>, never>) | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["byzantineValidators"][number], keyof import("../../../../tendermint/types/validator").Validator>, never>)[] & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["byzantineValidators"], keyof {
                            address?: Uint8Array | undefined;
                            pubKey?: {
                                ed25519?: Uint8Array | undefined;
                                secp256k1?: Uint8Array | undefined;
                            } | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        }[]>, never>) | undefined;
                        totalVotingPower?: bigint | undefined;
                        timestamp?: ({
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } & {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["timestamp"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                    } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"], keyof import("../../../../tendermint/types/evidence").LightClientAttackEvidence>, never>) | undefined;
                } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number], keyof import("../../../../tendermint/types/evidence").Evidence>, never>)[] & Record<Exclude<keyof I["block"]["evidence"]["evidence"], keyof {
                    duplicateVoteEvidence?: {
                        voteA?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        voteB?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        totalVotingPower?: bigint | undefined;
                        validatorPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                    lightClientAttackEvidence?: {
                        conflictingBlock?: {
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                        } | undefined;
                        commonHeight?: bigint | undefined;
                        byzantineValidators?: {
                            address?: Uint8Array | undefined;
                            pubKey?: {
                                ed25519?: Uint8Array | undefined;
                                secp256k1?: Uint8Array | undefined;
                            } | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        }[] | undefined;
                        totalVotingPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["block"]["evidence"], "evidence">, never>) | undefined;
            lastCommit?: ({
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
                    blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                    } & Record<Exclude<keyof I["block"]["lastCommit"]["blockId"]["partSetHeader"], keyof import("../../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                } & Record<Exclude<keyof I["block"]["lastCommit"]["blockId"], keyof BlockID>, never>) | undefined;
                signatures?: ({
                    blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    timestamp?: {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } | undefined;
                    signature?: Uint8Array | undefined;
                }[] & ({
                    blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    timestamp?: {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } | undefined;
                    signature?: Uint8Array | undefined;
                } & {
                    blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    timestamp?: ({
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } & {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } & Record<Exclude<keyof I["block"]["lastCommit"]["signatures"][number]["timestamp"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                    signature?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["block"]["lastCommit"]["signatures"][number], keyof import("../../../../tendermint/types/types").CommitSig>, never>)[] & Record<Exclude<keyof I["block"]["lastCommit"]["signatures"], keyof {
                    blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    timestamp?: {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } | undefined;
                    signature?: Uint8Array | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["block"]["lastCommit"], keyof import("../../../../tendermint/types/types").Commit>, never>) | undefined;
        } & Record<Exclude<keyof I["block"], keyof Block1>, never>) | undefined;
        sdkBlock?: ({
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
                proposerAddress?: string | undefined;
            } | undefined;
            data?: {
                txs?: Uint8Array[] | undefined;
            } | undefined;
            evidence?: {
                evidence?: {
                    duplicateVoteEvidence?: {
                        voteA?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        voteB?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        totalVotingPower?: bigint | undefined;
                        validatorPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                    lightClientAttackEvidence?: {
                        conflictingBlock?: {
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                        } | undefined;
                        commonHeight?: bigint | undefined;
                        byzantineValidators?: {
                            address?: Uint8Array | undefined;
                            pubKey?: {
                                ed25519?: Uint8Array | undefined;
                                secp256k1?: Uint8Array | undefined;
                            } | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        }[] | undefined;
                        totalVotingPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                }[] | undefined;
            } | undefined;
            lastCommit?: {
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
                    blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                proposerAddress?: string | undefined;
            } & {
                version?: ({
                    block?: bigint | undefined;
                    app?: bigint | undefined;
                } & {
                    block?: bigint | undefined;
                    app?: bigint | undefined;
                } & Record<Exclude<keyof I["sdkBlock"]["header"]["version"], keyof import("../../../../tendermint/version/types").Consensus>, never>) | undefined;
                chainId?: string | undefined;
                height?: bigint | undefined;
                time?: ({
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & Record<Exclude<keyof I["sdkBlock"]["header"]["time"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
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
                    } & Record<Exclude<keyof I["sdkBlock"]["header"]["lastBlockId"]["partSetHeader"], keyof import("../../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                } & Record<Exclude<keyof I["sdkBlock"]["header"]["lastBlockId"], keyof BlockID>, never>) | undefined;
                lastCommitHash?: Uint8Array | undefined;
                dataHash?: Uint8Array | undefined;
                validatorsHash?: Uint8Array | undefined;
                nextValidatorsHash?: Uint8Array | undefined;
                consensusHash?: Uint8Array | undefined;
                appHash?: Uint8Array | undefined;
                lastResultsHash?: Uint8Array | undefined;
                evidenceHash?: Uint8Array | undefined;
                proposerAddress?: string | undefined;
            } & Record<Exclude<keyof I["sdkBlock"]["header"], keyof import("./types").Header>, never>) | undefined;
            data?: ({
                txs?: Uint8Array[] | undefined;
            } & {
                txs?: (Uint8Array[] & Uint8Array[] & Record<Exclude<keyof I["sdkBlock"]["data"]["txs"], keyof Uint8Array[]>, never>) | undefined;
            } & Record<Exclude<keyof I["sdkBlock"]["data"], "txs">, never>) | undefined;
            evidence?: ({
                evidence?: {
                    duplicateVoteEvidence?: {
                        voteA?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        voteB?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        totalVotingPower?: bigint | undefined;
                        validatorPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                    lightClientAttackEvidence?: {
                        conflictingBlock?: {
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                        } | undefined;
                        commonHeight?: bigint | undefined;
                        byzantineValidators?: {
                            address?: Uint8Array | undefined;
                            pubKey?: {
                                ed25519?: Uint8Array | undefined;
                                secp256k1?: Uint8Array | undefined;
                            } | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        }[] | undefined;
                        totalVotingPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                }[] | undefined;
            } & {
                evidence?: ({
                    duplicateVoteEvidence?: {
                        voteA?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        voteB?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        totalVotingPower?: bigint | undefined;
                        validatorPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                    lightClientAttackEvidence?: {
                        conflictingBlock?: {
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                        } | undefined;
                        commonHeight?: bigint | undefined;
                        byzantineValidators?: {
                            address?: Uint8Array | undefined;
                            pubKey?: {
                                ed25519?: Uint8Array | undefined;
                                secp256k1?: Uint8Array | undefined;
                            } | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        }[] | undefined;
                        totalVotingPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                }[] & ({
                    duplicateVoteEvidence?: {
                        voteA?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        voteB?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        totalVotingPower?: bigint | undefined;
                        validatorPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                    lightClientAttackEvidence?: {
                        conflictingBlock?: {
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                        } | undefined;
                        commonHeight?: bigint | undefined;
                        byzantineValidators?: {
                            address?: Uint8Array | undefined;
                            pubKey?: {
                                ed25519?: Uint8Array | undefined;
                                secp256k1?: Uint8Array | undefined;
                            } | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        }[] | undefined;
                        totalVotingPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                } & {
                    duplicateVoteEvidence?: ({
                        voteA?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        voteB?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        totalVotingPower?: bigint | undefined;
                        validatorPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } & {
                        voteA?: ({
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                                } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteA"]["blockId"]["partSetHeader"], keyof import("../../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                            } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteA"]["blockId"], keyof BlockID>, never>) | undefined;
                            timestamp?: ({
                                seconds?: bigint | undefined;
                                nanos?: number | undefined;
                            } & {
                                seconds?: bigint | undefined;
                                nanos?: number | undefined;
                            } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteA"]["timestamp"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                            validatorAddress?: Uint8Array | undefined;
                            validatorIndex?: number | undefined;
                            signature?: Uint8Array | undefined;
                        } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteA"], keyof import("../../../../tendermint/types/types").Vote>, never>) | undefined;
                        voteB?: ({
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                                } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteB"]["blockId"]["partSetHeader"], keyof import("../../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                            } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteB"]["blockId"], keyof BlockID>, never>) | undefined;
                            timestamp?: ({
                                seconds?: bigint | undefined;
                                nanos?: number | undefined;
                            } & {
                                seconds?: bigint | undefined;
                                nanos?: number | undefined;
                            } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteB"]["timestamp"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                            validatorAddress?: Uint8Array | undefined;
                            validatorIndex?: number | undefined;
                            signature?: Uint8Array | undefined;
                        } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteB"], keyof import("../../../../tendermint/types/types").Vote>, never>) | undefined;
                        totalVotingPower?: bigint | undefined;
                        validatorPower?: bigint | undefined;
                        timestamp?: ({
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } & {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["timestamp"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                    } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["duplicateVoteEvidence"], keyof import("../../../../tendermint/types/evidence").DuplicateVoteEvidence>, never>) | undefined;
                    lightClientAttackEvidence?: ({
                        conflictingBlock?: {
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                        } | undefined;
                        commonHeight?: bigint | undefined;
                        byzantineValidators?: {
                            address?: Uint8Array | undefined;
                            pubKey?: {
                                ed25519?: Uint8Array | undefined;
                                secp256k1?: Uint8Array | undefined;
                            } | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        }[] | undefined;
                        totalVotingPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } & {
                        conflictingBlock?: ({
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                                    } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"]["version"], keyof import("../../../../tendermint/version/types").Consensus>, never>) | undefined;
                                    chainId?: string | undefined;
                                    height?: bigint | undefined;
                                    time?: ({
                                        seconds?: bigint | undefined;
                                        nanos?: number | undefined;
                                    } & {
                                        seconds?: bigint | undefined;
                                        nanos?: number | undefined;
                                    } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"]["time"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
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
                                        } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"]["lastBlockId"]["partSetHeader"], keyof import("../../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                                    } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"]["lastBlockId"], keyof BlockID>, never>) | undefined;
                                    lastCommitHash?: Uint8Array | undefined;
                                    dataHash?: Uint8Array | undefined;
                                    validatorsHash?: Uint8Array | undefined;
                                    nextValidatorsHash?: Uint8Array | undefined;
                                    consensusHash?: Uint8Array | undefined;
                                    appHash?: Uint8Array | undefined;
                                    lastResultsHash?: Uint8Array | undefined;
                                    evidenceHash?: Uint8Array | undefined;
                                    proposerAddress?: Uint8Array | undefined;
                                } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"], keyof import("../../../../tendermint/types/types").Header>, never>) | undefined;
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                                        } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["blockId"]["partSetHeader"], keyof import("../../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                                    } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["blockId"], keyof BlockID>, never>) | undefined;
                                    signatures?: ({
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                                        validatorAddress?: Uint8Array | undefined;
                                        timestamp?: {
                                            seconds?: bigint | undefined;
                                            nanos?: number | undefined;
                                        } | undefined;
                                        signature?: Uint8Array | undefined;
                                    }[] & ({
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                                        validatorAddress?: Uint8Array | undefined;
                                        timestamp?: {
                                            seconds?: bigint | undefined;
                                            nanos?: number | undefined;
                                        } | undefined;
                                        signature?: Uint8Array | undefined;
                                    } & {
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                                        validatorAddress?: Uint8Array | undefined;
                                        timestamp?: ({
                                            seconds?: bigint | undefined;
                                            nanos?: number | undefined;
                                        } & {
                                            seconds?: bigint | undefined;
                                            nanos?: number | undefined;
                                        } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["signatures"][number]["timestamp"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                                        signature?: Uint8Array | undefined;
                                    } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["signatures"][number], keyof import("../../../../tendermint/types/types").CommitSig>, never>)[] & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["signatures"], keyof {
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                                        validatorAddress?: Uint8Array | undefined;
                                        timestamp?: {
                                            seconds?: bigint | undefined;
                                            nanos?: number | undefined;
                                        } | undefined;
                                        signature?: Uint8Array | undefined;
                                    }[]>, never>) | undefined;
                                } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"], keyof import("../../../../tendermint/types/types").Commit>, never>) | undefined;
                            } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"], keyof import("../../../../tendermint/types/types").SignedHeader>, never>) | undefined;
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
                                    } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["validators"][number]["pubKey"], keyof import("../../../../tendermint/crypto/keys").PublicKey>, never>) | undefined;
                                    votingPower?: bigint | undefined;
                                    proposerPriority?: bigint | undefined;
                                } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["validators"][number], keyof import("../../../../tendermint/types/validator").Validator>, never>)[] & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["validators"], keyof {
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
                                    } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["proposer"]["pubKey"], keyof import("../../../../tendermint/crypto/keys").PublicKey>, never>) | undefined;
                                    votingPower?: bigint | undefined;
                                    proposerPriority?: bigint | undefined;
                                } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["proposer"], keyof import("../../../../tendermint/types/validator").Validator>, never>) | undefined;
                                totalVotingPower?: bigint | undefined;
                            } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"], keyof import("../../../../tendermint/types/validator").ValidatorSet>, never>) | undefined;
                        } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"], keyof import("../../../../tendermint/types/types").LightBlock>, never>) | undefined;
                        commonHeight?: bigint | undefined;
                        byzantineValidators?: ({
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
                            } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["byzantineValidators"][number]["pubKey"], keyof import("../../../../tendermint/crypto/keys").PublicKey>, never>) | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["byzantineValidators"][number], keyof import("../../../../tendermint/types/validator").Validator>, never>)[] & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["byzantineValidators"], keyof {
                            address?: Uint8Array | undefined;
                            pubKey?: {
                                ed25519?: Uint8Array | undefined;
                                secp256k1?: Uint8Array | undefined;
                            } | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        }[]>, never>) | undefined;
                        totalVotingPower?: bigint | undefined;
                        timestamp?: ({
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } & {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["timestamp"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                    } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number]["lightClientAttackEvidence"], keyof import("../../../../tendermint/types/evidence").LightClientAttackEvidence>, never>) | undefined;
                } & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"][number], keyof import("../../../../tendermint/types/evidence").Evidence>, never>)[] & Record<Exclude<keyof I["sdkBlock"]["evidence"]["evidence"], keyof {
                    duplicateVoteEvidence?: {
                        voteA?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        voteB?: {
                            type?: import("../../../../tendermint/types/types").SignedMsgType | undefined;
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
                        } | undefined;
                        totalVotingPower?: bigint | undefined;
                        validatorPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                    lightClientAttackEvidence?: {
                        conflictingBlock?: {
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
                                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                        } | undefined;
                        commonHeight?: bigint | undefined;
                        byzantineValidators?: {
                            address?: Uint8Array | undefined;
                            pubKey?: {
                                ed25519?: Uint8Array | undefined;
                                secp256k1?: Uint8Array | undefined;
                            } | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        }[] | undefined;
                        totalVotingPower?: bigint | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                    } | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["sdkBlock"]["evidence"], "evidence">, never>) | undefined;
            lastCommit?: ({
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
                    blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
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
                    } & Record<Exclude<keyof I["sdkBlock"]["lastCommit"]["blockId"]["partSetHeader"], keyof import("../../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                } & Record<Exclude<keyof I["sdkBlock"]["lastCommit"]["blockId"], keyof BlockID>, never>) | undefined;
                signatures?: ({
                    blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    timestamp?: {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } | undefined;
                    signature?: Uint8Array | undefined;
                }[] & ({
                    blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    timestamp?: {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } | undefined;
                    signature?: Uint8Array | undefined;
                } & {
                    blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    timestamp?: ({
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } & {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } & Record<Exclude<keyof I["sdkBlock"]["lastCommit"]["signatures"][number]["timestamp"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                    signature?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["sdkBlock"]["lastCommit"]["signatures"][number], keyof import("../../../../tendermint/types/types").CommitSig>, never>)[] & Record<Exclude<keyof I["sdkBlock"]["lastCommit"]["signatures"], keyof {
                    blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    timestamp?: {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } | undefined;
                    signature?: Uint8Array | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["sdkBlock"]["lastCommit"], keyof import("../../../../tendermint/types/types").Commit>, never>) | undefined;
        } & Record<Exclude<keyof I["sdkBlock"], keyof Block2>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof GetLatestBlockResponse>, never>>(object: I): GetLatestBlockResponse;
};
export declare const GetSyncingRequest: {
    typeUrl: string;
    encode(_: GetSyncingRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GetSyncingRequest;
    fromJSON(_: any): GetSyncingRequest;
    toJSON(_: GetSyncingRequest): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): GetSyncingRequest;
};
export declare const GetSyncingResponse: {
    typeUrl: string;
    encode(message: GetSyncingResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GetSyncingResponse;
    fromJSON(object: any): GetSyncingResponse;
    toJSON(message: GetSyncingResponse): unknown;
    fromPartial<I extends {
        syncing?: boolean | undefined;
    } & {
        syncing?: boolean | undefined;
    } & Record<Exclude<keyof I, "syncing">, never>>(object: I): GetSyncingResponse;
};
export declare const GetNodeInfoRequest: {
    typeUrl: string;
    encode(_: GetNodeInfoRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GetNodeInfoRequest;
    fromJSON(_: any): GetNodeInfoRequest;
    toJSON(_: GetNodeInfoRequest): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): GetNodeInfoRequest;
};
export declare const GetNodeInfoResponse: {
    typeUrl: string;
    encode(message: GetNodeInfoResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GetNodeInfoResponse;
    fromJSON(object: any): GetNodeInfoResponse;
    toJSON(message: GetNodeInfoResponse): unknown;
    fromPartial<I extends {
        defaultNodeInfo?: {
            protocolVersion?: {
                p2p?: bigint | undefined;
                block?: bigint | undefined;
                app?: bigint | undefined;
            } | undefined;
            defaultNodeId?: string | undefined;
            listenAddr?: string | undefined;
            network?: string | undefined;
            version?: string | undefined;
            channels?: Uint8Array | undefined;
            moniker?: string | undefined;
            other?: {
                txIndex?: string | undefined;
                rpcAddress?: string | undefined;
            } | undefined;
        } | undefined;
        applicationVersion?: {
            name?: string | undefined;
            appName?: string | undefined;
            version?: string | undefined;
            gitCommit?: string | undefined;
            buildTags?: string | undefined;
            goVersion?: string | undefined;
            buildDeps?: {
                path?: string | undefined;
                version?: string | undefined;
                sum?: string | undefined;
            }[] | undefined;
            cosmosSdkVersion?: string | undefined;
        } | undefined;
    } & {
        defaultNodeInfo?: ({
            protocolVersion?: {
                p2p?: bigint | undefined;
                block?: bigint | undefined;
                app?: bigint | undefined;
            } | undefined;
            defaultNodeId?: string | undefined;
            listenAddr?: string | undefined;
            network?: string | undefined;
            version?: string | undefined;
            channels?: Uint8Array | undefined;
            moniker?: string | undefined;
            other?: {
                txIndex?: string | undefined;
                rpcAddress?: string | undefined;
            } | undefined;
        } & {
            protocolVersion?: ({
                p2p?: bigint | undefined;
                block?: bigint | undefined;
                app?: bigint | undefined;
            } & {
                p2p?: bigint | undefined;
                block?: bigint | undefined;
                app?: bigint | undefined;
            } & Record<Exclude<keyof I["defaultNodeInfo"]["protocolVersion"], keyof import("../../../../tendermint/p2p/types").ProtocolVersion>, never>) | undefined;
            defaultNodeId?: string | undefined;
            listenAddr?: string | undefined;
            network?: string | undefined;
            version?: string | undefined;
            channels?: Uint8Array | undefined;
            moniker?: string | undefined;
            other?: ({
                txIndex?: string | undefined;
                rpcAddress?: string | undefined;
            } & {
                txIndex?: string | undefined;
                rpcAddress?: string | undefined;
            } & Record<Exclude<keyof I["defaultNodeInfo"]["other"], keyof import("../../../../tendermint/p2p/types").DefaultNodeInfoOther>, never>) | undefined;
        } & Record<Exclude<keyof I["defaultNodeInfo"], keyof DefaultNodeInfo>, never>) | undefined;
        applicationVersion?: ({
            name?: string | undefined;
            appName?: string | undefined;
            version?: string | undefined;
            gitCommit?: string | undefined;
            buildTags?: string | undefined;
            goVersion?: string | undefined;
            buildDeps?: {
                path?: string | undefined;
                version?: string | undefined;
                sum?: string | undefined;
            }[] | undefined;
            cosmosSdkVersion?: string | undefined;
        } & {
            name?: string | undefined;
            appName?: string | undefined;
            version?: string | undefined;
            gitCommit?: string | undefined;
            buildTags?: string | undefined;
            goVersion?: string | undefined;
            buildDeps?: ({
                path?: string | undefined;
                version?: string | undefined;
                sum?: string | undefined;
            }[] & ({
                path?: string | undefined;
                version?: string | undefined;
                sum?: string | undefined;
            } & {
                path?: string | undefined;
                version?: string | undefined;
                sum?: string | undefined;
            } & Record<Exclude<keyof I["applicationVersion"]["buildDeps"][number], keyof Module>, never>)[] & Record<Exclude<keyof I["applicationVersion"]["buildDeps"], keyof {
                path?: string | undefined;
                version?: string | undefined;
                sum?: string | undefined;
            }[]>, never>) | undefined;
            cosmosSdkVersion?: string | undefined;
        } & Record<Exclude<keyof I["applicationVersion"], keyof VersionInfo>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof GetNodeInfoResponse>, never>>(object: I): GetNodeInfoResponse;
};
export declare const VersionInfo: {
    typeUrl: string;
    encode(message: VersionInfo, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): VersionInfo;
    fromJSON(object: any): VersionInfo;
    toJSON(message: VersionInfo): unknown;
    fromPartial<I extends {
        name?: string | undefined;
        appName?: string | undefined;
        version?: string | undefined;
        gitCommit?: string | undefined;
        buildTags?: string | undefined;
        goVersion?: string | undefined;
        buildDeps?: {
            path?: string | undefined;
            version?: string | undefined;
            sum?: string | undefined;
        }[] | undefined;
        cosmosSdkVersion?: string | undefined;
    } & {
        name?: string | undefined;
        appName?: string | undefined;
        version?: string | undefined;
        gitCommit?: string | undefined;
        buildTags?: string | undefined;
        goVersion?: string | undefined;
        buildDeps?: ({
            path?: string | undefined;
            version?: string | undefined;
            sum?: string | undefined;
        }[] & ({
            path?: string | undefined;
            version?: string | undefined;
            sum?: string | undefined;
        } & {
            path?: string | undefined;
            version?: string | undefined;
            sum?: string | undefined;
        } & Record<Exclude<keyof I["buildDeps"][number], keyof Module>, never>)[] & Record<Exclude<keyof I["buildDeps"], keyof {
            path?: string | undefined;
            version?: string | undefined;
            sum?: string | undefined;
        }[]>, never>) | undefined;
        cosmosSdkVersion?: string | undefined;
    } & Record<Exclude<keyof I, keyof VersionInfo>, never>>(object: I): VersionInfo;
};
export declare const Module: {
    typeUrl: string;
    encode(message: Module, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Module;
    fromJSON(object: any): Module;
    toJSON(message: Module): unknown;
    fromPartial<I extends {
        path?: string | undefined;
        version?: string | undefined;
        sum?: string | undefined;
    } & {
        path?: string | undefined;
        version?: string | undefined;
        sum?: string | undefined;
    } & Record<Exclude<keyof I, keyof Module>, never>>(object: I): Module;
};
export declare const ABCIQueryRequest: {
    typeUrl: string;
    encode(message: ABCIQueryRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ABCIQueryRequest;
    fromJSON(object: any): ABCIQueryRequest;
    toJSON(message: ABCIQueryRequest): unknown;
    fromPartial<I extends {
        data?: Uint8Array | undefined;
        path?: string | undefined;
        height?: bigint | undefined;
        prove?: boolean | undefined;
    } & {
        data?: Uint8Array | undefined;
        path?: string | undefined;
        height?: bigint | undefined;
        prove?: boolean | undefined;
    } & Record<Exclude<keyof I, keyof ABCIQueryRequest>, never>>(object: I): ABCIQueryRequest;
};
export declare const ABCIQueryResponse: {
    typeUrl: string;
    encode(message: ABCIQueryResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ABCIQueryResponse;
    fromJSON(object: any): ABCIQueryResponse;
    toJSON(message: ABCIQueryResponse): unknown;
    fromPartial<I extends {
        code?: number | undefined;
        log?: string | undefined;
        info?: string | undefined;
        index?: bigint | undefined;
        key?: Uint8Array | undefined;
        value?: Uint8Array | undefined;
        proofOps?: {
            ops?: {
                type?: string | undefined;
                key?: Uint8Array | undefined;
                data?: Uint8Array | undefined;
            }[] | undefined;
        } | undefined;
        height?: bigint | undefined;
        codespace?: string | undefined;
    } & {
        code?: number | undefined;
        log?: string | undefined;
        info?: string | undefined;
        index?: bigint | undefined;
        key?: Uint8Array | undefined;
        value?: Uint8Array | undefined;
        proofOps?: ({
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
            } & Record<Exclude<keyof I["proofOps"]["ops"][number], keyof ProofOp>, never>)[] & Record<Exclude<keyof I["proofOps"]["ops"], keyof {
                type?: string | undefined;
                key?: Uint8Array | undefined;
                data?: Uint8Array | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["proofOps"], "ops">, never>) | undefined;
        height?: bigint | undefined;
        codespace?: string | undefined;
    } & Record<Exclude<keyof I, keyof ABCIQueryResponse>, never>>(object: I): ABCIQueryResponse;
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
/** Service defines the gRPC querier service for tendermint queries. */
export interface Service {
    /** GetNodeInfo queries the current node info. */
    GetNodeInfo(request?: GetNodeInfoRequest): Promise<GetNodeInfoResponse>;
    /** GetSyncing queries node syncing. */
    GetSyncing(request?: GetSyncingRequest): Promise<GetSyncingResponse>;
    /** GetLatestBlock returns the latest block. */
    GetLatestBlock(request?: GetLatestBlockRequest): Promise<GetLatestBlockResponse>;
    /** GetBlockByHeight queries block for given height. */
    GetBlockByHeight(request: GetBlockByHeightRequest): Promise<GetBlockByHeightResponse>;
    /** GetLatestValidatorSet queries latest validator-set. */
    GetLatestValidatorSet(request?: GetLatestValidatorSetRequest): Promise<GetLatestValidatorSetResponse>;
    /** GetValidatorSetByHeight queries validator-set at a given height. */
    GetValidatorSetByHeight(request: GetValidatorSetByHeightRequest): Promise<GetValidatorSetByHeightResponse>;
    /**
     * ABCIQuery defines a query handler that supports ABCI queries directly to the
     * application, bypassing Tendermint completely. The ABCI query must contain
     * a valid and supported path, including app, custom, p2p, and store.
     *
     * Since: cosmos-sdk 0.46
     */
    ABCIQuery(request: ABCIQueryRequest): Promise<ABCIQueryResponse>;
}
export declare class ServiceClientImpl implements Service {
    private readonly rpc;
    constructor(rpc: Rpc);
    GetNodeInfo(request?: GetNodeInfoRequest): Promise<GetNodeInfoResponse>;
    GetSyncing(request?: GetSyncingRequest): Promise<GetSyncingResponse>;
    GetLatestBlock(request?: GetLatestBlockRequest): Promise<GetLatestBlockResponse>;
    GetBlockByHeight(request: GetBlockByHeightRequest): Promise<GetBlockByHeightResponse>;
    GetLatestValidatorSet(request?: GetLatestValidatorSetRequest): Promise<GetLatestValidatorSetResponse>;
    GetValidatorSetByHeight(request: GetValidatorSetByHeightRequest): Promise<GetValidatorSetByHeightResponse>;
    ABCIQuery(request: ABCIQueryRequest): Promise<ABCIQueryResponse>;
}
