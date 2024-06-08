import { CometClient, HttpEndpoint } from "@cosmjs/tendermint-rpc";
import { MsgData } from "cosmjs-types/cosmos/base/abci/v1beta1/abci";
import { Coin } from "cosmjs-types/cosmos/base/v1beta1/coin";
import { Account, AccountParser } from "./accounts";
import { Event } from "./events";
import { AuthExtension, BankExtension, StakingExtension, TxExtension } from "./modules";
import { QueryClient } from "./queryclient";
import { SearchTxQuery } from "./search";
export declare class TimeoutError extends Error {
    readonly txId: string;
    constructor(message: string, txId: string);
}
export interface BlockHeader {
    readonly version: {
        readonly block: string;
        readonly app: string;
    };
    readonly height: number;
    readonly chainId: string;
    /** An RFC 3339 time string like e.g. '2020-02-15T10:39:10.4696305Z' */
    readonly time: string;
}
export interface Block {
    /** The ID is a hash of the block header (uppercase hex) */
    readonly id: string;
    readonly header: BlockHeader;
    /** Array of raw transactions */
    readonly txs: readonly Uint8Array[];
}
/** A transaction that is indexed as part of the transaction history */
export interface IndexedTx {
    readonly height: number;
    /** The position of the transaction within the block. This is a 0-based index. */
    readonly txIndex: number;
    /** Transaction hash (might be used as transaction ID). Guaranteed to be non-empty upper-case hex */
    readonly hash: string;
    /** Transaction execution error code. 0 on success. */
    readonly code: number;
    readonly events: readonly Event[];
    /**
     * A string-based log document.
     *
     * This currently seems to merge attributes of multiple events into one event per type
     * (https://github.com/tendermint/tendermint/issues/9595). You might want to use the `events`
     * field instead.
     *
     * @deprecated This field is not filled anymore in Cosmos SDK 0.50+ (https://github.com/cosmos/cosmos-sdk/pull/15845).
     * Please consider using `events` instead.
     */
    readonly rawLog: string;
    /**
     * Raw transaction bytes stored in Tendermint.
     *
     * If you hash this, you get the transaction hash (= transaction ID):
     *
     * ```js
     * import { sha256 } from "@cosmjs/crypto";
     * import { toHex } from "@cosmjs/encoding";
     *
     * const transactionId = toHex(sha256(indexTx.tx)).toUpperCase();
     * ```
     *
     * Use `decodeTxRaw` from @cosmjs/proto-signing to decode this.
     */
    readonly tx: Uint8Array;
    /**
     * The message responses of the [TxMsgData](https://github.com/cosmos/cosmos-sdk/blob/v0.46.3/proto/cosmos/base/abci/v1beta1/abci.proto#L128-L140)
     * as `Any`s.
     * This field is an empty list for chains running Cosmos SDK < 0.46.
     */
    readonly msgResponses: Array<{
        readonly typeUrl: string;
        readonly value: Uint8Array;
    }>;
    readonly gasUsed: bigint;
    readonly gasWanted: bigint;
}
export interface SequenceResponse {
    readonly accountNumber: number;
    readonly sequence: number;
}
/**
 * The response after successfully broadcasting a transaction.
 * Success or failure refer to the execution result.
 */
export interface DeliverTxResponse {
    readonly height: number;
    /** The position of the transaction within the block. This is a 0-based index. */
    readonly txIndex: number;
    /** Error code. The transaction suceeded if and only if code is 0. */
    readonly code: number;
    readonly transactionHash: string;
    readonly events: readonly Event[];
    /**
     * A string-based log document.
     *
     * This currently seems to merge attributes of multiple events into one event per type
     * (https://github.com/tendermint/tendermint/issues/9595). You might want to use the `events`
     * field instead.
     *
     * @deprecated This field is not filled anymore in Cosmos SDK 0.50+ (https://github.com/cosmos/cosmos-sdk/pull/15845).
     * Please consider using `events` instead.
     */
    readonly rawLog?: string;
    /** @deprecated Use `msgResponses` instead. */
    readonly data?: readonly MsgData[];
    /**
     * The message responses of the [TxMsgData](https://github.com/cosmos/cosmos-sdk/blob/v0.46.3/proto/cosmos/base/abci/v1beta1/abci.proto#L128-L140)
     * as `Any`s.
     * This field is an empty list for chains running Cosmos SDK < 0.46.
     */
    readonly msgResponses: Array<{
        readonly typeUrl: string;
        readonly value: Uint8Array;
    }>;
    readonly gasUsed: bigint;
    readonly gasWanted: bigint;
}
export declare function isDeliverTxFailure(result: DeliverTxResponse): boolean;
export declare function isDeliverTxSuccess(result: DeliverTxResponse): boolean;
/**
 * Ensures the given result is a success. Throws a detailed error message otherwise.
 */
export declare function assertIsDeliverTxSuccess(result: DeliverTxResponse): void;
/**
 * Ensures the given result is a failure. Throws a detailed error message otherwise.
 */
export declare function assertIsDeliverTxFailure(result: DeliverTxResponse): void;
/**
 * An error when broadcasting the transaction. This contains the CheckTx errors
 * from the blockchain. Once a transaction is included in a block no BroadcastTxError
 * is thrown, even if the execution fails (DeliverTx errors).
 */
export declare class BroadcastTxError extends Error {
    readonly code: number;
    readonly codespace: string;
    readonly log: string | undefined;
    constructor(code: number, codespace: string, log: string | undefined);
}
/** Use for testing only */
export interface PrivateStargateClient {
    readonly cometClient: CometClient | undefined;
}
export interface StargateClientOptions {
    readonly accountParser?: AccountParser;
}
export declare class StargateClient {
    private readonly cometClient;
    private readonly queryClient;
    private chainId;
    private readonly accountParser;
    /**
     * Creates an instance by connecting to the given CometBFT RPC endpoint.
     *
     * This uses auto-detection to decide between a CometBFT 0.38, Tendermint 0.37 and 0.34 client.
     * To set the Comet client explicitly, use `create`.
     */
    static connect(endpoint: string | HttpEndpoint, options?: StargateClientOptions): Promise<StargateClient>;
    /**
     * Creates an instance from a manually created Comet client.
     * Use this to use `Comet38Client` or `Tendermint37Client` instead of `Tendermint34Client`.
     */
    static create(cometClient: CometClient, options?: StargateClientOptions): Promise<StargateClient>;
    protected constructor(cometClient: CometClient | undefined, options: StargateClientOptions);
    protected getCometClient(): CometClient | undefined;
    protected forceGetCometClient(): CometClient;
    protected getQueryClient(): (QueryClient & AuthExtension & BankExtension & StakingExtension & TxExtension) | undefined;
    protected forceGetQueryClient(): QueryClient & AuthExtension & BankExtension & StakingExtension & TxExtension;
    getChainId(): Promise<string>;
    getHeight(): Promise<number>;
    getAccount(searchAddress: string): Promise<Account | null>;
    getSequence(address: string): Promise<SequenceResponse>;
    getBlock(height?: number): Promise<Block>;
    getBalance(address: string, searchDenom: string): Promise<Coin>;
    /**
     * Queries all balances for all denoms that belong to this address.
     *
     * Uses the grpc queries (which iterates over the store internally), and we cannot get
     * proofs from such a method.
     */
    getAllBalances(address: string): Promise<readonly Coin[]>;
    getBalanceStaked(address: string): Promise<Coin | null>;
    getDelegation(delegatorAddress: string, validatorAddress: string): Promise<Coin | null>;
    getTx(id: string): Promise<IndexedTx | null>;
    searchTx(query: SearchTxQuery): Promise<IndexedTx[]>;
    disconnect(): void;
    /**
     * Broadcasts a signed transaction to the network and monitors its inclusion in a block.
     *
     * If broadcasting is rejected by the node for some reason (e.g. because of a CheckTx failure),
     * an error is thrown.
     *
     * If the transaction is not included in a block before the provided timeout, this errors with a `TimeoutError`.
     *
     * If the transaction is included in a block, a `DeliverTxResponse` is returned. The caller then
     * usually needs to check for execution success or failure.
     */
    broadcastTx(tx: Uint8Array, timeoutMs?: number, pollIntervalMs?: number): Promise<DeliverTxResponse>;
    /**
     * Broadcasts a signed transaction to the network without monitoring it.
     *
     * If broadcasting is rejected by the node for some reason (e.g. because of a CheckTx failure),
     * an error is thrown.
     *
     * If the transaction is broadcasted, a `string` containing the hash of the transaction is returned. The caller then
     * usually needs to check if the transaction was included in a block and was successful.
     *
     * @returns Returns the hash of the transaction
     */
    broadcastTxSync(tx: Uint8Array): Promise<string>;
    private txsQuery;
}
