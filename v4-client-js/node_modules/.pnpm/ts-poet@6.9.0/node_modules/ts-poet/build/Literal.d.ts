import { Node } from "./Node";
import { ConditionalOutput } from "./ConditionalOutput";
/**
 * A literal source representation of the provided object.
 */
export declare class Literal extends Node {
    private readonly tokens;
    constructor(object: unknown);
    get childNodes(): unknown[];
    toCodeString(used: ConditionalOutput[]): string;
}
