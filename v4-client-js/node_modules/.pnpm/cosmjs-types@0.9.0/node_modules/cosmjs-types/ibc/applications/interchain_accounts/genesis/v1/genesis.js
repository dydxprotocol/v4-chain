"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.RegisteredInterchainAccount = exports.ActiveChannel = exports.HostGenesisState = exports.ControllerGenesisState = exports.GenesisState = exports.protobufPackage = void 0;
/* eslint-disable */
const controller_1 = require("../../controller/v1/controller");
const host_1 = require("../../host/v1/host");
const binary_1 = require("../../../../../binary");
const helpers_1 = require("../../../../../helpers");
exports.protobufPackage = "ibc.applications.interchain_accounts.genesis.v1";
function createBaseGenesisState() {
    return {
        controllerGenesisState: exports.ControllerGenesisState.fromPartial({}),
        hostGenesisState: exports.HostGenesisState.fromPartial({}),
    };
}
exports.GenesisState = {
    typeUrl: "/ibc.applications.interchain_accounts.genesis.v1.GenesisState",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.controllerGenesisState !== undefined) {
            exports.ControllerGenesisState.encode(message.controllerGenesisState, writer.uint32(10).fork()).ldelim();
        }
        if (message.hostGenesisState !== undefined) {
            exports.HostGenesisState.encode(message.hostGenesisState, writer.uint32(18).fork()).ldelim();
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
                    message.controllerGenesisState = exports.ControllerGenesisState.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.hostGenesisState = exports.HostGenesisState.decode(reader, reader.uint32());
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
        if ((0, helpers_1.isSet)(object.controllerGenesisState))
            obj.controllerGenesisState = exports.ControllerGenesisState.fromJSON(object.controllerGenesisState);
        if ((0, helpers_1.isSet)(object.hostGenesisState))
            obj.hostGenesisState = exports.HostGenesisState.fromJSON(object.hostGenesisState);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.controllerGenesisState !== undefined &&
            (obj.controllerGenesisState = message.controllerGenesisState
                ? exports.ControllerGenesisState.toJSON(message.controllerGenesisState)
                : undefined);
        message.hostGenesisState !== undefined &&
            (obj.hostGenesisState = message.hostGenesisState
                ? exports.HostGenesisState.toJSON(message.hostGenesisState)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGenesisState();
        if (object.controllerGenesisState !== undefined && object.controllerGenesisState !== null) {
            message.controllerGenesisState = exports.ControllerGenesisState.fromPartial(object.controllerGenesisState);
        }
        if (object.hostGenesisState !== undefined && object.hostGenesisState !== null) {
            message.hostGenesisState = exports.HostGenesisState.fromPartial(object.hostGenesisState);
        }
        return message;
    },
};
function createBaseControllerGenesisState() {
    return {
        activeChannels: [],
        interchainAccounts: [],
        ports: [],
        params: controller_1.Params.fromPartial({}),
    };
}
exports.ControllerGenesisState = {
    typeUrl: "/ibc.applications.interchain_accounts.genesis.v1.ControllerGenesisState",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.activeChannels) {
            exports.ActiveChannel.encode(v, writer.uint32(10).fork()).ldelim();
        }
        for (const v of message.interchainAccounts) {
            exports.RegisteredInterchainAccount.encode(v, writer.uint32(18).fork()).ldelim();
        }
        for (const v of message.ports) {
            writer.uint32(26).string(v);
        }
        if (message.params !== undefined) {
            controller_1.Params.encode(message.params, writer.uint32(34).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseControllerGenesisState();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.activeChannels.push(exports.ActiveChannel.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.interchainAccounts.push(exports.RegisteredInterchainAccount.decode(reader, reader.uint32()));
                    break;
                case 3:
                    message.ports.push(reader.string());
                    break;
                case 4:
                    message.params = controller_1.Params.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseControllerGenesisState();
        if (Array.isArray(object?.activeChannels))
            obj.activeChannels = object.activeChannels.map((e) => exports.ActiveChannel.fromJSON(e));
        if (Array.isArray(object?.interchainAccounts))
            obj.interchainAccounts = object.interchainAccounts.map((e) => exports.RegisteredInterchainAccount.fromJSON(e));
        if (Array.isArray(object?.ports))
            obj.ports = object.ports.map((e) => String(e));
        if ((0, helpers_1.isSet)(object.params))
            obj.params = controller_1.Params.fromJSON(object.params);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.activeChannels) {
            obj.activeChannels = message.activeChannels.map((e) => (e ? exports.ActiveChannel.toJSON(e) : undefined));
        }
        else {
            obj.activeChannels = [];
        }
        if (message.interchainAccounts) {
            obj.interchainAccounts = message.interchainAccounts.map((e) => e ? exports.RegisteredInterchainAccount.toJSON(e) : undefined);
        }
        else {
            obj.interchainAccounts = [];
        }
        if (message.ports) {
            obj.ports = message.ports.map((e) => e);
        }
        else {
            obj.ports = [];
        }
        message.params !== undefined &&
            (obj.params = message.params ? controller_1.Params.toJSON(message.params) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseControllerGenesisState();
        message.activeChannels = object.activeChannels?.map((e) => exports.ActiveChannel.fromPartial(e)) || [];
        message.interchainAccounts =
            object.interchainAccounts?.map((e) => exports.RegisteredInterchainAccount.fromPartial(e)) || [];
        message.ports = object.ports?.map((e) => e) || [];
        if (object.params !== undefined && object.params !== null) {
            message.params = controller_1.Params.fromPartial(object.params);
        }
        return message;
    },
};
function createBaseHostGenesisState() {
    return {
        activeChannels: [],
        interchainAccounts: [],
        port: "",
        params: host_1.Params.fromPartial({}),
    };
}
exports.HostGenesisState = {
    typeUrl: "/ibc.applications.interchain_accounts.genesis.v1.HostGenesisState",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.activeChannels) {
            exports.ActiveChannel.encode(v, writer.uint32(10).fork()).ldelim();
        }
        for (const v of message.interchainAccounts) {
            exports.RegisteredInterchainAccount.encode(v, writer.uint32(18).fork()).ldelim();
        }
        if (message.port !== "") {
            writer.uint32(26).string(message.port);
        }
        if (message.params !== undefined) {
            host_1.Params.encode(message.params, writer.uint32(34).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseHostGenesisState();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.activeChannels.push(exports.ActiveChannel.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.interchainAccounts.push(exports.RegisteredInterchainAccount.decode(reader, reader.uint32()));
                    break;
                case 3:
                    message.port = reader.string();
                    break;
                case 4:
                    message.params = host_1.Params.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseHostGenesisState();
        if (Array.isArray(object?.activeChannels))
            obj.activeChannels = object.activeChannels.map((e) => exports.ActiveChannel.fromJSON(e));
        if (Array.isArray(object?.interchainAccounts))
            obj.interchainAccounts = object.interchainAccounts.map((e) => exports.RegisteredInterchainAccount.fromJSON(e));
        if ((0, helpers_1.isSet)(object.port))
            obj.port = String(object.port);
        if ((0, helpers_1.isSet)(object.params))
            obj.params = host_1.Params.fromJSON(object.params);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.activeChannels) {
            obj.activeChannels = message.activeChannels.map((e) => (e ? exports.ActiveChannel.toJSON(e) : undefined));
        }
        else {
            obj.activeChannels = [];
        }
        if (message.interchainAccounts) {
            obj.interchainAccounts = message.interchainAccounts.map((e) => e ? exports.RegisteredInterchainAccount.toJSON(e) : undefined);
        }
        else {
            obj.interchainAccounts = [];
        }
        message.port !== undefined && (obj.port = message.port);
        message.params !== undefined &&
            (obj.params = message.params ? host_1.Params.toJSON(message.params) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseHostGenesisState();
        message.activeChannels = object.activeChannels?.map((e) => exports.ActiveChannel.fromPartial(e)) || [];
        message.interchainAccounts =
            object.interchainAccounts?.map((e) => exports.RegisteredInterchainAccount.fromPartial(e)) || [];
        message.port = object.port ?? "";
        if (object.params !== undefined && object.params !== null) {
            message.params = host_1.Params.fromPartial(object.params);
        }
        return message;
    },
};
function createBaseActiveChannel() {
    return {
        connectionId: "",
        portId: "",
        channelId: "",
        isMiddlewareEnabled: false,
    };
}
exports.ActiveChannel = {
    typeUrl: "/ibc.applications.interchain_accounts.genesis.v1.ActiveChannel",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.connectionId !== "") {
            writer.uint32(10).string(message.connectionId);
        }
        if (message.portId !== "") {
            writer.uint32(18).string(message.portId);
        }
        if (message.channelId !== "") {
            writer.uint32(26).string(message.channelId);
        }
        if (message.isMiddlewareEnabled === true) {
            writer.uint32(32).bool(message.isMiddlewareEnabled);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseActiveChannel();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.connectionId = reader.string();
                    break;
                case 2:
                    message.portId = reader.string();
                    break;
                case 3:
                    message.channelId = reader.string();
                    break;
                case 4:
                    message.isMiddlewareEnabled = reader.bool();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseActiveChannel();
        if ((0, helpers_1.isSet)(object.connectionId))
            obj.connectionId = String(object.connectionId);
        if ((0, helpers_1.isSet)(object.portId))
            obj.portId = String(object.portId);
        if ((0, helpers_1.isSet)(object.channelId))
            obj.channelId = String(object.channelId);
        if ((0, helpers_1.isSet)(object.isMiddlewareEnabled))
            obj.isMiddlewareEnabled = Boolean(object.isMiddlewareEnabled);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.connectionId !== undefined && (obj.connectionId = message.connectionId);
        message.portId !== undefined && (obj.portId = message.portId);
        message.channelId !== undefined && (obj.channelId = message.channelId);
        message.isMiddlewareEnabled !== undefined && (obj.isMiddlewareEnabled = message.isMiddlewareEnabled);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseActiveChannel();
        message.connectionId = object.connectionId ?? "";
        message.portId = object.portId ?? "";
        message.channelId = object.channelId ?? "";
        message.isMiddlewareEnabled = object.isMiddlewareEnabled ?? false;
        return message;
    },
};
function createBaseRegisteredInterchainAccount() {
    return {
        connectionId: "",
        portId: "",
        accountAddress: "",
    };
}
exports.RegisteredInterchainAccount = {
    typeUrl: "/ibc.applications.interchain_accounts.genesis.v1.RegisteredInterchainAccount",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.connectionId !== "") {
            writer.uint32(10).string(message.connectionId);
        }
        if (message.portId !== "") {
            writer.uint32(18).string(message.portId);
        }
        if (message.accountAddress !== "") {
            writer.uint32(26).string(message.accountAddress);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseRegisteredInterchainAccount();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.connectionId = reader.string();
                    break;
                case 2:
                    message.portId = reader.string();
                    break;
                case 3:
                    message.accountAddress = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseRegisteredInterchainAccount();
        if ((0, helpers_1.isSet)(object.connectionId))
            obj.connectionId = String(object.connectionId);
        if ((0, helpers_1.isSet)(object.portId))
            obj.portId = String(object.portId);
        if ((0, helpers_1.isSet)(object.accountAddress))
            obj.accountAddress = String(object.accountAddress);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.connectionId !== undefined && (obj.connectionId = message.connectionId);
        message.portId !== undefined && (obj.portId = message.portId);
        message.accountAddress !== undefined && (obj.accountAddress = message.accountAddress);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseRegisteredInterchainAccount();
        message.connectionId = object.connectionId ?? "";
        message.portId = object.portId ?? "";
        message.accountAddress = object.accountAddress ?? "";
        return message;
    },
};
//# sourceMappingURL=genesis.js.map