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
exports.optimize = void 0;
var lodash_1 = require("lodash");
var generator_1 = require("./generator");
var AST_1 = require("./types/AST");
var utils_1 = require("./utils");
function optimize(ast, options, processed) {
    if (processed === void 0) { processed = new Set(); }
    if (processed.has(ast)) {
        return ast;
    }
    processed.add(ast);
    switch (ast.type) {
        case 'INTERFACE':
            return Object.assign(ast, {
                params: ast.params.map(function (_) { return Object.assign(_, { ast: optimize(_.ast, options, processed) }); })
            });
        case 'INTERSECTION':
        case 'UNION':
            // Start with the leaves...
            var optimizedAST_1 = Object.assign(ast, {
                params: ast.params.map(function (_) { return optimize(_, options, processed); })
            });
            // [A, B, C, Any] -> Any
            if (optimizedAST_1.params.some(function (_) { return _.type === 'ANY'; })) {
                (0, utils_1.log)('cyan', 'optimizer', '[A, B, C, Any] -> Any', optimizedAST_1);
                return AST_1.T_ANY;
            }
            // [A, B, C, Unknown] -> Unknown
            if (optimizedAST_1.params.some(function (_) { return _.type === 'UNKNOWN'; })) {
                (0, utils_1.log)('cyan', 'optimizer', '[A, B, C, Unknown] -> Unknown', optimizedAST_1);
                return AST_1.T_UNKNOWN;
            }
            // [A (named), A] -> [A (named)]
            if (optimizedAST_1.params.every(function (_) {
                var a = (0, generator_1.generateType)(omitStandaloneName(_), options);
                var b = (0, generator_1.generateType)(omitStandaloneName(optimizedAST_1.params[0]), options);
                return a === b;
            }) &&
                optimizedAST_1.params.some(function (_) { return _.standaloneName !== undefined; })) {
                (0, utils_1.log)('cyan', 'optimizer', '[A (named), A] -> [A (named)]', optimizedAST_1);
                optimizedAST_1.params = optimizedAST_1.params.filter(function (_) { return _.standaloneName !== undefined; });
            }
            // [A, B, B] -> [A, B]
            var params = (0, lodash_1.uniqBy)(optimizedAST_1.params, function (_) { return (0, generator_1.generateType)(_, options); });
            if (params.length !== optimizedAST_1.params.length) {
                (0, utils_1.log)('cyan', 'optimizer', '[A, B, B] -> [A, B]', optimizedAST_1);
                optimizedAST_1.params = params;
            }
            return Object.assign(optimizedAST_1, {
                params: optimizedAST_1.params.map(function (_) { return optimize(_, options, processed); })
            });
        default:
            return ast;
    }
}
exports.optimize = optimize;
// TODO: More clearly disambiguate standalone names vs. aliased names instead.
function omitStandaloneName(ast) {
    switch (ast.type) {
        case 'ENUM':
            return ast;
        default:
            return __assign(__assign({}, ast), { standaloneName: undefined });
    }
}
//# sourceMappingURL=optimizer.js.map