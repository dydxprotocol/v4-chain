"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CommitID = exports.StoreInfo = exports.CommitInfo = exports.protobufPackage = void 0;
/* eslint-disable */
const timestamp_1 = require("../../../../google/protobuf/timestamp");
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "cosmos.base.store.v1beta1";
function createBaseCommitInfo() {
    return {
        version: BigInt(0),
        storeInfos: [],
        timestamp: timestamp_1.Timestamp.fromPartial({}),
    };
}
exports.CommitInfo = {
    typeUrl: "/cosmos.base.store.v1beta1.CommitInfo",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.version !== BigInt(0)) {
            writer.uint32(8).int64(message.version);
        }
        for (const v of message.storeInfos) {
            exports.StoreInfo.encode(v, writer.uint32(18).fork()).ldelim();
        }
        if (message.timestamp !== undefined) {
            timestamp_1.Timestamp.encode(message.timestamp, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseCommitInfo();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.version = reader.int64();
                    break;
                case 2:
                    message.storeInfos.push(exports.StoreInfo.decode(reader, reader.uint32()));
                    break;
                case 3:
                    message.timestamp = timestamp_1.Timestamp.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseCommitInfo();
        if ((0, helpers_1.isSet)(object.version))
            obj.version = BigInt(object.version.toString());
        if (Array.isArray(object?.storeInfos))
            obj.storeInfos = object.storeInfos.map((e) => exports.StoreInfo.fromJSON(e));
        if ((0, helpers_1.isSet)(object.timestamp))
            obj.timestamp = (0, helpers_1.fromJsonTimestamp)(object.timestamp);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.version !== undefined && (obj.version = (message.version || BigInt(0)).toString());
        if (message.storeInfos) {
            obj.storeInfos = message.storeInfos.map((e) => (e ? exports.StoreInfo.toJSON(e) : undefined));
        }
        else {
            obj.storeInfos = [];
        }
        message.timestamp !== undefined && (obj.timestamp = (0, helpers_1.fromTimestamp)(message.timestamp).toISOString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseCommitInfo();
        if (object.version !== undefined && object.version !== null) {
            message.version = BigInt(object.version.toString());
        }
        message.storeInfos = object.storeInfos?.map((e) => exports.StoreInfo.fromPartial(e)) || [];
        if (object.timestamp !== undefined && object.timestamp !== null) {
            message.timestamp = timestamp_1.Timestamp.fromPartial(object.timestamp);
        }
        return message;
    },
};
function createBaseStoreInfo() {
    return {
        name: "",
        commitId: exports.CommitID.fromPartial({}),
    };
}
exports.StoreInfo = {
    typeUrl: "/cosmos.base.store.v1beta1.StoreInfo",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.name !== "") {
            writer.uint32(10).string(message.name);
        }
        if (message.commitId !== undefined) {
            exports.CommitID.encode(message.commitId, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseStoreInfo();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.name = reader.string();
                    break;
                case 2:
                    message.commitId = exports.CommitID.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseStoreInfo();
        if ((0, helpers_1.isSet)(object.name))
            obj.name = String(object.name);
        if ((0, helpers_1.isSet)(object.commitId))
            obj.commitId = exports.CommitID.fromJSON(object.commitId);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.name !== undefined && (obj.name = message.name);
        message.commitId !== undefined &&
            (obj.commitId = message.commitId ? exports.CommitID.toJSON(message.commitId) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseStoreInfo();
        message.name = object.name ?? "";
        if (object.commitId !== undefined && object.commitId !== null) {
            message.commitId = exports.CommitID.fromPartial(object.commitId);
        }
        return message;
    },
};
function createBaseCommitID() {
    return {
        version: BigInt(0),
        hash: new Uint8Array(),
    };
}
exports.CommitID = {
    typeUrl: "/cosmos.base.store.v1beta1.CommitID",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.version !== BigInt(0)) {
            writer.uint32(8).int64(message.version);
        }
        if (message.hash.length !== 0) {
            writer.uint32(18).bytes(message.hash);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseCommitID();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.version = reader.int64();
                    break;
                case 2:
                    message.hash = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseCommitID();
        if ((0, helpers_1.isSet)(object.version))
            obj.version = BigInt(object.version.toString());
        if ((0, helpers_1.isSet)(object.hash))
            obj.hash = (0, helpers_1.bytesFromBase64)(object.hash);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.version !== undefined && (obj.version = (message.version || BigInt(0)).toString());
        message.hash !== undefined &&
            (obj.hash = (0, helpers_1.base64FromBytes)(message.hash !== undefined ? message.hash : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseCommitID();
        if (object.version !== undefined && object.version !== null) {
            message.version = BigInt(object.version.toString());
        }
        message.hash = object.hash ?? new Uint8Array();
        return message;
    },
};
//# sourceMappingURL=commit_info.js.map