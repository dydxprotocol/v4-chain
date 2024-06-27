"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.importClassesFromDirectories = void 0;
const path_1 = require("path");
const glob_1 = require("glob");
/**
 * Loads all exported classes from the given directory.
 */
function importClassesFromDirectories(directories, formats = ['.ts']) {
    const allFiles = directories.reduce((allDirs, dir) => {
        // glob docs says: Please only use forward-slashes in glob expressions.
        // therefore do not do any normalization of dir path
        return allDirs.concat((0, glob_1.sync)(dir));
    }, []);
    return allFiles.filter(file => {
        const dtsExtension = file.substring(file.length - 5, file.length);
        return formats.indexOf((0, path_1.extname)(file)) !== -1 && dtsExtension !== '.d.ts';
    });
}
exports.importClassesFromDirectories = importClassesFromDirectories;
//# sourceMappingURL=importClassesFromDirectories.js.map