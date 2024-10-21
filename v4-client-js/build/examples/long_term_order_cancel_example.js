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
async function sleep(ms) {
    return new Promise((resolve) => setTimeout(resolve, ms));
}
async function test() {
    const wallet = await local_wallet_1.default.fromMnemonic(constants_2.DYDX_TEST_MNEMONIC, src_1.BECH32_PREFIX);
    console.log(wallet);
    const network = constants_1.Network.testnet();
    const client = await composite_client_1.CompositeClient.connect(network);
    console.log('**Client**');
    console.log(client);
    const subaccount = new subaccount_1.SubaccountInfo(wallet, 0);
    /*
    Note this example places a stateful order.
    Programmatic traders should generally not use stateful orders for following reasons:
    - Stateful orders received out of order by validators will fail sequence number validation
      and be dropped.
    - Stateful orders have worse time priority since they are only matched after they are included
      on the block.
    - Stateful order rate limits are more restrictive than Short-Term orders, specifically max 2 per
      block / 20 per 100 blocks.
    - Stateful orders can only be canceled after they’ve been included in a block.
    */
    const longTermOrderClientId = (0, utils_1.randomInt)(constants_2.MAX_CLIENT_ID);
    try {
        // place a long term order
        const tx = await client.placeOrder(subaccount, 'ETH-USD', constants_1.OrderType.LIMIT, constants_1.OrderSide.SELL, 40000, 0.01, longTermOrderClientId, constants_1.OrderTimeInForce.GTT, 60, constants_1.OrderExecution.DEFAULT, false, false);
        console.log('**Long Term Order Tx**');
        console.log(tx.hash);
    }
    catch (error) {
        console.log('**Long Term Order Failed**');
        console.log(error.message);
    }
    await sleep(5000); // wait for placeOrder to complete
    try {
        // cancel the long term order
        const tx = await client.cancelOrder(subaccount, longTermOrderClientId, src_1.OrderFlags.LONG_TERM, 'ETH-USD', 0, 120);
        console.log('**Cancel Long Term Order Tx**');
        console.log(tx);
    }
    catch (error) {
        console.log('**Cancel Long Term Order Failed**');
        console.log(error.message);
    }
}
test().then(() => {
}).catch((error) => {
    console.log(error.message);
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoibG9uZ190ZXJtX29yZGVyX2NhbmNlbF9leGFtcGxlLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vZXhhbXBsZXMvbG9uZ190ZXJtX29yZGVyX2NhbmNlbF9leGFtcGxlLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7Ozs7O0FBQUEsZ0NBQW1EO0FBQ25ELHNFQUFrRTtBQUNsRSx3REFFa0M7QUFDbEMsdUZBQThEO0FBQzlELDBEQUEyRDtBQUMzRCw0Q0FBNkM7QUFDN0MsMkNBQWdFO0FBRWhFLEtBQUssVUFBVSxLQUFLLENBQUMsRUFBVTtJQUM3QixPQUFPLElBQUksT0FBTyxDQUFDLENBQUMsT0FBTyxFQUFFLEVBQUUsQ0FBQyxVQUFVLENBQUMsT0FBTyxFQUFFLEVBQUUsQ0FBQyxDQUFDLENBQUM7QUFDM0QsQ0FBQztBQUVELEtBQUssVUFBVSxJQUFJO0lBQ2pCLE1BQU0sTUFBTSxHQUFHLE1BQU0sc0JBQVcsQ0FBQyxZQUFZLENBQUMsOEJBQWtCLEVBQUUsbUJBQWEsQ0FBQyxDQUFDO0lBQ2pGLE9BQU8sQ0FBQyxHQUFHLENBQUMsTUFBTSxDQUFDLENBQUM7SUFDcEIsTUFBTSxPQUFPLEdBQUcsbUJBQU8sQ0FBQyxPQUFPLEVBQUUsQ0FBQztJQUNsQyxNQUFNLE1BQU0sR0FBRyxNQUFNLGtDQUFlLENBQUMsT0FBTyxDQUFDLE9BQU8sQ0FBQyxDQUFDO0lBQ3RELE9BQU8sQ0FBQyxHQUFHLENBQUMsWUFBWSxDQUFDLENBQUM7SUFDMUIsT0FBTyxDQUFDLEdBQUcsQ0FBQyxNQUFNLENBQUMsQ0FBQztJQUNwQixNQUFNLFVBQVUsR0FBRyxJQUFJLDJCQUFjLENBQUMsTUFBTSxFQUFFLENBQUMsQ0FBQyxDQUFDO0lBRWpEOzs7Ozs7Ozs7O01BVUU7SUFDRixNQUFNLHFCQUFxQixHQUFHLElBQUEsaUJBQVMsRUFBQyx5QkFBYSxDQUFDLENBQUM7SUFDdkQsSUFBSSxDQUFDO1FBQ0gsMEJBQTBCO1FBQzFCLE1BQU0sRUFBRSxHQUFHLE1BQU0sTUFBTSxDQUFDLFVBQVUsQ0FDaEMsVUFBVSxFQUNWLFNBQVMsRUFDVCxxQkFBUyxDQUFDLEtBQUssRUFDZixxQkFBUyxDQUFDLElBQUksRUFDZCxLQUFLLEVBQ0wsSUFBSSxFQUNKLHFCQUFxQixFQUNyQiw0QkFBZ0IsQ0FBQyxHQUFHLEVBQ3BCLEVBQUUsRUFDRiwwQkFBYyxDQUFDLE9BQU8sRUFDdEIsS0FBSyxFQUNMLEtBQUssQ0FDTixDQUFDO1FBQ0YsT0FBTyxDQUFDLEdBQUcsQ0FBQyx3QkFBd0IsQ0FBQyxDQUFDO1FBQ3RDLE9BQU8sQ0FBQyxHQUFHLENBQUMsRUFBRSxDQUFDLElBQUksQ0FBQyxDQUFDO0lBQ3ZCLENBQUM7SUFBQyxPQUFPLEtBQUssRUFBRSxDQUFDO1FBQ2YsT0FBTyxDQUFDLEdBQUcsQ0FBQyw0QkFBNEIsQ0FBQyxDQUFDO1FBQzFDLE9BQU8sQ0FBQyxHQUFHLENBQUMsS0FBSyxDQUFDLE9BQU8sQ0FBQyxDQUFDO0lBQzdCLENBQUM7SUFFRCxNQUFNLEtBQUssQ0FBQyxJQUFJLENBQUMsQ0FBQyxDQUFFLGtDQUFrQztJQUV0RCxJQUFJLENBQUM7UUFDSCw2QkFBNkI7UUFDN0IsTUFBTSxFQUFFLEdBQUcsTUFBTSxNQUFNLENBQUMsV0FBVyxDQUNqQyxVQUFVLEVBQ1YscUJBQXFCLEVBQ3JCLGdCQUFVLENBQUMsU0FBUyxFQUNwQixTQUFTLEVBQ1QsQ0FBQyxFQUNELEdBQUcsQ0FDSixDQUFDO1FBQ0YsT0FBTyxDQUFDLEdBQUcsQ0FBQywrQkFBK0IsQ0FBQyxDQUFDO1FBQzdDLE9BQU8sQ0FBQyxHQUFHLENBQUMsRUFBRSxDQUFDLENBQUM7SUFDbEIsQ0FBQztJQUFDLE9BQU8sS0FBSyxFQUFFLENBQUM7UUFDZixPQUFPLENBQUMsR0FBRyxDQUFDLG1DQUFtQyxDQUFDLENBQUM7UUFDakQsT0FBTyxDQUFDLEdBQUcsQ0FBQyxLQUFLLENBQUMsT0FBTyxDQUFDLENBQUM7SUFDN0IsQ0FBQztBQUNILENBQUM7QUFFRCxJQUFJLEVBQUUsQ0FBQyxJQUFJLENBQUMsR0FBRyxFQUFFO0FBQ2pCLENBQUMsQ0FBQyxDQUFDLEtBQUssQ0FBQyxDQUFDLEtBQUssRUFBRSxFQUFFO0lBQ2pCLE9BQU8sQ0FBQyxHQUFHLENBQUMsS0FBSyxDQUFDLE9BQU8sQ0FBQyxDQUFDO0FBQzdCLENBQUMsQ0FBQyxDQUFDIn0=