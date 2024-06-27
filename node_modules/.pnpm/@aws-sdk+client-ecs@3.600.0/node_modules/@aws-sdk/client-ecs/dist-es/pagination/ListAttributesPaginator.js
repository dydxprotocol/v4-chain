import { createPaginator } from "@smithy/core";
import { ListAttributesCommand, } from "../commands/ListAttributesCommand";
import { ECSClient } from "../ECSClient";
export const paginateListAttributes = createPaginator(ECSClient, ListAttributesCommand, "nextToken", "nextToken", "maxResults");
