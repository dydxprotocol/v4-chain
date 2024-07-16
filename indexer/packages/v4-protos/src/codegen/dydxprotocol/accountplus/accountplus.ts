import * as _m0 from "protobufjs/minimal";
import { DeepPartial, Long } from "../../helpers";
/** Account State */

export interface AccountState {
  address: string;
  timestampNonceDetails?: TimestampNonceDetails;
}
/** Account State */

export interface AccountStateSDKType {
  address: string;
  timestamp_nonce_details?: TimestampNonceDetailsSDKType;
}
/** Timestamp nonce details */

export interface TimestampNonceDetails {
  /** unsorted list of n most recent timestamp nonces */
  timestampNonces: Long[];
  /** max timestamp nonce that was ejected from list above */

  maxEjectedNonce: Long;
}
/** Timestamp nonce details */

export interface TimestampNonceDetailsSDKType {
  /** unsorted list of n most recent timestamp nonces */
  timestamp_nonces: Long[];
  /** max timestamp nonce that was ejected from list above */

  max_ejected_nonce: Long;
}

function createBaseAccountState(): AccountState {
  return {
    address: "",
    timestampNonceDetails: undefined
  };
}

export const AccountState = {
  encode(message: AccountState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }

    if (message.timestampNonceDetails !== undefined) {
      TimestampNonceDetails.encode(message.timestampNonceDetails, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AccountState {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAccountState();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;

        case 2:
          message.timestampNonceDetails = TimestampNonceDetails.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<AccountState>): AccountState {
    const message = createBaseAccountState();
    message.address = object.address ?? "";
    message.timestampNonceDetails = object.timestampNonceDetails !== undefined && object.timestampNonceDetails !== null ? TimestampNonceDetails.fromPartial(object.timestampNonceDetails) : undefined;
    return message;
  }

};

function createBaseTimestampNonceDetails(): TimestampNonceDetails {
  return {
    timestampNonces: [],
    maxEjectedNonce: Long.UZERO
  };
}

export const TimestampNonceDetails = {
  encode(message: TimestampNonceDetails, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    writer.uint32(10).fork();

    for (const v of message.timestampNonces) {
      writer.uint64(v);
    }

    writer.ldelim();

    if (!message.maxEjectedNonce.isZero()) {
      writer.uint32(16).uint64(message.maxEjectedNonce);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TimestampNonceDetails {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTimestampNonceDetails();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;

            while (reader.pos < end2) {
              message.timestampNonces.push((reader.uint64() as Long));
            }
          } else {
            message.timestampNonces.push((reader.uint64() as Long));
          }

          break;

        case 2:
          message.maxEjectedNonce = (reader.uint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<TimestampNonceDetails>): TimestampNonceDetails {
    const message = createBaseTimestampNonceDetails();
    message.timestampNonces = object.timestampNonces?.map(e => Long.fromValue(e)) || [];
    message.maxEjectedNonce = object.maxEjectedNonce !== undefined && object.maxEjectedNonce !== null ? Long.fromValue(object.maxEjectedNonce) : Long.UZERO;
    return message;
  }

};