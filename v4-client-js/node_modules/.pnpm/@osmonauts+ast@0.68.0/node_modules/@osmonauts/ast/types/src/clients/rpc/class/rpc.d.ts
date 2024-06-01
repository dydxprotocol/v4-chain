import * as t from '@babel/types';
import { ProtoService } from '@osmonauts/types';
import { GenericParseContext } from '../../../encoding';
export declare const createRpcClientInterface: (context: GenericParseContext, service: ProtoService) => t.ExportNamedDeclaration;
export declare const getRpcClassName: (service: ProtoService) => string;
export declare const createRpcClientClass: (context: GenericParseContext, service: ProtoService) => t.ExportNamedDeclaration;
export declare const createRpcInterface: (context: GenericParseContext, service: ProtoService) => t.TSInterfaceDeclaration;
