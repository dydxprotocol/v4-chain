"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const long_1 = __importDefault(require("long"));
const constants_1 = require("../src/clients/constants");
const local_wallet_1 = __importDefault(require("../src/clients/modules/local-wallet"));
const noble_client_1 = require("../src/clients/noble-client");
const validator_client_1 = require("../src/clients/validator-client");
const constants_2 = require("../src/lib/constants");
const utils_1 = require("../src/lib/utils");
const constants_3 = require("./constants");
async function test() {
    const dydxClient = await validator_client_1.ValidatorClient.connect(constants_1.Network.testnet().validatorConfig);
    const dydxWallet = await local_wallet_1.default.fromMnemonic(constants_3.DYDX_TEST_MNEMONIC, constants_2.BECH32_PREFIX);
    const nobleWallet = await local_wallet_1.default.fromMnemonic(constants_3.DYDX_TEST_MNEMONIC, constants_2.NOBLE_BECH32_PREFIX);
    const client = new noble_client_1.NobleClient('https://rpc.testnet.noble.strange.love');
    await client.connect(nobleWallet);
    if (nobleWallet.address === undefined || dydxWallet.address === undefined) {
        throw new Error('Wallet not found');
    }
    // IBC to noble
    const ibcToNobleMsg = {
        typeUrl: '/ibc.applications.transfer.v1.MsgTransfer',
        value: {
            sourcePort: 'transfer',
            sourceChannel: 'channel-0',
            token: {
                denom: 'ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5',
                amount: '1000000',
            },
            sender: dydxWallet.address,
            receiver: nobleWallet.address,
            timeoutTimestamp: long_1.default.fromNumber(Math.floor(Date.now() / 1000) * 1e9 + 10 * 60 * 1e9),
        },
    };
    const msgs = [ibcToNobleMsg];
    const encodeObjects = new Promise((resolve) => resolve(msgs));
    await dydxClient.post.send(dydxWallet, () => {
        return encodeObjects;
    }, false, undefined, undefined);
    await (0, utils_1.sleep)(30000);
    try {
        const coins = await client.getAccountBalances();
        console.log('Balances');
        console.log(JSON.stringify(coins));
        // IBC from noble
        const ibcFromNobleMsg = {
            typeUrl: '/ibc.applications.transfer.v1.MsgTransfer',
            value: {
                sourcePort: 'transfer',
                sourceChannel: 'channel-21',
                token: {
                    denom: 'uusdc',
                    amount: coins[0].amount,
                },
                sender: nobleWallet.address,
                receiver: dydxWallet.address,
                timeoutTimestamp: long_1.default.fromNumber(Math.floor(Date.now() / 1000) * 1e9 + 10 * 60 * 1e9),
            },
        };
        const fee = await client.simulateTransaction([ibcFromNobleMsg]);
        ibcFromNobleMsg.value.token.amount = (parseInt(ibcFromNobleMsg.value.token.amount, 10) -
            Math.floor(parseInt(fee.amount[0].amount, 10) * 1.4)).toString();
        await client.send([ibcFromNobleMsg]);
    }
    catch (error) {
        console.log(JSON.stringify(error.message));
    }
    await (0, utils_1.sleep)(30000);
    try {
        const coin = await client.getAccountBalance('uusdc');
        console.log('Balance');
        console.log(JSON.stringify(coin));
    }
    catch (error) {
        console.log(JSON.stringify(error.message));
    }
}
test()
    .then(() => { })
    .catch((error) => {
    console.log(error.message);
    console.log(error);
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoibm9ibGVfZXhhbXBsZS5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uL2V4YW1wbGVzL25vYmxlX2V4YW1wbGUudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7QUFDQSxnREFBd0I7QUFFeEIsd0RBQW1EO0FBQ25ELHVGQUE4RDtBQUM5RCw4REFBMEQ7QUFDMUQsc0VBQWtFO0FBQ2xFLG9EQUEwRTtBQUMxRSw0Q0FBeUM7QUFDekMsMkNBQWlEO0FBRWpELEtBQUssVUFBVSxJQUFJO0lBQ2pCLE1BQU0sVUFBVSxHQUFHLE1BQU0sa0NBQWUsQ0FBQyxPQUFPLENBQzlDLG1CQUFPLENBQUMsT0FBTyxFQUFFLENBQUMsZUFBZSxDQUNsQyxDQUFDO0lBRUYsTUFBTSxVQUFVLEdBQUcsTUFBTSxzQkFBVyxDQUFDLFlBQVksQ0FDL0MsOEJBQWtCLEVBQ2xCLHlCQUFhLENBQ2QsQ0FBQztJQUNGLE1BQU0sV0FBVyxHQUFHLE1BQU0sc0JBQVcsQ0FBQyxZQUFZLENBQ2hELDhCQUFrQixFQUNsQiwrQkFBbUIsQ0FDcEIsQ0FBQztJQUVGLE1BQU0sTUFBTSxHQUFHLElBQUksMEJBQVcsQ0FBQyx3Q0FBd0MsQ0FBQyxDQUFDO0lBQ3pFLE1BQU0sTUFBTSxDQUFDLE9BQU8sQ0FBQyxXQUFXLENBQUMsQ0FBQztJQUVsQyxJQUFJLFdBQVcsQ0FBQyxPQUFPLEtBQUssU0FBUyxJQUFJLFVBQVUsQ0FBQyxPQUFPLEtBQUssU0FBUyxFQUFFO1FBQ3pFLE1BQU0sSUFBSSxLQUFLLENBQUMsa0JBQWtCLENBQUMsQ0FBQztLQUNyQztJQUVELGVBQWU7SUFFZixNQUFNLGFBQWEsR0FBaUI7UUFDbEMsT0FBTyxFQUFFLDJDQUEyQztRQUNwRCxLQUFLLEVBQUU7WUFDTCxVQUFVLEVBQUUsVUFBVTtZQUN0QixhQUFhLEVBQUUsV0FBVztZQUMxQixLQUFLLEVBQUU7Z0JBQ0wsS0FBSyxFQUNILHNFQUFzRTtnQkFDeEUsTUFBTSxFQUFFLFNBQVM7YUFDbEI7WUFDRCxN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