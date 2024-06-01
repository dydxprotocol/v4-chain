import { Code } from "ts-poet";
import { Utils } from "./main";
/** Creates a function to transform a message Source to a Uint8Array Source. */
export declare function generateEncodeTransform(utils: Utils, fullName: string): Code;
/** Creates a function to transform a Uint8Array Source to a message Source. */
export declare function generateDecodeTransform(utils: Utils, fullName: string): Code;
