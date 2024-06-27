"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.PublicKey = exports.protobufPackage = void 0;
/* eslint-disable */
const binary_1 = require("../../binary");
const helpers_1 = require("../../helpers");
exports.protobufPackage = "tendermint.crypto";
function createBasePublicKey() {
    return {
        ed25519: undefined,
        secp256k1: undefined,
    };
}
exports.PublicKey = {
    typeUrl: "/tendermint.crypto.PublicKey",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.ed25519 !== undefined) {
            writer.uint32(10).bytes(message.ed25519);
        }
        if (message.secp256k1 !== undefined) {
            writer.uint32(18).bytes(message.secp256k1);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBasePublicKey();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.ed25519 = reader.bytes();
                    break;
                case 2:
                    message.secp256k1 = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBasePublicKey();
        if ((0, helpers_1.isSet)(object.ed25519))
            obj.ed25519 = (0, helpers_1.bytesFromBase64)(object.ed25519);
        if ((0, helpers_1.isSet)(object.secp256k1))
            obj.secp256k1 = (0, helpers_1.bytesFromBase64)(object.secp256k1);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.ed25519 !== undefined &&
            (obj.ed25519 = message.ed25519 !== undefined ? (0, helpers_1.base64FromBytes)(message.ed25519) : undefined);
        message.secp256k1 !== undefined &&
            (obj.secp256k1 = message.secp256k1 !== undefined ? (0, helpers_1.base64FromBytes)(message.secp256k1) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBasePublicKey();
        message.ed25519 = object.ed25519 ?? undefined;
        message.secp256k1 = object.secp256k1 ?? undefined;
        return message;
    },
};
//# sourceMappingURL=keys.js.map