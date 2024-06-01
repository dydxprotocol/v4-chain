"use strict";
var __assign = (this && this.__assign) || function () {
    __assign = Object.assign || function(t) {
        for (var s, i = 1, n = arguments.length; i < n; i++) {
            s = arguments[i];
            for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p))
                t[p] = s[p];
        }
        return t;
    };
    return __assign.apply(this, arguments);
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.parse = void 0;
var lodash_1 = require("lodash");
var util_1 = require("util");
var typesOfSchema_1 = require("./typesOfSchema");
var AST_1 = require("./types/AST");
var JSONSchema_1 = require("./types/JSONSchema");
var utils_1 = require("./utils");
function parse(schema, options, keyName, processed, usedNames) {
    if (processed === void 0) { processed = new Map(); }
    if (usedNames === void 0) { usedNames = new Set(); }
    if ((0, JSONSchema_1.isPrimitive)(schema)) {
        return parseLiteral(schema, keyName);
    }
    var types = (0, typesOfSchema_1.typesOfSchema)(schema);
    if (types.length === 1) {
        var ast_1 = parseAsTypeWithCache(schema, types[0], options, keyName, processed, usedNames);
        (0, utils_1.log)('blue', 'parser', 'Types:', types, 'Input:', schema, 'Output:', ast_1);
        return ast_1;
    }
    // Be careful to first process the intersection before processing its params,
    // so that it gets first pick for standalone name.
    var ast = parseAsTypeWithCache({
        $id: schema.$id,
        allOf: [],
        description: schema.description,
        title: schema.title
    }, 'ALL_OF', options, keyName, processed, usedNames);
    ast.params = types.map(function (type) {
        // We hoist description (for comment) and id/title (for standaloneName)
        // to the parent intersection type, so we remove it from the children.
        return parseAsTypeWithCache((0, utils_1.maybeStripNameHints)(schema), type, options, keyName, processed, usedNames);
    });
    (0, utils_1.log)('blue', 'parser', 'Types:', types, 'Input:', schema, 'Output:', ast);
    return ast;
}
exports.parse = parse;
function parseAsTypeWithCache(schema, type, options, keyName, processed, usedNames) {
    if (processed === void 0) { processed = new Map(); }
    if (usedNames === void 0) { usedNames = new Set(); }
    // If we've seen this node before, return it.
    var cachedTypeMap = processed.get(schema);
    if (!cachedTypeMap) {
        cachedTypeMap = new Map();
        processed.set(schema, cachedTypeMap);
    }
    var cachedAST = cachedTypeMap.get(type);
    if (cachedAST) {
        return cachedAST;
    }
    // Cache processed ASTs before they are actually computed, then update
    // them in place using set(). This is to avoid cycles.
    // TODO: Investigate alternative approaches (lazy-computing nodes, etc.)
    var ast = {};
    cachedTypeMap.set(type, ast);
    // Update the AST in place. This updates the `processed` cache, as well
    // as any nodes that directly reference the node.
    return Object.assign(ast, parseNonLiteral(schema, type, options, keyName, processed, usedNames));
}
function parseLiteral(schema, keyName) {
    return {
        keyName: keyName,
        params: schema,
        type: 'LITERAL'
    };
}
function parseNonLiteral(schema, type, options, keyName, processed, usedNames) {
    var definitions = getDefinitionsMemoized((0, JSONSchema_1.getRootSchema)(schema)); // TODO
    var keyNameFromDefinition = (0, lodash_1.findKey)(definitions, function (_) { return _ === schema; });
    switch (type) {
        case 'ALL_OF':
            return {
                comment: schema.description,
                keyName: keyName,
                standaloneName: standaloneName(schema, keyNameFromDefinition, usedNames),
                params: schema.allOf.map(function (_) { return parse(_, options, undefined, processed, usedNames); }),
                type: 'INTERSECTION'
            };
        case 'ANY':
            return __assign(__assign({}, (options.unknownAny ? AST_1.T_UNKNOWN : AST_1.T_ANY)), { comment: schema.description, keyName: keyName, standaloneName: standaloneName(schema, keyNameFromDefinition, usedNames) });
        case 'ANY_OF':
            return {
                comment: schema.description,
                keyName: keyName,
                standaloneName: standaloneName(schema, keyNameFromDefinition, usedNames),
                params: schema.anyOf.map(function (_) { return parse(_, options, undefined, processed, usedNames); }),
                type: 'UNION'
            };
        case 'BOOLEAN':
            return {
                comment: schema.description,
                keyName: keyName,
                standaloneName: standaloneName(schema, keyNameFromDefinition, usedNames),
                type: 'BOOLEAN'
            };
        case 'CUSTOM_TYPE':
            return {
                comment: schema.description,
                keyName: keyName,
                params: schema.tsType,
                standaloneName: standaloneName(schema, keyNameFromDefinition, usedNames),
                type: 'CUSTOM_TYPE'
            };
        case 'NAMED_ENUM':
            return {
                comment: schema.description,
                keyName: keyName,
                standaloneName: standaloneName(schema, keyNameFromDefinition !== null && keyNameFromDefinition !== void 0 ? keyNameFromDefinition : keyName, usedNames),
                params: schema.enum.map(function (_, n) { return ({
                    ast: parse(_, options, undefined, processed, usedNames),
                    keyName: schema.tsEnumNames[n]
                }); }),
                type: 'ENUM'
            };
        case 'NAMED_SCHEMA':
            return newInterface(schema, options, processed, usedNames, keyName);
        case 'NULL':
            return {
                comment: schema.description,
                keyName: keyName,
                standaloneName: standaloneName(schema, keyNameFromDefinition, usedNames),
                type: 'NULL'
            };
        case 'NUMBER':
            return {
                comment: schema.description,
                keyName: keyName,
                standaloneName: standaloneName(schema, keyNameFromDefinition, usedNames),
                type: 'NUMBER'
            };
        case 'OBJECT':
            return {
                comment: schema.description,
                keyName: keyName,
                standaloneName: standaloneName(schema, keyNameFromDefinition, usedNames),
                type: 'OBJECT'
            };
        case 'ONE_OF':
            return {
                comment: schema.description,
                keyName: keyName,
                standaloneName: standaloneName(schema, keyNameFromDefinition, usedNames),
                params: schema.oneOf.map(function (_) { return parse(_, options, undefined, processed, usedNames); }),
                type: 'UNION'
            };
        case 'REFERENCE':
            throw Error((0, util_1.format)('Refs should have been resolved by the resolver!', schema));
        case 'STRING':
            return {
                comment: schema.description,
                keyName: keyName,
                standaloneName: standaloneName(schema, keyNameFromDefinition, usedNames),
                type: 'STRING'
            };
        case 'TYPED_ARRAY':
            if (Array.isArray(schema.items)) {
                // normalised to not be undefined
                var minItems_1 = schema.minItems;
                var maxItems_1 = schema.maxItems;
                var arrayType = {
                    comment: schema.description,
                    keyName: keyName,
                    maxItems: maxItems_1,
                    minItems: minItems_1,
                    standaloneName: standaloneName(schema, keyNameFromDefinition, usedNames),
                    params: schema.items.map(function (_) { return parse(_, options, undefined, processed, usedNames); }),
                    type: 'TUPLE'
                };
                if (schema.additionalItems === true) {
                    arrayType.spreadParam = options.unknownAny ? AST_1.T_UNKNOWN : AST_1.T_ANY;
                }
                else if (schema.additionalItems) {
                    arrayType.spreadParam = parse(schema.additionalItems, options, undefined, processed, usedNames);
                }
                return arrayType;
            }
            else {
                return {
                    comment: schema.description,
                    keyName: keyName,
                    standaloneName: standaloneName(schema, keyNameFromDefinition, usedNames),
                    params: parse(schema.items, options, undefined, processed, usedNames),
                    type: 'ARRAY'
                };
            }
        case 'UNION':
            return {
                comment: schema.description,
                keyName: keyName,
                standaloneName: standaloneName(schema, keyNameFromDefinition, usedNames),
                params: schema.type.map(function (type) {
                    var member = __assign(__assign({}, (0, lodash_1.omit)(schema, '$id', 'description', 'title')), { type: type });
                    return parse((0, utils_1.maybeStripDefault)(member), options, undefined, processed, usedNames);
                }),
                type: 'UNION'
            };
        case 'UNNAMED_ENUM':
            return {
                comment: schema.description,
                keyName: keyName,
                standaloneName: standaloneName(schema, keyNameFromDefinition, usedNames),
                params: schema.enum.map(function (_) { return parse(_, options, undefined, processed, usedNames); }),
                type: 'UNION'
            };
        case 'UNNAMED_SCHEMA':
            return newInterface(schema, options, processed, usedNames, keyName, keyNameFromDefinition);
        case 'UNTYPED_ARRAY':
            // normalised to not be undefined
            var minItems = schema.minItems;
            var maxItems = typeof schema.maxItems === 'number' ? schema.maxItems : -1;
            var params = options.unknownAny ? AST_1.T_UNKNOWN : AST_1.T_ANY;
            if (minItems > 0 || maxItems >= 0) {
                return {
                    comment: schema.description,
                    keyName: keyName,
                    maxItems: schema.maxItems,
                    minItems: minItems,
                    // create a tuple of length N
                    params: Array(Math.max(maxItems, minItems) || 0).fill(params),
                    // if there is no maximum, then add a spread item to collect the rest
                    spreadParam: maxItems >= 0 ? undefined : params,
                    standaloneName: standaloneName(schema, keyNameFromDefinition, usedNames),
                    type: 'TUPLE'
                };
            }
            return {
                comment: schema.description,
                keyName: keyName,
                params: params,
                standaloneName: standaloneName(schema, keyNameFromDefinition, usedNames),
                type: 'ARRAY'
            };
    }
}
/**
 * Compute a schema name using a series of fallbacks
 */
function standaloneName(schema, keyNameFromDefinition, usedNames) {
    var name = schema.title || schema.$id || keyNameFromDefinition;
    if (name) {
        return (0, utils_1.generateName)(name, usedNames);
    }
}
function newInterface(schema, options, processed, usedNames, keyName, keyNameFromDefinition) {
    var name = standaloneName(schema, keyNameFromDefinition, usedNames);
    return {
        comment: schema.description,
        keyName: keyName,
        params: parseSchema(schema, options, processed, usedNames, name),
        standaloneName: name,
        superTypes: parseSuperTypes(schema, options, processed, usedNames),
        type: 'INTERFACE'
    };
}
function parseSuperTypes(schema, options, processed, usedNames) {
    // Type assertion needed because of dereferencing step
    // TODO: Type it upstream
    var superTypes = schema.extends;
    if (!superTypes) {
        return [];
    }
    return superTypes.map(function (_) { return parse(_, options, undefined, processed, usedNames); });
}
/**
 * Helper to parse schema properties into params on the parent schema's type
 */
function parseSchema(schema, options, processed, usedNames, parentSchemaName) {
    var asts = (0, lodash_1.map)(schema.properties, function (value, key) { return ({
        ast: parse(value, options, key, processed, usedNames),
        isPatternProperty: false,
        isRequired: (0, lodash_1.includes)(schema.required || [], key),
        isUnreachableDefinition: false,
        keyName: key
    }); });
    var singlePatternProperty = false;
    if (schema.patternProperties) {
        // partially support patternProperties. in the case that
        // additionalProperties is not set, and there is only a single
        // value definition, we can validate against that.
        singlePatternProperty = !schema.additionalProperties && Object.keys(schema.patternProperties).length === 1;
        asts = asts.concat((0, lodash_1.map)(schema.patternProperties, function (value, key) {
            var ast = parse(value, options, key, processed, usedNames);
            var comment = "This interface was referenced by `".concat(parentSchemaName, "`'s JSON-Schema definition\nvia the `patternProperty` \"").concat(key, "\".");
            ast.comment = ast.comment ? "".concat(ast.comment, "\n\n").concat(comment) : comment;
            return {
                ast: ast,
                isPatternProperty: !singlePatternProperty,
                isRequired: singlePatternProperty || (0, lodash_1.includes)(schema.required || [], key),
                isUnreachableDefinition: false,
                keyName: singlePatternProperty ? '[k: string]' : key
            };
        }));
    }
    if (options.unreachableDefinitions) {
        asts = asts.concat((0, lodash_1.map)(schema.$defs, function (value, key) {
            var ast = parse(value, options, key, processed, usedNames);
            var comment = "This interface was referenced by `".concat(parentSchemaName, "`'s JSON-Schema\nvia the `definition` \"").concat(key, "\".");
            ast.comment = ast.comment ? "".concat(ast.comment, "\n\n").concat(comment) : comment;
            return {
                ast: ast,
                isPatternProperty: false,
                isRequired: (0, lodash_1.includes)(schema.required || [], key),
                isUnreachableDefinition: true,
                keyName: key
            };
        }));
    }
    // handle additionalProperties
    switch (schema.additionalProperties) {
        case undefined:
        case true:
            if (singlePatternProperty) {
                return asts;
            }
            return asts.concat({
                ast: options.unknownAny ? AST_1.T_UNKNOWN_ADDITIONAL_PROPERTIES : AST_1.T_ANY_ADDITIONAL_PROPERTIES,
                isPatternProperty: false,
                isRequired: true,
                isUnreachableDefinition: false,
                keyName: '[k: string]'
            });
        case false:
            return asts;
        // pass "true" as the last param because in TS, properties
        // defined via index signatures are already optional
        default:
            return asts.concat({
                ast: parse(schema.additionalProperties, options, '[k: string]', processed, usedNames),
                isPatternProperty: false,
                isRequired: true,
                isUnreachableDefinition: false,
                keyName: '[k: string]'
            });
    }
}
function getDefinitions(schema, isSchema, processed) {
    if (isSchema === void 0) { isSchema = true; }
    if (processed === void 0) { processed = new Set(); }
    if (processed.has(schema)) {
        return {};
    }
    processed.add(schema);
    if (Array.isArray(schema)) {
        return schema.reduce(function (prev, cur) { return (__assign(__assign({}, prev), getDefinitions(cur, false, processed))); }, {});
    }
    if ((0, lodash_1.isPlainObject)(schema)) {
        return __assign(__assign({}, (isSchema && hasDefinitions(schema) ? schema.$defs : {})), Object.keys(schema).reduce(function (prev, cur) { return (__assign(__assign({}, prev), getDefinitions(schema[cur], false, processed))); }, {}));
    }
    return {};
}
var getDefinitionsMemoized = (0, lodash_1.memoize)(getDefinitions);
/**
 * TODO: Reduce rate of false positives
 */
function hasDefinitions(schema) {
    return '$defs' in schema;
}
//# sourceMappingURL=parser.js.map