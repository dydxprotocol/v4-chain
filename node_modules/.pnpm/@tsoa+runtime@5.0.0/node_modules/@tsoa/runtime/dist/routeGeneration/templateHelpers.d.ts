/// <reference types="node" />
import { AdditionalProps } from './additionalProps';
import { TsoaRoute } from './tsoa-route';
export declare function ValidateParam(property: TsoaRoute.PropertySchema, value: any, generatedModels: TsoaRoute.Models, name: string | undefined, fieldErrors: FieldErrors, parent: string | undefined, swaggerConfig: AdditionalProps): any;
export declare class ValidationService {
    private readonly models;
    constructor(models: TsoaRoute.Models);
    ValidateParam(property: TsoaRoute.PropertySchema, rawValue: any, name: string | undefined, fieldErrors: FieldErrors, parent: string | undefined, minimalSwaggerConfig: AdditionalProps): any;
    validateNestedObjectLiteral(name: string, value: any, fieldErrors: FieldErrors, swaggerConfig: AdditionalProps, nestedProperties: {
        [name: string]: TsoaRoute.PropertySchema;
    } | undefined, additionalProperties: TsoaRoute.PropertySchema | boolean | undefined, parent: string): any;
    validateInt(name: string, value: any, fieldErrors: FieldErrors, validators?: IntegerValidator, parent?: string): number | undefined;
    validateFloat(name: string, value: any, fieldErrors: FieldErrors, validators?: FloatValidator, parent?: string): number | undefined;
    validateEnum(name: string, value: unknown, fieldErrors: FieldErrors, members?: Array<string | number | boolean | null>, parent?: string): unknown;
    validateDate(name: string, value: any, fieldErrors: FieldErrors, validators?: DateValidator, parent?: string): Date | undefined;
    validateDateTime(name: string, value: any, fieldErrors: FieldErrors, validators?: DateTimeValidator, parent?: string): Date | undefined;
    validateString(name: string, value: any, fieldErrors: FieldErrors, validators?: StringValidator, parent?: string): string | undefined;
    validateBool(name: string, value: any, fieldErrors: FieldErrors, validators?: BooleanValidator, parent?: string): any;
    validateUndefined(name: string, value: any, fieldErrors: FieldErrors, parent?: string): undefined;
    validateArray(name: string, value: any[], fieldErrors: FieldErrors, swaggerConfig: AdditionalProps, schema?: TsoaRoute.PropertySchema, validators?: ArrayValidator, parent?: string): any[] | undefined;
    validateBuffer(_name: string, value: string): Buffer;
    validateUnion(name: string, value: any, fieldErrors: FieldErrors, swaggerConfig: AdditionalProps, property: TsoaRoute.PropertySchema, parent?: string): any;
    validateIntersection(name: string, value: any, fieldErrors: FieldErrors, swaggerConfig: AdditionalProps, subSchemas: TsoaRoute.PropertySchema[] | undefined, parent?: string): any;
    private toModelLike;
    /**
     * combine all schemas once without backwards combinations ie
     * input: [[value1], [value2]] should be [[value1, value2]]
     * not [[value1, value2],[value2, value1]]
     * and
     * input: [[value1], [value2], [value3]] should be [
     *   [value1, value2, value3],
     *   [value1, value2],
     *   [value1, value3],
     *   [value2, value3]
     * ]
     * @param modelSchemass
     */
    private selfIntersectionExcludingCombinations;
    private intersectRefObjectModelSchemas;
    private combineProperties;
    private getExcessPropertiesFor;
    validateModel(input: {
        name: string;
        value: any;
        modelDefinition: TsoaRoute.ModelSchema;
        fieldErrors: FieldErrors;
        parent?: string;
        minimalSwaggerConfig: AdditionalProps;
    }): any;
}
export interface IntegerValidator {
    isInt?: {
        errorMsg?: string;
    };
    isLong?: {
        errorMsg?: string;
    };
    minimum?: {
        value: number;
        errorMsg?: string;
    };
    maximum?: {
        value: number;
        errorMsg?: string;
    };
}
export interface FloatValidator {
    isFloat?: {
        errorMsg?: string;
    };
    isDouble?: {
        errorMsg?: string;
    };
    minimum?: {
        value: number;
        errorMsg?: string;
    };
    maximum?: {
        value: number;
        errorMsg?: string;
    };
}
export interface DateValidator {
    isDate?: {
        errorMsg?: string;
    };
    minDate?: {
        value: string;
        errorMsg?: string;
    };
    maxDate?: {
        value: string;
        errorMsg?: string;
    };
}
export interface DateTimeValidator {
    isDateTime?: {
        errorMsg?: string;
    };
    minDate?: {
        value: string;
        errorMsg?: string;
    };
    maxDate?: {
        value: string;
        errorMsg?: string;
    };
}
export interface StringValidator {
    isString?: {
        errorMsg?: string;
    };
    minLength?: {
        value: number;
        errorMsg?: string;
    };
    maxLength?: {
        value: number;
        errorMsg?: string;
    };
    pattern?: {
        value: string;
        errorMsg?: string;
    };
}
export interface BooleanValidator {
    isArray?: {
        errorMsg?: string;
    };
}
export interface ArrayValidator {
    isArray?: {
        errorMsg?: string;
    };
    minItems?: {
        value: number;
        errorMsg?: string;
    };
    maxItems?: {
        value: number;
        errorMsg?: string;
    };
    uniqueItems?: {
        errorMsg?: string;
    };
}
export type Validator = IntegerValidator | FloatValidator | DateValidator | DateTimeValidator | StringValidator | BooleanValidator | ArrayValidator;
export interface FieldErrors {
    [name: string]: {
        message: string;
        value?: any;
    };
}
export interface Exception extends Error {
    status: number;
}
export declare class ValidateError extends Error implements Exception {
    fields: FieldErrors;
    message: string;
    status: number;
    name: string;
    constructor(fields: FieldErrors, message: string);
}
