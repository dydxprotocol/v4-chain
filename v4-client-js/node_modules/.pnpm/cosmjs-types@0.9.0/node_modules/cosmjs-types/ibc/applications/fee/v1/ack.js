"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.IncentivizedAcknowledgement = exports.protobufPackage = void 0;
/* eslint-disable */
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "ibc.applications.fee.v1";
function createBaseIncentivizedAcknowledgement() {
    return {
        appAcknowledgement: new Uint8Array(),
        forwardRelayerAddress: "",
        underlyingAppSuccess: false,
    };
}
exports.IncentivizedAcknowledgement = {
    typeUrl: "/ibc.applications.fee.v1.IncentivizedAcknowledgement",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.appAcknowledgement.length !== 0) {
            writer.uint32(10).bytes(message.appAcknowledgement);
        }
        if (message.forwardRelayerAddress !== "") {
            writer.uint32(18).string(message.forwardRelayerAddress);
        }
        if (message.underlyingAppSuccess === true) {
            writer.uint32(24).bool(message.underlyingAppSuccess);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseIncentivizedAcknowledgement();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.appAcknowledgement = reader.bytes();
                    break;
                case 2:
                    message.forwardRelayerAddress = reader.string();
                    break;
                case 3:
                    message.underlyingAppSuccess = reader.bool();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseIncentivizedAcknowledgement();
        if ((0, helpers_1.isSet)(object.appAcknowledgement))
            obj.appAcknowledgement = (0, helpers_1.bytesFromBase64)(object.appAcknowledgement);
        if ((0, helpers_1.isSet)(object.forwardRelayerAddress))
            obj.forwardRelayerAddress = String(object.forwardRelayerAddress);
        if ((0, helpers_1.isSet)(object.underlyingAppSuccess))
            obj.underlyingAppSuccess = Boolean(object.underlyingAppSuccess);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.appAcknowledgement !== undefined &&
            (obj.appAcknowledgement = (0, helpers_1.base64FromBytes)(message.appAcknowledgement !== undefined ? message.appAcknowledgement : new Uint8Array()));
        message.forwardRelayerAddress !== undefined &&
            (obj.forwardRelayerAddress = message.forwardRelayerAddress);
        message.underlyingAppSuccess !== undefined && (obj.underlyingAppSuccess = message.underlyingAppSuccess);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseIncentivizedAcknowledgement();
        message.appAcknowledgement = object.appAcknowledgement ?? new Uint8Array();
        message.forwardRelayerAddress = object.forwardRelayerAddress ?? "";
        message.underlyingAppSuccess = object.underlyingAppSuccess ?? false;
        return message;
    },
};
//# sourceMappingURL=ack.js.map