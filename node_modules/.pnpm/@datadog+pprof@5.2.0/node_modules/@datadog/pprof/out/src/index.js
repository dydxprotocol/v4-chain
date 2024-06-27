"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.heap = exports.time = exports.getNativeThreadId = exports.setLogger = exports.SourceMapper = exports.encodeSync = exports.encode = void 0;
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
const fs_1 = require("fs");
const heapProfiler = __importStar(require("./heap-profiler"));
const profile_encoder_1 = require("./profile-encoder");
const timeProfiler = __importStar(require("./time-profiler"));
var profile_encoder_2 = require("./profile-encoder");
Object.defineProperty(exports, "encode", { enumerable: true, get: function () { return profile_encoder_2.encode; } });
Object.defineProperty(exports, "encodeSync", { enumerable: true, get: function () { return profile_encoder_2.encodeSync; } });
var sourcemapper_1 = require("./sourcemapper/sourcemapper");
Object.defineProperty(exports, "SourceMapper", { enumerable: true, get: function () { return sourcemapper_1.SourceMapper; } });
var logger_1 = require("./logger");
Object.defineProperty(exports, "setLogger", { enumerable: true, get: function () { return logger_1.setLogger; } });
var time_profiler_1 = require("./time-profiler");
Object.defineProperty(exports, "getNativeThreadId", { enumerable: true, get: function () { return time_profiler_1.getNativeThreadId; } });
exports.time = {
    profile: timeProfiler.profile,
    start: timeProfiler.start,
    stop: timeProfiler.stop,
    getContext: timeProfiler.getContext,
    setContext: timeProfiler.setContext,
    isStarted: timeProfiler.isStarted,
    v8ProfilerStuckEventLoopDetected: timeProfiler.v8ProfilerStuckEventLoopDetected,
    getState: timeProfiler.getState,
    constants: timeProfiler.constants,
};
exports.heap = {
    start: heapProfiler.start,
    stop: heapProfiler.stop,
    profile: heapProfiler.profile,
    convertProfile: heapProfiler.convertProfile,
    v8Profile: heapProfiler.v8Profile,
    monitorOutOfMemory: heapProfiler.monitorOutOfMemory,
    CallbackMode: heapProfiler.CallbackMode,
};
// If loaded with --require, start profiling.
if (module.parent && module.parent.id === 'internal/preload') {
    exports.time.start({});
    process.on('exit', () => {
        // The process is going to terminate imminently. All work here needs to
        // be synchronous.
        const profile = exports.time.stop();
        const buffer = (0, profile_encoder_1.encodeSync)(profile);
        (0, fs_1.writeFileSync)(`pprof-profile-${process.pid}.pb.gz`, buffer);
    });
}
//# sourceMappingURL=index.js.map