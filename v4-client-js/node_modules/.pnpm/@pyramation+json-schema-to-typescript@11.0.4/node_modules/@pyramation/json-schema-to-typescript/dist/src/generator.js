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
exports.generateType = exports.generate = void 0;
var lodash_1 = require("lodash");
var index_1 = require("./index");
var AST_1 = require("./types/AST");
var utils_1 = require("./utils");
function generate(ast, options) {
    if (options === void 0) { options = index_1.DEFAULT_OPTIONS; }
    return ([
        options.bannerComment,
        declareNamedTypes(ast, options, ast.standaloneName),
        declareNamedInterfaces(ast, options, ast.standaloneName),
        declareEnums(ast, options)
    ]
        .filter(Boolean)
        .join('\n\n') + '\n'); // trailing newline
}
exports.generate = generate;
function declareEnums(ast, options, processed) {
    if (processed === void 0) { processed = new Set(); }
    if (processed.has(ast)) {
        return '';
    }
    processed.add(ast);
    var type = '';
    switch (ast.type) {
        case 'ENUM':
            return generateStandaloneEnum(ast, options) + '\n';
        case 'ARRAY':
            return declareEnums(ast.params, options, processed);
        case 'UNION':
        case 'INTERSECTION':
            return ast.params.reduce(function (prev, ast) { return prev + declareEnums(ast, options, processed); }, '');
        case 'TUPLE':
            type = ast.params.reduce(function (prev, ast) { return prev + declareEnums(ast, options, processed); }, '');
            if (ast.spreadParam) {
                type += declareEnums(ast.spreadParam, options, processed);
            }
            return type;
        case 'INTERFACE':
            return getSuperTypesAndParams(ast).reduce(function (prev, ast) { return prev + declareEnums(ast, options, processed); }, '');
        default:
            return '';
    }
}
function declareNamedInterfaces(ast, options, rootASTName, processed) {
    if (processed === void 0) { processed = new Set(); }
    if (processed.has(ast)) {
        return '';
    }
    processed.add(ast);
    var type = '';
    switch (ast.type) {
        case 'ARRAY':
            type = declareNamedInterfaces(ast.params, options, rootASTName, processed);
            break;
        case 'INTERFACE':
            type = [
                (0, AST_1.hasStandaloneName)(ast) &&
                    (ast.standaloneName === rootASTName || options.declareExternallyReferenced) &&
                    generateStandaloneInterface(ast, options),
                getSuperTypesAndParams(ast)
                    .map(function (ast) { return declareNamedInterfaces(ast, options, rootASTName, processed); })
                    .filter(Boolean)
                    .join('\n')
            ]
                .filter(Boolean)
                .join('\n');
            break;
        case 'INTERSECTION':
        case 'TUPLE':
        case 'UNION':
            type = ast.params
                .map(function (_) { return declareNamedInterfaces(_, options, rootASTName, processed); })
                .filter(Boolean)
                .join('\n');
            if (ast.type === 'TUPLE' && ast.spreadParam) {
                type += declareNamedInterfaces(ast.spreadParam, options, rootASTName, processed);
            }
            break;
        default:
            type = '';
    }
    return type;
}
function declareNamedTypes(ast, options, rootASTName, processed) {
    if (processed === void 0) { processed = new Set(); }
    if (processed.has(ast)) {
        return '';
    }
    processed.add(ast);
    switch (ast.type) {
        case 'ARRAY':
            return [
                declareNamedTypes(ast.params, options, rootASTName, processed),
                (0, AST_1.hasStandaloneName)(ast) ? generateStandaloneType(ast, options) : undefined
            ]
                .filter(Boolean)
                .join('\n');
        case 'ENUM':
            return '';
        case 'INTERFACE':
            return getSuperTypesAndParams(ast)
                .map(function (ast) {
                return (ast.standaloneName === rootASTName || options.declareExternallyReferenced) &&
                    declareNamedTypes(ast, options, rootASTName, processed);
            })
                .filter(Boolean)
                .join('\n');
        case 'INTERSECTION':
        case 'TUPLE':
        case 'UNION':
            return [
                (0, AST_1.hasStandaloneName)(ast) ? generateStandaloneType(ast, options) : undefined,
                ast.params
                    .map(function (ast) { return declareNamedTypes(ast, options, rootASTName, processed); })
                    .filter(Boolean)
                    .join('\n'),
                'spreadParam' in ast && ast.spreadParam
                    ? declareNamedTypes(ast.spreadParam, options, rootASTName, processed)
                    : undefined
            ]
                .filter(Boolean)
                .join('\n');
        default:
            if ((0, AST_1.hasStandaloneName)(ast)) {
                return generateStandaloneType(ast, options);
            }
            return '';
    }
}
function generateTypeUnmemoized(ast, options) {
    var type = generateRawType(ast, options);
    if (options.strictIndexSignatures && ast.keyName === '[k: string]') {
        return "".concat(type, " | undefined");
    }
    return type;
}
exports.generateType = (0, lodash_1.memoize)(generateTypeUnmemoized);
function generateRawType(ast, options) {
    (0, utils_1.log)('magenta', 'generator', ast);
    if ((0, AST_1.hasStandaloneName)(ast)) {
        return (0, utils_1.toSafeString)(ast.standaloneName);
    }
    switch (ast.type) {
        case 'ANY':
            return 'any';
        case 'ARRAY':
            return (function () {
                var type = (0, exports.generateType)(ast.params, options);
                return type.endsWith('"') ? '(' + type + ')[]' : type + '[]';
            })();
        case 'BOOLEAN':
            return 'boolean';
        case 'INTERFACE':
            return generateInterface(ast, options);
        case 'INTERSECTION':
            return generateSetOperation(ast, options);
        case 'LITERAL':
            return JSON.stringify(ast.params);
        case 'NUMBER':
            return 'number';
        case 'NULL':
            return 'null';
        case 'OBJECT':
            return 'object';
        case 'REFERENCE':
            return ast.params;
        case 'STRING':
            return 'string';
        case 'TUPLE':
            return (function () {
                var minItems = ast.minItems;
                var maxItems = ast.maxItems || -1;
                var spreadParam = ast.spreadParam;
                var astParams = __spreadArray([], ast.params, true);
                if (minItems > 0 && minItems > astParams.length && ast.spreadParam === undefined) {
                    // this is a valid state, and JSONSchema doesn't care about the item type
                    if (maxItems < 0) {
                        // no max items and no spread param, so just spread any
                        spreadParam = options.unknownAny ? AST_1.T_UNKNOWN : AST_1.T_ANY;
                    }
                }
                if (maxItems > astParams.length && ast.spreadParam === undefined) {
                    // this is a valid state, and JSONSchema doesn't care about the item type
                    // fill the tuple with any elements
                    for (var i = astParams.length; i < maxItems; i += 1) {
                        astParams.push(options.unknownAny ? AST_1.T_UNKNOWN : AST_1.T_ANY);
                    }
                }
                function addSpreadParam(params) {
                    if (spreadParam) {
                        var spread = '...(' + (0, exports.generateType)(spreadParam, options) + ')[]';
                        params.push(spread);
                    }
                    return params;
                }
                function paramsToString(params) {
                    return '[' + params.join(', ') + ']';
                }
                var paramsList = astParams.map(function (param) { return (0, exports.generateType)(param, options); });
                if (paramsList.length > minItems) {
                    /*
                  if there are more items than the min, we return a union of tuples instead of
                  using the optional element operator. This is done because it is more typesafe.
          
                  // optional element operator
                  type A = [string, string?, string?]
                  const a: A = ['a', undefined, 'c'] // no error
          
                  // union of tuples
                  type B = [string] | [string, string] | [string, string, string]
                  const b: B = ['a', undefined, 'c'] // TS error
                  */
                    var cumulativeParamsList = paramsList.slice(0, minItems);
                    var typesToUnion = [];
                    if (cumulativeParamsList.length > 0) {
                        // actually has minItems, so add the initial state
                        typesToUnion.push(paramsToString(cumulativeParamsList));
                    }
                    else {
                        // no minItems means it's acceptable to have an empty tuple type
                        typesToUnion.push(paramsToString([]));
                    }
                    for (var i = minItems; i < paramsList.length; i += 1) {
                        cumulativeParamsList.push(paramsList[i]);
                        if (i === paramsList.length - 1) {
                            // only the last item in the union should have the spread parameter
                            addSpreadParam(cumulativeParamsList);
                        }
                        typesToUnion.push(paramsToString(cumulativeParamsList));
                    }
                    return typesToUnion.join('|');
                }
                // no max items so only need to return one type
                return paramsToString(addSpreadParam(paramsList));
            })();
        case 'UNION':
            return generateSetOperation(ast, options);
        case 'UNKNOWN':
            return 'unknown';
        case 'CUSTOM_TYPE':
            return ast.params;
    }
}
/**
 * Generate a Union or Intersection
 */
function generateSetOperation(ast, options) {
    var members = ast.params.map(function (_) { return (0, exports.generateType)(_, options); });
    var separator = ast.type === 'UNION' ? '|' : '&';
    return members.length === 1 ? members[0] : '(' + members.join(' ' + separator + ' ') + ')';
}
function generateInterface(ast, options) {
    return ("{" +
        '\n' +
        ast.params
            .filter(function (_) { return !_.isPatternProperty && !_.isUnreachableDefinition; })
            .map(function (_a) {
            var isRequired = _a.isRequired, keyName = _a.keyName, ast = _a.ast;
            return [isRequired, keyName, ast, (0, exports.generateType)(ast, options)];
        })
            .map(function (_a) {
            var isRequired = _a[0], keyName = _a[1], ast = _a[2], type = _a[3];
            return ((0, AST_1.hasComment)(ast) && !ast.standaloneName ? generateComment(ast.comment) + '\n' : '') +
                escapeKeyName(keyName) +
                (isRequired ? '' : '?') +
                ': ' +
                ((0, AST_1.hasStandaloneName)(ast) ? (0, utils_1.toSafeString)(type) : type);
        })
            .join('\n') +
        '\n' +
        '}');
}
function generateComment(comment) {
    return __spreadArray(__spreadArray(['/**'], comment.split('\n').map(function (_) { return ' * ' + _; }), true), [' */'], false).join('\n');
}
function generateStandaloneEnum(ast, options) {
    return (((0, AST_1.hasComment)(ast) ? generateComment(ast.comment) + '\n' : '') +
        'export ' +
        (options.enableConstEnums ? 'const ' : '') +
        "enum ".concat((0, utils_1.toSafeString)(ast.standaloneName), " {") +
        '\n' +
        ast.params.map(function (_a) {
            var ast = _a.ast, keyName = _a.keyName;
            return keyName + ' = ' + (0, exports.generateType)(ast, options);
        }).join(',\n') +
        '\n' +
        '}');
}
function generateStandaloneInterface(ast, options) {
    return (((0, AST_1.hasComment)(ast) ? generateComment(ast.comment) + '\n' : '') +
        "export interface ".concat((0, utils_1.toSafeString)(ast.standaloneName), " ") +
        (ast.superTypes.length > 0
            ? "extends ".concat(ast.superTypes.map(function (superType) { return (0, utils_1.toSafeString)(superType.standaloneName); }).join(', '), " ")
            : '') +
        generateInterface(ast, options));
}
function generateStandaloneType(ast, options) {
    return (((0, AST_1.hasComment)(ast) ? generateComment(ast.comment) + '\n' : '') +
        "export type ".concat((0, utils_1.toSafeString)(ast.standaloneName), " = ").concat((0, exports.generateType)((0, lodash_1.omit)(ast, 'standaloneName') /* TODO */, options)));
}
function escapeKeyName(keyName) {
    if (keyName.length && /[A-Za-z_$]/.test(keyName.charAt(0)) && /^[\w$]+$/.test(keyName)) {
        return keyName;
    }
    if (keyName === '[k: string]') {
        return keyName;
    }
    return JSON.stringify(keyName);
}
function getSuperTypesAndParams(ast) {
    return ast.params.map(function (param) { return param.ast; }).concat(ast.superTypes);
}
//# sourceMappingURL=generator.js.map