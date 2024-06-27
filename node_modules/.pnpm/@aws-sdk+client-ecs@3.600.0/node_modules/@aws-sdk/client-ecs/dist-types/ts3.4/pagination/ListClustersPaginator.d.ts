import { Paginator } from "@smithy/types";
import {
  ListClustersCommandInput,
  ListClustersCommandOutput,
} from "../commands/ListClustersCommand";
import { ECSPaginationConfiguration } from "./Interfaces";
export declare const paginateListClusters: (
  config: ECSPaginationConfiguration,
  input: ListClustersCommandInput,
  ...rest: any[]
) => Paginator<ListClustersCommandOutput>;
