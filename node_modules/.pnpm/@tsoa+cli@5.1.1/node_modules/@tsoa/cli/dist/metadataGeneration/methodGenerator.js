"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.MethodGenerator = void 0;
const ts = require("typescript");
const path = require("path");
const isVoidType_1 = require("../utils/isVoidType");
const decoratorUtils_1 = require("./../utils/decoratorUtils");
const jsDocUtils_1 = require("./../utils/jsDocUtils");
const extension_1 = require("./extension");
const exceptions_1 = require("./exceptions");
const parameterGenerator_1 = require("./parameterGenerator");
const typeResolver_1 = require("./typeResolver");
const headerTypeHelpers_1 = require("../utils/headerTypeHelpers");
class MethodGenerator {
    constructor(node, current, commonResponses, parentPath, parentTags, parentSecurity, isParentHidden) {
        this.node = node;
        this.current = current;
        this.commonResponses = commonResponses;
        this.parentPath = parentPath;
        this.parentTags = parentTags;
        this.parentSecurity = parentSecurity;
        this.isParentHidden = isParentHidden;
        this.processMethodDecorators();
    }
    IsValid() {
        return !!this.method;
    }
    Generate() {
        if (!this.IsValid()) {
            throw new exceptions_1.GenerateMetadataError("This isn't a valid a controller method.");
        }
        let nodeType = this.node.type;
        if (!nodeType) {
            const typeChecker = this.current.typeChecker;
            const signature = typeChecker.getSignatureFromDeclaration(this.node);
            const implicitType = typeChecker.getReturnTypeOfSignature(signature);
            nodeType = typeChecker.typeToTypeNode(implicitType, undefined, ts.NodeBuilderFlags.NoTruncation);
        }
        const type = new typeResolver_1.TypeResolver(nodeType, this.current).resolve();
        const responses = this.commonResponses.concat(this.getMethodResponses());
        const { response: successResponse, status: successStatus } = this.getMethodSuccessResponse(type);
        responses.push(successResponse);
        const parameters = this.buildParameters();
        const additionalResponses = parameters.filter((p) => p.in === 'res');
        responses.push(...additionalResponses);
        return {
            extensions: this.getExtensions(),
            deprecated: this.getIsDeprecated(),
            description: (0, jsDocUtils_1.getJSDocDescription)(this.node),
            isHidden: this.getIsHidden(),
            method: this.method,
            name: this.node.name.text,
            operationId: this.getOperationId(),
            parameters,
            path: this.path,
            produces: this.produces,
            consumes: this.consumes,
            responses,
            successStatus,
            security: this.getSecurity(),
            summary: (0, jsDocUtils_1.getJSDocComment)(this.node, 'summary'),
            tags: this.getTags(),
            type,
        };
    }
    buildParameters() {
        const fullPath = path.join(this.parentPath || '', this.path);
        const parameters = this.node.parameters
            .map(p => {
            try {
                return new parameterGenerator_1.ParameterGenerator(p, this.method, fullPath, this.current).Generate();
            }
            catch (e) {
                const methodId = this.node.name;
                const controllerId = this.node.parent.name;
                throw new exceptions_1.GenerateMetadataError(`${String(e.message)} \n in '${controllerId.text}.${methodId.text}'`);
            }
        })
            .reduce((flattened, params) => [...flattened, ...params], []);
        this.validateBodyParameters(parameters);
        this.validateQueryParameters(parameters);
        return parameters;
    }
    validateBodyParameters(parameters) {
        const bodyParameters = parameters.filter(p => p.in === 'body');
        const bodyProps = parameters.filter(p => p.in === 'body-prop');
        const hasFormDataParameters = parameters.some(p => p.in === 'formData');
        const hasBodyParameter = bodyProps.length + bodyParameters.length > 0;
        if (bodyParameters.length > 1) {
            throw new exceptions_1.GenerateMetadataError(`Only one body parameter allowed in '${this.getCurrentLocation()}' method.`);
        }
        if (bodyParameters.length > 0 && bodyProps.length > 0) {
            throw new exceptions_1.GenerateMetadataError(`Choose either during @Body or @BodyProp in '${this.getCurrentLocation()}' method.`);
        }
        if (hasBodyParameter && hasFormDataParameters) {
            throw new Error(`@Body or @BodyProp cannot be used with @FormField, @UploadedFile, or @UploadedFiles in '${this.getCurrentLocation()}' method.`);
        }
    }
    validateQueryParameters(parameters) {
        const queryParameters = parameters.filter(p => p.in === 'query');
        const queriesParameters = parameters.filter(p => p.in === 'queries');
        if (queriesParameters.length > 1) {
            throw new exceptions_1.GenerateMetadataError(`Only one queries parameter allowed in '${this.getCurrentLocation()}' method.`);
        }
        if (queriesParameters.length > 0 && queryParameters.length > 0) {
            throw new exceptions_1.GenerateMetadataError(`Choose either during @Query or @Queries in '${this.getCurrentLocation()}' method.`);
        }
    }
    getExtensions() {
        const extensionDecorators = this.getDecoratorsByIdentifier(this.node, 'Extension');
        if (!extensionDecorators || !extensionDecorators.length) {
            return [];
        }
        return (0, extension_1.getExtensions)(extensionDecorators, this.current);
    }
    getCurrentLocation() {
        const methodId = this.node.name;
        const controllerId = this.node.parent.name;
        return `${controllerId.text}.${methodId.text}`;
    }
    processMethodDecorators() {
        const pathDecorators = (0, decoratorUtils_1.getDecorators)(this.node, identifier => this.supportsPathMethod(identifier.text));
        if (!pathDecorators || !pathDecorators.length) {
            return;
        }
        if (pathDecorators.length > 1) {
            throw new exceptions_1.GenerateMetadataError(`Only one path decorator in '${this.getCurrentLocation()}' method, Found: ${pathDecorators.map(d => d.text).join(', ')}`);
        }
        const decorator = pathDecorators[0];
        this.method = decorator.text.toLowerCase();
        // if you don't pass in a path to the method decorator, we'll just use the base route
        // todo: what if someone has multiple no argument methods of the same type in a single controller?
        // we need to throw an error there
        this.path = (0, decoratorUtils_1.getPath)(decorator, this.current.typeChecker);
        this.produces = this.getProduces();
        this.consumes = this.getConsumes();
    }
    getProduces() {
        const produces = (0, decoratorUtils_1.getProduces)(this.node, this.current.typeChecker);
        return produces.length ? produces : undefined;
    }
    getConsumes() {
        const consumesDecorators = this.getDecoratorsByIdentifier(this.node, 'Consumes');
        if (!consumesDecorators || !consumesDecorators.length) {
            return;
        }
        if (consumesDecorators.length > 1) {
            throw new exceptions_1.GenerateMetadataError(`Only one Consumes decorator in '${this.getCurrentLocation()}' method, Found: ${consumesDecorators.map(d => d.text).join(', ')}`);
        }
        const [decorator] = consumesDecorators;
        const [consumes] = (0, decoratorUtils_1.getDecoratorValues)(decorator, this.current.typeChecker);
        return consumes;
    }
    getMethodResponses() {
        const decorators = this.getDecoratorsByIdentifier(this.node, 'Response');
        if (!decorators || !decorators.length) {
            return [];
        }
        return decorators.map(decorator => {
            const [name, description, example, produces] = (0, decoratorUtils_1.getDecoratorValues)(decorator, this.current.typeChecker);
            return {
                description: description || '',
                examples: example === undefined ? undefined : [example],
                name: name || '200',
                produces: this.getProducesAdapter(produces),
                schema: this.getSchemaFromDecorator(decorator, 0),
                headers: this.getHeadersFromDecorator(decorator, 1),
            };
        });
    }
    getMethodSuccessResponse(type) {
        const decorators = this.getDecoratorsByIdentifier(this.node, 'SuccessResponse');
        const examplesWithLabels = this.getMethodSuccessExamples();
        if (!decorators || !decorators.length) {
            const returnsDescription = (0, jsDocUtils_1.getJSDocComment)(this.node, 'returns') || 'Ok';
            return {
                response: {
                    description: (0, isVoidType_1.isVoidType)(type) ? 'No content' : returnsDescription,
                    examples: examplesWithLabels === null || examplesWithLabels === void 0 ? void 0 : examplesWithLabels.map(ex => ex.example),
                    exampleLabels: examplesWithLabels === null || examplesWithLabels === void 0 ? void 0 : examplesWithLabels.map(ex => ex.label),
                    name: (0, isVoidType_1.isVoidType)(type) ? '204' : '200',
                    produces: this.produces,
                    schema: type,
                },
            };
        }
        if (decorators.length > 1) {
            throw new exceptions_1.GenerateMetadataError(`Only one SuccessResponse decorator allowed in '${this.getCurrentLocation()}' method.`);
        }
        const [firstDecorator] = decorators;
        const [name, description, produces] = (0, decoratorUtils_1.getDecoratorValues)(firstDecorator, this.current.typeChecker);
        const headers = this.getHeadersFromDecorator(firstDecorator, 0);
        return {
            response: {
                description: description || '',
                examples: examplesWithLabels === null || examplesWithLabels === void 0 ? void 0 : examplesWithLabels.map(ex => ex.example),
                exampleLabels: examplesWithLabels === null || examplesWithLabels === void 0 ? void 0 : examplesWithLabels.map(ex => ex.label),
                name: name || '200',
                produces: this.getProducesAdapter(produces),
                schema: type,
                headers,
            },
            status: name && /^\d+$/.test(name) ? parseInt(name, 10) : undefined,
        };
    }
    getHeadersFromDecorator({ parent: expression }, headersIndex) {
        if (!ts.isCallExpression(expression)) {
            return undefined;
        }
        return (0, headerTypeHelpers_1.getHeaderType)(expression.typeArguments, headersIndex, this.current);
    }
    getSchemaFromDecorator({ parent: expression }, schemaIndex) {
        var _a;
        if (!ts.isCallExpression(expression) || !((_a = expression.typeArguments) === null || _a === void 0 ? void 0 : _a.length)) {
            return undefined;
        }
        return new typeResolver_1.TypeResolver(expression.typeArguments[schemaIndex], this.current).resolve();
    }
    getMethodSuccessExamples() {
        const exampleDecorators = this.getDecoratorsByIdentifier(this.node, 'Example');
        if (!exampleDecorators || !exampleDecorators.length) {
            return undefined;
        }
        const examples = exampleDecorators.map(exampleDecorator => {
            const [example, label] = (0, decoratorUtils_1.getDecoratorValues)(exampleDecorator, this.current.typeChecker);
            return { example, label };
        });
        return examples || undefined;
    }
    supportsPathMethod(method) {
        return ['options', 'get', 'post', 'put', 'patch', 'delete', 'head'].some(m => m === method.toLowerCase());
    }
    getIsDeprecated() {
        if ((0, jsDocUtils_1.isExistJSDocTag)(this.node, tag => tag.tagName.text === 'deprecated')) {
            return true;
        }
        const depDecorators = this.getDecoratorsByIdentifier(this.node, 'Deprecated');
        if (!depDecorators || !depDecorators.length) {
            return false;
        }
        if (depDecorators.length > 1) {
            throw new exceptions_1.GenerateMetadataError(`Only one Deprecated decorator allowed in '${this.getCurrentLocation()}' method.`);
        }
        return true;
    }
    getOperationId() {
        const opDecorators = this.getDecoratorsByIdentifier(this.node, 'OperationId');
        if (!opDecorators || !opDecorators.length) {
            return undefined;
        }
        if (opDecorators.length > 1) {
            throw new exceptions_1.GenerateMetadataError(`Only one OperationId decorator allowed in '${this.getCurrentLocation()}' method.`);
        }
        const values = (0, decoratorUtils_1.getDecoratorValues)(opDecorators[0], this.current.typeChecker);
        return values && values[0];
    }
    getTags() {
        const tagsDecorators = this.getDecoratorsByIdentifier(this.node, 'Tags');
        if (!tagsDecorators || !tagsDecorators.length) {
            return this.parentTags;
        }
        if (tagsDecorators.length > 1) {
            throw new exceptions_1.GenerateMetadataError(`Only one Tags decorator allowed in '${this.getCurrentLocation()}' method.`);
        }
        const tags = (0, decoratorUtils_1.getDecoratorValues)(tagsDecorators[0], this.current.typeChecker);
        if (tags && this.parentTags) {
            tags.push(...this.parentTags);
        }
        return tags;
    }
    getSecurity() {
        const noSecurityDecorators = this.getDecoratorsByIdentifier(this.node, 'NoSecurity');
        const securityDecorators = this.getDecoratorsByIdentifier(this.node, 'Security');
        if ((noSecurityDecorators === null || noSecurityDecorators === void 0 ? void 0 : noSecurityDecorators.length) > 1) {
            throw new exceptions_1.GenerateMetadataError(`Only one NoSecurity decorator allowed in '${this.getCurrentLocation()}' method.`);
        }
        if ((noSecurityDecorators === null || noSecurityDecorators === void 0 ? void 0 : noSecurityDecorators.length) && (securityDecorators === null || securityDecorators === void 0 ? void 0 : securityDecorators.length)) {
            throw new exceptions_1.GenerateMetadataError(`NoSecurity decorator cannot be used in conjunction with Security decorator in '${this.getCurrentLocation()}' method.`);
        }
        if (noSecurityDecorators === null || noSecurityDecorators === void 0 ? void 0 : noSecurityDecorators.length) {
            return [];
        }
        if (!securityDecorators || !securityDecorators.length) {
            return this.parentSecurity || [];
        }
        return securityDecorators.map(d => (0, decoratorUtils_1.getSecurites)(d, this.current.typeChecker));
    }
    getIsHidden() {
        const hiddenDecorators = this.getDecoratorsByIdentifier(this.node, 'Hidden');
        if (!hiddenDecorators || !hiddenDecorators.length) {
            return !!this.isParentHidden;
        }
        if (this.isParentHidden) {
            throw new exceptions_1.GenerateMetadataError(`Hidden decorator cannot be set on '${this.getCurrentLocation()}' it is already defined on the controller`);
        }
        if (hiddenDecorators.length > 1) {
            throw new exceptions_1.GenerateMetadataError(`Only one Hidden decorator allowed in '${this.getCurrentLocation()}' method.`);
        }
        return true;
    }
    getDecoratorsByIdentifier(node, id) {
        return (0, decoratorUtils_1.getDecorators)(node, identifier => identifier.text === id);
    }
    getProducesAdapter(produces) {
        if (Array.isArray(produces)) {
            return produces;
        }
        else if (typeof produces === 'string') {
            return [produces];
        }
        return;
    }
}
exports.MethodGenerator = MethodGenerator;
//# sourceMappingURL=methodGenerator.js.map