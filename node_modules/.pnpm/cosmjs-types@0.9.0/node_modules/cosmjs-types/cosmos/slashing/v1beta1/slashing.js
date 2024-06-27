"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Params = exports.ValidatorSigningInfo = exports.protobufPackage = void 0;
/* eslint-disable */
const timestamp_1 = require("../../../google/protobuf/timestamp");
const duration_1 = require("../../../google/protobuf/duration");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.slashing.v1beta1";
function createBaseValidatorSigningInfo() {
    return {
        address: "",
        startHeight: BigInt(0),
        indexOffset: BigInt(0),
        jailedUntil: timestamp_1.Timestamp.fromPartial({}),
        tombstoned: false,
        missedBlocksCounter: BigInt(0),
    };
}
exports.ValidatorSigningInfo = {
    typeUrl: "/cosmos.slashing.v1beta1.ValidatorSigningInfo",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.address !== "") {
            writer.uint32(10).string(message.address);
        }
        if (message.startHeight !== BigInt(0)) {
            writer.uint32(16).int64(message.startHeight);
        }
        if (message.indexOffset !== BigInt(0)) {
            writer.uint32(24).int64(message.indexOffset);
        }
        if (message.jailedUntil !== undefined) {
            timestamp_1.Timestamp.encode(message.jailedUntil, writer.uint32(34).fork()).ldelim();
        }
        if (message.tombstoned === true) {
            writer.uint32(40).bool(message.tombstoned);
        }
        if (message.missedBlocksCounter !== BigInt(0)) {
            writer.uint32(48).int64(message.missedBlocksCounter);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseValidatorSigningInfo();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.address = reader.string();
                    break;
                case 2:
                    message.startHeight = reader.int64();
                    break;
                case 3:
                    message.indexOffset = reader.int64();
                    break;
                case 4:
                    message.jailedUntil = timestamp_1.Timestamp.decode(reader, reader.uint32());
                    break;
                case 5:
                    message.tombstoned = reader.bool();
                    break;
                case 6:
                    message.missedBlocksCounter = reader.int64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseValidatorSigningInfo();
        if ((0, helpers_1.isSet)(object.address))
            obj.address = String(object.address);
        if ((0, helpers_1.isSet)(object.startHeight))
            obj.startHeight = BigInt(object.startHeight.toString());
        if ((0, helpers_1.isSet)(object.indexOffset))
            obj.indexOffset = BigInt(object.indexOffset.toString());
        if ((0, helpers_1.isSet)(object.jailedUntil))
            obj.jailedUntil = (0, helpers_1.fromJsonTimestamp)(object.jailedUntil);
        if ((0, helpers_1.isSet)(object.tombstoned))
            obj.tombstoned = Boolean(object.tombstoned);
        if ((0, helpers_1.isSet)(object.missedBlocksCounter))
            obj.missedBlocksCounter = BigInt(object.missedBlocksCounter.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.address !== undefined && (obj.address = message.address);
        message.startHeight !== undefined && (obj.startHeight = (message.startHeight || BigInt(0)).toString());
        message.indexOffset !== undefined && (obj.indexOffset = (message.indexOffset || BigInt(0)).toString());
        message.jailedUntil !== undefined && (obj.jailedUntil = (0, helpers_1.fromTimestamp)(message.jailedUntil).toISOString());
        message.tombstoned !== undefined && (obj.tombstoned = message.tombstoned);
        message.missedBlocksCounter !== undefined &&
            (obj.missedBlocksCounter = (message.missedBlocksCounter || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseValidatorSigningInfo();
        message.address = object.address ?? "";
        if (object.startHeight !== undefined && object.startHeight !== null) {
            message.startHeight = BigInt(object.startHeight.toString());
        }
        if (object.indexOffset !== undefined && object.indexOffset !== null) {
            message.indexOffset = BigInt(object.indexOffset.toString());
        }
        if (object.jailedUntil !== undefined && object.jailedUntil !== null) {
            message.jailedUntil = timestamp_1.Timestamp.fromPartial(object.jailedUntil);
        }
        message.tombstoned = object.tombstoned ?? false;
        if (object.missedBlocksCounter !== undefined && object.missedBlocksCounter !== null) {
            message.missedBlocksCounter = BigInt(object.missedBlocksCounter.toString());
        }
        return message;
    },
};
function createBaseParams() {
    return {
        signedBlocksWindow: BigInt(0),
        minSignedPerWindow: new Uint8Array(),
        downtimeJailDuration: duration_1.Duration.fromPartial({}),
        slashFractionDoubleSign: new Uint8Array(),
        slashFractionDowntime: new Uint8Array(),
    };
}
exports.Params = {
    typeUrl: "/cosmos.slashing.v1beta1.Params",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.signedBlocksWindow !== BigInt(0)) {
            writer.uint32(8).int64(message.signedBlocksWindow);
        }
        if (message.minSignedPerWindow.length !== 0) {
            writer.uint32(18).bytes(message.minSignedPerWindow);
        }
        if (message.downtimeJailDuration !== undefined) {
            duration_1.Duration.encode(message.downtimeJailDuration, writer.uint32(26).fork()).ldelim();
        }
        if (message.slashFractionDoubleSign.length !== 0) {
            writer.uint32(34).bytes(message.slashFractionDoubleSign);
        }
        if (message.slashFractionDowntime.length !== 0) {
            writer.uint32(42).bytes(message.slashFractionDowntime);
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
                    message.signedBlocksWindow = reader.int64();
                    break;
                case 2:
                    message.minSignedPerWindow = reader.bytes();
                    break;
                case 3:
                    message.downtimeJailDuration = duration_1.Duration.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.slashFractionDoubleSign = reader.bytes();
                    break;
                case 5:
                    message.slashFractionDowntime = reader.bytes();
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
        if ((0, helpers_1.isSet)(object.signedBlocksWindow))
            obj.signedBlocksWindow = BigInt(object.signedBlocksWindow.toString());
        if ((0, helpers_1.isSet)(object.minSignedPerWindow))
            obj.minSignedPerWindow = (0, helpers_1.bytesFromBase64)(object.minSignedPerWindow);
        if ((0, helpers_1.isSet)(object.downtimeJailDuration))
            obj.downtimeJailDuration = duration_1.Duration.fromJSON(object.downtimeJailDuration);
        if ((0, helpers_1.isSet)(object.slashFractionDoubleSign))
            obj.slashFractionDoubleSign = (0, helpers_1.bytesFromBase64)(object.slashFractionDoubleSign);
        if ((0, helpers_1.isSet)(object.slashFractionDowntime))
            obj.slashFractionDowntime = (0, helpers_1.bytesFromBase64)(object.slashFractionDowntime);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.signedBlocksWindow !== undefined &&
            (obj.signedBlocksWindow = (message.signedBlocksWindow || BigInt(0)).toString());
        message.minSignedPerWindow !== undefined &&
            (obj.minSignedPerWindow = (0, helpers_1.base64FromBytes)(message.minSignedPerWindow !== undefined ? message.minSignedPerWindow : new Uint8Array()));
        message.downtimeJailDuration !== undefined &&
            (obj.downtimeJailDuration = message.downtimeJailDuration
                ? duration_1.Duration.toJSON(message.downtimeJailDuration)
                : undefined);
        message.slashFractionDoubleSign !== undefined &&
            (obj.slashFractionDoubleSign = (0, helpers_1.base64FromBytes)(message.slashFractionDoubleSign !== undefined ? message.slashFractionDoubleSign : new Uint8Array()));
        message.slashFractionDowntime !== undefined &&
            (obj.slashFractionDowntime = (0, helpers_1.base64FromBytes)(message.slashFractionDowntime !== undefined ? message.slashFractionDowntime : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseParams();
        if (object.signedBlocksWindow !== undefined && object.signedBlocksWindow !== null) {
            message.signedBlocksWindow = BigInt(object.signedBlocksWindow.toString());
        }
        message.minSignedPerWindow = object.minSignedPerWindow ?? new Uint8Array();
        if (object.downtimeJailDuration !== undefined && object.downtimeJailDuration !== null) {
            message.downtimeJailDuration = duration_1.Duration.fromPartial(object.downtimeJailDuration);
        }
        message.slashFractionDoubleSign = object.slashFractionDoubleSign ?? new Uint8Array();
        message.slashFractionDowntime = object.slashFractionDowntime ?? new Uint8Array();
        return message;
    },
};
//# sourceMappingURL=slashing.js.map