import { DenseStore } from './store';
import { Mapping } from './mapping';
interface BaseSketchConfig {
    /** The mapping between values and indicies for the sketch */
    mapping: Mapping;
    /** Storage for positive values */
    store: DenseStore;
    /** Storage for negative values */
    negativeStore: DenseStore;
    /** The number of zeroes added to the sketch */
    zeroCount: number;
}
/** Base class for DDSketch*/
declare class BaseDDSketch {
    /** The mapping between values and indicies for the sketch */
    mapping: Mapping;
    /** Storage for positive values */
    store: DenseStore;
    /** Storage for negative values */
    negativeStore: DenseStore;
    /** The count of zero values */
    zeroCount: number;
    /** The minimum value seen by the sketch */
    min: number;
    /** The maximum value seen by the sketch */
    max: number;
    /** The total number of values seen by the sketch */
    count: number;
    /** The sum of the values seen by the sketch */
    sum: number;
    constructor({ mapping, store, negativeStore, zeroCount }: BaseSketchConfig);
    /**
     * Add a value to the sketch
     *
     * @param value The value to be added
     * @param weight The amount to weight the value (default 1.0)
     *
     * @throws Error if `weight` is 0 or negative
     */
    accept(value: number, weight?: number): void;
    /**
     * Retrieve a value from the sketch at the quantile
     *
     * @param quantile A number between `0` and `1` (inclusive)
     */
    getValueAtQuantile(quantile: number): number;
    /**
     * Merge the contents of the parameter `sketch` into this sketch
     *
     * @param sketch The sketch to merge into the caller sketch
     * @throws Error if the sketches were initialized with different `relativeAccuracy` values
     */
    merge(sketch: DDSketch): void;
    /**
     * Determine whether two sketches can be merged
     *
     * @param sketch The sketch to be merged into the caller sketch
     */
    mergeable(sketch: DDSketch): boolean;
    /**
     * Helper method to copy the contents of the parameter `store` into this store
     * @see DDSketch.merge to merge two sketches safely
     *
     * @param store The store to be copied into the caller store
     */
    _copy(sketch: DDSketch): void;
    /** Serialize a DDSketch to protobuf format */
    toProto(): Uint8Array;
    /**
     * Deserialize a DDSketch from protobuf data
     *
     * Note: `fromProto` currently loses summary statistics for the original
     * sketch (i.e. `min`, `max`)
     *
     * @param buffer Byte array containing DDSketch in protobuf format (from DDSketch.toProto)
     */
    static fromProto(buffer: Uint8Array): DDSketch;
}
interface SketchConfig {
    /** The accuracy guarantee of the sketch, between 0-1 (default 0.01) */
    relativeAccuracy?: number;
}
/** A quantile sketch with relative-error guarantees */
export declare class DDSketch extends BaseDDSketch {
    /**
     * Initialize a new DDSketch
     *
     * @param relativeAccuracy The accuracy guarantee of the sketch (default 0.01)
     */
    constructor({ relativeAccuracy }?: SketchConfig);
}
interface LogCollapsingSketchConfig {
    /** The accuracy guarantee of the sketch, between 0-1 (default 0.01) */
    relativeAccuracy?: number;
    binLimit?: number;
}
export declare class LogCollapsingLowestDenseDDSketch extends BaseDDSketch {
    /**
     * Initialize a new LogCollapsingLowestDenseDDSketch
     *
     * @param relativeAccuracy The accuracy guarantee of the sketch (default 0.01)
     * @param binLimit Number of bins before lowest indices are collapsed (default 2048)
     */
    constructor({ relativeAccuracy, binLimit }?: LogCollapsingSketchConfig);
}
export declare class LogCollapsingHighestDenseDDSketch extends BaseDDSketch {
    /**
     * Initialize a new LogCollapsingHighestDenseDDSketch
     *
     * @param relativeAccuracy The accuracy guarantee of the sketch (default 0.01)
     * @param binLimit Number of bins before highest indices are collapsed (default 2048)
     */
    constructor({ relativeAccuracy, binLimit }?: LogCollapsingSketchConfig);
}
export {};
