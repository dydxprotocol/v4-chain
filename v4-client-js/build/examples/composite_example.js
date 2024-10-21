"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const src_1 = require("../src");
const composite_client_1 = require("../src/clients/composite-client");
const constants_1 = require("../src/clients/constants");
const local_wallet_1 = __importDefault(require("../src/clients/modules/local-wallet"));
const subaccount_1 = require("../src/clients/subaccount");
const utils_1 = require("../src/lib/utils");
const constants_2 = require("./constants");
const human_readable_orders_json_1 = __importDefault(require("./human_readable_orders.json"));
async function sleep(ms) {
    return new Promise((resolve) => setTimeout(resolve, ms));
}
async function test() {
    var _a, _b, _c;
    const wallet = await local_wallet_1.default.fromMnemonic(constants_2.DYDX_TEST_MNEMONIC, src_1.BECH32_PREFIX);
    console.log(wallet);
    const network = constants_1.Network.testnet();
    const client = await composite_client_1.CompositeClient.connect(network);
    console.log('**Client**');
    console.log(client);
    const subaccount = new subaccount_1.SubaccountInfo(wallet, 0);
    for (const orderParams of human_readable_orders_json_1.default) {
        try {
            const type = constants_1.OrderType[orderParams.type];
            const side = constants_1.OrderSide[orderParams.side];
            const timeInForceString = (_a = orderParams.timeInForce) !== null && _a !== void 0 ? _a : 'GTT';
            const timeInForce = constants_1.OrderTimeInForce[timeInForceString];
            const price = (_b = orderParams.price) !== null && _b !== void 0 ? _b : 1350;
            const timeInForceSeconds = (timeInForce === constants_1.OrderTimeInForce.GTT) ? 60 : 0;
            const postOnly = (_c = orderParams.postOnly) !== null && _c !== void 0 ? _c : false;
            const tx = await client.placeOrder(subaccount, 'ETH-USD', type, side, price, 0.01, (0, utils_1.randomInt)(constants_2.MAX_CLIENT_ID), timeInForce, timeInForceSeconds, constants_1.OrderExecution.DEFAULT, postOnly, false);
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
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiY29tcG9zaXRlX2V4YW1wbGUuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi9leGFtcGxlcy9jb21wb3NpdGVfZXhhbXBsZS50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7OztBQUFBLGdDQUF1QztBQUN2QyxzRUFBa0U7QUFDbEUsd0RBRWtDO0FBQ2xDLHVGQUE4RDtBQUM5RCwwREFBMkQ7QUFDM0QsNENBQTZDO0FBQzdDLDJDQUFnRTtBQUNoRSw4RkFBd0Q7QUFFeEQsS0FBSyxVQUFVLEtBQUssQ0FBQyxFQUFVO0lBQzdCLE9BQU8sSUFBSSxPQUFPLENBQUMsQ0FBQyxPQUFPLEVBQUUsRUFBRSxDQUFDLFVBQVUsQ0FBQyxPQUFPLEVBQUUsRUFBRSxDQUFDLENBQUMsQ0FBQztBQUMzRCxDQUFDO0FBRUQsS0FBSyxVQUFVLElBQUk7O0lBQ2pCLE1BQU0sTUFBTSxHQUFHLE1BQU0sc0JBQVcsQ0FBQyxZQUFZLENBQUMsOEJBQWtCLEVBQUUsbUJBQWEsQ0FBQyxDQUFDO0lBQ2pGLE9BQU8sQ0FBQyxHQUFHLENBQUMsTUFBTSxDQUFDLENBQUM7SUFDcEIsTUFBTSxPQUFPLEdBQUcsbUJBQU8sQ0FBQyxPQUFPLEVBQUUsQ0FBQztJQUNsQyxNQUFNLE1BQU0sR0FBRyxNQUFNLGtDQUFlLENBQUMsT0FBTyxDQUFDLE9BQU8sQ0FBQyxDQUFDO0lBQ3RELE9BQU8sQ0FBQyxHQUFHLENBQUMsWUFBWSxDQUFDLENBQUM7SUFDMUIsT0FBTyxDQUFDLEdBQUcsQ0FBQyxNQUFNLENBQUMsQ0FBQztJQUNwQixNQUFNLFVBQVUsR0FBRyxJQUFJLDJCQUFjLENBQUMsTUFBTSxFQUFFLENBQUMsQ0FBQyxDQUFDO0lBQ2pELEtBQUssTUFBTSxXQUFXLElBQUksb0NBQVksRUFBRSxDQUFDO1FBQ3ZDLElBQUksQ0FBQztZQUNILE1BQU0sSUFBSSxHQUFHLHFCQUFTLENBQUMsV0FBVyxDQUFDLElBQThCLENBQUMsQ0FBQztZQUNuRSxNQUFNLElBQUksR0FBRyxxQkFBUyxDQUFDLFdBQVcsQ0FBQyxJQUE4QixDQUFDLENBQUM7WUFDbkUsTUFBTSxpQkFBaUIsR0FBRyxNQUFBLFdBQVcsQ0FBQyxXQUFXLG1DQUFJLEtBQUssQ0FBQztZQUMzRCxNQUFNLFdBQVcsR0FBRyw0QkFBZ0IsQ0FBQyxpQkFBa0QsQ0FBQyxDQUFDO1lBQ3pGLE1BQU0sS0FBSyxHQUFHLE1BQUEsV0FBVyxDQUFDLEtBQUssbUNBQUksSUFBSSxDQUFDO1lBQ3hDLE1BQU0sa0JBQWtCLEdBQUcsQ0FBQyxXQUFXLEtBQUssNEJBQWdCLENBQUMsR0FBRyxDQUFDLENBQUMsQ0FBQyxDQUFDLEVBQUUsQ0FBQyxDQUFDLENBQUMsQ0FBQyxDQUFDO1lBQzNFLE1BQU0sUUFBUSxHQUFHLE1BQUEsV0FBVyxDQUFDLFFBQVEsbUNBQUksS0FBSyxDQUFDO1lBQy9DLE1BQU0sRUFBRSxHQUFHLE1BQU0sTUFBTSxDQUFDLFVBQVUsQ0FDaEMsVUFBVSxFQUNWLFNBQVMsRUFDVCxJQUFJLEVBQ0osSUFBSSxFQUNKLEtBQUssRUFDTCxJQUFJLEVBQ0osSUFBQSxpQkFBUyxFQUFDLHlCQUFhLENBQUMsRUFDeEIsV0FBVyxFQUNYLGtCQUFrQixFQUNsQiwwQkFBYyxDQUFDLE9BQU8sRUFDdEIsUUFBUSxFQUNSLEtBQUssQ0FDTixDQUFDO1lBQ0YsT0FBTyxDQUFDLEdBQUcsQ0FBQyxjQUFjLENBQUMsQ0FBQztZQUM1QixPQUFPLENBQUMsR0FBRyxDQUFDLEVBQUUsQ0FBQyxDQUFDO1FBQ2xCLENBQUM7UUFBQyxPQUFPLEtBQUssRUFBRSxDQUFDO1lBQ2YsT0FBTyxDQUFDLEdBQUcsQ0FBQyxLQUFLLENBQUMsT0FBTyxDQUFDLENBQUM7UUFDN0IsQ0FBQztRQUVELE1BQU0sS0FBSyxDQUFDLElBQUksQ0FBQyxDQUFDLENBQUUsa0NBQWtDO0lBQ3hELENBQUM7QUFDSCxDQUFDO0FBRUQsSUFBSSxFQUFFLENBQUMsSUFBSSxDQUFDLEdBQUcsRUFBRTtBQUNqQixDQUFDLENBQUMsQ0FBQyxLQUFLLENBQUMsQ0FBQyxLQUFLLEVBQUUsRUFBRTtJQUNqQixPQUFPLENBQUMsR0FBRyxDQUFDLEtBQUssQ0FBQyxPQUFPLENBQUMsQ0FBQztBQUM3QixDQUFDLENBQUMsQ0FBQyJ9