"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const long_1 = __importDefault(require("long"));
const protobufjs_1 = __importDefault(require("protobufjs"));
const src_1 = require("../src");
const constants_1 = require("../src/clients/constants");
const local_wallet_1 = __importDefault(require("../src/clients/modules/local-wallet"));
const subaccount_1 = require("../src/clients/subaccount");
const validator_client_1 = require("../src/clients/validator-client");
const utils_1 = require("../src/lib/utils");
const constants_2 = require("./constants");
const raw_orders_json_1 = __importDefault(require("./raw_orders.json"));
// Required for encoding and decoding queries that are of type Long.
// Must be done once but since the individal modules should be usable
// - must be set in each module that encounters encoding/decoding Longs.
// Reference: https://github.com/protobufjs/protobuf.js/issues/921
protobufjs_1.default.util.Long = long_1.default;
protobufjs_1.default.configure();
function dummyOrder(height) {
    const placeOrder = constants_2.defaultOrder;
    placeOrder.clientId = (0, utils_1.randomInt)(1000000000);
    placeOrder.goodTilBlock = height + 3;
    return placeOrder;
}
async function sleep(ms) {
    return new Promise((resolve) => setTimeout(resolve, ms));
}
async function test() {
    const wallet = await local_wallet_1.default.fromMnemonic(constants_2.DYDX_TEST_MNEMONIC, src_1.BECH32_PREFIX);
    console.log(wallet);
    const client = await validator_client_1.ValidatorClient.connect(constants_1.Network.testnet().validatorConfig);
    console.log('**Client**');
    console.log(client);
    const value1 = long_1.default.fromNumber(400000000000);
    console.log(value1.toString());
    const subaccount = new subaccount_1.SubaccountInfo(wallet, 0);
    for (const orderParams of raw_orders_json_1.default) {
        const height = await client.get.latestBlockHeight();
        const placeOrder = dummyOrder(height);
        placeOrder.timeInForce = orderParams.timeInForce;
        placeOrder.reduceOnly = false; // reduceOnly is currently disabled
        placeOrder.orderFlags = orderParams.orderFlags;
        placeOrder.side = orderParams.side;
        placeOrder.quantums = long_1.default.fromNumber(orderParams.quantums);
        placeOrder.subticks = long_1.default.fromNumber(orderParams.subticks);
        try {
            if (placeOrder.orderFlags !== 0) {
                placeOrder.goodTilBlock = 0;
                const now = new Date();
                const millisecondsPerSecond = 1000;
                const interval = 60 * millisecondsPerSecond;
                const future = new Date(now.valueOf() + interval);
                placeOrder.goodTilBlockTime = Math.round(future.getTime() / 1000);
            }
            else {
                placeOrder.goodTilBlockTime = 0;
            }
            const tx = await client.post.placeOrderObject(subaccount, placeOrder);
            console.log('**Order Tx**');
            console.log(tx);
        }
        catch (error) {
            console.log(error.message);
        }
        await sleep(5000); // wait for placeOrder to complete
    }
}
test().then(() => {
}).catch((error) => {
    console.log(error.message);
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoidmFsaWRhdG9yX3Bvc3RfZXhhbXBsZS5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uL2V4YW1wbGVzL3ZhbGlkYXRvcl9wb3N0X2V4YW1wbGUudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7QUFBQSxnREFBd0I7QUFDeEIsNERBQWtDO0FBRWxDLGdDQUF1QztBQUN2Qyx3REFBbUQ7QUFDbkQsdUZBQThEO0FBQzlELDBEQUEyRDtBQUUzRCxzRUFBa0U7QUFDbEUsNENBQTZDO0FBQzdDLDJDQUErRDtBQUMvRCx3RUFBNkM7QUFFN0Msb0VBQW9FO0FBQ3BFLHFFQUFxRTtBQUNyRSx3RUFBd0U7QUFDeEUsa0VBQWtFO0FBQ2xFLG9CQUFRLENBQUMsSUFBSSxDQUFDLElBQUksR0FBRyxjQUFJLENBQUM7QUFDMUIsb0JBQVEsQ0FBQyxTQUFTLEVBQUUsQ0FBQztBQUVyQixTQUFTLFVBQVUsQ0FBQyxNQUFjO0lBQ2hDLE1BQU0sVUFBVSxHQUFHLHdCQUFZLENBQUM7SUFDaEMsVUFBVSxDQUFDLFFBQVEsR0FBRyxJQUFBLGlCQUFTLEVBQUMsVUFBVSxDQUFDLENBQUM7SUFDNUMsVUFBVSxDQUFDLFlBQVksR0FBRyxNQUFNLEdBQUcsQ0FBQyxDQUFDO0lBQ3JDLE9BQU8sVUFBVSxDQUFDO0FBQ3BCLENBQUM7QUFFRCxLQUFLLFVBQVUsS0FBSyxDQUFDLEVBQVU7SUFDN0IsT0FBTyxJQUFJLE9BQU8sQ0FBQyxDQUFDLE9BQU8sRUFBRSxFQUFFLENBQUMsVUFBVSxDQUFDLE9BQU8sRUFBRSxFQUFFLENBQUMsQ0FBQyxDQUFDO0FBQzNELENBQUM7QUFFRCxLQUFLLFVBQVUsSUFBSTtJQUNqQixNQUFNLE1BQU0sR0FBRyxNQUFNLHNCQUFXLENBQUMsWUFBWSxDQUFDLDhCQUFrQixFQUFFLG1CQUFhLENBQUMsQ0FBQztJQUNqRixPQUFPLENBQUMsR0FBRyxDQUFDLE1BQU0sQ0FBQyxDQUFDO0lBRXBCLE1BQU0sTUFBTSxHQUFHLE1BQU0sa0NBQWUs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