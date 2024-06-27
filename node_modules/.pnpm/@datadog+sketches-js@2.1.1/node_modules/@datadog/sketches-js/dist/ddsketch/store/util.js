"use strict";
/*
 * Unless explicitly stated otherwise all files in this repository are licensed
 * under the Apache 2.0 license (see LICENSE).
 * This product includes software developed at Datadog (https://www.datadoghq.com/).
 * Copyright 2020 Datadog, Inc.
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.sumOfRange = void 0;
/**
 * Return the sum of the values from range `start` to `end` in `array`
 */
var sumOfRange = function (array, start, end) {
    var sum = 0;
    for (var i = start; i <= end; i++) {
        sum += array[i];
    }
    return sum;
};
exports.sumOfRange = sumOfRange;
