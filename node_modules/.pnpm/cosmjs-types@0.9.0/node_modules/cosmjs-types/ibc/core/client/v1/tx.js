"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.MsgClientImpl = exports.MsgSubmitMisbehaviourResponse = exports.MsgSubmitMisbehaviour = exports.MsgUpgradeClientResponse = exports.MsgUpgradeClient = exports.MsgUpdateClientResponse = exports.MsgUpdateClient = exports.MsgCreateClientResponse = exports.MsgCreateClient = exports.protobufPackage = void 0;
/* eslint-disable */
const any_1 = require("../../../../google/protobuf/any");
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "ibc.core.client.v1";
function createBaseMsgCreateClient() {
    return {
        clientState: undefined,
        consensusState: undefined,
        signer: "",
    };
}
exports.MsgCreateClient = {
    typeUrl: "/ibc.core.client.v1.MsgCreateClient",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.clientState !== undefined) {
            any_1.Any.encode(message.clientState, writer.uint32(10).fork()).ldelim();
        }
        if (message.consensusState !== undefined) {
            any_1.Any.encode(message.consensusState, writer.uint32(18).fork()).ldelim();
        }
        if (message.signer !== "") {
            writer.uint32(26).string(message.signer);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgCreateClient();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.clientState = any_1.Any.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.consensusState = any_1.Any.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.signer = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgCreateClient();
        if ((0, helpers_1.isSet)(object.clientState))
            obj.clientState = any_1.Any.fromJSON(object.clientState);
        if ((0, helpers_1.isSet)(object.consensusState))
            obj.consensusState = any_1.Any.fromJSON(object.consensusState);
        if ((0, helpers_1.isSet)(object.signer))
            obj.signer = String(object.signer);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.clientState !== undefined &&
            (obj.clientState = message.clientState ? any_1.Any.toJSON(message.clientState) : undefined);
        message.consensusState !== undefined &&
            (obj.consensusState = message.consensusState ? any_1.Any.toJSON(message.consensusState) : undefined);
        message.signer !== undefined && (obj.signer = message.signer);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgCreateClient();
        if (object.clientState !== undefined && object.clientState !== null) {
            message.clientState = any_1.Any.fromPartial(object.clientState);
        }
        if (object.consensusState !== undefined && object.consensusState !== null) {
            message.consensusState = any_1.Any.fromPartial(object.consensusState);
        }
        message.signer = object.signer ?? "";
        return message;
    },
};
function createBaseMsgCreateClientResponse() {
    return {};
}
exports.MsgCreateClientResponse = {
    typeUrl: "/ibc.core.client.v1.MsgCreateClientResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgCreateClientResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(_) {
        const obj = createBaseMsgCreateClientResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgCreateClientResponse();
        return message;
    },
};
function createBaseMsgUpdateClient() {
    return {
        clientId: "",
        clientMessage: undefined,
        signer: "",
    };
}
exports.MsgUpdateClient = {
    typeUrl: "/ibc.core.client.v1.MsgUpdateClient",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.clientId !== "") {
            writer.uint32(10).string(message.clientId);
        }
        if (message.clientMessage !== undefined) {
            any_1.Any.encode(message.clientMessage, writer.uint32(18).fork()).ldelim();
        }
        if (message.signer !== "") {
            writer.uint32(26).string(message.signer);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgUpdateClient();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.clientId = reader.string();
                    break;
                case 2:
                    message.clientMessage = any_1.Any.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.signer = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgUpdateClient();
        if ((0, helpers_1.isSet)(object.clientId))
            obj.clientId = String(object.clientId);
        if ((0, helpers_1.isSet)(object.clientMessage))
            obj.clientMessage = any_1.Any.fromJSON(object.clientMessage);
        if ((0, helpers_1.isSet)(object.signer))
            obj.signer = String(object.signer);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.clientId !== undefined && (obj.clientId = message.clientId);
        message.clientMessage !== undefined &&
            (obj.clientMessage = message.clientMessage ? any_1.Any.toJSON(message.clientMessage) : undefined);
        message.signer !== undefined && (obj.signer = message.signer);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgUpdateClient();
        message.clientId = object.clientId ?? "";
        if (object.clientMessage !== undefined && object.clientMessage !== null) {
            message.clientMessage = any_1.Any.fromPartial(object.clientMessage);
        }
        message.signer = object.signer ?? "";
        return message;
    },
};
function createBaseMsgUpdateClientResponse() {
    return {};
}
exports.MsgUpdateClientResponse = {
    typeUrl: "/ibc.core.client.v1.MsgUpdateClientResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgUpdateClientResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(_) {
        const obj = createBaseMsgUpdateClientResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgUpdateClientResponse();
        return message;
    },
};
function createBaseMsgUpgradeClient() {
    return {
        clientId: "",
        clientState: undefined,
        consensusState: undefined,
        proofUpgradeClient: new Uint8Array(),
        proofUpgradeConsensusState: new Uint8Array(),
        signer: "",
    };
}
exports.MsgUpgradeClient = {
    typeUrl: "/ibc.core.client.v1.MsgUpgradeClient",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.clientId !== "") {
            writer.uint32(10).string(message.clientId);
        }
        if (message.clientState !== undefined) {
            any_1.Any.encode(message.clientState, writer.uint32(18).fork()).ldelim();
        }
        if (message.consensusState !== undefined) {
            any_1.Any.encode(message.consensusState, writer.uint32(26).fork()).ldelim();
        }
        if (message.proofUpgradeClient.length !== 0) {
            writer.uint32(34).bytes(message.proofUpgradeClient);
        }
        if (message.proofUpgradeConsensusState.length !== 0) {
            writer.uint32(42).bytes(message.proofUpgradeConsensusState);
        }
        if (message.signer !== "") {
            writer.uint32(50).string(message.signer);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgUpgradeClient();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.clientId = reader.string();
                    break;
                case 2:
                    message.clientState = any_1.Any.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.consensusState = any_1.Any.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.proofUpgradeClient = reader.bytes();
                    break;
                case 5:
                    message.proofUpgradeConsensusState = reader.bytes();
                    break;
                case 6:
                    message.signer = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgUpgradeClient();
        if ((0, helpers_1.isSet)(object.clientId))
            obj.clientId = String(object.clientId);
        if ((0, helpers_1.isSet)(object.clientState))
            obj.clientState = any_1.Any.fromJSON(object.clientState);
        if ((0, helpers_1.isSet)(object.consensusState))
            obj.consensusState = any_1.Any.fromJSON(object.consensusState);
        if ((0, helpers_1.isSet)(object.proofUpgradeClient))
            obj.proofUpgradeClient = (0, helpers_1.bytesFromBase64)(object.proofUpgradeClient);
        if ((0, helpers_1.isSet)(object.proofUpgradeConsensusState))
            obj.proofUpgradeConsensusState = (0, helpers_1.bytesFromBase64)(object.proofUpgradeConsensusState);
        if ((0, helpers_1.isSet)(object.signer))
            obj.signer = String(object.signer);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.clientId !== undefined && (obj.clientId = message.clientId);
        message.clientState !== undefined &&
            (obj.clientState = message.clientState ? any_1.Any.toJSON(message.clientState) : undefined);
        message.consensusState !== undefined &&
            (obj.consensusState = message.consensusState ? any_1.Any.toJSON(message.consensusState) : undefined);
        message.proofUpgradeClient !== undefined &&
            (obj.proofUpgradeClient = (0, helpers_1.base64FromBytes)(message.proofUpgradeClient !== undefined ? message.proofUpgradeClient : new Uint8Array()));
        message.proofUpgradeConsensusState !== undefined &&
            (obj.proofUpgradeConsensusState = (0, helpers_1.base64FromBytes)(message.proofUpgradeConsensusState !== undefined
                ? message.proofUpgradeConsensusState
                : new Uint8Array()));
        message.signer !== undefined && (obj.signer = message.signer);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgUpgradeClient();
        message.clientId = object.clientId ?? "";
        if (object.clientState !== undefined && object.clientState !== null) {
            message.clientState = any_1.Any.fromPartial(object.clientState);
        }
        if (object.consensusState !== undefined && object.consensusState !== null) {
            message.consensusState = any_1.Any.fromPartial(object.consensusState);
        }
        message.proofUpgradeClient = object.proofUpgradeClient ?? new Uint8Array();
        message.proofUpgradeConsensusState = object.proofUpgradeConsensusState ?? new Uint8Array();
        message.signer = object.signer ?? "";
        return message;
    },
};
function createBaseMsgUpgradeClientResponse() {
    return {};
}
exports.MsgUpgradeClientResponse = {
    typeUrl: "/ibc.core.client.v1.MsgUpgradeClientResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgUpgradeClientResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(_) {
        const obj = createBaseMsgUpgradeClientResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgUpgradeClientResponse();
        return message;
    },
};
function createBaseMsgSubmitMisbehaviour() {
    return {
        clientId: "",
        misbehaviour: undefined,
        signer: "",
    };
}
exports.MsgSubmitMisbehaviour = {
    typeUrl: "/ibc.core.client.v1.MsgSubmitMisbehaviour",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.clientId !== "") {
            writer.uint32(10).string(message.clientId);
        }
        if (message.misbehaviour !== undefined) {
            any_1.Any.encode(message.misbehaviour, writer.uint32(18).fork()).ldelim();
        }
        if (message.signer !== "") {
            writer.uint32(26).string(message.signer);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgSubmitMisbehaviour();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.clientId = reader.string();
                    break;
                case 2:
                    message.misbehaviour = any_1.Any.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.signer = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgSubmitMisbehaviour();
        if ((0, helpers_1.isSet)(object.clientId))
            obj.clientId = String(object.clientId);
        if ((0, helpers_1.isSet)(object.misbehaviour))
            obj.misbehaviour = any_1.Any.fromJSON(object.misbehaviour);
        if ((0, helpers_1.isSet)(object.signer))
            obj.signer = String(object.signer);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.clientId !== undefined && (obj.clientId = message.clientId);
        message.misbehaviour !== undefined &&
            (obj.misbehaviour = message.misbehaviour ? any_1.Any.toJSON(message.misbehaviour) : undefined);
        message.signer !== undefined && (obj.signer = message.signer);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgSubmitMisbehaviour();
        message.clientId = object.clientId ?? "";
        if (object.misbehaviour !== undefined && object.misbehaviour !== null) {
            message.misbehaviour = any_1.Any.fromPartial(object.misbehaviour);
        }
        message.signer = object.signer ?? "";
        return message;
    },
};
function createBaseMsgSubmitMisbehaviourResponse() {
    return {};
}
exports.MsgSubmitMisbehaviourResponse = {
    typeUrl: "/ibc.core.client.v1.MsgSubmitMisbehaviourResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgSubmitMisbehaviourResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(_) {
        const obj = createBaseMsgSubmitMisbehaviourResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgSubmitMisbehaviourResponse();
        return message;
    },
};
class MsgClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
        this.CreateClient = this.CreateClient.bind(this);
        this.UpdateClient = this.UpdateClient.bind(this);
        this.UpgradeClient = this.UpgradeClient.bind(this);
        this.SubmitMisbehaviour = this.SubmitMisbehaviour.bind(this);
    }
    CreateClient(request) {
        const data = exports.MsgCreateClient.encode(request).finish();
        const promise = this.rpc.request("ibc.core.client.v1.Msg", "CreateClient", data);
        return promise.then((data) => exports.MsgCreateClientResponse.decode(new binary_1.BinaryReader(data)));
    }
    UpdateClient(request) {
        const data = exports.MsgUpdateClient.encode(request).finish();
        const promise = this.rpc.request("ibc.core.client.v1.Msg", "UpdateClient", data);
        return promise.then((data) => exports.MsgUpdateClientResponse.decode(new binary_1.BinaryReader(data)));
    }
    UpgradeClient(request) {
        const data = exports.MsgUpgradeClient.encode(request).finish();
        const promise = this.rpc.request("ibc.core.client.v1.Msg", "UpgradeClient", data);
        return promise.then((data) => exports.MsgUpgradeClientResponse.decode(new binary_1.BinaryReader(data)));
    }
    SubmitMisbehaviour(request) {
        const data = exports.MsgSubmitMisbehaviour.encode(request).finish();
        const promise = this.rpc.request("ibc.core.client.v1.Msg", "SubmitMisbehaviour", data);
        return promise.then((data) => exports.MsgSubmitMisbehaviourResponse.decode(new binary_1.BinaryReader(data)));
    }
}
exports.MsgClientImpl = MsgClientImpl;
//# sourceMappingURL=tx.js.map