import { createPaginator } from "@smithy/core";
import { ListTasksCommand } from "../commands/ListTasksCommand";
import { ECSClient } from "../ECSClient";
export const paginateListTasks = createPaginator(ECSClient, ListTasksCommand, "nextToken", "nextToken", "maxResults");
