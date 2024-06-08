"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Params = exports.Minter = exports.protobufPackage = void 0;
/* eslint-disable */
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.mint.v1beta1";
function createBaseMinter() {
    return {
        inflation: "",
        annualProvisions: "",
    };
}
exports.Minter = {
    typeUrl: "/cosmos.mint.v1beta1.Minter",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.inflation !== "") {
            writer.uint32(10).string(message.inflation);
        }
        if (message.annualProvisions !== "") {
            writer.uint32(18).string(message.annualProvisions);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMinter();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.inflation = reader.string();
                    break;
                case 2:
                    message.annualProvisions = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMinter();
        if ((0, helpers_1.isSet)(object.inflation))
            obj.inflation = String(object.inflation);
        if ((0, helpers_1.isSet)(object.annualProvisions))
            obj.annualProvisions = String(object.annualProvisions);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.inflation !== undefined && (obj.inflation = message.inflation);
        message.annualProvisions !== undefined && (obj.annualProvisions = message.annualProvisions);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMinter();
        message.inflation = object.inflation ?? "";
        message.annualProvisions = object.annualProvisions ?? "";
        return message;
    },
};
function createBaseParams() {
    return {
        mintDenom: "",
        inflationRateChange: "",
        inflationMax: "",
        inflationMin: "",
        goalBonded: "",
        blocksPerYear: BigInt(0),
    };
}
exports.Params = {
    typeUrl: "/cosmos.mint.v1beta1.Params",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.mintDenom !== "") {
            writer.uint32(10).string(message.mintDenom);
        }
        if (message.inflationRateChange !== "") {
            writer.uint32(18).string(message.inflationRateChange);
        }
        if (message.inflationMax !== "") {
            writer.uint32(26).string(message.inflationMax);
        }
        if (message.inflationMin !== "") {
            writer.uint32(34).string(message.inflationMin);
        }
        if (message.goalBonded !== "") {
            writer.uint32(42).string(message.goalBonded);
        }
        if (message.blocksPerYear !== BigInt(0)) {
            writer.uint32(48).uint64(message.blocksPerYear);
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
                    message.mintDenom = reader.string();
                    break;
                case 2:
                    message.inflationRateChange = reader.string();
                    break;
                case 3:
                    message.inflationMax = reader.string();
                    break;
                case 4:
                    message.inflationMin = reader.string();
                    break;
                case 5:
                    message.goalBonded = reader.string();
                    break;
                case 6:
                    message.blocksPerYear = reader.uint64();
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
        if ((0, helpers_1.isSet)(object.mintDenom))
            obj.mintDenom = String(object.mintDenom);
        if ((0, helpers_1.isSet)(object.inflationRateChange))
            obj.inflationRateChange = String(object.inflationRateChange);
        if ((0, helpers_1.isSet)(object.inflationMax))
            obj.inflationMax = String(object.inflationMax);
        if ((0, helpers_1.isSet)(object.inflationMin))
            obj.inflationMin = String(object.inflationMin);
        if ((0, helpers_1.isSet)(object.goalBonded))
            obj.goalBonded = String(object.goalBonded);
        if ((0, helpers_1.isSet)(object.blocksPerYear))
            obj.blocksPerYear = BigInt(object.blocksPerYear.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.mintDenom !== undefined && (obj.mintDenom = message.mintDenom);
        message.inflationRateChange !== undefined && (obj.inflationRateChange = message.inflationRateChange);
        message.inflationMax !== undefined && (obj.inflationMax = message.inflationMax);
        message.inflationMin !== undefined && (obj.inflationMin = message.inflationMin);
        message.goalBonded !== undefined && (obj.goalBonded = message.goalBonded);
        message.blocksPerYear !== undefined &&
            (obj.blocksPerYear = (message.blocksPerYear || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseParams();
        message.mintDenom = object.mintDenom ?? "";
        message.inflationRateChange = object.inflationRateChange ?? "";
        message.inflationMax = object.inflationMax ?? "";
        message.inflationMin = object.inflationMin ?? "";
        message.goalBonded = object.goalBonded ?? "";
        if (object.blocksPerYear !== undefined && object.blocksPerYear !== null) {
            message.blocksPerYear = BigInt(object.blocksPerYear.toString());
        }
        return message;
    },
};
//# sourceMappingURL=mint.js.map