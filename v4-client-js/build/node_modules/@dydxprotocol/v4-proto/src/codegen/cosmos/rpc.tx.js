"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.createRPCMsgClient = void 0;
const createRPCMsgClient = async ({ rpc }) => ({
    cosmos: {
        auth: {
            v1beta1: new (await Promise.resolve().then(() => __importStar(require("./auth/v1beta1/tx.rpc.msg")))).MsgClientImpl(rpc)
        },
        authz: {
            v1beta1: new (await Promise.resolve().then(() => __importStar(require("./authz/v1beta1/tx.rpc.msg")))).MsgClientImpl(rpc)
        },
        bank: {
            v1beta1: new (await Promise.resolve().then(() => __importStar(require("./bank/v1beta1/tx.rpc.msg")))).MsgClientImpl(rpc)
        },
        circuit: {
            v1: new (await Promise.resolve().then(() => __importStar(require("./circuit/v1/tx.rpc.msg")))).MsgClientImpl(rpc)
        },
        consensus: {
            v1: new (await Promise.resolve().then(() => __importStar(require("./consensus/v1/tx.rpc.msg")))).MsgClientImpl(rpc)
        },
        crisis: {
            v1beta1: new (await Promise.resolve().then(() => __importStar(require("./crisis/v1beta1/tx.rpc.msg")))).MsgClientImpl(rpc)
        },
        distribution: {
            v1beta1: new (await Promise.resolve().then(() => __importStar(require("./distribution/v1beta1/tx.rpc.msg")))).MsgClientImpl(rpc)
        },
        evidence: {
            v1beta1: new (await Promise.resolve().then(() => __importStar(require("./evidence/v1beta1/tx.rpc.msg")))).MsgClientImpl(rpc)
        },
        feegrant: {
            v1beta1: new (await Promise.resolve().then(() => __importStar(require("./feegrant/v1beta1/tx.rpc.msg")))).MsgClientImpl(rpc)
        },
        gov: {
            v1: new (await Promise.resolve().then(() => __importStar(require("./gov/v1/tx.rpc.msg")))).MsgClientImpl(rpc),
            v1beta1: new (await Promise.resolve().then(() => __importStar(require("./gov/v1beta1/tx.rpc.msg")))).MsgClientImpl(rpc)
        },
        group: {
            v1: new (await Promise.resolve().then(() => __importStar(require("./group/v1/tx.rpc.msg")))).MsgClientImpl(rpc)
        },
        mint: {
            v1beta1: new (await Promise.resolve().then(() => __importStar(require("./mint/v1beta1/tx.rpc.msg")))).MsgClientImpl(rpc)
        },
        nft: {
            v1beta1: new (await Promise.resolve().then(() => __importStar(require("./nft/v1beta1/tx.rpc.msg")))).MsgClientImpl(rpc)
        },
        slashing: {
            v1beta1: new (await Promise.resolve().then(() => __importStar(require("./slashing/v1beta1/tx.rpc.msg")))).MsgClientImpl(rpc)
        },
        staking: {
            v1beta1: new (await Promise.resolve().then(() => __importStar(require("./staking/v1beta1/tx.rpc.msg")))).MsgClientImpl(rpc)
        },
        upgrade: {
            v1beta1: new (await Promise.resolve().then(() => __importStar(require("./upgrade/v1beta1/tx.rpc.msg")))).MsgClientImpl(rpc)
        },
        vesting: {
            v1beta1: new (await Promise.resolve().then(() => __importStar(require("./vesting/v1beta1/tx.rpc.msg")))).MsgClientImpl(rpc)
        }
    }
});
exports.createRPCMsgClient = createRPCMsgClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicnBjLnR4LmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vY29zbW9zL3JwYy50eC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7OztBQUNPLE1BQU0sa0JBQWtCLEdBQUcsS0FBSyxFQUFFLEVBQ3ZDLEdBQUcsRUFHSixFQUFFLEVBQUUsQ0FBQyxDQUFDO0lBQ0wsTUFBTSxFQUFFO1FBQ04sSUFBSSxFQUFFO1lBQ0osT0FBTyxFQUFFLElBQUksQ0FBQyx3REFBYSwyQkFBMkIsR0FBQyxDQUFDLENBQUMsYUFBYSxDQUFDLEdBQUcsQ0FBQztTQUM1RTtRQUNELEtBQUssRUFBRTtZQUNMLE9BQU8sRUFBRSxJQUFJLENBQUMsd0RBQWEsNEJBQTRCLEdBQUMsQ0FBQyxDQUFDLGFBQWEsQ0FBQyxHQUFHLENBQUM7U0FDN0U7UUFDRCxJQUFJLEVBQUU7WUFDSixPQUFPLEVBQUUsSUFBSSxDQUFDLHdEQUFhLDJCQUEyQixHQUFDLENBQUMsQ0FBQyxhQUFhLENBQUMsR0FBRyxDQUFDO1NBQzVFO1FBQ0QsT0FBTyxFQUFFO1lBQ1AsRUFBRSxFQUFFLElBQUksQ0FBQyx3REFBYSx5QkFBeUIsR0FBQyxDQUFDLENBQUMsYUFBYSxDQUFDLEdBQUcsQ0FBQztTQUNyRTtRQUNELFNBQVMsRUFBRTtZQUNULEVBQUUsRUFBRSxJQUFJLENBQUMsd0RBQWEsMkJBQTJCLEdBQUMsQ0FBQyxDQUFDLGFBQWEsQ0FBQyxHQUFHLENBQUM7U0FDdkU7UUFDRCxNQUFNLEVBQUU7WUFDTixPQUFPLEVBQUUsSUFBSSxDQUFDLHdEQUFhLDZCQUE2QixHQUFDLENBQUMsQ0FBQyxhQUFhLENBQUMsR0FBRyxDQUFDO1NBQzlFO1FBQ0QsWUFBWSxFQUFFO1lBQ1osT0FBTyxFQUFFLElBQUksQ0FBQyx3REFBYSxtQ0FBbUMsR0FBQyxDQUFDLENBQUMsYUFBYSxDQUFDLEdBQUcsQ0FBQztTQUNwRjtRQUNELFFBQVEsRUFBRTtZQUNSLE9BQU8sRUFBRSxJQUFJLENBQUMsd0RBQWEsK0JBQStCLEdBQUMsQ0FBQyxDQUFDLGFBQWEsQ0FBQyxHQUFHLENBQUM7U0FDaEY7UUFDRCxRQUFRLEVBQUU7WUFDUixPQUFPLEVBQUUsSUFBSSxDQUFDLHdEQUFhLCtCQUErQixHQUFDLENBQUMsQ0FBQyxhQUFhLENBQUMsR0FBRyxDQUFDO1NBQ2hGO1FBQ0QsR0FBRyxFQUFFO1lBQ0gsRUFBRSxFQUFFLElBQUksQ0FBQyx3REFBYSxxQkFBcUIsR0FBQyxDQUFDLENBQUMsYUFBYSxDQUFDLEdBQUcsQ0FBQztZQUNoRSxPQUFPLEVBQUUsSUFBSSxDQUFDLHdEQUFhLDBCQUEwQixHQUFDLENBQUMsQ0FBQyxhQUFhLENBQUMsR0FBRyxDQUFDO1NBQzNFO1FBQ0QsS0FBSyxFQUFFO1lBQ0wsRUFBRSxFQUFFLElBQUksQ0FBQyx3REFBYSx1QkFBdUIsR0FBQyxDQUFDLENBQUMsYUFBYSxDQUFDLEdBQUcsQ0FBQztTQUNuRTtRQUNELElBQUksRUFBRTtZQUNKLE9BQU8sRUFBRSxJQUFJLENBQUMsd0RBQWEsMkJBQTJCLEdBQUMsQ0FBQyxDQUFDLGFBQWEsQ0FBQyxHQUFHLENBQUM7U0FDNUU7UUFDRCxHQUFHLEVBQUU7WUFDSCxPQUFPLEVBQUUsSUFBSSxDQUFDLHdEQUFhLDBCQUEwQixHQUFDLENBQUMsQ0FBQyxhQUFhLENBQUMsR0FBRyxDQUFDO1NBQzNFO1FBQ0QsUUFBUSxFQUFFO1lBQ1IsT0FBTyxFQUFFLElBQUksQ0FBQyx3REFBYSwrQkFBK0IsR0FBQyxDQUFDLENBQUMsYUFBYSxDQUFDLEdBQUcsQ0FBQztTQUNoRjtRQUNELE9BQU8sRUFBRTtZQUNQLE9BQU8sRUFBRSxJQUFJLENBQUMsd0RBQWEsOEJBQThCLEdBQUMsQ0FBQyxDQUFDLGFBQWEsQ0FBQyxHQUFHLENBQUM7U0FDL0U7UUFDRCxPQUFPLEVBQUU7WUFDUCxPQUFPLEVBQUUsSUFBSSxDQUFDLHdEQUFhLDhCQUE4QixHQUFDLENBQUMsQ0FBQyxhQUFhLENBQUMsR0FBRyxDQUFDO1NBQy9FO1FBQ0QsT0FBTyxFQUFFO1lBQ1AsT0FBTyxFQUFFLElBQUksQ0FBQyx3REFBYSw4QkFBOEIsR0FBQyxDQUFDLENBQUMsYUFBYSxDQUFDLEdBQUcsQ0FBQztTQUMvRTtLQUNGO0NBQ0YsQ0FBQyxDQUFDO0FBM0RVLFFBQUEsa0JBQWtCLHNCQTJENUIifQ==