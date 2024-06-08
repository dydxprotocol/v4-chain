import { PageRequest, PageRequestSDKType, PageResponse, PageResponseSDKType } from "../../query/v1beta1/pagination";
import { Any, AnySDKType } from "../../../../google/protobuf/any";
import { BlockID, BlockIDSDKType } from "../../../../tendermint/types/types";
import { Block as Block1 } from "../../../../tendermint/types/block";
import { BlockSDKType as Block1SDKType } from "../../../../tendermint/types/block";
import { Block as Block2 } from "./types";
import { BlockSDKType as Block2SDKType } from "./types";
import { DefaultNodeInfo, DefaultNodeInfoSDKType } from "../../../../tendermint/p2p/types";
import { Long, DeepPartial } from "../../../../helpers";
import * as _m0 from "protobufjs/minimal";
/** GetValidatorSetByHeightRequest is the request type for the Query/GetValidatorSetByHeight RPC method. */

export interface GetValidatorSetByHeightRequest {
  height: Long;
  /** pagination defines an pagination for the request. */

  pagination?: PageRequest;
}
/** GetValidatorSetByHeightRequest is the request type for the Query/GetValidatorSetByHeight RPC method. */

export interface GetValidatorSetByHeightRequestSDKType {
  height: Long;
  pagination?: PageRequestSDKType;
}
/** GetValidatorSetByHeightResponse is the response type for the Query/GetValidatorSetByHeight RPC method. */

export interface GetValidatorSetByHeightResponse {
  blockHeight: Long;
  validators: Validator[];
  /** pagination defines an pagination for the response. */

  pagination?: PageResponse;
}
/** GetValidatorSetByHeightResponse is the response type for the Query/GetValidatorSetByHeight RPC method. */

export interface GetValidatorSetByHeightResponseSDKType {
  block_height: Long;
  validators: ValidatorSDKType[];
  pagination?: PageResponseSDKType;
}
/** GetLatestValidatorSetRequest is the request type for the Query/GetValidatorSetByHeight RPC method. */

export interface GetLatestValidatorSetRequest {
  /** pagination defines an pagination for the request. */
  pagination?: PageRequest;
}
/** GetLatestValidatorSetRequest is the request type for the Query/GetValidatorSetByHeight RPC method. */

export interface GetLatestValidatorSetRequestSDKType {
  pagination?: PageRequestSDKType;
}
/** GetLatestValidatorSetResponse is the response type for the Query/GetValidatorSetByHeight RPC method. */

export interface GetLatestValidatorSetResponse {
  blockHeight: Long;
  validators: Validator[];
  /** pagination defines an pagination for the response. */

  pagination?: PageResponse;
}
/** GetLatestValidatorSetResponse is the response type for the Query/GetValidatorSetByHeight RPC method. */

export interface GetLatestValidatorSetResponseSDKType {
  block_height: Long;
  validators: ValidatorSDKType[];
  pagination?: PageResponseSDKType;
}
/** Validator is the type for the validator-set. */

export interface Validator {
  address: string;
  pubKey?: Any;
  votingPower: Long;
  proposerPriority: Long;
}
/** Validator is the type for the validator-set. */

export interface ValidatorSDKType {
  address: string;
  pub_key?: AnySDKType;
  voting_power: Long;
  proposer_priority: Long;
}
/** GetBlockByHeightRequest is the request type for the Query/GetBlockByHeight RPC method. */

export interface GetBlockByHeightRequest {
  height: Long;
}
/** GetBlockByHeightRequest is the request type for the Query/GetBlockByHeight RPC method. */

export interface GetBlockByHeightRequestSDKType {
  height: Long;
}
/** GetBlockByHeightResponse is the response type for the Query/GetBlockByHeight RPC method. */

export interface GetBlockByHeightResponse {
  blockId?: BlockID;
  /** Deprecated: please use `sdk_block` instead */

  block?: Block1;
  /** Since: cosmos-sdk 0.47 */

  sdkBlock?: Block2;
}
/** GetBlockByHeightResponse is the response type for the Query/GetBlockByHeight RPC method. */

export interface GetBlockByHeightResponseSDKType {
  block_id?: BlockIDSDKType;
  block?: Block1SDKType;
  sdk_block?: Block2SDKType;
}
/** GetLatestBlockRequest is the request type for the Query/GetLatestBlock RPC method. */

export interface GetLatestBlockRequest {}
/** GetLatestBlockRequest is the request type for the Query/GetLatestBlock RPC method. */

export interface GetLatestBlockRequestSDKType {}
/** GetLatestBlockResponse is the response type for the Query/GetLatestBlock RPC method. */

export interface GetLatestBlockResponse {
  blockId?: BlockID;
  /** Deprecated: please use `sdk_block` instead */

  block?: Block1;
  /** Since: cosmos-sdk 0.47 */

  sdkBlock?: Block2;
}
/** GetLatestBlockResponse is the response type for the Query/GetLatestBlock RPC method. */

export interface GetLatestBlockResponseSDKType {
  block_id?: BlockIDSDKType;
  block?: Block1SDKType;
  sdk_block?: Block2SDKType;
}
/** GetSyncingRequest is the request type for the Query/GetSyncing RPC method. */

export interface GetSyncingRequest {}
/** GetSyncingRequest is the request type for the Query/GetSyncing RPC method. */

export interface GetSyncingRequestSDKType {}
/** GetSyncingResponse is the response type for the Query/GetSyncing RPC method. */

export interface GetSyncingResponse {
  syncing: boolean;
}
/** GetSyncingResponse is the response type for the Query/GetSyncing RPC method. */

export interface GetSyncingResponseSDKType {
  syncing: boolean;
}
/** GetNodeInfoRequest is the request type for the Query/GetNodeInfo RPC method. */

export interface GetNodeInfoRequest {}
/** GetNodeInfoRequest is the request type for the Query/GetNodeInfo RPC method. */

export interface GetNodeInfoRequestSDKType {}
/** GetNodeInfoResponse is the response type for the Query/GetNodeInfo RPC method. */

export interface GetNodeInfoResponse {
  defaultNodeInfo?: DefaultNodeInfo;
  applicationVersion?: VersionInfo;
}
/** GetNodeInfoResponse is the response type for the Query/GetNodeInfo RPC method. */

export interface GetNodeInfoResponseSDKType {
  default_node_info?: DefaultNodeInfoSDKType;
  application_version?: VersionInfoSDKType;
}
/** VersionInfo is the type for the GetNodeInfoResponse message. */

export interface VersionInfo {
  name: string;
  appName: string;
  version: string;
  gitCommit: string;
  buildTags: string;
  goVersion: string;
  buildDeps: Module[];
  /** Since: cosmos-sdk 0.43 */

  cosmosSdkVersion: string;
}
/** VersionInfo is the type for the GetNodeInfoResponse message. */

export interface VersionInfoSDKType {
  name: string;
  app_name: string;
  version: string;
  git_commit: string;
  build_tags: string;
  go_version: string;
  build_deps: ModuleSDKType[];
  cosmos_sdk_version: string;
}
/** Module is the type for VersionInfo */

export interface Module {
  /** module path */
  path: string;
  /** module version */

  version: string;
  /** checksum */

  sum: string;
}
/** Module is the type for VersionInfo */

export interface ModuleSDKType {
  path: string;
  version: string;
  sum: string;
}
/** ABCIQueryRequest defines the request structure for the ABCIQuery gRPC query. */

export interface ABCIQueryRequest {
  data: Uint8Array;
  path: string;
  height: Long;
  prove: boolean;
}
/** ABCIQueryRequest defines the request structure for the ABCIQuery gRPC query. */

export interface ABCIQueryRequestSDKType {
  data: Uint8Array;
  path: string;
  height: Long;
  prove: boolean;
}
/**
 * ABCIQueryResponse defines the response structure for the ABCIQuery gRPC query.
 * 
 * Note: This type is a duplicate of the ResponseQuery proto type defined in
 * Tendermint.
 */

export interface ABCIQueryResponse {
  code: number;
  /** nondeterministic */

  log: string;
  /** nondeterministic */

  info: string;
  index: Long;
  key: Uint8Array;
  value: Uint8Array;
  proofOps?: ProofOps;
  height: Long;
  codespace: string;
}
/**
 * ABCIQueryResponse defines the response structure for the ABCIQuery gRPC query.
 * 
 * Note: This type is a duplicate of the ResponseQuery proto type defined in
 * Tendermint.
 */

export interface ABCIQueryResponseSDKType {
  code: number;
  log: string;
  info: string;
  index: Long;
  key: Uint8Array;
  value: Uint8Array;
  proof_ops?: ProofOpsSDKType;
  height: Long;
  codespace: string;
}
/**
 * ProofOp defines an operation used for calculating Merkle root. The data could
 * be arbitrary format, providing necessary data for example neighbouring node
 * hash.
 * 
 * Note: This type is a duplicate of the ProofOp proto type defined in Tendermint.
 */

export interface ProofOp {
  type: string;
  key: Uint8Array;
  data: Uint8Array;
}
/**
 * ProofOp defines an operation used for calculating Merkle root. The data could
 * be arbitrary format, providing necessary data for example neighbouring node
 * hash.
 * 
 * Note: This type is a duplicate of the ProofOp proto type defined in Tendermint.
 */

export interface ProofOpSDKType {
  type: string;
  key: Uint8Array;
  data: Uint8Array;
}
/**
 * ProofOps is Merkle proof defined by the list of ProofOps.
 * 
 * Note: This type is a duplicate of the ProofOps proto type defined in Tendermint.
 */

export interface ProofOps {
  ops: ProofOp[];
}
/**
 * ProofOps is Merkle proof defined by the list of ProofOps.
 * 
 * Note: This type is a duplicate of the ProofOps proto type defined in Tendermint.
 */

export interface ProofOpsSDKType {
  ops: ProofOpSDKType[];
}

function createBaseGetValidatorSetByHeightRequest(): GetValidatorSetByHeightRequest {
  return {
    height: Long.ZERO,
    pagination: undefined
  };
}

export const GetValidatorSetByHeightRequest = {
  encode(message: GetValidatorSetByHeightRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (!message.height.isZero()) {
      writer.uint32(8).int64(message.height);
    }

    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetValidatorSetByHeightRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetValidatorSetByHeightRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.height = (reader.int64() as Long);
          break;

        case 2:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<GetValidatorSetByHeightRequest>): GetValidatorSetByHeightRequest {
    const message = createBaseGetValidatorSetByHeightRequest();
    message.height = object.height !== undefined && object.height !== null ? Long.fromValue(object.height) : Long.ZERO;
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseGetValidatorSetByHeightResponse(): GetValidatorSetByHeightResponse {
  return {
    blockHeight: Long.ZERO,
    validators: [],
    pagination: undefined
  };
}

export const GetValidatorSetByHeightResponse = {
  encode(message: GetValidatorSetByHeightResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (!message.blockHeight.isZero()) {
      writer.uint32(8).int64(message.blockHeight);
    }

    for (const v of message.validators) {
      Validator.encode(v!, writer.uint32(18).fork()).ldelim();
    }

    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(26).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetValidatorSetByHeightResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetValidatorSetByHeightResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.blockHeight = (reader.int64() as Long);
          break;

        case 2:
          message.validators.push(Validator.decode(reader, reader.uint32()));
          break;

        case 3:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<GetValidatorSetByHeightResponse>): GetValidatorSetByHeightResponse {
    const message = createBaseGetValidatorSetByHeightResponse();
    message.blockHeight = object.blockHeight !== undefined && object.blockHeight !== null ? Long.fromValue(object.blockHeight) : Long.ZERO;
    message.validators = object.validators?.map(e => Validator.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseGetLatestValidatorSetRequest(): GetLatestValidatorSetRequest {
  return {
    pagination: undefined
  };
}

export const GetLatestValidatorSetRequest = {
  encode(message: GetLatestValidatorSetRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetLatestValidatorSetRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetLatestValidatorSetRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<GetLatestValidatorSetRequest>): GetLatestValidatorSetRequest {
    const message = createBaseGetLatestValidatorSetRequest();
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseGetLatestValidatorSetResponse(): GetLatestValidatorSetResponse {
  return {
    blockHeight: Long.ZERO,
    validators: [],
    pagination: undefined
  };
}

export const GetLatestValidatorSetResponse = {
  encode(message: GetLatestValidatorSetResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (!message.blockHeight.isZero()) {
      writer.uint32(8).int64(message.blockHeight);
    }

    for (const v of message.validators) {
      Validator.encode(v!, writer.uint32(18).fork()).ldelim();
    }

    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(26).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetLatestValidatorSetResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetLatestValidatorSetResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.blockHeight = (reader.int64() as Long);
          break;

        case 2:
          message.validators.push(Validator.decode(reader, reader.uint32()));
          break;

        case 3:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<GetLatestValidatorSetResponse>): GetLatestValidatorSetResponse {
    const message = createBaseGetLatestValidatorSetResponse();
    message.blockHeight = object.blockHeight !== undefined && object.blockHeight !== null ? Long.fromValue(object.blockHeight) : Long.ZERO;
    message.validators = object.validators?.map(e => Validator.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseValidator(): Validator {
  return {
    address: "",
    pubKey: undefined,
    votingPower: Long.ZERO,
    proposerPriority: Long.ZERO
  };
}

export const Validator = {
  encode(message: Validator, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }

    if (message.pubKey !== undefined) {
      Any.encode(message.pubKey, writer.uint32(18).fork()).ldelim();
    }

    if (!message.votingPower.isZero()) {
      writer.uint32(24).int64(message.votingPower);
    }

    if (!message.proposerPriority.isZero()) {
      writer.uint32(32).int64(message.proposerPriority);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Validator {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseValidator();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;

        case 2:
          message.pubKey = Any.decode(reader, reader.uint32());
          break;

        case 3:
          message.votingPower = (reader.int64() as Long);
          break;

        case 4:
          message.proposerPriority = (reader.int64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<Validator>): Validator {
    const message = createBaseValidator();
    message.address = object.address ?? "";
    message.pubKey = object.pubKey !== undefined && object.pubKey !== null ? Any.fromPartial(object.pubKey) : undefined;
    message.votingPower = object.votingPower !== undefined && object.votingPower !== null ? Long.fromValue(object.votingPower) : Long.ZERO;
    message.proposerPriority = object.proposerPriority !== undefined && object.proposerPriority !== null ? Long.fromValue(object.proposerPriority) : Long.ZERO;
    return message;
  }

};

function createBaseGetBlockByHeightRequest(): GetBlockByHeightRequest {
  return {
    height: Long.ZERO
  };
}

export const GetBlockByHeightRequest = {
  encode(message: GetBlockByHeightRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (!message.height.isZero()) {
      writer.uint32(8).int64(message.height);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetBlockByHeightRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetBlockByHeightRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.height = (reader.int64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<GetBlockByHeightRequest>): GetBlockByHeightRequest {
    const message = createBaseGetBlockByHeightRequest();
    message.height = object.height !== undefined && object.height !== null ? Long.fromValue(object.height) : Long.ZERO;
    return message;
  }

};

function createBaseGetBlockByHeightResponse(): GetBlockByHeightResponse {
  return {
    blockId: undefined,
    block: undefined,
    sdkBlock: undefined
  };
}

export const GetBlockByHeightResponse = {
  encode(message: GetBlockByHeightResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.blockId !== undefined) {
      BlockID.encode(message.blockId, writer.uint32(10).fork()).ldelim();
    }

    if (message.block !== undefined) {
      Block1.encode(message.block, writer.uint32(18).fork()).ldelim();
    }

    if (message.sdkBlock !== undefined) {
      Block2.encode(message.sdkBlock, writer.uint32(26).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetBlockByHeightResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetBlockByHeightResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.blockId = BlockID.decode(reader, reader.uint32());
          break;

        case 2:
          message.block = Block1.decode(reader, reader.uint32());
          break;

        case 3:
          message.sdkBlock = Block2.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<GetBlockByHeightResponse>): GetBlockByHeightResponse {
    const message = createBaseGetBlockByHeightResponse();
    message.blockId = object.blockId !== undefined && object.blockId !== null ? BlockID.fromPartial(object.blockId) : undefined;
    message.block = object.block !== undefined && object.block !== null ? Block1.fromPartial(object.block) : undefined;
    message.sdkBlock = object.sdkBlock !== undefined && object.sdkBlock !== null ? Block2.fromPartial(object.sdkBlock) : undefined;
    return message;
  }

};

function createBaseGetLatestBlockRequest(): GetLatestBlockRequest {
  return {};
}

export const GetLatestBlockRequest = {
  encode(_: GetLatestBlockRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetLatestBlockRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetLatestBlockRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(_: DeepPartial<GetLatestBlockRequest>): GetLatestBlockRequest {
    const message = createBaseGetLatestBlockRequest();
    return message;
  }

};

function createBaseGetLatestBlockResponse(): GetLatestBlockResponse {
  return {
    blockId: undefined,
    block: undefined,
    sdkBlock: undefined
  };
}

export const GetLatestBlockResponse = {
  encode(message: GetLatestBlockResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.blockId !== undefined) {
      BlockID.encode(message.blockId, writer.uint32(10).fork()).ldelim();
    }

    if (message.block !== undefined) {
      Block1.encode(message.block, writer.uint32(18).fork()).ldelim();
    }

    if (message.sdkBlock !== undefined) {
      Block2.encode(message.sdkBlock, writer.uint32(26).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetLatestBlockResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetLatestBlockResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.blockId = BlockID.decode(reader, reader.uint32());
          break;

        case 2:
          message.block = Block1.decode(reader, reader.uint32());
          break;

        case 3:
          message.sdkBlock = Block2.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<GetLatestBlockResponse>): GetLatestBlockResponse {
    const message = createBaseGetLatestBlockResponse();
    message.blockId = object.blockId !== undefined && object.blockId !== null ? BlockID.fromPartial(object.blockId) : undefined;
    message.block = object.block !== undefined && object.block !== null ? Block1.fromPartial(object.block) : undefined;
    message.sdkBlock = object.sdkBlock !== undefined && object.sdkBlock !== null ? Block2.fromPartial(object.sdkBlock) : undefined;
    return message;
  }

};

function createBaseGetSyncingRequest(): GetSyncingRequest {
  return {};
}

export const GetSyncingRequest = {
  encode(_: GetSyncingRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetSyncingRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetSyncingRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(_: DeepPartial<GetSyncingRequest>): GetSyncingRequest {
    const message = createBaseGetSyncingRequest();
    return message;
  }

};

function createBaseGetSyncingResponse(): GetSyncingResponse {
  return {
    syncing: false
  };
}

export const GetSyncingResponse = {
  encode(message: GetSyncingResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.syncing === true) {
      writer.uint32(8).bool(message.syncing);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetSyncingResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetSyncingResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.syncing = reader.bool();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<GetSyncingResponse>): GetSyncingResponse {
    const message = createBaseGetSyncingResponse();
    message.syncing = object.syncing ?? false;
    return message;
  }

};

function createBaseGetNodeInfoRequest(): GetNodeInfoRequest {
  return {};
}

export const GetNodeInfoRequest = {
  encode(_: GetNodeInfoRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetNodeInfoRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetNodeInfoRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(_: DeepPartial<GetNodeInfoRequest>): GetNodeInfoRequest {
    const message = createBaseGetNodeInfoRequest();
    return message;
  }

};

function createBaseGetNodeInfoResponse(): GetNodeInfoResponse {
  return {
    defaultNodeInfo: undefined,
    applicationVersion: undefined
  };
}

export const GetNodeInfoResponse = {
  encode(message: GetNodeInfoResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.defaultNodeInfo !== undefined) {
      DefaultNodeInfo.encode(message.defaultNodeInfo, writer.uint32(10).fork()).ldelim();
    }

    if (message.applicationVersion !== undefined) {
      VersionInfo.encode(message.applicationVersion, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetNodeInfoResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetNodeInfoResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.defaultNodeInfo = DefaultNodeInfo.decode(reader, reader.uint32());
          break;

        case 2:
          message.applicationVersion = VersionInfo.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<GetNodeInfoResponse>): GetNodeInfoResponse {
    const message = createBaseGetNodeInfoResponse();
    message.defaultNodeInfo = object.defaultNodeInfo !== undefined && object.defaultNodeInfo !== null ? DefaultNodeInfo.fromPartial(object.defaultNodeInfo) : undefined;
    message.applicationVersion = object.applicationVersion !== undefined && object.applicationVersion !== null ? VersionInfo.fromPartial(object.applicationVersion) : undefined;
    return message;
  }

};

function createBaseVersionInfo(): VersionInfo {
  return {
    name: "",
    appName: "",
    version: "",
    gitCommit: "",
    buildTags: "",
    goVersion: "",
    buildDeps: [],
    cosmosSdkVersion: ""
  };
}

export const VersionInfo = {
  encode(message: VersionInfo, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }

    if (message.appName !== "") {
      writer.uint32(18).string(message.appName);
    }

    if (message.version !== "") {
      writer.uint32(26).string(message.version);
    }

    if (message.gitCommit !== "") {
      writer.uint32(34).string(message.gitCommit);
    }

    if (message.buildTags !== "") {
      writer.uint32(42).string(message.buildTags);
    }

    if (message.goVersion !== "") {
      writer.uint32(50).string(message.goVersion);
    }

    for (const v of message.buildDeps) {
      Module.encode(v!, writer.uint32(58).fork()).ldelim();
    }

    if (message.cosmosSdkVersion !== "") {
      writer.uint32(66).string(message.cosmosSdkVersion);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VersionInfo {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVersionInfo();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.name = reader.string();
          break;

        case 2:
          message.appName = reader.string();
          break;

        case 3:
          message.version = reader.string();
          break;

        case 4:
          message.gitCommit = reader.string();
          break;

        case 5:
          message.buildTags = reader.string();
          break;

        case 6:
          message.goVersion = reader.string();
          break;

        case 7:
          message.buildDeps.push(Module.decode(reader, reader.uint32()));
          break;

        case 8:
          message.cosmosSdkVersion = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<VersionInfo>): VersionInfo {
    const message = createBaseVersionInfo();
    message.name = object.name ?? "";
    message.appName = object.appName ?? "";
    message.version = object.version ?? "";
    message.gitCommit = object.gitCommit ?? "";
    message.buildTags = object.buildTags ?? "";
    message.goVersion = object.goVersion ?? "";
    message.buildDeps = object.buildDeps?.map(e => Module.fromPartial(e)) || [];
    message.cosmosSdkVersion = object.cosmosSdkVersion ?? "";
    return message;
  }

};

function createBaseModule(): Module {
  return {
    path: "",
    version: "",
    sum: ""
  };
}

export const Module = {
  encode(message: Module, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.path !== "") {
      writer.uint32(10).string(message.path);
    }

    if (message.version !== "") {
      writer.uint32(18).string(message.version);
    }

    if (message.sum !== "") {
      writer.uint32(26).string(message.sum);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Module {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseModule();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.path = reader.string();
          break;

        case 2:
          message.version = reader.string();
          break;

        case 3:
          message.sum = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<Module>): Module {
    const message = createBaseModule();
    message.path = object.path ?? "";
    message.version = object.version ?? "";
    message.sum = object.sum ?? "";
    return message;
  }

};

function createBaseABCIQueryRequest(): ABCIQueryRequest {
  return {
    data: new Uint8Array(),
    path: "",
    height: Long.ZERO,
    prove: false
  };
}

export const ABCIQueryRequest = {
  encode(message: ABCIQueryRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.data.length !== 0) {
      writer.uint32(10).bytes(message.data);
    }

    if (message.path !== "") {
      writer.uint32(18).string(message.path);
    }

    if (!message.height.isZero()) {
      writer.uint32(24).int64(message.height);
    }

    if (message.prove === true) {
      writer.uint32(32).bool(message.prove);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ABCIQueryRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseABCIQueryRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.data = reader.bytes();
          break;

        case 2:
          message.path = reader.string();
          break;

        case 3:
          message.height = (reader.int64() as Long);
          break;

        case 4:
          message.prove = reader.bool();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<ABCIQueryRequest>): ABCIQueryRequest {
    const message = createBaseABCIQueryRequest();
    message.data = object.data ?? new Uint8Array();
    message.path = object.path ?? "";
    message.height = object.height !== undefined && object.height !== null ? Long.fromValue(object.height) : Long.ZERO;
    message.prove = object.prove ?? false;
    return message;
  }

};

function createBaseABCIQueryResponse(): ABCIQueryResponse {
  return {
    code: 0,
    log: "",
    info: "",
    index: Long.ZERO,
    key: new Uint8Array(),
    value: new Uint8Array(),
    proofOps: undefined,
    height: Long.ZERO,
    codespace: ""
  };
}

export const ABCIQueryResponse = {
  encode(message: ABCIQueryResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.code !== 0) {
      writer.uint32(8).uint32(message.code);
    }

    if (message.log !== "") {
      writer.uint32(26).string(message.log);
    }

    if (message.info !== "") {
      writer.uint32(34).string(message.info);
    }

    if (!message.index.isZero()) {
      writer.uint32(40).int64(message.index);
    }

    if (message.key.length !== 0) {
      writer.uint32(50).bytes(message.key);
    }

    if (message.value.length !== 0) {
      writer.uint32(58).bytes(message.value);
    }

    if (message.proofOps !== undefined) {
      ProofOps.encode(message.proofOps, writer.uint32(66).fork()).ldelim();
    }

    if (!message.height.isZero()) {
      writer.uint32(72).int64(message.height);
    }

    if (message.codespace !== "") {
      writer.uint32(82).string(message.codespace);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ABCIQueryResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseABCIQueryResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.code = reader.uint32();
          break;

        case 3:
          message.log = reader.string();
          break;

        case 4:
          message.info = reader.string();
          break;

        case 5:
          message.index = (reader.int64() as Long);
          break;

        case 6:
          message.key = reader.bytes();
          break;

        case 7:
          message.value = reader.bytes();
          break;

        case 8:
          message.proofOps = ProofOps.decode(reader, reader.uint32());
          break;

        case 9:
          message.height = (reader.int64() as Long);
          break;

        case 10:
          message.codespace = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<ABCIQueryResponse>): ABCIQueryResponse {
    const message = createBaseABCIQueryResponse();
    message.code = object.code ?? 0;
    message.log = object.log ?? "";
    message.info = object.info ?? "";
    message.index = object.index !== undefined && object.index !== null ? Long.fromValue(object.index) : Long.ZERO;
    message.key = object.key ?? new Uint8Array();
    message.value = object.value ?? new Uint8Array();
    message.proofOps = object.proofOps !== undefined && object.proofOps !== null ? ProofOps.fromPartial(object.proofOps) : undefined;
    message.height = object.height !== undefined && object.height !== null ? Long.fromValue(object.height) : Long.ZERO;
    message.codespace = object.codespace ?? "";
    return message;
  }

};

function createBaseProofOp(): ProofOp {
  return {
    type: "",
    key: new Uint8Array(),
    data: new Uint8Array()
  };
}

export const ProofOp = {
  encode(message: ProofOp, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.type !== "") {
      writer.uint32(10).string(message.type);
    }

    if (message.key.length !== 0) {
      writer.uint32(18).bytes(message.key);
    }

    if (message.data.length !== 0) {
      writer.uint32(26).bytes(message.data);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ProofOp {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseProofOp();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.type = reader.string();
          break;

        case 2:
          message.key = reader.bytes();
          break;

        case 3:
          message.data = reader.bytes();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<ProofOp>): ProofOp {
    const message = createBaseProofOp();
    message.type = object.type ?? "";
    message.key = object.key ?? new Uint8Array();
    message.data = object.data ?? new Uint8Array();
    return message;
  }

};

function createBaseProofOps(): ProofOps {
  return {
    ops: []
  };
}

export const ProofOps = {
  encode(message: ProofOps, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.ops) {
      ProofOp.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ProofOps {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseProofOps();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.ops.push(ProofOp.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<ProofOps>): ProofOps {
    const message = createBaseProofOps();
    message.ops = object.ops?.map(e => ProofOp.fromPartial(e)) || [];
    return message;
  }

};