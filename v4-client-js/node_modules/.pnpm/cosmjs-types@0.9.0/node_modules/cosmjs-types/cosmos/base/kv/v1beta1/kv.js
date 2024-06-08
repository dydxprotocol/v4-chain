"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Pair = exports.Pairs = exports.protobufPackage = void 0;
/* eslint-disable */
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "cosmos.base.kv.v1beta1";
function createBasePairs() {
    return {
        pairs: [],
    };
}
exports.Pairs = {
    typeUrl: "/cosmos.base.kv.v1beta1.Pairs",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.pairs) {
            exports.Pair.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBasePairs();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.pairs.push(exports.Pair.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBasePairs();
        if (Array.isArray(object?.pairs))
            obj.pairs = object.pairs.map((e) => exports.Pair.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.pairs) {
            obj.pairs = message.pairs.map((e) => (e ? exports.Pair.toJSON(e) : undefined));
        }
        else {
            obj.pairs = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBasePairs();
        message.pairs = object.pairs?.map((e) => exports.Pair.fromPartial(e)) || [];
        return message;
    },
};
function createBasePair() {
    return {
        key: new Uint8Array(),
        value: new Uint8Array(),
    };
}
exports.Pair = {
    typeUrl: "/cosmos.base.kv.v1beta1.Pair",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.key.length !== 0) {
            writer.uint32(10).bytes(message.key);
        }
        if (message.value.length !== 0) {
            writer.uint32(18).bytes(message.value);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBasePair();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.key = reader.bytes();
                    break;
                case 2:
                    message.value = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBasePair();
        if ((0, helpers_1.isSet)(object.key))
            obj.key = (0, helpers_1.bytesFromBase64)(object.key);
        if ((0, helpers_1.isSet)(object.value))
            obj.value = (0, helpers_1.bytesFromBase64)(object.value);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.key !== undefined &&
            (obj.key = (0, helpers_1.base64FromBytes)(message.key !== undefined ? message.key : new Uint8Array()));
        message.value !== undefined &&
            (obj.value = (0, helpers_1.base64FromBytes)(message.value !== undefined ? message.value : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBasePair();
        message.key = object.key ?? new Uint8Array();
        message.value = object.value ?? new Uint8Array();
        return message;
    },
};
//# sourceMappingURL=kv.js.map