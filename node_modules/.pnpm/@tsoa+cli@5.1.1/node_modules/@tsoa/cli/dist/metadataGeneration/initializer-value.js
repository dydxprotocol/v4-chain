"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.getInitializerValue = void 0;
const ts = require("typescript");
const hasInitializer = (node) => Object.prototype.hasOwnProperty.call(node, 'initializer');
const extractInitializer = (decl) => (decl && hasInitializer(decl) && decl.initializer) || undefined;
const extractImportSpecifier = (symbol) => ((symbol === null || symbol === void 0 ? void 0 : symbol.declarations) && symbol.declarations.length > 0 && ts.isImportSpecifier(symbol.declarations[0]) && symbol.declarations[0]) || undefined;
const getInitializerValue = (initializer, typeChecker, type) => {
    if (!initializer || !typeChecker) {
        return;
    }
    switch (initializer.kind) {
        case ts.SyntaxKind.ArrayLiteralExpression: {
            const arrayLiteral = initializer;
            return arrayLiteral.elements.map(element => (0, exports.getInitializerValue)(element, typeChecker));
        }
        case ts.SyntaxKind.StringLiteral:
        case ts.SyntaxKind.NoSubstitutionTemplateLiteral:
            return initializer.text;
        case ts.SyntaxKind.TrueKeyword:
            return true;
        case ts.SyntaxKind.FalseKeyword:
            return false;
        case ts.SyntaxKind.PrefixUnaryExpression: {
            const prefixUnary = initializer;
            switch (prefixUnary.operator) {
                case ts.SyntaxKind.PlusToken:
                    return Number(prefixUnary.operand.text);
                case ts.SyntaxKind.MinusToken:
                    return Number(`-${prefixUnary.operand.text}`);
                default:
                    throw new Error(`Unsupport prefix operator token: ${prefixUnary.operator}`);
            }
        }
        case ts.SyntaxKind.NumberKeyword:
        case ts.SyntaxKind.FirstLiteralToken:
            return Number(initializer.text);
        case ts.SyntaxKind.NewExpression: {
            const newExpression = initializer;
            const ident = newExpression.expression;
            if (ident.text === 'Date') {
                let date = new Date();
                if (newExpression.arguments) {
                    const newArguments = newExpression.arguments.filter(args => args.kind !== undefined);
                    const argsValue = newArguments.map(args => (0, exports.getInitializerValue)(args, typeChecker));
                    if (argsValue.length > 0) {
                        date = new Date(argsValue);
                    }
                }
                const dateString = date.toISOString();
                if (type && type.dataType === 'date') {
                    return dateString.split('T')[0];
                }
                return dateString;
            }
            return;
        }
        case ts.SyntaxKind.NullKeyword:
            return null;
        case ts.SyntaxKind.ObjectLiteralExpression: {
            const objectLiteral = initializer;
            const nestedObject = {};
            objectLiteral.properties.forEach((p) => {
                nestedObject[p.name.text] = (0, exports.getInitializerValue)(p.initializer, typeChecker);
            });
            return nestedObject;
        }
        case ts.SyntaxKind.ImportSpecifier: {
            const importSpecifier = initializer;
            const importSymbol = typeChecker.getSymbolAtLocation(importSpecifier.name);
            if (!importSymbol)
                return;
            const aliasedSymbol = typeChecker.getAliasedSymbol(importSymbol);
            const declarations = aliasedSymbol.getDeclarations();
            const declaration = declarations && declarations.length > 0 ? declarations[0] : undefined;
            return (0, exports.getInitializerValue)(extractInitializer(declaration), typeChecker);
        }
        default: {
            const symbol = typeChecker.getSymbolAtLocation(initializer);
            return (0, exports.getInitializerValue)(extractInitializer(symbol === null || symbol === void 0 ? void 0 : symbol.valueDeclaration) || extractImportSpecifier(symbol), typeChecker);
        }
    }
};
exports.getInitializerValue = getInitializerValue;
//# sourceMappingURL=initializer-value.js.map