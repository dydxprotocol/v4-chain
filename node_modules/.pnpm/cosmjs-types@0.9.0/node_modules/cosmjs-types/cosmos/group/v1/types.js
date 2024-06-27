"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Vote = exports.TallyResult = exports.Proposal = exports.GroupPolicyInfo = exports.GroupMember = exports.GroupInfo = exports.DecisionPolicyWindows = exports.PercentageDecisionPolicy = exports.ThresholdDecisionPolicy = exports.MemberRequest = exports.Member = exports.proposalExecutorResultToJSON = exports.proposalExecutorResultFromJSON = exports.ProposalExecutorResult = exports.proposalStatusToJSON = exports.proposalStatusFromJSON = exports.ProposalStatus = exports.voteOptionToJSON = exports.voteOptionFromJSON = exports.VoteOption = exports.protobufPackage = void 0;
/* eslint-disable */
const timestamp_1 = require("../../../google/protobuf/timestamp");
const duration_1 = require("../../../google/protobuf/duration");
const any_1 = require("../../../google/protobuf/any");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.group.v1";
/** VoteOption enumerates the valid vote options for a given proposal. */
var VoteOption;
(function (VoteOption) {
    /**
     * VOTE_OPTION_UNSPECIFIED - VOTE_OPTION_UNSPECIFIED defines an unspecified vote option which will
     * return an error.
     */
    VoteOption[VoteOption["VOTE_OPTION_UNSPECIFIED"] = 0] = "VOTE_OPTION_UNSPECIFIED";
    /** VOTE_OPTION_YES - VOTE_OPTION_YES defines a yes vote option. */
    VoteOption[VoteOption["VOTE_OPTION_YES"] = 1] = "VOTE_OPTION_YES";
    /** VOTE_OPTION_ABSTAIN - VOTE_OPTION_ABSTAIN defines an abstain vote option. */
    VoteOption[VoteOption["VOTE_OPTION_ABSTAIN"] = 2] = "VOTE_OPTION_ABSTAIN";
    /** VOTE_OPTION_NO - VOTE_OPTION_NO defines a no vote option. */
    VoteOption[VoteOption["VOTE_OPTION_NO"] = 3] = "VOTE_OPTION_NO";
    /** VOTE_OPTION_NO_WITH_VETO - VOTE_OPTION_NO_WITH_VETO defines a no with veto vote option. */
    VoteOption[VoteOption["VOTE_OPTION_NO_WITH_VETO"] = 4] = "VOTE_OPTION_NO_WITH_VETO";
    VoteOption[VoteOption["UNRECOGNIZED"] = -1] = "UNRECOGNIZED";
})(VoteOption || (exports.VoteOption = VoteOption = {}));
function voteOptionFromJSON(object) {
    switch (object) {
        case 0:
        case "VOTE_OPTION_UNSPECIFIED":
            return VoteOption.VOTE_OPTION_UNSPECIFIED;
        case 1:
        case "VOTE_OPTION_YES":
            return VoteOption.VOTE_OPTION_YES;
        case 2:
        case "VOTE_OPTION_ABSTAIN":
            return VoteOption.VOTE_OPTION_ABSTAIN;
        case 3:
        case "VOTE_OPTION_NO":
            return VoteOption.VOTE_OPTION_NO;
        case 4:
        case "VOTE_OPTION_NO_WITH_VETO":
            return VoteOption.VOTE_OPTION_NO_WITH_VETO;
        case -1:
        case "UNRECOGNIZED":
        default:
            return VoteOption.UNRECOGNIZED;
    }
}
exports.voteOptionFromJSON = voteOptionFromJSON;
function voteOptionToJSON(object) {
    switch (object) {
        case VoteOption.VOTE_OPTION_UNSPECIFIED:
            return "VOTE_OPTION_UNSPECIFIED";
        case VoteOption.VOTE_OPTION_YES:
            return "VOTE_OPTION_YES";
        case VoteOption.VOTE_OPTION_ABSTAIN:
            return "VOTE_OPTION_ABSTAIN";
        case VoteOption.VOTE_OPTION_NO:
            return "VOTE_OPTION_NO";
        case VoteOption.VOTE_OPTION_NO_WITH_VETO:
            return "VOTE_OPTION_NO_WITH_VETO";
        case VoteOption.UNRECOGNIZED:
        default:
            return "UNRECOGNIZED";
    }
}
exports.voteOptionToJSON = voteOptionToJSON;
/** ProposalStatus defines proposal statuses. */
var ProposalStatus;
(function (ProposalStatus) {
    /** PROPOSAL_STATUS_UNSPECIFIED - An empty value is invalid and not allowed. */
    ProposalStatus[ProposalStatus["PROPOSAL_STATUS_UNSPECIFIED"] = 0] = "PROPOSAL_STATUS_UNSPECIFIED";
    /** PROPOSAL_STATUS_SUBMITTED - Initial status of a proposal when submitted. */
    ProposalStatus[ProposalStatus["PROPOSAL_STATUS_SUBMITTED"] = 1] = "PROPOSAL_STATUS_SUBMITTED";
    /**
     * PROPOSAL_STATUS_ACCEPTED - Final status of a proposal when the final tally is done and the outcome
     * passes the group policy's decision policy.
     */
    ProposalStatus[ProposalStatus["PROPOSAL_STATUS_ACCEPTED"] = 2] = "PROPOSAL_STATUS_ACCEPTED";
    /**
     * PROPOSAL_STATUS_REJECTED - Final status of a proposal when the final tally is done and the outcome
     * is rejected by the group policy's decision policy.
     */
    ProposalStatus[ProposalStatus["PROPOSAL_STATUS_REJECTED"] = 3] = "PROPOSAL_STATUS_REJECTED";
    /**
     * PROPOSAL_STATUS_ABORTED - Final status of a proposal when the group policy is modified before the
     * final tally.
     */
    ProposalStatus[ProposalStatus["PROPOSAL_STATUS_ABORTED"] = 4] = "PROPOSAL_STATUS_ABORTED";
    /**
     * PROPOSAL_STATUS_WITHDRAWN - A proposal can be withdrawn before the voting start time by the owner.
     * When this happens the final status is Withdrawn.
     */
    ProposalStatus[ProposalStatus["PROPOSAL_STATUS_WITHDRAWN"] = 5] = "PROPOSAL_STATUS_WITHDRAWN";
    ProposalStatus[ProposalStatus["UNRECOGNIZED"] = -1] = "UNRECOGNIZED";
})(ProposalStatus || (exports.ProposalStatus = ProposalStatus = {}));
function proposalStatusFromJSON(object) {
    switch (object) {
        case 0:
        case "PROPOSAL_STATUS_UNSPECIFIED":
            return ProposalStatus.PROPOSAL_STATUS_UNSPECIFIED;
        case 1:
        case "PROPOSAL_STATUS_SUBMITTED":
            return ProposalStatus.PROPOSAL_STATUS_SUBMITTED;
        case 2:
        case "PROPOSAL_STATUS_ACCEPTED":
            return ProposalStatus.PROPOSAL_STATUS_ACCEPTED;
        case 3:
        case "PROPOSAL_STATUS_REJECTED":
            return ProposalStatus.PROPOSAL_STATUS_REJECTED;
        case 4:
        case "PROPOSAL_STATUS_ABORTED":
            return ProposalStatus.PROPOSAL_STATUS_ABORTED;
        case 5:
        case "PROPOSAL_STATUS_WITHDRAWN":
            return ProposalStatus.PROPOSAL_STATUS_WITHDRAWN;
        case -1:
        case "UNRECOGNIZED":
        default:
            return ProposalStatus.UNRECOGNIZED;
    }
}
exports.proposalStatusFromJSON = proposalStatusFromJSON;
function proposalStatusToJSON(object) {
    switch (object) {
        case ProposalStatus.PROPOSAL_STATUS_UNSPECIFIED:
            return "PROPOSAL_STATUS_UNSPECIFIED";
        case ProposalStatus.PROPOSAL_STATUS_SUBMITTED:
            return "PROPOSAL_STATUS_SUBMITTED";
        case ProposalStatus.PROPOSAL_STATUS_ACCEPTED:
            return "PROPOSAL_STATUS_ACCEPTED";
        case ProposalStatus.PROPOSAL_STATUS_REJECTED:
            return "PROPOSAL_STATUS_REJECTED";
        case ProposalStatus.PROPOSAL_STATUS_ABORTED:
            return "PROPOSAL_STATUS_ABORTED";
        case ProposalStatus.PROPOSAL_STATUS_WITHDRAWN:
            return "PROPOSAL_STATUS_WITHDRAWN";
        case ProposalStatus.UNRECOGNIZED:
        default:
            return "UNRECOGNIZED";
    }
}
exports.proposalStatusToJSON = proposalStatusToJSON;
/** ProposalExecutorResult defines types of proposal executor results. */
var ProposalExecutorResult;
(function (ProposalExecutorResult) {
    /** PROPOSAL_EXECUTOR_RESULT_UNSPECIFIED - An empty value is not allowed. */
    ProposalExecutorResult[ProposalExecutorResult["PROPOSAL_EXECUTOR_RESULT_UNSPECIFIED"] = 0] = "PROPOSAL_EXECUTOR_RESULT_UNSPECIFIED";
    /** PROPOSAL_EXECUTOR_RESULT_NOT_RUN - We have not yet run the executor. */
    ProposalExecutorResult[ProposalExecutorResult["PROPOSAL_EXECUTOR_RESULT_NOT_RUN"] = 1] = "PROPOSAL_EXECUTOR_RESULT_NOT_RUN";
    /** PROPOSAL_EXECUTOR_RESULT_SUCCESS - The executor was successful and proposed action updated state. */
    ProposalExecutorResult[ProposalExecutorResult["PROPOSAL_EXECUTOR_RESULT_SUCCESS"] = 2] = "PROPOSAL_EXECUTOR_RESULT_SUCCESS";
    /** PROPOSAL_EXECUTOR_RESULT_FAILURE - The executor returned an error and proposed action didn't update state. */
    ProposalExecutorResult[ProposalExecutorResult["PROPOSAL_EXECUTOR_RESULT_FAILURE"] = 3] = "PROPOSAL_EXECUTOR_RESULT_FAILURE";
    ProposalExecutorResult[ProposalExecutorResult["UNRECOGNIZED"] = -1] = "UNRECOGNIZED";
})(ProposalExecutorResult || (exports.ProposalExecutorResult = ProposalExecutorResult = {}));
function proposalExecutorResultFromJSON(object) {
    switch (object) {
        case 0:
        case "PROPOSAL_EXECUTOR_RESULT_UNSPECIFIED":
            return ProposalExecutorResult.PROPOSAL_EXECUTOR_RESULT_UNSPECIFIED;
        case 1:
        case "PROPOSAL_EXECUTOR_RESULT_NOT_RUN":
            return ProposalExecutorResult.PROPOSAL_EXECUTOR_RESULT_NOT_RUN;
        case 2:
        case "PROPOSAL_EXECUTOR_RESULT_SUCCESS":
            return ProposalExecutorResult.PROPOSAL_EXECUTOR_RESULT_SUCCESS;
        case 3:
        case "PROPOSAL_EXECUTOR_RESULT_FAILURE":
            return ProposalExecutorResult.PROPOSAL_EXECUTOR_RESULT_FAILURE;
        case -1:
        case "UNRECOGNIZED":
        default:
            return ProposalExecutorResult.UNRECOGNIZED;
    }
}
exports.proposalExecutorResultFromJSON = proposalExecutorResultFromJSON;
function proposalExecutorResultToJSON(object) {
    switch (object) {
        case ProposalExecutorResult.PROPOSAL_EXECUTOR_RESULT_UNSPECIFIED:
            return "PROPOSAL_EXECUTOR_RESULT_UNSPECIFIED";
        case ProposalExecutorResult.PROPOSAL_EXECUTOR_RESULT_NOT_RUN:
            return "PROPOSAL_EXECUTOR_RESULT_NOT_RUN";
        case ProposalExecutorResult.PROPOSAL_EXECUTOR_RESULT_SUCCESS:
            return "PROPOSAL_EXECUTOR_RESULT_SUCCESS";
        case ProposalExecutorResult.PROPOSAL_EXECUTOR_RESULT_FAILURE:
            return "PROPOSAL_EXECUTOR_RESULT_FAILURE";
        case ProposalExecutorResult.UNRECOGNIZED:
        default:
            return "UNRECOGNIZED";
    }
}
exports.proposalExecutorResultToJSON = proposalExecutorResultToJSON;
function createBaseMember() {
    return {
        address: "",
        weight: "",
        metadata: "",
        addedAt: timestamp_1.Timestamp.fromPartial({}),
    };
}
exports.Member = {
    typeUrl: "/cosmos.group.v1.Member",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.address !== "") {
            writer.uint32(10).string(message.address);
        }
        if (message.weight !== "") {
            writer.uint32(18).string(message.weight);
        }
        if (message.metadata !== "") {
            writer.uint32(26).string(message.metadata);
        }
        if (message.addedAt !== undefined) {
            timestamp_1.Timestamp.encode(message.addedAt, writer.uint32(34).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMember();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.address = reader.string();
                    break;
                case 2:
                    message.weight = reader.string();
                    break;
                case 3:
                    message.metadata = reader.string();
                    break;
                case 4:
                    message.addedAt = timestamp_1.Timestamp.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMember();
        if ((0, helpers_1.isSet)(object.address))
            obj.address = String(object.address);
        if ((0, helpers_1.isSet)(object.weight))
            obj.weight = String(object.weight);
        if ((0, helpers_1.isSet)(object.metadata))
            obj.metadata = String(object.metadata);
        if ((0, helpers_1.isSet)(object.addedAt))
            obj.addedAt = (0, helpers_1.fromJsonTimestamp)(object.addedAt);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.address !== undefined && (obj.address = message.address);
        message.weight !== undefined && (obj.weight = message.weight);
        message.metadata !== undefined && (obj.metadata = message.metadata);
        message.addedAt !== undefined && (obj.addedAt = (0, helpers_1.fromTimestamp)(message.addedAt).toISOString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMember();
        message.address = object.address ?? "";
        message.weight = object.weight ?? "";
        message.metadata = object.metadata ?? "";
        if (object.addedAt !== undefined && object.addedAt !== null) {
            message.addedAt = timestamp_1.Timestamp.fromPartial(object.addedAt);
        }
        return message;
    },
};
function createBaseMemberRequest() {
    return {
        address: "",
        weight: "",
        metadata: "",
    };
}
exports.MemberRequest = {
    typeUrl: "/cosmos.group.v1.MemberRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.address !== "") {
            writer.uint32(10).string(message.address);
        }
        if (message.weight !== "") {
            writer.uint32(18).string(message.weight);
        }
        if (message.metadata !== "") {
            writer.uint32(26).string(message.metadata);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMemberRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.address = reader.string();
                    break;
                case 2:
                    message.weight = reader.string();
                    break;
                case 3:
                    message.metadata = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMemberRequest();
        if ((0, helpers_1.isSet)(object.address))
            obj.address = String(object.address);
        if ((0, helpers_1.isSet)(object.weight))
            obj.weight = String(object.weight);
        if ((0, helpers_1.isSet)(object.metadata))
            obj.metadata = String(object.metadata);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.address !== undefined && (obj.address = message.address);
        message.weight !== undefined && (obj.weight = message.weight);
        message.metadata !== undefined && (obj.metadata = message.metadata);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMemberRequest();
        message.address = object.address ?? "";
        message.weight = object.weight ?? "";
        message.metadata = object.metadata ?? "";
        return message;
    },
};
function createBaseThresholdDecisionPolicy() {
    return {
        threshold: "",
        windows: undefined,
    };
}
exports.ThresholdDecisionPolicy = {
    typeUrl: "/cosmos.group.v1.ThresholdDecisionPolicy",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.threshold !== "") {
            writer.uint32(10).string(message.threshold);
        }
        if (message.windows !== undefined) {
            exports.DecisionPolicyWindows.encode(message.windows, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseThresholdDecisionPolicy();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.threshold = reader.string();
                    break;
                case 2:
                    message.windows = exports.DecisionPolicyWindows.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseThresholdDecisionPolicy();
        if ((0, helpers_1.isSet)(object.threshold))
            obj.threshold = String(object.threshold);
        if ((0, helpers_1.isSet)(object.windows))
            obj.windows = exports.DecisionPolicyWindows.fromJSON(object.windows);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.threshold !== undefined && (obj.threshold = message.threshold);
        message.windows !== undefined &&
            (obj.windows = message.windows ? exports.DecisionPolicyWindows.toJSON(message.windows) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseThresholdDecisionPolicy();
        message.threshold = object.threshold ?? "";
        if (object.windows !== undefined && object.windows !== null) {
            message.windows = exports.DecisionPolicyWindows.fromPartial(object.windows);
        }
        return message;
    },
};
function createBasePercentageDecisionPolicy() {
    return {
        percentage: "",
        windows: undefined,
    };
}
exports.PercentageDecisionPolicy = {
    typeUrl: "/cosmos.group.v1.PercentageDecisionPolicy",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.percentage !== "") {
            writer.uint32(10).string(message.percentage);
        }
        if (message.windows !== undefined) {
            exports.DecisionPolicyWindows.encode(message.windows, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBasePercentageDecisionPolicy();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.percentage = reader.string();
                    break;
                case 2:
                    message.windows = exports.DecisionPolicyWindows.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBasePercentageDecisionPolicy();
        if ((0, helpers_1.isSet)(object.percentage))
            obj.percentage = String(object.percentage);
        if ((0, helpers_1.isSet)(object.windows))
            obj.windows = exports.DecisionPolicyWindows.fromJSON(object.windows);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.percentage !== undefined && (obj.percentage = message.percentage);
        message.windows !== undefined &&
            (obj.windows = message.windows ? exports.DecisionPolicyWindows.toJSON(message.windows) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBasePercentageDecisionPolicy();
        message.percentage = object.percentage ?? "";
        if (object.windows !== undefined && object.windows !== null) {
            message.windows = exports.DecisionPolicyWindows.fromPartial(object.windows);
        }
        return message;
    },
};
function createBaseDecisionPolicyWindows() {
    return {
        votingPeriod: duration_1.Duration.fromPartial({}),
        minExecutionPeriod: duration_1.Duration.fromPartial({}),
    };
}
exports.DecisionPolicyWindows = {
    typeUrl: "/cosmos.group.v1.DecisionPolicyWindows",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.votingPeriod !== undefined) {
            duration_1.Duration.encode(message.votingPeriod, writer.uint32(10).fork()).ldelim();
        }
        if (message.minExecutionPeriod !== undefined) {
            duration_1.Duration.encode(message.minExecutionPeriod, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseDecisionPolicyWindows();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.votingPeriod = duration_1.Duration.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.minExecutionPeriod = duration_1.Duration.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseDecisionPolicyWindows();
        if ((0, helpers_1.isSet)(object.votingPeriod))
            obj.votingPeriod = duration_1.Duration.fromJSON(object.votingPeriod);
        if ((0, helpers_1.isSet)(object.minExecutionPeriod))
            obj.minExecutionPeriod = duration_1.Duration.fromJSON(object.minExecutionPeriod);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.votingPeriod !== undefined &&
            (obj.votingPeriod = message.votingPeriod ? duration_1.Duration.toJSON(message.votingPeriod) : undefined);
        message.minExecutionPeriod !== undefined &&
            (obj.minExecutionPeriod = message.minExecutionPeriod
                ? duration_1.Duration.toJSON(message.minExecutionPeriod)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseDecisionPolicyWindows();
        if (object.votingPeriod !== undefined && object.votingPeriod !== null) {
            message.votingPeriod = duration_1.Duration.fromPartial(object.votingPeriod);
        }
        if (object.minExecutionPeriod !== undefined && object.minExecutionPeriod !== null) {
            message.minExecutionPeriod = duration_1.Duration.fromPartial(object.minExecutionPeriod);
        }
        return message;
    },
};
function createBaseGroupInfo() {
    return {
        id: BigInt(0),
        admin: "",
        metadata: "",
        version: BigInt(0),
        totalWeight: "",
        createdAt: timestamp_1.Timestamp.fromPartial({}),
    };
}
exports.GroupInfo = {
    typeUrl: "/cosmos.group.v1.GroupInfo",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.id !== BigInt(0)) {
            writer.uint32(8).uint64(message.id);
        }
        if (message.admin !== "") {
            writer.uint32(18).string(message.admin);
        }
        if (message.metadata !== "") {
            writer.uint32(26).string(message.metadata);
        }
        if (message.version !== BigInt(0)) {
            writer.uint32(32).uint64(message.version);
        }
        if (message.totalWeight !== "") {
            writer.uint32(42).string(message.totalWeight);
        }
        if (message.createdAt !== undefined) {
            timestamp_1.Timestamp.encode(message.createdAt, writer.uint32(50).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGroupInfo();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.id = reader.uint64();
                    break;
                case 2:
                    message.admin = reader.string();
                    break;
                case 3:
                    message.metadata = reader.string();
                    break;
                case 4:
                    message.version = reader.uint64();
                    break;
                case 5:
                    message.totalWeight = reader.string();
                    break;
                case 6:
                    message.createdAt = timestamp_1.Timestamp.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseGroupInfo();
        if ((0, helpers_1.isSet)(object.id))
            obj.id = BigInt(object.id.toString());
        if ((0, helpers_1.isSet)(object.admin))
            obj.admin = String(object.admin);
        if ((0, helpers_1.isSet)(object.metadata))
            obj.metadata = String(object.metadata);
        if ((0, helpers_1.isSet)(object.version))
            obj.version = BigInt(object.version.toString());
        if ((0, helpers_1.isSet)(object.totalWeight))
            obj.totalWeight = String(object.totalWeight);
        if ((0, helpers_1.isSet)(object.createdAt))
            obj.createdAt = (0, helpers_1.fromJsonTimestamp)(object.createdAt);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.id !== undefined && (obj.id = (message.id || BigInt(0)).toString());
        message.admin !== undefined && (obj.admin = message.admin);
        message.metadata !== undefined && (obj.metadata = message.metadata);
        message.version !== undefined && (obj.version = (message.version || BigInt(0)).toString());
        message.totalWeight !== undefined && (obj.totalWeight = message.totalWeight);
        message.createdAt !== undefined && (obj.createdAt = (0, helpers_1.fromTimestamp)(message.createdAt).toISOString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGroupInfo();
        if (object.id !== undefined && object.id !== null) {
            message.id = BigInt(object.id.toString());
        }
        message.admin = object.admin ?? "";
        message.metadata = object.metadata ?? "";
        if (object.version !== undefined && object.version !== null) {
            message.version = BigInt(object.version.toString());
        }
        message.totalWeight = object.totalWeight ?? "";
        if (object.createdAt !== undefined && object.createdAt !== null) {
            message.createdAt = timestamp_1.Timestamp.fromPartial(object.createdAt);
        }
        return message;
    },
};
function createBaseGroupMember() {
    return {
        groupId: BigInt(0),
        member: undefined,
    };
}
exports.GroupMember = {
    typeUrl: "/cosmos.group.v1.GroupMember",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.groupId !== BigInt(0)) {
            writer.uint32(8).uint64(message.groupId);
        }
        if (message.member !== undefined) {
            exports.Member.encode(message.member, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGroupMember();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.groupId = reader.uint64();
                    break;
                case 2:
                    message.member = exports.Member.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseGroupMember();
        if ((0, helpers_1.isSet)(object.groupId))
            obj.groupId = BigInt(object.groupId.toString());
        if ((0, helpers_1.isSet)(object.member))
            obj.member = exports.Member.fromJSON(object.member);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.groupId !== undefined && (obj.groupId = (message.groupId || BigInt(0)).toString());
        message.member !== undefined && (obj.member = message.member ? exports.Member.toJSON(message.member) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGroupMember();
        if (object.groupId !== undefined && object.groupId !== null) {
            message.groupId = BigInt(object.groupId.toString());
        }
        if (object.member !== undefined && object.member !== null) {
            message.member = exports.Member.fromPartial(object.member);
        }
        return message;
    },
};
function createBaseGroupPolicyInfo() {
    return {
        address: "",
        groupId: BigInt(0),
        admin: "",
        metadata: "",
        version: BigInt(0),
        decisionPolicy: undefined,
        createdAt: timestamp_1.Timestamp.fromPartial({}),
    };
}
exports.GroupPolicyInfo = {
    typeUrl: "/cosmos.group.v1.GroupPolicyInfo",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.address !== "") {
            writer.uint32(10).string(message.address);
        }
        if (message.groupId !== BigInt(0)) {
            writer.uint32(16).uint64(message.groupId);
        }
        if (message.admin !== "") {
            writer.uint32(26).string(message.admin);
        }
        if (message.metadata !== "") {
            writer.uint32(34).string(message.metadata);
        }
        if (message.version !== BigInt(0)) {
            writer.uint32(40).uint64(message.version);
        }
        if (message.decisionPolicy !== undefined) {
            any_1.Any.encode(message.decisionPolicy, writer.uint32(50).fork()).ldelim();
        }
        if (message.createdAt !== undefined) {
            timestamp_1.Timestamp.encode(message.createdAt, writer.uint32(58).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGroupPolicyInfo();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.address = reader.string();
                    break;
                case 2:
                    message.groupId = reader.uint64();
                    break;
                case 3:
                    message.admin = reader.string();
                    break;
                case 4:
                    message.metadata = reader.string();
                    break;
                case 5:
                    message.version = reader.uint64();
                    break;
                case 6:
                    message.decisionPolicy = any_1.Any.decode(reader, reader.uint32());
                    break;
                case 7:
                    message.createdAt = timestamp_1.Timestamp.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseGroupPolicyInfo();
        if ((0, helpers_1.isSet)(object.address))
            obj.address = String(object.address);
        if ((0, helpers_1.isSet)(object.groupId))
            obj.groupId = BigInt(object.groupId.toString());
        if ((0, helpers_1.isSet)(object.admin))
            obj.admin = String(object.admin);
        if ((0, helpers_1.isSet)(object.metadata))
            obj.metadata = String(object.metadata);
        if ((0, helpers_1.isSet)(object.version))
            obj.version = BigInt(object.version.toString());
        if ((0, helpers_1.isSet)(object.decisionPolicy))
            obj.decisionPolicy = any_1.Any.fromJSON(object.decisionPolicy);
        if ((0, helpers_1.isSet)(object.createdAt))
            obj.createdAt = (0, helpers_1.fromJsonTimestamp)(object.createdAt);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.address !== undefined && (obj.address = message.address);
        message.groupId !== undefined && (obj.groupId = (message.groupId || BigInt(0)).toString());
        message.admin !== undefined && (obj.admin = message.admin);
        message.metadata !== undefined && (obj.metadata = message.metadata);
        message.version !== undefined && (obj.version = (message.version || BigInt(0)).toString());
        message.decisionPolicy !== undefined &&
            (obj.decisionPolicy = message.decisionPolicy ? any_1.Any.toJSON(message.decisionPolicy) : undefined);
        message.createdAt !== undefined && (obj.createdAt = (0, helpers_1.fromTimestamp)(message.createdAt).toISOString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGroupPolicyInfo();
        message.address = object.address ?? "";
        if (object.groupId !== undefined && object.groupId !== null) {
            message.groupId = BigInt(object.groupId.toString());
        }
        message.admin = object.admin ?? "";
        message.metadata = object.metadata ?? "";
        if (object.version !== undefined && object.version !== null) {
            message.version = BigInt(object.version.toString());
        }
        if (object.decisionPolicy !== undefined && object.decisionPolicy !== null) {
            message.decisionPolicy = any_1.Any.fromPartial(object.decisionPolicy);
        }
        if (object.createdAt !== undefined && object.createdAt !== null) {
            message.createdAt = timestamp_1.Timestamp.fromPartial(object.createdAt);
        }
        return message;
    },
};
function createBaseProposal() {
    return {
        id: BigInt(0),
        groupPolicyAddress: "",
        metadata: "",
        proposers: [],
        submitTime: timestamp_1.Timestamp.fromPartial({}),
        groupVersion: BigInt(0),
        groupPolicyVersion: BigInt(0),
        status: 0,
        finalTallyResult: exports.TallyResult.fromPartial({}),
        votingPeriodEnd: timestamp_1.Timestamp.fromPartial({}),
        executorResult: 0,
        messages: [],
        title: "",
        summary: "",
    };
}
exports.Proposal = {
    typeUrl: "/cosmos.group.v1.Proposal",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.id !== BigInt(0)) {
            writer.uint32(8).uint64(message.id);
        }
        if (message.groupPolicyAddress !== "") {
            writer.uint32(18).string(message.groupPolicyAddress);
        }
        if (message.metadata !== "") {
            writer.uint32(26).string(message.metadata);
        }
        for (const v of message.proposers) {
            writer.uint32(34).string(v);
        }
        if (message.submitTime !== undefined) {
            timestamp_1.Timestamp.encode(message.submitTime, writer.uint32(42).fork()).ldelim();
        }
        if (message.groupVersion !== BigInt(0)) {
            writer.uint32(48).uint64(message.groupVersion);
        }
        if (message.groupPolicyVersion !== BigInt(0)) {
            writer.uint32(56).uint64(message.groupPolicyVersion);
        }
        if (message.status !== 0) {
            writer.uint32(64).int32(message.status);
        }
        if (message.finalTallyResult !== undefined) {
            exports.TallyResult.encode(message.finalTallyResult, writer.uint32(74).fork()).ldelim();
        }
        if (message.votingPeriodEnd !== undefined) {
            timestamp_1.Timestamp.encode(message.votingPeriodEnd, writer.uint32(82).fork()).ldelim();
        }
        if (message.executorResult !== 0) {
            writer.uint32(88).int32(message.executorResult);
        }
        for (const v of message.messages) {
            any_1.Any.encode(v, writer.uint32(98).fork()).ldelim();
        }
        if (message.title !== "") {
            writer.uint32(106).string(message.title);
        }
        if (message.summary !== "") {
            writer.uint32(114).string(message.summary);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseProposal();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.id = reader.uint64();
                    break;
                case 2:
                    message.groupPolicyAddress = reader.string();
                    break;
                case 3:
                    message.metadata = reader.string();
                    break;
                case 4:
                    message.proposers.push(reader.string());
                    break;
                case 5:
                    message.submitTime = timestamp_1.Timestamp.decode(reader, reader.uint32());
                    break;
                case 6:
                    message.groupVersion = reader.uint64();
                    break;
                case 7:
                    message.groupPolicyVersion = reader.uint64();
                    break;
                case 8:
                    message.status = reader.int32();
                    break;
                case 9:
                    message.finalTallyResult = exports.TallyResult.decode(reader, reader.uint32());
                    break;
                case 10:
                    message.votingPeriodEnd = timestamp_1.Timestamp.decode(reader, reader.uint32());
                    break;
                case 11:
                    message.executorResult = reader.int32();
                    break;
                case 12:
                    message.messages.push(any_1.Any.decode(reader, reader.uint32()));
                    break;
                case 13:
                    message.title = reader.string();
                    break;
                case 14:
                    message.summary = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseProposal();
        if ((0, helpers_1.isSet)(object.id))
            obj.id = BigInt(object.id.toString());
        if ((0, helpers_1.isSet)(object.groupPolicyAddress))
            obj.groupPolicyAddress = String(object.groupPolicyAddress);
        if ((0, helpers_1.isSet)(object.metadata))
            obj.metadata = String(object.metadata);
        if (Array.isArray(object?.proposers))
            obj.proposers = object.proposers.map((e) => String(e));
        if ((0, helpers_1.isSet)(object.submitTime))
            obj.submitTime = (0, helpers_1.fromJsonTimestamp)(object.submitTime);
        if ((0, helpers_1.isSet)(object.groupVersion))
            obj.groupVersion = BigInt(object.groupVersion.toString());
        if ((0, helpers_1.isSet)(object.groupPolicyVersion))
            obj.groupPolicyVersion = BigInt(object.groupPolicyVersion.toString());
        if ((0, helpers_1.isSet)(object.status))
            obj.status = proposalStatusFromJSON(object.status);
        if ((0, helpers_1.isSet)(object.finalTallyResult))
            obj.finalTallyResult = exports.TallyResult.fromJSON(object.finalTallyResult);
        if ((0, helpers_1.isSet)(object.votingPeriodEnd))
            obj.votingPeriodEnd = (0, helpers_1.fromJsonTimestamp)(object.votingPeriodEnd);
        if ((0, helpers_1.isSet)(object.executorResult))
            obj.executorResult = proposalExecutorResultFromJSON(object.executorResult);
        if (Array.isArray(object?.messages))
            obj.messages = object.messages.map((e) => any_1.Any.fromJSON(e));
        if ((0, helpers_1.isSet)(object.title))
            obj.title = String(object.title);
        if ((0, helpers_1.isSet)(object.summary))
            obj.summary = String(object.summary);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.id !== undefined && (obj.id = (message.id || BigInt(0)).toString());
        message.groupPolicyAddress !== undefined && (obj.groupPolicyAddress = message.groupPolicyAddress);
        message.metadata !== undefined && (obj.metadata = message.metadata);
        if (message.proposers) {
            obj.proposers = message.proposers.map((e) => e);
        }
        else {
            obj.proposers = [];
        }
        message.submitTime !== undefined && (obj.submitTime = (0, helpers_1.fromTimestamp)(message.submitTime).toISOString());
        message.groupVersion !== undefined && (obj.groupVersion = (message.groupVersion || BigInt(0)).toString());
        message.groupPolicyVersion !== undefined &&
            (obj.groupPolicyVersion = (message.groupPolicyVersion || BigInt(0)).toString());
        message.status !== undefined && (obj.status = proposalStatusToJSON(message.status));
        message.finalTallyResult !== undefined &&
            (obj.finalTallyResult = message.finalTallyResult
                ? exports.TallyResult.toJSON(message.finalTallyResult)
                : undefined);
        message.votingPeriodEnd !== undefined &&
            (obj.votingPeriodEnd = (0, helpers_1.fromTimestamp)(message.votingPeriodEnd).toISOString());
        message.executorResult !== undefined &&
            (obj.executorResult = proposalExecutorResultToJSON(message.executorResult));
        if (message.messages) {
            obj.messages = message.messages.map((e) => (e ? any_1.Any.toJSON(e) : undefined));
        }
        else {
            obj.messages = [];
        }
        message.title !== undefined && (obj.title = message.title);
        message.summary !== undefined && (obj.summary = message.summary);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseProposal();
        if (object.id !== undefined && object.id !== null) {
            message.id = BigInt(object.id.toString());
        }
        message.groupPolicyAddress = object.groupPolicyAddress ?? "";
        message.metadata = object.metadata ?? "";
        message.proposers = object.proposers?.map((e) => e) || [];
        if (object.submitTime !== undefined && object.submitTime !== null) {
            message.submitTime = timestamp_1.Timestamp.fromPartial(object.submitTime);
        }
        if (object.groupVersion !== undefined && object.groupVersion !== null) {
            message.groupVersion = BigInt(object.groupVersion.toString());
        }
        if (object.groupPolicyVersion !== undefined && object.groupPolicyVersion !== null) {
            message.groupPolicyVersion = BigInt(object.groupPolicyVersion.toString());
        }
        message.status = object.status ?? 0;
        if (object.finalTallyResult !== undefined && object.finalTallyResult !== null) {
            message.finalTallyResult = exports.TallyResult.fromPartial(object.finalTallyResult);
        }
        if (object.votingPeriodEnd !== undefined && object.votingPeriodEnd !== null) {
            message.votingPeriodEnd = timestamp_1.Timestamp.fromPartial(object.votingPeriodEnd);
        }
        message.executorResult = object.executorResult ?? 0;
        message.messages = object.messages?.map((e) => any_1.Any.fromPartial(e)) || [];
        message.title = object.title ?? "";
        message.summary = object.summary ?? "";
        return message;
    },
};
function createBaseTallyResult() {
    return {
        yesCount: "",
        abstainCount: "",
        noCount: "",
        noWithVetoCount: "",
    };
}
exports.TallyResult = {
    typeUrl: "/cosmos.group.v1.TallyResult",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.yesCount !== "") {
            writer.uint32(10).string(message.yesCount);
        }
        if (message.abstainCount !== "") {
            writer.uint32(18).string(message.abstainCount);
        }
        if (message.noCount !== "") {
            writer.uint32(26).string(message.noCount);
        }
        if (message.noWithVetoCount !== "") {
            writer.uint32(34).string(message.noWithVetoCount);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseTallyResult();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.yesCount = reader.string();
                    break;
                case 2:
                    message.abstainCount = reader.string();
                    break;
                case 3:
                    message.noCount = reader.string();
                    break;
                case 4:
                    message.noWithVetoCount = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseTallyResult();
        if ((0, helpers_1.isSet)(object.yesCount))
            obj.yesCount = String(object.yesCount);
        if ((0, helpers_1.isSet)(object.abstainCount))
            obj.abstainCount = String(object.abstainCount);
        if ((0, helpers_1.isSet)(object.noCount))
            obj.noCount = String(object.noCount);
        if ((0, helpers_1.isSet)(object.noWithVetoCount))
            obj.noWithVetoCount = String(object.noWithVetoCount);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.yesCount !== undefined && (obj.yesCount = message.yesCount);
        message.abstainCount !== undefined && (obj.abstainCount = message.abstainCount);
        message.noCount !== undefined && (obj.noCount = message.noCount);
        message.noWithVetoCount !== undefined && (obj.noWithVetoCount = message.noWithVetoCount);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseTallyResult();
        message.yesCount = object.yesCount ?? "";
        message.abstainCount = object.abstainCount ?? "";
        message.noCount = object.noCount ?? "";
        message.noWithVetoCount = object.noWithVetoCount ?? "";
        return message;
    },
};
function createBaseVote() {
    return {
        proposalId: BigInt(0),
        voter: "",
        option: 0,
        metadata: "",
        submitTime: timestamp_1.Timestamp.fromPartial({}),
    };
}
exports.Vote = {
    typeUrl: "/cosmos.group.v1.Vote",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.proposalId !== BigInt(0)) {
            writer.uint32(8).uint64(message.proposalId);
        }
        if (message.voter !== "") {
            writer.uint32(18).string(message.voter);
        }
        if (message.option !== 0) {
            writer.uint32(24).int32(message.option);
        }
        if (message.metadata !== "") {
            writer.uint32(34).string(message.metadata);
        }
        if (message.submitTime !== undefined) {
            timestamp_1.Timestamp.encode(message.submitTime, writer.uint32(42).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseVote();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.proposalId = reader.uint64();
                    break;
                case 2:
                    message.voter = reader.string();
                    break;
                case 3:
                    message.option = reader.int32();
                    break;
                case 4:
                    message.metadata = reader.string();
                    break;
                case 5:
                    message.submitTime = timestamp_1.Timestamp.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseVote();
        if ((0, helpers_1.isSet)(object.proposalId))
            obj.proposalId = BigInt(object.proposalId.toString());
        if ((0, helpers_1.isSet)(object.voter))
            obj.voter = String(object.voter);
        if ((0, helpers_1.isSet)(object.option))
            obj.option = voteOptionFromJSON(object.option);
        if ((0, helpers_1.isSet)(object.metadata))
            obj.metadata = String(object.metadata);
        if ((0, helpers_1.isSet)(object.submitTime))
            obj.submitTime = (0, helpers_1.fromJsonTimestamp)(object.submitTime);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.proposalId !== undefined && (obj.proposalId = (message.proposalId || BigInt(0)).toString());
        message.voter !== undefined && (obj.voter = message.voter);
        message.option !== undefined && (obj.option = voteOptionToJSON(message.option));
        message.metadata !== undefined && (obj.metadata = message.metadata);
        message.submitTime !== undefined && (obj.submitTime = (0, helpers_1.fromTimestamp)(message.submitTime).toISOString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseVote();
        if (object.proposalId !== undefined && object.proposalId !== null) {
            message.proposalId = BigInt(object.proposalId.toString());
        }
        message.voter = object.voter ?? "";
        message.option = object.option ?? 0;
        message.metadata = object.metadata ?? "";
        if (object.submitTime !== undefined && object.submitTime !== null) {
            message.submitTime = timestamp_1.Timestamp.fromPartial(object.submitTime);
        }
        return message;
    },
};
//# sourceMappingURL=types.js.map