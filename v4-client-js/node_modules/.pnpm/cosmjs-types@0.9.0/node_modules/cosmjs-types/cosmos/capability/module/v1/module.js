"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Module = exports.protobufPackage = void 0;
/* eslint-disable */
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "cosmos.capability.module.v1";
function createBaseModule() {
    return {
        sealKeeper: false,
    };
}
exports.Module = {
    typeUrl: "/cosmos.capability.module.v1.Module",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.sealKeeper === true) {
            writer.uint32(8).bool(message.sealKeeper);
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
                    message.sealKeeper = reader.bool();
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
        if ((0, helpers_1.isSet)(object.sealKeeper))
            obj.sealKeeper = Boolean(object.sealKeeper);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.sealKeeper !== undefined && (obj.sealKeeper = message.sealKeeper);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseModule();
        message.sealKeeper = object.sealKeeper ?? false;
        return message;
    },
};
//# sourceMappingURL=module.js.map