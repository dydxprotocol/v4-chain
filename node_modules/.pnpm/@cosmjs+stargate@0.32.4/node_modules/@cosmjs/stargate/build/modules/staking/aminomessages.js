"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.createStakingAminoConverters = exports.isAminoMsgCancelUnbondingDelegation = exports.isAminoMsgUndelegate = exports.isAminoMsgBeginRedelegate = exports.isAminoMsgDelegate = exports.isAminoMsgEditValidator = exports.isAminoMsgCreateValidator = exports.protoDecimalToJson = void 0;
const math_1 = require("@cosmjs/math");
const proto_signing_1 = require("@cosmjs/proto-signing");
const utils_1 = require("@cosmjs/utils");
function protoDecimalToJson(decimal) {
    const parsed = math_1.Decimal.fromAtomics(decimal, 18);
    const [whole, fractional] = parsed.toString().split(".");
    return `${whole}.${(fractional ?? "").padEnd(18, "0")}`;
}
exports.protoDecimalToJson = protoDecimalToJson;
function jsonDecimalToProto(decimal) {
    const parsed = math_1.Decimal.fromUserInput(decimal, 18);
    return parsed.atomics;
}
function isAminoMsgCreateValidator(msg) {
    return msg.type === "cosmos-sdk/MsgCreateValidator";
}
exports.isAminoMsgCreateValidator = isAminoMsgCreateValidator;
function isAminoMsgEditValidator(msg) {
    return msg.type === "cosmos-sdk/MsgEditValidator";
}
exports.isAminoMsgEditValidator = isAminoMsgEditValidator;
function isAminoMsgDelegate(msg) {
    return msg.type === "cosmos-sdk/MsgDelegate";
}
exports.isAminoMsgDelegate = isAminoMsgDelegate;
function isAminoMsgBeginRedelegate(msg) {
    return msg.type === "cosmos-sdk/MsgBeginRedelegate";
}
exports.isAminoMsgBeginRedelegate = isAminoMsgBeginRedelegate;
function isAminoMsgUndelegate(msg) {
    return msg.type === "cosmos-sdk/MsgUndelegate";
}
exports.isAminoMsgUndelegate = isAminoMsgUndelegate;
function isAminoMsgCancelUnbondingDelegation(msg) {
    return msg.type === "cosmos-sdk/MsgCancelUnbondingDelegation";
}
exports.isAminoMsgCancelUnbondingDelegation = isAminoMsgCancelUnbondingDelegation;
function createStakingAminoConverters() {
    return {
        "/cosmos.staking.v1beta1.MsgBeginRedelegate": {
            aminoType: "cosmos-sdk/MsgBeginRedelegate",
            toAmino: ({ delegatorAddress, validatorSrcAddress, validatorDstAddress, amount, }) => {
                (0, utils_1.assertDefinedAndNotNull)(amount, "missing amount");
                return {
                    delegator_address: delegatorAddress,
                    validator_src_address: validatorSrcAddress,
                    validator_dst_address: validatorDstAddress,
                    amount: amount,
                };
            },
            fromAmino: ({ delegator_address, validator_src_address, validator_dst_address, amount, }) => ({
                delegatorAddress: delegator_address,
                validatorSrcAddress: validator_src_address,
                validatorDstAddress: validator_dst_address,
                amount: amount,
            }),
        },
        "/cosmos.staking.v1beta1.MsgCreateValidator": {
            aminoType: "cosmos-sdk/MsgCreateValidator",
            toAmino: ({ description, commission, minSelfDelegation, delegatorAddress, validatorAddress, pubkey, value, }) => {
                (0, utils_1.assertDefinedAndNotNull)(description, "missing description");
                (0, utils_1.assertDefinedAndNotNull)(commission, "missing commission");
                (0, utils_1.assertDefinedAndNotNull)(pubkey, "missing pubkey");
                (0, utils_1.assertDefinedAndNotNull)(value, "missing value");
                return {
                    description: {
                        moniker: description.moniker,
                        identity: description.identity,
                        website: description.website,
                        security_contact: description.securityContact,
                        details: description.details,
                    },
                    commission: {
                        rate: protoDecimalToJson(commission.rate),
                        max_rate: protoDecimalToJson(commission.maxRate),
                        max_change_rate: protoDecimalToJson(commission.maxChangeRate),
                    },
                    min_self_delegation: minSelfDelegation,
                    delegator_address: delegatorAddress,
                    validator_address: validatorAddress,
                    pubkey: (0, proto_signing_1.decodePubkey)(pubkey),
                    value: value,
                };
            },
            fromAmino: ({ description, commission, min_self_delegation, delegator_address, validator_address, pubkey, value, }) => {
                return {
                    description: {
                        moniker: description.moniker,
                        identity: description.identity,
                        website: description.website,
                        securityContact: description.security_contact,
                        details: description.details,
                    },
                    commission: {
                        rate: jsonDecimalToProto(commission.rate),
                        maxRate: jsonDecimalToProto(commission.max_rate),
                        maxChangeRate: jsonDecimalToProto(commission.max_change_rate),
                    },
                    minSelfDelegation: min_self_delegation,
                    delegatorAddress: delegator_address,
                    validatorAddress: validator_address,
                    pubkey: (0, proto_signing_1.encodePubkey)(pubkey),
                    value: value,
                };
            },
        },
        "/cosmos.staking.v1beta1.MsgDelegate": {
            aminoType: "cosmos-sdk/MsgDelegate",
            toAmino: ({ delegatorAddress, validatorAddress, amount }) => {
                (0, utils_1.assertDefinedAndNotNull)(amount, "missing amount");
                return {
                    delegator_address: delegatorAddress,
                    validator_address: validatorAddress,
                    amount: amount,
                };
            },
            fromAmino: ({ delegator_address, validator_address, amount, }) => ({
                delegatorAddress: delegator_address,
                validatorAddress: validator_address,
                amount: amount,
            }),
        },
        "/cosmos.staking.v1beta1.MsgEditValidator": {
            aminoType: "cosmos-sdk/MsgEditValidator",
            toAmino: ({ description, commissionRate, minSelfDelegation, validatorAddress, }) => {
                (0, utils_1.assertDefinedAndNotNull)(description, "missing description");
                return {
                    description: {
                        moniker: description.moniker,
                        identity: description.identity,
                        website: description.website,
                        security_contact: description.securityContact,
                        details: description.details,
                    },
                    // empty string in the protobuf document means "do not change"
                    commission_rate: commissionRate ? protoDecimalToJson(commissionRate) : undefined,
                    // empty string in the protobuf document means "do not change"
                    min_self_delegation: minSelfDelegation ? minSelfDelegation : undefined,
                    validator_address: validatorAddress,
                };
            },
            fromAmino: ({ description, commission_rate, min_self_delegation, validator_address, }) => ({
                description: {
                    moniker: description.moniker,
                    identity: description.identity,
                    website: description.website,
                    securityContact: description.security_contact,
                    details: description.details,
                },
                // empty string in the protobuf document means "do not change"
                commissionRate: commission_rate ? jsonDecimalToProto(commission_rate) : "",
                // empty string in the protobuf document means "do not change"
                minSelfDelegation: min_self_delegation ?? "",
                validatorAddress: validator_address,
            }),
        },
        "/cosmos.staking.v1beta1.MsgUndelegate": {
            aminoType: "cosmos-sdk/MsgUndelegate",
            toAmino: ({ delegatorAddress, validatorAddress, amount, }) => {
                (0, utils_1.assertDefinedAndNotNull)(amount, "missing amount");
                return {
                    delegator_address: delegatorAddress,
                    validator_address: validatorAddress,
                    amount: amount,
                };
            },
            fromAmino: ({ delegator_address, validator_address, amount, }) => ({
                delegatorAddress: delegator_address,
                validatorAddress: validator_address,
                amount: amount,
            }),
        },
        "/cosmos.staking.v1beta1.MsgCancelUnbondingDelegation": {
            aminoType: "cosmos-sdk/MsgCancelUnbondingDelegation",
            toAmino: ({ delegatorAddress, validatorAddress, amount, creationHeight, }) => {
                (0, utils_1.assertDefinedAndNotNull)(amount, "missing amount");
                return {
                    delegator_address: delegatorAddress,
                    validator_address: validatorAddress,
                    amount: amount,
                    creation_height: creationHeight.toString(),
                };
            },
            fromAmino: ({ delegator_address, validator_address, amount, creation_height, }) => ({
                delegatorAddress: delegator_address,
                validatorAddress: validator_address,
                amount: amount,
                creationHeight: BigInt(creation_height),
            }),
        },
    };
}
exports.createStakingAminoConverters = createStakingAminoConverters;
//# sourceMappingURL=aminomessages.js.map