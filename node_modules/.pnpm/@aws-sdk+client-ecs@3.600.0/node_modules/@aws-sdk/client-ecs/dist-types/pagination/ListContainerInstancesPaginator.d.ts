import { Paginator } from "@smithy/types";
import { ListContainerInstancesCommandInput, ListContainerInstancesCommandOutput } from "../commands/ListContainerInstancesCommand";
import { ECSPaginationConfiguration } from "./Interfaces";
/**
 * @public
 */
export declare const paginateListContainerInstances: (config: ECSPaginationConfiguration, input: ListContainerInstancesCommandInput, ...rest: any[]) => Paginator<ListContainerInstancesCommandOutput>;
