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
export interface TimeProfile {
    /** Time in nanoseconds at which profile was stopped. */
    endTime: number;
    topDownRoot: TimeProfileNode;
    /** Time in nanoseconds at which profile was started. */
    startTime: number;
    hasCpuTime?: boolean;
    /** CPU time of non-JS threads, only reported for the main worker thread */
    nonJSThreadsCpuTime?: number;
}
export interface ProfileNode {
    name?: string;
    scriptName: string;
    scriptId?: number;
    lineNumber?: number;
    columnNumber?: number;
    children: ProfileNode[];
}
export interface TimeProfileNodeContext {
    context: object;
    timestamp: bigint;
    cpuTime: number;
}
export interface TimeProfileNode extends ProfileNode {
    hitCount: number;
    contexts?: TimeProfileNodeContext[];
}
export interface AllocationProfileNode extends ProfileNode {
    allocations: Allocation[];
}
export interface Allocation {
    sizeBytes: number;
    count: number;
}
export interface LabelSet {
    [key: string]: string | number;
}
export interface GenerateAllocationLabelsFunction {
    ({ node }: {
        node: AllocationProfileNode;
    }): LabelSet;
}
export interface GenerateTimeLabelsArgs {
    node: TimeProfileNode;
    context?: TimeProfileNodeContext;
}
export interface GenerateTimeLabelsFunction {
    (args: GenerateTimeLabelsArgs): LabelSet;
}
