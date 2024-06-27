import { Paginator } from "@smithy/types";
import {
  ListFunctionsCommandInput,
  ListFunctionsCommandOutput,
} from "../commands/ListFunctionsCommand";
import { LambdaPaginationConfiguration } from "./Interfaces";
export declare const paginateListFunctions: (
  config: LambdaPaginationConfiguration,
  input: ListFunctionsCommandInput,
  ...rest: any[]
) => Paginator<ListFunctionsCommandOutput>;
