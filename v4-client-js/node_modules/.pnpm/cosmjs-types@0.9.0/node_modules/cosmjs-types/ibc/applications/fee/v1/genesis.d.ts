import { IdentifiedPacketFees } from "./fee";
import { PacketId } from "../../../core/channel/v1/channel";
import { BinaryReader, BinaryWriter } from "../../../../binary";
export declare const protobufPackage = "ibc.applications.fee.v1";
/** GenesisState defines the ICS29 fee middleware genesis state */
export interface GenesisState {
    /** list of identified packet fees */
    identifiedFees: IdentifiedPacketFees[];
    /** list of fee enabled channels */
    feeEnabledChannels: FeeEnabledChannel[];
    /** list of registered payees */
    registeredPayees: RegisteredPayee[];
    /** list of registered counterparty payees */
    registeredCounterpartyPayees: RegisteredCounterpartyPayee[];
    /** list of forward relayer addresses */
    forwardRelayers: ForwardRelayerAddress[];
}
/** FeeEnabledChannel contains the PortID & ChannelID for a fee enabled channel */
export interface FeeEnabledChannel {
    /** unique port identifier */
    portId: string;
    /** unique channel identifier */
    channelId: string;
}
/** RegisteredPayee contains the relayer address and payee address for a specific channel */
export interface RegisteredPayee {
    /** unique channel identifier */
    channelId: string;
    /** the relayer address */
    relayer: string;
    /** the payee address */
    payee: string;
}
/**
 * RegisteredCounterpartyPayee contains the relayer address and counterparty payee address for a specific channel (used
 * for recv fee distribution)
 */
export interface RegisteredCounterpartyPayee {
    /** unique channel identifier */
    channelId: string;
    /** the relayer address */
    relayer: string;
    /** the counterparty payee address */
    counterpartyPayee: string;
}
/** ForwardRelayerAddress contains the forward relayer address and PacketId used for async acknowledgements */
export interface ForwardRelayerAddress {
    /** the forward relayer address */
    address: string;
    /** unique packet identifer comprised of the channel ID, port ID and sequence */
    packetId: PacketId;
}
export declare const GenesisState: {
    typeUrl: string;
    encode(message: GenesisState, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
    fromPartial<I extends {
        identifiedFees?: {
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
        feeEnabledChannels?: {
            portId?: string | undefined;
            channelId?: string | undefined;
        }[] | undefined;
        registeredPayees?: {
            channelId?: string | undefined;
            relayer?: string | undefined;
            payee?: string | undefined;
        }[] | undefined;
        registeredCounterpartyPayees?: {
            channelId?: string | undefined;
            relayer?: string | undefined;
            counterpartyPayee?: string | undefined;
        }[] | undefined;
        forwardRelayers?: {
            address?: string | undefined;
            packetId?: {
                portId?: string | undefined;
                channelId?: string | undefined;
                sequence?: bigint | undefined;
            } | undefined;
        }[] | undefined;
    } & {
        identifiedFees?: ({
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
            } & Record<Exclude<keyof I["identifiedFees"][number]["packetId"], keyof PacketId>, never>) | undefined;
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
                    } & Record<Exclude<keyof I["identifiedFees"][number]["packetFees"][number]["fee"]["recvFee"][number], keyof import("../../../../cosmos/base/v1beta1/coin").Coin>, never>)[] & Record<Exclude<keyof I["identifiedFees"][number]["packetFees"][number]["fee"]["recvFee"], keyof {
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
                    } & Record<Exclude<keyof I["identifiedFees"][number]["packetFees"][number]["fee"]["ackFee"][number], keyof import("../../../../cosmos/base/v1beta1/coin").Coin>, never>)[] & Record<Exclude<keyof I["identifiedFees"][number]["packetFees"][number]["fee"]["ackFee"], keyof {
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
                    } & Record<Exclude<keyof I["identifiedFees"][number]["packetFees"][number]["fee"]["timeoutFee"][number], keyof import("../../../../cosmos/base/v1beta1/coin").Coin>, never>)[] & Record<Exclude<keyof I["identifiedFees"][number]["packetFees"][number]["fee"]["timeoutFee"], keyof {
                        denom?: string | undefined;
                        amount?: string | undefined;
                    }[]>, never>) | undefined;
                } & Record<Exclude<keyof I["identifiedFees"][number]["packetFees"][number]["fee"], keyof import("./fee").Fee>, never>) | undefined;
                refundAddress?: string | undefined;
                relayers?: (string[] & string[] & Record<Exclude<keyof I["identifiedFees"][number]["packetFees"][number]["relayers"], keyof string[]>, never>) | undefined;
            } & Record<Exclude<keyof I["identifiedFees"][number]["packetFees"][number], keyof import("./fee").PacketFee>, never>)[] & Record<Exclude<keyof I["identifiedFees"][number]["packetFees"], keyof {
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
        } & Record<Exclude<keyof I["identifiedFees"][number], keyof IdentifiedPacketFees>, never>)[] & Record<Exclude<keyof I["identifiedFees"], keyof {
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
        registeredPayees?: ({
            channelId?: string | undefined;
            relayer?: string | undefined;
            payee?: string | undefined;
        }[] & ({
            channelId?: string | undefined;
            relayer?: string | undefined;
            payee?: string | undefined;
        } & {
            channelId?: string | undefined;
            relayer?: string | undefined;
            payee?: string | undefined;
        } & Record<Exclude<keyof I["registeredPayees"][number], keyof RegisteredPayee>, never>)[] & Record<Exclude<keyof I["registeredPayees"], keyof {
            channelId?: string | undefined;
            relayer?: string | undefined;
            payee?: string | undefined;
        }[]>, never>) | undefined;
        registeredCounterpartyPayees?: ({
            channelId?: string | undefined;
            relayer?: string | undefined;
            counterpartyPayee?: string | undefined;
        }[] & ({
            channelId?: string | undefined;
            relayer?: string | undefined;
            counterpartyPayee?: string | undefined;
        } & {
            channelId?: string | undefined;
            relayer?: string | undefined;
            counterpartyPayee?: string | undefined;
        } & Record<Exclude<keyof I["registeredCounterpartyPayees"][number], keyof RegisteredCounterpartyPayee>, never>)[] & Record<Exclude<keyof I["registeredCounterpartyPayees"], keyof {
            channelId?: string | undefined;
            relayer?: string | undefined;
            counterpartyPayee?: string | undefined;
        }[]>, never>) | undefined;
        forwardRelayers?: ({
            address?: string | undefined;
            packetId?: {
                portId?: string | undefined;
                channelId?: string | undefined;
                sequence?: bigint | undefined;
            } | undefined;
        }[] & ({
            address?: string | undefined;
            packetId?: {
                portId?: string | undefined;
                channelId?: string | undefined;
                sequence?: bigint | undefined;
            } | undefined;
        } & {
            address?: string | undefined;
            packetId?: ({
                portId?: string | undefined;
                channelId?: string | undefined;
                sequence?: bigint | undefined;
            } & {
                portId?: string | undefined;
                channelId?: string | undefined;
                sequence?: bigint | undefined;
            } & Record<Exclude<keyof I["forwardRelayers"][number]["packetId"], keyof PacketId>, never>) | undefined;
        } & Record<Exclude<keyof I["forwardRelayers"][number], keyof ForwardRelayerAddress>, never>)[] & Record<Exclude<keyof I["forwardRelayers"], keyof {
            address?: string | undefined;
            packetId?: {
                portId?: string | undefined;
                channelId?: string | undefined;
                sequence?: bigint | undefined;
            } | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof GenesisState>, never>>(object: I): GenesisState;
};
export declare const FeeEnabledChannel: {
    typeUrl: string;
    encode(message: FeeEnabledChannel, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): FeeEnabledChannel;
    fromJSON(object: any): FeeEnabledChannel;
    toJSON(message: FeeEnabledChannel): unknown;
    fromPartial<I extends {
        portId?: string | undefined;
        channelId?: string | undefined;
    } & {
        portId?: string | undefined;
        channelId?: string | undefined;
    } & Record<Exclude<keyof I, keyof FeeEnabledChannel>, never>>(object: I): FeeEnabledChannel;
};
export declare const RegisteredPayee: {
    typeUrl: string;
    encode(message: RegisteredPayee, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): RegisteredPayee;
    fromJSON(object: any): RegisteredPayee;
    toJSON(message: RegisteredPayee): unknown;
    fromPartial<I extends {
        channelId?: string | undefined;
        relayer?: string | undefined;
        payee?: string | undefined;
    } & {
        channelId?: string | undefined;
        relayer?: string | undefined;
        payee?: string | undefined;
    } & Record<Exclude<keyof I, keyof RegisteredPayee>, never>>(object: I): RegisteredPayee;
};
export declare const RegisteredCounterpartyPayee: {
    typeUrl: string;
    encode(message: RegisteredCounterpartyPayee, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): RegisteredCounterpartyPayee;
    fromJSON(object: any): RegisteredCounterpartyPayee;
    toJSON(message: RegisteredCounterpartyPayee): unknown;
    fromPartial<I extends {
        channelId?: string | undefined;
        relayer?: string | undefined;
        counterpartyPayee?: string | undefined;
    } & {
        channelId?: string | undefined;
        relayer?: string | undefined;
        counterpartyPayee?: string | undefined;
    } & Record<Exclude<keyof I, keyof RegisteredCounterpartyPayee>, never>>(object: I): RegisteredCounterpartyPayee;
};
export declare const ForwardRelayerAddress: {
    typeUrl: string;
    encode(message: ForwardRelayerAddress, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ForwardRelayerAddress;
    fromJSON(object: any): ForwardRelayerAddress;
    toJSON(message: ForwardRelayerAddress): unknown;
    fromPartial<I extends {
        address?: string | undefined;
        packetId?: {
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
        } | undefined;
    } & {
        address?: string | undefined;
        packetId?: ({
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
        } & {
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
        } & Record<Exclude<keyof I["packetId"], keyof PacketId>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof ForwardRelayerAddress>, never>>(object: I): ForwardRelayerAddress;
};
