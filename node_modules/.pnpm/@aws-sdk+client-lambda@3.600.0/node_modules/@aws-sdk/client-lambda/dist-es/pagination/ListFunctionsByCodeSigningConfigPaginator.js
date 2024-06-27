import { createPaginator } from "@smithy/core";
import { ListFunctionsByCodeSigningConfigCommand, } from "../commands/ListFunctionsByCodeSigningConfigCommand";
import { LambdaClient } from "../LambdaClient";
export const paginateListFunctionsByCodeSigningConfig = createPaginator(LambdaClient, ListFunctionsByCodeSigningConfigCommand, "Marker", "NextMarker", "MaxItems");
