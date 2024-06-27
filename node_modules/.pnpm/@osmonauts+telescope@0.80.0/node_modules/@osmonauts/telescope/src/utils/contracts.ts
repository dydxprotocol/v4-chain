import { readSchemas } from '@cosmwasm/ts-codegen';
import { pascal } from 'case';
import { basename, dirname, join } from 'path';
import { readFileSync, readdirSync } from 'fs';

export const getDirectories = source =>
    readdirSync(source, { withFileTypes: true })
        .filter(dirent => dirent.isDirectory())
        .map(dirent => dirent.name);

export const getContracts = () => {
    const contracts = getDirectories('./contracts')
        .map(contractDirname => {
            return {
                name: `${contractDirname}`,
                value: `./contracts/${contractDirname}`
            }
        });
    return contracts;
};

export const getContractSchemata = async (schemata: any[], out: string, argv) => {
    const s = [];
    for (let i = 0; i < schemata.length; i++) {
        const path = schemata[i];
        const pkg = JSON.parse(readFileSync(join(path, 'package.json'), 'utf-8'));
        const name = basename(path);
        const folder = basename(dirname(path));
        const contractName = pascal(pkg.contract) || pascal(name);
        const schemas = await readSchemas({ schemaDir: path, schemaOptions: argv });
        const outPath = join(out, folder);
        s.push({
            contractName, schemas, outPath
        });
    }
    return s;
}