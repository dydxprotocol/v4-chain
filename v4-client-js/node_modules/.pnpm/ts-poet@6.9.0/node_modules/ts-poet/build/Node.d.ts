import { ConditionalOutput } from "./ConditionalOutput";
export declare abstract class Node {
    /** Return the unformatted code for this node. */
    abstract toCodeString(used: ConditionalOutput[]): string;
    /** Any potentially string/SymbolSpec/Code nested nodes within us. */
    abstract get childNodes(): unknown[];
}
