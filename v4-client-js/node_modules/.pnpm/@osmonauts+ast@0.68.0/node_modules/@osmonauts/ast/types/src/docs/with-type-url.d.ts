import * as t from '@babel/types';
import { ServiceMutation } from "@osmonauts/types";
export declare const documentWithTypeUrl: (mutations: ServiceMutation[]) => t.VariableDeclaration;
export declare const documentWithTypeUrlReadme: (mutations: ServiceMutation[]) => string;
