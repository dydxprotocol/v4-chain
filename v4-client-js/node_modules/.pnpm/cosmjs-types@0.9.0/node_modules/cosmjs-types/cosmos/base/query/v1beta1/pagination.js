"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.PageResponse = exports.PageRequest = exports.protobufPackage = void 0;
/* eslint-disable */
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "cosmos.base.query.v1beta1";
function createBasePageRequest() {
    return {
        key: new Uint8Array(),
        offset: BigInt(0),
        limit: BigInt(0),
        countTotal: false,
        reverse: false,
    };
}
exports.PageRequest = {
    typeUrl: "/cosmos.base.query.v1beta1.PageRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.key.length !== 0) {
            writer.uint32(10).bytes(message.key);
        }
        if (message.offset !== BigInt(0)) {
            writer.uint32(16).uint64(message.offset);
        }
        if (message.limit !== BigInt(0)) {
            writer.uint32(24).uint64(message.limit);
        }
        if (message.countTotal === true) {
            writer.uint32(32).bool(message.countTotal);
        }
        if (message.reverse === true) {
            writer.uint32(40).bool(message.reverse);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBasePageRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.key = reader.bytes();
                    break;
                case 2:
                    message.offset = reader.uint64();
                    break;
                case 3:
                    message.limit = reader.uint64();
                    break;
                case 4:
                    message.countTotal = reader.bool();
                    break;
                case 5:
                    message.reverse = reader.bool();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBasePageRequest();
        if ((0, helpers_1.isSet)(object.key))
            obj.key = (0, helpers_1.bytesFromBase64)(object.key);
        if ((0, helpers_1.isSet)(object.offset))
            obj.offset = BigInt(object.offset.toString());
        if ((0, helpers_1.isSet)(object.limit))
            obj.limit = BigInt(object.limit.toString());
        if ((0, helpers_1.isSet)(object.countTotal))
            obj.countTotal = Boolean(object.countTotal);
        if ((0, helpers_1.isSet)(object.reverse))
            obj.reverse = Boolean(object.reverse);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.key !== undefined &&
            (obj.key = (0, helpers_1.base64FromBytes)(message.key !== undefined ? message.key : new Uint8Array()));
        message.offset !== undefined && (obj.offset = (message.offset || BigInt(0)).toString());
        message.limit !== undefined && (obj.limit = (message.limit || BigInt(0)).toString());
        message.countTotal !== undefined && (obj.countTotal = message.countTotal);
        message.reverse !== undefined && (obj.reverse = message.reverse);
        return obj;
    },
    fromPartial(object) {
        const message = createBasePageRequest();
        message.key = object.key ?? new Uint8Array();
        if (object.offset !== undefined && object.offset !== null) {
            message.offset = BigInt(object.offset.toString());
        }
        if (object.limit !== undefined && object.limit !== null) {
            message.limit = BigInt(object.limit.toString());
        }
        message.countTotal = object.countTotal ?? false;
        message.reverse = object.reverse ?? false;
        return message;
    },
};
function createBasePageResponse() {
    return {
        nextKey: new Uint8Array(),
        total: BigInt(0),
    };
}
exports.PageResponse = {
    typeUrl: "/cosmos.base.query.v1beta1.PageResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.nextKey.length !== 0) {
            writer.uint32(10).bytes(message.nextKey);
        }
        if (message.total !== BigInt(0)) {
            writer.uint32(16).uint64(message.total);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBasePageResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.nextKey = reader.bytes();
                    break;
                case 2:
                    message.total = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBasePageResponse();
        if ((0, helpers_1.isSet)(object.nextKey))
            obj.nextKey = (0, helpers_1.bytesFromBase64)(object.nextKey);
        if ((0, helpers_1.isSet)(object.total))
            obj.total = BigInt(object.total.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.nextKey !== undefined &&
            (obj.nextKey = (0, helpers_1.base64FromBytes)(message.nextKey !== undefined ? message.nextKey : new Uint8Array()));
        message.total !== undefined && (obj.total = (message.total || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBasePageResponse();
        message.nextKey = object.nextKey ?? new Uint8Array();
        if (object.total !== undefined && object.total !== null) {
            message.total = BigInt(object.total.toString());
        }
        return message;
    },
};
//# sourceMappingURL=pagination.js.map