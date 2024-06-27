import { Paginator } from "@smithy/types";
import {
  ListServicesByNamespaceCommandInput,
  ListServicesByNamespaceCommandOutput,
} from "../commands/ListServicesByNamespaceCommand";
import { ECSPaginationConfiguration } from "./Interfaces";
export declare const paginateListServicesByNamespace: (
  config: ECSPaginationConfiguration,
  input: ListServicesByNamespaceCommandInput,
  ...rest: any[]
) => Paginator<ListServicesByNamespaceCommandOutput>;
