"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.makeSignBytes = exports.makeSignDoc = exports.makeAuthInfoBytes = void 0;
/* eslint-disable @typescript-eslint/naming-convention */
const utils_1 = require("@cosmjs/utils");
const signing_1 = require("cosmjs-types/cosmos/tx/signing/v1beta1/signing");
const tx_1 = require("cosmjs-types/cosmos/tx/v1beta1/tx");
/**
 * Create signer infos from the provided signers.
 *
 * This implementation does not support different signing modes for the different signers.
 */
function makeSignerInfos(signers, signMode) {
    return signers.map(({ pubkey, sequence }) => ({
        publicKey: pubkey,
        modeInfo: {
            single: { mode: signMode },
        },
        sequence: BigInt(sequence),
    }));
}
/**
 * Creates and serializes an AuthInfo document.
 *
 * This implementation does not support different signing modes for the different signers.
 */
function makeAuthInfoBytes(signers, feeAmount, gasLimit, feeGranter, feePayer, signMode = signing_1.SignMode.SIGN_MODE_DIRECT) {
    // Required arguments 4 and 5 were added in CosmJS 0.29. Use runtime checks to help our non-TS users.
    (0, utils_1.assert)(feeGranter === undefined || typeof feeGranter === "string", "feeGranter must be undefined or string");
    (0, utils_1.assert)(feePayer === undefined || typeof feePayer === "string", "feePayer must be undefined or string");
    const authInfo = tx_1.AuthInfo.fromPartial({
        signerInfos: makeSignerInfos(signers, signMode),
        fee: {
            amount: [...feeAmount],
            gasLimit: BigInt(gasLimit),
            granter: feeGranter,
            payer: feePayer,
        },
    });
    return tx_1.AuthInfo.encode(authInfo).finish();
}
exports.makeAuthInfoBytes = makeAuthInfoBytes;
function makeSignDoc(bodyBytes, authInfoBytes, chainId, accountNumber) {
    return {
        bodyBytes: bodyBytes,
        authInfoBytes: authInfoBytes,
        chainId: chainId,
        accountNumber: BigInt(accountNumber),
    };
}
exports.makeSignDoc = makeSignDoc;
function makeSignBytes({ accountNumber, authInfoBytes, bodyBytes, chainId }) {
    const signDoc = tx_1.SignDoc.fromPartial({
        accountNumber: accountNumber,
        authInfoBytes: authInfoBytes,
        bodyBytes: bodyBytes,
        chainId: chainId,
    });
    return tx_1.SignDoc.encode(signDoc).finish();
}
exports.makeSignBytes = makeSignBytes;
//# sourceMappingURL=signing.js.map