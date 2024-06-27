import * as t from '@babel/types';
import { AminoParseContext } from '../encoding';
export interface ServiceMethod {
    methodName: string;
    typeUrl: string;
    TypeName: string;
}
export declare const createTypeRegistryObject: (mutation: ServiceMethod) => t.ObjectProperty;
export declare const createTypeRegistry: (context: AminoParseContext, mutations: ServiceMethod[]) => t.ExportNamedDeclaration;
export declare const createRegistryLoader: (context: AminoParseContext) => t.ExportNamedDeclaration;
