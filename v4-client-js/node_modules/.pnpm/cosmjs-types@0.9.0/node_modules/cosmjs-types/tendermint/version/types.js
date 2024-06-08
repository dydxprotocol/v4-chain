"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Consensus = exports.App = exports.protobufPackage = void 0;
/* eslint-disable */
const binary_1 = require("../../binary");
const helpers_1 = require("../../helpers");
exports.protobufPackage = "tendermint.version";
function createBaseApp() {
    return {
        protocol: BigInt(0),
        software: "",
    };
}
exports.App = {
    typeUrl: "/tendermint.version.App",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.protocol !== BigInt(0)) {
            writer.uint32(8).uint64(message.protocol);
        }
        if (message.software !== "") {
            writer.uint32(18).string(message.software);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseApp();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.protocol = reader.uint64();
                    break;
                case 2:
                    message.software = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseApp();
        if ((0, helpers_1.isSet)(object.protocol))
            obj.protocol = BigInt(object.protocol.toString());
        if ((0, helpers_1.isSet)(object.software))
            obj.software = String(object.software);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.protocol !== undefined && (obj.protocol = (message.protocol || BigInt(0)).toString());
        message.software !== undefined && (obj.software = message.software);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseApp();
        if (object.protocol !== undefined && object.protocol !== null) {
            message.protocol = BigInt(object.protocol.toString());
        }
        message.software = object.software ?? "";
        return message;
    },
};
function createBaseConsensus() {
    return {
        block: BigInt(0),
        app: BigInt(0),
    };
}
exports.Consensus = {
    typeUrl: "/tendermint.version.Consensus",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.block !== BigInt(0)) {
            writer.uint32(8).uint64(message.block);
        }
        if (message.app !== BigInt(0)) {
            writer.uint32(16).uint64(message.app);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseConsensus();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.block = reader.uint64();
                    break;
                case 2:
                    message.app = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseConsensus();
        if ((0, helpers_1.isSet)(object.block))
            obj.block = BigInt(object.block.toString());
        if ((0, helpers_1.isSet)(object.app))
            obj.app = BigInt(object.app.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.block !== undefined && (obj.block = (message.block || BigInt(0)).toString());
        message.app !== undefined && (obj.app = (message.app || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseConsensus();
        if (object.block !== undefined && object.block !== null) {
            message.block = BigInt(object.block.toString());
        }
        if (object.app !== undefined && object.app !== null) {
            message.app = BigInt(object.app.toString());
        }
        return message;
    },
};
//# sourceMappingURL=types.js.map