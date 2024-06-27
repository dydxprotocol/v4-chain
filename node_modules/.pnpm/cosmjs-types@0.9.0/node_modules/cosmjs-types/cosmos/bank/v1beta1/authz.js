"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SendAuthorization = exports.protobufPackage = void 0;
/* eslint-disable */
const coin_1 = require("../../base/v1beta1/coin");
const binary_1 = require("../../../binary");
exports.protobufPackage = "cosmos.bank.v1beta1";
function createBaseSendAuthorization() {
    return {
        spendLimit: [],
        allowList: [],
    };
}
exports.SendAuthorization = {
    typeUrl: "/cosmos.bank.v1beta1.SendAuthorization",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.spendLimit) {
            coin_1.Coin.encode(v, writer.uint32(10).fork()).ldelim();
        }
        for (const v of message.allowList) {
            writer.uint32(18).string(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSendAuthorization();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.spendLimit.push(coin_1.Coin.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.allowList.push(reader.string());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseSendAuthorization();
        if (Array.isArray(object?.spendLimit))
            obj.spendLimit = object.spendLimit.map((e) => coin_1.Coin.fromJSON(e));
        if (Array.isArray(object?.allowList))
            obj.allowList = object.allowList.map((e) => String(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.spendLimit) {
            obj.spendLimit = message.spendLimit.map((e) => (e ? coin_1.Coin.toJSON(e) : undefined));
        }
        else {
            obj.spendLimit = [];
        }
        if (message.allowList) {
            obj.allowList = message.allowList.map((e) => e);
        }
        else {
            obj.allowList = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseSendAuthorization();
        message.spendLimit = object.spendLimit?.map((e) => coin_1.Coin.fromPartial(e)) || [];
        message.allowList = object.allowList?.map((e) => e) || [];
        return message;
    },
};
//# sourceMappingURL=authz.js.map