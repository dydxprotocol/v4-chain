"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.PacketSequence = exports.GenesisState = exports.protobufPackage = void 0;
/* eslint-disable */
const channel_1 = require("./channel");
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "ibc.core.channel.v1";
function createBaseGenesisState() {
    return {
        channels: [],
        acknowledgements: [],
        commitments: [],
        receipts: [],
        sendSequences: [],
        recvSequences: [],
        ackSequences: [],
        nextChannelSequence: BigInt(0),
    };
}
exports.GenesisState = {
    typeUrl: "/ibc.core.channel.v1.GenesisState",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.channels) {
            channel_1.IdentifiedChannel.encode(v, writer.uint32(10).fork()).ldelim();
        }
        for (const v of message.acknowledgements) {
            channel_1.PacketState.encode(v, writer.uint32(18).fork()).ldelim();
        }
        for (const v of message.commitments) {
            channel_1.PacketState.encode(v, writer.uint32(26).fork()).ldelim();
        }
        for (const v of message.receipts) {
            channel_1.PacketState.encode(v, writer.uint32(34).fork()).ldelim();
        }
        for (const v of message.sendSequences) {
            exports.PacketSequence.encode(v, writer.uint32(42).fork()).ldelim();
        }
        for (const v of message.recvSequences) {
            exports.PacketSequence.encode(v, writer.uint32(50).fork()).ldelim();
        }
        for (const v of message.ackSequences) {
            exports.PacketSequence.encode(v, writer.uint32(58).fork()).ldelim();
        }
        if (message.nextChannelSequence !== BigInt(0)) {
            writer.uint32(64).uint64(message.nextChannelSequence);
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
                    message.channels.push(channel_1.IdentifiedChannel.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.acknowledgements.push(channel_1.PacketState.decode(reader, reader.uint32()));
                    break;
                case 3:
                    message.commitments.push(channel_1.PacketState.decode(reader, reader.uint32()));
                    break;
                case 4:
                    message.receipts.push(channel_1.PacketState.decode(reader, reader.uint32()));
                    break;
                case 5:
                    message.sendSequences.push(exports.PacketSequence.decode(reader, reader.uint32()));
                    break;
                case 6:
                    message.recvSequences.push(exports.PacketSequence.decode(reader, reader.uint32()));
                    break;
                case 7:
                    message.ackSequences.push(exports.PacketSequence.decode(reader, reader.uint32()));
                    break;
                case 8:
                    message.nextChannelSequence = reader.uint64();
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
        if (Array.isArray(object?.channels))
            obj.channels = object.channels.map((e) => channel_1.IdentifiedChannel.fromJSON(e));
        if (Array.isArray(object?.acknowledgements))
            obj.acknowledgements = object.acknowledgements.map((e) => channel_1.PacketState.fromJSON(e));
        if (Array.isArray(object?.commitments))
            obj.commitments = object.commitments.map((e) => channel_1.PacketState.fromJSON(e));
        if (Array.isArray(object?.receipts))
            obj.receipts = object.receipts.map((e) => channel_1.PacketState.fromJSON(e));
        if (Array.isArray(object?.sendSequences))
            obj.sendSequences = object.sendSequences.map((e) => exports.PacketSequence.fromJSON(e));
        if (Array.isArray(object?.recvSequences))
            obj.recvSequences = object.recvSequences.map((e) => exports.PacketSequence.fromJSON(e));
        if (Array.isArray(object?.ackSequences))
            obj.ackSequences = object.ackSequences.map((e) => exports.PacketSequence.fromJSON(e));
        if ((0, helpers_1.isSet)(object.nextChannelSequence))
            obj.nextChannelSequence = BigInt(object.nextChannelSequence.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.channels) {
            obj.channels = message.channels.map((e) => (e ? channel_1.IdentifiedChannel.toJSON(e) : undefined));
        }
        else {
            obj.channels = [];
        }
        if (message.acknowledgements) {
            obj.acknowledgements = message.acknowledgements.map((e) => (e ? channel_1.PacketState.toJSON(e) : undefined));
        }
        else {
            obj.acknowledgements = [];
        }
        if (message.commitments) {
            obj.commitments = message.commitments.map((e) => (e ? channel_1.PacketState.toJSON(e) : undefined));
        }
        else {
            obj.commitments = [];
        }
        if (message.receipts) {
            obj.receipts = message.receipts.map((e) => (e ? channel_1.PacketState.toJSON(e) : undefined));
        }
        else {
            obj.receipts = [];
        }
        if (message.sendSequences) {
            obj.sendSequences = message.sendSequences.map((e) => (e ? exports.PacketSequence.toJSON(e) : undefined));
        }
        else {
            obj.sendSequences = [];
        }
        if (message.recvSequences) {
            obj.recvSequences = message.recvSequences.map((e) => (e ? exports.PacketSequence.toJSON(e) : undefined));
        }
        else {
            obj.recvSequences = [];
        }
        if (message.ackSequences) {
            obj.ackSequences = message.ackSequences.map((e) => (e ? exports.PacketSequence.toJSON(e) : undefined));
        }
        else {
            obj.ackSequences = [];
        }
        message.nextChannelSequence !== undefined &&
            (obj.nextChannelSequence = (message.nextChannelSequence || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGenesisState();
        message.channels = object.channels?.map((e) => channel_1.IdentifiedChannel.fromPartial(e)) || [];
        message.acknowledgements = object.acknowledgements?.map((e) => channel_1.PacketState.fromPartial(e)) || [];
        message.commitments = object.commitments?.map((e) => channel_1.PacketState.fromPartial(e)) || [];
        message.receipts = object.receipts?.map((e) => channel_1.PacketState.fromPartial(e)) || [];
        message.sendSequences = object.sendSequences?.map((e) => exports.PacketSequence.fromPartial(e)) || [];
        message.recvSequences = object.recvSequences?.map((e) => exports.PacketSequence.fromPartial(e)) || [];
        message.ackSequences = object.ackSequences?.map((e) => exports.PacketSequence.fromPartial(e)) || [];
        if (object.nextChannelSequence !== undefined && object.nextChannelSequence !== null) {
            message.nextChannelSequence = BigInt(object.nextChannelSequence.toString());
        }
        return message;
    },
};
function createBasePacketSequence() {
    return {
        portId: "",
        channelId: "",
        sequence: BigInt(0),
    };
}
exports.PacketSequence = {
    typeUrl: "/ibc.core.channel.v1.PacketSequence",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.portId !== "") {
            writer.uint32(10).string(message.portId);
        }
        if (message.channelId !== "") {
            writer.uint32(18).string(message.channelId);
        }
        if (message.sequence !== BigInt(0)) {
            writer.uint32(24).uint64(message.sequence);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBasePacketSequence();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.portId = reader.string();
                    break;
                case 2:
                    message.channelId = reader.string();
                    break;
                case 3:
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
        const obj = createBasePacketSequence();
        if ((0, helpers_1.isSet)(object.portId))
            obj.portId = String(object.portId);
        if ((0, helpers_1.isSet)(object.channelId))
            obj.channelId = String(object.channelId);
        if ((0, helpers_1.isSet)(object.sequence))
            obj.sequence = BigInt(object.sequence.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.portId !== undefined && (obj.portId = message.portId);
        message.channelId !== undefined && (obj.channelId = message.channelId);
        message.sequence !== undefined && (obj.sequence = (message.sequence || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBasePacketSequence();
        message.portId = object.portId ?? "";
        message.channelId = object.channelId ?? "";
        if (object.sequence !== undefined && object.sequence !== null) {
            message.sequence = BigInt(object.sequence.toString());
        }
        return message;
    },
};
//# sourceMappingURL=genesis.js.map