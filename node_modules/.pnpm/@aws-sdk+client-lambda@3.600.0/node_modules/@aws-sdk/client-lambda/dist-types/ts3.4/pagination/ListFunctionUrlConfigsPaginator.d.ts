import { Paginator } from "@smithy/types";
import {
  ListFunctionUrlConfigsCommandInput,
  ListFunctionUrlConfigsCommandOutput,
} from "../commands/ListFunctionUrlConfigsCommand";
import { LambdaPaginationConfiguration } from "./Interfaces";
export declare const paginateListFunctionUrlConfigs: (
  config: LambdaPaginationConfiguration,
  input: ListFunctionUrlConfigsCommandInput,
  ...rest: any[]
) => Paginator<ListFunctionUrlConfigsCommandOutput>;
