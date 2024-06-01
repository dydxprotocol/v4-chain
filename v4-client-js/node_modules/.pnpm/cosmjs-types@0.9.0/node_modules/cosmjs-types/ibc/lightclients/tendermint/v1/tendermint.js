"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Fraction = exports.Header = exports.Misbehaviour = exports.ConsensusState = exports.ClientState = exports.protobufPackage = void 0;
/* eslint-disable */
const duration_1 = require("../../../../google/protobuf/duration");
const client_1 = require("../../../core/client/v1/client");
const proofs_1 = require("../../../../cosmos/ics23/v1/proofs");
const timestamp_1 = require("../../../../google/protobuf/timestamp");
const commitment_1 = require("../../../core/commitment/v1/commitment");
const types_1 = require("../../../../tendermint/types/types");
const validator_1 = require("../../../../tendermint/types/validator");
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "ibc.lightclients.tendermint.v1";
function createBaseClientState() {
    return {
        chainId: "",
        trustLevel: exports.Fraction.fromPartial({}),
        trustingPeriod: duration_1.Duration.fromPartial({}),
        unbondingPeriod: duration_1.Duration.fromPartial({}),
        maxClockDrift: duration_1.Duration.fromPartial({}),
        frozenHeight: client_1.Height.fromPartial({}),
        latestHeight: client_1.Height.fromPartial({}),
        proofSpecs: [],
        upgradePath: [],
        allowUpdateAfterExpiry: false,
        allowUpdateAfterMisbehaviour: false,
    };
}
exports.ClientState = {
    typeUrl: "/ibc.lightclients.tendermint.v1.ClientState",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.chainId !== "") {
            writer.uint32(10).string(message.chainId);
        }
        if (message.trustLevel !== undefined) {
            exports.Fraction.encode(message.trustLevel, writer.uint32(18).fork()).ldelim();
        }
        if (message.trustingPeriod !== undefined) {
            duration_1.Duration.encode(message.trustingPeriod, writer.uint32(26).fork()).ldelim();
        }
        if (message.unbondingPeriod !== undefined) {
            duration_1.Duration.encode(message.unbondingPeriod, writer.uint32(34).fork()).ldelim();
        }
        if (message.maxClockDrift !== undefined) {
            duration_1.Duration.encode(message.maxClockDrift, writer.uint32(42).fork()).ldelim();
        }
        if (message.frozenHeight !== undefined) {
            client_1.Height.encode(message.frozenHeight, writer.uint32(50).fork()).ldelim();
        }
        if (message.latestHeight !== undefined) {
            client_1.Height.encode(message.latestHeight, writer.uint32(58).fork()).ldelim();
        }
        for (const v of message.proofSpecs) {
            proofs_1.ProofSpec.encode(v, writer.uint32(66).fork()).ldelim();
        }
        for (const v of message.upgradePath) {
            writer.uint32(74).string(v);
        }
        if (message.allowUpdateAfterExpiry === true) {
            writer.uint32(80).bool(message.allowUpdateAfterExpiry);
        }
        if (message.allowUpdateAfterMisbehaviour === true) {
            writer.uint32(88).bool(message.allowUpdateAfterMisbehaviour);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseClientState();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.chainId = reader.string();
                    break;
                case 2:
                    message.trustLevel = exports.Fraction.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.trustingPeriod = duration_1.Duration.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.unbondingPeriod = duration_1.Duration.decode(reader, reader.uint32());
                    break;
                case 5:
                    message.maxClockDrift = duration_1.Duration.decode(reader, reader.uint32());
                    break;
                case 6:
                    message.frozenHeight = client_1.Height.decode(reader, reader.uint32());
                    break;
                case 7:
                    message.latestHeight = client_1.Height.decode(reader, reader.uint32());
                    break;
                case 8:
                    message.proofSpecs.push(proofs_1.ProofSpec.decode(reader, reader.uint32()));
                    break;
                case 9:
                    message.upgradePath.push(reader.string());
                    break;
                case 10:
                    message.allowUpdateAfterExpiry = reader.bool();
                    break;
                case 11:
                    message.allowUpdateAfterMisbehaviour = reader.bool();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseClientState();
        if ((0, helpers_1.isSet)(object.chainId))
            obj.chainId = String(object.chainId);
        if ((0, helpers_1.isSet)(object.trustLevel))
            obj.trustLevel = exports.Fraction.fromJSON(object.trustLevel);
        if ((0, helpers_1.isSet)(object.trustingPeriod))
            obj.trustingPeriod = duration_1.Duration.fromJSON(object.trustingPeriod);
        if ((0, helpers_1.isSet)(object.unbondingPeriod))
            obj.unbondingPeriod = duration_1.Duration.fromJSON(object.unbondingPeriod);
        if ((0, helpers_1.isSet)(object.maxClockDrift))
            obj.maxClockDrift = duration_1.Duration.fromJSON(object.maxClockDrift);
        if ((0, helpers_1.isSet)(object.frozenHeight))
            obj.frozenHeight = client_1.Height.fromJSON(object.frozenHeight);
        if ((0, helpers_1.isSet)(object.latestHeight))
            obj.latestHeight = client_1.Height.fromJSON(object.latestHeight);
        if (Array.isArray(object?.proofSpecs))
            obj.proofSpecs = object.proofSpecs.map((e) => proofs_1.ProofSpec.fromJSON(e));
        if (Array.isArray(object?.upgradePath))
            obj.upgradePath = object.upgradePath.map((e) => String(e));
        if ((0, helpers_1.isSet)(object.allowUpdateAfterExpiry))
            obj.allowUpdateAfterExpiry = Boolean(object.allowUpdateAfterExpiry);
        if ((0, helpers_1.isSet)(object.allowUpdateAfterMisbehaviour))
            obj.allowUpdateAfterMisbehaviour = Boolean(object.allowUpdateAfterMisbehaviour);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.chainId !== undefined && (obj.chainId = message.chainId);
        message.trustLevel !== undefined &&
            (obj.trustLevel = message.trustLevel ? exports.Fraction.toJSON(message.trustLevel) : undefined);
        message.trustingPeriod !== undefined &&
            (obj.trustingPeriod = message.trustingPeriod ? duration_1.Duration.toJSON(message.trustingPeriod) : undefined);
        message.unbondingPeriod !== undefined &&
            (obj.unbondingPeriod = message.unbondingPeriod ? duration_1.Duration.toJSON(message.unbondingPeriod) : undefined);
        message.maxClockDrift !== undefined &&
            (obj.maxClockDrift = message.maxClockDrift ? duration_1.Duration.toJSON(message.maxClockDrift) : undefined);
        message.frozenHeight !== undefined &&
            (obj.frozenHeight = message.frozenHeight ? client_1.Height.toJSON(message.frozenHeight) : undefined);
        message.latestHeight !== undefined &&
            (obj.latestHeight = message.latestHeight ? client_1.Height.toJSON(message.latestHeight) : undefined);
        if (message.proofSpecs) {
            obj.proofSpecs = message.proofSpecs.map((e) => (e ? proofs_1.ProofSpec.toJSON(e) : undefined));
        }
        else {
            obj.proofSpecs = [];
        }
        if (message.upgradePath) {
            obj.upgradePath = message.upgradePath.map((e) => e);
        }
        else {
            obj.upgradePath = [];
        }
        message.allowUpdateAfterExpiry !== undefined &&
            (obj.allowUpdateAfterExpiry = message.allowUpdateAfterExpiry);
        message.allowUpdateAfterMisbehaviour !== undefined &&
            (obj.allowUpdateAfterMisbehaviour = message.allowUpdateAfterMisbehaviour);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseClientState();
        message.chainId = object.chainId ?? "";
        if (object.trustLevel !== undefined && object.trustLevel !== null) {
            message.trustLevel = exports.Fraction.fromPartial(object.trustLevel);
        }
        if (object.trustingPeriod !== undefined && object.trustingPeriod !== null) {
            message.trustingPeriod = duration_1.Duration.fromPartial(object.trustingPeriod);
        }
        if (object.unbondingPeriod !== undefined && object.unbondingPeriod !== null) {
            message.unbondingPeriod = duration_1.Duration.fromPartial(object.unbondingPeriod);
        }
        if (object.maxClockDrift !== undefined && object.maxClockDrift !== null) {
            message.maxClockDrift = duration_1.Duration.fromPartial(object.maxClockDrift);
        }
        if (object.frozenHeight !== undefined && object.frozenHeight !== null) {
            message.frozenHeight = client_1.Height.fromPartial(object.frozenHeight);
        }
        if (object.latestHeight !== undefined && object.latestHeight !== null) {
            message.latestHeight = client_1.Height.fromPartial(object.latestHeight);
        }
        message.proofSpecs = object.proofSpecs?.map((e) => proofs_1.ProofSpec.fromPartial(e)) || [];
        message.upgradePath = object.upgradePath?.map((e) => e) || [];
        message.allowUpdateAfterExpiry = object.allowUpdateAfterExpiry ?? false;
        message.allowUpdateAfterMisbehaviour = object.allowUpdateAfterMisbehaviour ?? false;
        return message;
    },
};
function createBaseConsensusState() {
    return {
        timestamp: timestamp_1.Timestamp.fromPartial({}),
        root: commitment_1.MerkleRoot.fromPartial({}),
        nextValidatorsHash: new Uint8Array(),
    };
}
exports.ConsensusState = {
    typeUrl: "/ibc.lightclients.tendermint.v1.ConsensusState",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.timestamp !== undefined) {
            timestamp_1.Timestamp.encode(message.timestamp, writer.uint32(10).fork()).ldelim();
        }
        if (message.root !== undefined) {
            commitment_1.MerkleRoot.encode(message.root, writer.uint32(18).fork()).ldelim();
        }
        if (message.nextValidatorsHash.length !== 0) {
            writer.uint32(26).bytes(message.nextValidatorsHash);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseConsensusState();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.timestamp = timestamp_1.Timestamp.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.root = commitment_1.MerkleRoot.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.nextValidatorsHash = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseConsensusState();
        if ((0, helpers_1.isSet)(object.timestamp))
            obj.timestamp = (0, helpers_1.fromJsonTimestamp)(object.timestamp);
        if ((0, helpers_1.isSet)(object.root))
            obj.root = commitment_1.MerkleRoot.fromJSON(object.root);
        if ((0, helpers_1.isSet)(object.nextValidatorsHash))
            obj.nextValidatorsHash = (0, helpers_1.bytesFromBase64)(object.nextValidatorsHash);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.timestamp !== undefined && (obj.timestamp = (0, helpers_1.fromTimestamp)(message.timestamp).toISOString());
        message.root !== undefined && (obj.root = message.root ? commitment_1.MerkleRoot.toJSON(message.root) : undefined);
        message.nextValidatorsHash !== undefined &&
            (obj.nextValidatorsHash = (0, helpers_1.base64FromBytes)(message.nextValidatorsHash !== undefined ? message.nextValidatorsHash : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseConsensusState();
        if (object.timestamp !== undefined && object.timestamp !== null) {
            message.timestamp = timestamp_1.Timestamp.fromPartial(object.timestamp);
        }
        if (object.root !== undefined && object.root !== null) {
            message.root = commitment_1.MerkleRoot.fromPartial(object.root);
        }
        message.nextValidatorsHash = object.nextValidatorsHash ?? new Uint8Array();
        return message;
    },
};
function createBaseMisbehaviour() {
    return {
        clientId: "",
        header1: undefined,
        header2: undefined,
    };
}
exports.Misbehaviour = {
    typeUrl: "/ibc.lightclients.tendermint.v1.Misbehaviour",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.clientId !== "") {
            writer.uint32(10).string(message.clientId);
        }
        if (message.header1 !== undefined) {
            exports.Header.encode(message.header1, writer.uint32(18).fork()).ldelim();
        }
        if (message.header2 !== undefined) {
            exports.Header.encode(message.header2, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMisbehaviour();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.clientId = reader.string();
                    break;
                case 2:
                    message.header1 = exports.Header.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.header2 = exports.Header.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMisbehaviour();
        if ((0, helpers_1.isSet)(object.clientId))
            obj.clientId = String(object.clientId);
        if ((0, helpers_1.isSet)(object.header1))
            obj.header1 = exports.Header.fromJSON(object.header1);
        if ((0, helpers_1.isSet)(object.header2))
            obj.header2 = exports.Header.fromJSON(object.header2);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.clientId !== undefined && (obj.clientId = message.clientId);
        message.header1 !== undefined &&
            (obj.header1 = message.header1 ? exports.Header.toJSON(message.header1) : undefined);
        message.header2 !== undefined &&
            (obj.header2 = message.header2 ? exports.Header.toJSON(message.header2) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMisbehaviour();
        message.clientId = object.clientId ?? "";
        if (object.header1 !== undefined && object.header1 !== null) {
            message.header1 = exports.Header.fromPartial(object.header1);
        }
        if (object.header2 !== undefined && object.header2 !== null) {
            message.header2 = exports.Header.fromPartial(object.header2);
        }
        return message;
    },
};
function createBaseHeader() {
    return {
        signedHeader: undefined,
        validatorSet: undefined,
        trustedHeight: client_1.Height.fromPartial({}),
        trustedValidators: undefined,
    };
}
exports.Header = {
    typeUrl: "/ibc.lightclients.tendermint.v1.Header",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.signedHeader !== undefined) {
            types_1.SignedHeader.encode(message.signedHeader, writer.uint32(10).fork()).ldelim();
        }
        if (message.validatorSet !== undefined) {
            validator_1.ValidatorSet.encode(message.validatorSet, writer.uint32(18).fork()).ldelim();
        }
        if (message.trustedHeight !== undefined) {
            client_1.Height.encode(message.trustedHeight, writer.uint32(26).fork()).ldelim();
        }
        if (message.trustedValidators !== undefined) {
            validator_1.ValidatorSet.encode(message.trustedValidators, writer.uint32(34).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseHeader();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.signedHeader = types_1.SignedHeader.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.validatorSet = validator_1.ValidatorSet.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.trustedHeight = client_1.Height.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.trustedValidators = validator_1.ValidatorSet.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseHeader();
        if ((0, helpers_1.isSet)(object.signedHeader))
            obj.signedHeader = types_1.SignedHeader.fromJSON(object.signedHeader);
        if ((0, helpers_1.isSet)(object.validatorSet))
            obj.validatorSet = validator_1.ValidatorSet.fromJSON(object.validatorSet);
        if ((0, helpers_1.isSet)(object.trustedHeight))
            obj.trustedHeight = client_1.Height.fromJSON(object.trustedHeight);
        if ((0, helpers_1.isSet)(object.trustedValidators))
            obj.trustedValidators = validator_1.ValidatorSet.fromJSON(object.trustedValidators);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.signedHeader !== undefined &&
            (obj.signedHeader = message.signedHeader ? types_1.SignedHeader.toJSON(message.signedHeader) : undefined);
        message.validatorSet !== undefined &&
            (obj.validatorSet = message.validatorSet ? validator_1.ValidatorSet.toJSON(message.validatorSet) : undefined);
        message.trustedHeight !== undefined &&
            (obj.trustedHeight = message.trustedHeight ? client_1.Height.toJSON(message.trustedHeight) : undefined);
        message.trustedValidators !== undefined &&
            (obj.trustedValidators = message.trustedValidators
                ? validator_1.ValidatorSet.toJSON(message.trustedValidators)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseHeader();
        if (object.signedHeader !== undefined && object.signedHeader !== null) {
            message.signedHeader = types_1.SignedHeader.fromPartial(object.signedHeader);
        }
        if (object.validatorSet !== undefined && object.validatorSet !== null) {
            message.validatorSet = validator_1.ValidatorSet.fromPartial(object.validatorSet);
        }
        if (object.trustedHeight !== undefined && object.trustedHeight !== null) {
            message.trustedHeight = client_1.Height.fromPartial(object.trustedHeight);
        }
        if (object.trustedValidators !== undefined && object.trustedValidators !== null) {
            message.trustedValidators = validator_1.ValidatorSet.fromPartial(object.trustedValidators);
        }
        return message;
    },
};
function createBaseFraction() {
    return {
        numerator: BigInt(0),
        denominator: BigInt(0),
    };
}
exports.Fraction = {
    typeUrl: "/ibc.lightclients.tendermint.v1.Fraction",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.numerator !== BigInt(0)) {
            writer.uint32(8).uint64(message.numerator);
        }
        if (message.denominator !== BigInt(0)) {
            writer.uint32(16).uint64(message.denominator);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseFraction();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.numerator = reader.uint64();
                    break;
                case 2:
                    message.denominator = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseFraction();
        if ((0, helpers_1.isSet)(object.numerator))
            obj.numerator = BigInt(object.numerator.toString());
        if ((0, helpers_1.isSet)(object.denominator))
            obj.denominator = BigInt(object.denominator.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.numerator !== undefined && (obj.numerator = (message.numerator || BigInt(0)).toString());
        message.denominator !== undefined && (obj.denominator = (message.denominator || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseFraction();
        if (object.numerator !== undefined && object.numerator !== null) {
            message.numerator = BigInt(object.numerator.toString());
        }
        if (object.denominator !== undefined && object.denominator !== null) {
            message.denominator = BigInt(object.denominator.toString());
        }
        return message;
    },
};
//# sourceMappingURL=tendermint.js.map