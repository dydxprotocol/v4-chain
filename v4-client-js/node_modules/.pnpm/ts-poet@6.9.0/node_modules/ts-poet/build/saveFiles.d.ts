import { Code, ToStringOpts } from "./Code";
export interface CodegenFile {
    name: string;
    contents: Code | string;
    /** Whether to generate the file just once, or overwrite it each time. */
    overwrite: boolean;
    /** Whether to use a `hash` comment hint to conditionally generate files (avoids format cost). */
    hash?: boolean;
    /** File-specific `toString` opts. */
    toStringOpts?: ToStringOpts;
}
export interface SaveFilesOpts {
    /** The tool name, i.e. `joist-codegen` or `ts-proto`; used in file prefix for overwrite: true. Defaults to ts-poet. */
    toolName?: string;
    /** The base directory to output to, defaults to `./`. */
    directory?: string;
    /** The files to generate. */
    files: CodegenFile[];
    /** The default toString opts, i.e. dprint settings, etc. */
    toStringOpts?: ToStringOpts;
}
/**
 * Saves multiple {@link CodegenFile}s.
 */
export declare function saveFiles(opts: SaveFilesOpts): Promise<void>;
