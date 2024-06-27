import { Coin } from "../../base/v1beta1/coin";
import { Timestamp } from "../../../google/protobuf/timestamp";
import { Duration } from "../../../google/protobuf/duration";
import { Any } from "../../../google/protobuf/any";
import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmos.feegrant.v1beta1";
/**
 * BasicAllowance implements Allowance with a one-time grant of coins
 * that optionally expires. The grantee can use up to SpendLimit to cover fees.
 */
export interface BasicAllowance {
    /**
     * spend_limit specifies the maximum amount of coins that can be spent
     * by this allowance and will be updated as coins are spent. If it is
     * empty, there is no spend limit and any amount of coins can be spent.
     */
    spendLimit: Coin[];
    /** expiration specifies an optional time when this allowance expires */
    expiration?: Timestamp;
}
/**
 * PeriodicAllowance extends Allowance to allow for both a maximum cap,
 * as well as a limit per time period.
 */
export interface PeriodicAllowance {
    /** basic specifies a struct of `BasicAllowance` */
    basic: BasicAllowance;
    /**
     * period specifies the time duration in which period_spend_limit coins can
     * be spent before that allowance is reset
     */
    period: Duration;
    /**
     * period_spend_limit specifies the maximum number of coins that can be spent
     * in the period
     */
    periodSpendLimit: Coin[];
    /** period_can_spend is the number of coins left to be spent before the period_reset time */
    periodCanSpend: Coin[];
    /**
     * period_reset is the time at which this period resets and a new one begins,
     * it is calculated from the start time of the first transaction after the
     * last period ended
     */
    periodReset: Timestamp;
}
/** AllowedMsgAllowance creates allowance only for specified message types. */
export interface AllowedMsgAllowance {
    /** allowance can be any of basic and periodic fee allowance. */
    allowance?: Any;
    /** allowed_messages are the messages for which the grantee has the access. */
    allowedMessages: string[];
}
/** Grant is stored in the KVStore to record a grant with full context */
export interface Grant {
    /** granter is the address of the user granting an allowance of their funds. */
    granter: string;
    /** grantee is the address of the user being granted an allowance of another user's funds. */
    grantee: string;
    /** allowance can be any of basic, periodic, allowed fee allowance. */
    allowance?: Any;
}
export declare const BasicAllowance: {
    typeUrl: string;
    encode(message: BasicAllowance, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): BasicAllowance;
    fromJSON(object: any): BasicAllowance;
    toJSON(message: BasicAllowance): unknown;
    fromPartial<I extends {
        spendLimit?: {
            denom?: string | undefined;
            amount?: string | undefined;
        }[] | undefined;
        expiration?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
    } & {
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
        expiration?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["expiration"], keyof Timestamp>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof BasicAllowance>, never>>(object: I): BasicAllowance;
};
export declare const PeriodicAllowance: {
    typeUrl: string;
    encode(message: PeriodicAllowance, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): PeriodicAllowance;
    fromJSON(object: any): PeriodicAllowance;
    toJSON(message: PeriodicAllowance): unknown;
    fromPartial<I extends {
        basic?: {
            spendLimit?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            expiration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        } | undefined;
        period?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
        periodSpendLimit?: {
            denom?: string | undefined;
            amount?: string | undefined;
        }[] | undefined;
        periodCanSpend?: {
            denom?: string | undefined;
            amount?: string | undefined;
        }[] | undefined;
        periodReset?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
    } & {
        basic?: ({
            spendLimit?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            expiration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        } & {
            spendLimit?: ({
                denom?: string | undefined;
                amount?: string | undefined;
            }[] & ({
                denom?: string | undefined;
                amount?: string | undefined;
            } & {
                denom?: string | undefined;
                amount?: string | undefined;
            } & Record<Exclude<keyof I["basic"]["spendLimit"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["basic"]["spendLimit"], keyof {
                denom?: string | undefined;
                amount?: string | undefined;
            }[]>, never>) | undefined;
            expiration?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["basic"]["expiration"], keyof Timestamp>, never>) | undefined;
        } & Record<Exclude<keyof I["basic"], keyof BasicAllowance>, never>) | undefined;
        period?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["period"], keyof Duration>, never>) | undefined;
        periodSpendLimit?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        }[] & ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["periodSpendLimit"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["periodSpendLimit"], keyof {
            denom?: string | undefined;
            amount?: string | undefined;
        }[]>, never>) | undefined;
        periodCanSpend?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        }[] & ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["periodCanSpend"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["periodCanSpend"], keyof {
            denom?: string | undefined;
            amount?: string | undefined;
        }[]>, never>) | undefined;
        periodReset?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["periodReset"], keyof Timestamp>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof PeriodicAllowance>, never>>(object: I): PeriodicAllowance;
};
export declare const AllowedMsgAllowance: {
    typeUrl: string;
    encode(message: AllowedMsgAllowance, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): AllowedMsgAllowance;
    fromJSON(object: any): AllowedMsgAllowance;
    toJSON(message: AllowedMsgAllowance): unknown;
    fromPartial<I extends {
        allowance?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
        allowedMessages?: string[] | undefined;
    } & {
        allowance?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["allowance"], keyof Any>, never>) | undefined;
        allowedMessages?: (string[] & string[] & Record<Exclude<keyof I["allowedMessages"], keyof string[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof AllowedMsgAllowance>, never>>(object: I): AllowedMsgAllowance;
};
export declare const Grant: {
    typeUrl: string;
    encode(message: Grant, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Grant;
    fromJSON(object: any): Grant;
    toJSON(message: Grant): unknown;
    fromPartial<I extends {
        granter?: string | undefined;
        grantee?: string | undefined;
        allowance?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
    } & {
        granter?: string | undefined;
        grantee?: string | undefined;
        allowance?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["allowance"], keyof Any>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof Grant>, never>>(object: I): Grant;
};
