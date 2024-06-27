import { IIndexMapping, IndexMapping as ProtoIndexMapping } from '../proto/compiled';
import type { Mapping } from './types';
/**
 * A mapping between values and integer indices that imposes relative accuracy
 * guarantees. Specifically, for any value `minPossible() < value <
 * maxPossible` implementations of `KeyMapping` must be such that
 * `value(key(v))` is close to `v` with a relative error that is less than
 * `relativeAccuracy`.
 *
 * In implementations of KeyMapping, there is generally a trade-off between the
 * cost of computing the key and the number of keys that are required to cover a
 * given range of values (memory optimality). The most memory-optimal mapping is
 * the LogarithmicMapping, but it requires the costly evaluation of the logarithm
 * when computing the index. Other mappings can approximate the logarithmic
 * mapping, while being less computationally costly.
 */
export declare class KeyMapping implements Mapping {
    relativeAccuracy: number;
    /** The base for the exponential buckets. gamma = (1 + alpha) / (1 - alpha) */
    gamma: number;
    /** The smallest possible value the sketch can distinguish from 0 */
    minPossible: number;
    /** The largest possible value the sketch can handle */
    maxPossible: number;
    /** Used for calculating _logGamma(value). Initially, _multiplier = 1 / log(gamma) */
    _multiplier: number;
    /** An offset that can be used for shifting all keys */
    _offset: number;
    constructor(relativeAccuracy: number, offset?: number);
    static fromGammaOffset(gamma: number, indexOffset: number): KeyMapping;
    /** Retrieve the key specifying the bucket for a `value` */
    key(value: number): number;
    /** Retrieve the value represented by the bucket at `key` */
    value(key: number): number;
    toProto(): IIndexMapping;
    static fromProto(protoMapping?: IIndexMapping | null): KeyMapping;
    /** Return (an approximation of) the logarithm of the value base gamma */
    _logGamma(value: number): number;
    /** Return (an approximation of) gamma to the power value */
    _powGamma(value: number): number;
    _protoInterpolation(): ProtoIndexMapping.Interpolation;
}
