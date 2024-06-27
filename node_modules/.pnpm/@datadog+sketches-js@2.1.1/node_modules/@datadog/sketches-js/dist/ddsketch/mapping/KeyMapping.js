"use strict";
/*
 * Unless explicitly stated otherwise all files in this repository are licensed
 * under the Apache 2.0 license (see LICENSE).
 * This product includes software developed at Datadog (https://www.datadoghq.com/).
 * Copyright 2020 Datadog, Inc.
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.KeyMapping = void 0;
var index_1 = require("./index");
// 1.1125369292536007e-308
var MIN_SAFE_FLOAT = Math.pow(2, -1023);
var MAX_SAFE_FLOAT = Number.MAX_VALUE;
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
var KeyMapping = /** @class */ (function () {
    function KeyMapping(relativeAccuracy, offset) {
        if (offset === void 0) { offset = 0; }
        if (relativeAccuracy <= 0 || relativeAccuracy >= 1) {
            throw Error('Relative accuracy must be between 0 and 1 when initializing a KeyMapping');
        }
        this.relativeAccuracy = relativeAccuracy;
        this._offset = offset;
        var gammaMantissa = (2 * relativeAccuracy) / (1 - relativeAccuracy);
        this.gamma = 1 + gammaMantissa;
        this._multiplier = 1 / Math.log1p(gammaMantissa);
        this.minPossible = MIN_SAFE_FLOAT * this.gamma;
        this.maxPossible = MAX_SAFE_FLOAT / this.gamma;
    }
    KeyMapping.fromGammaOffset = function (gamma, indexOffset) {
        var relativeAccuracy = (gamma - 1) / (gamma + 1);
        return new this(relativeAccuracy, indexOffset);
    };
    /** Retrieve the key specifying the bucket for a `value` */
    KeyMapping.prototype.key = function (value) {
        return Math.ceil(this._logGamma(value)) + this._offset;
    };
    /** Retrieve the value represented by the bucket at `key` */
    KeyMapping.prototype.value = function (key) {
        return this._powGamma(key - this._offset) * (2 / (1 + this.gamma));
    };
    KeyMapping.prototype.toProto = function () {
        var ProtoIndexMapping = require('../proto/compiled').IndexMapping;
        return ProtoIndexMapping.create({
            gamma: this.gamma,
            indexOffset: this._offset,
            interpolation: this._protoInterpolation()
        });
    };
    KeyMapping.fromProto = function (protoMapping) {
        if (!protoMapping ||
            /* Double equals (==) is intentional here to check for
             * `null` | `undefined` without including `0` */
            protoMapping.gamma == null ||
            protoMapping.indexOffset == null) {
            throw Error('Failed to decode mapping from protobuf');
        }
        var Interpolation = require('../proto/compiled').IndexMapping.Interpolation;
        var interpolation = protoMapping.interpolation, gamma = protoMapping.gamma, indexOffset = protoMapping.indexOffset;
        switch (interpolation) {
            case Interpolation.NONE:
                return index_1.LogarithmicMapping.fromGammaOffset(gamma, indexOffset);
            case Interpolation.LINEAR:
                return index_1.LinearlyInterpolatedMapping.fromGammaOffset(gamma, indexOffset);
            case Interpolation.CUBIC:
                return index_1.CubicallyInterpolatedMapping.fromGammaOffset(gamma, indexOffset);
            default:
                throw Error('Unrecognized mapping when decoding from protobuf');
        }
    };
    /** Return (an approximation of) the logarithm of the value base gamma */
    KeyMapping.prototype._logGamma = function (value) {
        return Math.log2(value) * this._multiplier;
    };
    /** Return (an approximation of) gamma to the power value */
    KeyMapping.prototype._powGamma = function (value) {
        return Math.pow(2, value / this._multiplier);
    };
    KeyMapping.prototype._protoInterpolation = function () {
        var Interpolation = require('../proto/compiled').IndexMapping.Interpolation;
        return Interpolation.NONE;
    };
    return KeyMapping;
}());
exports.KeyMapping = KeyMapping;
