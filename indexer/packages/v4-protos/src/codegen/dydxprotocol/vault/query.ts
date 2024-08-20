import { VaultType, VaultTypeSDKType, VaultId, VaultIdSDKType } from "./vault";
import { PageRequest, PageRequestSDKType, PageResponse, PageResponseSDKType } from "../../cosmos/base/query/v1beta1/pagination";
import { Params, ParamsSDKType, QuotingParams, QuotingParamsSDKType, VaultParams, VaultParamsSDKType } from "./params";
import { SubaccountId, SubaccountIdSDKType } from "../subaccounts/subaccount";
import { NumShares, NumSharesSDKType, OwnerShare, OwnerShareSDKType } from "./share";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** QueryParamsRequest is a request type for the Params RPC method. */

export interface QueryParamsRequest {}
/** QueryParamsRequest is a request type for the Params RPC method. */

export interface QueryParamsRequestSDKType {}
/** QueryParamsResponse is a response type for the Params RPC method. */

export interface QueryParamsResponse {
  /** Deprecated since v6.x in favor of default_quoting_params. */

  /** @deprecated */
  params?: Params;
  defaultQuotingParams?: QuotingParams;
}
/** QueryParamsResponse is a response type for the Params RPC method. */

export interface QueryParamsResponseSDKType {
  /** Deprecated since v6.x in favor of default_quoting_params. */

  /** @deprecated */
  params?: ParamsSDKType;
  default_quoting_params?: QuotingParamsSDKType;
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
}
/** QueryVaultResponse is a response type for the Vault RPC method. */

export interface QueryVaultResponseSDKType {
  vault_id?: VaultIdSDKType;
  subaccount_id?: SubaccountIdSDKType;
  equity: Uint8Array;
  inventory: Uint8Array;
  vault_params?: VaultParamsSDKType;
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
  pagination?: PageRequest;
}
/**
 * QueryMegavaultOwnerSharesRequest is a request type for the
 * MegavaultOwnerShares RPC method.
 */

export interface QueryMegavaultOwnerSharesRequestSDKType {
  pagination?: PageRequestSDKType;
}
/**
 * QueryMegavaultOwnerSharesResponse is a response type for the
 * MegavaultOwnerShares RPC method.
 */

export interface QueryMegavaultOwnerSharesResponse {
  ownerShares: OwnerShare[];
  pagination?: PageResponse;
}
/**
 * QueryMegavaultOwnerSharesResponse is a response type for the
 * MegavaultOwnerShares RPC method.
 */

export interface QueryMegavaultOwnerSharesResponseSDKType {
  owner_shares: OwnerShareSDKType[];
  pagination?: PageResponseSDKType;
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
    params: undefined,
    defaultQuotingParams: undefined
  };
}

export const QueryParamsResponse = {
  encode(message: QueryParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }

    if (message.defaultQuotingParams !== undefined) {
      QuotingParams.encode(message.defaultQuotingParams, writer.uint32(18).fork()).ldelim();
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
          message.params = Params.decode(reader, reader.uint32());
          break;

        case 2:
          message.defaultQuotingParams = QuotingParams.decode(reader, reader.uint32());
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
    message.params = object.params !== undefined && object.params !== null ? Params.fromPartial(object.params) : undefined;
    message.defaultQuotingParams = object.defaultQuotingParams !== undefined && object.defaultQuotingParams !== null ? QuotingParams.fromPartial(object.defaultQuotingParams) : undefined;
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
    vaultParams: undefined
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
    pagination: undefined
  };
}

export const QueryMegavaultOwnerSharesRequest = {
  encode(message: QueryMegavaultOwnerSharesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(26).fork()).ldelim();
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
        case 3:
          message.pagination = PageRequest.decode(reader, reader.uint32());
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
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQueryMegavaultOwnerSharesResponse(): QueryMegavaultOwnerSharesResponse {
  return {
    ownerShares: [],
    pagination: undefined
  };
}

export const QueryMegavaultOwnerSharesResponse = {
  encode(message: QueryMegavaultOwnerSharesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.ownerShares) {
      OwnerShare.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
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

  fromPartial(object: DeepPartial<QueryMegavaultOwnerSharesResponse>): QueryMegavaultOwnerSharesResponse {
    const message = createBaseQueryMegavaultOwnerSharesResponse();
    message.ownerShares = object.ownerShares?.map(e => OwnerShare.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  }

};