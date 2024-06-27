import { StdFee } from "@cosmjs/amino";
import { EncodeObject, GeneratedType, OfflineSigner, Registry } from "@cosmjs/proto-signing";
import { CometClient, HttpEndpoint } from "@cosmjs/tendermint-rpc";
import { Coin } from "cosmjs-types/cosmos/base/v1beta1/coin";
import { TxRaw } from "cosmjs-types/cosmos/tx/v1beta1/tx";
import { Height } from "cosmjs-types/ibc/core/client/v1/client";
import { AminoConverters, AminoTypes } from "./aminotypes";
import { GasPrice } from "./fee";
import { DeliverTxResponse, StargateClient, StargateClientOptions } from "./stargateclient";
export declare const defaultRegistryTypes: ReadonlyArray<[string, GeneratedType]>;
/**
 * Signing information for a single signer that is not included in the transaction.
 *
 * @see https://github.com/cosmos/cosmos-sdk/blob/v0.42.2/x/auth/signing/sign_mode_handler.go#L23-L37
 */
export interface SignerData {
    readonly accountNumber: number;
    readonly sequence: number;
    readonly chainId: string;
}
/** Use for testing only */
export interface PrivateSigningStargateClient {
    readonly registry: Registry;
}
export interface SigningStargateClientOptions extends StargateClientOptions {
    readonly registry?: Registry;
    readonly aminoTypes?: AminoTypes;
    readonly broadcastTimeoutMs?: number;
    readonly broadcastPollIntervalMs?: number;
    readonly gasPrice?: GasPrice;
}
export declare function createDefaultAminoConverters(): AminoConverters;
export declare class SigningStargateClient extends StargateClient {
    readonly registry: Registry;
    readonly broadcastTimeoutMs: number | undefined;
    readonly broadcastPollIntervalMs: number | undefined;
    private readonly signer;
    private readonly aminoTypes;
    private readonly gasPrice;
    private readonly defaultGasMultiplier;
    /**
     * Creates an instance by connecting to the given CometBFT RPC endpoint.
     *
     * This uses auto-detection to decide between a CometBFT 0.38, Tendermint 0.37 and 0.34 client.
     * To set the Comet client explicitly, use `createWithSigner`.
     */
    static connectWithSigner(endpoint: string | HttpEndpoint, signer: OfflineSigner, options?: SigningStargateClientOptions): Promise<SigningStargateClient>;
    /**
     * Creates an instance from a manually created Comet client.
     * Use this to use `Comet38Client` or `Tendermint37Client` instead of `Tendermint34Client`.
     */
    static createWithSigner(cometClient: CometClient, signer: OfflineSigner, options?: SigningStargateClientOptions): Promise<SigningStargateClient>;
    /**
     * Creates a client in offline mode.
     *
     * This should only be used in niche cases where you know exactly what you're doing,
     * e.g. when building an offline signing application.
     *
     * When you try to use online functionality with such a signer, an
     * exception will be raised.
     */
    static offline(signer: OfflineSigner, options?: SigningStargateClientOptions): Promise<SigningStargateClient>;
    protected constructor(cometClient: CometClient | undefined, signer: OfflineSigner, options: SigningStargateClientOptions);
    simulate(signerAddress: string, messages: readonly EncodeObject[], memo: string | undefined): Promise<number>;
    sendTokens(senderAddress: string, recipientAddress: string, amount: readonly Coin[], fee: StdFee | "auto" | number, memo?: string): Promise<DeliverTxResponse>;
    delegateTokens(delegatorAddress: string, validatorAddress: string, amount: Coin, fee: StdFee | "auto" | number, memo?: string): Promise<DeliverTxResponse>;
    undelegateTokens(delegatorAddress: string, validatorAddress: string, amount: Coin, fee: StdFee | "auto" | number, memo?: string): Promise<DeliverTxResponse>;
    withdrawRewards(delegatorAddress: string, validatorAddress: string, fee: StdFee | "auto" | number, memo?: string): Promise<DeliverTxResponse>;
    /**
     * @deprecated This API does not support setting the memo field of `MsgTransfer` (only the transaction memo).
     * We'll remove this method at some point because trying to wrap the various message types is a losing strategy.
     * Please migrate to `signAndBroadcast` with an `MsgTransferEncodeObject` created in the caller code instead.
     * @see https://github.com/cosmos/cosmjs/issues/1493
     */
    sendIbcTokens(senderAddress: string, recipientAddress: string, transferAmount: Coin, sourcePort: string, sourceChannel: string, timeoutHeight: Height | undefined, 
    /** timeout in seconds */
    timeoutTimestamp: number | undefined, fee: StdFee | "auto" | number, memo?: string): Promise<DeliverTxResponse>;
    signAndBroadcast(signerAddress: string, messages: readonly EncodeObject[], fee: StdFee | "auto" | number, memo?: string, timeoutHeight?: bigint): Promise<DeliverTxResponse>;
    /**
     * This method is useful if you want to send a transaction in broadcast,
     * without waiting for it to be placed inside a block, because for example
     * I would like to receive the hash to later track the transaction with another tool.
     * @returns Returns the hash of the transaction
     */
    signAndBroadcastSync(signerAddress: string, messages: readonly EncodeObject[], fee: StdFee | "auto" | number, memo?: string, timeoutHeight?: bigint): Promise<string>;
    /**
     * Gets account number and sequence from the API, creates a sign doc,
     * creates a single signature and assembles the signed transaction.
     *
     * The sign mode (SIGN_MODE_DIRECT or SIGN_MODE_LEGACY_AMINO_JSON) is determined by this client's signer.
     *
     * You can pass signer data (account number, sequence and chain ID) explicitly instead of querying them
     * from the chain. This is needed when signing for a multisig account, but it also allows for offline signing
     * (See the SigningStargateClient.offline constructor).
     */
    sign(signerAddress: string, messages: readonly EncodeObject[], fee: StdFee, memo: string, explicitSignerData?: SignerData, timeoutHeight?: bigint): Promise<TxRaw>;
    private signAmino;
    private signDirect;
}
