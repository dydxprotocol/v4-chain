import { Paginator } from "@smithy/types";
import {
  ListFunctionsByCodeSigningConfigCommandInput,
  ListFunctionsByCodeSigningConfigCommandOutput,
} from "../commands/ListFunctionsByCodeSigningConfigCommand";
import { LambdaPaginationConfiguration } from "./Interfaces";
export declare const paginateListFunctionsByCodeSigningConfig: (
  config: LambdaPaginationConfiguration,
  input: ListFunctionsByCodeSigningConfigCommandInput,
  ...rest: any[]
) => Paginator<ListFunctionsByCodeSigningConfigCommandOutput>;
