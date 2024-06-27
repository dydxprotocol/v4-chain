"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Module = exports.protobufPackage = void 0;
/* eslint-disable */
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "cosmos.gov.module.v1";
function createBaseModule() {
    return {
        maxMetadataLen: BigInt(0),
        authority: "",
    };
}
exports.Module = {
    typeUrl: "/cosmos.gov.module.v1.Module",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.maxMetadataLen !== BigInt(0)) {
            writer.uint32(8).uint64(message.maxMetadataLen);
        }
        if (message.authority !== "") {
            writer.uint32(18).string(message.authority);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseModule();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.maxMetadataLen = reader.uint64();
                    break;
                case 2:
                    message.authority = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseModule();
        if ((0, helpers_1.isSet)(object.maxMetadataLen))
            obj.maxMetadataLen = BigInt(object.maxMetadataLen.toString());
        if ((0, helpers_1.isSet)(object.authority))
            obj.authority = String(object.authority);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.maxMetadataLen !== undefined &&
            (obj.maxMetadataLen = (message.maxMetadataLen || BigInt(0)).toString());
        message.authority !== undefined && (obj.authority = message.authority);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseModule();
        if (object.maxMetadataLen !== undefined && object.maxMetadataLen !== null) {
            message.maxMetadataLen = BigInt(object.maxMetadataLen.toString());
        }
        message.authority = object.authority ?? "";
        return message;
    },
};
//# sourceMappingURL=module.js.map