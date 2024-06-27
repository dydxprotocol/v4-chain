import { Paginator } from "@smithy/types";
import { ListVersionsByFunctionCommandInput, ListVersionsByFunctionCommandOutput } from "../commands/ListVersionsByFunctionCommand";
import { LambdaPaginationConfiguration } from "./Interfaces";
/**
 * @public
 */
export declare const paginateListVersionsByFunction: (config: LambdaPaginationConfiguration, input: ListVersionsByFunctionCommandInput, ...rest: any[]) => Paginator<ListVersionsByFunctionCommandOutput>;
