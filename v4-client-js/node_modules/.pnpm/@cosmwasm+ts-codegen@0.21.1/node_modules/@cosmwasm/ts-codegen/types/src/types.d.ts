import { JSONSchema } from "wasm-ast-types";
interface KeyedSchema {
    [key: string]: JSONSchema;
}
export interface IDLObject {
    contract_name: string;
    contract_version: string;
    idl_version: string;
    instantiate: JSONSchema;
    execute: JSONSchema;
    query: JSONSchema;
    migrate: JSONSchema;
    sudo: JSONSchema;
    responses: KeyedSchema;
}
export {};
