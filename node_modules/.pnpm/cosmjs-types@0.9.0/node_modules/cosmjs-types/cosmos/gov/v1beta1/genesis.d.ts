import { Deposit, Vote, Proposal, DepositParams, VotingParams, TallyParams } from "./gov";
import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmos.gov.v1beta1";
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
    /** params defines all the parameters of related to deposit. */
    depositParams: DepositParams;
    /** params defines all the parameters of related to voting. */
    votingParams: VotingParams;
    /** params defines all the parameters of related to tally. */
    tallyParams: TallyParams;
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
            option?: import("./gov").VoteOption | undefined;
            options?: {
                option?: import("./gov").VoteOption | undefined;
                weight?: string | undefined;
            }[] | undefined;
        }[] | undefined;
        proposals?: {
            proposalId?: bigint | undefined;
            content?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            status?: import("./gov").ProposalStatus | undefined;
            finalTallyResult?: {
                yes?: string | undefined;
                abstain?: string | undefined;
                no?: string | undefined;
                noWithVeto?: string | undefined;
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
            quorum?: Uint8Array | undefined;
            threshold?: Uint8Array | undefined;
            vetoThreshold?: Uint8Array | undefined;
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
            option?: import("./gov").VoteOption | undefined;
            options?: {
                option?: import("./gov").VoteOption | undefined;
                weight?: string | undefined;
            }[] | undefined;
        }[] & ({
            proposalId?: bigint | undefined;
            voter?: string | undefined;
            option?: import("./gov").VoteOption | undefined;
            options?: {
                option?: import("./gov").VoteOption | undefined;
                weight?: string | undefined;
            }[] | undefined;
        } & {
            proposalId?: bigint | undefined;
            voter?: string | undefined;
            option?: import("./gov").VoteOption | undefined;
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
        } & Record<Exclude<keyof I["votes"][number], keyof Vote>, never>)[] & Record<Exclude<keyof I["votes"], keyof {
            proposalId?: bigint | undefined;
            voter?: string | undefined;
            option?: import("./gov").VoteOption | undefined;
            options?: {
                option?: import("./gov").VoteOption | undefined;
                weight?: string | undefined;
            }[] | undefined;
        }[]>, never>) | undefined;
        proposals?: ({
            proposalId?: bigint | undefined;
            content?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            status?: import("./gov").ProposalStatus | undefined;
            finalTallyResult?: {
                yes?: string | undefined;
                abstain?: string | undefined;
                no?: string | undefined;
                noWithVeto?: string | undefined;
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
        }[] & ({
            proposalId?: bigint | undefined;
            content?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            status?: import("./gov").ProposalStatus | undefined;
            finalTallyResult?: {
                yes?: string | undefined;
                abstain?: string | undefined;
                no?: string | undefined;
                noWithVeto?: string | undefined;
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
        } & {
            proposalId?: bigint | undefined;
            content?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["proposals"][number]["content"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
            status?: import("./gov").ProposalStatus | undefined;
            finalTallyResult?: ({
                yes?: string | undefined;
                abstain?: string | undefined;
                no?: string | undefined;
                noWithVeto?: string | undefined;
            } & {
                yes?: string | undefined;
                abstain?: string | undefined;
                no?: string | undefined;
                noWithVeto?: string | undefined;
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
        } & Record<Exclude<keyof I["proposals"][number], keyof Proposal>, never>)[] & Record<Exclude<keyof I["proposals"], keyof {
            proposalId?: bigint | undefined;
            content?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            status?: import("./gov").ProposalStatus | undefined;
            finalTallyResult?: {
                yes?: string | undefined;
                abstain?: string | undefined;
                no?: string | undefined;
                noWithVeto?: string | undefined;
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
            quorum?: Uint8Array | undefined;
            threshold?: Uint8Array | undefined;
            vetoThreshold?: Uint8Array | undefined;
        } & {
            quorum?: Uint8Array | undefined;
            threshold?: Uint8Array | undefined;
            vetoThreshold?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["tallyParams"], keyof TallyParams>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof GenesisState>, never>>(object: I): GenesisState;
};
