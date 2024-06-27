"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Metadata = exports.protobufPackage = void 0;
/* eslint-disable */
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "ibc.applications.fee.v1";
function createBaseMetadata() {
    return {
        feeVersion: "",
        appVersion: "",
    };
}
exports.Metadata = {
    typeUrl: "/ibc.applications.fee.v1.Metadata",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.feeVersion !== "") {
            writer.uint32(10).string(message.feeVersion);
        }
        if (message.appVersion !== "") {
            writer.uint32(18).string(message.appVersion);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMetadata();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.feeVersion = reader.string();
                    break;
                case 2:
                    message.appVersion = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMetadata();
        if ((0, helpers_1.isSet)(object.feeVersion))
            obj.feeVersion = String(object.feeVersion);
        if ((0, helpers_1.isSet)(object.appVersion))
            obj.appVersion = String(object.appVersion);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.feeVersion !== undefined && (obj.feeVersion = message.feeVersion);
        message.appVersion !== undefined && (obj.appVersion = message.appVersion);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMetadata();
        message.feeVersion = object.feeVersion ?? "";
        message.appVersion = object.appVersion ?? "";
        return message;
    },
};
//# sourceMappingURL=metadata.js.map