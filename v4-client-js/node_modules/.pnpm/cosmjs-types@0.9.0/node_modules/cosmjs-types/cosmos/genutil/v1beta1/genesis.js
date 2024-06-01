"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GenesisState = exports.protobufPackage = void 0;
/* eslint-disable */
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.genutil.v1beta1";
function createBaseGenesisState() {
    return {
        genTxs: [],
    };
}
exports.GenesisState = {
    typeUrl: "/cosmos.genutil.v1beta1.GenesisState",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.genTxs) {
            writer.uint32(10).bytes(v);
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
                    message.genTxs.push(reader.bytes());
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
        if (Array.isArray(object?.genTxs))
            obj.genTxs = object.genTxs.map((e) => (0, helpers_1.bytesFromBase64)(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.genTxs) {
            obj.genTxs = message.genTxs.map((e) => (0, helpers_1.base64FromBytes)(e !== undefined ? e : new Uint8Array()));
        }
        else {
            obj.genTxs = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGenesisState();
        message.genTxs = object.genTxs?.map((e) => e) || [];
        return message;
    },
};
//# sourceMappingURL=genesis.js.map