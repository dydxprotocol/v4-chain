/**
 * Copyright 2017 Google Inc. All Rights Reserved.
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
import { SourceMapper } from './sourcemapper/sourcemapper';
import { getNativeThreadId } from './time-profiler-bindings';
import { GenerateTimeLabelsFunction } from './v8-types';
type Microseconds = number;
type Milliseconds = number;
export interface TimeProfilerOptions {
    /** time in milliseconds for which to collect profile. */
    durationMillis?: Milliseconds;
    /** average time in microseconds between samples */
    intervalMicros?: Microseconds;
    sourceMapper?: SourceMapper;
    /**
     * This configuration option is experimental.
     * When set to true, functions will be aggregated at the line level, rather
     * than at the function level.
     * This defaults to false.
     */
    lineNumbers?: boolean;
    withContexts?: boolean;
    workaroundV8Bug?: boolean;
    collectCpuTime?: boolean;
}
export declare function profile(options?: TimeProfilerOptions): Promise<import("pprof-format").Profile>;
export declare function start(options?: TimeProfilerOptions): void;
export declare function stop(restart?: boolean, generateLabels?: GenerateTimeLabelsFunction): import("pprof-format").Profile;
export declare function getState(): any;
export declare function setContext(context?: object): void;
export declare function getContext(): any;
export declare function isStarted(): boolean;
export declare function v8ProfilerStuckEventLoopDetected(): number;
export declare const constants: {
    kSampleCount: any;
    NON_JS_THREADS_FUNCTION_NAME: string;
};
export { getNativeThreadId };
