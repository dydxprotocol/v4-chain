import { TelescopeOptions, ProtoField, ProtoRef } from '@osmonauts/types';
import { ProtoStore } from '@osmonauts/proto-parser';
export interface ParseContext {
    options: TelescopeOptions;
    imports: ImportUsage[];
    utils: Record<string, boolean>;
    addUtil: Function;
}
export interface ImportUsage {
    type: 'typeImport' | 'toJSONEnum' | 'fromJSONEnum';
    name: string;
    import: string;
    importedAs?: string;
}
export declare class GenericParseContext implements ParseContext {
    options: TelescopeOptions;
    imports: ImportUsage[];
    utils: Record<string, boolean>;
    store: ProtoStore;
    ref: ProtoRef;
    constructor(ref: ProtoRef, store: ProtoStore, options: TelescopeOptions);
    pluginValue(name: any): any;
    isExcluded(): any;
    addUtil(util: any): void;
    addImport(imp: ImportUsage): void;
    getTypeName(field: ProtoField): string;
}
export declare class AminoParseContext extends GenericParseContext implements ParseContext {
    aminoCasingFn: Function;
    constructor(ref: ProtoRef, store: ProtoStore, options: TelescopeOptions);
    private setAminoCasingFn;
    aminoCaseField(field: ProtoField): string;
    private lookupTypeFromCurrentPath;
    getTypeFromCurrentPath(field: ProtoField, currentProtoPath: string): any;
    lookupEnumFromJson(field: ProtoField, currentProtoPath: string): string;
    lookupEnumToJson(field: ProtoField, currentProtoPath: string): string;
}
export declare class ProtoParseContext extends GenericParseContext implements ParseContext {
    constructor(ref: ProtoRef, store: ProtoStore, options: TelescopeOptions);
    getToEnum(field: ProtoField): string;
    getFromEnum(field: ProtoField): string;
}
