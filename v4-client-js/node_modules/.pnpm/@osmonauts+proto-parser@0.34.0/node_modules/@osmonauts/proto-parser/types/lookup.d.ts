import { ProtoStore } from './store';
import { ProtoRef } from '@osmonauts/types';
export interface Lookup {
    obj: any;
    name: string;
    import: string;
    importType: string;
    importedName: string;
    package: string;
}
export declare const recursiveLookup: (proto: any, name: string, scope?: string[], allowNested?: boolean) => any;
export declare const importLookup: (store: ProtoStore, ref: ProtoRef, name: string) => {
    name: string;
    importType: string;
    import: string;
    importedName: string;
    obj: any;
};
export declare const protoImportLookup: (store: ProtoStore, ref: ProtoRef, name: string) => Lookup;
export declare const protoScopeImportLookup: (store: ProtoStore, ref: ProtoRef, name: string) => Lookup;
export declare const lookup: (store: ProtoStore, ref: ProtoRef, name: string, allowNested?: boolean) => any;
export declare const lookupNested: (ref: ProtoRef, traversal: string[], name: string, allowNested?: boolean) => any;
export declare const lookupAny: (store: ProtoStore, ref: ProtoRef, name: string) => Lookup;
export declare const lookupAnyFromImports: (store: ProtoStore, ref: ProtoRef, name: string) => Lookup;
export declare const lookupLocal: (store: ProtoStore, ref: ProtoRef, name: string) => Lookup;
export declare const externalLookup: (store: ProtoStore, ref: ProtoRef, name: string) => Lookup;
