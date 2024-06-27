"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Module = exports.protobufPackage = void 0;
/* eslint-disable */
const binary_1 = require("../../../../binary");
exports.protobufPackage = "cosmos.nft.module.v1";
function createBaseModule() {
    return {};
}
exports.Module = {
    typeUrl: "/cosmos.nft.module.v1.Module",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseModule();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(_) {
        const obj = createBaseModule();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseModule();
        return message;
    },
};
//# sourceMappingURL=module.js.map