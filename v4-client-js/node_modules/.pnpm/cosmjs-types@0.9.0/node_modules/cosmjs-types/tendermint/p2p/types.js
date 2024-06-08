"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DefaultNodeInfoOther = exports.DefaultNodeInfo = exports.ProtocolVersion = exports.NetAddress = exports.protobufPackage = void 0;
/* eslint-disable */
const binary_1 = require("../../binary");
const helpers_1 = require("../../helpers");
exports.protobufPackage = "tendermint.p2p";
function createBaseNetAddress() {
    return {
        id: "",
        ip: "",
        port: 0,
    };
}
exports.NetAddress = {
    typeUrl: "/tendermint.p2p.NetAddress",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.id !== "") {
            writer.uint32(10).string(message.id);
        }
        if (message.ip !== "") {
            writer.uint32(18).string(message.ip);
        }
        if (message.port !== 0) {
            writer.uint32(24).uint32(message.port);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseNetAddress();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.id = reader.string();
                    break;
                case 2:
                    message.ip = reader.string();
                    break;
                case 3:
                    message.port = reader.uint32();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseNetAddress();
        if ((0, helpers_1.isSet)(object.id))
            obj.id = String(object.id);
        if ((0, helpers_1.isSet)(object.ip))
            obj.ip = String(object.ip);
        if ((0, helpers_1.isSet)(object.port))
            obj.port = Number(object.port);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.id !== undefined && (obj.id = message.id);
        message.ip !== undefined && (obj.ip = message.ip);
        message.port !== undefined && (obj.port = Math.round(message.port));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseNetAddress();
        message.id = object.id ?? "";
        message.ip = object.ip ?? "";
        message.port = object.port ?? 0;
        return message;
    },
};
function createBaseProtocolVersion() {
    return {
        p2p: BigInt(0),
        block: BigInt(0),
        app: BigInt(0),
    };
}
exports.ProtocolVersion = {
    typeUrl: "/tendermint.p2p.ProtocolVersion",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.p2p !== BigInt(0)) {
            writer.uint32(8).uint64(message.p2p);
        }
        if (message.block !== BigInt(0)) {
            writer.uint32(16).uint64(message.block);
        }
        if (message.app !== BigInt(0)) {
            writer.uint32(24).uint64(message.app);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseProtocolVersion();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.p2p = reader.uint64();
                    break;
                case 2:
                    message.block = reader.uint64();
                    break;
                case 3:
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
        const obj = createBaseProtocolVersion();
        if ((0, helpers_1.isSet)(object.p2p))
            obj.p2p = BigInt(object.p2p.toString());
        if ((0, helpers_1.isSet)(object.block))
            obj.block = BigInt(object.block.toString());
        if ((0, helpers_1.isSet)(object.app))
            obj.app = BigInt(object.app.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.p2p !== undefined && (obj.p2p = (message.p2p || BigInt(0)).toString());
        message.block !== undefined && (obj.block = (message.block || BigInt(0)).toString());
        message.app !== undefined && (obj.app = (message.app || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseProtocolVersion();
        if (object.p2p !== undefined && object.p2p !== null) {
            message.p2p = BigInt(object.p2p.toString());
        }
        if (object.block !== undefined && object.block !== null) {
            message.block = BigInt(object.block.toString());
        }
        if (object.app !== undefined && object.app !== null) {
            message.app = BigInt(object.app.toString());
        }
        return message;
    },
};
function createBaseDefaultNodeInfo() {
    return {
        protocolVersion: exports.ProtocolVersion.fromPartial({}),
        defaultNodeId: "",
        listenAddr: "",
        network: "",
        version: "",
        channels: new Uint8Array(),
        moniker: "",
        other: exports.DefaultNodeInfoOther.fromPartial({}),
    };
}
exports.DefaultNodeInfo = {
    typeUrl: "/tendermint.p2p.DefaultNodeInfo",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.protocolVersion !== undefined) {
            exports.ProtocolVersion.encode(message.protocolVersion, writer.uint32(10).fork()).ldelim();
        }
        if (message.defaultNodeId !== "") {
            writer.uint32(18).string(message.defaultNodeId);
        }
        if (message.listenAddr !== "") {
            writer.uint32(26).string(message.listenAddr);
        }
        if (message.network !== "") {
            writer.uint32(34).string(message.network);
        }
        if (message.version !== "") {
            writer.uint32(42).string(message.version);
        }
        if (message.channels.length !== 0) {
            writer.uint32(50).bytes(message.channels);
        }
        if (message.moniker !== "") {
            writer.uint32(58).string(message.moniker);
        }
        if (message.other !== undefined) {
            exports.DefaultNodeInfoOther.encode(message.other, writer.uint32(66).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseDefaultNodeInfo();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.protocolVersion = exports.ProtocolVersion.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.defaultNodeId = reader.string();
                    break;
                case 3:
                    message.listenAddr = reader.string();
                    break;
                case 4:
                    message.network = reader.string();
                    break;
                case 5:
                    message.version = reader.string();
                    break;
                case 6:
                    message.channels = reader.bytes();
                    break;
                case 7:
                    message.moniker = reader.string();
                    break;
                case 8:
                    message.other = exports.DefaultNodeInfoOther.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseDefaultNodeInfo();
        if ((0, helpers_1.isSet)(object.protocolVersion))
            obj.protocolVersion = exports.ProtocolVersion.fromJSON(object.protocolVersion);
        if ((0, helpers_1.isSet)(object.defaultNodeId))
            obj.defaultNodeId = String(object.defaultNodeId);
        if ((0, helpers_1.isSet)(object.listenAddr))
            obj.listenAddr = String(object.listenAddr);
        if ((0, helpers_1.isSet)(object.network))
            obj.network = String(object.network);
        if ((0, helpers_1.isSet)(object.version))
            obj.version = String(object.version);
        if ((0, helpers_1.isSet)(object.channels))
            obj.channels = (0, helpers_1.bytesFromBase64)(object.channels);
        if ((0, helpers_1.isSet)(object.moniker))
            obj.moniker = String(object.moniker);
        if ((0, helpers_1.isSet)(object.other))
            obj.other = exports.DefaultNodeInfoOther.fromJSON(object.other);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.protocolVersion !== undefined &&
            (obj.protocolVersion = message.protocolVersion
                ? exports.ProtocolVersion.toJSON(message.protocolVersion)
                : undefined);
        message.defaultNodeId !== undefined && (obj.defaultNodeId = message.defaultNodeId);
        message.listenAddr !== undefined && (obj.listenAddr = message.listenAddr);
        message.network !== undefined && (obj.network = message.network);
        message.version !== undefined && (obj.version = message.version);
        message.channels !== undefined &&
            (obj.channels = (0, helpers_1.base64FromBytes)(message.channels !== undefined ? message.channels : new Uint8Array()));
        message.moniker !== undefined && (obj.moniker = message.moniker);
        message.other !== undefined &&
            (obj.other = message.other ? exports.DefaultNodeInfoOther.toJSON(message.other) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseDefaultNodeInfo();
        if (object.protocolVersion !== undefined && object.protocolVersion !== null) {
            message.protocolVersion = exports.ProtocolVersion.fromPartial(object.protocolVersion);
        }
        message.defaultNodeId = object.defaultNodeId ?? "";
        message.listenAddr = object.listenAddr ?? "";
        message.network = object.network ?? "";
        message.version = object.version ?? "";
        message.channels = object.channels ?? new Uint8Array();
        message.moniker = object.moniker ?? "";
        if (object.other !== undefined && object.other !== null) {
            message.other = exports.DefaultNodeInfoOther.fromPartial(object.other);
        }
        return message;
    },
};
function createBaseDefaultNodeInfoOther() {
    return {
        txIndex: "",
        rpcAddress: "",
    };
}
exports.DefaultNodeInfoOther = {
    typeUrl: "/tendermint.p2p.DefaultNodeInfoOther",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.txIndex !== "") {
            writer.uint32(10).string(message.txIndex);
        }
        if (message.rpcAddress !== "") {
            writer.uint32(18).string(message.rpcAddress);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseDefaultNodeInfoOther();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.txIndex = reader.string();
                    break;
                case 2:
                    message.rpcAddress = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseDefaultNodeInfoOther();
        if ((0, helpers_1.isSet)(object.txIndex))
            obj.txIndex = String(object.txIndex);
        if ((0, helpers_1.isSet)(object.rpcAddress))
            obj.rpcAddress = String(object.rpcAddress);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.txIndex !== undefined && (obj.txIndex = message.txIndex);
        message.rpcAddress !== undefined && (obj.rpcAddress = message.rpcAddress);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseDefaultNodeInfoOther();
        message.txIndex = object.txIndex ?? "";
        message.rpcAddress = object.rpcAddress ?? "";
        return message;
    },
};
//# sourceMappingURL=types.js.map