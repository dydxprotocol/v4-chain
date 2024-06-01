import { Rpc } from "../../../helpers";
import { QueryClient } from "@cosmjs/stargate";
import { SimulateRequest, SimulateResponse, GetTxRequest, GetTxResponse, BroadcastTxRequest, BroadcastTxResponse, GetTxsEventRequest, GetTxsEventResponse, GetBlockWithTxsRequest, GetBlockWithTxsResponse, TxDecodeRequest, TxDecodeResponse, TxEncodeRequest, TxEncodeResponse, TxEncodeAminoRequest, TxEncodeAminoResponse, TxDecodeAminoRequest, TxDecodeAminoResponse } from "./service";
/** Service defines a gRPC service for interacting with transactions. */
export interface Service {
    /** Simulate simulates executing a transaction for estimating gas usage. */
    simulate(request: SimulateRequest): Promise<SimulateResponse>;
    /** GetTx fetches a tx by hash. */
    getTx(request: GetTxRequest): Promise<GetTxResponse>;
    /** BroadcastTx broadcast transaction. */
    broadcastTx(request: BroadcastTxRequest): Promise<BroadcastTxResponse>;
    /** GetTxsEvent fetches txs by event. */
    getTxsEvent(request: GetTxsEventRequest): Promise<GetTxsEventResponse>;
    /**
     * GetBlockWithTxs fetches a block with decoded txs.
     *
     * Since: cosmos-sdk 0.45.2
     */
    getBlockWithTxs(request: GetBlockWithTxsRequest): Promise<GetBlockWithTxsResponse>;
    /**
     * TxDecode decodes the transaction.
     *
     * Since: cosmos-sdk 0.47
     */
    txDecode(request: TxDecodeRequest): Promise<TxDecodeResponse>;
    /**
     * TxEncode encodes the transaction.
     *
     * Since: cosmos-sdk 0.47
     */
    txEncode(request: TxEncodeRequest): Promise<TxEncodeResponse>;
    /**
     * TxEncodeAmino encodes an Amino transaction from JSON to encoded bytes.
     *
     * Since: cosmos-sdk 0.47
     */
    txEncodeAmino(request: TxEncodeAminoRequest): Promise<TxEncodeAminoResponse>;
    /**
     * TxDecodeAmino decodes an Amino transaction from encoded bytes to JSON.
     *
     * Since: cosmos-sdk 0.47
     */
    txDecodeAmino(request: TxDecodeAminoRequest): Promise<TxDecodeAminoResponse>;
}
export declare class ServiceClientImpl implements Service {
    private readonly rpc;
    constructor(rpc: Rpc);
    simulate(request: SimulateRequest): Promise<SimulateResponse>;
    getTx(request: GetTxRequest): Promise<GetTxResponse>;
    broadcastTx(request: BroadcastTxRequest): Promise<BroadcastTxResponse>;
    getTxsEvent(request: GetTxsEventRequest): Promise<GetTxsEventResponse>;
    getBlockWithTxs(request: GetBlockWithTxsRequest): Promise<GetBlockWithTxsResponse>;
    txDecode(request: TxDecodeRequest): Promise<TxDecodeResponse>;
    txEncode(request: TxEncodeRequest): Promise<TxEncodeResponse>;
    txEncodeAmino(request: TxEncodeAminoRequest): Promise<TxEncodeAminoResponse>;
    txDecodeAmino(request: TxDecodeAminoRequest): Promise<TxDecodeAminoResponse>;
}
export declare const createRpcQueryExtension: (base: QueryClient) => {
    simulate(request: SimulateRequest): Promise<SimulateResponse>;
    getTx(request: GetTxRequest): Promise<GetTxResponse>;
    broadcastTx(request: BroadcastTxRequest): Promise<BroadcastTxResponse>;
    getTxsEvent(request: GetTxsEventRequest): Promise<GetTxsEventResponse>;
    getBlockWithTxs(request: GetBlockWithTxsRequest): Promise<GetBlockWithTxsResponse>;
    txDecode(request: TxDecodeRequest): Promise<TxDecodeResponse>;
    txEncode(request: TxEncodeRequest): Promise<TxEncodeResponse>;
    txEncodeAmino(request: TxEncodeAminoRequest): Promise<TxEncodeAminoResponse>;
    txDecodeAmino(request: TxDecodeAminoRequest): Promise<TxDecodeAminoResponse>;
};
