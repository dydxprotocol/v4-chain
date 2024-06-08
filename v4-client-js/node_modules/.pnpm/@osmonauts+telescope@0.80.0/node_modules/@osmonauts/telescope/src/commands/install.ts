import { resolve, join, dirname, basename, extname } from 'path';
import { sync as mkdirp } from 'mkdirp';
import { sync as glob } from 'glob';
import { sync as rimraf } from 'rimraf';
import { exec } from 'shelljs';
import { prompt } from '../prompt';
import { parse } from 'parse-package-name';
import { tmpdir } from 'os';
import { readFileSync, writeFileSync } from 'fs';

const TMPDIR = tmpdir();
const rnd = () =>
    Math.random().toString(36).substring(2, 15) +
    Math.random().toString(36).substring(2, 15);

const getPackages = (names) => {
    return names.map(pkg => {
        const { name, version } = parse(pkg);
        return `${name}@${version}`
    }).join(' ');
};

export default async (argv) => {

    // don't prompt if we got this...
    if (argv._.length) {
        argv.pkg = argv._;
    }

    // current dir/package
    const cur = process.cwd();
    let thisPackage;
    try {
        thisPackage = JSON.parse(readFileSync(join(cur, 'package.json'), 'utf-8'));
    } catch (e) {
        throw new Error('make sure you are inside of a telescope package!');
    }

    // what are we installing?
    let { pkg } = await prompt([
        {
            type: 'checkbox',
            name: 'pkg',
            message:
                'which packages do you want to support?',
            choices: [
                'akash',
                // 'axelar',
                'bcna',
                // 'bitsong', // Error: namespace and package not supported, yet.
                // 'canto',
                // 'cht',
                'comdex',
                'confio',
                'cosmos',
                'cosmos_proto',
                'cosmwasm',
                // 'cronos',
                // 'cyber',
                // 'desmos',
                // 'echelon',
                // 'em',
                // 'eve',
                // 'evmos',
                // 'galaxy',
                'gogoproto',
                'google',
                'ibc',
                'injective',
                'juno',
                // 'kava',
                // 'likechain',
                // 'lumenx',
                'omniflix',
                // 'oracle',
                'osmosis',
                // 'panacea',
                'passage3d',
                'persistence',
                // 'provenance',
                // 'pylons',
                'regen',
                // 'scavenge',
                // 'sentinel',
                // 'shentu',
                // 'sifnode',
                'stargaze',
                'tendermint',
                'umee',
                // 'uptake',
                // 'vidulum',
            ].map(name => {
                return {
                    name,
                    value: `@protobufs/${name}`
                }
            })
        }
    ], argv);

    // install
    if (!Array.isArray(pkg)) pkg = [pkg];
    const tmp = join(TMPDIR, rnd());
    mkdirp(tmp);
    process.chdir(tmp);
    exec(`npm install ${getPackages(pkg)} --production --prefix ./protobufs`);

    // protos
    const pkgs = glob('./protobufs/**/package.json');
    const cmds = pkgs
        .filter((f) => { return f !== './protobufs/package.json' })
        .map((f) => resolve(join(tmp, f)))
        .map((conf) => {
            const extDir = dirname(conf);
            const dir = extDir.split('node_modules/')[1];
            const dst = basename(dir);

            const files = glob(`${extDir}/**/*`, { nodir: true });
            files.forEach(f => {
                if (extname(f) === '.proto'
                    || f === 'package.json'
                    || /license/i.test(f)
                    || /readme/i.test(f)) return;
                rimraf(f);
            });
            return [extDir, resolve(join(cur, 'proto', dst)), dir];
        });

    // move protos 
    for (const [src, dst, pkg] of cmds) {
        rimraf(dst);
        console.log(`installing ${pkg}...`);
        mkdirp(dirname(dst));
        exec(`mv ${src} ${dst}`);
    }

    // package
    const packageInfo = JSON.parse(readFileSync('./protobufs/package.json', 'utf-8'));
    const deps = packageInfo.dependencies ?? {};
    thisPackage.devDependencies = thisPackage.devDependencies ?? {};
    thisPackage.devDependencies = {
        ...thisPackage.devDependencies,
        ...deps
    };

    thisPackage.devDependencies = Object.fromEntries(Object.entries(thisPackage.devDependencies).sort());

    writeFileSync(join(cur, 'package.json'), JSON.stringify(thisPackage, null, 2));

    // cleanup
    rimraf(tmp);
    process.chdir(cur);
};
