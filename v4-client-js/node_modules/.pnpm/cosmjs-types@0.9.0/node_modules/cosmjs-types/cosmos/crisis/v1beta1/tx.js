"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.MsgClientImpl = exports.MsgUpdateParamsResponse = exports.MsgUpdateParams = exports.MsgVerifyInvariantResponse = exports.MsgVerifyInvariant = exports.protobufPackage = void 0;
/* eslint-disable */
const coin_1 = require("../../base/v1beta1/coin");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.crisis.v1beta1";
function createBaseMsgVerifyInvariant() {
    return {
        sender: "",
        invariantModuleName: "",
        invariantRoute: "",
    };
}
exports.MsgVerifyInvariant = {
    typeUrl: "/cosmos.crisis.v1beta1.MsgVerifyInvariant",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.sender !== "") {
            writer.uint32(10).string(message.sender);
        }
        if (message.invariantModuleName !== "") {
            writer.uint32(18).string(message.invariantModuleName);
        }
        if (message.invariantRoute !== "") {
            writer.uint32(26).string(message.invariantRoute);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgVerifyInvariant();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.sender = reader.string();
                    break;
                case 2:
                    message.invariantModuleName = reader.string();
                    break;
                case 3:
                    message.invariantRoute = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgVerifyInvariant();
        if ((0, helpers_1.isSet)(object.sender))
            obj.sender = String(object.sender);
        if ((0, helpers_1.isSet)(object.invariantModuleName))
            obj.invariantModuleName = String(object.invariantModuleName);
        if ((0, helpers_1.isSet)(object.invariantRoute))
            obj.invariantRoute = String(object.invariantRoute);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.sender !== undefined && (obj.sender = message.sender);
        message.invariantModuleName !== undefined && (obj.invariantModuleName = message.invariantModuleName);
        message.invariantRoute !== undefined && (obj.invariantRoute = message.invariantRoute);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgVerifyInvariant();
        message.sender = object.sender ?? "";
        message.invariantModuleName = object.invariantModuleName ?? "";
        message.invariantRoute = object.invariantRoute ?? "";
        return message;
    },
};
function createBaseMsgVerifyInvariantResponse() {
    return {};
}
exports.MsgVerifyInvariantResponse = {
    typeUrl: "/cosmos.crisis.v1beta1.MsgVerifyInvariantResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgVerifyInvariantResponse();
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
        const obj = createBaseMsgVerifyInvariantResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgVerifyInvariantResponse();
        return message;
    },
};
function createBaseMsgUpdateParams() {
    return {
        authority: "",
        constantFee: coin_1.Coin.fromPartial({}),
    };
}
exports.MsgUpdateParams = {
    typeUrl: "/cosmos.crisis.v1beta1.MsgUpdateParams",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.authority !== "") {
            writer.uint32(10).string(message.authority);
        }
        if (message.constantFee !== undefined) {
            coin_1.Coin.encode(message.constantFee, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgUpdateParams();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.authority = reader.string();
                    break;
                case 2:
                    message.constantFee = coin_1.Coin.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgUpdateParams();
        if ((0, helpers_1.isSet)(object.authority))
            obj.authority = String(object.authority);
        if ((0, helpers_1.isSet)(object.constantFee))
            obj.constantFee = coin_1.Coin.fromJSON(object.constantFee);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.authority !== undefined && (obj.authority = message.authority);
        message.constantFee !== undefined &&
            (obj.constantFee = message.constantFee ? coin_1.Coin.toJSON(message.constantFee) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgUpdateParams();
        message.authority = object.authority ?? "";
        if (object.constantFee !== undefined && object.constantFee !== null) {
            message.constantFee = coin_1.Coin.fromPartial(object.constantFee);
        }
        return message;
    },
};
function createBaseMsgUpdateParamsResponse() {
    return {};
}
exports.MsgUpdateParamsResponse = {
    typeUrl: "/cosmos.crisis.v1beta1.MsgUpdateParamsResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgUpdateParamsResponse();
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
        const obj = createBaseMsgUpdateParamsResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgUpdateParamsResponse();
        return message;
    },
};
class MsgClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
        this.VerifyInvariant = this.VerifyInvariant.bind(this);
        this.UpdateParams = this.UpdateParams.bind(this);
    }
    VerifyInvariant(request) {
        const data = exports.MsgVerifyInvariant.encode(request).finish();
        const promise = this.rpc.request("cosmos.crisis.v1beta1.Msg", "VerifyInvariant", data);
        return promise.then((data) => exports.MsgVerifyInvariantResponse.decode(new binary_1.BinaryReader(data)));
    }
    UpdateParams(request) {
        const data = exports.MsgUpdateParams.encode(request).finish();
        const promise = this.rpc.request("cosmos.crisis.v1beta1.Msg", "UpdateParams", data);
        return promise.then((data) => exports.MsgUpdateParamsResponse.decode(new binary_1.BinaryReader(data)));
    }
}
exports.MsgClientImpl = MsgClientImpl;
//# sourceMappingURL=tx.js.map