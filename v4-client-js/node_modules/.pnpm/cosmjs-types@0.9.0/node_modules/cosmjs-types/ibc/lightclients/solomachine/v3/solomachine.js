"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.HeaderData = exports.SignBytes = exports.TimestampedSignatureData = exports.SignatureAndData = exports.Misbehaviour = exports.Header = exports.ConsensusState = exports.ClientState = exports.protobufPackage = void 0;
/* eslint-disable */
const any_1 = require("../../../../google/protobuf/any");
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "ibc.lightclients.solomachine.v3";
function createBaseClientState() {
    return {
        sequence: BigInt(0),
        isFrozen: false,
        consensusState: undefined,
    };
}
exports.ClientState = {
    typeUrl: "/ibc.lightclients.solomachine.v3.ClientState",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.sequence !== BigInt(0)) {
            writer.uint32(8).uint64(message.sequence);
        }
        if (message.isFrozen === true) {
            writer.uint32(16).bool(message.isFrozen);
        }
        if (message.consensusState !== undefined) {
            exports.ConsensusState.encode(message.consensusState, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseClientState();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.sequence = reader.uint64();
                    break;
                case 2:
                    message.isFrozen = reader.bool();
                    break;
                case 3:
                    message.consensusState = exports.ConsensusState.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseClientState();
        if ((0, helpers_1.isSet)(object.sequence))
            obj.sequence = BigInt(object.sequence.toString());
        if ((0, helpers_1.isSet)(object.isFrozen))
            obj.isFrozen = Boolean(object.isFrozen);
        if ((0, helpers_1.isSet)(object.consensusState))
            obj.consensusState = exports.ConsensusState.fromJSON(object.consensusState);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.sequence !== undefined && (obj.sequence = (message.sequence || BigInt(0)).toString());
        message.isFrozen !== undefined && (obj.isFrozen = message.isFrozen);
        message.consensusState !== undefined &&
            (obj.consensusState = message.consensusState
                ? exports.ConsensusState.toJSON(message.consensusState)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseClientState();
        if (object.sequence !== undefined && object.sequence !== null) {
            message.sequence = BigInt(object.sequence.toString());
        }
        message.isFrozen = object.isFrozen ?? false;
        if (object.consensusState !== undefined && object.consensusState !== null) {
            message.consensusState = exports.ConsensusState.fromPartial(object.consensusState);
        }
        return message;
    },
};
function createBaseConsensusState() {
    return {
        publicKey: undefined,
        diversifier: "",
        timestamp: BigInt(0),
    };
}
exports.ConsensusState = {
    typeUrl: "/ibc.lightclients.solomachine.v3.ConsensusState",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.publicKey !== undefined) {
            any_1.Any.encode(message.publicKey, writer.uint32(10).fork()).ldelim();
        }
        if (message.diversifier !== "") {
            writer.uint32(18).string(message.diversifier);
        }
        if (message.timestamp !== BigInt(0)) {
            writer.uint32(24).uint64(message.timestamp);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseConsensusState();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.publicKey = any_1.Any.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.diversifier = reader.string();
                    break;
                case 3:
                    message.timestamp = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseConsensusState();
        if ((0, helpers_1.isSet)(object.publicKey))
            obj.publicKey = any_1.Any.fromJSON(object.publicKey);
        if ((0, helpers_1.isSet)(object.diversifier))
            obj.diversifier = String(object.diversifier);
        if ((0, helpers_1.isSet)(object.timestamp))
            obj.timestamp = BigInt(object.timestamp.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.publicKey !== undefined &&
            (obj.publicKey = message.publicKey ? any_1.Any.toJSON(message.publicKey) : undefined);
        message.diversifier !== undefined && (obj.diversifier = message.diversifier);
        message.timestamp !== undefined && (obj.timestamp = (message.timestamp || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseConsensusState();
        if (object.publicKey !== undefined && object.publicKey !== null) {
            message.publicKey = any_1.Any.fromPartial(object.publicKey);
        }
        message.diversifier = object.diversifier ?? "";
        if (object.timestamp !== undefined && object.timestamp !== null) {
            message.timestamp = BigInt(object.timestamp.toString());
        }
        return message;
    },
};
function createBaseHeader() {
    return {
        timestamp: BigInt(0),
        signature: new Uint8Array(),
        newPublicKey: undefined,
        newDiversifier: "",
    };
}
exports.Header = {
    typeUrl: "/ibc.lightclients.solomachine.v3.Header",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.timestamp !== BigInt(0)) {
            writer.uint32(8).uint64(message.timestamp);
        }
        if (message.signature.length !== 0) {
            writer.uint32(18).bytes(message.signature);
        }
        if (message.newPublicKey !== undefined) {
            any_1.Any.encode(message.newPublicKey, writer.uint32(26).fork()).ldelim();
        }
        if (message.newDiversifier !== "") {
            writer.uint32(34).string(message.newDiversifier);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseHeader();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.timestamp = reader.uint64();
                    break;
                case 2:
                    message.signature = reader.bytes();
                    break;
                case 3:
                    message.newPublicKey = any_1.Any.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.newDiversifier = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseHeader();
        if ((0, helpers_1.isSet)(object.timestamp))
            obj.timestamp = BigInt(object.timestamp.toString());
        if ((0, helpers_1.isSet)(object.signature))
            obj.signature = (0, helpers_1.bytesFromBase64)(object.signature);
        if ((0, helpers_1.isSet)(object.newPublicKey))
            obj.newPublicKey = any_1.Any.fromJSON(object.newPublicKey);
        if ((0, helpers_1.isSet)(object.newDiversifier))
            obj.newDiversifier = String(object.newDiversifier);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.timestamp !== undefined && (obj.timestamp = (message.timestamp || BigInt(0)).toString());
        message.signature !== undefined &&
            (obj.signature = (0, helpers_1.base64FromBytes)(message.signature !== undefined ? message.signature : new Uint8Array()));
        message.newPublicKey !== undefined &&
            (obj.newPublicKey = message.newPublicKey ? any_1.Any.toJSON(message.newPublicKey) : undefined);
        message.newDiversifier !== undefined && (obj.newDiversifier = message.newDiversifier);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseHeader();
        if (object.timestamp !== undefined && object.timestamp !== null) {
            message.timestamp = BigInt(object.timestamp.toString());
        }
        message.signature = object.signature ?? new Uint8Array();
        if (object.newPublicKey !== undefined && object.newPublicKey !== null) {
            message.newPublicKey = any_1.Any.fromPartial(object.newPublicKey);
        }
        message.newDiversifier = object.newDiversifier ?? "";
        return message;
    },
};
function createBaseMisbehaviour() {
    return {
        sequence: BigInt(0),
        signatureOne: undefined,
        signatureTwo: undefined,
    };
}
exports.Misbehaviour = {
    typeUrl: "/ibc.lightclients.solomachine.v3.Misbehaviour",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.sequence !== BigInt(0)) {
            writer.uint32(8).uint64(message.sequence);
        }
        if (message.signatureOne !== undefined) {
            exports.SignatureAndData.encode(message.signatureOne, writer.uint32(18).fork()).ldelim();
        }
        if (message.signatureTwo !== undefined) {
            exports.SignatureAndData.encode(message.signatureTwo, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMisbehaviour();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.sequence = reader.uint64();
                    break;
                case 2:
                    message.signatureOne = exports.SignatureAndData.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.signatureTwo = exports.SignatureAndData.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMisbehaviour();
        if ((0, helpers_1.isSet)(object.sequence))
            obj.sequence = BigInt(object.sequence.toString());
        if ((0, helpers_1.isSet)(object.signatureOne))
            obj.signatureOne = exports.SignatureAndData.fromJSON(object.signatureOne);
        if ((0, helpers_1.isSet)(object.signatureTwo))
            obj.signatureTwo = exports.SignatureAndData.fromJSON(object.signatureTwo);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.sequence !== undefined && (obj.sequence = (message.sequence || BigInt(0)).toString());
        message.signatureOne !== undefined &&
            (obj.signatureOne = message.signatureOne ? exports.SignatureAndData.toJSON(message.signatureOne) : undefined);
        message.signatureTwo !== undefined &&
            (obj.signatureTwo = message.signatureTwo ? exports.SignatureAndData.toJSON(message.signatureTwo) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMisbehaviour();
        if (object.sequence !== undefined && object.sequence !== null) {
            message.sequence = BigInt(object.sequence.toString());
        }
        if (object.signatureOne !== undefined && object.signatureOne !== null) {
            message.signatureOne = exports.SignatureAndData.fromPartial(object.signatureOne);
        }
        if (object.signatureTwo !== undefined && object.signatureTwo !== null) {
            message.signatureTwo = exports.SignatureAndData.fromPartial(object.signatureTwo);
        }
        return message;
    },
};
function createBaseSignatureAndData() {
    return {
        signature: new Uint8Array(),
        path: new Uint8Array(),
        data: new Uint8Array(),
        timestamp: BigInt(0),
    };
}
exports.SignatureAndData = {
    typeUrl: "/ibc.lightclients.solomachine.v3.SignatureAndData",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.signature.length !== 0) {
            writer.uint32(10).bytes(message.signature);
        }
        if (message.path.length !== 0) {
            writer.uint32(18).bytes(message.path);
        }
        if (message.data.length !== 0) {
            writer.uint32(26).bytes(message.data);
        }
        if (message.timestamp !== BigInt(0)) {
            writer.uint32(32).uint64(message.timestamp);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSignatureAndData();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.signature = reader.bytes();
                    break;
                case 2:
                    message.path = reader.bytes();
                    break;
                case 3:
                    message.data = reader.bytes();
                    break;
                case 4:
                    message.timestamp = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseSignatureAndData();
        if ((0, helpers_1.isSet)(object.signature))
            obj.signature = (0, helpers_1.bytesFromBase64)(object.signature);
        if ((0, helpers_1.isSet)(object.path))
            obj.path = (0, helpers_1.bytesFromBase64)(object.path);
        if ((0, helpers_1.isSet)(object.data))
            obj.data = (0, helpers_1.bytesFromBase64)(object.data);
        if ((0, helpers_1.isSet)(object.timestamp))
            obj.timestamp = BigInt(object.timestamp.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.signature !== undefined &&
            (obj.signature = (0, helpers_1.base64FromBytes)(message.signature !== undefined ? message.signature : new Uint8Array()));
        message.path !== undefined &&
            (obj.path = (0, helpers_1.base64FromBytes)(message.path !== undefined ? message.path : new Uint8Array()));
        message.data !== undefined &&
            (obj.data = (0, helpers_1.base64FromBytes)(message.data !== undefined ? message.data : new Uint8Array()));
        message.timestamp !== undefined && (obj.timestamp = (message.timestamp || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseSignatureAndData();
        message.signature = object.signature ?? new Uint8Array();
        message.path = object.path ?? new Uint8Array();
        message.data = object.data ?? new Uint8Array();
        if (object.timestamp !== undefined && object.timestamp !== null) {
            message.timestamp = BigInt(object.timestamp.toString());
        }
        return message;
    },
};
function createBaseTimestampedSignatureData() {
    return {
        signatureData: new Uint8Array(),
        timestamp: BigInt(0),
    };
}
exports.TimestampedSignatureData = {
    typeUrl: "/ibc.lightclients.solomachine.v3.TimestampedSignatureData",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.signatureData.length !== 0) {
            writer.uint32(10).bytes(message.signatureData);
        }
        if (message.timestamp !== BigInt(0)) {
            writer.uint32(16).uint64(message.timestamp);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseTimestampedSignatureData();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.signatureData = reader.bytes();
                    break;
                case 2:
                    message.timestamp = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseTimestampedSignatureData();
        if ((0, helpers_1.isSet)(object.signatureData))
            obj.signatureData = (0, helpers_1.bytesFromBase64)(object.signatureData);
        if ((0, helpers_1.isSet)(object.timestamp))
            obj.timestamp = BigInt(object.timestamp.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.signatureData !== undefined &&
            (obj.signatureData = (0, helpers_1.base64FromBytes)(message.signatureData !== undefined ? message.signatureData : new Uint8Array()));
        message.timestamp !== undefined && (obj.timestamp = (message.timestamp || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseTimestampedSignatureData();
        message.signatureData = object.signatureData ?? new Uint8Array();
        if (object.timestamp !== undefined && object.timestamp !== null) {
            message.timestamp = BigInt(object.timestamp.toString());
        }
        return message;
    },
};
function createBaseSignBytes() {
    return {
        sequence: BigInt(0),
        timestamp: BigInt(0),
        diversifier: "",
        path: new Uint8Array(),
        data: new Uint8Array(),
    };
}
exports.SignBytes = {
    typeUrl: "/ibc.lightclients.solomachine.v3.SignBytes",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.sequence !== BigInt(0)) {
            writer.uint32(8).uint64(message.sequence);
        }
        if (message.timestamp !== BigInt(0)) {
            writer.uint32(16).uint64(message.timestamp);
        }
        if (message.diversifier !== "") {
            writer.uint32(26).string(message.diversifier);
        }
        if (message.path.length !== 0) {
            writer.uint32(34).bytes(message.path);
        }
        if (message.data.length !== 0) {
            writer.uint32(42).bytes(message.data);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSignBytes();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.sequence = reader.uint64();
                    break;
                case 2:
                    message.timestamp = reader.uint64();
                    break;
                case 3:
                    message.diversifier = reader.string();
                    break;
                case 4:
                    message.path = reader.bytes();
                    break;
                case 5:
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
        const obj = createBaseSignBytes();
        if ((0, helpers_1.isSet)(object.sequence))
            obj.sequence = BigInt(object.sequence.toString());
        if ((0, helpers_1.isSet)(object.timestamp))
            obj.timestamp = BigInt(object.timestamp.toString());
        if ((0, helpers_1.isSet)(object.diversifier))
            obj.diversifier = String(object.diversifier);
        if ((0, helpers_1.isSet)(object.path))
            obj.path = (0, helpers_1.bytesFromBase64)(object.path);
        if ((0, helpers_1.isSet)(object.data))
            obj.data = (0, helpers_1.bytesFromBase64)(object.data);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.sequence !== undefined && (obj.sequence = (message.sequence || BigInt(0)).toString());
        message.timestamp !== undefined && (obj.timestamp = (message.timestamp || BigInt(0)).toString());
        message.diversifier !== undefined && (obj.diversifier = message.diversifier);
        message.path !== undefined &&
            (obj.path = (0, helpers_1.base64FromBytes)(message.path !== undefined ? message.path : new Uint8Array()));
        message.data !== undefined &&
            (obj.data = (0, helpers_1.base64FromBytes)(message.data !== undefined ? message.data : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseSignBytes();
        if (object.sequence !== undefined && object.sequence !== null) {
            message.sequence = BigInt(object.sequence.toString());
        }
        if (object.timestamp !== undefined && object.timestamp !== null) {
            message.timestamp = BigInt(object.timestamp.toString());
        }
        message.diversifier = object.diversifier ?? "";
        message.path = object.path ?? new Uint8Array();
        message.data = object.data ?? new Uint8Array();
        return message;
    },
};
function createBaseHeaderData() {
    return {
        newPubKey: undefined,
        newDiversifier: "",
    };
}
exports.HeaderData = {
    typeUrl: "/ibc.lightclients.solomachine.v3.HeaderData",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.newPubKey !== undefined) {
            any_1.Any.encode(message.newPubKey, writer.uint32(10).fork()).ldelim();
        }
        if (message.newDiversifier !== "") {
            writer.uint32(18).string(message.newDiversifier);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseHeaderData();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.newPubKey = any_1.Any.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.newDiversifier = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseHeaderData();
        if ((0, helpers_1.isSet)(object.newPubKey))
            obj.newPubKey = any_1.Any.fromJSON(object.newPubKey);
        if ((0, helpers_1.isSet)(object.newDiversifier))
            obj.newDiversifier = String(object.newDiversifier);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.newPubKey !== undefined &&
            (obj.newPubKey = message.newPubKey ? any_1.Any.toJSON(message.newPubKey) : undefined);
        message.newDiversifier !== undefined && (obj.newDiversifier = message.newDiversifier);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseHeaderData();
        if (object.newPubKey !== undefined && object.newPubKey !== null) {
            message.newPubKey = any_1.Any.fromPartial(object.newPubKey);
        }
        message.newDiversifier = object.newDiversifier ?? "";
        return message;
    },
};
//# sourceMappingURL=solomachine.js.map