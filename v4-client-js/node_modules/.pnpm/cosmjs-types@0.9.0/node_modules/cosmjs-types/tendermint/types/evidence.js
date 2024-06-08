"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.EvidenceList = exports.LightClientAttackEvidence = exports.DuplicateVoteEvidence = exports.Evidence = exports.protobufPackage = void 0;
/* eslint-disable */
const types_1 = require("./types");
const timestamp_1 = require("../../google/protobuf/timestamp");
const validator_1 = require("./validator");
const binary_1 = require("../../binary");
const helpers_1 = require("../../helpers");
exports.protobufPackage = "tendermint.types";
function createBaseEvidence() {
    return {
        duplicateVoteEvidence: undefined,
        lightClientAttackEvidence: undefined,
    };
}
exports.Evidence = {
    typeUrl: "/tendermint.types.Evidence",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.duplicateVoteEvidence !== undefined) {
            exports.DuplicateVoteEvidence.encode(message.duplicateVoteEvidence, writer.uint32(10).fork()).ldelim();
        }
        if (message.lightClientAttackEvidence !== undefined) {
            exports.LightClientAttackEvidence.encode(message.lightClientAttackEvidence, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseEvidence();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.duplicateVoteEvidence = exports.DuplicateVoteEvidence.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.lightClientAttackEvidence = exports.LightClientAttackEvidence.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseEvidence();
        if ((0, helpers_1.isSet)(object.duplicateVoteEvidence))
            obj.duplicateVoteEvidence = exports.DuplicateVoteEvidence.fromJSON(object.duplicateVoteEvidence);
        if ((0, helpers_1.isSet)(object.lightClientAttackEvidence))
            obj.lightClientAttackEvidence = exports.LightClientAttackEvidence.fromJSON(object.lightClientAttackEvidence);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.duplicateVoteEvidence !== undefined &&
            (obj.duplicateVoteEvidence = message.duplicateVoteEvidence
                ? exports.DuplicateVoteEvidence.toJSON(message.duplicateVoteEvidence)
                : undefined);
        message.lightClientAttackEvidence !== undefined &&
            (obj.lightClientAttackEvidence = message.lightClientAttackEvidence
                ? exports.LightClientAttackEvidence.toJSON(message.lightClientAttackEvidence)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseEvidence();
        if (object.duplicateVoteEvidence !== undefined && object.duplicateVoteEvidence !== null) {
            message.duplicateVoteEvidence = exports.DuplicateVoteEvidence.fromPartial(object.duplicateVoteEvidence);
        }
        if (object.lightClientAttackEvidence !== undefined && object.lightClientAttackEvidence !== null) {
            message.lightClientAttackEvidence = exports.LightClientAttackEvidence.fromPartial(object.lightClientAttackEvidence);
        }
        return message;
    },
};
function createBaseDuplicateVoteEvidence() {
    return {
        voteA: undefined,
        voteB: undefined,
        totalVotingPower: BigInt(0),
        validatorPower: BigInt(0),
        timestamp: timestamp_1.Timestamp.fromPartial({}),
    };
}
exports.DuplicateVoteEvidence = {
    typeUrl: "/tendermint.types.DuplicateVoteEvidence",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.voteA !== undefined) {
            types_1.Vote.encode(message.voteA, writer.uint32(10).fork()).ldelim();
        }
        if (message.voteB !== undefined) {
            types_1.Vote.encode(message.voteB, writer.uint32(18).fork()).ldelim();
        }
        if (message.totalVotingPower !== BigInt(0)) {
            writer.uint32(24).int64(message.totalVotingPower);
        }
        if (message.validatorPower !== BigInt(0)) {
            writer.uint32(32).int64(message.validatorPower);
        }
        if (message.timestamp !== undefined) {
            timestamp_1.Timestamp.encode(message.timestamp, writer.uint32(42).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseDuplicateVoteEvidence();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.voteA = types_1.Vote.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.voteB = types_1.Vote.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.totalVotingPower = reader.int64();
                    break;
                case 4:
                    message.validatorPower = reader.int64();
                    break;
                case 5:
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
        const obj = createBaseDuplicateVoteEvidence();
        if ((0, helpers_1.isSet)(object.voteA))
            obj.voteA = types_1.Vote.fromJSON(object.voteA);
        if ((0, helpers_1.isSet)(object.voteB))
            obj.voteB = types_1.Vote.fromJSON(object.voteB);
        if ((0, helpers_1.isSet)(object.totalVotingPower))
            obj.totalVotingPower = BigInt(object.totalVotingPower.toString());
        if ((0, helpers_1.isSet)(object.validatorPower))
            obj.validatorPower = BigInt(object.validatorPower.toString());
        if ((0, helpers_1.isSet)(object.timestamp))
            obj.timestamp = (0, helpers_1.fromJsonTimestamp)(object.timestamp);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.voteA !== undefined && (obj.voteA = message.voteA ? types_1.Vote.toJSON(message.voteA) : undefined);
        message.voteB !== undefined && (obj.voteB = message.voteB ? types_1.Vote.toJSON(message.voteB) : undefined);
        message.totalVotingPower !== undefined &&
            (obj.totalVotingPower = (message.totalVotingPower || BigInt(0)).toString());
        message.validatorPower !== undefined &&
            (obj.validatorPower = (message.validatorPower || BigInt(0)).toString());
        message.timestamp !== undefined && (obj.timestamp = (0, helpers_1.fromTimestamp)(message.timestamp).toISOString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseDuplicateVoteEvidence();
        if (object.voteA !== undefined && object.voteA !== null) {
            message.voteA = types_1.Vote.fromPartial(object.voteA);
        }
        if (object.voteB !== undefined && object.voteB !== null) {
            message.voteB = types_1.Vote.fromPartial(object.voteB);
        }
        if (object.totalVotingPower !== undefined && object.totalVotingPower !== null) {
            message.totalVotingPower = BigInt(object.totalVotingPower.toString());
        }
        if (object.validatorPower !== undefined && object.validatorPower !== null) {
            message.validatorPower = BigInt(object.validatorPower.toString());
        }
        if (object.timestamp !== undefined && object.timestamp !== null) {
            message.timestamp = timestamp_1.Timestamp.fromPartial(object.timestamp);
        }
        return message;
    },
};
function createBaseLightClientAttackEvidence() {
    return {
        conflictingBlock: undefined,
        commonHeight: BigInt(0),
        byzantineValidators: [],
        totalVotingPower: BigInt(0),
        timestamp: timestamp_1.Timestamp.fromPartial({}),
    };
}
exports.LightClientAttackEvidence = {
    typeUrl: "/tendermint.types.LightClientAttackEvidence",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.conflictingBlock !== undefined) {
            types_1.LightBlock.encode(message.conflictingBlock, writer.uint32(10).fork()).ldelim();
        }
        if (message.commonHeight !== BigInt(0)) {
            writer.uint32(16).int64(message.commonHeight);
        }
        for (const v of message.byzantineValidators) {
            validator_1.Validator.encode(v, writer.uint32(26).fork()).ldelim();
        }
        if (message.totalVotingPower !== BigInt(0)) {
            writer.uint32(32).int64(message.totalVotingPower);
        }
        if (message.timestamp !== undefined) {
            timestamp_1.Timestamp.encode(message.timestamp, writer.uint32(42).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseLightClientAttackEvidence();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.conflictingBlock = types_1.LightBlock.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.commonHeight = reader.int64();
                    break;
                case 3:
                    message.byzantineValidators.push(validator_1.Validator.decode(reader, reader.uint32()));
                    break;
                case 4:
                    message.totalVotingPower = reader.int64();
                    break;
                case 5:
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
        const obj = createBaseLightClientAttackEvidence();
        if ((0, helpers_1.isSet)(object.conflictingBlock))
            obj.conflictingBlock = types_1.LightBlock.fromJSON(object.conflictingBlock);
        if ((0, helpers_1.isSet)(object.commonHeight))
            obj.commonHeight = BigInt(object.commonHeight.toString());
        if (Array.isArray(object?.byzantineValidators))
            obj.byzantineValidators = object.byzantineValidators.map((e) => validator_1.Validator.fromJSON(e));
        if ((0, helpers_1.isSet)(object.totalVotingPower))
            obj.totalVotingPower = BigInt(object.totalVotingPower.toString());
        if ((0, helpers_1.isSet)(object.timestamp))
            obj.timestamp = (0, helpers_1.fromJsonTimestamp)(object.timestamp);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.conflictingBlock !== undefined &&
            (obj.conflictingBlock = message.conflictingBlock
                ? types_1.LightBlock.toJSON(message.conflictingBlock)
                : undefined);
        message.commonHeight !== undefined && (obj.commonHeight = (message.commonHeight || BigInt(0)).toString());
        if (message.byzantineValidators) {
            obj.byzantineValidators = message.byzantineValidators.map((e) => (e ? validator_1.Validator.toJSON(e) : undefined));
        }
        else {
            obj.byzantineValidators = [];
        }
        message.totalVotingPower !== undefined &&
            (obj.totalVotingPower = (message.totalVotingPower || BigInt(0)).toString());
        message.timestamp !== undefined && (obj.timestamp = (0, helpers_1.fromTimestamp)(message.timestamp).toISOString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseLightClientAttackEvidence();
        if (object.conflictingBlock !== undefined && object.conflictingBlock !== null) {
            message.conflictingBlock = types_1.LightBlock.fromPartial(object.conflictingBlock);
        }
        if (object.commonHeight !== undefined && object.commonHeight !== null) {
            message.commonHeight = BigInt(object.commonHeight.toString());
        }
        message.byzantineValidators = object.byzantineValidators?.map((e) => validator_1.Validator.fromPartial(e)) || [];
        if (object.totalVotingPower !== undefined && object.totalVotingPower !== null) {
            message.totalVotingPower = BigInt(object.totalVotingPower.toString());
        }
        if (object.timestamp !== undefined && object.timestamp !== null) {
            message.timestamp = timestamp_1.Timestamp.fromPartial(object.timestamp);
        }
        return message;
    },
};
function createBaseEvidenceList() {
    return {
        evidence: [],
    };
}
exports.EvidenceList = {
    typeUrl: "/tendermint.types.EvidenceList",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.evidence) {
            exports.Evidence.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseEvidenceList();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.evidence.push(exports.Evidence.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseEvidenceList();
        if (Array.isArray(object?.evidence))
            obj.evidence = object.evidence.map((e) => exports.Evidence.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.evidence) {
            obj.evidence = message.evidence.map((e) => (e ? exports.Evidence.toJSON(e) : undefined));
        }
        else {
            obj.evidence = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseEvidenceList();
        message.evidence = object.evidence?.map((e) => exports.Evidence.fromPartial(e)) || [];
        return message;
    },
};
//# sourceMappingURL=evidence.js.map