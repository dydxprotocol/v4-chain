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
exports.createRPCQueryClient = void 0;
const tendermint_rpc_1 = require("@cosmjs/tendermint-rpc");
const stargate_1 = require("@cosmjs/stargate");
const createRPCQueryClient = async ({ rpcEndpoint }) => {
    const tmClient = await tendermint_rpc_1.Tendermint34Client.connect(rpcEndpoint);
    const client = new stargate_1.QueryClient(tmClient);
    return {
        cosmos: {
            app: {
                v1alpha1: (await Promise.resolve().then(() => __importStar(require("./app/v1alpha1/query.rpc.Query")))).createRpcQueryExtension(client)
            },
            auth: {
                v1beta1: (await Promise.resolve().then(() => __importStar(require("./auth/v1beta1/query.rpc.Query")))).createRpcQueryExtension(client)
            },
            authz: {
                v1beta1: (await Promise.resolve().then(() => __importStar(require("./authz/v1beta1/query.rpc.Query")))).createRpcQueryExtension(client)
            },
            bank: {
                v1beta1: (await Promise.resolve().then(() => __importStar(require("./bank/v1beta1/query.rpc.Query")))).createRpcQueryExtension(client)
            },
            base: {
                node: {
                    v1beta1: (await Promise.resolve().then(() => __importStar(require("./base/node/v1beta1/query.rpc.Service")))).createRpcQueryExtension(client)
                },
                tendermint: {
                    v1beta1: (await Promise.resolve().then(() => __importStar(require("./base/tendermint/v1beta1/query.rpc.Service")))).createRpcQueryExtension(client)
                }
            },
            circuit: {
                v1: (await Promise.resolve().then(() => __importStar(require("./circuit/v1/query.rpc.Query")))).createRpcQueryExtension(client)
            },
            consensus: {
                v1: (await Promise.resolve().then(() => __importStar(require("./consensus/v1/query.rpc.Query")))).createRpcQueryExtension(client)
            },
            distribution: {
                v1beta1: (await Promise.resolve().then(() => __importStar(require("./distribution/v1beta1/query.rpc.Query")))).createRpcQueryExtension(client)
            },
            evidence: {
                v1beta1: (await Promise.resolve().then(() => __importStar(require("./evidence/v1beta1/query.rpc.Query")))).createRpcQueryExtension(client)
            },
            feegrant: {
                v1beta1: (await Promise.resolve().then(() => __importStar(require("./feegrant/v1beta1/query.rpc.Query")))).createRpcQueryExtension(client)
            },
            gov: {
                v1: (await Promise.resolve().then(() => __importStar(require("./gov/v1/query.rpc.Query")))).createRpcQueryExtension(client),
                v1beta1: (await Promise.resolve().then(() => __importStar(require("./gov/v1beta1/query.rpc.Query")))).createRpcQueryExtension(client)
            },
            group: {
                v1: (await Promise.resolve().then(() => __importStar(require("./group/v1/query.rpc.Query")))).createRpcQueryExtension(client)
            },
            mint: {
                v1beta1: (await Promise.resolve().then(() => __importStar(require("./mint/v1beta1/query.rpc.Query")))).createRpcQueryExtension(client)
            },
            nft: {
                v1beta1: (await Promise.resolve().then(() => __importStar(require("./nft/v1beta1/query.rpc.Query")))).createRpcQueryExtension(client)
            },
            orm: {
                query: {
                    v1alpha1: (await Promise.resolve().then(() => __importStar(require("./orm/query/v1alpha1/query.rpc.Query")))).createRpcQueryExtension(client)
                }
            },
            params: {
                v1beta1: (await Promise.resolve().then(() => __importStar(require("./params/v1beta1/query.rpc.Query")))).createRpcQueryExtension(client)
            },
            slashing: {
                v1beta1: (await Promise.resolve().then(() => __importStar(require("./slashing/v1beta1/query.rpc.Query")))).createRpcQueryExtension(client)
            },
            staking: {
                v1beta1: (await Promise.resolve().then(() => __importStar(require("./staking/v1beta1/query.rpc.Query")))).createRpcQueryExtension(client)
            },
            tx: {
                v1beta1: (await Promise.resolve().then(() => __importStar(require("./tx/v1beta1/service.rpc.Service")))).createRpcQueryExtension(client)
            },
            upgrade: {
                v1beta1: (await Promise.resolve().then(() => __importStar(require("./upgrade/v1beta1/query.rpc.Query")))).createRpcQueryExtension(client)
            }
        }
    };
};
exports.createRPCQueryClient = createRPCQueryClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicnBjLnF1ZXJ5LmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vY29zbW9zL3JwYy5xdWVyeS50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7OztBQUFBLDJEQUEwRTtBQUMxRSwrQ0FBK0M7QUFDeEMsTUFBTSxvQkFBb0IsR0FBRyxLQUFLLEVBQUUsRUFDekMsV0FBVyxFQUdaLEVBQUUsRUFBRTtJQUNILE1BQU0sUUFBUSxHQUFHLE1BQU0sbUNBQWtCLENBQUMsT0FBTyxDQUFDLFdBQVcsQ0FBQyxDQUFDO0lBQy9ELE1BQU0sTUFBTSxHQUFHLElBQUksc0JBQVcsQ0FBQyxRQUFRLENBQUMsQ0FBQztJQUN6QyxPQUFPO1FBQ0wsTUFBTSxFQUFFO1lBQ04sR0FBRyxFQUFFO2dCQUNILFFBQVEsRUFBRSxDQUFDLHdEQUFhLGdDQUFnQyxHQUFDLENBQUMsQ0FBQyx1QkFBdUIsQ0FBQyxNQUFNLENBQUM7YUFDM0Y7WUFDRCxJQUFJLEVBQUU7Z0JBQ0osT0FBTyxFQUFFLENBQUMsd0RBQWEsZ0NBQWdDLEdBQUMsQ0FBQyxDQUFDLHVCQUF1QixDQUFDLE1BQU0sQ0FBQzthQUMxRjtZQUNELEtBQUssRUFBRTtnQkFDTCxPQUFPLEVBQUUsQ0FBQyx3REFBYSxpQ0FBaUMsR0FBQyxDQUFDLENBQUMsdUJBQXVCLENBQUMsTUFBTSxDQUFDO2FBQzNGO1lBQ0QsSUFBSSxFQUFFO2dCQUNKLE9BQU8sRUFBRSxDQUFDLHdEQUFhLGdDQUFnQyxHQUFDLENBQUMsQ0FBQyx1QkFBdUIsQ0FBQyxNQUFNLENBQUM7YUFDMUY7WUFDRCxJQUFJLEVBQUU7Z0JBQ0osSUFBSSxFQUFFO29CQUNKLE9BQU8sRUFBRSxDQUFDLHdEQUFhLHVDQUF1QyxHQUFDLENBQUMsQ0FBQyx1QkFBdUIsQ0FBQyxNQUFNLENBQUM7aUJBQ2pHO2dCQUNELFVBQVUsRUFBRTtvQkFDVixPQUFPLEVBQUUsQ0FBQyx3REFBYSw2Q0FBNkMsR0FBQyxDQUFDLENBQUMsdUJBQXVCLENBQUMsTUFBTSxDQUFDO2lCQUN2RzthQUNGO1lBQ0QsT0FBTyxFQUFFO2dCQUNQLEVBQUUsRUFBRSxDQUFDLHdEQUFhLDhCQUE4QixHQUFDLENBQUMsQ0FBQyx1QkFBdUIsQ0FBQyxNQUFNLENBQUM7YUFDbkY7WUFDRCxTQUFTLEVBQUU7Z0JBQ1QsRUFBRSxFQUFFLENBQUMsd0RBQWEsZ0NBQWdDLEdBQUMsQ0FBQyxDQUFDLHVCQUF1QixDQUFDLE1BQU0sQ0FBQzthQUNyRjtZQUNELFlBQVksRUFBRTtnQkFDWixPQUFPLEVBQUUsQ0FBQyx3REFBYSx3Q0FBd0MsR0FBQyxDQUFDLENBQUMsdUJBQXVCLENBQUMsTUFBTSxDQUFDO2FBQ2xHO1lBQ0QsUUFBUSxFQUFFO2dCQUNSLE9BQU8sRUFBRSxDQUFDLHdEQUFhLG9DQUFvQyxHQUFDLENBQUMsQ0FBQyx1QkFBdUIsQ0FBQyxNQUFNLENBQUM7YUFDOUY7WUFDRCxRQUFRLEVBQUU7Z0JBQ1IsT0FBTyxFQUFFLENBQUMsd0RBQWEsb0NBQW9DLEdBQUMsQ0FBQyxDQUFDLHVCQUF1QixDQUFDLE1BQU0sQ0FBQzthQUM5RjtZQUNELEdBQUcsRUFBRTtnQkFDSCxFQUFFLEVBQUUsQ0FBQyx3REFBYSwwQkFBMEIsR0FBQyxDQUFDLENBQUMsdUJBQXVCLENBQUMsTUFBTSxDQUFDO2dCQUM5RSxPQUFPLEVBQUUsQ0FBQyx3REFBYSwrQkFBK0IsR0FBQyxDQUFDLENBQUMsdUJBQXVCLENBQUMsTUFBTSxDQUFDO2FBQ3pGO1lBQ0QsS0FBSyxFQUFFO2dCQUNMLEVBQUUsRUFBRSxDQUFDLHdEQUFhLDRCQUE0QixHQUFDLENBQUMsQ0FBQyx1QkFBdUIsQ0FBQyxNQUFNLENBQUM7YUFDakY7WUFDRCxJQUFJLEVBQUU7Z0JBQ0osT0FBTyxFQUFFLENBQUMsd0RBQWEsZ0NBQWdDLEdBQUMsQ0FBQyxDQUFDLHVCQUF1QixDQUFDLE1BQU0sQ0FBQzthQUMxRjtZQUNELEdBQUcsRUFBRTtnQkFDSCxPQUFPLEVBQUUsQ0FBQyx3REFBYSwrQkFBK0IsR0FBQyxDQUFDLENBQUMsdUJBQXVCLENBQUMsTUFBTSxDQUFDO2FBQ3pGO1lBQ0QsR0FBRyxFQUFFO2dCQUNILEtBQUssRUFBRTtvQkFDTCxRQUFRLEVBQUUsQ0FBQyx3REFBYSxzQ0FBc0MsR0FBQyxDQUFDLENBQUMsdUJBQXVCLENBQUMsTUFBTSxDQUFDO2lCQUNqRzthQUNGO1lBQ0QsTUFBTSxFQUFFO2dCQUNOLE9BQU8sRUFBRSxDQUFDLHdEQUFhLGtDQUFrQyxHQUFDLENBQUMsQ0FBQyx1QkFBdUIsQ0FBQyxNQUFNLENBQUM7YUFDNUY7WUFDRCxRQUFRLEVBQUU7Z0JBQ1IsT0FBTyxFQUFFLENBQUMsd0RBQWEsb0NBQW9DLEdBQUMsQ0FBQyxDQUFDLHVCQUF1QixDQUFDLE1BQU0sQ0FBQzthQUM5RjtZQUNELE9BQU8sRUFBRTtnQkFDUCxPQUFPLEVBQUUsQ0FBQyx3REFBYSxtQ0FBbUMsR0FBQyxDQUFDLENBQUMsdUJBQXVCLENBQUMsTUFBTSxDQUFDO2FBQzdGO1lBQ0QsRUFBRSxFQUFFO2dCQUNGLE9BQU8sRUFBRSxDQUFDLHdEQUFhLGtDQUFrQyxHQUFDLENBQUMsQ0FBQyx1QkFBdUIsQ0FBQyxNQUFNLENBQUM7YUFDNUY7WUFDRCxPQUFPLEVBQUU7Z0JBQ1AsT0FBTyxFQUFFLENBQUMsd0RBQWEsbUNBQW1DLEdBQUMsQ0FBQyxDQUFDLHVCQUF1QixDQUFDLE1BQU0sQ0FBQzthQUM3RjtTQUNGO0tBQ0YsQ0FBQztBQUNKLENBQUMsQ0FBQztBQS9FVyxRQUFBLG9CQUFvQix3QkErRS9CIn0=