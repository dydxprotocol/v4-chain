import * as t from '@babel/types';
import { ExecuteMsg, QueryMsg } from '../types';
import { RenderContext } from '../context';
import { JSONSchema } from '../types';
interface ReactQueryHookQuery {
    context: RenderContext;
    hookName: string;
    hookParamsTypeName: string;
    hookKeyName: string;
    queryKeysName: string;
    responseType: string;
    methodName: string;
    jsonschema: any;
}
interface ReactQueryHooks {
    context: RenderContext;
    queryMsg: QueryMsg;
    contractName: string;
    QueryClient: string;
}
export declare const createReactQueryHooks: ({ context, queryMsg, contractName, QueryClient }: ReactQueryHooks) => any[];
export declare const createReactQueryHook: ({ context, hookName, hookParamsTypeName, responseType, hookKeyName, queryKeysName, methodName, jsonschema }: ReactQueryHookQuery) => t.ExportNamedDeclaration;
interface ReactQueryMutationHookInterface {
    context: RenderContext;
    ExecuteClient: string;
    mutationHookParamsTypeName: string;
    jsonschema: JSONSchema;
    useMutationTypeParameter: t.TSTypeParameterInstantiation;
}
/**
 * Example:
```
export interface Cw4UpdateMembersMutation {
  client: Cw4GroupClient
  args: {
    tokenId: string
    remove: string[]
  }
  options?: Omit<
    UseMutationOptions<ExecuteResult, Error, Pick<Cw4UpdateMembersMutation, 'args'>>,
    'mutationFn'
  >
}
```
 */
export declare const createReactQueryMutationArgsInterface: ({ context, ExecuteClient, mutationHookParamsTypeName, useMutationTypeParameter, jsonschema }: ReactQueryMutationHookInterface) => t.ExportNamedDeclaration;
interface ReactQueryMutationHooks {
    context: RenderContext;
    execMsg: ExecuteMsg;
    contractName: string;
    ExecuteClient: string;
}
export declare const createReactQueryMutationHooks: ({ context, execMsg, contractName, ExecuteClient }: ReactQueryMutationHooks) => any;
interface ReactQueryMutationHook {
    context: RenderContext;
    mutationHookName: string;
    mutationHookParamsTypeName: string;
    execMethodName: string;
    useMutationTypeParameter: t.TSTypeParameterInstantiation;
    hasMsg: boolean;
}
/**
 *
 * Example:
```
export const useCw4UpdateMembersMutation = ({ client, options }: Omit<Cw4UpdateMembersMutation, 'args'>) =>
  useMutation<ExecuteResult, Error, Pick<Cw4UpdateMembersMutation, 'args'>>(
    ({ args }) => client.updateMembers(args),
    options
  )
```
 */
export declare const createReactQueryMutationHook: ({ context, mutationHookName, mutationHookParamsTypeName, execMethodName, useMutationTypeParameter, hasMsg }: ReactQueryMutationHook) => t.ExportNamedDeclaration;
interface ReactQueryHookQueryInterface {
    context: RenderContext;
    QueryClient: string;
    hookParamsTypeName: string;
    queryInterfaceName: string;
    responseType: string;
    jsonschema: any;
}
export declare const createReactQueryHookInterface: ({ context, QueryClient, hookParamsTypeName, queryInterfaceName, responseType, jsonschema }: ReactQueryHookQueryInterface) => t.ExportNamedDeclaration;
export {};
