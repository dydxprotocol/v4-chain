import * as t from '@babel/types';
import { ProtoService, ProtoServiceMethodInfo } from '@osmonauts/types';
import { GenericParseContext } from '../../../encoding';
export declare const getUrlTemplateString: (url: string) => {
    strs: any[];
    atEnd: boolean;
};
export declare const makeTemplateTag: (info: ProtoServiceMethodInfo) => t.TemplateLiteral;
export declare const makeTemplateTagLegacy: (info: ProtoServiceMethodInfo) => t.TemplateLiteral;
export declare const createLCDClient: (context: GenericParseContext, service: ProtoService) => t.ExportNamedDeclaration;
export declare const createAggregatedLCDClient: (context: GenericParseContext, services: ProtoService[], clientName: string) => t.ExportNamedDeclaration;
