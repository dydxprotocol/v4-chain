import { createPaginator } from "@smithy/core";
import { ListAccountSettingsCommand, } from "../commands/ListAccountSettingsCommand";
import { ECSClient } from "../ECSClient";
export const paginateListAccountSettings = createPaginator(ECSClient, ListAccountSettingsCommand, "nextToken", "nextToken", "maxResults");
