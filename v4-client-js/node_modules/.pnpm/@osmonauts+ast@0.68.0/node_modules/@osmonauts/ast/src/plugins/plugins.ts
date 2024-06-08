import { TelescopeOptions, TelescopeOption } from '@osmonauts/types';
import * as dotty from 'dotty';

const getAllPackageParts = (name: string, list?: string[]) => {
    if (!list) list = [name];
    const newParts = name.split('.');
    newParts.pop();
    if (!newParts.length) return [...list];
    const newName = newParts.join('.');
    return getAllPackageParts(newName, [...list, newName]);
};

export const getPluginValue = (optionName: TelescopeOption | string, currentPkg: string, options: TelescopeOptions) => {
    const pkgOpts = options.packages;
    let value;
    getAllPackageParts(currentPkg).some((pkg, i) => {
        if (dotty.exists(pkgOpts, pkg)) {
            const obj = dotty.get(pkgOpts, pkg);
            if (dotty.exists(obj, optionName)) {
                value = dotty.get(obj, optionName);
                return true;
            }
        }
    });

    if (value === undefined) {
        const defaultValue = dotty.exists(options, optionName) ? dotty.get(options, optionName) : undefined;
        value = defaultValue;
    }
    return value;
};
