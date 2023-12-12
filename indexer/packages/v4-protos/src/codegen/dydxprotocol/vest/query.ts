import { VestEntry, VestEntryAmino, VestEntrySDKType } from "./vest_entry";
import { BinaryReader, BinaryWriter } from "../../binary";
/** QueryVestEntryRequest is a request type for the VestEntry RPC method. */
export interface QueryVestEntryRequest {
  /** QueryVestEntryRequest is a request type for the VestEntry RPC method. */
  vesterAccount: string;
}
export interface QueryVestEntryRequestProtoMsg {
  typeUrl: "/dydxprotocol.vest.QueryVestEntryRequest";
  value: Uint8Array;
}
/** QueryVestEntryRequest is a request type for the VestEntry RPC method. */
export interface QueryVestEntryRequestAmino {
  /** QueryVestEntryRequest is a request type for the VestEntry RPC method. */
  vester_account?: string;
}
export interface QueryVestEntryRequestAminoMsg {
  type: "/dydxprotocol.vest.QueryVestEntryRequest";
  value: QueryVestEntryRequestAmino;
}
/** QueryVestEntryRequest is a request type for the VestEntry RPC method. */
export interface QueryVestEntryRequestSDKType {
  vester_account: string;
}
/** QueryVestEntryResponse is a response type for the VestEntry RPC method. */
export interface QueryVestEntryResponse {
  entry: VestEntry;
}
export interface QueryVestEntryResponseProtoMsg {
  typeUrl: "/dydxprotocol.vest.QueryVestEntryResponse";
  value: Uint8Array;
}
/** QueryVestEntryResponse is a response type for the VestEntry RPC method. */
export interface QueryVestEntryResponseAmino {
  entry?: VestEntryAmino;
}
export interface QueryVestEntryResponseAminoMsg {
  type: "/dydxprotocol.vest.QueryVestEntryResponse";
  value: QueryVestEntryResponseAmino;
}
/** QueryVestEntryResponse is a response type for the VestEntry RPC method. */
export interface QueryVestEntryResponseSDKType {
  entry: VestEntrySDKType;
}
function createBaseQueryVestEntryRequest(): QueryVestEntryRequest {
  return {
    vesterAccount: ""
  };
}
export const QueryVestEntryRequest = {
  typeUrl: "/dydxprotocol.vest.QueryVestEntryRequest",
  encode(message: QueryVestEntryRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.vesterAccount !== "") {
      writer.uint32(10).string(message.vesterAccount);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryVestEntryRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryVestEntryRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.vesterAccount = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryVestEntryRequest>): QueryVestEntryRequest {
    const message = createBaseQueryVestEntryRequest();
    message.vesterAccount = object.vesterAccount ?? "";
    return message;
  },
  fromAmino(object: QueryVestEntryRequestAmino): QueryVestEntryRequest {
    const message = createBaseQueryVestEntryRequest();
    if (object.vester_account !== undefined && object.vester_account !== null) {
      message.vesterAccount = object.vester_account;
    }
    return message;
  },
  toAmino(message: QueryVestEntryRequest): QueryVestEntryRequestAmino {
    const obj: any = {};
    obj.vester_account = message.vesterAccount;
    return obj;
  },
  fromAminoMsg(object: QueryVestEntryRequestAminoMsg): QueryVestEntryRequest {
    return QueryVestEntryRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryVestEntryRequestProtoMsg): QueryVestEntryRequest {
    return QueryVestEntryRequest.decode(message.value);
  },
  toProto(message: QueryVestEntryRequest): Uint8Array {
    return QueryVestEntryRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryVestEntryRequest): QueryVestEntryRequestProtoMsg {
    return {
      typeUrl: "/dydxprotocol.vest.QueryVestEntryRequest",
      value: QueryVestEntryRequest.encode(message).finish()
    };
  }
};
function createBaseQueryVestEntryResponse(): QueryVestEntryResponse {
  return {
    entry: VestEntry.fromPartial({})
  };
}
export const QueryVestEntryResponse = {
  typeUrl: "/dydxprotocol.vest.QueryVestEntryResponse",
  encode(message: QueryVestEntryResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.entry !== undefined) {
      VestEntry.encode(message.entry, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryVestEntryResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryVestEntryResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.entry = VestEntry.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryVestEntryResponse>): QueryVestEntryResponse {
    const message = createBaseQueryVestEntryResponse();
    message.entry = object.entry !== undefined && object.entry !== null ? VestEntry.fromPartial(object.entry) : undefined;
    return message;
  },
  fromAmino(object: QueryVestEntryResponseAmino): QueryVestEntryResponse {
    const message = createBaseQueryVestEntryResponse();
    if (object.entry !== undefined && object.entry !== null) {
      message.entry = VestEntry.fromAmino(object.entry);
    }
    return message;
  },
  toAmino(message: QueryVestEntryResponse): QueryVestEntryResponseAmino {
    const obj: any = {};
    obj.entry = message.entry ? VestEntry.toAmino(message.entry) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryVestEntryResponseAminoMsg): QueryVestEntryResponse {
    return QueryVestEntryResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryVestEntryResponseProtoMsg): QueryVestEntryResponse {
    return QueryVestEntryResponse.decode(message.value);
  },
  toProto(message: QueryVestEntryResponse): Uint8Array {
    return QueryVestEntryResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryVestEntryResponse): QueryVestEntryResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.vest.QueryVestEntryResponse",
      value: QueryVestEntryResponse.encode(message).finish()
    };
  }
};