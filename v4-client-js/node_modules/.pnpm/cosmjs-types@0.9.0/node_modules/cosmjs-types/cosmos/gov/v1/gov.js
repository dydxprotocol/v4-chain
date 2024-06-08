"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Params = exports.TallyParams = exports.VotingParams = exports.DepositParams = exports.Vote = exports.TallyResult = exports.Proposal = exports.Deposit = exports.WeightedVoteOption = exports.proposalStatusToJSON = exports.proposalStatusFromJSON = exports.ProposalStatus = exports.voteOptionToJSON = exports.voteOptionFromJSON = exports.VoteOption = exports.protobufPackage = void 0;
/* eslint-disable */
const coin_1 = require("../../base/v1beta1/coin");
const any_1 = require("../../../google/protobuf/any");
const timestamp_1 = require("../../../google/protobuf/timestamp");
const duration_1 = require("../../../google/protobuf/duration");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.gov.v1";
/** VoteOption enumerates the valid vote options for a given governance proposal. */
var VoteOption;
(function (VoteOption) {
    /** VOTE_OPTION_UNSPECIFIED - VOTE_OPTION_UNSPECIFIED defines a no-op vote option. */
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
/** ProposalStatus enumerates the valid statuses of a proposal. */
var ProposalStatus;
(function (ProposalStatus) {
    /** PROPOSAL_STATUS_UNSPECIFIED - PROPOSAL_STATUS_UNSPECIFIED defines the default proposal status. */
    ProposalStatus[ProposalStatus["PROPOSAL_STATUS_UNSPECIFIED"] = 0] = "PROPOSAL_STATUS_UNSPECIFIED";
    /**
     * PROPOSAL_STATUS_DEPOSIT_PERIOD - PROPOSAL_STATUS_DEPOSIT_PERIOD defines a proposal status during the deposit
     * period.
     */
    ProposalStatus[ProposalStatus["PROPOSAL_STATUS_DEPOSIT_PERIOD"] = 1] = "PROPOSAL_STATUS_DEPOSIT_PERIOD";
    /**
     * PROPOSAL_STATUS_VOTING_PERIOD - PROPOSAL_STATUS_VOTING_PERIOD defines a proposal status during the voting
     * period.
     */
    ProposalStatus[ProposalStatus["PROPOSAL_STATUS_VOTING_PERIOD"] = 2] = "PROPOSAL_STATUS_VOTING_PERIOD";
    /**
     * PROPOSAL_STATUS_PASSED - PROPOSAL_STATUS_PASSED defines a proposal status of a proposal that has
     * passed.
     */
    ProposalStatus[ProposalStatus["PROPOSAL_STATUS_PASSED"] = 3] = "PROPOSAL_STATUS_PASSED";
    /**
     * PROPOSAL_STATUS_REJECTED - PROPOSAL_STATUS_REJECTED defines a proposal status of a proposal that has
     * been rejected.
     */
    ProposalStatus[ProposalStatus["PROPOSAL_STATUS_REJECTED"] = 4] = "PROPOSAL_STATUS_REJECTED";
    /**
     * PROPOSAL_STATUS_FAILED - PROPOSAL_STATUS_FAILED defines a proposal status of a proposal that has
     * failed.
     */
    ProposalStatus[ProposalStatus["PROPOSAL_STATUS_FAILED"] = 5] = "PROPOSAL_STATUS_FAILED";
    ProposalStatus[ProposalStatus["UNRECOGNIZED"] = -1] = "UNRECOGNIZED";
})(ProposalStatus || (exports.ProposalStatus = ProposalStatus = {}));
function proposalStatusFromJSON(object) {
    switch (object) {
        case 0:
        case "PROPOSAL_STATUS_UNSPECIFIED":
            return ProposalStatus.PROPOSAL_STATUS_UNSPECIFIED;
        case 1:
        case "PROPOSAL_STATUS_DEPOSIT_PERIOD":
            return ProposalStatus.PROPOSAL_STATUS_DEPOSIT_PERIOD;
        case 2:
        case "PROPOSAL_STATUS_VOTING_PERIOD":
            return ProposalStatus.PROPOSAL_STATUS_VOTING_PERIOD;
        case 3:
        case "PROPOSAL_STATUS_PASSED":
            return ProposalStatus.PROPOSAL_STATUS_PASSED;
        case 4:
        case "PROPOSAL_STATUS_REJECTED":
            return ProposalStatus.PROPOSAL_STATUS_REJECTED;
        case 5:
        case "PROPOSAL_STATUS_FAILED":
            return ProposalStatus.PROPOSAL_STATUS_FAILED;
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
        case ProposalStatus.PROPOSAL_STATUS_DEPOSIT_PERIOD:
            return "PROPOSAL_STATUS_DEPOSIT_PERIOD";
        case ProposalStatus.PROPOSAL_STATUS_VOTING_PERIOD:
            return "PROPOSAL_STATUS_VOTING_PERIOD";
        case ProposalStatus.PROPOSAL_STATUS_PASSED:
            return "PROPOSAL_STATUS_PASSED";
        case ProposalStatus.PROPOSAL_STATUS_REJECTED:
            return "PROPOSAL_STATUS_REJECTED";
        case ProposalStatus.PROPOSAL_STATUS_FAILED:
            return "PROPOSAL_STATUS_FAILED";
        case ProposalStatus.UNRECOGNIZED:
        default:
            return "UNRECOGNIZED";
    }
}
exports.proposalStatusToJSON = proposalStatusToJSON;
function createBaseWeightedVoteOption() {
    return {
        option: 0,
        weight: "",
    };
}
exports.WeightedVoteOption = {
    typeUrl: "/cosmos.gov.v1.WeightedVoteOption",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.option !== 0) {
            writer.uint32(8).int32(message.option);
        }
        if (message.weight !== "") {
            writer.uint32(18).string(message.weight);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseWeightedVoteOption();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.option = reader.int32();
                    break;
                case 2:
                    message.weight = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseWeightedVoteOption();
        if ((0, helpers_1.isSet)(object.option))
            obj.option = voteOptionFromJSON(object.option);
        if ((0, helpers_1.isSet)(object.weight))
            obj.weight = String(object.weight);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.option !== undefined && (obj.option = voteOptionToJSON(message.option));
        message.weight !== undefined && (obj.weight = message.weight);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseWeightedVoteOption();
        message.option = object.option ?? 0;
        message.weight = object.weight ?? "";
        return message;
    },
};
function createBaseDeposit() {
    return {
        proposalId: BigInt(0),
        depositor: "",
        amount: [],
    };
}
exports.Deposit = {
    typeUrl: "/cosmos.gov.v1.Deposit",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.proposalId !== BigInt(0)) {
            writer.uint32(8).uint64(message.proposalId);
        }
        if (message.depositor !== "") {
            writer.uint32(18).string(message.depositor);
        }
        for (const v of message.amount) {
            coin_1.Coin.encode(v, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseDeposit();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.proposalId = reader.uint64();
                    break;
                case 2:
                    message.depositor = reader.string();
                    break;
                case 3:
                    message.amount.push(coin_1.Coin.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseDeposit();
        if ((0, helpers_1.isSet)(object.proposalId))
            obj.proposalId = BigInt(object.proposalId.toString());
        if ((0, helpers_1.isSet)(object.depositor))
            obj.depositor = String(object.depositor);
        if (Array.isArray(object?.amount))
            obj.amount = object.amount.map((e) => coin_1.Coin.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.proposalId !== undefined && (obj.proposalId = (message.proposalId || BigInt(0)).toString());
        message.depositor !== undefined && (obj.depositor = message.depositor);
        if (message.amount) {
            obj.amount = message.amount.map((e) => (e ? coin_1.Coin.toJSON(e) : undefined));
        }
        else {
            obj.amount = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseDeposit();
        if (object.proposalId !== undefined && object.proposalId !== null) {
            message.proposalId = BigInt(object.proposalId.toString());
        }
        message.depositor = object.depositor ?? "";
        message.amount = object.amount?.map((e) => coin_1.Coin.fromPartial(e)) || [];
        return message;
    },
};
function createBaseProposal() {
    return {
        id: BigInt(0),
        messages: [],
        status: 0,
        finalTallyResult: undefined,
        submitTime: undefined,
        depositEndTime: undefined,
        totalDeposit: [],
        votingStartTime: undefined,
        votingEndTime: undefined,
        metadata: "",
        title: "",
        summary: "",
        proposer: "",
    };
}
exports.Proposal = {
    typeUrl: "/cosmos.gov.v1.Proposal",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.id !== BigInt(0)) {
            writer.uint32(8).uint64(message.id);
        }
        for (const v of message.messages) {
            any_1.Any.encode(v, writer.uint32(18).fork()).ldelim();
        }
        if (message.status !== 0) {
            writer.uint32(24).int32(message.status);
        }
        if (message.finalTallyResult !== undefined) {
            exports.TallyResult.encode(message.finalTallyResult, writer.uint32(34).fork()).ldelim();
        }
        if (message.submitTime !== undefined) {
            timestamp_1.Timestamp.encode(message.submitTime, writer.uint32(42).fork()).ldelim();
        }
        if (message.depositEndTime !== undefined) {
            timestamp_1.Timestamp.encode(message.depositEndTime, writer.uint32(50).fork()).ldelim();
        }
        for (const v of message.totalDeposit) {
            coin_1.Coin.encode(v, writer.uint32(58).fork()).ldelim();
        }
        if (message.votingStartTime !== undefined) {
            timestamp_1.Timestamp.encode(message.votingStartTime, writer.uint32(66).fork()).ldelim();
        }
        if (message.votingEndTime !== undefined) {
            timestamp_1.Timestamp.encode(message.votingEndTime, writer.uint32(74).fork()).ldelim();
        }
        if (message.metadata !== "") {
            writer.uint32(82).string(message.metadata);
        }
        if (message.title !== "") {
            writer.uint32(90).string(message.title);
        }
        if (message.summary !== "") {
            writer.uint32(98).string(message.summary);
        }
        if (message.proposer !== "") {
            writer.uint32(106).string(message.proposer);
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
                    message.messages.push(any_1.Any.decode(reader, reader.uint32()));
                    break;
                case 3:
                    message.status = reader.int32();
                    break;
                case 4:
                    message.finalTallyResult = exports.TallyResult.decode(reader, reader.uint32());
                    break;
                case 5:
                    message.submitTime = timestamp_1.Timestamp.decode(reader, reader.uint32());
                    break;
                case 6:
                    message.depositEndTime = timestamp_1.Timestamp.decode(reader, reader.uint32());
                    break;
                case 7:
                    message.totalDeposit.push(coin_1.Coin.decode(reader, reader.uint32()));
                    break;
                case 8:
                    message.votingStartTime = timestamp_1.Timestamp.decode(reader, reader.uint32());
                    break;
                case 9:
                    message.votingEndTime = timestamp_1.Timestamp.decode(reader, reader.uint32());
                    break;
                case 10:
                    message.metadata = reader.string();
                    break;
                case 11:
                    message.title = reader.string();
                    break;
                case 12:
                    message.summary = reader.string();
                    break;
                case 13:
                    message.proposer = reader.string();
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
        if (Array.isArray(object?.messages))
            obj.messages = object.messages.map((e) => any_1.Any.fromJSON(e));
        if ((0, helpers_1.isSet)(object.status))
            obj.status = proposalStatusFromJSON(object.status);
        if ((0, helpers_1.isSet)(object.finalTallyResult))
            obj.finalTallyResult = exports.TallyResult.fromJSON(object.finalTallyResult);
        if ((0, helpers_1.isSet)(object.submitTime))
            obj.submitTime = (0, helpers_1.fromJsonTimestamp)(object.submitTime);
        if ((0, helpers_1.isSet)(object.depositEndTime))
            obj.depositEndTime = (0, helpers_1.fromJsonTimestamp)(object.depositEndTime);
        if (Array.isArray(object?.totalDeposit))
            obj.totalDeposit = object.totalDeposit.map((e) => coin_1.Coin.fromJSON(e));
        if ((0, helpers_1.isSet)(object.votingStartTime))
            obj.votingStartTime = (0, helpers_1.fromJsonTimestamp)(object.votingStartTime);
        if ((0, helpers_1.isSet)(object.votingEndTime))
            obj.votingEndTime = (0, helpers_1.fromJsonTimestamp)(object.votingEndTime);
        if ((0, helpers_1.isSet)(object.metadata))
            obj.metadata = String(object.metadata);
        if ((0, helpers_1.isSet)(object.title))
            obj.title = String(object.title);
        if ((0, helpers_1.isSet)(object.summary))
            obj.summary = String(object.summary);
        if ((0, helpers_1.isSet)(object.proposer))
            obj.proposer = String(object.proposer);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.id !== undefined && (obj.id = (message.id || BigInt(0)).toString());
        if (message.messages) {
            obj.messages = message.messages.map((e) => (e ? any_1.Any.toJSON(e) : undefined));
        }
        else {
            obj.messages = [];
        }
        message.status !== undefined && (obj.status = proposalStatusToJSON(message.status));
        message.finalTallyResult !== undefined &&
            (obj.finalTallyResult = message.finalTallyResult
                ? exports.TallyResult.toJSON(message.finalTallyResult)
                : undefined);
        message.submitTime !== undefined && (obj.submitTime = (0, helpers_1.fromTimestamp)(message.submitTime).toISOString());
        message.depositEndTime !== undefined &&
            (obj.depositEndTime = (0, helpers_1.fromTimestamp)(message.depositEndTime).toISOString());
        if (message.totalDeposit) {
            obj.totalDeposit = message.totalDeposit.map((e) => (e ? coin_1.Coin.toJSON(e) : undefined));
        }
        else {
            obj.totalDeposit = [];
        }
        message.votingStartTime !== undefined &&
            (obj.votingStartTime = (0, helpers_1.fromTimestamp)(message.votingStartTime).toISOString());
        message.votingEndTime !== undefined &&
            (obj.votingEndTime = (0, helpers_1.fromTimestamp)(message.votingEndTime).toISOString());
        message.metadata !== undefined && (obj.metadata = message.metadata);
        message.title !== undefined && (obj.title = message.title);
        message.summary !== undefined && (obj.summary = message.summary);
        message.proposer !== undefined && (obj.proposer = message.proposer);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseProposal();
        if (object.id !== undefined && object.id !== null) {
            message.id = BigInt(object.id.toString());
        }
        message.messages = object.messages?.map((e) => any_1.Any.fromPartial(e)) || [];
        message.status = object.status ?? 0;
        if (object.finalTallyResult !== undefined && object.finalTallyResult !== null) {
            message.finalTallyResult = exports.TallyResult.fromPartial(object.finalTallyResult);
        }
        if (object.submitTime !== undefined && object.submitTime !== null) {
            message.submitTime = timestamp_1.Timestamp.fromPartial(object.submitTime);
        }
        if (object.depositEndTime !== undefined && object.depositEndTime !== null) {
            message.depositEndTime = timestamp_1.Timestamp.fromPartial(object.depositEndTime);
        }
        message.totalDeposit = object.totalDeposit?.map((e) => coin_1.Coin.fromPartial(e)) || [];
        if (object.votingStartTime !== undefined && object.votingStartTime !== null) {
            message.votingStartTime = timestamp_1.Timestamp.fromPartial(object.votingStartTime);
        }
        if (object.votingEndTime !== undefined && object.votingEndTime !== null) {
            message.votingEndTime = timestamp_1.Timestamp.fromPartial(object.votingEndTime);
        }
        message.metadata = object.metadata ?? "";
        message.title = object.title ?? "";
        message.summary = object.summary ?? "";
        message.proposer = object.proposer ?? "";
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
    typeUrl: "/cosmos.gov.v1.TallyResult",
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
        options: [],
        metadata: "",
    };
}
exports.Vote = {
    typeUrl: "/cosmos.gov.v1.Vote",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.proposalId !== BigInt(0)) {
            writer.uint32(8).uint64(message.proposalId);
        }
        if (message.voter !== "") {
            writer.uint32(18).string(message.voter);
        }
        for (const v of message.options) {
            exports.WeightedVoteOption.encode(v, writer.uint32(34).fork()).ldelim();
        }
        if (message.metadata !== "") {
            writer.uint32(42).string(message.metadata);
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
                case 4:
                    message.options.push(exports.WeightedVoteOption.decode(reader, reader.uint32()));
                    break;
                case 5:
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
        const obj = createBaseVote();
        if ((0, helpers_1.isSet)(object.proposalId))
            obj.proposalId = BigInt(object.proposalId.toString());
        if ((0, helpers_1.isSet)(object.voter))
            obj.voter = String(object.voter);
        if (Array.isArray(object?.options))
            obj.options = object.options.map((e) => exports.WeightedVoteOption.fromJSON(e));
        if ((0, helpers_1.isSet)(object.metadata))
            obj.metadata = String(object.metadata);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.proposalId !== undefined && (obj.proposalId = (message.proposalId || BigInt(0)).toString());
        message.voter !== undefined && (obj.voter = message.voter);
        if (message.options) {
            obj.options = message.options.map((e) => (e ? exports.WeightedVoteOption.toJSON(e) : undefined));
        }
        else {
            obj.options = [];
        }
        message.metadata !== undefined && (obj.metadata = message.metadata);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseVote();
        if (object.proposalId !== undefined && object.proposalId !== null) {
            message.proposalId = BigInt(object.proposalId.toString());
        }
        message.voter = object.voter ?? "";
        message.options = object.options?.map((e) => exports.WeightedVoteOption.fromPartial(e)) || [];
        message.metadata = object.metadata ?? "";
        return message;
    },
};
function createBaseDepositParams() {
    return {
        minDeposit: [],
        maxDepositPeriod: undefined,
    };
}
exports.DepositParams = {
    typeUrl: "/cosmos.gov.v1.DepositParams",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.minDeposit) {
            coin_1.Coin.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.maxDepositPeriod !== undefined) {
            duration_1.Duration.encode(message.maxDepositPeriod, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseDepositParams();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.minDeposit.push(coin_1.Coin.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.maxDepositPeriod = duration_1.Duration.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseDepositParams();
        if (Array.isArray(object?.minDeposit))
            obj.minDeposit = object.minDeposit.map((e) => coin_1.Coin.fromJSON(e));
        if ((0, helpers_1.isSet)(object.maxDepositPeriod))
            obj.maxDepositPeriod = duration_1.Duration.fromJSON(object.maxDepositPeriod);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.minDeposit) {
            obj.minDeposit = message.minDeposit.map((e) => (e ? coin_1.Coin.toJSON(e) : undefined));
        }
        else {
            obj.minDeposit = [];
        }
        message.maxDepositPeriod !== undefined &&
            (obj.maxDepositPeriod = message.maxDepositPeriod
                ? duration_1.Duration.toJSON(message.maxDepositPeriod)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseDepositParams();
        message.minDeposit = object.minDeposit?.map((e) => coin_1.Coin.fromPartial(e)) || [];
        if (object.maxDepositPeriod !== undefined && object.maxDepositPeriod !== null) {
            message.maxDepositPeriod = duration_1.Duration.fromPartial(object.maxDepositPeriod);
        }
        return message;
    },
};
function createBaseVotingParams() {
    return {
        votingPeriod: undefined,
    };
}
exports.VotingParams = {
    typeUrl: "/cosmos.gov.v1.VotingParams",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.votingPeriod !== undefined) {
            duration_1.Duration.encode(message.votingPeriod, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseVotingParams();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.votingPeriod = duration_1.Duration.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseVotingParams();
        if ((0, helpers_1.isSet)(object.votingPeriod))
            obj.votingPeriod = duration_1.Duration.fromJSON(object.votingPeriod);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.votingPeriod !== undefined &&
            (obj.votingPeriod = message.votingPeriod ? duration_1.Duration.toJSON(message.votingPeriod) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseVotingParams();
        if (object.votingPeriod !== undefined && object.votingPeriod !== null) {
            message.votingPeriod = duration_1.Duration.fromPartial(object.votingPeriod);
        }
        return message;
    },
};
function createBaseTallyParams() {
    return {
        quorum: "",
        threshold: "",
        vetoThreshold: "",
    };
}
exports.TallyParams = {
    typeUrl: "/cosmos.gov.v1.TallyParams",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.quorum !== "") {
            writer.uint32(10).string(message.quorum);
        }
        if (message.threshold !== "") {
            writer.uint32(18).string(message.threshold);
        }
        if (message.vetoThreshold !== "") {
            writer.uint32(26).string(message.vetoThreshold);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseTallyParams();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.quorum = reader.string();
                    break;
                case 2:
                    message.threshold = reader.string();
                    break;
                case 3:
                    message.vetoThreshold = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseTallyParams();
        if ((0, helpers_1.isSet)(object.quorum))
            obj.quorum = String(object.quorum);
        if ((0, helpers_1.isSet)(object.threshold))
            obj.threshold = String(object.threshold);
        if ((0, helpers_1.isSet)(object.vetoThreshold))
            obj.vetoThreshold = String(object.vetoThreshold);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.quorum !== undefined && (obj.quorum = message.quorum);
        message.threshold !== undefined && (obj.threshold = message.threshold);
        message.vetoThreshold !== undefined && (obj.vetoThreshold = message.vetoThreshold);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseTallyParams();
        message.quorum = object.quorum ?? "";
        message.threshold = object.threshold ?? "";
        message.vetoThreshold = object.vetoThreshold ?? "";
        return message;
    },
};
function createBaseParams() {
    return {
        minDeposit: [],
        maxDepositPeriod: undefined,
        votingPeriod: undefined,
        quorum: "",
        threshold: "",
        vetoThreshold: "",
        minInitialDepositRatio: "",
        burnVoteQuorum: false,
        burnProposalDepositPrevote: false,
        burnVoteVeto: false,
    };
}
exports.Params = {
    typeUrl: "/cosmos.gov.v1.Params",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.minDeposit) {
            coin_1.Coin.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.maxDepositPeriod !== undefined) {
            duration_1.Duration.encode(message.maxDepositPeriod, writer.uint32(18).fork()).ldelim();
        }
        if (message.votingPeriod !== undefined) {
            duration_1.Duration.encode(message.votingPeriod, writer.uint32(26).fork()).ldelim();
        }
        if (message.quorum !== "") {
            writer.uint32(34).string(message.quorum);
        }
        if (message.threshold !== "") {
            writer.uint32(42).string(message.threshold);
        }
        if (message.vetoThreshold !== "") {
            writer.uint32(50).string(message.vetoThreshold);
        }
        if (message.minInitialDepositRatio !== "") {
            writer.uint32(58).string(message.minInitialDepositRatio);
        }
        if (message.burnVoteQuorum === true) {
            writer.uint32(104).bool(message.burnVoteQuorum);
        }
        if (message.burnProposalDepositPrevote === true) {
            writer.uint32(112).bool(message.burnProposalDepositPrevote);
        }
        if (message.burnVoteVeto === true) {
            writer.uint32(120).bool(message.burnVoteVeto);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseParams();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.minDeposit.push(coin_1.Coin.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.maxDepositPeriod = duration_1.Duration.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.votingPeriod = duration_1.Duration.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.quorum = reader.string();
                    break;
                case 5:
                    message.threshold = reader.string();
                    break;
                case 6:
                    message.vetoThreshold = reader.string();
                    break;
                case 7:
                    message.minInitialDepositRatio = reader.string();
                    break;
                case 13:
                    message.burnVoteQuorum = reader.bool();
                    break;
                case 14:
                    message.burnProposalDepositPrevote = reader.bool();
                    break;
                case 15:
                    message.burnVoteVeto = reader.bool();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseParams();
        if (Array.isArray(object?.minDeposit))
            obj.minDeposit = object.minDeposit.map((e) => coin_1.Coin.fromJSON(e));
        if ((0, helpers_1.isSet)(object.maxDepositPeriod))
            obj.maxDepositPeriod = duration_1.Duration.fromJSON(object.maxDepositPeriod);
        if ((0, helpers_1.isSet)(object.votingPeriod))
            obj.votingPeriod = duration_1.Duration.fromJSON(object.votingPeriod);
        if ((0, helpers_1.isSet)(object.quorum))
            obj.quorum = String(object.quorum);
        if ((0, helpers_1.isSet)(object.threshold))
            obj.threshold = String(object.threshold);
        if ((0, helpers_1.isSet)(object.vetoThreshold))
            obj.vetoThreshold = String(object.vetoThreshold);
        if ((0, helpers_1.isSet)(object.minInitialDepositRatio))
            obj.minInitialDepositRatio = String(object.minInitialDepositRatio);
        if ((0, helpers_1.isSet)(object.burnVoteQuorum))
            obj.burnVoteQuorum = Boolean(object.burnVoteQuorum);
        if ((0, helpers_1.isSet)(object.burnProposalDepositPrevote))
            obj.burnProposalDepositPrevote = Boolean(object.burnProposalDepositPrevote);
        if ((0, helpers_1.isSet)(object.burnVoteVeto))
            obj.burnVoteVeto = Boolean(object.burnVoteVeto);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.minDeposit) {
            obj.minDeposit = message.minDeposit.map((e) => (e ? coin_1.Coin.toJSON(e) : undefined));
        }
        else {
            obj.minDeposit = [];
        }
        message.maxDepositPeriod !== undefined &&
            (obj.maxDepositPeriod = message.maxDepositPeriod
                ? duration_1.Duration.toJSON(message.maxDepositPeriod)
                : undefined);
        message.votingPeriod !== undefined &&
            (obj.votingPeriod = message.votingPeriod ? duration_1.Duration.toJSON(message.votingPeriod) : undefined);
        message.quorum !== undefined && (obj.quorum = message.quorum);
        message.threshold !== undefined && (obj.threshold = message.threshold);
        message.vetoThreshold !== undefined && (obj.vetoThreshold = message.vetoThreshold);
        message.minInitialDepositRatio !== undefined &&
            (obj.minInitialDepositRatio = message.minInitialDepositRatio);
        message.burnVoteQuorum !== undefined && (obj.burnVoteQuorum = message.burnVoteQuorum);
        message.burnProposalDepositPrevote !== undefined &&
            (obj.burnProposalDepositPrevote = message.burnProposalDepositPrevote);
        message.burnVoteVeto !== undefined && (obj.burnVoteVeto = message.burnVoteVeto);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseParams();
        message.minDeposit = object.minDeposit?.map((e) => coin_1.Coin.fromPartial(e)) || [];
        if (object.maxDepositPeriod !== undefined && object.maxDepositPeriod !== null) {
            message.maxDepositPeriod = duration_1.Duration.fromPartial(object.maxDepositPeriod);
        }
        if (object.votingPeriod !== undefined && object.votingPeriod !== null) {
            message.votingPeriod = duration_1.Duration.fromPartial(object.votingPeriod);
        }
        message.quorum = object.quorum ?? "";
        message.threshold = object.threshold ?? "";
        message.vetoThreshold = object.vetoThreshold ?? "";
        message.minInitialDepositRatio = object.minInitialDepositRatio ?? "";
        message.burnVoteQuorum = object.burnVoteQuorum ?? false;
        message.burnProposalDepositPrevote = object.burnProposalDepositPrevote ?? false;
        message.burnVoteVeto = object.burnVoteVeto ?? false;
        return message;
    },
};
//# sourceMappingURL=gov.js.map