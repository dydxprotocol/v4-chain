"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Sequence = exports.Contract = exports.Code = exports.GenesisState = exports.protobufPackage = void 0;
/* eslint-disable */
const types_1 = require("./types");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmwasm.wasm.v1";
function createBaseGenesisState() {
    return {
        params: types_1.Params.fromPartial({}),
        codes: [],
        contracts: [],
        sequences: [],
    };
}
exports.GenesisState = {
    typeUrl: "/cosmwasm.wasm.v1.GenesisState",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.params !== undefined) {
            types_1.Params.encode(message.params, writer.uint32(10).fork()).ldelim();
        }
        for (const v of message.codes) {
            exports.Code.encode(v, writer.uint32(18).fork()).ldelim();
        }
        for (const v of message.contracts) {
            exports.Contract.encode(v, writer.uint32(26).fork()).ldelim();
        }
        for (const v of message.sequences) {
            exports.Sequence.encode(v, writer.uint32(34).fork()).ldelim();
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
                    message.params = types_1.Params.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.codes.push(exports.Code.decode(reader, reader.uint32()));
                    break;
                case 3:
                    message.contracts.push(exports.Contract.decode(reader, reader.uint32()));
                    break;
                case 4:
                    message.sequences.push(exports.Sequence.decode(reader, reader.uint32()));
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
        if ((0, helpers_1.isSet)(object.params))
            obj.params = types_1.Params.fromJSON(object.params);
        if (Array.isArray(object?.codes))
            obj.codes = object.codes.map((e) => exports.Code.fromJSON(e));
        if (Array.isArray(object?.contracts))
            obj.contracts = object.contracts.map((e) => exports.Contract.fromJSON(e));
        if (Array.isArray(object?.sequences))
            obj.sequences = object.sequences.map((e) => exports.Sequence.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.params !== undefined && (obj.params = message.params ? types_1.Params.toJSON(message.params) : undefined);
        if (message.codes) {
            obj.codes = message.codes.map((e) => (e ? exports.Code.toJSON(e) : undefined));
        }
        else {
            obj.codes = [];
        }
        if (message.contracts) {
            obj.contracts = message.contracts.map((e) => (e ? exports.Contract.toJSON(e) : undefined));
        }
        else {
            obj.contracts = [];
        }
        if (message.sequences) {
            obj.sequences = message.sequences.map((e) => (e ? exports.Sequence.toJSON(e) : undefined));
        }
        else {
            obj.sequences = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGenesisState();
        if (object.params !== undefined && object.params !== null) {
            message.params = types_1.Params.fromPartial(object.params);
        }
        message.codes = object.codes?.map((e) => exports.Code.fromPartial(e)) || [];
        message.contracts = object.contracts?.map((e) => exports.Contract.fromPartial(e)) || [];
        message.sequences = object.sequences?.map((e) => exports.Sequence.fromPartial(e)) || [];
        return message;
    },
};
function createBaseCode() {
    return {
        codeId: BigInt(0),
        codeInfo: types_1.CodeInfo.fromPartial({}),
        codeBytes: new Uint8Array(),
        pinned: false,
    };
}
exports.Code = {
    typeUrl: "/cosmwasm.wasm.v1.Code",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.codeId !== BigInt(0)) {
            writer.uint32(8).uint64(message.codeId);
        }
        if (message.codeInfo !== undefined) {
            types_1.CodeInfo.encode(message.codeInfo, writer.uint32(18).fork()).ldelim();
        }
        if (message.codeBytes.length !== 0) {
            writer.uint32(26).bytes(message.codeBytes);
        }
        if (message.pinned === true) {
            writer.uint32(32).bool(message.pinned);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseCode();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.codeId = reader.uint64();
                    break;
                case 2:
                    message.codeInfo = types_1.CodeInfo.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.codeBytes = reader.bytes();
                    break;
                case 4:
                    message.pinned = reader.bool();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseCode();
        if ((0, helpers_1.isSet)(object.codeId))
            obj.codeId = BigInt(object.codeId.toString());
        if ((0, helpers_1.isSet)(object.codeInfo))
            obj.codeInfo = types_1.CodeInfo.fromJSON(object.codeInfo);
        if ((0, helpers_1.isSet)(object.codeBytes))
            obj.codeBytes = (0, helpers_1.bytesFromBase64)(object.codeBytes);
        if ((0, helpers_1.isSet)(object.pinned))
            obj.pinned = Boolean(object.pinned);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.codeId !== undefined && (obj.codeId = (message.codeId || BigInt(0)).toString());
        message.codeInfo !== undefined &&
            (obj.codeInfo = message.codeInfo ? types_1.CodeInfo.toJSON(message.codeInfo) : undefined);
        message.codeBytes !== undefined &&
            (obj.codeBytes = (0, helpers_1.base64FromBytes)(message.codeBytes !== undefined ? message.codeBytes : new Uint8Array()));
        message.pinned !== undefined && (obj.pinned = message.pinned);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseCode();
        if (object.codeId !== undefined && object.codeId !== null) {
            message.codeId = BigInt(object.codeId.toString());
        }
        if (object.codeInfo !== undefined && object.codeInfo !== null) {
            message.codeInfo = types_1.CodeInfo.fromPartial(object.codeInfo);
        }
        message.codeBytes = object.codeBytes ?? new Uint8Array();
        message.pinned = object.pinned ?? false;
        return message;
    },
};
function createBaseContract() {
    return {
        contractAddress: "",
        contractInfo: types_1.ContractInfo.fromPartial({}),
        contractState: [],
        contractCodeHistory: [],
    };
}
exports.Contract = {
    typeUrl: "/cosmwasm.wasm.v1.Contract",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.contractAddress !== "") {
            writer.uint32(10).string(message.contractAddress);
        }
        if (message.contractInfo !== undefined) {
            types_1.ContractInfo.encode(message.contractInfo, writer.uint32(18).fork()).ldelim();
        }
        for (const v of message.contractState) {
            types_1.Model.encode(v, writer.uint32(26).fork()).ldelim();
        }
        for (const v of message.contractCodeHistory) {
            types_1.ContractCodeHistoryEntry.encode(v, writer.uint32(34).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseContract();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.contractAddress = reader.string();
                    break;
                case 2:
                    message.contractInfo = types_1.ContractInfo.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.contractState.push(types_1.Model.decode(reader, reader.uint32()));
                    break;
                case 4:
                    message.contractCodeHistory.push(types_1.ContractCodeHistoryEntry.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseContract();
        if ((0, helpers_1.isSet)(object.contractAddress))
            obj.contractAddress = String(object.contractAddress);
        if ((0, helpers_1.isSet)(object.contractInfo))
            obj.contractInfo = types_1.ContractInfo.fromJSON(object.contractInfo);
        if (Array.isArray(object?.contractState))
            obj.contractState = object.contractState.map((e) => types_1.Model.fromJSON(e));
        if (Array.isArray(object?.contractCodeHistory))
            obj.contractCodeHistory = object.contractCodeHistory.map((e) => types_1.ContractCodeHistoryEntry.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.contractAddress !== undefined && (obj.contractAddress = message.contractAddress);
        message.contractInfo !== undefined &&
            (obj.contractInfo = message.contractInfo ? types_1.ContractInfo.toJSON(message.contractInfo) : undefined);
        if (message.contractState) {
            obj.contractState = message.contractState.map((e) => (e ? types_1.Model.toJSON(e) : undefined));
        }
        else {
            obj.contractState = [];
        }
        if (message.contractCodeHistory) {
            obj.contractCodeHistory = message.contractCodeHistory.map((e) => e ? types_1.ContractCodeHistoryEntry.toJSON(e) : undefined);
        }
        else {
            obj.contractCodeHistory = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseContract();
        message.contractAddress = object.contractAddress ?? "";
        if (object.contractInfo !== undefined && object.contractInfo !== null) {
            message.contractInfo = types_1.ContractInfo.fromPartial(object.contractInfo);
        }
        message.contractState = object.contractState?.map((e) => types_1.Model.fromPartial(e)) || [];
        message.contractCodeHistory =
            object.contractCodeHistory?.map((e) => types_1.ContractCodeHistoryEntry.fromPartial(e)) || [];
        return message;
    },
};
function createBaseSequence() {
    return {
        idKey: new Uint8Array(),
        value: BigInt(0),
    };
}
exports.Sequence = {
    typeUrl: "/cosmwasm.wasm.v1.Sequence",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.idKey.length !== 0) {
            writer.uint32(10).bytes(message.idKey);
        }
        if (message.value !== BigInt(0)) {
            writer.uint32(16).uint64(message.value);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSequence();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.idKey = reader.bytes();
                    break;
                case 2:
                    message.value = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseSequence();
        if ((0, helpers_1.isSet)(object.idKey))
            obj.idKey = (0, helpers_1.bytesFromBase64)(object.idKey);
        if ((0, helpers_1.isSet)(object.value))
            obj.value = BigInt(object.value.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.idKey !== undefined &&
            (obj.idKey = (0, helpers_1.base64FromBytes)(message.idKey !== undefined ? message.idKey : new Uint8Array()));
        message.value !== undefined && (obj.value = (message.value || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseSequence();
        message.idKey = object.idKey ?? new Uint8Array();
        if (object.value !== undefined && object.value !== null) {
            message.value = BigInt(object.value.toString());
        }
        return message;
    },
};
//# sourceMappingURL=genesis.js.map