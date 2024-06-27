"use strict";
/*
 * Unless explicitly stated otherwise all files in this repository are licensed
 * under the Apache 2.0 license (see LICENSE).
 * This product includes software developed at Datadog (https://www.datadoghq.com/).
 * Copyright 2020 Datadog, Inc.
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.ldexp = exports.frexp = void 0;
/**
 * Splits a double-precision floating-point number into a normalized fraction
 * and an integer power of two.
 */
function frexp(value) {
    if (value === 0 || !Number.isFinite(value))
        return [value, 0];
    var absValue = Math.abs(value);
    var exponent = Math.max(-1023, Math.floor(Math.log2(absValue)) + 1);
    var mantissa = absValue * Math.pow(2, -exponent);
    while (mantissa < 0.5) {
        mantissa *= 2;
        exponent--;
    }
    while (mantissa >= 1) {
        mantissa *= 0.5;
        exponent++;
    }
    if (value < 0) {
        mantissa = -mantissa;
    }
    return [mantissa, exponent];
}
exports.frexp = frexp;
/**
 * Multiplies a double-precision floating-point number by an integer power of
 * two; i.e., x = frac * 2^exp.
 */
function ldexp(mantissa, exponent) {
    var iterations = Math.min(3, Math.ceil(Math.abs(exponent) / 1023));
    var result = mantissa;
    for (var i = 0; i < iterations; i++) {
        result *= Math.pow(2, Math.floor((exponent + i) / iterations));
    }
    return result;
}
exports.ldexp = ldexp;
