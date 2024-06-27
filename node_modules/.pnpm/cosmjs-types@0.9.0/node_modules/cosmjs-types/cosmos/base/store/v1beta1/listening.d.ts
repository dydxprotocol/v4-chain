import { RequestDeliverTx, ResponseDeliverTx, RequestBeginBlock, ResponseBeginBlock, RequestEndBlock, ResponseEndBlock, ResponseCommit } from "../../../../tendermint/abci/types";
import { BinaryReader, BinaryWriter } from "../../../../binary";
export declare const protobufPackage = "cosmos.base.store.v1beta1";
/**
 * StoreKVPair is a KVStore KVPair used for listening to state changes (Sets and Deletes)
 * It optionally includes the StoreKey for the originating KVStore and a Boolean flag to distinguish between Sets and
 * Deletes
 *
 * Since: cosmos-sdk 0.43
 */
export interface StoreKVPair {
    /** the store key for the KVStore this pair originates from */
    storeKey: string;
    /** true indicates a delete operation, false indicates a set operation */
    delete: boolean;
    key: Uint8Array;
    value: Uint8Array;
}
/**
 * BlockMetadata contains all the abci event data of a block
 * the file streamer dump them into files together with the state changes.
 */
export interface BlockMetadata {
    requestBeginBlock?: RequestBeginBlock;
    responseBeginBlock?: ResponseBeginBlock;
    deliverTxs: BlockMetadata_DeliverTx[];
    requestEndBlock?: RequestEndBlock;
    responseEndBlock?: ResponseEndBlock;
    responseCommit?: ResponseCommit;
}
/** DeliverTx encapulate deliver tx request and response. */
export interface BlockMetadata_DeliverTx {
    request?: RequestDeliverTx;
    response?: ResponseDeliverTx;
}
export declare const StoreKVPair: {
    typeUrl: string;
    encode(message: StoreKVPair, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): StoreKVPair;
    fromJSON(object: any): StoreKVPair;
    toJSON(message: StoreKVPair): unknown;
    fromPartial<I extends {
        storeKey?: string | undefined;
        delete?: boolean | undefined;
        key?: Uint8Array | undefined;
        value?: Uint8Array | undefined;
    } & {
        storeKey?: string | undefined;
        delete?: boolean | undefined;
        key?: Uint8Array | undefined;
        value?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof StoreKVPair>, never>>(object: I): StoreKVPair;
};
export declare const BlockMetadata: {
    typeUrl: string;
    encode(message: BlockMetadata, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): BlockMetadata;
    fromJSON(object: any): BlockMetadata;
    toJSON(message: BlockMetadata): unknown;
    fromPartial<I extends {
        requestBeginBlock?: {
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
                type?: import("../../../../tendermint/abci/types").MisbehaviorType | undefined;
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
        responseBeginBlock?: {
            events?: {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[] | undefined;
        } | undefined;
        deliverTxs?: {
            request?: {
                tx?: Uint8Array | undefined;
            } | undefined;
            response?: {
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
        }[] | undefined;
        requestEndBlock?: {
            height?: bigint | undefined;
        } | undefined;
        responseEndBlock?: {
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
        responseCommit?: {
            data?: Uint8Array | undefined;
            retainHeight?: bigint | undefined;
        } | undefined;
    } & {
        requestBeginBlock?: ({
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
                type?: import("../../../../tendermint/abci/types").MisbehaviorType | undefined;
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
                } & Record<Exclude<keyof I["requestBeginBlock"]["header"]["version"], keyof import("../../../../tendermint/version/types").Consensus>, never>) | undefined;
                chainId?: string | undefined;
                height?: bigint | undefined;
                time?: ({
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & Record<Exclude<keyof I["requestBeginBlock"]["header"]["time"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
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
                    } & Record<Exclude<keyof I["requestBeginBlock"]["header"]["lastBlockId"]["partSetHeader"], keyof import("../../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                } & Record<Exclude<keyof I["requestBeginBlock"]["header"]["lastBlockId"], keyof import("../../../../tendermint/types/types").BlockID>, never>) | undefined;
                lastCommitHash?: Uint8Array | undefined;
                dataHash?: Uint8Array | undefined;
                validatorsHash?: Uint8Array | undefined;
                nextValidatorsHash?: Uint8Array | undefined;
                consensusHash?: Uint8Array | undefined;
                appHash?: Uint8Array | undefined;
                lastResultsHash?: Uint8Array | undefined;
                evidenceHash?: Uint8Array | undefined;
                proposerAddress?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["requestBeginBlock"]["header"], keyof import("../../../../tendermint/types/types").Header>, never>) | undefined;
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
                    } & Record<Exclude<keyof I["requestBeginBlock"]["lastCommitInfo"]["votes"][number]["validator"], keyof import("../../../../tendermint/abci/types").Validator>, never>) | undefined;
                    signedLastBlock?: boolean | undefined;
                } & Record<Exclude<keyof I["requestBeginBlock"]["lastCommitInfo"]["votes"][number], keyof import("../../../../tendermint/abci/types").VoteInfo>, never>)[] & Record<Exclude<keyof I["requestBeginBlock"]["lastCommitInfo"]["votes"], keyof {
                    validator?: {
                        address?: Uint8Array | undefined;
                        power?: bigint | undefined;
                    } | undefined;
                    signedLastBlock?: boolean | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["requestBeginBlock"]["lastCommitInfo"], keyof import("../../../../tendermint/abci/types").CommitInfo>, never>) | undefined;
            byzantineValidators?: ({
                type?: import("../../../../tendermint/abci/types").MisbehaviorType | undefined;
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
                type?: import("../../../../tendermint/abci/types").MisbehaviorType | undefined;
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
                type?: import("../../../../tendermint/abci/types").MisbehaviorType | undefined;
                validator?: ({
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } & {
                    address?: Uint8Array | undefined;
                    power?: bigint | undefined;
                } & Record<Exclude<keyof I["requestBeginBlock"]["byzantineValidators"][number]["validator"], keyof import("../../../../tendermint/abci/types").Validator>, never>) | undefined;
                height?: bigint | undefined;
                time?: ({
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & Record<Exclude<keyof I["requestBeginBlock"]["byzantineValidators"][number]["time"], keyof import("../../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                totalVotingPower?: bigint | undefined;
            } & Record<Exclude<keyof I["requestBeginBlock"]["byzantineValidators"][number], keyof import("../../../../tendermint/abci/types").Misbehavior>, never>)[] & Record<Exclude<keyof I["requestBeginBlock"]["byzantineValidators"], keyof {
                type?: import("../../../../tendermint/abci/types").MisbehaviorType | undefined;
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
        } & Record<Exclude<keyof I["requestBeginBlock"], keyof RequestBeginBlock>, never>) | undefined;
        responseBeginBlock?: ({
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
                } & Record<Exclude<keyof I["responseBeginBlock"]["events"][number]["attributes"][number], keyof import("../../../../tendermint/abci/types").EventAttribute>, never>)[] & Record<Exclude<keyof I["responseBeginBlock"]["events"][number]["attributes"], keyof {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["responseBeginBlock"]["events"][number], keyof import("../../../../tendermint/abci/types").Event>, never>)[] & Record<Exclude<keyof I["responseBeginBlock"]["events"], keyof {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["responseBeginBlock"], "events">, never>) | undefined;
        deliverTxs?: ({
            request?: {
                tx?: Uint8Array | undefined;
            } | undefined;
            response?: {
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
        }[] & ({
            request?: {
                tx?: Uint8Array | undefined;
            } | undefined;
            response?: {
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
            request?: ({
                tx?: Uint8Array | undefined;
            } & {
                tx?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["deliverTxs"][number]["request"], "tx">, never>) | undefined;
            response?: ({
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
                    } & Record<Exclude<keyof I["deliverTxs"][number]["response"]["events"][number]["attributes"][number], keyof import("../../../../tendermint/abci/types").EventAttribute>, never>)[] & Record<Exclude<keyof I["deliverTxs"][number]["response"]["events"][number]["attributes"], keyof {
                        key?: string | undefined;
                        value?: string | undefined;
                        index?: boolean | undefined;
                    }[]>, never>) | undefined;
                } & Record<Exclude<keyof I["deliverTxs"][number]["response"]["events"][number], keyof import("../../../../tendermint/abci/types").Event>, never>)[] & Record<Exclude<keyof I["deliverTxs"][number]["response"]["events"], keyof {
                    type?: string | undefined;
                    attributes?: {
                        key?: string | undefined;
                        value?: string | undefined;
                        index?: boolean | undefined;
                    }[] | undefined;
                }[]>, never>) | undefined;
                codespace?: string | undefined;
            } & Record<Exclude<keyof I["deliverTxs"][number]["response"], keyof ResponseDeliverTx>, never>) | undefined;
        } & Record<Exclude<keyof I["deliverTxs"][number], keyof BlockMetadata_DeliverTx>, never>)[] & Record<Exclude<keyof I["deliverTxs"], keyof {
            request?: {
                tx?: Uint8Array | undefined;
            } | undefined;
            response?: {
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
        }[]>, never>) | undefined;
        requestEndBlock?: ({
            height?: bigint | undefined;
        } & {
            height?: bigint | undefined;
        } & Record<Exclude<keyof I["requestEndBlock"], "height">, never>) | undefined;
        responseEndBlock?: ({
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
                } & Record<Exclude<keyof I["responseEndBlock"]["validatorUpdates"][number]["pubKey"], keyof import("../../../../tendermint/crypto/keys").PublicKey>, never>) | undefined;
                power?: bigint | undefined;
            } & Record<Exclude<keyof I["responseEndBlock"]["validatorUpdates"][number], keyof import("../../../../tendermint/abci/types").ValidatorUpdate>, never>)[] & Record<Exclude<keyof I["responseEndBlock"]["validatorUpdates"], keyof {
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
                } & Record<Exclude<keyof I["responseEndBlock"]["consensusParamUpdates"]["block"], keyof import("../../../../tendermint/types/params").BlockParams>, never>) | undefined;
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
                    } & Record<Exclude<keyof I["responseEndBlock"]["consensusParamUpdates"]["evidence"]["maxAgeDuration"], keyof import("../../../../google/protobuf/duration").Duration>, never>) | undefined;
                    maxBytes?: bigint | undefined;
                } & Record<Exclude<keyof I["responseEndBlock"]["consensusParamUpdates"]["evidence"], keyof import("../../../../tendermint/types/params").EvidenceParams>, never>) | undefined;
                validator?: ({
                    pubKeyTypes?: string[] | undefined;
                } & {
                    pubKeyTypes?: (string[] & string[] & Record<Exclude<keyof I["responseEndBlock"]["consensusParamUpdates"]["validator"]["pubKeyTypes"], keyof string[]>, never>) | undefined;
                } & Record<Exclude<keyof I["responseEndBlock"]["consensusParamUpdates"]["validator"], "pubKeyTypes">, never>) | undefined;
                version?: ({
                    app?: bigint | undefined;
                } & {
                    app?: bigint | undefined;
                } & Record<Exclude<keyof I["responseEndBlock"]["consensusParamUpdates"]["version"], "app">, never>) | undefined;
            } & Record<Exclude<keyof I["responseEndBlock"]["consensusParamUpdates"], keyof import("../../../../tendermint/types/params").ConsensusParams>, never>) | undefined;
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
                } & Record<Exclude<keyof I["responseEndBlock"]["events"][number]["attributes"][number], keyof import("../../../../tendermint/abci/types").EventAttribute>, never>)[] & Record<Exclude<keyof I["responseEndBlock"]["events"][number]["attributes"], keyof {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["responseEndBlock"]["events"][number], keyof import("../../../../tendermint/abci/types").Event>, never>)[] & Record<Exclude<keyof I["responseEndBlock"]["events"], keyof {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["responseEndBlock"], keyof ResponseEndBlock>, never>) | undefined;
        responseCommit?: ({
            data?: Uint8Array | undefined;
            retainHeight?: bigint | undefined;
        } & {
            data?: Uint8Array | undefined;
            retainHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["responseCommit"], keyof ResponseCommit>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof BlockMetadata>, never>>(object: I): BlockMetadata;
};
export declare const BlockMetadata_DeliverTx: {
    typeUrl: string;
    encode(message: BlockMetadata_DeliverTx, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): BlockMetadata_DeliverTx;
    fromJSON(object: any): BlockMetadata_DeliverTx;
    toJSON(message: BlockMetadata_DeliverTx): unknown;
    fromPartial<I extends {
        request?: {
            tx?: Uint8Array | undefined;
        } | undefined;
        response?: {
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
        request?: ({
            tx?: Uint8Array | undefined;
        } & {
            tx?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["request"], "tx">, never>) | undefined;
        response?: ({
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
                } & Record<Exclude<keyof I["response"]["events"][number]["attributes"][number], keyof import("../../../../tendermint/abci/types").EventAttribute>, never>)[] & Record<Exclude<keyof I["response"]["events"][number]["attributes"], keyof {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["response"]["events"][number], keyof import("../../../../tendermint/abci/types").Event>, never>)[] & Record<Exclude<keyof I["response"]["events"], keyof {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[]>, never>) | undefined;
            codespace?: string | undefined;
        } & Record<Exclude<keyof I["response"], keyof ResponseDeliverTx>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof BlockMetadata_DeliverTx>, never>>(object: I): BlockMetadata_DeliverTx;
};
