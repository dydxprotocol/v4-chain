"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.StoreKeyConfig = exports.Module = exports.protobufPackage = void 0;
/* eslint-disable */
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "cosmos.app.runtime.v1alpha1";
function createBaseModule() {
    return {
        appName: "",
        beginBlockers: [],
        endBlockers: [],
        initGenesis: [],
        exportGenesis: [],
        overrideStoreKeys: [],
    };
}
exports.Module = {
    typeUrl: "/cosmos.app.runtime.v1alpha1.Module",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.appName !== "") {
            writer.uint32(10).string(message.appName);
        }
        for (const v of message.beginBlockers) {
            writer.uint32(18).string(v);
        }
        for (const v of message.endBlockers) {
            writer.uint32(26).string(v);
        }
        for (const v of message.initGenesis) {
            writer.uint32(34).string(v);
        }
        for (const v of message.exportGenesis) {
            writer.uint32(42).string(v);
        }
        for (const v of message.overrideStoreKeys) {
            exports.StoreKeyConfig.encode(v, writer.uint32(50).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseModule();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.appName = reader.string();
                    break;
                case 2:
                    message.beginBlockers.push(reader.string());
                    break;
                case 3:
                    message.endBlockers.push(reader.string());
                    break;
                case 4:
                    message.initGenesis.push(reader.string());
                    break;
                case 5:
                    message.exportGenesis.push(reader.string());
                    break;
                case 6:
                    message.overrideStoreKeys.push(exports.StoreKeyConfig.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseModule();
        if ((0, helpers_1.isSet)(object.appName))
            obj.appName = String(object.appName);
        if (Array.isArray(object?.beginBlockers))
            obj.beginBlockers = object.beginBlockers.map((e) => String(e));
        if (Array.isArray(object?.endBlockers))
            obj.endBlockers = object.endBlockers.map((e) => String(e));
        if (Array.isArray(object?.initGenesis))
            obj.initGenesis = object.initGenesis.map((e) => String(e));
        if (Array.isArray(object?.exportGenesis))
            obj.exportGenesis = object.exportGenesis.map((e) => String(e));
        if (Array.isArray(object?.overrideStoreKeys))
            obj.overrideStoreKeys = object.overrideStoreKeys.map((e) => exports.StoreKeyConfig.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.appName !== undefined && (obj.appName = message.appName);
        if (message.beginBlockers) {
            obj.beginBlockers = message.beginBlockers.map((e) => e);
        }
        else {
            obj.beginBlockers = [];
        }
        if (message.endBlockers) {
            obj.endBlockers = message.endBlockers.map((e) => e);
        }
        else {
            obj.endBlockers = [];
        }
        if (message.initGenesis) {
            obj.initGenesis = message.initGenesis.map((e) => e);
        }
        else {
            obj.initGenesis = [];
        }
        if (message.exportGenesis) {
            obj.exportGenesis = message.exportGenesis.map((e) => e);
        }
        else {
            obj.exportGenesis = [];
        }
        if (message.overrideStoreKeys) {
            obj.overrideStoreKeys = message.overrideStoreKeys.map((e) => e ? exports.StoreKeyConfig.toJSON(e) : undefined);
        }
        else {
            obj.overrideStoreKeys = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseModule();
        message.appName = object.appName ?? "";
        message.beginBlockers = object.beginBlockers?.map((e) => e) || [];
        message.endBlockers = object.endBlockers?.map((e) => e) || [];
        message.initGenesis = object.initGenesis?.map((e) => e) || [];
        message.exportGenesis = object.exportGenesis?.map((e) => e) || [];
        message.overrideStoreKeys = object.overrideStoreKeys?.map((e) => exports.StoreKeyConfig.fromPartial(e)) || [];
        return message;
    },
};
function createBaseStoreKeyConfig() {
    return {
        moduleName: "",
        kvStoreKey: "",
    };
}
exports.StoreKeyConfig = {
    typeUrl: "/cosmos.app.runtime.v1alpha1.StoreKeyConfig",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.moduleName !== "") {
            writer.uint32(10).string(message.moduleName);
        }
        if (message.kvStoreKey !== "") {
            writer.uint32(18).string(message.kvStoreKey);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseStoreKeyConfig();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.moduleName = reader.string();
                    break;
                case 2:
                    message.kvStoreKey = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseStoreKeyConfig();
        if ((0, helpers_1.isSet)(object.moduleName))
            obj.moduleName = String(object.moduleName);
        if ((0, helpers_1.isSet)(object.kvStoreKey))
            obj.kvStoreKey = String(object.kvStoreKey);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.moduleName !== undefined && (obj.moduleName = message.moduleName);
        message.kvStoreKey !== undefined && (obj.kvStoreKey = message.kvStoreKey);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseStoreKeyConfig();
        message.moduleName = object.moduleName ?? "";
        message.kvStoreKey = object.kvStoreKey ?? "";
        return message;
    },
};
//# sourceMappingURL=module.js.map