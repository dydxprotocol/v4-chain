import { Any } from "../../../../google/protobuf/any";
import { Event } from "../../../../tendermint/abci/types";
import { BinaryReader, BinaryWriter } from "../../../../binary";
export declare const protobufPackage = "cosmos.base.abci.v1beta1";
/**
 * TxResponse defines a structure containing relevant tx data and metadata. The
 * tags are stringified and the log is JSON decoded.
 */
export interface TxResponse {
    /** The block height */
    height: bigint;
    /** The transaction hash. */
    txhash: string;
    /** Namespace for the Code */
    codespace: string;
    /** Response code. */
    code: number;
    /** Result bytes, if any. */
    data: string;
    /**
     * The output of the application's logger (raw string). May be
     * non-deterministic.
     */
    rawLog: string;
    /** The output of the application's logger (typed). May be non-deterministic. */
    logs: ABCIMessageLog[];
    /** Additional information. May be non-deterministic. */
    info: string;
    /** Amount of gas requested for transaction. */
    gasWanted: bigint;
    /** Amount of gas consumed by transaction. */
    gasUsed: bigint;
    /** The request transaction bytes. */
    tx?: Any;
    /**
     * Time of the previous block. For heights > 1, it's the weighted median of
     * the timestamps of the valid votes in the block.LastCommit. For height == 1,
     * it's genesis time.
     */
    timestamp: string;
    /**
     * Events defines all the events emitted by processing a transaction. Note,
     * these events include those emitted by processing all the messages and those
     * emitted from the ante. Whereas Logs contains the events, with
     * additional metadata, emitted only by processing the messages.
     *
     * Since: cosmos-sdk 0.42.11, 0.44.5, 0.45
     */
    events: Event[];
}
/** ABCIMessageLog defines a structure containing an indexed tx ABCI message log. */
export interface ABCIMessageLog {
    msgIndex: number;
    log: string;
    /**
     * Events contains a slice of Event objects that were emitted during some
     * execution.
     */
    events: StringEvent[];
}
/**
 * StringEvent defines en Event object wrapper where all the attributes
 * contain key/value pairs that are strings instead of raw bytes.
 */
export interface StringEvent {
    type: string;
    attributes: Attribute[];
}
/**
 * Attribute defines an attribute wrapper where the key and value are
 * strings instead of raw bytes.
 */
export interface Attribute {
    key: string;
    value: string;
}
/** GasInfo defines tx execution gas context. */
export interface GasInfo {
    /** GasWanted is the maximum units of work we allow this tx to perform. */
    gasWanted: bigint;
    /** GasUsed is the amount of gas actually consumed. */
    gasUsed: bigint;
}
/** Result is the union of ResponseFormat and ResponseCheckTx. */
export interface Result {
    /**
     * Data is any data returned from message or handler execution. It MUST be
     * length prefixed in order to separate data from multiple message executions.
     * Deprecated. This field is still populated, but prefer msg_response instead
     * because it also contains the Msg response typeURL.
     */
    /** @deprecated */
    data: Uint8Array;
    /** Log contains the log information from message or handler execution. */
    log: string;
    /**
     * Events contains a slice of Event objects that were emitted during message
     * or handler execution.
     */
    events: Event[];
    /**
     * msg_responses contains the Msg handler responses type packed in Anys.
     *
     * Since: cosmos-sdk 0.46
     */
    msgResponses: Any[];
}
/**
 * SimulationResponse defines the response generated when a transaction is
 * successfully simulated.
 */
export interface SimulationResponse {
    gasInfo: GasInfo;
    result?: Result;
}
/**
 * MsgData defines the data returned in a Result object during message
 * execution.
 */
/** @deprecated */
export interface MsgData {
    msgType: string;
    data: Uint8Array;
}
/**
 * TxMsgData defines a list of MsgData. A transaction will have a MsgData object
 * for each message.
 */
export interface TxMsgData {
    /** data field is deprecated and not populated. */
    /** @deprecated */
    data: MsgData[];
    /**
     * msg_responses contains the Msg handler responses packed into Anys.
     *
     * Since: cosmos-sdk 0.46
     */
    msgResponses: Any[];
}
/** SearchTxsResult defines a structure for querying txs pageable */
export interface SearchTxsResult {
    /** Count of all txs */
    totalCount: bigint;
    /** Count of txs in current page */
    count: bigint;
    /** Index of current page, start from 1 */
    pageNumber: bigint;
    /** Count of total pages */
    pageTotal: bigint;
    /** Max count txs per page */
    limit: bigint;
    /** List of txs in current page */
    txs: TxResponse[];
}
export declare const TxResponse: {
    typeUrl: string;
    encode(message: TxResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): TxResponse;
    fromJSON(object: any): TxResponse;
    toJSON(message: TxResponse): unknown;
    fromPartial<I extends {
        height?: bigint | undefined;
        txhash?: string | undefined;
        codespace?: string | undefined;
        code?: number | undefined;
        data?: string | undefined;
        rawLog?: string | undefined;
        logs?: {
            msgIndex?: number | undefined;
            log?: string | undefined;
            events?: {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                }[] | undefined;
            }[] | undefined;
        }[] | undefined;
        info?: string | undefined;
        gasWanted?: bigint | undefined;
        gasUsed?: bigint | undefined;
        tx?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
        timestamp?: string | undefined;
        events?: {
            type?: string | undefined;
            attributes?: {
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            }[] | undefined;
        }[] | undefined;
    } & {
        height?: bigint | undefined;
        txhash?: string | undefined;
        codespace?: string | undefined;
        code?: number | undefined;
        data?: string | undefined;
        rawLog?: string | undefined;
        logs?: ({
            msgIndex?: number | undefined;
            log?: string | undefined;
            events?: {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                }[] | undefined;
            }[] | undefined;
        }[] & ({
            msgIndex?: number | undefined;
            log?: string | undefined;
            events?: {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                }[] | undefined;
            }[] | undefined;
        } & {
            msgIndex?: number | undefined;
            log?: string | undefined;
            events?: ({
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                }[] | undefined;
            }[] & ({
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                }[] | undefined;
            } & {
                type?: string | undefined;
                attributes?: ({
                    key?: string | undefined;
                    value?: string | undefined;
                }[] & ({
                    key?: string | undefined;
                    value?: string | undefined;
                } & {
                    key?: string | undefined;
                    value?: string | undefined;
                } & Record<Exclude<keyof I["logs"][number]["events"][number]["attributes"][number], keyof Attribute>, never>)[] & Record<Exclude<keyof I["logs"][number]["events"][number]["attributes"], keyof {
                    key?: string | undefined;
                    value?: string | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["logs"][number]["events"][number], keyof StringEvent>, never>)[] & Record<Exclude<keyof I["logs"][number]["events"], keyof {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                }[] | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["logs"][number], keyof ABCIMessageLog>, never>)[] & Record<Exclude<keyof I["logs"], keyof {
            msgIndex?: number | undefined;
            log?: string | undefined;
            events?: {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                }[] | undefined;
            }[] | undefined;
        }[]>, never>) | undefined;
        info?: string | undefined;
        gasWanted?: bigint | undefined;
        gasUsed?: bigint | undefined;
        tx?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["tx"], keyof Any>, never>) | undefined;
        timestamp?: string | undefined;
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
            } & Record<Exclude<keyof I["events"][number]["attributes"][number], keyof import("../../../../tendermint/abci/types").EventAttribute>, never>)[] & Record<Exclude<keyof I["events"][number]["attributes"], keyof {
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
    } & Record<Exclude<keyof I, keyof TxResponse>, never>>(object: I): TxResponse;
};
export declare const ABCIMessageLog: {
    typeUrl: string;
    encode(message: ABCIMessageLog, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ABCIMessageLog;
    fromJSON(object: any): ABCIMessageLog;
    toJSON(message: ABCIMessageLog): unknown;
    fromPartial<I extends {
        msgIndex?: number | undefined;
        log?: string | undefined;
        events?: {
            type?: string | undefined;
            attributes?: {
                key?: string | undefined;
                value?: string | undefined;
            }[] | undefined;
        }[] | undefined;
    } & {
        msgIndex?: number | undefined;
        log?: string | undefined;
        events?: ({
            type?: string | undefined;
            attributes?: {
                key?: string | undefined;
                value?: string | undefined;
            }[] | undefined;
        }[] & ({
            type?: string | undefined;
            attributes?: {
                key?: string | undefined;
                value?: string | undefined;
            }[] | undefined;
        } & {
            type?: string | undefined;
            attributes?: ({
                key?: string | undefined;
                value?: string | undefined;
            }[] & ({
                key?: string | undefined;
                value?: string | undefined;
            } & {
                key?: string | undefined;
                value?: string | undefined;
            } & Record<Exclude<keyof I["events"][number]["attributes"][number], keyof Attribute>, never>)[] & Record<Exclude<keyof I["events"][number]["attributes"], keyof {
                key?: string | undefined;
                value?: string | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["events"][number], keyof StringEvent>, never>)[] & Record<Exclude<keyof I["events"], keyof {
            type?: string | undefined;
            attributes?: {
                key?: string | undefined;
                value?: string | undefined;
            }[] | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof ABCIMessageLog>, never>>(object: I): ABCIMessageLog;
};
export declare const StringEvent: {
    typeUrl: string;
    encode(message: StringEvent, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): StringEvent;
    fromJSON(object: any): StringEvent;
    toJSON(message: StringEvent): unknown;
    fromPartial<I extends {
        type?: string | undefined;
        attributes?: {
            key?: string | undefined;
            value?: string | undefined;
        }[] | undefined;
    } & {
        type?: string | undefined;
        attributes?: ({
            key?: string | undefined;
            value?: string | undefined;
        }[] & ({
            key?: string | undefined;
            value?: string | undefined;
        } & {
            key?: string | undefined;
            value?: string | undefined;
        } & Record<Exclude<keyof I["attributes"][number], keyof Attribute>, never>)[] & Record<Exclude<keyof I["attributes"], keyof {
            key?: string | undefined;
            value?: string | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof StringEvent>, never>>(object: I): StringEvent;
};
export declare const Attribute: {
    typeUrl: string;
    encode(message: Attribute, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Attribute;
    fromJSON(object: any): Attribute;
    toJSON(message: Attribute): unknown;
    fromPartial<I extends {
        key?: string | undefined;
        value?: string | undefined;
    } & {
        key?: string | undefined;
        value?: string | undefined;
    } & Record<Exclude<keyof I, keyof Attribute>, never>>(object: I): Attribute;
};
export declare const GasInfo: {
    typeUrl: string;
    encode(message: GasInfo, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GasInfo;
    fromJSON(object: any): GasInfo;
    toJSON(message: GasInfo): unknown;
    fromPartial<I extends {
        gasWanted?: bigint | undefined;
        gasUsed?: bigint | undefined;
    } & {
        gasWanted?: bigint | undefined;
        gasUsed?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof GasInfo>, never>>(object: I): GasInfo;
};
export declare const Result: {
    typeUrl: string;
    encode(message: Result, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Result;
    fromJSON(object: any): Result;
    toJSON(message: Result): unknown;
    fromPartial<I extends {
        data?: Uint8Array | undefined;
        log?: string | undefined;
        events?: {
            type?: string | undefined;
            attributes?: {
                key?: string | undefined;
                value?: string | undefined;
                index?: boolean | undefined;
            }[] | undefined;
        }[] | undefined;
        msgResponses?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        }[] | undefined;
    } & {
        data?: Uint8Array | undefined;
        log?: string | undefined;
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
            } & Record<Exclude<keyof I["events"][number]["attributes"][number], keyof import("../../../../tendermint/abci/types").EventAttribute>, never>)[] & Record<Exclude<keyof I["events"][number]["attributes"], keyof {
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
        msgResponses?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        }[] & ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["msgResponses"][number], keyof Any>, never>)[] & Record<Exclude<keyof I["msgResponses"], keyof {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof Result>, never>>(object: I): Result;
};
export declare const SimulationResponse: {
    typeUrl: string;
    encode(message: SimulationResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): SimulationResponse;
    fromJSON(object: any): SimulationResponse;
    toJSON(message: SimulationResponse): unknown;
    fromPartial<I extends {
        gasInfo?: {
            gasWanted?: bigint | undefined;
            gasUsed?: bigint | undefined;
        } | undefined;
        result?: {
            data?: Uint8Array | undefined;
            log?: string | undefined;
            events?: {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[] | undefined;
            msgResponses?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            }[] | undefined;
        } | undefined;
    } & {
        gasInfo?: ({
            gasWanted?: bigint | undefined;
            gasUsed?: bigint | undefined;
        } & {
            gasWanted?: bigint | undefined;
            gasUsed?: bigint | undefined;
        } & Record<Exclude<keyof I["gasInfo"], keyof GasInfo>, never>) | undefined;
        result?: ({
            data?: Uint8Array | undefined;
            log?: string | undefined;
            events?: {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[] | undefined;
            msgResponses?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            }[] | undefined;
        } & {
            data?: Uint8Array | undefined;
            log?: string | undefined;
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
                } & Record<Exclude<keyof I["result"]["events"][number]["attributes"][number], keyof import("../../../../tendermint/abci/types").EventAttribute>, never>)[] & Record<Exclude<keyof I["result"]["events"][number]["attributes"], keyof {
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
            msgResponses?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            }[] & ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["result"]["msgResponses"][number], keyof Any>, never>)[] & Record<Exclude<keyof I["result"]["msgResponses"], keyof {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["result"], keyof Result>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof SimulationResponse>, never>>(object: I): SimulationResponse;
};
export declare const MsgData: {
    typeUrl: string;
    encode(message: MsgData, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgData;
    fromJSON(object: any): MsgData;
    toJSON(message: MsgData): unknown;
    fromPartial<I extends {
        msgType?: string | undefined;
        data?: Uint8Array | undefined;
    } & {
        msgType?: string | undefined;
        data?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof MsgData>, never>>(object: I): MsgData;
};
export declare const TxMsgData: {
    typeUrl: string;
    encode(message: TxMsgData, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): TxMsgData;
    fromJSON(object: any): TxMsgData;
    toJSON(message: TxMsgData): unknown;
    fromPartial<I extends {
        data?: {
            msgType?: string | undefined;
            data?: Uint8Array | undefined;
        }[] | undefined;
        msgResponses?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        }[] | undefined;
    } & {
        data?: ({
            msgType?: string | undefined;
            data?: Uint8Array | undefined;
        }[] & ({
            msgType?: string | undefined;
            data?: Uint8Array | undefined;
        } & {
            msgType?: string | undefined;
            data?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["data"][number], keyof MsgData>, never>)[] & Record<Exclude<keyof I["data"], keyof {
            msgType?: string | undefined;
            data?: Uint8Array | undefined;
        }[]>, never>) | undefined;
        msgResponses?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        }[] & ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["msgResponses"][number], keyof Any>, never>)[] & Record<Exclude<keyof I["msgResponses"], keyof {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof TxMsgData>, never>>(object: I): TxMsgData;
};
export declare const SearchTxsResult: {
    typeUrl: string;
    encode(message: SearchTxsResult, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): SearchTxsResult;
    fromJSON(object: any): SearchTxsResult;
    toJSON(message: SearchTxsResult): unknown;
    fromPartial<I extends {
        totalCount?: bigint | undefined;
        count?: bigint | undefined;
        pageNumber?: bigint | undefined;
        pageTotal?: bigint | undefined;
        limit?: bigint | undefined;
        txs?: {
            height?: bigint | undefined;
            txhash?: string | undefined;
            codespace?: string | undefined;
            code?: number | undefined;
            data?: string | undefined;
            rawLog?: string | undefined;
            logs?: {
                msgIndex?: number | undefined;
                log?: string | undefined;
                events?: {
                    type?: string | undefined;
                    attributes?: {
                        key?: string | undefined;
                        value?: string | undefined;
                    }[] | undefined;
                }[] | undefined;
            }[] | undefined;
            info?: string | undefined;
            gasWanted?: bigint | undefined;
            gasUsed?: bigint | undefined;
            tx?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            timestamp?: string | undefined;
            events?: {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[] | undefined;
        }[] | undefined;
    } & {
        totalCount?: bigint | undefined;
        count?: bigint | undefined;
        pageNumber?: bigint | undefined;
        pageTotal?: bigint | undefined;
        limit?: bigint | undefined;
        txs?: ({
            height?: bigint | undefined;
            txhash?: string | undefined;
            codespace?: string | undefined;
            code?: number | undefined;
            data?: string | undefined;
            rawLog?: string | undefined;
            logs?: {
                msgIndex?: number | undefined;
                log?: string | undefined;
                events?: {
                    type?: string | undefined;
                    attributes?: {
                        key?: string | undefined;
                        value?: string | undefined;
                    }[] | undefined;
                }[] | undefined;
            }[] | undefined;
            info?: string | undefined;
            gasWanted?: bigint | undefined;
            gasUsed?: bigint | undefined;
            tx?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            timestamp?: string | undefined;
            events?: {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[] | undefined;
        }[] & ({
            height?: bigint | undefined;
            txhash?: string | undefined;
            codespace?: string | undefined;
            code?: number | undefined;
            data?: string | undefined;
            rawLog?: string | undefined;
            logs?: {
                msgIndex?: number | undefined;
                log?: string | undefined;
                events?: {
                    type?: string | undefined;
                    attributes?: {
                        key?: string | undefined;
                        value?: string | undefined;
                    }[] | undefined;
                }[] | undefined;
            }[] | undefined;
            info?: string | undefined;
            gasWanted?: bigint | undefined;
            gasUsed?: bigint | undefined;
            tx?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            timestamp?: string | undefined;
            events?: {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[] | undefined;
        } & {
            height?: bigint | undefined;
            txhash?: string | undefined;
            codespace?: string | undefined;
            code?: number | undefined;
            data?: string | undefined;
            rawLog?: string | undefined;
            logs?: ({
                msgIndex?: number | undefined;
                log?: string | undefined;
                events?: {
                    type?: string | undefined;
                    attributes?: {
                        key?: string | undefined;
                        value?: string | undefined;
                    }[] | undefined;
                }[] | undefined;
            }[] & ({
                msgIndex?: number | undefined;
                log?: string | undefined;
                events?: {
                    type?: string | undefined;
                    attributes?: {
                        key?: string | undefined;
                        value?: string | undefined;
                    }[] | undefined;
                }[] | undefined;
            } & {
                msgIndex?: number | undefined;
                log?: string | undefined;
                events?: ({
                    type?: string | undefined;
                    attributes?: {
                        key?: string | undefined;
                        value?: string | undefined;
                    }[] | undefined;
                }[] & ({
                    type?: string | undefined;
                    attributes?: {
                        key?: string | undefined;
                        value?: string | undefined;
                    }[] | undefined;
                } & {
                    type?: string | undefined;
                    attributes?: ({
                        key?: string | undefined;
                        value?: string | undefined;
                    }[] & ({
                        key?: string | undefined;
                        value?: string | undefined;
                    } & {
                        key?: string | undefined;
                        value?: string | undefined;
                    } & Record<Exclude<keyof I["txs"][number]["logs"][number]["events"][number]["attributes"][number], keyof Attribute>, never>)[] & Record<Exclude<keyof I["txs"][number]["logs"][number]["events"][number]["attributes"], keyof {
                        key?: string | undefined;
                        value?: string | undefined;
                    }[]>, never>) | undefined;
                } & Record<Exclude<keyof I["txs"][number]["logs"][number]["events"][number], keyof StringEvent>, never>)[] & Record<Exclude<keyof I["txs"][number]["logs"][number]["events"], keyof {
                    type?: string | undefined;
                    attributes?: {
                        key?: string | undefined;
                        value?: string | undefined;
                    }[] | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["txs"][number]["logs"][number], keyof ABCIMessageLog>, never>)[] & Record<Exclude<keyof I["txs"][number]["logs"], keyof {
                msgIndex?: number | undefined;
                log?: string | undefined;
                events?: {
                    type?: string | undefined;
                    attributes?: {
                        key?: string | undefined;
                        value?: string | undefined;
                    }[] | undefined;
                }[] | undefined;
            }[]>, never>) | undefined;
            info?: string | undefined;
            gasWanted?: bigint | undefined;
            gasUsed?: bigint | undefined;
            tx?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["txs"][number]["tx"], keyof Any>, never>) | undefined;
            timestamp?: string | undefined;
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
                } & Record<Exclude<keyof I["txs"][number]["events"][number]["attributes"][number], keyof import("../../../../tendermint/abci/types").EventAttribute>, never>)[] & Record<Exclude<keyof I["txs"][number]["events"][number]["attributes"], keyof {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["txs"][number]["events"][number], keyof Event>, never>)[] & Record<Exclude<keyof I["txs"][number]["events"], keyof {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["txs"][number], keyof TxResponse>, never>)[] & Record<Exclude<keyof I["txs"], keyof {
            height?: bigint | undefined;
            txhash?: string | undefined;
            codespace?: string | undefined;
            code?: number | undefined;
            data?: string | undefined;
            rawLog?: string | undefined;
            logs?: {
                msgIndex?: number | undefined;
                log?: string | undefined;
                events?: {
                    type?: string | undefined;
                    attributes?: {
                        key?: string | undefined;
                        value?: string | undefined;
                    }[] | undefined;
                }[] | undefined;
            }[] | undefined;
            info?: string | undefined;
            gasWanted?: bigint | undefined;
            gasUsed?: bigint | undefined;
            tx?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            timestamp?: string | undefined;
            events?: {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[] | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof SearchTxsResult>, never>>(object: I): SearchTxsResult;
};
