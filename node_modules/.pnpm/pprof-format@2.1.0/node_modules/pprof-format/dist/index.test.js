"use strict";
/**
 * Unless explicitly stated otherwise all files in this repository are licensed under the MIT License.
 *
 * This product includes software developed at Datadog (https://www.datadoghq.com/  Copyright 2022 Datadog, Inc.
 */
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
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const tap_1 = __importDefault(require("tap"));
const profile_1 = require("../testing/proto/profile");
const zlib_1 = require("zlib");
const fs = __importStar(require("fs"));
const { decode, toObject } = profile_1.perftools.profiles.Profile;
const index_js_1 = require("./index.js");
tap_1.default.Test.prototype.addAssert('constructs', 3, function (Type, data, encodings, message) {
    message = message || 'construction';
    return this.test(message, (t) => {
        const value = new Type(data);
        for (const { field } of encodings) {
            if (typeof data[field] === 'object') {
                t.has(value[field], data[field], `has given ${field}`);
            }
            else {
                t.equal(value[field], data[field], `has given ${field}`);
            }
        }
        t.end();
    });
});
tap_1.default.Test.prototype.addAssert('encodes', 3, function (Type, data, encodings, message) {
    message = message || 'encoding';
    return this.test(message, (t) => {
        t.test('per-field validation', (t2) => {
            for (const { field, value } of encodings) {
                const fun = new Type({
                    // Hack to exclude stringTable data from any checks except for the string table itself
                    stringTable: new index_js_1.StringTable(index_js_1.emptyTableToken),
                    [field]: data[field]
                });
                const msg = `has expected encoding of ${field} field`;
                t2.equal(bufToHex(fun.encode()), value, msg);
            }
            t2.end();
        });
        t.test('full object validation', (t2) => {
            const fun = new Type(data);
            t2.equal(bufToHex(fun.encode()), fullEncoding(encodings), 'has expected encoding of full object');
            t2.end();
        });
        t.end();
    });
});
tap_1.default.Test.prototype.addAssert('decodes', 3, function (Type, data, encodings, message) {
    message = message || 'decoding';
    return this.test(message, (t) => {
        t.test('per-field validation', (t2) => {
            for (const { field, value } of encodings) {
                if (!value)
                    continue;
                const fun = Type.decode(hexToBuf(value));
                const msg = `has expected decoding of ${field} field`;
                t2.has(fun, { [field]: data[field] }, msg);
            }
            t2.end();
        });
        t.test('full object validation', (t2) => {
            const fun = Type.decode(hexToBuf(fullEncoding(encodings)));
            t2.has(fun, data, 'has expected encoding of full object');
            t2.end();
        });
        t.end();
    });
});
const stringTable = new index_js_1.StringTable();
const functionData = {
    id: 123,
    name: stringTable.dedup('fn name'),
    systemName: stringTable.dedup('fn systemName'),
    filename: stringTable.dedup('fn filename'),
    startLine: 789
};
const functionEncodings = [
    { field: 'id', value: '087b' },
    { field: 'name', value: '1001' },
    { field: 'systemName', value: '1802' },
    { field: 'filename', value: '2003' },
    { field: 'startLine', value: '289506' }
];
tap_1.default.test('Function', (t) => {
    t.constructs(index_js_1.Function, functionData, functionEncodings);
    t.encodes(index_js_1.Function, functionData, functionEncodings);
    t.decodes(index_js_1.Function, functionData, functionEncodings);
    t.end();
});
const labelData = {
    key: stringTable.dedup('label key'),
    str: stringTable.dedup('label str'),
    num: 123,
    numUnit: stringTable.dedup('label numUnit')
};
const labelEncodings = [
    { field: 'key', value: '0804' },
    { field: 'str', value: '1005' },
    { field: 'num', value: '187b' },
    { field: 'numUnit', value: '2006' }
];
tap_1.default.test('Label', (t) => {
    t.constructs(index_js_1.Label, labelData, labelEncodings);
    t.encodes(index_js_1.Label, labelData, labelEncodings);
    t.decodes(index_js_1.Label, labelData, labelEncodings);
    t.end();
});
const lineData = {
    functionId: 1234,
    line: 5678
};
const lineEncodings = [
    { field: 'functionId', value: '08d209' },
    { field: 'line', value: '10ae2c' },
];
tap_1.default.test('Line', (t) => {
    t.constructs(index_js_1.Line, lineData, lineEncodings);
    t.encodes(index_js_1.Line, lineData, lineEncodings);
    t.decodes(index_js_1.Line, lineData, lineEncodings);
    t.end();
});
const locationData = {
    id: 12,
    mappingId: 34,
    address: 56,
    line: [lineData],
    isFolded: true
};
const locationEncodings = [
    { field: 'id', value: '080c' },
    { field: 'mappingId', value: '1022' },
    { field: 'address', value: '1838' },
    { field: 'line', value: embeddedField('22', lineEncodings) },
    { field: 'isFolded', value: '2801' },
];
tap_1.default.test('Location', (t) => {
    t.constructs(index_js_1.Location, locationData, locationEncodings);
    t.encodes(index_js_1.Location, locationData, locationEncodings);
    t.decodes(index_js_1.Location, locationData, locationEncodings);
    t.end();
});
const mappingData = {
    id: 1,
    memoryStart: 2,
    memoryLimit: 3,
    fileOffset: 4,
    filename: stringTable.dedup('mapping filename'),
    buildId: stringTable.dedup('mapping build id'),
    hasFunctions: true,
    hasFilenames: true,
    hasLineNumbers: true,
    hasInlineFrames: true,
};
const mappingEncodings = [
    { field: 'id', value: '0801' },
    { field: 'memoryStart', value: '1002' },
    { field: 'memoryLimit', value: '1803' },
    { field: 'fileOffset', value: '2004' },
    { field: 'filename', value: '2807' },
    { field: 'buildId', value: '3008' },
    { field: 'hasFunctions', value: '3801' },
    { field: 'hasFilenames', value: '4001' },
    { field: 'hasLineNumbers', value: '4801' },
    { field: 'hasInlineFrames', value: '5001' },
];
tap_1.default.test('Mapping', (t) => {
    t.constructs(index_js_1.Mapping, mappingData, mappingEncodings);
    t.encodes(index_js_1.Mapping, mappingData, mappingEncodings);
    t.decodes(index_js_1.Mapping, mappingData, mappingEncodings);
    t.end();
});
const sampleData = {
    locationId: [1, 2, 3],
    value: [...Array(180).keys()],
    label: [labelData]
};
const sampleEncodings = [
    { field: 'locationId', value: '0a03010203' },
    { field: 'value', value: embeddedField('12', sampleData.value.map(x => ({ field: '', value: hexVarInt(x) }))) },
    { field: 'label', value: embeddedField('1a', labelEncodings) },
];
tap_1.default.test('Sample', (t) => {
    t.constructs(index_js_1.Sample, sampleData, sampleEncodings);
    t.encodes(index_js_1.Sample, sampleData, sampleEncodings);
    t.decodes(index_js_1.Sample, sampleData, sampleEncodings);
    t.end();
});
const valueTypeData = {
    type: stringTable.dedup('value type type'),
    unit: stringTable.dedup('value type unit')
};
const valueTypeEncodings = [
    { field: 'type', value: '0809' },
    { field: 'unit', value: '100a' },
];
tap_1.default.test('ValueType', (t) => {
    t.constructs(index_js_1.ValueType, valueTypeData, valueTypeEncodings);
    t.encodes(index_js_1.ValueType, valueTypeData, valueTypeEncodings);
    t.decodes(index_js_1.ValueType, valueTypeData, valueTypeEncodings);
    t.end();
});
const profileData = {
    sampleType: [valueTypeData],
    sample: [sampleData],
    mapping: [mappingData],
    location: [locationData],
    function: [functionData],
    stringTable,
    timeNanos: 1000000n,
    durationNanos: 1234,
    periodType: valueTypeData,
    period: 1234 / 2,
    comment: [
        stringTable.dedup('some very very very very very very very very very very very very very very very very very very very very very very very very comment'),
        stringTable.dedup('another comment')
    ]
};
const profileEncodings = [
    { field: 'sampleType', value: embeddedField('0a', valueTypeEncodings) },
    { field: 'sample', value: embeddedField('12', sampleEncodings) },
    { field: 'mapping', value: embeddedField('1a', mappingEncodings) },
    { field: 'location', value: embeddedField('22', locationEncodings) },
    { field: 'function', value: embeddedField('2a', functionEncodings) },
    { field: 'stringTable', value: encodeStringTable(stringTable) },
    { field: 'timeNanos', value: '48c0843d' },
    { field: 'durationNanos', value: '50d209' },
    { field: 'periodType', value: embeddedField('5a', valueTypeEncodings) },
    { field: 'period', value: '60e904' },
    { field: 'comment', value: '6a020b0c' },
];
tap_1.default.test('Profile', (t) => {
    t.constructs(index_js_1.Profile, profileData, profileEncodings);
    t.encodes(index_js_1.Profile, profileData, profileEncodings);
    // Profiles additionally can be encoded asynchronously to break up
    // encoding into smaller chunks to have less latency impact.
    t.test('async encoding', (t) => {
        t.test('per-field validation', async (t2) => {
            for (const { field, value } of profileEncodings) {
                const fun = new index_js_1.Profile({
                    // Hack to exclude stringTable data from any checks except for the string table itself
                    stringTable: new index_js_1.StringTable(index_js_1.emptyTableToken),
                    [field]: profileData[field]
                });
                const msg = `has expected encoding of ${field} field`;
                t2.equal(bufToHex(await fun.encodeAsync()), value, msg);
            }
            t2.end();
        });
        t.test('full object validation', async (t2) => {
            const fun = new index_js_1.Profile(profileData);
            t2.equal(bufToHex(await fun.encodeAsync()), fullEncoding(profileEncodings), 'has expected encoding of full object');
            t2.end();
        });
        t.end();
    });
    t.decodes(index_js_1.Profile, profileData, profileEncodings);
    t.end();
});
function encodeStringTable(strings) {
    return strings.strings.map(s => {
        const buf = new TextEncoder().encode(s);
        return `32${hexVarInt(buf.length)}${bufToHex(buf)}`;
    }).join('');
}
function hexNum(d) {
    let hex = Number(d).toString(16);
    if (hex.length == 1) {
        hex = "0" + hex;
    }
    return hex;
}
function hexVarInt(num) {
    let n = BigInt(num);
    if (n < 0) {
        // take two's complement to encode negative number
        n = 2n ** 64n - n;
    }
    let str = '';
    const maxbits = 7n;
    const max = (1n << maxbits) - 1n;
    while (n > max) {
        str += hexNum(Number((n & max) | (1n << maxbits)));
        n >>= maxbits;
    }
    str += hexNum(Number(n));
    return str;
}
function embeddedField(fieldBit, data) {
    const encoded = fullEncoding(data);
    const size = hexVarInt(encoded.length / 2);
    return [fieldBit, size, encoded].join('');
}
function fullEncoding(encodings) {
    return encodings.map(e => e.value).join('');
}
function hexToBuf(hex) {
    return Uint8Array.from((hex.match(/.{2}/g) || []).map(v => parseInt(v, 16)));
}
function bufToHex(buf) {
    return Array.from(buf).map(hexNum).join('');
}
tap_1.default.test('StringTable', (t) => {
    t.test('encodes correctly', (t) => {
        const encodings = {
            '': '3200',
            'hello': '320568656c6c6f'
        };
        const table = new index_js_1.StringTable();
        t.equal(bufToHex(table.encode()), encodings['']);
        table.dedup('hello');
        t.equal(bufToHex(table.encode()), encodings[''] + encodings['hello']);
        t.end();
    });
    t.end();
});
function profileToObject(profile) {
    profile.stringTable = profile.stringTable.strings;
    return profile;
}
tap_1.default.test('Protobufjs compat', (t) => {
    t.test('encodes correctly', (t) => {
        const profile = new index_js_1.Profile(profileData);
        const encodedProfile = profile.encode();
        const decodedProfile = decode(encodedProfile);
        t.same(profileToObject(profile), toObject(decodedProfile, { longs: String, defaults: true }));
        t.end();
    });
    t.test('decodes correctly', (t) => {
        const buf = (0, zlib_1.gunzipSync)(fs.readFileSync('./testing/test.pprof'));
        t.same(profileToObject(index_js_1.Profile.decode(buf)), toObject(decode(buf), { longs: String, defaults: true }));
        t.end();
    });
    t.end();
});
//# sourceMappingURL=index.test.js.map