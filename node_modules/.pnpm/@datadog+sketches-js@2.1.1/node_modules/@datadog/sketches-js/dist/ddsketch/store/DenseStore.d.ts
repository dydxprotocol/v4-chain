import type { Store } from './types';
import { IStore } from '../proto/compiled';
/**
 * `DenseStore` is a store that keeps all the bins between
 * the bin for the `minKey` and the `maxKey`.
 */
export declare class DenseStore implements Store<DenseStore> {
    /** Storage for counts */
    bins: number[];
    /** The total number of values added to the store */
    count: number;
    /** The minimum key bin */
    minKey: number;
    /** The maximum key bin */
    maxKey: number;
    /** The number of bins to grow when necessary */
    chunkSize: number;
    /** The difference between the keys and the index in which they are stored */
    offset: number;
    /**
     * Initialize a new DenseStore
     *
     * @param chunkSize The number of bins to add each time the bins grow (default 128)
     */
    constructor(chunkSize?: number);
    /**
     * Update the counter at the specified index key, growing the number of bins if necessary
     *
     * @param key The key of the index to update
     * @param weight The amount to weight the key (default 1.0)
     */
    add(key: number, weight?: number): void;
    /**
     * Return the key for the value at the given rank
     *
     * E.g., if the non-zero bins are [1, 1] for keys a, b with no offset
     *
     * if lower = True:
     *     keyAtRank(x) = a for x in [0, 1)
     *     keyAtRank(x) = b for x in [1, 2)
     * if lower = False:
     *     keyAtRank(x) = a for x in (-1, 0]
     *     keyAtRank(x) = b for x in (0, 1]
     *
     * @param rank The rank at which to retrieve the key
     */
    keyAtRank(rank: number, lower?: boolean): number;
    /**
     * Merge the contents of the parameter `store` into this store
     *
     * @param store The store to merge into the caller store
     */
    merge(store: DenseStore): void;
    /**
     * Directly clone the contents of the parameter `store` into this store
     *
     * @param store The store to be copied into the caller store
     */
    copy(store: DenseStore): void;
    /**
     * Return the length of the underlying storage (`bins`)
     */
    length(): number;
    _getNewLength(newMinKey: number, newMaxKey: number): number;
    /**
     * Adjust the `bins`, the `offset`, the `minKey`, and the `maxKey`
     * without resizing the bins, in order to try to make it fit the specified range.
     * Collapse to the left if necessary
     */
    _adjust(newMinKey: number, newMaxKey: number): void;
    /** Shift the bins by `shift`. This changes the `offset` */
    _shiftBins(shift: number): void;
    /** Center the bins. This changes the `offset` */
    _centerBins(newMinKey: number, newMaxKey: number): void;
    /** Grow the bins as necessary, and call _adjust */
    _extendRange(key: number, secondKey?: number): void;
    /** Calculate the bin index for the key, extending the range if necessary */
    _getIndex(key: number): number;
    toProto(): IStore;
    static fromProto(protoStore?: IStore | null): DenseStore;
}
