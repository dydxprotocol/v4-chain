import { Duration } from "../../../../google/protobuf/duration";
import { Height } from "../../../core/client/v1/client";
import { ProofSpec } from "../../../../cosmos/ics23/v1/proofs";
import { Timestamp } from "../../../../google/protobuf/timestamp";
import { MerkleRoot } from "../../../core/commitment/v1/commitment";
import { SignedHeader } from "../../../../tendermint/types/types";
import { ValidatorSet } from "../../../../tendermint/types/validator";
import { BinaryReader, BinaryWriter } from "../../../../binary";
export declare const protobufPackage = "ibc.lightclients.tendermint.v1";
/**
 * ClientState from Tendermint tracks the current validator set, latest height,
 * and a possible frozen height.
 */
export interface ClientState {
    chainId: string;
    trustLevel: Fraction;
    /**
     * duration of the period since the LastestTimestamp during which the
     * submitted headers are valid for upgrade
     */
    trustingPeriod: Duration;
    /** duration of the staking unbonding period */
    unbondingPeriod: Duration;
    /** defines how much new (untrusted) header's Time can drift into the future. */
    maxClockDrift: Duration;
    /** Block height when the client was frozen due to a misbehaviour */
    frozenHeight: Height;
    /** Latest height the client was updated to */
    latestHeight: Height;
    /** Proof specifications used in verifying counterparty state */
    proofSpecs: ProofSpec[];
    /**
     * Path at which next upgraded client will be committed.
     * Each element corresponds to the key for a single CommitmentProof in the
     * chained proof. NOTE: ClientState must stored under
     * `{upgradePath}/{upgradeHeight}/clientState` ConsensusState must be stored
     * under `{upgradepath}/{upgradeHeight}/consensusState` For SDK chains using
     * the default upgrade module, upgrade_path should be []string{"upgrade",
     * "upgradedIBCState"}`
     */
    upgradePath: string[];
    /** allow_update_after_expiry is deprecated */
    /** @deprecated */
    allowUpdateAfterExpiry: boolean;
    /** allow_update_after_misbehaviour is deprecated */
    /** @deprecated */
    allowUpdateAfterMisbehaviour: boolean;
}
/** ConsensusState defines the consensus state from Tendermint. */
export interface ConsensusState {
    /**
     * timestamp that corresponds to the block height in which the ConsensusState
     * was stored.
     */
    timestamp: Timestamp;
    /** commitment root (i.e app hash) */
    root: MerkleRoot;
    nextValidatorsHash: Uint8Array;
}
/**
 * Misbehaviour is a wrapper over two conflicting Headers
 * that implements Misbehaviour interface expected by ICS-02
 */
export interface Misbehaviour {
    /** ClientID is deprecated */
    /** @deprecated */
    clientId: string;
    header1?: Header;
    header2?: Header;
}
/**
 * Header defines the Tendermint client consensus Header.
 * It encapsulates all the information necessary to update from a trusted
 * Tendermint ConsensusState. The inclusion of TrustedHeight and
 * TrustedValidators allows this update to process correctly, so long as the
 * ConsensusState for the TrustedHeight exists, this removes race conditions
 * among relayers The SignedHeader and ValidatorSet are the new untrusted update
 * fields for the client. The TrustedHeight is the height of a stored
 * ConsensusState on the client that will be used to verify the new untrusted
 * header. The Trusted ConsensusState must be within the unbonding period of
 * current time in order to correctly verify, and the TrustedValidators must
 * hash to TrustedConsensusState.NextValidatorsHash since that is the last
 * trusted validator set at the TrustedHeight.
 */
export interface Header {
    signedHeader?: SignedHeader;
    validatorSet?: ValidatorSet;
    trustedHeight: Height;
    trustedValidators?: ValidatorSet;
}
/**
 * Fraction defines the protobuf message type for tmmath.Fraction that only
 * supports positive values.
 */
export interface Fraction {
    numerator: bigint;
    denominator: bigint;
}
export declare const ClientState: {
    typeUrl: string;
    encode(message: ClientState, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ClientState;
    fromJSON(object: any): ClientState;
    toJSON(message: ClientState): unknown;
    fromPartial<I extends {
        chainId?: string | undefined;
        trustLevel?: {
            numerator?: bigint | undefined;
            denominator?: bigint | undefined;
        } | undefined;
        trustingPeriod?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
        unbondingPeriod?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
        maxClockDrift?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
        frozenHeight?: {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } | undefined;
        latestHeight?: {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } | undefined;
        proofSpecs?: {
            leafSpec?: {
                hash?: import("../../../../cosmos/ics23/v1/proofs").HashOp | undefined;
                prehashKey?: import("../../../../cosmos/ics23/v1/proofs").HashOp | undefined;
                prehashValue?: import("../../../../cosmos/ics23/v1/proofs").HashOp | undefined;
                length?: import("../../../../cosmos/ics23/v1/proofs").LengthOp | undefined;
                prefix?: Uint8Array | undefined;
            } | undefined;
            innerSpec?: {
                childOrder?: number[] | undefined;
                childSize?: number | undefined;
                minPrefixLength?: number | undefined;
                maxPrefixLength?: number | undefined;
                emptyChild?: Uint8Array | undefined;
                hash?: import("../../../../cosmos/ics23/v1/proofs").HashOp | undefined;
            } | undefined;
            maxDepth?: number | undefined;
            minDepth?: number | undefined;
        }[] | undefined;
        upgradePath?: string[] | undefined;
        allowUpdateAfterExpiry?: boolean | undefined;
        allowUpdateAfterMisbehaviour?: boolean | undefined;
    } & {
        chainId?: string | undefined;
        trustLevel?: ({
            numerator?: bigint | undefined;
            denominator?: bigint | undefined;
        } & {
            numerator?: bigint | undefined;
            denominator?: bigint | undefined;
        } & Record<Exclude<keyof I["trustLevel"], keyof Fraction>, never>) | undefined;
        trustingPeriod?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["trustingPeriod"], keyof Duration>, never>) | undefined;
        unbondingPeriod?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["unbondingPeriod"], keyof Duration>, never>) | undefined;
        maxClockDrift?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["maxClockDrift"], keyof Duration>, never>) | undefined;
        frozenHeight?: ({
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["frozenHeight"], keyof Height>, never>) | undefined;
        latestHeight?: ({
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["latestHeight"], keyof Height>, never>) | undefined;
        proofSpecs?: ({
            leafSpec?: {
                hash?: import("../../../../cosmos/ics23/v1/proofs").HashOp | undefined;
                prehashKey?: import("../../../../cosmos/ics23/v1/proofs").HashOp | undefined;
                prehashValue?: import("../../../../cosmos/ics23/v1/proofs").HashOp | undefined;
                length?: import("../../../../cosmos/ics23/v1/proofs").LengthOp | undefined;
                prefix?: Uint8Array | undefined;
            } | undefined;
            innerSpec?: {
                childOrder?: number[] | undefined;
                childSize?: number | undefined;
                minPrefixLength?: number | undefined;
                maxPrefixLength?: number | undefined;
                emptyChild?: Uint8Array | undefined;
                hash?: import("../../../../cosmos/ics23/v1/proofs").HashOp | undefined;
            } | undefined;
            maxDepth?: number | undefined;
            minDepth?: number | undefined;
        }[] & ({
            leafSpec?: {
                hash?: import("../../../../cosmos/ics23/v1/proofs").HashOp | undefined;
                prehashKey?: import("../../../../cosmos/ics23/v1/proofs").HashOp | undefined;
                prehashValue?: import("../../../../cosmos/ics23/v1/proofs").HashOp | undefined;
                length?: import("../../../../cosmos/ics23/v1/proofs").LengthOp | undefined;
                prefix?: Uint8Array | undefined;
            } | undefined;
            innerSpec?: {
                childOrder?: number[] | undefined;
                childSize?: number | undefined;
                minPrefixLength?: number | undefined;
                maxPrefixLength?: number | undefined;
                emptyChild?: Uint8Array | undefined;
                hash?: import("../../../../cosmos/ics23/v1/proofs").HashOp | undefined;
            } | undefined;
            maxDepth?: number | undefined;
            minDepth?: number | undefined;
        } & {
            leafSpec?: ({
                hash?: import("../../../../cosmos/ics23/v1/proofs").HashOp | undefined;
                prehashKey?: import("../../../../cosmos/ics23/v1/proofs").HashOp | undefined;
                prehashValue?: import("../../../../cosmos/ics23/v1/proofs").HashOp | undefined;
                length?: import("../../../../cosmos/ics23/v1/proofs").LengthOp | undefined;
                prefix?: Uint8Array | undefined;
            } & {
                hash?: import("../../../../cosmos/ics23/v1/proofs").HashOp | undefined;
                prehashKey?: import("../../../../cosmos/ics23/v1/proofs").HashOp | undefined;
                prehashValue?: import("../../../../cosmos/ics23/v1/proofs").HashOp | undefined;
                length?: import("../../../../cosmos/ics23/v1/proofs").LengthOp | undefined;
                prefix?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["proofSpecs"][number]["leafSpec"], keyof import("../../../../cosmos/ics23/v1/proofs").LeafOp>, never>) | undefined;
            innerSpec?: ({
                childOrder?: number[] | undefined;
                childSize?: number | undefined;
                minPrefixLength?: number | undefined;
                maxPrefixLength?: number | undefined;
                emptyChild?: Uint8Array | undefined;
                hash?: import("../../../../cosmos/ics23/v1/proofs").HashOp | undefined;
            } & {
                childOrder?: (number[] & number[] & Record<Exclude<keyof I["proofSpecs"][number]["innerSpec"]["childOrder"], keyof number[]>, never>) | undefined;
                childSize?: number | undefined;
                minPrefixLength?: number | undefined;
                maxPrefixLength?: number | undefined;
                emptyChild?: Uint8Array | undefined;
                hash?: import("../../../../cosmos/ics23/v1/proofs").HashOp | undefined;
            } & Record<Exclude<keyof I["proofSpecs"][number]["innerSpec"], keyof import("../../../../cosmos/ics23/v1/proofs").InnerSpec>, never>) | undefined;
            maxDepth?: number | undefined;
            minDepth?: number | undefined;
        } & Record<Exclude<keyof I["proofSpecs"][number], keyof ProofSpec>, never>)[] & Record<Exclude<keyof I["proofSpecs"], keyof {
            leafSpec?: {
                hash?: import("../../../../cosmos/ics23/v1/proofs").HashOp | undefined;
                prehashKey?: import("../../../../cosmos/ics23/v1/proofs").HashOp | undefined;
                prehashValue?: import("../../../../cosmos/ics23/v1/proofs").HashOp | undefined;
                length?: import("../../../../cosmos/ics23/v1/proofs").LengthOp | undefined;
                prefix?: Uint8Array | undefined;
            } | undefined;
            innerSpec?: {
                childOrder?: number[] | undefined;
                childSize?: number | undefined;
                minPrefixLength?: number | undefined;
                maxPrefixLength?: number | undefined;
                emptyChild?: Uint8Array | undefined;
                hash?: import("../../../../cosmos/ics23/v1/proofs").HashOp | undefined;
            } | undefined;
            maxDepth?: number | undefined;
            minDepth?: number | undefined;
        }[]>, never>) | undefined;
        upgradePath?: (string[] & string[] & Record<Exclude<keyof I["upgradePath"], keyof string[]>, never>) | undefined;
        allowUpdateAfterExpiry?: boolean | undefined;
        allowUpdateAfterMisbehaviour?: boolean | undefined;
    } & Record<Exclude<keyof I, keyof ClientState>, never>>(object: I): ClientState;
};
export declare const ConsensusState: {
    typeUrl: string;
    encode(message: ConsensusState, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ConsensusState;
    fromJSON(object: any): ConsensusState;
    toJSON(message: ConsensusState): unknown;
    fromPartial<I extends {
        timestamp?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
        root?: {
            hash?: Uint8Array | undefined;
        } | undefined;
        nextValidatorsHash?: Uint8Array | undefined;
    } & {
        timestamp?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["timestamp"], keyof Timestamp>, never>) | undefined;
        root?: ({
            hash?: Uint8Array | undefined;
        } & {
            hash?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["root"], "hash">, never>) | undefined;
        nextValidatorsHash?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof ConsensusState>, never>>(object: I): ConsensusState;
};
export declare const Misbehaviour: {
    typeUrl: string;
    encode(message: Misbehaviour, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Misbehaviour;
    fromJSON(object: any): Misbehaviour;
    toJSON(message: Misbehaviour): unknown;
    fromPartial<I extends {
        clientId?: string | undefined;
        header1?: {
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
            trustedHeight?: {
                revisionNumber?: bigint | undefined;
                revisionHeight?: bigint | undefined;
            } | undefined;
            trustedValidators?: {
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
        header2?: {
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
            trustedHeight?: {
                revisionNumber?: bigint | undefined;
                revisionHeight?: bigint | undefined;
            } | undefined;
            trustedValidators?: {
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
    } & {
        clientId?: string | undefined;
        header1?: ({
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
            trustedHeight?: {
                revisionNumber?: bigint | undefined;
                revisionHeight?: bigint | undefined;
            } | undefined;
            trustedValidators?: {
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
                    } & Record<Exclude<keyof I["header1"]["signedHeader"]["header"]["version"], keyof import("../../../../tendermint/version/types").Consensus>, never>) | undefined;
                    chainId?: string | undefined;
                    height?: bigint | undefined;
                    time?: ({
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } & {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } & Record<Exclude<keyof I["header1"]["signedHeader"]["header"]["time"], keyof Timestamp>, never>) | undefined;
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
                        } & Record<Exclude<keyof I["header1"]["signedHeader"]["header"]["lastBlockId"]["partSetHeader"], keyof import("../../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                    } & Record<Exclude<keyof I["header1"]["signedHeader"]["header"]["lastBlockId"], keyof import("../../../../tendermint/types/types").BlockID>, never>) | undefined;
                    lastCommitHash?: Uint8Array | undefined;
                    dataHash?: Uint8Array | undefined;
                    validatorsHash?: Uint8Array | undefined;
                    nextValidatorsHash?: Uint8Array | undefined;
                    consensusHash?: Uint8Array | undefined;
                    appHash?: Uint8Array | undefined;
                    lastResultsHash?: Uint8Array | undefined;
                    evidenceHash?: Uint8Array | undefined;
                    proposerAddress?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["header1"]["signedHeader"]["header"], keyof import("../../../../tendermint/types/types").Header>, never>) | undefined;
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
                        } & Record<Exclude<keyof I["header1"]["signedHeader"]["commit"]["blockId"]["partSetHeader"], keyof import("../../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                    } & Record<Exclude<keyof I["header1"]["signedHeader"]["commit"]["blockId"], keyof import("../../../../tendermint/types/types").BlockID>, never>) | undefined;
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
                        } & Record<Exclude<keyof I["header1"]["signedHeader"]["commit"]["signatures"][number]["timestamp"], keyof Timestamp>, never>) | undefined;
                        signature?: Uint8Array | undefined;
                    } & Record<Exclude<keyof I["header1"]["signedHeader"]["commit"]["signatures"][number], keyof import("../../../../tendermint/types/types").CommitSig>, never>)[] & Record<Exclude<keyof I["header1"]["signedHeader"]["commit"]["signatures"], keyof {
                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                        validatorAddress?: Uint8Array | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                        signature?: Uint8Array | undefined;
                    }[]>, never>) | undefined;
                } & Record<Exclude<keyof I["header1"]["signedHeader"]["commit"], keyof import("../../../../tendermint/types/types").Commit>, never>) | undefined;
            } & Record<Exclude<keyof I["header1"]["signedHeader"], keyof SignedHeader>, never>) | undefined;
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
                    } & Record<Exclude<keyof I["header1"]["validatorSet"]["validators"][number]["pubKey"], keyof import("../../../../tendermint/crypto/keys").PublicKey>, never>) | undefined;
                    votingPower?: bigint | undefined;
                    proposerPriority?: bigint | undefined;
                } & Record<Exclude<keyof I["header1"]["validatorSet"]["validators"][number], keyof import("../../../../tendermint/types/validator").Validator>, never>)[] & Record<Exclude<keyof I["header1"]["validatorSet"]["validators"], keyof {
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
                    } & Record<Exclude<keyof I["header1"]["validatorSet"]["proposer"]["pubKey"], keyof import("../../../../tendermint/crypto/keys").PublicKey>, never>) | undefined;
                    votingPower?: bigint | undefined;
                    proposerPriority?: bigint | undefined;
                } & Record<Exclude<keyof I["header1"]["validatorSet"]["proposer"], keyof import("../../../../tendermint/types/validator").Validator>, never>) | undefined;
                totalVotingPower?: bigint | undefined;
            } & Record<Exclude<keyof I["header1"]["validatorSet"], keyof ValidatorSet>, never>) | undefined;
            trustedHeight?: ({
                revisionNumber?: bigint | undefined;
                revisionHeight?: bigint | undefined;
            } & {
                revisionNumber?: bigint | undefined;
                revisionHeight?: bigint | undefined;
            } & Record<Exclude<keyof I["header1"]["trustedHeight"], keyof Height>, never>) | undefined;
            trustedValidators?: ({
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
                    } & Record<Exclude<keyof I["header1"]["trustedValidators"]["validators"][number]["pubKey"], keyof import("../../../../tendermint/crypto/keys").PublicKey>, never>) | undefined;
                    votingPower?: bigint | undefined;
                    proposerPriority?: bigint | undefined;
                } & Record<Exclude<keyof I["header1"]["trustedValidators"]["validators"][number], keyof import("../../../../tendermint/types/validator").Validator>, never>)[] & Record<Exclude<keyof I["header1"]["trustedValidators"]["validators"], keyof {
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
                    } & Record<Exclude<keyof I["header1"]["trustedValidators"]["proposer"]["pubKey"], keyof import("../../../../tendermint/crypto/keys").PublicKey>, never>) | undefined;
                    votingPower?: bigint | undefined;
                    proposerPriority?: bigint | undefined;
                } & Record<Exclude<keyof I["header1"]["trustedValidators"]["proposer"], keyof import("../../../../tendermint/types/validator").Validator>, never>) | undefined;
                totalVotingPower?: bigint | undefined;
            } & Record<Exclude<keyof I["header1"]["trustedValidators"], keyof ValidatorSet>, never>) | undefined;
        } & Record<Exclude<keyof I["header1"], keyof Header>, never>) | undefined;
        header2?: ({
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
            trustedHeight?: {
                revisionNumber?: bigint | undefined;
                revisionHeight?: bigint | undefined;
            } | undefined;
            trustedValidators?: {
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
                    } & Record<Exclude<keyof I["header2"]["signedHeader"]["header"]["version"], keyof import("../../../../tendermint/version/types").Consensus>, never>) | undefined;
                    chainId?: string | undefined;
                    height?: bigint | undefined;
                    time?: ({
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } & {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } & Record<Exclude<keyof I["header2"]["signedHeader"]["header"]["time"], keyof Timestamp>, never>) | undefined;
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
                        } & Record<Exclude<keyof I["header2"]["signedHeader"]["header"]["lastBlockId"]["partSetHeader"], keyof import("../../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                    } & Record<Exclude<keyof I["header2"]["signedHeader"]["header"]["lastBlockId"], keyof import("../../../../tendermint/types/types").BlockID>, never>) | undefined;
                    lastCommitHash?: Uint8Array | undefined;
                    dataHash?: Uint8Array | undefined;
                    validatorsHash?: Uint8Array | undefined;
                    nextValidatorsHash?: Uint8Array | undefined;
                    consensusHash?: Uint8Array | undefined;
                    appHash?: Uint8Array | undefined;
                    lastResultsHash?: Uint8Array | undefined;
                    evidenceHash?: Uint8Array | undefined;
                    proposerAddress?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["header2"]["signedHeader"]["header"], keyof import("../../../../tendermint/types/types").Header>, never>) | undefined;
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
                        } & Record<Exclude<keyof I["header2"]["signedHeader"]["commit"]["blockId"]["partSetHeader"], keyof import("../../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                    } & Record<Exclude<keyof I["header2"]["signedHeader"]["commit"]["blockId"], keyof import("../../../../tendermint/types/types").BlockID>, never>) | undefined;
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
                        } & Record<Exclude<keyof I["header2"]["signedHeader"]["commit"]["signatures"][number]["timestamp"], keyof Timestamp>, never>) | undefined;
                        signature?: Uint8Array | undefined;
                    } & Record<Exclude<keyof I["header2"]["signedHeader"]["commit"]["signatures"][number], keyof import("../../../../tendermint/types/types").CommitSig>, never>)[] & Record<Exclude<keyof I["header2"]["signedHeader"]["commit"]["signatures"], keyof {
                        blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                        validatorAddress?: Uint8Array | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                        signature?: Uint8Array | undefined;
                    }[]>, never>) | undefined;
                } & Record<Exclude<keyof I["header2"]["signedHeader"]["commit"], keyof import("../../../../tendermint/types/types").Commit>, never>) | undefined;
            } & Record<Exclude<keyof I["header2"]["signedHeader"], keyof SignedHeader>, never>) | undefined;
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
                    } & Record<Exclude<keyof I["header2"]["validatorSet"]["validators"][number]["pubKey"], keyof import("../../../../tendermint/crypto/keys").PublicKey>, never>) | undefined;
                    votingPower?: bigint | undefined;
                    proposerPriority?: bigint | undefined;
                } & Record<Exclude<keyof I["header2"]["validatorSet"]["validators"][number], keyof import("../../../../tendermint/types/validator").Validator>, never>)[] & Record<Exclude<keyof I["header2"]["validatorSet"]["validators"], keyof {
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
                    } & Record<Exclude<keyof I["header2"]["validatorSet"]["proposer"]["pubKey"], keyof import("../../../../tendermint/crypto/keys").PublicKey>, never>) | undefined;
                    votingPower?: bigint | undefined;
                    proposerPriority?: bigint | undefined;
                } & Record<Exclude<keyof I["header2"]["validatorSet"]["proposer"], keyof import("../../../../tendermint/types/validator").Validator>, never>) | undefined;
                totalVotingPower?: bigint | undefined;
            } & Record<Exclude<keyof I["header2"]["validatorSet"], keyof ValidatorSet>, never>) | undefined;
            trustedHeight?: ({
                revisionNumber?: bigint | undefined;
                revisionHeight?: bigint | undefined;
            } & {
                revisionNumber?: bigint | undefined;
                revisionHeight?: bigint | undefined;
            } & Record<Exclude<keyof I["header2"]["trustedHeight"], keyof Height>, never>) | undefined;
            trustedValidators?: ({
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
                    } & Record<Exclude<keyof I["header2"]["trustedValidators"]["validators"][number]["pubKey"], keyof import("../../../../tendermint/crypto/keys").PublicKey>, never>) | undefined;
                    votingPower?: bigint | undefined;
                    proposerPriority?: bigint | undefined;
                } & Record<Exclude<keyof I["header2"]["trustedValidators"]["validators"][number], keyof import("../../../../tendermint/types/validator").Validator>, never>)[] & Record<Exclude<keyof I["header2"]["trustedValidators"]["validators"], keyof {
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
                    } & Record<Exclude<keyof I["header2"]["trustedValidators"]["proposer"]["pubKey"], keyof import("../../../../tendermint/crypto/keys").PublicKey>, never>) | undefined;
                    votingPower?: bigint | undefined;
                    proposerPriority?: bigint | undefined;
                } & Record<Exclude<keyof I["header2"]["trustedValidators"]["proposer"], keyof import("../../../../tendermint/types/validator").Validator>, never>) | undefined;
                totalVotingPower?: bigint | undefined;
            } & Record<Exclude<keyof I["header2"]["trustedValidators"], keyof ValidatorSet>, never>) | undefined;
        } & Record<Exclude<keyof I["header2"], keyof Header>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof Misbehaviour>, never>>(object: I): Misbehaviour;
};
export declare const Header: {
    typeUrl: string;
    encode(message: Header, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Header;
    fromJSON(object: any): Header;
    toJSON(message: Header): unknown;
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
        trustedHeight?: {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } | undefined;
        trustedValidators?: {
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
                } & Record<Exclude<keyof I["signedHeader"]["header"]["version"], keyof import("../../../../tendermint/version/types").Consensus>, never>) | undefined;
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
                    } & Record<Exclude<keyof I["signedHeader"]["header"]["lastBlockId"]["partSetHeader"], keyof import("../../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                } & Record<Exclude<keyof I["signedHeader"]["header"]["lastBlockId"], keyof import("../../../../tendermint/types/types").BlockID>, never>) | undefined;
                lastCommitHash?: Uint8Array | undefined;
                dataHash?: Uint8Array | undefined;
                validatorsHash?: Uint8Array | undefined;
                nextValidatorsHash?: Uint8Array | undefined;
                consensusHash?: Uint8Array | undefined;
                appHash?: Uint8Array | undefined;
                lastResultsHash?: Uint8Array | undefined;
                evidenceHash?: Uint8Array | undefined;
                proposerAddress?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["signedHeader"]["header"], keyof import("../../../../tendermint/types/types").Header>, never>) | undefined;
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
                    } & Record<Exclude<keyof I["signedHeader"]["commit"]["blockId"]["partSetHeader"], keyof import("../../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                } & Record<Exclude<keyof I["signedHeader"]["commit"]["blockId"], keyof import("../../../../tendermint/types/types").BlockID>, never>) | undefined;
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
                    } & Record<Exclude<keyof I["signedHeader"]["commit"]["signatures"][number]["timestamp"], keyof Timestamp>, never>) | undefined;
                    signature?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["signedHeader"]["commit"]["signatures"][number], keyof import("../../../../tendermint/types/types").CommitSig>, never>)[] & Record<Exclude<keyof I["signedHeader"]["commit"]["signatures"], keyof {
                    blockIdFlag?: import("../../../../tendermint/types/types").BlockIDFlag | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    timestamp?: {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } | undefined;
                    signature?: Uint8Array | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["signedHeader"]["commit"], keyof import("../../../../tendermint/types/types").Commit>, never>) | undefined;
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
                } & Record<Exclude<keyof I["validatorSet"]["validators"][number]["pubKey"], keyof import("../../../../tendermint/crypto/keys").PublicKey>, never>) | undefined;
                votingPower?: bigint | undefined;
                proposerPriority?: bigint | undefined;
            } & Record<Exclude<keyof I["validatorSet"]["validators"][number], keyof import("../../../../tendermint/types/validator").Validator>, never>)[] & Record<Exclude<keyof I["validatorSet"]["validators"], keyof {
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
                } & Record<Exclude<keyof I["validatorSet"]["proposer"]["pubKey"], keyof import("../../../../tendermint/crypto/keys").PublicKey>, never>) | undefined;
                votingPower?: bigint | undefined;
                proposerPriority?: bigint | undefined;
            } & Record<Exclude<keyof I["validatorSet"]["proposer"], keyof import("../../../../tendermint/types/validator").Validator>, never>) | undefined;
            totalVotingPower?: bigint | undefined;
        } & Record<Exclude<keyof I["validatorSet"], keyof ValidatorSet>, never>) | undefined;
        trustedHeight?: ({
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["trustedHeight"], keyof Height>, never>) | undefined;
        trustedValidators?: ({
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
                } & Record<Exclude<keyof I["trustedValidators"]["validators"][number]["pubKey"], keyof import("../../../../tendermint/crypto/keys").PublicKey>, never>) | undefined;
                votingPower?: bigint | undefined;
                proposerPriority?: bigint | undefined;
            } & Record<Exclude<keyof I["trustedValidators"]["validators"][number], keyof import("../../../../tendermint/types/validator").Validator>, never>)[] & Record<Exclude<keyof I["trustedValidators"]["validators"], keyof {
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
                } & Record<Exclude<keyof I["trustedValidators"]["proposer"]["pubKey"], keyof import("../../../../tendermint/crypto/keys").PublicKey>, never>) | undefined;
                votingPower?: bigint | undefined;
                proposerPriority?: bigint | undefined;
            } & Record<Exclude<keyof I["trustedValidators"]["proposer"], keyof import("../../../../tendermint/types/validator").Validator>, never>) | undefined;
            totalVotingPower?: bigint | undefined;
        } & Record<Exclude<keyof I["trustedValidators"], keyof ValidatorSet>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof Header>, never>>(object: I): Header;
};
export declare const Fraction: {
    typeUrl: string;
    encode(message: Fraction, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Fraction;
    fromJSON(object: any): Fraction;
    toJSON(message: Fraction): unknown;
    fromPartial<I extends {
        numerator?: bigint | undefined;
        denominator?: bigint | undefined;
    } & {
        numerator?: bigint | undefined;
        denominator?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof Fraction>, never>>(object: I): Fraction;
};
