import { createPaginator } from "@smithy/core";
import { ListProvisionedConcurrencyConfigsCommand, } from "../commands/ListProvisionedConcurrencyConfigsCommand";
import { LambdaClient } from "../LambdaClient";
export const paginateListProvisionedConcurrencyConfigs = createPaginator(LambdaClient, ListProvisionedConcurrencyConfigsCommand, "Marker", "NextMarker", "MaxItems");
