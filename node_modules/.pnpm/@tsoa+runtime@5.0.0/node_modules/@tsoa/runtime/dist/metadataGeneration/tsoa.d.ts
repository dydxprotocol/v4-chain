import { ExtensionType } from '../decorators/extension';
import type { Swagger } from '../swagger/swagger';
export declare namespace Tsoa {
    interface Metadata {
        controllers: Controller[];
        referenceTypeMap: ReferenceTypeMap;
    }
    interface Controller {
        location: string;
        methods: Method[];
        name: string;
        path: string;
        produces?: string[];
    }
    interface Method {
        extensions: Extension[];
        deprecated?: boolean;
        description?: string;
        method: 'get' | 'post' | 'put' | 'delete' | 'options' | 'head' | 'patch';
        name: string;
        parameters: Parameter[];
        path: string;
        produces?: string[];
        consumes?: string;
        type: Type;
        tags?: string[];
        responses: Response[];
        successStatus?: number;
        security: Security[];
        summary?: string;
        isHidden: boolean;
        operationId?: string;
    }
    interface Parameter {
        parameterName: string;
        example?: Array<{
            [exampleName: string]: Swagger.Example3;
        }>;
        description?: string;
        in: 'query' | 'queries' | 'header' | 'path' | 'formData' | 'body' | 'body-prop' | 'request' | 'res';
        name: string;
        required?: boolean;
        type: Type;
        default?: unknown;
        validators: Validators;
        deprecated: boolean;
        exampleLabels?: Array<string | undefined>;
    }
    interface ResParameter extends Response, Parameter {
        in: 'res';
        description: string;
    }
    interface ArrayParameter extends Parameter {
        type: ArrayType;
        collectionFormat?: 'csv' | 'multi' | 'pipes' | 'ssv' | 'tsv';
    }
    interface Validators {
        [key: string]: {
            value?: unknown;
            errorMsg?: string;
        };
    }
    interface Security {
        [key: string]: string[];
    }
    interface Extension {
        key: string;
        value: ExtensionType | ExtensionType[];
    }
    interface Response {
        description: string;
        name: string;
        produces?: string[];
        schema?: Type;
        examples?: Array<{
            [exampleName: string]: Swagger.Example3;
        }>;
        exampleLabels?: Array<string | undefined>;
        headers?: HeaderType;
    }
    interface Property {
        default?: unknown;
        description?: string;
        format?: string;
        example?: unknown;
        name: string;
        type: Type;
        required: boolean;
        validators: Validators;
        deprecated: boolean;
        extensions?: Extension[];
    }
    type TypeStringLiteral = 'string' | 'boolean' | 'double' | 'float' | 'file' | 'integer' | 'long' | 'enum' | 'array' | 'datetime' | 'date' | 'binary' | 'buffer' | 'byte' | 'void' | 'object' | 'any' | 'refEnum' | 'refObject' | 'refAlias' | 'nestedObjectLiteral' | 'union' | 'intersection' | 'undefined';
    type RefTypeLiteral = 'refObject' | 'refEnum' | 'refAlias';
    type PrimitiveTypeLiteral = Exclude<TypeStringLiteral, RefTypeLiteral | 'enum' | 'array' | 'void' | 'undefined' | 'nestedObjectLiteral' | 'union' | 'intersection'>;
    interface TypeBase {
        dataType: TypeStringLiteral;
    }
    type PrimitiveType = StringType | BooleanType | DoubleType | FloatType | IntegerType | LongType | VoidType | UndefinedType;
    /**
     * This is one of the possible objects that tsoa creates that helps the code store information about the type it found in the code.
     */
    type Type = PrimitiveType | ObjectsNoPropsType | EnumType | ArrayType | FileType | DateTimeType | DateType | BinaryType | BufferType | ByteType | AnyType | RefEnumType | RefObjectType | RefAliasType | NestedObjectLiteralType | UnionType | IntersectionType;
    interface StringType extends TypeBase {
        dataType: 'string';
    }
    interface BooleanType extends TypeBase {
        dataType: 'boolean';
    }
    /**
     * This is the type that occurs when a developer writes `const foo: object = {}` since it can no longer have any properties added to it.
     */
    interface ObjectsNoPropsType extends TypeBase {
        dataType: 'object';
    }
    interface DoubleType extends TypeBase {
        dataType: 'double';
    }
    interface FloatType extends TypeBase {
        dataType: 'float';
    }
    interface IntegerType extends TypeBase {
        dataType: 'integer';
    }
    interface LongType extends TypeBase {
        dataType: 'long';
    }
    /**
     * Not to be confused with `RefEnumType` which is a reusable enum which has a $ref name generated for it. This however, is an inline enum.
     */
    interface EnumType extends TypeBase {
        dataType: 'enum';
        enums: Array<string | number | boolean | null>;
    }
    interface ArrayType extends TypeBase {
        dataType: 'array';
        elementType: Type;
    }
    interface DateType extends TypeBase {
        dataType: 'date';
    }
    interface FileType extends TypeBase {
        dataType: 'file';
    }
    interface DateTimeType extends TypeBase {
        dataType: 'datetime';
    }
    interface BinaryType extends TypeBase {
        dataType: 'binary';
    }
    interface BufferType extends TypeBase {
        dataType: 'buffer';
    }
    interface ByteType extends TypeBase {
        dataType: 'byte';
    }
    interface VoidType extends TypeBase {
        dataType: 'void';
    }
    interface UndefinedType extends TypeBase {
        dataType: 'undefined';
    }
    interface AnyType extends TypeBase {
        dataType: 'any';
    }
    interface NestedObjectLiteralType extends TypeBase {
        dataType: 'nestedObjectLiteral';
        properties: Property[];
        additionalProperties?: Type;
    }
    interface RefEnumType extends ReferenceTypeBase {
        dataType: 'refEnum';
        enums: Array<string | number>;
        enumVarnames?: string[];
    }
    interface RefObjectType extends ReferenceTypeBase {
        dataType: 'refObject';
        properties: Property[];
        additionalProperties?: Type;
    }
    interface RefAliasType extends Omit<Property, 'name' | 'required'>, ReferenceTypeBase {
        dataType: 'refAlias';
    }
    type ReferenceType = RefEnumType | RefObjectType | RefAliasType;
    interface ReferenceTypeBase extends TypeBase {
        description?: string;
        dataType: RefTypeLiteral;
        refName: string;
        example?: unknown;
        deprecated: boolean;
    }
    interface UnionType extends TypeBase {
        dataType: 'union';
        types: Type[];
    }
    interface IntersectionType extends TypeBase {
        dataType: 'intersection';
        types: Type[];
    }
    interface ReferenceTypeMap {
        [refName: string]: Tsoa.ReferenceType;
    }
    interface MethodsSignatureMap {
        [signature: string]: string[];
    }
    type HeaderType = Tsoa.NestedObjectLiteralType | Tsoa.RefObjectType;
}
