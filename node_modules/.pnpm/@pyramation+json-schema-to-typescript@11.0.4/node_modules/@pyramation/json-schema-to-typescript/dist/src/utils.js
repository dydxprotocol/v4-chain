"use strict";
var __spreadArray = (this && this.__spreadArray) || function (to, from, pack) {
    if (pack || arguments.length === 2) for (var i = 0, l = from.length, ar; i < l; i++) {
        if (ar || !(i in from)) {
            if (!ar) ar = Array.prototype.slice.call(from, 0, i);
            ar[i] = from[i];
        }
    }
    return to.concat(ar || Array.prototype.slice.call(from));
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.isSchemaLike = exports.appendToDescription = exports.maybeStripNameHints = exports.maybeStripDefault = exports.pathTransform = exports.escapeBlockComment = exports.log = exports.error = exports.generateName = exports.toSafeString = exports.stripExtension = exports.justName = exports.traverse = exports.Try = void 0;
var lodash_1 = require("lodash");
var path_1 = require("path");
var JSONSchema_1 = require("./types/JSONSchema");
// TODO: pull out into a separate package
function Try(fn, err) {
    try {
        return fn();
    }
    catch (e) {
        return err(e);
    }
}
exports.Try = Try;
// keys that shouldn't be traversed by the catchall step
var BLACKLISTED_KEYS = new Set([
    'id',
    '$defs',
    '$id',
    '$schema',
    'title',
    'description',
    'default',
    'multipleOf',
    'maximum',
    'exclusiveMaximum',
    'minimum',
    'exclusiveMinimum',
    'maxLength',
    'minLength',
    'pattern',
    'additionalItems',
    'items',
    'maxItems',
    'minItems',
    'uniqueItems',
    'maxProperties',
    'minProperties',
    'required',
    'additionalProperties',
    'definitions',
    'properties',
    'patternProperties',
    'dependencies',
    'enum',
    'type',
    'allOf',
    'anyOf',
    'oneOf',
    'not'
]);
function traverseObjectKeys(obj, callback, processed) {
    Object.keys(obj).forEach(function (k) {
        if (obj[k] && typeof obj[k] === 'object' && !Array.isArray(obj[k])) {
            traverse(obj[k], callback, processed, k);
        }
    });
}
function traverseArray(arr, callback, processed) {
    arr.forEach(function (s, k) { return traverse(s, callback, processed, k.toString()); });
}
function traverse(schema, callback, processed, key) {
    if (processed === void 0) { processed = new Set(); }
    // Handle recursive schemas
    if (processed.has(schema)) {
        return;
    }
    processed.add(schema);
    callback(schema, key !== null && key !== void 0 ? key : null);
    if (schema.anyOf) {
        traverseArray(schema.anyOf, callback, processed);
    }
    if (schema.allOf) {
        traverseArray(schema.allOf, callback, processed);
    }
    if (schema.oneOf) {
        traverseArray(schema.oneOf, callback, processed);
    }
    if (schema.properties) {
        traverseObjectKeys(schema.properties, callback, processed);
    }
    if (schema.patternProperties) {
        traverseObjectKeys(schema.patternProperties, callback, processed);
    }
    if (schema.additionalProperties && typeof schema.additionalProperties === 'object') {
        traverse(schema.additionalProperties, callback, processed);
    }
    if (schema.items) {
        var items = schema.items;
        if (Array.isArray(items)) {
            traverseArray(items, callback, processed);
        }
        else {
            traverse(items, callback, processed);
        }
    }
    if (schema.additionalItems && typeof schema.additionalItems === 'object') {
        traverse(schema.additionalItems, callback, processed);
    }
    if (schema.dependencies) {
        if (Array.isArray(schema.dependencies)) {
            traverseArray(schema.dependencies, callback, processed);
        }
        else {
            traverseObjectKeys(schema.dependencies, callback, processed);
        }
    }
    if (schema.definitions) {
        traverseObjectKeys(schema.definitions, callback, processed);
    }
    if (schema.$defs) {
        traverseObjectKeys(schema.$defs, callback, processed);
    }
    if (schema.not) {
        traverse(schema.not, callback, processed);
    }
    // technically you can put definitions on any key
    Object.keys(schema)
        .filter(function (key) { return !BLACKLISTED_KEYS.has(key); })
        .forEach(function (key) {
        var child = schema[key];
        if (child && typeof child === 'object') {
            traverseObjectKeys(child, callback, processed);
        }
    });
}
exports.traverse = traverse;
/**
 * Eg. `foo/bar/baz.json` => `baz`
 */
function justName(filename) {
    if (filename === void 0) { filename = ''; }
    return stripExtension((0, path_1.basename)(filename));
}
exports.justName = justName;
/**
 * Avoid appending "js" to top-level unnamed schemas
 */
function stripExtension(filename) {
    return filename.replace((0, path_1.extname)(filename), '');
}
exports.stripExtension = stripExtension;
/**
 * Convert a string that might contain spaces or special characters to one that
 * can safely be used as a TypeScript interface or enum name.
 */
function toSafeString(string) {
    // identifiers in javaScript/ts:
    // First character: a-zA-Z | _ | $
    // Rest: a-zA-Z | _ | $ | 0-9
    return (0, lodash_1.upperFirst)(
    // remove accents, umlauts, ... by their basic latin letters
    (0, lodash_1.deburr)(string)
        // replace chars which are not valid for typescript identifiers with whitespace
        .replace(/(^\s*[^a-zA-Z_$])|([^a-zA-Z_$\d])/g, ' ')
        // uppercase leading underscores followed by lowercase
        .replace(/^_[a-z]/g, function (match) { return match.toUpperCase(); })
        // remove non-leading underscores followed by lowercase (convert snake_case)
        .replace(/_[a-z]/g, function (match) { return match.substr(1, match.length).toUpperCase(); })
        // uppercase letters after digits, dollars
        .replace(/([\d$]+[a-zA-Z])/g, function (match) { return match.toUpperCase(); })
        // uppercase first letter after whitespace
        .replace(/\s+([a-zA-Z])/g, function (match) { return (0, lodash_1.trim)(match.toUpperCase()); })
        // remove remaining whitespace
        .replace(/\s/g, ''));
}
exports.toSafeString = toSafeString;
function generateName(from, usedNames) {
    var name = toSafeString(from);
    if (!name) {
        name = 'NoName';
    }
    // increment counter until we find a free name
    if (usedNames.has(name)) {
        var counter = 1;
        var nameWithCounter = "".concat(name).concat(counter);
        while (usedNames.has(nameWithCounter)) {
            nameWithCounter = "".concat(name).concat(counter);
            counter++;
        }
        name = nameWithCounter;
    }
    usedNames.add(name);
    return name;
}
exports.generateName = generateName;
function error() {
    var _a;
    var messages = [];
    for (var _i = 0; _i < arguments.length; _i++) {
        messages[_i] = arguments[_i];
    }
    if (!process.env.VERBOSE) {
        return console.error(messages);
    }
    console.error.apply(console, __spreadArray([(_a = getStyledTextForLogging('red')) === null || _a === void 0 ? void 0 : _a('error')], messages, false));
}
exports.error = error;
function log(style, title) {
    var _a;
    var messages = [];
    for (var _i = 2; _i < arguments.length; _i++) {
        messages[_i - 2] = arguments[_i];
    }
    if (!process.env.VERBOSE) {
        return;
    }
    var lastMessage = null;
    if (messages.length > 1 && typeof messages[messages.length - 1] !== 'string') {
        lastMessage = messages.splice(messages.length - 1, 1);
    }
    console.info.apply(console, __spreadArray([require('cli-color').whiteBright.bgCyan('debug'), (_a = getStyledTextForLogging(style)) === null || _a === void 0 ? void 0 : _a(title)], messages, false));
    if (lastMessage) {
        console.dir(lastMessage, { depth: 6, maxArrayLength: 6 });
    }
}
exports.log = log;
function getStyledTextForLogging(style) {
    if (!process.env.VERBOSE) {
        return;
    }
    switch (style) {
        case 'blue':
            return require('cli-color').whiteBright.bgBlue;
        case 'cyan':
            return require('cli-color').whiteBright.bgCyan;
        case 'green':
            return require('cli-color').whiteBright.bgGreen;
        case 'magenta':
            return require('cli-color').whiteBright.bgMagenta;
        case 'red':
            return require('cli-color').whiteBright.bgRedBright;
        case 'white':
            return require('cli-color').black.bgWhite;
        case 'yellow':
            return require('cli-color').whiteBright.bgYellow;
    }
}
/**
 * escape block comments in schema descriptions so that they don't unexpectedly close JSDoc comments in generated typescript interfaces
 */
function escapeBlockComment(schema) {
    var replacer = '* /';
    if (schema === null || typeof schema !== 'object') {
        return;
    }
    for (var _i = 0, _a = Object.keys(schema); _i < _a.length; _i++) {
        var key = _a[_i];
        if (key === 'description' && typeof schema[key] === 'string') {
            schema[key] = schema[key].replace(/\*\//g, replacer);
        }
    }
}
exports.escapeBlockComment = escapeBlockComment;
/*
the following logic determines the out path by comparing the in path to the users specified out path.
For example, if input directory MultiSchema looks like:
  MultiSchema/foo/a.json
  MultiSchema/bar/fuzz/c.json
  MultiSchema/bar/d.json
And the user wants the outputs to be in MultiSchema/Out, then this code will be able to map the inner directories foo, bar, and fuzz into the intended Out directory like so:
  MultiSchema/Out/foo/a.json
  MultiSchema/Out/bar/fuzz/c.json
  MultiSchema/Out/bar/d.json
*/
function pathTransform(outputPath, inputPath, filePath) {
    var inPathList = (0, path_1.normalize)(inputPath).split(path_1.sep);
    var filePathList = (0, path_1.dirname)((0, path_1.normalize)(filePath)).split(path_1.sep);
    var filePathRel = filePathList.filter(function (f, i) { return f !== inPathList[i]; });
    return path_1.posix.join.apply(path_1.posix, __spreadArray([path_1.posix.normalize(outputPath)], filePathRel, false));
}
exports.pathTransform = pathTransform;
/**
 * Removes the schema's `default` property if it doesn't match the schema's `type` property.
 * Useful when parsing unions.
 *
 * Mutates `schema`.
 */
function maybeStripDefault(schema) {
    if (!('default' in schema)) {
        return schema;
    }
    switch (schema.type) {
        case 'array':
            if (Array.isArray(schema.default)) {
                return schema;
            }
            break;
        case 'boolean':
            if (typeof schema.default === 'boolean') {
                return schema;
            }
            break;
        case 'integer':
        case 'number':
            if (typeof schema.default === 'number') {
                return schema;
            }
            break;
        case 'string':
            if (typeof schema.default === 'string') {
                return schema;
            }
            break;
        case 'null':
            if (schema.default === null) {
                return schema;
            }
            break;
        case 'object':
            if ((0, lodash_1.isPlainObject)(schema.default)) {
                return schema;
            }
            break;
    }
    delete schema.default;
    return schema;
}
exports.maybeStripDefault = maybeStripDefault;
/**
 * Removes the schema's `$id`, `name`, and `description` properties
 * if they exist.
 * Useful when parsing intersections.
 *
 * Mutates `schema`.
 */
function maybeStripNameHints(schema) {
    if ('$id' in schema) {
        delete schema.$id;
    }
    if ('description' in schema) {
        delete schema.description;
    }
    if ('name' in schema) {
        delete schema.name;
    }
    return schema;
}
exports.maybeStripNameHints = maybeStripNameHints;
function appendToDescription(existingDescription) {
    var values = [];
    for (var _i = 1; _i < arguments.length; _i++) {
        values[_i - 1] = arguments[_i];
    }
    if (existingDescription) {
        return "".concat(existingDescription, "\n\n").concat(values.join('\n'));
    }
    return values.join('\n');
}
exports.appendToDescription = appendToDescription;
function isSchemaLike(schema) {
    if (!(0, lodash_1.isPlainObject)(schema)) {
        return false;
    }
    var parent = schema[JSONSchema_1.Parent];
    if (parent === null) {
        return true;
    }
    var JSON_SCHEMA_KEYWORDS = [
        '$defs',
        'allOf',
        'anyOf',
        'definitions',
        'dependencies',
        'enum',
        'not',
        'oneOf',
        'patternProperties',
        'properties',
        'required'
    ];
    if (JSON_SCHEMA_KEYWORDS.some(function (_) { return parent[_] === schema; })) {
        return false;
    }
    return true;
}
exports.isSchemaLike = isSchemaLike;
//# sourceMappingURL=utils.js.map