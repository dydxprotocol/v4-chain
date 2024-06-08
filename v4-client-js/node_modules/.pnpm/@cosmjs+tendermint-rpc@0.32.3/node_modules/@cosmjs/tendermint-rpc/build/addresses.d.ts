export declare function rawEd25519PubkeyToRawAddress(pubkeyData: Uint8Array): Uint8Array;
export declare function rawSecp256k1PubkeyToRawAddress(pubkeyData: Uint8Array): Uint8Array;
/**
 * Returns Tendermint address as bytes.
 *
 * This is for addresses that are derived by the Tendermint keypair (typically Ed25519).
 * Sometimes those addresses are bech32-encoded and contain the term "cons" in the presix
 * ("cosmosvalcons1...").
 *
 * For secp256k1 this assumes we already have a compressed pubkey, which is the default in Cosmos.
 */
export declare function pubkeyToRawAddress(type: "ed25519" | "secp256k1", data: Uint8Array): Uint8Array;
/**
 * Returns Tendermint address in uppercase hex format.
 *
 * This is for addresses that are derived by the Tendermint keypair (typically Ed25519).
 * Sometimes those addresses are bech32-encoded and contain the term "cons" in the presix
 * ("cosmosvalcons1...").
 *
 * For secp256k1 this assumes we already have a compressed pubkey, which is the default in Cosmos.
 */
export declare function pubkeyToAddress(type: "ed25519" | "secp256k1", data: Uint8Array): string;
