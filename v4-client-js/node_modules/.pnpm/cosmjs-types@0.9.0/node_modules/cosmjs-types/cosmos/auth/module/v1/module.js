"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ModuleAccountPermission = exports.Module = exports.protobufPackage = void 0;
/* eslint-disable */
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "cosmos.auth.module.v1";
function createBaseModule() {
    return {
        bech32Prefix: "",
        moduleAccountPermissions: [],
        authority: "",
    };
}
exports.Module = {
    typeUrl: "/cosmos.auth.module.v1.Module",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.bech32Prefix !== "") {
            writer.uint32(10).string(message.bech32Prefix);
        }
        for (const v of message.moduleAccountPermissions) {
            exports.ModuleAccountPermission.encode(v, writer.uint32(18).fork()).ldelim();
        }
        if (message.authority !== "") {
            writer.uint32(26).string(message.authority);
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
                    message.bech32Prefix = reader.string();
                    break;
                case 2:
                    message.moduleAccountPermissions.push(exports.ModuleAccountPermission.decode(reader, reader.uint32()));
                    break;
                case 3:
                    message.authority = reader.string();
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
        if ((0, helpers_1.isSet)(object.bech32Prefix))
            obj.bech32Prefix = String(object.bech32Prefix);
        if (Array.isArray(object?.moduleAccountPermissions))
            obj.moduleAccountPermissions = object.moduleAccountPermissions.map((e) => exports.ModuleAccountPermission.fromJSON(e));
        if ((0, helpers_1.isSet)(object.authority))
            obj.authority = String(object.authority);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.bech32Prefix !== undefined && (obj.bech32Prefix = message.bech32Prefix);
        if (message.moduleAccountPermissions) {
            obj.moduleAccountPermissions = message.moduleAccountPermissions.map((e) => e ? exports.ModuleAccountPermission.toJSON(e) : undefined);
        }
        else {
            obj.moduleAccountPermissions = [];
        }
        message.authority !== undefined && (obj.authority = message.authority);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseModule();
        message.bech32Prefix = object.bech32Prefix ?? "";
        message.moduleAccountPermissions =
            object.moduleAccountPermissions?.map((e) => exports.ModuleAccountPermission.fromPartial(e)) || [];
        message.authority = object.authority ?? "";
        return message;
    },
};
function createBaseModuleAccountPermission() {
    return {
        account: "",
        permissions: [],
    };
}
exports.ModuleAccountPermission = {
    typeUrl: "/cosmos.auth.module.v1.ModuleAccountPermission",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.account !== "") {
            writer.uint32(10).string(message.account);
        }
        for (const v of message.permissions) {
            writer.uint32(18).string(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseModuleAccountPermission();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.account = reader.string();
                    break;
                case 2:
                    message.permissions.push(reader.string());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseModuleAccountPermission();
        if ((0, helpers_1.isSet)(object.account))
            obj.account = String(object.account);
        if (Array.isArray(object?.permissions))
            obj.permissions = object.permissions.map((e) => String(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.account !== undefined && (obj.account = message.account);
        if (message.permissions) {
            obj.permissions = message.permissions.map((e) => e);
        }
        else {
            obj.permissions = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseModuleAccountPermission();
        message.account = object.account ?? "";
        message.permissions = object.permissions?.map((e) => e) || [];
        return message;
    },
};
//# sourceMappingURL=module.js.map