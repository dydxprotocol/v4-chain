"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.HttpBatchClient = void 0;
const json_rpc_1 = require("@cosmjs/json-rpc");
const http_1 = require("./http");
const rpcclient_1 = require("./rpcclient");
// Those values are private and can change any time.
// Does a user need to know them? I don't think so. You either set
// a custom value or leave the option field unset.
const defaultHttpBatchClientOptions = {
    dispatchInterval: 20,
    batchSizeLimit: 20,
};
class HttpBatchClient {
    constructor(endpoint, options = {}) {
        this.queue = [];
        this.options = {
            batchSizeLimit: options.batchSizeLimit ?? defaultHttpBatchClientOptions.batchSizeLimit,
            dispatchInterval: options.dispatchInterval ?? defaultHttpBatchClientOptions.dispatchInterval,
        };
        if (typeof endpoint === "string") {
            if (!(0, rpcclient_1.hasProtocol)(endpoint)) {
                throw new Error("Endpoint URL is missing a protocol. Expected 'https://' or 'http://'.");
            }
            this.url = endpoint;
        }
        else {
            this.url = endpoint.url;
            this.headers = endpoint.headers;
        }
        this.timer = setInterval(() => this.tick(), options.dispatchInterval);
        this.validate();
    }
    disconnect() {
        this.timer && clearInterval(this.timer);
        this.timer = undefined;
    }
    async execute(request) {
        return new Promise((resolve, reject) => {
            this.queue.push({ request, resolve, reject });
            if (this.queue.length >= this.options.batchSizeLimit) {
                // this train is full, let's go
                this.tick();
            }
        });
    }
    validate() {
        if (!this.options.batchSizeLimit ||
            !Number.isSafeInteger(this.options.batchSizeLimit) ||
            this.options.batchSizeLimit < 1) {
            throw new Error("batchSizeLimit must be a safe integer >= 1");
        }
    }
    /**
     * This is called in an interval where promise rejections cannot be handled.
     * So this is not async and HTTP errors need to be handled by the queued promises.
     */
    tick() {
        // Avoid race conditions
        const batch = this.queue.splice(0, this.options.batchSizeLimit);
        if (!batch.length)
            return;
        const requests = batch.map((s) => s.request);
        const requestIds = requests.map((request) => request.id);
        (0, http_1.http)("POST", this.url, this.headers, requests).then((raw) => {
            // Requests with a single entry return as an object
            const arr = Array.isArray(raw) ? raw : [raw];
            arr.forEach((el) => {
                const req = batch.find((s) => s.request.id === el.id);
                if (!req)
                    return;
                const { reject, resolve } = req;
                const response = (0, json_rpc_1.parseJsonRpcResponse)(el);
                if ((0, json_rpc_1.isJsonRpcErrorResponse)(response)) {
                    reject(new Error(JSON.stringify(response.error)));
                }
                else {
                    resolve(response);
                }
            });
        }, (error) => {
            for (const requestId of requestIds) {
                const req = batch.find((s) => s.request.id === requestId);
                if (!req)
                    return;
                req.reject(error);
            }
        });
    }
}
exports.HttpBatchClient = HttpBatchClient;
//# sourceMappingURL=httpbatchclient.js.map