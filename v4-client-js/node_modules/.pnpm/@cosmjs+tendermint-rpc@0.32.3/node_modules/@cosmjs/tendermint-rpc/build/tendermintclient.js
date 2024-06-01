"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.connectComet = exports.isComet38Client = exports.isTendermint37Client = exports.isTendermint34Client = void 0;
const comet38_1 = require("./comet38");
const tendermint34_1 = require("./tendermint34");
const tendermint37_1 = require("./tendermint37");
function isTendermint34Client(client) {
    return client instanceof tendermint34_1.Tendermint34Client;
}
exports.isTendermint34Client = isTendermint34Client;
function isTendermint37Client(client) {
    return client instanceof tendermint37_1.Tendermint37Client;
}
exports.isTendermint37Client = isTendermint37Client;
function isComet38Client(client) {
    return client instanceof comet38_1.Comet38Client;
}
exports.isComet38Client = isComet38Client;
/**
 * Auto-detects the version of the backend and uses a suitable client.
 */
async function connectComet(endpoint) {
    // Tendermint/CometBFT 0.34/0.37/0.38 auto-detection. Starting with 0.37 we seem to get reliable versions again 🎉
    // Using 0.34 as the fallback.
    let out;
    const tm37Client = await tendermint37_1.Tendermint37Client.connect(endpoint);
    const version = (await tm37Client.status()).nodeInfo.version;
    if (version.startsWith("0.37.")) {
        out = tm37Client;
    }
    else if (version.startsWith("0.38.")) {
        tm37Client.disconnect();
        out = await comet38_1.Comet38Client.connect(endpoint);
    }
    else {
        tm37Client.disconnect();
        out = await tendermint34_1.Tendermint34Client.connect(endpoint);
    }
    return out;
}
exports.connectComet = connectComet;
//# sourceMappingURL=tendermintclient.js.map