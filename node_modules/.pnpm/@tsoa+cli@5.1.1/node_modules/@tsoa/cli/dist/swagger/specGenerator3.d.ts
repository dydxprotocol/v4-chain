import { ExtendedSpecConfig } from '../cli';
import { Tsoa, Swagger } from '@tsoa/runtime';
import { SpecGenerator } from './specGenerator';
/**
 * TODO:
 * Handle formData parameters
 * Handle requestBodies of type other than json
 * Handle requestBodies as reusable objects
 * Handle headers, examples, responses, etc.
 * Cleaner interface between SpecGenerator2 and SpecGenerator3
 * Also accept OpenAPI 3.0.0 metadata, like components/securitySchemes instead of securityDefinitions
 */
export declare class SpecGenerator3 extends SpecGenerator {
    protected readonly metadata: Tsoa.Metadata;
    protected readonly config: ExtendedSpecConfig;
    constructor(metadata: Tsoa.Metadata, config: ExtendedSpecConfig);
    GetSpec(): Swagger.Spec3;
    private buildInfo;
    private buildComponents;
    private translateSecurityDefinitions;
    private hasOAuthFlow;
    private hasOAuthFlows;
    private buildServers;
    private buildSchema;
    private buildPaths;
    private buildMethod;
    protected buildOperation(controllerName: string, method: Tsoa.Method, defaultProduces?: string[]): Swagger.Operation3;
    private buildRequestBodyWithFormData;
    private buildRequestBody;
    private buildMediaType;
    private buildQueriesParameter;
    private buildParameter;
    protected buildProperties(source: Tsoa.Property[]): {
        [propertyName: string]: Swagger.Schema3;
    };
    protected getSwaggerTypeForReferenceType(referenceType: Tsoa.ReferenceType): Swagger.BaseSchema;
    protected getSwaggerTypeForPrimitiveType(dataType: Tsoa.PrimitiveTypeLiteral): Swagger.Schema;
    private isNull;
    private groupEnums;
    protected removeDuplicateSwaggerTypes(types: Array<Swagger.Schema | Swagger.BaseSchema>): Array<Swagger.Schema | Swagger.BaseSchema>;
    protected getSwaggerTypeForUnionType(type: Tsoa.UnionType, title?: string): {
        allOf: (Swagger.Schema | Swagger.BaseSchema)[];
        nullable: true;
    } | {
        nullable: true;
        type?: Swagger.DataType | undefined;
        format?: Swagger.DataFormat | undefined;
        additionalProperties?: boolean | Swagger.BaseSchema | undefined;
        properties?: {
            [propertyName: string]: Swagger.Schema3;
        } | undefined;
        discriminator?: string | undefined;
        readOnly?: boolean | undefined;
        xml?: Swagger.XML | undefined;
        externalDocs?: Swagger.ExternalDocs | undefined;
        example?: unknown;
        required?: string[] | undefined;
        $ref?: string | undefined;
        title?: string | undefined;
        description?: string | undefined;
        default?: unknown;
        multipleOf?: number | undefined;
        maximum?: number | undefined;
        exclusiveMaximum?: number | undefined;
        minimum?: number | undefined;
        exclusiveMinimum?: number | undefined;
        maxLength?: number | undefined;
        minLength?: number | undefined;
        pattern?: string | undefined;
        maxItems?: number | undefined;
        minItems?: number | undefined;
        uniqueItems?: boolean | undefined;
        maxProperties?: number | undefined;
        minProperties?: number | undefined;
        enum?: (string | number | boolean | null)[] | undefined;
        'x-enum-varnames'?: string[] | undefined;
        items?: Swagger.BaseSchema | undefined;
        allOf?: undefined;
    } | {
        nullable: true;
        type?: string | undefined;
        format?: Swagger.DataFormat | undefined;
        $ref?: string | undefined;
        title?: string | undefined;
        description?: string | undefined;
        default?: unknown;
        multipleOf?: number | undefined;
        maximum?: number | undefined;
        exclusiveMaximum?: number | undefined;
        minimum?: number | undefined;
        exclusiveMinimum?: number | undefined;
        maxLength?: number | undefined;
        minLength?: number | undefined;
        pattern?: string | undefined;
        maxItems?: number | undefined;
        minItems?: number | undefined;
        uniqueItems?: boolean | undefined;
        maxProperties?: number | undefined;
        minProperties?: number | undefined;
        enum?: (string | number | boolean | null)[] | undefined;
        'x-enum-varnames'?: string[] | undefined;
        items?: Swagger.BaseSchema | undefined;
        allOf?: undefined;
    } | {
        anyOf: (Swagger.Schema | Swagger.BaseSchema)[];
        nullable: true;
        title?: string | undefined;
        allOf?: undefined;
    } | {
        type?: Swagger.DataType | undefined;
        format?: Swagger.DataFormat | undefined;
        additionalProperties?: boolean | Swagger.BaseSchema | undefined;
        properties?: {
            [propertyName: string]: Swagger.Schema3;
        } | undefined;
        discriminator?: string | undefined;
        readOnly?: boolean | undefined;
        xml?: Swagger.XML | undefined;
        externalDocs?: Swagger.ExternalDocs | undefined;
        example?: unknown;
        required?: string[] | undefined;
        $ref?: string | undefined;
        title?: string | undefined;
        description?: string | undefined;
        default?: unknown;
        multipleOf?: number | undefined;
        maximum?: number | undefined;
        exclusiveMaximum?: number | undefined;
        minimum?: number | undefined;
        exclusiveMinimum?: number | undefined;
        maxLength?: number | undefined;
        minLength?: number | undefined;
        pattern?: string | undefined;
        maxItems?: number | undefined;
        minItems?: number | undefined;
        uniqueItems?: boolean | undefined;
        maxProperties?: number | undefined;
        minProperties?: number | undefined;
        enum?: (string | number | boolean | null)[] | undefined;
        'x-enum-varnames'?: string[] | undefined;
        items?: Swagger.BaseSchema | undefined;
        allOf?: undefined;
        nullable?: undefined;
    } | {
        type?: string | undefined;
        format?: Swagger.DataFormat | undefined;
        $ref?: string | undefined;
        title?: string | undefined;
        description?: string | undefined;
        default?: unknown;
        multipleOf?: number | undefined;
        maximum?: number | undefined;
        exclusiveMaximum?: number | undefined;
        minimum?: number | undefined;
        exclusiveMinimum?: number | undefined;
        maxLength?: number | undefined;
        minLength?: number | undefined;
        pattern?: string | undefined;
        maxItems?: number | undefined;
        minItems?: number | undefined;
        uniqueItems?: boolean | undefined;
        maxProperties?: number | undefined;
        minProperties?: number | undefined;
        enum?: (string | number | boolean | null)[] | undefined;
        'x-enum-varnames'?: string[] | undefined;
        items?: Swagger.BaseSchema | undefined;
        allOf?: undefined;
        nullable?: undefined;
    } | {
        anyOf: (Swagger.Schema | Swagger.BaseSchema)[];
        title?: string | undefined;
        allOf?: undefined;
        nullable?: undefined;
    };
    protected getSwaggerTypeForIntersectionType(type: Tsoa.IntersectionType, title?: string): {
        title?: string | undefined;
        allOf: (Swagger.Schema | Swagger.BaseSchema)[];
    };
    protected getSwaggerTypeForEnumType(enumType: Tsoa.EnumType, title?: string): Swagger.Schema3;
}
