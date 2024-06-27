import { createPaginator } from "@smithy/core";
import { ListContainerInstancesCommand, } from "../commands/ListContainerInstancesCommand";
import { ECSClient } from "../ECSClient";
export const paginateListContainerInstances = createPaginator(ECSClient, ListContainerInstancesCommand, "nextToken", "nextToken", "maxResults");
