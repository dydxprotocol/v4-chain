import { VaultType, VaultTypeSDKType, VaultId, VaultIdSDKType } from "./vault";
import { PageRequest, PageRequestSDKType, PageResponse, PageResponseSDKType } from "../../cosmos/base/query/v1beta1/pagination";
import { NumShares, NumSharesSDKType, ShareUnlock, ShareUnlockSDKType, OwnerShare, OwnerShareSDKType } from "./share";
import { QuotingParams, QuotingParamsSDKType, OperatorParams, OperatorParamsSDKType, VaultParams, VaultParamsSDKType } from "./params";
import { SubaccountId, SubaccountIdSDKType } from "../subaccounts/subaccount";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** QueryParamsRequest is a request type for the Params RPC method. */

export interface QueryParamsRequest {}
/** QueryParamsRequest is a request type for the Params RPC method. */

export interface QueryParamsRequestSDKType {}
/** QueryParamsResponse is a response type for the Params RPC method. */

export interface QueryParamsResponse {
  defaultQuotingParams?: QuotingParams;
  operatorParams?: OperatorParams;
}
/** QueryParamsResponse is a response type for the Params RPC method. */

export interface QueryParamsResponseSDKType {
  default_quoting_params?: QuotingParamsSDKType;
  operator_params?: OperatorParamsSDKType;
}
/** QueryVaultRequest is a request type for the Vault RPC method. */

export interface QueryVaultRequest {
  type: VaultType;
  number: number;
}
/** QueryVaultRequest is a request type for the Vault RPC method. */

export interface QueryVaultRequestSDKType {
  type: VaultTypeSDKType;
  number: number;
}
/** QueryVaultResponse is a response type for the Vault RPC method. */

export interface QueryVaultResponse {
  vaultId?: VaultId;
  subaccountId?: SubaccountId;
  equity: Uint8Array;
  inventory: Uint8Array;
  vaultParams?: VaultParams;
  mostRecentClientIds: number[];
}
/** QueryVaultResponse is a response type for the Vault RPC method. */

export interface QueryVaultResponseSDKType {
  vault_id?: VaultIdSDKType;
  subaccount_id?: SubaccountIdSDKType;
  equity: Uint8Array;
  inventory: Uint8Array;
  vault_params?: VaultParamsSDKType;
  most_recent_client_ids: number[];
}
/** QueryAllVaultsRequest is a request type for the AllVaults RPC method. */

export interface QueryAllVaultsRequest {
  pagination?: PageRequest;
}
/** QueryAllVaultsRequest is a request type for the AllVaults RPC method. */

export interface QueryAllVaultsRequestSDKType {
  pagination?: PageRequestSDKType;
}
/** QueryAllVaultsResponse is a response type for the AllVaults RPC method. */

export interface QueryAllVaultsResponse {
  vaults: QueryVaultResponse[];
  pagination?: PageResponse;
}
/** QueryAllVaultsResponse is a response type for the AllVaults RPC method. */

export interface QueryAllVaultsResponseSDKType {
  vaults: QueryVaultResponseSDKType[];
  pagination?: PageResponseSDKType;
}
/**
 * QueryMegavaultTotalSharesRequest is a request type for the
 * MegavaultTotalShares RPC method.
 */

export interface QueryMegavaultTotalSharesRequest {}
/**
 * QueryMegavaultTotalSharesRequest is a request type for the
 * MegavaultTotalShares RPC method.
 */

export interface QueryMegavaultTotalSharesRequestSDKType {}
/**
 * QueryMegavaultTotalSharesResponse is a response type for the
 * MegavaultTotalShares RPC method.
 */

export interface QueryMegavaultTotalSharesResponse {
  /**
   * QueryMegavaultTotalSharesResponse is a response type for the
   * MegavaultTotalShares RPC method.
   */
  totalShares?: NumShares;
}
/**
 * QueryMegavaultTotalSharesResponse is a response type for the
 * MegavaultTotalShares RPC method.
 */

export interface QueryMegavaultTotalSharesResponseSDKType {
  /**
   * QueryMegavaultTotalSharesResponse is a response type for the
   * MegavaultTotalShares RPC method.
   */
  total_shares?: NumSharesSDKType;
}
/**
 * QueryMegavaultOwnerSharesRequest is a request type for the
 * MegavaultOwnerShares RPC method.
 */

export interface QueryMegavaultOwnerSharesRequest {
  address: string;
}
/**
 * QueryMegavaultOwnerSharesRequest is a request type for the
 * MegavaultOwnerShares RPC method.
 */

export interface QueryMegavaultOwnerSharesRequestSDKType {
  address: string;
}
/**
 * QueryMegavaultOwnerSharesResponse is a response type for the
 * MegavaultOwnerShares RPC method.
 */

export interface QueryMegavaultOwnerSharesResponse {
  /** Owner address. */
  address: string;
  /** Total number of shares that belong to the owner. */

  shares?: NumShares;
  /** All share unlocks. */

  shareUnlocks: ShareUnlock[];
  /** Owner equity in megavault (in quote quantums). */

  equity: Uint8Array;
  /**
   * Equity that owner can withdraw in quote quantums (as one cannot
   * withdraw locked shares).
   */

  withdrawableEquity: Uint8Array;
}
/**
 * QueryMegavaultOwnerSharesResponse is a response type for the
 * MegavaultOwnerShares RPC method.
 */

export interface QueryMegavaultOwnerSharesResponseSDKType {
  /** Owner address. */
  address: string;
  /** Total number of shares that belong to the owner. */

  shares?: NumSharesSDKType;
  /** All share unlocks. */

  share_unlocks: ShareUnlockSDKType[];
  /** Owner equity in megavault (in quote quantums). */

  equity: Uint8Array;
  /**
   * Equity that owner can withdraw in quote quantums (as one cannot
   * withdraw locked shares).
   */

  withdrawable_equity: Uint8Array;
}
/**
 * QueryMegavaultAllOwnerSharesRequest is a request type for the
 * MegavaultAllOwnerShares RPC method.
 */

export interface QueryMegavaultAllOwnerSharesRequest {
  pagination?: PageRequest;
}
/**
 * QueryMegavaultAllOwnerSharesRequest is a request type for the
 * MegavaultAllOwnerShares RPC method.
 */

export interface QueryMegavaultAllOwnerSharesRequestSDKType {
  pagination?: PageRequestSDKType;
}
/**
 * QueryMegavaultAllOwnerSharesResponse is a response type for the
 * MegavaultAllOwnerShares RPC method.
 */

export interface QueryMegavaultAllOwnerSharesResponse {
  ownerShares: OwnerShare[];
  pagination?: PageResponse;
}
/**
 * QueryMegavaultAllOwnerSharesResponse is a response type for the
 * MegavaultAllOwnerShares RPC method.
 */

export interface QueryMegavaultAllOwnerSharesResponseSDKType {
  owner_shares: OwnerShareSDKType[];
  pagination?: PageResponseSDKType;
}
/** QueryVaultParamsRequest is a request for the VaultParams RPC method. */

export interface QueryVaultParamsRequest {
  type: VaultType;
  number: number;
}
/** QueryVaultParamsRequest is a request for the VaultParams RPC method. */

export interface QueryVaultParamsRequestSDKType {
  type: VaultTypeSDKType;
  number: number;
}
/** QueryVaultParamsResponse is a response for the VaultParams RPC method. */

export interface QueryVaultParamsResponse {
  vaultId?: VaultId;
  vaultParams?: VaultParams;
}
/** QueryVaultParamsResponse is a response for the VaultParams RPC method. */

export interface QueryVaultParamsResponseSDKType {
  vault_id?: VaultIdSDKType;
  vault_params?: VaultParamsSDKType;
}
/**
 * QueryMegavaultWithdrawalInfoRequest is a request type for the
 * MegavaultWithdrawalInfo RPC method.
 */

export interface QueryMegavaultWithdrawalInfoRequest {
  /** Number of shares to withdraw. */
  sharesToWithdraw?: NumShares;
}
/**
 * QueryMegavaultWithdrawalInfoRequest is a request type for the
 * MegavaultWithdrawalInfo RPC method.
 */

export interface QueryMegavaultWithdrawalInfoRequestSDKType {
  /** Number of shares to withdraw. */
  shares_to_withdraw?: NumSharesSDKType;
}
/**
 * QueryMegavaultWithdrawalInfoResponse is a response type for the
 * MegavaultWithdrawalInfo RPC method.
 */

export interface QueryMegavaultWithdrawalInfoResponse {
  /** Number of shares to withdraw. */
  sharesToWithdraw?: NumShares;
  /**
   * Number of quote quantums above `shares` are expected to redeem.
   * Withdrawl slippage can be calculated by comparing
   * `expected_quote_quantums` with
   * `megavault_equity * shares_to_withdraw / total_shares`
   */

  expectedQuoteQuantums: Uint8Array;
  /** Equity of megavault (in quote quantums). */

  megavaultEquity: Uint8Array;
  /** Total shares in megavault. */

  totalShares?: NumShares;
}
/**
 * QueryMegavaultWithdrawalInfoResponse is a response type for the
 * MegavaultWithdrawalInfo RPC method.
 */

export interface QueryMegavaultWithdrawalInfoResponseSDKType {
  /** Number of shares to withdraw. */
  shares_to_withdraw?: NumSharesSDKType;
  /**
   * Number of quote quantums above `shares` are expected to redeem.
   * Withdrawl slippage can be calculated by comparing
   * `expected_quote_quantums` with
   * `megavault_equity * shares_to_withdraw / total_shares`
   */

  expected_quote_quantums: Uint8Array;
  /** Equity of megavault (in quote quantums). */

  megavault_equity: Uint8Array;
  /** Total shares in megavault. */

  total_shares?: NumSharesSDKType;
}

function createBaseQueryParamsRequest(): QueryParamsRequest {
  return {};
}

export const QueryParamsRequest = {
  encode(_: QueryParamsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryParamsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryParamsRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(_: DeepPartial<QueryParamsRequest>): QueryParamsRequest {
    const message = createBaseQueryParamsRequest();
    return message;
  }

};

function createBaseQueryParamsResponse(): QueryParamsResponse {
  return {
    defaultQuotingParams: undefined,
    operatorParams: undefined
  };
}

export const QueryParamsResponse = {
  encode(message: QueryParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.defaultQuotingParams !== undefined) {
      QuotingParams.encode(message.defaultQuotingParams, writer.uint32(10).fork()).ldelim();
    }

    if (message.operatorParams !== undefined) {
      OperatorParams.encode(message.operatorParams, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryParamsResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.defaultQuotingParams = QuotingParams.decode(reader, reader.uint32());
          break;

        case 2:
          message.operatorParams = OperatorParams.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryParamsResponse>): QueryParamsResponse {
    const message = createBaseQueryParamsResponse();
    message.defaultQuotingParams = object.defaultQuotingParams !== undefined && object.defaultQuotingParams !== null ? QuotingParams.fromPartial(object.defaultQuotingParams) : undefined;
    message.operatorParams = object.operatorParams !== undefined && object.operatorParams !== null ? OperatorParams.fromPartial(object.operatorParams) : undefined;
    return message;
  }

};

function createBaseQueryVaultRequest(): QueryVaultRequest {
  return {
    type: 0,
    number: 0
  };
}

export const QueryVaultRequest = {
  encode(message: QueryVaultRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.type !== 0) {
      writer.uint32(8).int32(message.type);
    }

    if (message.number !== 0) {
      writer.uint32(16).uint32(message.number);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryVaultRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryVaultRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.type = (reader.int32() as any);
          break;

        case 2:
          message.number = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryVaultRequest>): QueryVaultRequest {
    const message = createBaseQueryVaultRequest();
    message.type = object.type ?? 0;
    message.number = object.number ?? 0;
    return message;
  }

};

function createBaseQueryVaultResponse(): QueryVaultResponse {
  return {
    vaultId: undefined,
    subaccountId: undefined,
    equity: new Uint8Array(),
    inventory: new Uint8Array(),
    vaultParams: undefined,
    mostRecentClientIds: []
  };
}

export const QueryVaultResponse = {
  encode(message: QueryVaultResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.vaultId !== undefined) {
      VaultId.encode(message.vaultId, writer.uint32(10).fork()).ldelim();
    }

    if (message.subaccountId !== undefined) {
      SubaccountId.encode(message.subaccountId, writer.uint32(18).fork()).ldelim();
    }

    if (message.equity.length !== 0) {
      writer.uint32(26).bytes(message.equity);
    }

    if (message.inventory.length !== 0) {
      writer.uint32(34).bytes(message.inventory);
    }

    if (message.vaultParams !== undefined) {
      VaultParams.encode(message.vaultParams, writer.uint32(42).fork()).ldelim();
    }

    writer.uint32(50).fork();

    for (const v of message.mostRecentClientIds) {
      writer.uint32(v);
    }

    writer.ldelim();
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryVaultResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryVaultResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.vaultId = VaultId.decode(reader, reader.uint32());
          break;

        case 2:
          message.subaccountId = SubaccountId.decode(reader, reader.uint32());
          break;

        case 3:
          message.equity = reader.bytes();
          break;

        case 4:
          message.inventory = reader.bytes();
          break;

        case 5:
          message.vaultParams = VaultParams.decode(reader, reader.uint32());
          break;

        case 6:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;

            while (reader.pos < end2) {
              message.mostRecentClientIds.push(reader.uint32());
            }
          } else {
            message.mostRecentClientIds.push(reader.uint32());
          }

          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryVaultResponse>): QueryVaultResponse {
    const message = createBaseQueryVaultResponse();
    message.vaultId = object.vaultId !== undefined && object.vaultId !== null ? VaultId.fromPartial(object.vaultId) : undefined;
    message.subaccountId = object.subaccountId !== undefined && object.subaccountId !== null ? SubaccountId.fromPartial(object.subaccountId) : undefined;
    message.equity = object.equity ?? new Uint8Array();
    message.inventory = object.inventory ?? new Uint8Array();
    message.vaultParams = object.vaultParams !== undefined && object.vaultParams !== null ? VaultParams.fromPartial(object.vaultParams) : undefined;
    message.mostRecentClientIds = object.mostRecentClientIds?.map(e => e) || [];
    return message;
  }

};

function createBaseQueryAllVaultsRequest(): QueryAllVaultsRequest {
  return {
    pagination: undefined
  };
}

export const QueryAllVaultsRequest = {
  encode(message: QueryAllVaultsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllVaultsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllVaultsRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryAllVaultsRequest>): QueryAllVaultsRequest {
    const message = createBaseQueryAllVaultsRequest();
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQueryAllVaultsResponse(): QueryAllVaultsResponse {
  return {
    vaults: [],
    pagination: undefined
  };
}

export const QueryAllVaultsResponse = {
  encode(message: QueryAllVaultsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.vaults) {
      QueryVaultResponse.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllVaultsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllVaultsResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.vaults.push(QueryVaultResponse.decode(reader, reader.uint32()));
          break;

        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryAllVaultsResponse>): QueryAllVaultsResponse {
    const message = createBaseQueryAllVaultsResponse();
    message.vaults = object.vaults?.map(e => QueryVaultResponse.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQueryMegavaultTotalSharesRequest(): QueryMegavaultTotalSharesRequest {
  return {};
}

export const QueryMegavaultTotalSharesRequest = {
  encode(_: QueryMegavaultTotalSharesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryMegavaultTotalSharesRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryMegavaultTotalSharesRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(_: DeepPartial<QueryMegavaultTotalSharesRequest>): QueryMegavaultTotalSharesRequest {
    const message = createBaseQueryMegavaultTotalSharesRequest();
    return message;
  }

};

function createBaseQueryMegavaultTotalSharesResponse(): QueryMegavaultTotalSharesResponse {
  return {
    totalShares: undefined
  };
}

export const QueryMegavaultTotalSharesResponse = {
  encode(message: QueryMegavaultTotalSharesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.totalShares !== undefined) {
      NumShares.encode(message.totalShares, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryMegavaultTotalSharesResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryMegavaultTotalSharesResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.totalShares = NumShares.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryMegavaultTotalSharesResponse>): QueryMegavaultTotalSharesResponse {
    const message = createBaseQueryMegavaultTotalSharesResponse();
    message.totalShares = object.totalShares !== undefined && object.totalShares !== null ? NumShares.fromPartial(object.totalShares) : undefined;
    return message;
  }

};

function createBaseQueryMegavaultOwnerSharesRequest(): QueryMegavaultOwnerSharesRequest {
  return {
    address: ""
  };
}

export const QueryMegavaultOwnerSharesRequest = {
  encode(message: QueryMegavaultOwnerSharesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryMegavaultOwnerSharesRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryMegavaultOwnerSharesRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryMegavaultOwnerSharesRequest>): QueryMegavaultOwnerSharesRequest {
    const message = createBaseQueryMegavaultOwnerSharesRequest();
    message.address = object.address ?? "";
    return message;
  }

};

function createBaseQueryMegavaultOwnerSharesResponse(): QueryMegavaultOwnerSharesResponse {
  return {
    address: "",
    shares: undefined,
    shareUnlocks: [],
    equity: new Uint8Array(),
    withdrawableEquity: new Uint8Array()
  };
}

export const QueryMegavaultOwnerSharesResponse = {
  encode(message: QueryMegavaultOwnerSharesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }

    if (message.shares !== undefined) {
      NumShares.encode(message.shares, writer.uint32(18).fork()).ldelim();
    }

    for (const v of message.shareUnlocks) {
      ShareUnlock.encode(v!, writer.uint32(26).fork()).ldelim();
    }

    if (message.equity.length !== 0) {
      writer.uint32(34).bytes(message.equity);
    }

    if (message.withdrawableEquity.length !== 0) {
      writer.uint32(42).bytes(message.withdrawableEquity);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryMegavaultOwnerSharesResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryMegavaultOwnerSharesResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;

        case 2:
          message.shares = NumShares.decode(reader, reader.uint32());
          break;

        case 3:
          message.shareUnlocks.push(ShareUnlock.decode(reader, reader.uint32()));
          break;

        case 4:
          message.equity = reader.bytes();
          break;

        case 5:
          message.withdrawableEquity = reader.bytes();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryMegavaultOwnerSharesResponse>): QueryMegavaultOwnerSharesResponse {
    const message = createBaseQueryMegavaultOwnerSharesResponse();
    message.address = object.address ?? "";
    message.shares = object.shares !== undefined && object.shares !== null ? NumShares.fromPartial(object.shares) : undefined;
    message.shareUnlocks = object.shareUnlocks?.map(e => ShareUnlock.fromPartial(e)) || [];
    message.equity = object.equity ?? new Uint8Array();
    message.withdrawableEquity = object.withdrawableEquity ?? new Uint8Array();
    return message;
  }

};

function createBaseQueryMegavaultAllOwnerSharesRequest(): QueryMegavaultAllOwnerSharesRequest {
  return {
    pagination: undefined
  };
}

export const QueryMegavaultAllOwnerSharesRequest = {
  encode(message: QueryMegavaultAllOwnerSharesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryMegavaultAllOwnerSharesRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryMegavaultAllOwnerSharesRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryMegavaultAllOwnerSharesRequest>): QueryMegavaultAllOwnerSharesRequest {
    const message = createBaseQueryMegavaultAllOwnerSharesRequest();
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQueryMegavaultAllOwnerSharesResponse(): QueryMegavaultAllOwnerSharesResponse {
  return {
    ownerShares: [],
    pagination: undefined
  };
}

export const QueryMegavaultAllOwnerSharesResponse = {
  encode(message: QueryMegavaultAllOwnerSharesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.ownerShares) {
      OwnerShare.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryMegavaultAllOwnerSharesResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryMegavaultAllOwnerSharesResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.ownerShares.push(OwnerShare.decode(reader, reader.uint32()));
          break;

        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryMegavaultAllOwnerSharesResponse>): QueryMegavaultAllOwnerSharesResponse {
    const message = createBaseQueryMegavaultAllOwnerSharesResponse();
    message.ownerShares = object.ownerShares?.map(e => OwnerShare.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQueryVaultParamsRequest(): QueryVaultParamsRequest {
  return {
    type: 0,
    number: 0
  };
}

export const QueryVaultParamsRequest = {
  encode(message: QueryVaultParamsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.type !== 0) {
      writer.uint32(8).int32(message.type);
    }

    if (message.number !== 0) {
      writer.uint32(16).uint32(message.number);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryVaultParamsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryVaultParamsRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.type = (reader.int32() as any);
          break;

        case 2:
          message.number = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryVaultParamsRequest>): QueryVaultParamsRequest {
    const message = createBaseQueryVaultParamsRequest();
    message.type = object.type ?? 0;
    message.number = object.number ?? 0;
    return message;
  }

};

function createBaseQueryVaultParamsResponse(): QueryVaultParamsResponse {
  return {
    vaultId: undefined,
    vaultParams: undefined
  };
}

export const QueryVaultParamsResponse = {
  encode(message: QueryVaultParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.vaultId !== undefined) {
      VaultId.encode(message.vaultId, writer.uint32(10).fork()).ldelim();
    }

    if (message.vaultParams !== undefined) {
      VaultParams.encode(message.vaultParams, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryVaultParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryVaultParamsResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.vaultId = VaultId.decode(reader, reader.uint32());
          break;

        case 2:
          message.vaultParams = VaultParams.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryVaultParamsResponse>): QueryVaultParamsResponse {
    const message = createBaseQueryVaultParamsResponse();
    message.vaultId = object.vaultId !== undefined && object.vaultId !== null ? VaultId.fromPartial(object.vaultId) : undefined;
    message.vaultParams = object.vaultParams !== undefined && object.vaultParams !== null ? VaultParams.fromPartial(object.vaultParams) : undefined;
    return message;
  }

};

function createBaseQueryMegavaultWithdrawalInfoRequest(): QueryMegavaultWithdrawalInfoRequest {
  return {
    sharesToWithdraw: undefined
  };
}

export const QueryMegavaultWithdrawalInfoRequest = {
  encode(message: QueryMegavaultWithdrawalInfoRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.sharesToWithdraw !== undefined) {
      NumShares.encode(message.sharesToWithdraw, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryMegavaultWithdrawalInfoRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryMegavaultWithdrawalInfoRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.sharesToWithdraw = NumShares.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryMegavaultWithdrawalInfoRequest>): QueryMegavaultWithdrawalInfoRequest {
    const message = createBaseQueryMegavaultWithdrawalInfoRequest();
    message.sharesToWithdraw = object.sharesToWithdraw !== undefined && object.sharesToWithdraw !== null ? NumShares.fromPartial(object.sharesToWithdraw) : undefined;
    return message;
  }

};

function createBaseQueryMegavaultWithdrawalInfoResponse(): QueryMegavaultWithdrawalInfoResponse {
  return {
    sharesToWithdraw: undefined,
    expectedQuoteQuantums: new Uint8Array(),
    megavaultEquity: new Uint8Array(),
    totalShares: undefined
  };
}

export const QueryMegavaultWithdrawalInfoResponse = {
  encode(message: QueryMegavaultWithdrawalInfoResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.sharesToWithdraw !== undefined) {
      NumShares.encode(message.sharesToWithdraw, writer.uint32(10).fork()).ldelim();
    }

    if (message.expectedQuoteQuantums.length !== 0) {
      writer.uint32(18).bytes(message.expectedQuoteQuantums);
    }

    if (message.megavaultEquity.length !== 0) {
      writer.uint32(26).bytes(message.megavaultEquity);
    }

    if (message.totalShares !== undefined) {
      NumShares.encode(message.totalShares, writer.uint32(34).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryMegavaultWithdrawalInfoResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryMegavaultWithdrawalInfoResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.sharesToWithdraw = NumShares.decode(reader, reader.uint32());
          break;

        case 2:
          message.expectedQuoteQuantums = reader.bytes();
          break;

        case 3:
          message.megavaultEquity = reader.bytes();
          break;

        case 4:
          message.totalShares = NumShares.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryMegavaultWithdrawalInfoResponse>): QueryMegavaultWithdrawalInfoResponse {
    const message = createBaseQueryMegavaultWithdrawalInfoResponse();
    message.sharesToWithdraw = object.sharesToWithdraw !== undefined && object.sharesToWithdraw !== null ? NumShares.fromPartial(object.sharesToWithdraw) : undefined;
    message.expectedQuoteQuantums = object.expectedQuoteQuantums ?? new Uint8Array();
    message.megavaultEquity = object.megavaultEquity ?? new Uint8Array();
    message.totalShares = object.totalShares !== undefined && object.totalShares !== null ? NumShares.fromPartial(object.totalShares) : undefined;
    return message;
  }

};