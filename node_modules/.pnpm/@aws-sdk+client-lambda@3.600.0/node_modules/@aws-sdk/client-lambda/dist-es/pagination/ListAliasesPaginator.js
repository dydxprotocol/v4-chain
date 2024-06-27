import { createPaginator } from "@smithy/core";
import { ListAliasesCommand } from "../commands/ListAliasesCommand";
import { LambdaClient } from "../LambdaClient";
export const paginateListAliases = createPaginator(LambdaClient, ListAliasesCommand, "Marker", "NextMarker", "MaxItems");
