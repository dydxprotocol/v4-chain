import { ExtendedSpecConfig } from '../cli';
import { Tsoa, Swagger } from '@tsoa/runtime';
import { SpecGenerator } from './specGenerator';
export declare class SpecGenerator2 extends SpecGenerator {
    protected readonly metadata: Tsoa.Metadata;
    protected readonly config: ExtendedSpecConfig;
    constructor(metadata: Tsoa.Metadata, config: ExtendedSpecConfig);
    GetSpec(): Swagger.Spec2;
    private buildDefinitions;
    private buildPaths;
    private buildMethod;
    protected buildOperation(controllerName: string, method: Tsoa.Method, defaultProduces?: string[]): Swagger.Operation;
    private buildBodyPropParameter;
    private buildQueriesParameter;
    private buildParameter;
    protected buildProperties(source: Tsoa.Property[]): {
        [propertyName: string]: Swagger.Schema2;
    };
    protected getSwaggerTypeForUnionType(type: Tsoa.UnionType): Swagger.BaseSchema;
    protected getSwaggerTypeForIntersectionType(type: Tsoa.IntersectionType): {
        type: string;
        properties: {};
    };
    protected getSwaggerTypeForReferenceType(referenceType: Tsoa.ReferenceType): Swagger.BaseSchema;
    private decideEnumType;
    protected getSwaggerTypeForEnumType(enumType: Tsoa.EnumType): Swagger.Schema2;
}
