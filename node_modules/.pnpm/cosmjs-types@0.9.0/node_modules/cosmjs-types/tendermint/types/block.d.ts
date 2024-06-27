import { Header, Data, Commit } from "./types";
import { EvidenceList } from "./evidence";
import { BinaryReader, BinaryWriter } from "../../binary";
export declare const protobufPackage = "tendermint.types";
export interface Block {
    header: Header;
    data: Data;
    evidence: EvidenceList;
    lastCommit?: Commit;
}
export declare const Block: {
    typeUrl: string;
    encode(message: Block, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Block;
    fromJSON(object: any): Block;
    toJSON(message: Block): unknown;
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
        data?: {
            txs?: Uint8Array[] | undefined;
        } | undefined;
        evidence?: {
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
            } & Record<Exclude<keyof I["header"]["version"], keyof import("../version/types").Consensus>, never>) | undefined;
            chainId?: string | undefined;
            height?: bigint | undefined;
            time?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["header"]["time"], keyof import("../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
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
                } & Record<Exclude<keyof I["header"]["lastBlockId"]["partSetHeader"], keyof import("./types").PartSetHeader>, never>) | undefined;
            } & Record<Exclude<keyof I["header"]["lastBlockId"], keyof import("./types").BlockID>, never>) | undefined;
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
        data?: ({
            txs?: Uint8Array[] | undefined;
        } & {
            txs?: (Uint8Array[] & Uint8Array[] & Record<Exclude<keyof I["data"]["txs"], keyof Uint8Array[]>, never>) | undefined;
        } & Record<Exclude<keyof I["data"], "txs">, never>) | undefined;
        evidence?: ({
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
                            } & Record<Exclude<keyof I["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteA"]["blockId"]["partSetHeader"], keyof import("./types").PartSetHeader>, never>) | undefined;
                        } & Record<Exclude<keyof I["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteA"]["blockId"], keyof import("./types").BlockID>, never>) | undefined;
                        timestamp?: ({
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } & {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } & Record<Exclude<keyof I["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteA"]["timestamp"], keyof import("../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                        validatorAddress?: Uint8Array | undefined;
                        validatorIndex?: number | undefined;
                        signature?: Uint8Array | undefined;
                    } & Record<Exclude<keyof I["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteA"], keyof import("./types").Vote>, never>) | undefined;
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
                            } & Record<Exclude<keyof I["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteB"]["blockId"]["partSetHeader"], keyof import("./types").PartSetHeader>, never>) | undefined;
                        } & Record<Exclude<keyof I["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteB"]["blockId"], keyof import("./types").BlockID>, never>) | undefined;
                        timestamp?: ({
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } & {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } & Record<Exclude<keyof I["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteB"]["timestamp"], keyof import("../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                        validatorAddress?: Uint8Array | undefined;
                        validatorIndex?: number | undefined;
                        signature?: Uint8Array | undefined;
                    } & Record<Exclude<keyof I["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteB"], keyof import("./types").Vote>, never>) | undefined;
                    totalVotingPower?: bigint | undefined;
                    validatorPower?: bigint | undefined;
                    timestamp?: ({
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } & {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } & Record<Exclude<keyof I["evidence"]["evidence"][number]["duplicateVoteEvidence"]["timestamp"], keyof import("../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                } & Record<Exclude<keyof I["evidence"]["evidence"][number]["duplicateVoteEvidence"], keyof import("./evidence").DuplicateVoteEvidence>, never>) | undefined;
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
                                } & Record<Exclude<keyof I["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"]["version"], keyof import("../version/types").Consensus>, never>) | undefined;
                                chainId?: string | undefined;
                                height?: bigint | undefined;
                                time?: ({
                                    seconds?: bigint | undefined;
                                    nanos?: number | undefined;
                                } & {
                                    seconds?: bigint | undefined;
                                    nanos?: number | undefined;
                                } & Record<Exclude<keyof I["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"]["time"], keyof import("../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
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
                                    } & Record<Exclude<keyof I["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"]["lastBlockId"]["partSetHeader"], keyof import("./types").PartSetHeader>, never>) | undefined;
                                } & Record<Exclude<keyof I["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"]["lastBlockId"], keyof import("./types").BlockID>, never>) | undefined;
                                lastCommitHash?: Uint8Array | undefined;
                                dataHash?: Uint8Array | undefined;
                                validatorsHash?: Uint8Array | undefined;
                                nextValidatorsHash?: Uint8Array | undefined;
                                consensusHash?: Uint8Array | undefined;
                                appHash?: Uint8Array | undefined;
                                lastResultsHash?: Uint8Array | undefined;
                                evidenceHash?: Uint8Array | undefined;
                                proposerAddress?: Uint8Array | undefined;
                            } & Record<Exclude<keyof I["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"], keyof Header>, never>) | undefined;
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
                                    } & Record<Exclude<keyof I["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["blockId"]["partSetHeader"], keyof import("./types").PartSetHeader>, never>) | undefined;
                                } & Record<Exclude<keyof I["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["blockId"], keyof import("./types").BlockID>, never>) | undefined;
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
                                    } & Record<Exclude<keyof I["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["signatures"][number]["timestamp"], keyof import("../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                                    signature?: Uint8Array | undefined;
                                } & Record<Exclude<keyof I["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["signatures"][number], keyof import("./types").CommitSig>, never>)[] & Record<Exclude<keyof I["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["signatures"], keyof {
                                    blockIdFlag?: import("./types").BlockIDFlag | undefined;
                                    validatorAddress?: Uint8Array | undefined;
                                    timestamp?: {
                                        seconds?: bigint | undefined;
                                        nanos?: number | undefined;
                                    } | undefined;
                                    signature?: Uint8Array | undefined;
                                }[]>, never>) | undefined;
                            } & Record<Exclude<keyof I["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"], keyof Commit>, never>) | undefined;
                        } & Record<Exclude<keyof I["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"], keyof import("./types").SignedHeader>, never>) | undefined;
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
                                } & Record<Exclude<keyof I["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["validators"][number]["pubKey"], keyof import("../crypto/keys").PublicKey>, never>) | undefined;
                                votingPower?: bigint | undefined;
                                proposerPriority?: bigint | undefined;
                            } & Record<Exclude<keyof I["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["validators"][number], keyof import("./validator").Validator>, never>)[] & Record<Exclude<keyof I["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["validators"], keyof {
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
                                } & Record<Exclude<keyof I["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["proposer"]["pubKey"], keyof import("../crypto/keys").PublicKey>, never>) | undefined;
                                votingPower?: bigint | undefined;
                                proposerPriority?: bigint | undefined;
                            } & Record<Exclude<keyof I["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["proposer"], keyof import("./validator").Validator>, never>) | undefined;
                            totalVotingPower?: bigint | undefined;
                        } & Record<Exclude<keyof I["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"], keyof import("./validator").ValidatorSet>, never>) | undefined;
                    } & Record<Exclude<keyof I["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"], keyof import("./types").LightBlock>, never>) | undefined;
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
                        } & Record<Exclude<keyof I["evidence"]["evidence"][number]["lightClientAttackEvidence"]["byzantineValidators"][number]["pubKey"], keyof import("../crypto/keys").PublicKey>, never>) | undefined;
                        votingPower?: bigint | undefined;
                        proposerPriority?: bigint | undefined;
                    } & Record<Exclude<keyof I["evidence"]["evidence"][number]["lightClientAttackEvidence"]["byzantineValidators"][number], keyof import("./validator").Validator>, never>)[] & Record<Exclude<keyof I["evidence"]["evidence"][number]["lightClientAttackEvidence"]["byzantineValidators"], keyof {
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
                    } & Record<Exclude<keyof I["evidence"]["evidence"][number]["lightClientAttackEvidence"]["timestamp"], keyof import("../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                } & Record<Exclude<keyof I["evidence"]["evidence"][number]["lightClientAttackEvidence"], keyof import("./evidence").LightClientAttackEvidence>, never>) | undefined;
            } & Record<Exclude<keyof I["evidence"]["evidence"][number], keyof import("./evidence").Evidence>, never>)[] & Record<Exclude<keyof I["evidence"]["evidence"], keyof {
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
        } & Record<Exclude<keyof I["evidence"], "evidence">, never>) | undefined;
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
                } & Record<Exclude<keyof I["lastCommit"]["blockId"]["partSetHeader"], keyof import("./types").PartSetHeader>, never>) | undefined;
            } & Record<Exclude<keyof I["lastCommit"]["blockId"], keyof import("./types").BlockID>, never>) | undefined;
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
                } & Record<Exclude<keyof I["lastCommit"]["signatures"][number]["timestamp"], keyof import("../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                signature?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["lastCommit"]["signatures"][number], keyof import("./types").CommitSig>, never>)[] & Record<Exclude<keyof I["lastCommit"]["signatures"], keyof {
                blockIdFlag?: import("./types").BlockIDFlag | undefined;
                validatorAddress?: Uint8Array | undefined;
                timestamp?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                signature?: Uint8Array | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["lastCommit"], keyof Commit>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof Block>, never>>(object: I): Block;
};
