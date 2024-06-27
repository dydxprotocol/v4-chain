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
exports.CubicallyInterpolatedMapping = void 0;
var KeyMapping_1 = require("./KeyMapping");
var math_1 = require("../math");
/**
 * A fast KeyMapping that approximates the memory-optimal LogarithmicMapping by
 * extracting the floor value of the logarithm to the base 2 from the binary
 * representations of floating-point values and cubically interpolating the
 * logarithm in-between.
 *
 * More detailed documentation of this method can be found in:
 * <a href="https://github.com/DataDog/sketches-java/">sketches-java</a>
 */
var CubicallyInterpolatedMapping = /** @class */ (function (_super) {
    __extends(CubicallyInterpolatedMapping, _super);
    function CubicallyInterpolatedMapping(relativeAccuracy, offset) {
        if (offset === void 0) { offset = 0; }
        var _this = _super.call(this, relativeAccuracy, offset) || this;
        _this.A = 6 / 35;
        _this.B = -3 / 5;
        _this.C = 10 / 7;
        _this._multiplier /= _this.C;
        return _this;
    }
    /** Approximates log2 using a cubic polynomial */
    CubicallyInterpolatedMapping.prototype._cubicLog2Approx = function (value) {
        var _a = (0, math_1.frexp)(value), mantissa = _a[0], exponent = _a[1];
        var significand = 2 * mantissa - 1;
        return (((this.A * significand + this.B) * significand + this.C) *
            significand +
            (exponent - 1));
    };
    /** Derived from Cardano's formula */
    CubicallyInterpolatedMapping.prototype._cubicExp2Approx = function (value) {
        var exponent = Math.floor(value);
        var delta0 = this.B * this.B - 3 * this.A * this.C;
        var delta1 = 2 * this.B * this.B * this.B -
            9 * this.A * this.B * this.C -
            27 * this.A * this.A * (value - exponent);
        var cardano = Math.cbrt((delta1 -
            Math.sqrt(delta1 * delta1 - 4 * delta0 * delta0 * delta0)) /
            2);
        var significandPlusOne = -(this.B + cardano + delta0 / cardano) / (3 * this.A) + 1;
        var mantissa = significandPlusOne / 2;
        return (0, math_1.ldexp)(mantissa, exponent + 1);
    };
    CubicallyInterpolatedMapping.prototype._logGamma = function (value) {
        return this._cubicLog2Approx(value) * this._multiplier;
    };
    CubicallyInterpolatedMapping.prototype._powGamma = function (value) {
        return this._cubicExp2Approx(value / this._multiplier);
    };
    CubicallyInterpolatedMapping.prototype._protoInterpolation = function () {
        var Interpolation = require('../proto/compiled').IndexMapping.Interpolation;
        return Interpolation.CUBIC;
    };
    return CubicallyInterpolatedMapping;
}(KeyMapping_1.KeyMapping));
exports.CubicallyInterpolatedMapping = CubicallyInterpolatedMapping;
