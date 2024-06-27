"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.getExtensionsFromJSDocComments = exports.getExtensions = void 0;
const ts = require("typescript");
const initializer_value_1 = require("./initializer-value");
const jsonUtils_1 = require("../utils/jsonUtils");
function getExtensions(decorators, metadataGenerator) {
    const extensions = decorators.map(extensionDecorator => {
        if (!ts.isCallExpression(extensionDecorator.parent)) {
            throw new Error('The parent of the @Extension is not a CallExpression. Are you using it in the right place?');
        }
        const [decoratorKeyArg, decoratorValueArg] = extensionDecorator.parent.arguments;
        if (!ts.isStringLiteral(decoratorKeyArg)) {
            throw new Error('The first argument of @Extension must be a string');
        }
        const attributeKey = decoratorKeyArg.text;
        if (!decoratorValueArg) {
            throw new Error(`Extension '${attributeKey}' must contain a value`);
        }
        validateExtensionKey(attributeKey);
        const attributeValue = (0, initializer_value_1.getInitializerValue)(decoratorValueArg, metadataGenerator.typeChecker);
        return { key: attributeKey, value: attributeValue };
    });
    return extensions;
}
exports.getExtensions = getExtensions;
function getExtensionsFromJSDocComments(comments) {
    const extensions = [];
    comments.forEach(comment => {
        const extensionData = (0, jsonUtils_1.safeFromJson)(comment);
        if (extensionData) {
            const keys = Object.keys(extensionData);
            keys.forEach(key => {
                validateExtensionKey(key);
                extensions.push({ key: key, value: extensionData[key] });
            });
        }
    });
    return extensions;
}
exports.getExtensionsFromJSDocComments = getExtensionsFromJSDocComments;
function validateExtensionKey(key) {
    if (key.indexOf('x-') !== 0) {
        throw new Error('Extensions must begin with "x-" to be valid. Please see the following link for more information: https://swagger.io/docs/specification/openapi-extensions/');
    }
}
//# sourceMappingURL=extension.js.map