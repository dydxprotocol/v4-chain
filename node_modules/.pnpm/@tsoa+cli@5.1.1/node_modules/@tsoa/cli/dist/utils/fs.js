"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.fsReadFile = exports.fsWriteFile = exports.fsMkDir = exports.fsExists = void 0;
const fs = require("fs");
const util_1 = require("util");
exports.fsExists = (0, util_1.promisify)(fs.exists);
exports.fsMkDir = (0, util_1.promisify)(fs.mkdir);
exports.fsWriteFile = (0, util_1.promisify)(fs.writeFile);
exports.fsReadFile = (0, util_1.promisify)(fs.readFile);
//# sourceMappingURL=fs.js.map