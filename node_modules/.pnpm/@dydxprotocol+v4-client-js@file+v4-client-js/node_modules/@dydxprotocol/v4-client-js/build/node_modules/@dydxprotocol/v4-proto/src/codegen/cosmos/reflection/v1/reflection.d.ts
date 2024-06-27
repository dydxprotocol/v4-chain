import { FileDescriptorProto, FileDescriptorProtoSDKType } from "../../../google/protobuf/descriptor";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../../helpers";
/** FileDescriptorsRequest is the Query/FileDescriptors request type. */
export interface FileDescriptorsRequest {
}
/** FileDescriptorsRequest is the Query/FileDescriptors request type. */
export interface FileDescriptorsRequestSDKType {
}
/** FileDescriptorsResponse is the Query/FileDescriptors response type. */
export interface FileDescriptorsResponse {
    /** files is the file descriptors. */
    files: FileDescriptorProto[];
}
/** FileDescriptorsResponse is the Query/FileDescriptors response type. */
export interface FileDescriptorsResponseSDKType {
    files: FileDescriptorProtoSDKType[];
}
export declare const FileDescriptorsRequest: {
    encode(_: FileDescriptorsRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): FileDescriptorsRequest;
    fromPartial(_: DeepPartial<FileDescriptorsRequest>): FileDescriptorsRequest;
};
export declare const FileDescriptorsResponse: {
    encode(message: FileDescriptorsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): FileDescriptorsResponse;
    fromPartial(object: DeepPartial<FileDescriptorsResponse>): FileDescriptorsResponse;
};
