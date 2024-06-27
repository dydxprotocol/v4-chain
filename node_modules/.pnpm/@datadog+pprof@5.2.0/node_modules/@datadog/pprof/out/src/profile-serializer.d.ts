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
import { AllocationProfileNode, GenerateAllocationLabelsFunction, GenerateTimeLabelsFunction, TimeProfile } from './v8-types';
export declare const NON_JS_THREADS_FUNCTION_NAME = "(non-JS threads)";
/**
 * Converts v8 time profile into into a profile proto.
 * (https://github.com/google/pprof/blob/master/proto/profile.proto)
 *
 * @param prof - profile to be converted.
 * @param intervalMicros - average time (microseconds) between samples.
 */
export declare function serializeTimeProfile(prof: TimeProfile, intervalMicros: number, sourceMapper?: SourceMapper, recomputeSamplingInterval?: boolean, generateLabels?: GenerateTimeLabelsFunction): Profile;
/**
 * Converts v8 heap profile into into a profile proto.
 * (https://github.com/google/pprof/blob/master/proto/profile.proto)
 *
 * @param prof - profile to be converted.
 * @param startTimeNanos - start time of profile, in nanoseconds (POSIX time).
 * @param durationsNanos - duration of the profile (wall clock time) in
 * nanoseconds.
 * @param intervalBytes - bytes allocated between samples.
 */
export declare function serializeHeapProfile(prof: AllocationProfileNode, startTimeNanos: number, intervalBytes: number, ignoreSamplesPath?: string, sourceMapper?: SourceMapper, generateLabels?: GenerateAllocationLabelsFunction): Profile;
