import { RenderOptions, defaultOptions } from "wasm-ast-types";

import { header } from '../utils/header';
import { join } from "path";
import { writeFileSync } from 'fs';
import { sync as mkdirp } from "mkdirp";

import generateMessageComposer from '../generators/message-composer';
import generateTypes from '../generators/types';
import generateReactQuery from '../generators/react-query';
import generateRecoil from '../generators/recoil';
import generateClient from '../generators/client';

import { basename } from 'path';
import { readSchemas } from '../utils';

import deepmerge from 'deepmerge';
import { pascal } from "case";
import { createFileBundle, recursiveModuleBundle } from "../bundler";

import generate from '@babel/generator';
import * as t from '@babel/types';

const defaultOpts: TSBuilderOptions = {
    bundle: {
        enabled: true,
        scope: 'contracts',
        bundleFile: 'bundle.ts'
    }
}

export interface TSBuilderInput {
    contracts: Array<ContractFile | string>;
    outPath: string;
    options?: TSBuilderOptions;
};

export interface BundleOptions {
    enabled?: boolean;
    scope?: string;
    bundleFile?: string;
};

export type TSBuilderOptions = {
    bundle?: BundleOptions;
} & RenderOptions;

export interface BuilderFile {
    type: 'type' | 'client' | 'recoil' | 'react-query' | 'message-composer';
    contract: string;
    localname: string;
    filename: string;
};

export interface ContractFile {
    name: string;
    dir: string;
}
export class TSBuilder {
    contracts: Array<ContractFile | string>;
    outPath: string;
    options?: TSBuilderOptions;

    protected files: BuilderFile[] = [];

    constructor({ contracts, outPath, options }: TSBuilderInput) {
        this.contracts = contracts;
        this.outPath = outPath;
        this.options = deepmerge(
            deepmerge(
                defaultOptions,
                defaultOpts
            ),
            options ?? {}
        );
    }

    getContracts(): ContractFile[] {
        return this.contracts.map(contractOpt => {
            if (typeof contractOpt === 'string') {
                const name = basename(contractOpt);
                const contractName = pascal(name);
                return {
                    name: contractName,
                    dir: contractOpt
                }
            }
            return {
                name: pascal(contractOpt.name),
                dir: contractOpt.dir
            };
        });
    }

    async renderTypes(contract: ContractFile) {
        const { enabled, ...options } = this.options.types;
        if (!enabled) return;
        const contractInfo = await readSchemas({
            schemaDir: contract.dir
        });
        const files = await generateTypes(contract.name, contractInfo, this.outPath, options);
        [].push.apply(this.files, files);
    }

    async renderClient(contract: ContractFile) {
        const { enabled, ...options } = this.options.client;
        if (!enabled) return;
        const contractInfo = await readSchemas({
            schemaDir: contract.dir
        });
        const files = await generateClient(contract.name, contractInfo, this.outPath, options);
        [].push.apply(this.files, files);
    }

    async renderRecoil(contract: ContractFile) {
        const { enabled, ...options } = this.options.recoil;
        if (!enabled) return;
        const contractInfo = await readSchemas({
            schemaDir: contract.dir
        });
        const files = await generateRecoil(contract.name, contractInfo, this.outPath, options);
        [].push.apply(this.files, files);
    }

    async renderReactQuery(contract: ContractFile) {
        const { enabled, ...options } = this.options.reactQuery;
        if (!enabled) return;
        const contractInfo = await readSchemas({
            schemaDir: contract.dir
        });
        const files = await generateReactQuery(contract.name, contractInfo, this.outPath, options);
        [].push.apply(this.files, files);
    }

    async renderMessageComposer(contract: ContractFile) {
        const { enabled, ...options } = this.options.messageComposer;
        if (!enabled) return;
        const contractInfo = await readSchemas({
            schemaDir: contract.dir
        });
        const files = await generateMessageComposer(contract.name, contractInfo, this.outPath, options);
        [].push.apply(this.files, files);
    }

    async build() {
        const contracts = this.getContracts();
        for (let c = 0; c < contracts.length; c++) {
            const contract = contracts[c];
            await this.renderTypes(contract);
            await this.renderClient(contract);
            await this.renderMessageComposer(contract);
            await this.renderReactQuery(contract);
            await this.renderRecoil(contract);
        }
        if (this.options.bundle.enabled) {
            this.bundle();
        }
    }

    async bundle() {

        const allFiles = this.files;

        const bundleFile = this.options.bundle.bundleFile;
        const bundleVariables = {};
        const importPaths = [];

        allFiles.forEach(file => {
            createFileBundle(
                `${this.options.bundle.scope}.${file.contract}`,
                file.localname,
                bundleFile,
                importPaths,
                bundleVariables
            );

        });

        const ast = recursiveModuleBundle(bundleVariables);
        let code = generate(t.program(
            [
                ...importPaths,
                ...ast
            ]
        )).code;

        mkdirp(this.outPath);

        if (code.trim() === '') code = 'export {};'

        writeFileSync(join(this.outPath, bundleFile), header + code);

    }
}