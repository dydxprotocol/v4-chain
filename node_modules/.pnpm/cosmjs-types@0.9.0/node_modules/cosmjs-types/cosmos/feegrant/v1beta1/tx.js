"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.MsgClientImpl = exports.MsgRevokeAllowanceResponse = exports.MsgRevokeAllowance = exports.MsgGrantAllowanceResponse = exports.MsgGrantAllowance = exports.protobufPackage = void 0;
/* eslint-disable */
const any_1 = require("../../../google/protobuf/any");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.feegrant.v1beta1";
function createBaseMsgGrantAllowance() {
    return {
        granter: "",
        grantee: "",
        allowance: undefined,
    };
}
exports.MsgGrantAllowance = {
    typeUrl: "/cosmos.feegrant.v1beta1.MsgGrantAllowance",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.granter !== "") {
            writer.uint32(10).string(message.granter);
        }
        if (message.grantee !== "") {
            writer.uint32(18).string(message.grantee);
        }
        if (message.allowance !== undefined) {
            any_1.Any.encode(message.allowance, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgGrantAllowance();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.granter = reader.string();
                    break;
                case 2:
                    message.grantee = reader.string();
                    break;
                case 3:
                    message.allowance = any_1.Any.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgGrantAllowance();
        if ((0, helpers_1.isSet)(object.granter))
            obj.granter = String(object.granter);
        if ((0, helpers_1.isSet)(object.grantee))
            obj.grantee = String(object.grantee);
        if ((0, helpers_1.isSet)(object.allowance))
            obj.allowance = any_1.Any.fromJSON(object.allowance);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.granter !== undefined && (obj.granter = message.granter);
        message.grantee !== undefined && (obj.grantee = message.grantee);
        message.allowance !== undefined &&
            (obj.allowance = message.allowance ? any_1.Any.toJSON(message.allowance) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgGrantAllowance();
        message.granter = object.granter ?? "";
        message.grantee = object.grantee ?? "";
        if (object.allowance !== undefined && object.allowance !== null) {
            message.allowance = any_1.Any.fromPartial(object.allowance);
        }
        return message;
    },
};
function createBaseMsgGrantAllowanceResponse() {
    return {};
}
exports.MsgGrantAllowanceResponse = {
    typeUrl: "/cosmos.feegrant.v1beta1.MsgGrantAllowanceResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgGrantAllowanceResponse();
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
        const obj = createBaseMsgGrantAllowanceResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgGrantAllowanceResponse();
        return message;
    },
};
function createBaseMsgRevokeAllowance() {
    return {
        granter: "",
        grantee: "",
    };
}
exports.MsgRevokeAllowance = {
    typeUrl: "/cosmos.feegrant.v1beta1.MsgRevokeAllowance",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.granter !== "") {
            writer.uint32(10).string(message.granter);
        }
        if (message.grantee !== "") {
            writer.uint32(18).string(message.grantee);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgRevokeAllowance();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.granter = reader.string();
                    break;
                case 2:
                    message.grantee = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgRevokeAllowance();
        if ((0, helpers_1.isSet)(object.granter))
            obj.granter = String(object.granter);
        if ((0, helpers_1.isSet)(object.grantee))
            obj.grantee = String(object.grantee);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.granter !== undefined && (obj.granter = message.granter);
        message.grantee !== undefined && (obj.grantee = message.grantee);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgRevokeAllowance();
        message.granter = object.granter ?? "";
        message.grantee = object.grantee ?? "";
        return message;
    },
};
function createBaseMsgRevokeAllowanceResponse() {
    return {};
}
exports.MsgRevokeAllowanceResponse = {
    typeUrl: "/cosmos.feegrant.v1beta1.MsgRevokeAllowanceResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgRevokeAllowanceResponse();
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
        const obj = createBaseMsgRevokeAllowanceResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgRevokeAllowanceResponse();
        return message;
    },
};
class MsgClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
        this.GrantAllowance = this.GrantAllowance.bind(this);
        this.RevokeAllowance = this.RevokeAllowance.bind(this);
    }
    GrantAllowance(request) {
        const data = exports.MsgGrantAllowance.encode(request).finish();
        const promise = this.rpc.request("cosmos.feegrant.v1beta1.Msg", "GrantAllowance", data);
        return promise.then((data) => exports.MsgGrantAllowanceResponse.decode(new binary_1.BinaryReader(data)));
    }
    RevokeAllowance(request) {
        const data = exports.MsgRevokeAllowance.encode(request).finish();
        const promise = this.rpc.request("cosmos.feegrant.v1beta1.Msg", "RevokeAllowance", data);
        return promise.then((data) => exports.MsgRevokeAllowanceResponse.decode(new binary_1.BinaryReader(data)));
    }
}
exports.MsgClientImpl = MsgClientImpl;
//# sourceMappingURL=tx.js.map