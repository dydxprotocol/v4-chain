import generate from '@babel/generator';
import cases from 'jest-in-case'
import * as t from '@babel/types';
import {
    recursiveModuleBundle
} from './bundle'
import { getGenericParseContext } from '../../test-utils'


const preview = (ast) => {
    return generate(t.program(ast)).code;
}

const context = getGenericParseContext();

cases(`recursiveModuleBundle`, opts => {

    expect(preview(recursiveModuleBundle(context.options, opts.data))).toMatchSnapshot();
}, [
    {
        name: 'root',
        data: {
            osmosis: {
                __export: true,
                _0: true
            },
            tendermint: {
                __export: true,
                _0: true
            }
        }
    },
    {
        name: 'single',
        data: {
            ics23: {
                __export: true,
                _0: true
            }
        }
    },
    {
        name: 'tendermint',
        data: {
            tendermint: {
                abci: {
                    __export: true,
                    _223: true
                },
                crypto: {
                    __export: true,
                    _224: true,
                    _225: true
                },
                libs: {
                    bits: {
                        __export: true,
                        _226: true
                    }
                },
                p2p: {
                    __export: true,
                    _227: true
                },
                types: {
                    __export: true,
                    _228: true,
                    _229: true,
                    _230: true,
                    _231: true,
                    _232: true
                },
                version: {
                    __export: true,
                    _233: true
                }
            }
        }
    },
    {
        name: 'osmo',
        data: {
            osmosis: {
                claim: {
                    v1beta1: {
                        __export: true,
                        _177: true,
                        _178: true,
                        _179: true,
                        _180: true
                    }
                },
                epochs: {
                    v1beta1: {
                        __export: true,
                        _181: true,
                        _182: true
                    }
                },
                gamm: {
                    v1beta1: {
                        __export: true,
                        _183: true,
                        _184: true,
                        _185: true,
                        _186: true,
                        _187: true
                    }
                },
                incentives: {
                    __export: true,
                    _188: true,
                    _189: true,
                    _190: true,
                    _191: true,
                    _192: true
                },
                lockup: {
                    __export: true,
                    _193: true,
                    _194: true,
                    _195: true,
                    _196: true
                },
                mint: {
                    v1beta1: {
                        __export: true,
                        _197: true,
                        _198: true,
                        _199: true
                    }
                },
                poolincentives: {
                    v1beta1: {
                        __export: true,
                        _200: true,
                        _201: true,
                        _202: true,
                        _203: true
                    }
                },
                store: {
                    v1beta1: {
                        __export: true,
                        _204: true
                    }
                },
                superfluid: {
                    __export: true,
                    _205: true,
                    _206: true,
                    _207: true,
                    _208: true,
                    _209: true
                },
                txfees: {
                    v1beta1: {
                        __export: true,
                        _210: true,
                        _211: true,
                        _212: true,
                        _213: true
                    }
                }
            }
        }
    }
])