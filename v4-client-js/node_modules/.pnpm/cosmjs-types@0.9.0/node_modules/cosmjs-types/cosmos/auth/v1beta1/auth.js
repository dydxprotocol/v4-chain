"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Params = exports.ModuleCredential = exports.ModuleAccount = exports.BaseAccount = exports.protobufPackage = void 0;
/* eslint-disable */
const any_1 = require("../../../google/protobuf/any");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.auth.v1beta1";
function createBaseBaseAccount() {
    return {
        address: "",
        pubKey: undefined,
        accountNumber: BigInt(0),
        sequence: BigInt(0),
    };
}
exports.BaseAccount = {
    typeUrl: "/cosmos.auth.v1beta1.BaseAccount",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.address !== "") {
            writer.uint32(10).string(message.address);
        }
        if (message.pubKey !== undefined) {
            any_1.Any.encode(message.pubKey, writer.uint32(18).fork()).ldelim();
        }
        if (message.accountNumber !== BigInt(0)) {
            writer.uint32(24).uint64(message.accountNumber);
        }
        if (message.sequence !== BigInt(0)) {
            writer.uint32(32).uint64(message.sequence);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseBaseAccount();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.address = reader.string();
                    break;
                case 2:
                    message.pubKey = any_1.Any.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.accountNumber = reader.uint64();
                    break;
                case 4:
                    message.sequence = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseBaseAccount();
        if ((0, helpers_1.isSet)(object.address))
            obj.address = String(object.address);
        if ((0, helpers_1.isSet)(object.pubKey))
            obj.pubKey = any_1.Any.fromJSON(object.pubKey);
        if ((0, helpers_1.isSet)(object.accountNumber))
            obj.accountNumber = BigInt(object.accountNumber.toString());
        if ((0, helpers_1.isSet)(object.sequence))
            obj.sequence = BigInt(object.sequence.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.address !== undefined && (obj.address = message.address);
        message.pubKey !== undefined && (obj.pubKey = message.pubKey ? any_1.Any.toJSON(message.pubKey) : undefined);
        message.accountNumber !== undefined &&
            (obj.accountNumber = (message.accountNumber || BigInt(0)).toString());
        message.sequence !== undefined && (obj.sequence = (message.sequence || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseBaseAccount();
        message.address = object.address ?? "";
        if (object.pubKey !== undefined && object.pubKey !== null) {
            message.pubKey = any_1.Any.fromPartial(object.pubKey);
        }
        if (object.accountNumber !== undefined && object.accountNumber !== null) {
            message.accountNumber = BigInt(object.accountNumber.toString());
        }
        if (object.sequence !== undefined && object.sequence !== null) {
            message.sequence = BigInt(object.sequence.toString());
        }
        return message;
    },
};
function createBaseModuleAccount() {
    return {
        baseAccount: undefined,
        name: "",
        permissions: [],
    };
}
exports.ModuleAccount = {
    typeUrl: "/cosmos.auth.v1beta1.ModuleAccount",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.baseAccount !== undefined) {
            exports.BaseAccount.encode(message.baseAccount, writer.uint32(10).fork()).ldelim();
        }
        if (message.name !== "") {
            writer.uint32(18).string(message.name);
        }
        for (const v of message.permissions) {
            writer.uint32(26).string(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseModuleAccount();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.baseAccount = exports.BaseAccount.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.name = reader.string();
                    break;
                case 3:
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
        const obj = createBaseModuleAccount();
        if ((0, helpers_1.isSet)(object.baseAccount))
            obj.baseAccount = exports.BaseAccount.fromJSON(object.baseAccount);
        if ((0, helpers_1.isSet)(object.name))
            obj.name = String(object.name);
        if (Array.isArray(object?.permissions))
            obj.permissions = object.permissions.map((e) => String(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.baseAccount !== undefined &&
            (obj.baseAccount = message.baseAccount ? exports.BaseAccount.toJSON(message.baseAccount) : undefined);
        message.name !== undefined && (obj.name = message.name);
        if (message.permissions) {
            obj.permissions = message.permissions.map((e) => e);
        }
        else {
            obj.permissions = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseModuleAccount();
        if (object.baseAccount !== undefined && object.baseAccount !== null) {
            message.baseAccount = exports.BaseAccount.fromPartial(object.baseAccount);
        }
        message.name = object.name ?? "";
        message.permissions = object.permissions?.map((e) => e) || [];
        return message;
    },
};
function createBaseModuleCredential() {
    return {
        moduleName: "",
        derivationKeys: [],
    };
}
exports.ModuleCredential = {
    typeUrl: "/cosmos.auth.v1beta1.ModuleCredential",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.moduleName !== "") {
            writer.uint32(10).string(message.moduleName);
        }
        for (const v of message.derivationKeys) {
            writer.uint32(18).bytes(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseModuleCredential();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.moduleName = reader.string();
                    break;
                case 2:
                    message.derivationKeys.push(reader.bytes());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseModuleCredential();
        if ((0, helpers_1.isSet)(object.moduleName))
            obj.moduleName = String(object.moduleName);
        if (Array.isArray(object?.derivationKeys))
            obj.derivationKeys = object.derivationKeys.map((e) => (0, helpers_1.bytesFromBase64)(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.moduleName !== undefined && (obj.moduleName = message.moduleName);
        if (message.derivationKeys) {
            obj.derivationKeys = message.derivationKeys.map((e) => (0, helpers_1.base64FromBytes)(e !== undefined ? e : new Uint8Array()));
        }
        else {
            obj.derivationKeys = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseModuleCredential();
        message.moduleName = object.moduleName ?? "";
        message.derivationKeys = object.derivationKeys?.map((e) => e) || [];
        return message;
    },
};
function createBaseParams() {
    return {
        maxMemoCharacters: BigInt(0),
        txSigLimit: BigInt(0),
        txSizeCostPerByte: BigInt(0),
        sigVerifyCostEd25519: BigInt(0),
        sigVerifyCostSecp256k1: BigInt(0),
    };
}
exports.Params = {
    typeUrl: "/cosmos.auth.v1beta1.Params",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.maxMemoCharacters !== BigInt(0)) {
            writer.uint32(8).uint64(message.maxMemoCharacters);
        }
        if (message.txSigLimit !== BigInt(0)) {
            writer.uint32(16).uint64(message.txSigLimit);
        }
        if (message.txSizeCostPerByte !== BigInt(0)) {
            writer.uint32(24).uint64(message.txSizeCostPerByte);
        }
        if (message.sigVerifyCostEd25519 !== BigInt(0)) {
            writer.uint32(32).uint64(message.sigVerifyCostEd25519);
        }
        if (message.sigVerifyCostSecp256k1 !== BigInt(0)) {
            writer.uint32(40).uint64(message.sigVerifyCostSecp256k1);
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
                    message.maxMemoCharacters = reader.uint64();
                    break;
                case 2:
                    message.txSigLimit = reader.uint64();
                    break;
                case 3:
                    message.txSizeCostPerByte = reader.uint64();
                    break;
                case 4:
                    message.sigVerifyCostEd25519 = reader.uint64();
                    break;
                case 5:
                    message.sigVerifyCostSecp256k1 = reader.uint64();
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
        if ((0, helpers_1.isSet)(object.maxMemoCharacters))
            obj.maxMemoCharacters = BigInt(object.maxMemoCharacters.toString());
        if ((0, helpers_1.isSet)(object.txSigLimit))
            obj.txSigLimit = BigInt(object.txSigLimit.toString());
        if ((0, helpers_1.isSet)(object.txSizeCostPerByte))
            obj.txSizeCostPerByte = BigInt(object.txSizeCostPerByte.toString());
        if ((0, helpers_1.isSet)(object.sigVerifyCostEd25519))
            obj.sigVerifyCostEd25519 = BigInt(object.sigVerifyCostEd25519.toString());
        if ((0, helpers_1.isSet)(object.sigVerifyCostSecp256k1))
            obj.sigVerifyCostSecp256k1 = BigInt(object.sigVerifyCostSecp256k1.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.maxMemoCharacters !== undefined &&
            (obj.maxMemoCharacters = (message.maxMemoCharacters || BigInt(0)).toString());
        message.txSigLimit !== undefined && (obj.txSigLimit = (message.txSigLimit || BigInt(0)).toString());
        message.txSizeCostPerByte !== undefined &&
            (obj.txSizeCostPerByte = (message.txSizeCostPerByte || BigInt(0)).toString());
        message.sigVerifyCostEd25519 !== undefined &&
            (obj.sigVerifyCostEd25519 = (message.sigVerifyCostEd25519 || BigInt(0)).toString());
        message.sigVerifyCostSecp256k1 !== undefined &&
            (obj.sigVerifyCostSecp256k1 = (message.sigVerifyCostSecp256k1 || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseParams();
        if (object.maxMemoCharacters !== undefined && object.maxMemoCharacters !== null) {
            message.maxMemoCharacters = BigInt(object.maxMemoCharacters.toString());
        }
        if (object.txSigLimit !== undefined && object.txSigLimit !== null) {
            message.txSigLimit = BigInt(object.txSigLimit.toString());
        }
        if (object.txSizeCostPerByte !== undefined && object.txSizeCostPerByte !== null) {
            message.txSizeCostPerByte = BigInt(object.txSizeCostPerByte.toString());
        }
        if (object.sigVerifyCostEd25519 !== undefined && object.sigVerifyCostEd25519 !== null) {
            message.sigVerifyCostEd25519 = BigInt(object.sigVerifyCostEd25519.toString());
        }
        if (object.sigVerifyCostSecp256k1 !== undefined && object.sigVerifyCostSecp256k1 !== null) {
            message.sigVerifyCostSecp256k1 = BigInt(object.sigVerifyCostSecp256k1.toString());
        }
        return message;
    },
};
//# sourceMappingURL=auth.js.map