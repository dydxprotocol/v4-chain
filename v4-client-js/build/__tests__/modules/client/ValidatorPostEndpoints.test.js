"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const order_1 = require("@dydxprotocol/v4-proto/src/codegen/dydxprotocol/clob/order");
const constants_1 = require("../../../src/clients/constants");
const local_wallet_1 = __importDefault(require("../../../src/clients/modules/local-wallet"));
const subaccount_1 = require("../../../src/clients/subaccount");
const validator_client_1 = require("../../../src/clients/validator-client");
const utils_1 = require("../../../src/lib/utils");
const constants_2 = require("../../../examples/constants");
const src_1 = require("../../../src");
function dummyOrder(height) {
    const placeOrder = constants_2.defaultOrder;
    placeOrder.clientId = (0, utils_1.randomInt)(1000000000);
    placeOrder.goodTilBlock = height + 3;
    // placeOrder.goodTilBlockTime = height + 3;
    const random = (0, utils_1.randomInt)(1000);
    if ((random % 2) === 0) {
        placeOrder.side = order_1.Order_Side.SIDE_BUY;
    }
    else {
        placeOrder.side = order_1.Order_Side.SIDE_SELL;
    }
    return placeOrder;
}
describe('Validator Client', () => {
    let wallet;
    let client;
    describe('Post', () => {
        beforeEach(async () => {
            wallet = await local_wallet_1.default.fromMnemonic(constants_2.DYDX_TEST_MNEMONIC, src_1.BECH32_PREFIX);
            client = await validator_client_1.ValidatorClient.connect(constants_1.Network.testnet().validatorConfig);
        });
        it('PlaceOrder', async () => {
            console.log('**Client**');
            console.log(client);
            const address = wallet.address;
            const account = await client.get.getAccount(address);
            console.log('**Account**');
            console.log(account);
            const height = await client.get.latestBlockHeight();
            const subaccount = new subaccount_1.SubaccountInfo(wallet, 0);
            const placeOrder = dummyOrder(height);
            placeOrder.clientId = (0, utils_1.randomInt)(1000000000);
            const tx = await client.post.placeOrderObject(subaccount, placeOrder);
            console.log('**Order Tx**');
            console.log(tx);
        });
    });
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiVmFsaWRhdG9yUG9zdEVuZHBvaW50cy50ZXN0LmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vX190ZXN0c19fL21vZHVsZXMvY2xpZW50L1ZhbGlkYXRvclBvc3RFbmRwb2ludHMudGVzdC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7OztBQUFBLHNGQUF3RjtBQUV4Riw4REFBeUQ7QUFDekQsNkZBQW9FO0FBQ3BFLGdFQUFpRTtBQUVqRSw0RUFBd0U7QUFDeEUsa0RBQW1EO0FBQ25ELDJEQUErRTtBQUMvRSxzQ0FBNkM7QUFFN0MsU0FBUyxVQUFVLENBQUMsTUFBYztJQUNoQyxNQUFNLFVBQVUsR0FBRyx3QkFBWSxDQUFDO0lBQ2hDLFVBQVUsQ0FBQyxRQUFRLEdBQUcsSUFBQSxpQkFBUyxFQUFDLFVBQVUsQ0FBQyxDQUFDO0lBQzVDLFVBQVUsQ0FBQyxZQUFZLEdBQUcsTUFBTSxHQUFHLENBQUMsQ0FBQztJQUNyQyw0Q0FBNEM7SUFDNUMsTUFBTSxNQUFNLEdBQUcsSUFBQSxpQkFBUyxFQUFDLElBQUksQ0FBQyxDQUFDO0lBQy9CLElBQUksQ0FBQyxNQUFNLEdBQUcsQ0FBQyxDQUFDLEtBQUssQ0FBQyxFQUFFO1FBQ3RCLFVBQVUsQ0FBQyxJQUFJLEdBQUcsa0JBQVUsQ0FBQyxRQUFRLENBQUM7S0FDdkM7U0FBTTtRQUNMLFVBQVUsQ0FBQyxJQUFJLEdBQUcsa0JBQVUsQ0FBQyxTQUFTLENBQUM7S0FDeEM7SUFDRCxPQUFPLFVBQVUsQ0FBQztBQUNwQixDQUFDO0FBRUQsUUFBUSxDQUFDLGtCQUFrQixFQUFFLEdBQUcsRUFBRTtJQUNoQyxJQUFJLE1BQW1CLENBQUM7SUFDeEIsSUFBSSxNQUF1QixDQUFDO0lBRTVCLFFBQVEsQ0FBQyxNQUFNLEVBQUUsR0FBRyxFQUFFO1FBQ3BCLFVBQVUsQ0FBQyxLQUFLLElBQUksRUFBRTtZQUNwQixNQUFNLEdBQUcsTUFBTSxzQkFBVyxDQUFDLFlBQVksQ0FBQyw4QkFBa0IsRUFBRSxtQkFBYSxDQUFDLENBQUM7WUFDM0UsTUFBTSxHQUFHLE1BQU0sa0NBQWUsQ0FBQyxPQUFPLENBQUMsbUJBQU8sQ0FBQyxPQUFPLEVBQUUsQ0FBQyxlQUFlLENBQUMsQ0FBQztRQUM1RSxDQUFDLENBQUMsQ0FBQztRQUVILEVBQUUsQ0FBQyxZQUFZLEVBQUUsS0FBSyxJQUFJLEVBQUU7WUFDMUIsT0FBTyxDQUFDLEdBQUcsQ0FBQyxZQUFZLENBQUMsQ0FBQztZQUMxQixPQUFPLENBQUMsR0FBRyxDQUFDLE1BQU0sQ0FBQyxDQUFDO1lBQ3BCLE1BQU0sT0FBTyxHQUFHLE1BQU0sQ0FBQyxPQUFRLENBQUM7WUFDaEMsTUFBTSxPQUFPLEdBQUcsTUFBTSxNQUFNLENBQUMsR0FBRyxDQUFDLFVBQVUsQ0FBQyxPQUFPLENBQUMsQ0FBQztZQUNyRCxPQUFPLENBQUMsR0FBRyxDQUFDLGFBQWEsQ0FBQyxDQUFDO1lBQzNCLE9BQU8sQ0FBQyxHQUFHLENBQUMsT0FBTyxDQUFDLENBQUM7WUFDckIsTUFBTSxNQUFNLEdBQUcsTUFBTSxNQUFNLENBQUMsR0FBRyxDQUFDLGlCQUFpQixFQUFFLENBQUM7WUFDcEQsTUFBTSxVQUFVLEdBQUcsSUFBSSwyQkFBYyxDQUFDLE1BQU0sRUFBRSxDQUFDLENBQUMsQ0FBQztZQUNqRCxNQUFNLFVBQVUsR0FBRyxVQUFVLENBQUMsTUFBTSxDQUFDLENBQUM7WUFDdEMsVUFBVSxDQUFDLFFBQVEsR0FBRyxJQUFBLGlCQUFTLEVBQUMsVUFBYSxDQUFDLENBQUM7WUFDL0MsTUFBTSxFQUFFLEdBQUcsTUFBTSxNQUFNLENBQUMsSUFBSSxDQUFDLGdCQUFnQixDQUMzQyxVQUFVLEVBQ1YsVUFBVSxDQUNYLENBQUM7WUFDRixPQUFPLENBQUMsR0FBRyxDQUFDLGNBQWMsQ0FBQyxDQUFDO1lBQzVCLE9BQU8sQ0FBQyxHQUFHLENBQUMsRUFBRSxDQUFDLENBQUM7UUFDbEIsQ0FBQyxDQUFDLENBQUM7SUFDTCxDQUFDLENBQUMsQ0FBQztBQUNMLENBQUMsQ0FBQyxDQUFDIn0=