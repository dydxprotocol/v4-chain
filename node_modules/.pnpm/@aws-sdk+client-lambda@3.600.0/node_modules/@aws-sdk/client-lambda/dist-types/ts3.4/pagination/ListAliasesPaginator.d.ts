import { Paginator } from "@smithy/types";
import {
  ListAliasesCommandInput,
  ListAliasesCommandOutput,
} from "../commands/ListAliasesCommand";
import { LambdaPaginationConfiguration } from "./Interfaces";
export declare const paginateListAliases: (
  config: LambdaPaginationConfiguration,
  input: ListAliasesCommandInput,
  ...rest: any[]
) => Paginator<ListAliasesCommandOutput>;
