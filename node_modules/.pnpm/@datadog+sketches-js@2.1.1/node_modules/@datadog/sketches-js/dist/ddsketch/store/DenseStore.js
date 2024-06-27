"use strict";
/*
 * Unless explicitly stated otherwise all files in this repository are licensed
 * under the Apache 2.0 license (see LICENSE).
 * This product includes software developed at Datadog (https://www.datadoghq.com/).
 * Copyright 2020 Datadog, Inc.
 */
var __spreadArray = (this && this.__spreadArray) || function (to, from, pack) {
    if (pack || arguments.length === 2) for (var i = 0, l = from.length, ar; i < l; i++) {
        if (ar || !(i in from)) {
            if (!ar) ar = Array.prototype.slice.call(from, 0, i);
            ar[i] = from[i];
        }
    }
    return to.concat(ar || Array.prototype.slice.call(from));
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.DenseStore = void 0;
var util_1 = require("./util");
/** The default number of bins to grow when necessary */
var CHUNK_SIZE = 128;
/**
 * `DenseStore` is a store that keeps all the bins between
 * the bin for the `minKey` and the `maxKey`.
 */
var DenseStore = /** @class */ (function () {
    /**
     * Initialize a new DenseStore
     *
     * @param chunkSize The number of bins to add each time the bins grow (default 128)
     */
    function DenseStore(chunkSize) {
        if (chunkSize === void 0) { chunkSize = CHUNK_SIZE; }
        this.chunkSize = chunkSize;
        this.bins = [];
        this.count = 0;
        this.minKey = Infinity;
        this.maxKey = -Infinity;
        this.offset = 0;
    }
    /**
     * Update the counter at the specified index key, growing the number of bins if necessary
     *
     * @param key The key of the index to update
     * @param weight The amount to weight the key (default 1.0)
     */
    DenseStore.prototype.add = function (key, weight) {
        if (weight === void 0) { weight = 1; }
        var index = this._getIndex(key);
        this.bins[index] += weight;
        this.count += weight;
    };
    /**
     * Return the key for the value at the given rank
     *
     * E.g., if the non-zero bins are [1, 1] for keys a, b with no offset
     *
     * if lower = True:
     *     keyAtRank(x) = a for x in [0, 1)
     *     keyAtRank(x) = b for x in [1, 2)
     * if lower = False:
     *     keyAtRank(x) = a for x in (-1, 0]
     *     keyAtRank(x) = b for x in (0, 1]
     *
     * @param rank The rank at which to retrieve the key
     */
    DenseStore.prototype.keyAtRank = function (rank, lower) {
        if (lower === void 0) { lower = true; }
        var runningCount = 0;
        for (var i = 0; i < this.length(); i++) {
            var bin = this.bins[i];
            runningCount += bin;
            if ((lower && runningCount > rank) ||
                (!lower && runningCount >= rank + 1)) {
                return i + this.offset;
            }
        }
        return this.maxKey;
    };
    /**
     * Merge the contents of the parameter `store` into this store
     *
     * @param store The store to merge into the caller store
     */
    DenseStore.prototype.merge = function (store) {
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
    DenseStore.prototype.copy = function (store) {
        this.bins = __spreadArray([], store.bins, true);
        this.count = store.count;
        this.minKey = store.minKey;
        this.maxKey = store.maxKey;
        this.offset = store.offset;
    };
    /**
     * Return the length of the underlying storage (`bins`)
     */
    DenseStore.prototype.length = function () {
        return this.bins.length;
    };
    DenseStore.prototype._getNewLength = function (newMinKey, newMaxKey) {
        var desiredLength = newMaxKey - newMinKey + 1;
        return this.chunkSize * Math.ceil(desiredLength / this.chunkSize);
    };
    /**
     * Adjust the `bins`, the `offset`, the `minKey`, and the `maxKey`
     * without resizing the bins, in order to try to make it fit the specified range.
     * Collapse to the left if necessary
     */
    DenseStore.prototype._adjust = function (newMinKey, newMaxKey) {
        this._centerBins(newMinKey, newMaxKey);
        this.minKey = newMinKey;
        this.maxKey = newMaxKey;
    };
    /** Shift the bins by `shift`. This changes the `offset` */
    DenseStore.prototype._shiftBins = function (shift) {
        var _a, _b;
        if (shift > 0) {
            this.bins = this.bins.slice(0, -shift);
            (_a = this.bins).unshift.apply(_a, new Array(shift).fill(0));
        }
        else {
            this.bins = this.bins.slice(Math.abs(shift));
            (_b = this.bins).push.apply(_b, new Array(Math.abs(shift)).fill(0));
        }
        this.offset -= shift;
    };
    /** Center the bins. This changes the `offset` */
    DenseStore.prototype._centerBins = function (newMinKey, newMaxKey) {
        var middleKey = newMinKey + Math.floor((newMaxKey - newMinKey + 1) / 2);
        this._shiftBins(Math.floor(this.offset + this.length() / 2) - middleKey);
    };
    /** Grow the bins as necessary, and call _adjust */
    DenseStore.prototype._extendRange = function (key, secondKey) {
        var _a;
        secondKey = secondKey || key;
        var newMinKey = Math.min(key, secondKey, this.minKey);
        var newMaxKey = Math.max(key, secondKey, this.maxKey);
        if (this.length() === 0) {
            this.bins = new Array(this._getNewLength(newMinKey, newMaxKey)).fill(0);
            this.offset = newMinKey;
            this._adjust(newMinKey, newMaxKey);
        }
        else if (newMinKey >= this.minKey &&
            newMaxKey < this.offset + this.length()) {
            // No need to change the range, just update the min and max keys
            this.minKey = newMinKey;
            this.maxKey = newMaxKey;
        }
        else {
            // Grow the bins
            var newLength = this._getNewLength(newMinKey, newMaxKey);
            if (newLength > this.length()) {
                (_a = this.bins).push.apply(_a, new Array(newLength - this.length()).fill(0));
            }
            this._adjust(newMinKey, newMaxKey);
        }
    };
    /** Calculate the bin index for the key, extending the range if necessary */
    DenseStore.prototype._getIndex = function (key) {
        if (key < this.minKey) {
            this._extendRange(key);
        }
        else if (key > this.maxKey) {
            this._extendRange(key);
        }
        return key - this.offset;
    };
    DenseStore.prototype.toProto = function () {
        var ProtoStore = require('../proto/compiled').Store;
        return ProtoStore.create({
            contiguousBinCounts: this.bins,
            contiguousBinIndexOffset: this.offset
        });
    };
    DenseStore.fromProto = function (protoStore) {
        if (!protoStore ||
            /* Double equals (==) is intentional here to check for
             * `null` | `undefined` without including `0` */
            protoStore.contiguousBinCounts == null ||
            protoStore.contiguousBinIndexOffset == null) {
            throw Error('Failed to decode store from protobuf');
        }
        var store = new this();
        var index = protoStore.contiguousBinIndexOffset;
        store.offset = index;
        for (var _i = 0, _a = protoStore.contiguousBinCounts; _i < _a.length; _i++) {
            var count = _a[_i];
            store.add(index, count);
            index += 1;
        }
        return store;
    };
    return DenseStore;
}());
exports.DenseStore = DenseStore;
