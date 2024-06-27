"use strict";
/**
 * Copyright 2019 Google Inc. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.encodeSync = exports.encode = void 0;
const util_1 = require("util");
const zlib_1 = require("zlib");
const gzipPromise = (0, util_1.promisify)(zlib_1.gzip);
function encode(profile) {
    return profile.encodeAsync().then(gzipPromise);
}
exports.encode = encode;
function encodeSync(profile) {
    return (0, zlib_1.gzipSync)(profile.encode());
}
exports.encodeSync = encodeSync;
//# sourceMappingURL=profile-encoder.js.map