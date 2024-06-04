export declare class BaseTendermintClient {
    block(): Promise<void>;
    broadcastTxSync(): Promise<void>;
    broadcastTxAsync(): Promise<void>;
    txSearchAll(): Promise<void>;
}
export declare class BaseQueryClient {
    tx: {
        simulate(): Promise<void>;
    };
    queryUnverified(): Promise<void>;
}
export declare class BaseStargateSigningClient {
    sign(): Promise<void>;
}
export declare class BaseWallet {
    getAccounts(): Promise<void>;
}
