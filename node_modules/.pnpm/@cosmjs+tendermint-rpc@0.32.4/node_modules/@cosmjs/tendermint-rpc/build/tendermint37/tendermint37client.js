"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.Tendermint37Client = void 0;
const jsonrpc_1 = require("../jsonrpc");
const rpcclients_1 = require("../rpcclients");
const adaptor_1 = require("./adaptor");
const requests = __importStar(require("./requests"));
class Tendermint37Client {
    /**
     * Creates a new Tendermint client for the given endpoint.
     *
     * Uses HTTP when the URL schema is http or https. Uses WebSockets otherwise.
     */
    static async connect(endpoint) {
        let rpcClient;
        if (typeof endpoint === "object") {
            rpcClient = new rpcclients_1.HttpClient(endpoint);
        }
        else {
            const useHttp = endpoint.startsWith("http://") || endpoint.startsWith("https://");
            rpcClient = useHttp ? new rpcclients_1.HttpClient(endpoint) : new rpcclients_1.WebsocketClient(endpoint);
        }
        // For some very strange reason I don't understand, tests start to fail on some systems
        // (our CI) when skipping the status call before doing other queries. Sleeping a little
        // while did not help. Thus we query the version as a way to say "hi" to the backend,
        // even in cases where we don't use the result.
        const _version = await this.detectVersion(rpcClient);
        return Tendermint37Client.create(rpcClient);
    }
    /**
     * Creates a new Tendermint client given an RPC client.
     */
    static async create(rpcClient) {
        return new Tendermint37Client(rpcClient);
    }
    static async detectVersion(client) {
        const req = (0, jsonrpc_1.createJsonRpcRequest)(requests.Method.Status);
        const response = await client.execute(req);
        const result = response.result;
        if (!result || !result.node_info) {
            throw new Error("Unrecognized format for status response");
        }
        const version = result.node_info.version;
        if (typeof version !== "string") {
            throw new Error("Unrecognized version format: must be string");
        }
        return version;
    }
    /**
     * Use `Tendermint37Client.connect` or `Tendermint37Client.create` to create an instance.
     */
    constructor(client) {
        this.client = client;
    }
    disconnect() {
        this.client.disconnect();
    }
    async abciInfo() {
        const query = { method: requests.Method.AbciInfo };
        return this.doCall(query, adaptor_1.Params.encodeAbciInfo, adaptor_1.Responses.decodeAbciInfo);
    }
    async abciQuery(params) {
        const query = { params: params, method: requests.Method.AbciQuery };
        return this.doCall(query, adaptor_1.Params.encodeAbciQuery, adaptor_1.Responses.decodeAbciQuery);
    }
    async block(height) {
        const query = { method: requests.Method.Block, params: { height: height } };
        return this.doCall(query, adaptor_1.Params.encodeBlock, adaptor_1.Responses.decodeBlock);
    }
    async blockResults(height) {
        const query = {
            method: requests.Method.BlockResults,
            params: { height: height },
        };
        return this.doCall(query, adaptor_1.Params.encodeBlockResults, adaptor_1.Responses.decodeBlockResults);
    }
    /**
     * Search for events that are in a block.
     *
     * NOTE
     * This method will error on any node that is running a Tendermint version lower than 0.34.9.
     *
     * @see https://docs.tendermint.com/master/rpc/#/Info/block_search
     */
    async blockSearch(params) {
        const query = { params: params, method: requests.Method.BlockSearch };
        const resp = await this.doCall(query, adaptor_1.Params.encodeBlockSearch, adaptor_1.Responses.decodeBlockSearch);
        return {
            ...resp,
            // make sure we sort by height, as tendermint may be sorting by string value of the height
            blocks: [...resp.blocks].sort((a, b) => a.block.header.height - b.block.header.height),
        };
    }
    // this should paginate through all blockSearch options to ensure it returns all results.
    // starts with page 1 or whatever was provided (eg. to start on page 7)
    //
    // NOTE
    // This method will error on any node that is running a Tendermint version lower than 0.34.9.
    async blockSearchAll(params) {
        let page = params.page || 1;
        const blocks = [];
        let done = false;
        while (!done) {
            const resp = await this.blockSearch({ ...params, page: page });
            blocks.push(...resp.blocks);
            if (blocks.length < resp.totalCount) {
                page++;
            }
            else {
                done = true;
            }
        }
        // make sure we sort by height, as tendermint may be sorting by string value of the height
        // and the earlier items may be in a higher page than the later items
        blocks.sort((a, b) => a.block.header.height - b.block.header.height);
        return {
            totalCount: blocks.length,
            blocks: blocks,
        };
    }
    /**
     * Queries block headers filtered by minHeight <= height <= maxHeight.
     *
     * @param minHeight The minimum height to be included in the result. Defaults to 0.
     * @param maxHeight The maximum height to be included in the result. Defaults to infinity.
     */
    async blockchain(minHeight, maxHeight) {
        const query = {
            method: requests.Method.Blockchain,
            params: {
                minHeight: minHeight,
                maxHeight: maxHeight,
            },
        };
        return this.doCall(query, adaptor_1.Params.encodeBlockchain, adaptor_1.Responses.decodeBlockchain);
    }
    /**
     * Broadcast transaction to mempool and wait for response
     *
     * @see https://docs.tendermint.com/master/rpc/#/Tx/broadcast_tx_sync
     */
    async broadcastTxSync(params) {
        const query = { params: params, method: requests.Method.BroadcastTxSync };
        return this.doCall(query, adaptor_1.Params.encodeBroadcastTx, adaptor_1.Responses.decodeBroadcastTxSync);
    }
    /**
     * Broadcast transaction to mempool and do not wait for result
     *
     * @see https://docs.tendermint.com/master/rpc/#/Tx/broadcast_tx_async
     */
    async broadcastTxAsync(params) {
        const query = { params: params, method: requests.Method.BroadcastTxAsync };
        return this.doCall(query, adaptor_1.Params.encodeBroadcastTx, adaptor_1.Responses.decodeBroadcastTxAsync);
    }
    /**
     * Broadcast transaction to mempool and wait for block
     *
     * @see https://docs.tendermint.com/master/rpc/#/Tx/broadcast_tx_commit
     */
    async broadcastTxCommit(params) {
        const query = { params: params, method: requests.Method.BroadcastTxCommit };
        return this.doCall(query, adaptor_1.Params.encodeBroadcastTx, adaptor_1.Responses.decodeBroadcastTxCommit);
    }
    async commit(height) {
        const query = { method: requests.Method.Commit, params: { height: height } };
        return this.doCall(query, adaptor_1.Params.encodeCommit, adaptor_1.Responses.decodeCommit);
    }
    async genesis() {
        const query = { method: requests.Method.Genesis };
        return this.doCall(query, adaptor_1.Params.encodeGenesis, adaptor_1.Responses.decodeGenesis);
    }
    async health() {
        const query = { method: requests.Method.Health };
        return this.doCall(query, adaptor_1.Params.encodeHealth, adaptor_1.Responses.decodeHealth);
    }
    async numUnconfirmedTxs() {
        const query = { method: requests.Method.NumUnconfirmedTxs };
        return this.doCall(query, adaptor_1.Params.encodeNumUnconfirmedTxs, adaptor_1.Responses.decodeNumUnconfirmedTxs);
    }
    async status() {
        const query = { method: requests.Method.Status };
        return this.doCall(query, adaptor_1.Params.encodeStatus, adaptor_1.Responses.decodeStatus);
    }
    subscribeNewBlock() {
        const request = {
            method: requests.Method.Subscribe,
            query: { type: requests.SubscriptionEventType.NewBlock },
        };
        return this.subscribe(request, adaptor_1.Responses.decodeNewBlockEvent);
    }
    subscribeNewBlockHeader() {
        const request = {
            method: requests.Method.Subscribe,
            query: { type: requests.SubscriptionEventType.NewBlockHeader },
        };
        return this.subscribe(request, adaptor_1.Responses.decodeNewBlockHeaderEvent);
    }
    subscribeTx(query) {
        const request = {
            method: requests.Method.Subscribe,
            query: {
                type: requests.SubscriptionEventType.Tx,
                raw: query,
            },
        };
        return this.subscribe(request, adaptor_1.Responses.decodeTxEvent);
    }
    /**
     * Get a single transaction by hash
     *
     * @see https://docs.tendermint.com/master/rpc/#/Info/tx
     */
    async tx(params) {
        const query = { params: params, method: requests.Method.Tx };
        return this.doCall(query, adaptor_1.Params.encodeTx, adaptor_1.Responses.decodeTx);
    }
    /**
     * Search for transactions that are in a block
     *
     * @see https://docs.tendermint.com/master/rpc/#/Info/tx_search
     */
    async txSearch(params) {
        const query = { params: params, method: requests.Method.TxSearch };
        return this.doCall(query, adaptor_1.Params.encodeTxSearch, adaptor_1.Responses.decodeTxSearch);
    }
    // this should paginate through all txSearch options to ensure it returns all results.
    // starts with page 1 or whatever was provided (eg. to start on page 7)
    async txSearchAll(params) {
        let page = params.page || 1;
        const txs = [];
        let done = false;
        while (!done) {
            const resp = await this.txSearch({ ...params, page: page });
            txs.push(...resp.txs);
            if (txs.length < resp.totalCount) {
                page++;
            }
            else {
                done = true;
            }
        }
        return {
            totalCount: txs.length,
            txs: txs,
        };
    }
    async validators(params) {
        const query = {
            method: requests.Method.Validators,
            params: params,
        };
        return this.doCall(query, adaptor_1.Params.encodeValidators, adaptor_1.Responses.decodeValidators);
    }
    async validatorsAll(height) {
        const validators = [];
        let page = 1;
        let done = false;
        let blockHeight = height;
        while (!done) {
            const response = await this.validators({
                per_page: 50,
                height: blockHeight,
                page: page,
            });
            validators.push(...response.validators);
            blockHeight = blockHeight || response.blockHeight;
            if (validators.length < response.total) {
                page++;
            }
            else {
                done = true;
            }
        }
        return {
            // NOTE: Default value is for type safety but this should always be set
            blockHeight: blockHeight ?? 0,
            count: validators.length,
            total: validators.length,
            validators: validators,
        };
    }
    // doCall is a helper to handle the encode/call/decode logic
    async doCall(request, encode, decode) {
        const req = encode(request);
        const result = await this.client.execute(req);
        return decode(result);
    }
    subscribe(request, decode) {
        if (!(0, rpcclients_1.instanceOfRpcStreamingClient)(this.client)) {
            throw new Error("This RPC client type cannot subscribe to events");
        }
        const req = adaptor_1.Params.encodeSubscribe(request);
        const eventStream = this.client.listen(req);
        return eventStream.map((event) => {
            return decode(event);
        });
    }
}
exports.Tendermint37Client = Tendermint37Client;
//# sourceMappingURL=tendermint37client.js.map