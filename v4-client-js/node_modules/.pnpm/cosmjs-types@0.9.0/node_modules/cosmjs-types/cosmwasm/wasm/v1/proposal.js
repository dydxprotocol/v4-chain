"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.StoreAndInstantiateContractProposal = exports.UpdateInstantiateConfigProposal = exports.AccessConfigUpdate = exports.UnpinCodesProposal = exports.PinCodesProposal = exports.ClearAdminProposal = exports.UpdateAdminProposal = exports.ExecuteContractProposal = exports.SudoContractProposal = exports.MigrateContractProposal = exports.InstantiateContract2Proposal = exports.InstantiateContractProposal = exports.StoreCodeProposal = exports.protobufPackage = void 0;
/* eslint-disable */
const types_1 = require("./types");
const coin_1 = require("../../../cosmos/base/v1beta1/coin");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmwasm.wasm.v1";
function createBaseStoreCodeProposal() {
    return {
        title: "",
        description: "",
        runAs: "",
        wasmByteCode: new Uint8Array(),
        instantiatePermission: undefined,
        unpinCode: false,
        source: "",
        builder: "",
        codeHash: new Uint8Array(),
    };
}
exports.StoreCodeProposal = {
    typeUrl: "/cosmwasm.wasm.v1.StoreCodeProposal",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.title !== "") {
            writer.uint32(10).string(message.title);
        }
        if (message.description !== "") {
            writer.uint32(18).string(message.description);
        }
        if (message.runAs !== "") {
            writer.uint32(26).string(message.runAs);
        }
        if (message.wasmByteCode.length !== 0) {
            writer.uint32(34).bytes(message.wasmByteCode);
        }
        if (message.instantiatePermission !== undefined) {
            types_1.AccessConfig.encode(message.instantiatePermission, writer.uint32(58).fork()).ldelim();
        }
        if (message.unpinCode === true) {
            writer.uint32(64).bool(message.unpinCode);
        }
        if (message.source !== "") {
            writer.uint32(74).string(message.source);
        }
        if (message.builder !== "") {
            writer.uint32(82).string(message.builder);
        }
        if (message.codeHash.length !== 0) {
            writer.uint32(90).bytes(message.codeHash);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseStoreCodeProposal();
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
                    message.runAs = reader.string();
                    break;
                case 4:
                    message.wasmByteCode = reader.bytes();
                    break;
                case 7:
                    message.instantiatePermission = types_1.AccessConfig.decode(reader, reader.uint32());
                    break;
                case 8:
                    message.unpinCode = reader.bool();
                    break;
                case 9:
                    message.source = reader.string();
                    break;
                case 10:
                    message.builder = reader.string();
                    break;
                case 11:
                    message.codeHash = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseStoreCodeProposal();
        if ((0, helpers_1.isSet)(object.title))
            obj.title = String(object.title);
        if ((0, helpers_1.isSet)(object.description))
            obj.description = String(object.description);
        if ((0, helpers_1.isSet)(object.runAs))
            obj.runAs = String(object.runAs);
        if ((0, helpers_1.isSet)(object.wasmByteCode))
            obj.wasmByteCode = (0, helpers_1.bytesFromBase64)(object.wasmByteCode);
        if ((0, helpers_1.isSet)(object.instantiatePermission))
            obj.instantiatePermission = types_1.AccessConfig.fromJSON(object.instantiatePermission);
        if ((0, helpers_1.isSet)(object.unpinCode))
            obj.unpinCode = Boolean(object.unpinCode);
        if ((0, helpers_1.isSet)(object.source))
            obj.source = String(object.source);
        if ((0, helpers_1.isSet)(object.builder))
            obj.builder = String(object.builder);
        if ((0, helpers_1.isSet)(object.codeHash))
            obj.codeHash = (0, helpers_1.bytesFromBase64)(object.codeHash);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.title !== undefined && (obj.title = message.title);
        message.description !== undefined && (obj.description = message.description);
        message.runAs !== undefined && (obj.runAs = message.runAs);
        message.wasmByteCode !== undefined &&
            (obj.wasmByteCode = (0, helpers_1.base64FromBytes)(message.wasmByteCode !== undefined ? message.wasmByteCode : new Uint8Array()));
        message.instantiatePermission !== undefined &&
            (obj.instantiatePermission = message.instantiatePermission
                ? types_1.AccessConfig.toJSON(message.instantiatePermission)
                : undefined);
        message.unpinCode !== undefined && (obj.unpinCode = message.unpinCode);
        message.source !== undefined && (obj.source = message.source);
        message.builder !== undefined && (obj.builder = message.builder);
        message.codeHash !== undefined &&
            (obj.codeHash = (0, helpers_1.base64FromBytes)(message.codeHash !== undefined ? message.codeHash : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseStoreCodeProposal();
        message.title = object.title ?? "";
        message.description = object.description ?? "";
        message.runAs = object.runAs ?? "";
        message.wasmByteCode = object.wasmByteCode ?? new Uint8Array();
        if (object.instantiatePermission !== undefined && object.instantiatePermission !== null) {
            message.instantiatePermission = types_1.AccessConfig.fromPartial(object.instantiatePermission);
        }
        message.unpinCode = object.unpinCode ?? false;
        message.source = object.source ?? "";
        message.builder = object.builder ?? "";
        message.codeHash = object.codeHash ?? new Uint8Array();
        return message;
    },
};
function createBaseInstantiateContractProposal() {
    return {
        title: "",
        description: "",
        runAs: "",
        admin: "",
        codeId: BigInt(0),
        label: "",
        msg: new Uint8Array(),
        funds: [],
    };
}
exports.InstantiateContractProposal = {
    typeUrl: "/cosmwasm.wasm.v1.InstantiateContractProposal",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.title !== "") {
            writer.uint32(10).string(message.title);
        }
        if (message.description !== "") {
            writer.uint32(18).string(message.description);
        }
        if (message.runAs !== "") {
            writer.uint32(26).string(message.runAs);
        }
        if (message.admin !== "") {
            writer.uint32(34).string(message.admin);
        }
        if (message.codeId !== BigInt(0)) {
            writer.uint32(40).uint64(message.codeId);
        }
        if (message.label !== "") {
            writer.uint32(50).string(message.label);
        }
        if (message.msg.length !== 0) {
            writer.uint32(58).bytes(message.msg);
        }
        for (const v of message.funds) {
            coin_1.Coin.encode(v, writer.uint32(66).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseInstantiateContractProposal();
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
                    message.runAs = reader.string();
                    break;
                case 4:
                    message.admin = reader.string();
                    break;
                case 5:
                    message.codeId = reader.uint64();
                    break;
                case 6:
                    message.label = reader.string();
                    break;
                case 7:
                    message.msg = reader.bytes();
                    break;
                case 8:
                    message.funds.push(coin_1.Coin.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseInstantiateContractProposal();
        if ((0, helpers_1.isSet)(object.title))
            obj.title = String(object.title);
        if ((0, helpers_1.isSet)(object.description))
            obj.description = String(object.description);
        if ((0, helpers_1.isSet)(object.runAs))
            obj.runAs = String(object.runAs);
        if ((0, helpers_1.isSet)(object.admin))
            obj.admin = String(object.admin);
        if ((0, helpers_1.isSet)(object.codeId))
            obj.codeId = BigInt(object.codeId.toString());
        if ((0, helpers_1.isSet)(object.label))
            obj.label = String(object.label);
        if ((0, helpers_1.isSet)(object.msg))
            obj.msg = (0, helpers_1.bytesFromBase64)(object.msg);
        if (Array.isArray(object?.funds))
            obj.funds = object.funds.map((e) => coin_1.Coin.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.title !== undefined && (obj.title = message.title);
        message.description !== undefined && (obj.description = message.description);
        message.runAs !== undefined && (obj.runAs = message.runAs);
        message.admin !== undefined && (obj.admin = message.admin);
        message.codeId !== undefined && (obj.codeId = (message.codeId || BigInt(0)).toString());
        message.label !== undefined && (obj.label = message.label);
        message.msg !== undefined &&
            (obj.msg = (0, helpers_1.base64FromBytes)(message.msg !== undefined ? message.msg : new Uint8Array()));
        if (message.funds) {
            obj.funds = message.funds.map((e) => (e ? coin_1.Coin.toJSON(e) : undefined));
        }
        else {
            obj.funds = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseInstantiateContractProposal();
        message.title = object.title ?? "";
        message.description = object.description ?? "";
        message.runAs = object.runAs ?? "";
        message.admin = object.admin ?? "";
        if (object.codeId !== undefined && object.codeId !== null) {
            message.codeId = BigInt(object.codeId.toString());
        }
        message.label = object.label ?? "";
        message.msg = object.msg ?? new Uint8Array();
        message.funds = object.funds?.map((e) => coin_1.Coin.fromPartial(e)) || [];
        return message;
    },
};
function createBaseInstantiateContract2Proposal() {
    return {
        title: "",
        description: "",
        runAs: "",
        admin: "",
        codeId: BigInt(0),
        label: "",
        msg: new Uint8Array(),
        funds: [],
        salt: new Uint8Array(),
        fixMsg: false,
    };
}
exports.InstantiateContract2Proposal = {
    typeUrl: "/cosmwasm.wasm.v1.InstantiateContract2Proposal",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.title !== "") {
            writer.uint32(10).string(message.title);
        }
        if (message.description !== "") {
            writer.uint32(18).string(message.description);
        }
        if (message.runAs !== "") {
            writer.uint32(26).string(message.runAs);
        }
        if (message.admin !== "") {
            writer.uint32(34).string(message.admin);
        }
        if (message.codeId !== BigInt(0)) {
            writer.uint32(40).uint64(message.codeId);
        }
        if (message.label !== "") {
            writer.uint32(50).string(message.label);
        }
        if (message.msg.length !== 0) {
            writer.uint32(58).bytes(message.msg);
        }
        for (const v of message.funds) {
            coin_1.Coin.encode(v, writer.uint32(66).fork()).ldelim();
        }
        if (message.salt.length !== 0) {
            writer.uint32(74).bytes(message.salt);
        }
        if (message.fixMsg === true) {
            writer.uint32(80).bool(message.fixMsg);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseInstantiateContract2Proposal();
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
                    message.runAs = reader.string();
                    break;
                case 4:
                    message.admin = reader.string();
                    break;
                case 5:
                    message.codeId = reader.uint64();
                    break;
                case 6:
                    message.label = reader.string();
                    break;
                case 7:
                    message.msg = reader.bytes();
                    break;
                case 8:
                    message.funds.push(coin_1.Coin.decode(reader, reader.uint32()));
                    break;
                case 9:
                    message.salt = reader.bytes();
                    break;
                case 10:
                    message.fixMsg = reader.bool();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseInstantiateContract2Proposal();
        if ((0, helpers_1.isSet)(object.title))
            obj.title = String(object.title);
        if ((0, helpers_1.isSet)(object.description))
            obj.description = String(object.description);
        if ((0, helpers_1.isSet)(object.runAs))
            obj.runAs = String(object.runAs);
        if ((0, helpers_1.isSet)(object.admin))
            obj.admin = String(object.admin);
        if ((0, helpers_1.isSet)(object.codeId))
            obj.codeId = BigInt(object.codeId.toString());
        if ((0, helpers_1.isSet)(object.label))
            obj.label = String(object.label);
        if ((0, helpers_1.isSet)(object.msg))
            obj.msg = (0, helpers_1.bytesFromBase64)(object.msg);
        if (Array.isArray(object?.funds))
            obj.funds = object.funds.map((e) => coin_1.Coin.fromJSON(e));
        if ((0, helpers_1.isSet)(object.salt))
            obj.salt = (0, helpers_1.bytesFromBase64)(object.salt);
        if ((0, helpers_1.isSet)(object.fixMsg))
            obj.fixMsg = Boolean(object.fixMsg);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.title !== undefined && (obj.title = message.title);
        message.description !== undefined && (obj.description = message.description);
        message.runAs !== undefined && (obj.runAs = message.runAs);
        message.admin !== undefined && (obj.admin = message.admin);
        message.codeId !== undefined && (obj.codeId = (message.codeId || BigInt(0)).toString());
        message.label !== undefined && (obj.label = message.label);
        message.msg !== undefined &&
            (obj.msg = (0, helpers_1.base64FromBytes)(message.msg !== undefined ? message.msg : new Uint8Array()));
        if (message.funds) {
            obj.funds = message.funds.map((e) => (e ? coin_1.Coin.toJSON(e) : undefined));
        }
        else {
            obj.funds = [];
        }
        message.salt !== undefined &&
            (obj.salt = (0, helpers_1.base64FromBytes)(message.salt !== undefined ? message.salt : new Uint8Array()));
        message.fixMsg !== undefined && (obj.fixMsg = message.fixMsg);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseInstantiateContract2Proposal();
        message.title = object.title ?? "";
        message.description = object.description ?? "";
        message.runAs = object.runAs ?? "";
        message.admin = object.admin ?? "";
        if (object.codeId !== undefined && object.codeId !== null) {
            message.codeId = BigInt(object.codeId.toString());
        }
        message.label = object.label ?? "";
        message.msg = object.msg ?? new Uint8Array();
        message.funds = object.funds?.map((e) => coin_1.Coin.fromPartial(e)) || [];
        message.salt = object.salt ?? new Uint8Array();
        message.fixMsg = object.fixMsg ?? false;
        return message;
    },
};
function createBaseMigrateContractProposal() {
    return {
        title: "",
        description: "",
        contract: "",
        codeId: BigInt(0),
        msg: new Uint8Array(),
    };
}
exports.MigrateContractProposal = {
    typeUrl: "/cosmwasm.wasm.v1.MigrateContractProposal",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.title !== "") {
            writer.uint32(10).string(message.title);
        }
        if (message.description !== "") {
            writer.uint32(18).string(message.description);
        }
        if (message.contract !== "") {
            writer.uint32(34).string(message.contract);
        }
        if (message.codeId !== BigInt(0)) {
            writer.uint32(40).uint64(message.codeId);
        }
        if (message.msg.length !== 0) {
            writer.uint32(50).bytes(message.msg);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMigrateContractProposal();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.title = reader.string();
                    break;
                case 2:
                    message.description = reader.string();
                    break;
                case 4:
                    message.contract = reader.string();
                    break;
                case 5:
                    message.codeId = reader.uint64();
                    break;
                case 6:
                    message.msg = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMigrateContractProposal();
        if ((0, helpers_1.isSet)(object.title))
            obj.title = String(object.title);
        if ((0, helpers_1.isSet)(object.description))
            obj.description = String(object.description);
        if ((0, helpers_1.isSet)(object.contract))
            obj.contract = String(object.contract);
        if ((0, helpers_1.isSet)(object.codeId))
            obj.codeId = BigInt(object.codeId.toString());
        if ((0, helpers_1.isSet)(object.msg))
            obj.msg = (0, helpers_1.bytesFromBase64)(object.msg);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.title !== undefined && (obj.title = message.title);
        message.description !== undefined && (obj.description = message.description);
        message.contract !== undefined && (obj.contract = message.contract);
        message.codeId !== undefined && (obj.codeId = (message.codeId || BigInt(0)).toString());
        message.msg !== undefined &&
            (obj.msg = (0, helpers_1.base64FromBytes)(message.msg !== undefined ? message.msg : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMigrateContractProposal();
        message.title = object.title ?? "";
        message.description = object.description ?? "";
        message.contract = object.contract ?? "";
        if (object.codeId !== undefined && object.codeId !== null) {
            message.codeId = BigInt(object.codeId.toString());
        }
        message.msg = object.msg ?? new Uint8Array();
        return message;
    },
};
function createBaseSudoContractProposal() {
    return {
        title: "",
        description: "",
        contract: "",
        msg: new Uint8Array(),
    };
}
exports.SudoContractProposal = {
    typeUrl: "/cosmwasm.wasm.v1.SudoContractProposal",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.title !== "") {
            writer.uint32(10).string(message.title);
        }
        if (message.description !== "") {
            writer.uint32(18).string(message.description);
        }
        if (message.contract !== "") {
            writer.uint32(26).string(message.contract);
        }
        if (message.msg.length !== 0) {
            writer.uint32(34).bytes(message.msg);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSudoContractProposal();
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
                    message.contract = reader.string();
                    break;
                case 4:
                    message.msg = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseSudoContractProposal();
        if ((0, helpers_1.isSet)(object.title))
            obj.title = String(object.title);
        if ((0, helpers_1.isSet)(object.description))
            obj.description = String(object.description);
        if ((0, helpers_1.isSet)(object.contract))
            obj.contract = String(object.contract);
        if ((0, helpers_1.isSet)(object.msg))
            obj.msg = (0, helpers_1.bytesFromBase64)(object.msg);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.title !== undefined && (obj.title = message.title);
        message.description !== undefined && (obj.description = message.description);
        message.contract !== undefined && (obj.contract = message.contract);
        message.msg !== undefined &&
            (obj.msg = (0, helpers_1.base64FromBytes)(message.msg !== undefined ? message.msg : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseSudoContractProposal();
        message.title = object.title ?? "";
        message.description = object.description ?? "";
        message.contract = object.contract ?? "";
        message.msg = object.msg ?? new Uint8Array();
        return message;
    },
};
function createBaseExecuteContractProposal() {
    return {
        title: "",
        description: "",
        runAs: "",
        contract: "",
        msg: new Uint8Array(),
        funds: [],
    };
}
exports.ExecuteContractProposal = {
    typeUrl: "/cosmwasm.wasm.v1.ExecuteContractProposal",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.title !== "") {
            writer.uint32(10).string(message.title);
        }
        if (message.description !== "") {
            writer.uint32(18).string(message.description);
        }
        if (message.runAs !== "") {
            writer.uint32(26).string(message.runAs);
        }
        if (message.contract !== "") {
            writer.uint32(34).string(message.contract);
        }
        if (message.msg.length !== 0) {
            writer.uint32(42).bytes(message.msg);
        }
        for (const v of message.funds) {
            coin_1.Coin.encode(v, writer.uint32(50).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseExecuteContractProposal();
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
                    message.runAs = reader.string();
                    break;
                case 4:
                    message.contract = reader.string();
                    break;
                case 5:
                    message.msg = reader.bytes();
                    break;
                case 6:
                    message.funds.push(coin_1.Coin.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseExecuteContractProposal();
        if ((0, helpers_1.isSet)(object.title))
            obj.title = String(object.title);
        if ((0, helpers_1.isSet)(object.description))
            obj.description = String(object.description);
        if ((0, helpers_1.isSet)(object.runAs))
            obj.runAs = String(object.runAs);
        if ((0, helpers_1.isSet)(object.contract))
            obj.contract = String(object.contract);
        if ((0, helpers_1.isSet)(object.msg))
            obj.msg = (0, helpers_1.bytesFromBase64)(object.msg);
        if (Array.isArray(object?.funds))
            obj.funds = object.funds.map((e) => coin_1.Coin.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.title !== undefined && (obj.title = message.title);
        message.description !== undefined && (obj.description = message.description);
        message.runAs !== undefined && (obj.runAs = message.runAs);
        message.contract !== undefined && (obj.contract = message.contract);
        message.msg !== undefined &&
            (obj.msg = (0, helpers_1.base64FromBytes)(message.msg !== undefined ? message.msg : new Uint8Array()));
        if (message.funds) {
            obj.funds = message.funds.map((e) => (e ? coin_1.Coin.toJSON(e) : undefined));
        }
        else {
            obj.funds = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseExecuteContractProposal();
        message.title = object.title ?? "";
        message.description = object.description ?? "";
        message.runAs = object.runAs ?? "";
        message.contract = object.contract ?? "";
        message.msg = object.msg ?? new Uint8Array();
        message.funds = object.funds?.map((e) => coin_1.Coin.fromPartial(e)) || [];
        return message;
    },
};
function createBaseUpdateAdminProposal() {
    return {
        title: "",
        description: "",
        newAdmin: "",
        contract: "",
    };
}
exports.UpdateAdminProposal = {
    typeUrl: "/cosmwasm.wasm.v1.UpdateAdminProposal",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.title !== "") {
            writer.uint32(10).string(message.title);
        }
        if (message.description !== "") {
            writer.uint32(18).string(message.description);
        }
        if (message.newAdmin !== "") {
            writer.uint32(26).string(message.newAdmin);
        }
        if (message.contract !== "") {
            writer.uint32(34).string(message.contract);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseUpdateAdminProposal();
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
                    message.newAdmin = reader.string();
                    break;
                case 4:
                    message.contract = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseUpdateAdminProposal();
        if ((0, helpers_1.isSet)(object.title))
            obj.title = String(object.title);
        if ((0, helpers_1.isSet)(object.description))
            obj.description = String(object.description);
        if ((0, helpers_1.isSet)(object.newAdmin))
            obj.newAdmin = String(object.newAdmin);
        if ((0, helpers_1.isSet)(object.contract))
            obj.contract = String(object.contract);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.title !== undefined && (obj.title = message.title);
        message.description !== undefined && (obj.description = message.description);
        message.newAdmin !== undefined && (obj.newAdmin = message.newAdmin);
        message.contract !== undefined && (obj.contract = message.contract);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseUpdateAdminProposal();
        message.title = object.title ?? "";
        message.description = object.description ?? "";
        message.newAdmin = object.newAdmin ?? "";
        message.contract = object.contract ?? "";
        return message;
    },
};
function createBaseClearAdminProposal() {
    return {
        title: "",
        description: "",
        contract: "",
    };
}
exports.ClearAdminProposal = {
    typeUrl: "/cosmwasm.wasm.v1.ClearAdminProposal",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.title !== "") {
            writer.uint32(10).string(message.title);
        }
        if (message.description !== "") {
            writer.uint32(18).string(message.description);
        }
        if (message.contract !== "") {
            writer.uint32(26).string(message.contract);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseClearAdminProposal();
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
                    message.contract = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseClearAdminProposal();
        if ((0, helpers_1.isSet)(object.title))
            obj.title = String(object.title);
        if ((0, helpers_1.isSet)(object.description))
            obj.description = String(object.description);
        if ((0, helpers_1.isSet)(object.contract))
            obj.contract = String(object.contract);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.title !== undefined && (obj.title = message.title);
        message.description !== undefined && (obj.description = message.description);
        message.contract !== undefined && (obj.contract = message.contract);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseClearAdminProposal();
        message.title = object.title ?? "";
        message.description = object.description ?? "";
        message.contract = object.contract ?? "";
        return message;
    },
};
function createBasePinCodesProposal() {
    return {
        title: "",
        description: "",
        codeIds: [],
    };
}
exports.PinCodesProposal = {
    typeUrl: "/cosmwasm.wasm.v1.PinCodesProposal",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.title !== "") {
            writer.uint32(10).string(message.title);
        }
        if (message.description !== "") {
            writer.uint32(18).string(message.description);
        }
        writer.uint32(26).fork();
        for (const v of message.codeIds) {
            writer.uint64(v);
        }
        writer.ldelim();
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBasePinCodesProposal();
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
                    if ((tag & 7) === 2) {
                        const end2 = reader.uint32() + reader.pos;
                        while (reader.pos < end2) {
                            message.codeIds.push(reader.uint64());
                        }
                    }
                    else {
                        message.codeIds.push(reader.uint64());
                    }
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBasePinCodesProposal();
        if ((0, helpers_1.isSet)(object.title))
            obj.title = String(object.title);
        if ((0, helpers_1.isSet)(object.description))
            obj.description = String(object.description);
        if (Array.isArray(object?.codeIds))
            obj.codeIds = object.codeIds.map((e) => BigInt(e.toString()));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.title !== undefined && (obj.title = message.title);
        message.description !== undefined && (obj.description = message.description);
        if (message.codeIds) {
            obj.codeIds = message.codeIds.map((e) => (e || BigInt(0)).toString());
        }
        else {
            obj.codeIds = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBasePinCodesProposal();
        message.title = object.title ?? "";
        message.description = object.description ?? "";
        message.codeIds = object.codeIds?.map((e) => BigInt(e.toString())) || [];
        return message;
    },
};
function createBaseUnpinCodesProposal() {
    return {
        title: "",
        description: "",
        codeIds: [],
    };
}
exports.UnpinCodesProposal = {
    typeUrl: "/cosmwasm.wasm.v1.UnpinCodesProposal",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.title !== "") {
            writer.uint32(10).string(message.title);
        }
        if (message.description !== "") {
            writer.uint32(18).string(message.description);
        }
        writer.uint32(26).fork();
        for (const v of message.codeIds) {
            writer.uint64(v);
        }
        writer.ldelim();
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseUnpinCodesProposal();
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
                    if ((tag & 7) === 2) {
                        const end2 = reader.uint32() + reader.pos;
                        while (reader.pos < end2) {
                            message.codeIds.push(reader.uint64());
                        }
                    }
                    else {
                        message.codeIds.push(reader.uint64());
                    }
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseUnpinCodesProposal();
        if ((0, helpers_1.isSet)(object.title))
            obj.title = String(object.title);
        if ((0, helpers_1.isSet)(object.description))
            obj.description = String(object.description);
        if (Array.isArray(object?.codeIds))
            obj.codeIds = object.codeIds.map((e) => BigInt(e.toString()));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.title !== undefined && (obj.title = message.title);
        message.description !== undefined && (obj.description = message.description);
        if (message.codeIds) {
            obj.codeIds = message.codeIds.map((e) => (e || BigInt(0)).toString());
        }
        else {
            obj.codeIds = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseUnpinCodesProposal();
        message.title = object.title ?? "";
        message.description = object.description ?? "";
        message.codeIds = object.codeIds?.map((e) => BigInt(e.toString())) || [];
        return message;
    },
};
function createBaseAccessConfigUpdate() {
    return {
        codeId: BigInt(0),
        instantiatePermission: types_1.AccessConfig.fromPartial({}),
    };
}
exports.AccessConfigUpdate = {
    typeUrl: "/cosmwasm.wasm.v1.AccessConfigUpdate",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.codeId !== BigInt(0)) {
            writer.uint32(8).uint64(message.codeId);
        }
        if (message.instantiatePermission !== undefined) {
            types_1.AccessConfig.encode(message.instantiatePermission, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseAccessConfigUpdate();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.codeId = reader.uint64();
                    break;
                case 2:
                    message.instantiatePermission = types_1.AccessConfig.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseAccessConfigUpdate();
        if ((0, helpers_1.isSet)(object.codeId))
            obj.codeId = BigInt(object.codeId.toString());
        if ((0, helpers_1.isSet)(object.instantiatePermission))
            obj.instantiatePermission = types_1.AccessConfig.fromJSON(object.instantiatePermission);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.codeId !== undefined && (obj.codeId = (message.codeId || BigInt(0)).toString());
        message.instantiatePermission !== undefined &&
            (obj.instantiatePermission = message.instantiatePermission
                ? types_1.AccessConfig.toJSON(message.instantiatePermission)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseAccessConfigUpdate();
        if (object.codeId !== undefined && object.codeId !== null) {
            message.codeId = BigInt(object.codeId.toString());
        }
        if (object.instantiatePermission !== undefined && object.instantiatePermission !== null) {
            message.instantiatePermission = types_1.AccessConfig.fromPartial(object.instantiatePermission);
        }
        return message;
    },
};
function createBaseUpdateInstantiateConfigProposal() {
    return {
        title: "",
        description: "",
        accessConfigUpdates: [],
    };
}
exports.UpdateInstantiateConfigProposal = {
    typeUrl: "/cosmwasm.wasm.v1.UpdateInstantiateConfigProposal",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.title !== "") {
            writer.uint32(10).string(message.title);
        }
        if (message.description !== "") {
            writer.uint32(18).string(message.description);
        }
        for (const v of message.accessConfigUpdates) {
            exports.AccessConfigUpdate.encode(v, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseUpdateInstantiateConfigProposal();
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
                    message.accessConfigUpdates.push(exports.AccessConfigUpdate.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseUpdateInstantiateConfigProposal();
        if ((0, helpers_1.isSet)(object.title))
            obj.title = String(object.title);
        if ((0, helpers_1.isSet)(object.description))
            obj.description = String(object.description);
        if (Array.isArray(object?.accessConfigUpdates))
            obj.accessConfigUpdates = object.accessConfigUpdates.map((e) => exports.AccessConfigUpdate.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.title !== undefined && (obj.title = message.title);
        message.description !== undefined && (obj.description = message.description);
        if (message.accessConfigUpdates) {
            obj.accessConfigUpdates = message.accessConfigUpdates.map((e) => e ? exports.AccessConfigUpdate.toJSON(e) : undefined);
        }
        else {
            obj.accessConfigUpdates = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseUpdateInstantiateConfigProposal();
        message.title = object.title ?? "";
        message.description = object.description ?? "";
        message.accessConfigUpdates =
            object.accessConfigUpdates?.map((e) => exports.AccessConfigUpdate.fromPartial(e)) || [];
        return message;
    },
};
function createBaseStoreAndInstantiateContractProposal() {
    return {
        title: "",
        description: "",
        runAs: "",
        wasmByteCode: new Uint8Array(),
        instantiatePermission: undefined,
        unpinCode: false,
        admin: "",
        label: "",
        msg: new Uint8Array(),
        funds: [],
        source: "",
        builder: "",
        codeHash: new Uint8Array(),
    };
}
exports.StoreAndInstantiateContractProposal = {
    typeUrl: "/cosmwasm.wasm.v1.StoreAndInstantiateContractProposal",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.title !== "") {
            writer.uint32(10).string(message.title);
        }
        if (message.description !== "") {
            writer.uint32(18).string(message.description);
        }
        if (message.runAs !== "") {
            writer.uint32(26).string(message.runAs);
        }
        if (message.wasmByteCode.length !== 0) {
            writer.uint32(34).bytes(message.wasmByteCode);
        }
        if (message.instantiatePermission !== undefined) {
            types_1.AccessConfig.encode(message.instantiatePermission, writer.uint32(42).fork()).ldelim();
        }
        if (message.unpinCode === true) {
            writer.uint32(48).bool(message.unpinCode);
        }
        if (message.admin !== "") {
            writer.uint32(58).string(message.admin);
        }
        if (message.label !== "") {
            writer.uint32(66).string(message.label);
        }
        if (message.msg.length !== 0) {
            writer.uint32(74).bytes(message.msg);
        }
        for (const v of message.funds) {
            coin_1.Coin.encode(v, writer.uint32(82).fork()).ldelim();
        }
        if (message.source !== "") {
            writer.uint32(90).string(message.source);
        }
        if (message.builder !== "") {
            writer.uint32(98).string(message.builder);
        }
        if (message.codeHash.length !== 0) {
            writer.uint32(106).bytes(message.codeHash);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseStoreAndInstantiateContractProposal();
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
                    message.runAs = reader.string();
                    break;
                case 4:
                    message.wasmByteCode = reader.bytes();
                    break;
                case 5:
                    message.instantiatePermission = types_1.AccessConfig.decode(reader, reader.uint32());
                    break;
                case 6:
                    message.unpinCode = reader.bool();
                    break;
                case 7:
                    message.admin = reader.string();
                    break;
                case 8:
                    message.label = reader.string();
                    break;
                case 9:
                    message.msg = reader.bytes();
                    break;
                case 10:
                    message.funds.push(coin_1.Coin.decode(reader, reader.uint32()));
                    break;
                case 11:
                    message.source = reader.string();
                    break;
                case 12:
                    message.builder = reader.string();
                    break;
                case 13:
                    message.codeHash = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseStoreAndInstantiateContractProposal();
        if ((0, helpers_1.isSet)(object.title))
            obj.title = String(object.title);
        if ((0, helpers_1.isSet)(object.description))
            obj.description = String(object.description);
        if ((0, helpers_1.isSet)(object.runAs))
            obj.runAs = String(object.runAs);
        if ((0, helpers_1.isSet)(object.wasmByteCode))
            obj.wasmByteCode = (0, helpers_1.bytesFromBase64)(object.wasmByteCode);
        if ((0, helpers_1.isSet)(object.instantiatePermission))
            obj.instantiatePermission = types_1.AccessConfig.fromJSON(object.instantiatePermission);
        if ((0, helpers_1.isSet)(object.unpinCode))
            obj.unpinCode = Boolean(object.unpinCode);
        if ((0, helpers_1.isSet)(object.admin))
            obj.admin = String(object.admin);
        if ((0, helpers_1.isSet)(object.label))
            obj.label = String(object.label);
        if ((0, helpers_1.isSet)(object.msg))
            obj.msg = (0, helpers_1.bytesFromBase64)(object.msg);
        if (Array.isArray(object?.funds))
            obj.funds = object.funds.map((e) => coin_1.Coin.fromJSON(e));
        if ((0, helpers_1.isSet)(object.source))
            obj.source = String(object.source);
        if ((0, helpers_1.isSet)(object.builder))
            obj.builder = String(object.builder);
        if ((0, helpers_1.isSet)(object.codeHash))
            obj.codeHash = (0, helpers_1.bytesFromBase64)(object.codeHash);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.title !== undefined && (obj.title = message.title);
        message.description !== undefined && (obj.description = message.description);
        message.runAs !== undefined && (obj.runAs = message.runAs);
        message.wasmByteCode !== undefined &&
            (obj.wasmByteCode = (0, helpers_1.base64FromBytes)(message.wasmByteCode !== undefined ? message.wasmByteCode : new Uint8Array()));
        message.instantiatePermission !== undefined &&
            (obj.instantiatePermission = message.instantiatePermission
                ? types_1.AccessConfig.toJSON(message.instantiatePermission)
                : undefined);
        message.unpinCode !== undefined && (obj.unpinCode = message.unpinCode);
        message.admin !== undefined && (obj.admin = message.admin);
        message.label !== undefined && (obj.label = message.label);
        message.msg !== undefined &&
            (obj.msg = (0, helpers_1.base64FromBytes)(message.msg !== undefined ? message.msg : new Uint8Array()));
        if (message.funds) {
            obj.funds = message.funds.map((e) => (e ? coin_1.Coin.toJSON(e) : undefined));
        }
        else {
            obj.funds = [];
        }
        message.source !== undefined && (obj.source = message.source);
        message.builder !== undefined && (obj.builder = message.builder);
        message.codeHash !== undefined &&
            (obj.codeHash = (0, helpers_1.base64FromBytes)(message.codeHash !== undefined ? message.codeHash : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseStoreAndInstantiateContractProposal();
        message.title = object.title ?? "";
        message.description = object.description ?? "";
        message.runAs = object.runAs ?? "";
        message.wasmByteCode = object.wasmByteCode ?? new Uint8Array();
        if (object.instantiatePermission !== undefined && object.instantiatePermission !== null) {
            message.instantiatePermission = types_1.AccessConfig.fromPartial(object.instantiatePermission);
        }
        message.unpinCode = object.unpinCode ?? false;
        message.admin = object.admin ?? "";
        message.label = object.label ?? "";
        message.msg = object.msg ?? new Uint8Array();
        message.funds = object.funds?.map((e) => coin_1.Coin.fromPartial(e)) || [];
        message.source = object.source ?? "";
        message.builder = object.builder ?? "";
        message.codeHash = object.codeHash ?? new Uint8Array();
        return message;
    },
};
//# sourceMappingURL=proposal.js.map