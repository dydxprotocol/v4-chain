#!/usr/bin/env node
"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g;
    return g = { next: verb(0), "throw": verb(1), "return": verb(2) }, typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (_) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
var minimist = require("minimist");
var get_stdin_1 = __importDefault(require("get-stdin"));
var fs_1 = require("mz/fs");
var mkdirp = __importStar(require("mkdirp"));
var glob_promise_1 = __importDefault(require("glob-promise"));
var isGlob = require("is-glob");
var path_1 = require("path");
var index_1 = require("./index");
var utils_1 = require("./utils");
main(minimist(process.argv.slice(2), {
    alias: {
        help: ['h'],
        input: ['i'],
        output: ['o']
    }
}));
function main(argv) {
    return __awaiter(this, void 0, void 0, function () {
        var argIn, argOut, ISGLOB, ISDIR, result, e_1;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    if (argv.help) {
                        printHelp();
                        process.exit(0);
                    }
                    argIn = argv._[0] || argv.input;
                    argOut = argv._[1] || argv.output // the output can be omitted so this can be undefined
                    ;
                    ISGLOB = isGlob(argIn);
                    ISDIR = isDir(argIn);
                    if ((ISGLOB || ISDIR) && argOut && argOut.includes('.d.ts')) {
                        throw new ReferenceError("You have specified a single file ".concat(argOut, " output for a multi file input ").concat(argIn, ". This feature is not yet supported, refer to issue #272 (https://github.com/bcherny/json-schema-to-typescript/issues/272)"));
                    }
                    _a.label = 1;
                case 1:
                    _a.trys.push([1, 8, , 9]);
                    if (!ISGLOB) return [3 /*break*/, 3];
                    return [4 /*yield*/, processGlob(argIn, argOut, argv)];
                case 2:
                    _a.sent();
                    return [3 /*break*/, 7];
                case 3:
                    if (!ISDIR) return [3 /*break*/, 5];
                    return [4 /*yield*/, processDir(argIn, argOut, argv)];
                case 4:
                    _a.sent();
                    return [3 /*break*/, 7];
                case 5: return [4 /*yield*/, processFile(argIn, argv)];
                case 6:
                    result = _a.sent();
                    outputResult(result, argOut);
                    _a.label = 7;
                case 7: return [3 /*break*/, 9];
                case 8:
                    e_1 = _a.sent();
                    (0, utils_1.error)(e_1);
                    process.exit(1);
                    return [3 /*break*/, 9];
                case 9: return [2 /*return*/];
            }
        });
    });
}
// check if path is an existing directory
function isDir(path) {
    return (0, fs_1.existsSync)(path) && (0, fs_1.lstatSync)(path).isDirectory();
}
function processGlob(argIn, argOut, argv) {
    return __awaiter(this, void 0, void 0, function () {
        var files, results;
        var _this = this;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, (0, glob_promise_1.default)(argIn)]; // execute glob pattern match
                case 1:
                    files = _a.sent() // execute glob pattern match
                    ;
                    if (files.length === 0) {
                        throw ReferenceError("You passed a glob pattern \"".concat(argIn, "\", but there are no files that match that pattern in ").concat(process.cwd()));
                    }
                    return [4 /*yield*/, Promise.all(files.map(function (file) { return __awaiter(_this, void 0, void 0, function () {
                            var _a;
                            return __generator(this, function (_b) {
                                switch (_b.label) {
                                    case 0:
                                        _a = [file];
                                        return [4 /*yield*/, processFile(file, argv)];
                                    case 1: return [2 /*return*/, _a.concat([_b.sent()])];
                                }
                            });
                        }); }))
                        // careful to do this serially
                    ];
                case 2:
                    results = _a.sent();
                    // careful to do this serially
                    results.forEach(function (_a) {
                        var file = _a[0], result = _a[1];
                        var outputPath = argOut && "".concat(argOut, "/").concat((0, path_1.basename)(file, '.json'), ".d.ts");
                        outputResult(result, outputPath);
                    });
                    return [2 /*return*/];
            }
        });
    });
}
function processDir(argIn, argOut, argv) {
    return __awaiter(this, void 0, void 0, function () {
        var files, results;
        var _this = this;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    files = getPaths(argIn);
                    return [4 /*yield*/, Promise.all(files.map(function (file) { return __awaiter(_this, void 0, void 0, function () {
                            var _a, outputPath, _b;
                            return __generator(this, function (_c) {
                                switch (_c.label) {
                                    case 0:
                                        if (!!argOut) return [3 /*break*/, 2];
                                        _a = [file];
                                        return [4 /*yield*/, processFile(file, argv)];
                                    case 1: return [2 /*return*/, _a.concat([_c.sent()])];
                                    case 2:
                                        outputPath = (0, utils_1.pathTransform)(argOut, argIn, file);
                                        _b = [file];
                                        return [4 /*yield*/, processFile(file, argv)];
                                    case 3: return [2 /*return*/, _b.concat([_c.sent(), outputPath])];
                                }
                            });
                        }); }))
                        // careful to do this serially
                    ];
                case 1:
                    results = _a.sent();
                    // careful to do this serially
                    results.forEach(function (_a) {
                        var file = _a[0], result = _a[1], outputPath = _a[2];
                        return outputResult(result, outputPath ? "".concat(outputPath, "/").concat((0, path_1.basename)(file, '.json'), ".d.ts") : undefined);
                    });
                    return [2 /*return*/];
            }
        });
    });
}
function outputResult(result, outputPath) {
    return __awaiter(this, void 0, void 0, function () {
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    if (!!outputPath) return [3 /*break*/, 1];
                    process.stdout.write(result);
                    return [3 /*break*/, 3];
                case 1:
                    if (!isDir((0, path_1.dirname)(outputPath))) {
                        mkdirp.sync((0, path_1.dirname)(outputPath));
                    }
                    return [4 /*yield*/, (0, fs_1.writeFile)(outputPath, result)];
                case 2: return [2 /*return*/, _a.sent()];
                case 3: return [2 /*return*/];
            }
        });
    });
}
function processFile(argIn, argv) {
    return __awaiter(this, void 0, void 0, function () {
        var schema, _a, _b;
        return __generator(this, function (_c) {
            switch (_c.label) {
                case 0:
                    _b = (_a = JSON).parse;
                    return [4 /*yield*/, readInput(argIn)];
                case 1:
                    schema = _b.apply(_a, [_c.sent()]);
                    return [2 /*return*/, (0, index_1.compile)(schema, argIn, argv)];
            }
        });
    });
}
function getPaths(path, paths) {
    if (paths === void 0) { paths = []; }
    if ((0, fs_1.existsSync)(path) && (0, fs_1.lstatSync)(path).isDirectory()) {
        (0, fs_1.readdirSync)((0, path_1.resolve)(path)).forEach(function (item) { return getPaths((0, path_1.join)(path, item), paths); });
    }
    else {
        paths.push(path);
    }
    return paths;
}
function readInput(argIn) {
    if (!argIn) {
        return (0, get_stdin_1.default)();
    }
    return (0, fs_1.readFile)((0, path_1.resolve)(process.cwd(), argIn), 'utf-8');
}
function printHelp() {
    var pkg = require('../../package.json');
    process.stdout.write("\n".concat(pkg.name, " ").concat(pkg.version, "\nUsage: json2ts [--input, -i] [IN_FILE] [--output, -o] [OUT_FILE] [OPTIONS]\n\nWith no IN_FILE, or when IN_FILE is -, read standard input.\nWith no OUT_FILE and when IN_FILE is specified, create .d.ts file in the same directory.\nWith no OUT_FILE nor IN_FILE, write to standard output.\n\nYou can use any of the following options by adding them at the end.\nBoolean values can be set to false using the 'no-' prefix.\n\n  --additionalProperties\n      Default value for additionalProperties, when it is not explicitly set\n  --cwd=XXX\n      Root directory for resolving $ref\n  --declareExternallyReferenced\n      Declare external schemas referenced via '$ref'?\n  --enableConstEnums\n      Prepend enums with 'const'?\n  --format\n      Format code? Set this to false to improve performance.\n  --maxItems\n      Maximum number of unioned tuples to emit when representing bounded-size\n      array types, before falling back to emitting unbounded arrays. Increase\n      this to improve precision of emitted types, decrease it to improve\n      performance, or set it to -1 to ignore minItems and maxItems.\n  --style.XXX=YYY\n      Prettier configuration\n  --unknownAny\n      Output unknown type instead of any type\n  --unreachableDefinitions\n      Generates code for definitions that aren't referenced by the schema\n"));
}
//# sourceMappingURL=cli.js.map