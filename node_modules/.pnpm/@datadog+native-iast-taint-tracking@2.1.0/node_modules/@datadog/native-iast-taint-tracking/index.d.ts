/**
* Unless explicitly stated otherwise all files in this repository are licensed under the Apache-2.0 License.
* This product includes software developed at Datadog (https://www.datadoghq.com/). Copyright 2022 Datadog, Inc.
**/

declare module 'datadog-iast-taint-tracking' {

    export interface NativeInputInfo {
        parameterName: string;
        parameterValue: string;
        type: string;
        readonly ref?: string;
    }

    export interface NativeTaintedRange {
        start: number;
        end: number;
        iinfo: NativeInputInfo;
        secureMarks: number;
        readonly ref?: string;
    }

    export interface Metrics {
        requestCount: number;
    }

    export interface TaintedUtils {
        createTransaction(transactionId: string): string;
        newTaintedString(transactionId: string, original: string, paramName: string, type: string): string;
        newTaintedObject(transactionId: string, original: any, paramName: string, type: string): any;
        addSecureMarksToTaintedString(transactionId: string, taintedString: string, secureMarks: number): string;
        isTainted(transactionId: string, ...args: string[]): boolean;
        getMetrics(transactionId: string, telemetryVerbosity: number): Metrics;
        getRanges(transactionId: string, original: string): NativeTaintedRange[];
        removeTransaction(transactionId: string): void;
        setMaxTransactions(maxTransactions: number): void;
        concat(transactionId: string, result: string, op1: string, op2: string): string;
        trim(transactionId: string, result: string, thisArg: string): string;
        trimEnd(transactionId: string, result: string, thisArg: string): string;
        slice(transactionId: string, result: string, original: string, start: number, end: number): string;
        substring(transactionId: string, subject: string, result: string, start: number, end: number): string;
        substr(transactionId: string, subject: string, result: string, start: number, length: number): string;
        replace(transactionId: string, result: string, thisArg: string, matcher: unknown, replacer: unknown): string;
        stringCase(transactionId: string, result: string, thisArg: string): string;
        arrayJoin(transactionId: string, result: string, thisArg: any[], separator?: any): string;
    }
}
