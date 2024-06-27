/**
 * Unless explicitly stated otherwise all files in this repository are licensed under the MIT License.
 *
 * This product includes software developed at Datadog (https://www.datadoghq.com/  Copyright 2022 Datadog, Inc.
 */
/// <reference types="node" />
type Numeric = number | bigint;
export declare const emptyTableToken: unique symbol;
export declare class StringTable {
    #private;
    strings: string[];
    constructor(tok?: typeof emptyTableToken);
    get encodedLength(): number;
    _encodeToBuffer(buffer: Uint8Array, offset: number): number;
    encode(buffer?: Uint8Array): Uint8Array;
    static _encodeStringFromUtf8(stringBuffer: Uint8Array | Buffer): Uint8Array;
    static _encodeString(string: string): Uint8Array;
    dedup(string: string): number;
    _decodeString(buffer: Uint8Array): void;
}
export type ValueTypeInput = {
    type?: Numeric;
    unit?: Numeric;
};
export declare class ValueType {
    type: Numeric;
    unit: Numeric;
    constructor(data: ValueTypeInput);
    get length(): number;
    _encodeToBuffer(buffer: Uint8Array, offset?: number): number;
    encode(buffer?: Uint8Array): Uint8Array;
    static decodeValue(data: ValueTypeInput, field: number, buffer: Uint8Array): void;
    static decode(buffer: Uint8Array): ValueType;
}
export type LabelInput = {
    key?: Numeric;
    str?: Numeric;
    num?: Numeric;
    numUnit?: Numeric;
};
export declare class Label {
    key: Numeric;
    str: Numeric;
    num: Numeric;
    numUnit: Numeric;
    constructor(data: LabelInput);
    get length(): number;
    _encodeToBuffer(buffer: Uint8Array, offset?: number): number;
    encode(buffer?: Uint8Array): Uint8Array;
    static decodeValue(data: LabelInput, field: number, buffer: Uint8Array): void;
    static decode(buffer: Uint8Array): Label;
}
export type SampleInput = {
    locationId?: Array<Numeric>;
    value?: Array<Numeric>;
    label?: Array<LabelInput>;
};
export declare class Sample {
    locationId: Array<Numeric>;
    value: Array<Numeric>;
    label: Array<Label>;
    constructor(data: SampleInput);
    get length(): number;
    _encodeToBuffer(buffer: Uint8Array, offset?: number): number;
    encode(buffer?: Uint8Array): Uint8Array;
    static decodeValue(data: SampleInput, field: number, buffer: Uint8Array): void;
    static decode(buffer: Uint8Array): Sample;
}
export type MappingInput = {
    id?: Numeric;
    memoryStart?: Numeric;
    memoryLimit?: Numeric;
    fileOffset?: Numeric;
    filename?: Numeric;
    buildId?: Numeric;
    hasFunctions?: boolean;
    hasFilenames?: boolean;
    hasLineNumbers?: boolean;
    hasInlineFrames?: boolean;
};
export declare class Mapping {
    id: Numeric;
    memoryStart: Numeric;
    memoryLimit: Numeric;
    fileOffset: Numeric;
    filename: Numeric;
    buildId: Numeric;
    hasFunctions: boolean;
    hasFilenames: boolean;
    hasLineNumbers: boolean;
    hasInlineFrames: boolean;
    constructor(data: MappingInput);
    get length(): number;
    _encodeToBuffer(buffer: Uint8Array, offset?: number): number;
    encode(buffer?: Uint8Array): Uint8Array;
    static decodeValue(data: MappingInput, field: number, buffer: Uint8Array): void;
    static decode(buffer: Uint8Array): Mapping;
}
export type LineInput = {
    functionId?: Numeric;
    line?: Numeric;
};
export declare class Line {
    functionId: Numeric;
    line: Numeric;
    constructor(data: LineInput);
    get length(): number;
    _encodeToBuffer(buffer: Uint8Array, offset?: number): number;
    encode(buffer?: Uint8Array): Uint8Array;
    static decodeValue(data: LineInput, field: number, buffer: Uint8Array): void;
    static decode(buffer: Uint8Array): Line;
}
export type LocationInput = {
    id?: Numeric;
    mappingId?: Numeric;
    address?: Numeric;
    line?: Array<LineInput>;
    isFolded?: boolean;
};
export declare class Location {
    id: Numeric;
    mappingId: Numeric;
    address: Numeric;
    line: Array<Line>;
    isFolded: boolean;
    constructor(data: LocationInput);
    get length(): number;
    _encodeToBuffer(buffer: Uint8Array, offset?: number): number;
    encode(buffer?: Uint8Array): Uint8Array;
    static decodeValue(data: LocationInput, field: number, buffer: Uint8Array): void;
    static decode(buffer: Uint8Array): Location;
}
export type FunctionInput = {
    id?: Numeric;
    name?: Numeric;
    systemName?: Numeric;
    filename?: Numeric;
    startLine?: Numeric;
};
export declare class Function {
    id: Numeric;
    name: Numeric;
    systemName: Numeric;
    filename: Numeric;
    startLine: Numeric;
    constructor(data: FunctionInput);
    get length(): number;
    _encodeToBuffer(buffer: Uint8Array, offset?: number): number;
    encode(buffer?: Uint8Array): Uint8Array;
    static decodeValue(data: FunctionInput, field: number, buffer: Uint8Array): void;
    static decode(buffer: Uint8Array): Function;
}
export type ProfileInput = {
    sampleType?: Array<ValueTypeInput>;
    sample?: Array<SampleInput>;
    mapping?: Array<MappingInput>;
    location?: Array<LocationInput>;
    function?: Array<FunctionInput>;
    stringTable?: StringTable;
    dropFrames?: Numeric;
    keepFrames?: Numeric;
    timeNanos?: Numeric;
    durationNanos?: Numeric;
    periodType?: ValueTypeInput;
    period?: Numeric;
    comment?: Array<Numeric>;
    defaultSampleType?: Numeric;
};
export declare class Profile {
    sampleType: Array<ValueType>;
    sample: Array<Sample>;
    mapping: Array<Mapping>;
    location: Array<Location>;
    function: Array<Function>;
    stringTable: StringTable;
    dropFrames: Numeric;
    keepFrames: Numeric;
    timeNanos: Numeric;
    durationNanos: Numeric;
    periodType?: ValueType;
    period: Numeric;
    comment: Array<Numeric>;
    defaultSampleType: Numeric;
    constructor(data?: ProfileInput);
    get length(): number;
    _encodeSampleTypesToBuffer(buffer: Uint8Array, offset?: number): number;
    _encodeSamplesToBuffer(buffer: Uint8Array, offset?: number): number;
    _encodeMappingsToBuffer(buffer: Uint8Array, offset?: number): number;
    _encodeLocationsToBuffer(buffer: Uint8Array, offset?: number): number;
    _encodeFunctionsToBuffer(buffer: Uint8Array, offset?: number): number;
    _encodeBasicValuesToBuffer(buffer: Uint8Array, offset?: number): number;
    _encodeToBuffer(buffer: Uint8Array, offset?: number): number;
    _encodeToBufferAsync(buffer: Uint8Array, offset?: number): Promise<number>;
    encode(buffer?: Uint8Array): Uint8Array;
    encodeAsync(buffer?: Uint8Array): Promise<Uint8Array>;
    static decodeValue(data: ProfileInput, field: number, buffer: Uint8Array): void;
    static decode(buffer: Uint8Array): Profile;
}
export {};
