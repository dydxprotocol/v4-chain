"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SpecGenerator = void 0;
const runtime_1 = require("@tsoa/runtime");
const handlebars = require("handlebars");
class SpecGenerator {
    constructor(metadata, config) {
        this.metadata = metadata;
        this.config = config;
        this.determineImplicitAdditionalPropertiesValue = () => {
            if (this.config.noImplicitAdditionalProperties === 'silently-remove-extras') {
                return false;
            }
            else if (this.config.noImplicitAdditionalProperties === 'throw-on-extras') {
                return false;
            }
            else if (this.config.noImplicitAdditionalProperties === 'ignore') {
                return true;
            }
            else {
                return (0, runtime_1.assertNever)(this.config.noImplicitAdditionalProperties);
            }
        };
    }
    buildAdditionalProperties(type) {
        return this.getSwaggerType(type);
    }
    buildOperationIdTemplate(inlineTemplate) {
        handlebars.registerHelper('titleCase', (value) => (value ? value.charAt(0).toUpperCase() + value.slice(1) : value));
        handlebars.registerHelper('replace', (subject, searchValue, withValue = '') => (subject ? subject.replace(searchValue, withValue) : subject));
        return handlebars.compile(inlineTemplate, { noEscape: true });
    }
    getOperationId(controllerName, method) {
        var _a;
        return this.buildOperationIdTemplate((_a = this.config.operationIdTemplate) !== null && _a !== void 0 ? _a : '{{titleCase method.name}}')({
            method,
            controllerName,
        });
    }
    throwIfNotDataFormat(strToTest) {
        const guiltyUntilInnocent = strToTest;
        if (guiltyUntilInnocent === 'int32' ||
            guiltyUntilInnocent === 'int64' ||
            guiltyUntilInnocent === 'float' ||
            guiltyUntilInnocent === 'double' ||
            guiltyUntilInnocent === 'byte' ||
            guiltyUntilInnocent === 'binary' ||
            guiltyUntilInnocent === 'date' ||
            guiltyUntilInnocent === 'date-time' ||
            guiltyUntilInnocent === 'password') {
            return guiltyUntilInnocent;
        }
        else {
            return (0, runtime_1.assertNever)(guiltyUntilInnocent);
        }
    }
    throwIfNotDataType(strToTest) {
        const guiltyUntilInnocent = strToTest;
        if (guiltyUntilInnocent === 'array' ||
            guiltyUntilInnocent === 'boolean' ||
            guiltyUntilInnocent === 'integer' ||
            guiltyUntilInnocent === 'file' ||
            guiltyUntilInnocent === 'number' ||
            guiltyUntilInnocent === 'object' ||
            guiltyUntilInnocent === 'string' ||
            guiltyUntilInnocent === 'undefined') {
            return guiltyUntilInnocent;
        }
        else {
            return (0, runtime_1.assertNever)(guiltyUntilInnocent);
        }
    }
    getSwaggerType(type, title) {
        if (type.dataType === 'void' || type.dataType === 'undefined') {
            return this.getSwaggerTypeForVoid(type.dataType);
        }
        else if (type.dataType === 'refEnum' || type.dataType === 'refObject' || type.dataType === 'refAlias') {
            return this.getSwaggerTypeForReferenceType(type);
        }
        else if (type.dataType === 'any' ||
            type.dataType === 'binary' ||
            type.dataType === 'boolean' ||
            type.dataType === 'buffer' ||
            type.dataType === 'byte' ||
            type.dataType === 'date' ||
            type.dataType === 'datetime' ||
            type.dataType === 'double' ||
            type.dataType === 'float' ||
            type.dataType === 'file' ||
            type.dataType === 'integer' ||
            type.dataType === 'long' ||
            type.dataType === 'object' ||
            type.dataType === 'string') {
            return this.getSwaggerTypeForPrimitiveType(type.dataType);
        }
        else if (type.dataType === 'array') {
            return this.getSwaggerTypeForArrayType(type, title);
        }
        else if (type.dataType === 'enum') {
            return this.getSwaggerTypeForEnumType(type, title);
        }
        else if (type.dataType === 'union') {
            return this.getSwaggerTypeForUnionType(type, title);
        }
        else if (type.dataType === 'intersection') {
            return this.getSwaggerTypeForIntersectionType(type, title);
        }
        else if (type.dataType === 'nestedObjectLiteral') {
            return this.getSwaggerTypeForObjectLiteral(type, title);
        }
        else {
            return (0, runtime_1.assertNever)(type);
        }
    }
    getSwaggerTypeForObjectLiteral(objectLiteral, title) {
        const properties = this.buildProperties(objectLiteral.properties);
        const additionalProperties = objectLiteral.additionalProperties && this.getSwaggerType(objectLiteral.additionalProperties);
        const required = objectLiteral.properties.filter(prop => prop.required && !this.hasUndefined(prop)).map(prop => prop.name);
        // An empty list required: [] is not valid.
        // If all properties are optional, do not specify the required keyword.
        return {
            ...(title && { title }),
            properties,
            ...(additionalProperties && { additionalProperties }),
            ...(required && required.length && { required }),
            type: 'object',
        };
    }
    getSwaggerTypeForReferenceType(_referenceType) {
        return {
        // Don't set additionalProperties value here since it will be set within the $ref's model when that $ref gets created
        };
    }
    getSwaggerTypeForVoid(_dataType) {
        // Described here: https://swagger.io/docs/specification/describing-responses/#empty
        const voidSchema = {
        // isn't allowed to have additionalProperties at all (meaning not a boolean or object)
        };
        return voidSchema;
    }
    getSwaggerTypeForPrimitiveType(dataType) {
        if (dataType === 'object') {
            if (process.env.NODE_ENV !== 'tsoa_test') {
                // eslint-disable-next-line no-console
                console.warn(`The type Object is discouraged. Please consider using an interface such as:
          export interface IStringToStringDictionary {
            [key: string]: string;
          }
          // or
          export interface IRecordOfAny {
            [key: string]: any;
          }
        `);
            }
        }
        const map = {
            any: {
                // While the any type is discouraged, it does explicitly allows anything, so it should always allow additionalProperties
                additionalProperties: true,
            },
            binary: { type: 'string', format: 'binary' },
            boolean: { type: 'boolean' },
            buffer: { type: 'string', format: 'byte' },
            byte: { type: 'string', format: 'byte' },
            date: { type: 'string', format: 'date' },
            datetime: { type: 'string', format: 'date-time' },
            double: { type: 'number', format: 'double' },
            file: { type: 'file' },
            float: { type: 'number', format: 'float' },
            integer: { type: 'integer', format: 'int32' },
            long: { type: 'integer', format: 'int64' },
            object: {
                additionalProperties: this.determineImplicitAdditionalPropertiesValue(),
                type: 'object',
            },
            string: { type: 'string' },
        };
        return map[dataType];
    }
    getSwaggerTypeForArrayType(arrayType, title) {
        return {
            items: this.getSwaggerType(arrayType.elementType, title),
            type: 'array',
        };
    }
    determineTypesUsedInEnum(anEnum) {
        const typesUsedInEnum = anEnum.reduce((theSet, curr) => {
            const typeUsed = curr === null ? 'number' : typeof curr;
            theSet.add(typeUsed);
            return theSet;
        }, new Set());
        return typesUsedInEnum;
    }
    hasUndefined(property) {
        return property.type.dataType === 'undefined' || (property.type.dataType === 'union' && property.type.types.some(type => type.dataType === 'undefined'));
    }
    queriesPropertyToQueryParameter(property) {
        return {
            parameterName: property.name,
            example: [property.example],
            description: property.description,
            in: 'query',
            name: property.name,
            required: property.required,
            type: property.type,
            default: property.default,
            validators: property.validators,
            deprecated: property.deprecated,
        };
    }
}
exports.SpecGenerator = SpecGenerator;
//# sourceMappingURL=specGenerator.js.map