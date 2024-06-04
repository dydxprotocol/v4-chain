import { Coin } from "../../../../cosmos/base/v1beta1/coin";
import { PacketId } from "../../../core/channel/v1/channel";
import { BinaryReader, BinaryWriter } from "../../../../binary";
export declare const protobufPackage = "ibc.applications.fee.v1";
/** Fee defines the ICS29 receive, acknowledgement and timeout fees */
export interface Fee {
    /** the packet receive fee */
    recvFee: Coin[];
    /** the packet acknowledgement fee */
    ackFee: Coin[];
    /** the packet timeout fee */
    timeoutFee: Coin[];
}
/** PacketFee contains ICS29 relayer fees, refund address and optional list of permitted relayers */
export interface PacketFee {
    /** fee encapsulates the recv, ack and timeout fees associated with an IBC packet */
    fee: Fee;
    /** the refund address for unspent fees */
    refundAddress: string;
    /** optional list of relayers permitted to receive fees */
    relayers: string[];
}
/** PacketFees contains a list of type PacketFee */
export interface PacketFees {
    /** list of packet fees */
    packetFees: PacketFee[];
}
/** IdentifiedPacketFees contains a list of type PacketFee and associated PacketId */
export interface IdentifiedPacketFees {
    /** unique packet identifier comprised of the channel ID, port ID and sequence */
    packetId: PacketId;
    /** list of packet fees */
    packetFees: PacketFee[];
}
export declare const Fee: {
    typeUrl: string;
    encode(message: Fee, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Fee;
    fromJSON(object: any): Fee;
    toJSON(message: Fee): unknown;
    fromPartial<I extends {
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
        } & Record<Exclude<keyof I["recvFee"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["recvFee"], keyof {
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
        } & Record<Exclude<keyof I["ackFee"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["ackFee"], keyof {
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
        } & Record<Exclude<keyof I["timeoutFee"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["timeoutFee"], keyof {
            denom?: string | undefined;
            amount?: string | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof Fee>, never>>(object: I): Fee;
};
export declare const PacketFee: {
    typeUrl: string;
    encode(message: PacketFee, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): PacketFee;
    fromJSON(object: any): PacketFee;
    toJSON(message: PacketFee): unknown;
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
            } & Record<Exclude<keyof I["fee"]["recvFee"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["fee"]["recvFee"], keyof {
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
            } & Record<Exclude<keyof I["fee"]["ackFee"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["fee"]["ackFee"], keyof {
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
            } & Record<Exclude<keyof I["fee"]["timeoutFee"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["fee"]["timeoutFee"], keyof {
                denom?: string | undefined;
                amount?: string | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["fee"], keyof Fee>, never>) | undefined;
        refundAddress?: string | undefined;
        relayers?: (string[] & string[] & Record<Exclude<keyof I["relayers"], keyof string[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof PacketFee>, never>>(object: I): PacketFee;
};
export declare const PacketFees: {
    typeUrl: string;
    encode(message: PacketFees, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): PacketFees;
    fromJSON(object: any): PacketFees;
    toJSON(message: PacketFees): unknown;
    fromPartial<I extends {
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
                } & Record<Exclude<keyof I["packetFees"][number]["fee"]["recvFee"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["packetFees"][number]["fee"]["recvFee"], keyof {
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
                } & Record<Exclude<keyof I["packetFees"][number]["fee"]["ackFee"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["packetFees"][number]["fee"]["ackFee"], keyof {
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
                } & Record<Exclude<keyof I["packetFees"][number]["fee"]["timeoutFee"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["packetFees"][number]["fee"]["timeoutFee"], keyof {
                    denom?: string | undefined;
                    amount?: string | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["packetFees"][number]["fee"], keyof Fee>, never>) | undefined;
            refundAddress?: string | undefined;
            relayers?: (string[] & string[] & Record<Exclude<keyof I["packetFees"][number]["relayers"], keyof string[]>, never>) | undefined;
        } & Record<Exclude<keyof I["packetFees"][number], keyof PacketFee>, never>)[] & Record<Exclude<keyof I["packetFees"], keyof {
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
    } & Record<Exclude<keyof I, "packetFees">, never>>(object: I): PacketFees;
};
export declare const IdentifiedPacketFees: {
    typeUrl: string;
    encode(message: IdentifiedPacketFees, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): IdentifiedPacketFees;
    fromJSON(object: any): IdentifiedPacketFees;
    toJSON(message: IdentifiedPacketFees): unknown;
    fromPartial<I extends {
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
        } & Record<Exclude<keyof I["packetId"], keyof PacketId>, never>) | undefined;
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
                } & Record<Exclude<keyof I["packetFees"][number]["fee"]["recvFee"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["packetFees"][number]["fee"]["recvFee"], keyof {
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
                } & Record<Exclude<keyof I["packetFees"][number]["fee"]["ackFee"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["packetFees"][number]["fee"]["ackFee"], keyof {
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
                } & Record<Exclude<keyof I["packetFees"][number]["fee"]["timeoutFee"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["packetFees"][number]["fee"]["timeoutFee"], keyof {
                    denom?: string | undefined;
                    amount?: string | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["packetFees"][number]["fee"], keyof Fee>, never>) | undefined;
            refundAddress?: string | undefined;
            relayers?: (string[] & string[] & Record<Exclude<keyof I["packetFees"][number]["relayers"], keyof string[]>, never>) | undefined;
        } & Record<Exclude<keyof I["packetFees"][number], keyof PacketFee>, never>)[] & Record<Exclude<keyof I["packetFees"], keyof {
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
    } & Record<Exclude<keyof I, keyof IdentifiedPacketFees>, never>>(object: I): IdentifiedPacketFees;
};
