"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GenesisState = exports.protobufPackage = void 0;
/* eslint-disable */
const coin_1 = require("../../base/v1beta1/coin");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.crisis.v1beta1";
function createBaseGenesisState() {
    return {
        constantFee: coin_1.Coin.fromPartial({}),
    };
}
exports.GenesisState = {
    typeUrl: "/cosmos.crisis.v1beta1.GenesisState",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.constantFee !== undefined) {
            coin_1.Coin.encode(message.constantFee, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGenesisState();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 3:
                    message.constantFee = coin_1.Coin.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseGenesisState();
        if ((0, helpers_1.isSet)(object.constantFee))
            obj.constantFee = coin_1.Coin.fromJSON(object.constantFee);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.constantFee !== undefined &&
            (obj.constantFee = message.constantFee ? coin_1.Coin.toJSON(message.constantFee) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGenesisState();
        if (object.constantFee !== undefined && object.constantFee !== null) {
            message.constantFee = coin_1.Coin.fromPartial(object.constantFee);
        }
        return message;
    },
};
//# sourceMappingURL=genesis.js.map