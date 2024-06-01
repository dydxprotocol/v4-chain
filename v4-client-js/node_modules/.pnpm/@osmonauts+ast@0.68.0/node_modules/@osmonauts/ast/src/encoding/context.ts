import { TelescopeOptions, ProtoField, ProtoRef } from '@osmonauts/types';
import { ProtoStore, getObjectName } from '@osmonauts/proto-parser';
import { getEnumFromJsonName, getEnumToJsonName, getFieldsTypeName } from './proto';
import { getPluginValue } from '../plugins';
export interface ParseContext {
    options: TelescopeOptions;
    imports: ImportUsage[];
    utils: Record<string, boolean>;
    addUtil: Function;
}

export interface ImportUsage {
    type: 'typeImport' | 'toJSONEnum' | 'fromJSONEnum';
    name: string;
    import: string;
    importedAs?: string;
}
export class GenericParseContext implements ParseContext {
    options: TelescopeOptions;
    imports: ImportUsage[] = [];
    utils: Record<string, boolean> = {};
    store: ProtoStore;
    ref: ProtoRef;

    constructor(
        ref: ProtoRef,
        store: ProtoStore,
        options: TelescopeOptions
    ) {
        this.ref = ref;
        this.store = store;
        this.options = options;
        if (!this.options) {
            throw new Error('ParseContext requires options!');
        }

    }

    pluginValue(name) {
        return getPluginValue(name, this.ref.proto.package, this.options);
    }

    isExcluded() {
        const packages = this.pluginValue('prototypes.excluded.packages') ?? [];
        const protos = this.pluginValue('prototypes.excluded.protos') ?? [];
        return packages.includes(this.ref.proto.package) || protos.includes(this.ref.filename);
    }

    addUtil(util) {
        this.utils[util] = true;
    }

    addImport(imp: ImportUsage) {
        // some local lookups don't have an import (local proto-style lookups do)
        if (!imp.import) return;
        this.imports.push(imp)
    }

    getTypeName(field: ProtoField) {
        let name = getFieldsTypeName(field);
        let importedAs = name;
        const names = this.ref.traversed?.importNames;
        if (names
            && names.hasOwnProperty(field.import)
            && names[field.import].hasOwnProperty(name)
        ) {

            importedAs = names[field.import][name];
        }
        this.addImport({
            type: 'typeImport',
            name,
            importedAs,
            import: field.import
        })
        return importedAs;
    }

}

export class AminoParseContext extends GenericParseContext implements ParseContext {

    aminoCasingFn: Function;

    constructor(
        ref: ProtoRef,
        store: ProtoStore,
        options: TelescopeOptions
    ) {
        super(ref, store, options);
        this.ref = ref;
        this.store = store;
        this.options = options;

        this.setAminoCasingFn();

        if (!this.aminoCasingFn) {
            throw new Error('missing aminoCasingFn!')
        }
        this.aminoCaseField = this.aminoCaseField.bind(this);
    }

    private setAminoCasingFn() {
        if (this.aminoCasingFn) return this.aminoCasingFn;
        this.aminoCasingFn = this.pluginValue('aminoEncoding.casingFn');
        return this.aminoCasingFn;
    }

    aminoCaseField(field: ProtoField) {
        return field.options['(telescope:orig)'];
    }

    private lookupTypeFromCurrentPath(field: ProtoField, currentProtoPath: string) {
        const ref = this.store.findProto(currentProtoPath);
        let lookup = this.store.get(ref, field.parsedType.name);
        if (!lookup) {

            // if we can't find it, use field import
            if (field.import) {
                const importRef = this.store.findProto(field.import);
                if (!importRef) {
                    throw new Error(`bad import ${field.import}`);
                }
                lookup = this.store.get(importRef, field.parsedType.name);
            }

            if (!lookup) {
                throw new Error('Undefined Symbol: ' + field.parsedType.name);
            }
        }
        return lookup;
    }

    getTypeFromCurrentPath(field: ProtoField, currentProtoPath: string) {
        const ref = this.store.findProto(currentProtoPath);
        let lookup = this.store.get(ref, field.parsedType.name);
        if (!lookup) {

            // if we can't find it, use field import
            if (field.import) {
                const importRef = this.store.findProto(field.import);
                if (!importRef) {
                    throw new Error(`bad import ${field.import}`);
                }
                lookup = this.store.get(importRef, field.parsedType.name);
            }

            if (!lookup) {
                throw new Error('Undefined Symbol: ' + field.parsedType.name);
            }
        }

        this.addImport(
            {
                type: 'typeImport',
                name: lookup.importedName,
                import: lookup.import
            }
        )

        return lookup.obj;
    }

    lookupEnumFromJson(field: ProtoField, currentProtoPath: string) {
        const lookup = this.lookupTypeFromCurrentPath(field, currentProtoPath);
        const Enum = lookup.obj;
        const name = getEnumFromJsonName(getObjectName(Enum.name, Enum.scope));
        this.addImport({
            type: 'fromJSONEnum',
            name,
            import: lookup.import
        });
        return name;
    }

    lookupEnumToJson(field: ProtoField, currentProtoPath: string) {
        const lookup = this.lookupTypeFromCurrentPath(field, currentProtoPath);
        const Enum = lookup.obj;
        const name = getEnumToJsonName(getObjectName(Enum.name, Enum.scope));
        this.addImport({
            type: 'toJSONEnum',
            name,
            import: lookup.import
        });
        return name;
    }

}
export class ProtoParseContext extends GenericParseContext implements ParseContext {

    constructor(
        ref: ProtoRef,
        store: ProtoStore,
        options: TelescopeOptions
    ) {
        super(ref, store, options);
        this.ref = ref;
        this.store = store;
        this.options = options;
    }

    getToEnum(field: ProtoField) {
        const name = getEnumToJsonName(getFieldsTypeName(field));
        this.addImport({
            type: 'toJSONEnum',
            name,
            import: field.import
        });
        return name;
    }

    getFromEnum(field: ProtoField) {
        const fromJSONFuncName = getEnumFromJsonName(getFieldsTypeName(field));
        this.addImport({
            type: 'fromJSONEnum',
            name: fromJSONFuncName,
            import: field.import
        });
        return fromJSONFuncName;
    }

}

