import { Paginator } from "@smithy/types";
import { ListAttributesCommandInput, ListAttributesCommandOutput } from "../commands/ListAttributesCommand";
import { ECSPaginationConfiguration } from "./Interfaces";
/**
 * @public
 */
export declare const paginateListAttributes: (config: ECSPaginationConfiguration, input: ListAttributesCommandInput, ...rest: any[]) => Paginator<ListAttributesCommandOutput>;
