import { Paginator } from "@smithy/types";
import {
  ListTaskDefinitionsCommandInput,
  ListTaskDefinitionsCommandOutput,
} from "../commands/ListTaskDefinitionsCommand";
import { ECSPaginationConfiguration } from "./Interfaces";
export declare const paginateListTaskDefinitions: (
  config: ECSPaginationConfiguration,
  input: ListTaskDefinitionsCommandInput,
  ...rest: any[]
) => Paginator<ListTaskDefinitionsCommandOutput>;
