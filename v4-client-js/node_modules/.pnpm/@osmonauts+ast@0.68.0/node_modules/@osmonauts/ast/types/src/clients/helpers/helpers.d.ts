import * as t from '@babel/types';
import { Mutation } from '../../types';
import { GenericParseContext } from '../../encoding';
export declare const addMsgMethod: ({ methodName, typeUrl, TypeName, methodToCall }: {
    methodName: any;
    typeUrl: any;
    TypeName: any;
    methodToCall: any;
}) => t.ObjectMethod;
export declare const addFromJSONMethod: ({ methodName, typeUrl, TypeName }: {
    methodName: any;
    typeUrl: any;
    TypeName: any;
}) => t.ObjectMethod;
export declare const addFromPartialMethod: ({ methodName, typeUrl, TypeName }: {
    methodName: any;
    typeUrl: any;
    TypeName: any;
}) => t.ObjectMethod;
export declare const addToJSONMethod: ({ methodName, typeUrl, TypeName }: {
    methodName: any;
    typeUrl: any;
    TypeName: any;
}) => t.ObjectMethod;
export declare const addJsonMethod: ({ methodName, typeUrl, TypeName }: {
    methodName: any;
    typeUrl: any;
    TypeName: any;
}) => t.ObjectMethod;
export declare const addEncodedMethod: ({ methodName, typeUrl, TypeName }: {
    methodName: any;
    typeUrl: any;
    TypeName: any;
}) => t.ObjectMethod;
interface HelperObject {
    context: GenericParseContext;
    name: string;
    mutations: Mutation[];
}
export declare const createHelperObject: ({ context, name, mutations }: HelperObject) => t.ExportNamedDeclaration;
export {};
