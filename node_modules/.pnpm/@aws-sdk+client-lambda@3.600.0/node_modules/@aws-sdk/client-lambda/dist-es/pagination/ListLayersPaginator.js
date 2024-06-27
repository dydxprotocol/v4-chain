import { createPaginator } from "@smithy/core";
import { ListLayersCommand } from "../commands/ListLayersCommand";
import { LambdaClient } from "../LambdaClient";
export const paginateListLayers = createPaginator(LambdaClient, ListLayersCommand, "Marker", "NextMarker", "MaxItems");
