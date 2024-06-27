import { ProtoRef, ProtoType, ServiceInfo, ALLOWED_RPC_SERVICES } from '@osmonauts/types'
import { getObjectName } from '@osmonauts/proto-parser';
import { getKeyTypeEntryName } from '@osmonauts/ast';
import { getRoot } from './utils';
import { TelescopeParseContext } from './build';

export const parse = (
    context: TelescopeParseContext,
): void => {
    const root = getRoot(context.ref);
    parseRecur({
        context,
        obj: root.root,
        scope: [],
        isNested: false
    });
};

export const getParsedObjectName = (
    ref: ProtoRef,
    obj: any,
    scope: string[],
) => {
    const _scope = [...scope];
    const root = getRoot(ref);
    const allButPackage = _scope.splice(root.package.split('.').length);
    // pull off "this" name
    allButPackage.pop();
    return getObjectName(obj.name, [root.package, ...allButPackage]);
};

// TODO potentially move this back to ast or proto bc the ast lib references MapEntries...
const makeKeyTypeObj = (ref: ProtoRef, field: any, scope: string[]) => {
    const root = getRoot(ref);
    const scoped = [...scope].splice(root.package.split('.').length);
    const adhocObj: ProtoType = {
        type: 'Type',
        comment: undefined,
        fields: {
            key: {
                id: 1,
                type: field.keyType,
                scope: [...scoped],
                parsedType: {
                    name: field.keyType,
                    type: field.keyType
                },
                comment: undefined,
                options: undefined
            },
            value: {
                id: 2,
                type: field.parsedType.name,
                scope: [...scoped],
                parsedType: {
                    name: field.type,
                    type: field.parsedType.type
                },
                comment: undefined,
                options: undefined
            }
        }
    };
    return adhocObj;
}

export const parseType = (
    context: TelescopeParseContext,
    obj: any,
    // scope already has the name of "this" field at the end of it
    scope: string[],
    isNested: boolean = false
) => {

    obj.keyTypes.forEach(field => {
        const keyTypeObject = makeKeyTypeObj(context.ref, field, [...scope]);
        const name = getParsedObjectName(context.ref, {
            name: getKeyTypeEntryName(obj.name, field.name)
        }, [...scope]);
        context.addType(name, keyTypeObject, true);
    });

    // parse nested names
    let name = obj.name;
    if (isNested) {
        name = getParsedObjectName(context.ref, obj, [...scope]);
    }

    context.addType(name, obj, isNested);

    // render nested LAST
    if (obj.nested) {
        Object.keys(obj.nested).forEach(key => {
            // isNested = true;
            parseRecur({
                context,
                obj: obj.nested[key],
                scope: [...scope, key],
                isNested: true
            });
        });
    }

};

export const parseEnum = (
    context: TelescopeParseContext,
    obj: any,
    scope: string[],
    isNested: boolean = false
) => {
    let name = obj.name;
    // parse nested names
    if (isNested) {
        name = getParsedObjectName(context.ref, obj, scope);
    }
    context.addType(name, obj, isNested);
};

export const parseService = (
    context: TelescopeParseContext,
    obj: any,
    scope: string[],
    isNested: boolean = false
) => {

    const methodHash: Record<string, {
        requestType: string;
        responseType: string;
        comment?: string;
    }> = obj.methods;

    if (!ALLOWED_RPC_SERVICES.includes(obj.name)) {
        return;
    }

    Object.entries(methodHash)
        .forEach(([key, value]) => {
            const lookup = context.store.get(context.ref, value.requestType);
            if (!lookup) {
                console.warn(`cannot find ${value.requestType}`);
                throw new Error('undefined symbol for service.');
            }
            const lookupResponse = context.store.get(context.ref, value.responseType);
            if (!lookupResponse) {
                console.warn(`cannot find ${value.requestType}`);
                throw new Error('undefined symbol for service.');
            }
            const serviceInfo: ServiceInfo = {
                methodName: key,
                package: context.ref.proto.package,
                message: lookup.importedName,
                messageImport: lookup.import ?? context.ref.filename,
                response: lookupResponse.importedName,
                responseImport: lookupResponse.import ?? context.ref.filename,
                comment: value.comment
            };
            switch (obj.name) {
                case 'Msg':
                    context.addMutation(serviceInfo);
                    break;
                case 'Query':
                    context.addQuery(serviceInfo);
                    break;
                default:
                    context.addService(serviceInfo);
                    break;
            }
        });
};

interface ParseRecur {
    context: TelescopeParseContext,
    obj: any;
    scope: string[];
    isNested: boolean;
}
export const parseRecur = ({
    context,
    obj,
    scope,
    isNested
}: ParseRecur) => {
    switch (obj.type) {
        case 'Type':
            return parseType(
                context, obj, scope, isNested
            );
        case 'Enum':
            return parseEnum(
                context, obj, scope, isNested
            );
        case 'Service':
            return parseService(
                context, obj, scope, isNested
            );
        case 'Field':
            console.log(obj);
            return;
        case 'Root':
        case 'Namespace':
            if (obj.nested) {
                return Object.keys(obj.nested).forEach(key => {
                    parseRecur({
                        context,
                        obj: obj.nested[key],
                        scope: [...scope, key],
                        isNested
                    });
                });
            } else {
                throw new Error('parseRecur() cannot find protobufjs Type')
            }
        default:
        // if (obj.type === 'string') return;
        // if (obj.type === 'bool') return;
        // if (obj.type === 'HttpRule') return;
        // if (obj.type === 'InterfaceDescriptor') return;
        // if (obj.type === 'ScalarDescriptor') return;
        // if (obj.type === 'ModuleDescriptor') return;
        // if (obj.type === 'TableDescriptor') return;
        // if (obj.type === 'SingletonDescriptor') return;
        // if (obj.type === 'ModuleSchemaDescriptor') return;
        // if (obj.type === 'google.api.FieldBehavior') return;
        // if (obj.type === 'google.api.ResourceReference') return;
        // if (obj.type === 'google.api.ResourceDescriptor') return;
        // if (obj.type === 'google.api.RoutingRule') return;
        // if (obj.type === 'google.api.VisibilityRule') return;
        // if (obj.type === 'google.longrunning.OperationInfo') return;
        // throw new Error('parseRecur() cannot find protobufjs Type')
    }
};