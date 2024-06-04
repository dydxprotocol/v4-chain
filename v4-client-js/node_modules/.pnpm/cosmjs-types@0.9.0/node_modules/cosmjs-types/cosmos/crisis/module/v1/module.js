"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Module = exports.protobufPackage = void 0;
/* eslint-disable */
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "cosmos.crisis.module.v1";
function createBaseModule() {
    return {
        feeCollectorName: "",
        authority: "",
    };
}
exports.Module = {
    typeUrl: "/cosmos.crisis.module.v1.Module",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.feeCollectorName !== "") {
            writer.uint32(10).string(message.feeCollectorName);
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
                    message.feeCollectorName = reader.string();
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
        if ((0, helpers_1.isSet)(object.feeCollectorName))
            obj.feeCollectorName = String(object.feeCollectorName);
        if ((0, helpers_1.isSet)(object.authority))
            obj.authority = String(object.authority);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.feeCollectorName !== undefined && (obj.feeCollectorName = message.feeCollectorName);
        message.authority !== undefined && (obj.authority = message.authority);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseModule();
        message.feeCollectorName = object.feeCollectorName ?? "";
        message.authority = object.authority ?? "";
        return message;
    },
};
//# sourceMappingURL=module.js.map