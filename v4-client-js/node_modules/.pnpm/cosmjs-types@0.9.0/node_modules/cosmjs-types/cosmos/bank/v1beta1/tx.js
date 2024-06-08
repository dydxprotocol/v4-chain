"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.MsgClientImpl = exports.MsgSetSendEnabledResponse = exports.MsgSetSendEnabled = exports.MsgUpdateParamsResponse = exports.MsgUpdateParams = exports.MsgMultiSendResponse = exports.MsgMultiSend = exports.MsgSendResponse = exports.MsgSend = exports.protobufPackage = void 0;
/* eslint-disable */
const coin_1 = require("../../base/v1beta1/coin");
const bank_1 = require("./bank");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.bank.v1beta1";
function createBaseMsgSend() {
    return {
        fromAddress: "",
        toAddress: "",
        amount: [],
    };
}
exports.MsgSend = {
    typeUrl: "/cosmos.bank.v1beta1.MsgSend",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.fromAddress !== "") {
            writer.uint32(10).string(message.fromAddress);
        }
        if (message.toAddress !== "") {
            writer.uint32(18).string(message.toAddress);
        }
        for (const v of message.amount) {
            coin_1.Coin.encode(v, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgSend();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.fromAddress = reader.string();
                    break;
                case 2:
                    message.toAddress = reader.string();
                    break;
                case 3:
                    message.amount.push(coin_1.Coin.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgSend();
        if ((0, helpers_1.isSet)(object.fromAddress))
            obj.fromAddress = String(object.fromAddress);
        if ((0, helpers_1.isSet)(object.toAddress))
            obj.toAddress = String(object.toAddress);
        if (Array.isArray(object?.amount))
            obj.amount = object.amount.map((e) => coin_1.Coin.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.fromAddress !== undefined && (obj.fromAddress = message.fromAddress);
        message.toAddress !== undefined && (obj.toAddress = message.toAddress);
        if (message.amount) {
            obj.amount = message.amount.map((e) => (e ? coin_1.Coin.toJSON(e) : undefined));
        }
        else {
            obj.amount = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgSend();
        message.fromAddress = object.fromAddress ?? "";
        message.toAddress = object.toAddress ?? "";
        message.amount = object.amount?.map((e) => coin_1.Coin.fromPartial(e)) || [];
        return message;
    },
};
function createBaseMsgSendResponse() {
    return {};
}
exports.MsgSendResponse = {
    typeUrl: "/cosmos.bank.v1beta1.MsgSendResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgSendResponse();
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
        const obj = createBaseMsgSendResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgSendResponse();
        return message;
    },
};
function createBaseMsgMultiSend() {
    return {
        inputs: [],
        outputs: [],
    };
}
exports.MsgMultiSend = {
    typeUrl: "/cosmos.bank.v1beta1.MsgMultiSend",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.inputs) {
            bank_1.Input.encode(v, writer.uint32(10).fork()).ldelim();
        }
        for (const v of message.outputs) {
            bank_1.Output.encode(v, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgMultiSend();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.inputs.push(bank_1.Input.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.outputs.push(bank_1.Output.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgMultiSend();
        if (Array.isArray(object?.inputs))
            obj.inputs = object.inputs.map((e) => bank_1.Input.fromJSON(e));
        if (Array.isArray(object?.outputs))
            obj.outputs = object.outputs.map((e) => bank_1.Output.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.inputs) {
            obj.inputs = message.inputs.map((e) => (e ? bank_1.Input.toJSON(e) : undefined));
        }
        else {
            obj.inputs = [];
        }
        if (message.outputs) {
            obj.outputs = message.outputs.map((e) => (e ? bank_1.Output.toJSON(e) : undefined));
        }
        else {
            obj.outputs = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgMultiSend();
        message.inputs = object.inputs?.map((e) => bank_1.Input.fromPartial(e)) || [];
        message.outputs = object.outputs?.map((e) => bank_1.Output.fromPartial(e)) || [];
        return message;
    },
};
function createBaseMsgMultiSendResponse() {
    return {};
}
exports.MsgMultiSendResponse = {
    typeUrl: "/cosmos.bank.v1beta1.MsgMultiSendResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgMultiSendResponse();
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
        const obj = createBaseMsgMultiSendResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgMultiSendResponse();
        return message;
    },
};
function createBaseMsgUpdateParams() {
    return {
        authority: "",
        params: bank_1.Params.fromPartial({}),
    };
}
exports.MsgUpdateParams = {
    typeUrl: "/cosmos.bank.v1beta1.MsgUpdateParams",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.authority !== "") {
            writer.uint32(10).string(message.authority);
        }
        if (message.params !== undefined) {
            bank_1.Params.encode(message.params, writer.uint32(18).fork()).ldelim();
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
                    message.params = bank_1.Params.decode(reader, reader.uint32());
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
            obj.params = bank_1.Params.fromJSON(object.params);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.authority !== undefined && (obj.authority = message.authority);
        message.params !== undefined && (obj.params = message.params ? bank_1.Params.toJSON(message.params) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgUpdateParams();
        message.authority = object.authority ?? "";
        if (object.params !== undefined && object.params !== null) {
            message.params = bank_1.Params.fromPartial(object.params);
        }
        return message;
    },
};
function createBaseMsgUpdateParamsResponse() {
    return {};
}
exports.MsgUpdateParamsResponse = {
    typeUrl: "/cosmos.bank.v1beta1.MsgUpdateParamsResponse",
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
function createBaseMsgSetSendEnabled() {
    return {
        authority: "",
        sendEnabled: [],
        useDefaultFor: [],
    };
}
exports.MsgSetSendEnabled = {
    typeUrl: "/cosmos.bank.v1beta1.MsgSetSendEnabled",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.authority !== "") {
            writer.uint32(10).string(message.authority);
        }
        for (const v of message.sendEnabled) {
            bank_1.SendEnabled.encode(v, writer.uint32(18).fork()).ldelim();
        }
        for (const v of message.useDefaultFor) {
            writer.uint32(26).string(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgSetSendEnabled();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.authority = reader.string();
                    break;
                case 2:
                    message.sendEnabled.push(bank_1.SendEnabled.decode(reader, reader.uint32()));
                    break;
                case 3:
                    message.useDefaultFor.push(reader.string());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgSetSendEnabled();
        if ((0, helpers_1.isSet)(object.authority))
            obj.authority = String(object.authority);
        if (Array.isArray(object?.sendEnabled))
            obj.sendEnabled = object.sendEnabled.map((e) => bank_1.SendEnabled.fromJSON(e));
        if (Array.isArray(object?.useDefaultFor))
            obj.useDefaultFor = object.useDefaultFor.map((e) => String(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.authority !== undefined && (obj.authority = message.authority);
        if (message.sendEnabled) {
            obj.sendEnabled = message.sendEnabled.map((e) => (e ? bank_1.SendEnabled.toJSON(e) : undefined));
        }
        else {
            obj.sendEnabled = [];
        }
        if (message.useDefaultFor) {
            obj.useDefaultFor = message.useDefaultFor.map((e) => e);
        }
        else {
            obj.useDefaultFor = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgSetSendEnabled();
        message.authority = object.authority ?? "";
        message.sendEnabled = object.sendEnabled?.map((e) => bank_1.SendEnabled.fromPartial(e)) || [];
        message.useDefaultFor = object.useDefaultFor?.map((e) => e) || [];
        return message;
    },
};
function createBaseMsgSetSendEnabledResponse() {
    return {};
}
exports.MsgSetSendEnabledResponse = {
    typeUrl: "/cosmos.bank.v1beta1.MsgSetSendEnabledResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgSetSendEnabledResponse();
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
        const obj = createBaseMsgSetSendEnabledResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgSetSendEnabledResponse();
        return message;
    },
};
class MsgClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
        this.Send = this.Send.bind(this);
        this.MultiSend = this.MultiSend.bind(this);
        this.UpdateParams = this.UpdateParams.bind(this);
        this.SetSendEnabled = this.SetSendEnabled.bind(this);
    }
    Send(request) {
        const data = exports.MsgSend.encode(request).finish();
        const promise = this.rpc.request("cosmos.bank.v1beta1.Msg", "Send", data);
        return promise.then((data) => exports.MsgSendResponse.decode(new binary_1.BinaryReader(data)));
    }
    MultiSend(request) {
        const data = exports.MsgMultiSend.encode(request).finish();
        const promise = this.rpc.request("cosmos.bank.v1beta1.Msg", "MultiSend", data);
        return promise.then((data) => exports.MsgMultiSendResponse.decode(new binary_1.BinaryReader(data)));
    }
    UpdateParams(request) {
        const data = exports.MsgUpdateParams.encode(request).finish();
        const promise = this.rpc.request("cosmos.bank.v1beta1.Msg", "UpdateParams", data);
        return promise.then((data) => exports.MsgUpdateParamsResponse.decode(new binary_1.BinaryReader(data)));
    }
    SetSendEnabled(request) {
        const data = exports.MsgSetSendEnabled.encode(request).finish();
        const promise = this.rpc.request("cosmos.bank.v1beta1.Msg", "SetSendEnabled", data);
        return promise.then((data) => exports.MsgSetSendEnabledResponse.decode(new binary_1.BinaryReader(data)));
    }
}
exports.MsgClientImpl = MsgClientImpl;
//# sourceMappingURL=tx.js.map