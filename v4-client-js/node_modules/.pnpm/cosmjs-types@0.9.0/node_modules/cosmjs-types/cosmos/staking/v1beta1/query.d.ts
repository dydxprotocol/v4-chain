import { PageRequest, PageResponse } from "../../base/query/v1beta1/pagination";
import { Validator, DelegationResponse, UnbondingDelegation, RedelegationResponse, HistoricalInfo, Pool, Params } from "./staking";
import { BinaryReader, BinaryWriter } from "../../../binary";
import { Rpc } from "../../../helpers";
export declare const protobufPackage = "cosmos.staking.v1beta1";
/** QueryValidatorsRequest is request type for Query/Validators RPC method. */
export interface QueryValidatorsRequest {
    /** status enables to query for validators matching a given status. */
    status: string;
    /** pagination defines an optional pagination for the request. */
    pagination?: PageRequest;
}
/** QueryValidatorsResponse is response type for the Query/Validators RPC method */
export interface QueryValidatorsResponse {
    /** validators contains all the queried validators. */
    validators: Validator[];
    /** pagination defines the pagination in the response. */
    pagination?: PageResponse;
}
/** QueryValidatorRequest is response type for the Query/Validator RPC method */
export interface QueryValidatorRequest {
    /** validator_addr defines the validator address to query for. */
    validatorAddr: string;
}
/** QueryValidatorResponse is response type for the Query/Validator RPC method */
export interface QueryValidatorResponse {
    /** validator defines the validator info. */
    validator: Validator;
}
/**
 * QueryValidatorDelegationsRequest is request type for the
 * Query/ValidatorDelegations RPC method
 */
export interface QueryValidatorDelegationsRequest {
    /** validator_addr defines the validator address to query for. */
    validatorAddr: string;
    /** pagination defines an optional pagination for the request. */
    pagination?: PageRequest;
}
/**
 * QueryValidatorDelegationsResponse is response type for the
 * Query/ValidatorDelegations RPC method
 */
export interface QueryValidatorDelegationsResponse {
    delegationResponses: DelegationResponse[];
    /** pagination defines the pagination in the response. */
    pagination?: PageResponse;
}
/**
 * QueryValidatorUnbondingDelegationsRequest is required type for the
 * Query/ValidatorUnbondingDelegations RPC method
 */
export interface QueryValidatorUnbondingDelegationsRequest {
    /** validator_addr defines the validator address to query for. */
    validatorAddr: string;
    /** pagination defines an optional pagination for the request. */
    pagination?: PageRequest;
}
/**
 * QueryValidatorUnbondingDelegationsResponse is response type for the
 * Query/ValidatorUnbondingDelegations RPC method.
 */
export interface QueryValidatorUnbondingDelegationsResponse {
    unbondingResponses: UnbondingDelegation[];
    /** pagination defines the pagination in the response. */
    pagination?: PageResponse;
}
/** QueryDelegationRequest is request type for the Query/Delegation RPC method. */
export interface QueryDelegationRequest {
    /** delegator_addr defines the delegator address to query for. */
    delegatorAddr: string;
    /** validator_addr defines the validator address to query for. */
    validatorAddr: string;
}
/** QueryDelegationResponse is response type for the Query/Delegation RPC method. */
export interface QueryDelegationResponse {
    /** delegation_responses defines the delegation info of a delegation. */
    delegationResponse?: DelegationResponse;
}
/**
 * QueryUnbondingDelegationRequest is request type for the
 * Query/UnbondingDelegation RPC method.
 */
export interface QueryUnbondingDelegationRequest {
    /** delegator_addr defines the delegator address to query for. */
    delegatorAddr: string;
    /** validator_addr defines the validator address to query for. */
    validatorAddr: string;
}
/**
 * QueryDelegationResponse is response type for the Query/UnbondingDelegation
 * RPC method.
 */
export interface QueryUnbondingDelegationResponse {
    /** unbond defines the unbonding information of a delegation. */
    unbond: UnbondingDelegation;
}
/**
 * QueryDelegatorDelegationsRequest is request type for the
 * Query/DelegatorDelegations RPC method.
 */
export interface QueryDelegatorDelegationsRequest {
    /** delegator_addr defines the delegator address to query for. */
    delegatorAddr: string;
    /** pagination defines an optional pagination for the request. */
    pagination?: PageRequest;
}
/**
 * QueryDelegatorDelegationsResponse is response type for the
 * Query/DelegatorDelegations RPC method.
 */
export interface QueryDelegatorDelegationsResponse {
    /** delegation_responses defines all the delegations' info of a delegator. */
    delegationResponses: DelegationResponse[];
    /** pagination defines the pagination in the response. */
    pagination?: PageResponse;
}
/**
 * QueryDelegatorUnbondingDelegationsRequest is request type for the
 * Query/DelegatorUnbondingDelegations RPC method.
 */
export interface QueryDelegatorUnbondingDelegationsRequest {
    /** delegator_addr defines the delegator address to query for. */
    delegatorAddr: string;
    /** pagination defines an optional pagination for the request. */
    pagination?: PageRequest;
}
/**
 * QueryUnbondingDelegatorDelegationsResponse is response type for the
 * Query/UnbondingDelegatorDelegations RPC method.
 */
export interface QueryDelegatorUnbondingDelegationsResponse {
    unbondingResponses: UnbondingDelegation[];
    /** pagination defines the pagination in the response. */
    pagination?: PageResponse;
}
/**
 * QueryRedelegationsRequest is request type for the Query/Redelegations RPC
 * method.
 */
export interface QueryRedelegationsRequest {
    /** delegator_addr defines the delegator address to query for. */
    delegatorAddr: string;
    /** src_validator_addr defines the validator address to redelegate from. */
    srcValidatorAddr: string;
    /** dst_validator_addr defines the validator address to redelegate to. */
    dstValidatorAddr: string;
    /** pagination defines an optional pagination for the request. */
    pagination?: PageRequest;
}
/**
 * QueryRedelegationsResponse is response type for the Query/Redelegations RPC
 * method.
 */
export interface QueryRedelegationsResponse {
    redelegationResponses: RedelegationResponse[];
    /** pagination defines the pagination in the response. */
    pagination?: PageResponse;
}
/**
 * QueryDelegatorValidatorsRequest is request type for the
 * Query/DelegatorValidators RPC method.
 */
export interface QueryDelegatorValidatorsRequest {
    /** delegator_addr defines the delegator address to query for. */
    delegatorAddr: string;
    /** pagination defines an optional pagination for the request. */
    pagination?: PageRequest;
}
/**
 * QueryDelegatorValidatorsResponse is response type for the
 * Query/DelegatorValidators RPC method.
 */
export interface QueryDelegatorValidatorsResponse {
    /** validators defines the validators' info of a delegator. */
    validators: Validator[];
    /** pagination defines the pagination in the response. */
    pagination?: PageResponse;
}
/**
 * QueryDelegatorValidatorRequest is request type for the
 * Query/DelegatorValidator RPC method.
 */
export interface QueryDelegatorValidatorRequest {
    /** delegator_addr defines the delegator address to query for. */
    delegatorAddr: string;
    /** validator_addr defines the validator address to query for. */
    validatorAddr: string;
}
/**
 * QueryDelegatorValidatorResponse response type for the
 * Query/DelegatorValidator RPC method.
 */
export interface QueryDelegatorValidatorResponse {
    /** validator defines the validator info. */
    validator: Validator;
}
/**
 * QueryHistoricalInfoRequest is request type for the Query/HistoricalInfo RPC
 * method.
 */
export interface QueryHistoricalInfoRequest {
    /** height defines at which height to query the historical info. */
    height: bigint;
}
/**
 * QueryHistoricalInfoResponse is response type for the Query/HistoricalInfo RPC
 * method.
 */
export interface QueryHistoricalInfoResponse {
    /** hist defines the historical info at the given height. */
    hist?: HistoricalInfo;
}
/** QueryPoolRequest is request type for the Query/Pool RPC method. */
export interface QueryPoolRequest {
}
/** QueryPoolResponse is response type for the Query/Pool RPC method. */
export interface QueryPoolResponse {
    /** pool defines the pool info. */
    pool: Pool;
}
/** QueryParamsRequest is request type for the Query/Params RPC method. */
export interface QueryParamsRequest {
}
/** QueryParamsResponse is response type for the Query/Params RPC method. */
export interface QueryParamsResponse {
    /** params holds all the parameters of this module. */
    params: Params;
}
export declare const QueryValidatorsRequest: {
    typeUrl: string;
    encode(message: QueryValidatorsRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryValidatorsRequest;
    fromJSON(object: any): QueryValidatorsRequest;
    toJSON(message: QueryValidatorsRequest): unknown;
    fromPartial<I extends {
        status?: string | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        status?: string | undefined;
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
    } & Record<Exclude<keyof I, keyof QueryValidatorsRequest>, never>>(object: I): QueryValidatorsRequest;
};
export declare const QueryValidatorsResponse: {
    typeUrl: string;
    encode(message: QueryValidatorsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryValidatorsResponse;
    fromJSON(object: any): QueryValidatorsResponse;
    toJSON(message: QueryValidatorsResponse): unknown;
    fromPartial<I extends {
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
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
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
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryValidatorsResponse>, never>>(object: I): QueryValidatorsResponse;
};
export declare const QueryValidatorRequest: {
    typeUrl: string;
    encode(message: QueryValidatorRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryValidatorRequest;
    fromJSON(object: any): QueryValidatorRequest;
    toJSON(message: QueryValidatorRequest): unknown;
    fromPartial<I extends {
        validatorAddr?: string | undefined;
    } & {
        validatorAddr?: string | undefined;
    } & Record<Exclude<keyof I, "validatorAddr">, never>>(object: I): QueryValidatorRequest;
};
export declare const QueryValidatorResponse: {
    typeUrl: string;
    encode(message: QueryValidatorResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryValidatorResponse;
    fromJSON(object: any): QueryValidatorResponse;
    toJSON(message: QueryValidatorResponse): unknown;
    fromPartial<I extends {
        validator?: {
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
        } | undefined;
    } & {
        validator?: ({
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
            } & Record<Exclude<keyof I["validator"]["consensusPubkey"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
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
            } & Record<Exclude<keyof I["validator"]["description"], keyof import("./staking").Description>, never>) | undefined;
            unbondingHeight?: bigint | undefined;
            unbondingTime?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["validator"]["unbondingTime"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
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
                } & Record<Exclude<keyof I["validator"]["commission"]["commissionRates"], keyof import("./staking").CommissionRates>, never>) | undefined;
                updateTime?: ({
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & Record<Exclude<keyof I["validator"]["commission"]["updateTime"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
            } & Record<Exclude<keyof I["validator"]["commission"], keyof import("./staking").Commission>, never>) | undefined;
            minSelfDelegation?: string | undefined;
            unbondingOnHoldRefCount?: bigint | undefined;
            unbondingIds?: (bigint[] & bigint[] & Record<Exclude<keyof I["validator"]["unbondingIds"], keyof bigint[]>, never>) | undefined;
        } & Record<Exclude<keyof I["validator"], keyof Validator>, never>) | undefined;
    } & Record<Exclude<keyof I, "validator">, never>>(object: I): QueryValidatorResponse;
};
export declare const QueryValidatorDelegationsRequest: {
    typeUrl: string;
    encode(message: QueryValidatorDelegationsRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryValidatorDelegationsRequest;
    fromJSON(object: any): QueryValidatorDelegationsRequest;
    toJSON(message: QueryValidatorDelegationsRequest): unknown;
    fromPartial<I extends {
        validatorAddr?: string | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        validatorAddr?: string | undefined;
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
    } & Record<Exclude<keyof I, keyof QueryValidatorDelegationsRequest>, never>>(object: I): QueryValidatorDelegationsRequest;
};
export declare const QueryValidatorDelegationsResponse: {
    typeUrl: string;
    encode(message: QueryValidatorDelegationsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryValidatorDelegationsResponse;
    fromJSON(object: any): QueryValidatorDelegationsResponse;
    toJSON(message: QueryValidatorDelegationsResponse): unknown;
    fromPartial<I extends {
        delegationResponses?: {
            delegation?: {
                delegatorAddress?: string | undefined;
                validatorAddress?: string | undefined;
                shares?: string | undefined;
            } | undefined;
            balance?: {
                denom?: string | undefined;
                amount?: string | undefined;
            } | undefined;
        }[] | undefined;
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
        delegationResponses?: ({
            delegation?: {
                delegatorAddress?: string | undefined;
                validatorAddress?: string | undefined;
                shares?: string | undefined;
            } | undefined;
            balance?: {
                denom?: string | undefined;
                amount?: string | undefined;
            } | undefined;
        }[] & ({
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
            } & Record<Exclude<keyof I["delegationResponses"][number]["delegation"], keyof import("./staking").Delegation>, never>) | undefined;
            balance?: ({
                denom?: string | undefined;
                amount?: string | undefined;
            } & {
                denom?: string | undefined;
                amount?: string | undefined;
            } & Record<Exclude<keyof I["delegationResponses"][number]["balance"], keyof import("../../base/v1beta1/coin").Coin>, never>) | undefined;
        } & Record<Exclude<keyof I["delegationResponses"][number], keyof DelegationResponse>, never>)[] & Record<Exclude<keyof I["delegationResponses"], keyof {
            delegation?: {
                delegatorAddress?: string | undefined;
                validatorAddress?: string | undefined;
                shares?: string | undefined;
            } | undefined;
            balance?: {
                denom?: string | undefined;
                amount?: string | undefined;
            } | undefined;
        }[]>, never>) | undefined;
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryValidatorDelegationsResponse>, never>>(object: I): QueryValidatorDelegationsResponse;
};
export declare const QueryValidatorUnbondingDelegationsRequest: {
    typeUrl: string;
    encode(message: QueryValidatorUnbondingDelegationsRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryValidatorUnbondingDelegationsRequest;
    fromJSON(object: any): QueryValidatorUnbondingDelegationsRequest;
    toJSON(message: QueryValidatorUnbondingDelegationsRequest): unknown;
    fromPartial<I extends {
        validatorAddr?: string | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        validatorAddr?: string | undefined;
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
    } & Record<Exclude<keyof I, keyof QueryValidatorUnbondingDelegationsRequest>, never>>(object: I): QueryValidatorUnbondingDelegationsRequest;
};
export declare const QueryValidatorUnbondingDelegationsResponse: {
    typeUrl: string;
    encode(message: QueryValidatorUnbondingDelegationsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryValidatorUnbondingDelegationsResponse;
    fromJSON(object: any): QueryValidatorUnbondingDelegationsResponse;
    toJSON(message: QueryValidatorUnbondingDelegationsResponse): unknown;
    fromPartial<I extends {
        unbondingResponses?: {
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
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
        unbondingResponses?: ({
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
                } & Record<Exclude<keyof I["unbondingResponses"][number]["entries"][number]["completionTime"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                initialBalance?: string | undefined;
                balance?: string | undefined;
                unbondingId?: bigint | undefined;
                unbondingOnHoldRefCount?: bigint | undefined;
            } & Record<Exclude<keyof I["unbondingResponses"][number]["entries"][number], keyof import("./staking").UnbondingDelegationEntry>, never>)[] & Record<Exclude<keyof I["unbondingResponses"][number]["entries"], keyof {
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
        } & Record<Exclude<keyof I["unbondingResponses"][number], keyof UnbondingDelegation>, never>)[] & Record<Exclude<keyof I["unbondingResponses"], keyof {
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
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryValidatorUnbondingDelegationsResponse>, never>>(object: I): QueryValidatorUnbondingDelegationsResponse;
};
export declare const QueryDelegationRequest: {
    typeUrl: string;
    encode(message: QueryDelegationRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryDelegationRequest;
    fromJSON(object: any): QueryDelegationRequest;
    toJSON(message: QueryDelegationRequest): unknown;
    fromPartial<I extends {
        delegatorAddr?: string | undefined;
        validatorAddr?: string | undefined;
    } & {
        delegatorAddr?: string | undefined;
        validatorAddr?: string | undefined;
    } & Record<Exclude<keyof I, keyof QueryDelegationRequest>, never>>(object: I): QueryDelegationRequest;
};
export declare const QueryDelegationResponse: {
    typeUrl: string;
    encode(message: QueryDelegationResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryDelegationResponse;
    fromJSON(object: any): QueryDelegationResponse;
    toJSON(message: QueryDelegationResponse): unknown;
    fromPartial<I extends {
        delegationResponse?: {
            delegation?: {
                delegatorAddress?: string | undefined;
                validatorAddress?: string | undefined;
                shares?: string | undefined;
            } | undefined;
            balance?: {
                denom?: string | undefined;
                amount?: string | undefined;
            } | undefined;
        } | undefined;
    } & {
        delegationResponse?: ({
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
            } & Record<Exclude<keyof I["delegationResponse"]["delegation"], keyof import("./staking").Delegation>, never>) | undefined;
            balance?: ({
                denom?: string | undefined;
                amount?: string | undefined;
            } & {
                denom?: string | undefined;
                amount?: string | undefined;
            } & Record<Exclude<keyof I["delegationResponse"]["balance"], keyof import("../../base/v1beta1/coin").Coin>, never>) | undefined;
        } & Record<Exclude<keyof I["delegationResponse"], keyof DelegationResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, "delegationResponse">, never>>(object: I): QueryDelegationResponse;
};
export declare const QueryUnbondingDelegationRequest: {
    typeUrl: string;
    encode(message: QueryUnbondingDelegationRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryUnbondingDelegationRequest;
    fromJSON(object: any): QueryUnbondingDelegationRequest;
    toJSON(message: QueryUnbondingDelegationRequest): unknown;
    fromPartial<I extends {
        delegatorAddr?: string | undefined;
        validatorAddr?: string | undefined;
    } & {
        delegatorAddr?: string | undefined;
        validatorAddr?: string | undefined;
    } & Record<Exclude<keyof I, keyof QueryUnbondingDelegationRequest>, never>>(object: I): QueryUnbondingDelegationRequest;
};
export declare const QueryUnbondingDelegationResponse: {
    typeUrl: string;
    encode(message: QueryUnbondingDelegationResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryUnbondingDelegationResponse;
    fromJSON(object: any): QueryUnbondingDelegationResponse;
    toJSON(message: QueryUnbondingDelegationResponse): unknown;
    fromPartial<I extends {
        unbond?: {
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
        } | undefined;
    } & {
        unbond?: ({
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
                } & Record<Exclude<keyof I["unbond"]["entries"][number]["completionTime"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                initialBalance?: string | undefined;
                balance?: string | undefined;
                unbondingId?: bigint | undefined;
                unbondingOnHoldRefCount?: bigint | undefined;
            } & Record<Exclude<keyof I["unbond"]["entries"][number], keyof import("./staking").UnbondingDelegationEntry>, never>)[] & Record<Exclude<keyof I["unbond"]["entries"], keyof {
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
        } & Record<Exclude<keyof I["unbond"], keyof UnbondingDelegation>, never>) | undefined;
    } & Record<Exclude<keyof I, "unbond">, never>>(object: I): QueryUnbondingDelegationResponse;
};
export declare const QueryDelegatorDelegationsRequest: {
    typeUrl: string;
    encode(message: QueryDelegatorDelegationsRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryDelegatorDelegationsRequest;
    fromJSON(object: any): QueryDelegatorDelegationsRequest;
    toJSON(message: QueryDelegatorDelegationsRequest): unknown;
    fromPartial<I extends {
        delegatorAddr?: string | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        delegatorAddr?: string | undefined;
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
    } & Record<Exclude<keyof I, keyof QueryDelegatorDelegationsRequest>, never>>(object: I): QueryDelegatorDelegationsRequest;
};
export declare const QueryDelegatorDelegationsResponse: {
    typeUrl: string;
    encode(message: QueryDelegatorDelegationsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryDelegatorDelegationsResponse;
    fromJSON(object: any): QueryDelegatorDelegationsResponse;
    toJSON(message: QueryDelegatorDelegationsResponse): unknown;
    fromPartial<I extends {
        delegationResponses?: {
            delegation?: {
                delegatorAddress?: string | undefined;
                validatorAddress?: string | undefined;
                shares?: string | undefined;
            } | undefined;
            balance?: {
                denom?: string | undefined;
                amount?: string | undefined;
            } | undefined;
        }[] | undefined;
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
        delegationResponses?: ({
            delegation?: {
                delegatorAddress?: string | undefined;
                validatorAddress?: string | undefined;
                shares?: string | undefined;
            } | undefined;
            balance?: {
                denom?: string | undefined;
                amount?: string | undefined;
            } | undefined;
        }[] & ({
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
            } & Record<Exclude<keyof I["delegationResponses"][number]["delegation"], keyof import("./staking").Delegation>, never>) | undefined;
            balance?: ({
                denom?: string | undefined;
                amount?: string | undefined;
            } & {
                denom?: string | undefined;
                amount?: string | undefined;
            } & Record<Exclude<keyof I["delegationResponses"][number]["balance"], keyof import("../../base/v1beta1/coin").Coin>, never>) | undefined;
        } & Record<Exclude<keyof I["delegationResponses"][number], keyof DelegationResponse>, never>)[] & Record<Exclude<keyof I["delegationResponses"], keyof {
            delegation?: {
                delegatorAddress?: string | undefined;
                validatorAddress?: string | undefined;
                shares?: string | undefined;
            } | undefined;
            balance?: {
                denom?: string | undefined;
                amount?: string | undefined;
            } | undefined;
        }[]>, never>) | undefined;
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryDelegatorDelegationsResponse>, never>>(object: I): QueryDelegatorDelegationsResponse;
};
export declare const QueryDelegatorUnbondingDelegationsRequest: {
    typeUrl: string;
    encode(message: QueryDelegatorUnbondingDelegationsRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryDelegatorUnbondingDelegationsRequest;
    fromJSON(object: any): QueryDelegatorUnbondingDelegationsRequest;
    toJSON(message: QueryDelegatorUnbondingDelegationsRequest): unknown;
    fromPartial<I extends {
        delegatorAddr?: string | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        delegatorAddr?: string | undefined;
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
    } & Record<Exclude<keyof I, keyof QueryDelegatorUnbondingDelegationsRequest>, never>>(object: I): QueryDelegatorUnbondingDelegationsRequest;
};
export declare const QueryDelegatorUnbondingDelegationsResponse: {
    typeUrl: string;
    encode(message: QueryDelegatorUnbondingDelegationsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryDelegatorUnbondingDelegationsResponse;
    fromJSON(object: any): QueryDelegatorUnbondingDelegationsResponse;
    toJSON(message: QueryDelegatorUnbondingDelegationsResponse): unknown;
    fromPartial<I extends {
        unbondingResponses?: {
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
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
        unbondingResponses?: ({
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
                } & Record<Exclude<keyof I["unbondingResponses"][number]["entries"][number]["completionTime"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                initialBalance?: string | undefined;
                balance?: string | undefined;
                unbondingId?: bigint | undefined;
                unbondingOnHoldRefCount?: bigint | undefined;
            } & Record<Exclude<keyof I["unbondingResponses"][number]["entries"][number], keyof import("./staking").UnbondingDelegationEntry>, never>)[] & Record<Exclude<keyof I["unbondingResponses"][number]["entries"], keyof {
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
        } & Record<Exclude<keyof I["unbondingResponses"][number], keyof UnbondingDelegation>, never>)[] & Record<Exclude<keyof I["unbondingResponses"], keyof {
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
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryDelegatorUnbondingDelegationsResponse>, never>>(object: I): QueryDelegatorUnbondingDelegationsResponse;
};
export declare const QueryRedelegationsRequest: {
    typeUrl: string;
    encode(message: QueryRedelegationsRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryRedelegationsRequest;
    fromJSON(object: any): QueryRedelegationsRequest;
    toJSON(message: QueryRedelegationsRequest): unknown;
    fromPartial<I extends {
        delegatorAddr?: string | undefined;
        srcValidatorAddr?: string | undefined;
        dstValidatorAddr?: string | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        delegatorAddr?: string | undefined;
        srcValidatorAddr?: string | undefined;
        dstValidatorAddr?: string | undefined;
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
    } & Record<Exclude<keyof I, keyof QueryRedelegationsRequest>, never>>(object: I): QueryRedelegationsRequest;
};
export declare const QueryRedelegationsResponse: {
    typeUrl: string;
    encode(message: QueryRedelegationsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryRedelegationsResponse;
    fromJSON(object: any): QueryRedelegationsResponse;
    toJSON(message: QueryRedelegationsResponse): unknown;
    fromPartial<I extends {
        redelegationResponses?: {
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
        }[] | undefined;
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
        redelegationResponses?: ({
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
        }[] & ({
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
                    } & Record<Exclude<keyof I["redelegationResponses"][number]["redelegation"]["entries"][number]["completionTime"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                    initialBalance?: string | undefined;
                    sharesDst?: string | undefined;
                    unbondingId?: bigint | undefined;
                    unbondingOnHoldRefCount?: bigint | undefined;
                } & Record<Exclude<keyof I["redelegationResponses"][number]["redelegation"]["entries"][number], keyof import("./staking").RedelegationEntry>, never>)[] & Record<Exclude<keyof I["redelegationResponses"][number]["redelegation"]["entries"], keyof {
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
            } & Record<Exclude<keyof I["redelegationResponses"][number]["redelegation"], keyof import("./staking").Redelegation>, never>) | undefined;
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
                    } & Record<Exclude<keyof I["redelegationResponses"][number]["entries"][number]["redelegationEntry"]["completionTime"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                    initialBalance?: string | undefined;
                    sharesDst?: string | undefined;
                    unbondingId?: bigint | undefined;
                    unbondingOnHoldRefCount?: bigint | undefined;
                } & Record<Exclude<keyof I["redelegationResponses"][number]["entries"][number]["redelegationEntry"], keyof import("./staking").RedelegationEntry>, never>) | undefined;
                balance?: string | undefined;
            } & Record<Exclude<keyof I["redelegationResponses"][number]["entries"][number], keyof import("./staking").RedelegationEntryResponse>, never>)[] & Record<Exclude<keyof I["redelegationResponses"][number]["entries"], keyof {
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
        } & Record<Exclude<keyof I["redelegationResponses"][number], keyof RedelegationResponse>, never>)[] & Record<Exclude<keyof I["redelegationResponses"], keyof {
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
        }[]>, never>) | undefined;
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryRedelegationsResponse>, never>>(object: I): QueryRedelegationsResponse;
};
export declare const QueryDelegatorValidatorsRequest: {
    typeUrl: string;
    encode(message: QueryDelegatorValidatorsRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryDelegatorValidatorsRequest;
    fromJSON(object: any): QueryDelegatorValidatorsRequest;
    toJSON(message: QueryDelegatorValidatorsRequest): unknown;
    fromPartial<I extends {
        delegatorAddr?: string | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        delegatorAddr?: string | undefined;
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
    } & Record<Exclude<keyof I, keyof QueryDelegatorValidatorsRequest>, never>>(object: I): QueryDelegatorValidatorsRequest;
};
export declare const QueryDelegatorValidatorsResponse: {
    typeUrl: string;
    encode(message: QueryDelegatorValidatorsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryDelegatorValidatorsResponse;
    fromJSON(object: any): QueryDelegatorValidatorsResponse;
    toJSON(message: QueryDelegatorValidatorsResponse): unknown;
    fromPartial<I extends {
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
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
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
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryDelegatorValidatorsResponse>, never>>(object: I): QueryDelegatorValidatorsResponse;
};
export declare const QueryDelegatorValidatorRequest: {
    typeUrl: string;
    encode(message: QueryDelegatorValidatorRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryDelegatorValidatorRequest;
    fromJSON(object: any): QueryDelegatorValidatorRequest;
    toJSON(message: QueryDelegatorValidatorRequest): unknown;
    fromPartial<I extends {
        delegatorAddr?: string | undefined;
        validatorAddr?: string | undefined;
    } & {
        delegatorAddr?: string | undefined;
        validatorAddr?: string | undefined;
    } & Record<Exclude<keyof I, keyof QueryDelegatorValidatorRequest>, never>>(object: I): QueryDelegatorValidatorRequest;
};
export declare const QueryDelegatorValidatorResponse: {
    typeUrl: string;
    encode(message: QueryDelegatorValidatorResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryDelegatorValidatorResponse;
    fromJSON(object: any): QueryDelegatorValidatorResponse;
    toJSON(message: QueryDelegatorValidatorResponse): unknown;
    fromPartial<I extends {
        validator?: {
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
        } | undefined;
    } & {
        validator?: ({
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
            } & Record<Exclude<keyof I["validator"]["consensusPubkey"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
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
            } & Record<Exclude<keyof I["validator"]["description"], keyof import("./staking").Description>, never>) | undefined;
            unbondingHeight?: bigint | undefined;
            unbondingTime?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["validator"]["unbondingTime"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
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
                } & Record<Exclude<keyof I["validator"]["commission"]["commissionRates"], keyof import("./staking").CommissionRates>, never>) | undefined;
                updateTime?: ({
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & Record<Exclude<keyof I["validator"]["commission"]["updateTime"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
            } & Record<Exclude<keyof I["validator"]["commission"], keyof import("./staking").Commission>, never>) | undefined;
            minSelfDelegation?: string | undefined;
            unbondingOnHoldRefCount?: bigint | undefined;
            unbondingIds?: (bigint[] & bigint[] & Record<Exclude<keyof I["validator"]["unbondingIds"], keyof bigint[]>, never>) | undefined;
        } & Record<Exclude<keyof I["validator"], keyof Validator>, never>) | undefined;
    } & Record<Exclude<keyof I, "validator">, never>>(object: I): QueryDelegatorValidatorResponse;
};
export declare const QueryHistoricalInfoRequest: {
    typeUrl: string;
    encode(message: QueryHistoricalInfoRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryHistoricalInfoRequest;
    fromJSON(object: any): QueryHistoricalInfoRequest;
    toJSON(message: QueryHistoricalInfoRequest): unknown;
    fromPartial<I extends {
        height?: bigint | undefined;
    } & {
        height?: bigint | undefined;
    } & Record<Exclude<keyof I, "height">, never>>(object: I): QueryHistoricalInfoRequest;
};
export declare const QueryHistoricalInfoResponse: {
    typeUrl: string;
    encode(message: QueryHistoricalInfoResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryHistoricalInfoResponse;
    fromJSON(object: any): QueryHistoricalInfoResponse;
    toJSON(message: QueryHistoricalInfoResponse): unknown;
    fromPartial<I extends {
        hist?: {
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
        } | undefined;
    } & {
        hist?: ({
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
                } & Record<Exclude<keyof I["hist"]["header"]["version"], keyof import("../../../tendermint/version/types").Consensus>, never>) | undefined;
                chainId?: string | undefined;
                height?: bigint | undefined;
                time?: ({
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & Record<Exclude<keyof I["hist"]["header"]["time"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
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
                    } & Record<Exclude<keyof I["hist"]["header"]["lastBlockId"]["partSetHeader"], keyof import("../../../tendermint/types/types").PartSetHeader>, never>) | undefined;
                } & Record<Exclude<keyof I["hist"]["header"]["lastBlockId"], keyof import("../../../tendermint/types/types").BlockID>, never>) | undefined;
                lastCommitHash?: Uint8Array | undefined;
                dataHash?: Uint8Array | undefined;
                validatorsHash?: Uint8Array | undefined;
                nextValidatorsHash?: Uint8Array | undefined;
                consensusHash?: Uint8Array | undefined;
                appHash?: Uint8Array | undefined;
                lastResultsHash?: Uint8Array | undefined;
                evidenceHash?: Uint8Array | undefined;
                proposerAddress?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["hist"]["header"], keyof import("../../../tendermint/types/types").Header>, never>) | undefined;
            valset?: ({
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
                } & Record<Exclude<keyof I["hist"]["valset"][number]["consensusPubkey"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
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
                } & Record<Exclude<keyof I["hist"]["valset"][number]["description"], keyof import("./staking").Description>, never>) | undefined;
                unbondingHeight?: bigint | undefined;
                unbondingTime?: ({
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & Record<Exclude<keyof I["hist"]["valset"][number]["unbondingTime"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
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
                    } & Record<Exclude<keyof I["hist"]["valset"][number]["commission"]["commissionRates"], keyof import("./staking").CommissionRates>, never>) | undefined;
                    updateTime?: ({
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } & {
                        seconds?: bigint | undefined;
                        nanos?: number | undefined;
                    } & Record<Exclude<keyof I["hist"]["valset"][number]["commission"]["updateTime"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                } & Record<Exclude<keyof I["hist"]["valset"][number]["commission"], keyof import("./staking").Commission>, never>) | undefined;
                minSelfDelegation?: string | undefined;
                unbondingOnHoldRefCount?: bigint | undefined;
                unbondingIds?: (bigint[] & bigint[] & Record<Exclude<keyof I["hist"]["valset"][number]["unbondingIds"], keyof bigint[]>, never>) | undefined;
            } & Record<Exclude<keyof I["hist"]["valset"][number], keyof Validator>, never>)[] & Record<Exclude<keyof I["hist"]["valset"], keyof {
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
        } & Record<Exclude<keyof I["hist"], keyof HistoricalInfo>, never>) | undefined;
    } & Record<Exclude<keyof I, "hist">, never>>(object: I): QueryHistoricalInfoResponse;
};
export declare const QueryPoolRequest: {
    typeUrl: string;
    encode(_: QueryPoolRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryPoolRequest;
    fromJSON(_: any): QueryPoolRequest;
    toJSON(_: QueryPoolRequest): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): QueryPoolRequest;
};
export declare const QueryPoolResponse: {
    typeUrl: string;
    encode(message: QueryPoolResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryPoolResponse;
    fromJSON(object: any): QueryPoolResponse;
    toJSON(message: QueryPoolResponse): unknown;
    fromPartial<I extends {
        pool?: {
            notBondedTokens?: string | undefined;
            bondedTokens?: string | undefined;
        } | undefined;
    } & {
        pool?: ({
            notBondedTokens?: string | undefined;
            bondedTokens?: string | undefined;
        } & {
            notBondedTokens?: string | undefined;
            bondedTokens?: string | undefined;
        } & Record<Exclude<keyof I["pool"], keyof Pool>, never>) | undefined;
    } & Record<Exclude<keyof I, "pool">, never>>(object: I): QueryPoolResponse;
};
export declare const QueryParamsRequest: {
    typeUrl: string;
    encode(_: QueryParamsRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryParamsRequest;
    fromJSON(_: any): QueryParamsRequest;
    toJSON(_: QueryParamsRequest): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): QueryParamsRequest;
};
export declare const QueryParamsResponse: {
    typeUrl: string;
    encode(message: QueryParamsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryParamsResponse;
    fromJSON(object: any): QueryParamsResponse;
    toJSON(message: QueryParamsResponse): unknown;
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
    } & Record<Exclude<keyof I, "params">, never>>(object: I): QueryParamsResponse;
};
/** Query defines the gRPC querier service. */
export interface Query {
    /**
     * Validators queries all validators that match the given status.
     *
     * When called from another module, this query might consume a high amount of
     * gas if the pagination field is incorrectly set.
     */
    Validators(request: QueryValidatorsRequest): Promise<QueryValidatorsResponse>;
    /** Validator queries validator info for given validator address. */
    Validator(request: QueryValidatorRequest): Promise<QueryValidatorResponse>;
    /**
     * ValidatorDelegations queries delegate info for given validator.
     *
     * When called from another module, this query might consume a high amount of
     * gas if the pagination field is incorrectly set.
     */
    ValidatorDelegations(request: QueryValidatorDelegationsRequest): Promise<QueryValidatorDelegationsResponse>;
    /**
     * ValidatorUnbondingDelegations queries unbonding delegations of a validator.
     *
     * When called from another module, this query might consume a high amount of
     * gas if the pagination field is incorrectly set.
     */
    ValidatorUnbondingDelegations(request: QueryValidatorUnbondingDelegationsRequest): Promise<QueryValidatorUnbondingDelegationsResponse>;
    /** Delegation queries delegate info for given validator delegator pair. */
    Delegation(request: QueryDelegationRequest): Promise<QueryDelegationResponse>;
    /**
     * UnbondingDelegation queries unbonding info for given validator delegator
     * pair.
     */
    UnbondingDelegation(request: QueryUnbondingDelegationRequest): Promise<QueryUnbondingDelegationResponse>;
    /**
     * DelegatorDelegations queries all delegations of a given delegator address.
     *
     * When called from another module, this query might consume a high amount of
     * gas if the pagination field is incorrectly set.
     */
    DelegatorDelegations(request: QueryDelegatorDelegationsRequest): Promise<QueryDelegatorDelegationsResponse>;
    /**
     * DelegatorUnbondingDelegations queries all unbonding delegations of a given
     * delegator address.
     *
     * When called from another module, this query might consume a high amount of
     * gas if the pagination field is incorrectly set.
     */
    DelegatorUnbondingDelegations(request: QueryDelegatorUnbondingDelegationsRequest): Promise<QueryDelegatorUnbondingDelegationsResponse>;
    /**
     * Redelegations queries redelegations of given address.
     *
     * When called from another module, this query might consume a high amount of
     * gas if the pagination field is incorrectly set.
     */
    Redelegations(request: QueryRedelegationsRequest): Promise<QueryRedelegationsResponse>;
    /**
     * DelegatorValidators queries all validators info for given delegator
     * address.
     *
     * When called from another module, this query might consume a high amount of
     * gas if the pagination field is incorrectly set.
     */
    DelegatorValidators(request: QueryDelegatorValidatorsRequest): Promise<QueryDelegatorValidatorsResponse>;
    /**
     * DelegatorValidator queries validator info for given delegator validator
     * pair.
     */
    DelegatorValidator(request: QueryDelegatorValidatorRequest): Promise<QueryDelegatorValidatorResponse>;
    /** HistoricalInfo queries the historical info for given height. */
    HistoricalInfo(request: QueryHistoricalInfoRequest): Promise<QueryHistoricalInfoResponse>;
    /** Pool queries the pool info. */
    Pool(request?: QueryPoolRequest): Promise<QueryPoolResponse>;
    /** Parameters queries the staking parameters. */
    Params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    Validators(request: QueryValidatorsRequest): Promise<QueryValidatorsResponse>;
    Validator(request: QueryValidatorRequest): Promise<QueryValidatorResponse>;
    ValidatorDelegations(request: QueryValidatorDelegationsRequest): Promise<QueryValidatorDelegationsResponse>;
    ValidatorUnbondingDelegations(request: QueryValidatorUnbondingDelegationsRequest): Promise<QueryValidatorUnbondingDelegationsResponse>;
    Delegation(request: QueryDelegationRequest): Promise<QueryDelegationResponse>;
    UnbondingDelegation(request: QueryUnbondingDelegationRequest): Promise<QueryUnbondingDelegationResponse>;
    DelegatorDelegations(request: QueryDelegatorDelegationsRequest): Promise<QueryDelegatorDelegationsResponse>;
    DelegatorUnbondingDelegations(request: QueryDelegatorUnbondingDelegationsRequest): Promise<QueryDelegatorUnbondingDelegationsResponse>;
    Redelegations(request: QueryRedelegationsRequest): Promise<QueryRedelegationsResponse>;
    DelegatorValidators(request: QueryDelegatorValidatorsRequest): Promise<QueryDelegatorValidatorsResponse>;
    DelegatorValidator(request: QueryDelegatorValidatorRequest): Promise<QueryDelegatorValidatorResponse>;
    HistoricalInfo(request: QueryHistoricalInfoRequest): Promise<QueryHistoricalInfoResponse>;
    Pool(request?: QueryPoolRequest): Promise<QueryPoolResponse>;
    Params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
}
