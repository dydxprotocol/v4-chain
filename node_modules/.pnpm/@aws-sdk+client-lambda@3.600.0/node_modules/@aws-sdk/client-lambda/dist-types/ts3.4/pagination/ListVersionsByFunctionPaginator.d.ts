import { Paginator } from "@smithy/types";
import {
  ListVersionsByFunctionCommandInput,
  ListVersionsByFunctionCommandOutput,
} from "../commands/ListVersionsByFunctionCommand";
import { LambdaPaginationConfiguration } from "./Interfaces";
export declare const paginateListVersionsByFunction: (
  config: LambdaPaginationConfiguration,
  input: ListVersionsByFunctionCommandInput,
  ...rest: any[]
) => Paginator<ListVersionsByFunctionCommandOutput>;
