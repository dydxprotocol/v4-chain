import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../../helpers";
/** StorageType */
export declare enum StorageType {
    /**
     * STORAGE_TYPE_DEFAULT_UNSPECIFIED - STORAGE_TYPE_DEFAULT_UNSPECIFIED indicates the persistent storage where all
     * data is stored in the regular Merkle-tree backed KV-store.
     */
    STORAGE_TYPE_DEFAULT_UNSPECIFIED = 0,
    /**
     * STORAGE_TYPE_MEMORY - STORAGE_TYPE_MEMORY indicates in-memory storage that will be
     * reloaded every time an app restarts. Tables with this type of storage
     * will by default be ignored when importing and exporting a module's
     * state from JSON.
     */
    STORAGE_TYPE_MEMORY = 1,
    /**
     * STORAGE_TYPE_TRANSIENT - STORAGE_TYPE_TRANSIENT indicates transient storage that is reset
     * at the end of every block. Tables with this type of storage
     * will by default be ignored when importing and exporting a module's
     * state from JSON.
     */
    STORAGE_TYPE_TRANSIENT = 2,
    UNRECOGNIZED = -1
}
export declare const StorageTypeSDKType: typeof StorageType;
export declare function storageTypeFromJSON(object: any): StorageType;
export declare function storageTypeToJSON(object: StorageType): string;
/** ModuleSchemaDescriptor describe's a module's ORM schema. */
export interface ModuleSchemaDescriptor {
    schemaFile: ModuleSchemaDescriptor_FileEntry[];
    /**
     * prefix is an optional prefix that precedes all keys in this module's
     * store.
     */
    prefix: Uint8Array;
}
/** ModuleSchemaDescriptor describe's a module's ORM schema. */
export interface ModuleSchemaDescriptorSDKType {
    schema_file: ModuleSchemaDescriptor_FileEntrySDKType[];
    prefix: Uint8Array;
}
/** FileEntry describes an ORM file used in a module. */
export interface ModuleSchemaDescriptor_FileEntry {
    /**
     * id is a prefix that will be varint encoded and prepended to all the
     * table keys specified in the file's tables.
     */
    id: number;
    /**
     * proto_file_name is the name of a file .proto in that contains
     * table definitions. The .proto file must be in a package that the
     * module has referenced using cosmos.app.v1.ModuleDescriptor.use_package.
     */
    protoFileName: string;
    /**
     * storage_type optionally indicates the type of storage this file's
     * tables should used. If it is left unspecified, the default KV-storage
     * of the app will be used.
     */
    storageType: StorageType;
}
/** FileEntry describes an ORM file used in a module. */
export interface ModuleSchemaDescriptor_FileEntrySDKType {
    id: number;
    proto_file_name: string;
    storage_type: StorageType;
}
export declare const ModuleSchemaDescriptor: {
    encode(message: ModuleSchemaDescriptor, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ModuleSchemaDescriptor;
    fromPartial(object: DeepPartial<ModuleSchemaDescriptor>): ModuleSchemaDescriptor;
};
export declare const ModuleSchemaDescriptor_FileEntry: {
    encode(message: ModuleSchemaDescriptor_FileEntry, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ModuleSchemaDescriptor_FileEntry;
    fromPartial(object: DeepPartial<ModuleSchemaDescriptor_FileEntry>): ModuleSchemaDescriptor_FileEntry;
};
