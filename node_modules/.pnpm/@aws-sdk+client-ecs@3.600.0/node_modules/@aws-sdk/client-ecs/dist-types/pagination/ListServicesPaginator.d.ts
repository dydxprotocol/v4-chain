import { Paginator } from "@smithy/types";
import { ListServicesCommandInput, ListServicesCommandOutput } from "../commands/ListServicesCommand";
import { ECSPaginationConfiguration } from "./Interfaces";
/**
 * @public
 */
export declare const paginateListServices: (config: ECSPaginationConfiguration, input: ListServicesCommandInput, ...rest: any[]) => Paginator<ListServicesCommandOutput>;
