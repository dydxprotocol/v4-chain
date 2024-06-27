"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.MsgClientImpl = exports.MsgCommunityPoolSpendResponse = exports.MsgCommunityPoolSpend = exports.MsgUpdateParamsResponse = exports.MsgUpdateParams = exports.MsgFundCommunityPoolResponse = exports.MsgFundCommunityPool = exports.MsgWithdrawValidatorCommissionResponse = exports.MsgWithdrawValidatorCommission = exports.MsgWithdrawDelegatorRewardResponse = exports.MsgWithdrawDelegatorReward = exports.MsgSetWithdrawAddressResponse = exports.MsgSetWithdrawAddress = exports.protobufPackage = void 0;
/* eslint-disable */
const coin_1 = require("../../base/v1beta1/coin");
const distribution_1 = require("./distribution");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.distribution.v1beta1";
function createBaseMsgSetWithdrawAddress() {
    return {
        delegatorAddress: "",
        withdrawAddress: "",
    };
}
exports.MsgSetWithdrawAddress = {
    typeUrl: "/cosmos.distribution.v1beta1.MsgSetWithdrawAddress",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.delegatorAddress !== "") {
            writer.uint32(10).string(message.delegatorAddress);
        }
        if (message.withdrawAddress !== "") {
            writer.uint32(18).string(message.withdrawAddress);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgSetWithdrawAddress();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.delegatorAddress = reader.string();
                    break;
                case 2:
                    message.withdrawAddress = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgSetWithdrawAddress();
        if ((0, helpers_1.isSet)(object.delegatorAddress))
            obj.delegatorAddress = String(object.delegatorAddress);
        if ((0, helpers_1.isSet)(object.withdrawAddress))
            obj.withdrawAddress = String(object.withdrawAddress);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.delegatorAddress !== undefined && (obj.delegatorAddress = message.delegatorAddress);
        message.withdrawAddress !== undefined && (obj.withdrawAddress = message.withdrawAddress);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgSetWithdrawAddress();
        message.delegatorAddress = object.delegatorAddress ?? "";
        message.withdrawAddress = object.withdrawAddress ?? "";
        return message;
    },
};
function createBaseMsgSetWithdrawAddressResponse() {
    return {};
}
exports.MsgSetWithdrawAddressResponse = {
    typeUrl: "/cosmos.distribution.v1beta1.MsgSetWithdrawAddressResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgSetWithdrawAddressResponse();
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
        const obj = createBaseMsgSetWithdrawAddressResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgSetWithdrawAddressResponse();
        return message;
    },
};
function createBaseMsgWithdrawDelegatorReward() {
    return {
        delegatorAddress: "",
        validatorAddress: "",
    };
}
exports.MsgWithdrawDelegatorReward = {
    typeUrl: "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.delegatorAddress !== "") {
            writer.uint32(10).string(message.delegatorAddress);
        }
        if (message.validatorAddress !== "") {
            writer.uint32(18).string(message.validatorAddress);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgWithdrawDelegatorReward();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.delegatorAddress = reader.string();
                    break;
                case 2:
                    message.validatorAddress = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgWithdrawDelegatorReward();
        if ((0, helpers_1.isSet)(object.delegatorAddress))
            obj.delegatorAddress = String(object.delegatorAddress);
        if ((0, helpers_1.isSet)(object.validatorAddress))
            obj.validatorAddress = String(object.validatorAddress);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.delegatorAddress !== undefined && (obj.delegatorAddress = message.delegatorAddress);
        message.validatorAddress !== undefined && (obj.validatorAddress = message.validatorAddress);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgWithdrawDelegatorReward();
        message.delegatorAddress = object.delegatorAddress ?? "";
        message.validatorAddress = object.validatorAddress ?? "";
        return message;
    },
};
function createBaseMsgWithdrawDelegatorRewardResponse() {
    return {
        amount: [],
    };
}
exports.MsgWithdrawDelegatorRewardResponse = {
    typeUrl: "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorRewardResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.amount) {
            coin_1.Coin.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgWithdrawDelegatorRewardResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
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
        const obj = createBaseMsgWithdrawDelegatorRewardResponse();
        if (Array.isArray(object?.amount))
            obj.amount = object.amount.map((e) => coin_1.Coin.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.amount) {
            obj.amount = message.amount.map((e) => (e ? coin_1.Coin.toJSON(e) : undefined));
        }
        else {
            obj.amount = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgWithdrawDelegatorRewardResponse();
        message.amount = object.amount?.map((e) => coin_1.Coin.fromPartial(e)) || [];
        return message;
    },
};
function createBaseMsgWithdrawValidatorCommission() {
    return {
        validatorAddress: "",
    };
}
exports.MsgWithdrawValidatorCommission = {
    typeUrl: "/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.validatorAddress !== "") {
            writer.uint32(10).string(message.validatorAddress);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgWithdrawValidatorCommission();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.validatorAddress = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgWithdrawValidatorCommission();
        if ((0, helpers_1.isSet)(object.validatorAddress))
            obj.validatorAddress = String(object.validatorAddress);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.validatorAddress !== undefined && (obj.validatorAddress = message.validatorAddress);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgWithdrawValidatorCommission();
        message.validatorAddress = object.validatorAddress ?? "";
        return message;
    },
};
function createBaseMsgWithdrawValidatorCommissionResponse() {
    return {
        amount: [],
    };
}
exports.MsgWithdrawValidatorCommissionResponse = {
    typeUrl: "/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommissionResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.amount) {
            coin_1.Coin.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgWithdrawValidatorCommissionResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
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
        const obj = createBaseMsgWithdrawValidatorCommissionResponse();
        if (Array.isArray(object?.amount))
            obj.amount = object.amount.map((e) => coin_1.Coin.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.amount) {
            obj.amount = message.amount.map((e) => (e ? coin_1.Coin.toJSON(e) : undefined));
        }
        else {
            obj.amount = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgWithdrawValidatorCommissionResponse();
        message.amount = object.amount?.map((e) => coin_1.Coin.fromPartial(e)) || [];
        return message;
    },
};
function createBaseMsgFundCommunityPool() {
    return {
        amount: [],
        depositor: "",
    };
}
exports.MsgFundCommunityPool = {
    typeUrl: "/cosmos.distribution.v1beta1.MsgFundCommunityPool",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.amount) {
            coin_1.Coin.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.depositor !== "") {
            writer.uint32(18).string(message.depositor);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgFundCommunityPool();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.amount.push(coin_1.Coin.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.depositor = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgFundCommunityPool();
        if (Array.isArray(object?.amount))
            obj.amount = object.amount.map((e) => coin_1.Coin.fromJSON(e));
        if ((0, helpers_1.isSet)(object.depositor))
            obj.depositor = String(object.depositor);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.amount) {
            obj.amount = message.amount.map((e) => (e ? coin_1.Coin.toJSON(e) : undefined));
        }
        else {
            obj.amount = [];
        }
        message.depositor !== undefined && (obj.depositor = message.depositor);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgFundCommunityPool();
        message.amount = object.amount?.map((e) => coin_1.Coin.fromPartial(e)) || [];
        message.depositor = object.depositor ?? "";
        return message;
    },
};
function createBaseMsgFundCommunityPoolResponse() {
    return {};
}
exports.MsgFundCommunityPoolResponse = {
    typeUrl: "/cosmos.distribution.v1beta1.MsgFundCommunityPoolResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgFundCommunityPoolResponse();
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
        const obj = createBaseMsgFundCommunityPoolResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgFundCommunityPoolResponse();
        return message;
    },
};
function createBaseMsgUpdateParams() {
    return {
        authority: "",
        params: distribution_1.Params.fromPartial({}),
    };
}
exports.MsgUpdateParams = {
    typeUrl: "/cosmos.distribution.v1beta1.MsgUpdateParams",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.authority !== "") {
            writer.uint32(10).string(message.authority);
        }
        if (message.params !== undefined) {
            distribution_1.Params.encode(message.params, writer.uint32(18).fork()).ldelim();
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
                    message.params = distribution_1.Params.decode(reader, reader.uint32());
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
            obj.params = distribution_1.Params.fromJSON(object.params);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.authority !== undefined && (obj.authority = message.authority);
        message.params !== undefined && (obj.params = message.params ? distribution_1.Params.toJSON(message.params) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgUpdateParams();
        message.authority = object.authority ?? "";
        if (object.params !== undefined && object.params !== null) {
            message.params = distribution_1.Params.fromPartial(object.params);
        }
        return message;
    },
};
function createBaseMsgUpdateParamsResponse() {
    return {};
}
exports.MsgUpdateParamsResponse = {
    typeUrl: "/cosmos.distribution.v1beta1.MsgUpdateParamsResponse",
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
function createBaseMsgCommunityPoolSpend() {
    return {
        authority: "",
        recipient: "",
        amount: [],
    };
}
exports.MsgCommunityPoolSpend = {
    typeUrl: "/cosmos.distribution.v1beta1.MsgCommunityPoolSpend",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.authority !== "") {
            writer.uint32(10).string(message.authority);
        }
        if (message.recipient !== "") {
            writer.uint32(18).string(message.recipient);
        }
        for (const v of message.amount) {
            coin_1.Coin.encode(v, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgCommunityPoolSpend();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.authority = reader.string();
                    break;
                case 2:
                    message.recipient = reader.string();
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
        const obj = createBaseMsgCommunityPoolSpend();
        if ((0, helpers_1.isSet)(object.authority))
            obj.authority = String(object.authority);
        if ((0, helpers_1.isSet)(object.recipient))
            obj.recipient = String(object.recipient);
        if (Array.isArray(object?.amount))
            obj.amount = object.amount.map((e) => coin_1.Coin.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.authority !== undefined && (obj.authority = message.authority);
        message.recipient !== undefined && (obj.recipient = message.recipient);
        if (message.amount) {
            obj.amount = message.amount.map((e) => (e ? coin_1.Coin.toJSON(e) : undefined));
        }
        else {
            obj.amount = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgCommunityPoolSpend();
        message.authority = object.authority ?? "";
        message.recipient = object.recipient ?? "";
        message.amount = object.amount?.map((e) => coin_1.Coin.fromPartial(e)) || [];
        return message;
    },
};
function createBaseMsgCommunityPoolSpendResponse() {
    return {};
}
exports.MsgCommunityPoolSpendResponse = {
    typeUrl: "/cosmos.distribution.v1beta1.MsgCommunityPoolSpendResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgCommunityPoolSpendResponse();
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
        const obj = createBaseMsgCommunityPoolSpendResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgCommunityPoolSpendResponse();
        return message;
    },
};
class MsgClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
        this.SetWithdrawAddress = this.SetWithdrawAddress.bind(this);
        this.WithdrawDelegatorReward = this.WithdrawDelegatorReward.bind(this);
        this.WithdrawValidatorCommission = this.WithdrawValidatorCommission.bind(this);
        this.FundCommunityPool = this.FundCommunityPool.bind(this);
        this.UpdateParams = this.UpdateParams.bind(this);
        this.CommunityPoolSpend = this.CommunityPoolSpend.bind(this);
    }
    SetWithdrawAddress(request) {
        const data = exports.MsgSetWithdrawAddress.encode(request).finish();
        const promise = this.rpc.request("cosmos.distribution.v1beta1.Msg", "SetWithdrawAddress", data);
        return promise.then((data) => exports.MsgSetWithdrawAddressResponse.decode(new binary_1.BinaryReader(data)));
    }
    WithdrawDelegatorReward(request) {
        const data = exports.MsgWithdrawDelegatorReward.encode(request).finish();
        const promise = this.rpc.request("cosmos.distribution.v1beta1.Msg", "WithdrawDelegatorReward", data);
        return promise.then((data) => exports.MsgWithdrawDelegatorRewardResponse.decode(new binary_1.BinaryReader(data)));
    }
    WithdrawValidatorCommission(request) {
        const data = exports.MsgWithdrawValidatorCommission.encode(request).finish();
        const promise = this.rpc.request("cosmos.distribution.v1beta1.Msg", "WithdrawValidatorCommission", data);
        return promise.then((data) => exports.MsgWithdrawValidatorCommissionResponse.decode(new binary_1.BinaryReader(data)));
    }
    FundCommunityPool(request) {
        const data = exports.MsgFundCommunityPool.encode(request).finish();
        const promise = this.rpc.request("cosmos.distribution.v1beta1.Msg", "FundCommunityPool", data);
        return promise.then((data) => exports.MsgFundCommunityPoolResponse.decode(new binary_1.BinaryReader(data)));
    }
    UpdateParams(request) {
        const data = exports.MsgUpdateParams.encode(request).finish();
        const promise = this.rpc.request("cosmos.distribution.v1beta1.Msg", "UpdateParams", data);
        return promise.then((data) => exports.MsgUpdateParamsResponse.decode(new binary_1.BinaryReader(data)));
    }
    CommunityPoolSpend(request) {
        const data = exports.MsgCommunityPoolSpend.encode(request).finish();
        const promise = this.rpc.request("cosmos.distribution.v1beta1.Msg", "CommunityPoolSpend", data);
        return promise.then((data) => exports.MsgCommunityPoolSpendResponse.decode(new binary_1.BinaryReader(data)));
    }
}
exports.MsgClientImpl = MsgClientImpl;
//# sourceMappingURL=tx.js.map