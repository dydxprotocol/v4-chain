import { PageRequest, PageResponse } from "../../../base/query/v1beta1/pagination";
import { Any } from "../../../../google/protobuf/any";
import { Timestamp } from "../../../../google/protobuf/timestamp";
import { Duration } from "../../../../google/protobuf/duration";
import { BinaryReader, BinaryWriter } from "../../../../binary";
import { Rpc } from "../../../../helpers";
export declare const protobufPackage = "cosmos.orm.query.v1alpha1";
/** GetRequest is the Query/Get request type. */
export interface GetRequest {
    /** message_name is the fully-qualified message name of the ORM table being queried. */
    messageName: string;
    /**
     * index is the index fields expression used in orm definitions. If it
     * is empty, the table's primary key is assumed. If it is non-empty, it must
     * refer to an unique index.
     */
    index: string;
    /**
     * values are the values of the fields corresponding to the requested index.
     * There must be as many values provided as there are fields in the index and
     * these values must correspond to the index field types.
     */
    values: IndexValue[];
}
/** GetResponse is the Query/Get response type. */
export interface GetResponse {
    /**
     * result is the result of the get query. If no value is found, the gRPC
     * status code NOT_FOUND will be returned.
     */
    result?: Any;
}
/** ListRequest is the Query/List request type. */
export interface ListRequest {
    /** message_name is the fully-qualified message name of the ORM table being queried. */
    messageName: string;
    /**
     * index is the index fields expression used in orm definitions. If it
     * is empty, the table's primary key is assumed.
     */
    index: string;
    /** prefix defines a prefix query. */
    prefix?: ListRequest_Prefix;
    /** range defines a range query. */
    range?: ListRequest_Range;
    /** pagination is the pagination request. */
    pagination?: PageRequest;
}
/** Prefix specifies the arguments to a prefix query. */
export interface ListRequest_Prefix {
    /**
     * values specifies the index values for the prefix query.
     * It is valid to special a partial prefix with fewer values than
     * the number of fields in the index.
     */
    values: IndexValue[];
}
/** Range specifies the arguments to a range query. */
export interface ListRequest_Range {
    /**
     * start specifies the starting index values for the range query.
     * It is valid to provide fewer values than the number of fields in the
     * index.
     */
    start: IndexValue[];
    /**
     * end specifies the inclusive ending index values for the range query.
     * It is valid to provide fewer values than the number of fields in the
     * index.
     */
    end: IndexValue[];
}
/** ListResponse is the Query/List response type. */
export interface ListResponse {
    /** results are the results of the query. */
    results: Any[];
    /** pagination is the pagination response. */
    pagination?: PageResponse;
}
/** IndexValue represents the value of a field in an ORM index expression. */
export interface IndexValue {
    /**
     * uint specifies a value for an uint32, fixed32, uint64, or fixed64
     * index field.
     */
    uint?: bigint;
    /**
     * int64 specifies a value for an int32, sfixed32, int64, or sfixed64
     * index field.
     */
    int?: bigint;
    /** str specifies a value for a string index field. */
    str?: string;
    /** bytes specifies a value for a bytes index field. */
    bytes?: Uint8Array;
    /** enum specifies a value for an enum index field. */
    enum?: string;
    /** bool specifies a value for a bool index field. */
    bool?: boolean;
    /** timestamp specifies a value for a timestamp index field. */
    timestamp?: Timestamp;
    /** duration specifies a value for a duration index field. */
    duration?: Duration;
}
export declare const GetRequest: {
    typeUrl: string;
    encode(message: GetRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GetRequest;
    fromJSON(object: any): GetRequest;
    toJSON(message: GetRequest): unknown;
    fromPartial<I extends {
        messageName?: string | undefined;
        index?: string | undefined;
        values?: {
            uint?: bigint | undefined;
            int?: bigint | undefined;
            str?: string | undefined;
            bytes?: Uint8Array | undefined;
            enum?: string | undefined;
            bool?: boolean | undefined;
            timestamp?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            duration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        }[] | undefined;
    } & {
        messageName?: string | undefined;
        index?: string | undefined;
        values?: ({
            uint?: bigint | undefined;
            int?: bigint | undefined;
            str?: string | undefined;
            bytes?: Uint8Array | undefined;
            enum?: string | undefined;
            bool?: boolean | undefined;
            timestamp?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            duration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        }[] & ({
            uint?: bigint | undefined;
            int?: bigint | undefined;
            str?: string | undefined;
            bytes?: Uint8Array | undefined;
            enum?: string | undefined;
            bool?: boolean | undefined;
            timestamp?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            duration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        } & {
            uint?: bigint | undefined;
            int?: bigint | undefined;
            str?: string | undefined;
            bytes?: Uint8Array | undefined;
            enum?: string | undefined;
            bool?: boolean | undefined;
            timestamp?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["values"][number]["timestamp"], keyof Timestamp>, never>) | undefined;
            duration?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["values"][number]["duration"], keyof Duration>, never>) | undefined;
        } & Record<Exclude<keyof I["values"][number], keyof IndexValue>, never>)[] & Record<Exclude<keyof I["values"], keyof {
            uint?: bigint | undefined;
            int?: bigint | undefined;
            str?: string | undefined;
            bytes?: Uint8Array | undefined;
            enum?: string | undefined;
            bool?: boolean | undefined;
            timestamp?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            duration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof GetRequest>, never>>(object: I): GetRequest;
};
export declare const GetResponse: {
    typeUrl: string;
    encode(message: GetResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GetResponse;
    fromJSON(object: any): GetResponse;
    toJSON(message: GetResponse): unknown;
    fromPartial<I extends {
        result?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
    } & {
        result?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["result"], keyof Any>, never>) | undefined;
    } & Record<Exclude<keyof I, "result">, never>>(object: I): GetResponse;
};
export declare const ListRequest: {
    typeUrl: string;
    encode(message: ListRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ListRequest;
    fromJSON(object: any): ListRequest;
    toJSON(message: ListRequest): unknown;
    fromPartial<I extends {
        messageName?: string | undefined;
        index?: string | undefined;
        prefix?: {
            values?: {
                uint?: bigint | undefined;
                int?: bigint | undefined;
                str?: string | undefined;
                bytes?: Uint8Array | undefined;
                enum?: string | undefined;
                bool?: boolean | undefined;
                timestamp?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                duration?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
            }[] | undefined;
        } | undefined;
        range?: {
            start?: {
                uint?: bigint | undefined;
                int?: bigint | undefined;
                str?: string | undefined;
                bytes?: Uint8Array | undefined;
                enum?: string | undefined;
                bool?: boolean | undefined;
                timestamp?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                duration?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
            }[] | undefined;
            end?: {
                uint?: bigint | undefined;
                int?: bigint | undefined;
                str?: string | undefined;
                bytes?: Uint8Array | undefined;
                enum?: string | undefined;
                bool?: boolean | undefined;
                timestamp?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                duration?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
            }[] | undefined;
        } | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        messageName?: string | undefined;
        index?: string | undefined;
        prefix?: ({
            values?: {
                uint?: bigint | undefined;
                int?: bigint | undefined;
                str?: string | undefined;
                bytes?: Uint8Array | undefined;
                enum?: string | undefined;
                bool?: boolean | undefined;
                timestamp?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                duration?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
            }[] | undefined;
        } & {
            values?: ({
                uint?: bigint | undefined;
                int?: bigint | undefined;
                str?: string | undefined;
                bytes?: Uint8Array | undefined;
                enum?: string | undefined;
                bool?: boolean | undefined;
                timestamp?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                duration?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
            }[] & ({
                uint?: bigint | undefined;
                int?: bigint | undefined;
                str?: string | undefined;
                bytes?: Uint8Array | undefined;
                enum?: string | undefined;
                bool?: boolean | undefined;
                timestamp?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                duration?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
            } & {
                uint?: bigint | undefined;
                int?: bigint | undefined;
                str?: string | undefined;
                bytes?: Uint8Array | undefined;
                enum?: string | undefined;
                bool?: boolean | undefined;
                timestamp?: ({
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & Record<Exclude<keyof I["prefix"]["values"][number]["timestamp"], keyof Timestamp>, never>) | undefined;
                duration?: ({
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & Record<Exclude<keyof I["prefix"]["values"][number]["duration"], keyof Duration>, never>) | undefined;
            } & Record<Exclude<keyof I["prefix"]["values"][number], keyof IndexValue>, never>)[] & Record<Exclude<keyof I["prefix"]["values"], keyof {
                uint?: bigint | undefined;
                int?: bigint | undefined;
                str?: string | undefined;
                bytes?: Uint8Array | undefined;
                enum?: string | undefined;
                bool?: boolean | undefined;
                timestamp?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                duration?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["prefix"], "values">, never>) | undefined;
        range?: ({
            start?: {
                uint?: bigint | undefined;
                int?: bigint | undefined;
                str?: string | undefined;
                bytes?: Uint8Array | undefined;
                enum?: string | undefined;
                bool?: boolean | undefined;
                timestamp?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                duration?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
            }[] | undefined;
            end?: {
                uint?: bigint | undefined;
                int?: bigint | undefined;
                str?: string | undefined;
                bytes?: Uint8Array | undefined;
                enum?: string | undefined;
                bool?: boolean | undefined;
                timestamp?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                duration?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
            }[] | undefined;
        } & {
            start?: ({
                uint?: bigint | undefined;
                int?: bigint | undefined;
                str?: string | undefined;
                bytes?: Uint8Array | undefined;
                enum?: string | undefined;
                bool?: boolean | undefined;
                timestamp?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                duration?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
            }[] & ({
                uint?: bigint | undefined;
                int?: bigint | undefined;
                str?: string | undefined;
                bytes?: Uint8Array | undefined;
                enum?: string | undefined;
                bool?: boolean | undefined;
                timestamp?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                duration?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
            } & {
                uint?: bigint | undefined;
                int?: bigint | undefined;
                str?: string | undefined;
                bytes?: Uint8Array | undefined;
                enum?: string | undefined;
                bool?: boolean | undefined;
                timestamp?: ({
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & Record<Exclude<keyof I["range"]["start"][number]["timestamp"], keyof Timestamp>, never>) | undefined;
                duration?: ({
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & Record<Exclude<keyof I["range"]["start"][number]["duration"], keyof Duration>, never>) | undefined;
            } & Record<Exclude<keyof I["range"]["start"][number], keyof IndexValue>, never>)[] & Record<Exclude<keyof I["range"]["start"], keyof {
                uint?: bigint | undefined;
                int?: bigint | undefined;
                str?: string | undefined;
                bytes?: Uint8Array | undefined;
                enum?: string | undefined;
                bool?: boolean | undefined;
                timestamp?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                duration?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
            }[]>, never>) | undefined;
            end?: ({
                uint?: bigint | undefined;
                int?: bigint | undefined;
                str?: string | undefined;
                bytes?: Uint8Array | undefined;
                enum?: string | undefined;
                bool?: boolean | undefined;
                timestamp?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                duration?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
            }[] & ({
                uint?: bigint | undefined;
                int?: bigint | undefined;
                str?: string | undefined;
                bytes?: Uint8Array | undefined;
                enum?: string | undefined;
                bool?: boolean | undefined;
                timestamp?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                duration?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
            } & {
                uint?: bigint | undefined;
                int?: bigint | undefined;
                str?: string | undefined;
                bytes?: Uint8Array | undefined;
                enum?: string | undefined;
                bool?: boolean | undefined;
                timestamp?: ({
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & Record<Exclude<keyof I["range"]["end"][number]["timestamp"], keyof Timestamp>, never>) | undefined;
                duration?: ({
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & Record<Exclude<keyof I["range"]["end"][number]["duration"], keyof Duration>, never>) | undefined;
            } & Record<Exclude<keyof I["range"]["end"][number], keyof IndexValue>, never>)[] & Record<Exclude<keyof I["range"]["end"], keyof {
                uint?: bigint | undefined;
                int?: bigint | undefined;
                str?: string | undefined;
                bytes?: Uint8Array | undefined;
                enum?: string | undefined;
                bool?: boolean | undefined;
                timestamp?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                duration?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["range"], keyof ListRequest_Range>, never>) | undefined;
        pagination?: ({
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } & {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageRequest>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof ListRequest>, never>>(object: I): ListRequest;
};
export declare const ListRequest_Prefix: {
    typeUrl: string;
    encode(message: ListRequest_Prefix, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ListRequest_Prefix;
    fromJSON(object: any): ListRequest_Prefix;
    toJSON(message: ListRequest_Prefix): unknown;
    fromPartial<I extends {
        values?: {
            uint?: bigint | undefined;
            int?: bigint | undefined;
            str?: string | undefined;
            bytes?: Uint8Array | undefined;
            enum?: string | undefined;
            bool?: boolean | undefined;
            timestamp?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            duration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        }[] | undefined;
    } & {
        values?: ({
            uint?: bigint | undefined;
            int?: bigint | undefined;
            str?: string | undefined;
            bytes?: Uint8Array | undefined;
            enum?: string | undefined;
            bool?: boolean | undefined;
            timestamp?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            duration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        }[] & ({
            uint?: bigint | undefined;
            int?: bigint | undefined;
            str?: string | undefined;
            bytes?: Uint8Array | undefined;
            enum?: string | undefined;
            bool?: boolean | undefined;
            timestamp?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            duration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        } & {
            uint?: bigint | undefined;
            int?: bigint | undefined;
            str?: string | undefined;
            bytes?: Uint8Array | undefined;
            enum?: string | undefined;
            bool?: boolean | undefined;
            timestamp?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["values"][number]["timestamp"], keyof Timestamp>, never>) | undefined;
            duration?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["values"][number]["duration"], keyof Duration>, never>) | undefined;
        } & Record<Exclude<keyof I["values"][number], keyof IndexValue>, never>)[] & Record<Exclude<keyof I["values"], keyof {
            uint?: bigint | undefined;
            int?: bigint | undefined;
            str?: string | undefined;
            bytes?: Uint8Array | undefined;
            enum?: string | undefined;
            bool?: boolean | undefined;
            timestamp?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            duration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "values">, never>>(object: I): ListRequest_Prefix;
};
export declare const ListRequest_Range: {
    typeUrl: string;
    encode(message: ListRequest_Range, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ListRequest_Range;
    fromJSON(object: any): ListRequest_Range;
    toJSON(message: ListRequest_Range): unknown;
    fromPartial<I extends {
        start?: {
            uint?: bigint | undefined;
            int?: bigint | undefined;
            str?: string | undefined;
            bytes?: Uint8Array | undefined;
            enum?: string | undefined;
            bool?: boolean | undefined;
            timestamp?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            duration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        }[] | undefined;
        end?: {
            uint?: bigint | undefined;
            int?: bigint | undefined;
            str?: string | undefined;
            bytes?: Uint8Array | undefined;
            enum?: string | undefined;
            bool?: boolean | undefined;
            timestamp?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            duration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        }[] | undefined;
    } & {
        start?: ({
            uint?: bigint | undefined;
            int?: bigint | undefined;
            str?: string | undefined;
            bytes?: Uint8Array | undefined;
            enum?: string | undefined;
            bool?: boolean | undefined;
            timestamp?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            duration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        }[] & ({
            uint?: bigint | undefined;
            int?: bigint | undefined;
            str?: string | undefined;
            bytes?: Uint8Array | undefined;
            enum?: string | undefined;
            bool?: boolean | undefined;
            timestamp?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            duration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        } & {
            uint?: bigint | undefined;
            int?: bigint | undefined;
            str?: string | undefined;
            bytes?: Uint8Array | undefined;
            enum?: string | undefined;
            bool?: boolean | undefined;
            timestamp?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["start"][number]["timestamp"], keyof Timestamp>, never>) | undefined;
            duration?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["start"][number]["duration"], keyof Duration>, never>) | undefined;
        } & Record<Exclude<keyof I["start"][number], keyof IndexValue>, never>)[] & Record<Exclude<keyof I["start"], keyof {
            uint?: bigint | undefined;
            int?: bigint | undefined;
            str?: string | undefined;
            bytes?: Uint8Array | undefined;
            enum?: string | undefined;
            bool?: boolean | undefined;
            timestamp?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            duration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        }[]>, never>) | undefined;
        end?: ({
            uint?: bigint | undefined;
            int?: bigint | undefined;
            str?: string | undefined;
            bytes?: Uint8Array | undefined;
            enum?: string | undefined;
            bool?: boolean | undefined;
            timestamp?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            duration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        }[] & ({
            uint?: bigint | undefined;
            int?: bigint | undefined;
            str?: string | undefined;
            bytes?: Uint8Array | undefined;
            enum?: string | undefined;
            bool?: boolean | undefined;
            timestamp?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            duration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        } & {
            uint?: bigint | undefined;
            int?: bigint | undefined;
            str?: string | undefined;
            bytes?: Uint8Array | undefined;
            enum?: string | undefined;
            bool?: boolean | undefined;
            timestamp?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["end"][number]["timestamp"], keyof Timestamp>, never>) | undefined;
            duration?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["end"][number]["duration"], keyof Duration>, never>) | undefined;
        } & Record<Exclude<keyof I["end"][number], keyof IndexValue>, never>)[] & Record<Exclude<keyof I["end"], keyof {
            uint?: bigint | undefined;
            int?: bigint | undefined;
            str?: string | undefined;
            bytes?: Uint8Array | undefined;
            enum?: string | undefined;
            bool?: boolean | undefined;
            timestamp?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            duration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof ListRequest_Range>, never>>(object: I): ListRequest_Range;
};
export declare const ListResponse: {
    typeUrl: string;
    encode(message: ListResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ListResponse;
    fromJSON(object: any): ListResponse;
    toJSON(message: ListResponse): unknown;
    fromPartial<I extends {
        results?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        }[] | undefined;
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
        results?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        }[] & ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["results"][number], keyof Any>, never>)[] & Record<Exclude<keyof I["results"], keyof {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        }[]>, never>) | undefined;
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof ListResponse>, never>>(object: I): ListResponse;
};
export declare const IndexValue: {
    typeUrl: string;
    encode(message: IndexValue, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): IndexValue;
    fromJSON(object: any): IndexValue;
    toJSON(message: IndexValue): unknown;
    fromPartial<I extends {
        uint?: bigint | undefined;
        int?: bigint | undefined;
        str?: string | undefined;
        bytes?: Uint8Array | undefined;
        enum?: string | undefined;
        bool?: boolean | undefined;
        timestamp?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
        duration?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
    } & {
        uint?: bigint | undefined;
        int?: bigint | undefined;
        str?: string | undefined;
        bytes?: Uint8Array | undefined;
        enum?: string | undefined;
        bool?: boolean | undefined;
        timestamp?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["timestamp"], keyof Timestamp>, never>) | undefined;
        duration?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["duration"], keyof Duration>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof IndexValue>, never>>(object: I): IndexValue;
};
/** Query is a generic gRPC service for querying ORM data. */
export interface Query {
    /** Get queries an ORM table against an unique index. */
    Get(request: GetRequest): Promise<GetResponse>;
    /** List queries an ORM table against an index. */
    List(request: ListRequest): Promise<ListResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    Get(request: GetRequest): Promise<GetResponse>;
    List(request: ListRequest): Promise<ListResponse>;
}
