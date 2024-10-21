export declare const isTruthy: <T>(n?: T | false | null | undefined | 0) => n is T;
export declare class NetworkOptimizer {
    private validatorClients;
    private indexerClients;
    findOptimalNode(endpointUrls: string[], chainId: string): Promise<string>;
    findOptimalIndexer(endpointUrls: string[]): Promise<string>;
}
