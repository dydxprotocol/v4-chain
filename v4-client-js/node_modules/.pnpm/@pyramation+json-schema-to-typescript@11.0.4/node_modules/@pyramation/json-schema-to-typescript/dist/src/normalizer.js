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
exports.normalize = void 0;
var JSONSchema_1 = require("./types/JSONSchema");
var utils_1 = require("./utils");
var util_1 = require("util");
var rules = new Map();
function hasType(schema, type) {
    return schema.type === type || (Array.isArray(schema.type) && schema.type.includes(type));
}
function isObjectType(schema) {
    return schema.properties !== undefined || hasType(schema, 'object') || hasType(schema, 'any');
}
function isArrayType(schema) {
    return schema.items !== undefined || hasType(schema, 'array') || hasType(schema, 'any');
}
rules.set('Remove `type=["null"]` if `enum=[null]`', function (schema) {
    if (Array.isArray(schema.enum) &&
        schema.enum.some(function (e) { return e === null; }) &&
        Array.isArray(schema.type) &&
        schema.type.includes('null')) {
        schema.type = schema.type.filter(function (type) { return type !== 'null'; });
    }
});
rules.set('Destructure unary types', function (schema) {
    if (schema.type && Array.isArray(schema.type) && schema.type.length === 1) {
        schema.type = schema.type[0];
    }
});
rules.set('Add empty `required` property if none is defined', function (schema) {
    if (isObjectType(schema) && !('required' in schema)) {
        schema.required = [];
    }
});
rules.set('Transform `required`=false to `required`=[]', function (schema) {
    if (schema.required === false) {
        schema.required = [];
    }
});
rules.set('Default additionalProperties', function (schema, _, options) {
    if (isObjectType(schema) && !('additionalProperties' in schema) && schema.patternProperties === undefined) {
        schema.additionalProperties = options.additionalProperties;
    }
});
rules.set('Transform id to $id', function (schema, fileName) {
    if (!(0, utils_1.isSchemaLike)(schema)) {
        return;
    }
    if (schema.id && schema.$id && schema.id !== schema.$id) {
        throw ReferenceError("Schema must define either id or $id, not both. Given id=".concat(schema.id, ", $id=").concat(schema.$id, " in ").concat(fileName));
    }
    if (schema.id) {
        schema.$id = schema.id;
        delete schema.id;
    }
});
rules.set('Add an $id to anything that needs it', function (schema, fileName, _options, _key, dereferencedPaths) {
    if (!(0, utils_1.isSchemaLike)(schema)) {
        return;
    }
    // Top-level schema
    if (!schema.$id && !schema[JSONSchema_1.Parent]) {
        schema.$id = (0, utils_1.toSafeString)((0, utils_1.justName)(fileName));
        return;
    }
    // Sub-schemas with references
    if (!isArrayType(schema) && !isObjectType(schema)) {
        return;
    }
    // We'll infer from $id and title downstream
    // TODO: Normalize upstream
    var dereferencedName = dereferencedPaths.get(schema);
    if (!schema.$id && !schema.title && dereferencedName) {
        schema.$id = (0, utils_1.toSafeString)((0, utils_1.justName)(dereferencedName));
    }
    if (dereferencedName) {
        dereferencedPaths.delete(schema);
    }
});
rules.set('Escape closing JSDoc comment', function (schema) {
    (0, utils_1.escapeBlockComment)(schema);
});
rules.set('Add JSDoc comments for minItems and maxItems', function (schema) {
    if (!isArrayType(schema)) {
        return;
    }
    var commentsToAppend = [
        'minItems' in schema ? "@minItems ".concat(schema.minItems) : '',
        'maxItems' in schema ? "@maxItems ".concat(schema.maxItems) : ''
    ].filter(Boolean);
    if (commentsToAppend.length) {
        schema.description = utils_1.appendToDescription.apply(void 0, __spreadArray([schema.description], commentsToAppend, false));
    }
});
rules.set('Optionally remove maxItems and minItems', function (schema, _fileName, options) {
    if (!isArrayType(schema)) {
        return;
    }
    if ('minItems' in schema && options.ignoreMinAndMaxItems) {
        delete schema.minItems;
    }
    if ('maxItems' in schema && (options.ignoreMinAndMaxItems || options.maxItems === -1)) {
        delete schema.maxItems;
    }
});
rules.set('Normalize schema.minItems', function (schema, _fileName, options) {
    if (options.ignoreMinAndMaxItems) {
        return;
    }
    // make sure we only add the props onto array types
    if (!isArrayType(schema)) {
        return;
    }
    var minItems = schema.minItems;
    schema.minItems = typeof minItems === 'number' ? minItems : 0;
    // cannot normalize maxItems because maxItems = 0 has an actual meaning
});
rules.set('Remove maxItems if it is big enough to likely cause OOMs', function (schema, _fileName, options) {
    if (options.ignoreMinAndMaxItems || options.maxItems === -1) {
        return;
    }
    if (!isArrayType(schema)) {
        return;
    }
    var maxItems = schema.maxItems, minItems = schema.minItems;
    // minItems is guaranteed to be a number after the previous rule runs
    if (maxItems !== undefined && maxItems - minItems > options.maxItems) {
        delete schema.maxItems;
    }
});
rules.set('Normalize schema.items', function (schema, _fileName, options) {
    if (options.ignoreMinAndMaxItems) {
        return;
    }
    var maxItems = schema.maxItems, minItems = schema.minItems;
    var hasMaxItems = typeof maxItems === 'number' && maxItems >= 0;
    var hasMinItems = typeof minItems === 'number' && minItems > 0;
    if (schema.items && !Array.isArray(schema.items) && (hasMaxItems || hasMinItems)) {
        var items = schema.items;
        // create a tuple of length N
        var newItems = Array(maxItems || minItems || 0).fill(items);
        if (!hasMaxItems) {
            // if there is no maximum, then add a spread item to collect the rest
            schema.additionalItems = items;
        }
        schema.items = newItems;
    }
    if (Array.isArray(schema.items) && hasMaxItems && maxItems < schema.items.length) {
        // it's perfectly valid to provide 5 item defs but require maxItems 1
        // obviously we shouldn't emit a type for items that aren't expected
        schema.items = schema.items.slice(0, maxItems);
    }
    return schema;
});
rules.set('Remove extends, if it is empty', function (schema) {
    if (!schema.hasOwnProperty('extends')) {
        return;
    }
    if (schema.extends == null || (Array.isArray(schema.extends) && schema.extends.length === 0)) {
        delete schema.extends;
    }
});
rules.set('Make extends always an array, if it is defined', function (schema) {
    if (schema.extends == null) {
        return;
    }
    if (!Array.isArray(schema.extends)) {
        schema.extends = [schema.extends];
    }
});
rules.set('Transform definitions to $defs', function (schema, fileName) {
    if (schema.definitions && schema.$defs && !(0, util_1.isDeepStrictEqual)(schema.definitions, schema.$defs)) {
        throw ReferenceError("Schema must define either definitions or $defs, not both. Given id=".concat(schema.id, " in ").concat(fileName));
    }
    if (schema.definitions) {
        schema.$defs = schema.definitions;
        delete schema.definitions;
    }
});
rules.set('Transform const to singleton enum', function (schema) {
    if (schema.const !== undefined) {
        schema.enum = [schema.const];
        delete schema.const;
    }
});
function normalize(rootSchema, dereferencedPaths, filename, options) {
    rules.forEach(function (rule) { return (0, utils_1.traverse)(rootSchema, function (schema, key) { return rule(schema, filename, options, key, dereferencedPaths); }); });
    return rootSchema;
}
exports.normalize = normalize;
//# sourceMappingURL=normalizer.js.map