"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ParameterGenerator = void 0;
const ts = require("typescript");
const decoratorUtils_1 = require("./../utils/decoratorUtils");
const jsDocUtils_1 = require("./../utils/jsDocUtils");
const validatorUtils_1 = require("./../utils/validatorUtils");
const exceptions_1 = require("./exceptions");
const initializer_value_1 = require("./initializer-value");
const typeResolver_1 = require("./typeResolver");
const headerTypeHelpers_1 = require("../utils/headerTypeHelpers");
class ParameterGenerator {
    constructor(parameter, method, path, current) {
        this.parameter = parameter;
        this.method = method;
        this.path = path;
        this.current = current;
    }
    Generate() {
        const decoratorName = (0, decoratorUtils_1.getNodeFirstDecoratorName)(this.parameter, identifier => this.supportParameterDecorator(identifier.text));
        switch (decoratorName) {
            case 'Request':
                return [this.getRequestParameter(this.parameter)];
            case 'Body':
                return [this.getBodyParameter(this.parameter)];
            case 'BodyProp':
                return [this.getBodyPropParameter(this.parameter)];
            case 'FormField':
                return [this.getFormFieldParameter(this.parameter)];
            case 'Header':
                return [this.getHeaderParameter(this.parameter)];
            case 'Query':
                return this.getQueryParameters(this.parameter);
            case 'Queries':
                return [this.getQueriesParameters(this.parameter)];
            case 'Path':
                return [this.getPathParameter(this.parameter)];
            case 'Res':
                return this.getResParameters(this.parameter);
            case 'Inject':
                return [];
            case 'UploadedFile':
                return [this.getUploadedFileParameter(this.parameter)];
            case 'UploadedFiles':
                return [this.getUploadedFileParameter(this.parameter, true)];
            default:
                return [this.getPathParameter(this.parameter)];
        }
    }
    getRequestParameter(parameter) {
        const parameterName = parameter.name.text;
        return {
            description: this.getParameterDescription(parameter),
            in: 'request',
            name: parameterName,
            parameterName,
            required: !parameter.questionToken && !parameter.initializer,
            type: { dataType: 'object' },
            validators: (0, validatorUtils_1.getParameterValidators)(this.parameter, parameterName),
            deprecated: this.getParameterDeprecation(parameter),
        };
    }
    getResParameters(parameter) {
        const parameterName = parameter.name.text;
        const decorator = (0, decoratorUtils_1.getNodeFirstDecoratorValue)(this.parameter, this.current.typeChecker, ident => ident.text === 'Res') || parameterName;
        if (!decorator) {
            throw new exceptions_1.GenerateMetadataError('Could not find Decorator', parameter);
        }
        const typeNode = parameter.type;
        if (!typeNode || !ts.isTypeReferenceNode(typeNode) || typeNode.typeName.getText() !== 'TsoaResponse') {
            throw new exceptions_1.GenerateMetadataError('@Res() requires the type to be TsoaResponse<HTTPStatusCode, ResBody>', parameter);
        }
        if (!typeNode.typeArguments || !typeNode.typeArguments[0]) {
            throw new exceptions_1.GenerateMetadataError('@Res() requires the type to be TsoaResponse<HTTPStatusCode, ResBody>', parameter);
        }
        const statusArgument = typeNode.typeArguments[0];
        const bodyArgument = typeNode.typeArguments[1];
        // support a union of status codes, all with the same response body
        const statusArguments = ts.isUnionTypeNode(statusArgument) ? [...statusArgument.types] : [statusArgument];
        const statusArgumentTypes = statusArguments.map(a => this.current.typeChecker.getTypeAtLocation(a));
        const isNumberLiteralType = (tsType) => {
            // eslint-disable-next-line no-bitwise
            return (tsType.getFlags() & ts.TypeFlags.NumberLiteral) !== 0;
        };
        const headers = (0, headerTypeHelpers_1.getHeaderType)(typeNode.typeArguments, 2, this.current);
        return statusArgumentTypes.map(statusArgumentType => {
            if (!isNumberLiteralType(statusArgumentType)) {
                throw new exceptions_1.GenerateMetadataError('@Res() requires the type to be TsoaResponse<HTTPStatusCode, ResBody>', parameter);
            }
            const status = String(statusArgumentType.value);
            const type = new typeResolver_1.TypeResolver(bodyArgument, this.current, typeNode).resolve();
            const { examples, exampleLabels } = this.getParameterExample(parameter, parameterName);
            return {
                description: this.getParameterDescription(parameter) || '',
                in: 'res',
                name: status,
                produces: headers ? this.getProducesFromResHeaders(headers) : undefined,
                parameterName,
                examples,
                required: true,
                type,
                exampleLabels,
                schema: type,
                validators: {},
                headers,
                deprecated: this.getParameterDeprecation(parameter),
            };
        });
    }
    getProducesFromResHeaders(headers) {
        const { properties } = headers;
        const [contentTypeProp] = (properties || []).filter(p => p.name.toLowerCase() === 'content-type' && p.type.dataType === 'enum');
        if (contentTypeProp) {
            const type = contentTypeProp.type;
            return type.enums;
        }
        return;
    }
    getBodyPropParameter(parameter) {
        const parameterName = parameter.name.text;
        const type = this.getValidatedType(parameter);
        if (!this.supportBodyMethod(this.method)) {
            throw new exceptions_1.GenerateMetadataError(`@BodyProp('${parameterName}') Can't support in ${this.method.toUpperCase()} method.`);
        }
        const { examples: example, exampleLabels } = this.getParameterExample(parameter, parameterName);
        return {
            default: (0, initializer_value_1.getInitializerValue)(parameter.initializer, this.current.typeChecker, type),
            description: this.getParameterDescription(parameter),
            example,
            exampleLabels,
            in: 'body-prop',
            name: (0, decoratorUtils_1.getNodeFirstDecoratorValue)(this.parameter, this.current.typeChecker, ident => ident.text === 'BodyProp') || parameterName,
            parameterName,
            required: !parameter.questionToken && !parameter.initializer,
            type,
            validators: (0, validatorUtils_1.getParameterValidators)(this.parameter, parameterName),
            deprecated: this.getParameterDeprecation(parameter),
        };
    }
    getBodyParameter(parameter) {
        const parameterName = parameter.name.text;
        const type = this.getValidatedType(parameter);
        if (!this.supportBodyMethod(this.method)) {
            throw new exceptions_1.GenerateMetadataError(`@Body('${parameterName}') Can't support in ${this.method.toUpperCase()} method.`);
        }
        const { examples: example, exampleLabels } = this.getParameterExample(parameter, parameterName);
        return {
            description: this.getParameterDescription(parameter),
            in: 'body',
            name: parameterName,
            example,
            exampleLabels,
            parameterName,
            required: !parameter.questionToken && !parameter.initializer,
            type,
            validators: (0, validatorUtils_1.getParameterValidators)(this.parameter, parameterName),
            deprecated: this.getParameterDeprecation(parameter),
        };
    }
    getHeaderParameter(parameter) {
        const parameterName = parameter.name.text;
        const type = this.getValidatedType(parameter);
        if (!this.supportPathDataType(type)) {
            throw new exceptions_1.GenerateMetadataError(`@Header('${parameterName}') Can't support '${type.dataType}' type.`);
        }
        const { examples: example, exampleLabels } = this.getParameterExample(parameter, parameterName);
        return {
            default: (0, initializer_value_1.getInitializerValue)(parameter.initializer, this.current.typeChecker, type),
            description: this.getParameterDescription(parameter),
            example,
            exampleLabels,
            in: 'header',
            name: (0, decoratorUtils_1.getNodeFirstDecoratorValue)(this.parameter, this.current.typeChecker, ident => ident.text === 'Header') || parameterName,
            parameterName,
            required: !parameter.questionToken && !parameter.initializer,
            type,
            validators: (0, validatorUtils_1.getParameterValidators)(this.parameter, parameterName),
            deprecated: this.getParameterDeprecation(parameter),
        };
    }
    getUploadedFileParameter(parameter, isArray) {
        var _a;
        const parameterName = parameter.name.text;
        const elementType = { dataType: 'file' };
        let type;
        if (isArray) {
            type = { dataType: 'array', elementType };
        }
        else {
            type = elementType;
        }
        if (!this.supportPathDataType(elementType)) {
            throw new exceptions_1.GenerateMetadataError(`Parameter '${parameterName}:${type.dataType}' can't be passed as an uploaded file(s) parameter in '${this.method.toUpperCase()}'.`, parameter);
        }
        return {
            description: this.getParameterDescription(parameter),
            in: 'formData',
            name: (_a = (0, decoratorUtils_1.getNodeFirstDecoratorValue)(this.parameter, this.current.typeChecker, ident => {
                if (isArray) {
                    return ident.text === 'UploadedFiles';
                }
                return ident.text === 'UploadedFile';
            })) !== null && _a !== void 0 ? _a : parameterName,
            required: !parameter.questionToken && !parameter.initializer,
            type,
            parameterName,
            validators: (0, validatorUtils_1.getParameterValidators)(this.parameter, parameterName),
            deprecated: this.getParameterDeprecation(parameter),
        };
    }
    getFormFieldParameter(parameter) {
        var _a;
        const parameterName = parameter.name.text;
        const type = { dataType: 'string' };
        if (!this.supportPathDataType(type)) {
            throw new exceptions_1.GenerateMetadataError(`Parameter '${parameterName}:${type.dataType}' can't be passed as form field parameter in '${this.method.toUpperCase()}'.`, parameter);
        }
        return {
            description: this.getParameterDescription(parameter),
            in: 'formData',
            name: (_a = (0, decoratorUtils_1.getNodeFirstDecoratorValue)(this.parameter, this.current.typeChecker, ident => ident.text === 'FormField')) !== null && _a !== void 0 ? _a : parameterName,
            required: !parameter.questionToken && !parameter.initializer,
            type,
            parameterName,
            validators: (0, validatorUtils_1.getParameterValidators)(this.parameter, parameterName),
            deprecated: this.getParameterDeprecation(parameter),
        };
    }
    getQueriesParameters(parameter) {
        const parameterName = parameter.name.text;
        const type = this.getValidatedType(parameter);
        if (type.dataType !== 'refObject' && type.dataType !== 'nestedObjectLiteral') {
            throw new exceptions_1.GenerateMetadataError(`@Queries('${parameterName}') only support 'refObject' or 'nestedObjectLiteral' types. If you want only one query parameter, please use the '@Query' decorator.`);
        }
        for (const property of type.properties) {
            this.validateQueriesProperties(property, parameterName);
        }
        const { examples: example, exampleLabels } = this.getParameterExample(parameter, parameterName);
        return {
            description: this.getParameterDescription(parameter),
            in: 'queries',
            name: parameterName,
            example,
            exampleLabels,
            parameterName,
            required: !parameter.questionToken && !parameter.initializer,
            type,
            validators: (0, validatorUtils_1.getParameterValidators)(this.parameter, parameterName),
            deprecated: this.getParameterDeprecation(parameter),
        };
    }
    validateQueriesProperties(property, parentName) {
        if (property.type.dataType === 'array') {
            const arrayType = property.type;
            if (!this.supportPathDataType(arrayType.elementType)) {
                throw new exceptions_1.GenerateMetadataError(`@Queries('${parentName}') property '${property.name}' can't support array '${arrayType.elementType.dataType}' type.`);
            }
        }
        else if (!this.supportPathDataType(property.type)) {
            throw new exceptions_1.GenerateMetadataError(`@Queries('${parentName}') nested property '${property.name}' Can't support '${property.type.dataType}' type.`);
        }
    }
    getQueryParameters(parameter) {
        const parameterName = parameter.name.text;
        const type = this.getValidatedType(parameter);
        const { examples: example, exampleLabels } = this.getParameterExample(parameter, parameterName);
        const commonProperties = {
            default: (0, initializer_value_1.getInitializerValue)(parameter.initializer, this.current.typeChecker, type),
            description: this.getParameterDescription(parameter),
            example,
            exampleLabels,
            in: 'query',
            name: (0, decoratorUtils_1.getNodeFirstDecoratorValue)(this.parameter, this.current.typeChecker, ident => ident.text === 'Query') || parameterName,
            parameterName,
            required: !parameter.questionToken && !parameter.initializer,
            validators: (0, validatorUtils_1.getParameterValidators)(this.parameter, parameterName),
            deprecated: this.getParameterDeprecation(parameter),
        };
        if (this.getQueryParameterIsHidden(parameter)) {
            if (commonProperties.required) {
                throw new exceptions_1.GenerateMetadataError(`@Query('${parameterName}') Can't support @Hidden because it is required (does not allow undefined and does not have a default value).`);
            }
            return [];
        }
        if (type.dataType === 'array') {
            const arrayType = type;
            if (!this.supportPathDataType(arrayType.elementType)) {
                throw new exceptions_1.GenerateMetadataError(`@Query('${parameterName}') Can't support array '${arrayType.elementType.dataType}' type.`);
            }
            return [
                {
                    ...commonProperties,
                    collectionFormat: 'multi',
                    type: arrayType,
                },
            ];
        }
        if (!this.supportPathDataType(type)) {
            throw new exceptions_1.GenerateMetadataError(`@Query('${parameterName}') Can't support '${type.dataType}' type.`);
        }
        return [
            {
                ...commonProperties,
                type,
            },
        ];
    }
    getPathParameter(parameter) {
        const parameterName = parameter.name.text;
        const type = this.getValidatedType(parameter);
        const pathName = String((0, decoratorUtils_1.getNodeFirstDecoratorValue)(this.parameter, this.current.typeChecker, ident => ident.text === 'Path') || parameterName);
        if (!this.supportPathDataType(type)) {
            throw new exceptions_1.GenerateMetadataError(`@Path('${parameterName}') Can't support '${type.dataType}' type.`);
        }
        if (!this.path.includes(`{${pathName}}`) && !this.path.includes(`:${pathName}`)) {
            throw new exceptions_1.GenerateMetadataError(`@Path('${parameterName}') Can't match in URL: '${this.path}'.`);
        }
        const { examples, exampleLabels } = this.getParameterExample(parameter, parameterName);
        return {
            default: (0, initializer_value_1.getInitializerValue)(parameter.initializer, this.current.typeChecker, type),
            description: this.getParameterDescription(parameter),
            example: examples,
            exampleLabels,
            in: 'path',
            name: pathName,
            parameterName,
            required: true,
            type,
            validators: (0, validatorUtils_1.getParameterValidators)(this.parameter, parameterName),
            deprecated: this.getParameterDeprecation(parameter),
        };
    }
    getParameterDescription(node) {
        const symbol = this.current.typeChecker.getSymbolAtLocation(node.name);
        if (!symbol) {
            return undefined;
        }
        const comments = symbol.getDocumentationComment(this.current.typeChecker);
        if (comments.length) {
            return ts.displayPartsToString(comments);
        }
        return undefined;
    }
    getParameterDeprecation(node) {
        return (0, jsDocUtils_1.isExistJSDocTag)(node, tag => tag.tagName.text === 'deprecated') || (0, decoratorUtils_1.isDecorator)(node, identifier => identifier.text === 'Deprecated');
    }
    getParameterExample(node, parameterName) {
        const exampleLabels = [];
        const examples = (0, jsDocUtils_1.getJSDocTags)(node.parent, tag => {
            const comment = (0, jsDocUtils_1.commentToString)(tag.comment);
            const isExample = (tag.tagName.text === 'example' || tag.tagName.escapedText === 'example') && !!tag.comment && (comment === null || comment === void 0 ? void 0 : comment.startsWith(parameterName));
            if (isExample) {
                const hasExampleLabel = ((comment === null || comment === void 0 ? void 0 : comment.split(' ')[0].indexOf('.')) || -1) > 0;
                // custom example label is delimited by first '.' and the rest will all be included as example label
                exampleLabels.push(hasExampleLabel ? comment === null || comment === void 0 ? void 0 : comment.split(' ')[0].split('.').slice(1).join('.') : undefined);
            }
            return isExample !== null && isExample !== void 0 ? isExample : false;
        }).map(tag => { var _a; return ((0, jsDocUtils_1.commentToString)(tag.comment) || '').replace(`${((_a = (0, jsDocUtils_1.commentToString)(tag.comment)) === null || _a === void 0 ? void 0 : _a.split(' ')[0]) || ''}`, '').replace(/\r/g, ''); });
        if (examples.length === 0) {
            return {
                examples: undefined,
                exampleLabels: undefined,
            };
        }
        else {
            try {
                return {
                    examples: examples.map(example => JSON.parse(example)),
                    exampleLabels,
                };
            }
            catch (e) {
                throw new exceptions_1.GenerateMetadataError(`JSON format is incorrect: ${String(e.message)}`);
            }
        }
    }
    supportBodyMethod(method) {
        return ['post', 'put', 'patch', 'delete'].some(m => m === method.toLowerCase());
    }
    supportParameterDecorator(decoratorName) {
        return ['header', 'query', 'queries', 'path', 'body', 'bodyprop', 'request', 'res', 'inject', 'uploadedfile', 'uploadedfiles', 'formfield'].some(d => d === decoratorName.toLocaleLowerCase());
    }
    supportPathDataType(parameterType) {
        const supportedPathDataTypes = ['string', 'integer', 'long', 'float', 'double', 'date', 'datetime', 'buffer', 'boolean', 'enum', 'refEnum', 'file', 'any'];
        if (supportedPathDataTypes.find(t => t === parameterType.dataType)) {
            return true;
        }
        if (parameterType.dataType === 'refAlias') {
            return this.supportPathDataType(parameterType.type);
        }
        if (parameterType.dataType === 'union') {
            return !parameterType.types.map(t => this.supportPathDataType(t)).some(t => t === false);
        }
        return false;
    }
    getValidatedType(parameter) {
        let typeNode = parameter.type;
        if (!typeNode) {
            const type = this.current.typeChecker.getTypeAtLocation(parameter);
            typeNode = this.current.typeChecker.typeToTypeNode(type, undefined, ts.NodeBuilderFlags.NoTruncation);
        }
        return new typeResolver_1.TypeResolver(typeNode, this.current, parameter).resolve();
    }
    getQueryParameterIsHidden(parameter) {
        const hiddenDecorators = (0, decoratorUtils_1.getDecorators)(parameter, identifier => identifier.text === 'Hidden');
        if (!hiddenDecorators || !hiddenDecorators.length) {
            return false;
        }
        if (hiddenDecorators.length > 1) {
            const parameterName = parameter.name.text;
            throw new exceptions_1.GenerateMetadataError(`Only one Hidden decorator allowed on @Query('${parameterName}').`);
        }
        return true;
    }
}
exports.ParameterGenerator = ParameterGenerator;
//# sourceMappingURL=parameterGenerator.js.map