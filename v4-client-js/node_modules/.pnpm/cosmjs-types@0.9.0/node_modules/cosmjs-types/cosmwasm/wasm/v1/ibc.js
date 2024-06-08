"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.MsgIBCCloseChannel = exports.MsgIBCSendResponse = exports.MsgIBCSend = exports.protobufPackage = void 0;
/* eslint-disable */
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmwasm.wasm.v1";
function createBaseMsgIBCSend() {
    return {
        channel: "",
        timeoutHeight: BigInt(0),
        timeoutTimestamp: BigInt(0),
        data: new Uint8Array(),
    };
}
exports.MsgIBCSend = {
    typeUrl: "/cosmwasm.wasm.v1.MsgIBCSend",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.channel !== "") {
            writer.uint32(18).string(message.channel);
        }
        if (message.timeoutHeight !== BigInt(0)) {
            writer.uint32(32).uint64(message.timeoutHeight);
        }
        if (message.timeoutTimestamp !== BigInt(0)) {
            writer.uint32(40).uint64(message.timeoutTimestamp);
        }
        if (message.data.length !== 0) {
            writer.uint32(50).bytes(message.data);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgIBCSend();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 2:
                    message.channel = reader.string();
                    break;
                case 4:
                    message.timeoutHeight = reader.uint64();
                    break;
                case 5:
                    message.timeoutTimestamp = reader.uint64();
                    break;
                case 6:
                    message.data = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgIBCSend();
        if ((0, helpers_1.isSet)(object.channel))
            obj.channel = String(object.channel);
        if ((0, helpers_1.isSet)(object.timeoutHeight))
            obj.timeoutHeight = BigInt(object.timeoutHeight.toString());
        if ((0, helpers_1.isSet)(object.timeoutTimestamp))
            obj.timeoutTimestamp = BigInt(object.timeoutTimestamp.toString());
        if ((0, helpers_1.isSet)(object.data))
            obj.data = (0, helpers_1.bytesFromBase64)(object.data);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.channel !== undefined && (obj.channel = message.channel);
        message.timeoutHeight !== undefined &&
            (obj.timeoutHeight = (message.timeoutHeight || BigInt(0)).toString());
        message.timeoutTimestamp !== undefined &&
            (obj.timeoutTimestamp = (message.timeoutTimestamp || BigInt(0)).toString());
        message.data !== undefined &&
            (obj.data = (0, helpers_1.base64FromBytes)(message.data !== undefined ? message.data : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgIBCSend();
        message.channel = object.channel ?? "";
        if (object.timeoutHeight !== undefined && object.timeoutHeight !== null) {
            message.timeoutHeight = BigInt(object.timeoutHeight.toString());
        }
        if (object.timeoutTimestamp !== undefined && object.timeoutTimestamp !== null) {
            message.timeoutTimestamp = BigInt(object.timeoutTimestamp.toString());
        }
        message.data = object.data ?? new Uint8Array();
        return message;
    },
};
function createBaseMsgIBCSendResponse() {
    return {
        sequence: BigInt(0),
    };
}
exports.MsgIBCSendResponse = {
    typeUrl: "/cosmwasm.wasm.v1.MsgIBCSendResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.sequence !== BigInt(0)) {
            writer.uint32(8).uint64(message.sequence);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgIBCSendResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.sequence = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgIBCSendResponse();
        if ((0, helpers_1.isSet)(object.sequence))
            obj.sequence = BigInt(object.sequence.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.sequence !== undefined && (obj.sequence = (message.sequence || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgIBCSendResponse();
        if (object.sequence !== undefined && object.sequence !== null) {
            message.sequence = BigInt(object.sequence.toString());
        }
        return message;
    },
};
function createBaseMsgIBCCloseChannel() {
    return {
        channel: "",
    };
}
exports.MsgIBCCloseChannel = {
    typeUrl: "/cosmwasm.wasm.v1.MsgIBCCloseChannel",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.channel !== "") {
            writer.uint32(18).string(message.channel);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgIBCCloseChannel();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 2:
                    message.channel = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgIBCCloseChannel();
        if ((0, helpers_1.isSet)(object.channel))
            obj.channel = String(object.channel);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.channel !== undefined && (obj.channel = message.channel);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgIBCCloseChannel();
        message.channel = object.channel ?? "";
        return message;
    },
};
//# sourceMappingURL=ibc.js.map