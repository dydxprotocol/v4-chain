import { PageRequest, PageRequestSDKType, PageResponse, PageResponseSDKType } from "../../cosmos/base/query/v1beta1/pagination";
import { Controls, ControlsSDKType, Ledger, LedgerSDKType, PendingSettlement, PendingSettlementSDKType } from "./canonical_usdc";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial, Long } from "../../helpers";
/** QueryStateRequest requests canonical controls and ledger state. */

export interface QueryStateRequest {}
/** QueryStateRequest requests canonical controls and ledger state. */

export interface QueryStateRequestSDKType {}
/** QueryStateResponse contains canonical controls and ledger state. */

export interface QueryStateResponse {
  /** controls is the current canonical routing controls. */
  controls?: Controls;
  /** ledger is the current backing and in-flight accounting. */

  ledger?: Ledger;
}
/** QueryStateResponse contains canonical controls and ledger state. */

export interface QueryStateResponseSDKType {
  /** controls is the current canonical routing controls. */
  controls?: ControlsSDKType;
  /** ledger is the current backing and in-flight accounting. */

  ledger?: LedgerSDKType;
}
/** QueryPendingSettlementsRequest requests a page of in-flight settlements. */

export interface QueryPendingSettlementsRequest {
  /** pagination defines the page request. */
  pagination?: PageRequest;
}
/** QueryPendingSettlementsRequest requests a page of in-flight settlements. */

export interface QueryPendingSettlementsRequestSDKType {
  /** pagination defines the page request. */
  pagination?: PageRequestSDKType;
}
/** QueryPendingSettlementsResponse contains a page of in-flight settlements. */

export interface QueryPendingSettlementsResponse {
  /** pending_settlements is the page of in-flight settlements. */
  pendingSettlements: PendingSettlement[];
  /** pagination defines the page response. */

  pagination?: PageResponse;
}
/** QueryPendingSettlementsResponse contains a page of in-flight settlements. */

export interface QueryPendingSettlementsResponseSDKType {
  /** pending_settlements is the page of in-flight settlements. */
  pending_settlements: PendingSettlementSDKType[];
  /** pagination defines the page response. */

  pagination?: PageResponseSDKType;
}
/** QueryPendingSettlementRequest requests one in-flight settlement. */

export interface QueryPendingSettlementRequest {
  /** source_channel is the source channel of the IBC packet. */
  sourceChannel: string;
  /** sequence is the sequence of the IBC packet. */

  sequence: Long;
}
/** QueryPendingSettlementRequest requests one in-flight settlement. */

export interface QueryPendingSettlementRequestSDKType {
  /** source_channel is the source channel of the IBC packet. */
  source_channel: string;
  /** sequence is the sequence of the IBC packet. */

  sequence: Long;
}
/** QueryPendingSettlementResponse contains one in-flight settlement. */

export interface QueryPendingSettlementResponse {
  /** pending_settlement is the in-flight settlement for the packet. */
  pendingSettlement?: PendingSettlement;
}
/** QueryPendingSettlementResponse contains one in-flight settlement. */

export interface QueryPendingSettlementResponseSDKType {
  /** pending_settlement is the in-flight settlement for the packet. */
  pending_settlement?: PendingSettlementSDKType;
}

function createBaseQueryStateRequest(): QueryStateRequest {
  return {};
}

export const QueryStateRequest = {
  encode(_: QueryStateRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryStateRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryStateRequest();

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

  fromPartial(_: DeepPartial<QueryStateRequest>): QueryStateRequest {
    const message = createBaseQueryStateRequest();
    return message;
  }

};

function createBaseQueryStateResponse(): QueryStateResponse {
  return {
    controls: undefined,
    ledger: undefined
  };
}

export const QueryStateResponse = {
  encode(message: QueryStateResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.controls !== undefined) {
      Controls.encode(message.controls, writer.uint32(10).fork()).ldelim();
    }

    if (message.ledger !== undefined) {
      Ledger.encode(message.ledger, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryStateResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryStateResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.controls = Controls.decode(reader, reader.uint32());
          break;

        case 2:
          message.ledger = Ledger.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryStateResponse>): QueryStateResponse {
    const message = createBaseQueryStateResponse();
    message.controls = object.controls !== undefined && object.controls !== null ? Controls.fromPartial(object.controls) : undefined;
    message.ledger = object.ledger !== undefined && object.ledger !== null ? Ledger.fromPartial(object.ledger) : undefined;
    return message;
  }

};

function createBaseQueryPendingSettlementsRequest(): QueryPendingSettlementsRequest {
  return {
    pagination: undefined
  };
}

export const QueryPendingSettlementsRequest = {
  encode(message: QueryPendingSettlementsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryPendingSettlementsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryPendingSettlementsRequest();

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

  fromPartial(object: DeepPartial<QueryPendingSettlementsRequest>): QueryPendingSettlementsRequest {
    const message = createBaseQueryPendingSettlementsRequest();
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQueryPendingSettlementsResponse(): QueryPendingSettlementsResponse {
  return {
    pendingSettlements: [],
    pagination: undefined
  };
}

export const QueryPendingSettlementsResponse = {
  encode(message: QueryPendingSettlementsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.pendingSettlements) {
      PendingSettlement.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryPendingSettlementsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryPendingSettlementsResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.pendingSettlements.push(PendingSettlement.decode(reader, reader.uint32()));
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

  fromPartial(object: DeepPartial<QueryPendingSettlementsResponse>): QueryPendingSettlementsResponse {
    const message = createBaseQueryPendingSettlementsResponse();
    message.pendingSettlements = object.pendingSettlements?.map(e => PendingSettlement.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQueryPendingSettlementRequest(): QueryPendingSettlementRequest {
  return {
    sourceChannel: "",
    sequence: Long.UZERO
  };
}

export const QueryPendingSettlementRequest = {
  encode(message: QueryPendingSettlementRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.sourceChannel !== "") {
      writer.uint32(10).string(message.sourceChannel);
    }

    if (!message.sequence.isZero()) {
      writer.uint32(16).uint64(message.sequence);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryPendingSettlementRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryPendingSettlementRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.sourceChannel = reader.string();
          break;

        case 2:
          message.sequence = (reader.uint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryPendingSettlementRequest>): QueryPendingSettlementRequest {
    const message = createBaseQueryPendingSettlementRequest();
    message.sourceChannel = object.sourceChannel ?? "";
    message.sequence = object.sequence !== undefined && object.sequence !== null ? Long.fromValue(object.sequence) : Long.UZERO;
    return message;
  }

};

function createBaseQueryPendingSettlementResponse(): QueryPendingSettlementResponse {
  return {
    pendingSettlement: undefined
  };
}

export const QueryPendingSettlementResponse = {
  encode(message: QueryPendingSettlementResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pendingSettlement !== undefined) {
      PendingSettlement.encode(message.pendingSettlement, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryPendingSettlementResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryPendingSettlementResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.pendingSettlement = PendingSettlement.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryPendingSettlementResponse>): QueryPendingSettlementResponse {
    const message = createBaseQueryPendingSettlementResponse();
    message.pendingSettlement = object.pendingSettlement !== undefined && object.pendingSettlement !== null ? PendingSettlement.fromPartial(object.pendingSettlement) : undefined;
    return message;
  }

};