import { HttpRequest as __HttpRequest, HttpResponse as __HttpResponse } from "@smithy/protocol-http";
import { EventStreamSerdeContext as __EventStreamSerdeContext, SerdeContext as __SerdeContext } from "@smithy/types";
import { AddLayerVersionPermissionCommandInput, AddLayerVersionPermissionCommandOutput } from "../commands/AddLayerVersionPermissionCommand";
import { AddPermissionCommandInput, AddPermissionCommandOutput } from "../commands/AddPermissionCommand";
import { CreateAliasCommandInput, CreateAliasCommandOutput } from "../commands/CreateAliasCommand";
import { CreateCodeSigningConfigCommandInput, CreateCodeSigningConfigCommandOutput } from "../commands/CreateCodeSigningConfigCommand";
import { CreateEventSourceMappingCommandInput, CreateEventSourceMappingCommandOutput } from "../commands/CreateEventSourceMappingCommand";
import { CreateFunctionCommandInput, CreateFunctionCommandOutput } from "../commands/CreateFunctionCommand";
import { CreateFunctionUrlConfigCommandInput, CreateFunctionUrlConfigCommandOutput } from "../commands/CreateFunctionUrlConfigCommand";
import { DeleteAliasCommandInput, DeleteAliasCommandOutput } from "../commands/DeleteAliasCommand";
import { DeleteCodeSigningConfigCommandInput, DeleteCodeSigningConfigCommandOutput } from "../commands/DeleteCodeSigningConfigCommand";
import { DeleteEventSourceMappingCommandInput, DeleteEventSourceMappingCommandOutput } from "../commands/DeleteEventSourceMappingCommand";
import { DeleteFunctionCodeSigningConfigCommandInput, DeleteFunctionCodeSigningConfigCommandOutput } from "../commands/DeleteFunctionCodeSigningConfigCommand";
import { DeleteFunctionCommandInput, DeleteFunctionCommandOutput } from "../commands/DeleteFunctionCommand";
import { DeleteFunctionConcurrencyCommandInput, DeleteFunctionConcurrencyCommandOutput } from "../commands/DeleteFunctionConcurrencyCommand";
import { DeleteFunctionEventInvokeConfigCommandInput, DeleteFunctionEventInvokeConfigCommandOutput } from "../commands/DeleteFunctionEventInvokeConfigCommand";
import { DeleteFunctionUrlConfigCommandInput, DeleteFunctionUrlConfigCommandOutput } from "../commands/DeleteFunctionUrlConfigCommand";
import { DeleteLayerVersionCommandInput, DeleteLayerVersionCommandOutput } from "../commands/DeleteLayerVersionCommand";
import { DeleteProvisionedConcurrencyConfigCommandInput, DeleteProvisionedConcurrencyConfigCommandOutput } from "../commands/DeleteProvisionedConcurrencyConfigCommand";
import { GetAccountSettingsCommandInput, GetAccountSettingsCommandOutput } from "../commands/GetAccountSettingsCommand";
import { GetAliasCommandInput, GetAliasCommandOutput } from "../commands/GetAliasCommand";
import { GetCodeSigningConfigCommandInput, GetCodeSigningConfigCommandOutput } from "../commands/GetCodeSigningConfigCommand";
import { GetEventSourceMappingCommandInput, GetEventSourceMappingCommandOutput } from "../commands/GetEventSourceMappingCommand";
import { GetFunctionCodeSigningConfigCommandInput, GetFunctionCodeSigningConfigCommandOutput } from "../commands/GetFunctionCodeSigningConfigCommand";
import { GetFunctionCommandInput, GetFunctionCommandOutput } from "../commands/GetFunctionCommand";
import { GetFunctionConcurrencyCommandInput, GetFunctionConcurrencyCommandOutput } from "../commands/GetFunctionConcurrencyCommand";
import { GetFunctionConfigurationCommandInput, GetFunctionConfigurationCommandOutput } from "../commands/GetFunctionConfigurationCommand";
import { GetFunctionEventInvokeConfigCommandInput, GetFunctionEventInvokeConfigCommandOutput } from "../commands/GetFunctionEventInvokeConfigCommand";
import { GetFunctionUrlConfigCommandInput, GetFunctionUrlConfigCommandOutput } from "../commands/GetFunctionUrlConfigCommand";
import { GetLayerVersionByArnCommandInput, GetLayerVersionByArnCommandOutput } from "../commands/GetLayerVersionByArnCommand";
import { GetLayerVersionCommandInput, GetLayerVersionCommandOutput } from "../commands/GetLayerVersionCommand";
import { GetLayerVersionPolicyCommandInput, GetLayerVersionPolicyCommandOutput } from "../commands/GetLayerVersionPolicyCommand";
import { GetPolicyCommandInput, GetPolicyCommandOutput } from "../commands/GetPolicyCommand";
import { GetProvisionedConcurrencyConfigCommandInput, GetProvisionedConcurrencyConfigCommandOutput } from "../commands/GetProvisionedConcurrencyConfigCommand";
import { GetRuntimeManagementConfigCommandInput, GetRuntimeManagementConfigCommandOutput } from "../commands/GetRuntimeManagementConfigCommand";
import { InvokeAsyncCommandInput, InvokeAsyncCommandOutput } from "../commands/InvokeAsyncCommand";
import { InvokeCommandInput, InvokeCommandOutput } from "../commands/InvokeCommand";
import { InvokeWithResponseStreamCommandInput, InvokeWithResponseStreamCommandOutput } from "../commands/InvokeWithResponseStreamCommand";
import { ListAliasesCommandInput, ListAliasesCommandOutput } from "../commands/ListAliasesCommand";
import { ListCodeSigningConfigsCommandInput, ListCodeSigningConfigsCommandOutput } from "../commands/ListCodeSigningConfigsCommand";
import { ListEventSourceMappingsCommandInput, ListEventSourceMappingsCommandOutput } from "../commands/ListEventSourceMappingsCommand";
import { ListFunctionEventInvokeConfigsCommandInput, ListFunctionEventInvokeConfigsCommandOutput } from "../commands/ListFunctionEventInvokeConfigsCommand";
import { ListFunctionsByCodeSigningConfigCommandInput, ListFunctionsByCodeSigningConfigCommandOutput } from "../commands/ListFunctionsByCodeSigningConfigCommand";
import { ListFunctionsCommandInput, ListFunctionsCommandOutput } from "../commands/ListFunctionsCommand";
import { ListFunctionUrlConfigsCommandInput, ListFunctionUrlConfigsCommandOutput } from "../commands/ListFunctionUrlConfigsCommand";
import { ListLayersCommandInput, ListLayersCommandOutput } from "../commands/ListLayersCommand";
import { ListLayerVersionsCommandInput, ListLayerVersionsCommandOutput } from "../commands/ListLayerVersionsCommand";
import { ListProvisionedConcurrencyConfigsCommandInput, ListProvisionedConcurrencyConfigsCommandOutput } from "../commands/ListProvisionedConcurrencyConfigsCommand";
import { ListTagsCommandInput, ListTagsCommandOutput } from "../commands/ListTagsCommand";
import { ListVersionsByFunctionCommandInput, ListVersionsByFunctionCommandOutput } from "../commands/ListVersionsByFunctionCommand";
import { PublishLayerVersionCommandInput, PublishLayerVersionCommandOutput } from "../commands/PublishLayerVersionCommand";
import { PublishVersionCommandInput, PublishVersionCommandOutput } from "../commands/PublishVersionCommand";
import { PutFunctionCodeSigningConfigCommandInput, PutFunctionCodeSigningConfigCommandOutput } from "../commands/PutFunctionCodeSigningConfigCommand";
import { PutFunctionConcurrencyCommandInput, PutFunctionConcurrencyCommandOutput } from "../commands/PutFunctionConcurrencyCommand";
import { PutFunctionEventInvokeConfigCommandInput, PutFunctionEventInvokeConfigCommandOutput } from "../commands/PutFunctionEventInvokeConfigCommand";
import { PutProvisionedConcurrencyConfigCommandInput, PutProvisionedConcurrencyConfigCommandOutput } from "../commands/PutProvisionedConcurrencyConfigCommand";
import { PutRuntimeManagementConfigCommandInput, PutRuntimeManagementConfigCommandOutput } from "../commands/PutRuntimeManagementConfigCommand";
import { RemoveLayerVersionPermissionCommandInput, RemoveLayerVersionPermissionCommandOutput } from "../commands/RemoveLayerVersionPermissionCommand";
import { RemovePermissionCommandInput, RemovePermissionCommandOutput } from "../commands/RemovePermissionCommand";
import { TagResourceCommandInput, TagResourceCommandOutput } from "../commands/TagResourceCommand";
import { UntagResourceCommandInput, UntagResourceCommandOutput } from "../commands/UntagResourceCommand";
import { UpdateAliasCommandInput, UpdateAliasCommandOutput } from "../commands/UpdateAliasCommand";
import { UpdateCodeSigningConfigCommandInput, UpdateCodeSigningConfigCommandOutput } from "../commands/UpdateCodeSigningConfigCommand";
import { UpdateEventSourceMappingCommandInput, UpdateEventSourceMappingCommandOutput } from "../commands/UpdateEventSourceMappingCommand";
import { UpdateFunctionCodeCommandInput, UpdateFunctionCodeCommandOutput } from "../commands/UpdateFunctionCodeCommand";
import { UpdateFunctionConfigurationCommandInput, UpdateFunctionConfigurationCommandOutput } from "../commands/UpdateFunctionConfigurationCommand";
import { UpdateFunctionEventInvokeConfigCommandInput, UpdateFunctionEventInvokeConfigCommandOutput } from "../commands/UpdateFunctionEventInvokeConfigCommand";
import { UpdateFunctionUrlConfigCommandInput, UpdateFunctionUrlConfigCommandOutput } from "../commands/UpdateFunctionUrlConfigCommand";
/**
 * serializeAws_restJson1AddLayerVersionPermissionCommand
 */
export declare const se_AddLayerVersionPermissionCommand: (input: AddLayerVersionPermissionCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1AddPermissionCommand
 */
export declare const se_AddPermissionCommand: (input: AddPermissionCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1CreateAliasCommand
 */
export declare const se_CreateAliasCommand: (input: CreateAliasCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1CreateCodeSigningConfigCommand
 */
export declare const se_CreateCodeSigningConfigCommand: (input: CreateCodeSigningConfigCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1CreateEventSourceMappingCommand
 */
export declare const se_CreateEventSourceMappingCommand: (input: CreateEventSourceMappingCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1CreateFunctionCommand
 */
export declare const se_CreateFunctionCommand: (input: CreateFunctionCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1CreateFunctionUrlConfigCommand
 */
export declare const se_CreateFunctionUrlConfigCommand: (input: CreateFunctionUrlConfigCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1DeleteAliasCommand
 */
export declare const se_DeleteAliasCommand: (input: DeleteAliasCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1DeleteCodeSigningConfigCommand
 */
export declare const se_DeleteCodeSigningConfigCommand: (input: DeleteCodeSigningConfigCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1DeleteEventSourceMappingCommand
 */
export declare const se_DeleteEventSourceMappingCommand: (input: DeleteEventSourceMappingCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1DeleteFunctionCommand
 */
export declare const se_DeleteFunctionCommand: (input: DeleteFunctionCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1DeleteFunctionCodeSigningConfigCommand
 */
export declare const se_DeleteFunctionCodeSigningConfigCommand: (input: DeleteFunctionCodeSigningConfigCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1DeleteFunctionConcurrencyCommand
 */
export declare const se_DeleteFunctionConcurrencyCommand: (input: DeleteFunctionConcurrencyCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1DeleteFunctionEventInvokeConfigCommand
 */
export declare const se_DeleteFunctionEventInvokeConfigCommand: (input: DeleteFunctionEventInvokeConfigCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1DeleteFunctionUrlConfigCommand
 */
export declare const se_DeleteFunctionUrlConfigCommand: (input: DeleteFunctionUrlConfigCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1DeleteLayerVersionCommand
 */
export declare const se_DeleteLayerVersionCommand: (input: DeleteLayerVersionCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1DeleteProvisionedConcurrencyConfigCommand
 */
export declare const se_DeleteProvisionedConcurrencyConfigCommand: (input: DeleteProvisionedConcurrencyConfigCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1GetAccountSettingsCommand
 */
export declare const se_GetAccountSettingsCommand: (input: GetAccountSettingsCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1GetAliasCommand
 */
export declare const se_GetAliasCommand: (input: GetAliasCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1GetCodeSigningConfigCommand
 */
export declare const se_GetCodeSigningConfigCommand: (input: GetCodeSigningConfigCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1GetEventSourceMappingCommand
 */
export declare const se_GetEventSourceMappingCommand: (input: GetEventSourceMappingCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1GetFunctionCommand
 */
export declare const se_GetFunctionCommand: (input: GetFunctionCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1GetFunctionCodeSigningConfigCommand
 */
export declare const se_GetFunctionCodeSigningConfigCommand: (input: GetFunctionCodeSigningConfigCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1GetFunctionConcurrencyCommand
 */
export declare const se_GetFunctionConcurrencyCommand: (input: GetFunctionConcurrencyCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1GetFunctionConfigurationCommand
 */
export declare const se_GetFunctionConfigurationCommand: (input: GetFunctionConfigurationCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1GetFunctionEventInvokeConfigCommand
 */
export declare const se_GetFunctionEventInvokeConfigCommand: (input: GetFunctionEventInvokeConfigCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1GetFunctionUrlConfigCommand
 */
export declare const se_GetFunctionUrlConfigCommand: (input: GetFunctionUrlConfigCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1GetLayerVersionCommand
 */
export declare const se_GetLayerVersionCommand: (input: GetLayerVersionCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1GetLayerVersionByArnCommand
 */
export declare const se_GetLayerVersionByArnCommand: (input: GetLayerVersionByArnCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1GetLayerVersionPolicyCommand
 */
export declare const se_GetLayerVersionPolicyCommand: (input: GetLayerVersionPolicyCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1GetPolicyCommand
 */
export declare const se_GetPolicyCommand: (input: GetPolicyCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1GetProvisionedConcurrencyConfigCommand
 */
export declare const se_GetProvisionedConcurrencyConfigCommand: (input: GetProvisionedConcurrencyConfigCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1GetRuntimeManagementConfigCommand
 */
export declare const se_GetRuntimeManagementConfigCommand: (input: GetRuntimeManagementConfigCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1InvokeCommand
 */
export declare const se_InvokeCommand: (input: InvokeCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1InvokeAsyncCommand
 */
export declare const se_InvokeAsyncCommand: (input: InvokeAsyncCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1InvokeWithResponseStreamCommand
 */
export declare const se_InvokeWithResponseStreamCommand: (input: InvokeWithResponseStreamCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1ListAliasesCommand
 */
export declare const se_ListAliasesCommand: (input: ListAliasesCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1ListCodeSigningConfigsCommand
 */
export declare const se_ListCodeSigningConfigsCommand: (input: ListCodeSigningConfigsCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1ListEventSourceMappingsCommand
 */
export declare const se_ListEventSourceMappingsCommand: (input: ListEventSourceMappingsCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1ListFunctionEventInvokeConfigsCommand
 */
export declare const se_ListFunctionEventInvokeConfigsCommand: (input: ListFunctionEventInvokeConfigsCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1ListFunctionsCommand
 */
export declare const se_ListFunctionsCommand: (input: ListFunctionsCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1ListFunctionsByCodeSigningConfigCommand
 */
export declare const se_ListFunctionsByCodeSigningConfigCommand: (input: ListFunctionsByCodeSigningConfigCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1ListFunctionUrlConfigsCommand
 */
export declare const se_ListFunctionUrlConfigsCommand: (input: ListFunctionUrlConfigsCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1ListLayersCommand
 */
export declare const se_ListLayersCommand: (input: ListLayersCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1ListLayerVersionsCommand
 */
export declare const se_ListLayerVersionsCommand: (input: ListLayerVersionsCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1ListProvisionedConcurrencyConfigsCommand
 */
export declare const se_ListProvisionedConcurrencyConfigsCommand: (input: ListProvisionedConcurrencyConfigsCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1ListTagsCommand
 */
export declare const se_ListTagsCommand: (input: ListTagsCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1ListVersionsByFunctionCommand
 */
export declare const se_ListVersionsByFunctionCommand: (input: ListVersionsByFunctionCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1PublishLayerVersionCommand
 */
export declare const se_PublishLayerVersionCommand: (input: PublishLayerVersionCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1PublishVersionCommand
 */
export declare const se_PublishVersionCommand: (input: PublishVersionCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1PutFunctionCodeSigningConfigCommand
 */
export declare const se_PutFunctionCodeSigningConfigCommand: (input: PutFunctionCodeSigningConfigCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1PutFunctionConcurrencyCommand
 */
export declare const se_PutFunctionConcurrencyCommand: (input: PutFunctionConcurrencyCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1PutFunctionEventInvokeConfigCommand
 */
export declare const se_PutFunctionEventInvokeConfigCommand: (input: PutFunctionEventInvokeConfigCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1PutProvisionedConcurrencyConfigCommand
 */
export declare const se_PutProvisionedConcurrencyConfigCommand: (input: PutProvisionedConcurrencyConfigCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1PutRuntimeManagementConfigCommand
 */
export declare const se_PutRuntimeManagementConfigCommand: (input: PutRuntimeManagementConfigCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1RemoveLayerVersionPermissionCommand
 */
export declare const se_RemoveLayerVersionPermissionCommand: (input: RemoveLayerVersionPermissionCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1RemovePermissionCommand
 */
export declare const se_RemovePermissionCommand: (input: RemovePermissionCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1TagResourceCommand
 */
export declare const se_TagResourceCommand: (input: TagResourceCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1UntagResourceCommand
 */
export declare const se_UntagResourceCommand: (input: UntagResourceCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1UpdateAliasCommand
 */
export declare const se_UpdateAliasCommand: (input: UpdateAliasCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1UpdateCodeSigningConfigCommand
 */
export declare const se_UpdateCodeSigningConfigCommand: (input: UpdateCodeSigningConfigCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1UpdateEventSourceMappingCommand
 */
export declare const se_UpdateEventSourceMappingCommand: (input: UpdateEventSourceMappingCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1UpdateFunctionCodeCommand
 */
export declare const se_UpdateFunctionCodeCommand: (input: UpdateFunctionCodeCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1UpdateFunctionConfigurationCommand
 */
export declare const se_UpdateFunctionConfigurationCommand: (input: UpdateFunctionConfigurationCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1UpdateFunctionEventInvokeConfigCommand
 */
export declare const se_UpdateFunctionEventInvokeConfigCommand: (input: UpdateFunctionEventInvokeConfigCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * serializeAws_restJson1UpdateFunctionUrlConfigCommand
 */
export declare const se_UpdateFunctionUrlConfigCommand: (input: UpdateFunctionUrlConfigCommandInput, context: __SerdeContext) => Promise<__HttpRequest>;
/**
 * deserializeAws_restJson1AddLayerVersionPermissionCommand
 */
export declare const de_AddLayerVersionPermissionCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<AddLayerVersionPermissionCommandOutput>;
/**
 * deserializeAws_restJson1AddPermissionCommand
 */
export declare const de_AddPermissionCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<AddPermissionCommandOutput>;
/**
 * deserializeAws_restJson1CreateAliasCommand
 */
export declare const de_CreateAliasCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<CreateAliasCommandOutput>;
/**
 * deserializeAws_restJson1CreateCodeSigningConfigCommand
 */
export declare const de_CreateCodeSigningConfigCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<CreateCodeSigningConfigCommandOutput>;
/**
 * deserializeAws_restJson1CreateEventSourceMappingCommand
 */
export declare const de_CreateEventSourceMappingCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<CreateEventSourceMappingCommandOutput>;
/**
 * deserializeAws_restJson1CreateFunctionCommand
 */
export declare const de_CreateFunctionCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<CreateFunctionCommandOutput>;
/**
 * deserializeAws_restJson1CreateFunctionUrlConfigCommand
 */
export declare const de_CreateFunctionUrlConfigCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<CreateFunctionUrlConfigCommandOutput>;
/**
 * deserializeAws_restJson1DeleteAliasCommand
 */
export declare const de_DeleteAliasCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<DeleteAliasCommandOutput>;
/**
 * deserializeAws_restJson1DeleteCodeSigningConfigCommand
 */
export declare const de_DeleteCodeSigningConfigCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<DeleteCodeSigningConfigCommandOutput>;
/**
 * deserializeAws_restJson1DeleteEventSourceMappingCommand
 */
export declare const de_DeleteEventSourceMappingCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<DeleteEventSourceMappingCommandOutput>;
/**
 * deserializeAws_restJson1DeleteFunctionCommand
 */
export declare const de_DeleteFunctionCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<DeleteFunctionCommandOutput>;
/**
 * deserializeAws_restJson1DeleteFunctionCodeSigningConfigCommand
 */
export declare const de_DeleteFunctionCodeSigningConfigCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<DeleteFunctionCodeSigningConfigCommandOutput>;
/**
 * deserializeAws_restJson1DeleteFunctionConcurrencyCommand
 */
export declare const de_DeleteFunctionConcurrencyCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<DeleteFunctionConcurrencyCommandOutput>;
/**
 * deserializeAws_restJson1DeleteFunctionEventInvokeConfigCommand
 */
export declare const de_DeleteFunctionEventInvokeConfigCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<DeleteFunctionEventInvokeConfigCommandOutput>;
/**
 * deserializeAws_restJson1DeleteFunctionUrlConfigCommand
 */
export declare const de_DeleteFunctionUrlConfigCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<DeleteFunctionUrlConfigCommandOutput>;
/**
 * deserializeAws_restJson1DeleteLayerVersionCommand
 */
export declare const de_DeleteLayerVersionCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<DeleteLayerVersionCommandOutput>;
/**
 * deserializeAws_restJson1DeleteProvisionedConcurrencyConfigCommand
 */
export declare const de_DeleteProvisionedConcurrencyConfigCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<DeleteProvisionedConcurrencyConfigCommandOutput>;
/**
 * deserializeAws_restJson1GetAccountSettingsCommand
 */
export declare const de_GetAccountSettingsCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<GetAccountSettingsCommandOutput>;
/**
 * deserializeAws_restJson1GetAliasCommand
 */
export declare const de_GetAliasCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<GetAliasCommandOutput>;
/**
 * deserializeAws_restJson1GetCodeSigningConfigCommand
 */
export declare const de_GetCodeSigningConfigCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<GetCodeSigningConfigCommandOutput>;
/**
 * deserializeAws_restJson1GetEventSourceMappingCommand
 */
export declare const de_GetEventSourceMappingCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<GetEventSourceMappingCommandOutput>;
/**
 * deserializeAws_restJson1GetFunctionCommand
 */
export declare const de_GetFunctionCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<GetFunctionCommandOutput>;
/**
 * deserializeAws_restJson1GetFunctionCodeSigningConfigCommand
 */
export declare const de_GetFunctionCodeSigningConfigCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<GetFunctionCodeSigningConfigCommandOutput>;
/**
 * deserializeAws_restJson1GetFunctionConcurrencyCommand
 */
export declare const de_GetFunctionConcurrencyCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<GetFunctionConcurrencyCommandOutput>;
/**
 * deserializeAws_restJson1GetFunctionConfigurationCommand
 */
export declare const de_GetFunctionConfigurationCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<GetFunctionConfigurationCommandOutput>;
/**
 * deserializeAws_restJson1GetFunctionEventInvokeConfigCommand
 */
export declare const de_GetFunctionEventInvokeConfigCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<GetFunctionEventInvokeConfigCommandOutput>;
/**
 * deserializeAws_restJson1GetFunctionUrlConfigCommand
 */
export declare const de_GetFunctionUrlConfigCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<GetFunctionUrlConfigCommandOutput>;
/**
 * deserializeAws_restJson1GetLayerVersionCommand
 */
export declare const de_GetLayerVersionCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<GetLayerVersionCommandOutput>;
/**
 * deserializeAws_restJson1GetLayerVersionByArnCommand
 */
export declare const de_GetLayerVersionByArnCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<GetLayerVersionByArnCommandOutput>;
/**
 * deserializeAws_restJson1GetLayerVersionPolicyCommand
 */
export declare const de_GetLayerVersionPolicyCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<GetLayerVersionPolicyCommandOutput>;
/**
 * deserializeAws_restJson1GetPolicyCommand
 */
export declare const de_GetPolicyCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<GetPolicyCommandOutput>;
/**
 * deserializeAws_restJson1GetProvisionedConcurrencyConfigCommand
 */
export declare const de_GetProvisionedConcurrencyConfigCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<GetProvisionedConcurrencyConfigCommandOutput>;
/**
 * deserializeAws_restJson1GetRuntimeManagementConfigCommand
 */
export declare const de_GetRuntimeManagementConfigCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<GetRuntimeManagementConfigCommandOutput>;
/**
 * deserializeAws_restJson1InvokeCommand
 */
export declare const de_InvokeCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<InvokeCommandOutput>;
/**
 * deserializeAws_restJson1InvokeAsyncCommand
 */
export declare const de_InvokeAsyncCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<InvokeAsyncCommandOutput>;
/**
 * deserializeAws_restJson1InvokeWithResponseStreamCommand
 */
export declare const de_InvokeWithResponseStreamCommand: (output: __HttpResponse, context: __SerdeContext & __EventStreamSerdeContext) => Promise<InvokeWithResponseStreamCommandOutput>;
/**
 * deserializeAws_restJson1ListAliasesCommand
 */
export declare const de_ListAliasesCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<ListAliasesCommandOutput>;
/**
 * deserializeAws_restJson1ListCodeSigningConfigsCommand
 */
export declare const de_ListCodeSigningConfigsCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<ListCodeSigningConfigsCommandOutput>;
/**
 * deserializeAws_restJson1ListEventSourceMappingsCommand
 */
export declare const de_ListEventSourceMappingsCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<ListEventSourceMappingsCommandOutput>;
/**
 * deserializeAws_restJson1ListFunctionEventInvokeConfigsCommand
 */
export declare const de_ListFunctionEventInvokeConfigsCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<ListFunctionEventInvokeConfigsCommandOutput>;
/**
 * deserializeAws_restJson1ListFunctionsCommand
 */
export declare const de_ListFunctionsCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<ListFunctionsCommandOutput>;
/**
 * deserializeAws_restJson1ListFunctionsByCodeSigningConfigCommand
 */
export declare const de_ListFunctionsByCodeSigningConfigCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<ListFunctionsByCodeSigningConfigCommandOutput>;
/**
 * deserializeAws_restJson1ListFunctionUrlConfigsCommand
 */
export declare const de_ListFunctionUrlConfigsCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<ListFunctionUrlConfigsCommandOutput>;
/**
 * deserializeAws_restJson1ListLayersCommand
 */
export declare const de_ListLayersCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<ListLayersCommandOutput>;
/**
 * deserializeAws_restJson1ListLayerVersionsCommand
 */
export declare const de_ListLayerVersionsCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<ListLayerVersionsCommandOutput>;
/**
 * deserializeAws_restJson1ListProvisionedConcurrencyConfigsCommand
 */
export declare const de_ListProvisionedConcurrencyConfigsCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<ListProvisionedConcurrencyConfigsCommandOutput>;
/**
 * deserializeAws_restJson1ListTagsCommand
 */
export declare const de_ListTagsCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<ListTagsCommandOutput>;
/**
 * deserializeAws_restJson1ListVersionsByFunctionCommand
 */
export declare const de_ListVersionsByFunctionCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<ListVersionsByFunctionCommandOutput>;
/**
 * deserializeAws_restJson1PublishLayerVersionCommand
 */
export declare const de_PublishLayerVersionCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<PublishLayerVersionCommandOutput>;
/**
 * deserializeAws_restJson1PublishVersionCommand
 */
export declare const de_PublishVersionCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<PublishVersionCommandOutput>;
/**
 * deserializeAws_restJson1PutFunctionCodeSigningConfigCommand
 */
export declare const de_PutFunctionCodeSigningConfigCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<PutFunctionCodeSigningConfigCommandOutput>;
/**
 * deserializeAws_restJson1PutFunctionConcurrencyCommand
 */
export declare const de_PutFunctionConcurrencyCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<PutFunctionConcurrencyCommandOutput>;
/**
 * deserializeAws_restJson1PutFunctionEventInvokeConfigCommand
 */
export declare const de_PutFunctionEventInvokeConfigCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<PutFunctionEventInvokeConfigCommandOutput>;
/**
 * deserializeAws_restJson1PutProvisionedConcurrencyConfigCommand
 */
export declare const de_PutProvisionedConcurrencyConfigCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<PutProvisionedConcurrencyConfigCommandOutput>;
/**
 * deserializeAws_restJson1PutRuntimeManagementConfigCommand
 */
export declare const de_PutRuntimeManagementConfigCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<PutRuntimeManagementConfigCommandOutput>;
/**
 * deserializeAws_restJson1RemoveLayerVersionPermissionCommand
 */
export declare const de_RemoveLayerVersionPermissionCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<RemoveLayerVersionPermissionCommandOutput>;
/**
 * deserializeAws_restJson1RemovePermissionCommand
 */
export declare const de_RemovePermissionCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<RemovePermissionCommandOutput>;
/**
 * deserializeAws_restJson1TagResourceCommand
 */
export declare const de_TagResourceCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<TagResourceCommandOutput>;
/**
 * deserializeAws_restJson1UntagResourceCommand
 */
export declare const de_UntagResourceCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<UntagResourceCommandOutput>;
/**
 * deserializeAws_restJson1UpdateAliasCommand
 */
export declare const de_UpdateAliasCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<UpdateAliasCommandOutput>;
/**
 * deserializeAws_restJson1UpdateCodeSigningConfigCommand
 */
export declare const de_UpdateCodeSigningConfigCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<UpdateCodeSigningConfigCommandOutput>;
/**
 * deserializeAws_restJson1UpdateEventSourceMappingCommand
 */
export declare const de_UpdateEventSourceMappingCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<UpdateEventSourceMappingCommandOutput>;
/**
 * deserializeAws_restJson1UpdateFunctionCodeCommand
 */
export declare const de_UpdateFunctionCodeCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<UpdateFunctionCodeCommandOutput>;
/**
 * deserializeAws_restJson1UpdateFunctionConfigurationCommand
 */
export declare const de_UpdateFunctionConfigurationCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<UpdateFunctionConfigurationCommandOutput>;
/**
 * deserializeAws_restJson1UpdateFunctionEventInvokeConfigCommand
 */
export declare const de_UpdateFunctionEventInvokeConfigCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<UpdateFunctionEventInvokeConfigCommandOutput>;
/**
 * deserializeAws_restJson1UpdateFunctionUrlConfigCommand
 */
export declare const de_UpdateFunctionUrlConfigCommand: (output: __HttpResponse, context: __SerdeContext) => Promise<UpdateFunctionUrlConfigCommandOutput>;
