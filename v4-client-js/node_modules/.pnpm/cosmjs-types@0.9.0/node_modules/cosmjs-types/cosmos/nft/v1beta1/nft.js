"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.NFT = exports.Class = exports.protobufPackage = void 0;
/* eslint-disable */
const any_1 = require("../../../google/protobuf/any");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.nft.v1beta1";
function createBaseClass() {
    return {
        id: "",
        name: "",
        symbol: "",
        description: "",
        uri: "",
        uriHash: "",
        data: undefined,
    };
}
exports.Class = {
    typeUrl: "/cosmos.nft.v1beta1.Class",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.id !== "") {
            writer.uint32(10).string(message.id);
        }
        if (message.name !== "") {
            writer.uint32(18).string(message.name);
        }
        if (message.symbol !== "") {
            writer.uint32(26).string(message.symbol);
        }
        if (message.description !== "") {
            writer.uint32(34).string(message.description);
        }
        if (message.uri !== "") {
            writer.uint32(42).string(message.uri);
        }
        if (message.uriHash !== "") {
            writer.uint32(50).string(message.uriHash);
        }
        if (message.data !== undefined) {
            any_1.Any.encode(message.data, writer.uint32(58).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseClass();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.id = reader.string();
                    break;
                case 2:
                    message.name = reader.string();
                    break;
                case 3:
                    message.symbol = reader.string();
                    break;
                case 4:
                    message.description = reader.string();
                    break;
                case 5:
                    message.uri = reader.string();
                    break;
                case 6:
                    message.uriHash = reader.string();
                    break;
                case 7:
                    message.data = any_1.Any.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseClass();
        if ((0, helpers_1.isSet)(object.id))
            obj.id = String(object.id);
        if ((0, helpers_1.isSet)(object.name))
            obj.name = String(object.name);
        if ((0, helpers_1.isSet)(object.symbol))
            obj.symbol = String(object.symbol);
        if ((0, helpers_1.isSet)(object.description))
            obj.description = String(object.description);
        if ((0, helpers_1.isSet)(object.uri))
            obj.uri = String(object.uri);
        if ((0, helpers_1.isSet)(object.uriHash))
            obj.uriHash = String(object.uriHash);
        if ((0, helpers_1.isSet)(object.data))
            obj.data = any_1.Any.fromJSON(object.data);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.id !== undefined && (obj.id = message.id);
        message.name !== undefined && (obj.name = message.name);
        message.symbol !== undefined && (obj.symbol = message.symbol);
        message.description !== undefined && (obj.description = message.description);
        message.uri !== undefined && (obj.uri = message.uri);
        message.uriHash !== undefined && (obj.uriHash = message.uriHash);
        message.data !== undefined && (obj.data = message.data ? any_1.Any.toJSON(message.data) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseClass();
        message.id = object.id ?? "";
        message.name = object.name ?? "";
        message.symbol = object.symbol ?? "";
        message.description = object.description ?? "";
        message.uri = object.uri ?? "";
        message.uriHash = object.uriHash ?? "";
        if (object.data !== undefined && object.data !== null) {
            message.data = any_1.Any.fromPartial(object.data);
        }
        return message;
    },
};
function createBaseNFT() {
    return {
        classId: "",
        id: "",
        uri: "",
        uriHash: "",
        data: undefined,
    };
}
exports.NFT = {
    typeUrl: "/cosmos.nft.v1beta1.NFT",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.classId !== "") {
            writer.uint32(10).string(message.classId);
        }
        if (message.id !== "") {
            writer.uint32(18).string(message.id);
        }
        if (message.uri !== "") {
            writer.uint32(26).string(message.uri);
        }
        if (message.uriHash !== "") {
            writer.uint32(34).string(message.uriHash);
        }
        if (message.data !== undefined) {
            any_1.Any.encode(message.data, writer.uint32(82).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseNFT();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.classId = reader.string();
                    break;
                case 2:
                    message.id = reader.string();
                    break;
                case 3:
                    message.uri = reader.string();
                    break;
                case 4:
                    message.uriHash = reader.string();
                    break;
                case 10:
                    message.data = any_1.Any.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseNFT();
        if ((0, helpers_1.isSet)(object.classId))
            obj.classId = String(object.classId);
        if ((0, helpers_1.isSet)(object.id))
            obj.id = String(object.id);
        if ((0, helpers_1.isSet)(object.uri))
            obj.uri = String(object.uri);
        if ((0, helpers_1.isSet)(object.uriHash))
            obj.uriHash = String(object.uriHash);
        if ((0, helpers_1.isSet)(object.data))
            obj.data = any_1.Any.fromJSON(object.data);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.classId !== undefined && (obj.classId = message.classId);
        message.id !== undefined && (obj.id = message.id);
        message.uri !== undefined && (obj.uri = message.uri);
        message.uriHash !== undefined && (obj.uriHash = message.uriHash);
        message.data !== undefined && (obj.data = message.data ? any_1.Any.toJSON(message.data) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseNFT();
        message.classId = object.classId ?? "";
        message.id = object.id ?? "";
        message.uri = object.uri ?? "";
        message.uriHash = object.uriHash ?? "";
        if (object.data !== undefined && object.data !== null) {
            message.data = any_1.Any.fromPartial(object.data);
        }
        return message;
    },
};
//# sourceMappingURL=nft.js.map