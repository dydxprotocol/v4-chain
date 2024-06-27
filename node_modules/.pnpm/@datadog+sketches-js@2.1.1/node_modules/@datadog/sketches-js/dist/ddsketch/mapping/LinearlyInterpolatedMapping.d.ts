import { KeyMapping } from './KeyMapping';
import { IndexMapping as IndexMappingProto } from '../proto/compiled';
/**
 * A fast KeyMapping that approximates the memory-optimal one
 * (LogarithmicMapping) by extracting the floor value of the logarithm to the
 * base 2 from the binary representations of floating-point values and
 * linearly interpolating the logarithm in-between.
 */
export declare class LinearlyInterpolatedMapping extends KeyMapping {
    constructor(relativeAccuracy: number, offset?: number);
    /**
     * Approximates log2 by s + f
     * where v = (s+1) * 2 ** f  for s in [0, 1)
     *
     * frexp(v) returns m and e s.t.
     * v = m * 2 ** e ; (m in [0.5, 1) or 0.0)
     * so we adjust m and e accordingly
     */
    _log2Approx(value: number): number;
    /** Inverse of _log2Approx */
    _exp2Approx(value: number): number;
    _logGamma(value: number): number;
    _powGamma(value: number): number;
    _protoInterpolation(): IndexMappingProto.Interpolation;
}
