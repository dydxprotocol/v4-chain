import { Timestamp } from "../../google/protobuf/timestamp";
import { ConsensusParams } from "../types/params";
import { Header } from "../types/types";
import { ProofOps } from "../crypto/proof";
import { PublicKey } from "../crypto/keys";
import { BinaryReader, BinaryWriter } from "../../binary";
import { Rpc } from "../../helpers";
export declare const protobufPackage = "tendermint.abci";
export declare enum CheckTxType {
    NEW = 0,
    RECHECK = 1,
    UNRECOGNIZED = -1
}
export declare function checkTxTypeFromJSON(object: any): CheckTxType;
export declare function checkTxTypeToJSON(object: CheckTxType): string;
export declare enum ResponseOfferSnapshot_Result {
    /** UNKNOWN - Unknown result, abort all snapshot restoration */
    UNKNOWN = 0,
    /** ACCEPT - Snapshot accepted, apply chunks */
    ACCEPT = 1,
    /** ABORT - Abort all snapshot restoration */
    ABORT = 2,
    /** REJECT - Reject this specific snapshot, try others */
    REJECT = 3,
    /** REJECT_FORMAT - Reject all snapshots of this format, try others */
    REJECT_FORMAT = 4,
    /** REJECT_SENDER - Reject all snapshots from the sender(s), try others */
    REJECT_SENDER = 5,
    UNRECOGNIZED = -1
}
export declare function responseOfferSnapshot_ResultFromJSON(object: any): ResponseOfferSnapshot_Result;
export declare function responseOfferSnapshot_ResultToJSON(object: ResponseOfferSnapshot_Result): string;
export declare enum ResponseApplySnapshotChunk_Result {
    /** UNKNOWN - Unknown result, abort all snapshot restoration */
    UNKNOWN = 0,
    /** ACCEPT - Chunk successfully accepted */
    ACCEPT = 1,
    /** ABORT - Abort all snapshot restoration */
    ABORT = 2,
    /** RETRY - Retry chunk (combine with refetch and reject) */
    RETRY = 3,
    /** RETRY_SNAPSHOT - Retry snapshot (combine with refetch and reject) */
    RETRY_SNAPSHOT = 4,
    /** REJECT_SNAPSHOT - Reject this snapshot, try others */
    REJECT_SNAPSHOT = 5,
    UNRECOGNIZED = -1
}
export declare function responseApplySnapshotChunk_ResultFromJSON(object: any): ResponseApplySnapshotChunk_Result;
export declare function responseApplySnapshotChunk_ResultToJSON(object: ResponseApplySnapshotChunk_Result): string;
export declare enum ResponseProcessProposal_ProposalStatus {
    UNKNOWN = 0,
    ACCEPT = 1,
    REJECT = 2,
    UNRECOGNIZED = -1
}
export declare function responseProcessProposal_ProposalStatusFromJSON(object: any): ResponseProcessProposal_ProposalStatus;
export declare function responseProcessProposal_ProposalStatusToJSON(object: ResponseProcessProposal_ProposalStatus): string;
export declare enum MisbehaviorType {
    UNKNOWN = 0,
    DUPLICATE_VOTE = 1,
    LIGHT_CLIENT_ATTACK = 2,
    UNRECOGNIZED = -1
}
export declare function misbehaviorTypeFromJSON(object: any): MisbehaviorType;
export declare function misbehaviorTypeToJSON(object: MisbehaviorType): string;
export interface Request {
    echo?: RequestEcho;
    flush?: RequestFlush;
    info?: RequestInfo;
    initChain?: RequestInitChain;
    query?: RequestQuery;
    beginBlock?: RequestBeginBlock;
    checkTx?: RequestCheckTx;
    deliverTx?: RequestDeliverTx;
    endBlock?: RequestEndBlock;
    commit?: RequestCommit;
    listSnapshots?: RequestListSnapshots;
    offerSnapshot?: RequestOfferSnapshot;
    loadSnapshotChunk?: RequestLoadSnapshotChunk;
    applySnapshotChunk?: RequestApplySnapshotChunk;
    prepareProposal?: RequestPrepareProposal;
    processProposal?: RequestProcessProposal;
}
export interface RequestEcho {
    message: string;
}
export interface RequestFlush {
}
export interface RequestInfo {
    version: string;
    blockVersion: bigint;
    p2pVersion: bigint;
    abciVersion: string;
}
export interface RequestInitChain {
    time: Timestamp;
    chainId: string;
    consensusParams?: ConsensusParams;
    validators: ValidatorUpdate[];
    appStateBytes: Uint8Array;
    initialHeight: bigint;
}
export interface RequestQuery {
    data: Uint8Array;
    path: string;
    height: bigint;
    prove: boolean;
}
export interface RequestBeginBlock {
    hash: Uint8Array;
    header: Header;
    lastCommitInfo: CommitInfo;
    byzantineValidators: Misbehavior[];
}
export interface RequestCheckTx {
    tx: Uint8Array;
    type: CheckTxType;
}
export interface RequestDeliverTx {
    tx: Uint8Array;
}
export interface RequestEndBlock {
    height: bigint;
}
export interface RequestCommit {
}
/** lists available snapshots */
export interface RequestListSnapshots {
}
/** offers a snapshot to the application */
export interface RequestOfferSnapshot {
    /** snapshot offered by peers */
    snapshot?: Snapshot;
    /** light client-verified app hash for snapshot height */
    appHash: Uint8Array;
}
/** loads a snapshot chunk */
export interface RequestLoadSnapshotChunk {
    height: bigint;
    format: number;
    chunk: number;
}
/** Applies a snapshot chunk */
export interface RequestApplySnapshotChunk {
    index: number;
    chunk: Uint8Array;
    sender: string;
}
export interface RequestPrepareProposal {
    /** the modified transactions cannot exceed this size. */
    maxTxBytes: bigint;
    /**
     * txs is an array of transactions that will be included in a block,
     * sent to the app for possible modifications.
     */
    txs: Uint8Array[];
    localLastCommit: ExtendedCommitInfo;
    misbehavior: Misbehavior[];
    height: bigint;
    time: Timestamp;
    nextValidatorsHash: Uint8Array;
    /** address of the public key of the validator proposing the block. */
    proposerAddress: Uint8Array;
}
export interface RequestProcessProposal {
    txs: Uint8Array[];
    proposedLastCommit: CommitInfo;
    misbehavior: Misbehavior[];
    /** hash is the merkle root hash of the fields of the proposed block. */
    hash: Uint8Array;
    height: bigint;
    time: Timestamp;
    nextValidatorsHash: Uint8Array;
    /** address of the public key of the original proposer of the block. */
    proposerAddress: Uint8Array;
}
export interface Response {
    exception?: ResponseException;
    echo?: ResponseEcho;
    flush?: ResponseFlush;
    info?: ResponseInfo;
    initChain?: ResponseInitChain;
    query?: ResponseQuery;
    beginBlock?: ResponseBeginBlock;
    checkTx?: ResponseCheckTx;
    deliverTx?: ResponseDeliverTx;
    endBlock?: ResponseEndBlock;
    commit?: ResponseCommit;
    listSnapshots?: ResponseListSnapshots;
    offerSnapshot?: ResponseOfferSnapshot;
    loadSnapshotChunk?: ResponseLoadSnapshotChunk;
    applySnapshotChunk?: ResponseApplySnapshotChunk;
    prepareProposal?: ResponsePrepareProposal;
    processProposal?: ResponseProcessProposal;
}
/** nondeterministic */
export interface ResponseException {
    error: string;
}
export interface ResponseEcho {
    message: string;
}
export interface ResponseFlush {
}
export interface ResponseInfo {
    data: string;
    version: string;
    appVersion: bigint;
    lastBlockHeight: bigint;
    lastBlockAppHash: Uint8Array;
}
export interface ResponseInitChain {
    consensusParams?: ConsensusParams;
    validators: ValidatorUpdate[];
    appHash: Uint8Array;
}
export interface ResponseQuery {
    code: number;
    /** bytes data = 2; // use "value" instead. */
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
export interface ResponseBeginBlock {
    events: Event[];
}
export interface ResponseCheckTx {
    code: number;
    data: Uint8Array;
    /** nondeterministic */
    log: string;
    /** nondeterministic */
    info: string;
    gasWanted: bigint;
    gasUsed: bigint;
    events: Event[];
    codespace: string;
    sender: string;
    priority: bigint;
    /**
     * mempool_error is set by CometBFT.
     * ABCI applictions creating a ResponseCheckTX should not set mempool_error.
     */
    mempoolError: string;
}
export interface ResponseDeliverTx {
    code: number;
    data: Uint8Array;
    /** nondeterministic */
    log: string;
    /** nondeterministic */
    info: string;
    gasWanted: bigint;
    gasUsed: bigint;
    events: Event[];
    codespace: string;
}
export interface ResponseEndBlock {
    validatorUpdates: ValidatorUpdate[];
    consensusParamUpdates?: ConsensusParams;
    events: Event[];
}
export interface ResponseCommit {
    /** reserve 1 */
    data: Uint8Array;
    retainHeight: bigint;
}
export interface ResponseListSnapshots {
    snapshots: Snapshot[];
}
export interface ResponseOfferSnapshot {
    result: ResponseOfferSnapshot_Result;
}
export interface ResponseLoadSnapshotChunk {
    chunk: Uint8Array;
}
export interface ResponseApplySnapshotChunk {
    result: ResponseApplySnapshotChunk_Result;
    /** Chunks to refetch and reapply */
    refetchChunks: number[];
    /** Chunk senders to reject and ban */
    rejectSenders: string[];
}
export interface ResponsePrepareProposal {
    txs: Uint8Array[];
}
export interface ResponseProcessProposal {
    status: ResponseProcessProposal_ProposalStatus;
}
export interface CommitInfo {
    round: number;
    votes: VoteInfo[];
}
export interface ExtendedCommitInfo {
    /** The round at which the block proposer decided in the previous height. */
    round: number;
    /**
     * List of validators' addresses in the last validator set with their voting
     * information, including vote extensions.
     */
    votes: ExtendedVoteInfo[];
}
/**
 * Event allows application developers to attach additional information to
 * ResponseBeginBlock, ResponseEndBlock, ResponseCheckTx and ResponseDeliverTx.
 * Later, transactions may be queried using these events.
 */
export interface Event {
    type: string;
    attributes: EventAttribute[];
}
/** EventAttribute is a single key-value pair, associated with an event. */
export interface EventAttribute {
    key: string;
    value: string;
    /** nondeterministic */
    index: boolean;
}
/**
 * TxResult contains results of executing the transaction.
 *
 * One usage is indexing transaction results.
 */
export interface TxResult {
    height: bigint;
    index: number;
    tx: Uint8Array;
    result: ResponseDeliverTx;
}
/** Validator */
export interface Validator {
    /**
     * The first 20 bytes of SHA256(public key)
     * PubKey pub_key = 2 [(gogoproto.nullable)=false];
     */
    address: Uint8Array;
    /** The voting power */
    power: bigint;
}
/** ValidatorUpdate */
export interface ValidatorUpdate {
    pubKey: PublicKey;
    power: bigint;
}
/** VoteInfo */
export interface VoteInfo {
    validator: Validator;
    signedLastBlock: boolean;
}
export interface ExtendedVoteInfo {
    validator: Validator;
    signedLastBlock: boolean;
    /** Reserved for future use */
    voteExtension: Uint8Array;
}
export interface Misbehavior {
    type: MisbehaviorType;
    /** The offending validator */
    validator: Validator;
    /** The height when the offense occurred */
    height: bigint;
    /** The corresponding time where the offense occurred */
    time: Timestamp;
    /**
     * Total voting power of the validator set in case the ABCI application does
     * not store historical validators.
     * https://github.com/tendermint/tendermint/issues/4581
     */
    totalVotingPower: bigint;
}
export interface Snapshot {
    /** The height at which the snapshot was taken */
    height: bigint;
    /** The application-specific snapshot format */
    format: number;
    /** Number of chunks in the snapshot */
    chunks: number;
    /** Arbitrary snapshot hash, equal only if identical */
    hash: Uint8Array;
    /** Arbitrary application metadata */
    metadata: Uint8Array;
}
export declare const Request: {
    typeUrl: string;
    encode(message: Request, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Request;
    fromJSON(object: any): Request;
    toJSON(message: Request): unknown;
    fromPartial<I extends {
        echo?: {
            message?: string | undefined;
        } | undefined;
        flush?: {} | undefined;
        info?: {
            version?: string | undefined;
            blockVersion?: bigint | undefined;
            p2pVersion?: bigint | undefined;
            abciVersion?: string | undefined;
        } | undefined;
        initChain?: {
            time?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            chainId?: string | undefined;
            consensusParams?: {
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
            } | undefined;
            validators?: {
                pubKey?: {
                    ed25519?: Uint8Array | undefined;
                    secp256k1?: Uint8Array | undefined;
                } | undefined;
                power?: bigint | undefined;
            }[] | undefined;
            appStateBytes?: Uint8Array | undefined;
            initialHeight?: bigint | undefined;
        } | undefined;
        query?: {
            data?: Uint8Array | undefined;
            path?: string | undefined;
            height?: bigint | undefined;
            prove?: boolean | undefined;
        } | undefined;
        beginBlock?: {
            hash?: Uint8Array | undefined;
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
            lastCommitInfo?: {
                round?: number | undefined;
                votes?: {
                    validator?: {
                        address?: Uint8Array | undefined;
                        power?: bigint | undefined;
                    } | undefined;
                    signedLastBlock?: boolean | undefined;
                }[] | undefined;
            } | undefined;
            byzantineValidators?: {
                type?: MisbehaviorType | undefined;
                validator?: {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } | undefined;
                height?: bigint | undefined;
                time?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                totalVotingPower?: bigint | undefined;
            }[] | undefined;
        } | undefined;
        checkTx?: {
            tx?: Uint8Array | undefined;
            type?: CheckTxType | undefined;
        } | undefined;
        deliverTx?: {
            tx?: Uint8Array | undefined;
        } | undefined;
        endBlock?: {
            height?: bigint | undefined;
        } | undefined;
        commit?: {} | undefined;
        listSnapshots?: {} | undefined;
        offerSnapshot?: {
            snapshot?: {
                height?: bigint | undefined;
                format?: number | undefined;
                chunks?: number | undefined;
                hash?: Uint8Array | undefined;
                metadata?: Uint8Array | undefined;
            } | undefined;
            appHash?: Uint8Array | undefined;
        } | undefined;
        loadSnapshotChunk?: {
            height?: bigint | undefined;
            format?: number | undefined;
            chunk?: number | undefined;
        } | undefined;
        applySnapshotChunk?: {
            index?: number | undefined;
            chunk?: Uint8Array | undefined;
            sender?: string | undefined;
        } | undefined;
        prepareProposal?: {
            maxTxBytes?: bigint | undefined;
            txs?: Uint8Array[] | undefined;
            localLastCommit?: {
                round?: number | undefined;
                votes?: {
                    validator?: {
                        address?: Uint8Array | undefined;
                        power?: bigint | undefined;
                    } | undefined;
                    signedLastBlock?: boolean | undefined;
                    voteExtension?: Uint8Array | undefined;
                }[] | undefined;
            } | undefined;
            misbehavior?: {
                type?: MisbehaviorType | undefined;
                validator?: {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } | undefined;
                height?: bigint | undefined;
                time?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                totalVotingPower?: bigint | undefined;
            }[] | undefined;
            height?: bigint | undefined;
            time?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            nextValidatorsHash?: Uint8Array | undefined;
            proposerAddress?: Uint8Array | undefined;
        } | undefined;
        processProposal?: {
            txs?: Uint8Array[] | undefined;
            proposedLastCommit?: {
                round?: number | undefined;
                votes?: {
                    validator?: {
                        address?: Uint8Array | undefined;
                        power?: bigint | undefined;
                    } | undefined;
                    signedLastBlock?: boolean | undefined;
                }[] | undefined;
            } | undefined;
            misbehavior?: {
                type?: MisbehaviorType | undefined;
                validator?: {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } | undefined;
                height?: bigint | undefined;
                time?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                totalVotingPower?: bigint | undefined;
            }[] | undefined;
            hash?: Uint8Array | undefined;
            height?: bigint | undefined;
            time?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            nextValidatorsHash?: Uint8Array | undefined;
            proposerAddress?: Uint8Array | undefined;
        } | undefined;
    } & {
        echo?: ({
            message?: string | undefined;
        } & {
            message?: string | undefined;
        } & Record<Exclude<keyof I["echo"], "message">, never>) | undefined;
        flush?: ({} & {} & Record<Exclude<keyof I["flush"], never>, never>) | undefined;
        info?: ({
            version?: string | undefined;
            blockVersion?: bigint | undefined;
            p2pVersion?: bigint | undefined;
            abciVersion?: string | undefined;
        } & {
            version?: string | undefined;
            blockVersion?: bigint | undefined;
            p2pVersion?: bigint | undefined;
            abciVersion?: string | undefined;
        } & Record<Exclude<keyof I["info"], keyof RequestInfo>, never>) | undefined;
        initChain?: ({
            time?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            chainId?: string | undefined;
            consensusParams?: {
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
            } | undefined;
            validators?: {
                pubKey?: {
                    ed25519?: Uint8Array | undefined;
                    secp256k1?: Uint8Array | undefined;
                } | undefined;
                power?: bigint | undefined;
            }[] | undefined;
            appStateBytes?: Uint8Array | undefined;
            initialHeight?: bigint | undefined;
        } & {
            time?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["initChain"]["time"], keyof Timestamp>, never>) | undefined;
            chainId?: string | undefined;
            consensusParams?: ({
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
                } & Record<Exclude<keyof I["initChain"]["consensusParams"]["block"], keyof import("../types/params").BlockParams>, never>) | undefined;
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
                    } & Record<Exclude<keyof I["initChain"]["consensusParams"]["evidence"]["maxAgeDuration"], keyof import("../../google/protobuf/duration").Duration>, never>) | undefined;
                    maxBytes?: bigint | undefined;
                } & Record<Exclude<keyof I["initChain"]["consensusParams"]["evidence"], keyof import("../types/params").EvidenceParams>, never>) | undefined;
                validator?: ({
                    pubKeyTypes?: string[] | undefined;
                } & {
                    pubKeyTypes?: (string[] & string[] & Record<Exclude<keyof I["initChain"]["consensusParams"]["validator"]["pubKeyTypes"], keyof string[]>, never>) | undefined;
                } & Record<Exclude<keyof I["initChain"]["consensusParams"]["validator"], "pubKeyTypes">, never>) | undefined;
                version?: ({
                    app?: bigint | undefined;
                } & {
                    app?: bigint | undefined;
                } & Record<Exclude<keyof I["initChain"]["consensusParams"]["version"], "app">, never>) | undefined;
            } & Record<Exclude<keyof I["initChain"]["consensusParams"], keyof ConsensusParams>, never>) | undefined;
            validators?: ({
                pubKey?: {
                    ed25519?: Uint8Array | undefined;
                    secp256k1?: Uint8Array | undefined;
                } | undefined;
                power?: bigint | undefined;
            }[] & ({
                pubKey?: {
                    ed25519?: Uint8Array | undefined;
                    secp256k1?: Uint8Array | undefined;
                } | undefined;
                power?: bigint | undefined;
            } & {
                pubKey?: ({
                    ed25519?: Uint8Array | undefined;
                    secp256k1?: Uint8Array | undefined;
                } & {
                    ed25519?: Uint8Array | undefined;
                    secp256k1?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["initChain"]["validators"][number]["pubKey"], keyof PublicKey>, never>) | undefined;
                power?: bigint | undefined;
            } & Record<Exclude<keyof I["initChain"]["validators"][number], keyof ValidatorUpdate>, never>)[] & Record<Exclude<keyof I["initChain"]["validators"], keyof {
                pubKey?: {
                    ed25519?: Uint8Array | undefined;
                    secp256k1?: Uint8Array | undefined;
                } | undefined;
                power?: bigint | undefined;
            }[]>, never>) | undefined;
            appStateBytes?: Uint8Array | undefined;
            initialHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["initChain"], keyof RequestInitChain>, never>) | undefined;
        query?: ({
            data?: Uint8Array | undefined;
            path?: string | undefined;
            height?: bigint | undefined;
            prove?: boolean | undefined;
        } & {
            data?: Uint8Array | undefined;
            path?: string | undefined;
            height?: bigint | undefined;
            prove?: boolean | undefined;
        } & Record<Exclude<keyof I["query"], keyof RequestQuery>, never>) | undefined;
        beginBlock?: ({
            hash?: Uint8Array | undefined;
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
            lastCommitInfo?: {
                round?: number | undefined;
                votes?: {
                    validator?: {
                        address?: Uint8Array | undefined;
                        power?: bigint | undefined;
                    } | undefined;
                    signedLastBlock?: boolean | undefined;
                }[] | undefined;
            } | undefined;
            byzantineValidators?: {
                type?: MisbehaviorType | undefined;
                validator?: {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } | undefined;
                height?: bigint | undefined;
                time?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                totalVotingPower?: bigint | undefined;
            }[] | undefined;
        } & {
            hash?: Uint8Array | undefined;
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
                } & Record<Exclude<keyof I["beginBlock"]["header"]["version"], keyof import("../version/types").Consensus>, never>) | undefined;
                chainId?: string | undefined;
                height?: bigint | undefined;
                time?: ({
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & Record<Exclude<keyof I["beginBlock"]["header"]["time"], keyof Timestamp>, never>) | undefined;
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
                    } & Record<Exclude<keyof I["beginBlock"]["header"]["lastBlockId"]["partSetHeader"], keyof import("../types/types").PartSetHeader>, never>) | undefined;
                } & Record<Exclude<keyof I["beginBlock"]["header"]["lastBlockId"], keyof import("../types/types").BlockID>, never>) | undefined;
                lastCommitHash?: Uint8Array | undefined;
                dataHash?: Uint8Array | undefined;
                validatorsHash?: Uint8Array | undefined;
                nextValidatorsHash?: Uint8Array | undefined;
                consensusHash?: Uint8Array | undefined;
                appHash?: Uint8Array | undefined;
                lastResultsHash?: Uint8Array | undefined;
                evidenceHash?: Uint8Array | undefined;
                proposerAddress?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["beginBlock"]["header"], keyof Header>, never>) | undefined;
            lastCommitInfo?: ({
                round?: number | undefined;
                votes?: {
                    validator?: {
                        address?: Uint8Array | undefined;
                        power?: bigint | undefined;
                    } | undefined;
                    signedLastBlock?: boolean | undefined;
                }[] | undefined;
            } & {
                round?: number | undefined;
                votes?: ({
                    validator?: {
                        address?: Uint8Array | undefined;
                        power?: bigint | undefined;
                    } | undefined;
                    signedLastBlock?: boolean | undefined;
                }[] & ({
                    validator?: {
                        address?: Uint8Array | undefined;
                        power?: bigint | undefined;
                    } | undefined;
                    signedLastBlock?: boolean | undefined;
                } & {
                    validator?: ({
                        address?: Uint8Array | undefined;
                        power?: bigint | undefined;
                    } & {
                        address?: Uint8Array | undefined;
                        power?: bigint | undefined;
                    } & Record<Exclude<keyof I["beginBlock"]["lastCommitInfo"]["votes"][number]["validator"], keyof Validator>, never>) | undefined;
                    signedLastBlock?: boolean | undefined;
                } & Record<Exclude<keyof I["beginBlock"]["lastCommitInfo"]["votes"][number], keyof VoteInfo>, never>)[] & Record<Exclude<keyof I["beginBlock"]["lastCommitInfo"]["votes"], keyof {
                    validator?: {
                        address?: Uint8Array | undefined;
                        power?: bigint | undefined;
                    } | undefined;
                    signedLastBlock?: boolean | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["beginBlock"]["lastCommitInfo"], keyof CommitInfo>, never>) | undefined;
            byzantineValidators?: ({
                type?: MisbehaviorType | undefined;
                validator?: {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } | undefined;
                height?: bigint | undefined;
                time?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                totalVotingPower?: bigint | undefined;
            }[] & ({
                type?: MisbehaviorType | undefined;
                validator?: {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } | undefined;
                height?: bigint | undefined;
                time?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                totalVotingPower?: bigint | undefined;
            } & {
                type?: MisbehaviorType | undefined;
                validator?: ({
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } & {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } & Record<Exclude<keyof I["beginBlock"]["byzantineValidators"][number]["validator"], keyof Validator>, never>) | undefined;
                height?: bigint | undefined;
                time?: ({
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & Record<Exclude<keyof I["beginBlock"]["byzantineValidators"][number]["time"], keyof Timestamp>, never>) | undefined;
                totalVotingPower?: bigint | undefined;
            } & Record<Exclude<keyof I["beginBlock"]["byzantineValidators"][number], keyof Misbehavior>, never>)[] & Record<Exclude<keyof I["beginBlock"]["byzantineValidators"], keyof {
                type?: MisbehaviorType | undefined;
                validator?: {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } | undefined;
                height?: bigint | undefined;
                time?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                totalVotingPower?: bigint | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["beginBlock"], keyof RequestBeginBlock>, never>) | undefined;
        checkTx?: ({
            tx?: Uint8Array | undefined;
            type?: CheckTxType | undefined;
        } & {
            tx?: Uint8Array | undefined;
            type?: CheckTxType | undefined;
        } & Record<Exclude<keyof I["checkTx"], keyof RequestCheckTx>, never>) | undefined;
        deliverTx?: ({
            tx?: Uint8Array | undefined;
        } & {
            tx?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["deliverTx"], "tx">, never>) | undefined;
        endBlock?: ({
            height?: bigint | undefined;
        } & {
            height?: bigint | undefined;
        } & Record<Exclude<keyof I["endBlock"], "height">, never>) | undefined;
        commit?: ({} & {} & Record<Exclude<keyof I["commit"], never>, never>) | undefined;
        listSnapshots?: ({} & {} & Record<Exclude<keyof I["listSnapshots"], never>, never>) | undefined;
        offerSnapshot?: ({
            snapshot?: {
                height?: bigint | undefined;
                format?: number | undefined;
                chunks?: number | undefined;
                hash?: Uint8Array | undefined;
                metadata?: Uint8Array | undefined;
            } | undefined;
            appHash?: Uint8Array | undefined;
        } & {
            snapshot?: ({
                height?: bigint | undefined;
                format?: number | undefined;
                chunks?: number | undefined;
                hash?: Uint8Array | undefined;
                metadata?: Uint8Array | undefined;
            } & {
                height?: bigint | undefined;
                format?: number | undefined;
                chunks?: number | undefined;
                hash?: Uint8Array | undefined;
                metadata?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["offerSnapshot"]["snapshot"], keyof Snapshot>, never>) | undefined;
            appHash?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["offerSnapshot"], keyof RequestOfferSnapshot>, never>) | undefined;
        loadSnapshotChunk?: ({
            height?: bigint | undefined;
            format?: number | undefined;
            chunk?: number | undefined;
        } & {
            height?: bigint | undefined;
            format?: number | undefined;
            chunk?: number | undefined;
        } & Record<Exclude<keyof I["loadSnapshotChunk"], keyof RequestLoadSnapshotChunk>, never>) | undefined;
        applySnapshotChunk?: ({
            index?: number | undefined;
            chunk?: Uint8Array | undefined;
            sender?: string | undefined;
        } & {
            index?: number | undefined;
            chunk?: Uint8Array | undefined;
            sender?: string | undefined;
        } & Record<Exclude<keyof I["applySnapshotChunk"], keyof RequestApplySnapshotChunk>, never>) | undefined;
        prepareProposal?: ({
            maxTxBytes?: bigint | undefined;
            txs?: Uint8Array[] | undefined;
            localLastCommit?: {
                round?: number | undefined;
                votes?: {
                    validator?: {
                        address?: Uint8Array | undefined;
                        power?: bigint | undefined;
                    } | undefined;
                    signedLastBlock?: boolean | undefined;
                    voteExtension?: Uint8Array | undefined;
                }[] | undefined;
            } | undefined;
            misbehavior?: {
                type?: MisbehaviorType | undefined;
                validator?: {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } | undefined;
                height?: bigint | undefined;
                time?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                totalVotingPower?: bigint | undefined;
            }[] | undefined;
            height?: bigint | undefined;
            time?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            nextValidatorsHash?: Uint8Array | undefined;
            proposerAddress?: Uint8Array | undefined;
        } & {
            maxTxBytes?: bigint | undefined;
            txs?: (Uint8Array[] & Uint8Array[] & Record<Exclude<keyof I["prepareProposal"]["txs"], keyof Uint8Array[]>, never>) | undefined;
            localLastCommit?: ({
                round?: number | undefined;
                votes?: {
                    validator?: {
                        address?: Uint8Array | undefined;
                        power?: bigint | undefined;
                    } | undefined;
                    signedLastBlock?: boolean | undefined;
                    voteExtension?: Uint8Array | undefined;
                }[] | undefined;
            } & {
                round?: number | undefined;
                votes?: ({
                    validator?: {
                        address?: Uint8Array | undefined;
                        power?: bigint | undefined;
                    } | undefined;
                    signedLastBlock?: boolean | undefined;
                    voteExtension?: Uint8Array | undefined;
                }[] & ({
                    validator?: {
                        address?: Uint8Array | undefined;
                        power?: bigint | undefined;
                    } | undefined;
                    signedLastBlock?: boolean | undefined;
                    voteExtension?: Uint8Array | undefined;
                } & {
                    validator?: ({
                        address?: Uint8Array | undefined;
                        power?: bigint | undefined;
                    } & {
                        address?: Uint8Array | undefined;
                        power?: bigint | undefined;
                    } & Record<Exclude<keyof I["prepareProposal"]["localLastCommit"]["votes"][number]["validator"], keyof Validator>, never>) | undefined;
                    signedLastBlock?: boolean | undefined;
                    voteExtension?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["prepareProposal"]["localLastCommit"]["votes"][number], keyof ExtendedVoteInfo>, never>)[] & Record<Exclude<keyof I["prepareProposal"]["localLastCommit"]["votes"], keyof {
                    validator?: {
                        address?: Uint8Array | undefined;
                        power?: bigint | undefined;
                    } | undefined;
                    signedLastBlock?: boolean | undefined;
                    voteExtension?: Uint8Array | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["prepareProposal"]["localLastCommit"], keyof ExtendedCommitInfo>, never>) | undefined;
            misbehavior?: ({
                type?: MisbehaviorType | undefined;
                validator?: {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } | undefined;
                height?: bigint | undefined;
                time?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                totalVotingPower?: bigint | undefined;
            }[] & ({
                type?: MisbehaviorType | undefined;
                validator?: {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } | undefined;
                height?: bigint | undefined;
                time?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                totalVotingPower?: bigint | undefined;
            } & {
                type?: MisbehaviorType | undefined;
                validator?: ({
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } & {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } & Record<Exclude<keyof I["prepareProposal"]["misbehavior"][number]["validator"], keyof Validator>, never>) | undefined;
                height?: bigint | undefined;
                time?: ({
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & Record<Exclude<keyof I["prepareProposal"]["misbehavior"][number]["time"], keyof Timestamp>, never>) | undefined;
                totalVotingPower?: bigint | undefined;
            } & Record<Exclude<keyof I["prepareProposal"]["misbehavior"][number], keyof Misbehavior>, never>)[] & Record<Exclude<keyof I["prepareProposal"]["misbehavior"], keyof {
                type?: MisbehaviorType | undefined;
                validator?: {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } | undefined;
                height?: bigint | undefined;
                time?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                totalVotingPower?: bigint | undefined;
            }[]>, never>) | undefined;
            height?: bigint | undefined;
            time?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["prepareProposal"]["time"], keyof Timestamp>, never>) | undefined;
            nextValidatorsHash?: Uint8Array | undefined;
            proposerAddress?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["prepareProposal"], keyof RequestPrepareProposal>, never>) | undefined;
        processProposal?: ({
            txs?: Uint8Array[] | undefined;
            proposedLastCommit?: {
                round?: number | undefined;
                votes?: {
                    validator?: {
                        address?: Uint8Array | undefined;
                        power?: bigint | undefined;
                    } | undefined;
                    signedLastBlock?: boolean | undefined;
                }[] | undefined;
            } | undefined;
            misbehavior?: {
                type?: MisbehaviorType | undefined;
                validator?: {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } | undefined;
                height?: bigint | undefined;
                time?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                totalVotingPower?: bigint | undefined;
            }[] | undefined;
            hash?: Uint8Array | undefined;
            height?: bigint | undefined;
            time?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            nextValidatorsHash?: Uint8Array | undefined;
            proposerAddress?: Uint8Array | undefined;
        } & {
            txs?: (Uint8Array[] & Uint8Array[] & Record<Exclude<keyof I["processProposal"]["txs"], keyof Uint8Array[]>, never>) | undefined;
            proposedLastCommit?: ({
                round?: number | undefined;
                votes?: {
                    validator?: {
                        address?: Uint8Array | undefined;
                        power?: bigint | undefined;
                    } | undefined;
                    signedLastBlock?: boolean | undefined;
                }[] | undefined;
            } & {
                round?: number | undefined;
                votes?: ({
                    validator?: {
                        address?: Uint8Array | undefined;
                        power?: bigint | undefined;
                    } | undefined;
                    signedLastBlock?: boolean | undefined;
                }[] & ({
                    validator?: {
                        address?: Uint8Array | undefined;
                        power?: bigint | undefined;
                    } | undefined;
                    signedLastBlock?: boolean | undefined;
                } & {
                    validator?: ({
                        address?: Uint8Array | undefined;
                        power?: bigint | undefined;
                    } & {
                        address?: Uint8Array | undefined;
                        power?: bigint | undefined;
                    } & Record<Exclude<keyof I["processProposal"]["proposedLastCommit"]["votes"][number]["validator"], keyof Validator>, never>) | undefined;
                    signedLastBlock?: boolean | undefined;
                } & Record<Exclude<keyof I["processProposal"]["proposedLastCommit"]["votes"][number], keyof VoteInfo>, never>)[] & Record<Exclude<keyof I["processProposal"]["proposedLastCommit"]["votes"], keyof {
                    validator?: {
                        address?: Uint8Array | undefined;
                        power?: bigint | undefined;
                    } | undefined;
                    signedLastBlock?: boolean | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["processProposal"]["proposedLastCommit"], keyof CommitInfo>, never>) | undefined;
            misbehavior?: ({
                type?: MisbehaviorType | undefined;
                validator?: {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } | undefined;
                height?: bigint | undefined;
                time?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                totalVotingPower?: bigint | undefined;
            }[] & ({
                type?: MisbehaviorType | undefined;
                validator?: {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } | undefined;
                height?: bigint | undefined;
                time?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                totalVotingPower?: bigint | undefined;
            } & {
                type?: MisbehaviorType | undefined;
                validator?: ({
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } & {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } & Record<Exclude<keyof I["processProposal"]["misbehavior"][number]["validator"], keyof Validator>, never>) | undefined;
                height?: bigint | undefined;
                time?: ({
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & Record<Exclude<keyof I["processProposal"]["misbehavior"][number]["time"], keyof Timestamp>, never>) | undefined;
                totalVotingPower?: bigint | undefined;
            } & Record<Exclude<keyof I["processProposal"]["misbehavior"][number], keyof Misbehavior>, never>)[] & Record<Exclude<keyof I["processProposal"]["misbehavior"], keyof {
                type?: MisbehaviorType | undefined;
                validator?: {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } | undefined;
                height?: bigint | undefined;
                time?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                totalVotingPower?: bigint | undefined;
            }[]>, never>) | undefined;
            hash?: Uint8Array | undefined;
            height?: bigint | undefined;
            time?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["processProposal"]["time"], keyof Timestamp>, never>) | undefined;
            nextValidatorsHash?: Uint8Array | undefined;
            proposerAddress?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["processProposal"], keyof RequestProcessProposal>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof Request>, never>>(object: I): Request;
};
export declare const RequestEcho: {
    typeUrl: string;
    encode(message: RequestEcho, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): RequestEcho;
    fromJSON(object: any): RequestEcho;
    toJSON(message: RequestEcho): unknown;
    fromPartial<I extends {
        message?: string | undefined;
    } & {
        message?: string | undefined;
    } & Record<Exclude<keyof I, "message">, never>>(object: I): RequestEcho;
};
export declare const RequestFlush: {
    typeUrl: string;
    encode(_: RequestFlush, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): RequestFlush;
    fromJSON(_: any): RequestFlush;
    toJSON(_: RequestFlush): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): RequestFlush;
};
export declare const RequestInfo: {
    typeUrl: string;
    encode(message: RequestInfo, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): RequestInfo;
    fromJSON(object: any): RequestInfo;
    toJSON(message: RequestInfo): unknown;
    fromPartial<I extends {
        version?: string | undefined;
        blockVersion?: bigint | undefined;
        p2pVersion?: bigint | undefined;
        abciVersion?: string | undefined;
    } & {
        version?: string | undefined;
        blockVersion?: bigint | undefined;
        p2pVersion?: bigint | undefined;
        abciVersion?: string | undefined;
    } & Record<Exclude<keyof I, keyof RequestInfo>, never>>(object: I): RequestInfo;
};
export declare const RequestInitChain: {
    typeUrl: string;
    encode(message: RequestInitChain, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): RequestInitChain;
    fromJSON(object: any): RequestInitChain;
    toJSON(message: RequestInitChain): unknown;
    fromPartial<I extends {
        time?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
        chainId?: string | undefined;
        consensusParams?: {
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
        } | undefined;
        validators?: {
            pubKey?: {
                ed25519?: Uint8Array | undefined;
                secp256k1?: Uint8Array | undefined;
            } | undefined;
            power?: bigint | undefined;
        }[] | undefined;
        appStateBytes?: Uint8Array | undefined;
        initialHeight?: bigint | undefined;
    } & {
        time?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["time"], keyof Timestamp>, never>) | undefined;
        chainId?: string | undefined;
        consensusParams?: ({
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
            } & Record<Exclude<keyof I["consensusParams"]["block"], keyof import("../types/params").BlockParams>, never>) | undefined;
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
                } & Record<Exclude<keyof I["consensusParams"]["evidence"]["maxAgeDuration"], keyof import("../../google/protobuf/duration").Duration>, never>) | undefined;
                maxBytes?: bigint | undefined;
            } & Record<Exclude<keyof I["consensusParams"]["evidence"], keyof import("../types/params").EvidenceParams>, never>) | undefined;
            validator?: ({
                pubKeyTypes?: string[] | undefined;
            } & {
                pubKeyTypes?: (string[] & string[] & Record<Exclude<keyof I["consensusParams"]["validator"]["pubKeyTypes"], keyof string[]>, never>) | undefined;
            } & Record<Exclude<keyof I["consensusParams"]["validator"], "pubKeyTypes">, never>) | undefined;
            version?: ({
                app?: bigint | undefined;
            } & {
                app?: bigint | undefined;
            } & Record<Exclude<keyof I["consensusParams"]["version"], "app">, never>) | undefined;
        } & Record<Exclude<keyof I["consensusParams"], keyof ConsensusParams>, never>) | undefined;
        validators?: ({
            pubKey?: {
                ed25519?: Uint8Array | undefined;
                secp256k1?: Uint8Array | undefined;
            } | undefined;
            power?: bigint | undefined;
        }[] & ({
            pubKey?: {
                ed25519?: Uint8Array | undefined;
                secp256k1?: Uint8Array | undefined;
            } | undefined;
            power?: bigint | undefined;
        } & {
            pubKey?: ({
                ed25519?: Uint8Array | undefined;
                secp256k1?: Uint8Array | undefined;
            } & {
                ed25519?: Uint8Array | undefined;
                secp256k1?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["validators"][number]["pubKey"], keyof PublicKey>, never>) | undefined;
            power?: bigint | undefined;
        } & Record<Exclude<keyof I["validators"][number], keyof ValidatorUpdate>, never>)[] & Record<Exclude<keyof I["validators"], keyof {
            pubKey?: {
                ed25519?: Uint8Array | undefined;
                secp256k1?: Uint8Array | undefined;
            } | undefined;
            power?: bigint | undefined;
        }[]>, never>) | undefined;
        appStateBytes?: Uint8Array | undefined;
        initialHeight?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof RequestInitChain>, never>>(object: I): RequestInitChain;
};
export declare const RequestQuery: {
    typeUrl: string;
    encode(message: RequestQuery, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): RequestQuery;
    fromJSON(object: any): RequestQuery;
    toJSON(message: RequestQuery): unknown;
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
    } & Record<Exclude<keyof I, keyof RequestQuery>, never>>(object: I): RequestQuery;
};
export declare const RequestBeginBlock: {
    typeUrl: string;
    encode(message: RequestBeginBlock, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): RequestBeginBlock;
    fromJSON(object: any): RequestBeginBlock;
    toJSON(message: RequestBeginBlock): unknown;
    fromPartial<I extends {
        hash?: Uint8Array | undefined;
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
        lastCommitInfo?: {
            round?: number | undefined;
            votes?: {
                validator?: {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } | undefined;
                signedLastBlock?: boolean | undefined;
            }[] | undefined;
        } | undefined;
        byzantineValidators?: {
            type?: MisbehaviorType | undefined;
            validator?: {
                address?: Uint8Array | undefined;
                power?: bigint | undefined;
            } | undefined;
            height?: bigint | undefined;
            time?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            totalVotingPower?: bigint | undefined;
        }[] | undefined;
    } & {
        hash?: Uint8Array | undefined;
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
            } & Record<Exclude<keyof I["header"]["version"], keyof import("../version/types").Consensus>, never>) | undefined;
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
                } & Record<Exclude<keyof I["header"]["lastBlockId"]["partSetHeader"], keyof import("../types/types").PartSetHeader>, never>) | undefined;
            } & Record<Exclude<keyof I["header"]["lastBlockId"], keyof import("../types/types").BlockID>, never>) | undefined;
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
        lastCommitInfo?: ({
            round?: number | undefined;
            votes?: {
                validator?: {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } | undefined;
                signedLastBlock?: boolean | undefined;
            }[] | undefined;
        } & {
            round?: number | undefined;
            votes?: ({
                validator?: {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } | undefined;
                signedLastBlock?: boolean | undefined;
            }[] & ({
                validator?: {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } | undefined;
                signedLastBlock?: boolean | undefined;
            } & {
                validator?: ({
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } & {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } & Record<Exclude<keyof I["lastCommitInfo"]["votes"][number]["validator"], keyof Validator>, never>) | undefined;
                signedLastBlock?: boolean | undefined;
            } & Record<Exclude<keyof I["lastCommitInfo"]["votes"][number], keyof VoteInfo>, never>)[] & Record<Exclude<keyof I["lastCommitInfo"]["votes"], keyof {
                validator?: {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } | undefined;
                signedLastBlock?: boolean | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["lastCommitInfo"], keyof CommitInfo>, never>) | undefined;
        byzantineValidators?: ({
            type?: MisbehaviorType | undefined;
            validator?: {
                address?: Uint8Array | undefined;
                power?: bigint | undefined;
            } | undefined;
            height?: bigint | undefined;
            time?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            totalVotingPower?: bigint | undefined;
        }[] & ({
            type?: MisbehaviorType | undefined;
            validator?: {
                address?: Uint8Array | undefined;
                power?: bigint | undefined;
            } | undefined;
            height?: bigint | undefined;
            time?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            totalVotingPower?: bigint | undefined;
        } & {
            type?: MisbehaviorType | undefined;
            validator?: ({
                address?: Uint8Array | undefined;
                power?: bigint | undefined;
            } & {
                address?: Uint8Array | undefined;
                power?: bigint | undefined;
            } & Record<Exclude<keyof I["byzantineValidators"][number]["validator"], keyof Validator>, never>) | undefined;
            height?: bigint | undefined;
            time?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["byzantineValidators"][number]["time"], keyof Timestamp>, never>) | undefined;
            totalVotingPower?: bigint | undefined;
        } & Record<Exclude<keyof I["byzantineValidators"][number], keyof Misbehavior>, never>)[] & Record<Exclude<keyof I["byzantineValidators"], keyof {
            type?: MisbehaviorType | undefined;
            validator?: {
                address?: Uint8Array | undefined;
                power?: bigint | undefined;
            } | undefined;
            height?: bigint | undefined;
            time?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            totalVotingPower?: bigint | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof RequestBeginBlock>, never>>(object: I): RequestBeginBlock;
};
export declare const RequestCheckTx: {
    typeUrl: string;
    encode(message: RequestCheckTx, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): RequestCheckTx;
    fromJSON(object: any): RequestCheckTx;
    toJSON(message: RequestCheckTx): unknown;
    fromPartial<I extends {
        tx?: Uint8Array | undefined;
        type?: CheckTxType | undefined;
    } & {
        tx?: Uint8Array | undefined;
        type?: CheckTxType | undefined;
    } & Record<Exclude<keyof I, keyof RequestCheckTx>, never>>(object: I): RequestCheckTx;
};
export declare const RequestDeliverTx: {
    typeUrl: string;
    encode(message: RequestDeliverTx, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): RequestDeliverTx;
    fromJSON(object: any): RequestDeliverTx;
    toJSON(message: RequestDeliverTx): unknown;
    fromPartial<I extends {
        tx?: Uint8Array | undefined;
    } & {
        tx?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, "tx">, never>>(object: I): RequestDeliverTx;
};
export declare const RequestEndBlock: {
    typeUrl: string;
    encode(message: RequestEndBlock, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): RequestEndBlock;
    fromJSON(object: any): RequestEndBlock;
    toJSON(message: RequestEndBlock): unknown;
    fromPartial<I extends {
        height?: bigint | undefined;
    } & {
        height?: bigint | undefined;
    } & Record<Exclude<keyof I, "height">, never>>(object: I): RequestEndBlock;
};
export declare const RequestCommit: {
    typeUrl: string;
    encode(_: RequestCommit, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): RequestCommit;
    fromJSON(_: any): RequestCommit;
    toJSON(_: RequestCommit): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): RequestCommit;
};
export declare const RequestListSnapshots: {
    typeUrl: string;
    encode(_: RequestListSnapshots, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): RequestListSnapshots;
    fromJSON(_: any): RequestListSnapshots;
    toJSON(_: RequestListSnapshots): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): RequestListSnapshots;
};
export declare const RequestOfferSnapshot: {
    typeUrl: string;
    encode(message: RequestOfferSnapshot, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): RequestOfferSnapshot;
    fromJSON(object: any): RequestOfferSnapshot;
    toJSON(message: RequestOfferSnapshot): unknown;
    fromPartial<I extends {
        snapshot?: {
            height?: bigint | undefined;
            format?: number | undefined;
            chunks?: number | undefined;
            hash?: Uint8Array | undefined;
            metadata?: Uint8Array | undefined;
        } | undefined;
        appHash?: Uint8Array | undefined;
    } & {
        snapshot?: ({
            height?: bigint | undefined;
            format?: number | undefined;
            chunks?: number | undefined;
            hash?: Uint8Array | undefined;
            metadata?: Uint8Array | undefined;
        } & {
            height?: bigint | undefined;
            format?: number | undefined;
            chunks?: number | undefined;
            hash?: Uint8Array | undefined;
            metadata?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["snapshot"], keyof Snapshot>, never>) | undefined;
        appHash?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof RequestOfferSnapshot>, never>>(object: I): RequestOfferSnapshot;
};
export declare const RequestLoadSnapshotChunk: {
    typeUrl: string;
    encode(message: RequestLoadSnapshotChunk, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): RequestLoadSnapshotChunk;
    fromJSON(object: any): RequestLoadSnapshotChunk;
    toJSON(message: RequestLoadSnapshotChunk): unknown;
    fromPartial<I extends {
        height?: bigint | undefined;
        format?: number | undefined;
        chunk?: number | undefined;
    } & {
        height?: bigint | undefined;
        format?: number | undefined;
        chunk?: number | undefined;
    } & Record<Exclude<keyof I, keyof RequestLoadSnapshotChunk>, never>>(object: I): RequestLoadSnapshotChunk;
};
export declare const RequestApplySnapshotChunk: {
    typeUrl: string;
    encode(message: RequestApplySnapshotChunk, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): RequestApplySnapshotChunk;
    fromJSON(object: any): RequestApplySnapshotChunk;
    toJSON(message: RequestApplySnapshotChunk): unknown;
    fromPartial<I extends {
        index?: number | undefined;
        chunk?: Uint8Array | undefined;
        sender?: string | undefined;
    } & {
        index?: number | undefined;
        chunk?: Uint8Array | undefined;
        sender?: string | undefined;
    } & Record<Exclude<keyof I, keyof RequestApplySnapshotChunk>, never>>(object: I): RequestApplySnapshotChunk;
};
export declare const RequestPrepareProposal: {
    typeUrl: string;
    encode(message: RequestPrepareProposal, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): RequestPrepareProposal;
    fromJSON(object: any): RequestPrepareProposal;
    toJSON(message: RequestPrepareProposal): unknown;
    fromPartial<I extends {
        maxTxBytes?: bigint | undefined;
        txs?: Uint8Array[] | undefined;
        localLastCommit?: {
            round?: number | undefined;
            votes?: {
                validator?: {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } | undefined;
                signedLastBlock?: boolean | undefined;
                voteExtension?: Uint8Array | undefined;
            }[] | undefined;
        } | undefined;
        misbehavior?: {
            type?: MisbehaviorType | undefined;
            validator?: {
                address?: Uint8Array | undefined;
                power?: bigint | undefined;
            } | undefined;
            height?: bigint | undefined;
            time?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            totalVotingPower?: bigint | undefined;
        }[] | undefined;
        height?: bigint | undefined;
        time?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
        nextValidatorsHash?: Uint8Array | undefined;
        proposerAddress?: Uint8Array | undefined;
    } & {
        maxTxBytes?: bigint | undefined;
        txs?: (Uint8Array[] & Uint8Array[] & Record<Exclude<keyof I["txs"], keyof Uint8Array[]>, never>) | undefined;
        localLastCommit?: ({
            round?: number | undefined;
            votes?: {
                validator?: {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } | undefined;
                signedLastBlock?: boolean | undefined;
                voteExtension?: Uint8Array | undefined;
            }[] | undefined;
        } & {
            round?: number | undefined;
            votes?: ({
                validator?: {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } | undefined;
                signedLastBlock?: boolean | undefined;
                voteExtension?: Uint8Array | undefined;
            }[] & ({
                validator?: {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } | undefined;
                signedLastBlock?: boolean | undefined;
                voteExtension?: Uint8Array | undefined;
            } & {
                validator?: ({
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } & {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } & Record<Exclude<keyof I["localLastCommit"]["votes"][number]["validator"], keyof Validator>, never>) | undefined;
                signedLastBlock?: boolean | undefined;
                voteExtension?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["localLastCommit"]["votes"][number], keyof ExtendedVoteInfo>, never>)[] & Record<Exclude<keyof I["localLastCommit"]["votes"], keyof {
                validator?: {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } | undefined;
                signedLastBlock?: boolean | undefined;
                voteExtension?: Uint8Array | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["localLastCommit"], keyof ExtendedCommitInfo>, never>) | undefined;
        misbehavior?: ({
            type?: MisbehaviorType | undefined;
            validator?: {
                address?: Uint8Array | undefined;
                power?: bigint | undefined;
            } | undefined;
            height?: bigint | undefined;
            time?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            totalVotingPower?: bigint | undefined;
        }[] & ({
            type?: MisbehaviorType | undefined;
            validator?: {
                address?: Uint8Array | undefined;
                power?: bigint | undefined;
            } | undefined;
            height?: bigint | undefined;
            time?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            totalVotingPower?: bigint | undefined;
        } & {
            type?: MisbehaviorType | undefined;
            validator?: ({
                address?: Uint8Array | undefined;
                power?: bigint | undefined;
            } & {
                address?: Uint8Array | undefined;
                power?: bigint | undefined;
            } & Record<Exclude<keyof I["misbehavior"][number]["validator"], keyof Validator>, never>) | undefined;
            height?: bigint | undefined;
            time?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["misbehavior"][number]["time"], keyof Timestamp>, never>) | undefined;
            totalVotingPower?: bigint | undefined;
        } & Record<Exclude<keyof I["misbehavior"][number], keyof Misbehavior>, never>)[] & Record<Exclude<keyof I["misbehavior"], keyof {
            type?: MisbehaviorType | undefined;
            validator?: {
                address?: Uint8Array | undefined;
                power?: bigint | undefined;
            } | undefined;
            height?: bigint | undefined;
            time?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            totalVotingPower?: bigint | undefined;
        }[]>, never>) | undefined;
        height?: bigint | undefined;
        time?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["time"], keyof Timestamp>, never>) | undefined;
        nextValidatorsHash?: Uint8Array | undefined;
        proposerAddress?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof RequestPrepareProposal>, never>>(object: I): RequestPrepareProposal;
};
export declare const RequestProcessProposal: {
    typeUrl: string;
    encode(message: RequestProcessProposal, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): RequestProcessProposal;
    fromJSON(object: any): RequestProcessProposal;
    toJSON(message: RequestProcessProposal): unknown;
    fromPartial<I extends {
        txs?: Uint8Array[] | undefined;
        proposedLastCommit?: {
            round?: number | undefined;
            votes?: {
                validator?: {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } | undefined;
                signedLastBlock?: boolean | undefined;
            }[] | undefined;
        } | undefined;
        misbehavior?: {
            type?: MisbehaviorType | undefined;
            validator?: {
                address?: Uint8Array | undefined;
                power?: bigint | undefined;
            } | undefined;
            height?: bigint | undefined;
            time?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            totalVotingPower?: bigint | undefined;
        }[] | undefined;
        hash?: Uint8Array | undefined;
        height?: bigint | undefined;
        time?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
        nextValidatorsHash?: Uint8Array | undefined;
        proposerAddress?: Uint8Array | undefined;
    } & {
        txs?: (Uint8Array[] & Uint8Array[] & Record<Exclude<keyof I["txs"], keyof Uint8Array[]>, never>) | undefined;
        proposedLastCommit?: ({
            round?: number | undefined;
            votes?: {
                validator?: {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } | undefined;
                signedLastBlock?: boolean | undefined;
            }[] | undefined;
        } & {
            round?: number | undefined;
            votes?: ({
                validator?: {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } | undefined;
                signedLastBlock?: boolean | undefined;
            }[] & ({
                validator?: {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } | undefined;
                signedLastBlock?: boolean | undefined;
            } & {
                validator?: ({
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } & {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } & Record<Exclude<keyof I["proposedLastCommit"]["votes"][number]["validator"], keyof Validator>, never>) | undefined;
                signedLastBlock?: boolean | undefined;
            } & Record<Exclude<keyof I["proposedLastCommit"]["votes"][number], keyof VoteInfo>, never>)[] & Record<Exclude<keyof I["proposedLastCommit"]["votes"], keyof {
                validator?: {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } | undefined;
                signedLastBlock?: boolean | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["proposedLastCommit"], keyof CommitInfo>, never>) | undefined;
        misbehavior?: ({
            type?: MisbehaviorType | undefined;
            validator?: {
                address?: Uint8Array | undefined;
                power?: bigint | undefined;
            } | undefined;
            height?: bigint | undefined;
            time?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            totalVotingPower?: bigint | undefined;
        }[] & ({
            type?: MisbehaviorType | undefined;
            validator?: {
                address?: Uint8Array | undefined;
                power?: bigint | undefined;
            } | undefined;
            height?: bigint | undefined;
            time?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            totalVotingPower?: bigint | undefined;
        } & {
            type?: MisbehaviorType | undefined;
            validator?: ({
                address?: Uint8Array | undefined;
                power?: bigint | undefined;
            } & {
                address?: Uint8Array | undefined;
                power?: bigint | undefined;
            } & Record<Exclude<keyof I["misbehavior"][number]["validator"], keyof Validator>, never>) | undefined;
            height?: bigint | undefined;
            time?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["misbehavior"][number]["time"], keyof Timestamp>, never>) | undefined;
            totalVotingPower?: bigint | undefined;
        } & Record<Exclude<keyof I["misbehavior"][number], keyof Misbehavior>, never>)[] & Record<Exclude<keyof I["misbehavior"], keyof {
            type?: MisbehaviorType | undefined;
            validator?: {
                address?: Uint8Array | undefined;
                power?: bigint | undefined;
            } | undefined;
            height?: bigint | undefined;
            time?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            totalVotingPower?: bigint | undefined;
        }[]>, never>) | undefined;
        hash?: Uint8Array | undefined;
        height?: bigint | undefined;
        time?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["time"], keyof Timestamp>, never>) | undefined;
        nextValidatorsHash?: Uint8Array | undefined;
        proposerAddress?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof RequestProcessProposal>, never>>(object: I): RequestProcessProposal;
};
export declare const Response: {
    typeUrl: string;
    encode(message: Response, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Response;
    fromJSON(object: any): Response;
    toJSON(message: Response): unknown;
    fromPartial<I extends {
        exception?: {
            error?: string | undefined;
        } | undefined;
        echo?: {
            message?: string | undefined;
        } | undefined;
        flush?: {} | undefined;
        info?: {
            data?: string | undefined;
            version?: string | undefined;
            appVersion?: bigint | undefined;
            lastBlockHeight?: bigint | undefined;
            lastBlockAppHash?: Uint8Array | undefined;
        } | undefined;
        initChain?: {
            consensusParams?: {
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
            } | undefined;
            validators?: {
                pubKey?: {
                    ed25519?: Uint8Array | undefined;
                    secp256k1?: Uint8Array | undefined;
                } | undefined;
                power?: bigint | undefined;
            }[] | undefined;
            appHash?: Uint8Array | undefined;
        } | undefined;
        query?: {
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
        } | undefined;
        beginBlock?: {
            events?: {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[] | undefined;
        } | undefined;
        checkTx?: {
            code?: number | undefined;
            data?: Uint8Array | undefined;
            log?: string | undefined;
            info?: string | undefined;
            gasWanted?: bigint | undefined;
            gasUsed?: bigint | undefined;
            events?: {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[] | undefined;
            codespace?: string | undefined;
            sender?: string | undefined;
            priority?: bigint | undefined;
            mempoolError?: string | undefined;
        } | undefined;
        deliverTx?: {
            code?: number | undefined;
            data?: Uint8Array | undefined;
            log?: string | undefined;
            info?: string | undefined;
            gasWanted?: bigint | undefined;
            gasUsed?: bigint | undefined;
            events?: {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[] | undefined;
            codespace?: string | undefined;
        } | undefined;
        endBlock?: {
            validatorUpdates?: {
                pubKey?: {
                    ed25519?: Uint8Array | undefined;
                    secp256k1?: Uint8Array | undefined;
                } | undefined;
                power?: bigint | undefined;
            }[] | undefined;
            consensusParamUpdates?: {
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
            } | undefined;
            events?: {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[] | undefined;
        } | undefined;
        commit?: {
            data?: Uint8Array | undefined;
            retainHeight?: bigint | undefined;
        } | undefined;
        listSnapshots?: {
            snapshots?: {
                height?: bigint | undefined;
                format?: number | undefined;
                chunks?: number | undefined;
                hash?: Uint8Array | undefined;
                metadata?: Uint8Array | undefined;
            }[] | undefined;
        } | undefined;
        offerSnapshot?: {
            result?: ResponseOfferSnapshot_Result | undefined;
        } | undefined;
        loadSnapshotChunk?: {
            chunk?: Uint8Array | undefined;
        } | undefined;
        applySnapshotChunk?: {
            result?: ResponseApplySnapshotChunk_Result | undefined;
            refetchChunks?: number[] | undefined;
            rejectSenders?: string[] | undefined;
        } | undefined;
        prepareProposal?: {
            txs?: Uint8Array[] | undefined;
        } | undefined;
        processProposal?: {
            status?: ResponseProcessProposal_ProposalStatus | undefined;
        } | undefined;
    } & {
        exception?: ({
            error?: string | undefined;
        } & {
            error?: string | undefined;
        } & Record<Exclude<keyof I["exception"], "error">, never>) | undefined;
        echo?: ({
            message?: string | undefined;
        } & {
            message?: string | undefined;
        } & Record<Exclude<keyof I["echo"], "message">, never>) | undefined;
        flush?: ({} & {} & Record<Exclude<keyof I["flush"], never>, never>) | undefined;
        info?: ({
            data?: string | undefined;
            version?: string | undefined;
            appVersion?: bigint | undefined;
            lastBlockHeight?: bigint | undefined;
            lastBlockAppHash?: Uint8Array | undefined;
        } & {
            data?: string | undefined;
            version?: string | undefined;
            appVersion?: bigint | undefined;
            lastBlockHeight?: bigint | undefined;
            lastBlockAppHash?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["info"], keyof ResponseInfo>, never>) | undefined;
        initChain?: ({
            consensusParams?: {
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
            } | undefined;
            validators?: {
                pubKey?: {
                    ed25519?: Uint8Array | undefined;
                    secp256k1?: Uint8Array | undefined;
                } | undefined;
                power?: bigint | undefined;
            }[] | undefined;
            appHash?: Uint8Array | undefined;
        } & {
            consensusParams?: ({
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
                } & Record<Exclude<keyof I["initChain"]["consensusParams"]["block"], keyof import("../types/params").BlockParams>, never>) | undefined;
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
                    } & Record<Exclude<keyof I["initChain"]["consensusParams"]["evidence"]["maxAgeDuration"], keyof import("../../google/protobuf/duration").Duration>, never>) | undefined;
                    maxBytes?: bigint | undefined;
                } & Record<Exclude<keyof I["initChain"]["consensusParams"]["evidence"], keyof import("../types/params").EvidenceParams>, never>) | undefined;
                validator?: ({
                    pubKeyTypes?: string[] | undefined;
                } & {
                    pubKeyTypes?: (string[] & string[] & Record<Exclude<keyof I["initChain"]["consensusParams"]["validator"]["pubKeyTypes"], keyof string[]>, never>) | undefined;
                } & Record<Exclude<keyof I["initChain"]["consensusParams"]["validator"], "pubKeyTypes">, never>) | undefined;
                version?: ({
                    app?: bigint | undefined;
                } & {
                    app?: bigint | undefined;
                } & Record<Exclude<keyof I["initChain"]["consensusParams"]["version"], "app">, never>) | undefined;
            } & Record<Exclude<keyof I["initChain"]["consensusParams"], keyof ConsensusParams>, never>) | undefined;
            validators?: ({
                pubKey?: {
                    ed25519?: Uint8Array | undefined;
                    secp256k1?: Uint8Array | undefined;
                } | undefined;
                power?: bigint | undefined;
            }[] & ({
                pubKey?: {
                    ed25519?: Uint8Array | undefined;
                    secp256k1?: Uint8Array | undefined;
                } | undefined;
                power?: bigint | undefined;
            } & {
                pubKey?: ({
                    ed25519?: Uint8Array | undefined;
                    secp256k1?: Uint8Array | undefined;
                } & {
                    ed25519?: Uint8Array | undefined;
                    secp256k1?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["initChain"]["validators"][number]["pubKey"], keyof PublicKey>, never>) | undefined;
                power?: bigint | undefined;
            } & Record<Exclude<keyof I["initChain"]["validators"][number], keyof ValidatorUpdate>, never>)[] & Record<Exclude<keyof I["initChain"]["validators"], keyof {
                pubKey?: {
                    ed25519?: Uint8Array | undefined;
                    secp256k1?: Uint8Array | undefined;
                } | undefined;
                power?: bigint | undefined;
            }[]>, never>) | undefined;
            appHash?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["initChain"], keyof ResponseInitChain>, never>) | undefined;
        query?: ({
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
                } & Record<Exclude<keyof I["query"]["proofOps"]["ops"][number], keyof import("../crypto/proof").ProofOp>, never>)[] & Record<Exclude<keyof I["query"]["proofOps"]["ops"], keyof {
                    type?: string | undefined;
                    key?: Uint8Array | undefined;
                    data?: Uint8Array | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["query"]["proofOps"], "ops">, never>) | undefined;
            height?: bigint | undefined;
            codespace?: string | undefined;
        } & Record<Exclude<keyof I["query"], keyof ResponseQuery>, never>) | undefined;
        beginBlock?: ({
            events?: {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[] | undefined;
        } & {
            events?: ({
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[] & ({
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            } & {
                type?: string | undefined;
                attributes?: ({
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] & ({
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                } & {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                } & Record<Exclude<keyof I["beginBlock"]["events"][number]["attributes"][number], keyof EventAttribute>, never>)[] & Record<Exclude<keyof I["beginBlock"]["events"][number]["attributes"], keyof {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["beginBlock"]["events"][number], keyof Event>, never>)[] & Record<Exclude<keyof I["beginBlock"]["events"], keyof {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["beginBlock"], "events">, never>) | undefined;
        checkTx?: ({
            code?: number | undefined;
            data?: Uint8Array | undefined;
            log?: string | undefined;
            info?: string | undefined;
            gasWanted?: bigint | undefined;
            gasUsed?: bigint | undefined;
            events?: {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[] | undefined;
            codespace?: string | undefined;
            sender?: string | undefined;
            priority?: bigint | undefined;
            mempoolError?: string | undefined;
        } & {
            code?: number | undefined;
            data?: Uint8Array | undefined;
            log?: string | undefined;
            info?: string | undefined;
            gasWanted?: bigint | undefined;
            gasUsed?: bigint | undefined;
            events?: ({
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[] & ({
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            } & {
                type?: string | undefined;
                attributes?: ({
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] & ({
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                } & {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                } & Record<Exclude<keyof I["checkTx"]["events"][number]["attributes"][number], keyof EventAttribute>, never>)[] & Record<Exclude<keyof I["checkTx"]["events"][number]["attributes"], keyof {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["checkTx"]["events"][number], keyof Event>, never>)[] & Record<Exclude<keyof I["checkTx"]["events"], keyof {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[]>, never>) | undefined;
            codespace?: string | undefined;
            sender?: string | undefined;
            priority?: bigint | undefined;
            mempoolError?: string | undefined;
        } & Record<Exclude<keyof I["checkTx"], keyof ResponseCheckTx>, never>) | undefined;
        deliverTx?: ({
            code?: number | undefined;
            data?: Uint8Array | undefined;
            log?: string | undefined;
            info?: string | undefined;
            gasWanted?: bigint | undefined;
            gasUsed?: bigint | undefined;
            events?: {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[] | undefined;
            codespace?: string | undefined;
        } & {
            code?: number | undefined;
            data?: Uint8Array | undefined;
            log?: string | undefined;
            info?: string | undefined;
            gasWanted?: bigint | undefined;
            gasUsed?: bigint | undefined;
            events?: ({
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[] & ({
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            } & {
                type?: string | undefined;
                attributes?: ({
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] & ({
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                } & {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                } & Record<Exclude<keyof I["deliverTx"]["events"][number]["attributes"][number], keyof EventAttribute>, never>)[] & Record<Exclude<keyof I["deliverTx"]["events"][number]["attributes"], keyof {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["deliverTx"]["events"][number], keyof Event>, never>)[] & Record<Exclude<keyof I["deliverTx"]["events"], keyof {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[]>, never>) | undefined;
            codespace?: string | undefined;
        } & Record<Exclude<keyof I["deliverTx"], keyof ResponseDeliverTx>, never>) | undefined;
        endBlock?: ({
            validatorUpdates?: {
                pubKey?: {
                    ed25519?: Uint8Array | undefined;
                    secp256k1?: Uint8Array | undefined;
                } | undefined;
                power?: bigint | undefined;
            }[] | undefined;
            consensusParamUpdates?: {
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
            } | undefined;
            events?: {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[] | undefined;
        } & {
            validatorUpdates?: ({
                pubKey?: {
                    ed25519?: Uint8Array | undefined;
                    secp256k1?: Uint8Array | undefined;
                } | undefined;
                power?: bigint | undefined;
            }[] & ({
                pubKey?: {
                    ed25519?: Uint8Array | undefined;
                    secp256k1?: Uint8Array | undefined;
                } | undefined;
                power?: bigint | undefined;
            } & {
                pubKey?: ({
                    ed25519?: Uint8Array | undefined;
                    secp256k1?: Uint8Array | undefined;
                } & {
                    ed25519?: Uint8Array | undefined;
                    secp256k1?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["endBlock"]["validatorUpdates"][number]["pubKey"], keyof PublicKey>, never>) | undefined;
                power?: bigint | undefined;
            } & Record<Exclude<keyof I["endBlock"]["validatorUpdates"][number], keyof ValidatorUpdate>, never>)[] & Record<Exclude<keyof I["endBlock"]["validatorUpdates"], keyof {
                pubKey?: {
                    ed25519?: Uint8Array | undefined;
                    secp256k1?: Uint8Array | undefined;
                } | undefined;
                power?: bigint | undefined;
            }[]>, never>) | undefined;
            consensusParamUpdates?: ({
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
                } & Record<Exclude<keyof I["endBlock"]["consensusParamUpdates"]["block"], keyof import("../types/params").BlockParams>, never>) | undefined;
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
                    } & Record<Exclude<keyof I["endBlock"]["consensusParamUpdates"]["evidence"]["maxAgeDuration"], keyof import("../../google/protobuf/duration").Duration>, never>) | undefined;
                    maxBytes?: bigint | undefined;
                } & Record<Exclude<keyof I["endBlock"]["consensusParamUpdates"]["evidence"], keyof import("../types/params").EvidenceParams>, never>) | undefined;
                validator?: ({
                    pubKeyTypes?: string[] | undefined;
                } & {
                    pubKeyTypes?: (string[] & string[] & Record<Exclude<keyof I["endBlock"]["consensusParamUpdates"]["validator"]["pubKeyTypes"], keyof string[]>, never>) | undefined;
                } & Record<Exclude<keyof I["endBlock"]["consensusParamUpdates"]["validator"], "pubKeyTypes">, never>) | undefined;
                version?: ({
                    app?: bigint | undefined;
                } & {
                    app?: bigint | undefined;
                } & Record<Exclude<keyof I["endBlock"]["consensusParamUpdates"]["version"], "app">, never>) | undefined;
            } & Record<Exclude<keyof I["endBlock"]["consensusParamUpdates"], keyof ConsensusParams>, never>) | undefined;
            events?: ({
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[] & ({
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            } & {
                type?: string | undefined;
                attributes?: ({
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] & ({
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                } & {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                } & Record<Exclude<keyof I["endBlock"]["events"][number]["attributes"][number], keyof EventAttribute>, never>)[] & Record<Exclude<keyof I["endBlock"]["events"][number]["attributes"], keyof {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["endBlock"]["events"][number], keyof Event>, never>)[] & Record<Exclude<keyof I["endBlock"]["events"], keyof {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["endBlock"], keyof ResponseEndBlock>, never>) | undefined;
        commit?: ({
            data?: Uint8Array | undefined;
            retainHeight?: bigint | undefined;
        } & {
            data?: Uint8Array | undefined;
            retainHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["commit"], keyof ResponseCommit>, never>) | undefined;
        listSnapshots?: ({
            snapshots?: {
                height?: bigint | undefined;
                format?: number | undefined;
                chunks?: number | undefined;
                hash?: Uint8Array | undefined;
                metadata?: Uint8Array | undefined;
            }[] | undefined;
        } & {
            snapshots?: ({
                height?: bigint | undefined;
                format?: number | undefined;
                chunks?: number | undefined;
                hash?: Uint8Array | undefined;
                metadata?: Uint8Array | undefined;
            }[] & ({
                height?: bigint | undefined;
                format?: number | undefined;
                chunks?: number | undefined;
                hash?: Uint8Array | undefined;
                metadata?: Uint8Array | undefined;
            } & {
                height?: bigint | undefined;
                format?: number | undefined;
                chunks?: number | undefined;
                hash?: Uint8Array | undefined;
                metadata?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["listSnapshots"]["snapshots"][number], keyof Snapshot>, never>)[] & Record<Exclude<keyof I["listSnapshots"]["snapshots"], keyof {
                height?: bigint | undefined;
                format?: number | undefined;
                chunks?: number | undefined;
                hash?: Uint8Array | undefined;
                metadata?: Uint8Array | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["listSnapshots"], "snapshots">, never>) | undefined;
        offerSnapshot?: ({
            result?: ResponseOfferSnapshot_Result | undefined;
        } & {
            result?: ResponseOfferSnapshot_Result | undefined;
        } & Record<Exclude<keyof I["offerSnapshot"], "result">, never>) | undefined;
        loadSnapshotChunk?: ({
            chunk?: Uint8Array | undefined;
        } & {
            chunk?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["loadSnapshotChunk"], "chunk">, never>) | undefined;
        applySnapshotChunk?: ({
            result?: ResponseApplySnapshotChunk_Result | undefined;
            refetchChunks?: number[] | undefined;
            rejectSenders?: string[] | undefined;
        } & {
            result?: ResponseApplySnapshotChunk_Result | undefined;
            refetchChunks?: (number[] & number[] & Record<Exclude<keyof I["applySnapshotChunk"]["refetchChunks"], keyof number[]>, never>) | undefined;
            rejectSenders?: (string[] & string[] & Record<Exclude<keyof I["applySnapshotChunk"]["rejectSenders"], keyof string[]>, never>) | undefined;
        } & Record<Exclude<keyof I["applySnapshotChunk"], keyof ResponseApplySnapshotChunk>, never>) | undefined;
        prepareProposal?: ({
            txs?: Uint8Array[] | undefined;
        } & {
            txs?: (Uint8Array[] & Uint8Array[] & Record<Exclude<keyof I["prepareProposal"]["txs"], keyof Uint8Array[]>, never>) | undefined;
        } & Record<Exclude<keyof I["prepareProposal"], "txs">, never>) | undefined;
        processProposal?: ({
            status?: ResponseProcessProposal_ProposalStatus | undefined;
        } & {
            status?: ResponseProcessProposal_ProposalStatus | undefined;
        } & Record<Exclude<keyof I["processProposal"], "status">, never>) | undefined;
    } & Record<Exclude<keyof I, keyof Response>, never>>(object: I): Response;
};
export declare const ResponseException: {
    typeUrl: string;
    encode(message: ResponseException, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ResponseException;
    fromJSON(object: any): ResponseException;
    toJSON(message: ResponseException): unknown;
    fromPartial<I extends {
        error?: string | undefined;
    } & {
        error?: string | undefined;
    } & Record<Exclude<keyof I, "error">, never>>(object: I): ResponseException;
};
export declare const ResponseEcho: {
    typeUrl: string;
    encode(message: ResponseEcho, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ResponseEcho;
    fromJSON(object: any): ResponseEcho;
    toJSON(message: ResponseEcho): unknown;
    fromPartial<I extends {
        message?: string | undefined;
    } & {
        message?: string | undefined;
    } & Record<Exclude<keyof I, "message">, never>>(object: I): ResponseEcho;
};
export declare const ResponseFlush: {
    typeUrl: string;
    encode(_: ResponseFlush, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ResponseFlush;
    fromJSON(_: any): ResponseFlush;
    toJSON(_: ResponseFlush): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): ResponseFlush;
};
export declare const ResponseInfo: {
    typeUrl: string;
    encode(message: ResponseInfo, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ResponseInfo;
    fromJSON(object: any): ResponseInfo;
    toJSON(message: ResponseInfo): unknown;
    fromPartial<I extends {
        data?: string | undefined;
        version?: string | undefined;
        appVersion?: bigint | undefined;
        lastBlockHeight?: bigint | undefined;
        lastBlockAppHash?: Uint8Array | undefined;
    } & {
        data?: string | undefined;
        version?: string | undefined;
        appVersion?: bigint | undefined;
        lastBlockHeight?: bigint | undefined;
        lastBlockAppHash?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof ResponseInfo>, never>>(object: I): ResponseInfo;
};
export declare const ResponseInitChain: {
    typeUrl: string;
    encode(message: ResponseInitChain, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ResponseInitChain;
    fromJSON(object: any): ResponseInitChain;
    toJSON(message: ResponseInitChain): unknown;
    fromPartial<I extends {
        consensusParams?: {
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
        } | undefined;
        validators?: {
            pubKey?: {
                ed25519?: Uint8Array | undefined;
                secp256k1?: Uint8Array | undefined;
            } | undefined;
            power?: bigint | undefined;
        }[] | undefined;
        appHash?: Uint8Array | undefined;
    } & {
        consensusParams?: ({
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
            } & Record<Exclude<keyof I["consensusParams"]["block"], keyof import("../types/params").BlockParams>, never>) | undefined;
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
                } & Record<Exclude<keyof I["consensusParams"]["evidence"]["maxAgeDuration"], keyof import("../../google/protobuf/duration").Duration>, never>) | undefined;
                maxBytes?: bigint | undefined;
            } & Record<Exclude<keyof I["consensusParams"]["evidence"], keyof import("../types/params").EvidenceParams>, never>) | undefined;
            validator?: ({
                pubKeyTypes?: string[] | undefined;
            } & {
                pubKeyTypes?: (string[] & string[] & Record<Exclude<keyof I["consensusParams"]["validator"]["pubKeyTypes"], keyof string[]>, never>) | undefined;
            } & Record<Exclude<keyof I["consensusParams"]["validator"], "pubKeyTypes">, never>) | undefined;
            version?: ({
                app?: bigint | undefined;
            } & {
                app?: bigint | undefined;
            } & Record<Exclude<keyof I["consensusParams"]["version"], "app">, never>) | undefined;
        } & Record<Exclude<keyof I["consensusParams"], keyof ConsensusParams>, never>) | undefined;
        validators?: ({
            pubKey?: {
                ed25519?: Uint8Array | undefined;
                secp256k1?: Uint8Array | undefined;
            } | undefined;
            power?: bigint | undefined;
        }[] & ({
            pubKey?: {
                ed25519?: Uint8Array | undefined;
                secp256k1?: Uint8Array | undefined;
            } | undefined;
            power?: bigint | undefined;
        } & {
            pubKey?: ({
                ed25519?: Uint8Array | undefined;
                secp256k1?: Uint8Array | undefined;
            } & {
                ed25519?: Uint8Array | undefined;
                secp256k1?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["validators"][number]["pubKey"], keyof PublicKey>, never>) | undefined;
            power?: bigint | undefined;
        } & Record<Exclude<keyof I["validators"][number], keyof ValidatorUpdate>, never>)[] & Record<Exclude<keyof I["validators"], keyof {
            pubKey?: {
                ed25519?: Uint8Array | undefined;
                secp256k1?: Uint8Array | undefined;
            } | undefined;
            power?: bigint | undefined;
        }[]>, never>) | undefined;
        appHash?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof ResponseInitChain>, never>>(object: I): ResponseInitChain;
};
export declare const ResponseQuery: {
    typeUrl: string;
    encode(message: ResponseQuery, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ResponseQuery;
    fromJSON(object: any): ResponseQuery;
    toJSON(message: ResponseQuery): unknown;
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
            } & Record<Exclude<keyof I["proofOps"]["ops"][number], keyof import("../crypto/proof").ProofOp>, never>)[] & Record<Exclude<keyof I["proofOps"]["ops"], keyof {
                type?: string | undefined;
                key?: Uint8Array | undefined;
                data?: Uint8Array | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["proofOps"], "ops">, never>) | undefined;
        height?: bigint | undefined;
        codespace?: string | undefined;
    } & Record<Exclude<keyof I, keyof ResponseQuery>, never>>(object: I): ResponseQuery;
};
export declare const ResponseBeginBlock: {
    typeUrl: string;
    encode(message: ResponseBeginBlock, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ResponseBeginBlock;
    fromJSON(object: any): ResponseBeginBlock;
    toJSON(message: ResponseBeginBlock): unknown;
    fromPartial<I extends {
        events?: {
            type?: string | undefined;
            attributes?: {
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            }[] | undefined;
        }[] | undefined;
    } & {
        events?: ({
            type?: string | undefined;
            attributes?: {
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            }[] | undefined;
        }[] & ({
            type?: string | undefined;
            attributes?: {
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            }[] | undefined;
        } & {
            type?: string | undefined;
            attributes?: ({
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            }[] & ({
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            } & {
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            } & Record<Exclude<keyof I["events"][number]["attributes"][number], keyof EventAttribute>, never>)[] & Record<Exclude<keyof I["events"][number]["attributes"], keyof {
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["events"][number], keyof Event>, never>)[] & Record<Exclude<keyof I["events"], keyof {
            type?: string | undefined;
            attributes?: {
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            }[] | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "events">, never>>(object: I): ResponseBeginBlock;
};
export declare const ResponseCheckTx: {
    typeUrl: string;
    encode(message: ResponseCheckTx, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ResponseCheckTx;
    fromJSON(object: any): ResponseCheckTx;
    toJSON(message: ResponseCheckTx): unknown;
    fromPartial<I extends {
        code?: number | undefined;
        data?: Uint8Array | undefined;
        log?: string | undefined;
        info?: string | undefined;
        gasWanted?: bigint | undefined;
        gasUsed?: bigint | undefined;
        events?: {
            type?: string | undefined;
            attributes?: {
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            }[] | undefined;
        }[] | undefined;
        codespace?: string | undefined;
        sender?: string | undefined;
        priority?: bigint | undefined;
        mempoolError?: string | undefined;
    } & {
        code?: number | undefined;
        data?: Uint8Array | undefined;
        log?: string | undefined;
        info?: string | undefined;
        gasWanted?: bigint | undefined;
        gasUsed?: bigint | undefined;
        events?: ({
            type?: string | undefined;
            attributes?: {
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            }[] | undefined;
        }[] & ({
            type?: string | undefined;
            attributes?: {
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            }[] | undefined;
        } & {
            type?: string | undefined;
            attributes?: ({
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            }[] & ({
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            } & {
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            } & Record<Exclude<keyof I["events"][number]["attributes"][number], keyof EventAttribute>, never>)[] & Record<Exclude<keyof I["events"][number]["attributes"], keyof {
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["events"][number], keyof Event>, never>)[] & Record<Exclude<keyof I["events"], keyof {
            type?: string | undefined;
            attributes?: {
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            }[] | undefined;
        }[]>, never>) | undefined;
        codespace?: string | undefined;
        sender?: string | undefined;
        priority?: bigint | undefined;
        mempoolError?: string | undefined;
    } & Record<Exclude<keyof I, keyof ResponseCheckTx>, never>>(object: I): ResponseCheckTx;
};
export declare const ResponseDeliverTx: {
    typeUrl: string;
    encode(message: ResponseDeliverTx, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ResponseDeliverTx;
    fromJSON(object: any): ResponseDeliverTx;
    toJSON(message: ResponseDeliverTx): unknown;
    fromPartial<I extends {
        code?: number | undefined;
        data?: Uint8Array | undefined;
        log?: string | undefined;
        info?: string | undefined;
        gasWanted?: bigint | undefined;
        gasUsed?: bigint | undefined;
        events?: {
            type?: string | undefined;
            attributes?: {
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            }[] | undefined;
        }[] | undefined;
        codespace?: string | undefined;
    } & {
        code?: number | undefined;
        data?: Uint8Array | undefined;
        log?: string | undefined;
        info?: string | undefined;
        gasWanted?: bigint | undefined;
        gasUsed?: bigint | undefined;
        events?: ({
            type?: string | undefined;
            attributes?: {
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            }[] | undefined;
        }[] & ({
            type?: string | undefined;
            attributes?: {
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            }[] | undefined;
        } & {
            type?: string | undefined;
            attributes?: ({
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            }[] & ({
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            } & {
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            } & Record<Exclude<keyof I["events"][number]["attributes"][number], keyof EventAttribute>, never>)[] & Record<Exclude<keyof I["events"][number]["attributes"], keyof {
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["events"][number], keyof Event>, never>)[] & Record<Exclude<keyof I["events"], keyof {
            type?: string | undefined;
            attributes?: {
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            }[] | undefined;
        }[]>, never>) | undefined;
        codespace?: string | undefined;
    } & Record<Exclude<keyof I, keyof ResponseDeliverTx>, never>>(object: I): ResponseDeliverTx;
};
export declare const ResponseEndBlock: {
    typeUrl: string;
    encode(message: ResponseEndBlock, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ResponseEndBlock;
    fromJSON(object: any): ResponseEndBlock;
    toJSON(message: ResponseEndBlock): unknown;
    fromPartial<I extends {
        validatorUpdates?: {
            pubKey?: {
                ed25519?: Uint8Array | undefined;
                secp256k1?: Uint8Array | undefined;
            } | undefined;
            power?: bigint | undefined;
        }[] | undefined;
        consensusParamUpdates?: {
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
        } | undefined;
        events?: {
            type?: string | undefined;
            attributes?: {
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            }[] | undefined;
        }[] | undefined;
    } & {
        validatorUpdates?: ({
            pubKey?: {
                ed25519?: Uint8Array | undefined;
                secp256k1?: Uint8Array | undefined;
            } | undefined;
            power?: bigint | undefined;
        }[] & ({
            pubKey?: {
                ed25519?: Uint8Array | undefined;
                secp256k1?: Uint8Array | undefined;
            } | undefined;
            power?: bigint | undefined;
        } & {
            pubKey?: ({
                ed25519?: Uint8Array | undefined;
                secp256k1?: Uint8Array | undefined;
            } & {
                ed25519?: Uint8Array | undefined;
                secp256k1?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["validatorUpdates"][number]["pubKey"], keyof PublicKey>, never>) | undefined;
            power?: bigint | undefined;
        } & Record<Exclude<keyof I["validatorUpdates"][number], keyof ValidatorUpdate>, never>)[] & Record<Exclude<keyof I["validatorUpdates"], keyof {
            pubKey?: {
                ed25519?: Uint8Array | undefined;
                secp256k1?: Uint8Array | undefined;
            } | undefined;
            power?: bigint | undefined;
        }[]>, never>) | undefined;
        consensusParamUpdates?: ({
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
            } & Record<Exclude<keyof I["consensusParamUpdates"]["block"], keyof import("../types/params").BlockParams>, never>) | undefined;
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
                } & Record<Exclude<keyof I["consensusParamUpdates"]["evidence"]["maxAgeDuration"], keyof import("../../google/protobuf/duration").Duration>, never>) | undefined;
                maxBytes?: bigint | undefined;
            } & Record<Exclude<keyof I["consensusParamUpdates"]["evidence"], keyof import("../types/params").EvidenceParams>, never>) | undefined;
            validator?: ({
                pubKeyTypes?: string[] | undefined;
            } & {
                pubKeyTypes?: (string[] & string[] & Record<Exclude<keyof I["consensusParamUpdates"]["validator"]["pubKeyTypes"], keyof string[]>, never>) | undefined;
            } & Record<Exclude<keyof I["consensusParamUpdates"]["validator"], "pubKeyTypes">, never>) | undefined;
            version?: ({
                app?: bigint | undefined;
            } & {
                app?: bigint | undefined;
            } & Record<Exclude<keyof I["consensusParamUpdates"]["version"], "app">, never>) | undefined;
        } & Record<Exclude<keyof I["consensusParamUpdates"], keyof ConsensusParams>, never>) | undefined;
        events?: ({
            type?: string | undefined;
            attributes?: {
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            }[] | undefined;
        }[] & ({
            type?: string | undefined;
            attributes?: {
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            }[] | undefined;
        } & {
            type?: string | undefined;
            attributes?: ({
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            }[] & ({
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            } & {
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            } & Record<Exclude<keyof I["events"][number]["attributes"][number], keyof EventAttribute>, never>)[] & Record<Exclude<keyof I["events"][number]["attributes"], keyof {
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["events"][number], keyof Event>, never>)[] & Record<Exclude<keyof I["events"], keyof {
            type?: string | undefined;
            attributes?: {
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            }[] | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof ResponseEndBlock>, never>>(object: I): ResponseEndBlock;
};
export declare const ResponseCommit: {
    typeUrl: string;
    encode(message: ResponseCommit, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ResponseCommit;
    fromJSON(object: any): ResponseCommit;
    toJSON(message: ResponseCommit): unknown;
    fromPartial<I extends {
        data?: Uint8Array | undefined;
        retainHeight?: bigint | undefined;
    } & {
        data?: Uint8Array | undefined;
        retainHeight?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof ResponseCommit>, never>>(object: I): ResponseCommit;
};
export declare const ResponseListSnapshots: {
    typeUrl: string;
    encode(message: ResponseListSnapshots, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ResponseListSnapshots;
    fromJSON(object: any): ResponseListSnapshots;
    toJSON(message: ResponseListSnapshots): unknown;
    fromPartial<I extends {
        snapshots?: {
            height?: bigint | undefined;
            format?: number | undefined;
            chunks?: number | undefined;
            hash?: Uint8Array | undefined;
            metadata?: Uint8Array | undefined;
        }[] | undefined;
    } & {
        snapshots?: ({
            height?: bigint | undefined;
            format?: number | undefined;
            chunks?: number | undefined;
            hash?: Uint8Array | undefined;
            metadata?: Uint8Array | undefined;
        }[] & ({
            height?: bigint | undefined;
            format?: number | undefined;
            chunks?: number | undefined;
            hash?: Uint8Array | undefined;
            metadata?: Uint8Array | undefined;
        } & {
            height?: bigint | undefined;
            format?: number | undefined;
            chunks?: number | undefined;
            hash?: Uint8Array | undefined;
            metadata?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["snapshots"][number], keyof Snapshot>, never>)[] & Record<Exclude<keyof I["snapshots"], keyof {
            height?: bigint | undefined;
            format?: number | undefined;
            chunks?: number | undefined;
            hash?: Uint8Array | undefined;
            metadata?: Uint8Array | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "snapshots">, never>>(object: I): ResponseListSnapshots;
};
export declare const ResponseOfferSnapshot: {
    typeUrl: string;
    encode(message: ResponseOfferSnapshot, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ResponseOfferSnapshot;
    fromJSON(object: any): ResponseOfferSnapshot;
    toJSON(message: ResponseOfferSnapshot): unknown;
    fromPartial<I extends {
        result?: ResponseOfferSnapshot_Result | undefined;
    } & {
        result?: ResponseOfferSnapshot_Result | undefined;
    } & Record<Exclude<keyof I, "result">, never>>(object: I): ResponseOfferSnapshot;
};
export declare const ResponseLoadSnapshotChunk: {
    typeUrl: string;
    encode(message: ResponseLoadSnapshotChunk, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ResponseLoadSnapshotChunk;
    fromJSON(object: any): ResponseLoadSnapshotChunk;
    toJSON(message: ResponseLoadSnapshotChunk): unknown;
    fromPartial<I extends {
        chunk?: Uint8Array | undefined;
    } & {
        chunk?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, "chunk">, never>>(object: I): ResponseLoadSnapshotChunk;
};
export declare const ResponseApplySnapshotChunk: {
    typeUrl: string;
    encode(message: ResponseApplySnapshotChunk, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ResponseApplySnapshotChunk;
    fromJSON(object: any): ResponseApplySnapshotChunk;
    toJSON(message: ResponseApplySnapshotChunk): unknown;
    fromPartial<I extends {
        result?: ResponseApplySnapshotChunk_Result | undefined;
        refetchChunks?: number[] | undefined;
        rejectSenders?: string[] | undefined;
    } & {
        result?: ResponseApplySnapshotChunk_Result | undefined;
        refetchChunks?: (number[] & number[] & Record<Exclude<keyof I["refetchChunks"], keyof number[]>, never>) | undefined;
        rejectSenders?: (string[] & string[] & Record<Exclude<keyof I["rejectSenders"], keyof string[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof ResponseApplySnapshotChunk>, never>>(object: I): ResponseApplySnapshotChunk;
};
export declare const ResponsePrepareProposal: {
    typeUrl: string;
    encode(message: ResponsePrepareProposal, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ResponsePrepareProposal;
    fromJSON(object: any): ResponsePrepareProposal;
    toJSON(message: ResponsePrepareProposal): unknown;
    fromPartial<I extends {
        txs?: Uint8Array[] | undefined;
    } & {
        txs?: (Uint8Array[] & Uint8Array[] & Record<Exclude<keyof I["txs"], keyof Uint8Array[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "txs">, never>>(object: I): ResponsePrepareProposal;
};
export declare const ResponseProcessProposal: {
    typeUrl: string;
    encode(message: ResponseProcessProposal, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ResponseProcessProposal;
    fromJSON(object: any): ResponseProcessProposal;
    toJSON(message: ResponseProcessProposal): unknown;
    fromPartial<I extends {
        status?: ResponseProcessProposal_ProposalStatus | undefined;
    } & {
        status?: ResponseProcessProposal_ProposalStatus | undefined;
    } & Record<Exclude<keyof I, "status">, never>>(object: I): ResponseProcessProposal;
};
export declare const CommitInfo: {
    typeUrl: string;
    encode(message: CommitInfo, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): CommitInfo;
    fromJSON(object: any): CommitInfo;
    toJSON(message: CommitInfo): unknown;
    fromPartial<I extends {
        round?: number | undefined;
        votes?: {
            validator?: {
                address?: Uint8Array | undefined;
                power?: bigint | undefined;
            } | undefined;
            signedLastBlock?: boolean | undefined;
        }[] | undefined;
    } & {
        round?: number | undefined;
        votes?: ({
            validator?: {
                address?: Uint8Array | undefined;
                power?: bigint | undefined;
            } | undefined;
            signedLastBlock?: boolean | undefined;
        }[] & ({
            validator?: {
                address?: Uint8Array | undefined;
                power?: bigint | undefined;
            } | undefined;
            signedLastBlock?: boolean | undefined;
        } & {
            validator?: ({
                address?: Uint8Array | undefined;
                power?: bigint | undefined;
            } & {
                address?: Uint8Array | undefined;
                power?: bigint | undefined;
            } & Record<Exclude<keyof I["votes"][number]["validator"], keyof Validator>, never>) | undefined;
            signedLastBlock?: boolean | undefined;
        } & Record<Exclude<keyof I["votes"][number], keyof VoteInfo>, never>)[] & Record<Exclude<keyof I["votes"], keyof {
            validator?: {
                address?: Uint8Array | undefined;
                power?: bigint | undefined;
            } | undefined;
            signedLastBlock?: boolean | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof CommitInfo>, never>>(object: I): CommitInfo;
};
export declare const ExtendedCommitInfo: {
    typeUrl: string;
    encode(message: ExtendedCommitInfo, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ExtendedCommitInfo;
    fromJSON(object: any): ExtendedCommitInfo;
    toJSON(message: ExtendedCommitInfo): unknown;
    fromPartial<I extends {
        round?: number | undefined;
        votes?: {
            validator?: {
                address?: Uint8Array | undefined;
                power?: bigint | undefined;
            } | undefined;
            signedLastBlock?: boolean | undefined;
            voteExtension?: Uint8Array | undefined;
        }[] | undefined;
    } & {
        round?: number | undefined;
        votes?: ({
            validator?: {
                address?: Uint8Array | undefined;
                power?: bigint | undefined;
            } | undefined;
            signedLastBlock?: boolean | undefined;
            voteExtension?: Uint8Array | undefined;
        }[] & ({
            validator?: {
                address?: Uint8Array | undefined;
                power?: bigint | undefined;
            } | undefined;
            signedLastBlock?: boolean | undefined;
            voteExtension?: Uint8Array | undefined;
        } & {
            validator?: ({
                address?: Uint8Array | undefined;
                power?: bigint | undefined;
            } & {
                address?: Uint8Array | undefined;
                power?: bigint | undefined;
            } & Record<Exclude<keyof I["votes"][number]["validator"], keyof Validator>, never>) | undefined;
            signedLastBlock?: boolean | undefined;
            voteExtension?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["votes"][number], keyof ExtendedVoteInfo>, never>)[] & Record<Exclude<keyof I["votes"], keyof {
            validator?: {
                address?: Uint8Array | undefined;
                power?: bigint | undefined;
            } | undefined;
            signedLastBlock?: boolean | undefined;
            voteExtension?: Uint8Array | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof ExtendedCommitInfo>, never>>(object: I): ExtendedCommitInfo;
};
export declare const Event: {
    typeUrl: string;
    encode(message: Event, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Event;
    fromJSON(object: any): Event;
    toJSON(message: Event): unknown;
    fromPartial<I extends {
        type?: string | undefined;
        attributes?: {
            key?: string | undefined;
            value?: string | undefined;
            index?: boolean | undefined;
        }[] | undefined;
    } & {
        type?: string | undefined;
        attributes?: ({
            key?: string | undefined;
            value?: string | undefined;
            index?: boolean | undefined;
        }[] & ({
            key?: string | undefined;
            value?: string | undefined;
            index?: boolean | undefined;
        } & {
            key?: string | undefined;
            value?: string | undefined;
            index?: boolean | undefined;
        } & Record<Exclude<keyof I["attributes"][number], keyof EventAttribute>, never>)[] & Record<Exclude<keyof I["attributes"], keyof {
            key?: string | undefined;
            value?: string | undefined;
            index?: boolean | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof Event>, never>>(object: I): Event;
};
export declare const EventAttribute: {
    typeUrl: string;
    encode(message: EventAttribute, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): EventAttribute;
    fromJSON(object: any): EventAttribute;
    toJSON(message: EventAttribute): unknown;
    fromPartial<I extends {
        key?: string | undefined;
        value?: string | undefined;
        index?: boolean | undefined;
    } & {
        key?: string | undefined;
        value?: string | undefined;
        index?: boolean | undefined;
    } & Record<Exclude<keyof I, keyof EventAttribute>, never>>(object: I): EventAttribute;
};
export declare const TxResult: {
    typeUrl: string;
    encode(message: TxResult, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): TxResult;
    fromJSON(object: any): TxResult;
    toJSON(message: TxResult): unknown;
    fromPartial<I extends {
        height?: bigint | undefined;
        index?: number | undefined;
        tx?: Uint8Array | undefined;
        result?: {
            code?: number | undefined;
            data?: Uint8Array | undefined;
            log?: string | undefined;
            info?: string | undefined;
            gasWanted?: bigint | undefined;
            gasUsed?: bigint | undefined;
            events?: {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[] | undefined;
            codespace?: string | undefined;
        } | undefined;
    } & {
        height?: bigint | undefined;
        index?: number | undefined;
        tx?: Uint8Array | undefined;
        result?: ({
            code?: number | undefined;
            data?: Uint8Array | undefined;
            log?: string | undefined;
            info?: string | undefined;
            gasWanted?: bigint | undefined;
            gasUsed?: bigint | undefined;
            events?: {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[] | undefined;
            codespace?: string | undefined;
        } & {
            code?: number | undefined;
            data?: Uint8Array | undefined;
            log?: string | undefined;
            info?: string | undefined;
            gasWanted?: bigint | undefined;
            gasUsed?: bigint | undefined;
            events?: ({
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[] & ({
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            } & {
                type?: string | undefined;
                attributes?: ({
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] & ({
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                } & {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                } & Record<Exclude<keyof I["result"]["events"][number]["attributes"][number], keyof EventAttribute>, never>)[] & Record<Exclude<keyof I["result"]["events"][number]["attributes"], keyof {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["result"]["events"][number], keyof Event>, never>)[] & Record<Exclude<keyof I["result"]["events"], keyof {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[]>, never>) | undefined;
            codespace?: string | undefined;
        } & Record<Exclude<keyof I["result"], keyof ResponseDeliverTx>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof TxResult>, never>>(object: I): TxResult;
};
export declare const Validator: {
    typeUrl: string;
    encode(message: Validator, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Validator;
    fromJSON(object: any): Validator;
    toJSON(message: Validator): unknown;
    fromPartial<I extends {
        address?: Uint8Array | undefined;
        power?: bigint | undefined;
    } & {
        address?: Uint8Array | undefined;
        power?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof Validator>, never>>(object: I): Validator;
};
export declare const ValidatorUpdate: {
    typeUrl: string;
    encode(message: ValidatorUpdate, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ValidatorUpdate;
    fromJSON(object: any): ValidatorUpdate;
    toJSON(message: ValidatorUpdate): unknown;
    fromPartial<I extends {
        pubKey?: {
            ed25519?: Uint8Array | undefined;
            secp256k1?: Uint8Array | undefined;
        } | undefined;
        power?: bigint | undefined;
    } & {
        pubKey?: ({
            ed25519?: Uint8Array | undefined;
            secp256k1?: Uint8Array | undefined;
        } & {
            ed25519?: Uint8Array | undefined;
            secp256k1?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["pubKey"], keyof PublicKey>, never>) | undefined;
        power?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof ValidatorUpdate>, never>>(object: I): ValidatorUpdate;
};
export declare const VoteInfo: {
    typeUrl: string;
    encode(message: VoteInfo, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): VoteInfo;
    fromJSON(object: any): VoteInfo;
    toJSON(message: VoteInfo): unknown;
    fromPartial<I extends {
        validator?: {
            address?: Uint8Array | undefined;
            power?: bigint | undefined;
        } | undefined;
        signedLastBlock?: boolean | undefined;
    } & {
        validator?: ({
            address?: Uint8Array | undefined;
            power?: bigint | undefined;
        } & {
            address?: Uint8Array | undefined;
            power?: bigint | undefined;
        } & Record<Exclude<keyof I["validator"], keyof Validator>, never>) | undefined;
        signedLastBlock?: boolean | undefined;
    } & Record<Exclude<keyof I, keyof VoteInfo>, never>>(object: I): VoteInfo;
};
export declare const ExtendedVoteInfo: {
    typeUrl: string;
    encode(message: ExtendedVoteInfo, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ExtendedVoteInfo;
    fromJSON(object: any): ExtendedVoteInfo;
    toJSON(message: ExtendedVoteInfo): unknown;
    fromPartial<I extends {
        validator?: {
            address?: Uint8Array | undefined;
            power?: bigint | undefined;
        } | undefined;
        signedLastBlock?: boolean | undefined;
        voteExtension?: Uint8Array | undefined;
    } & {
        validator?: ({
            address?: Uint8Array | undefined;
            power?: bigint | undefined;
        } & {
            address?: Uint8Array | undefined;
            power?: bigint | undefined;
        } & Record<Exclude<keyof I["validator"], keyof Validator>, never>) | undefined;
        signedLastBlock?: boolean | undefined;
        voteExtension?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof ExtendedVoteInfo>, never>>(object: I): ExtendedVoteInfo;
};
export declare const Misbehavior: {
    typeUrl: string;
    encode(message: Misbehavior, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Misbehavior;
    fromJSON(object: any): Misbehavior;
    toJSON(message: Misbehavior): unknown;
    fromPartial<I extends {
        type?: MisbehaviorType | undefined;
        validator?: {
            address?: Uint8Array | undefined;
            power?: bigint | undefined;
        } | undefined;
        height?: bigint | undefined;
        time?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
        totalVotingPower?: bigint | undefined;
    } & {
        type?: MisbehaviorType | undefined;
        validator?: ({
            address?: Uint8Array | undefined;
            power?: bigint | undefined;
        } & {
            address?: Uint8Array | undefined;
            power?: bigint | undefined;
        } & Record<Exclude<keyof I["validator"], keyof Validator>, never>) | undefined;
        height?: bigint | undefined;
        time?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["time"], keyof Timestamp>, never>) | undefined;
        totalVotingPower?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof Misbehavior>, never>>(object: I): Misbehavior;
};
export declare const Snapshot: {
    typeUrl: string;
    encode(message: Snapshot, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Snapshot;
    fromJSON(object: any): Snapshot;
    toJSON(message: Snapshot): unknown;
    fromPartial<I extends {
        height?: bigint | undefined;
        format?: number | undefined;
        chunks?: number | undefined;
        hash?: Uint8Array | undefined;
        metadata?: Uint8Array | undefined;
    } & {
        height?: bigint | undefined;
        format?: number | undefined;
        chunks?: number | undefined;
        hash?: Uint8Array | undefined;
        metadata?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof Snapshot>, never>>(object: I): Snapshot;
};
export interface ABCIApplication {
    Echo(request: RequestEcho): Promise<ResponseEcho>;
    Flush(request?: RequestFlush): Promise<ResponseFlush>;
    Info(request: RequestInfo): Promise<ResponseInfo>;
    DeliverTx(request: RequestDeliverTx): Promise<ResponseDeliverTx>;
    CheckTx(request: RequestCheckTx): Promise<ResponseCheckTx>;
    Query(request: RequestQuery): Promise<ResponseQuery>;
    Commit(request?: RequestCommit): Promise<ResponseCommit>;
    InitChain(request: RequestInitChain): Promise<ResponseInitChain>;
    BeginBlock(request: RequestBeginBlock): Promise<ResponseBeginBlock>;
    EndBlock(request: RequestEndBlock): Promise<ResponseEndBlock>;
    ListSnapshots(request?: RequestListSnapshots): Promise<ResponseListSnapshots>;
    OfferSnapshot(request: RequestOfferSnapshot): Promise<ResponseOfferSnapshot>;
    LoadSnapshotChunk(request: RequestLoadSnapshotChunk): Promise<ResponseLoadSnapshotChunk>;
    ApplySnapshotChunk(request: RequestApplySnapshotChunk): Promise<ResponseApplySnapshotChunk>;
    PrepareProposal(request: RequestPrepareProposal): Promise<ResponsePrepareProposal>;
    ProcessProposal(request: RequestProcessProposal): Promise<ResponseProcessProposal>;
}
export declare class ABCIApplicationClientImpl implements ABCIApplication {
    private readonly rpc;
    constructor(rpc: Rpc);
    Echo(request: RequestEcho): Promise<ResponseEcho>;
    Flush(request?: RequestFlush): Promise<ResponseFlush>;
    Info(request: RequestInfo): Promise<ResponseInfo>;
    DeliverTx(request: RequestDeliverTx): Promise<ResponseDeliverTx>;
    CheckTx(request: RequestCheckTx): Promise<ResponseCheckTx>;
    Query(request: RequestQuery): Promise<ResponseQuery>;
    Commit(request?: RequestCommit): Promise<ResponseCommit>;
    InitChain(request: RequestInitChain): Promise<ResponseInitChain>;
    BeginBlock(request: RequestBeginBlock): Promise<ResponseBeginBlock>;
    EndBlock(request: RequestEndBlock): Promise<ResponseEndBlock>;
    ListSnapshots(request?: RequestListSnapshots): Promise<ResponseListSnapshots>;
    OfferSnapshot(request: RequestOfferSnapshot): Promise<ResponseOfferSnapshot>;
    LoadSnapshotChunk(request: RequestLoadSnapshotChunk): Promise<ResponseLoadSnapshotChunk>;
    ApplySnapshotChunk(request: RequestApplySnapshotChunk): Promise<ResponseApplySnapshotChunk>;
    PrepareProposal(request: RequestPrepareProposal): Promise<ResponsePrepareProposal>;
    ProcessProposal(request: RequestProcessProposal): Promise<ResponseProcessProposal>;
}
