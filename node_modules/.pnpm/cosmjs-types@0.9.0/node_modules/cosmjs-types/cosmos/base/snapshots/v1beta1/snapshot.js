"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SnapshotSchema = exports.SnapshotKVItem = exports.SnapshotExtensionPayload = exports.SnapshotExtensionMeta = exports.SnapshotIAVLItem = exports.SnapshotStoreItem = exports.SnapshotItem = exports.Metadata = exports.Snapshot = exports.protobufPackage = void 0;
/* eslint-disable */
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "cosmos.base.snapshots.v1beta1";
function createBaseSnapshot() {
    return {
        height: BigInt(0),
        format: 0,
        chunks: 0,
        hash: new Uint8Array(),
        metadata: exports.Metadata.fromPartial({}),
    };
}
exports.Snapshot = {
    typeUrl: "/cosmos.base.snapshots.v1beta1.Snapshot",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.height !== BigInt(0)) {
            writer.uint32(8).uint64(message.height);
        }
        if (message.format !== 0) {
            writer.uint32(16).uint32(message.format);
        }
        if (message.chunks !== 0) {
            writer.uint32(24).uint32(message.chunks);
        }
        if (message.hash.length !== 0) {
            writer.uint32(34).bytes(message.hash);
        }
        if (message.metadata !== undefined) {
            exports.Metadata.encode(message.metadata, writer.uint32(42).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSnapshot();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.height = reader.uint64();
                    break;
                case 2:
                    message.format = reader.uint32();
                    break;
                case 3:
                    message.chunks = reader.uint32();
                    break;
                case 4:
                    message.hash = reader.bytes();
                    break;
                case 5:
                    message.metadata = exports.Metadata.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseSnapshot();
        if ((0, helpers_1.isSet)(object.height))
            obj.height = BigInt(object.height.toString());
        if ((0, helpers_1.isSet)(object.format))
            obj.format = Number(object.format);
        if ((0, helpers_1.isSet)(object.chunks))
            obj.chunks = Number(object.chunks);
        if ((0, helpers_1.isSet)(object.hash))
            obj.hash = (0, helpers_1.bytesFromBase64)(object.hash);
        if ((0, helpers_1.isSet)(object.metadata))
            obj.metadata = exports.Metadata.fromJSON(object.metadata);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.height !== undefined && (obj.height = (message.height || BigInt(0)).toString());
        message.format !== undefined && (obj.format = Math.round(message.format));
        message.chunks !== undefined && (obj.chunks = Math.round(message.chunks));
        message.hash !== undefined &&
            (obj.hash = (0, helpers_1.base64FromBytes)(message.hash !== undefined ? message.hash : new Uint8Array()));
        message.metadata !== undefined &&
            (obj.metadata = message.metadata ? exports.Metadata.toJSON(message.metadata) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseSnapshot();
        if (object.height !== undefined && object.height !== null) {
            message.height = BigInt(object.height.toString());
        }
        message.format = object.format ?? 0;
        message.chunks = object.chunks ?? 0;
        message.hash = object.hash ?? new Uint8Array();
        if (object.metadata !== undefined && object.metadata !== null) {
            message.metadata = exports.Metadata.fromPartial(object.metadata);
        }
        return message;
    },
};
function createBaseMetadata() {
    return {
        chunkHashes: [],
    };
}
exports.Metadata = {
    typeUrl: "/cosmos.base.snapshots.v1beta1.Metadata",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.chunkHashes) {
            writer.uint32(10).bytes(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMetadata();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.chunkHashes.push(reader.bytes());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMetadata();
        if (Array.isArray(object?.chunkHashes))
            obj.chunkHashes = object.chunkHashes.map((e) => (0, helpers_1.bytesFromBase64)(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.chunkHashes) {
            obj.chunkHashes = message.chunkHashes.map((e) => (0, helpers_1.base64FromBytes)(e !== undefined ? e : new Uint8Array()));
        }
        else {
            obj.chunkHashes = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMetadata();
        message.chunkHashes = object.chunkHashes?.map((e) => e) || [];
        return message;
    },
};
function createBaseSnapshotItem() {
    return {
        store: undefined,
        iavl: undefined,
        extension: undefined,
        extensionPayload: undefined,
        kv: undefined,
        schema: undefined,
    };
}
exports.SnapshotItem = {
    typeUrl: "/cosmos.base.snapshots.v1beta1.SnapshotItem",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.store !== undefined) {
            exports.SnapshotStoreItem.encode(message.store, writer.uint32(10).fork()).ldelim();
        }
        if (message.iavl !== undefined) {
            exports.SnapshotIAVLItem.encode(message.iavl, writer.uint32(18).fork()).ldelim();
        }
        if (message.extension !== undefined) {
            exports.SnapshotExtensionMeta.encode(message.extension, writer.uint32(26).fork()).ldelim();
        }
        if (message.extensionPayload !== undefined) {
            exports.SnapshotExtensionPayload.encode(message.extensionPayload, writer.uint32(34).fork()).ldelim();
        }
        if (message.kv !== undefined) {
            exports.SnapshotKVItem.encode(message.kv, writer.uint32(42).fork()).ldelim();
        }
        if (message.schema !== undefined) {
            exports.SnapshotSchema.encode(message.schema, writer.uint32(50).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSnapshotItem();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.store = exports.SnapshotStoreItem.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.iavl = exports.SnapshotIAVLItem.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.extension = exports.SnapshotExtensionMeta.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.extensionPayload = exports.SnapshotExtensionPayload.decode(reader, reader.uint32());
                    break;
                case 5:
                    message.kv = exports.SnapshotKVItem.decode(reader, reader.uint32());
                    break;
                case 6:
                    message.schema = exports.SnapshotSchema.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseSnapshotItem();
        if ((0, helpers_1.isSet)(object.store))
            obj.store = exports.SnapshotStoreItem.fromJSON(object.store);
        if ((0, helpers_1.isSet)(object.iavl))
            obj.iavl = exports.SnapshotIAVLItem.fromJSON(object.iavl);
        if ((0, helpers_1.isSet)(object.extension))
            obj.extension = exports.SnapshotExtensionMeta.fromJSON(object.extension);
        if ((0, helpers_1.isSet)(object.extensionPayload))
            obj.extensionPayload = exports.SnapshotExtensionPayload.fromJSON(object.extensionPayload);
        if ((0, helpers_1.isSet)(object.kv))
            obj.kv = exports.SnapshotKVItem.fromJSON(object.kv);
        if ((0, helpers_1.isSet)(object.schema))
            obj.schema = exports.SnapshotSchema.fromJSON(object.schema);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.store !== undefined &&
            (obj.store = message.store ? exports.SnapshotStoreItem.toJSON(message.store) : undefined);
        message.iavl !== undefined &&
            (obj.iavl = message.iavl ? exports.SnapshotIAVLItem.toJSON(message.iavl) : undefined);
        message.extension !== undefined &&
            (obj.extension = message.extension ? exports.SnapshotExtensionMeta.toJSON(message.extension) : undefined);
        message.extensionPayload !== undefined &&
            (obj.extensionPayload = message.extensionPayload
                ? exports.SnapshotExtensionPayload.toJSON(message.extensionPayload)
                : undefined);
        message.kv !== undefined && (obj.kv = message.kv ? exports.SnapshotKVItem.toJSON(message.kv) : undefined);
        message.schema !== undefined &&
            (obj.schema = message.schema ? exports.SnapshotSchema.toJSON(message.schema) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseSnapshotItem();
        if (object.store !== undefined && object.store !== null) {
            message.store = exports.SnapshotStoreItem.fromPartial(object.store);
        }
        if (object.iavl !== undefined && object.iavl !== null) {
            message.iavl = exports.SnapshotIAVLItem.fromPartial(object.iavl);
        }
        if (object.extension !== undefined && object.extension !== null) {
            message.extension = exports.SnapshotExtensionMeta.fromPartial(object.extension);
        }
        if (object.extensionPayload !== undefined && object.extensionPayload !== null) {
            message.extensionPayload = exports.SnapshotExtensionPayload.fromPartial(object.extensionPayload);
        }
        if (object.kv !== undefined && object.kv !== null) {
            message.kv = exports.SnapshotKVItem.fromPartial(object.kv);
        }
        if (object.schema !== undefined && object.schema !== null) {
            message.schema = exports.SnapshotSchema.fromPartial(object.schema);
        }
        return message;
    },
};
function createBaseSnapshotStoreItem() {
    return {
        name: "",
    };
}
exports.SnapshotStoreItem = {
    typeUrl: "/cosmos.base.snapshots.v1beta1.SnapshotStoreItem",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.name !== "") {
            writer.uint32(10).string(message.name);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSnapshotStoreItem();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.name = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseSnapshotStoreItem();
        if ((0, helpers_1.isSet)(object.name))
            obj.name = String(object.name);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.name !== undefined && (obj.name = message.name);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseSnapshotStoreItem();
        message.name = object.name ?? "";
        return message;
    },
};
function createBaseSnapshotIAVLItem() {
    return {
        key: new Uint8Array(),
        value: new Uint8Array(),
        version: BigInt(0),
        height: 0,
    };
}
exports.SnapshotIAVLItem = {
    typeUrl: "/cosmos.base.snapshots.v1beta1.SnapshotIAVLItem",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.key.length !== 0) {
            writer.uint32(10).bytes(message.key);
        }
        if (message.value.length !== 0) {
            writer.uint32(18).bytes(message.value);
        }
        if (message.version !== BigInt(0)) {
            writer.uint32(24).int64(message.version);
        }
        if (message.height !== 0) {
            writer.uint32(32).int32(message.height);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSnapshotIAVLItem();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.key = reader.bytes();
                    break;
                case 2:
                    message.value = reader.bytes();
                    break;
                case 3:
                    message.version = reader.int64();
                    break;
                case 4:
                    message.height = reader.int32();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseSnapshotIAVLItem();
        if ((0, helpers_1.isSet)(object.key))
            obj.key = (0, helpers_1.bytesFromBase64)(object.key);
        if ((0, helpers_1.isSet)(object.value))
            obj.value = (0, helpers_1.bytesFromBase64)(object.value);
        if ((0, helpers_1.isSet)(object.version))
            obj.version = BigInt(object.version.toString());
        if ((0, helpers_1.isSet)(object.height))
            obj.height = Number(object.height);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.key !== undefined &&
            (obj.key = (0, helpers_1.base64FromBytes)(message.key !== undefined ? message.key : new Uint8Array()));
        message.value !== undefined &&
            (obj.value = (0, helpers_1.base64FromBytes)(message.value !== undefined ? message.value : new Uint8Array()));
        message.version !== undefined && (obj.version = (message.version || BigInt(0)).toString());
        message.height !== undefined && (obj.height = Math.round(message.height));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseSnapshotIAVLItem();
        message.key = object.key ?? new Uint8Array();
        message.value = object.value ?? new Uint8Array();
        if (object.version !== undefined && object.version !== null) {
            message.version = BigInt(object.version.toString());
        }
        message.height = object.height ?? 0;
        return message;
    },
};
function createBaseSnapshotExtensionMeta() {
    return {
        name: "",
        format: 0,
    };
}
exports.SnapshotExtensionMeta = {
    typeUrl: "/cosmos.base.snapshots.v1beta1.SnapshotExtensionMeta",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.name !== "") {
            writer.uint32(10).string(message.name);
        }
        if (message.format !== 0) {
            writer.uint32(16).uint32(message.format);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSnapshotExtensionMeta();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.name = reader.string();
                    break;
                case 2:
                    message.format = reader.uint32();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseSnapshotExtensionMeta();
        if ((0, helpers_1.isSet)(object.name))
            obj.name = String(object.name);
        if ((0, helpers_1.isSet)(object.format))
            obj.format = Number(object.format);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.name !== undefined && (obj.name = message.name);
        message.format !== undefined && (obj.format = Math.round(message.format));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseSnapshotExtensionMeta();
        message.name = object.name ?? "";
        message.format = object.format ?? 0;
        return message;
    },
};
function createBaseSnapshotExtensionPayload() {
    return {
        payload: new Uint8Array(),
    };
}
exports.SnapshotExtensionPayload = {
    typeUrl: "/cosmos.base.snapshots.v1beta1.SnapshotExtensionPayload",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.payload.length !== 0) {
            writer.uint32(10).bytes(message.payload);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSnapshotExtensionPayload();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.payload = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseSnapshotExtensionPayload();
        if ((0, helpers_1.isSet)(object.payload))
            obj.payload = (0, helpers_1.bytesFromBase64)(object.payload);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.payload !== undefined &&
            (obj.payload = (0, helpers_1.base64FromBytes)(message.payload !== undefined ? message.payload : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseSnapshotExtensionPayload();
        message.payload = object.payload ?? new Uint8Array();
        return message;
    },
};
function createBaseSnapshotKVItem() {
    return {
        key: new Uint8Array(),
        value: new Uint8Array(),
    };
}
exports.SnapshotKVItem = {
    typeUrl: "/cosmos.base.snapshots.v1beta1.SnapshotKVItem",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.key.length !== 0) {
            writer.uint32(10).bytes(message.key);
        }
        if (message.value.length !== 0) {
            writer.uint32(18).bytes(message.value);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSnapshotKVItem();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.key = reader.bytes();
                    break;
                case 2:
                    message.value = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseSnapshotKVItem();
        if ((0, helpers_1.isSet)(object.key))
            obj.key = (0, helpers_1.bytesFromBase64)(object.key);
        if ((0, helpers_1.isSet)(object.value))
            obj.value = (0, helpers_1.bytesFromBase64)(object.value);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.key !== undefined &&
            (obj.key = (0, helpers_1.base64FromBytes)(message.key !== undefined ? message.key : new Uint8Array()));
        message.value !== undefined &&
            (obj.value = (0, helpers_1.base64FromBytes)(message.value !== undefined ? message.value : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseSnapshotKVItem();
        message.key = object.key ?? new Uint8Array();
        message.value = object.value ?? new Uint8Array();
        return message;
    },
};
function createBaseSnapshotSchema() {
    return {
        keys: [],
    };
}
exports.SnapshotSchema = {
    typeUrl: "/cosmos.base.snapshots.v1beta1.SnapshotSchema",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.keys) {
            writer.uint32(10).bytes(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSnapshotSchema();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.keys.push(reader.bytes());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseSnapshotSchema();
        if (Array.isArray(object?.keys))
            obj.keys = object.keys.map((e) => (0, helpers_1.bytesFromBase64)(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.keys) {
            obj.keys = message.keys.map((e) => (0, helpers_1.base64FromBytes)(e !== undefined ? e : new Uint8Array()));
        }
        else {
            obj.keys = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseSnapshotSchema();
        message.keys = object.keys?.map((e) => e) || [];
        return message;
    },
};
//# sourceMappingURL=snapshot.js.map