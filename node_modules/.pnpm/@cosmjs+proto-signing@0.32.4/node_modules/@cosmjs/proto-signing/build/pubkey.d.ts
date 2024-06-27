import { Pubkey, SinglePubkey } from "@cosmjs/amino";
import { Any } from "cosmjs-types/google/protobuf/any";
/**
 * Takes a pubkey in the Amino JSON object style (type/value wrapper)
 * and convertes it into a protobuf `Any`.
 *
 * This is the reverse operation to `decodePubkey`.
 */
export declare function encodePubkey(pubkey: Pubkey): Any;
/**
 * Decodes a single pubkey (i.e. not a multisig pubkey) from `Any` into
 * `SinglePubkey`.
 *
 * In most cases you probably want to use `decodePubkey`.
 */
export declare function anyToSinglePubkey(pubkey: Any): SinglePubkey;
/**
 * Decodes a pubkey from a protobuf `Any` into `Pubkey`.
 * This supports single pubkeys such as Cosmos ed25519 and secp256k1 keys
 * as well as multisig threshold pubkeys.
 */
export declare function decodePubkey(pubkey: Any): Pubkey;
/**
 * Decodes an optional pubkey from a protobuf `Any` into `Pubkey | null`.
 * This supports single pubkeys such as Cosmos ed25519 and secp256k1 keys
 * as well as multisig threshold pubkeys.
 */
export declare function decodeOptionalPubkey(pubkey: Any | null | undefined): Pubkey | null;
