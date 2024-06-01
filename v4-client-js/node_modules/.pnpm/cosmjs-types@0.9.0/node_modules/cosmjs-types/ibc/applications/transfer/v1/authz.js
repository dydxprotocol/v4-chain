"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.TransferAuthorization = exports.Allocation = exports.protobufPackage = void 0;
/* eslint-disable */
const coin_1 = require("../../../../cosmos/base/v1beta1/coin");
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "ibc.applications.transfer.v1";
function createBaseAllocation() {
    return {
        sourcePort: "",
        sourceChannel: "",
        spendLimit: [],
        allowList: [],
    };
}
exports.Allocation = {
    typeUrl: "/ibc.applications.transfer.v1.Allocation",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.sourcePort !== "") {
            writer.uint32(10).string(message.sourcePort);
        }
        if (message.sourceChannel !== "") {
            writer.uint32(18).string(message.sourceChannel);
        }
        for (const v of message.spendLimit) {
            coin_1.Coin.encode(v, writer.uint32(26).fork()).ldelim();
        }
        for (const v of message.allowList) {
            writer.uint32(34).string(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseAllocation();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.sourcePort = reader.string();
                    break;
                case 2:
                    message.sourceChannel = reader.string();
                    break;
                case 3:
                    message.spendLimit.push(coin_1.Coin.decode(reader, reader.uint32()));
                    break;
                case 4:
                    message.allowList.push(reader.string());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseAllocation();
        if ((0, helpers_1.isSet)(object.sourcePort))
            obj.sourcePort = String(object.sourcePort);
        if ((0, helpers_1.isSet)(object.sourceChannel))
            obj.sourceChannel = String(object.sourceChannel);
        if (Array.isArray(object?.spendLimit))
            obj.spendLimit = object.spendLimit.map((e) => coin_1.Coin.fromJSON(e));
        if (Array.isArray(object?.allowList))
            obj.allowList = object.allowList.map((e) => String(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.sourcePort !== undefined && (obj.sourcePort = message.sourcePort);
        message.sourceChannel !== undefined && (obj.sourceChannel = message.sourceChannel);
        if (message.spendLimit) {
            obj.spendLimit = message.spendLimit.map((e) => (e ? coin_1.Coin.toJSON(e) : undefined));
        }
        else {
            obj.spendLimit = [];
        }
        if (message.allowList) {
            obj.allowList = message.allowList.map((e) => e);
        }
        else {
            obj.allowList = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseAllocation();
        message.sourcePort = object.sourcePort ?? "";
        message.sourceChannel = object.sourceChannel ?? "";
        message.spendLimit = object.spendLimit?.map((e) => coin_1.Coin.fromPartial(e)) || [];
        message.allowList = object.allowList?.map((e) => e) || [];
        return message;
    },
};
function createBaseTransferAuthorization() {
    return {
        allocations: [],
    };
}
exports.TransferAuthorization = {
    typeUrl: "/ibc.applications.transfer.v1.TransferAuthorization",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.allocations) {
            exports.Allocation.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseTransferAuthorization();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.allocations.push(exports.Allocation.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseTransferAuthorization();
        if (Array.isArray(object?.allocations))
            obj.allocations = object.allocations.map((e) => exports.Allocation.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.allocations) {
            obj.allocations = message.allocations.map((e) => (e ? exports.Allocation.toJSON(e) : undefined));
        }
        else {
            obj.allocations = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseTransferAuthorization();
        message.allocations = object.allocations?.map((e) => exports.Allocation.fromPartial(e)) || [];
        return message;
    },
};
//# sourceMappingURL=authz.js.map