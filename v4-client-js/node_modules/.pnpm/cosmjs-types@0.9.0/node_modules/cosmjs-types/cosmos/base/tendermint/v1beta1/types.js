"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Header = exports.Block = exports.protobufPackage = void 0;
/* eslint-disable */
const types_1 = require("../../../../tendermint/types/types");
const evidence_1 = require("../../../../tendermint/types/evidence");
const types_2 = require("../../../../tendermint/version/types");
const timestamp_1 = require("../../../../google/protobuf/timestamp");
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "cosmos.base.tendermint.v1beta1";
function createBaseBlock() {
    return {
        header: exports.Header.fromPartial({}),
        data: types_1.Data.fromPartial({}),
        evidence: evidence_1.EvidenceList.fromPartial({}),
        lastCommit: undefined,
    };
}
exports.Block = {
    typeUrl: "/cosmos.base.tendermint.v1beta1.Block",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.header !== undefined) {
            exports.Header.encode(message.header, writer.uint32(10).fork()).ldelim();
        }
        if (message.data !== undefined) {
            types_1.Data.encode(message.data, writer.uint32(18).fork()).ldelim();
        }
        if (message.evidence !== undefined) {
            evidence_1.EvidenceList.encode(message.evidence, writer.uint32(26).fork()).ldelim();
        }
        if (message.lastCommit !== undefined) {
            types_1.Commit.encode(message.lastCommit, writer.uint32(34).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseBlock();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.header = exports.Header.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.data = types_1.Data.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.evidence = evidence_1.EvidenceList.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.lastCommit = types_1.Commit.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseBlock();
        if ((0, helpers_1.isSet)(object.header))
            obj.header = exports.Header.fromJSON(object.header);
        if ((0, helpers_1.isSet)(object.data))
            obj.data = types_1.Data.fromJSON(object.data);
        if ((0, helpers_1.isSet)(object.evidence))
            obj.evidence = evidence_1.EvidenceList.fromJSON(object.evidence);
        if ((0, helpers_1.isSet)(object.lastCommit))
            obj.lastCommit = types_1.Commit.fromJSON(object.lastCommit);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.header !== undefined && (obj.header = message.header ? exports.Header.toJSON(message.header) : undefined);
        message.data !== undefined && (obj.data = message.data ? types_1.Data.toJSON(message.data) : undefined);
        message.evidence !== undefined &&
            (obj.evidence = message.evidence ? evidence_1.EvidenceList.toJSON(message.evidence) : undefined);
        message.lastCommit !== undefined &&
            (obj.lastCommit = message.lastCommit ? types_1.Commit.toJSON(message.lastCommit) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseBlock();
        if (object.header !== undefined && object.header !== null) {
            message.header = exports.Header.fromPartial(object.header);
        }
        if (object.data !== undefined && object.data !== null) {
            message.data = types_1.Data.fromPartial(object.data);
        }
        if (object.evidence !== undefined && object.evidence !== null) {
            message.evidence = evidence_1.EvidenceList.fromPartial(object.evidence);
        }
        if (object.lastCommit !== undefined && object.lastCommit !== null) {
            message.lastCommit = types_1.Commit.fromPartial(object.lastCommit);
        }
        return message;
    },
};
function createBaseHeader() {
    return {
        version: types_2.Consensus.fromPartial({}),
        chainId: "",
        height: BigInt(0),
        time: timestamp_1.Timestamp.fromPartial({}),
        lastBlockId: types_1.BlockID.fromPartial({}),
        lastCommitHash: new Uint8Array(),
        dataHash: new Uint8Array(),
        validatorsHash: new Uint8Array(),
        nextValidatorsHash: new Uint8Array(),
        consensusHash: new Uint8Array(),
        appHash: new Uint8Array(),
        lastResultsHash: new Uint8Array(),
        evidenceHash: new Uint8Array(),
        proposerAddress: "",
    };
}
exports.Header = {
    typeUrl: "/cosmos.base.tendermint.v1beta1.Header",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.version !== undefined) {
            types_2.Consensus.encode(message.version, writer.uint32(10).fork()).ldelim();
        }
        if (message.chainId !== "") {
            writer.uint32(18).string(message.chainId);
        }
        if (message.height !== BigInt(0)) {
            writer.uint32(24).int64(message.height);
        }
        if (message.time !== undefined) {
            timestamp_1.Timestamp.encode(message.time, writer.uint32(34).fork()).ldelim();
        }
        if (message.lastBlockId !== undefined) {
            types_1.BlockID.encode(message.lastBlockId, writer.uint32(42).fork()).ldelim();
        }
        if (message.lastCommitHash.length !== 0) {
            writer.uint32(50).bytes(message.lastCommitHash);
        }
        if (message.dataHash.length !== 0) {
            writer.uint32(58).bytes(message.dataHash);
        }
        if (message.validatorsHash.length !== 0) {
            writer.uint32(66).bytes(message.validatorsHash);
        }
        if (message.nextValidatorsHash.length !== 0) {
            writer.uint32(74).bytes(message.nextValidatorsHash);
        }
        if (message.consensusHash.length !== 0) {
            writer.uint32(82).bytes(message.consensusHash);
        }
        if (message.appHash.length !== 0) {
            writer.uint32(90).bytes(message.appHash);
        }
        if (message.lastResultsHash.length !== 0) {
            writer.uint32(98).bytes(message.lastResultsHash);
        }
        if (message.evidenceHash.length !== 0) {
            writer.uint32(106).bytes(message.evidenceHash);
        }
        if (message.proposerAddress !== "") {
            writer.uint32(114).string(message.proposerAddress);
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
                    message.version = types_2.Consensus.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.chainId = reader.string();
                    break;
                case 3:
                    message.height = reader.int64();
                    break;
                case 4:
                    message.time = timestamp_1.Timestamp.decode(reader, reader.uint32());
                    break;
                case 5:
                    message.lastBlockId = types_1.BlockID.decode(reader, reader.uint32());
                    break;
                case 6:
                    message.lastCommitHash = reader.bytes();
                    break;
                case 7:
                    message.dataHash = reader.bytes();
                    break;
                case 8:
                    message.validatorsHash = reader.bytes();
                    break;
                case 9:
                    message.nextValidatorsHash = reader.bytes();
                    break;
                case 10:
                    message.consensusHash = reader.bytes();
                    break;
                case 11:
                    message.appHash = reader.bytes();
                    break;
                case 12:
                    message.lastResultsHash = reader.bytes();
                    break;
                case 13:
                    message.evidenceHash = reader.bytes();
                    break;
                case 14:
                    message.proposerAddress = reader.string();
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
        if ((0, helpers_1.isSet)(object.version))
            obj.version = types_2.Consensus.fromJSON(object.version);
        if ((0, helpers_1.isSet)(object.chainId))
            obj.chainId = String(object.chainId);
        if ((0, helpers_1.isSet)(object.height))
            obj.height = BigInt(object.height.toString());
        if ((0, helpers_1.isSet)(object.time))
            obj.time = (0, helpers_1.fromJsonTimestamp)(object.time);
        if ((0, helpers_1.isSet)(object.lastBlockId))
            obj.lastBlockId = types_1.BlockID.fromJSON(object.lastBlockId);
        if ((0, helpers_1.isSet)(object.lastCommitHash))
            obj.lastCommitHash = (0, helpers_1.bytesFromBase64)(object.lastCommitHash);
        if ((0, helpers_1.isSet)(object.dataHash))
            obj.dataHash = (0, helpers_1.bytesFromBase64)(object.dataHash);
        if ((0, helpers_1.isSet)(object.validatorsHash))
            obj.validatorsHash = (0, helpers_1.bytesFromBase64)(object.validatorsHash);
        if ((0, helpers_1.isSet)(object.nextValidatorsHash))
            obj.nextValidatorsHash = (0, helpers_1.bytesFromBase64)(object.nextValidatorsHash);
        if ((0, helpers_1.isSet)(object.consensusHash))
            obj.consensusHash = (0, helpers_1.bytesFromBase64)(object.consensusHash);
        if ((0, helpers_1.isSet)(object.appHash))
            obj.appHash = (0, helpers_1.bytesFromBase64)(object.appHash);
        if ((0, helpers_1.isSet)(object.lastResultsHash))
            obj.lastResultsHash = (0, helpers_1.bytesFromBase64)(object.lastResultsHash);
        if ((0, helpers_1.isSet)(object.evidenceHash))
            obj.evidenceHash = (0, helpers_1.bytesFromBase64)(object.evidenceHash);
        if ((0, helpers_1.isSet)(object.proposerAddress))
            obj.proposerAddress = String(object.proposerAddress);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.version !== undefined &&
            (obj.version = message.version ? types_2.Consensus.toJSON(message.version) : undefined);
        message.chainId !== undefined && (obj.chainId = message.chainId);
        message.height !== undefined && (obj.height = (message.height || BigInt(0)).toString());
        message.time !== undefined && (obj.time = (0, helpers_1.fromTimestamp)(message.time).toISOString());
        message.lastBlockId !== undefined &&
            (obj.lastBlockId = message.lastBlockId ? types_1.BlockID.toJSON(message.lastBlockId) : undefined);
        message.lastCommitHash !== undefined &&
            (obj.lastCommitHash = (0, helpers_1.base64FromBytes)(message.lastCommitHash !== undefined ? message.lastCommitHash : new Uint8Array()));
        message.dataHash !== undefined &&
            (obj.dataHash = (0, helpers_1.base64FromBytes)(message.dataHash !== undefined ? message.dataHash : new Uint8Array()));
        message.validatorsHash !== undefined &&
            (obj.validatorsHash = (0, helpers_1.base64FromBytes)(message.validatorsHash !== undefined ? message.validatorsHash : new Uint8Array()));
        message.nextValidatorsHash !== undefined &&
            (obj.nextValidatorsHash = (0, helpers_1.base64FromBytes)(message.nextValidatorsHash !== undefined ? message.nextValidatorsHash : new Uint8Array()));
        message.consensusHash !== undefined &&
            (obj.consensusHash = (0, helpers_1.base64FromBytes)(message.consensusHash !== undefined ? message.consensusHash : new Uint8Array()));
        message.appHash !== undefined &&
            (obj.appHash = (0, helpers_1.base64FromBytes)(message.appHash !== undefined ? message.appHash : new Uint8Array()));
        message.lastResultsHash !== undefined &&
            (obj.lastResultsHash = (0, helpers_1.base64FromBytes)(message.lastResultsHash !== undefined ? message.lastResultsHash : new Uint8Array()));
        message.evidenceHash !== undefined &&
            (obj.evidenceHash = (0, helpers_1.base64FromBytes)(message.evidenceHash !== undefined ? message.evidenceHash : new Uint8Array()));
        message.proposerAddress !== undefined && (obj.proposerAddress = message.proposerAddress);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseHeader();
        if (object.version !== undefined && object.version !== null) {
            message.version = types_2.Consensus.fromPartial(object.version);
        }
        message.chainId = object.chainId ?? "";
        if (object.height !== undefined && object.height !== null) {
            message.height = BigInt(object.height.toString());
        }
        if (object.time !== undefined && object.time !== null) {
            message.time = timestamp_1.Timestamp.fromPartial(object.time);
        }
        if (object.lastBlockId !== undefined && object.lastBlockId !== null) {
            message.lastBlockId = types_1.BlockID.fromPartial(object.lastBlockId);
        }
        message.lastCommitHash = object.lastCommitHash ?? new Uint8Array();
        message.dataHash = object.dataHash ?? new Uint8Array();
        message.validatorsHash = object.validatorsHash ?? new Uint8Array();
        message.nextValidatorsHash = object.nextValidatorsHash ?? new Uint8Array();
        message.consensusHash = object.consensusHash ?? new Uint8Array();
        message.appHash = object.appHash ?? new Uint8Array();
        message.lastResultsHash = object.lastResultsHash ?? new Uint8Array();
        message.evidenceHash = object.evidenceHash ?? new Uint8Array();
        message.proposerAddress = object.proposerAddress ?? "";
        return message;
    },
};
//# sourceMappingURL=types.js.map