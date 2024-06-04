"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GenesisState = exports.protobufPackage = void 0;
/* eslint-disable */
const types_1 = require("./types");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.group.v1";
function createBaseGenesisState() {
    return {
        groupSeq: BigInt(0),
        groups: [],
        groupMembers: [],
        groupPolicySeq: BigInt(0),
        groupPolicies: [],
        proposalSeq: BigInt(0),
        proposals: [],
        votes: [],
    };
}
exports.GenesisState = {
    typeUrl: "/cosmos.group.v1.GenesisState",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.groupSeq !== BigInt(0)) {
            writer.uint32(8).uint64(message.groupSeq);
        }
        for (const v of message.groups) {
            types_1.GroupInfo.encode(v, writer.uint32(18).fork()).ldelim();
        }
        for (const v of message.groupMembers) {
            types_1.GroupMember.encode(v, writer.uint32(26).fork()).ldelim();
        }
        if (message.groupPolicySeq !== BigInt(0)) {
            writer.uint32(32).uint64(message.groupPolicySeq);
        }
        for (const v of message.groupPolicies) {
            types_1.GroupPolicyInfo.encode(v, writer.uint32(42).fork()).ldelim();
        }
        if (message.proposalSeq !== BigInt(0)) {
            writer.uint32(48).uint64(message.proposalSeq);
        }
        for (const v of message.proposals) {
            types_1.Proposal.encode(v, writer.uint32(58).fork()).ldelim();
        }
        for (const v of message.votes) {
            types_1.Vote.encode(v, writer.uint32(66).fork()).ldelim();
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
                    message.groupSeq = reader.uint64();
                    break;
                case 2:
                    message.groups.push(types_1.GroupInfo.decode(reader, reader.uint32()));
                    break;
                case 3:
                    message.groupMembers.push(types_1.GroupMember.decode(reader, reader.uint32()));
                    break;
                case 4:
                    message.groupPolicySeq = reader.uint64();
                    break;
                case 5:
                    message.groupPolicies.push(types_1.GroupPolicyInfo.decode(reader, reader.uint32()));
                    break;
                case 6:
                    message.proposalSeq = reader.uint64();
                    break;
                case 7:
                    message.proposals.push(types_1.Proposal.decode(reader, reader.uint32()));
                    break;
                case 8:
                    message.votes.push(types_1.Vote.decode(reader, reader.uint32()));
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
        if ((0, helpers_1.isSet)(object.groupSeq))
            obj.groupSeq = BigInt(object.groupSeq.toString());
        if (Array.isArray(object?.groups))
            obj.groups = object.groups.map((e) => types_1.GroupInfo.fromJSON(e));
        if (Array.isArray(object?.groupMembers))
            obj.groupMembers = object.groupMembers.map((e) => types_1.GroupMember.fromJSON(e));
        if ((0, helpers_1.isSet)(object.groupPolicySeq))
            obj.groupPolicySeq = BigInt(object.groupPolicySeq.toString());
        if (Array.isArray(object?.groupPolicies))
            obj.groupPolicies = object.groupPolicies.map((e) => types_1.GroupPolicyInfo.fromJSON(e));
        if ((0, helpers_1.isSet)(object.proposalSeq))
            obj.proposalSeq = BigInt(object.proposalSeq.toString());
        if (Array.isArray(object?.proposals))
            obj.proposals = object.proposals.map((e) => types_1.Proposal.fromJSON(e));
        if (Array.isArray(object?.votes))
            obj.votes = object.votes.map((e) => types_1.Vote.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.groupSeq !== undefined && (obj.groupSeq = (message.groupSeq || BigInt(0)).toString());
        if (message.groups) {
            obj.groups = message.groups.map((e) => (e ? types_1.GroupInfo.toJSON(e) : undefined));
        }
        else {
            obj.groups = [];
        }
        if (message.groupMembers) {
            obj.groupMembers = message.groupMembers.map((e) => (e ? types_1.GroupMember.toJSON(e) : undefined));
        }
        else {
            obj.groupMembers = [];
        }
        message.groupPolicySeq !== undefined &&
            (obj.groupPolicySeq = (message.groupPolicySeq || BigInt(0)).toString());
        if (message.groupPolicies) {
            obj.groupPolicies = message.groupPolicies.map((e) => (e ? types_1.GroupPolicyInfo.toJSON(e) : undefined));
        }
        else {
            obj.groupPolicies = [];
        }
        message.proposalSeq !== undefined && (obj.proposalSeq = (message.proposalSeq || BigInt(0)).toString());
        if (message.proposals) {
            obj.proposals = message.proposals.map((e) => (e ? types_1.Proposal.toJSON(e) : undefined));
        }
        else {
            obj.proposals = [];
        }
        if (message.votes) {
            obj.votes = message.votes.map((e) => (e ? types_1.Vote.toJSON(e) : undefined));
        }
        else {
            obj.votes = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGenesisState();
        if (object.groupSeq !== undefined && object.groupSeq !== null) {
            message.groupSeq = BigInt(object.groupSeq.toString());
        }
        message.groups = object.groups?.map((e) => types_1.GroupInfo.fromPartial(e)) || [];
        message.groupMembers = object.groupMembers?.map((e) => types_1.GroupMember.fromPartial(e)) || [];
        if (object.groupPolicySeq !== undefined && object.groupPolicySeq !== null) {
            message.groupPolicySeq = BigInt(object.groupPolicySeq.toString());
        }
        message.groupPolicies = object.groupPolicies?.map((e) => types_1.GroupPolicyInfo.fromPartial(e)) || [];
        if (object.proposalSeq !== undefined && object.proposalSeq !== null) {
            message.proposalSeq = BigInt(object.proposalSeq.toString());
        }
        message.proposals = object.proposals?.map((e) => types_1.Proposal.fromPartial(e)) || [];
        message.votes = object.votes?.map((e) => types_1.Vote.fromPartial(e)) || [];
        return message;
    },
};
//# sourceMappingURL=genesis.js.map