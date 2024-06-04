import { BaseAccount } from "../../auth/v1beta1/auth";
import { Coin } from "../../base/v1beta1/coin";
import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmos.vesting.v1beta1";
/**
 * BaseVestingAccount implements the VestingAccount interface. It contains all
 * the necessary fields needed for any vesting account implementation.
 */
export interface BaseVestingAccount {
    baseAccount?: BaseAccount;
    originalVesting: Coin[];
    delegatedFree: Coin[];
    delegatedVesting: Coin[];
    /** Vesting end time, as unix timestamp (in seconds). */
    endTime: bigint;
}
/**
 * ContinuousVestingAccount implements the VestingAccount interface. It
 * continuously vests by unlocking coins linearly with respect to time.
 */
export interface ContinuousVestingAccount {
    baseVestingAccount?: BaseVestingAccount;
    /** Vesting start time, as unix timestamp (in seconds). */
    startTime: bigint;
}
/**
 * DelayedVestingAccount implements the VestingAccount interface. It vests all
 * coins after a specific time, but non prior. In other words, it keeps them
 * locked until a specified time.
 */
export interface DelayedVestingAccount {
    baseVestingAccount?: BaseVestingAccount;
}
/** Period defines a length of time and amount of coins that will vest. */
export interface Period {
    /** Period duration in seconds. */
    length: bigint;
    amount: Coin[];
}
/**
 * PeriodicVestingAccount implements the VestingAccount interface. It
 * periodically vests by unlocking coins during each specified period.
 */
export interface PeriodicVestingAccount {
    baseVestingAccount?: BaseVestingAccount;
    startTime: bigint;
    vestingPeriods: Period[];
}
/**
 * PermanentLockedAccount implements the VestingAccount interface. It does
 * not ever release coins, locking them indefinitely. Coins in this account can
 * still be used for delegating and for governance votes even while locked.
 *
 * Since: cosmos-sdk 0.43
 */
export interface PermanentLockedAccount {
    baseVestingAccount?: BaseVestingAccount;
}
export declare const BaseVestingAccount: {
    typeUrl: string;
    encode(message: BaseVestingAccount, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): BaseVestingAccount;
    fromJSON(object: any): BaseVestingAccount;
    toJSON(message: BaseVestingAccount): unknown;
    fromPartial<I extends {
        baseAccount?: {
            address?: string | undefined;
            pubKey?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            accountNumber?: bigint | undefined;
            sequence?: bigint | undefined;
        } | undefined;
        originalVesting?: {
            denom?: string | undefined;
            amount?: string | undefined;
        }[] | undefined;
        delegatedFree?: {
            denom?: string | undefined;
            amount?: string | undefined;
        }[] | undefined;
        delegatedVesting?: {
            denom?: string | undefined;
            amount?: string | undefined;
        }[] | undefined;
        endTime?: bigint | undefined;
    } & {
        baseAccount?: ({
            address?: string | undefined;
            pubKey?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            accountNumber?: bigint | undefined;
            sequence?: bigint | undefined;
        } & {
            address?: string | undefined;
            pubKey?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["baseAccount"]["pubKey"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
            accountNumber?: bigint | undefined;
            sequence?: bigint | undefined;
        } & Record<Exclude<keyof I["baseAccount"], keyof BaseAccount>, never>) | undefined;
        originalVesting?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        }[] & ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["originalVesting"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["originalVesting"], keyof {
            denom?: string | undefined;
            amount?: string | undefined;
        }[]>, never>) | undefined;
        delegatedFree?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        }[] & ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["delegatedFree"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["delegatedFree"], keyof {
            denom?: string | undefined;
            amount?: string | undefined;
        }[]>, never>) | undefined;
        delegatedVesting?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        }[] & ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["delegatedVesting"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["delegatedVesting"], keyof {
            denom?: string | undefined;
            amount?: string | undefined;
        }[]>, never>) | undefined;
        endTime?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof BaseVestingAccount>, never>>(object: I): BaseVestingAccount;
};
export declare const ContinuousVestingAccount: {
    typeUrl: string;
    encode(message: ContinuousVestingAccount, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ContinuousVestingAccount;
    fromJSON(object: any): ContinuousVestingAccount;
    toJSON(message: ContinuousVestingAccount): unknown;
    fromPartial<I extends {
        baseVestingAccount?: {
            baseAccount?: {
                address?: string | undefined;
                pubKey?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } | undefined;
                accountNumber?: bigint | undefined;
                sequence?: bigint | undefined;
            } | undefined;
            originalVesting?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            delegatedFree?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            delegatedVesting?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            endTime?: bigint | undefined;
        } | undefined;
        startTime?: bigint | undefined;
    } & {
        baseVestingAccount?: ({
            baseAccount?: {
                address?: string | undefined;
                pubKey?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } | undefined;
                accountNumber?: bigint | undefined;
                sequence?: bigint | undefined;
            } | undefined;
            originalVesting?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            delegatedFree?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            delegatedVesting?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            endTime?: bigint | undefined;
        } & {
            baseAccount?: ({
                address?: string | undefined;
                pubKey?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } | undefined;
                accountNumber?: bigint | undefined;
                sequence?: bigint | undefined;
            } & {
                address?: string | undefined;
                pubKey?: ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["baseVestingAccount"]["baseAccount"]["pubKey"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
                accountNumber?: bigint | undefined;
                sequence?: bigint | undefined;
            } & Record<Exclude<keyof I["baseVestingAccount"]["baseAccount"], keyof BaseAccount>, never>) | undefined;
            originalVesting?: ({
                denom?: string | undefined;
                amount?: string | undefined;
            }[] & ({
                denom?: string | undefined;
                amount?: string | undefined;
            } & {
                denom?: string | undefined;
                amount?: string | undefined;
            } & Record<Exclude<keyof I["baseVestingAccount"]["originalVesting"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["baseVestingAccount"]["originalVesting"], keyof {
                denom?: string | undefined;
                amount?: string | undefined;
            }[]>, never>) | undefined;
            delegatedFree?: ({
                denom?: string | undefined;
                amount?: string | undefined;
            }[] & ({
                denom?: string | undefined;
                amount?: string | undefined;
            } & {
                denom?: string | undefined;
                amount?: string | undefined;
            } & Record<Exclude<keyof I["baseVestingAccount"]["delegatedFree"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["baseVestingAccount"]["delegatedFree"], keyof {
                denom?: string | undefined;
                amount?: string | undefined;
            }[]>, never>) | undefined;
            delegatedVesting?: ({
                denom?: string | undefined;
                amount?: string | undefined;
            }[] & ({
                denom?: string | undefined;
                amount?: string | undefined;
            } & {
                denom?: string | undefined;
                amount?: string | undefined;
            } & Record<Exclude<keyof I["baseVestingAccount"]["delegatedVesting"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["baseVestingAccount"]["delegatedVesting"], keyof {
                denom?: string | undefined;
                amount?: string | undefined;
            }[]>, never>) | undefined;
            endTime?: bigint | undefined;
        } & Record<Exclude<keyof I["baseVestingAccount"], keyof BaseVestingAccount>, never>) | undefined;
        startTime?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof ContinuousVestingAccount>, never>>(object: I): ContinuousVestingAccount;
};
export declare const DelayedVestingAccount: {
    typeUrl: string;
    encode(message: DelayedVestingAccount, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): DelayedVestingAccount;
    fromJSON(object: any): DelayedVestingAccount;
    toJSON(message: DelayedVestingAccount): unknown;
    fromPartial<I extends {
        baseVestingAccount?: {
            baseAccount?: {
                address?: string | undefined;
                pubKey?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } | undefined;
                accountNumber?: bigint | undefined;
                sequence?: bigint | undefined;
            } | undefined;
            originalVesting?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            delegatedFree?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            delegatedVesting?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            endTime?: bigint | undefined;
        } | undefined;
    } & {
        baseVestingAccount?: ({
            baseAccount?: {
                address?: string | undefined;
                pubKey?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } | undefined;
                accountNumber?: bigint | undefined;
                sequence?: bigint | undefined;
            } | undefined;
            originalVesting?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            delegatedFree?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            delegatedVesting?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            endTime?: bigint | undefined;
        } & {
            baseAccount?: ({
                address?: string | undefined;
                pubKey?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } | undefined;
                accountNumber?: bigint | undefined;
                sequence?: bigint | undefined;
            } & {
                address?: string | undefined;
                pubKey?: ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["baseVestingAccount"]["baseAccount"]["pubKey"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
                accountNumber?: bigint | undefined;
                sequence?: bigint | undefined;
            } & Record<Exclude<keyof I["baseVestingAccount"]["baseAccount"], keyof BaseAccount>, never>) | undefined;
            originalVesting?: ({
                denom?: string | undefined;
                amount?: string | undefined;
            }[] & ({
                denom?: string | undefined;
                amount?: string | undefined;
            } & {
                denom?: string | undefined;
                amount?: string | undefined;
            } & Record<Exclude<keyof I["baseVestingAccount"]["originalVesting"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["baseVestingAccount"]["originalVesting"], keyof {
                denom?: string | undefined;
                amount?: string | undefined;
            }[]>, never>) | undefined;
            delegatedFree?: ({
                denom?: string | undefined;
                amount?: string | undefined;
            }[] & ({
                denom?: string | undefined;
                amount?: string | undefined;
            } & {
                denom?: string | undefined;
                amount?: string | undefined;
            } & Record<Exclude<keyof I["baseVestingAccount"]["delegatedFree"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["baseVestingAccount"]["delegatedFree"], keyof {
                denom?: string | undefined;
                amount?: string | undefined;
            }[]>, never>) | undefined;
            delegatedVesting?: ({
                denom?: string | undefined;
                amount?: string | undefined;
            }[] & ({
                denom?: string | undefined;
                amount?: string | undefined;
            } & {
                denom?: string | undefined;
                amount?: string | undefined;
            } & Record<Exclude<keyof I["baseVestingAccount"]["delegatedVesting"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["baseVestingAccount"]["delegatedVesting"], keyof {
                denom?: string | undefined;
                amount?: string | undefined;
            }[]>, never>) | undefined;
            endTime?: bigint | undefined;
        } & Record<Exclude<keyof I["baseVestingAccount"], keyof BaseVestingAccount>, never>) | undefined;
    } & Record<Exclude<keyof I, "baseVestingAccount">, never>>(object: I): DelayedVestingAccount;
};
export declare const Period: {
    typeUrl: string;
    encode(message: Period, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Period;
    fromJSON(object: any): Period;
    toJSON(message: Period): unknown;
    fromPartial<I extends {
        length?: bigint | undefined;
        amount?: {
            denom?: string | undefined;
            amount?: string | undefined;
        }[] | undefined;
    } & {
        length?: bigint | undefined;
        amount?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        }[] & ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["amount"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["amount"], keyof {
            denom?: string | undefined;
            amount?: string | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof Period>, never>>(object: I): Period;
};
export declare const PeriodicVestingAccount: {
    typeUrl: string;
    encode(message: PeriodicVestingAccount, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): PeriodicVestingAccount;
    fromJSON(object: any): PeriodicVestingAccount;
    toJSON(message: PeriodicVestingAccount): unknown;
    fromPartial<I extends {
        baseVestingAccount?: {
            baseAccount?: {
                address?: string | undefined;
                pubKey?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } | undefined;
                accountNumber?: bigint | undefined;
                sequence?: bigint | undefined;
            } | undefined;
            originalVesting?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            delegatedFree?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            delegatedVesting?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            endTime?: bigint | undefined;
        } | undefined;
        startTime?: bigint | undefined;
        vestingPeriods?: {
            length?: bigint | undefined;
            amount?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
        }[] | undefined;
    } & {
        baseVestingAccount?: ({
            baseAccount?: {
                address?: string | undefined;
                pubKey?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } | undefined;
                accountNumber?: bigint | undefined;
                sequence?: bigint | undefined;
            } | undefined;
            originalVesting?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            delegatedFree?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            delegatedVesting?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            endTime?: bigint | undefined;
        } & {
            baseAccount?: ({
                address?: string | undefined;
                pubKey?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } | undefined;
                accountNumber?: bigint | undefined;
                sequence?: bigint | undefined;
            } & {
                address?: string | undefined;
                pubKey?: ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["baseVestingAccount"]["baseAccount"]["pubKey"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
                accountNumber?: bigint | undefined;
                sequence?: bigint | undefined;
            } & Record<Exclude<keyof I["baseVestingAccount"]["baseAccount"], keyof BaseAccount>, never>) | undefined;
            originalVesting?: ({
                denom?: string | undefined;
                amount?: string | undefined;
            }[] & ({
                denom?: string | undefined;
                amount?: string | undefined;
            } & {
                denom?: string | undefined;
                amount?: string | undefined;
            } & Record<Exclude<keyof I["baseVestingAccount"]["originalVesting"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["baseVestingAccount"]["originalVesting"], keyof {
                denom?: string | undefined;
                amount?: string | undefined;
            }[]>, never>) | undefined;
            delegatedFree?: ({
                denom?: string | undefined;
                amount?: string | undefined;
            }[] & ({
                denom?: string | undefined;
                amount?: string | undefined;
            } & {
                denom?: string | undefined;
                amount?: string | undefined;
            } & Record<Exclude<keyof I["baseVestingAccount"]["delegatedFree"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["baseVestingAccount"]["delegatedFree"], keyof {
                denom?: string | undefined;
                amount?: string | undefined;
            }[]>, never>) | undefined;
            delegatedVesting?: ({
                denom?: string | undefined;
                amount?: string | undefined;
            }[] & ({
                denom?: string | undefined;
                amount?: string | undefined;
            } & {
                denom?: string | undefined;
                amount?: string | undefined;
            } & Record<Exclude<keyof I["baseVestingAccount"]["delegatedVesting"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["baseVestingAccount"]["delegatedVesting"], keyof {
                denom?: string | undefined;
                amount?: string | undefined;
            }[]>, never>) | undefined;
            endTime?: bigint | undefined;
        } & Record<Exclude<keyof I["baseVestingAccount"], keyof BaseVestingAccount>, never>) | undefined;
        startTime?: bigint | undefined;
        vestingPeriods?: ({
            length?: bigint | undefined;
            amount?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
        }[] & ({
            length?: bigint | undefined;
            amount?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
        } & {
            length?: bigint | undefined;
            amount?: ({
                denom?: string | undefined;
                amount?: string | undefined;
            }[] & ({
                denom?: string | undefined;
                amount?: string | undefined;
            } & {
                denom?: string | undefined;
                amount?: string | undefined;
            } & Record<Exclude<keyof I["vestingPeriods"][number]["amount"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["vestingPeriods"][number]["amount"], keyof {
                denom?: string | undefined;
                amount?: string | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["vestingPeriods"][number], keyof Period>, never>)[] & Record<Exclude<keyof I["vestingPeriods"], keyof {
            length?: bigint | undefined;
            amount?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof PeriodicVestingAccount>, never>>(object: I): PeriodicVestingAccount;
};
export declare const PermanentLockedAccount: {
    typeUrl: string;
    encode(message: PermanentLockedAccount, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): PermanentLockedAccount;
    fromJSON(object: any): PermanentLockedAccount;
    toJSON(message: PermanentLockedAccount): unknown;
    fromPartial<I extends {
        baseVestingAccount?: {
            baseAccount?: {
                address?: string | undefined;
                pubKey?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } | undefined;
                accountNumber?: bigint | undefined;
                sequence?: bigint | undefined;
            } | undefined;
            originalVesting?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            delegatedFree?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            delegatedVesting?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            endTime?: bigint | undefined;
        } | undefined;
    } & {
        baseVestingAccount?: ({
            baseAccount?: {
                address?: string | undefined;
                pubKey?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } | undefined;
                accountNumber?: bigint | undefined;
                sequence?: bigint | undefined;
            } | undefined;
            originalVesting?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            delegatedFree?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            delegatedVesting?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            endTime?: bigint | undefined;
        } & {
            baseAccount?: ({
                address?: string | undefined;
                pubKey?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } | undefined;
                accountNumber?: bigint | undefined;
                sequence?: bigint | undefined;
            } & {
                address?: string | undefined;
                pubKey?: ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["baseVestingAccount"]["baseAccount"]["pubKey"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
                accountNumber?: bigint | undefined;
                sequence?: bigint | undefined;
            } & Record<Exclude<keyof I["baseVestingAccount"]["baseAccount"], keyof BaseAccount>, never>) | undefined;
            originalVesting?: ({
                denom?: string | undefined;
                amount?: string | undefined;
            }[] & ({
                denom?: string | undefined;
                amount?: string | undefined;
            } & {
                denom?: string | undefined;
                amount?: string | undefined;
            } & Record<Exclude<keyof I["baseVestingAccount"]["originalVesting"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["baseVestingAccount"]["originalVesting"], keyof {
                denom?: string | undefined;
                amount?: string | undefined;
            }[]>, never>) | undefined;
            delegatedFree?: ({
                denom?: string | undefined;
                amount?: string | undefined;
            }[] & ({
                denom?: string | undefined;
                amount?: string | undefined;
            } & {
                denom?: string | undefined;
                amount?: string | undefined;
            } & Record<Exclude<keyof I["baseVestingAccount"]["delegatedFree"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["baseVestingAccount"]["delegatedFree"], keyof {
                denom?: string | undefined;
                amount?: string | undefined;
            }[]>, never>) | undefined;
            delegatedVesting?: ({
                denom?: string | undefined;
                amount?: string | undefined;
            }[] & ({
                denom?: string | undefined;
                amount?: string | undefined;
            } & {
                denom?: string | undefined;
                amount?: string | undefined;
            } & Record<Exclude<keyof I["baseVestingAccount"]["delegatedVesting"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["baseVestingAccount"]["delegatedVesting"], keyof {
                denom?: string | undefined;
                amount?: string | undefined;
            }[]>, never>) | undefined;
            endTime?: bigint | undefined;
        } & Record<Exclude<keyof I["baseVestingAccount"], keyof BaseVestingAccount>, never>) | undefined;
    } & Record<Exclude<keyof I, "baseVestingAccount">, never>>(object: I): PermanentLockedAccount;
};
