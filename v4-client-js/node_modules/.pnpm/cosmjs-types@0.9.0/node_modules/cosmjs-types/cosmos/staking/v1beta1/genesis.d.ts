import { Params, Validator, Delegation, UnbondingDelegation, Redelegation } from "./staking";
import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmos.staking.v1beta1";
/** GenesisState defines the staking module's genesis state. */
export interface GenesisState {
    /** params defines all the parameters of related to deposit. */
    params: Params;
    /**
     * last_total_power tracks the total amounts of bonded tokens recorded during
     * the previous end block.
     */
    lastTotalPower: Uint8Array;
    /**
     * last_validator_powers is a special index that provides a historical list
     * of the last-block's bonded validators.
     */
    lastValidatorPowers: LastValidatorPower[];
    /** delegations defines the validator set at genesis. */
    validators: Validator[];
    /** delegations defines the delegations active at genesis. */
    delegations: Delegation[];
    /** unbonding_delegations defines the unbonding delegations active at genesis. */
    unbondingDelegations: UnbondingDelegation[];
    /** redelegations defines the redelegations active at genesis. */
    redelegations: Redelegation[];
    exported: boolean;
}
/** LastValidatorPower required for validator set update logic. */
export interface LastValidatorPower {
    /** address is the address of the validator. */
    address: string;
    /** power defines the power of the validator. */
    power: bigint;
}
export declare const GenesisState: {
    typeUrl: string;
    encode(message: GenesisState, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
    fromPartial<I extends {
        params?: {
            unbondingTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            maxValidators?: number | undefined;
            maxEntries?: number | undefined;
            historicalEntries?: number | undefined;
            bondDenom?: string | undefined;
            minCommissionRate?: string | undefined;
        } | undefined;
        lastTotalPower?: Uint8Array | undefined;
        lastValidatorPowers?: {
            address?: string | undefined;
            power?: bigint | undefined;
        }[] | undefined;
        validators?: {
            operatorAddress?: string | undefined;
            consensusPubkey?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            jailed?: boolean | undefined;
            status?: import("./staking").BondStatus | undefined;
            tokens?: string | undefined;
            delegatorShares?: string | undefined;
            description?: {
                moniker?: string | undefined;
                identity?: string | undefined;
                website?: string | undefined;
                securityContact?: string | undefined;
                details?: string | undefined;
            } | undefined;
            unbondingHeight?: bigint | undefined;
            unbondingTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            commission?: {
                commissionRates?: {
                    rate?: string | undefined;
                    maxRate?: string | undefined;
                    maxChangeRate?: string | undefined;
                } | undefined;
                updateTime?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
            } | undefined;
            minSelfDelegation?: string | undefined;
            unbondingOnHoldRefCount?: bigint | undefined;
            unbondingIds?: bigint[] | undefined;
        }[] | undefined;
        delegations?: {
            delegatorAddress?: string | undefined;
            validatorAddress?: string | undefined;
            shares?: string | undefined;
        }[] | undefined;
        unbondingDelegations?: {
            delegatorAddress?: string | undefined;
            validatorAddress?: string | undefined;
            entries?: {
                creationHeight?: bigint | undefined;
                completionTime?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                initialBalance?: string | undefined;
                balance?: string | undefined;
                unbondingId?: bigint | undefined;
                unbondingOnHoldRefCount?: bigint | undefined;
            }[] | undefined;
        }[] | undefined;
        redelegations?: {
            delegatorAddress?: string | undefined;
            validatorSrcAddress?: string | undefined;
            validatorDstAddress?: string | undefined;
            entries?: {
                creationHeight?: bigint | undefined;
                completionTime?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                initialBalance?: string | undefined;
                sharesDst?: string | undefined;
                unbondingId?: bigint | undefined;
                unbondingOnHoldRefCount?: bigint | undefined;
            }[] | undefined;
        }[] | undefined;
        exported?: boolean | undefined;
    } & {
        params?: ({
            unbondingTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            maxValidators?: number | undefined;
            maxEntries?: number | undefined;
            historicalEntries?: number | undefined;
            bondDenom?: string | undefined;
            minCommissionRate?: string | undefined;
        } & {
            unbondingTime?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["params"]["unbondingTime"], keyof import("../../../google/protobuf/duration").Duration>, never>) | undefined;
            maxValidators?: number | undefined;
            maxEntries?: number | undefined;
            historicalEntries?: number | undefined;
            bondDenom?: string | undefined;
            minCommissionRate?: string | undefined;
        } & Record<Exclude<keyof I["params"], keyof Params>, never>) | undefined;
        lastTotalPower?: Uint8Array | undefined;
        lastValidatorPowers?: ({
            address?: string | undefined;
            power?: bigint | undefined;
        }[] & ({
            address?: string | undefined;
            power?: bigint | undefined;
        } & {
            address?: string | undefined;
            power?: bigint | undefined;
        } & Record<Exclude<keyof I["lastValidatorPowers"][number], keyof LastValidatorPower>, never>)[] & Record<Exclude<keyof I["lastValidatorPowers"], keyof {
            address?: string | undefined;
            power?: bigint | undefined;
        }[]>, never>) | undefined;
        validators?: ({
            operatorAddress?: string | undefined;
            consensusPubkey?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            jailed?: boolean | undefined;
            status?: import("./staking").BondStatus | undefined;
            tokens?: string | undefined;
            delegatorShares?: string | undefined;
            description?: {
                moniker?: string | undefined;
                identity?: string | undefined;
                website?: string | undefined;
                securityContact?: string | undefined;
                details?: string | undefined;
            } | undefined;
            unbondingHeight?: bigint | undefined;
            unbondingTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            commission?: {
                commissionRates?: {
                    rate?: string | undefined;
                    maxRate?: string | undefined;
                    maxChangeRate?: string | undefined;
                } | undefined;
                updateTime?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
            } | undefined;
            minSelfDelegation?: string | undefined;
            unbondingOnHoldRefCount?: bigint | undefined;
            unbondingIds?: bigint[] | undefined;
        }[] & ({
            operatorAddress?: string | undefined;
            consensusPubkey?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            jailed?: boolean | undefined;
            status?: import("./staking").BondStatus | undefined;
            tokens?: string | undefined;
            delegatorShares?: string | undefined;
            description?: {
                moniker?: string | undefined;
                identity?: string | undefined;
                website?: string | undefined;
                securityContact?: string | undefined;
                details?: string | undefined;
            } | undefined;
            unbondingHeight?: bigint | undefined;
            unbondingTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            commission?: {
                commissionRates?: {
                    rate?: string | undefined;
                    maxRate?: string | undefined;
                    maxChangeRate?: string | undefined;
                } | undefined;
                updateTime?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
            } | undefined;
            minSelfDelegation?: string | undefined;
            unbondingOnHoldRefCount?: bigint | undefined;
            unbondingIds?: bigint[] | undefined;
        } & {
            operatorAddress?: string | undefined;
            consensusPubkey?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["validators"][number]["consensusPubkey"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
            jailed?: boolean | undefined;
            status?: import("./staking").BondStatus | undefined;
            tokens?: string | undefined;
            delegatorShares?: string | undefined;
            description?: ({
                moniker?: string | undefined;
                identity?: string | undefined;
                website?: string | undefined;
                securityContact?: string | undefined;
                details?: string | undefined;
            } & {
                moniker?: string | undefined;
                identity?: string | undefined;
                website?: string | undefined;
                securityContact?: string | undefined;
                details?: string | undefined;
            } & Record<Exclude<keyof I["validators"][number]["description"], keyof import("./staking").Description>, never>) | undefined;
            unbondingHeight?: bigint | undefined;
            unbondingTime?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["validators"][number]["unbondingTime"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
            commission?: ({
                commissionRates?: {
                    rate?: string | undefined;
                    maxRate?: string | undefined;
                    maxChangeRate?: string | undefined;
                } | undefined;
                updateTime?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
            } & {
                commissionRates?: ({
                    rate?: string | undefined;
                    maxRate?: string | undefined;
                    maxChangeRate?: string | undefined;
                } & {
                    rate?: string | undefined;
                    maxRate?: string | undefined;
                    maxChangeRate?: string | undefined;
                } & Record<Exclude<keyof I["validators"][number]["commission"]["commissionRates"], keyof import("./staking").CommissionRates>, never>) | undefined;
                updateTime?: ({
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & Record<Exclude<keyof I["validators"][number]["commission"]["updateTime"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
            } & Record<Exclude<keyof I["validators"][number]["commission"], keyof import("./staking").Commission>, never>) | undefined;
            minSelfDelegation?: string | undefined;
            unbondingOnHoldRefCount?: bigint | undefined;
            unbondingIds?: (bigint[] & bigint[] & Record<Exclude<keyof I["validators"][number]["unbondingIds"], keyof bigint[]>, never>) | undefined;
        } & Record<Exclude<keyof I["validators"][number], keyof Validator>, never>)[] & Record<Exclude<keyof I["validators"], keyof {
            operatorAddress?: string | undefined;
            consensusPubkey?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            jailed?: boolean | undefined;
            status?: import("./staking").BondStatus | undefined;
            tokens?: string | undefined;
            delegatorShares?: string | undefined;
            description?: {
                moniker?: string | undefined;
                identity?: string | undefined;
                website?: string | undefined;
                securityContact?: string | undefined;
                details?: string | undefined;
            } | undefined;
            unbondingHeight?: bigint | undefined;
            unbondingTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            commission?: {
                commissionRates?: {
                    rate?: string | undefined;
                    maxRate?: string | undefined;
                    maxChangeRate?: string | undefined;
                } | undefined;
                updateTime?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
            } | undefined;
            minSelfDelegation?: string | undefined;
            unbondingOnHoldRefCount?: bigint | undefined;
            unbondingIds?: bigint[] | undefined;
        }[]>, never>) | undefined;
        delegations?: ({
            delegatorAddress?: string | undefined;
            validatorAddress?: string | undefined;
            shares?: string | undefined;
        }[] & ({
            delegatorAddress?: string | undefined;
            validatorAddress?: string | undefined;
            shares?: string | undefined;
        } & {
            delegatorAddress?: string | undefined;
            validatorAddress?: string | undefined;
            shares?: string | undefined;
        } & Record<Exclude<keyof I["delegations"][number], keyof Delegation>, never>)[] & Record<Exclude<keyof I["delegations"], keyof {
            delegatorAddress?: string | undefined;
            validatorAddress?: string | undefined;
            shares?: string | undefined;
        }[]>, never>) | undefined;
        unbondingDelegations?: ({
            delegatorAddress?: string | undefined;
            validatorAddress?: string | undefined;
            entries?: {
                creationHeight?: bigint | undefined;
                completionTime?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                initialBalance?: string | undefined;
                balance?: string | undefined;
                unbondingId?: bigint | undefined;
                unbondingOnHoldRefCount?: bigint | undefined;
            }[] | undefined;
        }[] & ({
            delegatorAddress?: string | undefined;
            validatorAddress?: string | undefined;
            entries?: {
                creationHeight?: bigint | undefined;
                completionTime?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                initialBalance?: string | undefined;
                balance?: string | undefined;
                unbondingId?: bigint | undefined;
                unbondingOnHoldRefCount?: bigint | undefined;
            }[] | undefined;
        } & {
            delegatorAddress?: string | undefined;
            validatorAddress?: string | undefined;
            entries?: ({
                creationHeight?: bigint | undefined;
                completionTime?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                initialBalance?: string | undefined;
                balance?: string | undefined;
                unbondingId?: bigint | undefined;
                unbondingOnHoldRefCount?: bigint | undefined;
            }[] & ({
                creationHeight?: bigint | undefined;
                completionTime?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                initialBalance?: string | undefined;
                balance?: string | undefined;
                unbondingId?: bigint | undefined;
                unbondingOnHoldRefCount?: bigint | undefined;
            } & {
                creationHeight?: bigint | undefined;
                completionTime?: ({
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & Record<Exclude<keyof I["unbondingDelegations"][number]["entries"][number]["completionTime"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                initialBalance?: string | undefined;
                balance?: string | undefined;
                unbondingId?: bigint | undefined;
                unbondingOnHoldRefCount?: bigint | undefined;
            } & Record<Exclude<keyof I["unbondingDelegations"][number]["entries"][number], keyof import("./staking").UnbondingDelegationEntry>, never>)[] & Record<Exclude<keyof I["unbondingDelegations"][number]["entries"], keyof {
                creationHeight?: bigint | undefined;
                completionTime?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                initialBalance?: string | undefined;
                balance?: string | undefined;
                unbondingId?: bigint | undefined;
                unbondingOnHoldRefCount?: bigint | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["unbondingDelegations"][number], keyof UnbondingDelegation>, never>)[] & Record<Exclude<keyof I["unbondingDelegations"], keyof {
            delegatorAddress?: string | undefined;
            validatorAddress?: string | undefined;
            entries?: {
                creationHeight?: bigint | undefined;
                completionTime?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                initialBalance?: string | undefined;
                balance?: string | undefined;
                unbondingId?: bigint | undefined;
                unbondingOnHoldRefCount?: bigint | undefined;
            }[] | undefined;
        }[]>, never>) | undefined;
        redelegations?: ({
            delegatorAddress?: string | undefined;
            validatorSrcAddress?: string | undefined;
            validatorDstAddress?: string | undefined;
            entries?: {
                creationHeight?: bigint | undefined;
                completionTime?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                initialBalance?: string | undefined;
                sharesDst?: string | undefined;
                unbondingId?: bigint | undefined;
                unbondingOnHoldRefCount?: bigint | undefined;
            }[] | undefined;
        }[] & ({
            delegatorAddress?: string | undefined;
            validatorSrcAddress?: string | undefined;
            validatorDstAddress?: string | undefined;
            entries?: {
                creationHeight?: bigint | undefined;
                completionTime?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                initialBalance?: string | undefined;
                sharesDst?: string | undefined;
                unbondingId?: bigint | undefined;
                unbondingOnHoldRefCount?: bigint | undefined;
            }[] | undefined;
        } & {
            delegatorAddress?: string | undefined;
            validatorSrcAddress?: string | undefined;
            validatorDstAddress?: string | undefined;
            entries?: ({
                creationHeight?: bigint | undefined;
                completionTime?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                initialBalance?: string | undefined;
                sharesDst?: string | undefined;
                unbondingId?: bigint | undefined;
                unbondingOnHoldRefCount?: bigint | undefined;
            }[] & ({
                creationHeight?: bigint | undefined;
                completionTime?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                initialBalance?: string | undefined;
                sharesDst?: string | undefined;
                unbondingId?: bigint | undefined;
                unbondingOnHoldRefCount?: bigint | undefined;
            } & {
                creationHeight?: bigint | undefined;
                completionTime?: ({
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & Record<Exclude<keyof I["redelegations"][number]["entries"][number]["completionTime"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                initialBalance?: string | undefined;
                sharesDst?: string | undefined;
                unbondingId?: bigint | undefined;
                unbondingOnHoldRefCount?: bigint | undefined;
            } & Record<Exclude<keyof I["redelegations"][number]["entries"][number], keyof import("./staking").RedelegationEntry>, never>)[] & Record<Exclude<keyof I["redelegations"][number]["entries"], keyof {
                creationHeight?: bigint | undefined;
                completionTime?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                initialBalance?: string | undefined;
                sharesDst?: string | undefined;
                unbondingId?: bigint | undefined;
                unbondingOnHoldRefCount?: bigint | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["redelegations"][number], keyof Redelegation>, never>)[] & Record<Exclude<keyof I["redelegations"], keyof {
            delegatorAddress?: string | undefined;
            validatorSrcAddress?: string | undefined;
            validatorDstAddress?: string | undefined;
            entries?: {
                creationHeight?: bigint | undefined;
                completionTime?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                initialBalance?: string | undefined;
                sharesDst?: string | undefined;
                unbondingId?: bigint | undefined;
                unbondingOnHoldRefCount?: bigint | undefined;
            }[] | undefined;
        }[]>, never>) | undefined;
        exported?: boolean | undefined;
    } & Record<Exclude<keyof I, keyof GenesisState>, never>>(object: I): GenesisState;
};
export declare const LastValidatorPower: {
    typeUrl: string;
    encode(message: LastValidatorPower, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): LastValidatorPower;
    fromJSON(object: any): LastValidatorPower;
    toJSON(message: LastValidatorPower): unknown;
    fromPartial<I extends {
        address?: string | undefined;
        power?: bigint | undefined;
    } & {
        address?: string | undefined;
        power?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof LastValidatorPower>, never>>(object: I): LastValidatorPower;
};
