"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Config = exports.protobufPackage = void 0;
/* eslint-disable */
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "cosmos.tx.config.v1";
function createBaseConfig() {
    return {
        skipAnteHandler: false,
        skipPostHandler: false,
    };
}
exports.Config = {
    typeUrl: "/cosmos.tx.config.v1.Config",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.skipAnteHandler === true) {
            writer.uint32(8).bool(message.skipAnteHandler);
        }
        if (message.skipPostHandler === true) {
            writer.uint32(16).bool(message.skipPostHandler);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseConfig();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.skipAnteHandler = reader.bool();
                    break;
                case 2:
                    message.skipPostHandler = reader.bool();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseConfig();
        if ((0, helpers_1.isSet)(object.skipAnteHandler))
            obj.skipAnteHandler = Boolean(object.skipAnteHandler);
        if ((0, helpers_1.isSet)(object.skipPostHandler))
            obj.skipPostHandler = Boolean(object.skipPostHandler);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.skipAnteHandler !== undefined && (obj.skipAnteHandler = message.skipAnteHandler);
        message.skipPostHandler !== undefined && (obj.skipPostHandler = message.skipPostHandler);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseConfig();
        message.skipAnteHandler = object.skipAnteHandler ?? false;
        message.skipPostHandler = object.skipPostHandler ?? false;
        return message;
    },
};
//# sourceMappingURL=config.js.map