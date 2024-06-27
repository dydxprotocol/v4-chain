import { Vote, LightBlock } from "./types";
import { Timestamp } from "../../google/protobuf/timestamp";
import { Validator } from "./validator";
import { BinaryReader, BinaryWriter } from "../../binary";
export declare const protobufPackage = "tendermint.types";
export interface Evidence {
    duplicateVoteEvidence?: DuplicateVoteEvidence;
    lightClientAttackEvidence?: LightClientAttackEvidence;
}
/** DuplicateVoteEvidence contains evidence of a validator signed two conflicting votes. */
export interface DuplicateVoteEvidence {
    voteA?: Vote;
    voteB?: Vote;
    totalVotingPower: bigint;
    validatorPower: bigint;
    timestamp: Timestamp;
}
/** LightClientAttackEvidence contains evidence of a set of validators attempting to mislead a light client. */
export interface LightClientAttackEvidence {
    conflictingBlock?: LightBlock;
    commonHeight: bigint;
    byzantineValidators: Validator[];
    totalVotingPower: bigint;
    timestamp: Timestamp;
}
export interface EvidenceList {
    evidence: Evidence[];
}
export declare const Evidence: {
    typeUrl: string;
    encode(message: Evidence, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Evidence;
    fromJSON(object: any): Evidence;
    toJSON(message: Evidence): unknown;
    fromPartial<I extends {
        duplicateVoteEvidence?: {
            voteA?: {
                type?: import("./types").SignedMsgType | undefined;
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
                type?: import("./types").SignedMsgType | undefined;
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
                            blockIdFlag?: import("./types").BlockIDFlag | undefined;
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
                type?: import("./types").SignedMsgType | undefined;
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
                type?: import("./types").SignedMsgType | undefined;
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
                type?: import("./types").SignedMsgType | undefined;
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
                type?: import("./types").SignedMsgType | undefined;
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
                    } & Record<Exclude<keyof I["duplicateVoteEvidence"]["voteA"]["blockId"]["partSetHeader"], keyof import("./types").PartSetHeader>, never>) | undefined;
                } & Record<Exclude<keyof I["duplicateVoteEvidence"]["voteA"]["blockId"], keyof import("./types").BlockID>, never>) | undefined;
                timestamp?: ({
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & Record<Exclude<keyof I["duplicateVoteEvidence"]["voteA"]["timestamp"], keyof Timestamp>, never>) | undefined;
                validatorAddress?: Uint8Array | undefined;
                validatorIndex?: number | undefined;
                signature?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["duplicateVoteEvidence"]["voteA"], keyof Vote>, never>) | undefined;
            voteB?: ({
                type?: import("./types").SignedMsgType | undefined;
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
                type?: import("./types").SignedMsgType | undefined;
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
                    } & Record<Exclude<keyof I["duplicateVoteEvidence"]["voteB"]["blockId"]["partSetHeader"], keyof import("./types").PartSetHeader>, never>) | undefined;
                } & Record<Exclude<keyof I["duplicateVoteEvidence"]["voteB"]["blockId"], keyof import("./types").BlockID>, never>) | undefined;
                timestamp?: ({
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & Record<Exclude<keyof I["duplicateVoteEvidence"]["voteB"]["timestamp"], keyof Timestamp>, never>) | undefined;
                validatorAddress?: Uint8Array | undefined;
                validatorIndex?: number | undefined;
                signature?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["duplicateVoteEvidence"]["voteB"], keyof Vote>, never>) | undefined;
            totalVotingPower?: bigint | undefined;
            validatorPower?: bigint | undefined;
            timestamp?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["duplicateVoteEvidence"]["timestamp"], keyof Timestamp>, never>) | undefined;
        } & Record<Exclude<keyof I["duplicateVoteEvidence"], keyof DuplicateVoteEvidence>, never>) | undefined;
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
                            blockIdFlag?: import("./types").BlockIDFlag | undefined;
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
                            blockIdFlag?: import("./types").BlockIDFlag | undefined;
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
                            blockIdFlag?: import("./types").BlockIDFlag | undefined;
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
                        } & Record<Exclude<keyof I["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"]["version"], keyof import("../version/types").Consensus>, never>) | undefined;
                        chainId?: string | undefined;
                        height?: bigint | undefined;
                        time?: ({
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } & {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } & Record<Exclude<keyof I["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"]["time"], keyof Timestamp>, never>) | undefined;
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
                            } & Record<Exclude<keyof I["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"]["lastBlockId"]["partSetHeader"], keyof import("./types").PartSetHeader>, never>) | undefined;
                        } & Record<Exclude<keyof I["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"]["lastBlockId"], keyof import("./types").BlockID>, never>) | undefined;
                        lastCommitHash?: Uint8Array | undefined;
                        dataHash?: Uint8Array | undefined;
                        validatorsHash?: Uint8Array | undefined;
                        nextValidatorsHash?: Uint8Array | undefined;
                        consensusHash?: Uint8Array | undefined;
                        appHash?: Uint8Array | undefined;
                        lastResultsHash?: Uint8Array | undefined;
                        evidenceHash?: Uint8Array | undefined;
                        proposerAddress?: Uint8Array | undefined;
                    } & Record<Exclude<keyof I["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"], keyof import("./types").Header>, never>) | undefined;
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
                            blockIdFlag?: import("./types").BlockIDFlag | undefined;
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
                            } & Record<Exclude<keyof I["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["blockId"]["partSetHeader"], keyof import("./types").PartSetHeader>, never>) | undefined;
                        } & Record<Exclude<keyof I["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["blockId"], keyof import("./types").BlockID>, never>) | undefined;
                        signatures?: ({
                            blockIdFlag?: import("./types").BlockIDFlag | undefined;
                            validatorAddress?: Uint8Array | undefined;
                            timestamp?: {
                                seconds?: bigint | undefined;
                                nanos?: number | undefined;
                            } | undefined;
                            signature?: Uint8Array | undefined;
                        }[] & ({
                            blockIdFlag?: import("./types").BlockIDFlag | undefined;
                            validatorAddress?: Uint8Array | undefined;
                            timestamp?: {
                                seconds?: bigint | undefined;
                                nanos?: number | undefined;
                            } | undefined;
                            signature?: Uint8Array | undefined;
                        } & {
                            blockIdFlag?: import("./types").BlockIDFlag | undefined;
                            validatorAddress?: Uint8Array | undefined;
                            timestamp?: ({
                                seconds?: bigint | undefined;
                                nanos?: number | undefined;
                            } & {
                                seconds?: bigint | undefined;
                                nanos?: number | undefined;
                            } & Record<Exclude<keyof I["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["signatures"][number]["timestamp"], keyof Timestamp>, never>) | undefined;
                            signature?: Uint8Array | undefined;
                        } & Record<Exclude<keyof I["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["signatures"][number], keyof import("./types").CommitSig>, never>)[] & Record<Exclude<keyof I["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["signatures"], keyof {
                            blockIdFlag?: import("./types").BlockIDFlag | undefined;
                            validatorAddress?: Uint8Array | undefined;
                            timestamp?: {
                                seconds?: bigint | undefined;
                                nanos?: number | undefined;
                            } | undefined;
                            signature?: Uint8Array | undefined;
                        }[]>, never>) | undefined;
                    } & Record<Exclude<keyof I["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"], keyof import("./types").Commit>, never>) | undefined;
                } & Record<Exclude<keyof I["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"], keyof import("./types").SignedHeader>, never>) | undefined;
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
                        } & Record<Exclude<keyof I["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["validators"][number]["pubKey"], keyof import("../crypto/keys").PublicKey>, never>) | undefined;
                        votingPower?: bigint | undefined;
                        proposerPriority?: bigint | undefined;
                    } & Record<Exclude<keyof I["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["validators"][number], keyof Validator>, never>)[] & Record<Exclude<keyof I["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["validators"], keyof {
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
                        } & Record<Exclude<keyof I["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["proposer"]["pubKey"], keyof import("../crypto/keys").PublicKey>, never>) | undefined;
                        votingPower?: bigint | undefined;
                        proposerPriority?: bigint | undefined;
                    } & Record<Exclude<keyof I["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["proposer"], keyof Validator>, never>) | undefined;
                    totalVotingPower?: bigint | undefined;
                } & Record<Exclude<keyof I["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"], keyof import("./validator").ValidatorSet>, never>) | undefined;
            } & Record<Exclude<keyof I["lightClientAttackEvidence"]["conflictingBlock"], keyof LightBlock>, never>) | undefined;
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
                } & Record<Exclude<keyof I["lightClientAttackEvidence"]["byzantineValidators"][number]["pubKey"], keyof import("../crypto/keys").PublicKey>, never>) | undefined;
                votingPower?: bigint | undefined;
                proposerPriority?: bigint | undefined;
            } & Record<Exclude<keyof I["lightClientAttackEvidence"]["byzantineValidators"][number], keyof Validator>, never>)[] & Record<Exclude<keyof I["lightClientAttackEvidence"]["byzantineValidators"], keyof {
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
            } & Record<Exclude<keyof I["lightClientAttackEvidence"]["timestamp"], keyof Timestamp>, never>) | undefined;
        } & Record<Exclude<keyof I["lightClientAttackEvidence"], keyof LightClientAttackEvidence>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof Evidence>, never>>(object: I): Evidence;
};
export declare const DuplicateVoteEvidence: {
    typeUrl: string;
    encode(message: DuplicateVoteEvidence, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): DuplicateVoteEvidence;
    fromJSON(object: any): DuplicateVoteEvidence;
    toJSON(message: DuplicateVoteEvidence): unknown;
    fromPartial<I extends {
        voteA?: {
            type?: import("./types").SignedMsgType | undefined;
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
            type?: import("./types").SignedMsgType | undefined;
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
            type?: import("./types").SignedMsgType | undefined;
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
            type?: import("./types").SignedMsgType | undefined;
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
                } & Record<Exclude<keyof I["voteA"]["blockId"]["partSetHeader"], keyof import("./types").PartSetHeader>, never>) | undefined;
            } & Record<Exclude<keyof I["voteA"]["blockId"], keyof import("./types").BlockID>, never>) | undefined;
            timestamp?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["voteA"]["timestamp"], keyof Timestamp>, never>) | undefined;
            validatorAddress?: Uint8Array | undefined;
            validatorIndex?: number | undefined;
            signature?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["voteA"], keyof Vote>, never>) | undefined;
        voteB?: ({
            type?: import("./types").SignedMsgType | undefined;
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
            type?: import("./types").SignedMsgType | undefined;
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
                } & Record<Exclude<keyof I["voteB"]["blockId"]["partSetHeader"], keyof import("./types").PartSetHeader>, never>) | undefined;
            } & Record<Exclude<keyof I["voteB"]["blockId"], keyof import("./types").BlockID>, never>) | undefined;
            timestamp?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["voteB"]["timestamp"], keyof Timestamp>, never>) | undefined;
            validatorAddress?: Uint8Array | undefined;
            validatorIndex?: number | undefined;
            signature?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["voteB"], keyof Vote>, never>) | undefined;
        totalVotingPower?: bigint | undefined;
        validatorPower?: bigint | undefined;
        timestamp?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["timestamp"], keyof Timestamp>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof DuplicateVoteEvidence>, never>>(object: I): DuplicateVoteEvidence;
};
export declare const LightClientAttackEvidence: {
    typeUrl: string;
    encode(message: LightClientAttackEvidence, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): LightClientAttackEvidence;
    fromJSON(object: any): LightClientAttackEvidence;
    toJSON(message: LightClientAttackEvidence): unknown;
    fromPartial<I extends {
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
                        blockIdFlag?: import("./types").BlockIDFlag | undefined;
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
                        blockIdFlag?: import("./types").BlockIDFlag | undefined;
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
                        blockIdFlag?: import("./types").BlockIDFlag | undefined;
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
                    } & Record<Exclude<keyof I["conflictingBlock"]["signedHeader"]["header"]["version"], keyof import("../version/types").Consensus>, never>) | undefined;
                    chainId?: string | undefined;
                    height?: bigint | undefined;
                    time?: ({
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } & {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } & Record<Exclude<keyof I["conflictingBlock"]["signedHeader"]["header"]["time"], keyof Timestamp>, never>) | undefined;
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
                        } & Record<Exclude<keyof I["conflictingBlock"]["signedHeader"]["header"]["lastBlockId"]["partSetHeader"], keyof import("./types").PartSetHeader>, never>) | undefined;
                    } & Record<Exclude<keyof I["conflictingBlock"]["signedHeader"]["header"]["lastBlockId"], keyof import("./types").BlockID>, never>) | undefined;
                    lastCommitHash?: Uint8Array | undefined;
                    dataHash?: Uint8Array | undefined;
                    validatorsHash?: Uint8Array | undefined;
                    nextValidatorsHash?: Uint8Array | undefined;
                    consensusHash?: Uint8Array | undefined;
                    appHash?: Uint8Array | undefined;
                    lastResultsHash?: Uint8Array | undefined;
                    evidenceHash?: Uint8Array | undefined;
                    proposerAddress?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["conflictingBlock"]["signedHeader"]["header"], keyof import("./types").Header>, never>) | undefined;
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
                        blockIdFlag?: import("./types").BlockIDFlag | undefined;
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
                        } & Record<Exclude<keyof I["conflictingBlock"]["signedHeader"]["commit"]["blockId"]["partSetHeader"], keyof import("./types").PartSetHeader>, never>) | undefined;
                    } & Record<Exclude<keyof I["conflictingBlock"]["signedHeader"]["commit"]["blockId"], keyof import("./types").BlockID>, never>) | undefined;
                    signatures?: ({
                        blockIdFlag?: import("./types").BlockIDFlag | undefined;
                        validatorAddress?: Uint8Array | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                        signature?: Uint8Array | undefined;
                    }[] & ({
                        blockIdFlag?: import("./types").BlockIDFlag | undefined;
                        validatorAddress?: Uint8Array | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                        signature?: Uint8Array | undefined;
                    } & {
                        blockIdFlag?: import("./types").BlockIDFlag | undefined;
                        validatorAddress?: Uint8Array | undefined;
                        timestamp?: ({
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } & {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } & Record<Exclude<keyof I["conflictingBlock"]["signedHeader"]["commit"]["signatures"][number]["timestamp"], keyof Timestamp>, never>) | undefined;
                        signature?: Uint8Array | undefined;
                    } & Record<Exclude<keyof I["conflictingBlock"]["signedHeader"]["commit"]["signatures"][number], keyof import("./types").CommitSig>, never>)[] & Record<Exclude<keyof I["conflictingBlock"]["signedHeader"]["commit"]["signatures"], keyof {
                        blockIdFlag?: import("./types").BlockIDFlag | undefined;
                        validatorAddress?: Uint8Array | undefined;
                        timestamp?: {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } | undefined;
                        signature?: Uint8Array | undefined;
                    }[]>, never>) | undefined;
                } & Record<Exclude<keyof I["conflictingBlock"]["signedHeader"]["commit"], keyof import("./types").Commit>, never>) | undefined;
            } & Record<Exclude<keyof I["conflictingBlock"]["signedHeader"], keyof import("./types").SignedHeader>, never>) | undefined;
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
                    } & Record<Exclude<keyof I["conflictingBlock"]["validatorSet"]["validators"][number]["pubKey"], keyof import("../crypto/keys").PublicKey>, never>) | undefined;
                    votingPower?: bigint | undefined;
                    proposerPriority?: bigint | undefined;
                } & Record<Exclude<keyof I["conflictingBlock"]["validatorSet"]["validators"][number], keyof Validator>, never>)[] & Record<Exclude<keyof I["conflictingBlock"]["validatorSet"]["validators"], keyof {
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
                    } & Record<Exclude<keyof I["conflictingBlock"]["validatorSet"]["proposer"]["pubKey"], keyof import("../crypto/keys").PublicKey>, never>) | undefined;
                    votingPower?: bigint | undefined;
                    proposerPriority?: bigint | undefined;
                } & Record<Exclude<keyof I["conflictingBlock"]["validatorSet"]["proposer"], keyof Validator>, never>) | undefined;
                totalVotingPower?: bigint | undefined;
            } & Record<Exclude<keyof I["conflictingBlock"]["validatorSet"], keyof import("./validator").ValidatorSet>, never>) | undefined;
        } & Record<Exclude<keyof I["conflictingBlock"], keyof LightBlock>, never>) | undefined;
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
            } & Record<Exclude<keyof I["byzantineValidators"][number]["pubKey"], keyof import("../crypto/keys").PublicKey>, never>) | undefined;
            votingPower?: bigint | undefined;
            proposerPriority?: bigint | undefined;
        } & Record<Exclude<keyof I["byzantineValidators"][number], keyof Validator>, never>)[] & Record<Exclude<keyof I["byzantineValidators"], keyof {
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
        } & Record<Exclude<keyof I["timestamp"], keyof Timestamp>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof LightClientAttackEvidence>, never>>(object: I): LightClientAttackEvidence;
};
export declare const EvidenceList: {
    typeUrl: string;
    encode(message: EvidenceList, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): EvidenceList;
    fromJSON(object: any): EvidenceList;
    toJSON(message: EvidenceList): unknown;
    fromPartial<I extends {
        evidence?: {
            duplicateVoteEvidence?: {
                voteA?: {
                    type?: import("./types").SignedMsgType | undefined;
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
                    type?: import("./types").SignedMsgType | undefined;
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
                                blockIdFlag?: import("./types").BlockIDFlag | undefined;
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
                    type?: import("./types").SignedMsgType | undefined;
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
                    type?: import("./types").SignedMsgType | undefined;
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
                                blockIdFlag?: import("./types").BlockIDFlag | undefined;
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
                    type?: import("./types").SignedMsgType | undefined;
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
                    type?: import("./types").SignedMsgType | undefined;
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
                                blockIdFlag?: import("./types").BlockIDFlag | undefined;
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
                    type?: import("./types").SignedMsgType | undefined;
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
                    type?: import("./types").SignedMsgType | undefined;
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
                    type?: import("./types").SignedMsgType | undefined;
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
                    type?: import("./types").SignedMsgType | undefined;
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
                        } & Record<Exclude<keyof I["evidence"][number]["duplicateVoteEvidence"]["voteA"]["blockId"]["partSetHeader"], keyof import("./types").PartSetHeader>, never>) | undefined;
                    } & Record<Exclude<keyof I["evidence"][number]["duplicateVoteEvidence"]["voteA"]["blockId"], keyof import("./types").BlockID>, never>) | undefined;
                    timestamp?: ({
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } & {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } & Record<Exclude<keyof I["evidence"][number]["duplicateVoteEvidence"]["voteA"]["timestamp"], keyof Timestamp>, never>) | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    validatorIndex?: number | undefined;
                    signature?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["evidence"][number]["duplicateVoteEvidence"]["voteA"], keyof Vote>, never>) | undefined;
                voteB?: ({
                    type?: import("./types").SignedMsgType | undefined;
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
                    type?: import("./types").SignedMsgType | undefined;
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
                        } & Record<Exclude<keyof I["evidence"][number]["duplicateVoteEvidence"]["voteB"]["blockId"]["partSetHeader"], keyof import("./types").PartSetHeader>, never>) | undefined;
                    } & Record<Exclude<keyof I["evidence"][number]["duplicateVoteEvidence"]["voteB"]["blockId"], keyof import("./types").BlockID>, never>) | undefined;
                    timestamp?: ({
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } & {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } & Record<Exclude<keyof I["evidence"][number]["duplicateVoteEvidence"]["voteB"]["timestamp"], keyof Timestamp>, never>) | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    validatorIndex?: number | undefined;
                    signature?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["evidence"][number]["duplicateVoteEvidence"]["voteB"], keyof Vote>, never>) | undefined;
                totalVotingPower?: bigint | undefined;
                validatorPower?: bigint | undefined;
                timestamp?: ({
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & Record<Exclude<keyof I["evidence"][number]["duplicateVoteEvidence"]["timestamp"], keyof Timestamp>, never>) | undefined;
            } & Record<Exclude<keyof I["evidence"][number]["duplicateVoteEvidence"], keyof DuplicateVoteEvidence>, never>) | undefined;
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
                                blockIdFlag?: import("./types").BlockIDFlag | undefined;
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
                                blockIdFlag?: import("./types").BlockIDFlag | undefined;
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
                                blockIdFlag?: import("./types").BlockIDFlag | undefined;
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
                            } & Record<Exclude<keyof I["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"]["version"], keyof import("../version/types").Consensus>, never>) | undefined;
                            chainId?: string | undefined;
                            height?: bigint | undefined;
                            time?: ({
                                seconds?: bigint | undefined;
                                nanos?: number | undefined;
                            } & {
                                seconds?: bigint | undefined;
                                nanos?: number | undefined;
                            } & Record<Exclude<keyof I["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"]["time"], keyof Timestamp>, never>) | undefined;
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
                                } & Record<Exclude<keyof I["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"]["lastBlockId"]["partSetHeader"], keyof import("./types").PartSetHeader>, never>) | undefined;
                            } & Record<Exclude<keyof I["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"]["lastBlockId"], keyof import("./types").BlockID>, never>) | undefined;
                            lastCommitHash?: Uint8Array | undefined;
                            dataHash?: Uint8Array | undefined;
                            validatorsHash?: Uint8Array | undefined;
                            nextValidatorsHash?: Uint8Array | undefined;
                            consensusHash?: Uint8Array | undefined;
                            appHash?: Uint8Array | undefined;
                            lastResultsHash?: Uint8Array | undefined;
                            evidenceHash?: Uint8Array | undefined;
                            proposerAddress?: Uint8Array | undefined;
                        } & Record<Exclude<keyof I["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"], keyof import("./types").Header>, never>) | undefined;
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
                                blockIdFlag?: import("./types").BlockIDFlag | undefined;
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
                                } & Record<Exclude<keyof I["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["blockId"]["partSetHeader"], keyof import("./types").PartSetHeader>, never>) | undefined;
                            } & Record<Exclude<keyof I["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["blockId"], keyof import("./types").BlockID>, never>) | undefined;
                            signatures?: ({
                                blockIdFlag?: import("./types").BlockIDFlag | undefined;
                                validatorAddress?: Uint8Array | undefined;
                                timestamp?: {
                                    seconds?: bigint | undefined;
                                    nanos?: number | undefined;
                                } | undefined;
                                signature?: Uint8Array | undefined;
                            }[] & ({
                                blockIdFlag?: import("./types").BlockIDFlag | undefined;
                                validatorAddress?: Uint8Array | undefined;
                                timestamp?: {
                                    seconds?: bigint | undefined;
                                    nanos?: number | undefined;
                                } | undefined;
                                signature?: Uint8Array | undefined;
                            } & {
                                blockIdFlag?: import("./types").BlockIDFlag | undefined;
                                validatorAddress?: Uint8Array | undefined;
                                timestamp?: ({
                                    seconds?: bigint | undefined;
                                    nanos?: number | undefined;
                                } & {
                                    seconds?: bigint | undefined;
                                    nanos?: number | undefined;
                                } & Record<Exclude<keyof I["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["signatures"][number]["timestamp"], keyof Timestamp>, never>) | undefined;
                                signature?: Uint8Array | undefined;
                            } & Record<Exclude<keyof I["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["signatures"][number], keyof import("./types").CommitSig>, never>)[] & Record<Exclude<keyof I["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["signatures"], keyof {
                                blockIdFlag?: import("./types").BlockIDFlag | undefined;
                                validatorAddress?: Uint8Array | undefined;
                                timestamp?: {
                                    seconds?: bigint | undefined;
                                    nanos?: number | undefined;
                                } | undefined;
                                signature?: Uint8Array | undefined;
                            }[]>, never>) | undefined;
                        } & Record<Exclude<keyof I["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"], keyof import("./types").Commit>, never>) | undefined;
                    } & Record<Exclude<keyof I["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"], keyof import("./types").SignedHeader>, never>) | undefined;
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
                            } & Record<Exclude<keyof I["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["validators"][number]["pubKey"], keyof import("../crypto/keys").PublicKey>, never>) | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        } & Record<Exclude<keyof I["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["validators"][number], keyof Validator>, never>)[] & Record<Exclude<keyof I["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["validators"], keyof {
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
                            } & Record<Exclude<keyof I["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["proposer"]["pubKey"], keyof import("../crypto/keys").PublicKey>, never>) | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        } & Record<Exclude<keyof I["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["proposer"], keyof Validator>, never>) | undefined;
                        totalVotingPower?: bigint | undefined;
                    } & Record<Exclude<keyof I["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"], keyof import("./validator").ValidatorSet>, never>) | undefined;
                } & Record<Exclude<keyof I["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"], keyof LightBlock>, never>) | undefined;
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
                    } & Record<Exclude<keyof I["evidence"][number]["lightClientAttackEvidence"]["byzantineValidators"][number]["pubKey"], keyof import("../crypto/keys").PublicKey>, never>) | undefined;
                    votingPower?: bigint | undefined;
                    proposerPriority?: bigint | undefined;
                } & Record<Exclude<keyof I["evidence"][number]["lightClientAttackEvidence"]["byzantineValidators"][number], keyof Validator>, never>)[] & Record<Exclude<keyof I["evidence"][number]["lightClientAttackEvidence"]["byzantineValidators"], keyof {
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
                } & Record<Exclude<keyof I["evidence"][number]["lightClientAttackEvidence"]["timestamp"], keyof Timestamp>, never>) | undefined;
            } & Record<Exclude<keyof I["evidence"][number]["lightClientAttackEvidence"], keyof LightClientAttackEvidence>, never>) | undefined;
        } & Record<Exclude<keyof I["evidence"][number], keyof Evidence>, never>)[] & Record<Exclude<keyof I["evidence"], keyof {
            duplicateVoteEvidence?: {
                voteA?: {
                    type?: import("./types").SignedMsgType | undefined;
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
                    type?: import("./types").SignedMsgType | undefined;
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
                                blockIdFlag?: import("./types").BlockIDFlag | undefined;
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
    } & Record<Exclude<keyof I, "evidence">, never>>(object: I): EvidenceList;
};
