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
exports.LogarithmicMapping = void 0;
var KeyMapping_1 = require("./KeyMapping");
/**
 * A memory-optimal KeyMapping, i.e., given a targeted relative accuracy, it
 * requires the least number of keys to cover a given range of values. This is
 * done by logarithmically mapping floating-point values to integers.
 */
var LogarithmicMapping = /** @class */ (function (_super) {
    __extends(LogarithmicMapping, _super);
    function LogarithmicMapping(relativeAccuracy, offset) {
        if (offset === void 0) { offset = 0; }
        var _this = _super.call(this, relativeAccuracy, offset) || this;
        _this._multiplier *= Math.log(2);
        return _this;
    }
    LogarithmicMapping.prototype._logGamma = function (value) {
        return Math.log2(value) * this._multiplier;
    };
    LogarithmicMapping.prototype._powGamma = function (value) {
        return Math.pow(2, value / this._multiplier);
    };
    LogarithmicMapping.prototype._protoInterpolation = function () {
        var Interpolation = require('../proto/compiled').IndexMapping.Interpolation;
        return Interpolation.NONE;
    };
    return LogarithmicMapping;
}(KeyMapping_1.KeyMapping));
exports.LogarithmicMapping = LogarithmicMapping;
