import { Node } from "./Node";
import { Code } from "./Code";
/**
 * Helps output conditional helper methods.
 *
 * The `ConditionalOutput` concept is split into a usage site and a declaration
 * site, i.e. declaring a `function someHelper() { ... }`, and calling it
 * like `someHelper()`.
 *
 * While generating code, you can make usage sites by using `someHelper` as
 * a placeholder, and then output the declaration with `someHelper.ifUsed`
 * to output the declaration conditionally only if `someHelper` has been
 * seen in the tree.
 *
 * ```typescript
 * const someHelper = conditionalOutput(
 *   "someHelper",
 *   code`function someHelper(n: number) { return n * 2; } `
 * );
 *
 * const code = code`
 *   ${someHelper}(1);
 *
 *   ${someHelper.ifUsed}
 * `
 * ```
 *
 * In the above scenario, it's obvious that `someHelper` is being used, but in
 * code generators with misc configuration options and conditional output paths
 * (i.e. should I output a date helper if dates are even used for this file?)
 * it is harder to tell when exactly a helper should/should not be included.
 */
export declare class ConditionalOutput extends Node {
    usageSiteName: string;
    declarationSiteCode: Code;
    constructor(usageSiteName: string, declarationSiteCode: Code);
    /** Returns the declaration code, typically to be included near the bottom of your output as top-level scope. */
    get ifUsed(): MaybeOutput;
    get childNodes(): unknown[];
    toCodeString(): string;
}
export declare class MaybeOutput {
    parent: ConditionalOutput;
    code: Code;
    constructor(parent: ConditionalOutput, code: Code);
}
