import { Paginator } from "@smithy/types";
import { ListFunctionsCommandInput, ListFunctionsCommandOutput } from "../commands/ListFunctionsCommand";
import { LambdaPaginationConfiguration } from "./Interfaces";
/**
 * @public
 */
export declare const paginateListFunctions: (config: LambdaPaginationConfiguration, input: ListFunctionsCommandInput, ...rest: any[]) => Paginator<ListFunctionsCommandOutput>;
