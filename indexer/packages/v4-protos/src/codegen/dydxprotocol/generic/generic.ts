import * as _m0 from "protobufjs/minimal";
import { DeepPartial, Long } from "../../helpers";
/** Uint32 is a proto to hold a single uint32 value. */

export interface Uint32 {
  /** Uint32 is a proto to hold a single uint32 value. */
  value: number;
}
/** Uint32 is a proto to hold a single uint32 value. */

export interface Uint32SDKType {
  /** Uint32 is a proto to hold a single uint32 value. */
  value: number;
}
/** Uint64 is a proto to hold a single uint64 value. */

export interface Uint64 {
  /** Uint64 is a proto to hold a single uint64 value. */
  value: Long;
}
/** Uint64 is a proto to hold a single uint64 value. */

export interface Uint64SDKType {
  /** Uint64 is a proto to hold a single uint64 value. */
  value: Long;
}
/** Int32 is a proto to hold a single sint32 value. */

export interface Int32 {
  /** Int32 is a proto to hold a single sint32 value. */
  value: number;
}
/** Int32 is a proto to hold a single sint32 value. */

export interface Int32SDKType {
  /** Int32 is a proto to hold a single sint32 value. */
  value: number;
}
/** Int64 is a proto to hold a single sint64 value. */

export interface Int64 {
  /** Int64 is a proto to hold a single sint64 value. */
  value: Long;
}
/** Int64 is a proto to hold a single sint64 value. */

export interface Int64SDKType {
  /** Int64 is a proto to hold a single sint64 value. */
  value: Long;
}

function createBaseUint32(): Uint32 {
  return {
    value: 0
  };
}

export const Uint32 = {
  encode(message: Uint32, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.value !== 0) {
      writer.uint32(8).uint32(message.value);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Uint32 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUint32();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.value = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<Uint32>): Uint32 {
    const message = createBaseUint32();
    message.value = object.value ?? 0;
    return message;
  }

};

function createBaseUint64(): Uint64 {
  return {
    value: Long.UZERO
  };
}

export const Uint64 = {
  encode(message: Uint64, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (!message.value.isZero()) {
      writer.uint32(8).uint64(message.value);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Uint64 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUint64();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.value = (reader.uint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<Uint64>): Uint64 {
    const message = createBaseUint64();
    message.value = object.value !== undefined && object.value !== null ? Long.fromValue(object.value) : Long.UZERO;
    return message;
  }

};

function createBaseInt32(): Int32 {
  return {
    value: 0
  };
}

export const Int32 = {
  encode(message: Int32, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.value !== 0) {
      writer.uint32(8).sint32(message.value);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Int32 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInt32();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.value = reader.sint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<Int32>): Int32 {
    const message = createBaseInt32();
    message.value = object.value ?? 0;
    return message;
  }

};

function createBaseInt64(): Int64 {
  return {
    value: Long.ZERO
  };
}

export const Int64 = {
  encode(message: Int64, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (!message.value.isZero()) {
      writer.uint32(8).sint64(message.value);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Int64 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInt64();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.value = (reader.sint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<Int64>): Int64 {
    const message = createBaseInt64();
    message.value = object.value !== undefined && object.value !== null ? Long.fromValue(object.value) : Long.ZERO;
    return message;
  }

};