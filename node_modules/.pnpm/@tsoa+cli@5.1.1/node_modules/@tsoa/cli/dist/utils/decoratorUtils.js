"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.getProduces = exports.getPath = exports.isDecorator = exports.getSecurites = exports.getDecoratorValues = exports.getNodeFirstDecoratorValue = exports.getNodeFirstDecoratorName = exports.getDecorators = void 0;
const ts = require("typescript");
const initializer_value_1 = require("../metadataGeneration/initializer-value");
function tsHasDecorators(ts) {
    return typeof ts.canHaveDecorators === 'function';
}
function getDecorators(node, isMatching) {
    // beginning in ts4.8 node.decorator is undefined, use getDecorators instead.
    const decorators = tsHasDecorators(ts) && ts.canHaveDecorators(node) ? ts.getDecorators(node) : node.decorators;
    if (!decorators || !decorators.length) {
        return [];
    }
    return decorators
        .map((e) => {
        while (e.expression !== undefined) {
            e = e.expression;
        }
        return e;
    })
        .filter(isMatching);
}
exports.getDecorators = getDecorators;
function getNodeFirstDecoratorName(node, isMatching) {
    const decorators = getDecorators(node, isMatching);
    if (!decorators || !decorators.length) {
        return;
    }
    return decorators[0].text;
}
exports.getNodeFirstDecoratorName = getNodeFirstDecoratorName;
function getNodeFirstDecoratorValue(node, typeChecker, isMatching) {
    const decorators = getDecorators(node, isMatching);
    if (!decorators || !decorators.length) {
        return;
    }
    const values = getDecoratorValues(decorators[0], typeChecker);
    return values && values[0];
}
exports.getNodeFirstDecoratorValue = getNodeFirstDecoratorValue;
function getDecoratorValues(decorator, typeChecker) {
    const expression = decorator.parent;
    const expArguments = expression.arguments;
    if (!expArguments || !expArguments.length) {
        return [];
    }
    return expArguments.map(a => (0, initializer_value_1.getInitializerValue)(a, typeChecker));
}
exports.getDecoratorValues = getDecoratorValues;
function getSecurites(decorator, typeChecker) {
    const [first, second] = getDecoratorValues(decorator, typeChecker);
    if (isObject(first)) {
        return first;
    }
    return { [first]: second || [] };
}
exports.getSecurites = getSecurites;
function isDecorator(node, isMatching) {
    const decorators = getDecorators(node, isMatching);
    if (!decorators || !decorators.length) {
        return false;
    }
    return true;
}
exports.isDecorator = isDecorator;
function isObject(v) {
    return typeof v === 'object' && v !== null;
}
function getPath(decorator, typeChecker) {
    const [path] = getDecoratorValues(decorator, typeChecker);
    if (path === undefined) {
        return '';
    }
    return path;
}
exports.getPath = getPath;
function getProduces(node, typeChecker) {
    const producesDecorators = getDecorators(node, identifier => identifier.text === 'Produces');
    if (!producesDecorators || !producesDecorators.length) {
        return [];
    }
    return producesDecorators.map(decorator => getDecoratorValues(decorator, typeChecker)[0]);
}
exports.getProduces = getProduces;
//# sourceMappingURL=decoratorUtils.js.map