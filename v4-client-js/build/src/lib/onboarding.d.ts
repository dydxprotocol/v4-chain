/**
 * @description Get Mnemonic and priv/pub keys from privateKeyBytes and BIP44 HD path
 *
 * @url https://github.com/confio/cosmos-hd-key-derivation-spec#bip44
 *
 * @param entropy used to generate mnemonic
 *
 * @param path BIP44 HD Path. Default is The Cosmos Hub path
 *
 * @throws Error if the hdkey does not exist
 *
 * @returns Mnemonic and priv/pub keys
 */
export declare const exportMnemonicAndPrivateKey: (entropy: Uint8Array, path?: string) => {
    mnemonic: string;
    privateKey: Uint8Array | null;
    publicKey: Uint8Array | null;
};
/**
 * @description Get private information for onboarding using an Ethereum Signature.
 *
 * @returns Mnemonic and Public/Private HD keys
 */
export declare const deriveHDKeyFromEthereumSignature: (signature: string) => {
    mnemonic: string;
    privateKey: Uint8Array | null;
    publicKey: Uint8Array | null;
};
