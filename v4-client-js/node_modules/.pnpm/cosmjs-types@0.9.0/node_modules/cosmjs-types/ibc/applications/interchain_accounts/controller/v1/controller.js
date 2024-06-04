"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Params = exports.protobufPackage = void 0;
/* eslint-disable */
const binary_1 = require("../../../../../binary");
const helpers_1 = require("../../../../../helpers");
exports.protobufPackage = "ibc.applications.interchain_accounts.controller.v1";
function createBaseParams() {
    return {
        controllerEnabled: false,
    };
}
exports.Params = {
    typeUrl: "/ibc.applications.interchain_accounts.controller.v1.Params",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.controllerEnabled === true) {
            writer.uint32(8).bool(message.controllerEnabled);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseParams();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.controllerEnabled = reader.bool();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseParams();
        if ((0, helpers_1.isSet)(object.controllerEnabled))
            obj.controllerEnabled = Boolean(object.controllerEnabled);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.controllerEnabled !== undefined && (obj.controllerEnabled = message.controllerEnabled);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseParams();
        message.controllerEnabled = object.controllerEnabled ?? false;
        return message;
    },
};
//# sourceMappingURL=controller.js.map