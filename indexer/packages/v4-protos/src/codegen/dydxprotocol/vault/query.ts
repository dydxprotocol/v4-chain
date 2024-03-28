import { VaultType, VaultTypeSDKType, VaultId, VaultIdSDKType } from "./vault";
import { Params, ParamsSDKType } from "./params";
import { SubaccountId, SubaccountIdSDKType } from "../subaccounts/subaccount";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial, Long } from "../../helpers";
/** QueryParamsRequest is a request type for the Params RPC method. */

export interface QueryParamsRequest {}
/** QueryParamsRequest is a request type for the Params RPC method. */

export interface QueryParamsRequestSDKType {}
/** QueryParamsResponse is a response type for the Params RPC method. */

export interface QueryParamsResponse {
  params?: Params;
}
/** QueryParamsResponse is a response type for the Params RPC method. */

export interface QueryParamsResponseSDKType {
  params?: ParamsSDKType;
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
  equity: Long;
  inventory: Long;
  totalShares: Long;
  allOwnerShares: OwnerShares[];
}
/** QueryVaultResponse is a response type for the Vault RPC method. */

export interface QueryVaultResponseSDKType {
  vault_id?: VaultIdSDKType;
  subaccount_id?: SubaccountIdSDKType;
  equity: Long;
  inventory: Long;
  total_shares: Long;
  all_owner_shares: OwnerSharesSDKType[];
}
/** OwnerShares is a message type for an owner and their shares. */

export interface OwnerShares {
  owner: string;
  shares: Long;
}
/** OwnerShares is a message type for an owner and their shares. */

export interface OwnerSharesSDKType {
  owner: string;
  shares: Long;
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
    params: undefined
  };
}

export const QueryParamsResponse = {
  encode(message: QueryParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
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
    equity: Long.UZERO,
    inventory: Long.UZERO,
    totalShares: Long.UZERO,
    allOwnerShares: []
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

    if (!message.equity.isZero()) {
      writer.uint32(24).uint64(message.equity);
    }

    if (!message.inventory.isZero()) {
      writer.uint32(32).uint64(message.inventory);
    }

    if (!message.totalShares.isZero()) {
      writer.uint32(40).uint64(message.totalShares);
    }

    for (const v of message.allOwnerShares) {
      OwnerShares.encode(v!, writer.uint32(50).fork()).ldelim();
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
          message.equity = (reader.uint64() as Long);
          break;

        case 4:
          message.inventory = (reader.uint64() as Long);
          break;

        case 5:
          message.totalShares = (reader.uint64() as Long);
          break;

        case 6:
          message.allOwnerShares.push(OwnerShares.decode(reader, reader.uint32()));
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
    message.equity = object.equity !== undefined && object.equity !== null ? Long.fromValue(object.equity) : Long.UZERO;
    message.inventory = object.inventory !== undefined && object.inventory !== null ? Long.fromValue(object.inventory) : Long.UZERO;
    message.totalShares = object.totalShares !== undefined && object.totalShares !== null ? Long.fromValue(object.totalShares) : Long.UZERO;
    message.allOwnerShares = object.allOwnerShares?.map(e => OwnerShares.fromPartial(e)) || [];
    return message;
  }

};

function createBaseOwnerShares(): OwnerShares {
  return {
    owner: "",
    shares: Long.UZERO
  };
}

export const OwnerShares = {
  encode(message: OwnerShares, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.owner !== "") {
      writer.uint32(10).string(message.owner);
    }

    if (!message.shares.isZero()) {
      writer.uint32(16).uint64(message.shares);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): OwnerShares {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOwnerShares();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.owner = reader.string();
          break;

        case 2:
          message.shares = (reader.uint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<OwnerShares>): OwnerShares {
    const message = createBaseOwnerShares();
    message.owner = object.owner ?? "";
    message.shares = object.shares !== undefined && object.shares !== null ? Long.fromValue(object.shares) : Long.UZERO;
    return message;
  }

};