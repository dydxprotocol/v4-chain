import { createAggregatedClient } from "@smithy/smithy-client";
import { AddLayerVersionPermissionCommand, } from "./commands/AddLayerVersionPermissionCommand";
import { AddPermissionCommand, } from "./commands/AddPermissionCommand";
import { CreateAliasCommand } from "./commands/CreateAliasCommand";
import { CreateCodeSigningConfigCommand, } from "./commands/CreateCodeSigningConfigCommand";
import { CreateEventSourceMappingCommand, } from "./commands/CreateEventSourceMappingCommand";
import { CreateFunctionCommand, } from "./commands/CreateFunctionCommand";
import { CreateFunctionUrlConfigCommand, } from "./commands/CreateFunctionUrlConfigCommand";
import { DeleteAliasCommand } from "./commands/DeleteAliasCommand";
import { DeleteCodeSigningConfigCommand, } from "./commands/DeleteCodeSigningConfigCommand";
import { DeleteEventSourceMappingCommand, } from "./commands/DeleteEventSourceMappingCommand";
import { DeleteFunctionCodeSigningConfigCommand, } from "./commands/DeleteFunctionCodeSigningConfigCommand";
import { DeleteFunctionCommand, } from "./commands/DeleteFunctionCommand";
import { DeleteFunctionConcurrencyCommand, } from "./commands/DeleteFunctionConcurrencyCommand";
import { DeleteFunctionEventInvokeConfigCommand, } from "./commands/DeleteFunctionEventInvokeConfigCommand";
import { DeleteFunctionUrlConfigCommand, } from "./commands/DeleteFunctionUrlConfigCommand";
import { DeleteLayerVersionCommand, } from "./commands/DeleteLayerVersionCommand";
import { DeleteProvisionedConcurrencyConfigCommand, } from "./commands/DeleteProvisionedConcurrencyConfigCommand";
import { GetAccountSettingsCommand, } from "./commands/GetAccountSettingsCommand";
import { GetAliasCommand } from "./commands/GetAliasCommand";
import { GetCodeSigningConfigCommand, } from "./commands/GetCodeSigningConfigCommand";
import { GetEventSourceMappingCommand, } from "./commands/GetEventSourceMappingCommand";
import { GetFunctionCodeSigningConfigCommand, } from "./commands/GetFunctionCodeSigningConfigCommand";
import { GetFunctionCommand } from "./commands/GetFunctionCommand";
import { GetFunctionConcurrencyCommand, } from "./commands/GetFunctionConcurrencyCommand";
import { GetFunctionConfigurationCommand, } from "./commands/GetFunctionConfigurationCommand";
import { GetFunctionEventInvokeConfigCommand, } from "./commands/GetFunctionEventInvokeConfigCommand";
import { GetFunctionUrlConfigCommand, } from "./commands/GetFunctionUrlConfigCommand";
import { GetLayerVersionByArnCommand, } from "./commands/GetLayerVersionByArnCommand";
import { GetLayerVersionCommand, } from "./commands/GetLayerVersionCommand";
import { GetLayerVersionPolicyCommand, } from "./commands/GetLayerVersionPolicyCommand";
import { GetPolicyCommand } from "./commands/GetPolicyCommand";
import { GetProvisionedConcurrencyConfigCommand, } from "./commands/GetProvisionedConcurrencyConfigCommand";
import { GetRuntimeManagementConfigCommand, } from "./commands/GetRuntimeManagementConfigCommand";
import { InvokeAsyncCommand } from "./commands/InvokeAsyncCommand";
import { InvokeCommand } from "./commands/InvokeCommand";
import { InvokeWithResponseStreamCommand, } from "./commands/InvokeWithResponseStreamCommand";
import { ListAliasesCommand } from "./commands/ListAliasesCommand";
import { ListCodeSigningConfigsCommand, } from "./commands/ListCodeSigningConfigsCommand";
import { ListEventSourceMappingsCommand, } from "./commands/ListEventSourceMappingsCommand";
import { ListFunctionEventInvokeConfigsCommand, } from "./commands/ListFunctionEventInvokeConfigsCommand";
import { ListFunctionsByCodeSigningConfigCommand, } from "./commands/ListFunctionsByCodeSigningConfigCommand";
import { ListFunctionsCommand, } from "./commands/ListFunctionsCommand";
import { ListFunctionUrlConfigsCommand, } from "./commands/ListFunctionUrlConfigsCommand";
import { ListLayersCommand } from "./commands/ListLayersCommand";
import { ListLayerVersionsCommand, } from "./commands/ListLayerVersionsCommand";
import { ListProvisionedConcurrencyConfigsCommand, } from "./commands/ListProvisionedConcurrencyConfigsCommand";
import { ListTagsCommand } from "./commands/ListTagsCommand";
import { ListVersionsByFunctionCommand, } from "./commands/ListVersionsByFunctionCommand";
import { PublishLayerVersionCommand, } from "./commands/PublishLayerVersionCommand";
import { PublishVersionCommand, } from "./commands/PublishVersionCommand";
import { PutFunctionCodeSigningConfigCommand, } from "./commands/PutFunctionCodeSigningConfigCommand";
import { PutFunctionConcurrencyCommand, } from "./commands/PutFunctionConcurrencyCommand";
import { PutFunctionEventInvokeConfigCommand, } from "./commands/PutFunctionEventInvokeConfigCommand";
import { PutProvisionedConcurrencyConfigCommand, } from "./commands/PutProvisionedConcurrencyConfigCommand";
import { PutRuntimeManagementConfigCommand, } from "./commands/PutRuntimeManagementConfigCommand";
import { RemoveLayerVersionPermissionCommand, } from "./commands/RemoveLayerVersionPermissionCommand";
import { RemovePermissionCommand, } from "./commands/RemovePermissionCommand";
import { TagResourceCommand } from "./commands/TagResourceCommand";
import { UntagResourceCommand, } from "./commands/UntagResourceCommand";
import { UpdateAliasCommand } from "./commands/UpdateAliasCommand";
import { UpdateCodeSigningConfigCommand, } from "./commands/UpdateCodeSigningConfigCommand";
import { UpdateEventSourceMappingCommand, } from "./commands/UpdateEventSourceMappingCommand";
import { UpdateFunctionCodeCommand, } from "./commands/UpdateFunctionCodeCommand";
import { UpdateFunctionConfigurationCommand, } from "./commands/UpdateFunctionConfigurationCommand";
import { UpdateFunctionEventInvokeConfigCommand, } from "./commands/UpdateFunctionEventInvokeConfigCommand";
import { UpdateFunctionUrlConfigCommand, } from "./commands/UpdateFunctionUrlConfigCommand";
import { LambdaClient } from "./LambdaClient";
const commands = {
    AddLayerVersionPermissionCommand,
    AddPermissionCommand,
    CreateAliasCommand,
    CreateCodeSigningConfigCommand,
    CreateEventSourceMappingCommand,
    CreateFunctionCommand,
    CreateFunctionUrlConfigCommand,
    DeleteAliasCommand,
    DeleteCodeSigningConfigCommand,
    DeleteEventSourceMappingCommand,
    DeleteFunctionCommand,
    DeleteFunctionCodeSigningConfigCommand,
    DeleteFunctionConcurrencyCommand,
    DeleteFunctionEventInvokeConfigCommand,
    DeleteFunctionUrlConfigCommand,
    DeleteLayerVersionCommand,
    DeleteProvisionedConcurrencyConfigCommand,
    GetAccountSettingsCommand,
    GetAliasCommand,
    GetCodeSigningConfigCommand,
    GetEventSourceMappingCommand,
    GetFunctionCommand,
    GetFunctionCodeSigningConfigCommand,
    GetFunctionConcurrencyCommand,
    GetFunctionConfigurationCommand,
    GetFunctionEventInvokeConfigCommand,
    GetFunctionUrlConfigCommand,
    GetLayerVersionCommand,
    GetLayerVersionByArnCommand,
    GetLayerVersionPolicyCommand,
    GetPolicyCommand,
    GetProvisionedConcurrencyConfigCommand,
    GetRuntimeManagementConfigCommand,
    InvokeCommand,
    InvokeAsyncCommand,
    InvokeWithResponseStreamCommand,
    ListAliasesCommand,
    ListCodeSigningConfigsCommand,
    ListEventSourceMappingsCommand,
    ListFunctionEventInvokeConfigsCommand,
    ListFunctionsCommand,
    ListFunctionsByCodeSigningConfigCommand,
    ListFunctionUrlConfigsCommand,
    ListLayersCommand,
    ListLayerVersionsCommand,
    ListProvisionedConcurrencyConfigsCommand,
    ListTagsCommand,
    ListVersionsByFunctionCommand,
    PublishLayerVersionCommand,
    PublishVersionCommand,
    PutFunctionCodeSigningConfigCommand,
    PutFunctionConcurrencyCommand,
    PutFunctionEventInvokeConfigCommand,
    PutProvisionedConcurrencyConfigCommand,
    PutRuntimeManagementConfigCommand,
    RemoveLayerVersionPermissionCommand,
    RemovePermissionCommand,
    TagResourceCommand,
    UntagResourceCommand,
    UpdateAliasCommand,
    UpdateCodeSigningConfigCommand,
    UpdateEventSourceMappingCommand,
    UpdateFunctionCodeCommand,
    UpdateFunctionConfigurationCommand,
    UpdateFunctionEventInvokeConfigCommand,
    UpdateFunctionUrlConfigCommand,
};
export class Lambda extends LambdaClient {
}
createAggregatedClient(commands, Lambda);
