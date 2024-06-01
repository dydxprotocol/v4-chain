"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ModuleVersion = exports.CancelSoftwareUpgradeProposal = exports.SoftwareUpgradeProposal = exports.Plan = exports.protobufPackage = void 0;
/* eslint-disable */
const timestamp_1 = require("../../../google/protobuf/timestamp");
const any_1 = require("../../../google/protobuf/any");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.upgrade.v1beta1";
function createBasePlan() {
    return {
        name: "",
        time: timestamp_1.Timestamp.fromPartial({}),
        height: BigInt(0),
        info: "",
        upgradedClientState: undefined,
    };
}
exports.Plan = {
    typeUrl: "/cosmos.upgrade.v1beta1.Plan",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.name !== "") {
            writer.uint32(10).string(message.name);
        }
        if (message.time !== undefined) {
            timestamp_1.Timestamp.encode(message.time, writer.uint32(18).fork()).ldelim();
        }
        if (message.height !== BigInt(0)) {
            writer.uint32(24).int64(message.height);
        }
        if (message.info !== "") {
            writer.uint32(34).string(message.info);
        }
        if (message.upgradedClientState !== undefined) {
            any_1.Any.encode(message.upgradedClientState, writer.uint32(42).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBasePlan();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.name = reader.string();
                    break;
                case 2:
                    message.time = timestamp_1.Timestamp.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.height = reader.int64();
                    break;
                case 4:
                    message.info = reader.string();
                    break;
                case 5:
                    message.upgradedClientState = any_1.Any.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBasePlan();
        if ((0, helpers_1.isSet)(object.name))
            obj.name = String(object.name);
        if ((0, helpers_1.isSet)(object.time))
            obj.time = (0, helpers_1.fromJsonTimestamp)(object.time);
        if ((0, helpers_1.isSet)(object.height))
            obj.height = BigInt(object.height.toString());
        if ((0, helpers_1.isSet)(object.info))
            obj.info = String(object.info);
        if ((0, helpers_1.isSet)(object.upgradedClientState))
            obj.upgradedClientState = any_1.Any.fromJSON(object.upgradedClientState);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.name !== undefined && (obj.name = message.name);
        message.time !== undefined && (obj.time = (0, helpers_1.fromTimestamp)(message.time).toISOString());
        message.height !== undefined && (obj.height = (message.height || BigInt(0)).toString());
        message.info !== undefined && (obj.info = message.info);
        message.upgradedClientState !== undefined &&
            (obj.upgradedClientState = message.upgradedClientState
                ? any_1.Any.toJSON(message.upgradedClientState)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBasePlan();
        message.name = object.name ?? "";
        if (object.time !== undefined && object.time !== null) {
            message.time = timestamp_1.Timestamp.fromPartial(object.time);
        }
        if (object.height !== undefined && object.height !== null) {
            message.height = BigInt(object.height.toString());
        }
        message.info = object.info ?? "";
        if (object.upgradedClientState !== undefined && object.upgradedClientState !== null) {
            message.upgradedClientState = any_1.Any.fromPartial(object.upgradedClientState);
        }
        return message;
    },
};
function createBaseSoftwareUpgradeProposal() {
    return {
        title: "",
        description: "",
        plan: exports.Plan.fromPartial({}),
    };
}
exports.SoftwareUpgradeProposal = {
    typeUrl: "/cosmos.upgrade.v1beta1.SoftwareUpgradeProposal",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.title !== "") {
            writer.uint32(10).string(message.title);
        }
        if (message.description !== "") {
            writer.uint32(18).string(message.description);
        }
        if (message.plan !== undefined) {
            exports.Plan.encode(message.plan, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSoftwareUpgradeProposal();
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
                    message.plan = exports.Plan.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseSoftwareUpgradeProposal();
        if ((0, helpers_1.isSet)(object.title))
            obj.title = String(object.title);
        if ((0, helpers_1.isSet)(object.description))
            obj.description = String(object.description);
        if ((0, helpers_1.isSet)(object.plan))
            obj.plan = exports.Plan.fromJSON(object.plan);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.title !== undefined && (obj.title = message.title);
        message.description !== undefined && (obj.description = message.description);
        message.plan !== undefined && (obj.plan = message.plan ? exports.Plan.toJSON(message.plan) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseSoftwareUpgradeProposal();
        message.title = object.title ?? "";
        message.description = object.description ?? "";
        if (object.plan !== undefined && object.plan !== null) {
            message.plan = exports.Plan.fromPartial(object.plan);
        }
        return message;
    },
};
function createBaseCancelSoftwareUpgradeProposal() {
    return {
        title: "",
        description: "",
    };
}
exports.CancelSoftwareUpgradeProposal = {
    typeUrl: "/cosmos.upgrade.v1beta1.CancelSoftwareUpgradeProposal",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.title !== "") {
            writer.uint32(10).string(message.title);
        }
        if (message.description !== "") {
            writer.uint32(18).string(message.description);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseCancelSoftwareUpgradeProposal();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.title = reader.string();
                    break;
                case 2:
                    message.description = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseCancelSoftwareUpgradeProposal();
        if ((0, helpers_1.isSet)(object.title))
            obj.title = String(object.title);
        if ((0, helpers_1.isSet)(object.description))
            obj.description = String(object.description);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.title !== undefined && (obj.title = message.title);
        message.description !== undefined && (obj.description = message.description);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseCancelSoftwareUpgradeProposal();
        message.title = object.title ?? "";
        message.description = object.description ?? "";
        return message;
    },
};
function createBaseModuleVersion() {
    return {
        name: "",
        version: BigInt(0),
    };
}
exports.ModuleVersion = {
    typeUrl: "/cosmos.upgrade.v1beta1.ModuleVersion",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.name !== "") {
            writer.uint32(10).string(message.name);
        }
        if (message.version !== BigInt(0)) {
            writer.uint32(16).uint64(message.version);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseModuleVersion();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.name = reader.string();
                    break;
                case 2:
                    message.version = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseModuleVersion();
        if ((0, helpers_1.isSet)(object.name))
            obj.name = String(object.name);
        if ((0, helpers_1.isSet)(object.version))
            obj.version = BigInt(object.version.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.name !== undefined && (obj.name = message.name);
        message.version !== undefined && (obj.version = (message.version || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseModuleVersion();
        message.name = object.name ?? "";
        if (object.version !== undefined && object.version !== null) {
            message.version = BigInt(object.version.toString());
        }
        return message;
    },
};
//# sourceMappingURL=upgrade.js.map