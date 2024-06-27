import { createPaginator } from "@smithy/core";
import { ListFunctionEventInvokeConfigsCommand, } from "../commands/ListFunctionEventInvokeConfigsCommand";
import { LambdaClient } from "../LambdaClient";
export const paginateListFunctionEventInvokeConfigs = createPaginator(LambdaClient, ListFunctionEventInvokeConfigsCommand, "Marker", "NextMarker", "MaxItems");
