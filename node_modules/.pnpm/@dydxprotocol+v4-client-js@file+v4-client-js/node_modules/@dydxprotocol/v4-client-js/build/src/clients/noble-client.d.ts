import { EncodeObject, Coin } from '@cosmjs/proto-signing';
import { DeliverTxResponse, GasPrice, StdFee } from '@cosmjs/stargate';
import LocalWallet from './modules/local-wallet';
export declare class NobleClient {
    private wallet?;
    private restEndpoint;
    private stargateClient?;
    constructor(restEndpoint: string);
    get isConnected(): boolean;
    connect(wallet: LocalWallet): Promise<void>;
    getAccountBalances(): Promise<readonly Coin[]>;
    getAccountBalance(denom: string): Promise<Coin>;
    send(messages: EncodeObject[], gasPrice?: GasPrice, memo?: string): Promise<DeliverTxResponse>;
    simulateTransaction(messages: readonly EncodeObject[], gasPrice?: GasPrice, memo?: string): Promise<StdFee>;
}
