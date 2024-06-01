"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.AcceptedMessagesFilter = exports.AcceptedMessageKeysFilter = exports.AllowAllMessagesFilter = exports.CombinedLimit = exports.MaxFundsLimit = exports.MaxCallsLimit = exports.ContractGrant = exports.ContractMigrationAuthorization = exports.ContractExecutionAuthorization = exports.protobufPackage = void 0;
/* eslint-disable */
const any_1 = require("../../../google/protobuf/any");
const coin_1 = require("../../../cosmos/base/v1beta1/coin");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmwasm.wasm.v1";
function createBaseContractExecutionAuthorization() {
    return {
        grants: [],
    };
}
exports.ContractExecutionAuthorization = {
    typeUrl: "/cosmwasm.wasm.v1.ContractExecutionAuthorization",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.grants) {
            exports.ContractGrant.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseContractExecutionAuthorization();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.grants.push(exports.ContractGrant.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseContractExecutionAuthorization();
        if (Array.isArray(object?.grants))
            obj.grants = object.grants.map((e) => exports.ContractGrant.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.grants) {
            obj.grants = message.grants.map((e) => (e ? exports.ContractGrant.toJSON(e) : undefined));
        }
        else {
            obj.grants = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseContractExecutionAuthorization();
        message.grants = object.grants?.map((e) => exports.ContractGrant.fromPartial(e)) || [];
        return message;
    },
};
function createBaseContractMigrationAuthorization() {
    return {
        grants: [],
    };
}
exports.ContractMigrationAuthorization = {
    typeUrl: "/cosmwasm.wasm.v1.ContractMigrationAuthorization",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.grants) {
            exports.ContractGrant.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseContractMigrationAuthorization();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.grants.push(exports.ContractGrant.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseContractMigrationAuthorization();
        if (Array.isArray(object?.grants))
            obj.grants = object.grants.map((e) => exports.ContractGrant.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.grants) {
            obj.grants = message.grants.map((e) => (e ? exports.ContractGrant.toJSON(e) : undefined));
        }
        else {
            obj.grants = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseContractMigrationAuthorization();
        message.grants = object.grants?.map((e) => exports.ContractGrant.fromPartial(e)) || [];
        return message;
    },
};
function createBaseContractGrant() {
    return {
        contract: "",
        limit: undefined,
        filter: undefined,
    };
}
exports.ContractGrant = {
    typeUrl: "/cosmwasm.wasm.v1.ContractGrant",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.contract !== "") {
            writer.uint32(10).string(message.contract);
        }
        if (message.limit !== undefined) {
            any_1.Any.encode(message.limit, writer.uint32(18).fork()).ldelim();
        }
        if (message.filter !== undefined) {
            any_1.Any.encode(message.filter, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseContractGrant();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.contract = reader.string();
                    break;
                case 2:
                    message.limit = any_1.Any.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.filter = any_1.Any.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseContractGrant();
        if ((0, helpers_1.isSet)(object.contract))
            obj.contract = String(object.contract);
        if ((0, helpers_1.isSet)(object.limit))
            obj.limit = any_1.Any.fromJSON(object.limit);
        if ((0, helpers_1.isSet)(object.filter))
            obj.filter = any_1.Any.fromJSON(object.filter);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.contract !== undefined && (obj.contract = message.contract);
        message.limit !== undefined && (obj.limit = message.limit ? any_1.Any.toJSON(message.limit) : undefined);
        message.filter !== undefined && (obj.filter = message.filter ? any_1.Any.toJSON(message.filter) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseContractGrant();
        message.contract = object.contract ?? "";
        if (object.limit !== undefined && object.limit !== null) {
            message.limit = any_1.Any.fromPartial(object.limit);
        }
        if (object.filter !== undefined && object.filter !== null) {
            message.filter = any_1.Any.fromPartial(object.filter);
        }
        return message;
    },
};
function createBaseMaxCallsLimit() {
    return {
        remaining: BigInt(0),
    };
}
exports.MaxCallsLimit = {
    typeUrl: "/cosmwasm.wasm.v1.MaxCallsLimit",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.remaining !== BigInt(0)) {
            writer.uint32(8).uint64(message.remaining);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMaxCallsLimit();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.remaining = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMaxCallsLimit();
        if ((0, helpers_1.isSet)(object.remaining))
            obj.remaining = BigInt(object.remaining.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.remaining !== undefined && (obj.remaining = (message.remaining || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMaxCallsLimit();
        if (object.remaining !== undefined && object.remaining !== null) {
            message.remaining = BigInt(object.remaining.toString());
        }
        return message;
    },
};
function createBaseMaxFundsLimit() {
    return {
        amounts: [],
    };
}
exports.MaxFundsLimit = {
    typeUrl: "/cosmwasm.wasm.v1.MaxFundsLimit",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.amounts) {
            coin_1.Coin.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMaxFundsLimit();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.amounts.push(coin_1.Coin.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMaxFundsLimit();
        if (Array.isArray(object?.amounts))
            obj.amounts = object.amounts.map((e) => coin_1.Coin.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.amounts) {
            obj.amounts = message.amounts.map((e) => (e ? coin_1.Coin.toJSON(e) : undefined));
        }
        else {
            obj.amounts = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMaxFundsLimit();
        message.amounts = object.amounts?.map((e) => coin_1.Coin.fromPartial(e)) || [];
        return message;
    },
};
function createBaseCombinedLimit() {
    return {
        callsRemaining: BigInt(0),
        amounts: [],
    };
}
exports.CombinedLimit = {
    typeUrl: "/cosmwasm.wasm.v1.CombinedLimit",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.callsRemaining !== BigInt(0)) {
            writer.uint32(8).uint64(message.callsRemaining);
        }
        for (const v of message.amounts) {
            coin_1.Coin.encode(v, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseCombinedLimit();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.callsRemaining = reader.uint64();
                    break;
                case 2:
                    message.amounts.push(coin_1.Coin.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseCombinedLimit();
        if ((0, helpers_1.isSet)(object.callsRemaining))
            obj.callsRemaining = BigInt(object.callsRemaining.toString());
        if (Array.isArray(object?.amounts))
            obj.amounts = object.amounts.map((e) => coin_1.Coin.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.callsRemaining !== undefined &&
            (obj.callsRemaining = (message.callsRemaining || BigInt(0)).toString());
        if (message.amounts) {
            obj.amounts = message.amounts.map((e) => (e ? coin_1.Coin.toJSON(e) : undefined));
        }
        else {
            obj.amounts = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseCombinedLimit();
        if (object.callsRemaining !== undefined && object.callsRemaining !== null) {
            message.callsRemaining = BigInt(object.callsRemaining.toString());
        }
        message.amounts = object.amounts?.map((e) => coin_1.Coin.fromPartial(e)) || [];
        return message;
    },
};
function createBaseAllowAllMessagesFilter() {
    return {};
}
exports.AllowAllMessagesFilter = {
    typeUrl: "/cosmwasm.wasm.v1.AllowAllMessagesFilter",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseAllowAllMessagesFilter();
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
        const obj = createBaseAllowAllMessagesFilter();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseAllowAllMessagesFilter();
        return message;
    },
};
function createBaseAcceptedMessageKeysFilter() {
    return {
        keys: [],
    };
}
exports.AcceptedMessageKeysFilter = {
    typeUrl: "/cosmwasm.wasm.v1.AcceptedMessageKeysFilter",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.keys) {
            writer.uint32(10).string(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseAcceptedMessageKeysFilter();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.keys.push(reader.string());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseAcceptedMessageKeysFilter();
        if (Array.isArray(object?.keys))
            obj.keys = object.keys.map((e) => String(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.keys) {
            obj.keys = message.keys.map((e) => e);
        }
        else {
            obj.keys = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseAcceptedMessageKeysFilter();
        message.keys = object.keys?.map((e) => e) || [];
        return message;
    },
};
function createBaseAcceptedMessagesFilter() {
    return {
        messages: [],
    };
}
exports.AcceptedMessagesFilter = {
    typeUrl: "/cosmwasm.wasm.v1.AcceptedMessagesFilter",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.messages) {
            writer.uint32(10).bytes(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseAcceptedMessagesFilter();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.messages.push(reader.bytes());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseAcceptedMessagesFilter();
        if (Array.isArray(object?.messages))
            obj.messages = object.messages.map((e) => (0, helpers_1.bytesFromBase64)(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.messages) {
            obj.messages = message.messages.map((e) => (0, helpers_1.base64FromBytes)(e !== undefined ? e : new Uint8Array()));
        }
        else {
            obj.messages = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseAcceptedMessagesFilter();
        message.messages = object.messages?.map((e) => e) || [];
        return message;
    },
};
//# sourceMappingURL=authz.js.map