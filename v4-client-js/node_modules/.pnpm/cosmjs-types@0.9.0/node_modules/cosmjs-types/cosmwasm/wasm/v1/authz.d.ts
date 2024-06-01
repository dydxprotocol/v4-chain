import { Any } from "../../../google/protobuf/any";
import { Coin } from "../../../cosmos/base/v1beta1/coin";
import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmwasm.wasm.v1";
/**
 * ContractExecutionAuthorization defines authorization for wasm execute.
 * Since: wasmd 0.30
 */
export interface ContractExecutionAuthorization {
    /** Grants for contract executions */
    grants: ContractGrant[];
}
/**
 * ContractMigrationAuthorization defines authorization for wasm contract
 * migration. Since: wasmd 0.30
 */
export interface ContractMigrationAuthorization {
    /** Grants for contract migrations */
    grants: ContractGrant[];
}
/**
 * ContractGrant a granted permission for a single contract
 * Since: wasmd 0.30
 */
export interface ContractGrant {
    /** Contract is the bech32 address of the smart contract */
    contract: string;
    /**
     * Limit defines execution limits that are enforced and updated when the grant
     * is applied. When the limit lapsed the grant is removed.
     */
    limit?: Any;
    /**
     * Filter define more fine-grained control on the message payload passed
     * to the contract in the operation. When no filter applies on execution, the
     * operation is prohibited.
     */
    filter?: Any;
}
/**
 * MaxCallsLimit limited number of calls to the contract. No funds transferable.
 * Since: wasmd 0.30
 */
export interface MaxCallsLimit {
    /** Remaining number that is decremented on each execution */
    remaining: bigint;
}
/**
 * MaxFundsLimit defines the maximal amounts that can be sent to the contract.
 * Since: wasmd 0.30
 */
export interface MaxFundsLimit {
    /** Amounts is the maximal amount of tokens transferable to the contract. */
    amounts: Coin[];
}
/**
 * CombinedLimit defines the maximal amounts that can be sent to a contract and
 * the maximal number of calls executable. Both need to remain >0 to be valid.
 * Since: wasmd 0.30
 */
export interface CombinedLimit {
    /** Remaining number that is decremented on each execution */
    callsRemaining: bigint;
    /** Amounts is the maximal amount of tokens transferable to the contract. */
    amounts: Coin[];
}
/**
 * AllowAllMessagesFilter is a wildcard to allow any type of contract payload
 * message.
 * Since: wasmd 0.30
 */
export interface AllowAllMessagesFilter {
}
/**
 * AcceptedMessageKeysFilter accept only the specific contract message keys in
 * the json object to be executed.
 * Since: wasmd 0.30
 */
export interface AcceptedMessageKeysFilter {
    /** Messages is the list of unique keys */
    keys: string[];
}
/**
 * AcceptedMessagesFilter accept only the specific raw contract messages to be
 * executed.
 * Since: wasmd 0.30
 */
export interface AcceptedMessagesFilter {
    /** Messages is the list of raw contract messages */
    messages: Uint8Array[];
}
export declare const ContractExecutionAuthorization: {
    typeUrl: string;
    encode(message: ContractExecutionAuthorization, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ContractExecutionAuthorization;
    fromJSON(object: any): ContractExecutionAuthorization;
    toJSON(message: ContractExecutionAuthorization): unknown;
    fromPartial<I extends {
        grants?: {
            contract?: string | undefined;
            limit?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            filter?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        }[] | undefined;
    } & {
        grants?: ({
            contract?: string | undefined;
            limit?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            filter?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        }[] & ({
            contract?: string | undefined;
            limit?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            filter?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        } & {
            contract?: string | undefined;
            limit?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["grants"][number]["limit"], keyof Any>, never>) | undefined;
            filter?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["grants"][number]["filter"], keyof Any>, never>) | undefined;
        } & Record<Exclude<keyof I["grants"][number], keyof ContractGrant>, never>)[] & Record<Exclude<keyof I["grants"], keyof {
            contract?: string | undefined;
            limit?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            filter?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "grants">, never>>(object: I): ContractExecutionAuthorization;
};
export declare const ContractMigrationAuthorization: {
    typeUrl: string;
    encode(message: ContractMigrationAuthorization, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ContractMigrationAuthorization;
    fromJSON(object: any): ContractMigrationAuthorization;
    toJSON(message: ContractMigrationAuthorization): unknown;
    fromPartial<I extends {
        grants?: {
            contract?: string | undefined;
            limit?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            filter?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        }[] | undefined;
    } & {
        grants?: ({
            contract?: string | undefined;
            limit?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            filter?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        }[] & ({
            contract?: string | undefined;
            limit?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            filter?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        } & {
            contract?: string | undefined;
            limit?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["grants"][number]["limit"], keyof Any>, never>) | undefined;
            filter?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["grants"][number]["filter"], keyof Any>, never>) | undefined;
        } & Record<Exclude<keyof I["grants"][number], keyof ContractGrant>, never>)[] & Record<Exclude<keyof I["grants"], keyof {
            contract?: string | undefined;
            limit?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            filter?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "grants">, never>>(object: I): ContractMigrationAuthorization;
};
export declare const ContractGrant: {
    typeUrl: string;
    encode(message: ContractGrant, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ContractGrant;
    fromJSON(object: any): ContractGrant;
    toJSON(message: ContractGrant): unknown;
    fromPartial<I extends {
        contract?: string | undefined;
        limit?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
        filter?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
    } & {
        contract?: string | undefined;
        limit?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["limit"], keyof Any>, never>) | undefined;
        filter?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["filter"], keyof Any>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof ContractGrant>, never>>(object: I): ContractGrant;
};
export declare const MaxCallsLimit: {
    typeUrl: string;
    encode(message: MaxCallsLimit, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MaxCallsLimit;
    fromJSON(object: any): MaxCallsLimit;
    toJSON(message: MaxCallsLimit): unknown;
    fromPartial<I extends {
        remaining?: bigint | undefined;
    } & {
        remaining?: bigint | undefined;
    } & Record<Exclude<keyof I, "remaining">, never>>(object: I): MaxCallsLimit;
};
export declare const MaxFundsLimit: {
    typeUrl: string;
    encode(message: MaxFundsLimit, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MaxFundsLimit;
    fromJSON(object: any): MaxFundsLimit;
    toJSON(message: MaxFundsLimit): unknown;
    fromPartial<I extends {
        amounts?: {
            denom?: string | undefined;
            amount?: string | undefined;
        }[] | undefined;
    } & {
        amounts?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        }[] & ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["amounts"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["amounts"], keyof {
            denom?: string | undefined;
            amount?: string | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "amounts">, never>>(object: I): MaxFundsLimit;
};
export declare const CombinedLimit: {
    typeUrl: string;
    encode(message: CombinedLimit, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): CombinedLimit;
    fromJSON(object: any): CombinedLimit;
    toJSON(message: CombinedLimit): unknown;
    fromPartial<I extends {
        callsRemaining?: bigint | undefined;
        amounts?: {
            denom?: string | undefined;
            amount?: string | undefined;
        }[] | undefined;
    } & {
        callsRemaining?: bigint | undefined;
        amounts?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        }[] & ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["amounts"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["amounts"], keyof {
            denom?: string | undefined;
            amount?: string | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof CombinedLimit>, never>>(object: I): CombinedLimit;
};
export declare const AllowAllMessagesFilter: {
    typeUrl: string;
    encode(_: AllowAllMessagesFilter, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): AllowAllMessagesFilter;
    fromJSON(_: any): AllowAllMessagesFilter;
    toJSON(_: AllowAllMessagesFilter): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): AllowAllMessagesFilter;
};
export declare const AcceptedMessageKeysFilter: {
    typeUrl: string;
    encode(message: AcceptedMessageKeysFilter, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): AcceptedMessageKeysFilter;
    fromJSON(object: any): AcceptedMessageKeysFilter;
    toJSON(message: AcceptedMessageKeysFilter): unknown;
    fromPartial<I extends {
        keys?: string[] | undefined;
    } & {
        keys?: (string[] & string[] & Record<Exclude<keyof I["keys"], keyof string[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "keys">, never>>(object: I): AcceptedMessageKeysFilter;
};
export declare const AcceptedMessagesFilter: {
    typeUrl: string;
    encode(message: AcceptedMessagesFilter, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): AcceptedMessagesFilter;
    fromJSON(object: any): AcceptedMessagesFilter;
    toJSON(message: AcceptedMessagesFilter): unknown;
    fromPartial<I extends {
        messages?: Uint8Array[] | undefined;
    } & {
        messages?: (Uint8Array[] & Uint8Array[] & Record<Exclude<keyof I["messages"], keyof Uint8Array[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "messages">, never>>(object: I): AcceptedMessagesFilter;
};
