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
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoic2hvcnRfdGVybV9vcmRlcl9jYW5jZWxfZXhhbXBsZS5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uL2V4YW1wbGVzL3Nob3J0X3Rlcm1fb3JkZXJfY2FuY2VsX2V4YW1wbGUudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7QUFBQSxnQ0FBc0U7QUFDdEUsc0VBQWtFO0FBQ2xFLHdEQUVrQztBQUNsQyx1RkFBOEQ7QUFDOUQsMERBQTJEO0FBQzNELDRDQUFvRDtBQUNwRCwyQ0FBZ0U7QUFFaEUsS0FBSyxVQUFVLElBQUk7SUFDakIsTUFBTSxNQUFNLEdBQUcsTUFBTSxzQkFBVyxDQUFDLFlBQVksQ0FBQyw4QkFBa0IsRUFBRSxtQkFBYSxDQUFDLENBQUM7SUFDakYsT0FBTyxDQUFDLEdBQUcsQ0FBQyxNQUFNLENBQUMsQ0FBQztJQUNwQixNQUFNLE9BQU8sR0FBRyxtQkFBTyxDQUFDLE9BQU8sRUFBRSxDQUFDO0lBQ2xDLE1BQU0sTUFBTSxHQUFHLE1BQU0sa0NBQWUsQ0FBQyxPQUFPLENBQUMsT0FBTyxDQUFDLENBQUM7SUFDdEQsT0FBTyxDQUFDLEdBQUcsQ0FBQyxZQUFZLENBQUMsQ0FBQztJQUMxQixPQUFPLENBQUMsR0FBRyxDQUFDLE1BQU0sQ0FBQyxDQUFDO0lBQ3BCLE1BQU0sVUFBVSxHQUFHLElBQUksMkJBQWMsQ0FBQyxNQUFNLEVBQUUsQ0FBQyxDQUFDLENBQUM7SUFFakQsTUFBTSxZQUFZLEdBQUcsTUFBTSxNQUFNLENBQUMsZUFBZSxDQUFDLEdBQUcsQ0FBQyxpQkFBaUIsRUFBRSxDQUFDO0lBQzFFLE1BQU0sb0JBQW9CLEdBQUcsWUFBWSxHQUFHLENBQUMsQ0FBQztJQUM5Qyw0RUFBNEU7SUFDNUUsb0RBQW9EO0lBQ3BELE1BQU0sWUFBWSxHQUFHLG9CQUFvQixHQUFHLEVBQUUsQ0FBQztJQUMvQyxNQUFNLHNCQUFzQixHQUFHLElBQUEsaUJBQVMsRUFBQyx5QkFBYSxDQUFDLENBQUM7SUFDeEQsSUFBSTtRQUNGLDJCQUEyQjtRQUMzQixNQUFNLEVBQUUsR0FBRyxNQUFNLE1BQU0sQ0FBQyxtQkFBbUIsQ0FDekMsVUFBVSxFQUNWLFNBQVMsRUFDVCxxQkFBUyxDQUFDLElBQUksRUFDZCxLQUFLLEVBQ0wsSUFBSSxFQUNKLHNCQUFzQixFQUN0QixZQUFZLEVBQ1osdUJBQWlCLENBQUMseUJBQXlCLEVBQzNDLEtBQUssQ0FDTixDQUFDO1FBQ0YsT0FBTyxDQUFDLEdBQUcsQ0FBQyx5QkFBeUIsQ0FBQyxDQUFDO1FBQ3ZDLE9BQU8sQ0FBQyxHQUFHLENBQUMsRUFBRSxDQUFDLElBQUksQ0FBQyxDQUFDO0tBQ3RCO0lBQUMsT0FBTyxLQUFLLEVBQUU7UUFDZCxPQUFPLENBQUMsR0FBRyxDQUFDLDZCQUE2QixDQUFDLENBQUM7UUFDM0MsT0FBTyxDQUFDLEdBQUcsQ0FBQyxLQUFLLENBQUMsT0FBTyxDQUFDLENBQUM7S0FDNUI7SUFFRCxNQUFNLElBQUEsYUFBSyxFQUFDLElBQUksQ0FBQyxDQUFDLENBQUUsa0NBQWtDO0lBRXRELElBQUk7UUFDRiw4QkFBOEI7UUFDOUIsTUFBTSxFQUFFLEdBQUcsTUFBTSxNQUFNLENBQUMsV0FBVyxDQUNqQyxVQUFVLEVBQ1Ysc0JBQXNCLEVBQ3RCLGdCQUFVLENBQUMsVUFBVSxFQUNyQixTQUFTLEVBQ1QsWUFBWSxHQUFHLEVBQUUsRUFDakIsQ0FBQyxDQUNGLENBQUM7UUFDRixPQUFPLENBQUMsR0FBRyxDQUFDLGdDQUFnQyxDQUFDLENBQUM7UUFDOUMsT0FBTyxDQUFDLEdBQUcsQ0FBQyxFQUFFLENBQUMsQ0FBQztLQUNqQjtJQUFDLE9BQU8sS0FBSyxFQUFFO1FBQ2QsT0FBTyxDQUFDLEdBQUcsQ0FBQyxvQ0FBb0MsQ0FBQyxDQUFDO1FBQ2xELE9BQU8sQ0FBQyxHQUFHLENBQUMsS0FBSyxDQUFDLE9BQU8sQ0FBQyxDQUFDO0tBQzVCO0FBQ0gsQ0FBQztBQUVELElBQUksRUFBRSxDQUFDLElBQUksQ0FBQyxHQUFHLEVBQUU7QUFDakIsQ0FBQyxDQUFDLENBQUMsS0FBSyxDQUFDLENBQUMsS0FBSyxFQUFFLEVBQUU7SUFDakIsT0FBTyxDQUFDLEdBQUcsQ0FBQyxLQUFLLENBQUMsT0FBTyxDQUFDLENBQUM7QUFDN0IsQ0FBQyxDQUFDLENBQUMifQ==