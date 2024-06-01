var __assign = (this && this.__assign) || Object.assign || function(t) {
    for (var s, i = 1, n = arguments.length; i < n; i++) {
        s = arguments[i];
        for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p))
            t[p] = s[p];
    }
    return t;
};
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
        define(["require", "exports", "colors", "inquirer"], factory);
    }
})(function (require, exports) {
    "use strict";
    var _this = this;
    Object.defineProperty(exports, "__esModule", { value: true });
    require("colors");
    var inquirer = require("inquirer");
    inquirer.registerPrompt('autocomplete', require('inquirer-autocomplete-prompt'));
    exports.required = function (questions) {
        return questions.map(function (q) {
            if (q.required && !q.validate) {
                q.validate = function (value) {
                    if (!value) {
                        return q.name + " is required";
                    }
                    return true;
                };
            }
            return q;
        });
    };
    exports.names = function (questions) {
        return questions.map(function (q) {
            q.message = "" + '['.white + q.name.blue + ']'.white + " " + q.message.green;
            return q;
        });
    };
    exports.filter = function (questions, answers) {
        var A = questions.map(function (q) { return q.name; });
        var B = Object.keys(answers);
        var diff = A.filter(function (x) { return !B.includes(x); });
        return A.filter(function (n) { return diff.includes(n); }).map(function (name) {
            return questions.find(function (o) { return o.name === name; });
        });
    };
    // converts argv._ into the answers when question specifies it
    exports._filter = function (questions, answers) {
        var _Qs = questions.filter(function (q) { return q.hasOwnProperty('_'); });
        var A = _Qs.map(function (v, i) { return i + ''; });
        var B = Object.keys(answers._ || []);
        var includes = A.filter(function (x) { return B.includes(x); });
        for (var i = 0; i < includes.length; i++) {
            answers[_Qs[i].name] = answers._.shift();
        }
        // now run the filter command if on any questions
        questions.filter(function (q) { return q.hasOwnProperty('filter') && typeof q.filter === 'function'; }).forEach(function (question) {
            if (answers.hasOwnProperty(question.name)) {
                answers[question.name] = question.filter(answers[question.name]);
            }
        });
        return answers;
    };
    exports.prompt = function (questions, answers) { return __awaiter(_this, void 0, void 0, function () {
        var result;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    exports._filter(questions, answers);
                    return [4 /*yield*/, inquirer.prompt(exports.names(exports.required(exports.filter(questions, answers))))];
                case 1:
                    result = _a.sent();
                    return [2 /*return*/, __assign({}, result, answers)];
            }
        });
    }); };
});
//# sourceMappingURL=index.js.map