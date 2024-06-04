import { EncodeObject } from '@cosmjs/proto-signing';
import { SigningStargateClient, StdFee } from '@cosmjs/stargate';
import { TransactionOptions } from '../types';
export declare class TransactionSigner {
    readonly address: string;
    readonly stargateSigningClient: SigningStargateClient;
    constructor(address: string, stargateSigningClient: SigningStargateClient);
    /**
     * @description Get the encoded signed transaction or the promise is rejected if
     * no fee can be set for the transaction.
     *
     * @throws UserError if the fee is undefined.
     * @returns The signed and encoded transaction.
     */
    signTransaction(messages: EncodeObject[], transactionOptions: TransactionOptions, fee?: StdFee, memo?: string): Promise<Uint8Array>;
}
