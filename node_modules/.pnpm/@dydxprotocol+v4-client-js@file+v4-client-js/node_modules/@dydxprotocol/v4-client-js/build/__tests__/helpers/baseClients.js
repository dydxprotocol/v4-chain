"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BaseWallet = exports.BaseStargateSigningClient = exports.BaseQueryClient = exports.BaseTendermintClient = void 0;
/* eslint-disable @typescript-eslint/no-empty-function */
class BaseTendermintClient {
    async block() { }
    async broadcastTxSync() { }
    async broadcastTxAsync() { }
    async txSearchAll() { }
}
exports.BaseTendermintClient = BaseTendermintClient;
class BaseQueryClient {
    constructor() {
        this.tx = {
            async simulate() { },
        };
    }
    async queryUnverified() { }
}
exports.BaseQueryClient = BaseQueryClient;
class BaseStargateSigningClient {
    async sign() { }
}
exports.BaseStargateSigningClient = BaseStargateSigningClient;
class BaseWallet {
    async getAccounts() { }
}
exports.BaseWallet = BaseWallet;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiYmFzZUNsaWVudHMuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi9fX3Rlc3RzX18vaGVscGVycy9iYXNlQ2xpZW50cy50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7QUFBQSx5REFBeUQ7QUFDekQsTUFBYSxvQkFBb0I7SUFDL0IsS0FBSyxDQUFDLEtBQUssS0FBbUIsQ0FBQztJQUMvQixLQUFLLENBQUMsZUFBZSxLQUFtQixDQUFDO0lBQ3pDLEtBQUssQ0FBQyxnQkFBZ0IsS0FBbUIsQ0FBQztJQUMxQyxLQUFLLENBQUMsV0FBVyxLQUFtQixDQUFDO0NBQ3RDO0FBTEQsb0RBS0M7QUFFRCxNQUFhLGVBQWU7SUFBNUI7UUFDRSxPQUFFLEdBQUc7WUFDSCxLQUFLLENBQUMsUUFBUSxLQUFtQixDQUFDO1NBQ25DLENBQUM7SUFHSixDQUFDO0lBREMsS0FBSyxDQUFDLGVBQWUsS0FBbUIsQ0FBQztDQUMxQztBQU5ELDBDQU1DO0FBRUQsTUFBYSx5QkFBeUI7SUFDcEMsS0FBSyxDQUFDLElBQUksS0FBbUIsQ0FBQztDQUMvQjtBQUZELDhEQUVDO0FBRUQsTUFBYSxVQUFVO0lBQ3JCLEtBQUssQ0FBQyxXQUFXLEtBQW1CLENBQUM7Q0FDdEM7QUFGRCxnQ0FFQyJ9