"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.NobleClient = void 0;
const proto_signing_1 = require("@cosmjs/proto-signing");
const stargate_1 = require("@cosmjs/stargate");
const constants_1 = require("./constants");
const cctpProto_1 = require("./lib/cctpProto");
class NobleClient {
    constructor(restEndpoint) {
        this.restEndpoint = restEndpoint;
    }
    get isConnected() {
        return Boolean(this.stargateClient);
    }
    async connect(wallet) {
        if ((wallet === null || wallet === void 0 ? void 0 : wallet.offlineSigner) === undefined) {
            throw new Error('Wallet signer not found');
        }
        this.wallet = wallet;
        this.stargateClient = await stargate_1.SigningStargateClient.connectWithSigner(this.restEndpoint, wallet.offlineSigner, {
            registry: new proto_signing_1.Registry([
                ['/circle.cctp.v1.MsgDepositForBurn', cctpProto_1.MsgDepositForBurn],
                ...stargate_1.defaultRegistryTypes,
            ]),
        });
    }
    getAccountBalances() {
        var _a;
        if (!this.stargateClient || ((_a = this.wallet) === null || _a === void 0 ? void 0 : _a.address) === undefined) {
            throw new Error('stargateClient not initialized');
        }
        return this.stargateClient.getAllBalances(this.wallet.address);
    }
    getAccountBalance(denom) {
        var _a;
        if (!this.stargateClient || ((_a = this.wallet) === null || _a === void 0 ? void 0 : _a.address) === undefined) {
            throw new Error('stargateClient not initialized');
        }
        return this.stargateClient.getBalance(this.wallet.address, denom);
    }
    async send(messages, gasPrice = stargate_1.GasPrice.fromString('0.025uusdc'), memo) {
        var _a;
        if (!this.stargateClient) {
            throw new Error('NobleClient stargateClient not initialized');
        }
        if (((_a = this.wallet) === null || _a === void 0 ? void 0 : _a.address) === undefined) {
            throw new Error('NobleClient wallet not initialized');
        }
        // Simulate to get the gas estimate
        const fee = await this.simulateTransaction(messages, gasPrice, memo);
        // Sign and broadcast the transaction
        return this.stargateClient.signAndBroadcast(this.wallet.address, messages, fee, memo !== null && memo !== void 0 ? memo : '');
    }
    async simulateTransaction(messages, gasPrice = stargate_1.GasPrice.fromString('0.025uusdc'), memo) {
        var _a, _b;
        if (!this.stargateClient) {
            throw new Error('NobleClient stargateClient not initialized');
        }
        if (((_a = this.wallet) === null || _a === void 0 ? void 0 : _a.address) === undefined) {
            throw new Error('NobleClient wallet not initialized');
        }
        // Get simulated response
        const gasEstimate = await this.stargateClient.simulate((_b = this.wallet) === null || _b === void 0 ? void 0 : _b.address, messages, memo);
        // Calculate and return the fee
        return (0, stargate_1.calculateFee)(Math.floor(gasEstimate * constants_1.GAS_MULTIPLIER), gasPrice);
    }
}
exports.NobleClient = NobleClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoibm9ibGUtY2xpZW50LmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vc3JjL2NsaWVudHMvbm9ibGUtY2xpZW50LnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7OztBQUFBLHlEQUFxRTtBQUNyRSwrQ0FPMEI7QUFFMUIsMkNBQTZDO0FBQzdDLCtDQUFvRDtBQUdwRCxNQUFhLFdBQVc7SUFLdEIsWUFBWSxZQUFvQjtRQUM5QixJQUFJLENBQUMsWUFBWSxHQUFHLFlBQVksQ0FBQztJQUNuQyxDQUFDO0lBRUQsSUFBSSxXQUFXO1FBQ2IsT0FBTyxPQUFPLENBQUMsSUFBSSxDQUFDLGNBQWMsQ0FBQyxDQUFDO0lBQ3RDLENBQUM7SUFFRCxLQUFLLENBQUMsT0FBTyxDQUFDLE1BQW1CO1FBQy9CLElBQUksQ0FBQSxNQUFNLGFBQU4sTUFBTSx1QkFBTixNQUFNLENBQUUsYUFBYSxNQUFLLFNBQVMsRUFBRTtZQUN2QyxNQUFNLElBQUksS0FBSyxDQUFDLHlCQUF5QixDQUFDLENBQUM7U0FDNUM7UUFDRCxJQUFJLENBQUMsTUFBTSxHQUFHLE1BQU0sQ0FBQztRQUNyQixJQUFJLENBQUMsY0FBYyxHQUFHLE1BQU0sZ0NBQXFCLENBQUMsaUJBQWlCLENBQ2pFLElBQUksQ0FBQyxZQUFZLEVBQ2pCLE1BQU0sQ0FBQyxhQUFhLEVBQ3BCO1lBQ0UsUUFBUSxFQUFFLElBQUksd0JBQVEsQ0FBQztnQkFDckIsQ0FBQyxtQ0FBbUMsRUFBRSw2QkFBaUIsQ0FBQztnQkFDeEQsR0FBRywrQkFBb0I7YUFDeEIsQ0FBQztTQUNILENBQ0YsQ0FBQztJQUNKLENBQUM7SUFFRCxrQkFBa0I7O1FBQ2hCLElBQUksQ0FBQyxJQUFJLENBQUMsY0FBYyxJQUFJLENBQUEsTUFBQSxJQUFJLENBQUMsTUFBTSwwQ0FBRSxPQUFPLE1BQUssU0FBUyxFQUFFO1lBQzlELE1BQU0sSUFBSSxLQUFLLENBQUMsZ0NBQWdDLENBQUMsQ0FBQztTQUNuRDtRQUNELE9BQU8sSUFBSSxDQUFDLGNBQWMsQ0FBQyxjQUFjLENBQUMsSUFBSSxDQUFDLE1BQU0sQ0FBQyxPQUFPLENBQUMsQ0FBQztJQUNqRSxDQUFDO0lBRUQsaUJBQWlCLENBQUMsS0FBYTs7UUFDN0IsSUFBSSxDQUFDLElBQUksQ0FBQyxjQUFjLElBQUksQ0FBQSxNQUFBLElBQUksQ0FBQyxNQUFNLDBDQUFFLE9BQU8sTUFBSyxTQUFTLEVBQUU7WUFDOUQsTUFBTSxJQUFJLEtBQUssQ0FBQyxnQ0FBZ0MsQ0FBQyxDQUFDO1NBQ25EO1FBQ0QsT0FBTyxJQUFJLENBQUMsY0FBYyxDQUFDLFVBQVUsQ0FBQyxJQUFJLENBQUMsTUFBTSxDQUFDLE9BQU8sRUFBRSxLQUFLLENBQUMsQ0FBQztJQUNwRSxDQUFDO0lBRUQsS0FBSyxDQUFDLElBQUksQ0FDUixRQUF3QixFQUN4QixXQUFxQixtQkFBUSxDQUFDLFVBQVUsQ0FBQyxZQUFZLENBQUMsRUFDdEQsSUFBYTs7UUFFYixJQUFJLENBQUMsSUFBSSxDQUFDLGNBQWMsRUFBRTtZQUN4QixNQUFNLElBQUksS0FBSyxDQUFDLDRDQUE0QyxDQUFDLENBQUM7U0FDL0Q7UUFDRCxJQUFJLENBQUEsTUFBQSxJQUFJLENBQUMsTUFBTSwwQ0FBRSxPQUFPLE1BQUssU0FBUyxFQUFFO1lBQ3RDLE1BQU0sSUFBSSxLQUFLLENBQUMsb0NBQW9DLENBQUMsQ0FBQztTQUN2RDtRQUNELG1DQUFtQztRQUNuQyxNQUFNLEdBQUcsR0FBRyxNQUFNLElBQUksQ0FBQyxtQkFBbUIsQ0FBQyxRQUFRLEVBQUUsUUFBUSxFQUFFLElBQUksQ0FBQyxDQUFDO1FBRXJFLHFDQUFxQztRQUNyQyxPQUFPLElBQUksQ0FBQyxjQUFjLENBQUMsZ0JBQWdCLENBQ3pDLElBQUksQ0FBQyxNQUFNLENBQUMsT0FBTyxFQUNuQixRQUFRLEVBQ1IsR0FBRyxFQUNILElBQUksYUFBSixJQUFJLGNBQUosSUFBSSxHQUFJLEVBQUUsQ0FDWCxDQUFDO0lBQ0osQ0FBQztJQUVELEtBQUssQ0FBQyxtQkFBbUIsQ0FDdkIsUUFBaUMsRUFDakMsV0FBcUIsbUJBQVEsQ0FBQyxVQUFVLENBQUMsWUFBWSxDQUFDLEVBQ3RELElBQWE7O1FBRWIsSUFBSSxDQUFDLElBQUksQ0FBQyxjQUFjLEVBQUU7WUFDeEIsTUFBTSxJQUFJLEtBQUssQ0FBQyw0Q0FBNEMsQ0FBQyxDQUFDO1NBQy9EO1FBQ0QsSUFBSSxDQUFBLE1BQUEsSUFBSSxDQUFDLE1BQU0sMENBQUUsT0FBTyxNQUFLLFNBQVMsRUFBRTtZQUN0QyxNQUFNLElBQUksS0FBSyxDQUFDLG9DQUFvQyxDQUFDLENBQUM7U0FDdkQ7UUFDRCx5QkFBeUI7UUFDekIsTUFBTSxXQUFXLEdBQUcsTUFBTSxJQUFJLENBQUMsY0FBYyxDQUFDLFFBQVEsQ0FDcEQsTUFBQSxJQUFJLENBQUMsTUFBTSwwQ0FBRSxPQUFPLEVBQ3BCLFFBQVEsRUFDUixJQUFJLENBQ0wsQ0FBQztRQUVGLCtCQUErQjtRQUMvQixPQUFPLElBQUEsdUJBQVksRUFBQyxJQUFJLENBQUMsS0FBSyxDQUFDLFdBQVcsR0FBRywwQkFBYyxDQUFDLEVBQUUsUUFBUSxDQUFDLENBQUM7SUFDMUUsQ0FBQztDQUNGO0FBeEZELGtDQXdGQyJ9