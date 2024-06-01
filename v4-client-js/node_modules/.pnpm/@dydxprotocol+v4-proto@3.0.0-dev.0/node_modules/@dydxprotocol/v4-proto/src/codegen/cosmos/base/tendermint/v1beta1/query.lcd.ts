import { setPaginationParams } from "../../../../helpers";
import { LCDClient } from "@osmonauts/lcd";
import { GetNodeInfoRequest, GetNodeInfoResponseSDKType, GetSyncingRequest, GetSyncingResponseSDKType, GetLatestBlockRequest, GetLatestBlockResponseSDKType, GetBlockByHeightRequest, GetBlockByHeightResponseSDKType, GetLatestValidatorSetRequest, GetLatestValidatorSetResponseSDKType, GetValidatorSetByHeightRequest, GetValidatorSetByHeightResponseSDKType, ABCIQueryRequest, ABCIQueryResponseSDKType } from "./query";
export class LCDQueryClient {
  req: LCDClient;

  constructor({
    requestClient
  }: {
    requestClient: LCDClient;
  }) {
    this.req = requestClient;
    this.getNodeInfo = this.getNodeInfo.bind(this);
    this.getSyncing = this.getSyncing.bind(this);
    this.getLatestBlock = this.getLatestBlock.bind(this);
    this.getBlockByHeight = this.getBlockByHeight.bind(this);
    this.getLatestValidatorSet = this.getLatestValidatorSet.bind(this);
    this.getValidatorSetByHeight = this.getValidatorSetByHeight.bind(this);
    this.aBCIQuery = this.aBCIQuery.bind(this);
  }
  /* GetNodeInfo queries the current node info. */


  async getNodeInfo(_params: GetNodeInfoRequest = {}): Promise<GetNodeInfoResponseSDKType> {
    const endpoint = `cosmos/base/tendermint/v1beta1/node_info`;
    return await this.req.get<GetNodeInfoResponseSDKType>(endpoint);
  }
  /* GetSyncing queries node syncing. */


  async getSyncing(_params: GetSyncingRequest = {}): Promise<GetSyncingResponseSDKType> {
    const endpoint = `cosmos/base/tendermint/v1beta1/syncing`;
    return await this.req.get<GetSyncingResponseSDKType>(endpoint);
  }
  /* GetLatestBlock returns the latest block. */


  async getLatestBlock(_params: GetLatestBlockRequest = {}): Promise<GetLatestBlockResponseSDKType> {
    const endpoint = `cosmos/base/tendermint/v1beta1/blocks/latest`;
    return await this.req.get<GetLatestBlockResponseSDKType>(endpoint);
  }
  /* GetBlockByHeight queries block for given height. */


  async getBlockByHeight(params: GetBlockByHeightRequest): Promise<GetBlockByHeightResponseSDKType> {
    const endpoint = `cosmos/base/tendermint/v1beta1/blocks/${params.height}`;
    return await this.req.get<GetBlockByHeightResponseSDKType>(endpoint);
  }
  /* GetLatestValidatorSet queries latest validator-set. */


  async getLatestValidatorSet(params: GetLatestValidatorSetRequest = {
    pagination: undefined
  }): Promise<GetLatestValidatorSetResponseSDKType> {
    const options: any = {
      params: {}
    };

    if (typeof params?.pagination !== "undefined") {
      setPaginationParams(options, params.pagination);
    }

    const endpoint = `cosmos/base/tendermint/v1beta1/validatorsets/latest`;
    return await this.req.get<GetLatestValidatorSetResponseSDKType>(endpoint, options);
  }
  /* GetValidatorSetByHeight queries validator-set at a given height. */


  async getValidatorSetByHeight(params: GetValidatorSetByHeightRequest): Promise<GetValidatorSetByHeightResponseSDKType> {
    const options: any = {
      params: {}
    };

    if (typeof params?.pagination !== "undefined") {
      setPaginationParams(options, params.pagination);
    }

    const endpoint = `cosmos/base/tendermint/v1beta1/validatorsets/${params.height}`;
    return await this.req.get<GetValidatorSetByHeightResponseSDKType>(endpoint, options);
  }
  /* ABCIQuery defines a query handler that supports ABCI queries directly to the
   application, bypassing Tendermint completely. The ABCI query must contain
   a valid and supported path, including app, custom, p2p, and store.
  
   Since: cosmos-sdk 0.46 */


  async aBCIQuery(params: ABCIQueryRequest): Promise<ABCIQueryResponseSDKType> {
    const options: any = {
      params: {}
    };

    if (typeof params?.data !== "undefined") {
      options.params.data = params.data;
    }

    if (typeof params?.path !== "undefined") {
      options.params.path = params.path;
    }

    if (typeof params?.height !== "undefined") {
      options.params.height = params.height;
    }

    if (typeof params?.prove !== "undefined") {
      options.params.prove = params.prove;
    }

    const endpoint = `cosmos/base/tendermint/v1beta1/abci_query`;
    return await this.req.get<ABCIQueryResponseSDKType>(endpoint, options);
  }

}