/** A key value pair for searching transactions */
export interface SearchPair {
    readonly key: string;
    readonly value: string | number | bigint;
}
/**
 * This query type allows you to pass arbitrary key/value pairs to the backend.
 */
export type SearchTxQuery = string | readonly SearchPair[];
export declare function isSearchTxQueryArray(query: SearchTxQuery): query is readonly SearchPair[];
