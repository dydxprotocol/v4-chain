import * as ts from 'typescript';
import { MetadataGenerator } from './metadataGenerator';
import { Tsoa } from '@tsoa/runtime';
interface Context {
    [name: string]: ts.TypeReferenceNode | ts.TypeNode;
}
export declare class TypeResolver {
    private readonly typeNode;
    private readonly current;
    private readonly parentNode?;
    private context;
    private readonly referencer?;
    private readonly addToRefTypeMap;
    constructor(typeNode: ts.TypeNode, current: MetadataGenerator, parentNode?: ts.Node | undefined, context?: Context, referencer?: ts.TypeNode | undefined, addToRefTypeMap?: boolean);
    static clearCache(): void;
    resolve(): Tsoa.Type;
    private getLiteralValue;
    private getPrimitiveType;
    private getDateType;
    private getDesignatedModels;
    private hasFlag;
    private getEnumerateType;
    private getReferenceType;
    private getTypeAliasReference;
    private getModelReference;
    private getRefTypeName;
    private attemptToResolveKindToPrimitive;
    private contextualizedName;
    private handleCachingAndCircularReferences;
    private createCircularDependencyResolver;
    private nodeIsUsable;
    private getModelTypeDeclaration;
    private getSymbolAtLocation;
    private getModelProperties;
    private propertyFromSignature;
    private propertyFromDeclaration;
    private getModelAdditionalProperties;
    private typeArgumentsToContext;
    private getModelInheritedProperties;
    private hasPublicModifier;
    private hasStaticModifier;
    private isAccessibleParameter;
    private getNodeDescription;
    private getNodeFormat;
    private getNodeExample;
    private getNodeExtension;
    private getDecoratorsByIdentifier;
}
export {};
