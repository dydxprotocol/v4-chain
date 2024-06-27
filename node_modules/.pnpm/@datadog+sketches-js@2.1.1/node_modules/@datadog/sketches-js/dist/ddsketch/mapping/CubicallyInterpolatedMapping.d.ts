import { KeyMapping } from './KeyMapping';
import { IndexMapping as IndexMappingProto } from '../proto/compiled';
/**
 * A fast KeyMapping that approximates the memory-optimal LogarithmicMapping by
 * extracting the floor value of the logarithm to the base 2 from the binary
 * representations of floating-point values and cubically interpolating the
 * logarithm in-between.
 *
 * More detailed documentation of this method can be found in:
 * <a href="https://github.com/DataDog/sketches-java/">sketches-java</a>
 */
export declare class CubicallyInterpolatedMapping extends KeyMapping {
    A: number;
    B: number;
    C: number;
    constructor(relativeAccuracy: number, offset?: number);
    /** Approximates log2 using a cubic polynomial */
    _cubicLog2Approx(value: number): number;
    /** Derived from Cardano's formula */
    _cubicExp2Approx(value: number): number;
    _logGamma(value: number): number;
    _powGamma(value: number): number;
    _protoInterpolation(): IndexMappingProto.Interpolation;
}
