export declare const isTruthy: <T>(n?: false | 0 | T | null | undefined) => n is T;
export declare class NetworkOptimizer {
    private validatorClients;
    private indexerClients;
    findOptimalNode(endpointUrls: string[], chainId: string): Promise<string>;
    findOptimalIndexer(endpointUrls: string[]): Promise<string>;
}
