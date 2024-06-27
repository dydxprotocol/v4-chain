import { createPaginator } from "@smithy/core";
import { ListVersionsByFunctionCommand, } from "../commands/ListVersionsByFunctionCommand";
import { LambdaClient } from "../LambdaClient";
export const paginateListVersionsByFunction = createPaginator(LambdaClient, ListVersionsByFunctionCommand, "Marker", "NextMarker", "MaxItems");
