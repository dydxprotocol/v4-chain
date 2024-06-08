"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ForwardRelayerAddress = exports.RegisteredCounterpartyPayee = exports.RegisteredPayee = exports.FeeEnabledChannel = exports.GenesisState = exports.protobufPackage = void 0;
/* eslint-disable */
const fee_1 = require("./fee");
const channel_1 = require("../../../core/channel/v1/channel");
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "ibc.applications.fee.v1";
function createBaseGenesisState() {
    return {
        identifiedFees: [],
        feeEnabledChannels: [],
        registeredPayees: [],
        registeredCounterpartyPayees: [],
        forwardRelayers: [],
    };
}
exports.GenesisState = {
    typeUrl: "/ibc.applications.fee.v1.GenesisState",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.identifiedFees) {
            fee_1.IdentifiedPacketFees.encode(v, writer.uint32(10).fork()).ldelim();
        }
        for (const v of message.feeEnabledChannels) {
            exports.FeeEnabledChannel.encode(v, writer.uint32(18).fork()).ldelim();
        }
        for (const v of message.registeredPayees) {
            exports.RegisteredPayee.encode(v, writer.uint32(26).fork()).ldelim();
        }
        for (const v of message.registeredCounterpartyPayees) {
            exports.RegisteredCounterpartyPayee.encode(v, writer.uint32(34).fork()).ldelim();
        }
        for (const v of message.forwardRelayers) {
            exports.ForwardRelayerAddress.encode(v, writer.uint32(42).fork()).ldelim();
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
                    message.identifiedFees.push(fee_1.IdentifiedPacketFees.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.feeEnabledChannels.push(exports.FeeEnabledChannel.decode(reader, reader.uint32()));
                    break;
                case 3:
                    message.registeredPayees.push(exports.RegisteredPayee.decode(reader, reader.uint32()));
                    break;
                case 4:
                    message.registeredCounterpartyPayees.push(exports.RegisteredCounterpartyPayee.decode(reader, reader.uint32()));
                    break;
                case 5:
                    message.forwardRelayers.push(exports.ForwardRelayerAddress.decode(reader, reader.uint32()));
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
        if (Array.isArray(object?.identifiedFees))
            obj.identifiedFees = object.identifiedFees.map((e) => fee_1.IdentifiedPacketFees.fromJSON(e));
        if (Array.isArray(object?.feeEnabledChannels))
            obj.feeEnabledChannels = object.feeEnabledChannels.map((e) => exports.FeeEnabledChannel.fromJSON(e));
        if (Array.isArray(object?.registeredPayees))
            obj.registeredPayees = object.registeredPayees.map((e) => exports.RegisteredPayee.fromJSON(e));
        if (Array.isArray(object?.registeredCounterpartyPayees))
            obj.registeredCounterpartyPayees = object.registeredCounterpartyPayees.map((e) => exports.RegisteredCounterpartyPayee.fromJSON(e));
        if (Array.isArray(object?.forwardRelayers))
            obj.forwardRelayers = object.forwardRelayers.map((e) => exports.ForwardRelayerAddress.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.identifiedFees) {
            obj.identifiedFees = message.identifiedFees.map((e) => e ? fee_1.IdentifiedPacketFees.toJSON(e) : undefined);
        }
        else {
            obj.identifiedFees = [];
        }
        if (message.feeEnabledChannels) {
            obj.feeEnabledChannels = message.feeEnabledChannels.map((e) => e ? exports.FeeEnabledChannel.toJSON(e) : undefined);
        }
        else {
            obj.feeEnabledChannels = [];
        }
        if (message.registeredPayees) {
            obj.registeredPayees = message.registeredPayees.map((e) => (e ? exports.RegisteredPayee.toJSON(e) : undefined));
        }
        else {
            obj.registeredPayees = [];
        }
        if (message.registeredCounterpartyPayees) {
            obj.registeredCounterpartyPayees = message.registeredCounterpartyPayees.map((e) => e ? exports.RegisteredCounterpartyPayee.toJSON(e) : undefined);
        }
        else {
            obj.registeredCounterpartyPayees = [];
        }
        if (message.forwardRelayers) {
            obj.forwardRelayers = message.forwardRelayers.map((e) => e ? exports.ForwardRelayerAddress.toJSON(e) : undefined);
        }
        else {
            obj.forwardRelayers = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGenesisState();
        message.identifiedFees = object.identifiedFees?.map((e) => fee_1.IdentifiedPacketFees.fromPartial(e)) || [];
        message.feeEnabledChannels =
            object.feeEnabledChannels?.map((e) => exports.FeeEnabledChannel.fromPartial(e)) || [];
        message.registeredPayees = object.registeredPayees?.map((e) => exports.RegisteredPayee.fromPartial(e)) || [];
        message.registeredCounterpartyPayees =
            object.registeredCounterpartyPayees?.map((e) => exports.RegisteredCounterpartyPayee.fromPartial(e)) || [];
        message.forwardRelayers = object.forwardRelayers?.map((e) => exports.ForwardRelayerAddress.fromPartial(e)) || [];
        return message;
    },
};
function createBaseFeeEnabledChannel() {
    return {
        portId: "",
        channelId: "",
    };
}
exports.FeeEnabledChannel = {
    typeUrl: "/ibc.applications.fee.v1.FeeEnabledChannel",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.portId !== "") {
            writer.uint32(10).string(message.portId);
        }
        if (message.channelId !== "") {
            writer.uint32(18).string(message.channelId);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseFeeEnabledChannel();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.portId = reader.string();
                    break;
                case 2:
                    message.channelId = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseFeeEnabledChannel();
        if ((0, helpers_1.isSet)(object.portId))
            obj.portId = String(object.portId);
        if ((0, helpers_1.isSet)(object.channelId))
            obj.channelId = String(object.channelId);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.portId !== undefined && (obj.portId = message.portId);
        message.channelId !== undefined && (obj.channelId = message.channelId);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseFeeEnabledChannel();
        message.portId = object.portId ?? "";
        message.channelId = object.channelId ?? "";
        return message;
    },
};
function createBaseRegisteredPayee() {
    return {
        channelId: "",
        relayer: "",
        payee: "",
    };
}
exports.RegisteredPayee = {
    typeUrl: "/ibc.applications.fee.v1.RegisteredPayee",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.channelId !== "") {
            writer.uint32(10).string(message.channelId);
        }
        if (message.relayer !== "") {
            writer.uint32(18).string(message.relayer);
        }
        if (message.payee !== "") {
            writer.uint32(26).string(message.payee);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseRegisteredPayee();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.channelId = reader.string();
                    break;
                case 2:
                    message.relayer = reader.string();
                    break;
                case 3:
                    message.payee = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseRegisteredPayee();
        if ((0, helpers_1.isSet)(object.channelId))
            obj.channelId = String(object.channelId);
        if ((0, helpers_1.isSet)(object.relayer))
            obj.relayer = String(object.relayer);
        if ((0, helpers_1.isSet)(object.payee))
            obj.payee = String(object.payee);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.channelId !== undefined && (obj.channelId = message.channelId);
        message.relayer !== undefined && (obj.relayer = message.relayer);
        message.payee !== undefined && (obj.payee = message.payee);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseRegisteredPayee();
        message.channelId = object.channelId ?? "";
        message.relayer = object.relayer ?? "";
        message.payee = object.payee ?? "";
        return message;
    },
};
function createBaseRegisteredCounterpartyPayee() {
    return {
        channelId: "",
        relayer: "",
        counterpartyPayee: "",
    };
}
exports.RegisteredCounterpartyPayee = {
    typeUrl: "/ibc.applications.fee.v1.RegisteredCounterpartyPayee",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.channelId !== "") {
            writer.uint32(10).string(message.channelId);
        }
        if (message.relayer !== "") {
            writer.uint32(18).string(message.relayer);
        }
        if (message.counterpartyPayee !== "") {
            writer.uint32(26).string(message.counterpartyPayee);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseRegisteredCounterpartyPayee();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.channelId = reader.string();
                    break;
                case 2:
                    message.relayer = reader.string();
                    break;
                case 3:
                    message.counterpartyPayee = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseRegisteredCounterpartyPayee();
        if ((0, helpers_1.isSet)(object.channelId))
            obj.channelId = String(object.channelId);
        if ((0, helpers_1.isSet)(object.relayer))
            obj.relayer = String(object.relayer);
        if ((0, helpers_1.isSet)(object.counterpartyPayee))
            obj.counterpartyPayee = String(object.counterpartyPayee);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.channelId !== undefined && (obj.channelId = message.channelId);
        message.relayer !== undefined && (obj.relayer = message.relayer);
        message.counterpartyPayee !== undefined && (obj.counterpartyPayee = message.counterpartyPayee);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseRegisteredCounterpartyPayee();
        message.channelId = object.channelId ?? "";
        message.relayer = object.relayer ?? "";
        message.counterpartyPayee = object.counterpartyPayee ?? "";
        return message;
    },
};
function createBaseForwardRelayerAddress() {
    return {
        address: "",
        packetId: channel_1.PacketId.fromPartial({}),
    };
}
exports.ForwardRelayerAddress = {
    typeUrl: "/ibc.applications.fee.v1.ForwardRelayerAddress",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.address !== "") {
            writer.uint32(10).string(message.address);
        }
        if (message.packetId !== undefined) {
            channel_1.PacketId.encode(message.packetId, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseForwardRelayerAddress();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.address = reader.string();
                    break;
                case 2:
                    message.packetId = channel_1.PacketId.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseForwardRelayerAddress();
        if ((0, helpers_1.isSet)(object.address))
            obj.address = String(object.address);
        if ((0, helpers_1.isSet)(object.packetId))
            obj.packetId = channel_1.PacketId.fromJSON(object.packetId);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.address !== undefined && (obj.address = message.address);
        message.packetId !== undefined &&
            (obj.packetId = message.packetId ? channel_1.PacketId.toJSON(message.packetId) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseForwardRelayerAddress();
        message.address = object.address ?? "";
        if (object.packetId !== undefined && object.packetId !== null) {
            message.packetId = channel_1.PacketId.fromPartial(object.packetId);
        }
        return message;
    },
};
//# sourceMappingURL=genesis.js.map