import { PageRequest, PageResponse } from "../../base/query/v1beta1/pagination";
import { Coin } from "../../base/v1beta1/coin";
import { Params, Metadata, SendEnabled } from "./bank";
import { BinaryReader, BinaryWriter } from "../../../binary";
import { Rpc } from "../../../helpers";
export declare const protobufPackage = "cosmos.bank.v1beta1";
/** QueryBalanceRequest is the request type for the Query/Balance RPC method. */
export interface QueryBalanceRequest {
    /** address is the address to query balances for. */
    address: string;
    /** denom is the coin denom to query balances for. */
    denom: string;
}
/** QueryBalanceResponse is the response type for the Query/Balance RPC method. */
export interface QueryBalanceResponse {
    /** balance is the balance of the coin. */
    balance?: Coin;
}
/** QueryBalanceRequest is the request type for the Query/AllBalances RPC method. */
export interface QueryAllBalancesRequest {
    /** address is the address to query balances for. */
    address: string;
    /** pagination defines an optional pagination for the request. */
    pagination?: PageRequest;
}
/**
 * QueryAllBalancesResponse is the response type for the Query/AllBalances RPC
 * method.
 */
export interface QueryAllBalancesResponse {
    /** balances is the balances of all the coins. */
    balances: Coin[];
    /** pagination defines the pagination in the response. */
    pagination?: PageResponse;
}
/**
 * QuerySpendableBalancesRequest defines the gRPC request structure for querying
 * an account's spendable balances.
 *
 * Since: cosmos-sdk 0.46
 */
export interface QuerySpendableBalancesRequest {
    /** address is the address to query spendable balances for. */
    address: string;
    /** pagination defines an optional pagination for the request. */
    pagination?: PageRequest;
}
/**
 * QuerySpendableBalancesResponse defines the gRPC response structure for querying
 * an account's spendable balances.
 *
 * Since: cosmos-sdk 0.46
 */
export interface QuerySpendableBalancesResponse {
    /** balances is the spendable balances of all the coins. */
    balances: Coin[];
    /** pagination defines the pagination in the response. */
    pagination?: PageResponse;
}
/**
 * QuerySpendableBalanceByDenomRequest defines the gRPC request structure for
 * querying an account's spendable balance for a specific denom.
 *
 * Since: cosmos-sdk 0.47
 */
export interface QuerySpendableBalanceByDenomRequest {
    /** address is the address to query balances for. */
    address: string;
    /** denom is the coin denom to query balances for. */
    denom: string;
}
/**
 * QuerySpendableBalanceByDenomResponse defines the gRPC response structure for
 * querying an account's spendable balance for a specific denom.
 *
 * Since: cosmos-sdk 0.47
 */
export interface QuerySpendableBalanceByDenomResponse {
    /** balance is the balance of the coin. */
    balance?: Coin;
}
/**
 * QueryTotalSupplyRequest is the request type for the Query/TotalSupply RPC
 * method.
 */
export interface QueryTotalSupplyRequest {
    /**
     * pagination defines an optional pagination for the request.
     *
     * Since: cosmos-sdk 0.43
     */
    pagination?: PageRequest;
}
/**
 * QueryTotalSupplyResponse is the response type for the Query/TotalSupply RPC
 * method
 */
export interface QueryTotalSupplyResponse {
    /** supply is the supply of the coins */
    supply: Coin[];
    /**
     * pagination defines the pagination in the response.
     *
     * Since: cosmos-sdk 0.43
     */
    pagination?: PageResponse;
}
/** QuerySupplyOfRequest is the request type for the Query/SupplyOf RPC method. */
export interface QuerySupplyOfRequest {
    /** denom is the coin denom to query balances for. */
    denom: string;
}
/** QuerySupplyOfResponse is the response type for the Query/SupplyOf RPC method. */
export interface QuerySupplyOfResponse {
    /** amount is the supply of the coin. */
    amount: Coin;
}
/** QueryParamsRequest defines the request type for querying x/bank parameters. */
export interface QueryParamsRequest {
}
/** QueryParamsResponse defines the response type for querying x/bank parameters. */
export interface QueryParamsResponse {
    params: Params;
}
/** QueryDenomsMetadataRequest is the request type for the Query/DenomsMetadata RPC method. */
export interface QueryDenomsMetadataRequest {
    /** pagination defines an optional pagination for the request. */
    pagination?: PageRequest;
}
/**
 * QueryDenomsMetadataResponse is the response type for the Query/DenomsMetadata RPC
 * method.
 */
export interface QueryDenomsMetadataResponse {
    /** metadata provides the client information for all the registered tokens. */
    metadatas: Metadata[];
    /** pagination defines the pagination in the response. */
    pagination?: PageResponse;
}
/** QueryDenomMetadataRequest is the request type for the Query/DenomMetadata RPC method. */
export interface QueryDenomMetadataRequest {
    /** denom is the coin denom to query the metadata for. */
    denom: string;
}
/**
 * QueryDenomMetadataResponse is the response type for the Query/DenomMetadata RPC
 * method.
 */
export interface QueryDenomMetadataResponse {
    /** metadata describes and provides all the client information for the requested token. */
    metadata: Metadata;
}
/**
 * QueryDenomOwnersRequest defines the request type for the DenomOwners RPC query,
 * which queries for a paginated set of all account holders of a particular
 * denomination.
 */
export interface QueryDenomOwnersRequest {
    /** denom defines the coin denomination to query all account holders for. */
    denom: string;
    /** pagination defines an optional pagination for the request. */
    pagination?: PageRequest;
}
/**
 * DenomOwner defines structure representing an account that owns or holds a
 * particular denominated token. It contains the account address and account
 * balance of the denominated token.
 *
 * Since: cosmos-sdk 0.46
 */
export interface DenomOwner {
    /** address defines the address that owns a particular denomination. */
    address: string;
    /** balance is the balance of the denominated coin for an account. */
    balance: Coin;
}
/**
 * QueryDenomOwnersResponse defines the RPC response of a DenomOwners RPC query.
 *
 * Since: cosmos-sdk 0.46
 */
export interface QueryDenomOwnersResponse {
    denomOwners: DenomOwner[];
    /** pagination defines the pagination in the response. */
    pagination?: PageResponse;
}
/**
 * QuerySendEnabledRequest defines the RPC request for looking up SendEnabled entries.
 *
 * Since: cosmos-sdk 0.47
 */
export interface QuerySendEnabledRequest {
    /** denoms is the specific denoms you want look up. Leave empty to get all entries. */
    denoms: string[];
    /**
     * pagination defines an optional pagination for the request. This field is
     * only read if the denoms field is empty.
     */
    pagination?: PageRequest;
}
/**
 * QuerySendEnabledResponse defines the RPC response of a SendEnable query.
 *
 * Since: cosmos-sdk 0.47
 */
export interface QuerySendEnabledResponse {
    sendEnabled: SendEnabled[];
    /**
     * pagination defines the pagination in the response. This field is only
     * populated if the denoms field in the request is empty.
     */
    pagination?: PageResponse;
}
export declare const QueryBalanceRequest: {
    typeUrl: string;
    encode(message: QueryBalanceRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryBalanceRequest;
    fromJSON(object: any): QueryBalanceRequest;
    toJSON(message: QueryBalanceRequest): unknown;
    fromPartial<I extends {
        address?: string | undefined;
        denom?: string | undefined;
    } & {
        address?: string | undefined;
        denom?: string | undefined;
    } & Record<Exclude<keyof I, keyof QueryBalanceRequest>, never>>(object: I): QueryBalanceRequest;
};
export declare const QueryBalanceResponse: {
    typeUrl: string;
    encode(message: QueryBalanceResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryBalanceResponse;
    fromJSON(object: any): QueryBalanceResponse;
    toJSON(message: QueryBalanceResponse): unknown;
    fromPartial<I extends {
        balance?: {
            denom?: string | undefined;
            amount?: string | undefined;
        } | undefined;
    } & {
        balance?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["balance"], keyof Coin>, never>) | undefined;
    } & Record<Exclude<keyof I, "balance">, never>>(object: I): QueryBalanceResponse;
};
export declare const QueryAllBalancesRequest: {
    typeUrl: string;
    encode(message: QueryAllBalancesRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryAllBalancesRequest;
    fromJSON(object: any): QueryAllBalancesRequest;
    toJSON(message: QueryAllBalancesRequest): unknown;
    fromPartial<I extends {
        address?: string | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        address?: string | undefined;
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
    } & Record<Exclude<keyof I, keyof QueryAllBalancesRequest>, never>>(object: I): QueryAllBalancesRequest;
};
export declare const QueryAllBalancesResponse: {
    typeUrl: string;
    encode(message: QueryAllBalancesResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryAllBalancesResponse;
    fromJSON(object: any): QueryAllBalancesResponse;
    toJSON(message: QueryAllBalancesResponse): unknown;
    fromPartial<I extends {
        balances?: {
            denom?: string | undefined;
            amount?: string | undefined;
        }[] | undefined;
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
        balances?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        }[] & ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["balances"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["balances"], keyof {
            denom?: string | undefined;
            amount?: string | undefined;
        }[]>, never>) | undefined;
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryAllBalancesResponse>, never>>(object: I): QueryAllBalancesResponse;
};
export declare const QuerySpendableBalancesRequest: {
    typeUrl: string;
    encode(message: QuerySpendableBalancesRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QuerySpendableBalancesRequest;
    fromJSON(object: any): QuerySpendableBalancesRequest;
    toJSON(message: QuerySpendableBalancesRequest): unknown;
    fromPartial<I extends {
        address?: string | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        address?: string | undefined;
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
    } & Record<Exclude<keyof I, keyof QuerySpendableBalancesRequest>, never>>(object: I): QuerySpendableBalancesRequest;
};
export declare const QuerySpendableBalancesResponse: {
    typeUrl: string;
    encode(message: QuerySpendableBalancesResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QuerySpendableBalancesResponse;
    fromJSON(object: any): QuerySpendableBalancesResponse;
    toJSON(message: QuerySpendableBalancesResponse): unknown;
    fromPartial<I extends {
        balances?: {
            denom?: string | undefined;
            amount?: string | undefined;
        }[] | undefined;
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
        balances?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        }[] & ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["balances"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["balances"], keyof {
            denom?: string | undefined;
            amount?: string | undefined;
        }[]>, never>) | undefined;
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QuerySpendableBalancesResponse>, never>>(object: I): QuerySpendableBalancesResponse;
};
export declare const QuerySpendableBalanceByDenomRequest: {
    typeUrl: string;
    encode(message: QuerySpendableBalanceByDenomRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QuerySpendableBalanceByDenomRequest;
    fromJSON(object: any): QuerySpendableBalanceByDenomRequest;
    toJSON(message: QuerySpendableBalanceByDenomRequest): unknown;
    fromPartial<I extends {
        address?: string | undefined;
        denom?: string | undefined;
    } & {
        address?: string | undefined;
        denom?: string | undefined;
    } & Record<Exclude<keyof I, keyof QuerySpendableBalanceByDenomRequest>, never>>(object: I): QuerySpendableBalanceByDenomRequest;
};
export declare const QuerySpendableBalanceByDenomResponse: {
    typeUrl: string;
    encode(message: QuerySpendableBalanceByDenomResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QuerySpendableBalanceByDenomResponse;
    fromJSON(object: any): QuerySpendableBalanceByDenomResponse;
    toJSON(message: QuerySpendableBalanceByDenomResponse): unknown;
    fromPartial<I extends {
        balance?: {
            denom?: string | undefined;
            amount?: string | undefined;
        } | undefined;
    } & {
        balance?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["balance"], keyof Coin>, never>) | undefined;
    } & Record<Exclude<keyof I, "balance">, never>>(object: I): QuerySpendableBalanceByDenomResponse;
};
export declare const QueryTotalSupplyRequest: {
    typeUrl: string;
    encode(message: QueryTotalSupplyRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryTotalSupplyRequest;
    fromJSON(object: any): QueryTotalSupplyRequest;
    toJSON(message: QueryTotalSupplyRequest): unknown;
    fromPartial<I extends {
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
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
    } & Record<Exclude<keyof I, "pagination">, never>>(object: I): QueryTotalSupplyRequest;
};
export declare const QueryTotalSupplyResponse: {
    typeUrl: string;
    encode(message: QueryTotalSupplyResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryTotalSupplyResponse;
    fromJSON(object: any): QueryTotalSupplyResponse;
    toJSON(message: QueryTotalSupplyResponse): unknown;
    fromPartial<I extends {
        supply?: {
            denom?: string | undefined;
            amount?: string | undefined;
        }[] | undefined;
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
        supply?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        }[] & ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["supply"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["supply"], keyof {
            denom?: string | undefined;
            amount?: string | undefined;
        }[]>, never>) | undefined;
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryTotalSupplyResponse>, never>>(object: I): QueryTotalSupplyResponse;
};
export declare const QuerySupplyOfRequest: {
    typeUrl: string;
    encode(message: QuerySupplyOfRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QuerySupplyOfRequest;
    fromJSON(object: any): QuerySupplyOfRequest;
    toJSON(message: QuerySupplyOfRequest): unknown;
    fromPartial<I extends {
        denom?: string | undefined;
    } & {
        denom?: string | undefined;
    } & Record<Exclude<keyof I, "denom">, never>>(object: I): QuerySupplyOfRequest;
};
export declare const QuerySupplyOfResponse: {
    typeUrl: string;
    encode(message: QuerySupplyOfResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QuerySupplyOfResponse;
    fromJSON(object: any): QuerySupplyOfResponse;
    toJSON(message: QuerySupplyOfResponse): unknown;
    fromPartial<I extends {
        amount?: {
            denom?: string | undefined;
            amount?: string | undefined;
        } | undefined;
    } & {
        amount?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["amount"], keyof Coin>, never>) | undefined;
    } & Record<Exclude<keyof I, "amount">, never>>(object: I): QuerySupplyOfResponse;
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
            sendEnabled?: {
                denom?: string | undefined;
                enabled?: boolean | undefined;
            }[] | undefined;
            defaultSendEnabled?: boolean | undefined;
        } | undefined;
    } & {
        params?: ({
            sendEnabled?: {
                denom?: string | undefined;
                enabled?: boolean | undefined;
            }[] | undefined;
            defaultSendEnabled?: boolean | undefined;
        } & {
            sendEnabled?: ({
                denom?: string | undefined;
                enabled?: boolean | undefined;
            }[] & ({
                denom?: string | undefined;
                enabled?: boolean | undefined;
            } & {
                denom?: string | undefined;
                enabled?: boolean | undefined;
            } & Record<Exclude<keyof I["params"]["sendEnabled"][number], keyof SendEnabled>, never>)[] & Record<Exclude<keyof I["params"]["sendEnabled"], keyof {
                denom?: string | undefined;
                enabled?: boolean | undefined;
            }[]>, never>) | undefined;
            defaultSendEnabled?: boolean | undefined;
        } & Record<Exclude<keyof I["params"], keyof Params>, never>) | undefined;
    } & Record<Exclude<keyof I, "params">, never>>(object: I): QueryParamsResponse;
};
export declare const QueryDenomsMetadataRequest: {
    typeUrl: string;
    encode(message: QueryDenomsMetadataRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryDenomsMetadataRequest;
    fromJSON(object: any): QueryDenomsMetadataRequest;
    toJSON(message: QueryDenomsMetadataRequest): unknown;
    fromPartial<I extends {
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
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
    } & Record<Exclude<keyof I, "pagination">, never>>(object: I): QueryDenomsMetadataRequest;
};
export declare const QueryDenomsMetadataResponse: {
    typeUrl: string;
    encode(message: QueryDenomsMetadataResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryDenomsMetadataResponse;
    fromJSON(object: any): QueryDenomsMetadataResponse;
    toJSON(message: QueryDenomsMetadataResponse): unknown;
    fromPartial<I extends {
        metadatas?: {
            description?: string | undefined;
            denomUnits?: {
                denom?: string | undefined;
                exponent?: number | undefined;
                aliases?: string[] | undefined;
            }[] | undefined;
            base?: string | undefined;
            display?: string | undefined;
            name?: string | undefined;
            symbol?: string | undefined;
            uri?: string | undefined;
            uriHash?: string | undefined;
        }[] | undefined;
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
        metadatas?: ({
            description?: string | undefined;
            denomUnits?: {
                denom?: string | undefined;
                exponent?: number | undefined;
                aliases?: string[] | undefined;
            }[] | undefined;
            base?: string | undefined;
            display?: string | undefined;
            name?: string | undefined;
            symbol?: string | undefined;
            uri?: string | undefined;
            uriHash?: string | undefined;
        }[] & ({
            description?: string | undefined;
            denomUnits?: {
                denom?: string | undefined;
                exponent?: number | undefined;
                aliases?: string[] | undefined;
            }[] | undefined;
            base?: string | undefined;
            display?: string | undefined;
            name?: string | undefined;
            symbol?: string | undefined;
            uri?: string | undefined;
            uriHash?: string | undefined;
        } & {
            description?: string | undefined;
            denomUnits?: ({
                denom?: string | undefined;
                exponent?: number | undefined;
                aliases?: string[] | undefined;
            }[] & ({
                denom?: string | undefined;
                exponent?: number | undefined;
                aliases?: string[] | undefined;
            } & {
                denom?: string | undefined;
                exponent?: number | undefined;
                aliases?: (string[] & string[] & Record<Exclude<keyof I["metadatas"][number]["denomUnits"][number]["aliases"], keyof string[]>, never>) | undefined;
            } & Record<Exclude<keyof I["metadatas"][number]["denomUnits"][number], keyof import("./bank").DenomUnit>, never>)[] & Record<Exclude<keyof I["metadatas"][number]["denomUnits"], keyof {
                denom?: string | undefined;
                exponent?: number | undefined;
                aliases?: string[] | undefined;
            }[]>, never>) | undefined;
            base?: string | undefined;
            display?: string | undefined;
            name?: string | undefined;
            symbol?: string | undefined;
            uri?: string | undefined;
            uriHash?: string | undefined;
        } & Record<Exclude<keyof I["metadatas"][number], keyof Metadata>, never>)[] & Record<Exclude<keyof I["metadatas"], keyof {
            description?: string | undefined;
            denomUnits?: {
                denom?: string | undefined;
                exponent?: number | undefined;
                aliases?: string[] | undefined;
            }[] | undefined;
            base?: string | undefined;
            display?: string | undefined;
            name?: string | undefined;
            symbol?: string | undefined;
            uri?: string | undefined;
            uriHash?: string | undefined;
        }[]>, never>) | undefined;
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryDenomsMetadataResponse>, never>>(object: I): QueryDenomsMetadataResponse;
};
export declare const QueryDenomMetadataRequest: {
    typeUrl: string;
    encode(message: QueryDenomMetadataRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryDenomMetadataRequest;
    fromJSON(object: any): QueryDenomMetadataRequest;
    toJSON(message: QueryDenomMetadataRequest): unknown;
    fromPartial<I extends {
        denom?: string | undefined;
    } & {
        denom?: string | undefined;
    } & Record<Exclude<keyof I, "denom">, never>>(object: I): QueryDenomMetadataRequest;
};
export declare const QueryDenomMetadataResponse: {
    typeUrl: string;
    encode(message: QueryDenomMetadataResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryDenomMetadataResponse;
    fromJSON(object: any): QueryDenomMetadataResponse;
    toJSON(message: QueryDenomMetadataResponse): unknown;
    fromPartial<I extends {
        metadata?: {
            description?: string | undefined;
            denomUnits?: {
                denom?: string | undefined;
                exponent?: number | undefined;
                aliases?: string[] | undefined;
            }[] | undefined;
            base?: string | undefined;
            display?: string | undefined;
            name?: string | undefined;
            symbol?: string | undefined;
            uri?: string | undefined;
            uriHash?: string | undefined;
        } | undefined;
    } & {
        metadata?: ({
            description?: string | undefined;
            denomUnits?: {
                denom?: string | undefined;
                exponent?: number | undefined;
                aliases?: string[] | undefined;
            }[] | undefined;
            base?: string | undefined;
            display?: string | undefined;
            name?: string | undefined;
            symbol?: string | undefined;
            uri?: string | undefined;
            uriHash?: string | undefined;
        } & {
            description?: string | undefined;
            denomUnits?: ({
                denom?: string | undefined;
                exponent?: number | undefined;
                aliases?: string[] | undefined;
            }[] & ({
                denom?: string | undefined;
                exponent?: number | undefined;
                aliases?: string[] | undefined;
            } & {
                denom?: string | undefined;
                exponent?: number | undefined;
                aliases?: (string[] & string[] & Record<Exclude<keyof I["metadata"]["denomUnits"][number]["aliases"], keyof string[]>, never>) | undefined;
            } & Record<Exclude<keyof I["metadata"]["denomUnits"][number], keyof import("./bank").DenomUnit>, never>)[] & Record<Exclude<keyof I["metadata"]["denomUnits"], keyof {
                denom?: string | undefined;
                exponent?: number | undefined;
                aliases?: string[] | undefined;
            }[]>, never>) | undefined;
            base?: string | undefined;
            display?: string | undefined;
            name?: string | undefined;
            symbol?: string | undefined;
            uri?: string | undefined;
            uriHash?: string | undefined;
        } & Record<Exclude<keyof I["metadata"], keyof Metadata>, never>) | undefined;
    } & Record<Exclude<keyof I, "metadata">, never>>(object: I): QueryDenomMetadataResponse;
};
export declare const QueryDenomOwnersRequest: {
    typeUrl: string;
    encode(message: QueryDenomOwnersRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryDenomOwnersRequest;
    fromJSON(object: any): QueryDenomOwnersRequest;
    toJSON(message: QueryDenomOwnersRequest): unknown;
    fromPartial<I extends {
        denom?: string | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        denom?: string | undefined;
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
    } & Record<Exclude<keyof I, keyof QueryDenomOwnersRequest>, never>>(object: I): QueryDenomOwnersRequest;
};
export declare const DenomOwner: {
    typeUrl: string;
    encode(message: DenomOwner, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): DenomOwner;
    fromJSON(object: any): DenomOwner;
    toJSON(message: DenomOwner): unknown;
    fromPartial<I extends {
        address?: string | undefined;
        balance?: {
            denom?: string | undefined;
            amount?: string | undefined;
        } | undefined;
    } & {
        address?: string | undefined;
        balance?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["balance"], keyof Coin>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof DenomOwner>, never>>(object: I): DenomOwner;
};
export declare const QueryDenomOwnersResponse: {
    typeUrl: string;
    encode(message: QueryDenomOwnersResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryDenomOwnersResponse;
    fromJSON(object: any): QueryDenomOwnersResponse;
    toJSON(message: QueryDenomOwnersResponse): unknown;
    fromPartial<I extends {
        denomOwners?: {
            address?: string | undefined;
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
        denomOwners?: ({
            address?: string | undefined;
            balance?: {
                denom?: string | undefined;
                amount?: string | undefined;
            } | undefined;
        }[] & ({
            address?: string | undefined;
            balance?: {
                denom?: string | undefined;
                amount?: string | undefined;
            } | undefined;
        } & {
            address?: string | undefined;
            balance?: ({
                denom?: string | undefined;
                amount?: string | undefined;
            } & {
                denom?: string | undefined;
                amount?: string | undefined;
            } & Record<Exclude<keyof I["denomOwners"][number]["balance"], keyof Coin>, never>) | undefined;
        } & Record<Exclude<keyof I["denomOwners"][number], keyof DenomOwner>, never>)[] & Record<Exclude<keyof I["denomOwners"], keyof {
            address?: string | undefined;
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
    } & Record<Exclude<keyof I, keyof QueryDenomOwnersResponse>, never>>(object: I): QueryDenomOwnersResponse;
};
export declare const QuerySendEnabledRequest: {
    typeUrl: string;
    encode(message: QuerySendEnabledRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QuerySendEnabledRequest;
    fromJSON(object: any): QuerySendEnabledRequest;
    toJSON(message: QuerySendEnabledRequest): unknown;
    fromPartial<I extends {
        denoms?: string[] | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        denoms?: (string[] & string[] & Record<Exclude<keyof I["denoms"], keyof string[]>, never>) | undefined;
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
    } & Record<Exclude<keyof I, keyof QuerySendEnabledRequest>, never>>(object: I): QuerySendEnabledRequest;
};
export declare const QuerySendEnabledResponse: {
    typeUrl: string;
    encode(message: QuerySendEnabledResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QuerySendEnabledResponse;
    fromJSON(object: any): QuerySendEnabledResponse;
    toJSON(message: QuerySendEnabledResponse): unknown;
    fromPartial<I extends {
        sendEnabled?: {
            denom?: string | undefined;
            enabled?: boolean | undefined;
        }[] | undefined;
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
        sendEnabled?: ({
            denom?: string | undefined;
            enabled?: boolean | undefined;
        }[] & ({
            denom?: string | undefined;
            enabled?: boolean | undefined;
        } & {
            denom?: string | undefined;
            enabled?: boolean | undefined;
        } & Record<Exclude<keyof I["sendEnabled"][number], keyof SendEnabled>, never>)[] & Record<Exclude<keyof I["sendEnabled"], keyof {
            denom?: string | undefined;
            enabled?: boolean | undefined;
        }[]>, never>) | undefined;
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QuerySendEnabledResponse>, never>>(object: I): QuerySendEnabledResponse;
};
/** Query defines the gRPC querier service. */
export interface Query {
    /** Balance queries the balance of a single coin for a single account. */
    Balance(request: QueryBalanceRequest): Promise<QueryBalanceResponse>;
    /**
     * AllBalances queries the balance of all coins for a single account.
     *
     * When called from another module, this query might consume a high amount of
     * gas if the pagination field is incorrectly set.
     */
    AllBalances(request: QueryAllBalancesRequest): Promise<QueryAllBalancesResponse>;
    /**
     * SpendableBalances queries the spendable balance of all coins for a single
     * account.
     *
     * When called from another module, this query might consume a high amount of
     * gas if the pagination field is incorrectly set.
     *
     * Since: cosmos-sdk 0.46
     */
    SpendableBalances(request: QuerySpendableBalancesRequest): Promise<QuerySpendableBalancesResponse>;
    /**
     * SpendableBalanceByDenom queries the spendable balance of a single denom for
     * a single account.
     *
     * When called from another module, this query might consume a high amount of
     * gas if the pagination field is incorrectly set.
     *
     * Since: cosmos-sdk 0.47
     */
    SpendableBalanceByDenom(request: QuerySpendableBalanceByDenomRequest): Promise<QuerySpendableBalanceByDenomResponse>;
    /**
     * TotalSupply queries the total supply of all coins.
     *
     * When called from another module, this query might consume a high amount of
     * gas if the pagination field is incorrectly set.
     */
    TotalSupply(request?: QueryTotalSupplyRequest): Promise<QueryTotalSupplyResponse>;
    /**
     * SupplyOf queries the supply of a single coin.
     *
     * When called from another module, this query might consume a high amount of
     * gas if the pagination field is incorrectly set.
     */
    SupplyOf(request: QuerySupplyOfRequest): Promise<QuerySupplyOfResponse>;
    /** Params queries the parameters of x/bank module. */
    Params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
    /** DenomsMetadata queries the client metadata of a given coin denomination. */
    DenomMetadata(request: QueryDenomMetadataRequest): Promise<QueryDenomMetadataResponse>;
    /**
     * DenomsMetadata queries the client metadata for all registered coin
     * denominations.
     */
    DenomsMetadata(request?: QueryDenomsMetadataRequest): Promise<QueryDenomsMetadataResponse>;
    /**
     * DenomOwners queries for all account addresses that own a particular token
     * denomination.
     *
     * When called from another module, this query might consume a high amount of
     * gas if the pagination field is incorrectly set.
     *
     * Since: cosmos-sdk 0.46
     */
    DenomOwners(request: QueryDenomOwnersRequest): Promise<QueryDenomOwnersResponse>;
    /**
     * SendEnabled queries for SendEnabled entries.
     *
     * This query only returns denominations that have specific SendEnabled settings.
     * Any denomination that does not have a specific setting will use the default
     * params.default_send_enabled, and will not be returned by this query.
     *
     * Since: cosmos-sdk 0.47
     */
    SendEnabled(request: QuerySendEnabledRequest): Promise<QuerySendEnabledResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    Balance(request: QueryBalanceRequest): Promise<QueryBalanceResponse>;
    AllBalances(request: QueryAllBalancesRequest): Promise<QueryAllBalancesResponse>;
    SpendableBalances(request: QuerySpendableBalancesRequest): Promise<QuerySpendableBalancesResponse>;
    SpendableBalanceByDenom(request: QuerySpendableBalanceByDenomRequest): Promise<QuerySpendableBalanceByDenomResponse>;
    TotalSupply(request?: QueryTotalSupplyRequest): Promise<QueryTotalSupplyResponse>;
    SupplyOf(request: QuerySupplyOfRequest): Promise<QuerySupplyOfResponse>;
    Params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
    DenomMetadata(request: QueryDenomMetadataRequest): Promise<QueryDenomMetadataResponse>;
    DenomsMetadata(request?: QueryDenomsMetadataRequest): Promise<QueryDenomsMetadataResponse>;
    DenomOwners(request: QueryDenomOwnersRequest): Promise<QueryDenomOwnersResponse>;
    SendEnabled(request: QuerySendEnabledRequest): Promise<QuerySendEnabledResponse>;
}
