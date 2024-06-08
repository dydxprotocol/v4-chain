import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../../../helpers";
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
/** Config is the config object of the x/auth/tx package. */
export interface ConfigSDKType {
    skip_ante_handler: boolean;
    skip_post_handler: boolean;
}
export declare const Config: {
    encode(message: Config, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Config;
    fromPartial(object: DeepPartial<Config>): Config;
};
