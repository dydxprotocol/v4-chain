"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.MsgClientImpl = exports.MsgCancelUpgradeResponse = exports.MsgCancelUpgrade = exports.MsgSoftwareUpgradeResponse = exports.MsgSoftwareUpgrade = exports.protobufPackage = void 0;
/* eslint-disable */
const upgrade_1 = require("./upgrade");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.upgrade.v1beta1";
function createBaseMsgSoftwareUpgrade() {
    return {
        authority: "",
        plan: upgrade_1.Plan.fromPartial({}),
    };
}
exports.MsgSoftwareUpgrade = {
    typeUrl: "/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.authority !== "") {
            writer.uint32(10).string(message.authority);
        }
        if (message.plan !== undefined) {
            upgrade_1.Plan.encode(message.plan, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgSoftwareUpgrade();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.authority = reader.string();
                    break;
                case 2:
                    message.plan = upgrade_1.Plan.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgSoftwareUpgrade();
        if ((0, helpers_1.isSet)(object.authority))
            obj.authority = String(object.authority);
        if ((0, helpers_1.isSet)(object.plan))
            obj.plan = upgrade_1.Plan.fromJSON(object.plan);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.authority !== undefined && (obj.authority = message.authority);
        message.plan !== undefined && (obj.plan = message.plan ? upgrade_1.Plan.toJSON(message.plan) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgSoftwareUpgrade();
        message.authority = object.authority ?? "";
        if (object.plan !== undefined && object.plan !== null) {
            message.plan = upgrade_1.Plan.fromPartial(object.plan);
        }
        return message;
    },
};
function createBaseMsgSoftwareUpgradeResponse() {
    return {};
}
exports.MsgSoftwareUpgradeResponse = {
    typeUrl: "/cosmos.upgrade.v1beta1.MsgSoftwareUpgradeResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgSoftwareUpgradeResponse();
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
        const obj = createBaseMsgSoftwareUpgradeResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgSoftwareUpgradeResponse();
        return message;
    },
};
function createBaseMsgCancelUpgrade() {
    return {
        authority: "",
    };
}
exports.MsgCancelUpgrade = {
    typeUrl: "/cosmos.upgrade.v1beta1.MsgCancelUpgrade",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.authority !== "") {
            writer.uint32(10).string(message.authority);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgCancelUpgrade();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
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
        const obj = createBaseMsgCancelUpgrade();
        if ((0, helpers_1.isSet)(object.authority))
            obj.authority = String(object.authority);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.authority !== undefined && (obj.authority = message.authority);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgCancelUpgrade();
        message.authority = object.authority ?? "";
        return message;
    },
};
function createBaseMsgCancelUpgradeResponse() {
    return {};
}
exports.MsgCancelUpgradeResponse = {
    typeUrl: "/cosmos.upgrade.v1beta1.MsgCancelUpgradeResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgCancelUpgradeResponse();
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
        const obj = createBaseMsgCancelUpgradeResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgCancelUpgradeResponse();
        return message;
    },
};
class MsgClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
        this.SoftwareUpgrade = this.SoftwareUpgrade.bind(this);
        this.CancelUpgrade = this.CancelUpgrade.bind(this);
    }
    SoftwareUpgrade(request) {
        const data = exports.MsgSoftwareUpgrade.encode(request).finish();
        const promise = this.rpc.request("cosmos.upgrade.v1beta1.Msg", "SoftwareUpgrade", data);
        return promise.then((data) => exports.MsgSoftwareUpgradeResponse.decode(new binary_1.BinaryReader(data)));
    }
    CancelUpgrade(request) {
        const data = exports.MsgCancelUpgrade.encode(request).finish();
        const promise = this.rpc.request("cosmos.upgrade.v1beta1.Msg", "CancelUpgrade", data);
        return promise.then((data) => exports.MsgCancelUpgradeResponse.decode(new binary_1.BinaryReader(data)));
    }
}
exports.MsgClientImpl = MsgClientImpl;
//# sourceMappingURL=tx.js.map