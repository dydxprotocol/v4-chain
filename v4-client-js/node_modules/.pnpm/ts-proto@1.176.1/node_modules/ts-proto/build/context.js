"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.createFileContext = void 0;
function createFileContext(file) {
    return { isProto3Syntax: file.syntax === "proto3" };
}
exports.createFileContext = createFileContext;
