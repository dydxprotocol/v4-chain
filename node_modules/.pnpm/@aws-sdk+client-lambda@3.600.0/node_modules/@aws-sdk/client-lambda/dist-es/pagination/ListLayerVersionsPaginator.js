import { createPaginator } from "@smithy/core";
import { ListLayerVersionsCommand, } from "../commands/ListLayerVersionsCommand";
import { LambdaClient } from "../LambdaClient";
export const paginateListLayerVersions = createPaginator(LambdaClient, ListLayerVersionsCommand, "Marker", "NextMarker", "MaxItems");
