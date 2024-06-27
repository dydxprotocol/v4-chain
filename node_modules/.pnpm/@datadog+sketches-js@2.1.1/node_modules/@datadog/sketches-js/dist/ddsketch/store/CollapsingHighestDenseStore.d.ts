import { DenseStore } from './DenseStore';
/**
 * `CollapsingHighestDenseStore` is a dense store that keeps all the bins between
 * the bin for the `minKey` and the `maxKey`, but collapsing the left-most bins
 * if the number of bins exceeds `binLimit`
 */
export declare class CollapsingHighestDenseStore extends DenseStore {
    /** The maximum number of bins */
    binLimit: number;
    /** Whether the store has been collapsed to make room for additional keys */
    isCollapsed: boolean;
    /**
     * Initialize a new CollapsingHighestDenseStore
     *
     * @param binLimit The maximum number of bins
     * @param chunkSize The number of bins to add each time the bins grow (default 128)
     */
    constructor(binLimit: number, chunkSize?: number);
    /**
     * Merge the contents of the parameter `store` into this store
     *
     * @param store The store to merge into the caller store
     */
    merge(store: CollapsingHighestDenseStore): void;
    /**
     * Directly clone the contents of the parameter `store` into this store
     *
     * @param store The store to be copied into the caller store
     */
    copy(store: CollapsingHighestDenseStore): void;
    _getNewLength(newMinKey: number, newMaxKey: number): number;
    /**
     * Adjust the `bins`, the `offset`, the `minKey`, and the `maxKey`
     * without resizing the bins, in order to try to make it fit the specified range.
     * Collapse to the left if necessary
     */
    _adjust(newMinKey: number, newMaxKey: number): void;
    /** Calculate the bin index for the key, extending the range if necessary */
    _getIndex(key: number): number;
}
