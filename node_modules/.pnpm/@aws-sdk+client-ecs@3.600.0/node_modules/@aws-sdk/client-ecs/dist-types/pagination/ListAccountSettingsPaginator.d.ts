import { Paginator } from "@smithy/types";
import { ListAccountSettingsCommandInput, ListAccountSettingsCommandOutput } from "../commands/ListAccountSettingsCommand";
import { ECSPaginationConfiguration } from "./Interfaces";
/**
 * @public
 */
export declare const paginateListAccountSettings: (config: ECSPaginationConfiguration, input: ListAccountSettingsCommandInput, ...rest: any[]) => Paginator<ListAccountSettingsCommandOutput>;
