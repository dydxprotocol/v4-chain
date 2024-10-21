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
    // Use sDAI denom for ibc transfers
    const ibcToNobleMsg = {
        typeUrl: '/ibc.applications.transfer.v1.MsgTransfer',
        value: {
            sourcePort: 'transfer',
            sourceChannel: 'channel-0',
            token: {
                denom: 'ibc/DEEFE2DEFDC8EA8879923C4CCA42BB888C3CD03FF7ECFEFB1C2FEC27A732ACC8',
                amount: '1000000000000000000',
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
                    denom: 'utdai',
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
        const coin = await client.getAccountBalance('utdai');
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
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoibm9ibGVfZXhhbXBsZS5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uL2V4YW1wbGVzL25vYmxlX2V4YW1wbGUudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7QUFDQSxnREFBd0I7QUFFeEIsd0RBQW1EO0FBQ25ELHVGQUE4RDtBQUM5RCw4REFBMEQ7QUFDMUQsc0VBQWtFO0FBQ2xFLG9EQUEwRTtBQUMxRSw0Q0FBeUM7QUFDekMsMkNBQWlEO0FBRWpELEtBQUssVUFBVSxJQUFJO0lBQ2pCLE1BQU0sVUFBVSxHQUFHLE1BQU0sa0NBQWUsQ0FBQyxPQUFPLENBQzlDLG1CQUFPLENBQUMsT0FBTyxFQUFFLENBQUMsZUFBZSxDQUNsQyxDQUFDO0lBRUYsTUFBTSxVQUFVLEdBQUcsTUFBTSxzQkFBVyxDQUFDLFlBQVksQ0FDL0MsOEJBQWtCLEVBQ2xCLHlCQUFhLENBQ2QsQ0FBQztJQUNGLE1BQU0sV0FBVyxHQUFHLE1BQU0sc0JBQVcsQ0FBQyxZQUFZLENBQ2hELDhCQUFrQixFQUNsQiwrQkFBbUIsQ0FDcEIsQ0FBQztJQUVGLE1BQU0sTUFBTSxHQUFHLElBQUksMEJBQVcsQ0FBQyx3Q0FBd0MsQ0FBQyxDQUFDO0lBQ3pFLE1BQU0sTUFBTSxDQUFDLE9BQU8sQ0FBQyxXQUFXLENBQUMsQ0FBQztJQUVsQyxJQUFJLFdBQVcsQ0FBQyxPQUFPLEtBQUssU0FBUyxJQUFJLFVBQVUsQ0FBQyxPQUFPLEtBQUssU0FBUyxFQUFFLENBQUM7UUFDMUUsTUFBTSxJQUFJLEtBQUssQ0FBQyxrQkFBa0IsQ0FBQyxDQUFDO0lBQ3RDLENBQUM7SUFFRCxlQUFlO0lBQ2YsbUNBQW1DO0lBQ25DLE1BQU0sYUFBYSxHQUFpQjtRQUNsQyxPQUFPLEVBQUUsMkNBQTJDO1FBQ3BELEtBQUssRUFBRTtZQUNMLFVBQVUsRUFBRSxVQUFVO1lBQ3RCLGFBQWEsRUFBRSxXQUFXO1lBQzFCLEtBQUssRUFBRTtnQkFDTCxLQUFLLEVBQ0gsc0VBQXNFO2dCQUN4RSxNQUFN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