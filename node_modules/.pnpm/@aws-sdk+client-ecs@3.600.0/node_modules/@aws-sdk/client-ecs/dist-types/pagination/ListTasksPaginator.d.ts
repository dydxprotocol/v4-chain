import { Paginator } from "@smithy/types";
import { ListTasksCommandInput, ListTasksCommandOutput } from "../commands/ListTasksCommand";
import { ECSPaginationConfiguration } from "./Interfaces";
/**
 * @public
 */
export declare const paginateListTasks: (config: ECSPaginationConfiguration, input: ListTasksCommandInput, ...rest: any[]) => Paginator<ListTasksCommandOutput>;
