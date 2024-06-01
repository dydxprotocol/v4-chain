import { ContractInfo } from 'wasm-ast-types';
interface ReadSchemaOpts {
    schemaDir: string;
    clean?: boolean;
}
export declare const readSchemas: ({ schemaDir, clean }: ReadSchemaOpts) => Promise<ContractInfo>;
export declare const findQueryMsg: (schemas: any) => any;
export declare const findExecuteMsg: (schemas: any) => any;
export declare const findAndParseTypes: (schemas: any) => Promise<{}>;
export {};
