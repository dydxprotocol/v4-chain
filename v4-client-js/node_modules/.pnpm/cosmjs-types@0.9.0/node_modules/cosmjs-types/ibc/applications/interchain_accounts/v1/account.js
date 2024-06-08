"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.InterchainAccount = exports.protobufPackage = void 0;
/* eslint-disable */
const auth_1 = require("../../../../cosmos/auth/v1beta1/auth");
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "ibc.applications.interchain_accounts.v1";
function createBaseInterchainAccount() {
    return {
        baseAccount: undefined,
        accountOwner: "",
    };
}
exports.InterchainAccount = {
    typeUrl: "/ibc.applications.interchain_accounts.v1.InterchainAccount",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.baseAccount !== undefined) {
            auth_1.BaseAccount.encode(message.baseAccount, writer.uint32(10).fork()).ldelim();
        }
        if (message.accountOwner !== "") {
            writer.uint32(18).string(message.accountOwner);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseInterchainAccount();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.baseAccount = auth_1.BaseAccount.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.accountOwner = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseInterchainAccount();
        if ((0, helpers_1.isSet)(object.baseAccount))
            obj.baseAccount = auth_1.BaseAccount.fromJSON(object.baseAccount);
        if ((0, helpers_1.isSet)(object.accountOwner))
            obj.accountOwner = String(object.accountOwner);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.baseAccount !== undefined &&
            (obj.baseAccount = message.baseAccount ? auth_1.BaseAccount.toJSON(message.baseAccount) : undefined);
        message.accountOwner !== undefined && (obj.accountOwner = message.accountOwner);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseInterchainAccount();
        if (object.baseAccount !== undefined && object.baseAccount !== null) {
            message.baseAccount = auth_1.BaseAccount.fromPartial(object.baseAccount);
        }
        message.accountOwner = object.accountOwner ?? "";
        return message;
    },
};
//# sourceMappingURL=account.js.map