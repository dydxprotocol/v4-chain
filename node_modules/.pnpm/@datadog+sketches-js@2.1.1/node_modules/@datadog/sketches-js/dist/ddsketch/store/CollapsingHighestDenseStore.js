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
exports.CollapsingHighestDenseStore = void 0;
var DenseStore_1 = require("./DenseStore");
var util_1 = require("./util");
/**
 * `CollapsingHighestDenseStore` is a dense store that keeps all the bins between
 * the bin for the `minKey` and the `maxKey`, but collapsing the left-most bins
 * if the number of bins exceeds `binLimit`
 */
var CollapsingHighestDenseStore = /** @class */ (function (_super) {
    __extends(CollapsingHighestDenseStore, _super);
    /**
     * Initialize a new CollapsingHighestDenseStore
     *
     * @param binLimit The maximum number of bins
     * @param chunkSize The number of bins to add each time the bins grow (default 128)
     */
    function CollapsingHighestDenseStore(binLimit, chunkSize) {
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
    CollapsingHighestDenseStore.prototype.merge = function (store) {
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
        var collapseEndIndex = store.maxKey - store.offset + 1;
        var collapseStartIndex = Math.max(this.maxKey + 1, store.minKey) - store.offset;
        if (collapseEndIndex > collapseStartIndex) {
            var collapseCount = (0, util_1.sumOfRange)(store.bins, collapseStartIndex, collapseEndIndex);
            this.bins[this.length() - 1] += collapseCount;
        }
        else {
            collapseStartIndex = collapseEndIndex;
        }
        for (var key = store.minKey; key < collapseStartIndex + store.offset; key++) {
            this.bins[key - this.offset] += store.bins[key - store.offset];
        }
        this.count += store.count;
    };
    /**
     * Directly clone the contents of the parameter `store` into this store
     *
     * @param store The store to be copied into the caller store
     */
    CollapsingHighestDenseStore.prototype.copy = function (store) {
        _super.prototype.copy.call(this, store);
        this.isCollapsed = store.isCollapsed;
    };
    CollapsingHighestDenseStore.prototype._getNewLength = function (newMinKey, newMaxKey) {
        var desiredLength = newMaxKey - newMinKey + 1;
        return Math.min(this.chunkSize * Math.ceil(desiredLength / this.chunkSize), this.binLimit);
    };
    /**
     * Adjust the `bins`, the `offset`, the `minKey`, and the `maxKey`
     * without resizing the bins, in order to try to make it fit the specified range.
     * Collapse to the left if necessary
     */
    CollapsingHighestDenseStore.prototype._adjust = function (newMinKey, newMaxKey) {
        if (newMaxKey - newMinKey + 1 > this.length()) {
            // The range of keys is too wide, the lowest bins need to be collapsed
            newMaxKey = newMinKey + this.length() + 1;
            if (newMaxKey <= this.minKey) {
                // Put everything in the first bin
                this.offset = newMinKey;
                this.maxKey = newMaxKey;
                this.bins.fill(0);
                this.bins[this.length() - 1] = this.count;
            }
            else {
                var shift = this.offset - newMinKey;
                if (shift > 0) {
                    var collapseStartIndex = newMaxKey - this.offset + 1;
                    var collapseEndIndex = this.maxKey - this.offset + 1;
                    var collapsedCount = (0, util_1.sumOfRange)(this.bins, collapseStartIndex, collapseEndIndex);
                    this.bins.fill(0, collapseStartIndex, collapseEndIndex);
                    this.bins[collapseStartIndex - 1] += collapsedCount;
                    this.maxKey = newMaxKey;
                    this._shiftBins(shift);
                }
                else {
                    this.maxKey = newMaxKey;
                    // Shift the buckets to make room for newMinKey
                    this._shiftBins(shift);
                }
                this.minKey = newMinKey;
                this.isCollapsed = true;
            }
        }
        else {
            this._centerBins(newMinKey, newMaxKey);
            this.minKey = newMinKey;
            this.maxKey = newMaxKey;
        }
    };
    /** Calculate the bin index for the key, extending the range if necessary */
    CollapsingHighestDenseStore.prototype._getIndex = function (key) {
        if (key < this.minKey) {
            if (this.isCollapsed) {
                return this.length() - 1;
            }
            this._extendRange(key);
            if (this.isCollapsed) {
                return this.length() - 1;
            }
        }
        else if (key > this.maxKey) {
            this._extendRange(key);
        }
        return key - this.offset;
    };
    return CollapsingHighestDenseStore;
}(DenseStore_1.DenseStore));
exports.CollapsingHighestDenseStore = CollapsingHighestDenseStore;
