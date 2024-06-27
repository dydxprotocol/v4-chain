/**
 * Copyright 2016 Google Inc. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
import * as sourceMap from 'source-map';
export interface MapInfoCompiled {
    mapFileDir: string;
    mapConsumer: sourceMap.RawSourceMap;
}
export interface GeneratedLocation {
    file: string;
    name?: string;
    line: number;
    column: number;
}
export interface SourceLocation {
    file?: string;
    name?: string;
    line?: number;
    column?: number;
}
export declare class SourceMapper {
    infoMap: Map<string, MapInfoCompiled>;
    debug: boolean;
    static create(searchDirs: string[], debug?: boolean): Promise<SourceMapper>;
    /**
     * @param {Array.<string>} sourceMapPaths An array of paths to .map source map
     *  files that should be processed.  The paths should be relative to the
     *  current process's current working directory
     * @param {Logger} logger A logger that reports errors that occurred while
     *  processing the given source map files
     * @constructor
     */
    constructor(debug?: boolean);
    /**
     * Used to get the information about the transpiled file from a given input
     * source file provided there isn't any ambiguity with associating the input
     * path to exactly one output transpiled file.
     *
     * @param inputPath The (possibly relative) path to the original source file.
     * @return The `MapInfoCompiled` object that describes the transpiled file
     *  associated with the specified input path.  `null` is returned if either
     *  zero files are associated with the input path or if more than one file
     *  could possibly be associated with the given input path.
     */
    private getMappingInfo;
    /**
     * Used to determine if the source file specified by the given path has
     * a .map file and an output file associated with it.
     *
     * If there is no such mapping, it could be because the input file is not
     * the input to a transpilation process or it is the input to a transpilation
     * process but its corresponding .map file was not given to the constructor
     * of this mapper.
     *
     * @param {string} inputPath The path to an input file that could
     *  possibly be the input to a transpilation process.  The path should be
     *  relative to the process's current working directory.
     */
    hasMappingInfo(inputPath: string): boolean;
    /**
     * @param {string} inputPath The path to an input file that could possibly
     *  be the input to a transpilation process.  The path should be relative to
     *  the process's current working directory
     * @param {number} The line number in the input file where the line number is
     *   zero-based.
     * @param {number} (Optional) The column number in the line of the file
     *   specified where the column number is zero-based.
     * @return {Object} The object returned has a "file" attribute for the
     *   path of the output file associated with the given input file (where the
     *   path is relative to the process's current working directory),
     *   a "line" attribute of the line number in the output file associated with
     *   the given line number for the input file, and an optional "column" number
     *   of the column number of the output file associated with the given file
     *   and line information.
     *
     *   If the given input file does not have mapping information associated
     *   with it then the input location is returned.
     */
    mappingInfo(location: GeneratedLocation): SourceLocation;
}
