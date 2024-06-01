"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Params = exports.protobufPackage = void 0;
/* eslint-disable */
const binary_1 = require("../../../../../binary");
const helpers_1 = require("../../../../../helpers");
exports.protobufPackage = "ibc.applications.interchain_accounts.host.v1";
function createBaseParams() {
    return {
        hostEnabled: false,
        allowMessages: [],
    };
}
exports.Params = {
    typeUrl: "/ibc.applications.interchain_accounts.host.v1.Params",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.hostEnabled === true) {
            writer.uint32(8).bool(message.hostEnabled);
        }
        for (const v of message.allowMessages) {
            writer.uint32(18).string(v);
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
                    message.hostEnabled = reader.bool();
                    break;
                case 2:
                    message.allowMessages.push(reader.string());
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
        if ((0, helpers_1.isSet)(object.hostEnabled))
            obj.hostEnabled = Boolean(object.hostEnabled);
        if (Array.isArray(object?.allowMessages))
            obj.allowMessages = object.allowMessages.map((e) => String(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.hostEnabled !== undefined && (obj.hostEnabled = message.hostEnabled);
        if (message.allowMessages) {
            obj.allowMessages = message.allowMessages.map((e) => e);
        }
        else {
            obj.allowMessages = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseParams();
        message.hostEnabled = object.hostEnabled ?? false;
        message.allowMessages = object.allowMessages?.map((e) => e) || [];
        return message;
    },
};
//# sourceMappingURL=host.js.map