import { Tx } from "./tx";
import { PageRequest, PageResponse } from "../../base/query/v1beta1/pagination";
import { TxResponse, GasInfo, Result } from "../../base/abci/v1beta1/abci";
import { BlockID } from "../../../tendermint/types/types";
import { Block } from "../../../tendermint/types/block";
import { BinaryReader, BinaryWriter } from "../../../binary";
import { Rpc } from "../../../helpers";
export declare const protobufPackage = "cosmos.tx.v1beta1";
/** OrderBy defines the sorting order */
export declare enum OrderBy {
    /** ORDER_BY_UNSPECIFIED - ORDER_BY_UNSPECIFIED specifies an unknown sorting order. OrderBy defaults to ASC in this case. */
    ORDER_BY_UNSPECIFIED = 0,
    /** ORDER_BY_ASC - ORDER_BY_ASC defines ascending order */
    ORDER_BY_ASC = 1,
    /** ORDER_BY_DESC - ORDER_BY_DESC defines descending order */
    ORDER_BY_DESC = 2,
    UNRECOGNIZED = -1
}
export declare function orderByFromJSON(object: any): OrderBy;
export declare function orderByToJSON(object: OrderBy): string;
/** BroadcastMode specifies the broadcast mode for the TxService.Broadcast RPC method. */
export declare enum BroadcastMode {
    /** BROADCAST_MODE_UNSPECIFIED - zero-value for mode ordering */
    BROADCAST_MODE_UNSPECIFIED = 0,
    /**
     * BROADCAST_MODE_BLOCK - DEPRECATED: use BROADCAST_MODE_SYNC instead,
     * BROADCAST_MODE_BLOCK is not supported by the SDK from v0.47.x onwards.
     */
    BROADCAST_MODE_BLOCK = 1,
    /**
     * BROADCAST_MODE_SYNC - BROADCAST_MODE_SYNC defines a tx broadcasting mode where the client waits for
     * a CheckTx execution response only.
     */
    BROADCAST_MODE_SYNC = 2,
    /**
     * BROADCAST_MODE_ASYNC - BROADCAST_MODE_ASYNC defines a tx broadcasting mode where the client returns
     * immediately.
     */
    BROADCAST_MODE_ASYNC = 3,
    UNRECOGNIZED = -1
}
export declare function broadcastModeFromJSON(object: any): BroadcastMode;
export declare function broadcastModeToJSON(object: BroadcastMode): string;
/**
 * GetTxsEventRequest is the request type for the Service.TxsByEvents
 * RPC method.
 */
export interface GetTxsEventRequest {
    /** events is the list of transaction event type. */
    events: string[];
    /**
     * pagination defines a pagination for the request.
     * Deprecated post v0.46.x: use page and limit instead.
     */
    /** @deprecated */
    pagination?: PageRequest;
    orderBy: OrderBy;
    /** page is the page number to query, starts at 1. If not provided, will default to first page. */
    page: bigint;
    /**
     * limit is the total number of results to be returned in the result page.
     * If left empty it will default to a value to be set by each app.
     */
    limit: bigint;
}
/**
 * GetTxsEventResponse is the response type for the Service.TxsByEvents
 * RPC method.
 */
export interface GetTxsEventResponse {
    /** txs is the list of queried transactions. */
    txs: Tx[];
    /** tx_responses is the list of queried TxResponses. */
    txResponses: TxResponse[];
    /**
     * pagination defines a pagination for the response.
     * Deprecated post v0.46.x: use total instead.
     */
    /** @deprecated */
    pagination?: PageResponse;
    /** total is total number of results available */
    total: bigint;
}
/**
 * BroadcastTxRequest is the request type for the Service.BroadcastTxRequest
 * RPC method.
 */
export interface BroadcastTxRequest {
    /** tx_bytes is the raw transaction. */
    txBytes: Uint8Array;
    mode: BroadcastMode;
}
/**
 * BroadcastTxResponse is the response type for the
 * Service.BroadcastTx method.
 */
export interface BroadcastTxResponse {
    /** tx_response is the queried TxResponses. */
    txResponse?: TxResponse;
}
/**
 * SimulateRequest is the request type for the Service.Simulate
 * RPC method.
 */
export interface SimulateRequest {
    /**
     * tx is the transaction to simulate.
     * Deprecated. Send raw tx bytes instead.
     */
    /** @deprecated */
    tx?: Tx;
    /**
     * tx_bytes is the raw transaction.
     *
     * Since: cosmos-sdk 0.43
     */
    txBytes: Uint8Array;
}
/**
 * SimulateResponse is the response type for the
 * Service.SimulateRPC method.
 */
export interface SimulateResponse {
    /** gas_info is the information about gas used in the simulation. */
    gasInfo?: GasInfo;
    /** result is the result of the simulation. */
    result?: Result;
}
/**
 * GetTxRequest is the request type for the Service.GetTx
 * RPC method.
 */
export interface GetTxRequest {
    /** hash is the tx hash to query, encoded as a hex string. */
    hash: string;
}
/** GetTxResponse is the response type for the Service.GetTx method. */
export interface GetTxResponse {
    /** tx is the queried transaction. */
    tx?: Tx;
    /** tx_response is the queried TxResponses. */
    txResponse?: TxResponse;
}
/**
 * GetBlockWithTxsRequest is the request type for the Service.GetBlockWithTxs
 * RPC method.
 *
 * Since: cosmos-sdk 0.45.2
 */
export interface GetBlockWithTxsRequest {
    /** height is the height of the block to query. */
    height: bigint;
    /** pagination defines a pagination for the request. */
    pagination?: PageRequest;
}
/**
 * GetBlockWithTxsResponse is the response type for the Service.GetBlockWithTxs method.
 *
 * Since: cosmos-sdk 0.45.2
 */
export interface GetBlockWithTxsResponse {
    /** txs are the transactions in the block. */
    txs: Tx[];
    blockId?: BlockID;
    block?: Block;
    /** pagination defines a pagination for the response. */
    pagination?: PageResponse;
}
/**
 * TxDecodeRequest is the request type for the Service.TxDecode
 * RPC method.
 *
 * Since: cosmos-sdk 0.47
 */
export interface TxDecodeRequest {
    /** tx_bytes is the raw transaction. */
    txBytes: Uint8Array;
}
/**
 * TxDecodeResponse is the response type for the
 * Service.TxDecode method.
 *
 * Since: cosmos-sdk 0.47
 */
export interface TxDecodeResponse {
    /** tx is the decoded transaction. */
    tx?: Tx;
}
/**
 * TxEncodeRequest is the request type for the Service.TxEncode
 * RPC method.
 *
 * Since: cosmos-sdk 0.47
 */
export interface TxEncodeRequest {
    /** tx is the transaction to encode. */
    tx?: Tx;
}
/**
 * TxEncodeResponse is the response type for the
 * Service.TxEncode method.
 *
 * Since: cosmos-sdk 0.47
 */
export interface TxEncodeResponse {
    /** tx_bytes is the encoded transaction bytes. */
    txBytes: Uint8Array;
}
/**
 * TxEncodeAminoRequest is the request type for the Service.TxEncodeAmino
 * RPC method.
 *
 * Since: cosmos-sdk 0.47
 */
export interface TxEncodeAminoRequest {
    aminoJson: string;
}
/**
 * TxEncodeAminoResponse is the response type for the Service.TxEncodeAmino
 * RPC method.
 *
 * Since: cosmos-sdk 0.47
 */
export interface TxEncodeAminoResponse {
    aminoBinary: Uint8Array;
}
/**
 * TxDecodeAminoRequest is the request type for the Service.TxDecodeAmino
 * RPC method.
 *
 * Since: cosmos-sdk 0.47
 */
export interface TxDecodeAminoRequest {
    aminoBinary: Uint8Array;
}
/**
 * TxDecodeAminoResponse is the response type for the Service.TxDecodeAmino
 * RPC method.
 *
 * Since: cosmos-sdk 0.47
 */
export interface TxDecodeAminoResponse {
    aminoJson: string;
}
export declare const GetTxsEventRequest: {
    typeUrl: string;
    encode(message: GetTxsEventRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GetTxsEventRequest;
    fromJSON(object: any): GetTxsEventRequest;
    toJSON(message: GetTxsEventRequest): unknown;
    fromPartial<I extends {
        events?: string[] | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
        orderBy?: OrderBy | undefined;
        page?: bigint | undefined;
        limit?: bigint | undefined;
    } & {
        events?: (string[] & string[] & Record<Exclude<keyof I["events"], keyof string[]>, never>) | undefined;
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
        orderBy?: OrderBy | undefined;
        page?: bigint | undefined;
        limit?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof GetTxsEventRequest>, never>>(object: I): GetTxsEventRequest;
};
export declare const GetTxsEventResponse: {
    typeUrl: string;
    encode(message: GetTxsEventResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GetTxsEventResponse;
    fromJSON(object: any): GetTxsEventResponse;
    toJSON(message: GetTxsEventResponse): unknown;
    fromPartial<I extends {
        txs?: {
            body?: {
                messages?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                memo?: string | undefined;
                timeoutHeight?: bigint | undefined;
                extensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                nonCriticalExtensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
            } | undefined;
            authInfo?: {
                signerInfos?: {
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[] | undefined;
                fee?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } | undefined;
                tip?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    tipper?: string | undefined;
                } | undefined;
            } | undefined;
            signatures?: Uint8Array[] | undefined;
        }[] | undefined;
        txResponses?: {
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
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
        total?: bigint | undefined;
    } & {
        txs?: ({
            body?: {
                messages?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                memo?: string | undefined;
                timeoutHeight?: bigint | undefined;
                extensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                nonCriticalExtensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
            } | undefined;
            authInfo?: {
                signerInfos?: {
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[] | undefined;
                fee?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } | undefined;
                tip?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    tipper?: string | undefined;
                } | undefined;
            } | undefined;
            signatures?: Uint8Array[] | undefined;
        }[] & ({
            body?: {
                messages?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                memo?: string | undefined;
                timeoutHeight?: bigint | undefined;
                extensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                nonCriticalExtensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
            } | undefined;
            authInfo?: {
                signerInfos?: {
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[] | undefined;
                fee?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } | undefined;
                tip?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    tipper?: string | undefined;
                } | undefined;
            } | undefined;
            signatures?: Uint8Array[] | undefined;
        } & {
            body?: ({
                messages?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                memo?: string | undefined;
                timeoutHeight?: bigint | undefined;
                extensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                nonCriticalExtensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
            } & {
                messages?: ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] & ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["txs"][number]["body"]["messages"][number], keyof import("../../../google/protobuf/any").Any>, never>)[] & Record<Exclude<keyof I["txs"][number]["body"]["messages"], keyof {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[]>, never>) | undefined;
                memo?: string | undefined;
                timeoutHeight?: bigint | undefined;
                extensionOptions?: ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] & ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["txs"][number]["body"]["extensionOptions"][number], keyof import("../../../google/protobuf/any").Any>, never>)[] & Record<Exclude<keyof I["txs"][number]["body"]["extensionOptions"], keyof {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[]>, never>) | undefined;
                nonCriticalExtensionOptions?: ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] & ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["txs"][number]["body"]["nonCriticalExtensionOptions"][number], keyof import("../../../google/protobuf/any").Any>, never>)[] & Record<Exclude<keyof I["txs"][number]["body"]["nonCriticalExtensionOptions"], keyof {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["txs"][number]["body"], keyof import("./tx").TxBody>, never>) | undefined;
            authInfo?: ({
                signerInfos?: {
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[] | undefined;
                fee?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } | undefined;
                tip?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    tipper?: string | undefined;
                } | undefined;
            } & {
                signerInfos?: ({
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[] & ({
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                } & {
                    publicKey?: ({
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } & {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["publicKey"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
                    modeInfo?: ({
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } & {
                        single?: ({
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } & {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["single"], "mode">, never>) | undefined;
                        multi?: ({
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: {
                                single?: {
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } | undefined;
                                multi?: any | undefined;
                            }[] | undefined;
                        } & {
                            bitarray?: ({
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } & {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["bitarray"], keyof import("../../crypto/multisig/v1beta1/multisig").CompactBitArray>, never>) | undefined;
                            modeInfos?: ({
                                single?: {
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } | undefined;
                                multi?: {
                                    bitarray?: {
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } | undefined;
                                    modeInfos?: any[] | undefined;
                                } | undefined;
                            }[] & ({
                                single?: {
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } | undefined;
                                multi?: {
                                    bitarray?: {
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } | undefined;
                                    modeInfos?: any[] | undefined;
                                } | undefined;
                            } & {
                                single?: ({
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } & {
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["single"], "mode">, never>) | undefined;
                                multi?: ({
                                    bitarray?: {
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } | undefined;
                                    modeInfos?: {
                                        single?: {
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } | undefined;
                                        multi?: any | undefined;
                                    }[] | undefined;
                                } & {
                                    bitarray?: ({
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } & {
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["bitarray"], keyof import("../../crypto/multisig/v1beta1/multisig").CompactBitArray>, never>) | undefined;
                                    modeInfos?: ({
                                        single?: {
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } | undefined;
                                        multi?: {
                                            bitarray?: {
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } | undefined;
                                            modeInfos?: any[] | undefined;
                                        } | undefined;
                                    }[] & ({
                                        single?: {
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } | undefined;
                                        multi?: {
                                            bitarray?: {
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } | undefined;
                                            modeInfos?: any[] | undefined;
                                        } | undefined;
                                    } & {
                                        single?: ({
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } & {
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["single"], "mode">, never>) | undefined;
                                        multi?: ({
                                            bitarray?: {
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } | undefined;
                                            modeInfos?: {
                                                single?: {
                                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                } | undefined;
                                                multi?: any | undefined;
                                            }[] | undefined;
                                        } & {
                                            bitarray?: ({
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } & {
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"]["bitarray"], keyof import("../../crypto/multisig/v1beta1/multisig").CompactBitArray>, never>) | undefined;
                                            modeInfos?: ({
                                                single?: {
                                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                } | undefined;
                                                multi?: {
                                                    bitarray?: {
                                                        extraBitsStored?: number | undefined;
                                                        elems?: Uint8Array | undefined;
                                                    } | undefined;
                                                    modeInfos?: any[] | undefined;
                                                } | undefined;
                                            }[] & ({
                                                single?: {
                                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                } | undefined;
                                                multi?: {
                                                    bitarray?: {
                                                        extraBitsStored?: number | undefined;
                                                        elems?: Uint8Array | undefined;
                                                    } | undefined;
                                                    modeInfos?: any[] | undefined;
                                                } | undefined;
                                            } & {
                                                single?: ({
                                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                } & any & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["single"], "mode">, never>) | undefined;
                                                multi?: ({
                                                    bitarray?: {
                                                        extraBitsStored?: number | undefined;
                                                        elems?: Uint8Array | undefined;
                                                    } | undefined;
                                                    modeInfos?: {
                                                        single?: {
                                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                        } | undefined;
                                                        multi?: any | undefined;
                                                    }[] | undefined;
                                                } & any & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"], keyof import("./tx").ModeInfo_Multi>, never>) | undefined;
                                            } & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number], keyof import("./tx").ModeInfo>, never>)[] & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"]["modeInfos"], keyof {
                                                single?: {
                                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                } | undefined;
                                                multi?: {
                                                    bitarray?: {
                                                        extraBitsStored?: number | undefined;
                                                        elems?: Uint8Array | undefined;
                                                    } | undefined;
                                                    modeInfos?: any[] | undefined;
                                                } | undefined;
                                            }[]>, never>) | undefined;
                                        } & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"], keyof import("./tx").ModeInfo_Multi>, never>) | undefined;
                                    } & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number], keyof import("./tx").ModeInfo>, never>)[] & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"], keyof {
                                        single?: {
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } | undefined;
                                        multi?: {
                                            bitarray?: {
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } | undefined;
                                            modeInfos?: any[] | undefined;
                                        } | undefined;
                                    }[]>, never>) | undefined;
                                } & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"], keyof import("./tx").ModeInfo_Multi>, never>) | undefined;
                            } & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number], keyof import("./tx").ModeInfo>, never>)[] & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"], keyof {
                                single?: {
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } | undefined;
                                multi?: {
                                    bitarray?: {
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } | undefined;
                                    modeInfos?: any[] | undefined;
                                } | undefined;
                            }[]>, never>) | undefined;
                        } & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"], keyof import("./tx").ModeInfo_Multi>, never>) | undefined;
                    } & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"], keyof import("./tx").ModeInfo>, never>) | undefined;
                    sequence?: bigint | undefined;
                } & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number], keyof import("./tx").SignerInfo>, never>)[] & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"], keyof {
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[]>, never>) | undefined;
                fee?: ({
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } & {
                    amount?: ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] & ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & Record<Exclude<keyof I["txs"][number]["authInfo"]["fee"]["amount"][number], keyof import("../../base/v1beta1/coin").Coin>, never>)[] & Record<Exclude<keyof I["txs"][number]["authInfo"]["fee"]["amount"], keyof {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[]>, never>) | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } & Record<Exclude<keyof I["txs"][number]["authInfo"]["fee"], keyof import("./tx").Fee>, never>) | undefined;
                tip?: ({
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    tipper?: string | undefined;
                } & {
                    amount?: ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] & ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & Record<Exclude<keyof I["txs"][number]["authInfo"]["tip"]["amount"][number], keyof import("../../base/v1beta1/coin").Coin>, never>)[] & Record<Exclude<keyof I["txs"][number]["authInfo"]["tip"]["amount"], keyof {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[]>, never>) | undefined;
                    tipper?: string | undefined;
                } & Record<Exclude<keyof I["txs"][number]["authInfo"]["tip"], keyof import("./tx").Tip>, never>) | undefined;
            } & Record<Exclude<keyof I["txs"][number]["authInfo"], keyof import("./tx").AuthInfo>, never>) | undefined;
            signatures?: (Uint8Array[] & Uint8Array[] & Record<Exclude<keyof I["txs"][number]["signatures"], keyof Uint8Array[]>, never>) | undefined;
        } & Record<Exclude<keyof I["txs"][number], keyof Tx>, never>)[] & Record<Exclude<keyof I["txs"], keyof {
            body?: {
                messages?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                memo?: string | undefined;
                timeoutHeight?: bigint | undefined;
                extensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                nonCriticalExtensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
            } | undefined;
            authInfo?: {
                signerInfos?: {
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[] | undefined;
                fee?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } | undefined;
                tip?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    tipper?: string | undefined;
                } | undefined;
            } | undefined;
            signatures?: Uint8Array[] | undefined;
        }[]>, never>) | undefined;
        txResponses?: ({
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
                    } & Record<Exclude<keyof I["txResponses"][number]["logs"][number]["events"][number]["attributes"][number], keyof import("../../base/abci/v1beta1/abci").Attribute>, never>)[] & Record<Exclude<keyof I["txResponses"][number]["logs"][number]["events"][number]["attributes"], keyof {
                        key?: string | undefined;
                        value?: string | undefined;
                    }[]>, never>) | undefined;
                } & Record<Exclude<keyof I["txResponses"][number]["logs"][number]["events"][number], keyof import("../../base/abci/v1beta1/abci").StringEvent>, never>)[] & Record<Exclude<keyof I["txResponses"][number]["logs"][number]["events"], keyof {
                    type?: string | undefined;
                    attributes?: {
                        key?: string | undefined;
                        value?: string | undefined;
                    }[] | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["txResponses"][number]["logs"][number], keyof import("../../base/abci/v1beta1/abci").ABCIMessageLog>, never>)[] & Record<Exclude<keyof I["txResponses"][number]["logs"], keyof {
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
            } & Record<Exclude<keyof I["txResponses"][number]["tx"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
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
                } & Record<Exclude<keyof I["txResponses"][number]["events"][number]["attributes"][number], keyof import("../../../tendermint/abci/types").EventAttribute>, never>)[] & Record<Exclude<keyof I["txResponses"][number]["events"][number]["attributes"], keyof {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["txResponses"][number]["events"][number], keyof import("../../../tendermint/abci/types").Event>, never>)[] & Record<Exclude<keyof I["txResponses"][number]["events"], keyof {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["txResponses"][number], keyof TxResponse>, never>)[] & Record<Exclude<keyof I["txResponses"], keyof {
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
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
        total?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof GetTxsEventResponse>, never>>(object: I): GetTxsEventResponse;
};
export declare const BroadcastTxRequest: {
    typeUrl: string;
    encode(message: BroadcastTxRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): BroadcastTxRequest;
    fromJSON(object: any): BroadcastTxRequest;
    toJSON(message: BroadcastTxRequest): unknown;
    fromPartial<I extends {
        txBytes?: Uint8Array | undefined;
        mode?: BroadcastMode | undefined;
    } & {
        txBytes?: Uint8Array | undefined;
        mode?: BroadcastMode | undefined;
    } & Record<Exclude<keyof I, keyof BroadcastTxRequest>, never>>(object: I): BroadcastTxRequest;
};
export declare const BroadcastTxResponse: {
    typeUrl: string;
    encode(message: BroadcastTxResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): BroadcastTxResponse;
    fromJSON(object: any): BroadcastTxResponse;
    toJSON(message: BroadcastTxResponse): unknown;
    fromPartial<I extends {
        txResponse?: {
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
        } | undefined;
    } & {
        txResponse?: ({
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
                    } & Record<Exclude<keyof I["txResponse"]["logs"][number]["events"][number]["attributes"][number], keyof import("../../base/abci/v1beta1/abci").Attribute>, never>)[] & Record<Exclude<keyof I["txResponse"]["logs"][number]["events"][number]["attributes"], keyof {
                        key?: string | undefined;
                        value?: string | undefined;
                    }[]>, never>) | undefined;
                } & Record<Exclude<keyof I["txResponse"]["logs"][number]["events"][number], keyof import("../../base/abci/v1beta1/abci").StringEvent>, never>)[] & Record<Exclude<keyof I["txResponse"]["logs"][number]["events"], keyof {
                    type?: string | undefined;
                    attributes?: {
                        key?: string | undefined;
                        value?: string | undefined;
                    }[] | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["txResponse"]["logs"][number], keyof import("../../base/abci/v1beta1/abci").ABCIMessageLog>, never>)[] & Record<Exclude<keyof I["txResponse"]["logs"], keyof {
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
            } & Record<Exclude<keyof I["txResponse"]["tx"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
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
                } & Record<Exclude<keyof I["txResponse"]["events"][number]["attributes"][number], keyof import("../../../tendermint/abci/types").EventAttribute>, never>)[] & Record<Exclude<keyof I["txResponse"]["events"][number]["attributes"], keyof {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["txResponse"]["events"][number], keyof import("../../../tendermint/abci/types").Event>, never>)[] & Record<Exclude<keyof I["txResponse"]["events"], keyof {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["txResponse"], keyof TxResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, "txResponse">, never>>(object: I): BroadcastTxResponse;
};
export declare const SimulateRequest: {
    typeUrl: string;
    encode(message: SimulateRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): SimulateRequest;
    fromJSON(object: any): SimulateRequest;
    toJSON(message: SimulateRequest): unknown;
    fromPartial<I extends {
        tx?: {
            body?: {
                messages?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                memo?: string | undefined;
                timeoutHeight?: bigint | undefined;
                extensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                nonCriticalExtensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
            } | undefined;
            authInfo?: {
                signerInfos?: {
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[] | undefined;
                fee?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } | undefined;
                tip?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    tipper?: string | undefined;
                } | undefined;
            } | undefined;
            signatures?: Uint8Array[] | undefined;
        } | undefined;
        txBytes?: Uint8Array | undefined;
    } & {
        tx?: ({
            body?: {
                messages?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                memo?: string | undefined;
                timeoutHeight?: bigint | undefined;
                extensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                nonCriticalExtensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
            } | undefined;
            authInfo?: {
                signerInfos?: {
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[] | undefined;
                fee?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } | undefined;
                tip?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    tipper?: string | undefined;
                } | undefined;
            } | undefined;
            signatures?: Uint8Array[] | undefined;
        } & {
            body?: ({
                messages?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                memo?: string | undefined;
                timeoutHeight?: bigint | undefined;
                extensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                nonCriticalExtensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
            } & {
                messages?: ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] & ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["tx"]["body"]["messages"][number], keyof import("../../../google/protobuf/any").Any>, never>)[] & Record<Exclude<keyof I["tx"]["body"]["messages"], keyof {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[]>, never>) | undefined;
                memo?: string | undefined;
                timeoutHeight?: bigint | undefined;
                extensionOptions?: ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] & ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["tx"]["body"]["extensionOptions"][number], keyof import("../../../google/protobuf/any").Any>, never>)[] & Record<Exclude<keyof I["tx"]["body"]["extensionOptions"], keyof {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[]>, never>) | undefined;
                nonCriticalExtensionOptions?: ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] & ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["tx"]["body"]["nonCriticalExtensionOptions"][number], keyof import("../../../google/protobuf/any").Any>, never>)[] & Record<Exclude<keyof I["tx"]["body"]["nonCriticalExtensionOptions"], keyof {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["tx"]["body"], keyof import("./tx").TxBody>, never>) | undefined;
            authInfo?: ({
                signerInfos?: {
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[] | undefined;
                fee?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } | undefined;
                tip?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    tipper?: string | undefined;
                } | undefined;
            } & {
                signerInfos?: ({
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[] & ({
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                } & {
                    publicKey?: ({
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } & {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["publicKey"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
                    modeInfo?: ({
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } & {
                        single?: ({
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } & {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["single"], "mode">, never>) | undefined;
                        multi?: ({
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: {
                                single?: {
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } | undefined;
                                multi?: any | undefined;
                            }[] | undefined;
                        } & {
                            bitarray?: ({
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } & {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["bitarray"], keyof import("../../crypto/multisig/v1beta1/multisig").CompactBitArray>, never>) | undefined;
                            modeInfos?: ({
                                single?: {
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } | undefined;
                                multi?: {
                                    bitarray?: {
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } | undefined;
                                    modeInfos?: any[] | undefined;
                                } | undefined;
                            }[] & ({
                                single?: {
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } | undefined;
                                multi?: {
                                    bitarray?: {
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } | undefined;
                                    modeInfos?: any[] | undefined;
                                } | undefined;
                            } & {
                                single?: ({
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } & {
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["single"], "mode">, never>) | undefined;
                                multi?: ({
                                    bitarray?: {
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } | undefined;
                                    modeInfos?: {
                                        single?: {
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } | undefined;
                                        multi?: any | undefined;
                                    }[] | undefined;
                                } & {
                                    bitarray?: ({
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } & {
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["bitarray"], keyof import("../../crypto/multisig/v1beta1/multisig").CompactBitArray>, never>) | undefined;
                                    modeInfos?: ({
                                        single?: {
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } | undefined;
                                        multi?: {
                                            bitarray?: {
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } | undefined;
                                            modeInfos?: any[] | undefined;
                                        } | undefined;
                                    }[] & ({
                                        single?: {
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } | undefined;
                                        multi?: {
                                            bitarray?: {
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } | undefined;
                                            modeInfos?: any[] | undefined;
                                        } | undefined;
                                    } & {
                                        single?: ({
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } & {
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["single"], "mode">, never>) | undefined;
                                        multi?: ({
                                            bitarray?: {
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } | undefined;
                                            modeInfos?: {
                                                single?: {
                                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                } | undefined;
                                                multi?: any | undefined;
                                            }[] | undefined;
                                        } & {
                                            bitarray?: ({
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } & {
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"]["bitarray"], keyof import("../../crypto/multisig/v1beta1/multisig").CompactBitArray>, never>) | undefined;
                                            modeInfos?: ({
                                                single?: {
                                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                } | undefined;
                                                multi?: {
                                                    bitarray?: {
                                                        extraBitsStored?: number | undefined;
                                                        elems?: Uint8Array | undefined;
                                                    } | undefined;
                                                    modeInfos?: any[] | undefined;
                                                } | undefined;
                                            }[] & ({
                                                single?: {
                                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                } | undefined;
                                                multi?: {
                                                    bitarray?: {
                                                        extraBitsStored?: number | undefined;
                                                        elems?: Uint8Array | undefined;
                                                    } | undefined;
                                                    modeInfos?: any[] | undefined;
                                                } | undefined;
                                            } & {
                                                single?: ({
                                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                } & any & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["single"], "mode">, never>) | undefined;
                                                multi?: ({
                                                    bitarray?: {
                                                        extraBitsStored?: number | undefined;
                                                        elems?: Uint8Array | undefined;
                                                    } | undefined;
                                                    modeInfos?: {
                                                        single?: {
                                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                        } | undefined;
                                                        multi?: any | undefined;
                                                    }[] | undefined;
                                                } & any & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"], keyof import("./tx").ModeInfo_Multi>, never>) | undefined;
                                            } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number], keyof import("./tx").ModeInfo>, never>)[] & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"]["modeInfos"], keyof {
                                                single?: {
                                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                } | undefined;
                                                multi?: {
                                                    bitarray?: {
                                                        extraBitsStored?: number | undefined;
                                                        elems?: Uint8Array | undefined;
                                                    } | undefined;
                                                    modeInfos?: any[] | undefined;
                                                } | undefined;
                                            }[]>, never>) | undefined;
                                        } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"], keyof import("./tx").ModeInfo_Multi>, never>) | undefined;
                                    } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number], keyof import("./tx").ModeInfo>, never>)[] & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"], keyof {
                                        single?: {
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } | undefined;
                                        multi?: {
                                            bitarray?: {
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } | undefined;
                                            modeInfos?: any[] | undefined;
                                        } | undefined;
                                    }[]>, never>) | undefined;
                                } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"], keyof import("./tx").ModeInfo_Multi>, never>) | undefined;
                            } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number], keyof import("./tx").ModeInfo>, never>)[] & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"], keyof {
                                single?: {
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } | undefined;
                                multi?: {
                                    bitarray?: {
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } | undefined;
                                    modeInfos?: any[] | undefined;
                                } | undefined;
                            }[]>, never>) | undefined;
                        } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"], keyof import("./tx").ModeInfo_Multi>, never>) | undefined;
                    } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"], keyof import("./tx").ModeInfo>, never>) | undefined;
                    sequence?: bigint | undefined;
                } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number], keyof import("./tx").SignerInfo>, never>)[] & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"], keyof {
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[]>, never>) | undefined;
                fee?: ({
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } & {
                    amount?: ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] & ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & Record<Exclude<keyof I["tx"]["authInfo"]["fee"]["amount"][number], keyof import("../../base/v1beta1/coin").Coin>, never>)[] & Record<Exclude<keyof I["tx"]["authInfo"]["fee"]["amount"], keyof {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[]>, never>) | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } & Record<Exclude<keyof I["tx"]["authInfo"]["fee"], keyof import("./tx").Fee>, never>) | undefined;
                tip?: ({
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    tipper?: string | undefined;
                } & {
                    amount?: ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] & ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & Record<Exclude<keyof I["tx"]["authInfo"]["tip"]["amount"][number], keyof import("../../base/v1beta1/coin").Coin>, never>)[] & Record<Exclude<keyof I["tx"]["authInfo"]["tip"]["amount"], keyof {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[]>, never>) | undefined;
                    tipper?: string | undefined;
                } & Record<Exclude<keyof I["tx"]["authInfo"]["tip"], keyof import("./tx").Tip>, never>) | undefined;
            } & Record<Exclude<keyof I["tx"]["authInfo"], keyof import("./tx").AuthInfo>, never>) | undefined;
            signatures?: (Uint8Array[] & Uint8Array[] & Record<Exclude<keyof I["tx"]["signatures"], keyof Uint8Array[]>, never>) | undefined;
        } & Record<Exclude<keyof I["tx"], keyof Tx>, never>) | undefined;
        txBytes?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof SimulateRequest>, never>>(object: I): SimulateRequest;
};
export declare const SimulateResponse: {
    typeUrl: string;
    encode(message: SimulateResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): SimulateResponse;
    fromJSON(object: any): SimulateResponse;
    toJSON(message: SimulateResponse): unknown;
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
                } & Record<Exclude<keyof I["result"]["events"][number]["attributes"][number], keyof import("../../../tendermint/abci/types").EventAttribute>, never>)[] & Record<Exclude<keyof I["result"]["events"][number]["attributes"], keyof {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["result"]["events"][number], keyof import("../../../tendermint/abci/types").Event>, never>)[] & Record<Exclude<keyof I["result"]["events"], keyof {
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
            } & Record<Exclude<keyof I["result"]["msgResponses"][number], keyof import("../../../google/protobuf/any").Any>, never>)[] & Record<Exclude<keyof I["result"]["msgResponses"], keyof {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["result"], keyof Result>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof SimulateResponse>, never>>(object: I): SimulateResponse;
};
export declare const GetTxRequest: {
    typeUrl: string;
    encode(message: GetTxRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GetTxRequest;
    fromJSON(object: any): GetTxRequest;
    toJSON(message: GetTxRequest): unknown;
    fromPartial<I extends {
        hash?: string | undefined;
    } & {
        hash?: string | undefined;
    } & Record<Exclude<keyof I, "hash">, never>>(object: I): GetTxRequest;
};
export declare const GetTxResponse: {
    typeUrl: string;
    encode(message: GetTxResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GetTxResponse;
    fromJSON(object: any): GetTxResponse;
    toJSON(message: GetTxResponse): unknown;
    fromPartial<I extends {
        tx?: {
            body?: {
                messages?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                memo?: string | undefined;
                timeoutHeight?: bigint | undefined;
                extensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                nonCriticalExtensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
            } | undefined;
            authInfo?: {
                signerInfos?: {
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[] | undefined;
                fee?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } | undefined;
                tip?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    tipper?: string | undefined;
                } | undefined;
            } | undefined;
            signatures?: Uint8Array[] | undefined;
        } | undefined;
        txResponse?: {
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
        } | undefined;
    } & {
        tx?: ({
            body?: {
                messages?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                memo?: string | undefined;
                timeoutHeight?: bigint | undefined;
                extensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                nonCriticalExtensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
            } | undefined;
            authInfo?: {
                signerInfos?: {
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[] | undefined;
                fee?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } | undefined;
                tip?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    tipper?: string | undefined;
                } | undefined;
            } | undefined;
            signatures?: Uint8Array[] | undefined;
        } & {
            body?: ({
                messages?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                memo?: string | undefined;
                timeoutHeight?: bigint | undefined;
                extensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                nonCriticalExtensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
            } & {
                messages?: ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] & ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["tx"]["body"]["messages"][number], keyof import("../../../google/protobuf/any").Any>, never>)[] & Record<Exclude<keyof I["tx"]["body"]["messages"], keyof {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[]>, never>) | undefined;
                memo?: string | undefined;
                timeoutHeight?: bigint | undefined;
                extensionOptions?: ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] & ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["tx"]["body"]["extensionOptions"][number], keyof import("../../../google/protobuf/any").Any>, never>)[] & Record<Exclude<keyof I["tx"]["body"]["extensionOptions"], keyof {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[]>, never>) | undefined;
                nonCriticalExtensionOptions?: ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] & ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["tx"]["body"]["nonCriticalExtensionOptions"][number], keyof import("../../../google/protobuf/any").Any>, never>)[] & Record<Exclude<keyof I["tx"]["body"]["nonCriticalExtensionOptions"], keyof {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["tx"]["body"], keyof import("./tx").TxBody>, never>) | undefined;
            authInfo?: ({
                signerInfos?: {
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[] | undefined;
                fee?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } | undefined;
                tip?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    tipper?: string | undefined;
                } | undefined;
            } & {
                signerInfos?: ({
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[] & ({
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                } & {
                    publicKey?: ({
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } & {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["publicKey"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
                    modeInfo?: ({
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } & {
                        single?: ({
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } & {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["single"], "mode">, never>) | undefined;
                        multi?: ({
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: {
                                single?: {
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } | undefined;
                                multi?: any | undefined;
                            }[] | undefined;
                        } & {
                            bitarray?: ({
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } & {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["bitarray"], keyof import("../../crypto/multisig/v1beta1/multisig").CompactBitArray>, never>) | undefined;
                            modeInfos?: ({
                                single?: {
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } | undefined;
                                multi?: {
                                    bitarray?: {
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } | undefined;
                                    modeInfos?: any[] | undefined;
                                } | undefined;
                            }[] & ({
                                single?: {
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } | undefined;
                                multi?: {
                                    bitarray?: {
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } | undefined;
                                    modeInfos?: any[] | undefined;
                                } | undefined;
                            } & {
                                single?: ({
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } & {
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["single"], "mode">, never>) | undefined;
                                multi?: ({
                                    bitarray?: {
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } | undefined;
                                    modeInfos?: {
                                        single?: {
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } | undefined;
                                        multi?: any | undefined;
                                    }[] | undefined;
                                } & {
                                    bitarray?: ({
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } & {
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["bitarray"], keyof import("../../crypto/multisig/v1beta1/multisig").CompactBitArray>, never>) | undefined;
                                    modeInfos?: ({
                                        single?: {
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } | undefined;
                                        multi?: {
                                            bitarray?: {
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } | undefined;
                                            modeInfos?: any[] | undefined;
                                        } | undefined;
                                    }[] & ({
                                        single?: {
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } | undefined;
                                        multi?: {
                                            bitarray?: {
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } | undefined;
                                            modeInfos?: any[] | undefined;
                                        } | undefined;
                                    } & {
                                        single?: ({
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } & {
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["single"], "mode">, never>) | undefined;
                                        multi?: ({
                                            bitarray?: {
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } | undefined;
                                            modeInfos?: {
                                                single?: {
                                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                } | undefined;
                                                multi?: any | undefined;
                                            }[] | undefined;
                                        } & {
                                            bitarray?: ({
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } & {
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"]["bitarray"], keyof import("../../crypto/multisig/v1beta1/multisig").CompactBitArray>, never>) | undefined;
                                            modeInfos?: ({
                                                single?: {
                                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                } | undefined;
                                                multi?: {
                                                    bitarray?: {
                                                        extraBitsStored?: number | undefined;
                                                        elems?: Uint8Array | undefined;
                                                    } | undefined;
                                                    modeInfos?: any[] | undefined;
                                                } | undefined;
                                            }[] & ({
                                                single?: {
                                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                } | undefined;
                                                multi?: {
                                                    bitarray?: {
                                                        extraBitsStored?: number | undefined;
                                                        elems?: Uint8Array | undefined;
                                                    } | undefined;
                                                    modeInfos?: any[] | undefined;
                                                } | undefined;
                                            } & {
                                                single?: ({
                                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                } & any & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["single"], "mode">, never>) | undefined;
                                                multi?: ({
                                                    bitarray?: {
                                                        extraBitsStored?: number | undefined;
                                                        elems?: Uint8Array | undefined;
                                                    } | undefined;
                                                    modeInfos?: {
                                                        single?: {
                                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                        } | undefined;
                                                        multi?: any | undefined;
                                                    }[] | undefined;
                                                } & any & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"], keyof import("./tx").ModeInfo_Multi>, never>) | undefined;
                                            } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number], keyof import("./tx").ModeInfo>, never>)[] & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"]["modeInfos"], keyof {
                                                single?: {
                                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                } | undefined;
                                                multi?: {
                                                    bitarray?: {
                                                        extraBitsStored?: number | undefined;
                                                        elems?: Uint8Array | undefined;
                                                    } | undefined;
                                                    modeInfos?: any[] | undefined;
                                                } | undefined;
                                            }[]>, never>) | undefined;
                                        } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"], keyof import("./tx").ModeInfo_Multi>, never>) | undefined;
                                    } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number], keyof import("./tx").ModeInfo>, never>)[] & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"], keyof {
                                        single?: {
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } | undefined;
                                        multi?: {
                                            bitarray?: {
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } | undefined;
                                            modeInfos?: any[] | undefined;
                                        } | undefined;
                                    }[]>, never>) | undefined;
                                } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"], keyof import("./tx").ModeInfo_Multi>, never>) | undefined;
                            } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number], keyof import("./tx").ModeInfo>, never>)[] & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"], keyof {
                                single?: {
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } | undefined;
                                multi?: {
                                    bitarray?: {
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } | undefined;
                                    modeInfos?: any[] | undefined;
                                } | undefined;
                            }[]>, never>) | undefined;
                        } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"], keyof import("./tx").ModeInfo_Multi>, never>) | undefined;
                    } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"], keyof import("./tx").ModeInfo>, never>) | undefined;
                    sequence?: bigint | undefined;
                } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number], keyof import("./tx").SignerInfo>, never>)[] & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"], keyof {
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[]>, never>) | undefined;
                fee?: ({
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } & {
                    amount?: ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] & ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & Record<Exclude<keyof I["tx"]["authInfo"]["fee"]["amount"][number], keyof import("../../base/v1beta1/coin").Coin>, never>)[] & Record<Exclude<keyof I["tx"]["authInfo"]["fee"]["amount"], keyof {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[]>, never>) | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } & Record<Exclude<keyof I["tx"]["authInfo"]["fee"], keyof import("./tx").Fee>, never>) | undefined;
                tip?: ({
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    tipper?: string | undefined;
                } & {
                    amount?: ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] & ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & Record<Exclude<keyof I["tx"]["authInfo"]["tip"]["amount"][number], keyof import("../../base/v1beta1/coin").Coin>, never>)[] & Record<Exclude<keyof I["tx"]["authInfo"]["tip"]["amount"], keyof {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[]>, never>) | undefined;
                    tipper?: string | undefined;
                } & Record<Exclude<keyof I["tx"]["authInfo"]["tip"], keyof import("./tx").Tip>, never>) | undefined;
            } & Record<Exclude<keyof I["tx"]["authInfo"], keyof import("./tx").AuthInfo>, never>) | undefined;
            signatures?: (Uint8Array[] & Uint8Array[] & Record<Exclude<keyof I["tx"]["signatures"], keyof Uint8Array[]>, never>) | undefined;
        } & Record<Exclude<keyof I["tx"], keyof Tx>, never>) | undefined;
        txResponse?: ({
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
                    } & Record<Exclude<keyof I["txResponse"]["logs"][number]["events"][number]["attributes"][number], keyof import("../../base/abci/v1beta1/abci").Attribute>, never>)[] & Record<Exclude<keyof I["txResponse"]["logs"][number]["events"][number]["attributes"], keyof {
                        key?: string | undefined;
                        value?: string | undefined;
                    }[]>, never>) | undefined;
                } & Record<Exclude<keyof I["txResponse"]["logs"][number]["events"][number], keyof import("../../base/abci/v1beta1/abci").StringEvent>, never>)[] & Record<Exclude<keyof I["txResponse"]["logs"][number]["events"], keyof {
                    type?: string | undefined;
                    attributes?: {
                        key?: string | undefined;
                        value?: string | undefined;
                    }[] | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["txResponse"]["logs"][number], keyof import("../../base/abci/v1beta1/abci").ABCIMessageLog>, never>)[] & Record<Exclude<keyof I["txResponse"]["logs"], keyof {
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
            } & Record<Exclude<keyof I["txResponse"]["tx"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
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
                } & Record<Exclude<keyof I["txResponse"]["events"][number]["attributes"][number], keyof import("../../../tendermint/abci/types").EventAttribute>, never>)[] & Record<Exclude<keyof I["txResponse"]["events"][number]["attributes"], keyof {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["txResponse"]["events"][number], keyof import("../../../tendermint/abci/types").Event>, never>)[] & Record<Exclude<keyof I["txResponse"]["events"], keyof {
                type?: string | undefined;
                attributes?: {
                    key?: string | undefined;
                    value?: string | undefined;
                    index?: boolean | undefined;
                }[] | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["txResponse"], keyof TxResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof GetTxResponse>, never>>(object: I): GetTxResponse;
};
export declare const GetBlockWithTxsRequest: {
    typeUrl: string;
    encode(message: GetBlockWithTxsRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GetBlockWithTxsRequest;
    fromJSON(object: any): GetBlockWithTxsRequest;
    toJSON(message: GetBlockWithTxsRequest): unknown;
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
    } & Record<Exclude<keyof I, keyof GetBlockWithTxsRequest>, never>>(object: I): GetBlockWithTxsRequest;
};
export declare const GetBlockWithTxsResponse: {
    typeUrl: string;
    encode(message: GetBlockWithTxsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GetBlockWithTxsResponse;
    fromJSON(object: any): GetBlockWithTxsResponse;
    toJSON(message: GetBlockWithTxsResponse): unknown;
    fromPartial<I extends {
        txs?: {
            body?: {
                messages?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                memo?: string | undefined;
                timeoutHeight?: bigint | undefined;
                extensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                nonCriticalExtensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
            } | undefined;
            authInfo?: {
                signerInfos?: {
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[] | undefined;
                fee?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } | undefined;
                tip?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    tipper?: string | undefined;
                } | undefined;
            } | undefined;
            signatures?: Uint8Array[] | undefined;
        }[] | undefined;
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
                            type?: import("../../../tendermint/types/types").SignedMsgType | undefined;
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
                            type?: import("../../../tendermint/types/types").SignedMsgType | undefined;
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
                                        blockIdFlag?: import("../../../tendermint/types/types").BlockIDFlag | undefined;
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
                    blockIdFlag?: import("../../../tendermint/types/types").BlockIDFlag | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    timestamp?: {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } | undefined;
                    signature?: Uint8Array | undefined;
                }[] | undefined;
            } | undefined;
        } | undefined;
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
        txs?: ({
            body?: {
                messages?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                memo?: string | undefined;
                timeoutHeight?: bigint | undefined;
                extensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                nonCriticalExtensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
            } | undefined;
            authInfo?: {
                signerInfos?: {
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[] | undefined;
                fee?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } | undefined;
                tip?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    tipper?: string | undefined;
                } | undefined;
            } | undefined;
            signatures?: Uint8Array[] | undefined;
        }[] & ({
            body?: {
                messages?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                memo?: string | undefined;
                timeoutHeight?: bigint | undefined;
                extensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                nonCriticalExtensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
            } | undefined;
            authInfo?: {
                signerInfos?: {
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[] | undefined;
                fee?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } | undefined;
                tip?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    tipper?: string | undefined;
                } | undefined;
            } | undefined;
            signatures?: Uint8Array[] | undefined;
        } & {
            body?: ({
                messages?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                memo?: string | undefined;
                timeoutHeight?: bigint | undefined;
                extensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                nonCriticalExtensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
            } & {
                messages?: ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] & ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["txs"][number]["body"]["messages"][number], keyof import("../../../google/protobuf/any").Any>, never>)[] & Record<Exclude<keyof I["txs"][number]["body"]["messages"], keyof {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[]>, never>) | undefined;
                memo?: string | undefined;
                timeoutHeight?: bigint | undefined;
                extensionOptions?: ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] & ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["txs"][number]["body"]["extensionOptions"][number], keyof import("../../../google/protobuf/any").Any>, never>)[] & Record<Exclude<keyof I["txs"][number]["body"]["extensionOptions"], keyof {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[]>, never>) | undefined;
                nonCriticalExtensionOptions?: ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] & ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["txs"][number]["body"]["nonCriticalExtensionOptions"][number], keyof import("../../../google/protobuf/any").Any>, never>)[] & Record<Exclude<keyof I["txs"][number]["body"]["nonCriticalExtensionOptions"], keyof {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["txs"][number]["body"], keyof import("./tx").TxBody>, never>) | undefined;
            authInfo?: ({
                signerInfos?: {
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[] | undefined;
                fee?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } | undefined;
                tip?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    tipper?: string | undefined;
                } | undefined;
            } & {
                signerInfos?: ({
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[] & ({
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                } & {
                    publicKey?: ({
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } & {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["publicKey"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
                    modeInfo?: ({
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } & {
                        single?: ({
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } & {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["single"], "mode">, never>) | undefined;
                        multi?: ({
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: {
                                single?: {
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } | undefined;
                                multi?: any | undefined;
                            }[] | undefined;
                        } & {
                            bitarray?: ({
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } & {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["bitarray"], keyof import("../../crypto/multisig/v1beta1/multisig").CompactBitArray>, never>) | undefined;
                            modeInfos?: ({
                                single?: {
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } | undefined;
                                multi?: {
                                    bitarray?: {
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } | undefined;
                                    modeInfos?: any[] | undefined;
                                } | undefined;
                            }[] & ({
                                single?: {
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } | undefined;
                                multi?: {
                                    bitarray?: {
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } | undefined;
                                    modeInfos?: any[] | undefined;
                                } | undefined;
                            } & {
                                single?: ({
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } & {
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["single"], "mode">, never>) | undefined;
                                multi?: ({
                                    bitarray?: {
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } | undefined;
                                    modeInfos?: {
                                        single?: {
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } | undefined;
                                        multi?: any | undefined;
                                    }[] | undefined;
                                } & {
                                    bitarray?: ({
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } & {
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["bitarray"], keyof import("../../crypto/multisig/v1beta1/multisig").CompactBitArray>, never>) | undefined;
                                    modeInfos?: ({
                                        single?: {
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } | undefined;
                                        multi?: {
                                            bitarray?: {
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } | undefined;
                                            modeInfos?: any[] | undefined;
                                        } | undefined;
                                    }[] & ({
                                        single?: {
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } | undefined;
                                        multi?: {
                                            bitarray?: {
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } | undefined;
                                            modeInfos?: any[] | undefined;
                                        } | undefined;
                                    } & {
                                        single?: ({
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } & {
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["single"], "mode">, never>) | undefined;
                                        multi?: ({
                                            bitarray?: {
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } | undefined;
                                            modeInfos?: {
                                                single?: {
                                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                } | undefined;
                                                multi?: any | undefined;
                                            }[] | undefined;
                                        } & {
                                            bitarray?: ({
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } & {
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"]["bitarray"], keyof import("../../crypto/multisig/v1beta1/multisig").CompactBitArray>, never>) | undefined;
                                            modeInfos?: ({
                                                single?: {
                                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                } | undefined;
                                                multi?: {
                                                    bitarray?: {
                                                        extraBitsStored?: number | undefined;
                                                        elems?: Uint8Array | undefined;
                                                    } | undefined;
                                                    modeInfos?: any[] | undefined;
                                                } | undefined;
                                            }[] & ({
                                                single?: {
                                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                } | undefined;
                                                multi?: {
                                                    bitarray?: {
                                                        extraBitsStored?: number | undefined;
                                                        elems?: Uint8Array | undefined;
                                                    } | undefined;
                                                    modeInfos?: any[] | undefined;
                                                } | undefined;
                                            } & {
                                                single?: ({
                                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                } & any & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["single"], "mode">, never>) | undefined;
                                                multi?: ({
                                                    bitarray?: {
                                                        extraBitsStored?: number | undefined;
                                                        elems?: Uint8Array | undefined;
                                                    } | undefined;
                                                    modeInfos?: {
                                                        single?: {
                                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                        } | undefined;
                                                        multi?: any | undefined;
                                                    }[] | undefined;
                                                } & any & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"], keyof import("./tx").ModeInfo_Multi>, never>) | undefined;
                                            } & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number], keyof import("./tx").ModeInfo>, never>)[] & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"]["modeInfos"], keyof {
                                                single?: {
                                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                } | undefined;
                                                multi?: {
                                                    bitarray?: {
                                                        extraBitsStored?: number | undefined;
                                                        elems?: Uint8Array | undefined;
                                                    } | undefined;
                                                    modeInfos?: any[] | undefined;
                                                } | undefined;
                                            }[]>, never>) | undefined;
                                        } & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"], keyof import("./tx").ModeInfo_Multi>, never>) | undefined;
                                    } & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number], keyof import("./tx").ModeInfo>, never>)[] & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"], keyof {
                                        single?: {
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } | undefined;
                                        multi?: {
                                            bitarray?: {
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } | undefined;
                                            modeInfos?: any[] | undefined;
                                        } | undefined;
                                    }[]>, never>) | undefined;
                                } & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"], keyof import("./tx").ModeInfo_Multi>, never>) | undefined;
                            } & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number], keyof import("./tx").ModeInfo>, never>)[] & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"], keyof {
                                single?: {
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } | undefined;
                                multi?: {
                                    bitarray?: {
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } | undefined;
                                    modeInfos?: any[] | undefined;
                                } | undefined;
                            }[]>, never>) | undefined;
                        } & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"], keyof import("./tx").ModeInfo_Multi>, never>) | undefined;
                    } & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number]["modeInfo"], keyof import("./tx").ModeInfo>, never>) | undefined;
                    sequence?: bigint | undefined;
                } & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"][number], keyof import("./tx").SignerInfo>, never>)[] & Record<Exclude<keyof I["txs"][number]["authInfo"]["signerInfos"], keyof {
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[]>, never>) | undefined;
                fee?: ({
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } & {
                    amount?: ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] & ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & Record<Exclude<keyof I["txs"][number]["authInfo"]["fee"]["amount"][number], keyof import("../../base/v1beta1/coin").Coin>, never>)[] & Record<Exclude<keyof I["txs"][number]["authInfo"]["fee"]["amount"], keyof {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[]>, never>) | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } & Record<Exclude<keyof I["txs"][number]["authInfo"]["fee"], keyof import("./tx").Fee>, never>) | undefined;
                tip?: ({
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    tipper?: string | undefined;
                } & {
                    amount?: ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] & ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & Record<Exclude<keyof I["txs"][number]["authInfo"]["tip"]["amount"][number], keyof import("../../base/v1beta1/coin").Coin>, never>)[] & Record<Exclude<keyof I["txs"][number]["authInfo"]["tip"]["amount"], keyof {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[]>, never>) | undefined;
                    tipper?: string | undefined;
                } & Record<Exclude<keyof I["txs"][number]["authInfo"]["tip"], keyof import("./tx").Tip>, never>) | undefined;
            } & Record<Exclude<keyof I["txs"][number]["authInfo"], keyof import("./tx").AuthInfo>, never>) | undefined;
            signatures?: (Uint8Array[] & Uint8Array[] & Record<Exclude<keyof I["txs"][number]["signatures"], keyof Uint8Array[]>, never>) | undefined;
        } & Record<Exclude<keyof I["txs"][number], keyof Tx>, never>)[] & Record<Exclude<keyof I["txs"], keyof {
            body?: {
                messages?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                memo?: string | undefined;
                timeoutHeight?: bigint | undefined;
                extensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                nonCriticalExtensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
            } | undefined;
            authInfo?: {
                signerInfos?: {
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[] | undefined;
                fee?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } | undefined;
                tip?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    tipper?: string | undefined;
                } | undefined;
            } | undefined;
            signatures?: Uint8Array[] | undefined;
        }[]>, never>) | undefined;
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
            } & Record<Exclude<keyof I["blockId"]["partSetHeader"], keyof import("../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
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
                            type?: import("../../../tendermint/types/types").SignedMsgType | undefined;
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
                            type?: import("../../../tendermint/types/types").SignedMsgType | undefined;
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
                                        blockIdFlag?: import("../../../tendermint/types/types").BlockIDFlag | undefined;
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
                    blockIdFlag?: import("../../../tendermint/types/types").BlockIDFlag | undefined;
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
                } & Record<Exclude<keyof I["block"]["header"]["version"], keyof import("../../../tendermint/version/types").Consensus>, never>) | undefined;
                chainId?: string | undefined;
                height?: bigint | undefined;
                time?: ({
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & Record<Exclude<keyof I["block"]["header"]["time"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
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
                    } & Record<Exclude<keyof I["block"]["header"]["lastBlockId"]["partSetHeader"], keyof import("../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
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
            } & Record<Exclude<keyof I["block"]["header"], keyof import("../../../tendermint/types/types").Header>, never>) | undefined;
            data?: ({
                txs?: Uint8Array[] | undefined;
            } & {
                txs?: (Uint8Array[] & Uint8Array[] & Record<Exclude<keyof I["block"]["data"]["txs"], keyof Uint8Array[]>, never>) | undefined;
            } & Record<Exclude<keyof I["block"]["data"], "txs">, never>) | undefined;
            evidence?: ({
                evidence?: {
                    duplicateVoteEvidence?: {
                        voteA?: {
                            type?: import("../../../tendermint/types/types").SignedMsgType | undefined;
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
                            type?: import("../../../tendermint/types/types").SignedMsgType | undefined;
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
                                        blockIdFlag?: import("../../../tendermint/types/types").BlockIDFlag | undefined;
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
                            type?: import("../../../tendermint/types/types").SignedMsgType | undefined;
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
                            type?: import("../../../tendermint/types/types").SignedMsgType | undefined;
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
                                        blockIdFlag?: import("../../../tendermint/types/types").BlockIDFlag | undefined;
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
                            type?: import("../../../tendermint/types/types").SignedMsgType | undefined;
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
                            type?: import("../../../tendermint/types/types").SignedMsgType | undefined;
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
                                        blockIdFlag?: import("../../../tendermint/types/types").BlockIDFlag | undefined;
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
                            type?: import("../../../tendermint/types/types").SignedMsgType | undefined;
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
                            type?: import("../../../tendermint/types/types").SignedMsgType | undefined;
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
                            type?: import("../../../tendermint/types/types").SignedMsgType | undefined;
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
                            type?: import("../../../tendermint/types/types").SignedMsgType | undefined;
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
                                } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteA"]["blockId"]["partSetHeader"], keyof import("../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                            } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteA"]["blockId"], keyof BlockID>, never>) | undefined;
                            timestamp?: ({
                                seconds?: bigint | undefined;
                                nanos?: number | undefined;
                            } & {
                                seconds?: bigint | undefined;
                                nanos?: number | undefined;
                            } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteA"]["timestamp"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                            validatorAddress?: Uint8Array | undefined;
                            validatorIndex?: number | undefined;
                            signature?: Uint8Array | undefined;
                        } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteA"], keyof import("../../../tendermint/types/types").Vote>, never>) | undefined;
                        voteB?: ({
                            type?: import("../../../tendermint/types/types").SignedMsgType | undefined;
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
                            type?: import("../../../tendermint/types/types").SignedMsgType | undefined;
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
                                } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteB"]["blockId"]["partSetHeader"], keyof import("../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                            } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteB"]["blockId"], keyof BlockID>, never>) | undefined;
                            timestamp?: ({
                                seconds?: bigint | undefined;
                                nanos?: number | undefined;
                            } & {
                                seconds?: bigint | undefined;
                                nanos?: number | undefined;
                            } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteB"]["timestamp"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                            validatorAddress?: Uint8Array | undefined;
                            validatorIndex?: number | undefined;
                            signature?: Uint8Array | undefined;
                        } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["voteB"], keyof import("../../../tendermint/types/types").Vote>, never>) | undefined;
                        totalVotingPower?: bigint | undefined;
                        validatorPower?: bigint | undefined;
                        timestamp?: ({
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } & {
                            seconds?: bigint | undefined;
                            nanos?: number | undefined;
                        } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["duplicateVoteEvidence"]["timestamp"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                    } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["duplicateVoteEvidence"], keyof import("../../../tendermint/types/evidence").DuplicateVoteEvidence>, never>) | undefined;
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
                                        blockIdFlag?: import("../../../tendermint/types/types").BlockIDFlag | undefined;
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
                                        blockIdFlag?: import("../../../tendermint/types/types").BlockIDFlag | undefined;
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
                                        blockIdFlag?: import("../../../tendermint/types/types").BlockIDFlag | undefined;
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
                                    } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"]["version"], keyof import("../../../tendermint/version/types").Consensus>, never>) | undefined;
                                    chainId?: string | undefined;
                                    height?: bigint | undefined;
                                    time?: ({
                                        seconds?: bigint | undefined;
                                        nanos?: number | undefined;
                                    } & {
                                        seconds?: bigint | undefined;
                                        nanos?: number | undefined;
                                    } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"]["time"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
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
                                        } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"]["lastBlockId"]["partSetHeader"], keyof import("../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
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
                                } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["header"], keyof import("../../../tendermint/types/types").Header>, never>) | undefined;
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
                                        blockIdFlag?: import("../../../tendermint/types/types").BlockIDFlag | undefined;
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
                                        } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["blockId"]["partSetHeader"], keyof import("../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                                    } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["blockId"], keyof BlockID>, never>) | undefined;
                                    signatures?: ({
                                        blockIdFlag?: import("../../../tendermint/types/types").BlockIDFlag | undefined;
                                        validatorAddress?: Uint8Array | undefined;
                                        timestamp?: {
                                            seconds?: bigint | undefined;
                                            nanos?: number | undefined;
                                        } | undefined;
                                        signature?: Uint8Array | undefined;
                                    }[] & ({
                                        blockIdFlag?: import("../../../tendermint/types/types").BlockIDFlag | undefined;
                                        validatorAddress?: Uint8Array | undefined;
                                        timestamp?: {
                                            seconds?: bigint | undefined;
                                            nanos?: number | undefined;
                                        } | undefined;
                                        signature?: Uint8Array | undefined;
                                    } & {
                                        blockIdFlag?: import("../../../tendermint/types/types").BlockIDFlag | undefined;
                                        validatorAddress?: Uint8Array | undefined;
                                        timestamp?: ({
                                            seconds?: bigint | undefined;
                                            nanos?: number | undefined;
                                        } & {
                                            seconds?: bigint | undefined;
                                            nanos?: number | undefined;
                                        } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["signatures"][number]["timestamp"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                                        signature?: Uint8Array | undefined;
                                    } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["signatures"][number], keyof import("../../../tendermint/types/types").CommitSig>, never>)[] & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"]["signatures"], keyof {
                                        blockIdFlag?: import("../../../tendermint/types/types").BlockIDFlag | undefined;
                                        validatorAddress?: Uint8Array | undefined;
                                        timestamp?: {
                                            seconds?: bigint | undefined;
                                            nanos?: number | undefined;
                                        } | undefined;
                                        signature?: Uint8Array | undefined;
                                    }[]>, never>) | undefined;
                                } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"]["commit"], keyof import("../../../tendermint/types/types").Commit>, never>) | undefined;
                            } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["signedHeader"], keyof import("../../../tendermint/types/types").SignedHeader>, never>) | undefined;
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
                                    } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["validators"][number]["pubKey"], keyof import("../../../tendermint/crypto/keys").PublicKey>, never>) | undefined;
                                    votingPower?: bigint | undefined;
                                    proposerPriority?: bigint | undefined;
                                } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["validators"][number], keyof import("../../../tendermint/types/validator").Validator>, never>)[] & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["validators"], keyof {
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
                                    } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["proposer"]["pubKey"], keyof import("../../../tendermint/crypto/keys").PublicKey>, never>) | undefined;
                                    votingPower?: bigint | undefined;
                                    proposerPriority?: bigint | undefined;
                                } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"]["proposer"], keyof import("../../../tendermint/types/validator").Validator>, never>) | undefined;
                                totalVotingPower?: bigint | undefined;
                            } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"]["validatorSet"], keyof import("../../../tendermint/types/validator").ValidatorSet>, never>) | undefined;
                        } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["conflictingBlock"], keyof import("../../../tendermint/types/types").LightBlock>, never>) | undefined;
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
                            } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["byzantineValidators"][number]["pubKey"], keyof import("../../../tendermint/crypto/keys").PublicKey>, never>) | undefined;
                            votingPower?: bigint | undefined;
                            proposerPriority?: bigint | undefined;
                        } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["byzantineValidators"][number], keyof import("../../../tendermint/types/validator").Validator>, never>)[] & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["byzantineValidators"], keyof {
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
                        } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"]["timestamp"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                    } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number]["lightClientAttackEvidence"], keyof import("../../../tendermint/types/evidence").LightClientAttackEvidence>, never>) | undefined;
                } & Record<Exclude<keyof I["block"]["evidence"]["evidence"][number], keyof import("../../../tendermint/types/evidence").Evidence>, never>)[] & Record<Exclude<keyof I["block"]["evidence"]["evidence"], keyof {
                    duplicateVoteEvidence?: {
                        voteA?: {
                            type?: import("../../../tendermint/types/types").SignedMsgType | undefined;
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
                            type?: import("../../../tendermint/types/types").SignedMsgType | undefined;
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
                                        blockIdFlag?: import("../../../tendermint/types/types").BlockIDFlag | undefined;
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
                    blockIdFlag?: import("../../../tendermint/types/types").BlockIDFlag | undefined;
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
                    } & Record<Exclude<keyof I["block"]["lastCommit"]["blockId"]["partSetHeader"], keyof import("../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                } & Record<Exclude<keyof I["block"]["lastCommit"]["blockId"], keyof BlockID>, never>) | undefined;
                signatures?: ({
                    blockIdFlag?: import("../../../tendermint/types/types").BlockIDFlag | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    timestamp?: {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } | undefined;
                    signature?: Uint8Array | undefined;
                }[] & ({
                    blockIdFlag?: import("../../../tendermint/types/types").BlockIDFlag | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    timestamp?: {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } | undefined;
                    signature?: Uint8Array | undefined;
                } & {
                    blockIdFlag?: import("../../../tendermint/types/types").BlockIDFlag | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    timestamp?: ({
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } & {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } & Record<Exclude<keyof I["block"]["lastCommit"]["signatures"][number]["timestamp"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                    signature?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["block"]["lastCommit"]["signatures"][number], keyof import("../../../tendermint/types/types").CommitSig>, never>)[] & Record<Exclude<keyof I["block"]["lastCommit"]["signatures"], keyof {
                    blockIdFlag?: import("../../../tendermint/types/types").BlockIDFlag | undefined;
                    validatorAddress?: Uint8Array | undefined;
                    timestamp?: {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } | undefined;
                    signature?: Uint8Array | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["block"]["lastCommit"], keyof import("../../../tendermint/types/types").Commit>, never>) | undefined;
        } & Record<Exclude<keyof I["block"], keyof Block>, never>) | undefined;
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof GetBlockWithTxsResponse>, never>>(object: I): GetBlockWithTxsResponse;
};
export declare const TxDecodeRequest: {
    typeUrl: string;
    encode(message: TxDecodeRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): TxDecodeRequest;
    fromJSON(object: any): TxDecodeRequest;
    toJSON(message: TxDecodeRequest): unknown;
    fromPartial<I extends {
        txBytes?: Uint8Array | undefined;
    } & {
        txBytes?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, "txBytes">, never>>(object: I): TxDecodeRequest;
};
export declare const TxDecodeResponse: {
    typeUrl: string;
    encode(message: TxDecodeResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): TxDecodeResponse;
    fromJSON(object: any): TxDecodeResponse;
    toJSON(message: TxDecodeResponse): unknown;
    fromPartial<I extends {
        tx?: {
            body?: {
                messages?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                memo?: string | undefined;
                timeoutHeight?: bigint | undefined;
                extensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                nonCriticalExtensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
            } | undefined;
            authInfo?: {
                signerInfos?: {
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[] | undefined;
                fee?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } | undefined;
                tip?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    tipper?: string | undefined;
                } | undefined;
            } | undefined;
            signatures?: Uint8Array[] | undefined;
        } | undefined;
    } & {
        tx?: ({
            body?: {
                messages?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                memo?: string | undefined;
                timeoutHeight?: bigint | undefined;
                extensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                nonCriticalExtensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
            } | undefined;
            authInfo?: {
                signerInfos?: {
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[] | undefined;
                fee?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } | undefined;
                tip?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    tipper?: string | undefined;
                } | undefined;
            } | undefined;
            signatures?: Uint8Array[] | undefined;
        } & {
            body?: ({
                messages?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                memo?: string | undefined;
                timeoutHeight?: bigint | undefined;
                extensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                nonCriticalExtensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
            } & {
                messages?: ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] & ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["tx"]["body"]["messages"][number], keyof import("../../../google/protobuf/any").Any>, never>)[] & Record<Exclude<keyof I["tx"]["body"]["messages"], keyof {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[]>, never>) | undefined;
                memo?: string | undefined;
                timeoutHeight?: bigint | undefined;
                extensionOptions?: ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] & ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["tx"]["body"]["extensionOptions"][number], keyof import("../../../google/protobuf/any").Any>, never>)[] & Record<Exclude<keyof I["tx"]["body"]["extensionOptions"], keyof {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[]>, never>) | undefined;
                nonCriticalExtensionOptions?: ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] & ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["tx"]["body"]["nonCriticalExtensionOptions"][number], keyof import("../../../google/protobuf/any").Any>, never>)[] & Record<Exclude<keyof I["tx"]["body"]["nonCriticalExtensionOptions"], keyof {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["tx"]["body"], keyof import("./tx").TxBody>, never>) | undefined;
            authInfo?: ({
                signerInfos?: {
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[] | undefined;
                fee?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } | undefined;
                tip?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    tipper?: string | undefined;
                } | undefined;
            } & {
                signerInfos?: ({
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[] & ({
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                } & {
                    publicKey?: ({
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } & {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["publicKey"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
                    modeInfo?: ({
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } & {
                        single?: ({
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } & {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["single"], "mode">, never>) | undefined;
                        multi?: ({
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: {
                                single?: {
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } | undefined;
                                multi?: any | undefined;
                            }[] | undefined;
                        } & {
                            bitarray?: ({
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } & {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["bitarray"], keyof import("../../crypto/multisig/v1beta1/multisig").CompactBitArray>, never>) | undefined;
                            modeInfos?: ({
                                single?: {
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } | undefined;
                                multi?: {
                                    bitarray?: {
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } | undefined;
                                    modeInfos?: any[] | undefined;
                                } | undefined;
                            }[] & ({
                                single?: {
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } | undefined;
                                multi?: {
                                    bitarray?: {
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } | undefined;
                                    modeInfos?: any[] | undefined;
                                } | undefined;
                            } & {
                                single?: ({
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } & {
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["single"], "mode">, never>) | undefined;
                                multi?: ({
                                    bitarray?: {
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } | undefined;
                                    modeInfos?: {
                                        single?: {
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } | undefined;
                                        multi?: any | undefined;
                                    }[] | undefined;
                                } & {
                                    bitarray?: ({
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } & {
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["bitarray"], keyof import("../../crypto/multisig/v1beta1/multisig").CompactBitArray>, never>) | undefined;
                                    modeInfos?: ({
                                        single?: {
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } | undefined;
                                        multi?: {
                                            bitarray?: {
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } | undefined;
                                            modeInfos?: any[] | undefined;
                                        } | undefined;
                                    }[] & ({
                                        single?: {
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } | undefined;
                                        multi?: {
                                            bitarray?: {
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } | undefined;
                                            modeInfos?: any[] | undefined;
                                        } | undefined;
                                    } & {
                                        single?: ({
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } & {
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["single"], "mode">, never>) | undefined;
                                        multi?: ({
                                            bitarray?: {
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } | undefined;
                                            modeInfos?: {
                                                single?: {
                                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                } | undefined;
                                                multi?: any | undefined;
                                            }[] | undefined;
                                        } & {
                                            bitarray?: ({
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } & {
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"]["bitarray"], keyof import("../../crypto/multisig/v1beta1/multisig").CompactBitArray>, never>) | undefined;
                                            modeInfos?: ({
                                                single?: {
                                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                } | undefined;
                                                multi?: {
                                                    bitarray?: {
                                                        extraBitsStored?: number | undefined;
                                                        elems?: Uint8Array | undefined;
                                                    } | undefined;
                                                    modeInfos?: any[] | undefined;
                                                } | undefined;
                                            }[] & ({
                                                single?: {
                                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                } | undefined;
                                                multi?: {
                                                    bitarray?: {
                                                        extraBitsStored?: number | undefined;
                                                        elems?: Uint8Array | undefined;
                                                    } | undefined;
                                                    modeInfos?: any[] | undefined;
                                                } | undefined;
                                            } & {
                                                single?: ({
                                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                } & any & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["single"], "mode">, never>) | undefined;
                                                multi?: ({
                                                    bitarray?: {
                                                        extraBitsStored?: number | undefined;
                                                        elems?: Uint8Array | undefined;
                                                    } | undefined;
                                                    modeInfos?: {
                                                        single?: {
                                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                        } | undefined;
                                                        multi?: any | undefined;
                                                    }[] | undefined;
                                                } & any & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"], keyof import("./tx").ModeInfo_Multi>, never>) | undefined;
                                            } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number], keyof import("./tx").ModeInfo>, never>)[] & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"]["modeInfos"], keyof {
                                                single?: {
                                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                } | undefined;
                                                multi?: {
                                                    bitarray?: {
                                                        extraBitsStored?: number | undefined;
                                                        elems?: Uint8Array | undefined;
                                                    } | undefined;
                                                    modeInfos?: any[] | undefined;
                                                } | undefined;
                                            }[]>, never>) | undefined;
                                        } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"], keyof import("./tx").ModeInfo_Multi>, never>) | undefined;
                                    } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number], keyof import("./tx").ModeInfo>, never>)[] & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"], keyof {
                                        single?: {
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } | undefined;
                                        multi?: {
                                            bitarray?: {
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } | undefined;
                                            modeInfos?: any[] | undefined;
                                        } | undefined;
                                    }[]>, never>) | undefined;
                                } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"], keyof import("./tx").ModeInfo_Multi>, never>) | undefined;
                            } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number], keyof import("./tx").ModeInfo>, never>)[] & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"], keyof {
                                single?: {
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } | undefined;
                                multi?: {
                                    bitarray?: {
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } | undefined;
                                    modeInfos?: any[] | undefined;
                                } | undefined;
                            }[]>, never>) | undefined;
                        } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"], keyof import("./tx").ModeInfo_Multi>, never>) | undefined;
                    } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"], keyof import("./tx").ModeInfo>, never>) | undefined;
                    sequence?: bigint | undefined;
                } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number], keyof import("./tx").SignerInfo>, never>)[] & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"], keyof {
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[]>, never>) | undefined;
                fee?: ({
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } & {
                    amount?: ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] & ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & Record<Exclude<keyof I["tx"]["authInfo"]["fee"]["amount"][number], keyof import("../../base/v1beta1/coin").Coin>, never>)[] & Record<Exclude<keyof I["tx"]["authInfo"]["fee"]["amount"], keyof {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[]>, never>) | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } & Record<Exclude<keyof I["tx"]["authInfo"]["fee"], keyof import("./tx").Fee>, never>) | undefined;
                tip?: ({
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    tipper?: string | undefined;
                } & {
                    amount?: ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] & ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & Record<Exclude<keyof I["tx"]["authInfo"]["tip"]["amount"][number], keyof import("../../base/v1beta1/coin").Coin>, never>)[] & Record<Exclude<keyof I["tx"]["authInfo"]["tip"]["amount"], keyof {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[]>, never>) | undefined;
                    tipper?: string | undefined;
                } & Record<Exclude<keyof I["tx"]["authInfo"]["tip"], keyof import("./tx").Tip>, never>) | undefined;
            } & Record<Exclude<keyof I["tx"]["authInfo"], keyof import("./tx").AuthInfo>, never>) | undefined;
            signatures?: (Uint8Array[] & Uint8Array[] & Record<Exclude<keyof I["tx"]["signatures"], keyof Uint8Array[]>, never>) | undefined;
        } & Record<Exclude<keyof I["tx"], keyof Tx>, never>) | undefined;
    } & Record<Exclude<keyof I, "tx">, never>>(object: I): TxDecodeResponse;
};
export declare const TxEncodeRequest: {
    typeUrl: string;
    encode(message: TxEncodeRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): TxEncodeRequest;
    fromJSON(object: any): TxEncodeRequest;
    toJSON(message: TxEncodeRequest): unknown;
    fromPartial<I extends {
        tx?: {
            body?: {
                messages?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                memo?: string | undefined;
                timeoutHeight?: bigint | undefined;
                extensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                nonCriticalExtensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
            } | undefined;
            authInfo?: {
                signerInfos?: {
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[] | undefined;
                fee?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } | undefined;
                tip?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    tipper?: string | undefined;
                } | undefined;
            } | undefined;
            signatures?: Uint8Array[] | undefined;
        } | undefined;
    } & {
        tx?: ({
            body?: {
                messages?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                memo?: string | undefined;
                timeoutHeight?: bigint | undefined;
                extensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                nonCriticalExtensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
            } | undefined;
            authInfo?: {
                signerInfos?: {
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[] | undefined;
                fee?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } | undefined;
                tip?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    tipper?: string | undefined;
                } | undefined;
            } | undefined;
            signatures?: Uint8Array[] | undefined;
        } & {
            body?: ({
                messages?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                memo?: string | undefined;
                timeoutHeight?: bigint | undefined;
                extensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
                nonCriticalExtensionOptions?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] | undefined;
            } & {
                messages?: ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] & ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["tx"]["body"]["messages"][number], keyof import("../../../google/protobuf/any").Any>, never>)[] & Record<Exclude<keyof I["tx"]["body"]["messages"], keyof {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[]>, never>) | undefined;
                memo?: string | undefined;
                timeoutHeight?: bigint | undefined;
                extensionOptions?: ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] & ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["tx"]["body"]["extensionOptions"][number], keyof import("../../../google/protobuf/any").Any>, never>)[] & Record<Exclude<keyof I["tx"]["body"]["extensionOptions"], keyof {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[]>, never>) | undefined;
                nonCriticalExtensionOptions?: ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[] & ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["tx"]["body"]["nonCriticalExtensionOptions"][number], keyof import("../../../google/protobuf/any").Any>, never>)[] & Record<Exclude<keyof I["tx"]["body"]["nonCriticalExtensionOptions"], keyof {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["tx"]["body"], keyof import("./tx").TxBody>, never>) | undefined;
            authInfo?: ({
                signerInfos?: {
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[] | undefined;
                fee?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } | undefined;
                tip?: {
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    tipper?: string | undefined;
                } | undefined;
            } & {
                signerInfos?: ({
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[] & ({
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                } & {
                    publicKey?: ({
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } & {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["publicKey"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
                    modeInfo?: ({
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } & {
                        single?: ({
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } & {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["single"], "mode">, never>) | undefined;
                        multi?: ({
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: {
                                single?: {
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } | undefined;
                                multi?: any | undefined;
                            }[] | undefined;
                        } & {
                            bitarray?: ({
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } & {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["bitarray"], keyof import("../../crypto/multisig/v1beta1/multisig").CompactBitArray>, never>) | undefined;
                            modeInfos?: ({
                                single?: {
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } | undefined;
                                multi?: {
                                    bitarray?: {
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } | undefined;
                                    modeInfos?: any[] | undefined;
                                } | undefined;
                            }[] & ({
                                single?: {
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } | undefined;
                                multi?: {
                                    bitarray?: {
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } | undefined;
                                    modeInfos?: any[] | undefined;
                                } | undefined;
                            } & {
                                single?: ({
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } & {
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["single"], "mode">, never>) | undefined;
                                multi?: ({
                                    bitarray?: {
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } | undefined;
                                    modeInfos?: {
                                        single?: {
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } | undefined;
                                        multi?: any | undefined;
                                    }[] | undefined;
                                } & {
                                    bitarray?: ({
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } & {
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["bitarray"], keyof import("../../crypto/multisig/v1beta1/multisig").CompactBitArray>, never>) | undefined;
                                    modeInfos?: ({
                                        single?: {
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } | undefined;
                                        multi?: {
                                            bitarray?: {
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } | undefined;
                                            modeInfos?: any[] | undefined;
                                        } | undefined;
                                    }[] & ({
                                        single?: {
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } | undefined;
                                        multi?: {
                                            bitarray?: {
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } | undefined;
                                            modeInfos?: any[] | undefined;
                                        } | undefined;
                                    } & {
                                        single?: ({
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } & {
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["single"], "mode">, never>) | undefined;
                                        multi?: ({
                                            bitarray?: {
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } | undefined;
                                            modeInfos?: {
                                                single?: {
                                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                } | undefined;
                                                multi?: any | undefined;
                                            }[] | undefined;
                                        } & {
                                            bitarray?: ({
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } & {
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"]["bitarray"], keyof import("../../crypto/multisig/v1beta1/multisig").CompactBitArray>, never>) | undefined;
                                            modeInfos?: ({
                                                single?: {
                                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                } | undefined;
                                                multi?: {
                                                    bitarray?: {
                                                        extraBitsStored?: number | undefined;
                                                        elems?: Uint8Array | undefined;
                                                    } | undefined;
                                                    modeInfos?: any[] | undefined;
                                                } | undefined;
                                            }[] & ({
                                                single?: {
                                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                } | undefined;
                                                multi?: {
                                                    bitarray?: {
                                                        extraBitsStored?: number | undefined;
                                                        elems?: Uint8Array | undefined;
                                                    } | undefined;
                                                    modeInfos?: any[] | undefined;
                                                } | undefined;
                                            } & {
                                                single?: ({
                                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                } & any & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["single"], "mode">, never>) | undefined;
                                                multi?: ({
                                                    bitarray?: {
                                                        extraBitsStored?: number | undefined;
                                                        elems?: Uint8Array | undefined;
                                                    } | undefined;
                                                    modeInfos?: {
                                                        single?: {
                                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                        } | undefined;
                                                        multi?: any | undefined;
                                                    }[] | undefined;
                                                } & any & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"], keyof import("./tx").ModeInfo_Multi>, never>) | undefined;
                                            } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number], keyof import("./tx").ModeInfo>, never>)[] & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"]["modeInfos"], keyof {
                                                single?: {
                                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                                } | undefined;
                                                multi?: {
                                                    bitarray?: {
                                                        extraBitsStored?: number | undefined;
                                                        elems?: Uint8Array | undefined;
                                                    } | undefined;
                                                    modeInfos?: any[] | undefined;
                                                } | undefined;
                                            }[]>, never>) | undefined;
                                        } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number]["multi"], keyof import("./tx").ModeInfo_Multi>, never>) | undefined;
                                    } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"][number], keyof import("./tx").ModeInfo>, never>)[] & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"]["modeInfos"], keyof {
                                        single?: {
                                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                        } | undefined;
                                        multi?: {
                                            bitarray?: {
                                                extraBitsStored?: number | undefined;
                                                elems?: Uint8Array | undefined;
                                            } | undefined;
                                            modeInfos?: any[] | undefined;
                                        } | undefined;
                                    }[]>, never>) | undefined;
                                } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number]["multi"], keyof import("./tx").ModeInfo_Multi>, never>) | undefined;
                            } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"][number], keyof import("./tx").ModeInfo>, never>)[] & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"]["modeInfos"], keyof {
                                single?: {
                                    mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                                } | undefined;
                                multi?: {
                                    bitarray?: {
                                        extraBitsStored?: number | undefined;
                                        elems?: Uint8Array | undefined;
                                    } | undefined;
                                    modeInfos?: any[] | undefined;
                                } | undefined;
                            }[]>, never>) | undefined;
                        } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"]["multi"], keyof import("./tx").ModeInfo_Multi>, never>) | undefined;
                    } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number]["modeInfo"], keyof import("./tx").ModeInfo>, never>) | undefined;
                    sequence?: bigint | undefined;
                } & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"][number], keyof import("./tx").SignerInfo>, never>)[] & Record<Exclude<keyof I["tx"]["authInfo"]["signerInfos"], keyof {
                    publicKey?: {
                        typeUrl?: string | undefined;
                        value?: Uint8Array | undefined;
                    } | undefined;
                    modeInfo?: {
                        single?: {
                            mode?: import("../signing/v1beta1/signing").SignMode | undefined;
                        } | undefined;
                        multi?: {
                            bitarray?: {
                                extraBitsStored?: number | undefined;
                                elems?: Uint8Array | undefined;
                            } | undefined;
                            modeInfos?: any[] | undefined;
                        } | undefined;
                    } | undefined;
                    sequence?: bigint | undefined;
                }[]>, never>) | undefined;
                fee?: ({
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } & {
                    amount?: ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] & ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & Record<Exclude<keyof I["tx"]["authInfo"]["fee"]["amount"][number], keyof import("../../base/v1beta1/coin").Coin>, never>)[] & Record<Exclude<keyof I["tx"]["authInfo"]["fee"]["amount"], keyof {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[]>, never>) | undefined;
                    gasLimit?: bigint | undefined;
                    payer?: string | undefined;
                    granter?: string | undefined;
                } & Record<Exclude<keyof I["tx"]["authInfo"]["fee"], keyof import("./tx").Fee>, never>) | undefined;
                tip?: ({
                    amount?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    tipper?: string | undefined;
                } & {
                    amount?: ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] & ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & Record<Exclude<keyof I["tx"]["authInfo"]["tip"]["amount"][number], keyof import("../../base/v1beta1/coin").Coin>, never>)[] & Record<Exclude<keyof I["tx"]["authInfo"]["tip"]["amount"], keyof {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[]>, never>) | undefined;
                    tipper?: string | undefined;
                } & Record<Exclude<keyof I["tx"]["authInfo"]["tip"], keyof import("./tx").Tip>, never>) | undefined;
            } & Record<Exclude<keyof I["tx"]["authInfo"], keyof import("./tx").AuthInfo>, never>) | undefined;
            signatures?: (Uint8Array[] & Uint8Array[] & Record<Exclude<keyof I["tx"]["signatures"], keyof Uint8Array[]>, never>) | undefined;
        } & Record<Exclude<keyof I["tx"], keyof Tx>, never>) | undefined;
    } & Record<Exclude<keyof I, "tx">, never>>(object: I): TxEncodeRequest;
};
export declare const TxEncodeResponse: {
    typeUrl: string;
    encode(message: TxEncodeResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): TxEncodeResponse;
    fromJSON(object: any): TxEncodeResponse;
    toJSON(message: TxEncodeResponse): unknown;
    fromPartial<I extends {
        txBytes?: Uint8Array | undefined;
    } & {
        txBytes?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, "txBytes">, never>>(object: I): TxEncodeResponse;
};
export declare const TxEncodeAminoRequest: {
    typeUrl: string;
    encode(message: TxEncodeAminoRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): TxEncodeAminoRequest;
    fromJSON(object: any): TxEncodeAminoRequest;
    toJSON(message: TxEncodeAminoRequest): unknown;
    fromPartial<I extends {
        aminoJson?: string | undefined;
    } & {
        aminoJson?: string | undefined;
    } & Record<Exclude<keyof I, "aminoJson">, never>>(object: I): TxEncodeAminoRequest;
};
export declare const TxEncodeAminoResponse: {
    typeUrl: string;
    encode(message: TxEncodeAminoResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): TxEncodeAminoResponse;
    fromJSON(object: any): TxEncodeAminoResponse;
    toJSON(message: TxEncodeAminoResponse): unknown;
    fromPartial<I extends {
        aminoBinary?: Uint8Array | undefined;
    } & {
        aminoBinary?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, "aminoBinary">, never>>(object: I): TxEncodeAminoResponse;
};
export declare const TxDecodeAminoRequest: {
    typeUrl: string;
    encode(message: TxDecodeAminoRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): TxDecodeAminoRequest;
    fromJSON(object: any): TxDecodeAminoRequest;
    toJSON(message: TxDecodeAminoRequest): unknown;
    fromPartial<I extends {
        aminoBinary?: Uint8Array | undefined;
    } & {
        aminoBinary?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, "aminoBinary">, never>>(object: I): TxDecodeAminoRequest;
};
export declare const TxDecodeAminoResponse: {
    typeUrl: string;
    encode(message: TxDecodeAminoResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): TxDecodeAminoResponse;
    fromJSON(object: any): TxDecodeAminoResponse;
    toJSON(message: TxDecodeAminoResponse): unknown;
    fromPartial<I extends {
        aminoJson?: string | undefined;
    } & {
        aminoJson?: string | undefined;
    } & Record<Exclude<keyof I, "aminoJson">, never>>(object: I): TxDecodeAminoResponse;
};
/** Service defines a gRPC service for interacting with transactions. */
export interface Service {
    /** Simulate simulates executing a transaction for estimating gas usage. */
    Simulate(request: SimulateRequest): Promise<SimulateResponse>;
    /** GetTx fetches a tx by hash. */
    GetTx(request: GetTxRequest): Promise<GetTxResponse>;
    /** BroadcastTx broadcast transaction. */
    BroadcastTx(request: BroadcastTxRequest): Promise<BroadcastTxResponse>;
    /** GetTxsEvent fetches txs by event. */
    GetTxsEvent(request: GetTxsEventRequest): Promise<GetTxsEventResponse>;
    /**
     * GetBlockWithTxs fetches a block with decoded txs.
     *
     * Since: cosmos-sdk 0.45.2
     */
    GetBlockWithTxs(request: GetBlockWithTxsRequest): Promise<GetBlockWithTxsResponse>;
    /**
     * TxDecode decodes the transaction.
     *
     * Since: cosmos-sdk 0.47
     */
    TxDecode(request: TxDecodeRequest): Promise<TxDecodeResponse>;
    /**
     * TxEncode encodes the transaction.
     *
     * Since: cosmos-sdk 0.47
     */
    TxEncode(request: TxEncodeRequest): Promise<TxEncodeResponse>;
    /**
     * TxEncodeAmino encodes an Amino transaction from JSON to encoded bytes.
     *
     * Since: cosmos-sdk 0.47
     */
    TxEncodeAmino(request: TxEncodeAminoRequest): Promise<TxEncodeAminoResponse>;
    /**
     * TxDecodeAmino decodes an Amino transaction from encoded bytes to JSON.
     *
     * Since: cosmos-sdk 0.47
     */
    TxDecodeAmino(request: TxDecodeAminoRequest): Promise<TxDecodeAminoResponse>;
}
export declare class ServiceClientImpl implements Service {
    private readonly rpc;
    constructor(rpc: Rpc);
    Simulate(request: SimulateRequest): Promise<SimulateResponse>;
    GetTx(request: GetTxRequest): Promise<GetTxResponse>;
    BroadcastTx(request: BroadcastTxRequest): Promise<BroadcastTxResponse>;
    GetTxsEvent(request: GetTxsEventRequest): Promise<GetTxsEventResponse>;
    GetBlockWithTxs(request: GetBlockWithTxsRequest): Promise<GetBlockWithTxsResponse>;
    TxDecode(request: TxDecodeRequest): Promise<TxDecodeResponse>;
    TxEncode(request: TxEncodeRequest): Promise<TxEncodeResponse>;
    TxEncodeAmino(request: TxEncodeAminoRequest): Promise<TxEncodeAminoResponse>;
    TxDecodeAmino(request: TxDecodeAminoRequest): Promise<TxDecodeAminoResponse>;
}
