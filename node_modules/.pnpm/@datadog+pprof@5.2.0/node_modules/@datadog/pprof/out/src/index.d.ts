import * as heapProfiler from './heap-profiler';
import * as timeProfiler from './time-profiler';
export { AllocationProfileNode, TimeProfileNode, ProfileNode, LabelSet, } from './v8-types';
export { encode, encodeSync } from './profile-encoder';
export { SourceMapper } from './sourcemapper/sourcemapper';
export { setLogger } from './logger';
export { getNativeThreadId } from './time-profiler';
export declare const time: {
    profile: typeof timeProfiler.profile;
    start: typeof timeProfiler.start;
    stop: typeof timeProfiler.stop;
    getContext: typeof timeProfiler.getContext;
    setContext: typeof timeProfiler.setContext;
    isStarted: typeof timeProfiler.isStarted;
    v8ProfilerStuckEventLoopDetected: typeof timeProfiler.v8ProfilerStuckEventLoopDetected;
    getState: typeof timeProfiler.getState;
    constants: {
        kSampleCount: any;
        NON_JS_THREADS_FUNCTION_NAME: string;
    };
};
export declare const heap: {
    start: typeof heapProfiler.start;
    stop: typeof heapProfiler.stop;
    profile: typeof heapProfiler.profile;
    convertProfile: typeof heapProfiler.convertProfile;
    v8Profile: typeof heapProfiler.v8Profile;
    monitorOutOfMemory: typeof heapProfiler.monitorOutOfMemory;
    CallbackMode: {
        Async: number;
        Interrupt: number;
        Both: number;
    };
};
