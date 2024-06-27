/**
 * Splits a double-precision floating-point number into a normalized fraction
 * and an integer power of two.
 */
export declare function frexp(value: number): [number, number];
/**
 * Multiplies a double-precision floating-point number by an integer power of
 * two; i.e., x = frac * 2^exp.
 */
export declare function ldexp(mantissa: number, exponent: number): number;
