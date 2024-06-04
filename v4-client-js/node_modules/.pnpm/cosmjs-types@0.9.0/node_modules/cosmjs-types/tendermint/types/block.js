"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Block = exports.protobufPackage = void 0;
/* eslint-disable */
const types_1 = require("./types");
const evidence_1 = require("./evidence");
const binary_1 = require("../../binary");
const helpers_1 = require("../../helpers");
exports.protobufPackage = "tendermint.types";
function createBaseBlock() {
    return {
        header: types_1.Header.fromPartial({}),
        data: types_1.Data.fromPartial({}),
        evidence: evidence_1.EvidenceList.fromPartial({}),
        lastCommit: undefined,
    };
}
exports.Block = {
    typeUrl: "/tendermint.types.Block",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.header !== undefined) {
            types_1.Header.encode(message.header, writer.uint32(10).fork()).ldelim();
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
                    message.header = types_1.Header.decode(reader, reader.uint32());
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
            obj.header = types_1.Header.fromJSON(object.header);
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
        message.header !== undefined && (obj.header = message.header ? types_1.Header.toJSON(message.header) : undefined);
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
            message.header = types_1.Header.fromPartial(object.header);
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
//# sourceMappingURL=block.js.map