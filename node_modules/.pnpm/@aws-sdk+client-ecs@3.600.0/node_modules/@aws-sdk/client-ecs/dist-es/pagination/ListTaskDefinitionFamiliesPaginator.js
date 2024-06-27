import { createPaginator } from "@smithy/core";
import { ListTaskDefinitionFamiliesCommand, } from "../commands/ListTaskDefinitionFamiliesCommand";
import { ECSClient } from "../ECSClient";
export const paginateListTaskDefinitionFamilies = createPaginator(ECSClient, ListTaskDefinitionFamiliesCommand, "nextToken", "nextToken", "maxResults");
