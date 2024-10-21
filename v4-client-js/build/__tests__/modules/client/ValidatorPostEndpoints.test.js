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
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiVmFsaWRhdG9yUG9zdEVuZHBvaW50cy50ZXN0LmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vX190ZXN0c19fL21vZHVsZXMvY2xpZW50L1ZhbGlkYXRvclBvc3RFbmRwb2ludHMudGVzdC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7OztBQUFBLHNGQUF3RjtBQUV4Riw4REFBeUQ7QUFDekQsNkZBQW9FO0FBQ3BFLGdFQUFpRTtBQUVqRSw0RUFBd0U7QUFDeEUsa0RBQW1EO0FBQ25ELDJEQUErRTtBQUMvRSxzQ0FBNkM7QUFFN0MsU0FBUyxVQUFVLENBQUMsTUFBYztJQUNoQyxNQUFNLFVBQVUsR0FBRyx3QkFBWSxDQUFDO0lBQ2hDLFVBQVUsQ0FBQyxRQUFRLEdBQUcsSUFBQSxpQkFBUyxFQUFDLFVBQVUsQ0FBQyxDQUFDO0lBQzVDLFVBQVUsQ0FBQyxZQUFZLEdBQUcsTUFBTSxHQUFHLENBQUMsQ0FBQztJQUNyQyw0Q0FBNEM7SUFDNUMsTUFBTSxNQUFNLEdBQUcsSUFBQSxpQkFBUyxFQUFDLElBQUksQ0FBQyxDQUFDO0lBQy9CLElBQUksQ0FBQyxNQUFNLEdBQUcsQ0FBQyxDQUFDLEtBQUssQ0FBQyxFQUFFLENBQUM7UUFDdkIsVUFBVSxDQUFDLElBQUksR0FBRyxrQkFBVSxDQUFDLFFBQVEsQ0FBQztJQUN4QyxDQUFDO1NBQU0sQ0FBQztRQUNOLFVBQVUsQ0FBQyxJQUFJLEdBQUcsa0JBQVUsQ0FBQyxTQUFTLENBQUM7SUFDekMsQ0FBQztJQUNELE9BQU8sVUFBVSxDQUFDO0FBQ3BCLENBQUM7QUFFRCxRQUFRLENBQUMsa0JBQWtCLEVBQUUsR0FBRyxFQUFFO0lBQ2hDLElBQUksTUFBbUIsQ0FBQztJQUN4QixJQUFJLE1BQXVCLENBQUM7SUFFNUIsUUFBUSxDQUFDLE1BQU0sRUFBRSxHQUFHLEVBQUU7UUFDcEIsVUFBVSxDQUFDLEtBQUssSUFBSSxFQUFFO1lBQ3BCLE1BQU0sR0FBRyxNQUFNLHNCQUFXLENBQUMsWUFBWSxDQUFDLDhCQUFrQixFQUFFLG1CQUFhLENBQUMsQ0FBQztZQUMzRSxNQUFNLEdBQUcsTUFBTSxrQ0FBZSxDQUFDLE9BQU8sQ0FBQyxtQkFBTyxDQUFDLE9BQU8sRUFBRSxDQUFDLGVBQWUsQ0FBQyxDQUFDO1FBQzVFLENBQUMsQ0FBQyxDQUFDO1FBRUgsRUFBRSxDQUFDLFlBQVksRUFBRSxLQUFLLElBQUksRUFBRTtZQUMxQixPQUFPLENBQUMsR0FBRyxDQUFDLFlBQVksQ0FBQyxDQUFDO1lBQzFCLE9BQU8sQ0FBQyxHQUFHLENBQUMsTUFBTSxDQUFDLENBQUM7WUFDcEIsTUFBTSxPQUFPLEdBQUcsTUFBTSxDQUFDLE9BQVEsQ0FBQztZQUNoQyxNQUFNLE9BQU8sR0FBRyxNQUFNLE1BQU0sQ0FBQyxHQUFHLENBQUMsVUFBVSxDQUFDLE9BQU8sQ0FBQyxDQUFDO1lBQ3JELE9BQU8sQ0FBQyxHQUFHLENBQUMsYUFBYSxDQUFDLENBQUM7WUFDM0IsT0FBTyxDQUFDLEdBQUcsQ0FBQyxPQUFPLENBQUMsQ0FBQztZQUNyQixNQUFNLE1BQU0sR0FBRyxNQUFNLE1BQU0sQ0FBQyxHQUFHLENBQUMsaUJBQWlCLEVBQUUsQ0FBQztZQUNwRCxNQUFNLFVBQVUsR0FBRyxJQUFJLDJCQUFjLENBQUMsTUFBTSxFQUFFLENBQUMsQ0FBQyxDQUFDO1lBQ2pELE1BQU0sVUFBVSxHQUFHLFVBQVUsQ0FBQyxNQUFNLENBQUMsQ0FBQztZQUN0QyxVQUFVLENBQUMsUUFBUSxHQUFHLElBQUEsaUJBQVMsRUFBQyxVQUFhLENBQUMsQ0FBQztZQUMvQyxNQUFNLEVBQUUsR0FBRyxNQUFNLE1BQU0sQ0FBQyxJQUFJLENBQUMsZ0JBQWdCLENBQzNDLFVBQVUsRUFDVixVQUFVLENBQ1gsQ0FBQztZQUNGLE9BQU8sQ0FBQyxHQUFHLENBQUMsY0FBYyxDQUFDLENBQUM7WUFDNUIsT0FBTyxDQUFDLEdBQUcsQ0FBQyxFQUFFLENBQUMsQ0FBQztRQUNsQixDQUFDLENBQUMsQ0FBQztJQUNMLENBQUMsQ0FBQyxDQUFDO0FBQ0wsQ0FBQyxDQUFDLENBQUMifQ==