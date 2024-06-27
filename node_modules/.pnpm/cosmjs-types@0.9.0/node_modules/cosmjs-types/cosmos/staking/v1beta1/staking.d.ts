import { Header } from "../../../tendermint/types/types";
import { Timestamp } from "../../../google/protobuf/timestamp";
import { Any } from "../../../google/protobuf/any";
import { Duration } from "../../../google/protobuf/duration";
import { Coin } from "../../base/v1beta1/coin";
import { ValidatorUpdate } from "../../../tendermint/abci/types";
import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmos.staking.v1beta1";
/** BondStatus is the status of a validator. */
export declare enum BondStatus {
    /** BOND_STATUS_UNSPECIFIED - UNSPECIFIED defines an invalid validator status. */
    BOND_STATUS_UNSPECIFIED = 0,
    /** BOND_STATUS_UNBONDED - UNBONDED defines a validator that is not bonded. */
    BOND_STATUS_UNBONDED = 1,
    /** BOND_STATUS_UNBONDING - UNBONDING defines a validator that is unbonding. */
    BOND_STATUS_UNBONDING = 2,
    /** BOND_STATUS_BONDED - BONDED defines a validator that is bonded. */
    BOND_STATUS_BONDED = 3,
    UNRECOGNIZED = -1
}
export declare function bondStatusFromJSON(object: any): BondStatus;
export declare function bondStatusToJSON(object: BondStatus): string;
/** Infraction indicates the infraction a validator commited. */
export declare enum Infraction {
    /** INFRACTION_UNSPECIFIED - UNSPECIFIED defines an empty infraction. */
    INFRACTION_UNSPECIFIED = 0,
    /** INFRACTION_DOUBLE_SIGN - DOUBLE_SIGN defines a validator that double-signs a block. */
    INFRACTION_DOUBLE_SIGN = 1,
    /** INFRACTION_DOWNTIME - DOWNTIME defines a validator that missed signing too many blocks. */
    INFRACTION_DOWNTIME = 2,
    UNRECOGNIZED = -1
}
export declare function infractionFromJSON(object: any): Infraction;
export declare function infractionToJSON(object: Infraction): string;
/**
 * HistoricalInfo contains header and validator information for a given block.
 * It is stored as part of staking module's state, which persists the `n` most
 * recent HistoricalInfo
 * (`n` is set by the staking module's `historical_entries` parameter).
 */
export interface HistoricalInfo {
    header: Header;
    valset: Validator[];
}
/**
 * CommissionRates defines the initial commission rates to be used for creating
 * a validator.
 */
export interface CommissionRates {
    /** rate is the commission rate charged to delegators, as a fraction. */
    rate: string;
    /** max_rate defines the maximum commission rate which validator can ever charge, as a fraction. */
    maxRate: string;
    /** max_change_rate defines the maximum daily increase of the validator commission, as a fraction. */
    maxChangeRate: string;
}
/** Commission defines commission parameters for a given validator. */
export interface Commission {
    /** commission_rates defines the initial commission rates to be used for creating a validator. */
    commissionRates: CommissionRates;
    /** update_time is the last time the commission rate was changed. */
    updateTime: Timestamp;
}
/** Description defines a validator description. */
export interface Description {
    /** moniker defines a human-readable name for the validator. */
    moniker: string;
    /** identity defines an optional identity signature (ex. UPort or Keybase). */
    identity: string;
    /** website defines an optional website link. */
    website: string;
    /** security_contact defines an optional email for security contact. */
    securityContact: string;
    /** details define other optional details. */
    details: string;
}
/**
 * Validator defines a validator, together with the total amount of the
 * Validator's bond shares and their exchange rate to coins. Slashing results in
 * a decrease in the exchange rate, allowing correct calculation of future
 * undelegations without iterating over delegators. When coins are delegated to
 * this validator, the validator is credited with a delegation whose number of
 * bond shares is based on the amount of coins delegated divided by the current
 * exchange rate. Voting power can be calculated as total bonded shares
 * multiplied by exchange rate.
 */
export interface Validator {
    /** operator_address defines the address of the validator's operator; bech encoded in JSON. */
    operatorAddress: string;
    /** consensus_pubkey is the consensus public key of the validator, as a Protobuf Any. */
    consensusPubkey?: Any;
    /** jailed defined whether the validator has been jailed from bonded status or not. */
    jailed: boolean;
    /** status is the validator status (bonded/unbonding/unbonded). */
    status: BondStatus;
    /** tokens define the delegated tokens (incl. self-delegation). */
    tokens: string;
    /** delegator_shares defines total shares issued to a validator's delegators. */
    delegatorShares: string;
    /** description defines the description terms for the validator. */
    description: Description;
    /** unbonding_height defines, if unbonding, the height at which this validator has begun unbonding. */
    unbondingHeight: bigint;
    /** unbonding_time defines, if unbonding, the min time for the validator to complete unbonding. */
    unbondingTime: Timestamp;
    /** commission defines the commission parameters. */
    commission: Commission;
    /**
     * min_self_delegation is the validator's self declared minimum self delegation.
     *
     * Since: cosmos-sdk 0.46
     */
    minSelfDelegation: string;
    /** strictly positive if this validator's unbonding has been stopped by external modules */
    unbondingOnHoldRefCount: bigint;
    /** list of unbonding ids, each uniquely identifing an unbonding of this validator */
    unbondingIds: bigint[];
}
/** ValAddresses defines a repeated set of validator addresses. */
export interface ValAddresses {
    addresses: string[];
}
/**
 * DVPair is struct that just has a delegator-validator pair with no other data.
 * It is intended to be used as a marshalable pointer. For example, a DVPair can
 * be used to construct the key to getting an UnbondingDelegation from state.
 */
export interface DVPair {
    delegatorAddress: string;
    validatorAddress: string;
}
/** DVPairs defines an array of DVPair objects. */
export interface DVPairs {
    pairs: DVPair[];
}
/**
 * DVVTriplet is struct that just has a delegator-validator-validator triplet
 * with no other data. It is intended to be used as a marshalable pointer. For
 * example, a DVVTriplet can be used to construct the key to getting a
 * Redelegation from state.
 */
export interface DVVTriplet {
    delegatorAddress: string;
    validatorSrcAddress: string;
    validatorDstAddress: string;
}
/** DVVTriplets defines an array of DVVTriplet objects. */
export interface DVVTriplets {
    triplets: DVVTriplet[];
}
/**
 * Delegation represents the bond with tokens held by an account. It is
 * owned by one delegator, and is associated with the voting power of one
 * validator.
 */
export interface Delegation {
    /** delegator_address is the bech32-encoded address of the delegator. */
    delegatorAddress: string;
    /** validator_address is the bech32-encoded address of the validator. */
    validatorAddress: string;
    /** shares define the delegation shares received. */
    shares: string;
}
/**
 * UnbondingDelegation stores all of a single delegator's unbonding bonds
 * for a single validator in an time-ordered list.
 */
export interface UnbondingDelegation {
    /** delegator_address is the bech32-encoded address of the delegator. */
    delegatorAddress: string;
    /** validator_address is the bech32-encoded address of the validator. */
    validatorAddress: string;
    /** entries are the unbonding delegation entries. */
    entries: UnbondingDelegationEntry[];
}
/** UnbondingDelegationEntry defines an unbonding object with relevant metadata. */
export interface UnbondingDelegationEntry {
    /** creation_height is the height which the unbonding took place. */
    creationHeight: bigint;
    /** completion_time is the unix time for unbonding completion. */
    completionTime: Timestamp;
    /** initial_balance defines the tokens initially scheduled to receive at completion. */
    initialBalance: string;
    /** balance defines the tokens to receive at completion. */
    balance: string;
    /** Incrementing id that uniquely identifies this entry */
    unbondingId: bigint;
    /** Strictly positive if this entry's unbonding has been stopped by external modules */
    unbondingOnHoldRefCount: bigint;
}
/** RedelegationEntry defines a redelegation object with relevant metadata. */
export interface RedelegationEntry {
    /** creation_height  defines the height which the redelegation took place. */
    creationHeight: bigint;
    /** completion_time defines the unix time for redelegation completion. */
    completionTime: Timestamp;
    /** initial_balance defines the initial balance when redelegation started. */
    initialBalance: string;
    /** shares_dst is the amount of destination-validator shares created by redelegation. */
    sharesDst: string;
    /** Incrementing id that uniquely identifies this entry */
    unbondingId: bigint;
    /** Strictly positive if this entry's unbonding has been stopped by external modules */
    unbondingOnHoldRefCount: bigint;
}
/**
 * Redelegation contains the list of a particular delegator's redelegating bonds
 * from a particular source validator to a particular destination validator.
 */
export interface Redelegation {
    /** delegator_address is the bech32-encoded address of the delegator. */
    delegatorAddress: string;
    /** validator_src_address is the validator redelegation source operator address. */
    validatorSrcAddress: string;
    /** validator_dst_address is the validator redelegation destination operator address. */
    validatorDstAddress: string;
    /** entries are the redelegation entries. */
    entries: RedelegationEntry[];
}
/** Params defines the parameters for the x/staking module. */
export interface Params {
    /** unbonding_time is the time duration of unbonding. */
    unbondingTime: Duration;
    /** max_validators is the maximum number of validators. */
    maxValidators: number;
    /** max_entries is the max entries for either unbonding delegation or redelegation (per pair/trio). */
    maxEntries: number;
    /** historical_entries is the number of historical entries to persist. */
    historicalEntries: number;
    /** bond_denom defines the bondable coin denomination. */
    bondDenom: string;
    /** min_commission_rate is the chain-wide minimum commission rate that a validator can charge their delegators */
    minCommissionRate: string;
}
/**
 * DelegationResponse is equivalent to Delegation except that it contains a
 * balance in addition to shares which is more suitable for client responses.
 */
export interface DelegationResponse {
    delegation: Delegation;
    balance: Coin;
}
/**
 * RedelegationEntryResponse is equivalent to a RedelegationEntry except that it
 * contains a balance in addition to shares which is more suitable for client
 * responses.
 */
export interface RedelegationEntryResponse {
    redelegationEntry: RedelegationEntry;
    balance: string;
}
/**
 * RedelegationResponse is equivalent to a Redelegation except that its entries
 * contain a balance in addition to shares which is more suitable for client
 * responses.
 */
export interface RedelegationResponse {
    redelegation: Redelegation;
    entries: RedelegationEntryResponse[];
}
/**
 * Pool is used for tracking bonded and not-bonded token supply of the bond
 * denomination.
 */
export interface Pool {
    notBondedTokens: string;
    bondedTokens: string;
}
/**
 * ValidatorUpdates defines an array of abci.ValidatorUpdate objects.
 * TODO: explore moving this to proto/cosmos/base to separate modules from tendermint dependence
 */
export interface ValidatorUpdates {
    updates: ValidatorUpdate[];
}
export declare const HistoricalInfo: {
    typeUrl: string;
    encode(message: HistoricalInfo, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): HistoricalInfo;
    fromJSON(object: any): HistoricalInfo;
    toJSON(message: HistoricalInfo): unknown;
    fromPartial<I extends {
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
        valset?: {
            operatorAddress?: string | undefined;
            consensusPubkey?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            jailed?: boolean | undefined;
            status?: BondStatus | undefined;
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
            } & Record<Exclude<keyof I["header"]["version"], keyof import("../../../tendermint/version/types").Consensus>, never>) | undefined;
            chainId?: string | undefined;
            height?: bigint | undefined;
            time?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["header"]["time"], keyof Timestamp>, never>) | undefined;
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
                } & Record<Exclude<keyof I["header"]["lastBlockId"]["partSetHeader"], keyof import("../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
            } & Record<Exclude<keyof I["header"]["lastBlockId"], keyof import("../../../tendermint/types/types").BlockID>, never>) | undefined;
            lastCommitHash?: Uint8Array | undefined;
            dataHash?: Uint8Array | undefined;
            validatorsHash?: Uint8Array | undefined;
            nextValidatorsHash?: Uint8Array | undefined;
            consensusHash?: Uint8Array | undefined;
            appHash?: Uint8Array | undefined;
            lastResultsHash?: Uint8Array | undefined;
            evidenceHash?: Uint8Array | undefined;
            proposerAddress?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["header"], keyof Header>, never>) | undefined;
        valset?: ({
            operatorAddress?: string | undefined;
            consensusPubkey?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            jailed?: boolean | undefined;
            status?: BondStatus | undefined;
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
            status?: BondStatus | undefined;
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
            } & Record<Exclude<keyof I["valset"][number]["consensusPubkey"], keyof Any>, never>) | undefined;
            jailed?: boolean | undefined;
            status?: BondStatus | undefined;
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
            } & Record<Exclude<keyof I["valset"][number]["description"], keyof Description>, never>) | undefined;
            unbondingHeight?: bigint | undefined;
            unbondingTime?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["valset"][number]["unbondingTime"], keyof Timestamp>, never>) | undefined;
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
                } & Record<Exclude<keyof I["valset"][number]["commission"]["commissionRates"], keyof CommissionRates>, never>) | undefined;
                updateTime?: ({
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & Record<Exclude<keyof I["valset"][number]["commission"]["updateTime"], keyof Timestamp>, never>) | undefined;
            } & Record<Exclude<keyof I["valset"][number]["commission"], keyof Commission>, never>) | undefined;
            minSelfDelegation?: string | undefined;
            unbondingOnHoldRefCount?: bigint | undefined;
            unbondingIds?: (bigint[] & bigint[] & Record<Exclude<keyof I["valset"][number]["unbondingIds"], keyof bigint[]>, never>) | undefined;
        } & Record<Exclude<keyof I["valset"][number], keyof Validator>, never>)[] & Record<Exclude<keyof I["valset"], keyof {
            operatorAddress?: string | undefined;
            consensusPubkey?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            jailed?: boolean | undefined;
            status?: BondStatus | undefined;
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
    } & Record<Exclude<keyof I, keyof HistoricalInfo>, never>>(object: I): HistoricalInfo;
};
export declare const CommissionRates: {
    typeUrl: string;
    encode(message: CommissionRates, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): CommissionRates;
    fromJSON(object: any): CommissionRates;
    toJSON(message: CommissionRates): unknown;
    fromPartial<I extends {
        rate?: string | undefined;
        maxRate?: string | undefined;
        maxChangeRate?: string | undefined;
    } & {
        rate?: string | undefined;
        maxRate?: string | undefined;
        maxChangeRate?: string | undefined;
    } & Record<Exclude<keyof I, keyof CommissionRates>, never>>(object: I): CommissionRates;
};
export declare const Commission: {
    typeUrl: string;
    encode(message: Commission, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Commission;
    fromJSON(object: any): Commission;
    toJSON(message: Commission): unknown;
    fromPartial<I extends {
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
        } & Record<Exclude<keyof I["commissionRates"], keyof CommissionRates>, never>) | undefined;
        updateTime?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["updateTime"], keyof Timestamp>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof Commission>, never>>(object: I): Commission;
};
export declare const Description: {
    typeUrl: string;
    encode(message: Description, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Description;
    fromJSON(object: any): Description;
    toJSON(message: Description): unknown;
    fromPartial<I extends {
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
    } & Record<Exclude<keyof I, keyof Description>, never>>(object: I): Description;
};
export declare const Validator: {
    typeUrl: string;
    encode(message: Validator, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Validator;
    fromJSON(object: any): Validator;
    toJSON(message: Validator): unknown;
    fromPartial<I extends {
        operatorAddress?: string | undefined;
        consensusPubkey?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
        jailed?: boolean | undefined;
        status?: BondStatus | undefined;
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
        } & Record<Exclude<keyof I["consensusPubkey"], keyof Any>, never>) | undefined;
        jailed?: boolean | undefined;
        status?: BondStatus | undefined;
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
        } & Record<Exclude<keyof I["description"], keyof Description>, never>) | undefined;
        unbondingHeight?: bigint | undefined;
        unbondingTime?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["unbondingTime"], keyof Timestamp>, never>) | undefined;
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
            } & Record<Exclude<keyof I["commission"]["commissionRates"], keyof CommissionRates>, never>) | undefined;
            updateTime?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["commission"]["updateTime"], keyof Timestamp>, never>) | undefined;
        } & Record<Exclude<keyof I["commission"], keyof Commission>, never>) | undefined;
        minSelfDelegation?: string | undefined;
        unbondingOnHoldRefCount?: bigint | undefined;
        unbondingIds?: (bigint[] & bigint[] & Record<Exclude<keyof I["unbondingIds"], keyof bigint[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof Validator>, never>>(object: I): Validator;
};
export declare const ValAddresses: {
    typeUrl: string;
    encode(message: ValAddresses, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ValAddresses;
    fromJSON(object: any): ValAddresses;
    toJSON(message: ValAddresses): unknown;
    fromPartial<I extends {
        addresses?: string[] | undefined;
    } & {
        addresses?: (string[] & string[] & Record<Exclude<keyof I["addresses"], keyof string[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "addresses">, never>>(object: I): ValAddresses;
};
export declare const DVPair: {
    typeUrl: string;
    encode(message: DVPair, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): DVPair;
    fromJSON(object: any): DVPair;
    toJSON(message: DVPair): unknown;
    fromPartial<I extends {
        delegatorAddress?: string | undefined;
        validatorAddress?: string | undefined;
    } & {
        delegatorAddress?: string | undefined;
        validatorAddress?: string | undefined;
    } & Record<Exclude<keyof I, keyof DVPair>, never>>(object: I): DVPair;
};
export declare const DVPairs: {
    typeUrl: string;
    encode(message: DVPairs, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): DVPairs;
    fromJSON(object: any): DVPairs;
    toJSON(message: DVPairs): unknown;
    fromPartial<I extends {
        pairs?: {
            delegatorAddress?: string | undefined;
            validatorAddress?: string | undefined;
        }[] | undefined;
    } & {
        pairs?: ({
            delegatorAddress?: string | undefined;
            validatorAddress?: string | undefined;
        }[] & ({
            delegatorAddress?: string | undefined;
            validatorAddress?: string | undefined;
        } & {
            delegatorAddress?: string | undefined;
            validatorAddress?: string | undefined;
        } & Record<Exclude<keyof I["pairs"][number], keyof DVPair>, never>)[] & Record<Exclude<keyof I["pairs"], keyof {
            delegatorAddress?: string | undefined;
            validatorAddress?: string | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "pairs">, never>>(object: I): DVPairs;
};
export declare const DVVTriplet: {
    typeUrl: string;
    encode(message: DVVTriplet, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): DVVTriplet;
    fromJSON(object: any): DVVTriplet;
    toJSON(message: DVVTriplet): unknown;
    fromPartial<I extends {
        delegatorAddress?: string | undefined;
        validatorSrcAddress?: string | undefined;
        validatorDstAddress?: string | undefined;
    } & {
        delegatorAddress?: string | undefined;
        validatorSrcAddress?: string | undefined;
        validatorDstAddress?: string | undefined;
    } & Record<Exclude<keyof I, keyof DVVTriplet>, never>>(object: I): DVVTriplet;
};
export declare const DVVTriplets: {
    typeUrl: string;
    encode(message: DVVTriplets, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): DVVTriplets;
    fromJSON(object: any): DVVTriplets;
    toJSON(message: DVVTriplets): unknown;
    fromPartial<I extends {
        triplets?: {
            delegatorAddress?: string | undefined;
            validatorSrcAddress?: string | undefined;
            validatorDstAddress?: string | undefined;
        }[] | undefined;
    } & {
        triplets?: ({
            delegatorAddress?: string | undefined;
            validatorSrcAddress?: string | undefined;
            validatorDstAddress?: string | undefined;
        }[] & ({
            delegatorAddress?: string | undefined;
            validatorSrcAddress?: string | undefined;
            validatorDstAddress?: string | undefined;
        } & {
            delegatorAddress?: string | undefined;
            validatorSrcAddress?: string | undefined;
            validatorDstAddress?: string | undefined;
        } & Record<Exclude<keyof I["triplets"][number], keyof DVVTriplet>, never>)[] & Record<Exclude<keyof I["triplets"], keyof {
            delegatorAddress?: string | undefined;
            validatorSrcAddress?: string | undefined;
            validatorDstAddress?: string | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "triplets">, never>>(object: I): DVVTriplets;
};
export declare const Delegation: {
    typeUrl: string;
    encode(message: Delegation, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Delegation;
    fromJSON(object: any): Delegation;
    toJSON(message: Delegation): unknown;
    fromPartial<I extends {
        delegatorAddress?: string | undefined;
        validatorAddress?: string | undefined;
        shares?: string | undefined;
    } & {
        delegatorAddress?: string | undefined;
        validatorAddress?: string | undefined;
        shares?: string | undefined;
    } & Record<Exclude<keyof I, keyof Delegation>, never>>(object: I): Delegation;
};
export declare const UnbondingDelegation: {
    typeUrl: string;
    encode(message: UnbondingDelegation, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): UnbondingDelegation;
    fromJSON(object: any): UnbondingDelegation;
    toJSON(message: UnbondingDelegation): unknown;
    fromPartial<I extends {
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
            } & Record<Exclude<keyof I["entries"][number]["completionTime"], keyof Timestamp>, never>) | undefined;
            initialBalance?: string | undefined;
            balance?: string | undefined;
            unbondingId?: bigint | undefined;
            unbondingOnHoldRefCount?: bigint | undefined;
        } & Record<Exclude<keyof I["entries"][number], keyof UnbondingDelegationEntry>, never>)[] & Record<Exclude<keyof I["entries"], keyof {
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
    } & Record<Exclude<keyof I, keyof UnbondingDelegation>, never>>(object: I): UnbondingDelegation;
};
export declare const UnbondingDelegationEntry: {
    typeUrl: string;
    encode(message: UnbondingDelegationEntry, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): UnbondingDelegationEntry;
    fromJSON(object: any): UnbondingDelegationEntry;
    toJSON(message: UnbondingDelegationEntry): unknown;
    fromPartial<I extends {
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
        } & Record<Exclude<keyof I["completionTime"], keyof Timestamp>, never>) | undefined;
        initialBalance?: string | undefined;
        balance?: string | undefined;
        unbondingId?: bigint | undefined;
        unbondingOnHoldRefCount?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof UnbondingDelegationEntry>, never>>(object: I): UnbondingDelegationEntry;
};
export declare const RedelegationEntry: {
    typeUrl: string;
    encode(message: RedelegationEntry, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): RedelegationEntry;
    fromJSON(object: any): RedelegationEntry;
    toJSON(message: RedelegationEntry): unknown;
    fromPartial<I extends {
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
        } & Record<Exclude<keyof I["completionTime"], keyof Timestamp>, never>) | undefined;
        initialBalance?: string | undefined;
        sharesDst?: string | undefined;
        unbondingId?: bigint | undefined;
        unbondingOnHoldRefCount?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof RedelegationEntry>, never>>(object: I): RedelegationEntry;
};
export declare const Redelegation: {
    typeUrl: string;
    encode(message: Redelegation, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Redelegation;
    fromJSON(object: any): Redelegation;
    toJSON(message: Redelegation): unknown;
    fromPartial<I extends {
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
            } & Record<Exclude<keyof I["entries"][number]["completionTime"], keyof Timestamp>, never>) | undefined;
            initialBalance?: string | undefined;
            sharesDst?: string | undefined;
            unbondingId?: bigint | undefined;
            unbondingOnHoldRefCount?: bigint | undefined;
        } & Record<Exclude<keyof I["entries"][number], keyof RedelegationEntry>, never>)[] & Record<Exclude<keyof I["entries"], keyof {
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
    } & Record<Exclude<keyof I, keyof Redelegation>, never>>(object: I): Redelegation;
};
export declare const Params: {
    typeUrl: string;
    encode(message: Params, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Params;
    fromJSON(object: any): Params;
    toJSON(message: Params): unknown;
    fromPartial<I extends {
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
        } & Record<Exclude<keyof I["unbondingTime"], keyof Duration>, never>) | undefined;
        maxValidators?: number | undefined;
        maxEntries?: number | undefined;
        historicalEntries?: number | undefined;
        bondDenom?: string | undefined;
        minCommissionRate?: string | undefined;
    } & Record<Exclude<keyof I, keyof Params>, never>>(object: I): Params;
};
export declare const DelegationResponse: {
    typeUrl: string;
    encode(message: DelegationResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): DelegationResponse;
    fromJSON(object: any): DelegationResponse;
    toJSON(message: DelegationResponse): unknown;
    fromPartial<I extends {
        delegation?: {
            delegatorAddress?: string | undefined;
            validatorAddress?: string | undefined;
            shares?: string | undefined;
        } | undefined;
        balance?: {
            denom?: string | undefined;
            amount?: string | undefined;
        } | undefined;
    } & {
        delegation?: ({
            delegatorAddress?: string | undefined;
            validatorAddress?: string | undefined;
            shares?: string | undefined;
        } & {
            delegatorAddress?: string | undefined;
            validatorAddress?: string | undefined;
            shares?: string | undefined;
        } & Record<Exclude<keyof I["delegation"], keyof Delegation>, never>) | undefined;
        balance?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["balance"], keyof Coin>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof DelegationResponse>, never>>(object: I): DelegationResponse;
};
export declare const RedelegationEntryResponse: {
    typeUrl: string;
    encode(message: RedelegationEntryResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): RedelegationEntryResponse;
    fromJSON(object: any): RedelegationEntryResponse;
    toJSON(message: RedelegationEntryResponse): unknown;
    fromPartial<I extends {
        redelegationEntry?: {
            creationHeight?: bigint | undefined;
            completionTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            initialBalance?: string | undefined;
            sharesDst?: string | undefined;
            unbondingId?: bigint | undefined;
            unbondingOnHoldRefCount?: bigint | undefined;
        } | undefined;
        balance?: string | undefined;
    } & {
        redelegationEntry?: ({
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
            } & Record<Exclude<keyof I["redelegationEntry"]["completionTime"], keyof Timestamp>, never>) | undefined;
            initialBalance?: string | undefined;
            sharesDst?: string | undefined;
            unbondingId?: bigint | undefined;
            unbondingOnHoldRefCount?: bigint | undefined;
        } & Record<Exclude<keyof I["redelegationEntry"], keyof RedelegationEntry>, never>) | undefined;
        balance?: string | undefined;
    } & Record<Exclude<keyof I, keyof RedelegationEntryResponse>, never>>(object: I): RedelegationEntryResponse;
};
export declare const RedelegationResponse: {
    typeUrl: string;
    encode(message: RedelegationResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): RedelegationResponse;
    fromJSON(object: any): RedelegationResponse;
    toJSON(message: RedelegationResponse): unknown;
    fromPartial<I extends {
        redelegation?: {
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
        } | undefined;
        entries?: {
            redelegationEntry?: {
                creationHeight?: bigint | undefined;
                completionTime?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                initialBalance?: string | undefined;
                sharesDst?: string | undefined;
                unbondingId?: bigint | undefined;
                unbondingOnHoldRefCount?: bigint | undefined;
            } | undefined;
            balance?: string | undefined;
        }[] | undefined;
    } & {
        redelegation?: ({
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
                } & Record<Exclude<keyof I["redelegation"]["entries"][number]["completionTime"], keyof Timestamp>, never>) | undefined;
                initialBalance?: string | undefined;
                sharesDst?: string | undefined;
                unbondingId?: bigint | undefined;
                unbondingOnHoldRefCount?: bigint | undefined;
            } & Record<Exclude<keyof I["redelegation"]["entries"][number], keyof RedelegationEntry>, never>)[] & Record<Exclude<keyof I["redelegation"]["entries"], keyof {
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
        } & Record<Exclude<keyof I["redelegation"], keyof Redelegation>, never>) | undefined;
        entries?: ({
            redelegationEntry?: {
                creationHeight?: bigint | undefined;
                completionTime?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                initialBalance?: string | undefined;
                sharesDst?: string | undefined;
                unbondingId?: bigint | undefined;
                unbondingOnHoldRefCount?: bigint | undefined;
            } | undefined;
            balance?: string | undefined;
        }[] & ({
            redelegationEntry?: {
                creationHeight?: bigint | undefined;
                completionTime?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                initialBalance?: string | undefined;
                sharesDst?: string | undefined;
                unbondingId?: bigint | undefined;
                unbondingOnHoldRefCount?: bigint | undefined;
            } | undefined;
            balance?: string | undefined;
        } & {
            redelegationEntry?: ({
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
                } & Record<Exclude<keyof I["entries"][number]["redelegationEntry"]["completionTime"], keyof Timestamp>, never>) | undefined;
                initialBalance?: string | undefined;
                sharesDst?: string | undefined;
                unbondingId?: bigint | undefined;
                unbondingOnHoldRefCount?: bigint | undefined;
            } & Record<Exclude<keyof I["entries"][number]["redelegationEntry"], keyof RedelegationEntry>, never>) | undefined;
            balance?: string | undefined;
        } & Record<Exclude<keyof I["entries"][number], keyof RedelegationEntryResponse>, never>)[] & Record<Exclude<keyof I["entries"], keyof {
            redelegationEntry?: {
                creationHeight?: bigint | undefined;
                completionTime?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                initialBalance?: string | undefined;
                sharesDst?: string | undefined;
                unbondingId?: bigint | undefined;
                unbondingOnHoldRefCount?: bigint | undefined;
            } | undefined;
            balance?: string | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof RedelegationResponse>, never>>(object: I): RedelegationResponse;
};
export declare const Pool: {
    typeUrl: string;
    encode(message: Pool, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Pool;
    fromJSON(object: any): Pool;
    toJSON(message: Pool): unknown;
    fromPartial<I extends {
        notBondedTokens?: string | undefined;
        bondedTokens?: string | undefined;
    } & {
        notBondedTokens?: string | undefined;
        bondedTokens?: string | undefined;
    } & Record<Exclude<keyof I, keyof Pool>, never>>(object: I): Pool;
};
export declare const ValidatorUpdates: {
    typeUrl: string;
    encode(message: ValidatorUpdates, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ValidatorUpdates;
    fromJSON(object: any): ValidatorUpdates;
    toJSON(message: ValidatorUpdates): unknown;
    fromPartial<I extends {
        updates?: {
            pubKey?: {
                ed25519?: Uint8Array | undefined;
                secp256k1?: Uint8Array | undefined;
            } | undefined;
            power?: bigint | undefined;
        }[] | undefined;
    } & {
        updates?: ({
            pubKey?: {
                ed25519?: Uint8Array | undefined;
                secp256k1?: Uint8Array | undefined;
            } | undefined;
            power?: bigint | undefined;
        }[] & ({
            pubKey?: {
                ed25519?: Uint8Array | undefined;
                secp256k1?: Uint8Array | undefined;
            } | undefined;
            power?: bigint | undefined;
        } & {
            pubKey?: ({
                ed25519?: Uint8Array | undefined;
                secp256k1?: Uint8Array | undefined;
            } & {
                ed25519?: Uint8Array | undefined;
                secp256k1?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["updates"][number]["pubKey"], keyof import("../../../tendermint/crypto/keys").PublicKey>, never>) | undefined;
            power?: bigint | undefined;
        } & Record<Exclude<keyof I["updates"][number], keyof ValidatorUpdate>, never>)[] & Record<Exclude<keyof I["updates"], keyof {
            pubKey?: {
                ed25519?: Uint8Array | undefined;
                secp256k1?: Uint8Array | undefined;
            } | undefined;
            power?: bigint | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "updates">, never>>(object: I): ValidatorUpdates;
};
