"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.TransactionSigner = void 0;
const tx_1 = require("cosmjs-types/cosmos/tx/v1beta1/tx");
const long_1 = __importDefault(require("long"));
const protobufjs_1 = __importDefault(require("protobufjs"));
const errors_1 = require("../lib/errors");
// Required for encoding and decoding queries that are of type Long.
// Must be done once but since the individal modules should be usable
// - must be set in each module that encounters encoding/decoding Longs.
// Reference: https://github.com/protobufjs/protobuf.js/issues/921
protobufjs_1.default.util.Long = long_1.default;
protobufjs_1.default.configure();
class TransactionSigner {
    constructor(address, stargateSigningClient) {
        this.address = address;
        this.stargateSigningClient = stargateSigningClient;
    }
    /**
     * @description Get the encoded signed transaction or the promise is rejected if
     * no fee can be set for the transaction.
     *
     * @throws UserError if the fee is undefined.
     * @returns The signed and encoded transaction.
     */
    async signTransaction(messages, transactionOptions, fee, memo = '') {
        // Verify there is either a fee or a path to getting the fee present.
        if (fee === undefined) {
            throw new errors_1.UserError('fee cannot be undefined');
        }
        // Sign, encode and return the transaction.
        const rawTx = await this.stargateSigningClient.sign(this.address, messages, fee, memo, {
            accountNumber: transactionOptions.accountNumber,
            sequence: transactionOptions.sequence,
            chainId: transactionOptions.chainId,
        });
        return Uint8Array.from(tx_1.TxRaw.encode(rawTx).finish());
    }
}
exports.TransactionSigner = TransactionSigner;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoic2lnbmVyLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vc3JjL2NsaWVudHMvbW9kdWxlcy9zaWduZXIudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7O0FBS0EsMERBQTBEO0FBQzFELGdEQUF3QjtBQUN4Qiw0REFBa0M7QUFFbEMsMENBQTBDO0FBSzFDLG9FQUFvRTtBQUNwRSxxRUFBcUU7QUFDckUsd0VBQXdFO0FBQ3hFLGtFQUFrRTtBQUNsRSxvQkFBUSxDQUFDLElBQUksQ0FBQyxJQUFJLEdBQUcsY0FBSSxDQUFDO0FBQzFCLG9CQUFRLENBQUMsU0FBUyxFQUFFLENBQUM7QUFFckIsTUFBYSxpQkFBaUI7SUFJNUIsWUFDRSxPQUFlLEVBQ2YscUJBQTRDO1FBRTVDLElBQUksQ0FBQyxPQUFPLEdBQUcsT0FBTyxDQUFDO1FBQ3ZCLElBQUksQ0FBQyxxQkFBcUIsR0FBRyxxQkFBcUIsQ0FBQztJQUNyRCxDQUFDO0lBRUQ7Ozs7OztPQU1HO0lBQ0gsS0FBSyxDQUFDLGVBQWUsQ0FDbkIsUUFBd0IsRUFDeEIsa0JBQXNDLEVBQ3RDLEdBQVksRUFDWixPQUFlLEVBQUU7UUFFakIscUVBQXFFO1FBQ3JFLElBQUksR0FBRyxLQUFLLFNBQVMsRUFBRTtZQUNyQixNQUFNLElBQUksa0JBQVMsQ0FBQyx5QkFBeUIsQ0FBQyxDQUFDO1NBQ2hEO1FBRUQsMkNBQTJDO1FBQzNDLE1BQU0sS0FBSyxHQUFVLE1BQU0sSUFBSSxDQUFDLHFCQUFxQixDQUFDLElBQUksQ0FDeEQsSUFBSSxDQUFDLE9BQU8sRUFDWixRQUFRLEVBQ1IsR0FBRyxFQUNILElBQUksRUFDSjtZQUNFLGFBQWEsRUFBRSxrQkFBa0IsQ0FBQyxhQUFhO1lBQy9DLFFBQVEsRUFBRSxrQkFBa0IsQ0FBQyxRQUFRO1lBQ3JDLE9BQU8sRUFBRSxrQkFBa0IsQ0FBQyxPQUFPO1NBQ3BDLENBQ0YsQ0FBQztRQUNGLE9BQU8sVUFBVSxDQUFDLElBQUksQ0FBQyxVQUFLLENBQUMsTUFBTSxDQUFDLEtBQUssQ0FBQyxDQUFDLE1BQU0sRUFBRSxDQUFDLENBQUM7SUFDdkQsQ0FBQztDQUNGO0FBNUNELDhDQTRDQyJ9