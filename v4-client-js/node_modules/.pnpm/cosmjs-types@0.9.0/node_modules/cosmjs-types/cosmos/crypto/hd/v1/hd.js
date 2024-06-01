"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BIP44Params = exports.protobufPackage = void 0;
/* eslint-disable */
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "cosmos.crypto.hd.v1";
function createBaseBIP44Params() {
    return {
        purpose: 0,
        coinType: 0,
        account: 0,
        change: false,
        addressIndex: 0,
    };
}
exports.BIP44Params = {
    typeUrl: "/cosmos.crypto.hd.v1.BIP44Params",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.purpose !== 0) {
            writer.uint32(8).uint32(message.purpose);
        }
        if (message.coinType !== 0) {
            writer.uint32(16).uint32(message.coinType);
        }
        if (message.account !== 0) {
            writer.uint32(24).uint32(message.account);
        }
        if (message.change === true) {
            writer.uint32(32).bool(message.change);
        }
        if (message.addressIndex !== 0) {
            writer.uint32(40).uint32(message.addressIndex);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseBIP44Params();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.purpose = reader.uint32();
                    break;
                case 2:
                    message.coinType = reader.uint32();
                    break;
                case 3:
                    message.account = reader.uint32();
                    break;
                case 4:
                    message.change = reader.bool();
                    break;
                case 5:
                    message.addressIndex = reader.uint32();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseBIP44Params();
        if ((0, helpers_1.isSet)(object.purpose))
            obj.purpose = Number(object.purpose);
        if ((0, helpers_1.isSet)(object.coinType))
            obj.coinType = Number(object.coinType);
        if ((0, helpers_1.isSet)(object.account))
            obj.account = Number(object.account);
        if ((0, helpers_1.isSet)(object.change))
            obj.change = Boolean(object.change);
        if ((0, helpers_1.isSet)(object.addressIndex))
            obj.addressIndex = Number(object.addressIndex);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.purpose !== undefined && (obj.purpose = Math.round(message.purpose));
        message.coinType !== undefined && (obj.coinType = Math.round(message.coinType));
        message.account !== undefined && (obj.account = Math.round(message.account));
        message.change !== undefined && (obj.change = message.change);
        message.addressIndex !== undefined && (obj.addressIndex = Math.round(message.addressIndex));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseBIP44Params();
        message.purpose = object.purpose ?? 0;
        message.coinType = object.coinType ?? 0;
        message.account = object.account ?? 0;
        message.change = object.change ?? false;
        message.addressIndex = object.addressIndex ?? 0;
        return message;
    },
};
//# sourceMappingURL=hd.js.map