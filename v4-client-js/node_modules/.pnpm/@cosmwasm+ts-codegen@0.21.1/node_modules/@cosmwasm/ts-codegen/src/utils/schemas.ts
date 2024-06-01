import { sync as glob } from 'glob';
import { readFileSync } from 'fs';
import { cleanse } from './cleanse';
import { compile } from '@pyramation/json-schema-to-typescript';
import { parser } from './parse';
import { ContractInfo, JSONSchema } from 'wasm-ast-types';
interface ReadSchemaOpts {
    schemaDir: string;
    clean?: boolean;
};

export const readSchemas = async ({
    schemaDir, clean = true
}: ReadSchemaOpts): Promise<ContractInfo> => {
    const fn = clean ? cleanse : (str) => str;
    const files = glob(schemaDir + '/**/*.json');
    const schemas = files
        .map(file => JSON.parse(readFileSync(file, 'utf-8')));

    if (schemas.length > 1) {
        // legacy
        // TODO add console.warn here
        return {
            schemas: fn(schemas)
        };
    }

    if (schemas.length === 0) {
        throw new Error('Error [too few files]: requires one schema file per contract');
    }

    if (schemas.length !== 1) {
        throw new Error('Error [too many files]: CosmWasm v1.1 schemas supports one file');
    }

    const idlObject = schemas[0];
    const {
        contract_name,
        contract_version,
        idl_version,
        responses,
        instantiate,
        execute,
        query,
        migrate,
        sudo
    } = idlObject;

    if (typeof idl_version !== 'string') {
        // legacy
        return {
            schemas: fn(schemas)
        };
    }

    // TODO use contract_name, etc.
    return {
        schemas: [
            ...Object.values(fn({
                instantiate,
                execute,
                query,
                migrate,
                sudo
            })).filter(Boolean),
            ...Object.values(fn({ ...responses })).filter(Boolean)
        ],
        responses,
        idlObject
    };
};

export const findQueryMsg = (schemas) => {
    const QueryMsg = schemas.find(schema => schema.title === 'QueryMsg');
    return QueryMsg;
};

export const findExecuteMsg = (schemas) => {
    const ExecuteMsg = schemas.find(schema =>
        schema.title === 'ExecuteMsg' ||
        schema.title === 'ExecuteMsg_for_Empty' || // if cleanse is used, this is never
        schema.title === 'ExecuteMsgForEmpty'
    );
    return ExecuteMsg;
};

export const findAndParseTypes = async (schemas) => {
    const Types = schemas;
    const allTypes = [];
    for (const typ in Types) {
        if (Types[typ].definitions) {
            for (const key of Object.keys(Types[typ].definitions)) {
                // set title
                Types[typ].definitions[key].title = key;
            }
        }
        const result = await compile(Types[typ], Types[typ].title);
        allTypes.push(result);
    }
    const typeHash = parser(allTypes);
    return typeHash;
};
