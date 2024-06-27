"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.decodeOptionalPubkey = exports.decodePubkey = exports.anyToSinglePubkey = exports.encodePubkey = void 0;
/* eslint-disable @typescript-eslint/naming-convention */
const amino_1 = require("@cosmjs/amino");
const encoding_1 = require("@cosmjs/encoding");
const math_1 = require("@cosmjs/math");
const keys_1 = require("cosmjs-types/cosmos/crypto/ed25519/keys");
const keys_2 = require("cosmjs-types/cosmos/crypto/multisig/keys");
const keys_3 = require("cosmjs-types/cosmos/crypto/secp256k1/keys");
const any_1 = require("cosmjs-types/google/protobuf/any");
/**
 * Takes a pubkey in the Amino JSON object style (type/value wrapper)
 * and convertes it into a protobuf `Any`.
 *
 * This is the reverse operation to `decodePubkey`.
 */
function encodePubkey(pubkey) {
    if ((0, amino_1.isSecp256k1Pubkey)(pubkey)) {
        const pubkeyProto = keys_3.PubKey.fromPartial({
            key: (0, encoding_1.fromBase64)(pubkey.value),
        });
        return any_1.Any.fromPartial({
            typeUrl: "/cosmos.crypto.secp256k1.PubKey",
            value: Uint8Array.from(keys_3.PubKey.encode(pubkeyProto).finish()),
        });
    }
    else if ((0, amino_1.isEd25519Pubkey)(pubkey)) {
        const pubkeyProto = keys_1.PubKey.fromPartial({
            key: (0, encoding_1.fromBase64)(pubkey.value),
        });
        return any_1.Any.fromPartial({
            typeUrl: "/cosmos.crypto.ed25519.PubKey",
            value: Uint8Array.from(keys_1.PubKey.encode(pubkeyProto).finish()),
        });
    }
    else if ((0, amino_1.isMultisigThresholdPubkey)(pubkey)) {
        const pubkeyProto = keys_2.LegacyAminoPubKey.fromPartial({
            threshold: math_1.Uint53.fromString(pubkey.value.threshold).toNumber(),
            publicKeys: pubkey.value.pubkeys.map(encodePubkey),
        });
        return any_1.Any.fromPartial({
            typeUrl: "/cosmos.crypto.multisig.LegacyAminoPubKey",
            value: Uint8Array.from(keys_2.LegacyAminoPubKey.encode(pubkeyProto).finish()),
        });
    }
    else {
        throw new Error(`Pubkey type ${pubkey.type} not recognized`);
    }
}
exports.encodePubkey = encodePubkey;
/**
 * Decodes a single pubkey (i.e. not a multisig pubkey) from `Any` into
 * `SinglePubkey`.
 *
 * In most cases you probably want to use `decodePubkey`.
 */
function anyToSinglePubkey(pubkey) {
    switch (pubkey.typeUrl) {
        case "/cosmos.crypto.secp256k1.PubKey": {
            const { key } = keys_3.PubKey.decode(pubkey.value);
            return (0, amino_1.encodeSecp256k1Pubkey)(key);
        }
        case "/cosmos.crypto.ed25519.PubKey": {
            const { key } = keys_1.PubKey.decode(pubkey.value);
            return (0, amino_1.encodeEd25519Pubkey)(key);
        }
        default:
            throw new Error(`Pubkey type_url ${pubkey.typeUrl} not recognized as single public key type`);
    }
}
exports.anyToSinglePubkey = anyToSinglePubkey;
/**
 * Decodes a pubkey from a protobuf `Any` into `Pubkey`.
 * This supports single pubkeys such as Cosmos ed25519 and secp256k1 keys
 * as well as multisig threshold pubkeys.
 */
function decodePubkey(pubkey) {
    switch (pubkey.typeUrl) {
        case "/cosmos.crypto.secp256k1.PubKey":
        case "/cosmos.crypto.ed25519.PubKey": {
            return anyToSinglePubkey(pubkey);
        }
        case "/cosmos.crypto.multisig.LegacyAminoPubKey": {
            const { threshold, publicKeys } = keys_2.LegacyAminoPubKey.decode(pubkey.value);
            const out = {
                type: "tendermint/PubKeyMultisigThreshold",
                value: {
                    threshold: threshold.toString(),
                    pubkeys: publicKeys.map(anyToSinglePubkey),
                },
            };
            return out;
        }
        default:
            throw new Error(`Pubkey type URL '${pubkey.typeUrl}' not recognized`);
    }
}
exports.decodePubkey = decodePubkey;
/**
 * Decodes an optional pubkey from a protobuf `Any` into `Pubkey | null`.
 * This supports single pubkeys such as Cosmos ed25519 and secp256k1 keys
 * as well as multisig threshold pubkeys.
 */
function decodeOptionalPubkey(pubkey) {
    if (!pubkey)
        return null;
    if (pubkey.typeUrl) {
        if (pubkey.value.length) {
            // both set
            return decodePubkey(pubkey);
        }
        else {
            throw new Error(`Pubkey is an Any with type URL '${pubkey.typeUrl}' but an empty value`);
        }
    }
    else {
        if (pubkey.value.length) {
            throw new Error(`Pubkey is an Any with an empty type URL but a value set`);
        }
        else {
            // both unset, assuming this empty instance means null
            return null;
        }
    }
}
exports.decodeOptionalPubkey = decodeOptionalPubkey;
//# sourceMappingURL=pubkey.js.map