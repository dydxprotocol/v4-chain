"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SpecGenerator3 = void 0;
const runtime_1 = require("@tsoa/runtime");
const isVoidType_1 = require("../utils/isVoidType");
const pathUtils_1 = require("./../utils/pathUtils");
const swaggerUtils_1 = require("./../utils/swaggerUtils");
const specGenerator_1 = require("./specGenerator");
/**
 * TODO:
 * Handle formData parameters
 * Handle requestBodies of type other than json
 * Handle requestBodies as reusable objects
 * Handle headers, examples, responses, etc.
 * Cleaner interface between SpecGenerator2 and SpecGenerator3
 * Also accept OpenAPI 3.0.0 metadata, like components/securitySchemes instead of securityDefinitions
 */
class SpecGenerator3 extends specGenerator_1.SpecGenerator {
    constructor(metadata, config) {
        super(metadata, config);
        this.metadata = metadata;
        this.config = config;
    }
    GetSpec() {
        let spec = {
            components: this.buildComponents(),
            info: this.buildInfo(),
            openapi: '3.0.0',
            paths: this.buildPaths(),
            servers: this.buildServers(),
            tags: this.config.tags,
        };
        if (this.config.spec) {
            this.config.specMerging = this.config.specMerging || 'immediate';
            const mergeFuncs = {
                immediate: Object.assign,
                recursive: require('merge').recursive,
                deepmerge: (spec, merge) => require('deepmerge').all([spec, merge]),
            };
            spec = mergeFuncs[this.config.specMerging](spec, this.config.spec);
        }
        return spec;
    }
    buildInfo() {
        const info = {
            title: this.config.name || '',
        };
        if (this.config.version) {
            info.version = this.config.version;
        }
        if (this.config.description) {
            info.description = this.config.description;
        }
        if (this.config.termsOfService) {
            info.termsOfService = this.config.termsOfService;
        }
        if (this.config.license) {
            info.license = { name: this.config.license };
        }
        if (this.config.contact) {
            info.contact = this.config.contact;
        }
        return info;
    }
    buildComponents() {
        const components = {
            examples: {},
            headers: {},
            parameters: {},
            requestBodies: {},
            responses: {},
            schemas: this.buildSchema(),
            securitySchemes: {},
        };
        if (this.config.securityDefinitions) {
            components.securitySchemes = this.translateSecurityDefinitions(this.config.securityDefinitions);
        }
        return components;
    }
    translateSecurityDefinitions(definitions) {
        const defs = {};
        Object.keys(definitions).forEach(key => {
            if (definitions[key].type === 'basic') {
                defs[key] = {
                    scheme: 'basic',
                    type: 'http',
                };
            }
            else if (definitions[key].type === 'oauth2') {
                const definition = definitions[key];
                const oauth = (defs[key] || {
                    type: 'oauth2',
                    description: definitions[key].description,
                    flows: (this.hasOAuthFlows(definition) && definition.flows) || {},
                });
                if (this.hasOAuthFlow(definition) && definition.flow === 'password') {
                    oauth.flows.password = { tokenUrl: definition.tokenUrl, scopes: definition.scopes || {} };
                }
                else if (this.hasOAuthFlow(definition) && definition.flow === 'accessCode') {
                    oauth.flows.authorizationCode = { tokenUrl: definition.tokenUrl, authorizationUrl: definition.authorizationUrl, scopes: definition.scopes || {} };
                }
                else if (this.hasOAuthFlow(definition) && definition.flow === 'application') {
                    oauth.flows.clientCredentials = { tokenUrl: definition.tokenUrl, scopes: definition.scopes || {} };
                }
                else if (this.hasOAuthFlow(definition) && definition.flow === 'implicit') {
                    oauth.flows.implicit = { authorizationUrl: definition.authorizationUrl, scopes: definition.scopes || {} };
                }
                defs[key] = oauth;
            }
            else {
                defs[key] = definitions[key];
            }
        });
        return defs;
    }
    hasOAuthFlow(definition) {
        return !!definition.flow;
    }
    hasOAuthFlows(definition) {
        return !!definition.flows;
    }
    buildServers() {
        const basePath = (0, pathUtils_1.normalisePath)(this.config.basePath, '/', undefined, false);
        const scheme = this.config.schemes ? this.config.schemes[0] : 'https';
        const url = this.config.host ? `${scheme}://${this.config.host}${basePath}` : basePath;
        return [
            {
                url,
            },
        ];
    }
    buildSchema() {
        const schema = {};
        Object.keys(this.metadata.referenceTypeMap).map(typeName => {
            const referenceType = this.metadata.referenceTypeMap[typeName];
            if (referenceType.dataType === 'refObject') {
                const required = referenceType.properties.filter(p => p.required && !this.hasUndefined(p)).map(p => p.name);
                schema[referenceType.refName] = {
                    description: referenceType.description,
                    properties: this.buildProperties(referenceType.properties),
                    required: required && required.length > 0 ? Array.from(new Set(required)) : undefined,
                    type: 'object',
                };
                if (referenceType.additionalProperties) {
                    schema[referenceType.refName].additionalProperties = this.buildAdditionalProperties(referenceType.additionalProperties);
                }
                else {
                    // Since additionalProperties was not explicitly set in the TypeScript interface for this model
                    //      ...we need to make a decision
                    schema[referenceType.refName].additionalProperties = this.determineImplicitAdditionalPropertiesValue();
                }
                if (referenceType.example) {
                    schema[referenceType.refName].example = referenceType.example;
                }
            }
            else if (referenceType.dataType === 'refEnum') {
                const enumTypes = this.determineTypesUsedInEnum(referenceType.enums);
                if (enumTypes.size === 1) {
                    schema[referenceType.refName] = {
                        description: referenceType.description,
                        enum: referenceType.enums,
                        type: enumTypes.has('string') ? 'string' : 'number',
                    };
                    if (this.config.xEnumVarnames && referenceType.enumVarnames !== undefined && referenceType.enums.length === referenceType.enumVarnames.length) {
                        schema[referenceType.refName]['x-enum-varnames'] = referenceType.enumVarnames;
                    }
                }
                else {
                    schema[referenceType.refName] = {
                        description: referenceType.description,
                        anyOf: [
                            {
                                type: 'number',
                                enum: referenceType.enums.filter(e => typeof e === 'number'),
                            },
                            {
                                type: 'string',
                                enum: referenceType.enums.filter(e => typeof e === 'string'),
                            },
                        ],
                    };
                }
            }
            else if (referenceType.dataType === 'refAlias') {
                const swaggerType = this.getSwaggerType(referenceType.type);
                const format = referenceType.format;
                const validators = Object.keys(referenceType.validators)
                    .filter(key => {
                    return !key.startsWith('is') && key !== 'minDate' && key !== 'maxDate';
                })
                    .reduce((acc, key) => {
                    return {
                        ...acc,
                        [key]: referenceType.validators[key].value,
                    };
                }, {});
                schema[referenceType.refName] = {
                    ...swaggerType,
                    default: referenceType.default || swaggerType.default,
                    example: referenceType.example,
                    format: format || swaggerType.format,
                    description: referenceType.description,
                    ...validators,
                };
            }
            else {
                (0, runtime_1.assertNever)(referenceType);
            }
            if (referenceType.deprecated) {
                schema[referenceType.refName].deprecated = true;
            }
        });
        return schema;
    }
    buildPaths() {
        const paths = {};
        this.metadata.controllers.forEach(controller => {
            const normalisedControllerPath = (0, pathUtils_1.normalisePath)(controller.path, '/');
            // construct documentation using all methods except @Hidden
            controller.methods
                .filter(method => !method.isHidden)
                .forEach(method => {
                const normalisedMethodPath = (0, pathUtils_1.normalisePath)(method.path, '/');
                let path = (0, pathUtils_1.normalisePath)(`${normalisedControllerPath}${normalisedMethodPath}`, '/', '', false);
                path = (0, pathUtils_1.convertColonPathParams)(path);
                paths[path] = paths[path] || {};
                this.buildMethod(controller.name, method, paths[path], controller.produces);
            });
        });
        return paths;
    }
    buildMethod(controllerName, method, pathObject, defaultProduces) {
        const pathMethod = (pathObject[method.method] = this.buildOperation(controllerName, method, defaultProduces));
        pathMethod.description = method.description;
        pathMethod.summary = method.summary;
        pathMethod.tags = method.tags;
        // Use operationId tag otherwise fallback to generated. Warning: This doesn't check uniqueness.
        pathMethod.operationId = method.operationId || pathMethod.operationId;
        if (method.deprecated) {
            pathMethod.deprecated = method.deprecated;
        }
        if (method.security) {
            pathMethod.security = method.security;
        }
        const bodyParams = method.parameters.filter(p => p.in === 'body');
        const bodyPropParams = method.parameters.filter(p => p.in === 'body-prop');
        const formParams = method.parameters.filter(p => p.in === 'formData');
        const queriesParams = method.parameters.filter(p => p.in === 'queries');
        pathMethod.parameters = method.parameters
            .filter(p => {
            return ['body', 'formData', 'request', 'body-prop', 'res', 'queries'].indexOf(p.in) === -1;
        })
            .map(p => this.buildParameter(p));
        if (queriesParams.length > 1) {
            throw new Error('Only one queries parameter allowed per controller method.');
        }
        if (queriesParams.length === 1) {
            pathMethod.parameters.push(...this.buildQueriesParameter(queriesParams[0]));
        }
        if (bodyParams.length > 1) {
            throw new Error('Only one body parameter allowed per controller method.');
        }
        if (bodyParams.length > 0 && formParams.length > 0) {
            throw new Error('Either body parameter or form parameters allowed per controller method - not both.');
        }
        if (bodyPropParams.length > 0) {
            if (!bodyParams.length) {
                bodyParams.push({
                    in: 'body',
                    name: 'body',
                    parameterName: 'body',
                    required: true,
                    type: {
                        dataType: 'nestedObjectLiteral',
                        properties: [],
                    },
                    validators: {},
                    deprecated: false,
                });
            }
            const type = bodyParams[0].type;
            bodyPropParams.forEach((bodyParam) => {
                type.properties.push(bodyParam);
            });
        }
        if (bodyParams.length > 0) {
            pathMethod.requestBody = this.buildRequestBody(controllerName, method, bodyParams[0]);
        }
        else if (formParams.length > 0) {
            pathMethod.requestBody = this.buildRequestBodyWithFormData(controllerName, method, formParams);
        }
        method.extensions.forEach(ext => (pathMethod[ext.key] = ext.value));
    }
    buildOperation(controllerName, method, defaultProduces) {
        const swaggerResponses = {};
        method.responses.forEach((res) => {
            swaggerResponses[res.name] = {
                description: res.description,
            };
            if (res.schema && !(0, isVoidType_1.isVoidType)(res.schema)) {
                swaggerResponses[res.name].content = {};
                const produces = res.produces || defaultProduces || [swaggerUtils_1.DEFAULT_RESPONSE_MEDIA_TYPE];
                for (const p of produces) {
                    const { content } = swaggerResponses[res.name];
                    swaggerResponses[res.name].content = {
                        ...content,
                        [p]: {
                            schema: this.getSwaggerType(res.schema, this.config.useTitleTagsForInlineObjects ? this.getOperationId(controllerName, method) + 'Response' : undefined),
                        },
                    };
                }
                if (res.examples) {
                    let exampleCounter = 1;
                    const examples = res.examples.reduce((acc, ex, currentIndex) => {
                        var _a;
                        const exampleLabel = (_a = res.exampleLabels) === null || _a === void 0 ? void 0 : _a[currentIndex];
                        return { ...acc, [exampleLabel === undefined ? `Example ${exampleCounter++}` : exampleLabel]: { value: ex } };
                    }, {});
                    for (const p of produces) {
                        /* eslint-disable @typescript-eslint/dot-notation */
                        (swaggerResponses[res.name].content || {})[p]['examples'] = examples;
                    }
                }
            }
            if (res.headers) {
                const headers = {};
                if (res.headers.dataType === 'refObject') {
                    headers[res.headers.refName] = {
                        schema: this.getSwaggerTypeForReferenceType(res.headers),
                        description: res.headers.description,
                    };
                }
                else if (res.headers.dataType === 'nestedObjectLiteral') {
                    res.headers.properties.forEach((each) => {
                        headers[each.name] = {
                            schema: this.getSwaggerType(each.type),
                            description: each.description,
                            required: each.required,
                        };
                    });
                }
                else {
                    (0, runtime_1.assertNever)(res.headers);
                }
                swaggerResponses[res.name].headers = headers;
            }
        });
        const operation = {
            operationId: this.getOperationId(controllerName, method),
            responses: swaggerResponses,
        };
        return operation;
    }
    buildRequestBodyWithFormData(controllerName, method, parameters) {
        const required = [];
        const properties = {};
        for (const parameter of parameters) {
            const mediaType = this.buildMediaType(controllerName, method, parameter);
            properties[parameter.name] = mediaType.schema;
            if (parameter.required) {
                required.push(parameter.name);
            }
            if (parameter.deprecated) {
                properties[parameter.name].deprecated = parameter.deprecated;
            }
        }
        const requestBody = {
            required: required.length > 0,
            content: {
                'multipart/form-data': {
                    schema: {
                        type: 'object',
                        properties,
                        // An empty list required: [] is not valid.
                        // If all properties are optional, do not specify the required keyword.
                        ...(required && required.length && { required }),
                    },
                },
            },
        };
        return requestBody;
    }
    buildRequestBody(controllerName, method, parameter) {
        const mediaType = this.buildMediaType(controllerName, method, parameter);
        const consumes = method.consumes || swaggerUtils_1.DEFAULT_REQUEST_MEDIA_TYPE;
        const requestBody = {
            description: parameter.description,
            required: parameter.required,
            content: {
                [consumes]: mediaType,
            },
        };
        return requestBody;
    }
    buildMediaType(controllerName, method, parameter) {
        const validators = Object.keys(parameter.validators)
            .filter(key => {
            return !key.startsWith('is') && key !== 'minDate' && key !== 'maxDate';
        })
            .reduce((acc, key) => {
            return {
                ...acc,
                [key]: validators[key].value,
            };
        }, {});
        const mediaType = {
            schema: {
                ...this.getSwaggerType(parameter.type, this.config.useTitleTagsForInlineObjects ? this.getOperationId(controllerName, method) + 'RequestBody' : undefined),
                ...validators,
                ...(parameter.description && { description: parameter.description }),
            },
        };
        const parameterExamples = parameter.example;
        const parameterExampleLabels = parameter.exampleLabels;
        if (parameterExamples === undefined) {
            mediaType.example = parameterExamples;
        }
        else if (parameterExamples.length === 1) {
            mediaType.example = parameterExamples[0];
        }
        else {
            let exampleCounter = 1;
            mediaType.examples = parameterExamples.reduce((acc, ex, currentIndex) => {
                const exampleLabel = parameterExampleLabels === null || parameterExampleLabels === void 0 ? void 0 : parameterExampleLabels[currentIndex];
                return { ...acc, [exampleLabel === undefined ? `Example ${exampleCounter++}` : exampleLabel]: { value: ex } };
            }, {});
        }
        return mediaType;
    }
    buildQueriesParameter(source) {
        if (source.type.dataType === 'refObject' || source.type.dataType === 'nestedObjectLiteral') {
            const properties = source.type.properties;
            return properties.map(property => this.buildParameter(this.queriesPropertyToQueryParameter(property)));
        }
        throw new Error(`Queries '${source.name}' parameter must be an object.`);
    }
    buildParameter(source) {
        const parameter = {
            description: source.description,
            in: source.in,
            name: source.name,
            required: source.required,
            schema: {
                default: source.default,
                format: undefined,
            },
        };
        if (source.deprecated) {
            parameter.deprecated = true;
        }
        const parameterType = this.getSwaggerType(source.type);
        if (parameterType.format) {
            parameter.schema.format = this.throwIfNotDataFormat(parameterType.format);
        }
        if (parameterType.$ref) {
            parameter.schema = parameterType;
            return parameter;
        }
        const validatorObjs = {};
        Object.keys(source.validators)
            .filter(key => {
            return !key.startsWith('is') && key !== 'minDate' && key !== 'maxDate';
        })
            .forEach((key) => {
            validatorObjs[key] = source.validators[key].value;
        });
        if (source.type.dataType === 'any') {
            parameter.schema.type = 'string';
        }
        else {
            if (parameterType.type) {
                parameter.schema.type = this.throwIfNotDataType(parameterType.type);
            }
            parameter.schema.items = parameterType.items;
            parameter.schema.enum = parameterType.enum;
        }
        parameter.schema = Object.assign({}, parameter.schema, validatorObjs);
        const parameterExamples = source.example;
        const parameterExampleLabels = source.exampleLabels;
        if (parameterExamples === undefined) {
            parameter.example = parameterExamples;
        }
        else if (parameterExamples.length === 1) {
            parameter.example = parameterExamples[0];
        }
        else {
            let exampleCounter = 1;
            parameter.examples = parameterExamples.reduce((acc, ex, currentIndex) => {
                const exampleLabel = parameterExampleLabels === null || parameterExampleLabels === void 0 ? void 0 : parameterExampleLabels[currentIndex];
                return { ...acc, [exampleLabel === undefined ? `Example ${exampleCounter++}` : exampleLabel]: { value: ex } };
            }, {});
        }
        return parameter;
    }
    buildProperties(source) {
        const properties = {};
        source.forEach(property => {
            const swaggerType = this.getSwaggerType(property.type);
            const format = property.format;
            swaggerType.description = property.description;
            swaggerType.example = property.example;
            swaggerType.format = format || swaggerType.format;
            if (!swaggerType.$ref) {
                swaggerType.default = property.default;
                Object.keys(property.validators)
                    .filter(key => {
                    return !key.startsWith('is') && key !== 'minDate' && key !== 'maxDate';
                })
                    .forEach(key => {
                    swaggerType[key] = property.validators[key].value;
                });
            }
            if (property.deprecated) {
                swaggerType.deprecated = true;
            }
            if (property.extensions) {
                property.extensions.forEach(property => {
                    swaggerType[property.key] = property.value;
                });
            }
            properties[property.name] = swaggerType;
        });
        return properties;
    }
    getSwaggerTypeForReferenceType(referenceType) {
        return { $ref: `#/components/schemas/${referenceType.refName}` };
    }
    getSwaggerTypeForPrimitiveType(dataType) {
        if (dataType === 'any') {
            // Setting additionalProperties causes issues with code generators for OpenAPI 3
            // Therefore, we avoid setting it explicitly (since it's the implicit default already)
            return {};
        }
        else if (dataType === 'file') {
            return { type: 'string', format: 'binary' };
        }
        return super.getSwaggerTypeForPrimitiveType(dataType);
    }
    isNull(type) {
        return type.dataType === 'enum' && type.enums.length === 1 && type.enums[0] === null;
    }
    // Join disparate enums with the same type into one.
    //
    // grouping enums is helpful because it makes the spec more readable and it
    // bypasses a failure in openapi-generator caused by using anyOf with
    // duplicate types.
    groupEnums(types) {
        const returnTypes = [];
        const enumValuesByType = {};
        for (const type of types) {
            if (type.enum && type.type) {
                for (const enumValue of type.enum) {
                    if (!enumValuesByType[type.type]) {
                        enumValuesByType[type.type] = [];
                    }
                    enumValuesByType[type.type][enumValue] = enumValue;
                }
            }
            // preserve non-enum types
            else {
                returnTypes.push(type);
            }
        }
        Object.keys(enumValuesByType).forEach(dataType => returnTypes.push({
            type: dataType,
            enum: Object.values(enumValuesByType[dataType]),
        }));
        return returnTypes;
    }
    removeDuplicateSwaggerTypes(types) {
        if (types.length === 1) {
            return types;
        }
        else {
            const typesSet = new Set();
            for (const type of types) {
                typesSet.add(JSON.stringify(type));
            }
            return Array.from(typesSet).map(typeString => JSON.parse(typeString));
        }
    }
    getSwaggerTypeForUnionType(type, title) {
        // Filter out nulls and undefineds
        const actualSwaggerTypes = this.removeDuplicateSwaggerTypes(this.groupEnums(type.types
            .filter(x => !this.isNull(x))
            .filter(x => x.dataType !== 'undefined')
            .map(x => this.getSwaggerType(x))));
        const nullable = type.types.some(x => this.isNull(x));
        if (nullable) {
            if (actualSwaggerTypes.length === 1) {
                const [swaggerType] = actualSwaggerTypes;
                // for ref union with null, use an allOf with a single
                // element since you can't attach nullable directly to a ref.
                // https://swagger.io/docs/specification/using-ref/#syntax
                if (swaggerType.$ref) {
                    return { allOf: [swaggerType], nullable };
                }
                return { ...(title && { title }), ...swaggerType, nullable };
            }
            else {
                return { ...(title && { title }), anyOf: actualSwaggerTypes, nullable };
            }
        }
        else {
            if (actualSwaggerTypes.length === 1) {
                return { ...(title && { title }), ...actualSwaggerTypes[0] };
            }
            else {
                return { ...(title && { title }), anyOf: actualSwaggerTypes };
            }
        }
    }
    getSwaggerTypeForIntersectionType(type, title) {
        return { allOf: type.types.map(x => this.getSwaggerType(x)), ...(title && { title }) };
    }
    getSwaggerTypeForEnumType(enumType, title) {
        const types = this.determineTypesUsedInEnum(enumType.enums);
        if (types.size === 1) {
            const type = types.values().next().value;
            const nullable = enumType.enums.includes(null) ? true : false;
            return { ...(title && { title }), type, enum: enumType.enums.map(member => (0, swaggerUtils_1.getValue)(type, member)), nullable };
        }
        else {
            const valuesDelimited = Array.from(types).join(',');
            throw new Error(`Enums can only have string or number values, but enum had ${valuesDelimited}`);
        }
    }
}
exports.SpecGenerator3 = SpecGenerator3;
//# sourceMappingURL=specGenerator3.js.map