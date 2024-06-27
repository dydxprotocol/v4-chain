"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ValidatorUpdates = exports.Pool = exports.RedelegationResponse = exports.RedelegationEntryResponse = exports.DelegationResponse = exports.Params = exports.Redelegation = exports.RedelegationEntry = exports.UnbondingDelegationEntry = exports.UnbondingDelegation = exports.Delegation = exports.DVVTriplets = exports.DVVTriplet = exports.DVPairs = exports.DVPair = exports.ValAddresses = exports.Validator = exports.Description = exports.Commission = exports.CommissionRates = exports.HistoricalInfo = exports.infractionToJSON = exports.infractionFromJSON = exports.Infraction = exports.bondStatusToJSON = exports.bondStatusFromJSON = exports.BondStatus = exports.protobufPackage = void 0;
/* eslint-disable */
const types_1 = require("../../../tendermint/types/types");
const timestamp_1 = require("../../../google/protobuf/timestamp");
const any_1 = require("../../../google/protobuf/any");
const duration_1 = require("../../../google/protobuf/duration");
const coin_1 = require("../../base/v1beta1/coin");
const types_2 = require("../../../tendermint/abci/types");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.staking.v1beta1";
/** BondStatus is the status of a validator. */
var BondStatus;
(function (BondStatus) {
    /** BOND_STATUS_UNSPECIFIED - UNSPECIFIED defines an invalid validator status. */
    BondStatus[BondStatus["BOND_STATUS_UNSPECIFIED"] = 0] = "BOND_STATUS_UNSPECIFIED";
    /** BOND_STATUS_UNBONDED - UNBONDED defines a validator that is not bonded. */
    BondStatus[BondStatus["BOND_STATUS_UNBONDED"] = 1] = "BOND_STATUS_UNBONDED";
    /** BOND_STATUS_UNBONDING - UNBONDING defines a validator that is unbonding. */
    BondStatus[BondStatus["BOND_STATUS_UNBONDING"] = 2] = "BOND_STATUS_UNBONDING";
    /** BOND_STATUS_BONDED - BONDED defines a validator that is bonded. */
    BondStatus[BondStatus["BOND_STATUS_BONDED"] = 3] = "BOND_STATUS_BONDED";
    BondStatus[BondStatus["UNRECOGNIZED"] = -1] = "UNRECOGNIZED";
})(BondStatus || (exports.BondStatus = BondStatus = {}));
function bondStatusFromJSON(object) {
    switch (object) {
        case 0:
        case "BOND_STATUS_UNSPECIFIED":
            return BondStatus.BOND_STATUS_UNSPECIFIED;
        case 1:
        case "BOND_STATUS_UNBONDED":
            return BondStatus.BOND_STATUS_UNBONDED;
        case 2:
        case "BOND_STATUS_UNBONDING":
            return BondStatus.BOND_STATUS_UNBONDING;
        case 3:
        case "BOND_STATUS_BONDED":
            return BondStatus.BOND_STATUS_BONDED;
        case -1:
        case "UNRECOGNIZED":
        default:
            return BondStatus.UNRECOGNIZED;
    }
}
exports.bondStatusFromJSON = bondStatusFromJSON;
function bondStatusToJSON(object) {
    switch (object) {
        case BondStatus.BOND_STATUS_UNSPECIFIED:
            return "BOND_STATUS_UNSPECIFIED";
        case BondStatus.BOND_STATUS_UNBONDED:
            return "BOND_STATUS_UNBONDED";
        case BondStatus.BOND_STATUS_UNBONDING:
            return "BOND_STATUS_UNBONDING";
        case BondStatus.BOND_STATUS_BONDED:
            return "BOND_STATUS_BONDED";
        case BondStatus.UNRECOGNIZED:
        default:
            return "UNRECOGNIZED";
    }
}
exports.bondStatusToJSON = bondStatusToJSON;
/** Infraction indicates the infraction a validator commited. */
var Infraction;
(function (Infraction) {
    /** INFRACTION_UNSPECIFIED - UNSPECIFIED defines an empty infraction. */
    Infraction[Infraction["INFRACTION_UNSPECIFIED"] = 0] = "INFRACTION_UNSPECIFIED";
    /** INFRACTION_DOUBLE_SIGN - DOUBLE_SIGN defines a validator that double-signs a block. */
    Infraction[Infraction["INFRACTION_DOUBLE_SIGN"] = 1] = "INFRACTION_DOUBLE_SIGN";
    /** INFRACTION_DOWNTIME - DOWNTIME defines a validator that missed signing too many blocks. */
    Infraction[Infraction["INFRACTION_DOWNTIME"] = 2] = "INFRACTION_DOWNTIME";
    Infraction[Infraction["UNRECOGNIZED"] = -1] = "UNRECOGNIZED";
})(Infraction || (exports.Infraction = Infraction = {}));
function infractionFromJSON(object) {
    switch (object) {
        case 0:
        case "INFRACTION_UNSPECIFIED":
            return Infraction.INFRACTION_UNSPECIFIED;
        case 1:
        case "INFRACTION_DOUBLE_SIGN":
            return Infraction.INFRACTION_DOUBLE_SIGN;
        case 2:
        case "INFRACTION_DOWNTIME":
            return Infraction.INFRACTION_DOWNTIME;
        case -1:
        case "UNRECOGNIZED":
        default:
            return Infraction.UNRECOGNIZED;
    }
}
exports.infractionFromJSON = infractionFromJSON;
function infractionToJSON(object) {
    switch (object) {
        case Infraction.INFRACTION_UNSPECIFIED:
            return "INFRACTION_UNSPECIFIED";
        case Infraction.INFRACTION_DOUBLE_SIGN:
            return "INFRACTION_DOUBLE_SIGN";
        case Infraction.INFRACTION_DOWNTIME:
            return "INFRACTION_DOWNTIME";
        case Infraction.UNRECOGNIZED:
        default:
            return "UNRECOGNIZED";
    }
}
exports.infractionToJSON = infractionToJSON;
function createBaseHistoricalInfo() {
    return {
        header: types_1.Header.fromPartial({}),
        valset: [],
    };
}
exports.HistoricalInfo = {
    typeUrl: "/cosmos.staking.v1beta1.HistoricalInfo",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.header !== undefined) {
            types_1.Header.encode(message.header, writer.uint32(10).fork()).ldelim();
        }
        for (const v of message.valset) {
            exports.Validator.encode(v, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseHistoricalInfo();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.header = types_1.Header.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.valset.push(exports.Validator.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseHistoricalInfo();
        if ((0, helpers_1.isSet)(object.header))
            obj.header = types_1.Header.fromJSON(object.header);
        if (Array.isArray(object?.valset))
            obj.valset = object.valset.map((e) => exports.Validator.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.header !== undefined && (obj.header = message.header ? types_1.Header.toJSON(message.header) : undefined);
        if (message.valset) {
            obj.valset = message.valset.map((e) => (e ? exports.Validator.toJSON(e) : undefined));
        }
        else {
            obj.valset = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseHistoricalInfo();
        if (object.header !== undefined && object.header !== null) {
            message.header = types_1.Header.fromPartial(object.header);
        }
        message.valset = object.valset?.map((e) => exports.Validator.fromPartial(e)) || [];
        return message;
    },
};
function createBaseCommissionRates() {
    return {
        rate: "",
        maxRate: "",
        maxChangeRate: "",
    };
}
exports.CommissionRates = {
    typeUrl: "/cosmos.staking.v1beta1.CommissionRates",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.rate !== "") {
            writer.uint32(10).string(message.rate);
        }
        if (message.maxRate !== "") {
            writer.uint32(18).string(message.maxRate);
        }
        if (message.maxChangeRate !== "") {
            writer.uint32(26).string(message.maxChangeRate);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseCommissionRates();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.rate = reader.string();
                    break;
                case 2:
                    message.maxRate = reader.string();
                    break;
                case 3:
                    message.maxChangeRate = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseCommissionRates();
        if ((0, helpers_1.isSet)(object.rate))
            obj.rate = String(object.rate);
        if ((0, helpers_1.isSet)(object.maxRate))
            obj.maxRate = String(object.maxRate);
        if ((0, helpers_1.isSet)(object.maxChangeRate))
            obj.maxChangeRate = String(object.maxChangeRate);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.rate !== undefined && (obj.rate = message.rate);
        message.maxRate !== undefined && (obj.maxRate = message.maxRate);
        message.maxChangeRate !== undefined && (obj.maxChangeRate = message.maxChangeRate);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseCommissionRates();
        message.rate = object.rate ?? "";
        message.maxRate = object.maxRate ?? "";
        message.maxChangeRate = object.maxChangeRate ?? "";
        return message;
    },
};
function createBaseCommission() {
    return {
        commissionRates: exports.CommissionRates.fromPartial({}),
        updateTime: timestamp_1.Timestamp.fromPartial({}),
    };
}
exports.Commission = {
    typeUrl: "/cosmos.staking.v1beta1.Commission",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.commissionRates !== undefined) {
            exports.CommissionRates.encode(message.commissionRates, writer.uint32(10).fork()).ldelim();
        }
        if (message.updateTime !== undefined) {
            timestamp_1.Timestamp.encode(message.updateTime, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseCommission();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.commissionRates = exports.CommissionRates.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.updateTime = timestamp_1.Timestamp.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseCommission();
        if ((0, helpers_1.isSet)(object.commissionRates))
            obj.commissionRates = exports.CommissionRates.fromJSON(object.commissionRates);
        if ((0, helpers_1.isSet)(object.updateTime))
            obj.updateTime = (0, helpers_1.fromJsonTimestamp)(object.updateTime);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.commissionRates !== undefined &&
            (obj.commissionRates = message.commissionRates
                ? exports.CommissionRates.toJSON(message.commissionRates)
                : undefined);
        message.updateTime !== undefined && (obj.updateTime = (0, helpers_1.fromTimestamp)(message.updateTime).toISOString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseCommission();
        if (object.commissionRates !== undefined && object.commissionRates !== null) {
            message.commissionRates = exports.CommissionRates.fromPartial(object.commissionRates);
        }
        if (object.updateTime !== undefined && object.updateTime !== null) {
            message.updateTime = timestamp_1.Timestamp.fromPartial(object.updateTime);
        }
        return message;
    },
};
function createBaseDescription() {
    return {
        moniker: "",
        identity: "",
        website: "",
        securityContact: "",
        details: "",
    };
}
exports.Description = {
    typeUrl: "/cosmos.staking.v1beta1.Description",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.moniker !== "") {
            writer.uint32(10).string(message.moniker);
        }
        if (message.identity !== "") {
            writer.uint32(18).string(message.identity);
        }
        if (message.website !== "") {
            writer.uint32(26).string(message.website);
        }
        if (message.securityContact !== "") {
            writer.uint32(34).string(message.securityContact);
        }
        if (message.details !== "") {
            writer.uint32(42).string(message.details);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseDescription();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.moniker = reader.string();
                    break;
                case 2:
                    message.identity = reader.string();
                    break;
                case 3:
                    message.website = reader.string();
                    break;
                case 4:
                    message.securityContact = reader.string();
                    break;
                case 5:
                    message.details = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseDescription();
        if ((0, helpers_1.isSet)(object.moniker))
            obj.moniker = String(object.moniker);
        if ((0, helpers_1.isSet)(object.identity))
            obj.identity = String(object.identity);
        if ((0, helpers_1.isSet)(object.website))
            obj.website = String(object.website);
        if ((0, helpers_1.isSet)(object.securityContact))
            obj.securityContact = String(object.securityContact);
        if ((0, helpers_1.isSet)(object.details))
            obj.details = String(object.details);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.moniker !== undefined && (obj.moniker = message.moniker);
        message.identity !== undefined && (obj.identity = message.identity);
        message.website !== undefined && (obj.website = message.website);
        message.securityContact !== undefined && (obj.securityContact = message.securityContact);
        message.details !== undefined && (obj.details = message.details);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseDescription();
        message.moniker = object.moniker ?? "";
        message.identity = object.identity ?? "";
        message.website = object.website ?? "";
        message.securityContact = object.securityContact ?? "";
        message.details = object.details ?? "";
        return message;
    },
};
function createBaseValidator() {
    return {
        operatorAddress: "",
        consensusPubkey: undefined,
        jailed: false,
        status: 0,
        tokens: "",
        delegatorShares: "",
        description: exports.Description.fromPartial({}),
        unbondingHeight: BigInt(0),
        unbondingTime: timestamp_1.Timestamp.fromPartial({}),
        commission: exports.Commission.fromPartial({}),
        minSelfDelegation: "",
        unbondingOnHoldRefCount: BigInt(0),
        unbondingIds: [],
    };
}
exports.Validator = {
    typeUrl: "/cosmos.staking.v1beta1.Validator",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.operatorAddress !== "") {
            writer.uint32(10).string(message.operatorAddress);
        }
        if (message.consensusPubkey !== undefined) {
            any_1.Any.encode(message.consensusPubkey, writer.uint32(18).fork()).ldelim();
        }
        if (message.jailed === true) {
            writer.uint32(24).bool(message.jailed);
        }
        if (message.status !== 0) {
            writer.uint32(32).int32(message.status);
        }
        if (message.tokens !== "") {
            writer.uint32(42).string(message.tokens);
        }
        if (message.delegatorShares !== "") {
            writer.uint32(50).string(message.delegatorShares);
        }
        if (message.description !== undefined) {
            exports.Description.encode(message.description, writer.uint32(58).fork()).ldelim();
        }
        if (message.unbondingHeight !== BigInt(0)) {
            writer.uint32(64).int64(message.unbondingHeight);
        }
        if (message.unbondingTime !== undefined) {
            timestamp_1.Timestamp.encode(message.unbondingTime, writer.uint32(74).fork()).ldelim();
        }
        if (message.commission !== undefined) {
            exports.Commission.encode(message.commission, writer.uint32(82).fork()).ldelim();
        }
        if (message.minSelfDelegation !== "") {
            writer.uint32(90).string(message.minSelfDelegation);
        }
        if (message.unbondingOnHoldRefCount !== BigInt(0)) {
            writer.uint32(96).int64(message.unbondingOnHoldRefCount);
        }
        writer.uint32(106).fork();
        for (const v of message.unbondingIds) {
            writer.uint64(v);
        }
        writer.ldelim();
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseValidator();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.operatorAddress = reader.string();
                    break;
                case 2:
                    message.consensusPubkey = any_1.Any.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.jailed = reader.bool();
                    break;
                case 4:
                    message.status = reader.int32();
                    break;
                case 5:
                    message.tokens = reader.string();
                    break;
                case 6:
                    message.delegatorShares = reader.string();
                    break;
                case 7:
                    message.description = exports.Description.decode(reader, reader.uint32());
                    break;
                case 8:
                    message.unbondingHeight = reader.int64();
                    break;
                case 9:
                    message.unbondingTime = timestamp_1.Timestamp.decode(reader, reader.uint32());
                    break;
                case 10:
                    message.commission = exports.Commission.decode(reader, reader.uint32());
                    break;
                case 11:
                    message.minSelfDelegation = reader.string();
                    break;
                case 12:
                    message.unbondingOnHoldRefCount = reader.int64();
                    break;
                case 13:
                    if ((tag & 7) === 2) {
                        const end2 = reader.uint32() + reader.pos;
                        while (reader.pos < end2) {
                            message.unbondingIds.push(reader.uint64());
                        }
                    }
                    else {
                        message.unbondingIds.push(reader.uint64());
                    }
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseValidator();
        if ((0, helpers_1.isSet)(object.operatorAddress))
            obj.operatorAddress = String(object.operatorAddress);
        if ((0, helpers_1.isSet)(object.consensusPubkey))
            obj.consensusPubkey = any_1.Any.fromJSON(object.consensusPubkey);
        if ((0, helpers_1.isSet)(object.jailed))
            obj.jailed = Boolean(object.jailed);
        if ((0, helpers_1.isSet)(object.status))
            obj.status = bondStatusFromJSON(object.status);
        if ((0, helpers_1.isSet)(object.tokens))
            obj.tokens = String(object.tokens);
        if ((0, helpers_1.isSet)(object.delegatorShares))
            obj.delegatorShares = String(object.delegatorShares);
        if ((0, helpers_1.isSet)(object.description))
            obj.description = exports.Description.fromJSON(object.description);
        if ((0, helpers_1.isSet)(object.unbondingHeight))
            obj.unbondingHeight = BigInt(object.unbondingHeight.toString());
        if ((0, helpers_1.isSet)(object.unbondingTime))
            obj.unbondingTime = (0, helpers_1.fromJsonTimestamp)(object.unbondingTime);
        if ((0, helpers_1.isSet)(object.commission))
            obj.commission = exports.Commission.fromJSON(object.commission);
        if ((0, helpers_1.isSet)(object.minSelfDelegation))
            obj.minSelfDelegation = String(object.minSelfDelegation);
        if ((0, helpers_1.isSet)(object.unbondingOnHoldRefCount))
            obj.unbondingOnHoldRefCount = BigInt(object.unbondingOnHoldRefCount.toString());
        if (Array.isArray(object?.unbondingIds))
            obj.unbondingIds = object.unbondingIds.map((e) => BigInt(e.toString()));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.operatorAddress !== undefined && (obj.operatorAddress = message.operatorAddress);
        message.consensusPubkey !== undefined &&
            (obj.consensusPubkey = message.consensusPubkey ? any_1.Any.toJSON(message.consensusPubkey) : undefined);
        message.jailed !== undefined && (obj.jailed = message.jailed);
        message.status !== undefined && (obj.status = bondStatusToJSON(message.status));
        message.tokens !== undefined && (obj.tokens = message.tokens);
        message.delegatorShares !== undefined && (obj.delegatorShares = message.delegatorShares);
        message.description !== undefined &&
            (obj.description = message.description ? exports.Description.toJSON(message.description) : undefined);
        message.unbondingHeight !== undefined &&
            (obj.unbondingHeight = (message.unbondingHeight || BigInt(0)).toString());
        message.unbondingTime !== undefined &&
            (obj.unbondingTime = (0, helpers_1.fromTimestamp)(message.unbondingTime).toISOString());
        message.commission !== undefined &&
            (obj.commission = message.commission ? exports.Commission.toJSON(message.commission) : undefined);
        message.minSelfDelegation !== undefined && (obj.minSelfDelegation = message.minSelfDelegation);
        message.unbondingOnHoldRefCount !== undefined &&
            (obj.unbondingOnHoldRefCount = (message.unbondingOnHoldRefCount || BigInt(0)).toString());
        if (message.unbondingIds) {
            obj.unbondingIds = message.unbondingIds.map((e) => (e || BigInt(0)).toString());
        }
        else {
            obj.unbondingIds = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseValidator();
        message.operatorAddress = object.operatorAddress ?? "";
        if (object.consensusPubkey !== undefined && object.consensusPubkey !== null) {
            message.consensusPubkey = any_1.Any.fromPartial(object.consensusPubkey);
        }
        message.jailed = object.jailed ?? false;
        message.status = object.status ?? 0;
        message.tokens = object.tokens ?? "";
        message.delegatorShares = object.delegatorShares ?? "";
        if (object.description !== undefined && object.description !== null) {
            message.description = exports.Description.fromPartial(object.description);
        }
        if (object.unbondingHeight !== undefined && object.unbondingHeight !== null) {
            message.unbondingHeight = BigInt(object.unbondingHeight.toString());
        }
        if (object.unbondingTime !== undefined && object.unbondingTime !== null) {
            message.unbondingTime = timestamp_1.Timestamp.fromPartial(object.unbondingTime);
        }
        if (object.commission !== undefined && object.commission !== null) {
            message.commission = exports.Commission.fromPartial(object.commission);
        }
        message.minSelfDelegation = object.minSelfDelegation ?? "";
        if (object.unbondingOnHoldRefCount !== undefined && object.unbondingOnHoldRefCount !== null) {
            message.unbondingOnHoldRefCount = BigInt(object.unbondingOnHoldRefCount.toString());
        }
        message.unbondingIds = object.unbondingIds?.map((e) => BigInt(e.toString())) || [];
        return message;
    },
};
function createBaseValAddresses() {
    return {
        addresses: [],
    };
}
exports.ValAddresses = {
    typeUrl: "/cosmos.staking.v1beta1.ValAddresses",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.addresses) {
            writer.uint32(10).string(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseValAddresses();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.addresses.push(reader.string());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseValAddresses();
        if (Array.isArray(object?.addresses))
            obj.addresses = object.addresses.map((e) => String(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.addresses) {
            obj.addresses = message.addresses.map((e) => e);
        }
        else {
            obj.addresses = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseValAddresses();
        message.addresses = object.addresses?.map((e) => e) || [];
        return message;
    },
};
function createBaseDVPair() {
    return {
        delegatorAddress: "",
        validatorAddress: "",
    };
}
exports.DVPair = {
    typeUrl: "/cosmos.staking.v1beta1.DVPair",
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
        const message = createBaseDVPair();
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
        const obj = createBaseDVPair();
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
        const message = createBaseDVPair();
        message.delegatorAddress = object.delegatorAddress ?? "";
        message.validatorAddress = object.validatorAddress ?? "";
        return message;
    },
};
function createBaseDVPairs() {
    return {
        pairs: [],
    };
}
exports.DVPairs = {
    typeUrl: "/cosmos.staking.v1beta1.DVPairs",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.pairs) {
            exports.DVPair.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseDVPairs();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.pairs.push(exports.DVPair.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseDVPairs();
        if (Array.isArray(object?.pairs))
            obj.pairs = object.pairs.map((e) => exports.DVPair.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.pairs) {
            obj.pairs = message.pairs.map((e) => (e ? exports.DVPair.toJSON(e) : undefined));
        }
        else {
            obj.pairs = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseDVPairs();
        message.pairs = object.pairs?.map((e) => exports.DVPair.fromPartial(e)) || [];
        return message;
    },
};
function createBaseDVVTriplet() {
    return {
        delegatorAddress: "",
        validatorSrcAddress: "",
        validatorDstAddress: "",
    };
}
exports.DVVTriplet = {
    typeUrl: "/cosmos.staking.v1beta1.DVVTriplet",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.delegatorAddress !== "") {
            writer.uint32(10).string(message.delegatorAddress);
        }
        if (message.validatorSrcAddress !== "") {
            writer.uint32(18).string(message.validatorSrcAddress);
        }
        if (message.validatorDstAddress !== "") {
            writer.uint32(26).string(message.validatorDstAddress);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseDVVTriplet();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.delegatorAddress = reader.string();
                    break;
                case 2:
                    message.validatorSrcAddress = reader.string();
                    break;
                case 3:
                    message.validatorDstAddress = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseDVVTriplet();
        if ((0, helpers_1.isSet)(object.delegatorAddress))
            obj.delegatorAddress = String(object.delegatorAddress);
        if ((0, helpers_1.isSet)(object.validatorSrcAddress))
            obj.validatorSrcAddress = String(object.validatorSrcAddress);
        if ((0, helpers_1.isSet)(object.validatorDstAddress))
            obj.validatorDstAddress = String(object.validatorDstAddress);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.delegatorAddress !== undefined && (obj.delegatorAddress = message.delegatorAddress);
        message.validatorSrcAddress !== undefined && (obj.validatorSrcAddress = message.validatorSrcAddress);
        message.validatorDstAddress !== undefined && (obj.validatorDstAddress = message.validatorDstAddress);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseDVVTriplet();
        message.delegatorAddress = object.delegatorAddress ?? "";
        message.validatorSrcAddress = object.validatorSrcAddress ?? "";
        message.validatorDstAddress = object.validatorDstAddress ?? "";
        return message;
    },
};
function createBaseDVVTriplets() {
    return {
        triplets: [],
    };
}
exports.DVVTriplets = {
    typeUrl: "/cosmos.staking.v1beta1.DVVTriplets",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.triplets) {
            exports.DVVTriplet.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseDVVTriplets();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.triplets.push(exports.DVVTriplet.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseDVVTriplets();
        if (Array.isArray(object?.triplets))
            obj.triplets = object.triplets.map((e) => exports.DVVTriplet.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.triplets) {
            obj.triplets = message.triplets.map((e) => (e ? exports.DVVTriplet.toJSON(e) : undefined));
        }
        else {
            obj.triplets = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseDVVTriplets();
        message.triplets = object.triplets?.map((e) => exports.DVVTriplet.fromPartial(e)) || [];
        return message;
    },
};
function createBaseDelegation() {
    return {
        delegatorAddress: "",
        validatorAddress: "",
        shares: "",
    };
}
exports.Delegation = {
    typeUrl: "/cosmos.staking.v1beta1.Delegation",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.delegatorAddress !== "") {
            writer.uint32(10).string(message.delegatorAddress);
        }
        if (message.validatorAddress !== "") {
            writer.uint32(18).string(message.validatorAddress);
        }
        if (message.shares !== "") {
            writer.uint32(26).string(message.shares);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseDelegation();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.delegatorAddress = reader.string();
                    break;
                case 2:
                    message.validatorAddress = reader.string();
                    break;
                case 3:
                    message.shares = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseDelegation();
        if ((0, helpers_1.isSet)(object.delegatorAddress))
            obj.delegatorAddress = String(object.delegatorAddress);
        if ((0, helpers_1.isSet)(object.validatorAddress))
            obj.validatorAddress = String(object.validatorAddress);
        if ((0, helpers_1.isSet)(object.shares))
            obj.shares = String(object.shares);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.delegatorAddress !== undefined && (obj.delegatorAddress = message.delegatorAddress);
        message.validatorAddress !== undefined && (obj.validatorAddress = message.validatorAddress);
        message.shares !== undefined && (obj.shares = message.shares);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseDelegation();
        message.delegatorAddress = object.delegatorAddress ?? "";
        message.validatorAddress = object.validatorAddress ?? "";
        message.shares = object.shares ?? "";
        return message;
    },
};
function createBaseUnbondingDelegation() {
    return {
        delegatorAddress: "",
        validatorAddress: "",
        entries: [],
    };
}
exports.UnbondingDelegation = {
    typeUrl: "/cosmos.staking.v1beta1.UnbondingDelegation",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.delegatorAddress !== "") {
            writer.uint32(10).string(message.delegatorAddress);
        }
        if (message.validatorAddress !== "") {
            writer.uint32(18).string(message.validatorAddress);
        }
        for (const v of message.entries) {
            exports.UnbondingDelegationEntry.encode(v, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseUnbondingDelegation();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.delegatorAddress = reader.string();
                    break;
                case 2:
                    message.validatorAddress = reader.string();
                    break;
                case 3:
                    message.entries.push(exports.UnbondingDelegationEntry.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseUnbondingDelegation();
        if ((0, helpers_1.isSet)(object.delegatorAddress))
            obj.delegatorAddress = String(object.delegatorAddress);
        if ((0, helpers_1.isSet)(object.validatorAddress))
            obj.validatorAddress = String(object.validatorAddress);
        if (Array.isArray(object?.entries))
            obj.entries = object.entries.map((e) => exports.UnbondingDelegationEntry.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.delegatorAddress !== undefined && (obj.delegatorAddress = message.delegatorAddress);
        message.validatorAddress !== undefined && (obj.validatorAddress = message.validatorAddress);
        if (message.entries) {
            obj.entries = message.entries.map((e) => (e ? exports.UnbondingDelegationEntry.toJSON(e) : undefined));
        }
        else {
            obj.entries = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseUnbondingDelegation();
        message.delegatorAddress = object.delegatorAddress ?? "";
        message.validatorAddress = object.validatorAddress ?? "";
        message.entries = object.entries?.map((e) => exports.UnbondingDelegationEntry.fromPartial(e)) || [];
        return message;
    },
};
function createBaseUnbondingDelegationEntry() {
    return {
        creationHeight: BigInt(0),
        completionTime: timestamp_1.Timestamp.fromPartial({}),
        initialBalance: "",
        balance: "",
        unbondingId: BigInt(0),
        unbondingOnHoldRefCount: BigInt(0),
    };
}
exports.UnbondingDelegationEntry = {
    typeUrl: "/cosmos.staking.v1beta1.UnbondingDelegationEntry",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.creationHeight !== BigInt(0)) {
            writer.uint32(8).int64(message.creationHeight);
        }
        if (message.completionTime !== undefined) {
            timestamp_1.Timestamp.encode(message.completionTime, writer.uint32(18).fork()).ldelim();
        }
        if (message.initialBalance !== "") {
            writer.uint32(26).string(message.initialBalance);
        }
        if (message.balance !== "") {
            writer.uint32(34).string(message.balance);
        }
        if (message.unbondingId !== BigInt(0)) {
            writer.uint32(40).uint64(message.unbondingId);
        }
        if (message.unbondingOnHoldRefCount !== BigInt(0)) {
            writer.uint32(48).int64(message.unbondingOnHoldRefCount);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseUnbondingDelegationEntry();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.creationHeight = reader.int64();
                    break;
                case 2:
                    message.completionTime = timestamp_1.Timestamp.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.initialBalance = reader.string();
                    break;
                case 4:
                    message.balance = reader.string();
                    break;
                case 5:
                    message.unbondingId = reader.uint64();
                    break;
                case 6:
                    message.unbondingOnHoldRefCount = reader.int64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseUnbondingDelegationEntry();
        if ((0, helpers_1.isSet)(object.creationHeight))
            obj.creationHeight = BigInt(object.creationHeight.toString());
        if ((0, helpers_1.isSet)(object.completionTime))
            obj.completionTime = (0, helpers_1.fromJsonTimestamp)(object.completionTime);
        if ((0, helpers_1.isSet)(object.initialBalance))
            obj.initialBalance = String(object.initialBalance);
        if ((0, helpers_1.isSet)(object.balance))
            obj.balance = String(object.balance);
        if ((0, helpers_1.isSet)(object.unbondingId))
            obj.unbondingId = BigInt(object.unbondingId.toString());
        if ((0, helpers_1.isSet)(object.unbondingOnHoldRefCount))
            obj.unbondingOnHoldRefCount = BigInt(object.unbondingOnHoldRefCount.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.creationHeight !== undefined &&
            (obj.creationHeight = (message.creationHeight || BigInt(0)).toString());
        message.completionTime !== undefined &&
            (obj.completionTime = (0, helpers_1.fromTimestamp)(message.completionTime).toISOString());
        message.initialBalance !== undefined && (obj.initialBalance = message.initialBalance);
        message.balance !== undefined && (obj.balance = message.balance);
        message.unbondingId !== undefined && (obj.unbondingId = (message.unbondingId || BigInt(0)).toString());
        message.unbondingOnHoldRefCount !== undefined &&
            (obj.unbondingOnHoldRefCount = (message.unbondingOnHoldRefCount || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseUnbondingDelegationEntry();
        if (object.creationHeight !== undefined && object.creationHeight !== null) {
            message.creationHeight = BigInt(object.creationHeight.toString());
        }
        if (object.completionTime !== undefined && object.completionTime !== null) {
            message.completionTime = timestamp_1.Timestamp.fromPartial(object.completionTime);
        }
        message.initialBalance = object.initialBalance ?? "";
        message.balance = object.balance ?? "";
        if (object.unbondingId !== undefined && object.unbondingId !== null) {
            message.unbondingId = BigInt(object.unbondingId.toString());
        }
        if (object.unbondingOnHoldRefCount !== undefined && object.unbondingOnHoldRefCount !== null) {
            message.unbondingOnHoldRefCount = BigInt(object.unbondingOnHoldRefCount.toString());
        }
        return message;
    },
};
function createBaseRedelegationEntry() {
    return {
        creationHeight: BigInt(0),
        completionTime: timestamp_1.Timestamp.fromPartial({}),
        initialBalance: "",
        sharesDst: "",
        unbondingId: BigInt(0),
        unbondingOnHoldRefCount: BigInt(0),
    };
}
exports.RedelegationEntry = {
    typeUrl: "/cosmos.staking.v1beta1.RedelegationEntry",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.creationHeight !== BigInt(0)) {
            writer.uint32(8).int64(message.creationHeight);
        }
        if (message.completionTime !== undefined) {
            timestamp_1.Timestamp.encode(message.completionTime, writer.uint32(18).fork()).ldelim();
        }
        if (message.initialBalance !== "") {
            writer.uint32(26).string(message.initialBalance);
        }
        if (message.sharesDst !== "") {
            writer.uint32(34).string(message.sharesDst);
        }
        if (message.unbondingId !== BigInt(0)) {
            writer.uint32(40).uint64(message.unbondingId);
        }
        if (message.unbondingOnHoldRefCount !== BigInt(0)) {
            writer.uint32(48).int64(message.unbondingOnHoldRefCount);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseRedelegationEntry();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.creationHeight = reader.int64();
                    break;
                case 2:
                    message.completionTime = timestamp_1.Timestamp.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.initialBalance = reader.string();
                    break;
                case 4:
                    message.sharesDst = reader.string();
                    break;
                case 5:
                    message.unbondingId = reader.uint64();
                    break;
                case 6:
                    message.unbondingOnHoldRefCount = reader.int64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseRedelegationEntry();
        if ((0, helpers_1.isSet)(object.creationHeight))
            obj.creationHeight = BigInt(object.creationHeight.toString());
        if ((0, helpers_1.isSet)(object.completionTime))
            obj.completionTime = (0, helpers_1.fromJsonTimestamp)(object.completionTime);
        if ((0, helpers_1.isSet)(object.initialBalance))
            obj.initialBalance = String(object.initialBalance);
        if ((0, helpers_1.isSet)(object.sharesDst))
            obj.sharesDst = String(object.sharesDst);
        if ((0, helpers_1.isSet)(object.unbondingId))
            obj.unbondingId = BigInt(object.unbondingId.toString());
        if ((0, helpers_1.isSet)(object.unbondingOnHoldRefCount))
            obj.unbondingOnHoldRefCount = BigInt(object.unbondingOnHoldRefCount.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.creationHeight !== undefined &&
            (obj.creationHeight = (message.creationHeight || BigInt(0)).toString());
        message.completionTime !== undefined &&
            (obj.completionTime = (0, helpers_1.fromTimestamp)(message.completionTime).toISOString());
        message.initialBalance !== undefined && (obj.initialBalance = message.initialBalance);
        message.sharesDst !== undefined && (obj.sharesDst = message.sharesDst);
        message.unbondingId !== undefined && (obj.unbondingId = (message.unbondingId || BigInt(0)).toString());
        message.unbondingOnHoldRefCount !== undefined &&
            (obj.unbondingOnHoldRefCount = (message.unbondingOnHoldRefCount || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseRedelegationEntry();
        if (object.creationHeight !== undefined && object.creationHeight !== null) {
            message.creationHeight = BigInt(object.creationHeight.toString());
        }
        if (object.completionTime !== undefined && object.completionTime !== null) {
            message.completionTime = timestamp_1.Timestamp.fromPartial(object.completionTime);
        }
        message.initialBalance = object.initialBalance ?? "";
        message.sharesDst = object.sharesDst ?? "";
        if (object.unbondingId !== undefined && object.unbondingId !== null) {
            message.unbondingId = BigInt(object.unbondingId.toString());
        }
        if (object.unbondingOnHoldRefCount !== undefined && object.unbondingOnHoldRefCount !== null) {
            message.unbondingOnHoldRefCount = BigInt(object.unbondingOnHoldRefCount.toString());
        }
        return message;
    },
};
function createBaseRedelegation() {
    return {
        delegatorAddress: "",
        validatorSrcAddress: "",
        validatorDstAddress: "",
        entries: [],
    };
}
exports.Redelegation = {
    typeUrl: "/cosmos.staking.v1beta1.Redelegation",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.delegatorAddress !== "") {
            writer.uint32(10).string(message.delegatorAddress);
        }
        if (message.validatorSrcAddress !== "") {
            writer.uint32(18).string(message.validatorSrcAddress);
        }
        if (message.validatorDstAddress !== "") {
            writer.uint32(26).string(message.validatorDstAddress);
        }
        for (const v of message.entries) {
            exports.RedelegationEntry.encode(v, writer.uint32(34).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseRedelegation();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.delegatorAddress = reader.string();
                    break;
                case 2:
                    message.validatorSrcAddress = reader.string();
                    break;
                case 3:
                    message.validatorDstAddress = reader.string();
                    break;
                case 4:
                    message.entries.push(exports.RedelegationEntry.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseRedelegation();
        if ((0, helpers_1.isSet)(object.delegatorAddress))
            obj.delegatorAddress = String(object.delegatorAddress);
        if ((0, helpers_1.isSet)(object.validatorSrcAddress))
            obj.validatorSrcAddress = String(object.validatorSrcAddress);
        if ((0, helpers_1.isSet)(object.validatorDstAddress))
            obj.validatorDstAddress = String(object.validatorDstAddress);
        if (Array.isArray(object?.entries))
            obj.entries = object.entries.map((e) => exports.RedelegationEntry.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.delegatorAddress !== undefined && (obj.delegatorAddress = message.delegatorAddress);
        message.validatorSrcAddress !== undefined && (obj.validatorSrcAddress = message.validatorSrcAddress);
        message.validatorDstAddress !== undefined && (obj.validatorDstAddress = message.validatorDstAddress);
        if (message.entries) {
            obj.entries = message.entries.map((e) => (e ? exports.RedelegationEntry.toJSON(e) : undefined));
        }
        else {
            obj.entries = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseRedelegation();
        message.delegatorAddress = object.delegatorAddress ?? "";
        message.validatorSrcAddress = object.validatorSrcAddress ?? "";
        message.validatorDstAddress = object.validatorDstAddress ?? "";
        message.entries = object.entries?.map((e) => exports.RedelegationEntry.fromPartial(e)) || [];
        return message;
    },
};
function createBaseParams() {
    return {
        unbondingTime: duration_1.Duration.fromPartial({}),
        maxValidators: 0,
        maxEntries: 0,
        historicalEntries: 0,
        bondDenom: "",
        minCommissionRate: "",
    };
}
exports.Params = {
    typeUrl: "/cosmos.staking.v1beta1.Params",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.unbondingTime !== undefined) {
            duration_1.Duration.encode(message.unbondingTime, writer.uint32(10).fork()).ldelim();
        }
        if (message.maxValidators !== 0) {
            writer.uint32(16).uint32(message.maxValidators);
        }
        if (message.maxEntries !== 0) {
            writer.uint32(24).uint32(message.maxEntries);
        }
        if (message.historicalEntries !== 0) {
            writer.uint32(32).uint32(message.historicalEntries);
        }
        if (message.bondDenom !== "") {
            writer.uint32(42).string(message.bondDenom);
        }
        if (message.minCommissionRate !== "") {
            writer.uint32(50).string(message.minCommissionRate);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseParams();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.unbondingTime = duration_1.Duration.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.maxValidators = reader.uint32();
                    break;
                case 3:
                    message.maxEntries = reader.uint32();
                    break;
                case 4:
                    message.historicalEntries = reader.uint32();
                    break;
                case 5:
                    message.bondDenom = reader.string();
                    break;
                case 6:
                    message.minCommissionRate = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseParams();
        if ((0, helpers_1.isSet)(object.unbondingTime))
            obj.unbondingTime = duration_1.Duration.fromJSON(object.unbondingTime);
        if ((0, helpers_1.isSet)(object.maxValidators))
            obj.maxValidators = Number(object.maxValidators);
        if ((0, helpers_1.isSet)(object.maxEntries))
            obj.maxEntries = Number(object.maxEntries);
        if ((0, helpers_1.isSet)(object.historicalEntries))
            obj.historicalEntries = Number(object.historicalEntries);
        if ((0, helpers_1.isSet)(object.bondDenom))
            obj.bondDenom = String(object.bondDenom);
        if ((0, helpers_1.isSet)(object.minCommissionRate))
            obj.minCommissionRate = String(object.minCommissionRate);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.unbondingTime !== undefined &&
            (obj.unbondingTime = message.unbondingTime ? duration_1.Duration.toJSON(message.unbondingTime) : undefined);
        message.maxValidators !== undefined && (obj.maxValidators = Math.round(message.maxValidators));
        message.maxEntries !== undefined && (obj.maxEntries = Math.round(message.maxEntries));
        message.historicalEntries !== undefined &&
            (obj.historicalEntries = Math.round(message.historicalEntries));
        message.bondDenom !== undefined && (obj.bondDenom = message.bondDenom);
        message.minCommissionRate !== undefined && (obj.minCommissionRate = message.minCommissionRate);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseParams();
        if (object.unbondingTime !== undefined && object.unbondingTime !== null) {
            message.unbondingTime = duration_1.Duration.fromPartial(object.unbondingTime);
        }
        message.maxValidators = object.maxValidators ?? 0;
        message.maxEntries = object.maxEntries ?? 0;
        message.historicalEntries = object.historicalEntries ?? 0;
        message.bondDenom = object.bondDenom ?? "";
        message.minCommissionRate = object.minCommissionRate ?? "";
        return message;
    },
};
function createBaseDelegationResponse() {
    return {
        delegation: exports.Delegation.fromPartial({}),
        balance: coin_1.Coin.fromPartial({}),
    };
}
exports.DelegationResponse = {
    typeUrl: "/cosmos.staking.v1beta1.DelegationResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.delegation !== undefined) {
            exports.Delegation.encode(message.delegation, writer.uint32(10).fork()).ldelim();
        }
        if (message.balance !== undefined) {
            coin_1.Coin.encode(message.balance, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseDelegationResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.delegation = exports.Delegation.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.balance = coin_1.Coin.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseDelegationResponse();
        if ((0, helpers_1.isSet)(object.delegation))
            obj.delegation = exports.Delegation.fromJSON(object.delegation);
        if ((0, helpers_1.isSet)(object.balance))
            obj.balance = coin_1.Coin.fromJSON(object.balance);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.delegation !== undefined &&
            (obj.delegation = message.delegation ? exports.Delegation.toJSON(message.delegation) : undefined);
        message.balance !== undefined &&
            (obj.balance = message.balance ? coin_1.Coin.toJSON(message.balance) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseDelegationResponse();
        if (object.delegation !== undefined && object.delegation !== null) {
            message.delegation = exports.Delegation.fromPartial(object.delegation);
        }
        if (object.balance !== undefined && object.balance !== null) {
            message.balance = coin_1.Coin.fromPartial(object.balance);
        }
        return message;
    },
};
function createBaseRedelegationEntryResponse() {
    return {
        redelegationEntry: exports.RedelegationEntry.fromPartial({}),
        balance: "",
    };
}
exports.RedelegationEntryResponse = {
    typeUrl: "/cosmos.staking.v1beta1.RedelegationEntryResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.redelegationEntry !== undefined) {
            exports.RedelegationEntry.encode(message.redelegationEntry, writer.uint32(10).fork()).ldelim();
        }
        if (message.balance !== "") {
            writer.uint32(34).string(message.balance);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseRedelegationEntryResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.redelegationEntry = exports.RedelegationEntry.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.balance = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseRedelegationEntryResponse();
        if ((0, helpers_1.isSet)(object.redelegationEntry))
            obj.redelegationEntry = exports.RedelegationEntry.fromJSON(object.redelegationEntry);
        if ((0, helpers_1.isSet)(object.balance))
            obj.balance = String(object.balance);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.redelegationEntry !== undefined &&
            (obj.redelegationEntry = message.redelegationEntry
                ? exports.RedelegationEntry.toJSON(message.redelegationEntry)
                : undefined);
        message.balance !== undefined && (obj.balance = message.balance);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseRedelegationEntryResponse();
        if (object.redelegationEntry !== undefined && object.redelegationEntry !== null) {
            message.redelegationEntry = exports.RedelegationEntry.fromPartial(object.redelegationEntry);
        }
        message.balance = object.balance ?? "";
        return message;
    },
};
function createBaseRedelegationResponse() {
    return {
        redelegation: exports.Redelegation.fromPartial({}),
        entries: [],
    };
}
exports.RedelegationResponse = {
    typeUrl: "/cosmos.staking.v1beta1.RedelegationResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.redelegation !== undefined) {
            exports.Redelegation.encode(message.redelegation, writer.uint32(10).fork()).ldelim();
        }
        for (const v of message.entries) {
            exports.RedelegationEntryResponse.encode(v, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseRedelegationResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.redelegation = exports.Redelegation.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.entries.push(exports.RedelegationEntryResponse.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseRedelegationResponse();
        if ((0, helpers_1.isSet)(object.redelegation))
            obj.redelegation = exports.Redelegation.fromJSON(object.redelegation);
        if (Array.isArray(object?.entries))
            obj.entries = object.entries.map((e) => exports.RedelegationEntryResponse.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.redelegation !== undefined &&
            (obj.redelegation = message.redelegation ? exports.Redelegation.toJSON(message.redelegation) : undefined);
        if (message.entries) {
            obj.entries = message.entries.map((e) => (e ? exports.RedelegationEntryResponse.toJSON(e) : undefined));
        }
        else {
            obj.entries = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseRedelegationResponse();
        if (object.redelegation !== undefined && object.redelegation !== null) {
            message.redelegation = exports.Redelegation.fromPartial(object.redelegation);
        }
        message.entries = object.entries?.map((e) => exports.RedelegationEntryResponse.fromPartial(e)) || [];
        return message;
    },
};
function createBasePool() {
    return {
        notBondedTokens: "",
        bondedTokens: "",
    };
}
exports.Pool = {
    typeUrl: "/cosmos.staking.v1beta1.Pool",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.notBondedTokens !== "") {
            writer.uint32(10).string(message.notBondedTokens);
        }
        if (message.bondedTokens !== "") {
            writer.uint32(18).string(message.bondedTokens);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBasePool();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.notBondedTokens = reader.string();
                    break;
                case 2:
                    message.bondedTokens = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBasePool();
        if ((0, helpers_1.isSet)(object.notBondedTokens))
            obj.notBondedTokens = String(object.notBondedTokens);
        if ((0, helpers_1.isSet)(object.bondedTokens))
            obj.bondedTokens = String(object.bondedTokens);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.notBondedTokens !== undefined && (obj.notBondedTokens = message.notBondedTokens);
        message.bondedTokens !== undefined && (obj.bondedTokens = message.bondedTokens);
        return obj;
    },
    fromPartial(object) {
        const message = createBasePool();
        message.notBondedTokens = object.notBondedTokens ?? "";
        message.bondedTokens = object.bondedTokens ?? "";
        return message;
    },
};
function createBaseValidatorUpdates() {
    return {
        updates: [],
    };
}
exports.ValidatorUpdates = {
    typeUrl: "/cosmos.staking.v1beta1.ValidatorUpdates",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.updates) {
            types_2.ValidatorUpdate.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseValidatorUpdates();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.updates.push(types_2.ValidatorUpdate.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseValidatorUpdates();
        if (Array.isArray(object?.updates))
            obj.updates = object.updates.map((e) => types_2.ValidatorUpdate.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.updates) {
            obj.updates = message.updates.map((e) => (e ? types_2.ValidatorUpdate.toJSON(e) : undefined));
        }
        else {
            obj.updates = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseValidatorUpdates();
        message.updates = object.updates?.map((e) => types_2.ValidatorUpdate.fromPartial(e)) || [];
        return message;
    },
};
//# sourceMappingURL=staking.js.map