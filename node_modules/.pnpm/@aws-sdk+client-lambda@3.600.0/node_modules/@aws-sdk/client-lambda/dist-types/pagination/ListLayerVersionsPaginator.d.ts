import { Paginator } from "@smithy/types";
import { ListLayerVersionsCommandInput, ListLayerVersionsCommandOutput } from "../commands/ListLayerVersionsCommand";
import { LambdaPaginationConfiguration } from "./Interfaces";
/**
 * @public
 */
export declare const paginateListLayerVersions: (config: LambdaPaginationConfiguration, input: ListLayerVersionsCommandInput, ...rest: any[]) => Paginator<ListLayerVersionsCommandOutput>;
