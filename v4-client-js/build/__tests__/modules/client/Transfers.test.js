"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const constants_1 = require("../../../src/clients/constants");
const local_wallet_1 = __importDefault(require("../../../src/clients/modules/local-wallet"));
const subaccount_1 = require("../../../src/clients/subaccount");
const validator_client_1 = require("../../../src/clients/validator-client");
const constants_2 = require("../../../examples/constants");
const long_1 = __importDefault(require("long"));
const src_1 = require("../../../src");
async function sleep(ms) {
    return new Promise((resolve) => setTimeout(resolve, ms));
}
describe('Validator Client', () => {
    let wallet;
    let subaccount;
    let client;
    describe('Transfers', () => {
        beforeEach(async () => {
            wallet = await local_wallet_1.default.fromMnemonic(constants_2.DYDX_TEST_MNEMONIC, src_1.BECH32_PREFIX);
            subaccount = new subaccount_1.SubaccountInfo(wallet, 0);
            client = await validator_client_1.ValidatorClient.connect(constants_1.Network.testnet().validatorConfig);
            await sleep(5000); // wait for withdraw to complete
        });
        it('Withdraw', async () => {
            const tx = await client.post.withdraw(subaccount, 0, new long_1.default(100000000), undefined);
            console.log('**Withdraw Tx**');
            console.log(tx);
        });
        it('Deposit', async () => {
            const tx = await client.post.deposit(subaccount, 0, new long_1.default(1000000));
            console.log('**Deposit Tx**');
            console.log(tx);
        });
        it('Transfer', async () => {
            const tx = await client.post.transfer(subaccount, subaccount.address, 1, 0, new long_1.default(1000));
            console.log('**Transfer Tx**');
            console.log(tx);
        });
    });
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiVHJhbnNmZXJzLnRlc3QuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi8uLi9fX3Rlc3RzX18vbW9kdWxlcy9jbGllbnQvVHJhbnNmZXJzLnRlc3QudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7QUFBQSw4REFBeUQ7QUFDekQsNkZBQW9FO0FBQ3BFLGdFQUFpRTtBQUNqRSw0RUFBd0U7QUFDeEUsMkRBQWlFO0FBQ2pFLGdEQUF3QjtBQUN4QixzQ0FBNkM7QUFFN0MsS0FBSyxVQUFVLEtBQUssQ0FBQyxFQUFVO0lBQzdCLE9BQU8sSUFBSSxPQUFPLENBQUMsQ0FBQyxPQUFPLEVBQUUsRUFBRSxDQUFDLFVBQVUsQ0FBQyxPQUFPLEVBQUUsRUFBRSxDQUFDLENBQUMsQ0FBQztBQUMzRCxDQUFDO0FBRUQsUUFBUSxDQUFDLGtCQUFrQixFQUFFLEdBQUcsRUFBRTtJQUNoQyxJQUFJLE1BQW1CLENBQUM7SUFDeEIsSUFBSSxVQUEwQixDQUFDO0lBQy9CLElBQUksTUFBdUIsQ0FBQztJQUU1QixRQUFRLENBQUMsV0FBVyxFQUFFLEdBQUcsRUFBRTtRQUN6QixVQUFVLENBQUMsS0FBSyxJQUFJLEVBQUU7WUFDcEIsTUFBTSxHQUFHLE1BQU0sc0JBQVcsQ0FBQyxZQUFZLENBQUMsOEJBQWtCLEVBQUUsbUJBQWEsQ0FBQyxDQUFDO1lBQzNFLFVBQVUsR0FBRyxJQUFJLDJCQUFjLENBQUMsTUFBTSxFQUFFLENBQUMsQ0FBQyxDQUFDO1lBQzNDLE1BQU0sR0FBRyxNQUFNLGtDQUFlLENBQUMsT0FBTyxDQUFDLG1CQUFPLENBQUMsT0FBTyxFQUFFLENBQUMsZUFBZSxDQUFDLENBQUM7WUFDMUUsTUFBTSxLQUFLLENBQUMsSUFBSSxDQUFDLENBQUMsQ0FBRSxnQ0FBZ0M7UUFDdEQsQ0FBQyxDQUFDLENBQUM7UUFFSCxFQUFFLENBQUMsVUFBVSxFQUFFLEtBQUssSUFBSSxFQUFFO1lBQ3hCLE1BQU0sRUFBRSxHQUFHLE1BQU0sTUFBTSxDQUFDLElBQUksQ0FBQyxRQUFRLENBQ25DLFVBQVUsRUFDVixDQUFDLEVBQ0QsSUFBSSxjQUFJLENBQUMsU0FBWSxDQUFDLEVBQ3RCLFNBQVMsQ0FDVixDQUFDO1lBQ0YsT0FBTyxDQUFDLEdBQUcsQ0FBQyxpQkFBaUIsQ0FBQyxDQUFDO1lBQy9CLE9BQU8sQ0FBQyxHQUFHLENBQUMsRUFBRSxDQUFDLENBQUM7UUFDbEIsQ0FBQyxDQUFDLENBQUM7UUFFSCxFQUFFLENBQUMsU0FBUyxFQUFFLEtBQUssSUFBSSxFQUFFO1lBQ3ZCLE1BQU0sRUFBRSxHQUFHLE1BQU0sTUFBTSxDQUFDLElBQUksQ0FBQyxPQUFPLENBQ2xDLFVBQVUsRUFDVixDQUFDLEVBQ0QsSUFBSSxjQUFJLENBQUMsT0FBUyxDQUFDLENBQ3BCLENBQUM7WUFDRixPQUFPLENBQUMsR0FBRyxDQUFDLGdCQUFnQixDQUFDLENBQUM7WUFDOUIsT0FBTyxDQUFDLEdBQUcsQ0FBQyxFQUFFLENBQUMsQ0FBQztRQUNsQixDQUFDLENBQUMsQ0FBQztRQUVILEVBQUUsQ0FBQyxVQUFVLEVBQUUsS0FBSyxJQUFJLEVBQUU7WUFDeEIsTUFBTSxFQUFFLEdBQUcsTUFBTSxNQUFNLENBQUMsSUFBSSxDQUFDLFFBQVEsQ0FDbkMsVUFBVSxFQUNWLFVBQVUsQ0FBQyxPQUFPLEVBQ2xCLENBQUMsRUFDRCxDQUFDLEVBQ0QsSUFBSSxjQUFJLENBQUMsSUFBSyxDQUFDLENBQ2hCLENBQUM7WUFDRixPQUFPLENBQUMsR0FBRyxDQUFDLGlCQUFpQixDQUFDLENBQUM7WUFDL0IsT0FBTyxDQUFDLEdBQUcsQ0FBQyxFQUFFLENBQUMsQ0FBQztRQUNsQixDQUFDLENBQUMsQ0FBQztJQUNMLENBQUMsQ0FBQyxDQUFDO0FBQ0wsQ0FBQyxDQUFDLENBQUMifQ==