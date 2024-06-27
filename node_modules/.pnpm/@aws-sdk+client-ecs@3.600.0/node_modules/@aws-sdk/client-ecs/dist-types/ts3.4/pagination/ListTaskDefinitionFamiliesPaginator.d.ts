import { Paginator } from "@smithy/types";
import {
  ListTaskDefinitionFamiliesCommandInput,
  ListTaskDefinitionFamiliesCommandOutput,
} from "../commands/ListTaskDefinitionFamiliesCommand";
import { ECSPaginationConfiguration } from "./Interfaces";
export declare const paginateListTaskDefinitionFamilies: (
  config: ECSPaginationConfiguration,
  input: ListTaskDefinitionFamiliesCommandInput,
  ...rest: any[]
) => Paginator<ListTaskDefinitionFamiliesCommandOutput>;
