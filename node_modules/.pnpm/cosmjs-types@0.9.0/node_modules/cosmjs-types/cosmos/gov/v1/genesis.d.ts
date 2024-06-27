import { Deposit, Vote, Proposal, DepositParams, VotingParams, TallyParams, Params } from "./gov";
import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmos.gov.v1";
/** GenesisState defines the gov module's genesis state. */
export interface GenesisState {
    /** starting_proposal_id is the ID of the starting proposal. */
    startingProposalId: bigint;
    /** deposits defines all the deposits present at genesis. */
    deposits: Deposit[];
    /** votes defines all the votes present at genesis. */
    votes: Vote[];
    /** proposals defines all the proposals present at genesis. */
    proposals: Proposal[];
    /**
     * Deprecated: Prefer to use `params` instead.
     * deposit_params defines all the paramaters of related to deposit.
     */
    /** @deprecated */
    depositParams?: DepositParams;
    /**
     * Deprecated: Prefer to use `params` instead.
     * voting_params defines all the paramaters of related to voting.
     */
    /** @deprecated */
    votingParams?: VotingParams;
    /**
     * Deprecated: Prefer to use `params` instead.
     * tally_params defines all the paramaters of related to tally.
     */
    /** @deprecated */
    tallyParams?: TallyParams;
    /**
     * params defines all the paramaters of x/gov module.
     *
     * Since: cosmos-sdk 0.47
     */
    params?: Params;
}
export declare const GenesisState: {
    typeUrl: string;
    encode(message: GenesisState, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
    fromPartial<I extends {
        startingProposalId?: bigint | undefined;
        deposits?: {
            proposalId?: bigint | undefined;
            depositor?: string | undefined;
            amount?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
        }[] | undefined;
        votes?: {
            proposalId?: bigint | undefined;
            voter?: string | undefined;
            options?: {
                option?: import("./gov").VoteOption | undefined;
                weight?: string | undefined;
            }[] | undefined;
            metadata?: string | undefined;
        }[] | undefined;
        proposals?: {
            id?: bigint | undefined;
            messages?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            }[] | undefined;
            status?: import("./gov").ProposalStatus | undefined;
            finalTallyResult?: {
                yesCount?: string | undefined;
                abstainCount?: string | undefined;
                noCount?: string | undefined;
                noWithVetoCount?: string | undefined;
            } | undefined;
            submitTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            depositEndTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            totalDeposit?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            votingStartTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            votingEndTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            metadata?: string | undefined;
            title?: string | undefined;
            summary?: string | undefined;
            proposer?: string | undefined;
        }[] | undefined;
        depositParams?: {
            minDeposit?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            maxDepositPeriod?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        } | undefined;
        votingParams?: {
            votingPeriod?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        } | undefined;
        tallyParams?: {
            quorum?: string | undefined;
            threshold?: string | undefined;
            vetoThreshold?: string | undefined;
        } | undefined;
        params?: {
            minDeposit?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            maxDepositPeriod?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            votingPeriod?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            quorum?: string | undefined;
            threshold?: string | undefined;
            vetoThreshold?: string | undefined;
            minInitialDepositRatio?: string | undefined;
            burnVoteQuorum?: boolean | undefined;
            burnProposalDepositPrevote?: boolean | undefined;
            burnVoteVeto?: boolean | undefined;
        } | undefined;
    } & {
        startingProposalId?: bigint | undefined;
        deposits?: ({
            proposalId?: bigint | undefined;
            depositor?: string | undefined;
            amount?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
        }[] & ({
            proposalId?: bigint | undefined;
            depositor?: string | undefined;
            amount?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
        } & {
            proposalId?: bigint | undefined;
            depositor?: string | undefined;
            amount?: ({
                denom?: string | undefined;
                amount?: string | undefined;
            }[] & ({
                denom?: string | undefined;
                amount?: string | undefined;
            } & {
                denom?: string | undefined;
                amount?: string | undefined;
            } & Record<Exclude<keyof I["deposits"][number]["amount"][number], keyof import("../../base/v1beta1/coin").Coin>, never>)[] & Record<Exclude<keyof I["deposits"][number]["amount"], keyof {
                denom?: string | undefined;
                amount?: string | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["deposits"][number], keyof Deposit>, never>)[] & Record<Exclude<keyof I["deposits"], keyof {
            proposalId?: bigint | undefined;
            depositor?: string | undefined;
            amount?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
        }[]>, never>) | undefined;
        votes?: ({
            proposalId?: bigint | undefined;
            voter?: string | undefined;
            options?: {
                option?: import("./gov").VoteOption | undefined;
                weight?: string | undefined;
            }[] | undefined;
            metadata?: string | undefined;
        }[] & ({
            proposalId?: bigint | undefined;
            voter?: string | undefined;
            options?: {
                option?: import("./gov").VoteOption | undefined;
                weight?: string | undefined;
            }[] | undefined;
            metadata?: string | undefined;
        } & {
            proposalId?: bigint | undefined;
            voter?: string | undefined;
            options?: ({
                option?: import("./gov").VoteOption | undefined;
                weight?: string | undefined;
            }[] & ({
                option?: import("./gov").VoteOption | undefined;
                weight?: string | undefined;
            } & {
                option?: import("./gov").VoteOption | undefined;
                weight?: string | undefined;
            } & Record<Exclude<keyof I["votes"][number]["options"][number], keyof import("./gov").WeightedVoteOption>, never>)[] & Record<Exclude<keyof I["votes"][number]["options"], keyof {
                option?: import("./gov").VoteOption | undefined;
                weight?: string | undefined;
            }[]>, never>) | undefined;
            metadata?: string | undefined;
        } & Record<Exclude<keyof I["votes"][number], keyof Vote>, never>)[] & Record<Exclude<keyof I["votes"], keyof {
            proposalId?: bigint | undefined;
            voter?: string | undefined;
            options?: {
                option?: import("./gov").VoteOption | undefined;
                weight?: string | undefined;
            }[] | undefined;
            metadata?: string | undefined;
        }[]>, never>) | undefined;
        proposals?: ({
            id?: bigint | undefined;
            messages?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            }[] | undefined;
            status?: import("./gov").ProposalStatus | undefined;
            finalTallyResult?: {
                yesCount?: string | undefined;
                abstainCount?: string | undefined;
                noCount?: string | undefined;
                noWithVetoCount?: string | undefined;
            } | undefined;
            submitTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            depositEndTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            totalDeposit?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            votingStartTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            votingEndTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            metadata?: string | undefined;
            title?: string | undefined;
            summary?: string | undefined;
            proposer?: string | undefined;
        }[] & ({
            id?: bigint | undefined;
            messages?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            }[] | undefined;
            status?: import("./gov").ProposalStatus | undefined;
            finalTallyResult?: {
                yesCount?: string | undefined;
                abstainCount?: string | undefined;
                noCount?: string | undefined;
                noWithVetoCount?: string | undefined;
            } | undefined;
            submitTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            depositEndTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            totalDeposit?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            votingStartTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            votingEndTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            metadata?: string | undefined;
            title?: string | undefined;
            summary?: string | undefined;
            proposer?: string | undefined;
        } & {
            id?: bigint | undefined;
            messages?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            }[] & ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["proposals"][number]["messages"][number], keyof import("../../../google/protobuf/any").Any>, never>)[] & Record<Exclude<keyof I["proposals"][number]["messages"], keyof {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            }[]>, never>) | undefined;
            status?: import("./gov").ProposalStatus | undefined;
            finalTallyResult?: ({
                yesCount?: string | undefined;
                abstainCount?: string | undefined;
                noCount?: string | undefined;
                noWithVetoCount?: string | undefined;
            } & {
                yesCount?: string | undefined;
                abstainCount?: string | undefined;
                noCount?: string | undefined;
                noWithVetoCount?: string | undefined;
            } & Record<Exclude<keyof I["proposals"][number]["finalTallyResult"], keyof import("./gov").TallyResult>, never>) | undefined;
            submitTime?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["proposals"][number]["submitTime"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
            depositEndTime?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["proposals"][number]["depositEndTime"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
            totalDeposit?: ({
                denom?: string | undefined;
                amount?: string | undefined;
            }[] & ({
                denom?: string | undefined;
                amount?: string | undefined;
            } & {
                denom?: string | undefined;
                amount?: string | undefined;
            } & Record<Exclude<keyof I["proposals"][number]["totalDeposit"][number], keyof import("../../base/v1beta1/coin").Coin>, never>)[] & Record<Exclude<keyof I["proposals"][number]["totalDeposit"], keyof {
                denom?: string | undefined;
                amount?: string | undefined;
            }[]>, never>) | undefined;
            votingStartTime?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["proposals"][number]["votingStartTime"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
            votingEndTime?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["proposals"][number]["votingEndTime"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
            metadata?: string | undefined;
            title?: string | undefined;
            summary?: string | undefined;
            proposer?: string | undefined;
        } & Record<Exclude<keyof I["proposals"][number], keyof Proposal>, never>)[] & Record<Exclude<keyof I["proposals"], keyof {
            id?: bigint | undefined;
            messages?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            }[] | undefined;
            status?: import("./gov").ProposalStatus | undefined;
            finalTallyResult?: {
                yesCount?: string | undefined;
                abstainCount?: string | undefined;
                noCount?: string | undefined;
                noWithVetoCount?: string | undefined;
            } | undefined;
            submitTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            depositEndTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            totalDeposit?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            votingStartTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            votingEndTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            metadata?: string | undefined;
            title?: string | undefined;
            summary?: string | undefined;
            proposer?: string | undefined;
        }[]>, never>) | undefined;
        depositParams?: ({
            minDeposit?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            maxDepositPeriod?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        } & {
            minDeposit?: ({
                denom?: string | undefined;
                amount?: string | undefined;
            }[] & ({
                denom?: string | undefined;
                amount?: string | undefined;
            } & {
                denom?: string | undefined;
                amount?: string | undefined;
            } & Record<Exclude<keyof I["depositParams"]["minDeposit"][number], keyof import("../../base/v1beta1/coin").Coin>, never>)[] & Record<Exclude<keyof I["depositParams"]["minDeposit"], keyof {
                denom?: string | undefined;
                amount?: string | undefined;
            }[]>, never>) | undefined;
            maxDepositPeriod?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["depositParams"]["maxDepositPeriod"], keyof import("../../../google/protobuf/duration").Duration>, never>) | undefined;
        } & Record<Exclude<keyof I["depositParams"], keyof DepositParams>, never>) | undefined;
        votingParams?: ({
            votingPeriod?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        } & {
            votingPeriod?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["votingParams"]["votingPeriod"], keyof import("../../../google/protobuf/duration").Duration>, never>) | undefined;
        } & Record<Exclude<keyof I["votingParams"], "votingPeriod">, never>) | undefined;
        tallyParams?: ({
            quorum?: string | undefined;
            threshold?: string | undefined;
            vetoThreshold?: string | undefined;
        } & {
            quorum?: string | undefined;
            threshold?: string | undefined;
            vetoThreshold?: string | undefined;
        } & Record<Exclude<keyof I["tallyParams"], keyof TallyParams>, never>) | undefined;
        params?: ({
            minDeposit?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            maxDepositPeriod?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            votingPeriod?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            quorum?: string | undefined;
            threshold?: string | undefined;
            vetoThreshold?: string | undefined;
            minInitialDepositRatio?: string | undefined;
            burnVoteQuorum?: boolean | undefined;
            burnProposalDepositPrevote?: boolean | undefined;
            burnVoteVeto?: boolean | undefined;
        } & {
            minDeposit?: ({
                denom?: string | undefined;
                amount?: string | undefined;
            }[] & ({
                denom?: string | undefined;
                amount?: string | undefined;
            } & {
                denom?: string | undefined;
                amount?: string | undefined;
            } & Record<Exclude<keyof I["params"]["minDeposit"][number], keyof import("../../base/v1beta1/coin").Coin>, never>)[] & Record<Exclude<keyof I["params"]["minDeposit"], keyof {
                denom?: string | undefined;
                amount?: string | undefined;
            }[]>, never>) | undefined;
            maxDepositPeriod?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["params"]["maxDepositPeriod"], keyof import("../../../google/protobuf/duration").Duration>, never>) | undefined;
            votingPeriod?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["params"]["votingPeriod"], keyof import("../../../google/protobuf/duration").Duration>, never>) | undefined;
            quorum?: string | undefined;
            threshold?: string | undefined;
            vetoThreshold?: string | undefined;
            minInitialDepositRatio?: string | undefined;
            burnVoteQuorum?: boolean | undefined;
            burnProposalDepositPrevote?: boolean | undefined;
            burnVoteVeto?: boolean | undefined;
        } & Record<Exclude<keyof I["params"], keyof Params>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof GenesisState>, never>>(object: I): GenesisState;
};
