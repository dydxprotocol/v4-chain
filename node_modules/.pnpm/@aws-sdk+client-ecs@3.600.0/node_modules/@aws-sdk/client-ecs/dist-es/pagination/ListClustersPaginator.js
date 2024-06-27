import { createPaginator } from "@smithy/core";
import { ListClustersCommand, } from "../commands/ListClustersCommand";
import { ECSClient } from "../ECSClient";
export const paginateListClusters = createPaginator(ECSClient, ListClustersCommand, "nextToken", "nextToken", "maxResults");
