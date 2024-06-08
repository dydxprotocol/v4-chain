import { FileDescriptorProto } from "ts-proto-descriptors";
/** This type is expecting a value from the Fields constant. */
export type FieldID = number;
/**
 * The field values here represent the proto field IDs associated with the types
 * (file,message,enum,service).
 *
 * For more information read the comments for SourceCodeInfo declared in
 * google's 'descriptor.proto' file, see:
 * https://github.com/protocolbuffers/protobuf/blob/master/src/google/protobuf/descriptor.proto#L730
 */
export declare const Fields: {
    file: {
        syntax: number;
        message_type: number;
        enum_type: number;
        service: number;
        extension: number;
    };
    message: {
        field: number;
        nested_type: number;
        enum_type: number;
        oneof_decl: number;
    };
    enum: {
        value: number;
    };
    service: {
        method: number;
    };
};
/**
 * This type is simply an interface on the SourceCodeInfo.Location message.
 */
export interface SourceDescription {
    readonly span: number[];
    readonly leadingComments: string;
    readonly trailingComments: string;
    readonly leadingDetachedComments: string[];
}
/**
 * Mapping from a string of dotted notation `path` parts to efficiently
 * lookup the related source information.
 */
export type SourceInfoMap = {
    [key: string]: SourceDescription;
};
/**
 * This class provides direct lookup and navigation through the type
 * system by the use of lookup/open to access the source info for types
 * defined in a protocol buffer.
 */
export default class SourceInfo implements SourceDescription {
    private readonly sourceCode;
    private readonly selfDescription;
    /** Returns an empty SourceInfo */
    static empty(): SourceInfo;
    /**
     * Creates the SourceInfo from the FileDescriptorProto given to you
     * by the protoc compiler. It indexes file.sourceCodeInfo by dotted
     * path notation and returns the root SourceInfo.
     */
    static fromDescriptor(file: FileDescriptorProto): SourceInfo;
    private constructor();
    /** Returns the code span [start line, start column, end line] */
    get span(): number[];
    /** Leading consecutive comment lines prior to the current element */
    get leadingComments(): string;
    /** Documentation is unclear about what exactly this is */
    get trailingComments(): string;
    /** Detached comments are those preceeding but separated by a blank non-comment line */
    get leadingDetachedComments(): string[];
    /** Return the source info for the field id and index specficied */
    lookup(type: FieldID, index?: number): SourceDescription;
    /** Returns a new SourceInfo class representing the field id and index specficied */
    open(type: FieldID, index: number): SourceInfo;
}
