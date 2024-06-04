"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GenesisState = exports.protobufPackage = void 0;
/* eslint-disable */
const gov_1 = require("./gov");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.gov.v1";
function createBaseGenesisState() {
    return {
        startingProposalId: BigInt(0),
        deposits: [],
        votes: [],
        proposals: [],
        depositParams: undefined,
        votingParams: undefined,
        tallyParams: undefined,
        params: undefined,
    };
}
exports.GenesisState = {
    typeUrl: "/cosmos.gov.v1.GenesisState",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.startingProposalId !== BigInt(0)) {
            writer.uint32(8).uint64(message.startingProposalId);
        }
        for (const v of message.deposits) {
            gov_1.Deposit.encode(v, writer.uint32(18).fork()).ldelim();
        }
        for (const v of message.votes) {
            gov_1.Vote.encode(v, writer.uint32(26).fork()).ldelim();
        }
        for (const v of message.proposals) {
            gov_1.Proposal.encode(v, writer.uint32(34).fork()).ldelim();
        }
        if (message.depositParams !== undefined) {
            gov_1.DepositParams.encode(message.depositParams, writer.uint32(42).fork()).ldelim();
        }
        if (message.votingParams !== undefined) {
            gov_1.VotingParams.encode(message.votingParams, writer.uint32(50).fork()).ldelim();
        }
        if (message.tallyParams !== undefined) {
            gov_1.TallyParams.encode(message.tallyParams, writer.uint32(58).fork()).ldelim();
        }
        if (message.params !== undefined) {
            gov_1.Params.encode(message.params, writer.uint32(66).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGenesisState();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.startingProposalId = reader.uint64();
                    break;
                case 2:
                    message.deposits.push(gov_1.Deposit.decode(reader, reader.uint32()));
                    break;
                case 3:
                    message.votes.push(gov_1.Vote.decode(reader, reader.uint32()));
                    break;
                case 4:
                    message.proposals.push(gov_1.Proposal.decode(reader, reader.uint32()));
                    break;
                case 5:
                    message.depositParams = gov_1.DepositParams.decode(reader, reader.uint32());
                    break;
                case 6:
                    message.votingParams = gov_1.VotingParams.decode(reader, reader.uint32());
                    break;
                case 7:
                    message.tallyParams = gov_1.TallyParams.decode(reader, reader.uint32());
                    break;
                case 8:
                    message.params = gov_1.Params.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseGenesisState();
        if ((0, helpers_1.isSet)(object.startingProposalId))
            obj.startingProposalId = BigInt(object.startingProposalId.toString());
        if (Array.isArray(object?.deposits))
            obj.deposits = object.deposits.map((e) => gov_1.Deposit.fromJSON(e));
        if (Array.isArray(object?.votes))
            obj.votes = object.votes.map((e) => gov_1.Vote.fromJSON(e));
        if (Array.isArray(object?.proposals))
            obj.proposals = object.proposals.map((e) => gov_1.Proposal.fromJSON(e));
        if ((0, helpers_1.isSet)(object.depositParams))
            obj.depositParams = gov_1.DepositParams.fromJSON(object.depositParams);
        if ((0, helpers_1.isSet)(object.votingParams))
            obj.votingParams = gov_1.VotingParams.fromJSON(object.votingParams);
        if ((0, helpers_1.isSet)(object.tallyParams))
            obj.tallyParams = gov_1.TallyParams.fromJSON(object.tallyParams);
        if ((0, helpers_1.isSet)(object.params))
            obj.params = gov_1.Params.fromJSON(object.params);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.startingProposalId !== undefined &&
            (obj.startingProposalId = (message.startingProposalId || BigInt(0)).toString());
        if (message.deposits) {
            obj.deposits = message.deposits.map((e) => (e ? gov_1.Deposit.toJSON(e) : undefined));
        }
        else {
            obj.deposits = [];
        }
        if (message.votes) {
            obj.votes = message.votes.map((e) => (e ? gov_1.Vote.toJSON(e) : undefined));
        }
        else {
            obj.votes = [];
        }
        if (message.proposals) {
            obj.proposals = message.proposals.map((e) => (e ? gov_1.Proposal.toJSON(e) : undefined));
        }
        else {
            obj.proposals = [];
        }
        message.depositParams !== undefined &&
            (obj.depositParams = message.depositParams ? gov_1.DepositParams.toJSON(message.depositParams) : undefined);
        message.votingParams !== undefined &&
            (obj.votingParams = message.votingParams ? gov_1.VotingParams.toJSON(message.votingParams) : undefined);
        message.tallyParams !== undefined &&
            (obj.tallyParams = message.tallyParams ? gov_1.TallyParams.toJSON(message.tallyParams) : undefined);
        message.params !== undefined && (obj.params = message.params ? gov_1.Params.toJSON(message.params) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGenesisState();
        if (object.startingProposalId !== undefined && object.startingProposalId !== null) {
            message.startingProposalId = BigInt(object.startingProposalId.toString());
        }
        message.deposits = object.deposits?.map((e) => gov_1.Deposit.fromPartial(e)) || [];
        message.votes = object.votes?.map((e) => gov_1.Vote.fromPartial(e)) || [];
        message.proposals = object.proposals?.map((e) => gov_1.Proposal.fromPartial(e)) || [];
        if (object.depositParams !== undefined && object.depositParams !== null) {
            message.depositParams = gov_1.DepositParams.fromPartial(object.depositParams);
        }
        if (object.votingParams !== undefined && object.votingParams !== null) {
            message.votingParams = gov_1.VotingParams.fromPartial(object.votingParams);
        }
        if (object.tallyParams !== undefined && object.tallyParams !== null) {
            message.tallyParams = gov_1.TallyParams.fromPartial(object.tallyParams);
        }
        if (object.params !== undefined && object.params !== null) {
            message.params = gov_1.Params.fromPartial(object.params);
        }
        return message;
    },
};
//# sourceMappingURL=genesis.js.map