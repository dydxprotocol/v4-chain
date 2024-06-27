"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.MsgClientImpl = exports.MsgPayPacketFeeAsyncResponse = exports.MsgPayPacketFeeAsync = exports.MsgPayPacketFeeResponse = exports.MsgPayPacketFee = exports.MsgRegisterCounterpartyPayeeResponse = exports.MsgRegisterCounterpartyPayee = exports.MsgRegisterPayeeResponse = exports.MsgRegisterPayee = exports.protobufPackage = void 0;
/* eslint-disable */
const fee_1 = require("./fee");
const channel_1 = require("../../../core/channel/v1/channel");
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "ibc.applications.fee.v1";
function createBaseMsgRegisterPayee() {
    return {
        portId: "",
        channelId: "",
        relayer: "",
        payee: "",
    };
}
exports.MsgRegisterPayee = {
    typeUrl: "/ibc.applications.fee.v1.MsgRegisterPayee",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.portId !== "") {
            writer.uint32(10).string(message.portId);
        }
        if (message.channelId !== "") {
            writer.uint32(18).string(message.channelId);
        }
        if (message.relayer !== "") {
            writer.uint32(26).string(message.relayer);
        }
        if (message.payee !== "") {
            writer.uint32(34).string(message.payee);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgRegisterPayee();
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
                    message.relayer = reader.string();
                    break;
                case 4:
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
        const obj = createBaseMsgRegisterPayee();
        if ((0, helpers_1.isSet)(object.portId))
            obj.portId = String(object.portId);
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
        message.portId !== undefined && (obj.portId = message.portId);
        message.channelId !== undefined && (obj.channelId = message.channelId);
        message.relayer !== undefined && (obj.relayer = message.relayer);
        message.payee !== undefined && (obj.payee = message.payee);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgRegisterPayee();
        message.portId = object.portId ?? "";
        message.channelId = object.channelId ?? "";
        message.relayer = object.relayer ?? "";
        message.payee = object.payee ?? "";
        return message;
    },
};
function createBaseMsgRegisterPayeeResponse() {
    return {};
}
exports.MsgRegisterPayeeResponse = {
    typeUrl: "/ibc.applications.fee.v1.MsgRegisterPayeeResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgRegisterPayeeResponse();
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
        const obj = createBaseMsgRegisterPayeeResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgRegisterPayeeResponse();
        return message;
    },
};
function createBaseMsgRegisterCounterpartyPayee() {
    return {
        portId: "",
        channelId: "",
        relayer: "",
        counterpartyPayee: "",
    };
}
exports.MsgRegisterCounterpartyPayee = {
    typeUrl: "/ibc.applications.fee.v1.MsgRegisterCounterpartyPayee",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.portId !== "") {
            writer.uint32(10).string(message.portId);
        }
        if (message.channelId !== "") {
            writer.uint32(18).string(message.channelId);
        }
        if (message.relayer !== "") {
            writer.uint32(26).string(message.relayer);
        }
        if (message.counterpartyPayee !== "") {
            writer.uint32(34).string(message.counterpartyPayee);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgRegisterCounterpartyPayee();
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
                    message.relayer = reader.string();
                    break;
                case 4:
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
        const obj = createBaseMsgRegisterCounterpartyPayee();
        if ((0, helpers_1.isSet)(object.portId))
            obj.portId = String(object.portId);
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
        message.portId !== undefined && (obj.portId = message.portId);
        message.channelId !== undefined && (obj.channelId = message.channelId);
        message.relayer !== undefined && (obj.relayer = message.relayer);
        message.counterpartyPayee !== undefined && (obj.counterpartyPayee = message.counterpartyPayee);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgRegisterCounterpartyPayee();
        message.portId = object.portId ?? "";
        message.channelId = object.channelId ?? "";
        message.relayer = object.relayer ?? "";
        message.counterpartyPayee = object.counterpartyPayee ?? "";
        return message;
    },
};
function createBaseMsgRegisterCounterpartyPayeeResponse() {
    return {};
}
exports.MsgRegisterCounterpartyPayeeResponse = {
    typeUrl: "/ibc.applications.fee.v1.MsgRegisterCounterpartyPayeeResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgRegisterCounterpartyPayeeResponse();
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
        const obj = createBaseMsgRegisterCounterpartyPayeeResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgRegisterCounterpartyPayeeResponse();
        return message;
    },
};
function createBaseMsgPayPacketFee() {
    return {
        fee: fee_1.Fee.fromPartial({}),
        sourcePortId: "",
        sourceChannelId: "",
        signer: "",
        relayers: [],
    };
}
exports.MsgPayPacketFee = {
    typeUrl: "/ibc.applications.fee.v1.MsgPayPacketFee",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.fee !== undefined) {
            fee_1.Fee.encode(message.fee, writer.uint32(10).fork()).ldelim();
        }
        if (message.sourcePortId !== "") {
            writer.uint32(18).string(message.sourcePortId);
        }
        if (message.sourceChannelId !== "") {
            writer.uint32(26).string(message.sourceChannelId);
        }
        if (message.signer !== "") {
            writer.uint32(34).string(message.signer);
        }
        for (const v of message.relayers) {
            writer.uint32(42).string(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgPayPacketFee();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.fee = fee_1.Fee.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.sourcePortId = reader.string();
                    break;
                case 3:
                    message.sourceChannelId = reader.string();
                    break;
                case 4:
                    message.signer = reader.string();
                    break;
                case 5:
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
        const obj = createBaseMsgPayPacketFee();
        if ((0, helpers_1.isSet)(object.fee))
            obj.fee = fee_1.Fee.fromJSON(object.fee);
        if ((0, helpers_1.isSet)(object.sourcePortId))
            obj.sourcePortId = String(object.sourcePortId);
        if ((0, helpers_1.isSet)(object.sourceChannelId))
            obj.sourceChannelId = String(object.sourceChannelId);
        if ((0, helpers_1.isSet)(object.signer))
            obj.signer = String(object.signer);
        if (Array.isArray(object?.relayers))
            obj.relayers = object.relayers.map((e) => String(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.fee !== undefined && (obj.fee = message.fee ? fee_1.Fee.toJSON(message.fee) : undefined);
        message.sourcePortId !== undefined && (obj.sourcePortId = message.sourcePortId);
        message.sourceChannelId !== undefined && (obj.sourceChannelId = message.sourceChannelId);
        message.signer !== undefined && (obj.signer = message.signer);
        if (message.relayers) {
            obj.relayers = message.relayers.map((e) => e);
        }
        else {
            obj.relayers = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgPayPacketFee();
        if (object.fee !== undefined && object.fee !== null) {
            message.fee = fee_1.Fee.fromPartial(object.fee);
        }
        message.sourcePortId = object.sourcePortId ?? "";
        message.sourceChannelId = object.sourceChannelId ?? "";
        message.signer = object.signer ?? "";
        message.relayers = object.relayers?.map((e) => e) || [];
        return message;
    },
};
function createBaseMsgPayPacketFeeResponse() {
    return {};
}
exports.MsgPayPacketFeeResponse = {
    typeUrl: "/ibc.applications.fee.v1.MsgPayPacketFeeResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgPayPacketFeeResponse();
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
        const obj = createBaseMsgPayPacketFeeResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgPayPacketFeeResponse();
        return message;
    },
};
function createBaseMsgPayPacketFeeAsync() {
    return {
        packetId: channel_1.PacketId.fromPartial({}),
        packetFee: fee_1.PacketFee.fromPartial({}),
    };
}
exports.MsgPayPacketFeeAsync = {
    typeUrl: "/ibc.applications.fee.v1.MsgPayPacketFeeAsync",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.packetId !== undefined) {
            channel_1.PacketId.encode(message.packetId, writer.uint32(10).fork()).ldelim();
        }
        if (message.packetFee !== undefined) {
            fee_1.PacketFee.encode(message.packetFee, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgPayPacketFeeAsync();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.packetId = channel_1.PacketId.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.packetFee = fee_1.PacketFee.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgPayPacketFeeAsync();
        if ((0, helpers_1.isSet)(object.packetId))
            obj.packetId = channel_1.PacketId.fromJSON(object.packetId);
        if ((0, helpers_1.isSet)(object.packetFee))
            obj.packetFee = fee_1.PacketFee.fromJSON(object.packetFee);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.packetId !== undefined &&
            (obj.packetId = message.packetId ? channel_1.PacketId.toJSON(message.packetId) : undefined);
        message.packetFee !== undefined &&
            (obj.packetFee = message.packetFee ? fee_1.PacketFee.toJSON(message.packetFee) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgPayPacketFeeAsync();
        if (object.packetId !== undefined && object.packetId !== null) {
            message.packetId = channel_1.PacketId.fromPartial(object.packetId);
        }
        if (object.packetFee !== undefined && object.packetFee !== null) {
            message.packetFee = fee_1.PacketFee.fromPartial(object.packetFee);
        }
        return message;
    },
};
function createBaseMsgPayPacketFeeAsyncResponse() {
    return {};
}
exports.MsgPayPacketFeeAsyncResponse = {
    typeUrl: "/ibc.applications.fee.v1.MsgPayPacketFeeAsyncResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgPayPacketFeeAsyncResponse();
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
        const obj = createBaseMsgPayPacketFeeAsyncResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgPayPacketFeeAsyncResponse();
        return message;
    },
};
class MsgClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
        this.RegisterPayee = this.RegisterPayee.bind(this);
        this.RegisterCounterpartyPayee = this.RegisterCounterpartyPayee.bind(this);
        this.PayPacketFee = this.PayPacketFee.bind(this);
        this.PayPacketFeeAsync = this.PayPacketFeeAsync.bind(this);
    }
    RegisterPayee(request) {
        const data = exports.MsgRegisterPayee.encode(request).finish();
        const promise = this.rpc.request("ibc.applications.fee.v1.Msg", "RegisterPayee", data);
        return promise.then((data) => exports.MsgRegisterPayeeResponse.decode(new binary_1.BinaryReader(data)));
    }
    RegisterCounterpartyPayee(request) {
        const data = exports.MsgRegisterCounterpartyPayee.encode(request).finish();
        const promise = this.rpc.request("ibc.applications.fee.v1.Msg", "RegisterCounterpartyPayee", data);
        return promise.then((data) => exports.MsgRegisterCounterpartyPayeeResponse.decode(new binary_1.BinaryReader(data)));
    }
    PayPacketFee(request) {
        const data = exports.MsgPayPacketFee.encode(request).finish();
        const promise = this.rpc.request("ibc.applications.fee.v1.Msg", "PayPacketFee", data);
        return promise.then((data) => exports.MsgPayPacketFeeResponse.decode(new binary_1.BinaryReader(data)));
    }
    PayPacketFeeAsync(request) {
        const data = exports.MsgPayPacketFeeAsync.encode(request).finish();
        const promise = this.rpc.request("ibc.applications.fee.v1.Msg", "PayPacketFeeAsync", data);
        return promise.then((data) => exports.MsgPayPacketFeeAsyncResponse.decode(new binary_1.BinaryReader(data)));
    }
}
exports.MsgClientImpl = MsgClientImpl;
//# sourceMappingURL=tx.js.map