import { AminoMsg, Coin } from "@cosmjs/amino";
import { AminoConverters } from "../../aminotypes";
/** Supports submitting arbitrary proposal content. */
export interface AminoMsgSubmitProposal extends AminoMsg {
    readonly type: "cosmos-sdk/MsgSubmitProposal";
    readonly value: {
        /**
         * A proposal structure, e.g.
         *
         * ```
         * {
         *   type: 'cosmos-sdk/TextProposal',
         *   value: {
         *     description: 'This proposal proposes to test whether this proposal passes',
         *     title: 'Test Proposal'
         *   }
         * }
         * ```
         */
        readonly content: {
            readonly type: string;
            readonly value: any;
        };
        readonly initial_deposit: readonly Coin[];
        /** Bech32 account address */
        readonly proposer: string;
    };
}
export declare function isAminoMsgSubmitProposal(msg: AminoMsg): msg is AminoMsgSubmitProposal;
/** Casts a vote */
export interface AminoMsgVote extends AminoMsg {
    readonly type: "cosmos-sdk/MsgVote";
    readonly value: {
        readonly proposal_id: string;
        /** Bech32 account address */
        readonly voter: string;
        /**
         * VoteOption as integer from 0 to 4 🤷‍
         *
         * @see https://github.com/cosmos/cosmos-sdk/blob/v0.42.9/x/gov/types/gov.pb.go#L38-L49
         */
        readonly option: number;
    };
}
export declare function isAminoMsgVote(msg: AminoMsg): msg is AminoMsgVote;
/**
 * @see https://github.com/cosmos/cosmos-sdk/blob/v0.44.5/x/gov/types/tx.pb.go#L196-L203
 * @see https://github.com/cosmos/cosmos-sdk/blob/v0.44.5/x/gov/types/gov.pb.go#L124-L130
 */
export interface AminoMsgVoteWeighted extends AminoMsg {
    readonly type: "cosmos-sdk/MsgVoteWeighted";
    readonly value: {
        readonly proposal_id: string;
        /** Bech32 account address */
        readonly voter: string;
        readonly options: Array<{
            /**
             * VoteOption as integer from 0 to 4 🤷‍
             *
             * @see https://github.com/cosmos/cosmos-sdk/blob/v0.44.5/x/gov/types/gov.pb.go#L35-L49
             */
            readonly option: number;
            readonly weight: string;
        }>;
    };
}
export declare function isAminoMsgVoteWeighted(msg: AminoMsg): msg is AminoMsgVoteWeighted;
/** Submits a deposit to an existing proposal */
export interface AminoMsgDeposit extends AminoMsg {
    readonly type: "cosmos-sdk/MsgDeposit";
    readonly value: {
        readonly proposal_id: string;
        /** Bech32 account address */
        readonly depositor: string;
        readonly amount: readonly Coin[];
    };
}
export declare function isAminoMsgDeposit(msg: AminoMsg): msg is AminoMsgDeposit;
export declare function createGovAminoConverters(): AminoConverters;
