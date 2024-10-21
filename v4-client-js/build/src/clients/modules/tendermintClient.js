"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.TendermintClient = void 0;
const encoding_1 = require("@cosmjs/encoding");
const math_1 = require("@cosmjs/math");
const stargate_1 = require("@cosmjs/stargate");
const tendermint_rpc_1 = require("@cosmjs/tendermint-rpc");
const utils_1 = require("@cosmjs/utils");
const errors_1 = require("../lib/errors");
class TendermintClient {
    constructor(baseClient, broadcastOptions) {
        this.baseClient = baseClient;
        this.broadcastOptions = broadcastOptions;
    }
    /**
     * @description Get a specific block if height is specified. Otherwise, get the most recent block.
     *
     * @returns Information about the block queried.
     */
    async getBlock(height) {
        const response = await this.baseClient.block(height);
        return {
            id: (0, encoding_1.toHex)(response.blockId.hash).toUpperCase(),
            header: {
                version: {
                    block: new math_1.Uint53(response.block.header.version.block).toString(),
                    app: new math_1.Uint53(response.block.header.version.app).toString(),
                },
                height: response.block.header.height,
                chainId: response.block.header.chainId,
                time: (0, tendermint_rpc_1.toRfc3339WithNanoseconds)(response.block.header.time),
            },
            txs: response.block.txs,
        };
    }
    /**
      * @description Broadcast a signed transaction with a specific mode.
      * @throws BroadcastErrorObject when result code is not zero. TypeError when mode is invalid.
      * @returns Differs depending on the BroadcastMode used.
      * See https://docs.cosmos.network/master/run-node/txs.html for more information.
      */
    async broadcastTransaction(tx, mode) {
        switch (mode) {
            case tendermint_rpc_1.Method.BroadcastTxAsync:
                return this.broadcastTransactionAsync(tx);
            case tendermint_rpc_1.Method.BroadcastTxSync:
                return this.broadcastTransactionSync(tx);
            case tendermint_rpc_1.Method.BroadcastTxCommit:
                return this.broadcastTransactionCommit(tx);
            default:
                throw new TypeError('broadcastTransaction: invalid BroadcastMode');
        }
    }
    /**
     * @description Broadcast a signed transaction.
     * @returns The transaction hash.
     */
    broadcastTransactionAsync(tx) {
        return this.baseClient.broadcastTxAsync({ tx });
    }
    /**
     * @description Broadcast a signed transaction and await the response.
     * @throws BroadcastErrorObject when result code is not zero.
     * @returns The response from the node once the transaction is processed by `CheckTx`.
     */
    async broadcastTransactionSync(tx) {
        const result = await this.baseClient.broadcastTxSync({ tx });
        if (result.code !== 0) {
            throw new errors_1.BroadcastErrorObject(`Broadcasting transaction failed: ${result.log}`, result);
        }
        return result;
    }
    /**
     * @description Broadcast a signed transaction and await for it to be included in the blockchain.
     * @throws BroadcastErrorObject when result code is not zero.
     * @returns The result of the transaction once included in the blockchain.
     */
    async broadcastTransactionCommit(tx) {
        const result = await this.broadcastTransactionSync(tx);
        return this.queryHash(result.hash);
    }
    /**
     * @description Using tx method, query for a transaction on-chain with retries specified by
     * the client BroadcastOptions.
     *
     * @throws TimeoutError if the transaction is not committed on-chain within the timeout limit.
     * @returns An indexed transaction containing information about the transaction when committed.
     */
    async queryHash(hash, time = 0) {
        const now = Date.now();
        const transactionId = (0, encoding_1.toHex)(hash).toUpperCase();
        if (time >= this.broadcastOptions.broadcastTimeoutMs) {
            throw new stargate_1.TimeoutError(`Transaction with hash [${hash}] was submitted but was not yet found on the chain. You might want to check later. Query timed out after ${this.broadcastOptions.broadcastTimeoutMs / 1000} seconds.`, transactionId);
        }
        await (0, utils_1.sleep)(this.broadcastOptions.broadcastPollIntervalMs);
        // If the transaction is not found, the tx method will throw an Internal Error.
        try {
            const tx = await this.baseClient.tx({ hash });
            return {
                height: tx.height,
                hash: (0, encoding_1.toHex)(tx.hash).toUpperCase(),
                code: tx.result.code,
                rawLog: tx.result.log !== undefined ? tx.result.log : '',
                tx: tx.tx,
                txIndex: tx.index,
                gasUsed: tx.result.gasUsed,
                gasWanted: tx.result.gasWanted,
                // Convert stargate events to tendermint events.
                events: tx.result.events.map((event) => {
                    return {
                        ...event,
                        attributes: event.attributes.map((attr) => {
                            return {
                                ...attr,
                                key: Buffer.from(attr.key).toString(),
                                value: Buffer.from(attr.value).toString(),
                            };
                        }),
                    };
                }),
                // @ts-ignore
                msgResponses: [],
            };
        }
        catch (error) {
            return this.queryHash(hash, time + Date.now() - now);
        }
    }
    /**
     * @description Set the broadcast options for this module.
     */
    setBroadcastOptions(broadcastOptions) {
        this.broadcastOptions = broadcastOptions;
    }
}
exports.TendermintClient = TendermintClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoidGVuZGVybWludENsaWVudC5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uL3NyYy9jbGllbnRzL21vZHVsZXMvdGVuZGVybWludENsaWVudC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7QUFBQSwrQ0FBeUM7QUFDekMsdUNBQXNDO0FBQ3RDLCtDQUkwQjtBQUMxQiwyREFJZ0M7QUFTaEMseUNBQXNDO0FBRXRDLDBDQUFxRDtBQUdyRCxNQUFhLGdCQUFnQjtJQUkzQixZQUNFLFVBQThCLEVBQzlCLGdCQUFrQztRQUVsQyxJQUFJLENBQUMsVUFBVSxHQUFHLFVBQVUsQ0FBQztRQUM3QixJQUFJLENBQUMsZ0JBQWdCLEdBQUcsZ0JBQWdCLENBQUM7SUFDM0MsQ0FBQztJQUVEOzs7O09BSUc7SUFDSCxLQUFLLENBQUMsUUFBUSxDQUFDLE1BQWU7UUFDNUIsTUFBTSxRQUFRLEdBQWtCLE1BQU0sSUFBSSxDQUFDLFVBQVUsQ0FBQyxLQUFLLENBQUMsTUFBTSxDQUFDLENBQUM7UUFDcEUsT0FBTztZQUNMLEVBQUUsRUFBRSxJQUFBLGdCQUFLLEVBQUMsUUFBUSxDQUFDLE9BQU8sQ0FBQyxJQUFJLENBQUMsQ0FBQyxXQUFXLEVBQUU7WUFDOUMsTUFBTSxFQUFFO2dCQUNOLE9BQU8sRUFBRTtvQkFDUCxLQUFLLEVBQUUsSUFBSSxhQUFNLENBQUMsUUFBUSxDQUFDLEtBQUssQ0FBQyxNQUFNLENBQUMsT0FBTyxDQUFDLEtBQUssQ0FBQyxDQUFDLFFBQVEsRUFBRTtvQkFDakUsR0FBRyxFQUFFLElBQUksYUFBTSxDQUFDLFFBQVEsQ0FBQyxLQUFLLENBQUMsTUFBTSxDQUFDLE9BQU8sQ0FBQyxHQUFHLENBQUMsQ0FBQyxRQUFRLEVBQUU7aUJBQzlEO2dCQUNELE1BQU0sRUFBRSxRQUFRLENBQUMsS0FBSyxDQUFDLE1BQU0sQ0FBQyxNQUFNO2dCQUNwQyxPQUFPLEVBQUUsUUFBUSxDQUFDLEtBQUssQ0FBQyxNQUFNLENBQUMsT0FBTztn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