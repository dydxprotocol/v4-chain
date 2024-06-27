import { createPaginator } from "@smithy/core";
import { ListServicesCommand, } from "../commands/ListServicesCommand";
import { ECSClient } from "../ECSClient";
export const paginateListServices = createPaginator(ECSClient, ListServicesCommand, "nextToken", "nextToken", "maxResults");
