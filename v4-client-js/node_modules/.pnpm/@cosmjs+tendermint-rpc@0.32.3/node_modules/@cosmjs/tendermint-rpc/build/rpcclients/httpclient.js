"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.HttpClient = void 0;
const json_rpc_1 = require("@cosmjs/json-rpc");
const http_1 = require("./http");
const rpcclient_1 = require("./rpcclient");
class HttpClient {
    constructor(endpoint) {
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
    }
    disconnect() {
        // nothing to be done
    }
    async execute(request) {
        const response = (0, json_rpc_1.parseJsonRpcResponse)(await (0, http_1.http)("POST", this.url, this.headers, request));
        if ((0, json_rpc_1.isJsonRpcErrorResponse)(response)) {
            throw new Error(JSON.stringify(response.error));
        }
        return response;
    }
}
exports.HttpClient = HttpClient;
//# sourceMappingURL=httpclient.js.map