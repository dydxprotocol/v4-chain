import { Paginator } from "@smithy/types";
import {
  ListFunctionEventInvokeConfigsCommandInput,
  ListFunctionEventInvokeConfigsCommandOutput,
} from "../commands/ListFunctionEventInvokeConfigsCommand";
import { LambdaPaginationConfiguration } from "./Interfaces";
export declare const paginateListFunctionEventInvokeConfigs: (
  config: LambdaPaginationConfiguration,
  input: ListFunctionEventInvokeConfigsCommandInput,
  ...rest: any[]
) => Paginator<ListFunctionEventInvokeConfigsCommandOutput>;
