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
exports.CollapsingLowestDenseStore = void 0;
var DenseStore_1 = require("./DenseStore");
var util_1 = require("./util");
/**
 * `CollapsingLowestDenseStore` is a dense store that keeps all the bins between
 * the bin for the `minKey` and the `maxKey`, but collapsing the left-most bins
 * if the number of bins exceeds `binLimit`
 */
var CollapsingLowestDenseStore = /** @class */ (function (_super) {
    __extends(CollapsingLowestDenseStore, _super);
    /**
     * Initialize a new CollapsingLowestDenseStore
     *
     * @param binLimit The maximum number of bins
     * @param chunkSize The number of bins to add each time the bins grow (default 128)
     */
    function CollapsingLowestDenseStore(binLimit, chunkSize) {
        var _this = _super.call(this, chunkSize) || this;
        _this.binLimit = binLimit;
        _this.isCollapsed = false;
        return _this;
    }
    /**
     * Merge the contents of the parameter `store` into this store
     *
     * @param store The store to merge into the caller store
     */
    CollapsingLowestDenseStore.prototype.merge = function (store) {
        if (store.count === 0) {
            return;
        }
        if (this.count === 0) {
            this.copy(store);
            return;
        }
        if (store.minKey < this.minKey || store.maxKey > this.maxKey) {
            this._extendRange(store.minKey, store.maxKey);
        }
        var collapseStartIndex = store.minKey - store.offset;
        var collapseEndIndex = Math.min(this.minKey, store.maxKey + 1) - store.offset;
        if (collapseEndIndex > collapseStartIndex) {
            var collapseCount = (0, util_1.sumOfRange)(store.bins, collapseStartIndex, collapseEndIndex);
            this.bins[0] += collapseCount;
        }
        else {
            collapseEndIndex = collapseStartIndex;
        }
        for (var key = collapseEndIndex + store.offset; key < store.maxKey + 1; key++) {
            this.bins[key - this.offset] += store.bins[key - store.offset];
        }
        this.count += store.count;
    };
    /**
     * Directly clone the contents of the parameter `store` into this store
     *
     * @param store The store to be copied into the caller store
     */
    CollapsingLowestDenseStore.prototype.copy = function (store) {
        _super.prototype.copy.call(this, store);
        this.isCollapsed = store.isCollapsed;
    };
    CollapsingLowestDenseStore.prototype._getNewLength = function (newMinKey, newMaxKey) {
        var desiredLength = newMaxKey - newMinKey + 1;
        return Math.min(this.chunkSize * Math.ceil(desiredLength / this.chunkSize), this.binLimit);
    };
    /**
     * Adjust the `bins`, the `offset`, the `minKey`, and the `maxKey`
     * without resizing the bins, in order to try to make it fit the specified range.
     * Collapse to the left if necessary
     */
    CollapsingLowestDenseStore.prototype._adjust = function (newMinKey, newMaxKey) {
        if (newMaxKey - newMinKey + 1 > this.length()) {
            // The range of keys is too wide, the lowest bins need to be collapsed
            newMinKey = newMaxKey - this.length() + 1;
            if (newMinKey >= this.maxKey) {
                // Put everything in the first bin
                this.offset = newMinKey;
                this.minKey = newMinKey;
                this.bins.fill(0);
                this.bins[0] = this.count;
            }
            else {
                var shift = this.offset - newMinKey;
                if (shift < 0) {
                    var collapseStartIndex = this.minKey - this.offset;
                    var collapseEndIndex = newMinKey - this.offset;
                    var collapsedCount = (0, util_1.sumOfRange)(this.bins, collapseStartIndex, collapseEndIndex);
                    this.bins.fill(0, collapseStartIndex, collapseEndIndex);
                    this.bins[collapseEndIndex] += collapsedCount;
                    this.minKey = newMinKey;
                    this._shiftBins(shift);
                }
                else {
                    this.minKey = newMinKey;
                    // Shift the buckets to make room for newMinKey
                    this._shiftBins(shift);
                }
            }
            this.maxKey = newMaxKey;
            this.isCollapsed = true;
        }
        else {
            this._centerBins(newMinKey, newMaxKey);
            this.minKey = newMinKey;
            this.maxKey = newMaxKey;
        }
    };
    /** Calculate the bin index for the key, extending the range if necessary */
    CollapsingLowestDenseStore.prototype._getIndex = function (key) {
        if (key < this.minKey) {
            if (this.isCollapsed) {
                return 0;
            }
            this._extendRange(key);
            if (this.isCollapsed) {
                return 0;
            }
        }
        else if (key > this.maxKey) {
            this._extendRange(key);
        }
        return key - this.offset;
    };
    return CollapsingLowestDenseStore;
}(DenseStore_1.DenseStore));
exports.CollapsingLowestDenseStore = CollapsingLowestDenseStore;
