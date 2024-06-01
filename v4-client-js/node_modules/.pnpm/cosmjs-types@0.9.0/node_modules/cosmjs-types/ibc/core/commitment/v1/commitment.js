"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.MerkleProof = exports.MerklePath = exports.MerklePrefix = exports.MerkleRoot = exports.protobufPackage = void 0;
/* eslint-disable */
const proofs_1 = require("../../../../cosmos/ics23/v1/proofs");
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "ibc.core.commitment.v1";
function createBaseMerkleRoot() {
    return {
        hash: new Uint8Array(),
    };
}
exports.MerkleRoot = {
    typeUrl: "/ibc.core.commitment.v1.MerkleRoot",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.hash.length !== 0) {
            writer.uint32(10).bytes(message.hash);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMerkleRoot();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
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
        const obj = createBaseMerkleRoot();
        if ((0, helpers_1.isSet)(object.hash))
            obj.hash = (0, helpers_1.bytesFromBase64)(object.hash);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.hash !== undefined &&
            (obj.hash = (0, helpers_1.base64FromBytes)(message.hash !== undefined ? message.hash : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMerkleRoot();
        message.hash = object.hash ?? new Uint8Array();
        return message;
    },
};
function createBaseMerklePrefix() {
    return {
        keyPrefix: new Uint8Array(),
    };
}
exports.MerklePrefix = {
    typeUrl: "/ibc.core.commitment.v1.MerklePrefix",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.keyPrefix.length !== 0) {
            writer.uint32(10).bytes(message.keyPrefix);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMerklePrefix();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.keyPrefix = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMerklePrefix();
        if ((0, helpers_1.isSet)(object.keyPrefix))
            obj.keyPrefix = (0, helpers_1.bytesFromBase64)(object.keyPrefix);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.keyPrefix !== undefined &&
            (obj.keyPrefix = (0, helpers_1.base64FromBytes)(message.keyPrefix !== undefined ? message.keyPrefix : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMerklePrefix();
        message.keyPrefix = object.keyPrefix ?? new Uint8Array();
        return message;
    },
};
function createBaseMerklePath() {
    return {
        keyPath: [],
    };
}
exports.MerklePath = {
    typeUrl: "/ibc.core.commitment.v1.MerklePath",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.keyPath) {
            writer.uint32(10).string(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMerklePath();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.keyPath.push(reader.string());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMerklePath();
        if (Array.isArray(object?.keyPath))
            obj.keyPath = object.keyPath.map((e) => String(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.keyPath) {
            obj.keyPath = message.keyPath.map((e) => e);
        }
        else {
            obj.keyPath = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMerklePath();
        message.keyPath = object.keyPath?.map((e) => e) || [];
        return message;
    },
};
function createBaseMerkleProof() {
    return {
        proofs: [],
    };
}
exports.MerkleProof = {
    typeUrl: "/ibc.core.commitment.v1.MerkleProof",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.proofs) {
            proofs_1.CommitmentProof.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMerkleProof();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.proofs.push(proofs_1.CommitmentProof.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMerkleProof();
        if (Array.isArray(object?.proofs))
            obj.proofs = object.proofs.map((e) => proofs_1.CommitmentProof.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.proofs) {
            obj.proofs = message.proofs.map((e) => (e ? proofs_1.CommitmentProof.toJSON(e) : undefined));
        }
        else {
            obj.proofs = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMerkleProof();
        message.proofs = object.proofs?.map((e) => proofs_1.CommitmentProof.fromPartial(e)) || [];
        return message;
    },
};
//# sourceMappingURL=commitment.js.map