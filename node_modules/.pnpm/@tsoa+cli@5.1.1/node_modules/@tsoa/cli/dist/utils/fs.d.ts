/// <reference types="node" />
import * as fs from 'fs';
export declare const fsExists: typeof fs.exists.__promisify__;
export declare const fsMkDir: typeof fs.mkdir.__promisify__;
export declare const fsWriteFile: typeof fs.writeFile.__promisify__;
export declare const fsReadFile: typeof fs.readFile.__promisify__;
