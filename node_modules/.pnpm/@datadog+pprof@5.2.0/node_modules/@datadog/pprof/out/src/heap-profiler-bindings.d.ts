/**
 * Copyright 2018 Google Inc. All Rights Reserved.
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
import { AllocationProfileNode } from './v8-types';
export declare function startSamplingHeapProfiler(heapIntervalBytes: number, heapStackDepth: number): void;
export declare function stopSamplingHeapProfiler(): void;
export declare function getAllocationProfile(): AllocationProfileNode;
export type NearHeapLimitCallback = (profile: AllocationProfileNode) => void;
export declare function monitorOutOfMemory(heapLimitExtensionSize: number, maxHeapLimitExtensionCount: number, dumpHeapProfileOnSdterr: boolean, exportCommand: Array<String> | undefined, callback: NearHeapLimitCallback | undefined, callbackMode: number, isMainThread: boolean): void;
