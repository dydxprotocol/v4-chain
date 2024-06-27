import { createPaginator } from "@smithy/core";
import { ListFunctionUrlConfigsCommand, } from "../commands/ListFunctionUrlConfigsCommand";
import { LambdaClient } from "../LambdaClient";
export const paginateListFunctionUrlConfigs = createPaginator(LambdaClient, ListFunctionUrlConfigsCommand, "Marker", "NextMarker", "MaxItems");
