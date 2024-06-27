import { createPaginator } from "@smithy/core";
import { ListTaskDefinitionsCommand, } from "../commands/ListTaskDefinitionsCommand";
import { ECSClient } from "../ECSClient";
export const paginateListTaskDefinitions = createPaginator(ECSClient, ListTaskDefinitionsCommand, "nextToken", "nextToken", "maxResults");
