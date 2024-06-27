import { Paginator } from "@smithy/types";
import { ListEventSourceMappingsCommandInput, ListEventSourceMappingsCommandOutput } from "../commands/ListEventSourceMappingsCommand";
import { LambdaPaginationConfiguration } from "./Interfaces";
/**
 * @public
 */
export declare const paginateListEventSourceMappings: (config: LambdaPaginationConfiguration, input: ListEventSourceMappingsCommandInput, ...rest: any[]) => Paginator<ListEventSourceMappingsCommandOutput>;
