import { Paginator } from "@smithy/types";
import { ListAliasesCommandInput, ListAliasesCommandOutput } from "../commands/ListAliasesCommand";
import { LambdaPaginationConfiguration } from "./Interfaces";
/**
 * @public
 */
export declare const paginateListAliases: (config: LambdaPaginationConfiguration, input: ListAliasesCommandInput, ...rest: any[]) => Paginator<ListAliasesCommandOutput>;
