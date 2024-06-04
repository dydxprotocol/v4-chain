"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Fields = void 0;
/**
 * The field values here represent the proto field IDs associated with the types
 * (file,message,enum,service).
 *
 * For more information read the comments for SourceCodeInfo declared in
 * google's 'descriptor.proto' file, see:
 * https://github.com/protocolbuffers/protobuf/blob/master/src/google/protobuf/descriptor.proto#L730
 */
exports.Fields = {
    file: {
        syntax: 12,
        message_type: 4,
        enum_type: 5,
        service: 6,
        extension: 7,
    },
    message: {
        field: 2,
        nested_type: 3,
        enum_type: 4,
        oneof_decl: 8,
    },
    enum: {
        value: 2,
    },
    service: {
        method: 2,
    },
};
/** An empty SourceDescription for when one is not available. */
class EmptyDescription {
    constructor() {
        this.span = [];
        this.leadingComments = "";
        this.trailingComments = "";
        this.leadingDetachedComments = [];
    }
}
/**
 * This class provides direct lookup and navigation through the type
 * system by the use of lookup/open to access the source info for types
 * defined in a protocol buffer.
 */
class SourceInfo {
    /** Returns an empty SourceInfo */
    static empty() {
        return new SourceInfo({}, new EmptyDescription());
    }
    /**
     * Creates the SourceInfo from the FileDescriptorProto given to you
     * by the protoc compiler. It indexes file.sourceCodeInfo by dotted
     * path notation and returns the root SourceInfo.
     */
    static fromDescriptor(file) {
        let map = {};
        if (file.sourceCodeInfo && file.sourceCodeInfo.location) {
            file.sourceCodeInfo.location.forEach((loc) => {
                map[loc.path.join(".")] = loc;
            });
        }
        return new SourceInfo(map, new EmptyDescription());
    }
    // Private
    constructor(sourceCode, selfDescription) {
        this.sourceCode = sourceCode;
        this.selfDescription = selfDescription;
    }
    /** Returns the code span [start line, start column, end line] */
    get span() {
        return this.selfDescription.span;
    }
    /** Leading consecutive comment lines prior to the current element */
    get leadingComments() {
        return this.selfDescription.leadingComments;
    }
    /** Documentation is unclear about what exactly this is */
    get trailingComments() {
        return this.selfDescription.trailingComments;
    }
    /** Detached comments are those preceeding but separated by a blank non-comment line */
    get leadingDetachedComments() {
        return this.selfDescription.leadingDetachedComments;
    }
    /** Return the source info for the field id and index specficied */
    lookup(type, index) {
        if (index === undefined) {
            return this.sourceCode[`${type}`] || new EmptyDescription();
        }
        return this.sourceCode[`${type}.${index}`] || new EmptyDescription();
    }
    /** Returns a new SourceInfo class representing the field id and index specficied */
    open(type, index) {
        const prefix = `${type}.${index}.`;
        const map = {};
        Object.keys(this.sourceCode)
            .filter((key) => key.startsWith(prefix))
            .forEach((key) => {
            map[key.substr(prefix.length)] = this.sourceCode[key];
        });
        return new SourceInfo(map, this.lookup(type, index));
    }
}
exports.default = SourceInfo;
