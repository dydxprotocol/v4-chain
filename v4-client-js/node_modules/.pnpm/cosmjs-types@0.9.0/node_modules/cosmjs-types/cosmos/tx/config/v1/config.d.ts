import { BinaryReader, BinaryWriter } from "../../../../binary";
export declare const protobufPackage = "cosmos.tx.config.v1";
/** Config is the config object of the x/auth/tx package. */
export interface Config {
    /**
     * skip_ante_handler defines whether the ante handler registration should be skipped in case an app wants to override
     * this functionality.
     */
    skipAnteHandler: boolean;
    /**
     * skip_post_handler defines whether the post handler registration should be skipped in case an app wants to override
     * this functionality.
     */
    skipPostHandler: boolean;
}
export declare const Config: {
    typeUrl: string;
    encode(message: Config, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Config;
    fromJSON(object: any): Config;
    toJSON(message: Config): unknown;
    fromPartial<I extends {
        skipAnteHandler?: boolean | undefined;
        skipPostHandler?: boolean | undefined;
    } & {
        skipAnteHandler?: boolean | undefined;
        skipPostHandler?: boolean | undefined;
    } & Record<Exclude<keyof I, keyof Config>, never>>(object: I): Config;
};
