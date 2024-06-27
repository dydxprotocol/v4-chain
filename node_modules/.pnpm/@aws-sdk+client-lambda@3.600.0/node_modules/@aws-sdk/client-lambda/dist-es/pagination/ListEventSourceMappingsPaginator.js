import { createPaginator } from "@smithy/core";
import { ListEventSourceMappingsCommand, } from "../commands/ListEventSourceMappingsCommand";
import { LambdaClient } from "../LambdaClient";
export const paginateListEventSourceMappings = createPaginator(LambdaClient, ListEventSourceMappingsCommand, "Marker", "NextMarker", "MaxItems");
