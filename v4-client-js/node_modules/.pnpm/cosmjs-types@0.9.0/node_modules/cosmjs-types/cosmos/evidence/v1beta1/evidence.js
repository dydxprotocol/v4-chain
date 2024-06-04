"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Equivocation = exports.protobufPackage = void 0;
/* eslint-disable */
const timestamp_1 = require("../../../google/protobuf/timestamp");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.evidence.v1beta1";
function createBaseEquivocation() {
    return {
        height: BigInt(0),
        time: timestamp_1.Timestamp.fromPartial({}),
        power: BigInt(0),
        consensusAddress: "",
    };
}
exports.Equivocation = {
    typeUrl: "/cosmos.evidence.v1beta1.Equivocation",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.height !== BigInt(0)) {
            writer.uint32(8).int64(message.height);
        }
        if (message.time !== undefined) {
            timestamp_1.Timestamp.encode(message.time, writer.uint32(18).fork()).ldelim();
        }
        if (message.power !== BigInt(0)) {
            writer.uint32(24).int64(message.power);
        }
        if (message.consensusAddress !== "") {
            writer.uint32(34).string(message.consensusAddress);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseEquivocation();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.height = reader.int64();
                    break;
                case 2:
                    message.time = timestamp_1.Timestamp.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.power = reader.int64();
                    break;
                case 4:
                    message.consensusAddress = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseEquivocation();
        if ((0, helpers_1.isSet)(object.height))
            obj.height = BigInt(object.height.toString());
        if ((0, helpers_1.isSet)(object.time))
            obj.time = (0, helpers_1.fromJsonTimestamp)(object.time);
        if ((0, helpers_1.isSet)(object.power))
            obj.power = BigInt(object.power.toString());
        if ((0, helpers_1.isSet)(object.consensusAddress))
            obj.consensusAddress = String(object.consensusAddress);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.height !== undefined && (obj.height = (message.height || BigInt(0)).toString());
        message.time !== undefined && (obj.time = (0, helpers_1.fromTimestamp)(message.time).toISOString());
        message.power !== undefined && (obj.power = (message.power || BigInt(0)).toString());
        message.consensusAddress !== undefined && (obj.consensusAddress = message.consensusAddress);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseEquivocation();
        if (object.height !== undefined && object.height !== null) {
            message.height = BigInt(object.height.toString());
        }
        if (object.time !== undefined && object.time !== null) {
            message.time = timestamp_1.Timestamp.fromPartial(object.time);
        }
        if (object.power !== undefined && object.power !== null) {
            message.power = BigInt(object.power.toString());
        }
        message.consensusAddress = object.consensusAddress ?? "";
        return message;
    },
};
//# sourceMappingURL=evidence.js.map