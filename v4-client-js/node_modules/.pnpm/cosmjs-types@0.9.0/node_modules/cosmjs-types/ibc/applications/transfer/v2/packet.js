"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.FungibleTokenPacketData = exports.protobufPackage = void 0;
/* eslint-disable */
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "ibc.applications.transfer.v2";
function createBaseFungibleTokenPacketData() {
    return {
        denom: "",
        amount: "",
        sender: "",
        receiver: "",
        memo: "",
    };
}
exports.FungibleTokenPacketData = {
    typeUrl: "/ibc.applications.transfer.v2.FungibleTokenPacketData",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.denom !== "") {
            writer.uint32(10).string(message.denom);
        }
        if (message.amount !== "") {
            writer.uint32(18).string(message.amount);
        }
        if (message.sender !== "") {
            writer.uint32(26).string(message.sender);
        }
        if (message.receiver !== "") {
            writer.uint32(34).string(message.receiver);
        }
        if (message.memo !== "") {
            writer.uint32(42).string(message.memo);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseFungibleTokenPacketData();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.denom = reader.string();
                    break;
                case 2:
                    message.amount = reader.string();
                    break;
                case 3:
                    message.sender = reader.string();
                    break;
                case 4:
                    message.receiver = reader.string();
                    break;
                case 5:
                    message.memo = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseFungibleTokenPacketData();
        if ((0, helpers_1.isSet)(object.denom))
            obj.denom = String(object.denom);
        if ((0, helpers_1.isSet)(object.amount))
            obj.amount = String(object.amount);
        if ((0, helpers_1.isSet)(object.sender))
            obj.sender = String(object.sender);
        if ((0, helpers_1.isSet)(object.receiver))
            obj.receiver = String(object.receiver);
        if ((0, helpers_1.isSet)(object.memo))
            obj.memo = String(object.memo);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.denom !== undefined && (obj.denom = message.denom);
        message.amount !== undefined && (obj.amount = message.amount);
        message.sender !== undefined && (obj.sender = message.sender);
        message.receiver !== undefined && (obj.receiver = message.receiver);
        message.memo !== undefined && (obj.memo = message.memo);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseFungibleTokenPacketData();
        message.denom = object.denom ?? "";
        message.amount = object.amount ?? "";
        message.sender = object.sender ?? "";
        message.receiver = object.receiver ?? "";
        message.memo = object.memo ?? "";
        return message;
    },
};
//# sourceMappingURL=packet.js.map