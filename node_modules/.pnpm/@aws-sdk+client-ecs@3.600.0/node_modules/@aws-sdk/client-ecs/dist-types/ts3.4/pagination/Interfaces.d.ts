import { PaginationConfiguration } from "@smithy/types";
import { ECSClient } from "../ECSClient";
export interface ECSPaginationConfiguration extends PaginationConfiguration {
  client: ECSClient;
}
