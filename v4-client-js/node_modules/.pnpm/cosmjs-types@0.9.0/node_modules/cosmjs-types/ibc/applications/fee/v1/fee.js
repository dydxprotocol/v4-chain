"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.IdentifiedPacketFees = exports.PacketFees = exports.PacketFee = exports.Fee = exports.protobufPackage = void 0;
/* eslint-disable */
const coin_1 = require("../../../../cosmos/base/v1beta1/coin");
const channel_1 = require("../../../core/channel/v1/channel");
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "ibc.applications.fee.v1";
function createBaseFee() {
    return {
        recvFee: [],
        ackFee: [],
        timeoutFee: [],
    };
}
exports.Fee = {
    typeUrl: "/ibc.applications.fee.v1.Fee",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.recvFee) {
            coin_1.Coin.encode(v, writer.uint32(10).fork()).ldelim();
        }
        for (const v of message.ackFee) {
            coin_1.Coin.encode(v, writer.uint32(18).fork()).ldelim();
        }
        for (const v of message.timeoutFee) {
            coin_1.Coin.encode(v, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseFee();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.recvFee.push(coin_1.Coin.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.ackFee.push(coin_1.Coin.decode(reader, reader.uint32()));
                    break;
                case 3:
                    message.timeoutFee.push(coin_1.Coin.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseFee();
        if (Array.isArray(object?.recvFee))
            obj.recvFee = object.recvFee.map((e) => coin_1.Coin.fromJSON(e));
        if (Array.isArray(object?.ackFee))
            obj.ackFee = object.ackFee.map((e) => coin_1.Coin.fromJSON(e));
        if (Array.isArray(object?.timeoutFee))
            obj.timeoutFee = object.timeoutFee.map((e) => coin_1.Coin.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.recvFee) {
            obj.recvFee = message.recvFee.map((e) => (e ? coin_1.Coin.toJSON(e) : undefined));
        }
        else {
            obj.recvFee = [];
        }
        if (message.ackFee) {
            obj.ackFee = message.ackFee.map((e) => (e ? coin_1.Coin.toJSON(e) : undefined));
        }
        else {
            obj.ackFee = [];
        }
        if (message.timeoutFee) {
            obj.timeoutFee = message.timeoutFee.map((e) => (e ? coin_1.Coin.toJSON(e) : undefined));
        }
        else {
            obj.timeoutFee = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseFee();
        message.recvFee = object.recvFee?.map((e) => coin_1.Coin.fromPartial(e)) || [];
        message.ackFee = object.ackFee?.map((e) => coin_1.Coin.fromPartial(e)) || [];
        message.timeoutFee = object.timeoutFee?.map((e) => coin_1.Coin.fromPartial(e)) || [];
        return message;
    },
};
function createBasePacketFee() {
    return {
        fee: exports.Fee.fromPartial({}),
        refundAddress: "",
        relayers: [],
    };
}
exports.PacketFee = {
    typeUrl: "/ibc.applications.fee.v1.PacketFee",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.fee !== undefined) {
            exports.Fee.encode(message.fee, writer.uint32(10).fork()).ldelim();
        }
        if (message.refundAddress !== "") {
            writer.uint32(18).string(message.refundAddress);
        }
        for (const v of message.relayers) {
            writer.uint32(26).string(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBasePacketFee();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.fee = exports.Fee.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.refundAddress = reader.string();
                    break;
                case 3:
                    message.relayers.push(reader.string());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBasePacketFee();
        if ((0, helpers_1.isSet)(object.fee))
            obj.fee = exports.Fee.fromJSON(object.fee);
        if ((0, helpers_1.isSet)(object.refundAddress))
            obj.refundAddress = String(object.refundAddress);
        if (Array.isArray(object?.relayers))
            obj.relayers = object.relayers.map((e) => String(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.fee !== undefined && (obj.fee = message.fee ? exports.Fee.toJSON(message.fee) : undefined);
        message.refundAddress !== undefined && (obj.refundAddress = message.refundAddress);
        if (message.relayers) {
            obj.relayers = message.relayers.map((e) => e);
        }
        else {
            obj.relayers = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBasePacketFee();
        if (object.fee !== undefined && object.fee !== null) {
            message.fee = exports.Fee.fromPartial(object.fee);
        }
        message.refundAddress = object.refundAddress ?? "";
        message.relayers = object.relayers?.map((e) => e) || [];
        return message;
    },
};
function createBasePacketFees() {
    return {
        packetFees: [],
    };
}
exports.PacketFees = {
    typeUrl: "/ibc.applications.fee.v1.PacketFees",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.packetFees) {
            exports.PacketFee.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBasePacketFees();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.packetFees.push(exports.PacketFee.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBasePacketFees();
        if (Array.isArray(object?.packetFees))
            obj.packetFees = object.packetFees.map((e) => exports.PacketFee.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.packetFees) {
            obj.packetFees = message.packetFees.map((e) => (e ? exports.PacketFee.toJSON(e) : undefined));
        }
        else {
            obj.packetFees = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBasePacketFees();
        message.packetFees = object.packetFees?.map((e) => exports.PacketFee.fromPartial(e)) || [];
        return message;
    },
};
function createBaseIdentifiedPacketFees() {
    return {
        packetId: channel_1.PacketId.fromPartial({}),
        packetFees: [],
    };
}
exports.IdentifiedPacketFees = {
    typeUrl: "/ibc.applications.fee.v1.IdentifiedPacketFees",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.packetId !== undefined) {
            channel_1.PacketId.encode(message.packetId, writer.uint32(10).fork()).ldelim();
        }
        for (const v of message.packetFees) {
            exports.PacketFee.encode(v, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseIdentifiedPacketFees();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.packetId = channel_1.PacketId.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.packetFees.push(exports.PacketFee.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseIdentifiedPacketFees();
        if ((0, helpers_1.isSet)(object.packetId))
            obj.packetId = channel_1.PacketId.fromJSON(object.packetId);
        if (Array.isArray(object?.packetFees))
            obj.packetFees = object.packetFees.map((e) => exports.PacketFee.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.packetId !== undefined &&
            (obj.packetId = message.packetId ? channel_1.PacketId.toJSON(message.packetId) : undefined);
        if (message.packetFees) {
            obj.packetFees = message.packetFees.map((e) => (e ? exports.PacketFee.toJSON(e) : undefined));
        }
        else {
            obj.packetFees = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseIdentifiedPacketFees();
        if (object.packetId !== undefined && object.packetId !== null) {
            message.packetId = channel_1.PacketId.fromPartial(object.packetId);
        }
        message.packetFees = object.packetFees?.map((e) => exports.PacketFee.fromPartial(e)) || [];
        return message;
    },
};
//# sourceMappingURL=fee.js.map