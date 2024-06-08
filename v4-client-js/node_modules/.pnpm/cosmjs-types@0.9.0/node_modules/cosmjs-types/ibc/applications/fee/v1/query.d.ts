import { PageRequest } from "../../../../cosmos/base/query/v1beta1/pagination";
import { PacketId } from "../../../core/channel/v1/channel";
import { IdentifiedPacketFees } from "./fee";
import { Coin } from "../../../../cosmos/base/v1beta1/coin";
import { FeeEnabledChannel } from "./genesis";
import { BinaryReader, BinaryWriter } from "../../../../binary";
import { Rpc } from "../../../../helpers";
export declare const protobufPackage = "ibc.applications.fee.v1";
/** QueryIncentivizedPacketsRequest defines the request type for the IncentivizedPackets rpc */
export interface QueryIncentivizedPacketsRequest {
    /** pagination defines an optional pagination for the request. */
    pagination?: PageRequest;
    /** block height at which to query */
    queryHeight: bigint;
}
/** QueryIncentivizedPacketsResponse defines the response type for the IncentivizedPackets rpc */
export interface QueryIncentivizedPacketsResponse {
    /** list of identified fees for incentivized packets */
    incentivizedPackets: IdentifiedPacketFees[];
}
/** QueryIncentivizedPacketRequest defines the request type for the IncentivizedPacket rpc */
export interface QueryIncentivizedPacketRequest {
    /** unique packet identifier comprised of channel ID, port ID and sequence */
    packetId: PacketId;
    /** block height at which to query */
    queryHeight: bigint;
}
/** QueryIncentivizedPacketsResponse defines the response type for the IncentivizedPacket rpc */
export interface QueryIncentivizedPacketResponse {
    /** the identified fees for the incentivized packet */
    incentivizedPacket: IdentifiedPacketFees;
}
/**
 * QueryIncentivizedPacketsForChannelRequest defines the request type for querying for all incentivized packets
 * for a specific channel
 */
export interface QueryIncentivizedPacketsForChannelRequest {
    /** pagination defines an optional pagination for the request. */
    pagination?: PageRequest;
    portId: string;
    channelId: string;
    /** Height to query at */
    queryHeight: bigint;
}
/** QueryIncentivizedPacketsResponse defines the response type for the incentivized packets RPC */
export interface QueryIncentivizedPacketsForChannelResponse {
    /** Map of all incentivized_packets */
    incentivizedPackets: IdentifiedPacketFees[];
}
/** QueryTotalRecvFeesRequest defines the request type for the TotalRecvFees rpc */
export interface QueryTotalRecvFeesRequest {
    /** the packet identifier for the associated fees */
    packetId: PacketId;
}
/** QueryTotalRecvFeesResponse defines the response type for the TotalRecvFees rpc */
export interface QueryTotalRecvFeesResponse {
    /** the total packet receive fees */
    recvFees: Coin[];
}
/** QueryTotalAckFeesRequest defines the request type for the TotalAckFees rpc */
export interface QueryTotalAckFeesRequest {
    /** the packet identifier for the associated fees */
    packetId: PacketId;
}
/** QueryTotalAckFeesResponse defines the response type for the TotalAckFees rpc */
export interface QueryTotalAckFeesResponse {
    /** the total packet acknowledgement fees */
    ackFees: Coin[];
}
/** QueryTotalTimeoutFeesRequest defines the request type for the TotalTimeoutFees rpc */
export interface QueryTotalTimeoutFeesRequest {
    /** the packet identifier for the associated fees */
    packetId: PacketId;
}
/** QueryTotalTimeoutFeesResponse defines the response type for the TotalTimeoutFees rpc */
export interface QueryTotalTimeoutFeesResponse {
    /** the total packet timeout fees */
    timeoutFees: Coin[];
}
/** QueryPayeeRequest defines the request type for the Payee rpc */
export interface QueryPayeeRequest {
    /** unique channel identifier */
    channelId: string;
    /** the relayer address to which the distribution address is registered */
    relayer: string;
}
/** QueryPayeeResponse defines the response type for the Payee rpc */
export interface QueryPayeeResponse {
    /** the payee address to which packet fees are paid out */
    payeeAddress: string;
}
/** QueryCounterpartyPayeeRequest defines the request type for the CounterpartyPayee rpc */
export interface QueryCounterpartyPayeeRequest {
    /** unique channel identifier */
    channelId: string;
    /** the relayer address to which the counterparty is registered */
    relayer: string;
}
/** QueryCounterpartyPayeeResponse defines the response type for the CounterpartyPayee rpc */
export interface QueryCounterpartyPayeeResponse {
    /** the counterparty payee address used to compensate forward relaying */
    counterpartyPayee: string;
}
/** QueryFeeEnabledChannelsRequest defines the request type for the FeeEnabledChannels rpc */
export interface QueryFeeEnabledChannelsRequest {
    /** pagination defines an optional pagination for the request. */
    pagination?: PageRequest;
    /** block height at which to query */
    queryHeight: bigint;
}
/** QueryFeeEnabledChannelsResponse defines the response type for the FeeEnabledChannels rpc */
export interface QueryFeeEnabledChannelsResponse {
    /** list of fee enabled channels */
    feeEnabledChannels: FeeEnabledChannel[];
}
/** QueryFeeEnabledChannelRequest defines the request type for the FeeEnabledChannel rpc */
export interface QueryFeeEnabledChannelRequest {
    /** unique port identifier */
    portId: string;
    /** unique channel identifier */
    channelId: string;
}
/** QueryFeeEnabledChannelResponse defines the response type for the FeeEnabledChannel rpc */
export interface QueryFeeEnabledChannelResponse {
    /** boolean flag representing the fee enabled channel status */
    feeEnabled: boolean;
}
export declare const QueryIncentivizedPacketsRequest: {
    typeUrl: string;
    encode(message: QueryIncentivizedPacketsRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryIncentivizedPacketsRequest;
    fromJSON(object: any): QueryIncentivizedPacketsRequest;
    toJSON(message: QueryIncentivizedPacketsRequest): unknown;
    fromPartial<I extends {
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
        queryHeight?: bigint | undefined;
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
        queryHeight?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof QueryIncentivizedPacketsRequest>, never>>(object: I): QueryIncentivizedPacketsRequest;
};
export declare const QueryIncentivizedPacketsResponse: {
    typeUrl: string;
    encode(message: QueryIncentivizedPacketsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryIncentivizedPacketsResponse;
    fromJSON(object: any): QueryIncentivizedPacketsResponse;
    toJSON(message: QueryIncentivizedPacketsResponse): unknown;
    fromPartial<I extends {
        incentivizedPackets?: {
            packetId?: {
                portId?: string | undefined;
                channelId?: string | undefined;
                sequence?: bigint | undefined;
            } | undefined;
            packetFees?: {
                fee?: {
                    recvFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    ackFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    timeoutFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                } | undefined;
                refundAddress?: string | undefined;
                relayers?: string[] | undefined;
            }[] | undefined;
        }[] | undefined;
    } & {
        incentivizedPackets?: ({
            packetId?: {
                portId?: string | undefined;
                channelId?: string | undefined;
                sequence?: bigint | undefined;
            } | undefined;
            packetFees?: {
                fee?: {
                    recvFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    ackFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    timeoutFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                } | undefined;
                refundAddress?: string | undefined;
                relayers?: string[] | undefined;
            }[] | undefined;
        }[] & ({
            packetId?: {
                portId?: string | undefined;
                channelId?: string | undefined;
                sequence?: bigint | undefined;
            } | undefined;
            packetFees?: {
                fee?: {
                    recvFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    ackFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    timeoutFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                } | undefined;
                refundAddress?: string | undefined;
                relayers?: string[] | undefined;
            }[] | undefined;
        } & {
            packetId?: ({
                portId?: string | undefined;
                channelId?: string | undefined;
                sequence?: bigint | undefined;
            } & {
                portId?: string | undefined;
                channelId?: string | undefined;
                sequence?: bigint | undefined;
            } & Record<Exclude<keyof I["incentivizedPackets"][number]["packetId"], keyof PacketId>, never>) | undefined;
            packetFees?: ({
                fee?: {
                    recvFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    ackFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    timeoutFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                } | undefined;
                refundAddress?: string | undefined;
                relayers?: string[] | undefined;
            }[] & ({
                fee?: {
                    recvFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    ackFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    timeoutFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                } | undefined;
                refundAddress?: string | undefined;
                relayers?: string[] | undefined;
            } & {
                fee?: ({
                    recvFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    ackFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    timeoutFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                } & {
                    recvFee?: ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] & ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & Record<Exclude<keyof I["incentivizedPackets"][number]["packetFees"][number]["fee"]["recvFee"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["incentivizedPackets"][number]["packetFees"][number]["fee"]["recvFee"], keyof {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[]>, never>) | undefined;
                    ackFee?: ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] & ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & Record<Exclude<keyof I["incentivizedPackets"][number]["packetFees"][number]["fee"]["ackFee"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["incentivizedPackets"][number]["packetFees"][number]["fee"]["ackFee"], keyof {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[]>, never>) | undefined;
                    timeoutFee?: ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] & ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & Record<Exclude<keyof I["incentivizedPackets"][number]["packetFees"][number]["fee"]["timeoutFee"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["incentivizedPackets"][number]["packetFees"][number]["fee"]["timeoutFee"], keyof {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[]>, never>) | undefined;
                } & Record<Exclude<keyof I["incentivizedPackets"][number]["packetFees"][number]["fee"], keyof import("./fee").Fee>, never>) | undefined;
                refundAddress?: string | undefined;
                relayers?: (string[] & string[] & Record<Exclude<keyof I["incentivizedPackets"][number]["packetFees"][number]["relayers"], keyof string[]>, never>) | undefined;
            } & Record<Exclude<keyof I["incentivizedPackets"][number]["packetFees"][number], keyof import("./fee").PacketFee>, never>)[] & Record<Exclude<keyof I["incentivizedPackets"][number]["packetFees"], keyof {
                fee?: {
                    recvFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    ackFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    timeoutFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                } | undefined;
                refundAddress?: string | undefined;
                relayers?: string[] | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["incentivizedPackets"][number], keyof IdentifiedPacketFees>, never>)[] & Record<Exclude<keyof I["incentivizedPackets"], keyof {
            packetId?: {
                portId?: string | undefined;
                channelId?: string | undefined;
                sequence?: bigint | undefined;
            } | undefined;
            packetFees?: {
                fee?: {
                    recvFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    ackFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    timeoutFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                } | undefined;
                refundAddress?: string | undefined;
                relayers?: string[] | undefined;
            }[] | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "incentivizedPackets">, never>>(object: I): QueryIncentivizedPacketsResponse;
};
export declare const QueryIncentivizedPacketRequest: {
    typeUrl: string;
    encode(message: QueryIncentivizedPacketRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryIncentivizedPacketRequest;
    fromJSON(object: any): QueryIncentivizedPacketRequest;
    toJSON(message: QueryIncentivizedPacketRequest): unknown;
    fromPartial<I extends {
        packetId?: {
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
        } | undefined;
        queryHeight?: bigint | undefined;
    } & {
        packetId?: ({
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
        } & {
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
        } & Record<Exclude<keyof I["packetId"], keyof PacketId>, never>) | undefined;
        queryHeight?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof QueryIncentivizedPacketRequest>, never>>(object: I): QueryIncentivizedPacketRequest;
};
export declare const QueryIncentivizedPacketResponse: {
    typeUrl: string;
    encode(message: QueryIncentivizedPacketResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryIncentivizedPacketResponse;
    fromJSON(object: any): QueryIncentivizedPacketResponse;
    toJSON(message: QueryIncentivizedPacketResponse): unknown;
    fromPartial<I extends {
        incentivizedPacket?: {
            packetId?: {
                portId?: string | undefined;
                channelId?: string | undefined;
                sequence?: bigint | undefined;
            } | undefined;
            packetFees?: {
                fee?: {
                    recvFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    ackFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    timeoutFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                } | undefined;
                refundAddress?: string | undefined;
                relayers?: string[] | undefined;
            }[] | undefined;
        } | undefined;
    } & {
        incentivizedPacket?: ({
            packetId?: {
                portId?: string | undefined;
                channelId?: string | undefined;
                sequence?: bigint | undefined;
            } | undefined;
            packetFees?: {
                fee?: {
                    recvFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    ackFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    timeoutFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                } | undefined;
                refundAddress?: string | undefined;
                relayers?: string[] | undefined;
            }[] | undefined;
        } & {
            packetId?: ({
                portId?: string | undefined;
                channelId?: string | undefined;
                sequence?: bigint | undefined;
            } & {
                portId?: string | undefined;
                channelId?: string | undefined;
                sequence?: bigint | undefined;
            } & Record<Exclude<keyof I["incentivizedPacket"]["packetId"], keyof PacketId>, never>) | undefined;
            packetFees?: ({
                fee?: {
                    recvFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    ackFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    timeoutFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                } | undefined;
                refundAddress?: string | undefined;
                relayers?: string[] | undefined;
            }[] & ({
                fee?: {
                    recvFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    ackFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    timeoutFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                } | undefined;
                refundAddress?: string | undefined;
                relayers?: string[] | undefined;
            } & {
                fee?: ({
                    recvFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    ackFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    timeoutFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                } & {
                    recvFee?: ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] & ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & Record<Exclude<keyof I["incentivizedPacket"]["packetFees"][number]["fee"]["recvFee"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["incentivizedPacket"]["packetFees"][number]["fee"]["recvFee"], keyof {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[]>, never>) | undefined;
                    ackFee?: ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] & ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & Record<Exclude<keyof I["incentivizedPacket"]["packetFees"][number]["fee"]["ackFee"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["incentivizedPacket"]["packetFees"][number]["fee"]["ackFee"], keyof {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[]>, never>) | undefined;
                    timeoutFee?: ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] & ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & Record<Exclude<keyof I["incentivizedPacket"]["packetFees"][number]["fee"]["timeoutFee"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["incentivizedPacket"]["packetFees"][number]["fee"]["timeoutFee"], keyof {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[]>, never>) | undefined;
                } & Record<Exclude<keyof I["incentivizedPacket"]["packetFees"][number]["fee"], keyof import("./fee").Fee>, never>) | undefined;
                refundAddress?: string | undefined;
                relayers?: (string[] & string[] & Record<Exclude<keyof I["incentivizedPacket"]["packetFees"][number]["relayers"], keyof string[]>, never>) | undefined;
            } & Record<Exclude<keyof I["incentivizedPacket"]["packetFees"][number], keyof import("./fee").PacketFee>, never>)[] & Record<Exclude<keyof I["incentivizedPacket"]["packetFees"], keyof {
                fee?: {
                    recvFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    ackFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    timeoutFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                } | undefined;
                refundAddress?: string | undefined;
                relayers?: string[] | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["incentivizedPacket"], keyof IdentifiedPacketFees>, never>) | undefined;
    } & Record<Exclude<keyof I, "incentivizedPacket">, never>>(object: I): QueryIncentivizedPacketResponse;
};
export declare const QueryIncentivizedPacketsForChannelRequest: {
    typeUrl: string;
    encode(message: QueryIncentivizedPacketsForChannelRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryIncentivizedPacketsForChannelRequest;
    fromJSON(object: any): QueryIncentivizedPacketsForChannelRequest;
    toJSON(message: QueryIncentivizedPacketsForChannelRequest): unknown;
    fromPartial<I extends {
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
        portId?: string | undefined;
        channelId?: string | undefined;
        queryHeight?: bigint | undefined;
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
        portId?: string | undefined;
        channelId?: string | undefined;
        queryHeight?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof QueryIncentivizedPacketsForChannelRequest>, never>>(object: I): QueryIncentivizedPacketsForChannelRequest;
};
export declare const QueryIncentivizedPacketsForChannelResponse: {
    typeUrl: string;
    encode(message: QueryIncentivizedPacketsForChannelResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryIncentivizedPacketsForChannelResponse;
    fromJSON(object: any): QueryIncentivizedPacketsForChannelResponse;
    toJSON(message: QueryIncentivizedPacketsForChannelResponse): unknown;
    fromPartial<I extends {
        incentivizedPackets?: {
            packetId?: {
                portId?: string | undefined;
                channelId?: string | undefined;
                sequence?: bigint | undefined;
            } | undefined;
            packetFees?: {
                fee?: {
                    recvFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    ackFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    timeoutFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                } | undefined;
                refundAddress?: string | undefined;
                relayers?: string[] | undefined;
            }[] | undefined;
        }[] | undefined;
    } & {
        incentivizedPackets?: ({
            packetId?: {
                portId?: string | undefined;
                channelId?: string | undefined;
                sequence?: bigint | undefined;
            } | undefined;
            packetFees?: {
                fee?: {
                    recvFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    ackFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    timeoutFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                } | undefined;
                refundAddress?: string | undefined;
                relayers?: string[] | undefined;
            }[] | undefined;
        }[] & ({
            packetId?: {
                portId?: string | undefined;
                channelId?: string | undefined;
                sequence?: bigint | undefined;
            } | undefined;
            packetFees?: {
                fee?: {
                    recvFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    ackFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    timeoutFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                } | undefined;
                refundAddress?: string | undefined;
                relayers?: string[] | undefined;
            }[] | undefined;
        } & {
            packetId?: ({
                portId?: string | undefined;
                channelId?: string | undefined;
                sequence?: bigint | undefined;
            } & {
                portId?: string | undefined;
                channelId?: string | undefined;
                sequence?: bigint | undefined;
            } & Record<Exclude<keyof I["incentivizedPackets"][number]["packetId"], keyof PacketId>, never>) | undefined;
            packetFees?: ({
                fee?: {
                    recvFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    ackFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    timeoutFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                } | undefined;
                refundAddress?: string | undefined;
                relayers?: string[] | undefined;
            }[] & ({
                fee?: {
                    recvFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    ackFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    timeoutFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                } | undefined;
                refundAddress?: string | undefined;
                relayers?: string[] | undefined;
            } & {
                fee?: ({
                    recvFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    ackFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    timeoutFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                } & {
                    recvFee?: ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] & ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & Record<Exclude<keyof I["incentivizedPackets"][number]["packetFees"][number]["fee"]["recvFee"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["incentivizedPackets"][number]["packetFees"][number]["fee"]["recvFee"], keyof {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[]>, never>) | undefined;
                    ackFee?: ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] & ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & Record<Exclude<keyof I["incentivizedPackets"][number]["packetFees"][number]["fee"]["ackFee"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["incentivizedPackets"][number]["packetFees"][number]["fee"]["ackFee"], keyof {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[]>, never>) | undefined;
                    timeoutFee?: ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] & ({
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    } & Record<Exclude<keyof I["incentivizedPackets"][number]["packetFees"][number]["fee"]["timeoutFee"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["incentivizedPackets"][number]["packetFees"][number]["fee"]["timeoutFee"], keyof {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[]>, never>) | undefined;
                } & Record<Exclude<keyof I["incentivizedPackets"][number]["packetFees"][number]["fee"], keyof import("./fee").Fee>, never>) | undefined;
                refundAddress?: string | undefined;
                relayers?: (string[] & string[] & Record<Exclude<keyof I["incentivizedPackets"][number]["packetFees"][number]["relayers"], keyof string[]>, never>) | undefined;
            } & Record<Exclude<keyof I["incentivizedPackets"][number]["packetFees"][number], keyof import("./fee").PacketFee>, never>)[] & Record<Exclude<keyof I["incentivizedPackets"][number]["packetFees"], keyof {
                fee?: {
                    recvFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    ackFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    timeoutFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                } | undefined;
                refundAddress?: string | undefined;
                relayers?: string[] | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["incentivizedPackets"][number], keyof IdentifiedPacketFees>, never>)[] & Record<Exclude<keyof I["incentivizedPackets"], keyof {
            packetId?: {
                portId?: string | undefined;
                channelId?: string | undefined;
                sequence?: bigint | undefined;
            } | undefined;
            packetFees?: {
                fee?: {
                    recvFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    ackFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                    timeoutFee?: {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[] | undefined;
                } | undefined;
                refundAddress?: string | undefined;
                relayers?: string[] | undefined;
            }[] | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "incentivizedPackets">, never>>(object: I): QueryIncentivizedPacketsForChannelResponse;
};
export declare const QueryTotalRecvFeesRequest: {
    typeUrl: string;
    encode(message: QueryTotalRecvFeesRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryTotalRecvFeesRequest;
    fromJSON(object: any): QueryTotalRecvFeesRequest;
    toJSON(message: QueryTotalRecvFeesRequest): unknown;
    fromPartial<I extends {
        packetId?: {
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
        } | undefined;
    } & {
        packetId?: ({
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
        } & {
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
        } & Record<Exclude<keyof I["packetId"], keyof PacketId>, never>) | undefined;
    } & Record<Exclude<keyof I, "packetId">, never>>(object: I): QueryTotalRecvFeesRequest;
};
export declare const QueryTotalRecvFeesResponse: {
    typeUrl: string;
    encode(message: QueryTotalRecvFeesResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryTotalRecvFeesResponse;
    fromJSON(object: any): QueryTotalRecvFeesResponse;
    toJSON(message: QueryTotalRecvFeesResponse): unknown;
    fromPartial<I extends {
        recvFees?: {
            denom?: string | undefined;
            amount?: string | undefined;
        }[] | undefined;
    } & {
        recvFees?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        }[] & ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["recvFees"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["recvFees"], keyof {
            denom?: string | undefined;
            amount?: string | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "recvFees">, never>>(object: I): QueryTotalRecvFeesResponse;
};
export declare const QueryTotalAckFeesRequest: {
    typeUrl: string;
    encode(message: QueryTotalAckFeesRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryTotalAckFeesRequest;
    fromJSON(object: any): QueryTotalAckFeesRequest;
    toJSON(message: QueryTotalAckFeesRequest): unknown;
    fromPartial<I extends {
        packetId?: {
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
        } | undefined;
    } & {
        packetId?: ({
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
        } & {
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
        } & Record<Exclude<keyof I["packetId"], keyof PacketId>, never>) | undefined;
    } & Record<Exclude<keyof I, "packetId">, never>>(object: I): QueryTotalAckFeesRequest;
};
export declare const QueryTotalAckFeesResponse: {
    typeUrl: string;
    encode(message: QueryTotalAckFeesResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryTotalAckFeesResponse;
    fromJSON(object: any): QueryTotalAckFeesResponse;
    toJSON(message: QueryTotalAckFeesResponse): unknown;
    fromPartial<I extends {
        ackFees?: {
            denom?: string | undefined;
            amount?: string | undefined;
        }[] | undefined;
    } & {
        ackFees?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        }[] & ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["ackFees"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["ackFees"], keyof {
            denom?: string | undefined;
            amount?: string | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "ackFees">, never>>(object: I): QueryTotalAckFeesResponse;
};
export declare const QueryTotalTimeoutFeesRequest: {
    typeUrl: string;
    encode(message: QueryTotalTimeoutFeesRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryTotalTimeoutFeesRequest;
    fromJSON(object: any): QueryTotalTimeoutFeesRequest;
    toJSON(message: QueryTotalTimeoutFeesRequest): unknown;
    fromPartial<I extends {
        packetId?: {
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
        } | undefined;
    } & {
        packetId?: ({
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
        } & {
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
        } & Record<Exclude<keyof I["packetId"], keyof PacketId>, never>) | undefined;
    } & Record<Exclude<keyof I, "packetId">, never>>(object: I): QueryTotalTimeoutFeesRequest;
};
export declare const QueryTotalTimeoutFeesResponse: {
    typeUrl: string;
    encode(message: QueryTotalTimeoutFeesResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryTotalTimeoutFeesResponse;
    fromJSON(object: any): QueryTotalTimeoutFeesResponse;
    toJSON(message: QueryTotalTimeoutFeesResponse): unknown;
    fromPartial<I extends {
        timeoutFees?: {
            denom?: string | undefined;
            amount?: string | undefined;
        }[] | undefined;
    } & {
        timeoutFees?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        }[] & ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["timeoutFees"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["timeoutFees"], keyof {
            denom?: string | undefined;
            amount?: string | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "timeoutFees">, never>>(object: I): QueryTotalTimeoutFeesResponse;
};
export declare const QueryPayeeRequest: {
    typeUrl: string;
    encode(message: QueryPayeeRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryPayeeRequest;
    fromJSON(object: any): QueryPayeeRequest;
    toJSON(message: QueryPayeeRequest): unknown;
    fromPartial<I extends {
        channelId?: string | undefined;
        relayer?: string | undefined;
    } & {
        channelId?: string | undefined;
        relayer?: string | undefined;
    } & Record<Exclude<keyof I, keyof QueryPayeeRequest>, never>>(object: I): QueryPayeeRequest;
};
export declare const QueryPayeeResponse: {
    typeUrl: string;
    encode(message: QueryPayeeResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryPayeeResponse;
    fromJSON(object: any): QueryPayeeResponse;
    toJSON(message: QueryPayeeResponse): unknown;
    fromPartial<I extends {
        payeeAddress?: string | undefined;
    } & {
        payeeAddress?: string | undefined;
    } & Record<Exclude<keyof I, "payeeAddress">, never>>(object: I): QueryPayeeResponse;
};
export declare const QueryCounterpartyPayeeRequest: {
    typeUrl: string;
    encode(message: QueryCounterpartyPayeeRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryCounterpartyPayeeRequest;
    fromJSON(object: any): QueryCounterpartyPayeeRequest;
    toJSON(message: QueryCounterpartyPayeeRequest): unknown;
    fromPartial<I extends {
        channelId?: string | undefined;
        relayer?: string | undefined;
    } & {
        channelId?: string | undefined;
        relayer?: string | undefined;
    } & Record<Exclude<keyof I, keyof QueryCounterpartyPayeeRequest>, never>>(object: I): QueryCounterpartyPayeeRequest;
};
export declare const QueryCounterpartyPayeeResponse: {
    typeUrl: string;
    encode(message: QueryCounterpartyPayeeResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryCounterpartyPayeeResponse;
    fromJSON(object: any): QueryCounterpartyPayeeResponse;
    toJSON(message: QueryCounterpartyPayeeResponse): unknown;
    fromPartial<I extends {
        counterpartyPayee?: string | undefined;
    } & {
        counterpartyPayee?: string | undefined;
    } & Record<Exclude<keyof I, "counterpartyPayee">, never>>(object: I): QueryCounterpartyPayeeResponse;
};
export declare const QueryFeeEnabledChannelsRequest: {
    typeUrl: string;
    encode(message: QueryFeeEnabledChannelsRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryFeeEnabledChannelsRequest;
    fromJSON(object: any): QueryFeeEnabledChannelsRequest;
    toJSON(message: QueryFeeEnabledChannelsRequest): unknown;
    fromPartial<I extends {
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
        queryHeight?: bigint | undefined;
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
        queryHeight?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof QueryFeeEnabledChannelsRequest>, never>>(object: I): QueryFeeEnabledChannelsRequest;
};
export declare const QueryFeeEnabledChannelsResponse: {
    typeUrl: string;
    encode(message: QueryFeeEnabledChannelsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryFeeEnabledChannelsResponse;
    fromJSON(object: any): QueryFeeEnabledChannelsResponse;
    toJSON(message: QueryFeeEnabledChannelsResponse): unknown;
    fromPartial<I extends {
        feeEnabledChannels?: {
            portId?: string | undefined;
            channelId?: string | undefined;
        }[] | undefined;
    } & {
        feeEnabledChannels?: ({
            portId?: string | undefined;
            channelId?: string | undefined;
        }[] & ({
            portId?: string | undefined;
            channelId?: string | undefined;
        } & {
            portId?: string | undefined;
            channelId?: string | undefined;
        } & Record<Exclude<keyof I["feeEnabledChannels"][number], keyof FeeEnabledChannel>, never>)[] & Record<Exclude<keyof I["feeEnabledChannels"], keyof {
            portId?: string | undefined;
            channelId?: string | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "feeEnabledChannels">, never>>(object: I): QueryFeeEnabledChannelsResponse;
};
export declare const QueryFeeEnabledChannelRequest: {
    typeUrl: string;
    encode(message: QueryFeeEnabledChannelRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryFeeEnabledChannelRequest;
    fromJSON(object: any): QueryFeeEnabledChannelRequest;
    toJSON(message: QueryFeeEnabledChannelRequest): unknown;
    fromPartial<I extends {
        portId?: string | undefined;
        channelId?: string | undefined;
    } & {
        portId?: string | undefined;
        channelId?: string | undefined;
    } & Record<Exclude<keyof I, keyof QueryFeeEnabledChannelRequest>, never>>(object: I): QueryFeeEnabledChannelRequest;
};
export declare const QueryFeeEnabledChannelResponse: {
    typeUrl: string;
    encode(message: QueryFeeEnabledChannelResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryFeeEnabledChannelResponse;
    fromJSON(object: any): QueryFeeEnabledChannelResponse;
    toJSON(message: QueryFeeEnabledChannelResponse): unknown;
    fromPartial<I extends {
        feeEnabled?: boolean | undefined;
    } & {
        feeEnabled?: boolean | undefined;
    } & Record<Exclude<keyof I, "feeEnabled">, never>>(object: I): QueryFeeEnabledChannelResponse;
};
/** Query defines the ICS29 gRPC querier service. */
export interface Query {
    /** IncentivizedPackets returns all incentivized packets and their associated fees */
    IncentivizedPackets(request: QueryIncentivizedPacketsRequest): Promise<QueryIncentivizedPacketsResponse>;
    /** IncentivizedPacket returns all packet fees for a packet given its identifier */
    IncentivizedPacket(request: QueryIncentivizedPacketRequest): Promise<QueryIncentivizedPacketResponse>;
    /** Gets all incentivized packets for a specific channel */
    IncentivizedPacketsForChannel(request: QueryIncentivizedPacketsForChannelRequest): Promise<QueryIncentivizedPacketsForChannelResponse>;
    /** TotalRecvFees returns the total receive fees for a packet given its identifier */
    TotalRecvFees(request: QueryTotalRecvFeesRequest): Promise<QueryTotalRecvFeesResponse>;
    /** TotalAckFees returns the total acknowledgement fees for a packet given its identifier */
    TotalAckFees(request: QueryTotalAckFeesRequest): Promise<QueryTotalAckFeesResponse>;
    /** TotalTimeoutFees returns the total timeout fees for a packet given its identifier */
    TotalTimeoutFees(request: QueryTotalTimeoutFeesRequest): Promise<QueryTotalTimeoutFeesResponse>;
    /** Payee returns the registered payee address for a specific channel given the relayer address */
    Payee(request: QueryPayeeRequest): Promise<QueryPayeeResponse>;
    /** CounterpartyPayee returns the registered counterparty payee for forward relaying */
    CounterpartyPayee(request: QueryCounterpartyPayeeRequest): Promise<QueryCounterpartyPayeeResponse>;
    /** FeeEnabledChannels returns a list of all fee enabled channels */
    FeeEnabledChannels(request: QueryFeeEnabledChannelsRequest): Promise<QueryFeeEnabledChannelsResponse>;
    /** FeeEnabledChannel returns true if the provided port and channel identifiers belong to a fee enabled channel */
    FeeEnabledChannel(request: QueryFeeEnabledChannelRequest): Promise<QueryFeeEnabledChannelResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    IncentivizedPackets(request: QueryIncentivizedPacketsRequest): Promise<QueryIncentivizedPacketsResponse>;
    IncentivizedPacket(request: QueryIncentivizedPacketRequest): Promise<QueryIncentivizedPacketResponse>;
    IncentivizedPacketsForChannel(request: QueryIncentivizedPacketsForChannelRequest): Promise<QueryIncentivizedPacketsForChannelResponse>;
    TotalRecvFees(request: QueryTotalRecvFeesRequest): Promise<QueryTotalRecvFeesResponse>;
    TotalAckFees(request: QueryTotalAckFeesRequest): Promise<QueryTotalAckFeesResponse>;
    TotalTimeoutFees(request: QueryTotalTimeoutFeesRequest): Promise<QueryTotalTimeoutFeesResponse>;
    Payee(request: QueryPayeeRequest): Promise<QueryPayeeResponse>;
    CounterpartyPayee(request: QueryCounterpartyPayeeRequest): Promise<QueryCounterpartyPayeeResponse>;
    FeeEnabledChannels(request: QueryFeeEnabledChannelsRequest): Promise<QueryFeeEnabledChannelsResponse>;
    FeeEnabledChannel(request: QueryFeeEnabledChannelRequest): Promise<QueryFeeEnabledChannelResponse>;
}
