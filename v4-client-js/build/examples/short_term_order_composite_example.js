"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const order_1 = require("@dydxprotocol/v4-proto/src/codegen/dydxprotocol/clob/order");
const src_1 = require("../src");
const composite_client_1 = require("../src/clients/composite-client");
const constants_1 = require("../src/clients/constants");
const local_wallet_1 = __importDefault(require("../src/clients/modules/local-wallet"));
const subaccount_1 = require("../src/clients/subaccount");
const utils_1 = require("../src/lib/utils");
const constants_2 = require("./constants");
const human_readable_short_term_orders_json_1 = __importDefault(require("./human_readable_short_term_orders.json"));
async function sleep(ms) {
    return new Promise((resolve) => setTimeout(resolve, ms));
}
async function test() {
    var _a;
    const wallet = await local_wallet_1.default.fromMnemonic(constants_2.DYDX_TEST_MNEMONIC, src_1.BECH32_PREFIX);
    console.log(wallet);
    const network = constants_1.Network.testnet();
    const client = await composite_client_1.CompositeClient.connect(network);
    console.log('**Client**');
    console.log(client);
    const subaccount = new subaccount_1.SubaccountInfo(wallet, 0);
    for (const orderParams of human_readable_short_term_orders_json_1.default) {
        try {
            const side = constants_1.OrderSide[orderParams.side];
            const price = (_a = orderParams.price) !== null && _a !== void 0 ? _a : 1350;
            const currentBlock = await client.validatorClient.get.latestBlockHeight();
            const nextValidBlockHeight = currentBlock + 1;
            // Note, you can change this to any number between `next_valid_block_height`
            // to `next_valid_block_height + SHORT_BLOCK_WINDOW`
            const goodTilBlock = nextValidBlockHeight + 10;
            const timeInForce = orderExecutionToTimeInForce(orderParams.timeInForce);
            // uint32
            const clientId = (0, utils_1.randomInt)(2 ** 32 - 1);
            const tx = await client.placeShortTermOrder(subaccount, 'ETH-USD', side, price, 0.01, clientId, goodTilBlock, timeInForce, false);
            console.log('**Order Tx**');
            console.log(tx.hash.toString());
        }
        catch (error) {
            console.log(error.message);
        }
        await sleep(5000); // wait for placeOrder to complete
    }
}
function orderExecutionToTimeInForce(orderExecution) {
    switch (orderExecution) {
        case constants_1.OrderExecution.DEFAULT:
            return order_1.Order_TimeInForce.TIME_IN_FORCE_UNSPECIFIED;
        case constants_1.OrderExecution.FOK:
            return order_1.Order_TimeInForce.TIME_IN_FORCE_FILL_OR_KILL;
        case constants_1.OrderExecution.IOC:
            return order_1.Order_TimeInForce.TIME_IN_FORCE_IOC;
        case constants_1.OrderExecution.POST_ONLY:
            return order_1.Order_TimeInForce.TIME_IN_FORCE_POST_ONLY;
        default:
            throw new Error('Unrecognized order execution');
    }
}
test().then(() => {
}).catch((error) => {
    console.log(error.message);
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoic2hvcnRfdGVybV9vcmRlcl9jb21wb3NpdGVfZXhhbXBsZS5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uL2V4YW1wbGVzL3Nob3J0X3Rlcm1fb3JkZXJfY29tcG9zaXRlX2V4YW1wbGUudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7QUFBQSxzRkFBK0Y7QUFFL0YsZ0NBQXVDO0FBQ3ZDLHNFQUFrRTtBQUNsRSx3REFFa0M7QUFDbEMsdUZBQThEO0FBQzlELDBEQUEyRDtBQUMzRCw0Q0FBNkM7QUFDN0MsMkNBQWlEO0FBQ2pELG9IQUFtRTtBQUVuRSxLQUFLLFVBQVUsS0FBSyxDQUFDLEVBQVU7SUFDN0IsT0FBTyxJQUFJLE9BQU8sQ0FBQyxDQUFDLE9BQU8sRUFBRSxFQUFFLENBQUMsVUFBVSxDQUFDLE9BQU8sRUFBRSxFQUFFLENBQUMsQ0FBQyxDQUFDO0FBQzNELENBQUM7QUFFRCxLQUFLLFVBQVUsSUFBSTs7SUFDakIsTUFBTSxNQUFNLEdBQUcsTUFBTSxzQkFBVyxDQUFDLFlBQVksQ0FBQyw4QkFBa0IsRUFBRSxtQkFBYSxDQUFDLENBQUM7SUFDakYsT0FBTyxDQUFDLEdBQUcsQ0FBQyxNQUFNLENBQUMsQ0FBQztJQUNwQixNQUFNLE9BQU8sR0FBRyxtQkFBTyxDQUFDLE9BQU8sRUFBRSxDQUFDO0lBQ2xDLE1BQU0sTUFBTSxHQUFHLE1BQU0sa0NBQWUsQ0FBQyxPQUFPLENBQUMsT0FBTyxDQUFDLENBQUM7SUFDdEQsT0FBTyxDQUFDLEdBQUcsQ0FBQyxZQUFZLENBQUMsQ0FBQztJQUMxQixPQUFPLENBQUMsR0FBRyxDQUFDLE1BQU0sQ0FBQyxDQUFDO0lBQ3BCLE1BQU0sVUFBVSxHQUFHLElBQUksMkJBQWMsQ0FBQyxNQUFNLEVBQUUsQ0FBQyxDQUFDLENBQUM7SUFDakQsS0FBSyxNQUFNLFdBQVcsSUFBSSwrQ0FBWSxFQUFFLENBQUM7UUFDdkMsSUFBSSxDQUFDO1lBQ0gsTUFBTSxJQUFJLEdBQUcscUJBQVMsQ0FBQyxXQUFXLENBQUMsSUFBOEIsQ0FBQyxDQUFDO1lBQ25FLE1BQU0sS0FBSyxHQUFHLE1BQUEsV0FBVyxDQUFDLEtBQUssbUNBQUksSUFBSSxDQUFDO1lBRXhDLE1BQU0sWUFBWSxHQUFHLE1BQU0sTUFBTSxDQUFDLGVBQWUsQ0FBQyxHQUFHLENBQUMsaUJBQWlCLEVBQUUsQ0FBQztZQUMxRSxNQUFNLG9CQUFvQixHQUFHLFlBQVksR0FBRyxDQUFDLENBQUM7WUFDOUMsNEVBQTRFO1lBQzVFLG9EQUFvRDtZQUNwRCxNQUFNLFlBQVksR0FBRyxvQkFBb0IsR0FBRyxFQUFFLENBQUM7WUFFL0MsTUFBTSxXQUFXLEdBQUcsMkJBQTJCLENBQUMsV0FBVyxDQUFDLFdBQVcsQ0FBQyxDQUFDO1lBRXpFLFNBQVM7WUFDVCxNQUFNLFFBQVEsR0FBRyxJQUFBLGlCQUFTLEVBQUMsQ0FBQyxJQUFJLEVBQUUsR0FBRyxDQUFDLENBQUMsQ0FBQztZQUV4QyxNQUFNLEVBQUUsR0FBRyxNQUFNLE1BQU0sQ0FBQyxtQkFBbUIsQ0FDekMsVUFBVSxFQUNWLFNBQVMsRUFDVCxJQUFJLEVBQ0osS0FBSyxFQUNMLElBQUksRUFDSixRQUFRLEVBQ1IsWUFBWSxFQUNaLFdBQVcsRUFDWCxLQUFLLENBQ04sQ0FBQztZQUNGLE9BQU8sQ0FBQyxHQUFHLENBQUMsY0FBYyxDQUFDLENBQUM7WUFDNUIsT0FBTyxDQUFDLEdBQUcsQ0FBQyxFQUFFLENBQUMsSUFBSSxDQUFDLFFBQVEsRUFBRSxDQUFDLENBQUM7UUFDbEMsQ0FBQztRQUFDLE9BQU8sS0FBSyxFQUFFLENBQUM7WUFDZixPQUFPLENBQUMsR0FBRyxDQUFDLEtBQUssQ0FBQyxPQUFPLENBQUMsQ0FBQztRQUM3QixDQUFDO1FBRUQsTUFBTSxLQUFLLENBQUMsSUFBSSxDQUFDLENBQUMsQ0FBRSxrQ0FBa0M7SUFDeEQsQ0FBQztBQUNILENBQUM7QUFFRCxTQUFTLDJCQUEyQixDQUFDLGNBQXNCO0lBQ3pELFFBQVEsY0FBYyxFQUFFLENBQUM7UUFDdkIsS0FBSywwQkFBYyxDQUFDLE9BQU87WUFDekIsT0FBTyx5QkFBaUIsQ0FBQyx5QkFBeUIsQ0FBQztRQUNyRCxLQUFLLDBCQUFjLENBQUMsR0FBRztZQUNyQixPQUFPLHlCQUFpQixDQUFDLDBCQUEwQixDQUFDO1FBQ3RELEtBQUssMEJBQWMsQ0FBQyxHQUFHO1lBQ3JCLE9BQU8seUJBQWlCLENBQUMsaUJBQWlCLENBQUM7UUFDN0MsS0FBSywwQkFBYyxDQUFDLFNBQVM7WUFDM0IsT0FBTyx5QkFBaUIsQ0FBQyx1QkFBdUIsQ0FBQztRQUNuRDtZQUNFLE1BQU0sSUFBSSxLQUFLLENBQUMsOEJBQThCLENBQUMsQ0FBQztJQUNwRCxDQUFDO0FBQ0gsQ0FBQztBQUVELElBQUksRUFBRSxDQUFDLElBQUksQ0FBQyxHQUFHLEVBQUU7QUFDakIsQ0FBQyxDQUFDLENBQUMsS0FBSyxDQUFDLENBQUMsS0FBSyxFQUFFLEVBQUU7SUFDakIsT0FBTyxDQUFDLEdBQUcsQ0FBQyxLQUFLLENBQUMsT0FBTyxDQUFDLENBQUM7QUFDN0IsQ0FBQyxDQUFDLENBQUMifQ==