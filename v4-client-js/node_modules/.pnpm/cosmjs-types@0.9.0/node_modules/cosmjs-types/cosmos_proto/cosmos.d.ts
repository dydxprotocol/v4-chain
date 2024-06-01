import { BinaryReader, BinaryWriter } from "../binary";
export declare const protobufPackage = "cosmos_proto";
export declare enum ScalarType {
    SCALAR_TYPE_UNSPECIFIED = 0,
    SCALAR_TYPE_STRING = 1,
    SCALAR_TYPE_BYTES = 2,
    UNRECOGNIZED = -1
}
export declare function scalarTypeFromJSON(object: any): ScalarType;
export declare function scalarTypeToJSON(object: ScalarType): string;
/**
 * InterfaceDescriptor describes an interface type to be used with
 * accepts_interface and implements_interface and declared by declare_interface.
 */
export interface InterfaceDescriptor {
    /**
     * name is the name of the interface. It should be a short-name (without
     * a period) such that the fully qualified name of the interface will be
     * package.name, ex. for the package a.b and interface named C, the
     * fully-qualified name will be a.b.C.
     */
    name: string;
    /**
     * description is a human-readable description of the interface and its
     * purpose.
     */
    description: string;
}
/**
 * ScalarDescriptor describes an scalar type to be used with
 * the scalar field option and declared by declare_scalar.
 * Scalars extend simple protobuf built-in types with additional
 * syntax and semantics, for instance to represent big integers.
 * Scalars should ideally define an encoding such that there is only one
 * valid syntactical representation for a given semantic meaning,
 * i.e. the encoding should be deterministic.
 */
export interface ScalarDescriptor {
    /**
     * name is the name of the scalar. It should be a short-name (without
     * a period) such that the fully qualified name of the scalar will be
     * package.name, ex. for the package a.b and scalar named C, the
     * fully-qualified name will be a.b.C.
     */
    name: string;
    /**
     * description is a human-readable description of the scalar and its
     * encoding format. For instance a big integer or decimal scalar should
     * specify precisely the expected encoding format.
     */
    description: string;
    /**
     * field_type is the type of field with which this scalar can be used.
     * Scalars can be used with one and only one type of field so that
     * encoding standards and simple and clear. Currently only string and
     * bytes fields are supported for scalars.
     */
    fieldType: ScalarType[];
}
export declare const InterfaceDescriptor: {
    typeUrl: string;
    encode(message: InterfaceDescriptor, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): InterfaceDescriptor;
    fromJSON(object: any): InterfaceDescriptor;
    toJSON(message: InterfaceDescriptor): unknown;
    fromPartial<I extends {
        name?: string | undefined;
        description?: string | undefined;
    } & {
        name?: string | undefined;
        description?: string | undefined;
    } & Record<Exclude<keyof I, keyof InterfaceDescriptor>, never>>(object: I): InterfaceDescriptor;
};
export declare const ScalarDescriptor: {
    typeUrl: string;
    encode(message: ScalarDescriptor, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ScalarDescriptor;
    fromJSON(object: any): ScalarDescriptor;
    toJSON(message: ScalarDescriptor): unknown;
    fromPartial<I extends {
        name?: string | undefined;
        description?: string | undefined;
        fieldType?: ScalarType[] | undefined;
    } & {
        name?: string | undefined;
        description?: string | undefined;
        fieldType?: (ScalarType[] & ScalarType[] & Record<Exclude<keyof I["fieldType"], keyof ScalarType[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof ScalarDescriptor>, never>>(object: I): ScalarDescriptor;
};
