import { Paginator } from "@smithy/types";
import { ListProvisionedConcurrencyConfigsCommandInput, ListProvisionedConcurrencyConfigsCommandOutput } from "../commands/ListProvisionedConcurrencyConfigsCommand";
import { LambdaPaginationConfiguration } from "./Interfaces";
/**
 * @public
 */
export declare const paginateListProvisionedConcurrencyConfigs: (config: LambdaPaginationConfiguration, input: ListProvisionedConcurrencyConfigsCommandInput, ...rest: any[]) => Paginator<ListProvisionedConcurrencyConfigsCommandOutput>;
