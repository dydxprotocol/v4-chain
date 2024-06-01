import { Description, CommissionRates, Params } from "./staking";
import { Any } from "../../../google/protobuf/any";
import { Coin } from "../../base/v1beta1/coin";
import { Timestamp } from "../../../google/protobuf/timestamp";
import { BinaryReader, BinaryWriter } from "../../../binary";
import { Rpc } from "../../../helpers";
export declare const protobufPackage = "cosmos.staking.v1beta1";
/** MsgCreateValidator defines a SDK message for creating a new validator. */
export interface MsgCreateValidator {
    description: Description;
    commission: CommissionRates;
    minSelfDelegation: string;
    delegatorAddress: string;
    validatorAddress: string;
    pubkey?: Any;
    value: Coin;
}
/** MsgCreateValidatorResponse defines the Msg/CreateValidator response type. */
export interface MsgCreateValidatorResponse {
}
/** MsgEditValidator defines a SDK message for editing an existing validator. */
export interface MsgEditValidator {
    description: Description;
    validatorAddress: string;
    /**
     * We pass a reference to the new commission rate and min self delegation as
     * it's not mandatory to update. If not updated, the deserialized rate will be
     * zero with no way to distinguish if an update was intended.
     * REF: #2373
     */
    commissionRate: string;
    minSelfDelegation: string;
}
/** MsgEditValidatorResponse defines the Msg/EditValidator response type. */
export interface MsgEditValidatorResponse {
}
/**
 * MsgDelegate defines a SDK message for performing a delegation of coins
 * from a delegator to a validator.
 */
export interface MsgDelegate {
    delegatorAddress: string;
    validatorAddress: string;
    amount: Coin;
}
/** MsgDelegateResponse defines the Msg/Delegate response type. */
export interface MsgDelegateResponse {
}
/**
 * MsgBeginRedelegate defines a SDK message for performing a redelegation
 * of coins from a delegator and source validator to a destination validator.
 */
export interface MsgBeginRedelegate {
    delegatorAddress: string;
    validatorSrcAddress: string;
    validatorDstAddress: string;
    amount: Coin;
}
/** MsgBeginRedelegateResponse defines the Msg/BeginRedelegate response type. */
export interface MsgBeginRedelegateResponse {
    completionTime: Timestamp;
}
/**
 * MsgUndelegate defines a SDK message for performing an undelegation from a
 * delegate and a validator.
 */
export interface MsgUndelegate {
    delegatorAddress: string;
    validatorAddress: string;
    amount: Coin;
}
/** MsgUndelegateResponse defines the Msg/Undelegate response type. */
export interface MsgUndelegateResponse {
    completionTime: Timestamp;
}
/**
 * MsgCancelUnbondingDelegation defines the SDK message for performing a cancel unbonding delegation for delegator
 *
 * Since: cosmos-sdk 0.46
 */
export interface MsgCancelUnbondingDelegation {
    delegatorAddress: string;
    validatorAddress: string;
    /** amount is always less than or equal to unbonding delegation entry balance */
    amount: Coin;
    /** creation_height is the height which the unbonding took place. */
    creationHeight: bigint;
}
/**
 * MsgCancelUnbondingDelegationResponse
 *
 * Since: cosmos-sdk 0.46
 */
export interface MsgCancelUnbondingDelegationResponse {
}
/**
 * MsgUpdateParams is the Msg/UpdateParams request type.
 *
 * Since: cosmos-sdk 0.47
 */
export interface MsgUpdateParams {
    /** authority is the address that controls the module (defaults to x/gov unless overwritten). */
    authority: string;
    /**
     * params defines the x/staking parameters to update.
     *
     * NOTE: All parameters must be supplied.
     */
    params: Params;
}
/**
 * MsgUpdateParamsResponse defines the response structure for executing a
 * MsgUpdateParams message.
 *
 * Since: cosmos-sdk 0.47
 */
export interface MsgUpdateParamsResponse {
}
export declare const MsgCreateValidator: {
    typeUrl: string;
    encode(message: MsgCreateValidator, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgCreateValidator;
    fromJSON(object: any): MsgCreateValidator;
    toJSON(message: MsgCreateValidator): unknown;
    fromPartial<I extends {
        description?: {
            moniker?: string | undefined;
            identity?: string | undefined;
            website?: string | undefined;
            securityContact?: string | undefined;
            details?: string | undefined;
        } | undefined;
        commission?: {
            rate?: string | undefined;
            maxRate?: string | undefined;
            maxChangeRate?: string | undefined;
        } | undefined;
        minSelfDelegation?: string | undefined;
        delegatorAddress?: string | undefined;
        validatorAddress?: string | undefined;
        pubkey?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
        value?: {
            denom?: string | undefined;
            amount?: string | undefined;
        } | undefined;
    } & {
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
        commission?: ({
            rate?: string | undefined;
            maxRate?: string | undefined;
            maxChangeRate?: string | undefined;
        } & {
            rate?: string | undefined;
            maxRate?: string | undefined;
            maxChangeRate?: string | undefined;
        } & Record<Exclude<keyof I["commission"], keyof CommissionRates>, never>) | undefined;
        minSelfDelegation?: string | undefined;
        delegatorAddress?: string | undefined;
        validatorAddress?: string | undefined;
        pubkey?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["pubkey"], keyof Any>, never>) | undefined;
        value?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["value"], keyof Coin>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof MsgCreateValidator>, never>>(object: I): MsgCreateValidator;
};
export declare const MsgCreateValidatorResponse: {
    typeUrl: string;
    encode(_: MsgCreateValidatorResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgCreateValidatorResponse;
    fromJSON(_: any): MsgCreateValidatorResponse;
    toJSON(_: MsgCreateValidatorResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgCreateValidatorResponse;
};
export declare const MsgEditValidator: {
    typeUrl: string;
    encode(message: MsgEditValidator, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgEditValidator;
    fromJSON(object: any): MsgEditValidator;
    toJSON(message: MsgEditValidator): unknown;
    fromPartial<I extends {
        description?: {
            moniker?: string | undefined;
            identity?: string | undefined;
            website?: string | undefined;
            securityContact?: string | undefined;
            details?: string | undefined;
        } | undefined;
        validatorAddress?: string | undefined;
        commissionRate?: string | undefined;
        minSelfDelegation?: string | undefined;
    } & {
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
        validatorAddress?: string | undefined;
        commissionRate?: string | undefined;
        minSelfDelegation?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgEditValidator>, never>>(object: I): MsgEditValidator;
};
export declare const MsgEditValidatorResponse: {
    typeUrl: string;
    encode(_: MsgEditValidatorResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgEditValidatorResponse;
    fromJSON(_: any): MsgEditValidatorResponse;
    toJSON(_: MsgEditValidatorResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgEditValidatorResponse;
};
export declare const MsgDelegate: {
    typeUrl: string;
    encode(message: MsgDelegate, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgDelegate;
    fromJSON(object: any): MsgDelegate;
    toJSON(message: MsgDelegate): unknown;
    fromPartial<I extends {
        delegatorAddress?: string | undefined;
        validatorAddress?: string | undefined;
        amount?: {
            denom?: string | undefined;
            amount?: string | undefined;
        } | undefined;
    } & {
        delegatorAddress?: string | undefined;
        validatorAddress?: string | undefined;
        amount?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["amount"], keyof Coin>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof MsgDelegate>, never>>(object: I): MsgDelegate;
};
export declare const MsgDelegateResponse: {
    typeUrl: string;
    encode(_: MsgDelegateResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgDelegateResponse;
    fromJSON(_: any): MsgDelegateResponse;
    toJSON(_: MsgDelegateResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgDelegateResponse;
};
export declare const MsgBeginRedelegate: {
    typeUrl: string;
    encode(message: MsgBeginRedelegate, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgBeginRedelegate;
    fromJSON(object: any): MsgBeginRedelegate;
    toJSON(message: MsgBeginRedelegate): unknown;
    fromPartial<I extends {
        delegatorAddress?: string | undefined;
        validatorSrcAddress?: string | undefined;
        validatorDstAddress?: string | undefined;
        amount?: {
            denom?: string | undefined;
            amount?: string | undefined;
        } | undefined;
    } & {
        delegatorAddress?: string | undefined;
        validatorSrcAddress?: string | undefined;
        validatorDstAddress?: string | undefined;
        amount?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["amount"], keyof Coin>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof MsgBeginRedelegate>, never>>(object: I): MsgBeginRedelegate;
};
export declare const MsgBeginRedelegateResponse: {
    typeUrl: string;
    encode(message: MsgBeginRedelegateResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgBeginRedelegateResponse;
    fromJSON(object: any): MsgBeginRedelegateResponse;
    toJSON(message: MsgBeginRedelegateResponse): unknown;
    fromPartial<I extends {
        completionTime?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
    } & {
        completionTime?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["completionTime"], keyof Timestamp>, never>) | undefined;
    } & Record<Exclude<keyof I, "completionTime">, never>>(object: I): MsgBeginRedelegateResponse;
};
export declare const MsgUndelegate: {
    typeUrl: string;
    encode(message: MsgUndelegate, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgUndelegate;
    fromJSON(object: any): MsgUndelegate;
    toJSON(message: MsgUndelegate): unknown;
    fromPartial<I extends {
        delegatorAddress?: string | undefined;
        validatorAddress?: string | undefined;
        amount?: {
            denom?: string | undefined;
            amount?: string | undefined;
        } | undefined;
    } & {
        delegatorAddress?: string | undefined;
        validatorAddress?: string | undefined;
        amount?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["amount"], keyof Coin>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof MsgUndelegate>, never>>(object: I): MsgUndelegate;
};
export declare const MsgUndelegateResponse: {
    typeUrl: string;
    encode(message: MsgUndelegateResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgUndelegateResponse;
    fromJSON(object: any): MsgUndelegateResponse;
    toJSON(message: MsgUndelegateResponse): unknown;
    fromPartial<I extends {
        completionTime?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
    } & {
        completionTime?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["completionTime"], keyof Timestamp>, never>) | undefined;
    } & Record<Exclude<keyof I, "completionTime">, never>>(object: I): MsgUndelegateResponse;
};
export declare const MsgCancelUnbondingDelegation: {
    typeUrl: string;
    encode(message: MsgCancelUnbondingDelegation, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgCancelUnbondingDelegation;
    fromJSON(object: any): MsgCancelUnbondingDelegation;
    toJSON(message: MsgCancelUnbondingDelegation): unknown;
    fromPartial<I extends {
        delegatorAddress?: string | undefined;
        validatorAddress?: string | undefined;
        amount?: {
            denom?: string | undefined;
            amount?: string | undefined;
        } | undefined;
        creationHeight?: bigint | undefined;
    } & {
        delegatorAddress?: string | undefined;
        validatorAddress?: string | undefined;
        amount?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["amount"], keyof Coin>, never>) | undefined;
        creationHeight?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof MsgCancelUnbondingDelegation>, never>>(object: I): MsgCancelUnbondingDelegation;
};
export declare const MsgCancelUnbondingDelegationResponse: {
    typeUrl: string;
    encode(_: MsgCancelUnbondingDelegationResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgCancelUnbondingDelegationResponse;
    fromJSON(_: any): MsgCancelUnbondingDelegationResponse;
    toJSON(_: MsgCancelUnbondingDelegationResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgCancelUnbondingDelegationResponse;
};
export declare const MsgUpdateParams: {
    typeUrl: string;
    encode(message: MsgUpdateParams, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateParams;
    fromJSON(object: any): MsgUpdateParams;
    toJSON(message: MsgUpdateParams): unknown;
    fromPartial<I extends {
        authority?: string | undefined;
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
    } & {
        authority?: string | undefined;
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
    } & Record<Exclude<keyof I, keyof MsgUpdateParams>, never>>(object: I): MsgUpdateParams;
};
export declare const MsgUpdateParamsResponse: {
    typeUrl: string;
    encode(_: MsgUpdateParamsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateParamsResponse;
    fromJSON(_: any): MsgUpdateParamsResponse;
    toJSON(_: MsgUpdateParamsResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgUpdateParamsResponse;
};
/** Msg defines the staking Msg service. */
export interface Msg {
    /** CreateValidator defines a method for creating a new validator. */
    CreateValidator(request: MsgCreateValidator): Promise<MsgCreateValidatorResponse>;
    /** EditValidator defines a method for editing an existing validator. */
    EditValidator(request: MsgEditValidator): Promise<MsgEditValidatorResponse>;
    /**
     * Delegate defines a method for performing a delegation of coins
     * from a delegator to a validator.
     */
    Delegate(request: MsgDelegate): Promise<MsgDelegateResponse>;
    /**
     * BeginRedelegate defines a method for performing a redelegation
     * of coins from a delegator and source validator to a destination validator.
     */
    BeginRedelegate(request: MsgBeginRedelegate): Promise<MsgBeginRedelegateResponse>;
    /**
     * Undelegate defines a method for performing an undelegation from a
     * delegate and a validator.
     */
    Undelegate(request: MsgUndelegate): Promise<MsgUndelegateResponse>;
    /**
     * CancelUnbondingDelegation defines a method for performing canceling the unbonding delegation
     * and delegate back to previous validator.
     *
     * Since: cosmos-sdk 0.46
     */
    CancelUnbondingDelegation(request: MsgCancelUnbondingDelegation): Promise<MsgCancelUnbondingDelegationResponse>;
    /**
     * UpdateParams defines an operation for updating the x/staking module
     * parameters.
     * Since: cosmos-sdk 0.47
     */
    UpdateParams(request: MsgUpdateParams): Promise<MsgUpdateParamsResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    CreateValidator(request: MsgCreateValidator): Promise<MsgCreateValidatorResponse>;
    EditValidator(request: MsgEditValidator): Promise<MsgEditValidatorResponse>;
    Delegate(request: MsgDelegate): Promise<MsgDelegateResponse>;
    BeginRedelegate(request: MsgBeginRedelegate): Promise<MsgBeginRedelegateResponse>;
    Undelegate(request: MsgUndelegate): Promise<MsgUndelegateResponse>;
    CancelUnbondingDelegation(request: MsgCancelUnbondingDelegation): Promise<MsgCancelUnbondingDelegationResponse>;
    UpdateParams(request: MsgUpdateParams): Promise<MsgUpdateParamsResponse>;
}
