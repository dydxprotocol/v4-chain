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
async function test() {
    const wallet = await local_wallet_1.default.fromMnemonic(constants_2.DYDX_TEST_MNEMONIC, src_1.BECH32_PREFIX);
    console.log(wallet);
    const network = constants_1.Network.testnet();
    const client = await composite_client_1.CompositeClient.connect(network);
    console.log('**Client**');
    console.log(client);
    const subaccount = new subaccount_1.SubaccountInfo(wallet, 0);
    const currentBlock = await client.validatorClient.get.latestBlockHeight();
    const nextValidBlockHeight = currentBlock + 1;
    // Note, you can change this to any number between `next_valid_block_height`
    // to `next_valid_block_height + SHORT_BLOCK_WINDOW`
    const goodTilBlock = nextValidBlockHeight + 10;
    const shortTermOrderClientId = (0, utils_1.randomInt)(constants_2.MAX_CLIENT_ID);
    try {
        // place a short term order
        const tx = await client.placeShortTermOrder(subaccount, 'ETH-USD', constants_1.OrderSide.SELL, 40000, 0.01, shortTermOrderClientId, goodTilBlock, src_1.Order_TimeInForce.TIME_IN_FORCE_UNSPECIFIED, false);
        console.log('**Short Term Order Tx**');
        console.log(tx.hash);
    }
    catch (error) {
        console.log('**Short Term Order Failed**');
        console.log(error.message);
    }
    await (0, utils_1.sleep)(5000); // wait for placeOrder to complete
    try {
        // cancel the short term order
        const tx = await client.cancelOrder(subaccount, shortTermOrderClientId, src_1.OrderFlags.SHORT_TERM, 'ETH-USD', goodTilBlock + 10, 0);
        console.log('**Cancel Short Term Order Tx**');
        console.log(tx);
    }
    catch (error) {
        console.log('**Cancel Short Term Order Failed**');
        console.log(error.message);
    }
}
test().then(() => {
}).catch((error) => {
    console.log(error.message);
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoic2hvcnRfdGVybV9vcmRlcl9jYW5jZWxfZXhhbXBsZS5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uL2V4YW1wbGVzL3Nob3J0X3Rlcm1fb3JkZXJfY2FuY2VsX2V4YW1wbGUudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7QUFBQSxnQ0FBc0U7QUFDdEUsc0VBQWtFO0FBQ2xFLHdEQUVrQztBQUNsQyx1RkFBOEQ7QUFDOUQsMERBQTJEO0FBQzNELDRDQUFvRDtBQUNwRCwyQ0FBZ0U7QUFFaEUsS0FBSyxVQUFVLElBQUk7SUFDakIsTUFBTSxNQUFNLEdBQUcsTUFBTSxzQkFBVyxDQUFDLFlBQVksQ0FBQyw4QkFBa0IsRUFBRSxtQkFBYSxDQUFDLENBQUM7SUFDakYsT0FBTyxDQUFDLEdBQUcsQ0FBQyxNQUFNLENBQUMsQ0FBQztJQUNwQixNQUFNLE9BQU8sR0FBRyxtQkFBTyxDQUFDLE9BQU8sRUFBRSxDQUFDO0lBQ2xDLE1BQU0sTUFBTSxHQUFHLE1BQU0sa0NBQWUsQ0FBQyxPQUFPLENBQUMsT0FBTyxDQUFDLENBQUM7SUFDdEQsT0FBTyxDQUFDLEdBQUcsQ0FBQyxZQUFZLENBQUMsQ0FBQztJQUMxQixPQUFPLENBQUMsR0FBRyxDQUFDLE1BQU0sQ0FBQyxDQUFDO0lBQ3BCLE1BQU0sVUFBVSxHQUFHLElBQUksMkJBQWMsQ0FBQyxNQUFNLEVBQUUsQ0FBQyxDQUFDLENBQUM7SUFFakQsTUFBTSxZQUFZLEdBQUcsTUFBTSxNQUFNLENBQUMsZUFBZSxDQUFDLEdBQUcsQ0FBQyxpQkFBaUIsRUFBRSxDQUFDO0lBQzFFLE1BQU0sb0JBQW9CLEdBQUcsWUFBWSxHQUFHLENBQUMsQ0FBQztJQUM5Qyw0RUFBNEU7SUFDNUUsb0RBQW9EO0lBQ3BELE1BQU0sWUFBWSxHQUFHLG9CQUFvQixHQUFHLEVBQUUsQ0FBQztJQUMvQyxNQUFNLHNCQUFzQixHQUFHLElBQUEsaUJBQVMsRUFBQyx5QkFBYSxDQUFDLENBQUM7SUFDeEQsSUFBSSxDQUFDO1FBQ0gsMkJBQTJCO1FBQzNCLE1BQU0sRUFBRSxHQUFHLE1BQU0sTUFBTSxDQUFDLG1CQUFtQixDQUN6QyxVQUFVLEVBQ1YsU0FBUyxFQUNULHFCQUFTLENBQUMsSUFBSSxFQUNkLEtBQUssRUFDTCxJQUFJLEVBQ0osc0JBQXNCLEVBQ3RCLFlBQVksRUFDWix1QkFBaUIsQ0FBQyx5QkFBeUIsRUFDM0MsS0FBSyxDQUNOLENBQUM7UUFDRixPQUFPLENBQUMsR0FBRyxDQUFDLHlCQUF5QixDQUFDLENBQUM7UUFDdkMsT0FBTyxDQUFDLEdBQUcsQ0FBQyxFQUFFLENBQUMsSUFBSSxDQUFDLENBQUM7SUFDdkIsQ0FBQztJQUFDLE9BQU8sS0FBSyxFQUFFLENBQUM7UUFDZixPQUFPLENBQUMsR0FBRyxDQUFDLDZCQUE2QixDQUFDLENBQUM7UUFDM0MsT0FBTyxDQUFDLEdBQUcsQ0FBQyxLQUFLLENBQUMsT0FBTyxDQUFDLENBQUM7SUFDN0IsQ0FBQztJQUVELE1BQU0sSUFBQSxhQUFLLEVBQUMsSUFBSSxDQUFDLENBQUMsQ0FBRSxrQ0FBa0M7SUFFdEQsSUFBSSxDQUFDO1FBQ0gsOEJBQThCO1FBQzlCLE1BQU0sRUFBRSxHQUFHLE1BQU0sTUFBTSxDQUFDLFdBQVcsQ0FDakMsVUFBVSxFQUNWLHNCQUFzQixFQUN0QixnQkFBVSxDQUFDLFVBQVUsRUFDckIsU0FBUyxFQUNULFlBQVksR0FBRyxFQUFFLEVBQ2pCLENBQUMsQ0FDRixDQUFDO1FBQ0YsT0FBTyxDQUFDLEdBQUcsQ0FBQyxnQ0FBZ0MsQ0FBQyxDQUFDO1FBQzlDLE9BQU8sQ0FBQyxHQUFHLENBQUMsRUFBRSxDQUFDLENBQUM7SUFDbEIsQ0FBQztJQUFDLE9BQU8sS0FBSyxFQUFFLENBQUM7UUFDZixPQUFPLENBQUMsR0FBRyxDQUFDLG9DQUFvQyxDQUFDLENBQUM7UUFDbEQsT0FBTyxDQUFDLEdBQUcsQ0FBQyxLQUFLLENBQUMsT0FBTyxDQUFDLENBQUM7SUFDN0IsQ0FBQztBQUNILENBQUM7QUFFRCxJQUFJLEVBQUUsQ0FBQyxJQUFJLENBQUMsR0FBRyxFQUFFO0FBQ2pCLENBQUMsQ0FBQyxDQUFDLEtBQUssQ0FBQyxDQUFDLEtBQUssRUFBRSxFQUFFO0lBQ2pCLE9BQU8sQ0FBQyxHQUFHLENBQUMsS0FBSyxDQUFDLE9BQU8sQ0FBQyxDQUFDO0FBQzdCLENBQUMsQ0FBQyxDQUFDIn0=