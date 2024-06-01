"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GolangBinding = exports.ModuleConfig = exports.Config = exports.protobufPackage = void 0;
/* eslint-disable */
const any_1 = require("../../../google/protobuf/any");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.app.v1alpha1";
function createBaseConfig() {
    return {
        modules: [],
        golangBindings: [],
    };
}
exports.Config = {
    typeUrl: "/cosmos.app.v1alpha1.Config",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.modules) {
            exports.ModuleConfig.encode(v, writer.uint32(10).fork()).ldelim();
        }
        for (const v of message.golangBindings) {
            exports.GolangBinding.encode(v, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseConfig();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.modules.push(exports.ModuleConfig.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.golangBindings.push(exports.GolangBinding.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseConfig();
        if (Array.isArray(object?.modules))
            obj.modules = object.modules.map((e) => exports.ModuleConfig.fromJSON(e));
        if (Array.isArray(object?.golangBindings))
            obj.golangBindings = object.golangBindings.map((e) => exports.GolangBinding.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.modules) {
            obj.modules = message.modules.map((e) => (e ? exports.ModuleConfig.toJSON(e) : undefined));
        }
        else {
            obj.modules = [];
        }
        if (message.golangBindings) {
            obj.golangBindings = message.golangBindings.map((e) => (e ? exports.GolangBinding.toJSON(e) : undefined));
        }
        else {
            obj.golangBindings = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseConfig();
        message.modules = object.modules?.map((e) => exports.ModuleConfig.fromPartial(e)) || [];
        message.golangBindings = object.golangBindings?.map((e) => exports.GolangBinding.fromPartial(e)) || [];
        return message;
    },
};
function createBaseModuleConfig() {
    return {
        name: "",
        config: undefined,
        golangBindings: [],
    };
}
exports.ModuleConfig = {
    typeUrl: "/cosmos.app.v1alpha1.ModuleConfig",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.name !== "") {
            writer.uint32(10).string(message.name);
        }
        if (message.config !== undefined) {
            any_1.Any.encode(message.config, writer.uint32(18).fork()).ldelim();
        }
        for (const v of message.golangBindings) {
            exports.GolangBinding.encode(v, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseModuleConfig();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.name = reader.string();
                    break;
                case 2:
                    message.config = any_1.Any.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.golangBindings.push(exports.GolangBinding.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseModuleConfig();
        if ((0, helpers_1.isSet)(object.name))
            obj.name = String(object.name);
        if ((0, helpers_1.isSet)(object.config))
            obj.config = any_1.Any.fromJSON(object.config);
        if (Array.isArray(object?.golangBindings))
            obj.golangBindings = object.golangBindings.map((e) => exports.GolangBinding.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.name !== undefined && (obj.name = message.name);
        message.config !== undefined && (obj.config = message.config ? any_1.Any.toJSON(message.config) : undefined);
        if (message.golangBindings) {
            obj.golangBindings = message.golangBindings.map((e) => (e ? exports.GolangBinding.toJSON(e) : undefined));
        }
        else {
            obj.golangBindings = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseModuleConfig();
        message.name = object.name ?? "";
        if (object.config !== undefined && object.config !== null) {
            message.config = any_1.Any.fromPartial(object.config);
        }
        message.golangBindings = object.golangBindings?.map((e) => exports.GolangBinding.fromPartial(e)) || [];
        return message;
    },
};
function createBaseGolangBinding() {
    return {
        interfaceType: "",
        implementation: "",
    };
}
exports.GolangBinding = {
    typeUrl: "/cosmos.app.v1alpha1.GolangBinding",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.interfaceType !== "") {
            writer.uint32(10).string(message.interfaceType);
        }
        if (message.implementation !== "") {
            writer.uint32(18).string(message.implementation);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGolangBinding();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.interfaceType = reader.string();
                    break;
                case 2:
                    message.implementation = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseGolangBinding();
        if ((0, helpers_1.isSet)(object.interfaceType))
            obj.interfaceType = String(object.interfaceType);
        if ((0, helpers_1.isSet)(object.implementation))
            obj.implementation = String(object.implementation);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.interfaceType !== undefined && (obj.interfaceType = message.interfaceType);
        message.implementation !== undefined && (obj.implementation = message.implementation);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGolangBinding();
        message.interfaceType = object.interfaceType ?? "";
        message.implementation = object.implementation ?? "";
        return message;
    },
};
//# sourceMappingURL=config.js.map