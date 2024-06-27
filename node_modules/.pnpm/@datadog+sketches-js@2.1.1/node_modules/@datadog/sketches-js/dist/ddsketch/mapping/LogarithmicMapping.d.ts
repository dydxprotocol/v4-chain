import { KeyMapping } from './KeyMapping';
import { IndexMapping as IndexMappingProto } from '../proto/compiled';
/**
 * A memory-optimal KeyMapping, i.e., given a targeted relative accuracy, it
 * requires the least number of keys to cover a given range of values. This is
 * done by logarithmically mapping floating-point values to integers.
 */
export declare class LogarithmicMapping extends KeyMapping {
    constructor(relativeAccuracy: number, offset?: number);
    _logGamma(value: number): number;
    _powGamma(value: number): number;
    _protoInterpolation(): IndexMappingProto.Interpolation;
}
