import { ExtendedSpecConfig } from '../cli';
import { Tsoa, Swagger } from '@tsoa/runtime';
export declare abstract class SpecGenerator {
    protected readonly metadata: Tsoa.Metadata;
    protected readonly config: ExtendedSpecConfig;
    constructor(metadata: Tsoa.Metadata, config: ExtendedSpecConfig);
    protected buildAdditionalProperties(type: Tsoa.Type): Swagger.Schema | Swagger.BaseSchema;
    protected buildOperationIdTemplate(inlineTemplate: string): HandlebarsTemplateDelegate<any>;
    protected getOperationId(controllerName: string, method: Tsoa.Method): string;
    throwIfNotDataFormat(strToTest: string): Swagger.DataFormat;
    throwIfNotDataType(strToTest: string): Swagger.DataType;
    protected getSwaggerType(type: Tsoa.Type, title?: string): Swagger.Schema | Swagger.BaseSchema;
    protected abstract getSwaggerTypeForUnionType(type: Tsoa.UnionType, title?: string): any;
    protected abstract getSwaggerTypeForIntersectionType(type: Tsoa.IntersectionType, title?: string): any;
    protected abstract buildProperties(properties: Tsoa.Property[]): {
        [propertyName: string]: Swagger.Schema | Swagger.Schema3;
    };
    getSwaggerTypeForObjectLiteral(objectLiteral: Tsoa.NestedObjectLiteralType, title?: string): Swagger.Schema;
    protected getSwaggerTypeForReferenceType(_referenceType: Tsoa.ReferenceType): Swagger.BaseSchema;
    protected getSwaggerTypeForVoid(_dataType: 'void' | 'undefined'): Swagger.BaseSchema;
    protected determineImplicitAdditionalPropertiesValue: () => boolean;
    protected getSwaggerTypeForPrimitiveType(dataType: Tsoa.PrimitiveTypeLiteral): Swagger.Schema;
    protected getSwaggerTypeForArrayType(arrayType: Tsoa.ArrayType, title?: string): Swagger.Schema;
    protected determineTypesUsedInEnum(anEnum: Array<string | number | boolean | null>): Set<"string" | "number" | "boolean">;
    protected abstract getSwaggerTypeForEnumType(enumType: Tsoa.EnumType, title?: string): Swagger.Schema2 | Swagger.Schema3;
    protected hasUndefined(property: Tsoa.Property): boolean;
    protected queriesPropertyToQueryParameter(property: Tsoa.Property): Tsoa.Parameter;
}
