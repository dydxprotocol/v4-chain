import { Order_Side, Order_TimeInForce } from '@dydxprotocol/v4-proto/src/codegen/dydxprotocol/clob/order';
import Long from 'long';
import { CompositeClient } from './composite-client';
import { Network, OrderType, OrderSide, OrderTimeInForce, OrderExecution } from './constants';
import { FaucetClient } from './faucet-client';
import LocalWallet from './modules/local-wallet';
import { NobleClient } from './noble-client';
import { OrderFlags } from './types';
declare global {
    var client: CompositeClient;
    var faucetClient: FaucetClient | null;
    var wallet: LocalWallet;
    var nobleClient: NobleClient | undefined;
    var nobleWallet: LocalWallet | undefined;
}
export declare function connectClient(network: Network): Promise<string>;
export declare function connectNetwork(paramsJSON: string): Promise<string>;
export declare function connectWallet(mnemonic: string): Promise<string>;
export declare function connect(network: Network, mnemonic: string): Promise<string>;
export declare function deriveMnemomicFromEthereumSignature(signature: string): Promise<string>;
export declare function getHeight(): Promise<string>;
export declare function getFeeTiers(): Promise<string>;
export declare function getUserFeeTier(address: string): Promise<string>;
export declare function getEquityTiers(): Promise<string>;
export declare function getPerpetualMarkets(): Promise<string>;
export declare function placeOrder(payload: string): Promise<string>;
export declare function wrappedError(error: Error): string;
export declare function cancelOrder(payload: string): Promise<string>;
export declare function deposit(payload: string): Promise<string>;
export declare function withdraw(payload: string): Promise<string>;
export declare function faucet(payload: string): Promise<string>;
export declare function withdrawToIBC(subaccountNumber: number, amount: string, payload: string): Promise<string>;
export declare function transferNativeToken(payload: string): Promise<string>;
export declare function getAccountBalance(): Promise<String>;
export declare function getAccountBalances(): Promise<String>;
export declare function getUserStats(payload: string): Promise<String>;
export declare function simulateDeposit(payload: string): Promise<string>;
export declare function simulateWithdraw(payload: string): Promise<string>;
export declare function simulateTransferNativeToken(payload: string): Promise<string>;
export declare function signRawPlaceOrder(subaccountNumber: number, clientId: number, clobPairId: number, side: Order_Side, quantums: Long, subticks: Long, timeInForce: Order_TimeInForce, orderFlags: number, reduceOnly: boolean, goodTilBlock: number, goodTilBlockTime: number, clientMetadata: number, routerFeePpm?: number, routerFeeSubaccountOwner?: string, routerFeeSubaccountNumber?: number): Promise<string>;
export declare function signPlaceOrder(subaccountNumber: number, marketId: string, type: OrderType, side: OrderSide, price: number, size: number, clientId: number, timeInForce: OrderTimeInForce, goodTilTimeInSeconds: number, execution: OrderExecution, postOnly: boolean, reduceOnly: boolean, routerFeePpm?: number, routerFeeSubaccountOwner?: string, routerFeeSubaccountNumber?: number): Promise<string>;
export declare function signCancelOrder(subaccountNumber: number, clientId: number, orderFlags: OrderFlags, clobPairId: number, goodTilBlock: number, goodTilBlockTime: number): Promise<string>;
export declare function encodeAccountRequestData(address: string): Promise<string>;
export declare function decodeAccountResponseValue(value: string): Promise<string>;
export declare function getOptimalNode(endpointUrlsAsJson: string): Promise<string>;
export declare function getOptimalIndexer(endpointUrlsAsJson: string): Promise<string>;
export declare function getRewardsParams(): Promise<string>;
export declare function getDelegatorDelegations(payload: string): Promise<string>;
export declare function getDelegatorUnbondingDelegations(payload: string): Promise<string>;
export declare function getMarketPrice(payload: string): Promise<string>;
export declare function getNobleBalance(): Promise<String>;
export declare function sendNobleIBC(squidPayload: string): Promise<String>;
export declare function withdrawToNobleIBC(payload: string): Promise<String>;
export declare function cctpWithdraw(squidPayload: string): Promise<String>;
