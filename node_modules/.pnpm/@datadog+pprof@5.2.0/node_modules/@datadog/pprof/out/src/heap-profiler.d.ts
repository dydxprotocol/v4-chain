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
import { Profile } from 'pprof-format';
import { SourceMapper } from './sourcemapper/sourcemapper';
import { AllocationProfileNode, GenerateAllocationLabelsFunction } from './v8-types';
export declare function v8Profile(): AllocationProfileNode;
/**
 * Collects a profile and returns it serialized in pprof format.
 * Throws if heap profiler is not enabled.
 *
 * @param ignoreSamplePath
 * @param sourceMapper
 */
export declare function profile(ignoreSamplePath?: string, sourceMapper?: SourceMapper, generateLabels?: GenerateAllocationLabelsFunction): Profile;
export declare function convertProfile(rootNode: AllocationProfileNode, ignoreSamplePath?: string, sourceMapper?: SourceMapper, generateLabels?: GenerateAllocationLabelsFunction): Profile;
/**
 * Starts heap profiling. If heap profiling has already been started with
 * the same parameters, this is a noop. If heap profiler has already been
 * started with different parameters, this throws an error.
 *
 * @param intervalBytes - average number of bytes between samples.
 * @param stackDepth - maximum stack depth for samples collected.
 */
export declare function start(intervalBytes: number, stackDepth: number): void;
export declare function stop(): void;
export type NearHeapLimitCallback = (profile: Profile) => void;
export declare const CallbackMode: {
    Async: number;
    Interrupt: number;
    Both: number;
};
/**
 * Add monitoring for v8 heap, heap profiler must already be started.
 * When an out of heap memory event occurs:
 *  - an extension of heap memory of |heapLimitExtensionSize| bytes is
 *    requested to v8. This extension can occur |maxHeapLimitExtensionCount|
 *    number of times. If the extension amount is not enough to satisfy
 *    memory allocation that triggers GC and OOM, process will abort.
 *  - heap profile is dumped as folded stacks on stderr if
 *    |dumpHeapProfileOnSdterr| is true
 *  - heap profile is dumped in temporary file and a new process is spawned
 *    with |exportCommand| arguments and profile path appended at the end.
 *  - |callback| is called. Callback can be invoked only if
 *    heapLimitExtensionSize is enough for the process to continue. Invocation
 *    will be done by a RequestInterrupt if |callbackMode| is Interrupt or Both,
 *    this might be unsafe since Isolate should not be reentered
 *    from RequestInterrupt, but this allows to interrupt synchronous code.
 *    Otherwise the callback is scheduled to be called asynchronously.
 * @param heapLimitExtensionSize - amount of bytes heap should be expanded
 *  with upon OOM
 * @param maxHeapLimitExtensionCount - maximum number of times heap size
 *  extension can occur
 * @param dumpHeapProfileOnSdterr - dump heap profile on stderr upon OOM
 * @param exportCommand - command to execute upon OOM, filepath of a
 *  temporary file containing heap profile will be appended
 * @param callback - callback to call when OOM occurs
 * @param callbackMode
 */
export declare function monitorOutOfMemory(heapLimitExtensionSize: number, maxHeapLimitExtensionCount: number, dumpHeapProfileOnSdterr: boolean, exportCommand?: Array<String>, callback?: NearHeapLimitCallback, callbackMode?: number): void;
