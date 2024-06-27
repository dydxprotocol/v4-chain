"use strict";
/*
 * Unless explicitly stated otherwise all files in this repository are licensed
 * under the Apache 2.0 license (see LICENSE).
 * This product includes software developed at Datadog (https://www.datadoghq.com/).
 * Copyright 2020 Datadog, Inc.
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.CollapsingHighestDenseStore = exports.CollapsingLowestDenseStore = exports.DenseStore = void 0;
var DenseStore_1 = require("./DenseStore");
Object.defineProperty(exports, "DenseStore", { enumerable: true, get: function () { return DenseStore_1.DenseStore; } });
var CollapsingLowestDenseStore_1 = require("./CollapsingLowestDenseStore");
Object.defineProperty(exports, "CollapsingLowestDenseStore", { enumerable: true, get: function () { return CollapsingLowestDenseStore_1.CollapsingLowestDenseStore; } });
var CollapsingHighestDenseStore_1 = require("./CollapsingHighestDenseStore");
Object.defineProperty(exports, "CollapsingHighestDenseStore", { enumerable: true, get: function () { return CollapsingHighestDenseStore_1.CollapsingHighestDenseStore; } });
