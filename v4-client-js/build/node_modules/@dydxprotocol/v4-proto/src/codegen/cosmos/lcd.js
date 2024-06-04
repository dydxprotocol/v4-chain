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
exports.createLCDClient = void 0;
const lcd_1 = require("@osmonauts/lcd");
const createLCDClient = async ({ restEndpoint }) => {
    const requestClient = new lcd_1.LCDClient({
        restEndpoint
    });
    return {
        cosmos: {
            auth: {
                v1beta1: new (await Promise.resolve().then(() => __importStar(require("./auth/v1beta1/query.lcd")))).LCDQueryClient({
                    requestClient
                })
            },
            authz: {
                v1beta1: new (await Promise.resolve().then(() => __importStar(require("./authz/v1beta1/query.lcd")))).LCDQueryClient({
                    requestClient
                })
            },
            bank: {
                v1beta1: new (await Promise.resolve().then(() => __importStar(require("./bank/v1beta1/query.lcd")))).LCDQueryClient({
                    requestClient
                })
            },
            base: {
                node: {
                    v1beta1: new (await Promise.resolve().then(() => __importStar(require("./base/node/v1beta1/query.lcd")))).LCDQueryClient({
                        requestClient
                    })
                },
                tendermint: {
                    v1beta1: new (await Promise.resolve().then(() => __importStar(require("./base/tendermint/v1beta1/query.lcd")))).LCDQueryClient({
                        requestClient
                    })
                }
            },
            circuit: {
                v1: new (await Promise.resolve().then(() => __importStar(require("./circuit/v1/query.lcd")))).LCDQueryClient({
                    requestClient
                })
            },
            consensus: {
                v1: new (await Promise.resolve().then(() => __importStar(require("./consensus/v1/query.lcd")))).LCDQueryClient({
                    requestClient
                })
            },
            distribution: {
                v1beta1: new (await Promise.resolve().then(() => __importStar(require("./distribution/v1beta1/query.lcd")))).LCDQueryClient({
                    requestClient
                })
            },
            evidence: {
                v1beta1: new (await Promise.resolve().then(() => __importStar(require("./evidence/v1beta1/query.lcd")))).LCDQueryClient({
                    requestClient
                })
            },
            feegrant: {
                v1beta1: new (await Promise.resolve().then(() => __importStar(require("./feegrant/v1beta1/query.lcd")))).LCDQueryClient({
                    requestClient
                })
            },
            gov: {
                v1: new (await Promise.resolve().then(() => __importStar(require("./gov/v1/query.lcd")))).LCDQueryClient({
                    requestClient
                }),
                v1beta1: new (await Promise.resolve().then(() => __importStar(require("./gov/v1beta1/query.lcd")))).LCDQueryClient({
                    requestClient
                })
            },
            group: {
                v1: new (await Promise.resolve().then(() => __importStar(require("./group/v1/query.lcd")))).LCDQueryClient({
                    requestClient
                })
            },
            mint: {
                v1beta1: new (await Promise.resolve().then(() => __importStar(require("./mint/v1beta1/query.lcd")))).LCDQueryClient({
                    requestClient
                })
            },
            nft: {
                v1beta1: new (await Promise.resolve().then(() => __importStar(require("./nft/v1beta1/query.lcd")))).LCDQueryClient({
                    requestClient
                })
            },
            params: {
                v1beta1: new (await Promise.resolve().then(() => __importStar(require("./params/v1beta1/query.lcd")))).LCDQueryClient({
                    requestClient
                })
            },
            slashing: {
                v1beta1: new (await Promise.resolve().then(() => __importStar(require("./slashing/v1beta1/query.lcd")))).LCDQueryClient({
                    requestClient
                })
            },
            staking: {
                v1beta1: new (await Promise.resolve().then(() => __importStar(require("./staking/v1beta1/query.lcd")))).LCDQueryClient({
                    requestClient
                })
            },
            tx: {
                v1beta1: new (await Promise.resolve().then(() => __importStar(require("./tx/v1beta1/service.lcd")))).LCDQueryClient({
                    requestClient
                })
            },
            upgrade: {
                v1beta1: new (await Promise.resolve().then(() => __importStar(require("./upgrade/v1beta1/query.lcd")))).LCDQueryClient({
                    requestClient
                })
            }
        }
    };
};
exports.createLCDClient = createLCDClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoibGNkLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vY29zbW9zL2xjZC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7OztBQUFBLHdDQUEyQztBQUNwQyxNQUFNLGVBQWUsR0FBRyxLQUFLLEVBQUUsRUFDcEMsWUFBWSxFQUdiLEVBQUUsRUFBRTtJQUNILE1BQU0sYUFBYSxHQUFHLElBQUksZUFBUyxDQUFDO1FBQ2xDLFlBQVk7S0FDYixDQUFDLENBQUM7SUFDSCxPQUFPO1FBQ0wsTUFBTSxFQUFFO1lBQ04sSUFBSSxFQUFFO2dCQUNKLE9BQU8sRUFBRSxJQUFJLENBQUMsd0RBQWEsMEJBQTBCLEdBQUMsQ0FBQyxDQUFDLGNBQWMsQ0FBQztvQkFDckUsYUFBYTtpQkFDZCxDQUFDO2FBQ0g7WUFDRCxLQUFLLEVBQUU7Z0JBQ0wsT0FBTyxFQUFFLElBQUksQ0FBQyx3REFBYSwyQkFBMkIsR0FBQyxDQUFDLENBQUMsY0FBYyxDQUFDO29CQUN0RSxhQUFhO2lCQUNkLENBQUM7YUFDSDtZQUNELElBQUksRUFBRTtnQkFDSixPQUFPLEVBQUUsSUFBSSxDQUFDLHdEQUFhLDBCQUEwQixHQUFDLENBQUMsQ0FBQyxjQUFjLENBQUM7b0JBQ3JFLGFBQWE7aUJBQ2QsQ0FBQzthQUNIO1lBQ0QsSUFBSSxFQUFFO2dCQUNKLElBQUksRUFBRTtvQkFDSixPQUFPLEVBQUUsSUFBSSxDQUFDLHdEQUFhLCtCQUErQixHQUFDLENBQUMsQ0FBQyxjQUFjLENBQUM7d0JBQzFFLGFBQWE7cUJBQ2QsQ0FBQztpQkFDSDtnQkFDRCxVQUFVLEVBQUU7b0JBQ1YsT0FBTyxFQUFFLElBQUksQ0FBQyx3REFBYSxxQ0FBcUMsR0FBQyxDQUFDLENBQUMsY0FBYyxDQUFDO3dCQUNoRixhQUFhO3FCQUNkLENBQUM7aUJBQ0g7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