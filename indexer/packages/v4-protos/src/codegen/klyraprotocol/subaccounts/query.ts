import { PageRequest, PageRequestSDKType, PageResponse, PageResponseSDKType } from "../../cosmos/base/query/v1beta1/pagination";
import { Subaccount, SubaccountSDKType } from "./subaccount";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** QueryGetSubaccountRequest is request type for the Query RPC method. */

export interface QueryGetSubaccountRequest {
  owner: string;
  number: number;
}
/** QueryGetSubaccountRequest is request type for the Query RPC method. */

export interface QueryGetSubaccountRequestSDKType {
  owner: string;
  number: number;
}
/** QuerySubaccountResponse is response type for the Query RPC method. */

export interface QuerySubaccountResponse {
  subaccount?: Subaccount;
}
/** QuerySubaccountResponse is response type for the Query RPC method. */

export interface QuerySubaccountResponseSDKType {
  subaccount?: SubaccountSDKType;
}
/** QueryAllSubaccountRequest is request type for the Query RPC method. */

export interface QueryAllSubaccountRequest {
  pagination?: PageRequest;
}
/** QueryAllSubaccountRequest is request type for the Query RPC method. */

export interface QueryAllSubaccountRequestSDKType {
  pagination?: PageRequestSDKType;
}
/** QuerySubaccountAllResponse is response type for the Query RPC method. */

export interface QuerySubaccountAllResponse {
  subaccount: Subaccount[];
  pagination?: PageResponse;
}
/** QuerySubaccountAllResponse is response type for the Query RPC method. */

export interface QuerySubaccountAllResponseSDKType {
  subaccount: SubaccountSDKType[];
  pagination?: PageResponseSDKType;
}
/**
 * QueryGetWithdrawalAndTransfersBlockedInfoRequest is a request type for
 * fetching information about whether withdrawals and transfers are blocked for
 * a collateral pool associated with the passed in perpetual id.
 */

export interface QueryGetWithdrawalAndTransfersBlockedInfoRequest {
  perpetualId: number;
}
/**
 * QueryGetWithdrawalAndTransfersBlockedInfoRequest is a request type for
 * fetching information about whether withdrawals and transfers are blocked for
 * a collateral pool associated with the passed in perpetual id.
 */

export interface QueryGetWithdrawalAndTransfersBlockedInfoRequestSDKType {
  perpetual_id: number;
}
/**
 * QueryGetWithdrawalAndTransfersBlockedInfoRequest is a response type for
 * fetching information about whether withdrawals and transfers are blocked.
 */

export interface QueryGetWithdrawalAndTransfersBlockedInfoResponse {
  negativeTncSubaccountSeenAtBlock: number;
  chainOutageSeenAtBlock: number;
  withdrawalsAndTransfersUnblockedAtBlock: number;
}
/**
 * QueryGetWithdrawalAndTransfersBlockedInfoRequest is a response type for
 * fetching information about whether withdrawals and transfers are blocked.
 */

export interface QueryGetWithdrawalAndTransfersBlockedInfoResponseSDKType {
  negative_tnc_subaccount_seen_at_block: number;
  chain_outage_seen_at_block: number;
  withdrawals_and_transfers_unblocked_at_block: number;
}
/**
 * QueryCollateralPoolAddressRequest is the request type for fetching the
 * account address of the collateral pool associated with the passed in
 * perpetual id.
 */

export interface QueryCollateralPoolAddressRequest {
  /**
   * QueryCollateralPoolAddressRequest is the request type for fetching the
   * account address of the collateral pool associated with the passed in
   * perpetual id.
   */
  perpetualId: number;
}
/**
 * QueryCollateralPoolAddressRequest is the request type for fetching the
 * account address of the collateral pool associated with the passed in
 * perpetual id.
 */

export interface QueryCollateralPoolAddressRequestSDKType {
  /**
   * QueryCollateralPoolAddressRequest is the request type for fetching the
   * account address of the collateral pool associated with the passed in
   * perpetual id.
   */
  perpetual_id: number;
}
/**
 * QueryCollateralPoolAddressResponse is a response type for fetching the
 * account address of the collateral pool associated with the passed in
 * perpetual id.
 */

export interface QueryCollateralPoolAddressResponse {
  collateralPoolAddress: string;
}
/**
 * QueryCollateralPoolAddressResponse is a response type for fetching the
 * account address of the collateral pool associated with the passed in
 * perpetual id.
 */

export interface QueryCollateralPoolAddressResponseSDKType {
  collateral_pool_address: string;
}

function createBaseQueryGetSubaccountRequest(): QueryGetSubaccountRequest {
  return {
    owner: "",
    number: 0
  };
}

export const QueryGetSubaccountRequest = {
  encode(message: QueryGetSubaccountRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.owner !== "") {
      writer.uint32(10).string(message.owner);
    }

    if (message.number !== 0) {
      writer.uint32(16).uint32(message.number);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetSubaccountRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetSubaccountRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.owner = reader.string();
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

  fromPartial(object: DeepPartial<QueryGetSubaccountRequest>): QueryGetSubaccountRequest {
    const message = createBaseQueryGetSubaccountRequest();
    message.owner = object.owner ?? "";
    message.number = object.number ?? 0;
    return message;
  }

};

function createBaseQuerySubaccountResponse(): QuerySubaccountResponse {
  return {
    subaccount: undefined
  };
}

export const QuerySubaccountResponse = {
  encode(message: QuerySubaccountResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.subaccount !== undefined) {
      Subaccount.encode(message.subaccount, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QuerySubaccountResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQuerySubaccountResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.subaccount = Subaccount.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QuerySubaccountResponse>): QuerySubaccountResponse {
    const message = createBaseQuerySubaccountResponse();
    message.subaccount = object.subaccount !== undefined && object.subaccount !== null ? Subaccount.fromPartial(object.subaccount) : undefined;
    return message;
  }

};

function createBaseQueryAllSubaccountRequest(): QueryAllSubaccountRequest {
  return {
    pagination: undefined
  };
}

export const QueryAllSubaccountRequest = {
  encode(message: QueryAllSubaccountRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllSubaccountRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllSubaccountRequest();

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

  fromPartial(object: DeepPartial<QueryAllSubaccountRequest>): QueryAllSubaccountRequest {
    const message = createBaseQueryAllSubaccountRequest();
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQuerySubaccountAllResponse(): QuerySubaccountAllResponse {
  return {
    subaccount: [],
    pagination: undefined
  };
}

export const QuerySubaccountAllResponse = {
  encode(message: QuerySubaccountAllResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.subaccount) {
      Subaccount.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QuerySubaccountAllResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQuerySubaccountAllResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.subaccount.push(Subaccount.decode(reader, reader.uint32()));
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

  fromPartial(object: DeepPartial<QuerySubaccountAllResponse>): QuerySubaccountAllResponse {
    const message = createBaseQuerySubaccountAllResponse();
    message.subaccount = object.subaccount?.map(e => Subaccount.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQueryGetWithdrawalAndTransfersBlockedInfoRequest(): QueryGetWithdrawalAndTransfersBlockedInfoRequest {
  return {
    perpetualId: 0
  };
}

export const QueryGetWithdrawalAndTransfersBlockedInfoRequest = {
  encode(message: QueryGetWithdrawalAndTransfersBlockedInfoRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.perpetualId !== 0) {
      writer.uint32(8).uint32(message.perpetualId);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetWithdrawalAndTransfersBlockedInfoRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetWithdrawalAndTransfersBlockedInfoRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.perpetualId = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryGetWithdrawalAndTransfersBlockedInfoRequest>): QueryGetWithdrawalAndTransfersBlockedInfoRequest {
    const message = createBaseQueryGetWithdrawalAndTransfersBlockedInfoRequest();
    message.perpetualId = object.perpetualId ?? 0;
    return message;
  }

};

function createBaseQueryGetWithdrawalAndTransfersBlockedInfoResponse(): QueryGetWithdrawalAndTransfersBlockedInfoResponse {
  return {
    negativeTncSubaccountSeenAtBlock: 0,
    chainOutageSeenAtBlock: 0,
    withdrawalsAndTransfersUnblockedAtBlock: 0
  };
}

export const QueryGetWithdrawalAndTransfersBlockedInfoResponse = {
  encode(message: QueryGetWithdrawalAndTransfersBlockedInfoResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.negativeTncSubaccountSeenAtBlock !== 0) {
      writer.uint32(8).uint32(message.negativeTncSubaccountSeenAtBlock);
    }

    if (message.chainOutageSeenAtBlock !== 0) {
      writer.uint32(16).uint32(message.chainOutageSeenAtBlock);
    }

    if (message.withdrawalsAndTransfersUnblockedAtBlock !== 0) {
      writer.uint32(24).uint32(message.withdrawalsAndTransfersUnblockedAtBlock);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetWithdrawalAndTransfersBlockedInfoResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetWithdrawalAndTransfersBlockedInfoResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.negativeTncSubaccountSeenAtBlock = reader.uint32();
          break;

        case 2:
          message.chainOutageSeenAtBlock = reader.uint32();
          break;

        case 3:
          message.withdrawalsAndTransfersUnblockedAtBlock = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryGetWithdrawalAndTransfersBlockedInfoResponse>): QueryGetWithdrawalAndTransfersBlockedInfoResponse {
    const message = createBaseQueryGetWithdrawalAndTransfersBlockedInfoResponse();
    message.negativeTncSubaccountSeenAtBlock = object.negativeTncSubaccountSeenAtBlock ?? 0;
    message.chainOutageSeenAtBlock = object.chainOutageSeenAtBlock ?? 0;
    message.withdrawalsAndTransfersUnblockedAtBlock = object.withdrawalsAndTransfersUnblockedAtBlock ?? 0;
    return message;
  }

};

function createBaseQueryCollateralPoolAddressRequest(): QueryCollateralPoolAddressRequest {
  return {
    perpetualId: 0
  };
}

export const QueryCollateralPoolAddressRequest = {
  encode(message: QueryCollateralPoolAddressRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.perpetualId !== 0) {
      writer.uint32(8).uint32(message.perpetualId);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryCollateralPoolAddressRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryCollateralPoolAddressRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.perpetualId = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryCollateralPoolAddressRequest>): QueryCollateralPoolAddressRequest {
    const message = createBaseQueryCollateralPoolAddressRequest();
    message.perpetualId = object.perpetualId ?? 0;
    return message;
  }

};

function createBaseQueryCollateralPoolAddressResponse(): QueryCollateralPoolAddressResponse {
  return {
    collateralPoolAddress: ""
  };
}

export const QueryCollateralPoolAddressResponse = {
  encode(message: QueryCollateralPoolAddressResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.collateralPoolAddress !== "") {
      writer.uint32(10).string(message.collateralPoolAddress);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryCollateralPoolAddressResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryCollateralPoolAddressResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.collateralPoolAddress = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryCollateralPoolAddressResponse>): QueryCollateralPoolAddressResponse {
    const message = createBaseQueryCollateralPoolAddressResponse();
    message.collateralPoolAddress = object.collateralPoolAddress ?? "";
    return message;
  }

};