"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Metadata = exports.protobufPackage = void 0;
/* eslint-disable */
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "ibc.applications.interchain_accounts.v1";
function createBaseMetadata() {
    return {
        version: "",
        controllerConnectionId: "",
        hostConnectionId: "",
        address: "",
        encoding: "",
        txType: "",
    };
}
exports.Metadata = {
    typeUrl: "/ibc.applications.interchain_accounts.v1.Metadata",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.version !== "") {
            writer.uint32(10).string(message.version);
        }
        if (message.controllerConnectionId !== "") {
            writer.uint32(18).string(message.controllerConnectionId);
        }
        if (message.hostConnectionId !== "") {
            writer.uint32(26).string(message.hostConnectionId);
        }
        if (message.address !== "") {
            writer.uint32(34).string(message.address);
        }
        if (message.encoding !== "") {
            writer.uint32(42).string(message.encoding);
        }
        if (message.txType !== "") {
            writer.uint32(50).string(message.txType);
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
                    message.version = reader.string();
                    break;
                case 2:
                    message.controllerConnectionId = reader.string();
                    break;
                case 3:
                    message.hostConnectionId = reader.string();
                    break;
                case 4:
                    message.address = reader.string();
                    break;
                case 5:
                    message.encoding = reader.string();
                    break;
                case 6:
                    message.txType = reader.string();
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
        if ((0, helpers_1.isSet)(object.version))
            obj.version = String(object.version);
        if ((0, helpers_1.isSet)(object.controllerConnectionId))
            obj.controllerConnectionId = String(object.controllerConnectionId);
        if ((0, helpers_1.isSet)(object.hostConnectionId))
            obj.hostConnectionId = String(object.hostConnectionId);
        if ((0, helpers_1.isSet)(object.address))
            obj.address = String(object.address);
        if ((0, helpers_1.isSet)(object.encoding))
            obj.encoding = String(object.encoding);
        if ((0, helpers_1.isSet)(object.txType))
            obj.txType = String(object.txType);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.version !== undefined && (obj.version = message.version);
        message.controllerConnectionId !== undefined &&
            (obj.controllerConnectionId = message.controllerConnectionId);
        message.hostConnectionId !== undefined && (obj.hostConnectionId = message.hostConnectionId);
        message.address !== undefined && (obj.address = message.address);
        message.encoding !== undefined && (obj.encoding = message.encoding);
        message.txType !== undefined && (obj.txType = message.txType);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMetadata();
        message.version = object.version ?? "";
        message.controllerConnectionId = object.controllerConnectionId ?? "";
        message.hostConnectionId = object.hostConnectionId ?? "";
        message.address = object.address ?? "";
        message.encoding = object.encoding ?? "";
        message.txType = object.txType ?? "";
        return message;
    },
};
//# sourceMappingURL=metadata.js.map