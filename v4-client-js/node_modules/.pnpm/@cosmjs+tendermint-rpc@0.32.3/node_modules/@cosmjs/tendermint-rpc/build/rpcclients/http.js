"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.http = void 0;
const axios_1 = __importDefault(require("axios"));
function filterBadStatus(res) {
    if (res.status >= 400) {
        throw new Error(`Bad status on response: ${res.status}`);
    }
    return res;
}
/**
 * Node.js 18 comes with exprimental fetch support (https://nodejs.org/de/blog/announcements/v18-release-announce/).
 * This is nice, but the implementation does not yet work wekk for us. We
 * can just stick with axios on those systems for now.
 */
// eslint-disable-next-line @typescript-eslint/ban-types
function isExperimental(nodeJsFunc) {
    // This works because we get this info in node 18:
    //
    // > fetch.toString()
    // 'async function fetch(input, init = undefined) {\n' +
    // "    emitExperimentalWarning('The Fetch API');\n" +
    // '    return lazyUndici().fetch(input, init);\n' +
    // '  }'
    return nodeJsFunc.toString().includes("emitExperimentalWarning");
}
/**
 * Helper to work around missing CORS support in Tendermint (https://github.com/tendermint/tendermint/pull/2800)
 *
 * For some reason, fetch does not complain about missing server-side CORS support.
 */
// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types
async function http(method, url, headers, request) {
    if (typeof fetch === "function" && !isExperimental(fetch)) {
        const settings = {
            method: method,
            body: request ? JSON.stringify(request) : undefined,
            headers: {
                // eslint-disable-next-line @typescript-eslint/naming-convention
                "Content-Type": "application/json",
                ...headers,
            },
        };
        return fetch(url, settings)
            .then(filterBadStatus)
            .then((res) => res.json());
    }
    else {
        return axios_1.default
            .request({ url: url, method: method, data: request, headers: headers })
            .then((res) => res.data);
    }
}
exports.http = http;
//# sourceMappingURL=http.js.map