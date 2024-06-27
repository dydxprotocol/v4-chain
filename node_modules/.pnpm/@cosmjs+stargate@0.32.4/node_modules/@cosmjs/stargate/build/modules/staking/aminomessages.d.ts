import { AminoMsg, Coin, Pubkey } from "@cosmjs/amino";
import { AminoConverter } from "../..";
/** The initial commission rates to be used for creating a validator */
interface CommissionRates {
    readonly rate: string;
    readonly max_rate: string;
    readonly max_change_rate: string;
}
/** A validator description. */
interface Description {
    readonly moniker: string;
    readonly identity: string;
    readonly website: string;
    readonly security_contact: string;
    readonly details: string;
}
export declare function protoDecimalToJson(decimal: string): string;
/** Creates a new validator. */
export interface AminoMsgCreateValidator extends AminoMsg {
    readonly type: "cosmos-sdk/MsgCreateValidator";
    readonly value: {
        readonly description: Description;
        readonly commission: CommissionRates;
        readonly min_self_delegation: string;
        /** Bech32 encoded delegator address */
        readonly delegator_address: string;
        /** Bech32 encoded validator address */
        readonly validator_address: string;
        /** Public key */
        readonly pubkey: Pubkey;
        readonly value: Coin;
    };
}
export declare function isAminoMsgCreateValidator(msg: AminoMsg): msg is AminoMsgCreateValidator;
/** Edits an existing validator. */
export interface AminoMsgEditValidator extends AminoMsg {
    readonly type: "cosmos-sdk/MsgEditValidator";
    readonly value: {
        readonly description: Description;
        /** Bech32 encoded validator address */
        readonly validator_address: string;
        /**
         * The new value for the comission rate.
         *
         * An empty string in the protobuf document means "do not change".
         * In Amino JSON this empty string becomes undefined (omitempty)
         */
        readonly commission_rate: string | undefined;
        /**
         * The new value for the comission rate.
         *
         * An empty string in the protobuf document means "do not change".
         * In Amino JSON this empty string becomes undefined (omitempty)
         */
        readonly min_self_delegation: string | undefined;
    };
}
export declare function isAminoMsgEditValidator(msg: AminoMsg): msg is AminoMsgEditValidator;
/**
 * Performs a delegation from a delegate to a validator.
 *
 * @see https://docs.cosmos.network/master/modules/staking/03_messages.html#msgdelegate
 */
export interface AminoMsgDelegate extends AminoMsg {
    readonly type: "cosmos-sdk/MsgDelegate";
    readonly value: {
        /** Bech32 encoded delegator address */
        readonly delegator_address: string;
        /** Bech32 encoded validator address */
        readonly validator_address: string;
        readonly amount: Coin;
    };
}
export declare function isAminoMsgDelegate(msg: AminoMsg): msg is AminoMsgDelegate;
/** Performs a redelegation from a delegate and source validator to a destination validator */
export interface AminoMsgBeginRedelegate extends AminoMsg {
    readonly type: "cosmos-sdk/MsgBeginRedelegate";
    readonly value: {
        /** Bech32 encoded delegator address */
        readonly delegator_address: string;
        /** Bech32 encoded source validator address */
        readonly validator_src_address: string;
        /** Bech32 encoded destination validator address */
        readonly validator_dst_address: string;
        readonly amount: Coin;
    };
}
export declare function isAminoMsgBeginRedelegate(msg: AminoMsg): msg is AminoMsgBeginRedelegate;
/** Performs an undelegation from a delegate and a validator */
export interface AminoMsgUndelegate extends AminoMsg {
    readonly type: "cosmos-sdk/MsgUndelegate";
    readonly value: {
        /** Bech32 encoded delegator address */
        readonly delegator_address: string;
        /** Bech32 encoded validator address */
        readonly validator_address: string;
        readonly amount: Coin;
    };
}
export declare function isAminoMsgUndelegate(msg: AminoMsg): msg is AminoMsgUndelegate;
export interface AminoMsgCancelUnbondingDelegation extends AminoMsg {
    readonly type: "cosmos-sdk/MsgCancelUnbondingDelegation";
    readonly value: {
        readonly delegator_address: string;
        readonly validator_address: string;
        readonly amount: Coin;
        readonly creation_height: string;
    };
}
export declare function isAminoMsgCancelUnbondingDelegation(msg: AminoMsg): msg is AminoMsgCancelUnbondingDelegation;
export declare function createStakingAminoConverters(): Record<string, AminoConverter>;
export {};
