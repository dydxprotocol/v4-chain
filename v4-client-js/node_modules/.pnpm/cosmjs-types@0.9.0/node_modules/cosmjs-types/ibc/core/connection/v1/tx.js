"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.MsgClientImpl = exports.MsgConnectionOpenConfirmResponse = exports.MsgConnectionOpenConfirm = exports.MsgConnectionOpenAckResponse = exports.MsgConnectionOpenAck = exports.MsgConnectionOpenTryResponse = exports.MsgConnectionOpenTry = exports.MsgConnectionOpenInitResponse = exports.MsgConnectionOpenInit = exports.protobufPackage = void 0;
/* eslint-disable */
const connection_1 = require("./connection");
const any_1 = require("../../../../google/protobuf/any");
const client_1 = require("../../client/v1/client");
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "ibc.core.connection.v1";
function createBaseMsgConnectionOpenInit() {
    return {
        clientId: "",
        counterparty: connection_1.Counterparty.fromPartial({}),
        version: undefined,
        delayPeriod: BigInt(0),
        signer: "",
    };
}
exports.MsgConnectionOpenInit = {
    typeUrl: "/ibc.core.connection.v1.MsgConnectionOpenInit",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.clientId !== "") {
            writer.uint32(10).string(message.clientId);
        }
        if (message.counterparty !== undefined) {
            connection_1.Counterparty.encode(message.counterparty, writer.uint32(18).fork()).ldelim();
        }
        if (message.version !== undefined) {
            connection_1.Version.encode(message.version, writer.uint32(26).fork()).ldelim();
        }
        if (message.delayPeriod !== BigInt(0)) {
            writer.uint32(32).uint64(message.delayPeriod);
        }
        if (message.signer !== "") {
            writer.uint32(42).string(message.signer);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgConnectionOpenInit();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.clientId = reader.string();
                    break;
                case 2:
                    message.counterparty = connection_1.Counterparty.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.version = connection_1.Version.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.delayPeriod = reader.uint64();
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
        const obj = createBaseMsgConnectionOpenInit();
        if ((0, helpers_1.isSet)(object.clientId))
            obj.clientId = String(object.clientId);
        if ((0, helpers_1.isSet)(object.counterparty))
            obj.counterparty = connection_1.Counterparty.fromJSON(object.counterparty);
        if ((0, helpers_1.isSet)(object.version))
            obj.version = connection_1.Version.fromJSON(object.version);
        if ((0, helpers_1.isSet)(object.delayPeriod))
            obj.delayPeriod = BigInt(object.delayPeriod.toString());
        if ((0, helpers_1.isSet)(object.signer))
            obj.signer = String(object.signer);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.clientId !== undefined && (obj.clientId = message.clientId);
        message.counterparty !== undefined &&
            (obj.counterparty = message.counterparty ? connection_1.Counterparty.toJSON(message.counterparty) : undefined);
        message.version !== undefined &&
            (obj.version = message.version ? connection_1.Version.toJSON(message.version) : undefined);
        message.delayPeriod !== undefined && (obj.delayPeriod = (message.delayPeriod || BigInt(0)).toString());
        message.signer !== undefined && (obj.signer = message.signer);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgConnectionOpenInit();
        message.clientId = object.clientId ?? "";
        if (object.counterparty !== undefined && object.counterparty !== null) {
            message.counterparty = connection_1.Counterparty.fromPartial(object.counterparty);
        }
        if (object.version !== undefined && object.version !== null) {
            message.version = connection_1.Version.fromPartial(object.version);
        }
        if (object.delayPeriod !== undefined && object.delayPeriod !== null) {
            message.delayPeriod = BigInt(object.delayPeriod.toString());
        }
        message.signer = object.signer ?? "";
        return message;
    },
};
function createBaseMsgConnectionOpenInitResponse() {
    return {};
}
exports.MsgConnectionOpenInitResponse = {
    typeUrl: "/ibc.core.connection.v1.MsgConnectionOpenInitResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgConnectionOpenInitResponse();
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
        const obj = createBaseMsgConnectionOpenInitResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgConnectionOpenInitResponse();
        return message;
    },
};
function createBaseMsgConnectionOpenTry() {
    return {
        clientId: "",
        previousConnectionId: "",
        clientState: undefined,
        counterparty: connection_1.Counterparty.fromPartial({}),
        delayPeriod: BigInt(0),
        counterpartyVersions: [],
        proofHeight: client_1.Height.fromPartial({}),
        proofInit: new Uint8Array(),
        proofClient: new Uint8Array(),
        proofConsensus: new Uint8Array(),
        consensusHeight: client_1.Height.fromPartial({}),
        signer: "",
        hostConsensusStateProof: new Uint8Array(),
    };
}
exports.MsgConnectionOpenTry = {
    typeUrl: "/ibc.core.connection.v1.MsgConnectionOpenTry",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.clientId !== "") {
            writer.uint32(10).string(message.clientId);
        }
        if (message.previousConnectionId !== "") {
            writer.uint32(18).string(message.previousConnectionId);
        }
        if (message.clientState !== undefined) {
            any_1.Any.encode(message.clientState, writer.uint32(26).fork()).ldelim();
        }
        if (message.counterparty !== undefined) {
            connection_1.Counterparty.encode(message.counterparty, writer.uint32(34).fork()).ldelim();
        }
        if (message.delayPeriod !== BigInt(0)) {
            writer.uint32(40).uint64(message.delayPeriod);
        }
        for (const v of message.counterpartyVersions) {
            connection_1.Version.encode(v, writer.uint32(50).fork()).ldelim();
        }
        if (message.proofHeight !== undefined) {
            client_1.Height.encode(message.proofHeight, writer.uint32(58).fork()).ldelim();
        }
        if (message.proofInit.length !== 0) {
            writer.uint32(66).bytes(message.proofInit);
        }
        if (message.proofClient.length !== 0) {
            writer.uint32(74).bytes(message.proofClient);
        }
        if (message.proofConsensus.length !== 0) {
            writer.uint32(82).bytes(message.proofConsensus);
        }
        if (message.consensusHeight !== undefined) {
            client_1.Height.encode(message.consensusHeight, writer.uint32(90).fork()).ldelim();
        }
        if (message.signer !== "") {
            writer.uint32(98).string(message.signer);
        }
        if (message.hostConsensusStateProof.length !== 0) {
            writer.uint32(106).bytes(message.hostConsensusStateProof);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgConnectionOpenTry();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.clientId = reader.string();
                    break;
                case 2:
                    message.previousConnectionId = reader.string();
                    break;
                case 3:
                    message.clientState = any_1.Any.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.counterparty = connection_1.Counterparty.decode(reader, reader.uint32());
                    break;
                case 5:
                    message.delayPeriod = reader.uint64();
                    break;
                case 6:
                    message.counterpartyVersions.push(connection_1.Version.decode(reader, reader.uint32()));
                    break;
                case 7:
                    message.proofHeight = client_1.Height.decode(reader, reader.uint32());
                    break;
                case 8:
                    message.proofInit = reader.bytes();
                    break;
                case 9:
                    message.proofClient = reader.bytes();
                    break;
                case 10:
                    message.proofConsensus = reader.bytes();
                    break;
                case 11:
                    message.consensusHeight = client_1.Height.decode(reader, reader.uint32());
                    break;
                case 12:
                    message.signer = reader.string();
                    break;
                case 13:
                    message.hostConsensusStateProof = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgConnectionOpenTry();
        if ((0, helpers_1.isSet)(object.clientId))
            obj.clientId = String(object.clientId);
        if ((0, helpers_1.isSet)(object.previousConnectionId))
            obj.previousConnectionId = String(object.previousConnectionId);
        if ((0, helpers_1.isSet)(object.clientState))
            obj.clientState = any_1.Any.fromJSON(object.clientState);
        if ((0, helpers_1.isSet)(object.counterparty))
            obj.counterparty = connection_1.Counterparty.fromJSON(object.counterparty);
        if ((0, helpers_1.isSet)(object.delayPeriod))
            obj.delayPeriod = BigInt(object.delayPeriod.toString());
        if (Array.isArray(object?.counterpartyVersions))
            obj.counterpartyVersions = object.counterpartyVersions.map((e) => connection_1.Version.fromJSON(e));
        if ((0, helpers_1.isSet)(object.proofHeight))
            obj.proofHeight = client_1.Height.fromJSON(object.proofHeight);
        if ((0, helpers_1.isSet)(object.proofInit))
            obj.proofInit = (0, helpers_1.bytesFromBase64)(object.proofInit);
        if ((0, helpers_1.isSet)(object.proofClient))
            obj.proofClient = (0, helpers_1.bytesFromBase64)(object.proofClient);
        if ((0, helpers_1.isSet)(object.proofConsensus))
            obj.proofConsensus = (0, helpers_1.bytesFromBase64)(object.proofConsensus);
        if ((0, helpers_1.isSet)(object.consensusHeight))
            obj.consensusHeight = client_1.Height.fromJSON(object.consensusHeight);
        if ((0, helpers_1.isSet)(object.signer))
            obj.signer = String(object.signer);
        if ((0, helpers_1.isSet)(object.hostConsensusStateProof))
            obj.hostConsensusStateProof = (0, helpers_1.bytesFromBase64)(object.hostConsensusStateProof);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.clientId !== undefined && (obj.clientId = message.clientId);
        message.previousConnectionId !== undefined && (obj.previousConnectionId = message.previousConnectionId);
        message.clientState !== undefined &&
            (obj.clientState = message.clientState ? any_1.Any.toJSON(message.clientState) : undefined);
        message.counterparty !== undefined &&
            (obj.counterparty = message.counterparty ? connection_1.Counterparty.toJSON(message.counterparty) : undefined);
        message.delayPeriod !== undefined && (obj.delayPeriod = (message.delayPeriod || BigInt(0)).toString());
        if (message.counterpartyVersions) {
            obj.counterpartyVersions = message.counterpartyVersions.map((e) => (e ? connection_1.Version.toJSON(e) : undefined));
        }
        else {
            obj.counterpartyVersions = [];
        }
        message.proofHeight !== undefined &&
            (obj.proofHeight = message.proofHeight ? client_1.Height.toJSON(message.proofHeight) : undefined);
        message.proofInit !== undefined &&
            (obj.proofInit = (0, helpers_1.base64FromBytes)(message.proofInit !== undefined ? message.proofInit : new Uint8Array()));
        message.proofClient !== undefined &&
            (obj.proofClient = (0, helpers_1.base64FromBytes)(message.proofClient !== undefined ? message.proofClient : new Uint8Array()));
        message.proofConsensus !== undefined &&
            (obj.proofConsensus = (0, helpers_1.base64FromBytes)(message.proofConsensus !== undefined ? message.proofConsensus : new Uint8Array()));
        message.consensusHeight !== undefined &&
            (obj.consensusHeight = message.consensusHeight ? client_1.Height.toJSON(message.consensusHeight) : undefined);
        message.signer !== undefined && (obj.signer = message.signer);
        message.hostConsensusStateProof !== undefined &&
            (obj.hostConsensusStateProof = (0, helpers_1.base64FromBytes)(message.hostConsensusStateProof !== undefined ? message.hostConsensusStateProof : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgConnectionOpenTry();
        message.clientId = object.clientId ?? "";
        message.previousConnectionId = object.previousConnectionId ?? "";
        if (object.clientState !== undefined && object.clientState !== null) {
            message.clientState = any_1.Any.fromPartial(object.clientState);
        }
        if (object.counterparty !== undefined && object.counterparty !== null) {
            message.counterparty = connection_1.Counterparty.fromPartial(object.counterparty);
        }
        if (object.delayPeriod !== undefined && object.delayPeriod !== null) {
            message.delayPeriod = BigInt(object.delayPeriod.toString());
        }
        message.counterpartyVersions = object.counterpartyVersions?.map((e) => connection_1.Version.fromPartial(e)) || [];
        if (object.proofHeight !== undefined && object.proofHeight !== null) {
            message.proofHeight = client_1.Height.fromPartial(object.proofHeight);
        }
        message.proofInit = object.proofInit ?? new Uint8Array();
        message.proofClient = object.proofClient ?? new Uint8Array();
        message.proofConsensus = object.proofConsensus ?? new Uint8Array();
        if (object.consensusHeight !== undefined && object.consensusHeight !== null) {
            message.consensusHeight = client_1.Height.fromPartial(object.consensusHeight);
        }
        message.signer = object.signer ?? "";
        message.hostConsensusStateProof = object.hostConsensusStateProof ?? new Uint8Array();
        return message;
    },
};
function createBaseMsgConnectionOpenTryResponse() {
    return {};
}
exports.MsgConnectionOpenTryResponse = {
    typeUrl: "/ibc.core.connection.v1.MsgConnectionOpenTryResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgConnectionOpenTryResponse();
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
        const obj = createBaseMsgConnectionOpenTryResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgConnectionOpenTryResponse();
        return message;
    },
};
function createBaseMsgConnectionOpenAck() {
    return {
        connectionId: "",
        counterpartyConnectionId: "",
        version: undefined,
        clientState: undefined,
        proofHeight: client_1.Height.fromPartial({}),
        proofTry: new Uint8Array(),
        proofClient: new Uint8Array(),
        proofConsensus: new Uint8Array(),
        consensusHeight: client_1.Height.fromPartial({}),
        signer: "",
        hostConsensusStateProof: new Uint8Array(),
    };
}
exports.MsgConnectionOpenAck = {
    typeUrl: "/ibc.core.connection.v1.MsgConnectionOpenAck",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.connectionId !== "") {
            writer.uint32(10).string(message.connectionId);
        }
        if (message.counterpartyConnectionId !== "") {
            writer.uint32(18).string(message.counterpartyConnectionId);
        }
        if (message.version !== undefined) {
            connection_1.Version.encode(message.version, writer.uint32(26).fork()).ldelim();
        }
        if (message.clientState !== undefined) {
            any_1.Any.encode(message.clientState, writer.uint32(34).fork()).ldelim();
        }
        if (message.proofHeight !== undefined) {
            client_1.Height.encode(message.proofHeight, writer.uint32(42).fork()).ldelim();
        }
        if (message.proofTry.length !== 0) {
            writer.uint32(50).bytes(message.proofTry);
        }
        if (message.proofClient.length !== 0) {
            writer.uint32(58).bytes(message.proofClient);
        }
        if (message.proofConsensus.length !== 0) {
            writer.uint32(66).bytes(message.proofConsensus);
        }
        if (message.consensusHeight !== undefined) {
            client_1.Height.encode(message.consensusHeight, writer.uint32(74).fork()).ldelim();
        }
        if (message.signer !== "") {
            writer.uint32(82).string(message.signer);
        }
        if (message.hostConsensusStateProof.length !== 0) {
            writer.uint32(90).bytes(message.hostConsensusStateProof);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgConnectionOpenAck();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.connectionId = reader.string();
                    break;
                case 2:
                    message.counterpartyConnectionId = reader.string();
                    break;
                case 3:
                    message.version = connection_1.Version.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.clientState = any_1.Any.decode(reader, reader.uint32());
                    break;
                case 5:
                    message.proofHeight = client_1.Height.decode(reader, reader.uint32());
                    break;
                case 6:
                    message.proofTry = reader.bytes();
                    break;
                case 7:
                    message.proofClient = reader.bytes();
                    break;
                case 8:
                    message.proofConsensus = reader.bytes();
                    break;
                case 9:
                    message.consensusHeight = client_1.Height.decode(reader, reader.uint32());
                    break;
                case 10:
                    message.signer = reader.string();
                    break;
                case 11:
                    message.hostConsensusStateProof = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgConnectionOpenAck();
        if ((0, helpers_1.isSet)(object.connectionId))
            obj.connectionId = String(object.connectionId);
        if ((0, helpers_1.isSet)(object.counterpartyConnectionId))
            obj.counterpartyConnectionId = String(object.counterpartyConnectionId);
        if ((0, helpers_1.isSet)(object.version))
            obj.version = connection_1.Version.fromJSON(object.version);
        if ((0, helpers_1.isSet)(object.clientState))
            obj.clientState = any_1.Any.fromJSON(object.clientState);
        if ((0, helpers_1.isSet)(object.proofHeight))
            obj.proofHeight = client_1.Height.fromJSON(object.proofHeight);
        if ((0, helpers_1.isSet)(object.proofTry))
            obj.proofTry = (0, helpers_1.bytesFromBase64)(object.proofTry);
        if ((0, helpers_1.isSet)(object.proofClient))
            obj.proofClient = (0, helpers_1.bytesFromBase64)(object.proofClient);
        if ((0, helpers_1.isSet)(object.proofConsensus))
            obj.proofConsensus = (0, helpers_1.bytesFromBase64)(object.proofConsensus);
        if ((0, helpers_1.isSet)(object.consensusHeight))
            obj.consensusHeight = client_1.Height.fromJSON(object.consensusHeight);
        if ((0, helpers_1.isSet)(object.signer))
            obj.signer = String(object.signer);
        if ((0, helpers_1.isSet)(object.hostConsensusStateProof))
            obj.hostConsensusStateProof = (0, helpers_1.bytesFromBase64)(object.hostConsensusStateProof);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.connectionId !== undefined && (obj.connectionId = message.connectionId);
        message.counterpartyConnectionId !== undefined &&
            (obj.counterpartyConnectionId = message.counterpartyConnectionId);
        message.version !== undefined &&
            (obj.version = message.version ? connection_1.Version.toJSON(message.version) : undefined);
        message.clientState !== undefined &&
            (obj.clientState = message.clientState ? any_1.Any.toJSON(message.clientState) : undefined);
        message.proofHeight !== undefined &&
            (obj.proofHeight = message.proofHeight ? client_1.Height.toJSON(message.proofHeight) : undefined);
        message.proofTry !== undefined &&
            (obj.proofTry = (0, helpers_1.base64FromBytes)(message.proofTry !== undefined ? message.proofTry : new Uint8Array()));
        message.proofClient !== undefined &&
            (obj.proofClient = (0, helpers_1.base64FromBytes)(message.proofClient !== undefined ? message.proofClient : new Uint8Array()));
        message.proofConsensus !== undefined &&
            (obj.proofConsensus = (0, helpers_1.base64FromBytes)(message.proofConsensus !== undefined ? message.proofConsensus : new Uint8Array()));
        message.consensusHeight !== undefined &&
            (obj.consensusHeight = message.consensusHeight ? client_1.Height.toJSON(message.consensusHeight) : undefined);
        message.signer !== undefined && (obj.signer = message.signer);
        message.hostConsensusStateProof !== undefined &&
            (obj.hostConsensusStateProof = (0, helpers_1.base64FromBytes)(message.hostConsensusStateProof !== undefined ? message.hostConsensusStateProof : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgConnectionOpenAck();
        message.connectionId = object.connectionId ?? "";
        message.counterpartyConnectionId = object.counterpartyConnectionId ?? "";
        if (object.version !== undefined && object.version !== null) {
            message.version = connection_1.Version.fromPartial(object.version);
        }
        if (object.clientState !== undefined && object.clientState !== null) {
            message.clientState = any_1.Any.fromPartial(object.clientState);
        }
        if (object.proofHeight !== undefined && object.proofHeight !== null) {
            message.proofHeight = client_1.Height.fromPartial(object.proofHeight);
        }
        message.proofTry = object.proofTry ?? new Uint8Array();
        message.proofClient = object.proofClient ?? new Uint8Array();
        message.proofConsensus = object.proofConsensus ?? new Uint8Array();
        if (object.consensusHeight !== undefined && object.consensusHeight !== null) {
            message.consensusHeight = client_1.Height.fromPartial(object.consensusHeight);
        }
        message.signer = object.signer ?? "";
        message.hostConsensusStateProof = object.hostConsensusStateProof ?? new Uint8Array();
        return message;
    },
};
function createBaseMsgConnectionOpenAckResponse() {
    return {};
}
exports.MsgConnectionOpenAckResponse = {
    typeUrl: "/ibc.core.connection.v1.MsgConnectionOpenAckResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgConnectionOpenAckResponse();
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
        const obj = createBaseMsgConnectionOpenAckResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgConnectionOpenAckResponse();
        return message;
    },
};
function createBaseMsgConnectionOpenConfirm() {
    return {
        connectionId: "",
        proofAck: new Uint8Array(),
        proofHeight: client_1.Height.fromPartial({}),
        signer: "",
    };
}
exports.MsgConnectionOpenConfirm = {
    typeUrl: "/ibc.core.connection.v1.MsgConnectionOpenConfirm",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.connectionId !== "") {
            writer.uint32(10).string(message.connectionId);
        }
        if (message.proofAck.length !== 0) {
            writer.uint32(18).bytes(message.proofAck);
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
        const message = createBaseMsgConnectionOpenConfirm();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.connectionId = reader.string();
                    break;
                case 2:
                    message.proofAck = reader.bytes();
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
        const obj = createBaseMsgConnectionOpenConfirm();
        if ((0, helpers_1.isSet)(object.connectionId))
            obj.connectionId = String(object.connectionId);
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
        message.connectionId !== undefined && (obj.connectionId = message.connectionId);
        message.proofAck !== undefined &&
            (obj.proofAck = (0, helpers_1.base64FromBytes)(message.proofAck !== undefined ? message.proofAck : new Uint8Array()));
        message.proofHeight !== undefined &&
            (obj.proofHeight = message.proofHeight ? client_1.Height.toJSON(message.proofHeight) : undefined);
        message.signer !== undefined && (obj.signer = message.signer);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgConnectionOpenConfirm();
        message.connectionId = object.connectionId ?? "";
        message.proofAck = object.proofAck ?? new Uint8Array();
        if (object.proofHeight !== undefined && object.proofHeight !== null) {
            message.proofHeight = client_1.Height.fromPartial(object.proofHeight);
        }
        message.signer = object.signer ?? "";
        return message;
    },
};
function createBaseMsgConnectionOpenConfirmResponse() {
    return {};
}
exports.MsgConnectionOpenConfirmResponse = {
    typeUrl: "/ibc.core.connection.v1.MsgConnectionOpenConfirmResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgConnectionOpenConfirmResponse();
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
        const obj = createBaseMsgConnectionOpenConfirmResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgConnectionOpenConfirmResponse();
        return message;
    },
};
class MsgClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
        this.ConnectionOpenInit = this.ConnectionOpenInit.bind(this);
        this.ConnectionOpenTry = this.ConnectionOpenTry.bind(this);
        this.ConnectionOpenAck = this.ConnectionOpenAck.bind(this);
        this.ConnectionOpenConfirm = this.ConnectionOpenConfirm.bind(this);
    }
    ConnectionOpenInit(request) {
        const data = exports.MsgConnectionOpenInit.encode(request).finish();
        const promise = this.rpc.request("ibc.core.connection.v1.Msg", "ConnectionOpenInit", data);
        return promise.then((data) => exports.MsgConnectionOpenInitResponse.decode(new binary_1.BinaryReader(data)));
    }
    ConnectionOpenTry(request) {
        const data = exports.MsgConnectionOpenTry.encode(request).finish();
        const promise = this.rpc.request("ibc.core.connection.v1.Msg", "ConnectionOpenTry", data);
        return promise.then((data) => exports.MsgConnectionOpenTryResponse.decode(new binary_1.BinaryReader(data)));
    }
    ConnectionOpenAck(request) {
        const data = exports.MsgConnectionOpenAck.encode(request).finish();
        const promise = this.rpc.request("ibc.core.connection.v1.Msg", "ConnectionOpenAck", data);
        return promise.then((data) => exports.MsgConnectionOpenAckResponse.decode(new binary_1.BinaryReader(data)));
    }
    ConnectionOpenConfirm(request) {
        const data = exports.MsgConnectionOpenConfirm.encode(request).finish();
        const promise = this.rpc.request("ibc.core.connection.v1.Msg", "ConnectionOpenConfirm", data);
        return promise.then((data) => exports.MsgConnectionOpenConfirmResponse.decode(new binary_1.BinaryReader(data)));
    }
}
exports.MsgClientImpl = MsgClientImpl;
//# sourceMappingURL=tx.js.map