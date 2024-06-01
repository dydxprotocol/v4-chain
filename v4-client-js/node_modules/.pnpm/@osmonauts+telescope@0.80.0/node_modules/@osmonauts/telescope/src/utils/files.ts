import * as t from '@babel/types';
import { parse, ParserPlugin } from '@babel/parser';
import { TelescopeOptions } from '@osmonauts/types';
import { sync as mkdirp } from 'mkdirp';
import { writeFileSync } from 'fs';
import { dirname } from 'path';
import minimatch from 'minimatch';
import generate from '@babel/generator';
import { unused } from './unused';
import traverse from '@babel/traverse';

export const writeAstToFile = (
    outPath: string,
    options: TelescopeOptions,
    program: t.Statement[],
    filename: string
) => {
    const ast = t.program(program);
    const content = generate(ast).code;

    if (options.removeUnusedImports) {
        const plugins: ParserPlugin[] = [
            'typescript'
        ];
        const newAst = parse(content, {
            sourceType: 'module',
            plugins
        });
        traverse(newAst, unused);
        const content2 = generate(newAst).code;
        writeContentToFile(outPath, options, content2, filename);
    } else {
        writeContentToFile(outPath, options, content, filename);
    }
}


export const writeContentToFile = (
    outPath: string,
    options: TelescopeOptions,
    content: string,
    filename: string
) => {
    let esLintPrefix = '';
    let tsLintPrefix = '';

    let nameWithoutPath = filename.replace(outPath, '');
    // strip off leading slash
    if (nameWithoutPath.startsWith('/')) nameWithoutPath = nameWithoutPath.replace(/^\//, '');

    options.tsDisable.patterns.forEach(pattern => {
        if (minimatch(nameWithoutPath, pattern)) {
            tsLintPrefix = `//@ts-nocheck\n`;
        }
    });
    options.eslintDisable.patterns.forEach(pattern => {
        if (minimatch(nameWithoutPath, pattern)) {
            esLintPrefix = `/* eslint-disable */\n`;
        }
    });

    if (
        options.tsDisable.files.includes(nameWithoutPath) ||
        options.tsDisable.disableAll
    ) {
        tsLintPrefix = `//@ts-nocheck\n`;
    }

    if (
        options.eslintDisable.files.includes(nameWithoutPath) ||
        options.eslintDisable.disableAll
    ) {
        esLintPrefix = `/* eslint-disable */\n`;
    }

    const text = tsLintPrefix + esLintPrefix + content;
    mkdirp(dirname(filename));
    writeFileSync(filename, text);
}
