"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Params = exports.Version = exports.ConnectionPaths = exports.ClientPaths = exports.Counterparty = exports.IdentifiedConnection = exports.ConnectionEnd = exports.stateToJSON = exports.stateFromJSON = exports.State = exports.protobufPackage = void 0;
/* eslint-disable */
const commitment_1 = require("../../commitment/v1/commitment");
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "ibc.core.connection.v1";
/**
 * State defines if a connection is in one of the following states:
 * INIT, TRYOPEN, OPEN or UNINITIALIZED.
 */
var State;
(function (State) {
    /** STATE_UNINITIALIZED_UNSPECIFIED - Default State */
    State[State["STATE_UNINITIALIZED_UNSPECIFIED"] = 0] = "STATE_UNINITIALIZED_UNSPECIFIED";
    /** STATE_INIT - A connection end has just started the opening handshake. */
    State[State["STATE_INIT"] = 1] = "STATE_INIT";
    /**
     * STATE_TRYOPEN - A connection end has acknowledged the handshake step on the counterparty
     * chain.
     */
    State[State["STATE_TRYOPEN"] = 2] = "STATE_TRYOPEN";
    /** STATE_OPEN - A connection end has completed the handshake. */
    State[State["STATE_OPEN"] = 3] = "STATE_OPEN";
    State[State["UNRECOGNIZED"] = -1] = "UNRECOGNIZED";
})(State || (exports.State = State = {}));
function stateFromJSON(object) {
    switch (object) {
        case 0:
        case "STATE_UNINITIALIZED_UNSPECIFIED":
            return State.STATE_UNINITIALIZED_UNSPECIFIED;
        case 1:
        case "STATE_INIT":
            return State.STATE_INIT;
        case 2:
        case "STATE_TRYOPEN":
            return State.STATE_TRYOPEN;
        case 3:
        case "STATE_OPEN":
            return State.STATE_OPEN;
        case -1:
        case "UNRECOGNIZED":
        default:
            return State.UNRECOGNIZED;
    }
}
exports.stateFromJSON = stateFromJSON;
function stateToJSON(object) {
    switch (object) {
        case State.STATE_UNINITIALIZED_UNSPECIFIED:
            return "STATE_UNINITIALIZED_UNSPECIFIED";
        case State.STATE_INIT:
            return "STATE_INIT";
        case State.STATE_TRYOPEN:
            return "STATE_TRYOPEN";
        case State.STATE_OPEN:
            return "STATE_OPEN";
        case State.UNRECOGNIZED:
        default:
            return "UNRECOGNIZED";
    }
}
exports.stateToJSON = stateToJSON;
function createBaseConnectionEnd() {
    return {
        clientId: "",
        versions: [],
        state: 0,
        counterparty: exports.Counterparty.fromPartial({}),
        delayPeriod: BigInt(0),
    };
}
exports.ConnectionEnd = {
    typeUrl: "/ibc.core.connection.v1.ConnectionEnd",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.clientId !== "") {
            writer.uint32(10).string(message.clientId);
        }
        for (const v of message.versions) {
            exports.Version.encode(v, writer.uint32(18).fork()).ldelim();
        }
        if (message.state !== 0) {
            writer.uint32(24).int32(message.state);
        }
        if (message.counterparty !== undefined) {
            exports.Counterparty.encode(message.counterparty, writer.uint32(34).fork()).ldelim();
        }
        if (message.delayPeriod !== BigInt(0)) {
            writer.uint32(40).uint64(message.delayPeriod);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseConnectionEnd();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.clientId = reader.string();
                    break;
                case 2:
                    message.versions.push(exports.Version.decode(reader, reader.uint32()));
                    break;
                case 3:
                    message.state = reader.int32();
                    break;
                case 4:
                    message.counterparty = exports.Counterparty.decode(reader, reader.uint32());
                    break;
                case 5:
                    message.delayPeriod = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseConnectionEnd();
        if ((0, helpers_1.isSet)(object.clientId))
            obj.clientId = String(object.clientId);
        if (Array.isArray(object?.versions))
            obj.versions = object.versions.map((e) => exports.Version.fromJSON(e));
        if ((0, helpers_1.isSet)(object.state))
            obj.state = stateFromJSON(object.state);
        if ((0, helpers_1.isSet)(object.counterparty))
            obj.counterparty = exports.Counterparty.fromJSON(object.counterparty);
        if ((0, helpers_1.isSet)(object.delayPeriod))
            obj.delayPeriod = BigInt(object.delayPeriod.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.clientId !== undefined && (obj.clientId = message.clientId);
        if (message.versions) {
            obj.versions = message.versions.map((e) => (e ? exports.Version.toJSON(e) : undefined));
        }
        else {
            obj.versions = [];
        }
        message.state !== undefined && (obj.state = stateToJSON(message.state));
        message.counterparty !== undefined &&
            (obj.counterparty = message.counterparty ? exports.Counterparty.toJSON(message.counterparty) : undefined);
        message.delayPeriod !== undefined && (obj.delayPeriod = (message.delayPeriod || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseConnectionEnd();
        message.clientId = object.clientId ?? "";
        message.versions = object.versions?.map((e) => exports.Version.fromPartial(e)) || [];
        message.state = object.state ?? 0;
        if (object.counterparty !== undefined && object.counterparty !== null) {
            message.counterparty = exports.Counterparty.fromPartial(object.counterparty);
        }
        if (object.delayPeriod !== undefined && object.delayPeriod !== null) {
            message.delayPeriod = BigInt(object.delayPeriod.toString());
        }
        return message;
    },
};
function createBaseIdentifiedConnection() {
    return {
        id: "",
        clientId: "",
        versions: [],
        state: 0,
        counterparty: exports.Counterparty.fromPartial({}),
        delayPeriod: BigInt(0),
    };
}
exports.IdentifiedConnection = {
    typeUrl: "/ibc.core.connection.v1.IdentifiedConnection",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.id !== "") {
            writer.uint32(10).string(message.id);
        }
        if (message.clientId !== "") {
            writer.uint32(18).string(message.clientId);
        }
        for (const v of message.versions) {
            exports.Version.encode(v, writer.uint32(26).fork()).ldelim();
        }
        if (message.state !== 0) {
            writer.uint32(32).int32(message.state);
        }
        if (message.counterparty !== undefined) {
            exports.Counterparty.encode(message.counterparty, writer.uint32(42).fork()).ldelim();
        }
        if (message.delayPeriod !== BigInt(0)) {
            writer.uint32(48).uint64(message.delayPeriod);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseIdentifiedConnection();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.id = reader.string();
                    break;
                case 2:
                    message.clientId = reader.string();
                    break;
                case 3:
                    message.versions.push(exports.Version.decode(reader, reader.uint32()));
                    break;
                case 4:
                    message.state = reader.int32();
                    break;
                case 5:
                    message.counterparty = exports.Counterparty.decode(reader, reader.uint32());
                    break;
                case 6:
                    message.delayPeriod = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseIdentifiedConnection();
        if ((0, helpers_1.isSet)(object.id))
            obj.id = String(object.id);
        if ((0, helpers_1.isSet)(object.clientId))
            obj.clientId = String(object.clientId);
        if (Array.isArray(object?.versions))
            obj.versions = object.versions.map((e) => exports.Version.fromJSON(e));
        if ((0, helpers_1.isSet)(object.state))
            obj.state = stateFromJSON(object.state);
        if ((0, helpers_1.isSet)(object.counterparty))
            obj.counterparty = exports.Counterparty.fromJSON(object.counterparty);
        if ((0, helpers_1.isSet)(object.delayPeriod))
            obj.delayPeriod = BigInt(object.delayPeriod.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.id !== undefined && (obj.id = message.id);
        message.clientId !== undefined && (obj.clientId = message.clientId);
        if (message.versions) {
            obj.versions = message.versions.map((e) => (e ? exports.Version.toJSON(e) : undefined));
        }
        else {
            obj.versions = [];
        }
        message.state !== undefined && (obj.state = stateToJSON(message.state));
        message.counterparty !== undefined &&
            (obj.counterparty = message.counterparty ? exports.Counterparty.toJSON(message.counterparty) : undefined);
        message.delayPeriod !== undefined && (obj.delayPeriod = (message.delayPeriod || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseIdentifiedConnection();
        message.id = object.id ?? "";
        message.clientId = object.clientId ?? "";
        message.versions = object.versions?.map((e) => exports.Version.fromPartial(e)) || [];
        message.state = object.state ?? 0;
        if (object.counterparty !== undefined && object.counterparty !== null) {
            message.counterparty = exports.Counterparty.fromPartial(object.counterparty);
        }
        if (object.delayPeriod !== undefined && object.delayPeriod !== null) {
            message.delayPeriod = BigInt(object.delayPeriod.toString());
        }
        return message;
    },
};
function createBaseCounterparty() {
    return {
        clientId: "",
        connectionId: "",
        prefix: commitment_1.MerklePrefix.fromPartial({}),
    };
}
exports.Counterparty = {
    typeUrl: "/ibc.core.connection.v1.Counterparty",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.clientId !== "") {
            writer.uint32(10).string(message.clientId);
        }
        if (message.connectionId !== "") {
            writer.uint32(18).string(message.connectionId);
        }
        if (message.prefix !== undefined) {
            commitment_1.MerklePrefix.encode(message.prefix, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseCounterparty();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.clientId = reader.string();
                    break;
                case 2:
                    message.connectionId = reader.string();
                    break;
                case 3:
                    message.prefix = commitment_1.MerklePrefix.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseCounterparty();
        if ((0, helpers_1.isSet)(object.clientId))
            obj.clientId = String(object.clientId);
        if ((0, helpers_1.isSet)(object.connectionId))
            obj.connectionId = String(object.connectionId);
        if ((0, helpers_1.isSet)(object.prefix))
            obj.prefix = commitment_1.MerklePrefix.fromJSON(object.prefix);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.clientId !== undefined && (obj.clientId = message.clientId);
        message.connectionId !== undefined && (obj.connectionId = message.connectionId);
        message.prefix !== undefined &&
            (obj.prefix = message.prefix ? commitment_1.MerklePrefix.toJSON(message.prefix) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseCounterparty();
        message.clientId = object.clientId ?? "";
        message.connectionId = object.connectionId ?? "";
        if (object.prefix !== undefined && object.prefix !== null) {
            message.prefix = commitment_1.MerklePrefix.fromPartial(object.prefix);
        }
        return message;
    },
};
function createBaseClientPaths() {
    return {
        paths: [],
    };
}
exports.ClientPaths = {
    typeUrl: "/ibc.core.connection.v1.ClientPaths",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.paths) {
            writer.uint32(10).string(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseClientPaths();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.paths.push(reader.string());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseClientPaths();
        if (Array.isArray(object?.paths))
            obj.paths = object.paths.map((e) => String(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.paths) {
            obj.paths = message.paths.map((e) => e);
        }
        else {
            obj.paths = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseClientPaths();
        message.paths = object.paths?.map((e) => e) || [];
        return message;
    },
};
function createBaseConnectionPaths() {
    return {
        clientId: "",
        paths: [],
    };
}
exports.ConnectionPaths = {
    typeUrl: "/ibc.core.connection.v1.ConnectionPaths",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.clientId !== "") {
            writer.uint32(10).string(message.clientId);
        }
        for (const v of message.paths) {
            writer.uint32(18).string(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseConnectionPaths();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.clientId = reader.string();
                    break;
                case 2:
                    message.paths.push(reader.string());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseConnectionPaths();
        if ((0, helpers_1.isSet)(object.clientId))
            obj.clientId = String(object.clientId);
        if (Array.isArray(object?.paths))
            obj.paths = object.paths.map((e) => String(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.clientId !== undefined && (obj.clientId = message.clientId);
        if (message.paths) {
            obj.paths = message.paths.map((e) => e);
        }
        else {
            obj.paths = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseConnectionPaths();
        message.clientId = object.clientId ?? "";
        message.paths = object.paths?.map((e) => e) || [];
        return message;
    },
};
function createBaseVersion() {
    return {
        identifier: "",
        features: [],
    };
}
exports.Version = {
    typeUrl: "/ibc.core.connection.v1.Version",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.identifier !== "") {
            writer.uint32(10).string(message.identifier);
        }
        for (const v of message.features) {
            writer.uint32(18).string(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseVersion();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.identifier = reader.string();
                    break;
                case 2:
                    message.features.push(reader.string());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseVersion();
        if ((0, helpers_1.isSet)(object.identifier))
            obj.identifier = String(object.identifier);
        if (Array.isArray(object?.features))
            obj.features = object.features.map((e) => String(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.identifier !== undefined && (obj.identifier = message.identifier);
        if (message.features) {
            obj.features = message.features.map((e) => e);
        }
        else {
            obj.features = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseVersion();
        message.identifier = object.identifier ?? "";
        message.features = object.features?.map((e) => e) || [];
        return message;
    },
};
function createBaseParams() {
    return {
        maxExpectedTimePerBlock: BigInt(0),
    };
}
exports.Params = {
    typeUrl: "/ibc.core.connection.v1.Params",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.maxExpectedTimePerBlock !== BigInt(0)) {
            writer.uint32(8).uint64(message.maxExpectedTimePerBlock);
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
                    message.maxExpectedTimePerBlock = reader.uint64();
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
        if ((0, helpers_1.isSet)(object.maxExpectedTimePerBlock))
            obj.maxExpectedTimePerBlock = BigInt(object.maxExpectedTimePerBlock.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.maxExpectedTimePerBlock !== undefined &&
            (obj.maxExpectedTimePerBlock = (message.maxExpectedTimePerBlock || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseParams();
        if (object.maxExpectedTimePerBlock !== undefined && object.maxExpectedTimePerBlock !== null) {
            message.maxExpectedTimePerBlock = BigInt(object.maxExpectedTimePerBlock.toString());
        }
        return message;
    },
};
//# sourceMappingURL=connection.js.map