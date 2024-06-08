"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.deriveHDKeyFromEthereumSignature = exports.exportMnemonicAndPrivateKey = void 0;
const bip32_1 = require("@scure/bip32");
const bip39_1 = require("@scure/bip39");
const english_1 = require("@scure/bip39/wordlists/english");
const keccak_1 = require("ethereum-cryptography/keccak");
const helpers_1 = require("./helpers");
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
const exportMnemonicAndPrivateKey = (entropy, path = "m/44'/118'/0'/0/0") => {
    const mnemonic = (0, bip39_1.entropyToMnemonic)(entropy, english_1.wordlist);
    const seed = (0, bip39_1.mnemonicToSeedSync)(mnemonic);
    const hdkey = bip32_1.HDKey.fromMasterSeed(seed);
    const derivedHdkey = hdkey.derive(path);
    if (!hdkey.privateKey) {
        throw new Error('null hd key');
    }
    return {
        mnemonic,
        privateKey: derivedHdkey.privateKey,
        publicKey: derivedHdkey.publicKey,
    };
};
exports.exportMnemonicAndPrivateKey = exportMnemonicAndPrivateKey;
/**
 * @description Get private information for onboarding using an Ethereum Signature.
 *
 * @returns Mnemonic and Public/Private HD keys
 */
const deriveHDKeyFromEthereumSignature = (signature) => {
    const buffer = Buffer.from((0, helpers_1.stripHexPrefix)(signature), 'hex');
    if (buffer.length !== 65) {
        throw new Error('Signature must be 65 bytes');
    }
    // Remove the 'v' value by taking only the first 64 bytes of the signature
    const rsValues = buffer.subarray(0, 64);
    // Hash the 'r' and 's' values down to 32 bytes (256 bits) using Keccak-256
    const entropy = (0, keccak_1.keccak256)(rsValues);
    return (0, exports.exportMnemonicAndPrivateKey)(entropy);
};
exports.deriveHDKeyFromEthereumSignature = deriveHDKeyFromEthereumSignature;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoib25ib2FyZGluZy5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uL3NyYy9saWIvb25ib2FyZGluZy50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7QUFBQSx3Q0FBcUM7QUFDckMsd0NBQXFFO0FBQ3JFLDREQUEwRDtBQUMxRCx5REFBeUQ7QUFFekQsdUNBQTJDO0FBRTNDOzs7Ozs7Ozs7Ozs7R0FZRztBQUNJLE1BQU0sMkJBQTJCLEdBQUcsQ0FBQyxPQUFtQixFQUFFLE9BQWUsbUJBQW1CLEVBSWpHLEVBQUU7SUFDRixNQUFNLFFBQVEsR0FBRyxJQUFBLHlCQUFpQixFQUFDLE9BQU8sRUFBRSxrQkFBUSxDQUFDLENBQUM7SUFDdEQsTUFBTSxJQUFJLEdBQUcsSUFBQSwwQkFBa0IsRUFBQyxRQUFRLENBQUMsQ0FBQztJQUUxQyxNQUFNLEtBQUssR0FBRyxhQUFLLENBQUMsY0FBYyxDQUFDLElBQUksQ0FBQyxDQUFDO0lBQ3pDLE1BQU0sWUFBWSxHQUFHLEtBQUssQ0FBQyxNQUFNLENBQUMsSUFBSSxDQUFDLENBQUM7SUFFeEMsSUFBSSxDQUFDLEtBQUssQ0FBQyxVQUFVLEVBQUU7UUFDckIsTUFBTSxJQUFJLEtBQUssQ0FBQyxhQUFhLENBQUMsQ0FBQztLQUNoQztJQUVELE9BQU87UUFDTCxRQUFRO1FBQ1IsVUFBVSxFQUFFLFlBQVksQ0FBQyxVQUFVO1FBQ25DLFNBQVMsRUFBRSxZQUFZLENBQUMsU0FBUztLQUNsQyxDQUFDO0FBQ0osQ0FBQyxDQUFDO0FBcEJXLFFBQUEsMkJBQTJCLCtCQW9CdEM7QUFFRjs7OztHQUlHO0FBQ0ksTUFBTSxnQ0FBZ0MsR0FBRyxDQUFDLFNBQWlCLEVBSWhFLEVBQUU7SUFDRixNQUFNLE1BQU0sR0FBRyxNQUFNLENBQUMsSUFBSSxDQUFDLElBQUEsd0JBQWMsRUFBQyxTQUFTLENBQUMsRUFBRSxLQUFLLENBQUMsQ0FBQztJQUU3RCxJQUFJLE1BQU0sQ0FBQyxNQUFNLEtBQUssRUFBRSxFQUFFO1FBQ3hCLE1BQU0sSUFBSSxLQUFLLENBQUMsNEJBQTRCLENBQUMsQ0FBQztLQUMvQztJQUVELDBFQUEwRTtJQUMxRSxNQUFNLFFBQVEsR0FBRyxNQUFNLENBQUMsUUFBUSxDQUFDLENBQUMsRUFBRSxFQUFFLENBQUMsQ0FBQztJQUN4QywyRUFBMkU7SUFDM0UsTUFBTSxPQUFPLEdBQUcsSUFBQSxrQkFBUyxFQUFDLFFBQVEsQ0FBQyxDQUFDO0lBQ3BDLE9BQU8sSUFBQSxtQ0FBMkIsRUFBQyxPQUFPLENBQUMsQ0FBQztBQUM5QyxDQUFDLENBQUM7QUFoQlcsUUFBQSxnQ0FBZ0Msb0NBZ0IzQyJ9