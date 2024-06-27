import { Coin } from "../../../../cosmos/base/v1beta1/coin";
import { BinaryReader, BinaryWriter } from "../../../../binary";
export declare const protobufPackage = "ibc.applications.transfer.v1";
/** Allocation defines the spend limit for a particular port and channel */
export interface Allocation {
    /** the port on which the packet will be sent */
    sourcePort: string;
    /** the channel by which the packet will be sent */
    sourceChannel: string;
    /** spend limitation on the channel */
    spendLimit: Coin[];
    /** allow list of receivers, an empty allow list permits any receiver address */
    allowList: string[];
}
/**
 * TransferAuthorization allows the grantee to spend up to spend_limit coins from
 * the granter's account for ibc transfer on a specific channel
 */
export interface TransferAuthorization {
    /** port and channel amounts */
    allocations: Allocation[];
}
export declare const Allocation: {
    typeUrl: string;
    encode(message: Allocation, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Allocation;
    fromJSON(object: any): Allocation;
    toJSON(message: Allocation): unknown;
    fromPartial<I extends {
        sourcePort?: string | undefined;
        sourceChannel?: string | undefined;
        spendLimit?: {
            denom?: string | undefined;
            amount?: string | undefined;
        }[] | undefined;
        allowList?: string[] | undefined;
    } & {
        sourcePort?: string | undefined;
        sourceChannel?: string | undefined;
        spendLimit?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        }[] & ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["spendLimit"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["spendLimit"], keyof {
            denom?: string | undefined;
            amount?: string | undefined;
        }[]>, never>) | undefined;
        allowList?: (string[] & string[] & Record<Exclude<keyof I["allowList"], keyof string[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof Allocation>, never>>(object: I): Allocation;
};
export declare const TransferAuthorization: {
    typeUrl: string;
    encode(message: TransferAuthorization, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): TransferAuthorization;
    fromJSON(object: any): TransferAuthorization;
    toJSON(message: TransferAuthorization): unknown;
    fromPartial<I extends {
        allocations?: {
            sourcePort?: string | undefined;
            sourceChannel?: string | undefined;
            spendLimit?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            allowList?: string[] | undefined;
        }[] | undefined;
    } & {
        allocations?: ({
            sourcePort?: string | undefined;
            sourceChannel?: string | undefined;
            spendLimit?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            allowList?: string[] | undefined;
        }[] & ({
            sourcePort?: string | undefined;
            sourceChannel?: string | undefined;
            spendLimit?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            allowList?: string[] | undefined;
        } & {
            sourcePort?: string | undefined;
            sourceChannel?: string | undefined;
            spendLimit?: ({
                denom?: string | undefined;
                amount?: string | undefined;
            }[] & ({
                denom?: string | undefined;
                amount?: string | undefined;
            } & {
                denom?: string | undefined;
                amount?: string | undefined;
            } & Record<Exclude<keyof I["allocations"][number]["spendLimit"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["allocations"][number]["spendLimit"], keyof {
                denom?: string | undefined;
                amount?: string | undefined;
            }[]>, never>) | undefined;
            allowList?: (string[] & string[] & Record<Exclude<keyof I["allocations"][number]["allowList"], keyof string[]>, never>) | undefined;
        } & Record<Exclude<keyof I["allocations"][number], keyof Allocation>, never>)[] & Record<Exclude<keyof I["allocations"], keyof {
            sourcePort?: string | undefined;
            sourceChannel?: string | undefined;
            spendLimit?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            allowList?: string[] | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "allocations">, never>>(object: I): TransferAuthorization;
};
