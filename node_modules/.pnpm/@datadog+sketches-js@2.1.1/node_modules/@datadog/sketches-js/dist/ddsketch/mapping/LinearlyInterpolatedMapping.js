"use strict";
/*
 * Unless explicitly stated otherwise all files in this repository are licensed
 * under the Apache 2.0 license (see LICENSE).
 * This product includes software developed at Datadog (https://www.datadoghq.com/).
 * Copyright 2020 Datadog, Inc.
 */
var __extends = (this && this.__extends) || (function () {
    var extendStatics = function (d, b) {
        extendStatics = Object.setPrototypeOf ||
            ({ __proto__: [] } instanceof Array && function (d, b) { d.__proto__ = b; }) ||
            function (d, b) { for (var p in b) if (Object.prototype.hasOwnProperty.call(b, p)) d[p] = b[p]; };
        return extendStatics(d, b);
    };
    return function (d, b) {
        if (typeof b !== "function" && b !== null)
            throw new TypeError("Class extends value " + String(b) + " is not a constructor or null");
        extendStatics(d, b);
        function __() { this.constructor = d; }
        d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
    };
})();
Object.defineProperty(exports, "__esModule", { value: true });
exports.LinearlyInterpolatedMapping = void 0;
var KeyMapping_1 = require("./KeyMapping");
var math_1 = require("../math");
/**
 * A fast KeyMapping that approximates the memory-optimal one
 * (LogarithmicMapping) by extracting the floor value of the logarithm to the
 * base 2 from the binary representations of floating-point values and
 * linearly interpolating the logarithm in-between.
 */
var LinearlyInterpolatedMapping = /** @class */ (function (_super) {
    __extends(LinearlyInterpolatedMapping, _super);
    function LinearlyInterpolatedMapping(relativeAccuracy, offset) {
        if (offset === void 0) { offset = 0; }
        return _super.call(this, relativeAccuracy, offset) || this;
    }
    /**
     * Approximates log2 by s + f
     * where v = (s+1) * 2 ** f  for s in [0, 1)
     *
     * frexp(v) returns m and e s.t.
     * v = m * 2 ** e ; (m in [0.5, 1) or 0.0)
     * so we adjust m and e accordingly
     */
    LinearlyInterpolatedMapping.prototype._log2Approx = function (value) {
        var _a = (0, math_1.frexp)(value), mantissa = _a[0], exponent = _a[1];
        var significand = 2 * mantissa - 1;
        return significand + (exponent - 1);
    };
    /** Inverse of _log2Approx */
    LinearlyInterpolatedMapping.prototype._exp2Approx = function (value) {
        var exponent = Math.floor(value) + 1;
        var mantissa = (value - exponent + 2) / 2;
        return (0, math_1.ldexp)(mantissa, exponent);
    };
    LinearlyInterpolatedMapping.prototype._logGamma = function (value) {
        return Math.log2(value) * this._multiplier;
    };
    LinearlyInterpolatedMapping.prototype._powGamma = function (value) {
        return Math.pow(2, value / this._multiplier);
    };
    LinearlyInterpolatedMapping.prototype._protoInterpolation = function () {
        var Interpolation = require('../proto/compiled').IndexMapping.Interpolation;
        return Interpolation.LINEAR;
    };
    return LinearlyInterpolatedMapping;
}(KeyMapping_1.KeyMapping));
exports.LinearlyInterpolatedMapping = LinearlyInterpolatedMapping;
