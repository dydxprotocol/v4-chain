var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : new P(function (resolve) { resolve(result.value); }).then(fulfilled, rejected); }
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
            if (f = 1, y && (t = y[op[0] & 2 ? "return" : op[0] ? "throw" : "next"]) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [0, t.value];
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
(function (factory) {
    if (typeof module === "object" && typeof module.exports === "object") {
        var v = factory(require, exports);
        if (v !== undefined) module.exports = v;
    }
    else if (typeof define === "function" && define.amd) {
        define(["require", "exports", "../index"], factory);
    }
})(function (require, exports) {
    "use strict";
    var _this = this;
    Object.defineProperty(exports, "__esModule", { value: true });
    var index_1 = require("../index");
    describe('arguments', function () {
        it('empty when all args supplied', function () {
            var questions = [
                {
                    name: 'hello',
                    message: 'hello',
                },
                {
                    name: 'world',
                    message: 'world',
                },
            ];
            var argv = {
                hello: 1,
                world: 2,
            };
            expect(index_1.filter(questions, argv)).toEqual([]);
            expect(argv).toEqual({
                hello: 1,
                world: 2,
            });
        });
        it('empty when all args supplied II', function () {
            var questions = [
                {
                    _: true,
                    name: 'foo',
                    message: 'foo',
                },
                {
                    name: 'bar',
                    message: 'bar',
                },
                {
                    _: true,
                    name: 'baz',
                    message: 'baz',
                },
            ];
            var argv = {
                _: [1, 3],
                bar: 2,
            };
            var _1 = index_1.filter(questions, argv);
            var _2 = index_1._filter(questions, argv);
            expect(_2).toEqual({ _: [], bar: 2, baz: 3, foo: 1 });
            expect(argv).toEqual({
                _: [],
                foo: 1,
                bar: 2,
                baz: 3,
            });
        });
        it('init example', function () { return __awaiter(_this, void 0, void 0, function () {
            var questions, argv, _1, _2;
            return __generator(this, function (_a) {
                questions = [
                    {
                        _: true,
                        name: 'foo',
                        message: '',
                    },
                ];
                argv = {
                    _: [],
                    bar: 2,
                };
                _1 = index_1._filter(questions, argv);
                _2 = index_1.filter(questions, argv);
                expect(_2).toEqual([
                    {
                        _: true,
                        name: 'foo',
                        message: '',
                    },
                ]);
                expect(argv).toEqual({
                    _: [],
                    bar: 2,
                });
                return [2 /*return*/];
            });
        }); });
        it('basic example', function () { return __awaiter(_this, void 0, void 0, function () {
            var questions, argv, _2;
            return __generator(this, function (_a) {
                questions = [
                    {
                        name: 'name',
                        message: 'project name (e.g., flipr)',
                        required: true,
                    },
                ];
                argv = { _: [], cmd: 'init' };
                index_1._filter(questions, argv);
                _2 = index_1.filter(questions, argv);
                expect(_2).toEqual([
                    {
                        name: 'name',
                        message: 'project name (e.g., flipr)',
                        required: true,
                    },
                ]);
                expect(argv).toEqual({ _: [], cmd: 'init' });
                return [2 /*return*/];
            });
        }); });
    });
    describe('prompt', function () {
        it('empty when all args supplied', function () { return __awaiter(_this, void 0, void 0, function () {
            var questions, argv, value;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        questions = [
                            {
                                name: 'hello',
                                message: '',
                            },
                            {
                                name: 'world',
                                message: '',
                            },
                        ];
                        argv = {
                            hello: 1,
                            world: 2,
                        };
                        return [4 /*yield*/, index_1.prompt(questions, argv)];
                    case 1:
                        value = _a.sent();
                        expect(value).toEqual({
                            hello: 1,
                            world: 2,
                        });
                        return [2 /*return*/];
                }
            });
        }); });
        it('empty when all args supplied', function () { return __awaiter(_this, void 0, void 0, function () {
            var questions, argv, value;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        questions = [
                            {
                                _: true,
                                name: 'foo',
                                message: '',
                            },
                            {
                                name: 'bar',
                                message: '',
                            },
                            {
                                _: true,
                                name: 'baz',
                                message: '',
                            },
                        ];
                        argv = {
                            _: [1, 3],
                            bar: 2,
                        };
                        return [4 /*yield*/, index_1.prompt(questions, argv)];
                    case 1:
                        value = _a.sent();
                        expect(argv).toEqual({
                            _: [],
                            foo: 1,
                            bar: 2,
                            baz: 3,
                        });
                        expect(value).toEqual({
                            _: [],
                            foo: 1,
                            bar: 2,
                            baz: 3,
                        });
                        return [2 /*return*/];
                }
            });
        }); });
        it('basic example', function () { return __awaiter(_this, void 0, void 0, function () {
            var questions, argv, value;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        questions = [
                            {
                                name: 'cmd',
                                message: 'project name (e.g., flipr)',
                                required: true,
                            },
                        ];
                        argv = { _: [], cmd: 'init' };
                        return [4 /*yield*/, index_1.prompt(questions, argv)];
                    case 1:
                        value = _a.sent();
                        expect(value).toEqual(argv);
                        return [2 /*return*/];
                }
            });
        }); });
    });
    describe('filter', function () {
        it('runs filter without _', function () { return __awaiter(_this, void 0, void 0, function () {
            var questions, argv, value;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        questions = [
                            {
                                name: 'hello',
                                message: '',
                                filter: function (val) {
                                    return val + '!';
                                }
                            },
                            {
                                name: 'world',
                                message: '',
                                filter: function (val) {
                                    return val + '!';
                                }
                            },
                        ];
                        argv = {
                            hello: 1,
                            world: 2,
                        };
                        return [4 /*yield*/, index_1.prompt(questions, argv)];
                    case 1:
                        value = _a.sent();
                        expect(value).toEqual({
                            hello: '1!',
                            world: '2!',
                        });
                        return [2 /*return*/];
                }
            });
        }); });
        it('runs filter with _', function () { return __awaiter(_this, void 0, void 0, function () {
            var questions, argv, value;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        questions = [
                            {
                                _: true,
                                name: 'hello',
                                message: '',
                                filter: function (val) {
                                    return val + '!';
                                }
                            },
                            {
                                name: 'world',
                                message: '',
                                filter: function (val) {
                                    return val + '!';
                                }
                            },
                        ];
                        argv = {
                            _: [1],
                            world: 2,
                        };
                        return [4 /*yield*/, index_1.prompt(questions, argv)];
                    case 1:
                        value = _a.sent();
                        expect(value).toEqual({
                            _: [],
                            hello: '1!',
                            world: '2!',
                        });
                        return [2 /*return*/];
                }
            });
        }); });
    });
});
//# sourceMappingURL=prompt.test.js.map