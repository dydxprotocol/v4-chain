"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.setupTxExtension = void 0;
const proto_signing_1 = require("@cosmjs/proto-signing");
const signing_1 = require("cosmjs-types/cosmos/tx/signing/v1beta1/signing");
const service_1 = require("cosmjs-types/cosmos/tx/v1beta1/service");
const tx_1 = require("cosmjs-types/cosmos/tx/v1beta1/tx");
const queryclient_1 = require("../../queryclient");
function setupTxExtension(base) {
    // Use this service to get easy typed access to query methods
    // This cannot be used for proof verification
    const rpc = (0, queryclient_1.createProtobufRpcClient)(base);
    const queryService = new service_1.ServiceClientImpl(rpc);
    return {
        tx: {
            getTx: async (txId) => {
                const request = {
                    hash: txId,
                };
                const response = await queryService.GetTx(request);
                return response;
            },
            simulate: async (messages, memo, signer, sequence) => {
                const tx = tx_1.Tx.fromPartial({
                    authInfo: tx_1.AuthInfo.fromPartial({
                        fee: tx_1.Fee.fromPartial({}),
                        signerInfos: [
                            {
                                publicKey: (0, proto_signing_1.encodePubkey)(signer),
                                sequence: BigInt(sequence),
                                modeInfo: { single: { mode: signing_1.SignMode.SIGN_MODE_UNSPECIFIED } },
                            },
                        ],
                    }),
                    body: tx_1.TxBody.fromPartial({
                        messages: Array.from(messages),
                        memo: memo,
                    }),
                    signatures: [new Uint8Array()],
                });
                const request = service_1.SimulateRequest.fromPartial({
                    txBytes: tx_1.Tx.encode(tx).finish(),
                });
                const response = await queryService.Simulate(request);
                return response;
            },
        },
    };
}
exports.setupTxExtension = setupTxExtension;
//# sourceMappingURL=queries.js.map