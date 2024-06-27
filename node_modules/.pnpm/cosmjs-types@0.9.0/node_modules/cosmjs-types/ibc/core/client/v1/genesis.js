"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.IdentifiedGenesisMetadata = exports.GenesisMetadata = exports.GenesisState = exports.protobufPackage = void 0;
/* eslint-disable */
const client_1 = require("./client");
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "ibc.core.client.v1";
function createBaseGenesisState() {
    return {
        clients: [],
        clientsConsensus: [],
        clientsMetadata: [],
        params: client_1.Params.fromPartial({}),
        createLocalhost: false,
        nextClientSequence: BigInt(0),
    };
}
exports.GenesisState = {
    typeUrl: "/ibc.core.client.v1.GenesisState",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.clients) {
            client_1.IdentifiedClientState.encode(v, writer.uint32(10).fork()).ldelim();
        }
        for (const v of message.clientsConsensus) {
            client_1.ClientConsensusStates.encode(v, writer.uint32(18).fork()).ldelim();
        }
        for (const v of message.clientsMetadata) {
            exports.IdentifiedGenesisMetadata.encode(v, writer.uint32(26).fork()).ldelim();
        }
        if (message.params !== undefined) {
            client_1.Params.encode(message.params, writer.uint32(34).fork()).ldelim();
        }
        if (message.createLocalhost === true) {
            writer.uint32(40).bool(message.createLocalhost);
        }
        if (message.nextClientSequence !== BigInt(0)) {
            writer.uint32(48).uint64(message.nextClientSequence);
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
                    message.clients.push(client_1.IdentifiedClientState.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.clientsConsensus.push(client_1.ClientConsensusStates.decode(reader, reader.uint32()));
                    break;
                case 3:
                    message.clientsMetadata.push(exports.IdentifiedGenesisMetadata.decode(reader, reader.uint32()));
                    break;
                case 4:
                    message.params = client_1.Params.decode(reader, reader.uint32());
                    break;
                case 5:
                    message.createLocalhost = reader.bool();
                    break;
                case 6:
                    message.nextClientSequence = reader.uint64();
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
        if (Array.isArray(object?.clients))
            obj.clients = object.clients.map((e) => client_1.IdentifiedClientState.fromJSON(e));
        if (Array.isArray(object?.clientsConsensus))
            obj.clientsConsensus = object.clientsConsensus.map((e) => client_1.ClientConsensusStates.fromJSON(e));
        if (Array.isArray(object?.clientsMetadata))
            obj.clientsMetadata = object.clientsMetadata.map((e) => exports.IdentifiedGenesisMetadata.fromJSON(e));
        if ((0, helpers_1.isSet)(object.params))
            obj.params = client_1.Params.fromJSON(object.params);
        if ((0, helpers_1.isSet)(object.createLocalhost))
            obj.createLocalhost = Boolean(object.createLocalhost);
        if ((0, helpers_1.isSet)(object.nextClientSequence))
            obj.nextClientSequence = BigInt(object.nextClientSequence.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.clients) {
            obj.clients = message.clients.map((e) => (e ? client_1.IdentifiedClientState.toJSON(e) : undefined));
        }
        else {
            obj.clients = [];
        }
        if (message.clientsConsensus) {
            obj.clientsConsensus = message.clientsConsensus.map((e) => e ? client_1.ClientConsensusStates.toJSON(e) : undefined);
        }
        else {
            obj.clientsConsensus = [];
        }
        if (message.clientsMetadata) {
            obj.clientsMetadata = message.clientsMetadata.map((e) => e ? exports.IdentifiedGenesisMetadata.toJSON(e) : undefined);
        }
        else {
            obj.clientsMetadata = [];
        }
        message.params !== undefined && (obj.params = message.params ? client_1.Params.toJSON(message.params) : undefined);
        message.createLocalhost !== undefined && (obj.createLocalhost = message.createLocalhost);
        message.nextClientSequence !== undefined &&
            (obj.nextClientSequence = (message.nextClientSequence || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGenesisState();
        message.clients = object.clients?.map((e) => client_1.IdentifiedClientState.fromPartial(e)) || [];
        message.clientsConsensus =
            object.clientsConsensus?.map((e) => client_1.ClientConsensusStates.fromPartial(e)) || [];
        message.clientsMetadata =
            object.clientsMetadata?.map((e) => exports.IdentifiedGenesisMetadata.fromPartial(e)) || [];
        if (object.params !== undefined && object.params !== null) {
            message.params = client_1.Params.fromPartial(object.params);
        }
        message.createLocalhost = object.createLocalhost ?? false;
        if (object.nextClientSequence !== undefined && object.nextClientSequence !== null) {
            message.nextClientSequence = BigInt(object.nextClientSequence.toString());
        }
        return message;
    },
};
function createBaseGenesisMetadata() {
    return {
        key: new Uint8Array(),
        value: new Uint8Array(),
    };
}
exports.GenesisMetadata = {
    typeUrl: "/ibc.core.client.v1.GenesisMetadata",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.key.length !== 0) {
            writer.uint32(10).bytes(message.key);
        }
        if (message.value.length !== 0) {
            writer.uint32(18).bytes(message.value);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGenesisMetadata();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.key = reader.bytes();
                    break;
                case 2:
                    message.value = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseGenesisMetadata();
        if ((0, helpers_1.isSet)(object.key))
            obj.key = (0, helpers_1.bytesFromBase64)(object.key);
        if ((0, helpers_1.isSet)(object.value))
            obj.value = (0, helpers_1.bytesFromBase64)(object.value);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.key !== undefined &&
            (obj.key = (0, helpers_1.base64FromBytes)(message.key !== undefined ? message.key : new Uint8Array()));
        message.value !== undefined &&
            (obj.value = (0, helpers_1.base64FromBytes)(message.value !== undefined ? message.value : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGenesisMetadata();
        message.key = object.key ?? new Uint8Array();
        message.value = object.value ?? new Uint8Array();
        return message;
    },
};
function createBaseIdentifiedGenesisMetadata() {
    return {
        clientId: "",
        clientMetadata: [],
    };
}
exports.IdentifiedGenesisMetadata = {
    typeUrl: "/ibc.core.client.v1.IdentifiedGenesisMetadata",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.clientId !== "") {
            writer.uint32(10).string(message.clientId);
        }
        for (const v of message.clientMetadata) {
            exports.GenesisMetadata.encode(v, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseIdentifiedGenesisMetadata();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.clientId = reader.string();
                    break;
                case 2:
                    message.clientMetadata.push(exports.GenesisMetadata.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseIdentifiedGenesisMetadata();
        if ((0, helpers_1.isSet)(object.clientId))
            obj.clientId = String(object.clientId);
        if (Array.isArray(object?.clientMetadata))
            obj.clientMetadata = object.clientMetadata.map((e) => exports.GenesisMetadata.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.clientId !== undefined && (obj.clientId = message.clientId);
        if (message.clientMetadata) {
            obj.clientMetadata = message.clientMetadata.map((e) => (e ? exports.GenesisMetadata.toJSON(e) : undefined));
        }
        else {
            obj.clientMetadata = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseIdentifiedGenesisMetadata();
        message.clientId = object.clientId ?? "";
        message.clientMetadata = object.clientMetadata?.map((e) => exports.GenesisMetadata.fromPartial(e)) || [];
        return message;
    },
};
//# sourceMappingURL=genesis.js.map