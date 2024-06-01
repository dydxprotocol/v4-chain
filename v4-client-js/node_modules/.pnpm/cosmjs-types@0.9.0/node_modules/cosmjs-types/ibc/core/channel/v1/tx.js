"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.MsgClientImpl = exports.MsgAcknowledgementResponse = exports.MsgAcknowledgement = exports.MsgTimeoutOnCloseResponse = exports.MsgTimeoutOnClose = exports.MsgTimeoutResponse = exports.MsgTimeout = exports.MsgRecvPacketResponse = exports.MsgRecvPacket = exports.MsgChannelCloseConfirmResponse = exports.MsgChannelCloseConfirm = exports.MsgChannelCloseInitResponse = exports.MsgChannelCloseInit = exports.MsgChannelOpenConfirmResponse = exports.MsgChannelOpenConfirm = exports.MsgChannelOpenAckResponse = exports.MsgChannelOpenAck = exports.MsgChannelOpenTryResponse = exports.MsgChannelOpenTry = exports.MsgChannelOpenInitResponse = exports.MsgChannelOpenInit = exports.responseResultTypeToJSON = exports.responseResultTypeFromJSON = exports.ResponseResultType = exports.protobufPackage = void 0;
/* eslint-disable */
const channel_1 = require("./channel");
const client_1 = require("../../client/v1/client");
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "ibc.core.channel.v1";
/** ResponseResultType defines the possible outcomes of the execution of a message */
var ResponseResultType;
(function (ResponseResultType) {
    /** RESPONSE_RESULT_TYPE_UNSPECIFIED - Default zero value enumeration */
    ResponseResultType[ResponseResultType["RESPONSE_RESULT_TYPE_UNSPECIFIED"] = 0] = "RESPONSE_RESULT_TYPE_UNSPECIFIED";
    /** RESPONSE_RESULT_TYPE_NOOP - The message did not call the IBC application callbacks (because, for example, the packet had already been relayed) */
    ResponseResultType[ResponseResultType["RESPONSE_RESULT_TYPE_NOOP"] = 1] = "RESPONSE_RESULT_TYPE_NOOP";
    /** RESPONSE_RESULT_TYPE_SUCCESS - The message was executed successfully */
    ResponseResultType[ResponseResultType["RESPONSE_RESULT_TYPE_SUCCESS"] = 2] = "RESPONSE_RESULT_TYPE_SUCCESS";
    ResponseResultType[ResponseResultType["UNRECOGNIZED"] = -1] = "UNRECOGNIZED";
})(ResponseResultType || (exports.ResponseResultType = ResponseResultType = {}));
function responseResultTypeFromJSON(object) {
    switch (object) {
        case 0:
        case "RESPONSE_RESULT_TYPE_UNSPECIFIED":
            return ResponseResultType.RESPONSE_RESULT_TYPE_UNSPECIFIED;
        case 1:
        case "RESPONSE_RESULT_TYPE_NOOP":
            return ResponseResultType.RESPONSE_RESULT_TYPE_NOOP;
        case 2:
        case "RESPONSE_RESULT_TYPE_SUCCESS":
            return ResponseResultType.RESPONSE_RESULT_TYPE_SUCCESS;
        case -1:
        case "UNRECOGNIZED":
        default:
            return ResponseResultType.UNRECOGNIZED;
    }
}
exports.responseResultTypeFromJSON = responseResultTypeFromJSON;
function responseResultTypeToJSON(object) {
    switch (object) {
        case ResponseResultType.RESPONSE_RESULT_TYPE_UNSPECIFIED:
            return "RESPONSE_RESULT_TYPE_UNSPECIFIED";
        case ResponseResultType.RESPONSE_RESULT_TYPE_NOOP:
            return "RESPONSE_RESULT_TYPE_NOOP";
        case ResponseResultType.RESPONSE_RESULT_TYPE_SUCCESS:
            return "RESPONSE_RESULT_TYPE_SUCCESS";
        case ResponseResultType.UNRECOGNIZED:
        default:
            return "UNRECOGNIZED";
    }
}
exports.responseResultTypeToJSON = responseResultTypeToJSON;
function createBaseMsgChannelOpenInit() {
    return {
        portId: "",
        channel: channel_1.Channel.fromPartial({}),
        signer: "",
    };
}
exports.MsgChannelOpenInit = {
    typeUrl: "/ibc.core.channel.v1.MsgChannelOpenInit",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.portId !== "") {
            writer.uint32(10).string(message.portId);
        }
        if (message.channel !== undefined) {
            channel_1.Channel.encode(message.channel, writer.uint32(18).fork()).ldelim();
        }
        if (message.signer !== "") {
            writer.uint32(26).string(message.signer);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgChannelOpenInit();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.portId = reader.string();
                    break;
                case 2:
                    message.channel = channel_1.Channel.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.signer = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgChannelOpenInit();
        if ((0, helpers_1.isSet)(object.portId))
            obj.portId = String(object.portId);
        if ((0, helpers_1.isSet)(object.channel))
            obj.channel = channel_1.Channel.fromJSON(object.channel);
        if ((0, helpers_1.isSet)(object.signer))
            obj.signer = String(object.signer);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.portId !== undefined && (obj.portId = message.portId);
        message.channel !== undefined &&
            (obj.channel = message.channel ? channel_1.Channel.toJSON(message.channel) : undefined);
        message.signer !== undefined && (obj.signer = message.signer);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgChannelOpenInit();
        message.portId = object.portId ?? "";
        if (object.channel !== undefined && object.channel !== null) {
            message.channel = channel_1.Channel.fromPartial(object.channel);
        }
        message.signer = object.signer ?? "";
        return message;
    },
};
function createBaseMsgChannelOpenInitResponse() {
    return {
        channelId: "",
        version: "",
    };
}
exports.MsgChannelOpenInitResponse = {
    typeUrl: "/ibc.core.channel.v1.MsgChannelOpenInitResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.channelId !== "") {
            writer.uint32(10).string(message.channelId);
        }
        if (message.version !== "") {
            writer.uint32(18).string(message.version);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgChannelOpenInitResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.channelId = reader.string();
                    break;
                case 2:
                    message.version = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgChannelOpenInitResponse();
        if ((0, helpers_1.isSet)(object.channelId))
            obj.channelId = String(object.channelId);
        if ((0, helpers_1.isSet)(object.version))
            obj.version = String(object.version);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.channelId !== undefined && (obj.channelId = message.channelId);
        message.version !== undefined && (obj.version = message.version);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgChannelOpenInitResponse();
        message.channelId = object.channelId ?? "";
        message.version = object.version ?? "";
        return message;
    },
};
function createBaseMsgChannelOpenTry() {
    return {
        portId: "",
        previousChannelId: "",
        channel: channel_1.Channel.fromPartial({}),
        counterpartyVersion: "",
        proofInit: new Uint8Array(),
        proofHeight: client_1.Height.fromPartial({}),
        signer: "",
    };
}
exports.MsgChannelOpenTry = {
    typeUrl: "/ibc.core.channel.v1.MsgChannelOpenTry",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.portId !== "") {
            writer.uint32(10).string(message.portId);
        }
        if (message.previousChannelId !== "") {
            writer.uint32(18).string(message.previousChannelId);
        }
        if (message.channel !== undefined) {
            channel_1.Channel.encode(message.channel, writer.uint32(26).fork()).ldelim();
        }
        if (message.counterpartyVersion !== "") {
            writer.uint32(34).string(message.counterpartyVersion);
        }
        if (message.proofInit.length !== 0) {
            writer.uint32(42).bytes(message.proofInit);
        }
        if (message.proofHeight !== undefined) {
            client_1.Height.encode(message.proofHeight, writer.uint32(50).fork()).ldelim();
        }
        if (message.signer !== "") {
            writer.uint32(58).string(message.signer);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgChannelOpenTry();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.portId = reader.string();
                    break;
                case 2:
                    message.previousChannelId = reader.string();
                    break;
                case 3:
                    message.channel = channel_1.Channel.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.counterpartyVersion = reader.string();
                    break;
                case 5:
                    message.proofInit = reader.bytes();
                    break;
                case 6:
                    message.proofHeight = client_1.Height.decode(reader, reader.uint32());
                    break;
                case 7:
                    message.signer = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgChannelOpenTry();
        if ((0, helpers_1.isSet)(object.portId))
            obj.portId = String(object.portId);
        if ((0, helpers_1.isSet)(object.previousChannelId))
            obj.previousChannelId = String(object.previousChannelId);
        if ((0, helpers_1.isSet)(object.channel))
            obj.channel = channel_1.Channel.fromJSON(object.channel);
        if ((0, helpers_1.isSet)(object.counterpartyVersion))
            obj.counterpartyVersion = String(object.counterpartyVersion);
        if ((0, helpers_1.isSet)(object.proofInit))
            obj.proofInit = (0, helpers_1.bytesFromBase64)(object.proofInit);
        if ((0, helpers_1.isSet)(object.proofHeight))
            obj.proofHeight = client_1.Height.fromJSON(object.proofHeight);
        if ((0, helpers_1.isSet)(object.signer))
            obj.signer = String(object.signer);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.portId !== undefined && (obj.portId = message.portId);
        message.previousChannelId !== undefined && (obj.previousChannelId = message.previousChannelId);
        message.channel !== undefined &&
            (obj.channel = message.channel ? channel_1.Channel.toJSON(message.channel) : undefined);
        message.counterpartyVersion !== undefined && (obj.counterpartyVersion = message.counterpartyVersion);
        message.proofInit !== undefined &&
            (obj.proofInit = (0, helpers_1.base64FromBytes)(message.proofInit !== undefined ? message.proofInit : new Uint8Array()));
        message.proofHeight !== undefined &&
            (obj.proofHeight = message.proofHeight ? client_1.Height.toJSON(message.proofHeight) : undefined);
        message.signer !== undefined && (obj.signer = message.signer);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgChannelOpenTry();
        message.portId = object.portId ?? "";
        message.previousChannelId = object.previousChannelId ?? "";
        if (object.channel !== undefined && object.channel !== null) {
            message.channel = channel_1.Channel.fromPartial(object.channel);
        }
        message.counterpartyVersion = object.counterpartyVersion ?? "";
        message.proofInit = object.proofInit ?? new Uint8Array();
        if (object.proofHeight !== undefined && object.proofHeight !== null) {
            message.proofHeight = client_1.Height.fromPartial(object.proofHeight);
        }
        message.signer = object.signer ?? "";
        return message;
    },
};
function createBaseMsgChannelOpenTryResponse() {
    return {
        version: "",
    };
}
exports.MsgChannelOpenTryResponse = {
    typeUrl: "/ibc.core.channel.v1.MsgChannelOpenTryResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.version !== "") {
            writer.uint32(10).string(message.version);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgChannelOpenTryResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.version = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgChannelOpenTryResponse();
        if ((0, helpers_1.isSet)(object.version))
            obj.version = String(object.version);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.version !== undefined && (obj.version = message.version);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgChannelOpenTryResponse();
        message.version = object.version ?? "";
        return message;
    },
};
function createBaseMsgChannelOpenAck() {
    return {
        portId: "",
        channelId: "",
        counterpartyChannelId: "",
        counterpartyVersion: "",
        proofTry: new Uint8Array(),
        proofHeight: client_1.Height.fromPartial({}),
        signer: "",
    };
}
exports.MsgChannelOpenAck = {
    typeUrl: "/ibc.core.channel.v1.MsgChannelOpenAck",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.portId !== "") {
            writer.uint32(10).string(message.portId);
        }
        if (message.channelId !== "") {
            writer.uint32(18).string(message.channelId);
        }
        if (message.counterpartyChannelId !== "") {
            writer.uint32(26).string(message.counterpartyChannelId);
        }
        if (message.counterpartyVersion !== "") {
            writer.uint32(34).string(message.counterpartyVersion);
        }
        if (message.proofTry.length !== 0) {
            writer.uint32(42).bytes(message.proofTry);
        }
        if (message.proofHeight !== undefined) {
            client_1.Height.encode(message.proofHeight, writer.uint32(50).fork()).ldelim();
        }
        if (message.signer !== "") {
            writer.uint32(58).string(message.signer);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgChannelOpenAck();
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
                    message.counterpartyChannelId = reader.string();
                    break;
                case 4:
                    message.counterpartyVersion = reader.string();
                    break;
                case 5:
                    message.proofTry = reader.bytes();
                    break;
                case 6:
                    message.proofHeight = client_1.Height.decode(reader, reader.uint32());
                    break;
                case 7:
                    message.signer = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgChannelOpenAck();
        if ((0, helpers_1.isSet)(object.portId))
            obj.portId = String(object.portId);
        if ((0, helpers_1.isSet)(object.channelId))
            obj.channelId = String(object.channelId);
        if ((0, helpers_1.isSet)(object.counterpartyChannelId))
            obj.counterpartyChannelId = String(object.counterpartyChannelId);
        if ((0, helpers_1.isSet)(object.counterpartyVersion))
            obj.counterpartyVersion = String(object.counterpartyVersion);
        if ((0, helpers_1.isSet)(object.proofTry))
            obj.proofTry = (0, helpers_1.bytesFromBase64)(object.proofTry);
        if ((0, helpers_1.isSet)(object.proofHeight))
            obj.proofHeight = client_1.Height.fromJSON(object.proofHeight);
        if ((0, helpers_1.isSet)(object.signer))
            obj.signer = String(object.signer);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.portId !== undefined && (obj.portId = message.portId);
        message.channelId !== undefined && (obj.channelId = message.channelId);
        message.counterpartyChannelId !== undefined &&
            (obj.counterpartyChannelId = message.counterpartyChannelId);
        message.counterpartyVersion !== undefined && (obj.counterpartyVersion = message.counterpartyVersion);
        message.proofTry !== undefined &&
            (obj.proofTry = (0, helpers_1.base64FromBytes)(message.proofTry !== undefined ? message.proofTry : new Uint8Array()));
        message.proofHeight !== undefined &&
            (obj.proofHeight = message.proofHeight ? client_1.Height.toJSON(message.proofHeight) : undefined);
        message.signer !== undefined && (obj.signer = message.signer);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgChannelOpenAck();
        message.portId = object.portId ?? "";
        message.channelId = object.channelId ?? "";
        message.counterpartyChannelId = object.counterpartyChannelId ?? "";
        message.counterpartyVersion = object.counterpartyVersion ?? "";
        message.proofTry = object.proofTry ?? new Uint8Array();
        if (object.proofHeight !== undefined && object.proofHeight !== null) {
            message.proofHeight = client_1.Height.fromPartial(object.proofHeight);
        }
        message.signer = object.signer ?? "";
        return message;
    },
};
function createBaseMsgChannelOpenAckResponse() {
    return {};
}
exports.MsgChannelOpenAckResponse = {
    typeUrl: "/ibc.core.channel.v1.MsgChannelOpenAckResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgChannelOpenAckResponse();
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
        const obj = createBaseMsgChannelOpenAckResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgChannelOpenAckResponse();
        return message;
    },
};
function createBaseMsgChannelOpenConfirm() {
    return {
        portId: "",
        channelId: "",
        proofAck: new Uint8Array(),
        proofHeight: client_1.Height.fromPartial({}),
        signer: "",
    };
}
exports.MsgChannelOpenConfirm = {
    typeUrl: "/ibc.core.channel.v1.MsgChannelOpenConfirm",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.portId !== "") {
            writer.uint32(10).string(message.portId);
        }
        if (message.channelId !== "") {
            writer.uint32(18).string(message.channelId);
        }
        if (message.proofAck.length !== 0) {
            writer.uint32(26).bytes(message.proofAck);
        }
        if (message.proofHeight !== undefined) {
            client_1.Height.encode(message.proofHeight, writer.uint32(34).fork()).ldelim();
        }
        if (message.signer !== "") {
            writer.uint32(42).string(message.signer);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgChannelOpenConfirm();
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
                    message.proofAck = reader.bytes();
                    break;
                case 4:
                    message.proofHeight = client_1.Height.decode(reader, reader.uint32());
                    break;
                case 5:
                    message.signer = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgChannelOpenConfirm();
        if ((0, helpers_1.isSet)(object.portId))
            obj.portId = String(object.portId);
        if ((0, helpers_1.isSet)(object.channelId))
            obj.channelId = String(object.channelId);
        if ((0, helpers_1.isSet)(object.proofAck))
            obj.proofAck = (0, helpers_1.bytesFromBase64)(object.proofAck);
        if ((0, helpers_1.isSet)(object.proofHeight))
            obj.proofHeight = client_1.Height.fromJSON(object.proofHeight);
        if ((0, helpers_1.isSet)(object.signer))
            obj.signer = String(object.signer);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.portId !== undefined && (obj.portId = message.portId);
        message.channelId !== undefined && (obj.channelId = message.channelId);
        message.proofAck !== undefined &&
            (obj.proofAck = (0, helpers_1.base64FromBytes)(message.proofAck !== undefined ? message.proofAck : new Uint8Array()));
        message.proofHeight !== undefined &&
            (obj.proofHeight = message.proofHeight ? client_1.Height.toJSON(message.proofHeight) : undefined);
        message.signer !== undefined && (obj.signer = message.signer);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgChannelOpenConfirm();
        message.portId = object.portId ?? "";
        message.channelId = object.channelId ?? "";
        message.proofAck = object.proofAck ?? new Uint8Array();
        if (object.proofHeight !== undefined && object.proofHeight !== null) {
            message.proofHeight = client_1.Height.fromPartial(object.proofHeight);
        }
        message.signer = object.signer ?? "";
        return message;
    },
};
function createBaseMsgChannelOpenConfirmResponse() {
    return {};
}
exports.MsgChannelOpenConfirmResponse = {
    typeUrl: "/ibc.core.channel.v1.MsgChannelOpenConfirmResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgChannelOpenConfirmResponse();
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
        const obj = createBaseMsgChannelOpenConfirmResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgChannelOpenConfirmResponse();
        return message;
    },
};
function createBaseMsgChannelCloseInit() {
    return {
        portId: "",
        channelId: "",
        signer: "",
    };
}
exports.MsgChannelCloseInit = {
    typeUrl: "/ibc.core.channel.v1.MsgChannelCloseInit",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.portId !== "") {
            writer.uint32(10).string(message.portId);
        }
        if (message.channelId !== "") {
            writer.uint32(18).string(message.channelId);
        }
        if (message.signer !== "") {
            writer.uint32(26).string(message.signer);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgChannelCloseInit();
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
                    message.signer = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgChannelCloseInit();
        if ((0, helpers_1.isSet)(object.portId))
            obj.portId = String(object.portId);
        if ((0, helpers_1.isSet)(object.channelId))
            obj.channelId = String(object.channelId);
        if ((0, helpers_1.isSet)(object.signer))
            obj.signer = String(object.signer);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.portId !== undefined && (obj.portId = message.portId);
        message.channelId !== undefined && (obj.channelId = message.channelId);
        message.signer !== undefined && (obj.signer = message.signer);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgChannelCloseInit();
        message.portId = object.portId ?? "";
        message.channelId = object.channelId ?? "";
        message.signer = object.signer ?? "";
        return message;
    },
};
function createBaseMsgChannelCloseInitResponse() {
    return {};
}
exports.MsgChannelCloseInitResponse = {
    typeUrl: "/ibc.core.channel.v1.MsgChannelCloseInitResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgChannelCloseInitResponse();
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
        const obj = createBaseMsgChannelCloseInitResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgChannelCloseInitResponse();
        return message;
    },
};
function createBaseMsgChannelCloseConfirm() {
    return {
        portId: "",
        channelId: "",
        proofInit: new Uint8Array(),
        proofHeight: client_1.Height.fromPartial({}),
        signer: "",
    };
}
exports.MsgChannelCloseConfirm = {
    typeUrl: "/ibc.core.channel.v1.MsgChannelCloseConfirm",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.portId !== "") {
            writer.uint32(10).string(message.portId);
        }
        if (message.channelId !== "") {
            writer.uint32(18).string(message.channelId);
        }
        if (message.proofInit.length !== 0) {
            writer.uint32(26).bytes(message.proofInit);
        }
        if (message.proofHeight !== undefined) {
            client_1.Height.encode(message.proofHeight, writer.uint32(34).fork()).ldelim();
        }
        if (message.signer !== "") {
            writer.uint32(42).string(message.signer);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgChannelCloseConfirm();
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
                    message.proofInit = reader.bytes();
                    break;
                case 4:
                    message.proofHeight = client_1.Height.decode(reader, reader.uint32());
                    break;
                case 5:
                    message.signer = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgChannelCloseConfirm();
        if ((0, helpers_1.isSet)(object.portId))
            obj.portId = String(object.portId);
        if ((0, helpers_1.isSet)(object.channelId))
            obj.channelId = String(object.channelId);
        if ((0, helpers_1.isSet)(object.proofInit))
            obj.proofInit = (0, helpers_1.bytesFromBase64)(object.proofInit);
        if ((0, helpers_1.isSet)(object.proofHeight))
            obj.proofHeight = client_1.Height.fromJSON(object.proofHeight);
        if ((0, helpers_1.isSet)(object.signer))
            obj.signer = String(object.signer);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.portId !== undefined && (obj.portId = message.portId);
        message.channelId !== undefined && (obj.channelId = message.channelId);
        message.proofInit !== undefined &&
            (obj.proofInit = (0, helpers_1.base64FromBytes)(message.proofInit !== undefined ? message.proofInit : new Uint8Array()));
        message.proofHeight !== undefined &&
            (obj.proofHeight = message.proofHeight ? client_1.Height.toJSON(message.proofHeight) : undefined);
        message.signer !== undefined && (obj.signer = message.signer);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgChannelCloseConfirm();
        message.portId = object.portId ?? "";
        message.channelId = object.channelId ?? "";
        message.proofInit = object.proofInit ?? new Uint8Array();
        if (object.proofHeight !== undefined && object.proofHeight !== null) {
            message.proofHeight = client_1.Height.fromPartial(object.proofHeight);
        }
        message.signer = object.signer ?? "";
        return message;
    },
};
function createBaseMsgChannelCloseConfirmResponse() {
    return {};
}
exports.MsgChannelCloseConfirmResponse = {
    typeUrl: "/ibc.core.channel.v1.MsgChannelCloseConfirmResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgChannelCloseConfirmResponse();
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
        const obj = createBaseMsgChannelCloseConfirmResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgChannelCloseConfirmResponse();
        return message;
    },
};
function createBaseMsgRecvPacket() {
    return {
        packet: channel_1.Packet.fromPartial({}),
        proofCommitment: new Uint8Array(),
        proofHeight: client_1.Height.fromPartial({}),
        signer: "",
    };
}
exports.MsgRecvPacket = {
    typeUrl: "/ibc.core.channel.v1.MsgRecvPacket",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.packet !== undefined) {
            channel_1.Packet.encode(message.packet, writer.uint32(10).fork()).ldelim();
        }
        if (message.proofCommitment.length !== 0) {
            writer.uint32(18).bytes(message.proofCommitment);
        }
        if (message.proofHeight !== undefined) {
            client_1.Height.encode(message.proofHeight, writer.uint32(26).fork()).ldelim();
        }
        if (message.signer !== "") {
            writer.uint32(34).string(message.signer);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgRecvPacket();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.packet = channel_1.Packet.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.proofCommitment = reader.bytes();
                    break;
                case 3:
                    message.proofHeight = client_1.Height.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.signer = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgRecvPacket();
        if ((0, helpers_1.isSet)(object.packet))
            obj.packet = channel_1.Packet.fromJSON(object.packet);
        if ((0, helpers_1.isSet)(object.proofCommitment))
            obj.proofCommitment = (0, helpers_1.bytesFromBase64)(object.proofCommitment);
        if ((0, helpers_1.isSet)(object.proofHeight))
            obj.proofHeight = client_1.Height.fromJSON(object.proofHeight);
        if ((0, helpers_1.isSet)(object.signer))
            obj.signer = String(object.signer);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.packet !== undefined && (obj.packet = message.packet ? channel_1.Packet.toJSON(message.packet) : undefined);
        message.proofCommitment !== undefined &&
            (obj.proofCommitment = (0, helpers_1.base64FromBytes)(message.proofCommitment !== undefined ? message.proofCommitment : new Uint8Array()));
        message.proofHeight !== undefined &&
            (obj.proofHeight = message.proofHeight ? client_1.Height.toJSON(message.proofHeight) : undefined);
        message.signer !== undefined && (obj.signer = message.signer);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgRecvPacket();
        if (object.packet !== undefined && object.packet !== null) {
            message.packet = channel_1.Packet.fromPartial(object.packet);
        }
        message.proofCommitment = object.proofCommitment ?? new Uint8Array();
        if (object.proofHeight !== undefined && object.proofHeight !== null) {
            message.proofHeight = client_1.Height.fromPartial(object.proofHeight);
        }
        message.signer = object.signer ?? "";
        return message;
    },
};
function createBaseMsgRecvPacketResponse() {
    return {
        result: 0,
    };
}
exports.MsgRecvPacketResponse = {
    typeUrl: "/ibc.core.channel.v1.MsgRecvPacketResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.result !== 0) {
            writer.uint32(8).int32(message.result);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgRecvPacketResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.result = reader.int32();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgRecvPacketResponse();
        if ((0, helpers_1.isSet)(object.result))
            obj.result = responseResultTypeFromJSON(object.result);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.result !== undefined && (obj.result = responseResultTypeToJSON(message.result));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgRecvPacketResponse();
        message.result = object.result ?? 0;
        return message;
    },
};
function createBaseMsgTimeout() {
    return {
        packet: channel_1.Packet.fromPartial({}),
        proofUnreceived: new Uint8Array(),
        proofHeight: client_1.Height.fromPartial({}),
        nextSequenceRecv: BigInt(0),
        signer: "",
    };
}
exports.MsgTimeout = {
    typeUrl: "/ibc.core.channel.v1.MsgTimeout",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.packet !== undefined) {
            channel_1.Packet.encode(message.packet, writer.uint32(10).fork()).ldelim();
        }
        if (message.proofUnreceived.length !== 0) {
            writer.uint32(18).bytes(message.proofUnreceived);
        }
        if (message.proofHeight !== undefined) {
            client_1.Height.encode(message.proofHeight, writer.uint32(26).fork()).ldelim();
        }
        if (message.nextSequenceRecv !== BigInt(0)) {
            writer.uint32(32).uint64(message.nextSequenceRecv);
        }
        if (message.signer !== "") {
            writer.uint32(42).string(message.signer);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgTimeout();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.packet = channel_1.Packet.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.proofUnreceived = reader.bytes();
                    break;
                case 3:
                    message.proofHeight = client_1.Height.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.nextSequenceRecv = reader.uint64();
                    break;
                case 5:
                    message.signer = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgTimeout();
        if ((0, helpers_1.isSet)(object.packet))
            obj.packet = channel_1.Packet.fromJSON(object.packet);
        if ((0, helpers_1.isSet)(object.proofUnreceived))
            obj.proofUnreceived = (0, helpers_1.bytesFromBase64)(object.proofUnreceived);
        if ((0, helpers_1.isSet)(object.proofHeight))
            obj.proofHeight = client_1.Height.fromJSON(object.proofHeight);
        if ((0, helpers_1.isSet)(object.nextSequenceRecv))
            obj.nextSequenceRecv = BigInt(object.nextSequenceRecv.toString());
        if ((0, helpers_1.isSet)(object.signer))
            obj.signer = String(object.signer);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.packet !== undefined && (obj.packet = message.packet ? channel_1.Packet.toJSON(message.packet) : undefined);
        message.proofUnreceived !== undefined &&
            (obj.proofUnreceived = (0, helpers_1.base64FromBytes)(message.proofUnreceived !== undefined ? message.proofUnreceived : new Uint8Array()));
        message.proofHeight !== undefined &&
            (obj.proofHeight = message.proofHeight ? client_1.Height.toJSON(message.proofHeight) : undefined);
        message.nextSequenceRecv !== undefined &&
            (obj.nextSequenceRecv = (message.nextSequenceRecv || BigInt(0)).toString());
        message.signer !== undefined && (obj.signer = message.signer);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgTimeout();
        if (object.packet !== undefined && object.packet !== null) {
            message.packet = channel_1.Packet.fromPartial(object.packet);
        }
        message.proofUnreceived = object.proofUnreceived ?? new Uint8Array();
        if (object.proofHeight !== undefined && object.proofHeight !== null) {
            message.proofHeight = client_1.Height.fromPartial(object.proofHeight);
        }
        if (object.nextSequenceRecv !== undefined && object.nextSequenceRecv !== null) {
            message.nextSequenceRecv = BigInt(object.nextSequenceRecv.toString());
        }
        message.signer = object.signer ?? "";
        return message;
    },
};
function createBaseMsgTimeoutResponse() {
    return {
        result: 0,
    };
}
exports.MsgTimeoutResponse = {
    typeUrl: "/ibc.core.channel.v1.MsgTimeoutResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.result !== 0) {
            writer.uint32(8).int32(message.result);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgTimeoutResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.result = reader.int32();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgTimeoutResponse();
        if ((0, helpers_1.isSet)(object.result))
            obj.result = responseResultTypeFromJSON(object.result);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.result !== undefined && (obj.result = responseResultTypeToJSON(message.result));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgTimeoutResponse();
        message.result = object.result ?? 0;
        return message;
    },
};
function createBaseMsgTimeoutOnClose() {
    return {
        packet: channel_1.Packet.fromPartial({}),
        proofUnreceived: new Uint8Array(),
        proofClose: new Uint8Array(),
        proofHeight: client_1.Height.fromPartial({}),
        nextSequenceRecv: BigInt(0),
        signer: "",
    };
}
exports.MsgTimeoutOnClose = {
    typeUrl: "/ibc.core.channel.v1.MsgTimeoutOnClose",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.packet !== undefined) {
            channel_1.Packet.encode(message.packet, writer.uint32(10).fork()).ldelim();
        }
        if (message.proofUnreceived.length !== 0) {
            writer.uint32(18).bytes(message.proofUnreceived);
        }
        if (message.proofClose.length !== 0) {
            writer.uint32(26).bytes(message.proofClose);
        }
        if (message.proofHeight !== undefined) {
            client_1.Height.encode(message.proofHeight, writer.uint32(34).fork()).ldelim();
        }
        if (message.nextSequenceRecv !== BigInt(0)) {
            writer.uint32(40).uint64(message.nextSequenceRecv);
        }
        if (message.signer !== "") {
            writer.uint32(50).string(message.signer);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgTimeoutOnClose();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.packet = channel_1.Packet.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.proofUnreceived = reader.bytes();
                    break;
                case 3:
                    message.proofClose = reader.bytes();
                    break;
                case 4:
                    message.proofHeight = client_1.Height.decode(reader, reader.uint32());
                    break;
                case 5:
                    message.nextSequenceRecv = reader.uint64();
                    break;
                case 6:
                    message.signer = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgTimeoutOnClose();
        if ((0, helpers_1.isSet)(object.packet))
            obj.packet = channel_1.Packet.fromJSON(object.packet);
        if ((0, helpers_1.isSet)(object.proofUnreceived))
            obj.proofUnreceived = (0, helpers_1.bytesFromBase64)(object.proofUnreceived);
        if ((0, helpers_1.isSet)(object.proofClose))
            obj.proofClose = (0, helpers_1.bytesFromBase64)(object.proofClose);
        if ((0, helpers_1.isSet)(object.proofHeight))
            obj.proofHeight = client_1.Height.fromJSON(object.proofHeight);
        if ((0, helpers_1.isSet)(object.nextSequenceRecv))
            obj.nextSequenceRecv = BigInt(object.nextSequenceRecv.toString());
        if ((0, helpers_1.isSet)(object.signer))
            obj.signer = String(object.signer);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.packet !== undefined && (obj.packet = message.packet ? channel_1.Packet.toJSON(message.packet) : undefined);
        message.proofUnreceived !== undefined &&
            (obj.proofUnreceived = (0, helpers_1.base64FromBytes)(message.proofUnreceived !== undefined ? message.proofUnreceived : new Uint8Array()));
        message.proofClose !== undefined &&
            (obj.proofClose = (0, helpers_1.base64FromBytes)(message.proofClose !== undefined ? message.proofClose : new Uint8Array()));
        message.proofHeight !== undefined &&
            (obj.proofHeight = message.proofHeight ? client_1.Height.toJSON(message.proofHeight) : undefined);
        message.nextSequenceRecv !== undefined &&
            (obj.nextSequenceRecv = (message.nextSequenceRecv || BigInt(0)).toString());
        message.signer !== undefined && (obj.signer = message.signer);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgTimeoutOnClose();
        if (object.packet !== undefined && object.packet !== null) {
            message.packet = channel_1.Packet.fromPartial(object.packet);
        }
        message.proofUnreceived = object.proofUnreceived ?? new Uint8Array();
        message.proofClose = object.proofClose ?? new Uint8Array();
        if (object.proofHeight !== undefined && object.proofHeight !== null) {
            message.proofHeight = client_1.Height.fromPartial(object.proofHeight);
        }
        if (object.nextSequenceRecv !== undefined && object.nextSequenceRecv !== null) {
            message.nextSequenceRecv = BigInt(object.nextSequenceRecv.toString());
        }
        message.signer = object.signer ?? "";
        return message;
    },
};
function createBaseMsgTimeoutOnCloseResponse() {
    return {
        result: 0,
    };
}
exports.MsgTimeoutOnCloseResponse = {
    typeUrl: "/ibc.core.channel.v1.MsgTimeoutOnCloseResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.result !== 0) {
            writer.uint32(8).int32(message.result);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgTimeoutOnCloseResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.result = reader.int32();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgTimeoutOnCloseResponse();
        if ((0, helpers_1.isSet)(object.result))
            obj.result = responseResultTypeFromJSON(object.result);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.result !== undefined && (obj.result = responseResultTypeToJSON(message.result));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgTimeoutOnCloseResponse();
        message.result = object.result ?? 0;
        return message;
    },
};
function createBaseMsgAcknowledgement() {
    return {
        packet: channel_1.Packet.fromPartial({}),
        acknowledgement: new Uint8Array(),
        proofAcked: new Uint8Array(),
        proofHeight: client_1.Height.fromPartial({}),
        signer: "",
    };
}
exports.MsgAcknowledgement = {
    typeUrl: "/ibc.core.channel.v1.MsgAcknowledgement",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.packet !== undefined) {
            channel_1.Packet.encode(message.packet, writer.uint32(10).fork()).ldelim();
        }
        if (message.acknowledgement.length !== 0) {
            writer.uint32(18).bytes(message.acknowledgement);
        }
        if (message.proofAcked.length !== 0) {
            writer.uint32(26).bytes(message.proofAcked);
        }
        if (message.proofHeight !== undefined) {
            client_1.Height.encode(message.proofHeight, writer.uint32(34).fork()).ldelim();
        }
        if (message.signer !== "") {
            writer.uint32(42).string(message.signer);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgAcknowledgement();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.packet = channel_1.Packet.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.acknowledgement = reader.bytes();
                    break;
                case 3:
                    message.proofAcked = reader.bytes();
                    break;
                case 4:
                    message.proofHeight = client_1.Height.decode(reader, reader.uint32());
                    break;
                case 5:
                    message.signer = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgAcknowledgement();
        if ((0, helpers_1.isSet)(object.packet))
            obj.packet = channel_1.Packet.fromJSON(object.packet);
        if ((0, helpers_1.isSet)(object.acknowledgement))
            obj.acknowledgement = (0, helpers_1.bytesFromBase64)(object.acknowledgement);
        if ((0, helpers_1.isSet)(object.proofAcked))
            obj.proofAcked = (0, helpers_1.bytesFromBase64)(object.proofAcked);
        if ((0, helpers_1.isSet)(object.proofHeight))
            obj.proofHeight = client_1.Height.fromJSON(object.proofHeight);
        if ((0, helpers_1.isSet)(object.signer))
            obj.signer = String(object.signer);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.packet !== undefined && (obj.packet = message.packet ? channel_1.Packet.toJSON(message.packet) : undefined);
        message.acknowledgement !== undefined &&
            (obj.acknowledgement = (0, helpers_1.base64FromBytes)(message.acknowledgement !== undefined ? message.acknowledgement : new Uint8Array()));
        message.proofAcked !== undefined &&
            (obj.proofAcked = (0, helpers_1.base64FromBytes)(message.proofAcked !== undefined ? message.proofAcked : new Uint8Array()));
        message.proofHeight !== undefined &&
            (obj.proofHeight = message.proofHeight ? client_1.Height.toJSON(message.proofHeight) : undefined);
        message.signer !== undefined && (obj.signer = message.signer);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgAcknowledgement();
        if (object.packet !== undefined && object.packet !== null) {
            message.packet = channel_1.Packet.fromPartial(object.packet);
        }
        message.acknowledgement = object.acknowledgement ?? new Uint8Array();
        message.proofAcked = object.proofAcked ?? new Uint8Array();
        if (object.proofHeight !== undefined && object.proofHeight !== null) {
            message.proofHeight = client_1.Height.fromPartial(object.proofHeight);
        }
        message.signer = object.signer ?? "";
        return message;
    },
};
function createBaseMsgAcknowledgementResponse() {
    return {
        result: 0,
    };
}
exports.MsgAcknowledgementResponse = {
    typeUrl: "/ibc.core.channel.v1.MsgAcknowledgementResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.result !== 0) {
            writer.uint32(8).int32(message.result);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgAcknowledgementResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.result = reader.int32();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgAcknowledgementResponse();
        if ((0, helpers_1.isSet)(object.result))
            obj.result = responseResultTypeFromJSON(object.result);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.result !== undefined && (obj.result = responseResultTypeToJSON(message.result));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgAcknowledgementResponse();
        message.result = object.result ?? 0;
        return message;
    },
};
class MsgClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
        this.ChannelOpenInit = this.ChannelOpenInit.bind(this);
        this.ChannelOpenTry = this.ChannelOpenTry.bind(this);
        this.ChannelOpenAck = this.ChannelOpenAck.bind(this);
        this.ChannelOpenConfirm = this.ChannelOpenConfirm.bind(this);
        this.ChannelCloseInit = this.ChannelCloseInit.bind(this);
        this.ChannelCloseConfirm = this.ChannelCloseConfirm.bind(this);
        this.RecvPacket = this.RecvPacket.bind(this);
        this.Timeout = this.Timeout.bind(this);
        this.TimeoutOnClose = this.TimeoutOnClose.bind(this);
        this.Acknowledgement = this.Acknowledgement.bind(this);
    }
    ChannelOpenInit(request) {
        const data = exports.MsgChannelOpenInit.encode(request).finish();
        const promise = this.rpc.request("ibc.core.channel.v1.Msg", "ChannelOpenInit", data);
        return promise.then((data) => exports.MsgChannelOpenInitResponse.decode(new binary_1.BinaryReader(data)));
    }
    ChannelOpenTry(request) {
        const data = exports.MsgChannelOpenTry.encode(request).finish();
        const promise = this.rpc.request("ibc.core.channel.v1.Msg", "ChannelOpenTry", data);
        return promise.then((data) => exports.MsgChannelOpenTryResponse.decode(new binary_1.BinaryReader(data)));
    }
    ChannelOpenAck(request) {
        const data = exports.MsgChannelOpenAck.encode(request).finish();
        const promise = this.rpc.request("ibc.core.channel.v1.Msg", "ChannelOpenAck", data);
        return promise.then((data) => exports.MsgChannelOpenAckResponse.decode(new binary_1.BinaryReader(data)));
    }
    ChannelOpenConfirm(request) {
        const data = exports.MsgChannelOpenConfirm.encode(request).finish();
        const promise = this.rpc.request("ibc.core.channel.v1.Msg", "ChannelOpenConfirm", data);
        return promise.then((data) => exports.MsgChannelOpenConfirmResponse.decode(new binary_1.BinaryReader(data)));
    }
    ChannelCloseInit(request) {
        const data = exports.MsgChannelCloseInit.encode(request).finish();
        const promise = this.rpc.request("ibc.core.channel.v1.Msg", "ChannelCloseInit", data);
        return promise.then((data) => exports.MsgChannelCloseInitResponse.decode(new binary_1.BinaryReader(data)));
    }
    ChannelCloseConfirm(request) {
        const data = exports.MsgChannelCloseConfirm.encode(request).finish();
        const promise = this.rpc.request("ibc.core.channel.v1.Msg", "ChannelCloseConfirm", data);
        return promise.then((data) => exports.MsgChannelCloseConfirmResponse.decode(new binary_1.BinaryReader(data)));
    }
    RecvPacket(request) {
        const data = exports.MsgRecvPacket.encode(request).finish();
        const promise = this.rpc.request("ibc.core.channel.v1.Msg", "RecvPacket", data);
        return promise.then((data) => exports.MsgRecvPacketResponse.decode(new binary_1.BinaryReader(data)));
    }
    Timeout(request) {
        const data = exports.MsgTimeout.encode(request).finish();
        const promise = this.rpc.request("ibc.core.channel.v1.Msg", "Timeout", data);
        return promise.then((data) => exports.MsgTimeoutResponse.decode(new binary_1.BinaryReader(data)));
    }
    TimeoutOnClose(request) {
        const data = exports.MsgTimeoutOnClose.encode(request).finish();
        const promise = this.rpc.request("ibc.core.channel.v1.Msg", "TimeoutOnClose", data);
        return promise.then((data) => exports.MsgTimeoutOnCloseResponse.decode(new binary_1.BinaryReader(data)));
    }
    Acknowledgement(request) {
        const data = exports.MsgAcknowledgement.encode(request).finish();
        const promise = this.rpc.request("ibc.core.channel.v1.Msg", "Acknowledgement", data);
        return promise.then((data) => exports.MsgAcknowledgementResponse.decode(new binary_1.BinaryReader(data)));
    }
}
exports.MsgClientImpl = MsgClientImpl;
//# sourceMappingURL=tx.js.map