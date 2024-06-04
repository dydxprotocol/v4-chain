"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.QueryClientImpl = exports.QueryFeeEnabledChannelResponse = exports.QueryFeeEnabledChannelRequest = exports.QueryFeeEnabledChannelsResponse = exports.QueryFeeEnabledChannelsRequest = exports.QueryCounterpartyPayeeResponse = exports.QueryCounterpartyPayeeRequest = exports.QueryPayeeResponse = exports.QueryPayeeRequest = exports.QueryTotalTimeoutFeesResponse = exports.QueryTotalTimeoutFeesRequest = exports.QueryTotalAckFeesResponse = exports.QueryTotalAckFeesRequest = exports.QueryTotalRecvFeesResponse = exports.QueryTotalRecvFeesRequest = exports.QueryIncentivizedPacketsForChannelResponse = exports.QueryIncentivizedPacketsForChannelRequest = exports.QueryIncentivizedPacketResponse = exports.QueryIncentivizedPacketRequest = exports.QueryIncentivizedPacketsResponse = exports.QueryIncentivizedPacketsRequest = exports.protobufPackage = void 0;
/* eslint-disable */
const pagination_1 = require("../../../../cosmos/base/query/v1beta1/pagination");
const channel_1 = require("../../../core/channel/v1/channel");
const fee_1 = require("./fee");
const coin_1 = require("../../../../cosmos/base/v1beta1/coin");
const genesis_1 = require("./genesis");
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "ibc.applications.fee.v1";
function createBaseQueryIncentivizedPacketsRequest() {
    return {
        pagination: undefined,
        queryHeight: BigInt(0),
    };
}
exports.QueryIncentivizedPacketsRequest = {
    typeUrl: "/ibc.applications.fee.v1.QueryIncentivizedPacketsRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.pagination !== undefined) {
            pagination_1.PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
        }
        if (message.queryHeight !== BigInt(0)) {
            writer.uint32(16).uint64(message.queryHeight);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryIncentivizedPacketsRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.pagination = pagination_1.PageRequest.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.queryHeight = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryIncentivizedPacketsRequest();
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageRequest.fromJSON(object.pagination);
        if ((0, helpers_1.isSet)(object.queryHeight))
            obj.queryHeight = BigInt(object.queryHeight.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageRequest.toJSON(message.pagination) : undefined);
        message.queryHeight !== undefined && (obj.queryHeight = (message.queryHeight || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryIncentivizedPacketsRequest();
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageRequest.fromPartial(object.pagination);
        }
        if (object.queryHeight !== undefined && object.queryHeight !== null) {
            message.queryHeight = BigInt(object.queryHeight.toString());
        }
        return message;
    },
};
function createBaseQueryIncentivizedPacketsResponse() {
    return {
        incentivizedPackets: [],
    };
}
exports.QueryIncentivizedPacketsResponse = {
    typeUrl: "/ibc.applications.fee.v1.QueryIncentivizedPacketsResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.incentivizedPackets) {
            fee_1.IdentifiedPacketFees.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryIncentivizedPacketsResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.incentivizedPackets.push(fee_1.IdentifiedPacketFees.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryIncentivizedPacketsResponse();
        if (Array.isArray(object?.incentivizedPackets))
            obj.incentivizedPackets = object.incentivizedPackets.map((e) => fee_1.IdentifiedPacketFees.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.incentivizedPackets) {
            obj.incentivizedPackets = message.incentivizedPackets.map((e) => e ? fee_1.IdentifiedPacketFees.toJSON(e) : undefined);
        }
        else {
            obj.incentivizedPackets = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryIncentivizedPacketsResponse();
        message.incentivizedPackets =
            object.incentivizedPackets?.map((e) => fee_1.IdentifiedPacketFees.fromPartial(e)) || [];
        return message;
    },
};
function createBaseQueryIncentivizedPacketRequest() {
    return {
        packetId: channel_1.PacketId.fromPartial({}),
        queryHeight: BigInt(0),
    };
}
exports.QueryIncentivizedPacketRequest = {
    typeUrl: "/ibc.applications.fee.v1.QueryIncentivizedPacketRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.packetId !== undefined) {
            channel_1.PacketId.encode(message.packetId, writer.uint32(10).fork()).ldelim();
        }
        if (message.queryHeight !== BigInt(0)) {
            writer.uint32(16).uint64(message.queryHeight);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryIncentivizedPacketRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.packetId = channel_1.PacketId.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.queryHeight = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryIncentivizedPacketRequest();
        if ((0, helpers_1.isSet)(object.packetId))
            obj.packetId = channel_1.PacketId.fromJSON(object.packetId);
        if ((0, helpers_1.isSet)(object.queryHeight))
            obj.queryHeight = BigInt(object.queryHeight.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.packetId !== undefined &&
            (obj.packetId = message.packetId ? channel_1.PacketId.toJSON(message.packetId) : undefined);
        message.queryHeight !== undefined && (obj.queryHeight = (message.queryHeight || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryIncentivizedPacketRequest();
        if (object.packetId !== undefined && object.packetId !== null) {
            message.packetId = channel_1.PacketId.fromPartial(object.packetId);
        }
        if (object.queryHeight !== undefined && object.queryHeight !== null) {
            message.queryHeight = BigInt(object.queryHeight.toString());
        }
        return message;
    },
};
function createBaseQueryIncentivizedPacketResponse() {
    return {
        incentivizedPacket: fee_1.IdentifiedPacketFees.fromPartial({}),
    };
}
exports.QueryIncentivizedPacketResponse = {
    typeUrl: "/ibc.applications.fee.v1.QueryIncentivizedPacketResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.incentivizedPacket !== undefined) {
            fee_1.IdentifiedPacketFees.encode(message.incentivizedPacket, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryIncentivizedPacketResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.incentivizedPacket = fee_1.IdentifiedPacketFees.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryIncentivizedPacketResponse();
        if ((0, helpers_1.isSet)(object.incentivizedPacket))
            obj.incentivizedPacket = fee_1.IdentifiedPacketFees.fromJSON(object.incentivizedPacket);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.incentivizedPacket !== undefined &&
            (obj.incentivizedPacket = message.incentivizedPacket
                ? fee_1.IdentifiedPacketFees.toJSON(message.incentivizedPacket)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryIncentivizedPacketResponse();
        if (object.incentivizedPacket !== undefined && object.incentivizedPacket !== null) {
            message.incentivizedPacket = fee_1.IdentifiedPacketFees.fromPartial(object.incentivizedPacket);
        }
        return message;
    },
};
function createBaseQueryIncentivizedPacketsForChannelRequest() {
    return {
        pagination: undefined,
        portId: "",
        channelId: "",
        queryHeight: BigInt(0),
    };
}
exports.QueryIncentivizedPacketsForChannelRequest = {
    typeUrl: "/ibc.applications.fee.v1.QueryIncentivizedPacketsForChannelRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.pagination !== undefined) {
            pagination_1.PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
        }
        if (message.portId !== "") {
            writer.uint32(18).string(message.portId);
        }
        if (message.channelId !== "") {
            writer.uint32(26).string(message.channelId);
        }
        if (message.queryHeight !== BigInt(0)) {
            writer.uint32(32).uint64(message.queryHeight);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryIncentivizedPacketsForChannelRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.pagination = pagination_1.PageRequest.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.portId = reader.string();
                    break;
                case 3:
                    message.channelId = reader.string();
                    break;
                case 4:
                    message.queryHeight = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryIncentivizedPacketsForChannelRequest();
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageRequest.fromJSON(object.pagination);
        if ((0, helpers_1.isSet)(object.portId))
            obj.portId = String(object.portId);
        if ((0, helpers_1.isSet)(object.channelId))
            obj.channelId = String(object.channelId);
        if ((0, helpers_1.isSet)(object.queryHeight))
            obj.queryHeight = BigInt(object.queryHeight.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageRequest.toJSON(message.pagination) : undefined);
        message.portId !== undefined && (obj.portId = message.portId);
        message.channelId !== undefined && (obj.channelId = message.channelId);
        message.queryHeight !== undefined && (obj.queryHeight = (message.queryHeight || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryIncentivizedPacketsForChannelRequest();
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageRequest.fromPartial(object.pagination);
        }
        message.portId = object.portId ?? "";
        message.channelId = object.channelId ?? "";
        if (object.queryHeight !== undefined && object.queryHeight !== null) {
            message.queryHeight = BigInt(object.queryHeight.toString());
        }
        return message;
    },
};
function createBaseQueryIncentivizedPacketsForChannelResponse() {
    return {
        incentivizedPackets: [],
    };
}
exports.QueryIncentivizedPacketsForChannelResponse = {
    typeUrl: "/ibc.applications.fee.v1.QueryIncentivizedPacketsForChannelResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.incentivizedPackets) {
            fee_1.IdentifiedPacketFees.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryIncentivizedPacketsForChannelResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.incentivizedPackets.push(fee_1.IdentifiedPacketFees.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryIncentivizedPacketsForChannelResponse();
        if (Array.isArray(object?.incentivizedPackets))
            obj.incentivizedPackets = object.incentivizedPackets.map((e) => fee_1.IdentifiedPacketFees.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.incentivizedPackets) {
            obj.incentivizedPackets = message.incentivizedPackets.map((e) => e ? fee_1.IdentifiedPacketFees.toJSON(e) : undefined);
        }
        else {
            obj.incentivizedPackets = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryIncentivizedPacketsForChannelResponse();
        message.incentivizedPackets =
            object.incentivizedPackets?.map((e) => fee_1.IdentifiedPacketFees.fromPartial(e)) || [];
        return message;
    },
};
function createBaseQueryTotalRecvFeesRequest() {
    return {
        packetId: channel_1.PacketId.fromPartial({}),
    };
}
exports.QueryTotalRecvFeesRequest = {
    typeUrl: "/ibc.applications.fee.v1.QueryTotalRecvFeesRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.packetId !== undefined) {
            channel_1.PacketId.encode(message.packetId, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryTotalRecvFeesRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
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
        const obj = createBaseQueryTotalRecvFeesRequest();
        if ((0, helpers_1.isSet)(object.packetId))
            obj.packetId = channel_1.PacketId.fromJSON(object.packetId);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.packetId !== undefined &&
            (obj.packetId = message.packetId ? channel_1.PacketId.toJSON(message.packetId) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryTotalRecvFeesRequest();
        if (object.packetId !== undefined && object.packetId !== null) {
            message.packetId = channel_1.PacketId.fromPartial(object.packetId);
        }
        return message;
    },
};
function createBaseQueryTotalRecvFeesResponse() {
    return {
        recvFees: [],
    };
}
exports.QueryTotalRecvFeesResponse = {
    typeUrl: "/ibc.applications.fee.v1.QueryTotalRecvFeesResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.recvFees) {
            coin_1.Coin.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryTotalRecvFeesResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.recvFees.push(coin_1.Coin.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryTotalRecvFeesResponse();
        if (Array.isArray(object?.recvFees))
            obj.recvFees = object.recvFees.map((e) => coin_1.Coin.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.recvFees) {
            obj.recvFees = message.recvFees.map((e) => (e ? coin_1.Coin.toJSON(e) : undefined));
        }
        else {
            obj.recvFees = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryTotalRecvFeesResponse();
        message.recvFees = object.recvFees?.map((e) => coin_1.Coin.fromPartial(e)) || [];
        return message;
    },
};
function createBaseQueryTotalAckFeesRequest() {
    return {
        packetId: channel_1.PacketId.fromPartial({}),
    };
}
exports.QueryTotalAckFeesRequest = {
    typeUrl: "/ibc.applications.fee.v1.QueryTotalAckFeesRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.packetId !== undefined) {
            channel_1.PacketId.encode(message.packetId, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryTotalAckFeesRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
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
        const obj = createBaseQueryTotalAckFeesRequest();
        if ((0, helpers_1.isSet)(object.packetId))
            obj.packetId = channel_1.PacketId.fromJSON(object.packetId);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.packetId !== undefined &&
            (obj.packetId = message.packetId ? channel_1.PacketId.toJSON(message.packetId) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryTotalAckFeesRequest();
        if (object.packetId !== undefined && object.packetId !== null) {
            message.packetId = channel_1.PacketId.fromPartial(object.packetId);
        }
        return message;
    },
};
function createBaseQueryTotalAckFeesResponse() {
    return {
        ackFees: [],
    };
}
exports.QueryTotalAckFeesResponse = {
    typeUrl: "/ibc.applications.fee.v1.QueryTotalAckFeesResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.ackFees) {
            coin_1.Coin.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryTotalAckFeesResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.ackFees.push(coin_1.Coin.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryTotalAckFeesResponse();
        if (Array.isArray(object?.ackFees))
            obj.ackFees = object.ackFees.map((e) => coin_1.Coin.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.ackFees) {
            obj.ackFees = message.ackFees.map((e) => (e ? coin_1.Coin.toJSON(e) : undefined));
        }
        else {
            obj.ackFees = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryTotalAckFeesResponse();
        message.ackFees = object.ackFees?.map((e) => coin_1.Coin.fromPartial(e)) || [];
        return message;
    },
};
function createBaseQueryTotalTimeoutFeesRequest() {
    return {
        packetId: channel_1.PacketId.fromPartial({}),
    };
}
exports.QueryTotalTimeoutFeesRequest = {
    typeUrl: "/ibc.applications.fee.v1.QueryTotalTimeoutFeesRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.packetId !== undefined) {
            channel_1.PacketId.encode(message.packetId, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryTotalTimeoutFeesRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
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
        const obj = createBaseQueryTotalTimeoutFeesRequest();
        if ((0, helpers_1.isSet)(object.packetId))
            obj.packetId = channel_1.PacketId.fromJSON(object.packetId);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.packetId !== undefined &&
            (obj.packetId = message.packetId ? channel_1.PacketId.toJSON(message.packetId) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryTotalTimeoutFeesRequest();
        if (object.packetId !== undefined && object.packetId !== null) {
            message.packetId = channel_1.PacketId.fromPartial(object.packetId);
        }
        return message;
    },
};
function createBaseQueryTotalTimeoutFeesResponse() {
    return {
        timeoutFees: [],
    };
}
exports.QueryTotalTimeoutFeesResponse = {
    typeUrl: "/ibc.applications.fee.v1.QueryTotalTimeoutFeesResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.timeoutFees) {
            coin_1.Coin.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryTotalTimeoutFeesResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.timeoutFees.push(coin_1.Coin.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryTotalTimeoutFeesResponse();
        if (Array.isArray(object?.timeoutFees))
            obj.timeoutFees = object.timeoutFees.map((e) => coin_1.Coin.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.timeoutFees) {
            obj.timeoutFees = message.timeoutFees.map((e) => (e ? coin_1.Coin.toJSON(e) : undefined));
        }
        else {
            obj.timeoutFees = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryTotalTimeoutFeesResponse();
        message.timeoutFees = object.timeoutFees?.map((e) => coin_1.Coin.fromPartial(e)) || [];
        return message;
    },
};
function createBaseQueryPayeeRequest() {
    return {
        channelId: "",
        relayer: "",
    };
}
exports.QueryPayeeRequest = {
    typeUrl: "/ibc.applications.fee.v1.QueryPayeeRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.channelId !== "") {
            writer.uint32(10).string(message.channelId);
        }
        if (message.relayer !== "") {
            writer.uint32(18).string(message.relayer);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryPayeeRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.channelId = reader.string();
                    break;
                case 2:
                    message.relayer = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryPayeeRequest();
        if ((0, helpers_1.isSet)(object.channelId))
            obj.channelId = String(object.channelId);
        if ((0, helpers_1.isSet)(object.relayer))
            obj.relayer = String(object.relayer);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.channelId !== undefined && (obj.channelId = message.channelId);
        message.relayer !== undefined && (obj.relayer = message.relayer);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryPayeeRequest();
        message.channelId = object.channelId ?? "";
        message.relayer = object.relayer ?? "";
        return message;
    },
};
function createBaseQueryPayeeResponse() {
    return {
        payeeAddress: "",
    };
}
exports.QueryPayeeResponse = {
    typeUrl: "/ibc.applications.fee.v1.QueryPayeeResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.payeeAddress !== "") {
            writer.uint32(10).string(message.payeeAddress);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryPayeeResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.payeeAddress = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryPayeeResponse();
        if ((0, helpers_1.isSet)(object.payeeAddress))
            obj.payeeAddress = String(object.payeeAddress);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.payeeAddress !== undefined && (obj.payeeAddress = message.payeeAddress);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryPayeeResponse();
        message.payeeAddress = object.payeeAddress ?? "";
        return message;
    },
};
function createBaseQueryCounterpartyPayeeRequest() {
    return {
        channelId: "",
        relayer: "",
    };
}
exports.QueryCounterpartyPayeeRequest = {
    typeUrl: "/ibc.applications.fee.v1.QueryCounterpartyPayeeRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.channelId !== "") {
            writer.uint32(10).string(message.channelId);
        }
        if (message.relayer !== "") {
            writer.uint32(18).string(message.relayer);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryCounterpartyPayeeRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.channelId = reader.string();
                    break;
                case 2:
                    message.relayer = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryCounterpartyPayeeRequest();
        if ((0, helpers_1.isSet)(object.channelId))
            obj.channelId = String(object.channelId);
        if ((0, helpers_1.isSet)(object.relayer))
            obj.relayer = String(object.relayer);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.channelId !== undefined && (obj.channelId = message.channelId);
        message.relayer !== undefined && (obj.relayer = message.relayer);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryCounterpartyPayeeRequest();
        message.channelId = object.channelId ?? "";
        message.relayer = object.relayer ?? "";
        return message;
    },
};
function createBaseQueryCounterpartyPayeeResponse() {
    return {
        counterpartyPayee: "",
    };
}
exports.QueryCounterpartyPayeeResponse = {
    typeUrl: "/ibc.applications.fee.v1.QueryCounterpartyPayeeResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.counterpartyPayee !== "") {
            writer.uint32(10).string(message.counterpartyPayee);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryCounterpartyPayeeResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
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
        const obj = createBaseQueryCounterpartyPayeeResponse();
        if ((0, helpers_1.isSet)(object.counterpartyPayee))
            obj.counterpartyPayee = String(object.counterpartyPayee);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.counterpartyPayee !== undefined && (obj.counterpartyPayee = message.counterpartyPayee);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryCounterpartyPayeeResponse();
        message.counterpartyPayee = object.counterpartyPayee ?? "";
        return message;
    },
};
function createBaseQueryFeeEnabledChannelsRequest() {
    return {
        pagination: undefined,
        queryHeight: BigInt(0),
    };
}
exports.QueryFeeEnabledChannelsRequest = {
    typeUrl: "/ibc.applications.fee.v1.QueryFeeEnabledChannelsRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.pagination !== undefined) {
            pagination_1.PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
        }
        if (message.queryHeight !== BigInt(0)) {
            writer.uint32(16).uint64(message.queryHeight);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryFeeEnabledChannelsRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.pagination = pagination_1.PageRequest.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.queryHeight = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryFeeEnabledChannelsRequest();
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageRequest.fromJSON(object.pagination);
        if ((0, helpers_1.isSet)(object.queryHeight))
            obj.queryHeight = BigInt(object.queryHeight.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageRequest.toJSON(message.pagination) : undefined);
        message.queryHeight !== undefined && (obj.queryHeight = (message.queryHeight || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryFeeEnabledChannelsRequest();
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageRequest.fromPartial(object.pagination);
        }
        if (object.queryHeight !== undefined && object.queryHeight !== null) {
            message.queryHeight = BigInt(object.queryHeight.toString());
        }
        return message;
    },
};
function createBaseQueryFeeEnabledChannelsResponse() {
    return {
        feeEnabledChannels: [],
    };
}
exports.QueryFeeEnabledChannelsResponse = {
    typeUrl: "/ibc.applications.fee.v1.QueryFeeEnabledChannelsResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.feeEnabledChannels) {
            genesis_1.FeeEnabledChannel.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryFeeEnabledChannelsResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.feeEnabledChannels.push(genesis_1.FeeEnabledChannel.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryFeeEnabledChannelsResponse();
        if (Array.isArray(object?.feeEnabledChannels))
            obj.feeEnabledChannels = object.feeEnabledChannels.map((e) => genesis_1.FeeEnabledChannel.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.feeEnabledChannels) {
            obj.feeEnabledChannels = message.feeEnabledChannels.map((e) => e ? genesis_1.FeeEnabledChannel.toJSON(e) : undefined);
        }
        else {
            obj.feeEnabledChannels = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryFeeEnabledChannelsResponse();
        message.feeEnabledChannels =
            object.feeEnabledChannels?.map((e) => genesis_1.FeeEnabledChannel.fromPartial(e)) || [];
        return message;
    },
};
function createBaseQueryFeeEnabledChannelRequest() {
    return {
        portId: "",
        channelId: "",
    };
}
exports.QueryFeeEnabledChannelRequest = {
    typeUrl: "/ibc.applications.fee.v1.QueryFeeEnabledChannelRequest",
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
        const message = createBaseQueryFeeEnabledChannelRequest();
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
        const obj = createBaseQueryFeeEnabledChannelRequest();
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
        const message = createBaseQueryFeeEnabledChannelRequest();
        message.portId = object.portId ?? "";
        message.channelId = object.channelId ?? "";
        return message;
    },
};
function createBaseQueryFeeEnabledChannelResponse() {
    return {
        feeEnabled: false,
    };
}
exports.QueryFeeEnabledChannelResponse = {
    typeUrl: "/ibc.applications.fee.v1.QueryFeeEnabledChannelResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.feeEnabled === true) {
            writer.uint32(8).bool(message.feeEnabled);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryFeeEnabledChannelResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.feeEnabled = reader.bool();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryFeeEnabledChannelResponse();
        if ((0, helpers_1.isSet)(object.feeEnabled))
            obj.feeEnabled = Boolean(object.feeEnabled);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.feeEnabled !== undefined && (obj.feeEnabled = message.feeEnabled);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryFeeEnabledChannelResponse();
        message.feeEnabled = object.feeEnabled ?? false;
        return message;
    },
};
class QueryClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
        this.IncentivizedPackets = this.IncentivizedPackets.bind(this);
        this.IncentivizedPacket = this.IncentivizedPacket.bind(this);
        this.IncentivizedPacketsForChannel = this.IncentivizedPacketsForChannel.bind(this);
        this.TotalRecvFees = this.TotalRecvFees.bind(this);
        this.TotalAckFees = this.TotalAckFees.bind(this);
        this.TotalTimeoutFees = this.TotalTimeoutFees.bind(this);
        this.Payee = this.Payee.bind(this);
        this.CounterpartyPayee = this.CounterpartyPayee.bind(this);
        this.FeeEnabledChannels = this.FeeEnabledChannels.bind(this);
        this.FeeEnabledChannel = this.FeeEnabledChannel.bind(this);
    }
    IncentivizedPackets(request) {
        const data = exports.QueryIncentivizedPacketsRequest.encode(request).finish();
        const promise = this.rpc.request("ibc.applications.fee.v1.Query", "IncentivizedPackets", data);
        return promise.then((data) => exports.QueryIncentivizedPacketsResponse.decode(new binary_1.BinaryReader(data)));
    }
    IncentivizedPacket(request) {
        const data = exports.QueryIncentivizedPacketRequest.encode(request).finish();
        const promise = this.rpc.request("ibc.applications.fee.v1.Query", "IncentivizedPacket", data);
        return promise.then((data) => exports.QueryIncentivizedPacketResponse.decode(new binary_1.BinaryReader(data)));
    }
    IncentivizedPacketsForChannel(request) {
        const data = exports.QueryIncentivizedPacketsForChannelRequest.encode(request).finish();
        const promise = this.rpc.request("ibc.applications.fee.v1.Query", "IncentivizedPacketsForChannel", data);
        return promise.then((data) => exports.QueryIncentivizedPacketsForChannelResponse.decode(new binary_1.BinaryReader(data)));
    }
    TotalRecvFees(request) {
        const data = exports.QueryTotalRecvFeesRequest.encode(request).finish();
        const promise = this.rpc.request("ibc.applications.fee.v1.Query", "TotalRecvFees", data);
        return promise.then((data) => exports.QueryTotalRecvFeesResponse.decode(new binary_1.BinaryReader(data)));
    }
    TotalAckFees(request) {
        const data = exports.QueryTotalAckFeesRequest.encode(request).finish();
        const promise = this.rpc.request("ibc.applications.fee.v1.Query", "TotalAckFees", data);
        return promise.then((data) => exports.QueryTotalAckFeesResponse.decode(new binary_1.BinaryReader(data)));
    }
    TotalTimeoutFees(request) {
        const data = exports.QueryTotalTimeoutFeesRequest.encode(request).finish();
        const promise = this.rpc.request("ibc.applications.fee.v1.Query", "TotalTimeoutFees", data);
        return promise.then((data) => exports.QueryTotalTimeoutFeesResponse.decode(new binary_1.BinaryReader(data)));
    }
    Payee(request) {
        const data = exports.QueryPayeeRequest.encode(request).finish();
        const promise = this.rpc.request("ibc.applications.fee.v1.Query", "Payee", data);
        return promise.then((data) => exports.QueryPayeeResponse.decode(new binary_1.BinaryReader(data)));
    }
    CounterpartyPayee(request) {
        const data = exports.QueryCounterpartyPayeeRequest.encode(request).finish();
        const promise = this.rpc.request("ibc.applications.fee.v1.Query", "CounterpartyPayee", data);
        return promise.then((data) => exports.QueryCounterpartyPayeeResponse.decode(new binary_1.BinaryReader(data)));
    }
    FeeEnabledChannels(request) {
        const data = exports.QueryFeeEnabledChannelsRequest.encode(request).finish();
        const promise = this.rpc.request("ibc.applications.fee.v1.Query", "FeeEnabledChannels", data);
        return promise.then((data) => exports.QueryFeeEnabledChannelsResponse.decode(new binary_1.BinaryReader(data)));
    }
    FeeEnabledChannel(request) {
        const data = exports.QueryFeeEnabledChannelRequest.encode(request).finish();
        const promise = this.rpc.request("ibc.applications.fee.v1.Query", "FeeEnabledChannel", data);
        return promise.then((data) => exports.QueryFeeEnabledChannelResponse.decode(new binary_1.BinaryReader(data)));
    }
}
exports.QueryClientImpl = QueryClientImpl;
//# sourceMappingURL=query.js.map