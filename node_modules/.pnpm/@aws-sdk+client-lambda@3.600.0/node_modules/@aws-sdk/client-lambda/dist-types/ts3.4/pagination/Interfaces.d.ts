import { PaginationConfiguration } from "@smithy/types";
import { LambdaClient } from "../LambdaClient";
export interface LambdaPaginationConfiguration extends PaginationConfiguration {
  client: LambdaClient;
}
