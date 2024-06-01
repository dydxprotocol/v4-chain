/**
 * Helper to work around missing CORS support in Tendermint (https://github.com/tendermint/tendermint/pull/2800)
 *
 * For some reason, fetch does not complain about missing server-side CORS support.
 */
export declare function http(method: "POST", url: string, headers: Record<string, string> | undefined, request?: any): Promise<any>;
