import { createPaginator } from "@smithy/core";
import { ListServicesByNamespaceCommand, } from "../commands/ListServicesByNamespaceCommand";
import { ECSClient } from "../ECSClient";
export const paginateListServicesByNamespace = createPaginator(ECSClient, ListServicesByNamespaceCommand, "nextToken", "nextToken", "maxResults");
