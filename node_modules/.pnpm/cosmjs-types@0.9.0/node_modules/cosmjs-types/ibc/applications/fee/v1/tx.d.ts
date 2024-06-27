import { Fee, PacketFee } from "./fee";
import { PacketId } from "../../../core/channel/v1/channel";
import { BinaryReader, BinaryWriter } from "../../../../binary";
import { Rpc } from "../../../../helpers";
export declare const protobufPackage = "ibc.applications.fee.v1";
/** MsgRegisterPayee defines the request type for the RegisterPayee rpc */
export interface MsgRegisterPayee {
    /** unique port identifier */
    portId: string;
    /** unique channel identifier */
    channelId: string;
    /** the relayer address */
    relayer: string;
    /** the payee address */
    payee: string;
}
/** MsgRegisterPayeeResponse defines the response type for the RegisterPayee rpc */
export interface MsgRegisterPayeeResponse {
}
/** MsgRegisterCounterpartyPayee defines the request type for the RegisterCounterpartyPayee rpc */
export interface MsgRegisterCounterpartyPayee {
    /** unique port identifier */
    portId: string;
    /** unique channel identifier */
    channelId: string;
    /** the relayer address */
    relayer: string;
    /** the counterparty payee address */
    counterpartyPayee: string;
}
/** MsgRegisterCounterpartyPayeeResponse defines the response type for the RegisterCounterpartyPayee rpc */
export interface MsgRegisterCounterpartyPayeeResponse {
}
/**
 * MsgPayPacketFee defines the request type for the PayPacketFee rpc
 * This Msg can be used to pay for a packet at the next sequence send & should be combined with the Msg that will be
 * paid for
 */
export interface MsgPayPacketFee {
    /** fee encapsulates the recv, ack and timeout fees associated with an IBC packet */
    fee: Fee;
    /** the source port unique identifier */
    sourcePortId: string;
    /** the source channel unique identifer */
    sourceChannelId: string;
    /** account address to refund fee if necessary */
    signer: string;
    /** optional list of relayers permitted to the receive packet fees */
    relayers: string[];
}
/** MsgPayPacketFeeResponse defines the response type for the PayPacketFee rpc */
export interface MsgPayPacketFeeResponse {
}
/**
 * MsgPayPacketFeeAsync defines the request type for the PayPacketFeeAsync rpc
 * This Msg can be used to pay for a packet at a specified sequence (instead of the next sequence send)
 */
export interface MsgPayPacketFeeAsync {
    /** unique packet identifier comprised of the channel ID, port ID and sequence */
    packetId: PacketId;
    /** the packet fee associated with a particular IBC packet */
    packetFee: PacketFee;
}
/** MsgPayPacketFeeAsyncResponse defines the response type for the PayPacketFeeAsync rpc */
export interface MsgPayPacketFeeAsyncResponse {
}
export declare const MsgRegisterPayee: {
    typeUrl: string;
    encode(message: MsgRegisterPayee, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgRegisterPayee;
    fromJSON(object: any): MsgRegisterPayee;
    toJSON(message: MsgRegisterPayee): unknown;
    fromPartial<I extends {
        portId?: string | undefined;
        channelId?: string | undefined;
        relayer?: string | undefined;
        payee?: string | undefined;
    } & {
        portId?: string | undefined;
        channelId?: string | undefined;
        relayer?: string | undefined;
        payee?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgRegisterPayee>, never>>(object: I): MsgRegisterPayee;
};
export declare const MsgRegisterPayeeResponse: {
    typeUrl: string;
    encode(_: MsgRegisterPayeeResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgRegisterPayeeResponse;
    fromJSON(_: any): MsgRegisterPayeeResponse;
    toJSON(_: MsgRegisterPayeeResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgRegisterPayeeResponse;
};
export declare const MsgRegisterCounterpartyPayee: {
    typeUrl: string;
    encode(message: MsgRegisterCounterpartyPayee, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgRegisterCounterpartyPayee;
    fromJSON(object: any): MsgRegisterCounterpartyPayee;
    toJSON(message: MsgRegisterCounterpartyPayee): unknown;
    fromPartial<I extends {
        portId?: string | undefined;
        channelId?: string | undefined;
        relayer?: string | undefined;
        counterpartyPayee?: string | undefined;
    } & {
        portId?: string | undefined;
        channelId?: string | undefined;
        relayer?: string | undefined;
        counterpartyPayee?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgRegisterCounterpartyPayee>, never>>(object: I): MsgRegisterCounterpartyPayee;
};
export declare const MsgRegisterCounterpartyPayeeResponse: {
    typeUrl: string;
    encode(_: MsgRegisterCounterpartyPayeeResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgRegisterCounterpartyPayeeResponse;
    fromJSON(_: any): MsgRegisterCounterpartyPayeeResponse;
    toJSON(_: MsgRegisterCounterpartyPayeeResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgRegisterCounterpartyPayeeResponse;
};
export declare const MsgPayPacketFee: {
    typeUrl: string;
    encode(message: MsgPayPacketFee, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgPayPacketFee;
    fromJSON(object: any): MsgPayPacketFee;
    toJSON(message: MsgPayPacketFee): unknown;
    fromPartial<I extends {
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
        sourcePortId?: string | undefined;
        sourceChannelId?: string | undefined;
        signer?: string | undefined;
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
            } & Record<Exclude<keyof I["fee"]["recvFee"][number], keyof import("../../../../cosmos/base/v1beta1/coin").Coin>, never>)[] & Record<Exclude<keyof I["fee"]["recvFee"], keyof {
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
            } & Record<Exclude<keyof I["fee"]["ackFee"][number], keyof import("../../../../cosmos/base/v1beta1/coin").Coin>, never>)[] & Record<Exclude<keyof I["fee"]["ackFee"], keyof {
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
            } & Record<Exclude<keyof I["fee"]["timeoutFee"][number], keyof import("../../../../cosmos/base/v1beta1/coin").Coin>, never>)[] & Record<Exclude<keyof I["fee"]["timeoutFee"], keyof {
                denom?: string | undefined;
                amount?: string | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["fee"], keyof Fee>, never>) | undefined;
        sourcePortId?: string | undefined;
        sourceChannelId?: string | undefined;
        signer?: string | undefined;
        relayers?: (string[] & string[] & Record<Exclude<keyof I["relayers"], keyof string[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof MsgPayPacketFee>, never>>(object: I): MsgPayPacketFee;
};
export declare const MsgPayPacketFeeResponse: {
    typeUrl: string;
    encode(_: MsgPayPacketFeeResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgPayPacketFeeResponse;
    fromJSON(_: any): MsgPayPacketFeeResponse;
    toJSON(_: MsgPayPacketFeeResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgPayPacketFeeResponse;
};
export declare const MsgPayPacketFeeAsync: {
    typeUrl: string;
    encode(message: MsgPayPacketFeeAsync, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgPayPacketFeeAsync;
    fromJSON(object: any): MsgPayPacketFeeAsync;
    toJSON(message: MsgPayPacketFeeAsync): unknown;
    fromPartial<I extends {
        packetId?: {
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
        } | undefined;
        packetFee?: {
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
        packetFee?: ({
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
                } & Record<Exclude<keyof I["packetFee"]["fee"]["recvFee"][number], keyof import("../../../../cosmos/base/v1beta1/coin").Coin>, never>)[] & Record<Exclude<keyof I["packetFee"]["fee"]["recvFee"], keyof {
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
                } & Record<Exclude<keyof I["packetFee"]["fee"]["ackFee"][number], keyof import("../../../../cosmos/base/v1beta1/coin").Coin>, never>)[] & Record<Exclude<keyof I["packetFee"]["fee"]["ackFee"], keyof {
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
                } & Record<Exclude<keyof I["packetFee"]["fee"]["timeoutFee"][number], keyof import("../../../../cosmos/base/v1beta1/coin").Coin>, never>)[] & Record<Exclude<keyof I["packetFee"]["fee"]["timeoutFee"], keyof {
                    denom?: string | undefined;
                    amount?: string | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["packetFee"]["fee"], keyof Fee>, never>) | undefined;
            refundAddress?: string | undefined;
            relayers?: (string[] & string[] & Record<Exclude<keyof I["packetFee"]["relayers"], keyof string[]>, never>) | undefined;
        } & Record<Exclude<keyof I["packetFee"], keyof PacketFee>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof MsgPayPacketFeeAsync>, never>>(object: I): MsgPayPacketFeeAsync;
};
export declare const MsgPayPacketFeeAsyncResponse: {
    typeUrl: string;
    encode(_: MsgPayPacketFeeAsyncResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgPayPacketFeeAsyncResponse;
    fromJSON(_: any): MsgPayPacketFeeAsyncResponse;
    toJSON(_: MsgPayPacketFeeAsyncResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgPayPacketFeeAsyncResponse;
};
/** Msg defines the ICS29 Msg service. */
export interface Msg {
    /**
     * RegisterPayee defines a rpc handler method for MsgRegisterPayee
     * RegisterPayee is called by the relayer on each channelEnd and allows them to set an optional
     * payee to which reverse and timeout relayer packet fees will be paid out. The payee should be registered on
     * the source chain from which packets originate as this is where fee distribution takes place. This function may be
     * called more than once by a relayer, in which case, the latest payee is always used.
     */
    RegisterPayee(request: MsgRegisterPayee): Promise<MsgRegisterPayeeResponse>;
    /**
     * RegisterCounterpartyPayee defines a rpc handler method for MsgRegisterCounterpartyPayee
     * RegisterCounterpartyPayee is called by the relayer on each channelEnd and allows them to specify the counterparty
     * payee address before relaying. This ensures they will be properly compensated for forward relaying since
     * the destination chain must include the registered counterparty payee address in the acknowledgement. This function
     * may be called more than once by a relayer, in which case, the latest counterparty payee address is always used.
     */
    RegisterCounterpartyPayee(request: MsgRegisterCounterpartyPayee): Promise<MsgRegisterCounterpartyPayeeResponse>;
    /**
     * PayPacketFee defines a rpc handler method for MsgPayPacketFee
     * PayPacketFee is an open callback that may be called by any module/user that wishes to escrow funds in order to
     * incentivize the relaying of the packet at the next sequence
     * NOTE: This method is intended to be used within a multi msg transaction, where the subsequent msg that follows
     * initiates the lifecycle of the incentivized packet
     */
    PayPacketFee(request: MsgPayPacketFee): Promise<MsgPayPacketFeeResponse>;
    /**
     * PayPacketFeeAsync defines a rpc handler method for MsgPayPacketFeeAsync
     * PayPacketFeeAsync is an open callback that may be called by any module/user that wishes to escrow funds in order to
     * incentivize the relaying of a known packet (i.e. at a particular sequence)
     */
    PayPacketFeeAsync(request: MsgPayPacketFeeAsync): Promise<MsgPayPacketFeeAsyncResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    RegisterPayee(request: MsgRegisterPayee): Promise<MsgRegisterPayeeResponse>;
    RegisterCounterpartyPayee(request: MsgRegisterCounterpartyPayee): Promise<MsgRegisterCounterpartyPayeeResponse>;
    PayPacketFee(request: MsgPayPacketFee): Promise<MsgPayPacketFeeResponse>;
    PayPacketFeeAsync(request: MsgPayPacketFeeAsync): Promise<MsgPayPacketFeeAsyncResponse>;
}
