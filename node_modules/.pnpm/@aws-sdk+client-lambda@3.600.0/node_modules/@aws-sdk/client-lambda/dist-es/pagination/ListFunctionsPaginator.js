import { createPaginator } from "@smithy/core";
import { ListFunctionsCommand, } from "../commands/ListFunctionsCommand";
import { LambdaClient } from "../LambdaClient";
export const paginateListFunctions = createPaginator(LambdaClient, ListFunctionsCommand, "Marker", "NextMarker", "MaxItems");
