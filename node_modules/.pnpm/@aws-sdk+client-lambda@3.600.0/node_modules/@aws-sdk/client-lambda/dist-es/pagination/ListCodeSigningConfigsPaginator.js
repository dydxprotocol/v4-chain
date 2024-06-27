import { createPaginator } from "@smithy/core";
import { ListCodeSigningConfigsCommand, } from "../commands/ListCodeSigningConfigsCommand";
import { LambdaClient } from "../LambdaClient";
export const paginateListCodeSigningConfigs = createPaginator(LambdaClient, ListCodeSigningConfigsCommand, "Marker", "NextMarker", "MaxItems");
