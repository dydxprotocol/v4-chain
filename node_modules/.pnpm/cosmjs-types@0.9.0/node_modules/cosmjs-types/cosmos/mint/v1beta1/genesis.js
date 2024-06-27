"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GenesisState = exports.protobufPackage = void 0;
/* eslint-disable */
const mint_1 = require("./mint");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.mint.v1beta1";
function createBaseGenesisState() {
    return {
        minter: mint_1.Minter.fromPartial({}),
        params: mint_1.Params.fromPartial({}),
    };
}
exports.GenesisState = {
    typeUrl: "/cosmos.mint.v1beta1.GenesisState",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.minter !== undefined) {
            mint_1.Minter.encode(message.minter, writer.uint32(10).fork()).ldelim();
        }
        if (message.params !== undefined) {
            mint_1.Params.encode(message.params, writer.uint32(18).fork()).ldelim();
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
                case 1:
                    message.minter = mint_1.Minter.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.params = mint_1.Params.decode(reader, reader.uint32());
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
        if ((0, helpers_1.isSet)(object.minter))
            obj.minter = mint_1.Minter.fromJSON(object.minter);
        if ((0, helpers_1.isSet)(object.params))
            obj.params = mint_1.Params.fromJSON(object.params);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.minter !== undefined && (obj.minter = message.minter ? mint_1.Minter.toJSON(message.minter) : undefined);
        message.params !== undefined && (obj.params = message.params ? mint_1.Params.toJSON(message.params) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGenesisState();
        if (object.minter !== undefined && object.minter !== null) {
            message.minter = mint_1.Minter.fromPartial(object.minter);
        }
        if (object.params !== undefined && object.params !== null) {
            message.params = mint_1.Params.fromPartial(object.params);
        }
        return message;
    },
};
//# sourceMappingURL=genesis.js.map