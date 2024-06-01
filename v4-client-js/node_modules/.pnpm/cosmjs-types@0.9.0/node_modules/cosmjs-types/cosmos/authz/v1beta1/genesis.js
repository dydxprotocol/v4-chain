"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GenesisState = exports.protobufPackage = void 0;
/* eslint-disable */
const authz_1 = require("./authz");
const binary_1 = require("../../../binary");
exports.protobufPackage = "cosmos.authz.v1beta1";
function createBaseGenesisState() {
    return {
        authorization: [],
    };
}
exports.GenesisState = {
    typeUrl: "/cosmos.authz.v1beta1.GenesisState",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.authorization) {
            authz_1.GrantAuthorization.encode(v, writer.uint32(10).fork()).ldelim();
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
                    message.authorization.push(authz_1.GrantAuthorization.decode(reader, reader.uint32()));
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
        if (Array.isArray(object?.authorization))
            obj.authorization = object.authorization.map((e) => authz_1.GrantAuthorization.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.authorization) {
            obj.authorization = message.authorization.map((e) => (e ? authz_1.GrantAuthorization.toJSON(e) : undefined));
        }
        else {
            obj.authorization = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGenesisState();
        message.authorization = object.authorization?.map((e) => authz_1.GrantAuthorization.fromPartial(e)) || [];
        return message;
    },
};
//# sourceMappingURL=genesis.js.map