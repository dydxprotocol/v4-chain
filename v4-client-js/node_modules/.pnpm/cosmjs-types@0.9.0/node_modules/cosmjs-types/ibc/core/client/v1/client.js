"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Params = exports.Height = exports.UpgradeProposal = exports.ClientUpdateProposal = exports.ClientConsensusStates = exports.ConsensusStateWithHeight = exports.IdentifiedClientState = exports.protobufPackage = void 0;
/* eslint-disable */
const any_1 = require("../../../../google/protobuf/any");
const upgrade_1 = require("../../../../cosmos/upgrade/v1beta1/upgrade");
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "ibc.core.client.v1";
function createBaseIdentifiedClientState() {
    return {
        clientId: "",
        clientState: undefined,
    };
}
exports.IdentifiedClientState = {
    typeUrl: "/ibc.core.client.v1.IdentifiedClientState",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.clientId !== "") {
            writer.uint32(10).string(message.clientId);
        }
        if (message.clientState !== undefined) {
            any_1.Any.encode(message.clientState, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseIdentifiedClientState();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.clientId = reader.string();
                    break;
                case 2:
                    message.clientState = any_1.Any.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseIdentifiedClientState();
        if ((0, helpers_1.isSet)(object.clientId))
            obj.clientId = String(object.clientId);
        if ((0, helpers_1.isSet)(object.clientState))
            obj.clientState = any_1.Any.fromJSON(object.clientState);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.clientId !== undefined && (obj.clientId = message.clientId);
        message.clientState !== undefined &&
            (obj.clientState = message.clientState ? any_1.Any.toJSON(message.clientState) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseIdentifiedClientState();
        message.clientId = object.clientId ?? "";
        if (object.clientState !== undefined && object.clientState !== null) {
            message.clientState = any_1.Any.fromPartial(object.clientState);
        }
        return message;
    },
};
function createBaseConsensusStateWithHeight() {
    return {
        height: exports.Height.fromPartial({}),
        consensusState: undefined,
    };
}
exports.ConsensusStateWithHeight = {
    typeUrl: "/ibc.core.client.v1.ConsensusStateWithHeight",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.height !== undefined) {
            exports.Height.encode(message.height, writer.uint32(10).fork()).ldelim();
        }
        if (message.consensusState !== undefined) {
            any_1.Any.encode(message.consensusState, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseConsensusStateWithHeight();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.height = exports.Height.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.consensusState = any_1.Any.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseConsensusStateWithHeight();
        if ((0, helpers_1.isSet)(object.height))
            obj.height = exports.Height.fromJSON(object.height);
        if ((0, helpers_1.isSet)(object.consensusState))
            obj.consensusState = any_1.Any.fromJSON(object.consensusState);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.height !== undefined && (obj.height = message.height ? exports.Height.toJSON(message.height) : undefined);
        message.consensusState !== undefined &&
            (obj.consensusState = message.consensusState ? any_1.Any.toJSON(message.consensusState) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseConsensusStateWithHeight();
        if (object.height !== undefined && object.height !== null) {
            message.height = exports.Height.fromPartial(object.height);
        }
        if (object.consensusState !== undefined && object.consensusState !== null) {
            message.consensusState = any_1.Any.fromPartial(object.consensusState);
        }
        return message;
    },
};
function createBaseClientConsensusStates() {
    return {
        clientId: "",
        consensusStates: [],
    };
}
exports.ClientConsensusStates = {
    typeUrl: "/ibc.core.client.v1.ClientConsensusStates",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.clientId !== "") {
            writer.uint32(10).string(message.clientId);
        }
        for (const v of message.consensusStates) {
            exports.ConsensusStateWithHeight.encode(v, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseClientConsensusStates();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.clientId = reader.string();
                    break;
                case 2:
                    message.consensusStates.push(exports.ConsensusStateWithHeight.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseClientConsensusStates();
        if ((0, helpers_1.isSet)(object.clientId))
            obj.clientId = String(object.clientId);
        if (Array.isArray(object?.consensusStates))
            obj.consensusStates = object.consensusStates.map((e) => exports.ConsensusStateWithHeight.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.clientId !== undefined && (obj.clientId = message.clientId);
        if (message.consensusStates) {
            obj.consensusStates = message.consensusStates.map((e) => e ? exports.ConsensusStateWithHeight.toJSON(e) : undefined);
        }
        else {
            obj.consensusStates = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseClientConsensusStates();
        message.clientId = object.clientId ?? "";
        message.consensusStates =
            object.consensusStates?.map((e) => exports.ConsensusStateWithHeight.fromPartial(e)) || [];
        return message;
    },
};
function createBaseClientUpdateProposal() {
    return {
        title: "",
        description: "",
        subjectClientId: "",
        substituteClientId: "",
    };
}
exports.ClientUpdateProposal = {
    typeUrl: "/ibc.core.client.v1.ClientUpdateProposal",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.title !== "") {
            writer.uint32(10).string(message.title);
        }
        if (message.description !== "") {
            writer.uint32(18).string(message.description);
        }
        if (message.subjectClientId !== "") {
            writer.uint32(26).string(message.subjectClientId);
        }
        if (message.substituteClientId !== "") {
            writer.uint32(34).string(message.substituteClientId);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseClientUpdateProposal();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.title = reader.string();
                    break;
                case 2:
                    message.description = reader.string();
                    break;
                case 3:
                    message.subjectClientId = reader.string();
                    break;
                case 4:
                    message.substituteClientId = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseClientUpdateProposal();
        if ((0, helpers_1.isSet)(object.title))
            obj.title = String(object.title);
        if ((0, helpers_1.isSet)(object.description))
            obj.description = String(object.description);
        if ((0, helpers_1.isSet)(object.subjectClientId))
            obj.subjectClientId = String(object.subjectClientId);
        if ((0, helpers_1.isSet)(object.substituteClientId))
            obj.substituteClientId = String(object.substituteClientId);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.title !== undefined && (obj.title = message.title);
        message.description !== undefined && (obj.description = message.description);
        message.subjectClientId !== undefined && (obj.subjectClientId = message.subjectClientId);
        message.substituteClientId !== undefined && (obj.substituteClientId = message.substituteClientId);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseClientUpdateProposal();
        message.title = object.title ?? "";
        message.description = object.description ?? "";
        message.subjectClientId = object.subjectClientId ?? "";
        message.substituteClientId = object.substituteClientId ?? "";
        return message;
    },
};
function createBaseUpgradeProposal() {
    return {
        title: "",
        description: "",
        plan: upgrade_1.Plan.fromPartial({}),
        upgradedClientState: undefined,
    };
}
exports.UpgradeProposal = {
    typeUrl: "/ibc.core.client.v1.UpgradeProposal",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.title !== "") {
            writer.uint32(10).string(message.title);
        }
        if (message.description !== "") {
            writer.uint32(18).string(message.description);
        }
        if (message.plan !== undefined) {
            upgrade_1.Plan.encode(message.plan, writer.uint32(26).fork()).ldelim();
        }
        if (message.upgradedClientState !== undefined) {
            any_1.Any.encode(message.upgradedClientState, writer.uint32(34).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseUpgradeProposal();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.title = reader.string();
                    break;
                case 2:
                    message.description = reader.string();
                    break;
                case 3:
                    message.plan = upgrade_1.Plan.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.upgradedClientState = any_1.Any.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseUpgradeProposal();
        if ((0, helpers_1.isSet)(object.title))
            obj.title = String(object.title);
        if ((0, helpers_1.isSet)(object.description))
            obj.description = String(object.description);
        if ((0, helpers_1.isSet)(object.plan))
            obj.plan = upgrade_1.Plan.fromJSON(object.plan);
        if ((0, helpers_1.isSet)(object.upgradedClientState))
            obj.upgradedClientState = any_1.Any.fromJSON(object.upgradedClientState);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.title !== undefined && (obj.title = message.title);
        message.description !== undefined && (obj.description = message.description);
        message.plan !== undefined && (obj.plan = message.plan ? upgrade_1.Plan.toJSON(message.plan) : undefined);
        message.upgradedClientState !== undefined &&
            (obj.upgradedClientState = message.upgradedClientState
                ? any_1.Any.toJSON(message.upgradedClientState)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseUpgradeProposal();
        message.title = object.title ?? "";
        message.description = object.description ?? "";
        if (object.plan !== undefined && object.plan !== null) {
            message.plan = upgrade_1.Plan.fromPartial(object.plan);
        }
        if (object.upgradedClientState !== undefined && object.upgradedClientState !== null) {
            message.upgradedClientState = any_1.Any.fromPartial(object.upgradedClientState);
        }
        return message;
    },
};
function createBaseHeight() {
    return {
        revisionNumber: BigInt(0),
        revisionHeight: BigInt(0),
    };
}
exports.Height = {
    typeUrl: "/ibc.core.client.v1.Height",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.revisionNumber !== BigInt(0)) {
            writer.uint32(8).uint64(message.revisionNumber);
        }
        if (message.revisionHeight !== BigInt(0)) {
            writer.uint32(16).uint64(message.revisionHeight);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseHeight();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.revisionNumber = reader.uint64();
                    break;
                case 2:
                    message.revisionHeight = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseHeight();
        if ((0, helpers_1.isSet)(object.revisionNumber))
            obj.revisionNumber = BigInt(object.revisionNumber.toString());
        if ((0, helpers_1.isSet)(object.revisionHeight))
            obj.revisionHeight = BigInt(object.revisionHeight.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.revisionNumber !== undefined &&
            (obj.revisionNumber = (message.revisionNumber || BigInt(0)).toString());
        message.revisionHeight !== undefined &&
            (obj.revisionHeight = (message.revisionHeight || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseHeight();
        if (object.revisionNumber !== undefined && object.revisionNumber !== null) {
            message.revisionNumber = BigInt(object.revisionNumber.toString());
        }
        if (object.revisionHeight !== undefined && object.revisionHeight !== null) {
            message.revisionHeight = BigInt(object.revisionHeight.toString());
        }
        return message;
    },
};
function createBaseParams() {
    return {
        allowedClients: [],
    };
}
exports.Params = {
    typeUrl: "/ibc.core.client.v1.Params",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.allowedClients) {
            writer.uint32(10).string(v);
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
                    message.allowedClients.push(reader.string());
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
        if (Array.isArray(object?.allowedClients))
            obj.allowedClients = object.allowedClients.map((e) => String(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.allowedClients) {
            obj.allowedClients = message.allowedClients.map((e) => e);
        }
        else {
            obj.allowedClients = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseParams();
        message.allowedClients = object.allowedClients?.map((e) => e) || [];
        return message;
    },
};
//# sourceMappingURL=client.js.map