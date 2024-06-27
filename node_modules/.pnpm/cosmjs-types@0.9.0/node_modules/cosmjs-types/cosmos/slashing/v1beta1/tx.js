"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.MsgClientImpl = exports.MsgUpdateParamsResponse = exports.MsgUpdateParams = exports.MsgUnjailResponse = exports.MsgUnjail = exports.protobufPackage = void 0;
/* eslint-disable */
const slashing_1 = require("./slashing");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.slashing.v1beta1";
function createBaseMsgUnjail() {
    return {
        validatorAddr: "",
    };
}
exports.MsgUnjail = {
    typeUrl: "/cosmos.slashing.v1beta1.MsgUnjail",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.validatorAddr !== "") {
            writer.uint32(10).string(message.validatorAddr);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgUnjail();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.validatorAddr = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgUnjail();
        if ((0, helpers_1.isSet)(object.validatorAddr))
            obj.validatorAddr = String(object.validatorAddr);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.validatorAddr !== undefined && (obj.validatorAddr = message.validatorAddr);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgUnjail();
        message.validatorAddr = object.validatorAddr ?? "";
        return message;
    },
};
function createBaseMsgUnjailResponse() {
    return {};
}
exports.MsgUnjailResponse = {
    typeUrl: "/cosmos.slashing.v1beta1.MsgUnjailResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgUnjailResponse();
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
        const obj = createBaseMsgUnjailResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgUnjailResponse();
        return message;
    },
};
function createBaseMsgUpdateParams() {
    return {
        authority: "",
        params: slashing_1.Params.fromPartial({}),
    };
}
exports.MsgUpdateParams = {
    typeUrl: "/cosmos.slashing.v1beta1.MsgUpdateParams",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.authority !== "") {
            writer.uint32(10).string(message.authority);
        }
        if (message.params !== undefined) {
            slashing_1.Params.encode(message.params, writer.uint32(18).fork()).ldelim();
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
                    message.params = slashing_1.Params.decode(reader, reader.uint32());
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
        if ((0, helpers_1.isSet)(object.params))
            obj.params = slashing_1.Params.fromJSON(object.params);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.authority !== undefined && (obj.authority = message.authority);
        message.params !== undefined && (obj.params = message.params ? slashing_1.Params.toJSON(message.params) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgUpdateParams();
        message.authority = object.authority ?? "";
        if (object.params !== undefined && object.params !== null) {
            message.params = slashing_1.Params.fromPartial(object.params);
        }
        return message;
    },
};
function createBaseMsgUpdateParamsResponse() {
    return {};
}
exports.MsgUpdateParamsResponse = {
    typeUrl: "/cosmos.slashing.v1beta1.MsgUpdateParamsResponse",
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
        this.Unjail = this.Unjail.bind(this);
        this.UpdateParams = this.UpdateParams.bind(this);
    }
    Unjail(request) {
        const data = exports.MsgUnjail.encode(request).finish();
        const promise = this.rpc.request("cosmos.slashing.v1beta1.Msg", "Unjail", data);
        return promise.then((data) => exports.MsgUnjailResponse.decode(new binary_1.BinaryReader(data)));
    }
    UpdateParams(request) {
        const data = exports.MsgUpdateParams.encode(request).finish();
        const promise = this.rpc.request("cosmos.slashing.v1beta1.Msg", "UpdateParams", data);
        return promise.then((data) => exports.MsgUpdateParamsResponse.decode(new binary_1.BinaryReader(data)));
    }
}
exports.MsgClientImpl = MsgClientImpl;
//# sourceMappingURL=tx.js.map