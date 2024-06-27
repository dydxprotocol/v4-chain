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
exports.LogCollapsingHighestDenseDDSketch = exports.LogCollapsingLowestDenseDDSketch = exports.DDSketch = void 0;
var store_1 = require("./store");
var mapping_1 = require("./mapping");
var DEFAULT_RELATIVE_ACCURACY = 0.01;
var DEFAULT_BIN_LIMIT = 2048;
/** Base class for DDSketch*/
var BaseDDSketch = /** @class */ (function () {
    function BaseDDSketch(_a) {
        var mapping = _a.mapping, store = _a.store, negativeStore = _a.negativeStore, zeroCount = _a.zeroCount;
        this.mapping = mapping;
        this.store = store;
        this.negativeStore = negativeStore;
        this.zeroCount = zeroCount;
        this.count =
            this.negativeStore.count + this.zeroCount + this.store.count;
        this.min = Infinity;
        this.max = -Infinity;
        this.sum = 0;
    }
    /**
     * Add a value to the sketch
     *
     * @param value The value to be added
     * @param weight The amount to weight the value (default 1.0)
     *
     * @throws Error if `weight` is 0 or negative
     */
    BaseDDSketch.prototype.accept = function (value, weight) {
        if (weight === void 0) { weight = 1; }
        if (weight <= 0) {
            throw Error('Weight must be a positive number');
        }
        if (value > this.mapping.minPossible) {
            var key = this.mapping.key(value);
            this.store.add(key, weight);
        }
        else if (value < -this.mapping.minPossible) {
            var key = this.mapping.key(-value);
            this.negativeStore.add(key, weight);
        }
        else {
            this.zeroCount += weight;
        }
        /* Keep track of summary stats */
        this.count += weight;
        this.sum += value * weight;
        if (value < this.min) {
            this.min = value;
        }
        if (value > this.max) {
            this.max = value;
        }
    };
    /**
     * Retrieve a value from the sketch at the quantile
     *
     * @param quantile A number between `0` and `1` (inclusive)
     */
    BaseDDSketch.prototype.getValueAtQuantile = function (quantile) {
        if (quantile < 0 || quantile > 1 || this.count === 0) {
            return NaN;
        }
        var rank = quantile * (this.count - 1);
        var quantileValue = 0;
        if (rank < this.negativeStore.count) {
            var reversedRank = this.negativeStore.count - rank - 1;
            var key = this.negativeStore.keyAtRank(reversedRank, false);
            quantileValue = -this.mapping.value(key);
        }
        else if (rank < this.zeroCount + this.negativeStore.count) {
            return 0;
        }
        else {
            var key = this.store.keyAtRank(rank - this.zeroCount - this.negativeStore.count);
            quantileValue = this.mapping.value(key);
        }
        return quantileValue;
    };
    /**
     * Merge the contents of the parameter `sketch` into this sketch
     *
     * @param sketch The sketch to merge into the caller sketch
     * @throws Error if the sketches were initialized with different `relativeAccuracy` values
     */
    BaseDDSketch.prototype.merge = function (sketch) {
        if (!this.mergeable(sketch)) {
            throw new Error('Cannot merge two DDSketches with different `relativeAccuracy` parameters');
        }
        if (sketch.count === 0) {
            return;
        }
        if (this.count === 0) {
            this._copy(sketch);
            return;
        }
        this.store.merge(sketch.store);
        /* Merge summary stats */
        this.zeroCount += sketch.zeroCount;
        this.count += sketch.count;
        this.sum += sketch.sum;
        if (sketch.min < this.min) {
            this.min = sketch.min;
        }
        if (sketch.max > this.max) {
            this.max = sketch.max;
        }
    };
    /**
     * Determine whether two sketches can be merged
     *
     * @param sketch The sketch to be merged into the caller sketch
     */
    BaseDDSketch.prototype.mergeable = function (sketch) {
        return this.mapping.gamma === sketch.mapping.gamma;
    };
    /**
     * Helper method to copy the contents of the parameter `store` into this store
     * @see DDSketch.merge to merge two sketches safely
     *
     * @param store The store to be copied into the caller store
     */
    BaseDDSketch.prototype._copy = function (sketch) {
        this.store.copy(sketch.store);
        this.negativeStore.copy(sketch.negativeStore);
        this.zeroCount = sketch.zeroCount;
        this.min = sketch.min;
        this.max = sketch.max;
        this.count = sketch.count;
        this.sum = sketch.sum;
    };
    /** Serialize a DDSketch to protobuf format */
    BaseDDSketch.prototype.toProto = function () {
        var ProtoDDSketch = require('./proto/compiled').DDSketch;
        var message = ProtoDDSketch.create({
            mapping: this.mapping.toProto(),
            positiveValues: this.store.toProto(),
            negativeValues: this.negativeStore.toProto(),
            zeroCount: this.zeroCount
        });
        return ProtoDDSketch.encode(message).finish();
    };
    /**
     * Deserialize a DDSketch from protobuf data
     *
     * Note: `fromProto` currently loses summary statistics for the original
     * sketch (i.e. `min`, `max`)
     *
     * @param buffer Byte array containing DDSketch in protobuf format (from DDSketch.toProto)
     */
    BaseDDSketch.fromProto = function (buffer) {
        var ProtoDDSketch = require('./proto/compiled').DDSketch;
        var decoded = ProtoDDSketch.decode(buffer);
        var mapping = mapping_1.KeyMapping.fromProto(decoded.mapping);
        var store = store_1.DenseStore.fromProto(decoded.positiveValues);
        var negativeStore = store_1.DenseStore.fromProto(decoded.negativeValues);
        var zeroCount = decoded.zeroCount;
        return new BaseDDSketch({ mapping: mapping, store: store, negativeStore: negativeStore, zeroCount: zeroCount });
    };
    return BaseDDSketch;
}());
var defaultConfig = {
    relativeAccuracy: DEFAULT_RELATIVE_ACCURACY
};
/** A quantile sketch with relative-error guarantees */
var DDSketch = /** @class */ (function (_super) {
    __extends(DDSketch, _super);
    /**
     * Initialize a new DDSketch
     *
     * @param relativeAccuracy The accuracy guarantee of the sketch (default 0.01)
     */
    function DDSketch(_a) {
        var _b = _a === void 0 ? defaultConfig : _a, _c = _b.relativeAccuracy, relativeAccuracy = _c === void 0 ? DEFAULT_RELATIVE_ACCURACY : _c;
        var mapping = new mapping_1.LogarithmicMapping(relativeAccuracy);
        var store = new store_1.DenseStore();
        var negativeStore = new store_1.DenseStore();
        return _super.call(this, { mapping: mapping, store: store, negativeStore: negativeStore, zeroCount: 0 }) || this;
    }
    return DDSketch;
}(BaseDDSketch));
exports.DDSketch = DDSketch;
var LogCollapsingLowestDenseDDSketch = /** @class */ (function (_super) {
    __extends(LogCollapsingLowestDenseDDSketch, _super);
    /**
     * Initialize a new LogCollapsingLowestDenseDDSketch
     *
     * @param relativeAccuracy The accuracy guarantee of the sketch (default 0.01)
     * @param binLimit Number of bins before lowest indices are collapsed (default 2048)
     */
    function LogCollapsingLowestDenseDDSketch(_a) {
        var _b = _a === void 0 ? defaultConfig : _a, _c = _b.relativeAccuracy, relativeAccuracy = _c === void 0 ? DEFAULT_RELATIVE_ACCURACY : _c, _d = _b.binLimit, binLimit = _d === void 0 ? DEFAULT_BIN_LIMIT : _d;
        var mapping = new mapping_1.LogarithmicMapping(relativeAccuracy);
        var store = new store_1.CollapsingLowestDenseStore(binLimit);
        var negativeStore = new store_1.CollapsingLowestDenseStore(binLimit);
        return _super.call(this, { mapping: mapping, store: store, negativeStore: negativeStore, zeroCount: 0 }) || this;
    }
    return LogCollapsingLowestDenseDDSketch;
}(BaseDDSketch));
exports.LogCollapsingLowestDenseDDSketch = LogCollapsingLowestDenseDDSketch;
var LogCollapsingHighestDenseDDSketch = /** @class */ (function (_super) {
    __extends(LogCollapsingHighestDenseDDSketch, _super);
    /**
     * Initialize a new LogCollapsingHighestDenseDDSketch
     *
     * @param relativeAccuracy The accuracy guarantee of the sketch (default 0.01)
     * @param binLimit Number of bins before highest indices are collapsed (default 2048)
     */
    function LogCollapsingHighestDenseDDSketch(_a) {
        var _b = _a === void 0 ? defaultConfig : _a, _c = _b.relativeAccuracy, relativeAccuracy = _c === void 0 ? DEFAULT_RELATIVE_ACCURACY : _c, _d = _b.binLimit, binLimit = _d === void 0 ? DEFAULT_BIN_LIMIT : _d;
        var mapping = new mapping_1.LogarithmicMapping(relativeAccuracy);
        var store = new store_1.CollapsingHighestDenseStore(binLimit);
        var negativeStore = new store_1.CollapsingHighestDenseStore(binLimit);
        return _super.call(this, { mapping: mapping, store: store, negativeStore: negativeStore, zeroCount: 0 }) || this;
    }
    return LogCollapsingHighestDenseDDSketch;
}(BaseDDSketch));
exports.LogCollapsingHighestDenseDDSketch = LogCollapsingHighestDenseDDSketch;
