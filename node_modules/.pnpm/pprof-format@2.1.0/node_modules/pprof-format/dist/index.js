"use strict";
/**
 * Unless explicitly stated otherwise all files in this repository are licensed under the MIT License.
 *
 * This product includes software developed at Datadog (https://www.datadoghq.com/  Copyright 2022 Datadog, Inc.
 */
var __classPrivateFieldGet = (this && this.__classPrivateFieldGet) || function (receiver, state, kind, f) {
    if (kind === "a" && !f) throw new TypeError("Private accessor was defined without a getter");
    if (typeof state === "function" ? receiver !== state || !f : !state.has(receiver)) throw new TypeError("Cannot read private member from an object whose class did not declare it");
    return kind === "m" ? f : kind === "a" ? f.call(receiver) : f ? f.value : state.get(receiver);
};
var _StringTable_encodings, _StringTable_positions;
Object.defineProperty(exports, "__esModule", { value: true });
exports.Profile = exports.Function = exports.Location = exports.Line = exports.Mapping = exports.Sample = exports.Label = exports.ValueType = exports.StringTable = exports.emptyTableToken = void 0;
/*!
 * Private helpers. These are only used by other helpers.
 */
const lowMaxBig = 2n ** 32n - 1n;
const lowMax = 2 ** 32 - 1;
const lowMaxPlus1 = lowMax + 1;
// Buffer.from(string, 'utf8') is faster, when available
const toUtf8 = typeof Buffer === 'undefined'
    ? (value) => new TextEncoder().encode(value)
    : (value) => Buffer.from(value, 'utf8');
function countNumberBytes(buffer) {
    if (!buffer.length)
        return 0;
    let i = 0;
    while (i < buffer.length && buffer[i++] >= 0b10000000)
        ;
    return i;
}
function decodeBigNumber(buffer) {
    if (!buffer.length)
        return BigInt(0);
    let value = BigInt(buffer[0] & 0b01111111);
    let i = 0;
    while (buffer[i++] >= 0b10000000) {
        value |= BigInt(buffer[i] & 0b01111111) << BigInt(7 * i);
    }
    return value;
}
function makeValue(value, offset = 0) {
    return { value, offset };
}
function getValue(mode, buffer) {
    switch (mode) {
        case kTypeVarInt:
            for (let i = 0; i < buffer.length; i++) {
                if (!(buffer[i] & 0b10000000)) {
                    return makeValue(buffer.slice(0, i + 1));
                }
            }
            return makeValue(buffer);
        case kTypeLengthDelim: {
            const offset = countNumberBytes(buffer);
            const size = decodeNumber(buffer);
            return makeValue(buffer.slice(offset, Number(size) + offset), offset);
        }
        default:
            throw new Error(`Unrecognized value type: ${mode}`);
    }
}
function lowBits(number) {
    return typeof number !== 'bigint'
        ? (number >>> 0) % lowMaxPlus1
        : Number(number & lowMaxBig);
}
function highBits(number) {
    return typeof number !== 'bigint'
        ? (number / lowMaxPlus1) >>> 0
        : Number(number >> 32n & lowMaxBig);
}
function long(number) {
    const sign = number < 0;
    if (sign)
        number = -number;
    let lo = lowBits(number);
    let hi = highBits(number);
    if (sign) {
        hi = ~hi >>> 0;
        lo = ~lo >>> 0;
        if (++lo > lowMax) {
            lo = 0;
            if (++hi > lowMax) {
                hi = 0;
            }
        }
    }
    return [hi, lo];
}
/**
 * Public helpers. These are used in the type definitions.
 */
const kTypeVarInt = 0;
const kTypeLengthDelim = 2;
function decodeNumber(buffer) {
    const size = countNumberBytes(buffer);
    if (size > 4)
        return decodeBigNumber(buffer);
    if (!buffer.length)
        return 0;
    let value = buffer[0] & 0b01111111;
    let i = 0;
    while (buffer[i++] >= 0b10000000) {
        value |= (buffer[i] & 0b01111111) << (7 * i);
    }
    return value;
}
function decodeNumbers(buffer) {
    const values = [];
    let start = 0;
    for (let i = 0; i < buffer.length; i++) {
        if ((buffer[i] & 0b10000000) === 0) {
            values.push(decodeNumber(buffer.slice(start, i + 1)));
            start = i + 1;
        }
    }
    return values;
}
function push(value, list) {
    if (list == null) {
        return [value];
    }
    list.push(value);
    return list;
}
function measureNumber(number) {
    if (number === 0 || number === 0n)
        return 0;
    const [hi, lo] = long(number);
    const a = lo;
    const b = (lo >>> 28 | hi << 4) >>> 0;
    const c = hi >>> 24;
    if (c !== 0) {
        return c < 128 ? 9 : 10;
    }
    if (b !== 0) {
        if (b < 16384) {
            return b < 128 ? 5 : 6;
        }
        return b < 2097152 ? 7 : 8;
    }
    if (a < 16384) {
        return a < 128 ? 1 : 2;
    }
    return a < 2097152 ? 3 : 4;
}
function measureValue(value) {
    if (typeof value === 'undefined')
        return 0;
    if (typeof value === 'number' || typeof value === 'bigint') {
        return measureNumber(value) || 1;
    }
    return value.length;
}
function measureArray(list) {
    let size = 0;
    for (const item of list) {
        size += measureValue(item);
    }
    return size;
}
function measureNumberField(number) {
    const length = measureNumber(number);
    return length ? 1 + length : 0;
}
function measureNumberArrayField(values) {
    let total = 0;
    for (const value of values) {
        // Arrays should always include zeros to keep positions consistent
        total += measureNumber(value) || 1;
    }
    // Packed arrays are encoded as Tag,Len,ConcatenatedElements
    // Tag is only one byte because field number is always < 16 in pprof
    return total ? 1 + measureNumber(total) + total : 0;
}
function measureLengthDelimField(value) {
    const length = measureValue(value);
    // Length delimited records / submessages are encoded as Tag,Len,EncodedRecord
    // Tag is only one byte because field number is always < 16 in pprof
    return length ? 1 + measureNumber(length) + length : 0;
}
function measureLengthDelimArrayField(values) {
    let total = 0;
    for (const value of values) {
        total += measureLengthDelimField(value);
    }
    return total;
}
function encodeNumber(buffer, i, number) {
    if (number === 0 || number === 0n) {
        buffer[i++] = 0;
        return i;
    }
    let [hi, lo] = long(number);
    while (hi) {
        buffer[i++] = lo & 127 | 128;
        lo = (lo >>> 7 | hi << 25) >>> 0;
        hi >>>= 7;
    }
    while (lo > 127) {
        buffer[i++] = lo & 127 | 128;
        lo = lo >>> 7;
    }
    buffer[i++] = lo;
    return i;
}
exports.emptyTableToken = Symbol();
class StringTable {
    constructor(tok) {
        this.strings = new Array();
        _StringTable_encodings.set(this, new Array());
        _StringTable_positions.set(this, new Map());
        if (tok !== exports.emptyTableToken) {
            this.dedup('');
        }
    }
    get encodedLength() {
        let size = 0;
        for (const encoded of __classPrivateFieldGet(this, _StringTable_encodings, "f")) {
            size += encoded.length;
        }
        return size;
    }
    _encodeToBuffer(buffer, offset) {
        for (const encoded of __classPrivateFieldGet(this, _StringTable_encodings, "f")) {
            buffer.set(encoded, offset);
            offset += encoded.length;
        }
        return offset;
    }
    encode(buffer = new Uint8Array(this.encodedLength)) {
        this._encodeToBuffer(buffer, 0);
        return buffer;
    }
    static _encodeStringFromUtf8(stringBuffer) {
        const buffer = new Uint8Array(1 + stringBuffer.length + (measureNumber(stringBuffer.length) || 1));
        let offset = 0;
        buffer[offset++] = 50; // (6 << 3) + kTypeLengthDelim
        offset = encodeNumber(buffer, offset, stringBuffer.length);
        if (stringBuffer.length > 0) {
            buffer.set(stringBuffer, offset++);
        }
        return buffer;
    }
    static _encodeString(string) {
        return StringTable._encodeStringFromUtf8(toUtf8(string));
    }
    dedup(string) {
        if (typeof string === 'number')
            return string;
        if (!__classPrivateFieldGet(this, _StringTable_positions, "f").has(string)) {
            const pos = this.strings.push(string) - 1;
            __classPrivateFieldGet(this, _StringTable_positions, "f").set(string, pos);
            // Encode strings on insertion
            __classPrivateFieldGet(this, _StringTable_encodings, "f").push(StringTable._encodeString(string));
        }
        return __classPrivateFieldGet(this, _StringTable_positions, "f").get(string);
    }
    _decodeString(buffer) {
        const string = new TextDecoder().decode(buffer);
        __classPrivateFieldGet(this, _StringTable_positions, "f").set(string, this.strings.push(string) - 1);
        __classPrivateFieldGet(this, _StringTable_encodings, "f").push(StringTable._encodeStringFromUtf8(buffer));
    }
}
exports.StringTable = StringTable;
_StringTable_encodings = new WeakMap(), _StringTable_positions = new WeakMap();
function decode(buffer, decoder) {
    const data = {};
    let index = 0;
    while (index < buffer.length) {
        const field = buffer[index] >> 3;
        const mode = buffer[index] & 0b111;
        index++;
        const { offset, value } = getValue(mode, buffer.slice(index));
        index += value.length + offset;
        decoder(data, field, value);
    }
    return data;
}
class ValueType {
    constructor(data) {
        this.type = data.type || 0;
        this.unit = data.unit || 0;
    }
    get length() {
        let total = 0;
        total += measureNumberField(this.type);
        total += measureNumberField(this.unit);
        return total;
    }
    _encodeToBuffer(buffer, offset = 0) {
        if (this.type) {
            buffer[offset++] = 8; // (1 << 3) + kTypeVarInt
            offset = encodeNumber(buffer, offset, this.type);
        }
        if (this.unit) {
            buffer[offset++] = 16; // (2 << 3) + kTypeVarInt
            offset = encodeNumber(buffer, offset, this.unit);
        }
        return offset;
    }
    encode(buffer = new Uint8Array(this.length)) {
        this._encodeToBuffer(buffer, 0);
        return buffer;
    }
    static decodeValue(data, field, buffer) {
        switch (field) {
            case 1:
                data.type = decodeNumber(buffer);
                break;
            case 2:
                data.unit = decodeNumber(buffer);
                break;
        }
    }
    static decode(buffer) {
        return new this(decode(buffer, this.decodeValue));
    }
}
exports.ValueType = ValueType;
class Label {
    constructor(data) {
        this.key = data.key || 0;
        this.str = data.str || 0;
        this.num = data.num || 0;
        this.numUnit = data.numUnit || 0;
    }
    get length() {
        let total = 0;
        total += measureNumberField(this.key);
        total += measureNumberField(this.str);
        total += measureNumberField(this.num);
        total += measureNumberField(this.numUnit);
        return total;
    }
    _encodeToBuffer(buffer, offset = 0) {
        if (this.key) {
            buffer[offset++] = 8; // (1 << 3) + kTypeVarInt
            offset = encodeNumber(buffer, offset, this.key);
        }
        if (this.str) {
            buffer[offset++] = 16; // (2 << 3) + kTypeVarInt
            offset = encodeNumber(buffer, offset, this.str);
        }
        if (this.num) {
            buffer[offset++] = 24; // (3 << 3) + kTypeVarInt
            offset = encodeNumber(buffer, offset, this.num);
        }
        if (this.numUnit) {
            buffer[offset++] = 32; // (4 << 3) + kTypeVarInt
            offset = encodeNumber(buffer, offset, this.numUnit);
        }
        return offset;
    }
    encode(buffer = new Uint8Array(this.length)) {
        this._encodeToBuffer(buffer, 0);
        return buffer;
    }
    static decodeValue(data, field, buffer) {
        switch (field) {
            case 1:
                data.key = decodeNumber(buffer);
                break;
            case 2:
                data.str = decodeNumber(buffer);
                break;
            case 3:
                data.num = decodeNumber(buffer);
                break;
            case 4:
                data.numUnit = decodeNumber(buffer);
                break;
        }
    }
    static decode(buffer) {
        return new this(decode(buffer, this.decodeValue));
    }
}
exports.Label = Label;
class Sample {
    constructor(data) {
        this.locationId = data.locationId || [];
        this.value = data.value || [];
        this.label = (data.label || []).map(l => new Label(l));
    }
    get length() {
        let total = 0;
        total += measureNumberArrayField(this.locationId);
        total += measureNumberArrayField(this.value);
        total += measureLengthDelimArrayField(this.label);
        return total;
    }
    _encodeToBuffer(buffer, offset = 0) {
        if (this.locationId.length) {
            buffer[offset++] = 10; // (1 << 3) + kTypeLengthDelim
            offset = encodeNumber(buffer, offset, measureArray(this.locationId));
            for (const locationId of this.locationId) {
                offset = encodeNumber(buffer, offset, locationId);
            }
        }
        if (this.value.length) {
            buffer[offset++] = 18; // (2 << 3) + kTypeLengthDelim
            offset = encodeNumber(buffer, offset, measureArray(this.value));
            for (const value of this.value) {
                offset = encodeNumber(buffer, offset, value);
            }
        }
        for (const label of this.label) {
            buffer[offset++] = 26; // (3 << 3) + kTypeLengthDelim
            offset = encodeNumber(buffer, offset, label.length);
            offset = label._encodeToBuffer(buffer, offset);
        }
        return offset;
    }
    encode(buffer = new Uint8Array(this.length)) {
        this._encodeToBuffer(buffer, 0);
        return buffer;
    }
    static decodeValue(data, field, buffer) {
        switch (field) {
            case 1:
                data.locationId = decodeNumbers(buffer);
                break;
            case 2:
                data.value = decodeNumbers(buffer);
                break;
            case 3:
                data.label = push(Label.decode(buffer), data.label);
                break;
        }
    }
    static decode(buffer) {
        return new this(decode(buffer, this.decodeValue));
    }
}
exports.Sample = Sample;
class Mapping {
    constructor(data) {
        this.id = data.id || 0;
        this.memoryStart = data.memoryStart || 0;
        this.memoryLimit = data.memoryLimit || 0;
        this.fileOffset = data.fileOffset || 0;
        this.filename = data.filename || 0;
        this.buildId = data.buildId || 0;
        this.hasFunctions = !!data.hasFunctions;
        this.hasFilenames = !!data.hasFilenames;
        this.hasLineNumbers = !!data.hasLineNumbers;
        this.hasInlineFrames = !!data.hasInlineFrames;
    }
    get length() {
        let total = 0;
        total += measureNumberField(this.id);
        total += measureNumberField(this.memoryStart);
        total += measureNumberField(this.memoryLimit);
        total += measureNumberField(this.fileOffset);
        total += measureNumberField(this.filename);
        total += measureNumberField(this.buildId);
        total += measureNumberField(this.hasFunctions ? 1 : 0);
        total += measureNumberField(this.hasFilenames ? 1 : 0);
        total += measureNumberField(this.hasLineNumbers ? 1 : 0);
        total += measureNumberField(this.hasInlineFrames ? 1 : 0);
        return total;
    }
    _encodeToBuffer(buffer, offset = 0) {
        if (this.id) {
            buffer[offset++] = 8; // (1 << 3) + kTypeVarInt
            offset = encodeNumber(buffer, offset, this.id);
        }
        if (this.memoryStart) {
            buffer[offset++] = 16; // (2 << 3) + kTypeVarInt
            offset = encodeNumber(buffer, offset, this.memoryStart);
        }
        if (this.memoryLimit) {
            buffer[offset++] = 24; // (3 << 3) + kTypeVarInt
            offset = encodeNumber(buffer, offset, this.memoryLimit);
        }
        if (this.fileOffset) {
            buffer[offset++] = 32; // (4 << 3) + kTypeVarInt
            offset = encodeNumber(buffer, offset, this.fileOffset);
        }
        if (this.filename) {
            buffer[offset++] = 40; // (5 << 3) + kTypeVarInt
            offset = encodeNumber(buffer, offset, this.filename);
        }
        if (this.buildId) {
            buffer[offset++] = 48; // (6 << 3) + kTypeVarInt
            offset = encodeNumber(buffer, offset, this.buildId);
        }
        if (this.hasFunctions) {
            buffer[offset++] = 56; // (7 << 3) + kTypeVarInt
            offset = encodeNumber(buffer, offset, 1);
        }
        if (this.hasFilenames) {
            buffer[offset++] = 64; // (8 << 3) + kTypeVarInt
            offset = encodeNumber(buffer, offset, 1);
        }
        if (this.hasLineNumbers) {
            buffer[offset++] = 72; // (9 << 3) + kTypeVarInt
            offset = encodeNumber(buffer, offset, 1);
        }
        if (this.hasInlineFrames) {
            buffer[offset++] = 80; // (10 << 3) + kTypeVarInt
            offset = encodeNumber(buffer, offset, 1);
        }
        return offset;
    }
    encode(buffer = new Uint8Array(this.length)) {
        this._encodeToBuffer(buffer, 0);
        return buffer;
    }
    static decodeValue(data, field, buffer) {
        switch (field) {
            case 1:
                data.id = decodeNumber(buffer);
                break;
            case 2:
                data.memoryStart = decodeNumber(buffer);
                break;
            case 3:
                data.memoryLimit = decodeNumber(buffer);
                break;
            case 4:
                data.fileOffset = decodeNumber(buffer);
                break;
            case 5:
                data.filename = decodeNumber(buffer);
                break;
            case 6:
                data.buildId = decodeNumber(buffer);
                break;
            case 7:
                data.hasFunctions = !!decodeNumber(buffer);
                break;
            case 8:
                data.hasFilenames = !!decodeNumber(buffer);
                break;
            case 9:
                data.hasLineNumbers = !!decodeNumber(buffer);
                break;
            case 10:
                data.hasInlineFrames = !!decodeNumber(buffer);
                break;
        }
    }
    static decode(buffer) {
        return new this(decode(buffer, this.decodeValue));
    }
}
exports.Mapping = Mapping;
class Line {
    constructor(data) {
        this.functionId = data.functionId || 0;
        this.line = data.line || 0;
    }
    get length() {
        let total = 0;
        total += measureNumberField(this.functionId);
        total += measureNumberField(this.line);
        return total;
    }
    _encodeToBuffer(buffer, offset = 0) {
        if (this.functionId) {
            buffer[offset++] = 8; // (1 << 3) + kTypeVarInt
            offset = encodeNumber(buffer, offset, this.functionId);
        }
        if (this.line) {
            buffer[offset++] = 16; // (2 << 3) + kTypeVarInt
            offset = encodeNumber(buffer, offset, this.line);
        }
        return offset;
    }
    encode(buffer = new Uint8Array(this.length)) {
        this._encodeToBuffer(buffer, 0);
        return buffer;
    }
    static decodeValue(data, field, buffer) {
        switch (field) {
            case 1:
                data.functionId = decodeNumber(buffer);
                break;
            case 2:
                data.line = decodeNumber(buffer);
                break;
        }
    }
    static decode(buffer) {
        return new this(decode(buffer, this.decodeValue));
    }
}
exports.Line = Line;
class Location {
    constructor(data) {
        this.id = data.id || 0;
        this.mappingId = data.mappingId || 0;
        this.address = data.address || 0;
        this.line = (data.line || []).map(l => new Line(l));
        this.isFolded = !!data.isFolded;
    }
    get length() {
        let total = 0;
        total += measureNumberField(this.id);
        total += measureNumberField(this.mappingId);
        total += measureNumberField(this.address);
        total += measureLengthDelimArrayField(this.line);
        total += measureNumberField(this.isFolded ? 1 : 0);
        return total;
    }
    _encodeToBuffer(buffer, offset = 0) {
        if (this.id) {
            buffer[offset++] = 8; // (1 << 3) + kTypeVarInt
            offset = encodeNumber(buffer, offset, this.id);
        }
        if (this.mappingId) {
            buffer[offset++] = 16; // (2 << 3) + kTypeVarInt
            offset = encodeNumber(buffer, offset, this.mappingId);
        }
        if (this.address) {
            buffer[offset++] = 24; // (3 << 3) + kTypeVarInt
            offset = encodeNumber(buffer, offset, this.address);
        }
        for (const line of this.line) {
            buffer[offset++] = 34; // (4 << 3) + kTypeLengthDelim
            offset = encodeNumber(buffer, offset, line.length);
            offset = line._encodeToBuffer(buffer, offset);
        }
        if (this.isFolded) {
            buffer[offset++] = 40; // (5 << 3) + kTypeVarInt
            offset = encodeNumber(buffer, offset, 1);
        }
        return offset;
    }
    encode(buffer = new Uint8Array(this.length)) {
        this._encodeToBuffer(buffer, 0);
        return buffer;
    }
    static decodeValue(data, field, buffer) {
        switch (field) {
            case 1:
                data.id = decodeNumber(buffer);
                break;
            case 2:
                data.mappingId = decodeNumber(buffer);
                break;
            case 3:
                data.address = decodeNumber(buffer);
                break;
            case 4:
                data.line = push(Line.decode(buffer), data.line);
                break;
            case 5:
                data.isFolded = !!decodeNumber(buffer);
                break;
        }
    }
    static decode(buffer) {
        return new this(decode(buffer, this.decodeValue));
    }
}
exports.Location = Location;
class Function {
    constructor(data) {
        this.id = data.id || 0;
        this.name = data.name || 0;
        this.systemName = data.systemName || 0;
        this.filename = data.filename || 0;
        this.startLine = data.startLine || 0;
    }
    get length() {
        let total = 0;
        total += measureNumberField(this.id);
        total += measureNumberField(this.name);
        total += measureNumberField(this.systemName);
        total += measureNumberField(this.filename);
        total += measureNumberField(this.startLine);
        return total;
    }
    _encodeToBuffer(buffer, offset = 0) {
        if (this.id) {
            buffer[offset++] = 8; // (1 << 3) + kTypeVarInt
            offset = encodeNumber(buffer, offset, this.id);
        }
        if (this.name) {
            buffer[offset++] = 16; // (2 << 3) + kTypeVarInt
            offset = encodeNumber(buffer, offset, this.name);
        }
        if (this.systemName) {
            buffer[offset++] = 24; // (3 << 3) + kTypeVarInt
            offset = encodeNumber(buffer, offset, this.systemName);
        }
        if (this.filename) {
            buffer[offset++] = 32; // (4 << 3) + kTypeVarInt
            offset = encodeNumber(buffer, offset, this.filename);
        }
        if (this.startLine) {
            buffer[offset++] = 40; // (5 << 3) + kTypeVarInt
            offset = encodeNumber(buffer, offset, this.startLine);
        }
        return offset;
    }
    encode(buffer = new Uint8Array(this.length)) {
        this._encodeToBuffer(buffer, 0);
        return buffer;
    }
    static decodeValue(data, field, buffer) {
        switch (field) {
            case 1:
                data.id = decodeNumber(buffer);
                break;
            case 2:
                data.name = decodeNumber(buffer);
                break;
            case 3:
                data.systemName = decodeNumber(buffer);
                break;
            case 4:
                data.filename = decodeNumber(buffer);
                break;
            case 5:
                data.startLine = decodeNumber(buffer);
                break;
        }
    }
    static decode(buffer) {
        return new this(decode(buffer, this.decodeValue));
    }
}
exports.Function = Function;
class Profile {
    constructor(data = {}) {
        this.sampleType = (data.sampleType || []).map(v => new ValueType(v));
        this.sample = (data.sample || []).map(v => new Sample(v));
        this.mapping = (data.mapping || []).map(v => new Mapping(v));
        this.location = (data.location || []).map(v => new Location(v));
        this.function = (data.function || []).map(v => new Function(v));
        this.stringTable = data.stringTable || new StringTable();
        this.dropFrames = data.dropFrames || 0;
        this.keepFrames = data.keepFrames || 0;
        this.timeNanos = data.timeNanos || 0;
        this.durationNanos = data.durationNanos || 0;
        this.periodType = data.periodType ? new ValueType(data.periodType) : undefined;
        this.period = data.period || 0;
        this.comment = data.comment || [];
        this.defaultSampleType = data.defaultSampleType || 0;
    }
    get length() {
        let total = 0;
        total += measureLengthDelimArrayField(this.sampleType);
        total += measureLengthDelimArrayField(this.sample);
        total += measureLengthDelimArrayField(this.mapping);
        total += measureLengthDelimArrayField(this.location);
        total += measureLengthDelimArrayField(this.function);
        total += this.stringTable.encodedLength;
        total += measureNumberField(this.dropFrames);
        total += measureNumberField(this.keepFrames);
        total += measureNumberField(this.timeNanos);
        total += measureNumberField(this.durationNanos);
        total += measureLengthDelimField(this.periodType);
        total += measureNumberField(this.period);
        total += measureNumberArrayField(this.comment);
        total += measureNumberField(this.defaultSampleType);
        return total;
    }
    _encodeSampleTypesToBuffer(buffer, offset = 0) {
        for (const sampleType of this.sampleType) {
            buffer[offset++] = 10; // (1 << 3) + kTypeLengthDelim
            offset = encodeNumber(buffer, offset, sampleType.length);
            offset = sampleType._encodeToBuffer(buffer, offset);
        }
        return offset;
    }
    _encodeSamplesToBuffer(buffer, offset = 0) {
        for (const sample of this.sample) {
            buffer[offset++] = 18; // (2 << 3) + kTypeLengthDelim
            offset = encodeNumber(buffer, offset, sample.length);
            offset = sample._encodeToBuffer(buffer, offset);
        }
        return offset;
    }
    _encodeMappingsToBuffer(buffer, offset = 0) {
        for (const mapping of this.mapping) {
            buffer[offset++] = 26; // (3 << 3) + kTypeLengthDelim
            offset = encodeNumber(buffer, offset, mapping.length);
            offset = mapping._encodeToBuffer(buffer, offset);
        }
        return offset;
    }
    _encodeLocationsToBuffer(buffer, offset = 0) {
        for (const location of this.location) {
            buffer[offset++] = 34; // (4 << 3) + kTypeLengthDelim
            offset = encodeNumber(buffer, offset, location.length);
            offset = location._encodeToBuffer(buffer, offset);
        }
        return offset;
    }
    _encodeFunctionsToBuffer(buffer, offset = 0) {
        for (const fun of this.function) {
            buffer[offset++] = 42; // (5 << 3) + kTypeLengthDelim
            offset = encodeNumber(buffer, offset, fun.length);
            offset = fun._encodeToBuffer(buffer, offset);
        }
        return offset;
    }
    _encodeBasicValuesToBuffer(buffer, offset = 0) {
        if (this.dropFrames) {
            buffer[offset++] = 56; // (7 << 3) + kTypeVarInt
            offset = encodeNumber(buffer, offset, this.dropFrames);
        }
        if (this.keepFrames) {
            buffer[offset++] = 64; // (8 << 3) + kTypeVarInt
            offset = encodeNumber(buffer, offset, this.keepFrames);
        }
        if (this.timeNanos) {
            buffer[offset++] = 72; // (9 << 3) + kTypeVarInt
            offset = encodeNumber(buffer, offset, this.timeNanos);
        }
        if (this.durationNanos) {
            buffer[offset++] = 80; // (10 << 3) + kTypeVarInt
            offset = encodeNumber(buffer, offset, this.durationNanos);
        }
        if (typeof this.periodType !== 'undefined') {
            buffer[offset++] = 90; // (11 << 3) + kTypeLengthDelim
            offset = encodeNumber(buffer, offset, this.periodType.length);
            offset = this.periodType._encodeToBuffer(buffer, offset);
        }
        if (this.period) {
            buffer[offset++] = 96; // (12 << 3) + kTypeVarInt
            offset = encodeNumber(buffer, offset, this.period);
        }
        if (this.comment.length) {
            buffer[offset++] = 106; // (13 << 3) + kTypeLengthDelim
            offset = encodeNumber(buffer, offset, measureArray(this.comment));
            for (const comment of this.comment) {
                offset = encodeNumber(buffer, offset, comment);
            }
        }
        if (this.defaultSampleType) {
            buffer[offset++] = 112; // (14 << 3) + kTypeVarInt
            offset = encodeNumber(buffer, offset, this.defaultSampleType);
        }
        return offset;
    }
    _encodeToBuffer(buffer, offset = 0) {
        offset = this._encodeSampleTypesToBuffer(buffer, offset);
        offset = this._encodeSamplesToBuffer(buffer, offset);
        offset = this._encodeMappingsToBuffer(buffer, offset);
        offset = this._encodeLocationsToBuffer(buffer, offset);
        offset = this._encodeFunctionsToBuffer(buffer, offset);
        offset = this.stringTable._encodeToBuffer(buffer, offset);
        offset = this._encodeBasicValuesToBuffer(buffer, offset);
        return offset;
    }
    async _encodeToBufferAsync(buffer, offset = 0) {
        offset = this._encodeSampleTypesToBuffer(buffer, offset);
        await new Promise(setImmediate);
        offset = this._encodeSamplesToBuffer(buffer, offset);
        await new Promise(setImmediate);
        offset = this._encodeMappingsToBuffer(buffer, offset);
        await new Promise(setImmediate);
        offset = this._encodeLocationsToBuffer(buffer, offset);
        await new Promise(setImmediate);
        offset = this._encodeFunctionsToBuffer(buffer, offset);
        await new Promise(setImmediate);
        offset = this.stringTable._encodeToBuffer(buffer, offset);
        await new Promise(setImmediate);
        offset = this._encodeBasicValuesToBuffer(buffer, offset);
        return offset;
    }
    encode(buffer = new Uint8Array(this.length)) {
        this._encodeToBuffer(buffer, 0);
        return buffer;
    }
    async encodeAsync(buffer = new Uint8Array(this.length)) {
        await this._encodeToBufferAsync(buffer, 0);
        return buffer;
    }
    static decodeValue(data, field, buffer) {
        switch (field) {
            case 1:
                data.sampleType = push(ValueType.decode(buffer), data.sampleType);
                break;
            case 2:
                data.sample = push(Sample.decode(buffer), data.sample);
                break;
            case 3:
                data.mapping = push(Mapping.decode(buffer), data.mapping);
                break;
            case 4:
                data.location = push(Location.decode(buffer), data.location);
                break;
            case 5:
                data.function = push(Function.decode(buffer), data.function);
                break;
            case 6: {
                if (data.stringTable === undefined) {
                    data.stringTable = new StringTable(exports.emptyTableToken);
                }
                data.stringTable._decodeString(buffer);
                break;
            }
            case 7:
                data.dropFrames = decodeNumber(buffer);
                break;
            case 8:
                data.keepFrames = decodeNumber(buffer);
                break;
            case 9:
                data.timeNanos = decodeNumber(buffer);
                break;
            case 10:
                data.durationNanos = decodeNumber(buffer);
                break;
            case 11:
                data.periodType = ValueType.decode(buffer);
                break;
            case 12:
                data.period = decodeNumber(buffer);
                break;
            case 13:
                data.comment = decodeNumbers(buffer);
                break;
            case 14:
                data.defaultSampleType = decodeNumber(buffer);
                break;
        }
    }
    static decode(buffer) {
        return new this(decode(buffer, this.decodeValue));
    }
}
exports.Profile = Profile;
//# sourceMappingURL=index.js.map