import { DecCoin, Coin } from "../../base/v1beta1/coin";
import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmos.distribution.v1beta1";
/** Params defines the set of params for the distribution module. */
export interface Params {
    communityTax: string;
    /**
     * Deprecated: The base_proposer_reward field is deprecated and is no longer used
     * in the x/distribution module's reward mechanism.
     */
    /** @deprecated */
    baseProposerReward: string;
    /**
     * Deprecated: The bonus_proposer_reward field is deprecated and is no longer used
     * in the x/distribution module's reward mechanism.
     */
    /** @deprecated */
    bonusProposerReward: string;
    withdrawAddrEnabled: boolean;
}
/**
 * ValidatorHistoricalRewards represents historical rewards for a validator.
 * Height is implicit within the store key.
 * Cumulative reward ratio is the sum from the zeroeth period
 * until this period of rewards / tokens, per the spec.
 * The reference count indicates the number of objects
 * which might need to reference this historical entry at any point.
 * ReferenceCount =
 *    number of outstanding delegations which ended the associated period (and
 *    might need to read that record)
 *  + number of slashes which ended the associated period (and might need to
 *  read that record)
 *  + one per validator for the zeroeth period, set on initialization
 */
export interface ValidatorHistoricalRewards {
    cumulativeRewardRatio: DecCoin[];
    referenceCount: number;
}
/**
 * ValidatorCurrentRewards represents current rewards and current
 * period for a validator kept as a running counter and incremented
 * each block as long as the validator's tokens remain constant.
 */
export interface ValidatorCurrentRewards {
    rewards: DecCoin[];
    period: bigint;
}
/**
 * ValidatorAccumulatedCommission represents accumulated commission
 * for a validator kept as a running counter, can be withdrawn at any time.
 */
export interface ValidatorAccumulatedCommission {
    commission: DecCoin[];
}
/**
 * ValidatorOutstandingRewards represents outstanding (un-withdrawn) rewards
 * for a validator inexpensive to track, allows simple sanity checks.
 */
export interface ValidatorOutstandingRewards {
    rewards: DecCoin[];
}
/**
 * ValidatorSlashEvent represents a validator slash event.
 * Height is implicit within the store key.
 * This is needed to calculate appropriate amount of staking tokens
 * for delegations which are withdrawn after a slash has occurred.
 */
export interface ValidatorSlashEvent {
    validatorPeriod: bigint;
    fraction: string;
}
/** ValidatorSlashEvents is a collection of ValidatorSlashEvent messages. */
export interface ValidatorSlashEvents {
    validatorSlashEvents: ValidatorSlashEvent[];
}
/** FeePool is the global fee pool for distribution. */
export interface FeePool {
    communityPool: DecCoin[];
}
/**
 * CommunityPoolSpendProposal details a proposal for use of community funds,
 * together with how many coins are proposed to be spent, and to which
 * recipient account.
 *
 * Deprecated: Do not use. As of the Cosmos SDK release v0.47.x, there is no
 * longer a need for an explicit CommunityPoolSpendProposal. To spend community
 * pool funds, a simple MsgCommunityPoolSpend can be invoked from the x/gov
 * module via a v1 governance proposal.
 */
/** @deprecated */
export interface CommunityPoolSpendProposal {
    title: string;
    description: string;
    recipient: string;
    amount: Coin[];
}
/**
 * DelegatorStartingInfo represents the starting info for a delegator reward
 * period. It tracks the previous validator period, the delegation's amount of
 * staking token, and the creation height (to check later on if any slashes have
 * occurred). NOTE: Even though validators are slashed to whole staking tokens,
 * the delegators within the validator may be left with less than a full token,
 * thus sdk.Dec is used.
 */
export interface DelegatorStartingInfo {
    previousPeriod: bigint;
    stake: string;
    height: bigint;
}
/**
 * DelegationDelegatorReward represents the properties
 * of a delegator's delegation reward.
 */
export interface DelegationDelegatorReward {
    validatorAddress: string;
    reward: DecCoin[];
}
/**
 * CommunityPoolSpendProposalWithDeposit defines a CommunityPoolSpendProposal
 * with a deposit
 */
export interface CommunityPoolSpendProposalWithDeposit {
    title: string;
    description: string;
    recipient: string;
    amount: string;
    deposit: string;
}
export declare const Params: {
    typeUrl: string;
    encode(message: Params, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Params;
    fromJSON(object: any): Params;
    toJSON(message: Params): unknown;
    fromPartial<I extends {
        communityTax?: string | undefined;
        baseProposerReward?: string | undefined;
        bonusProposerReward?: string | undefined;
        withdrawAddrEnabled?: boolean | undefined;
    } & {
        communityTax?: string | undefined;
        baseProposerReward?: string | undefined;
        bonusProposerReward?: string | undefined;
        withdrawAddrEnabled?: boolean | undefined;
    } & Record<Exclude<keyof I, keyof Params>, never>>(object: I): Params;
};
export declare const ValidatorHistoricalRewards: {
    typeUrl: string;
    encode(message: ValidatorHistoricalRewards, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ValidatorHistoricalRewards;
    fromJSON(object: any): ValidatorHistoricalRewards;
    toJSON(message: ValidatorHistoricalRewards): unknown;
    fromPartial<I extends {
        cumulativeRewardRatio?: {
            denom?: string | undefined;
            amount?: string | undefined;
        }[] | undefined;
        referenceCount?: number | undefined;
    } & {
        cumulativeRewardRatio?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        }[] & ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["cumulativeRewardRatio"][number], keyof DecCoin>, never>)[] & Record<Exclude<keyof I["cumulativeRewardRatio"], keyof {
            denom?: string | undefined;
            amount?: string | undefined;
        }[]>, never>) | undefined;
        referenceCount?: number | undefined;
    } & Record<Exclude<keyof I, keyof ValidatorHistoricalRewards>, never>>(object: I): ValidatorHistoricalRewards;
};
export declare const ValidatorCurrentRewards: {
    typeUrl: string;
    encode(message: ValidatorCurrentRewards, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ValidatorCurrentRewards;
    fromJSON(object: any): ValidatorCurrentRewards;
    toJSON(message: ValidatorCurrentRewards): unknown;
    fromPartial<I extends {
        rewards?: {
            denom?: string | undefined;
            amount?: string | undefined;
        }[] | undefined;
        period?: bigint | undefined;
    } & {
        rewards?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        }[] & ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["rewards"][number], keyof DecCoin>, never>)[] & Record<Exclude<keyof I["rewards"], keyof {
            denom?: string | undefined;
            amount?: string | undefined;
        }[]>, never>) | undefined;
        period?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof ValidatorCurrentRewards>, never>>(object: I): ValidatorCurrentRewards;
};
export declare const ValidatorAccumulatedCommission: {
    typeUrl: string;
    encode(message: ValidatorAccumulatedCommission, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ValidatorAccumulatedCommission;
    fromJSON(object: any): ValidatorAccumulatedCommission;
    toJSON(message: ValidatorAccumulatedCommission): unknown;
    fromPartial<I extends {
        commission?: {
            denom?: string | undefined;
            amount?: string | undefined;
        }[] | undefined;
    } & {
        commission?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        }[] & ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["commission"][number], keyof DecCoin>, never>)[] & Record<Exclude<keyof I["commission"], keyof {
            denom?: string | undefined;
            amount?: string | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "commission">, never>>(object: I): ValidatorAccumulatedCommission;
};
export declare const ValidatorOutstandingRewards: {
    typeUrl: string;
    encode(message: ValidatorOutstandingRewards, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ValidatorOutstandingRewards;
    fromJSON(object: any): ValidatorOutstandingRewards;
    toJSON(message: ValidatorOutstandingRewards): unknown;
    fromPartial<I extends {
        rewards?: {
            denom?: string | undefined;
            amount?: string | undefined;
        }[] | undefined;
    } & {
        rewards?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        }[] & ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["rewards"][number], keyof DecCoin>, never>)[] & Record<Exclude<keyof I["rewards"], keyof {
            denom?: string | undefined;
            amount?: string | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "rewards">, never>>(object: I): ValidatorOutstandingRewards;
};
export declare const ValidatorSlashEvent: {
    typeUrl: string;
    encode(message: ValidatorSlashEvent, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ValidatorSlashEvent;
    fromJSON(object: any): ValidatorSlashEvent;
    toJSON(message: ValidatorSlashEvent): unknown;
    fromPartial<I extends {
        validatorPeriod?: bigint | undefined;
        fraction?: string | undefined;
    } & {
        validatorPeriod?: bigint | undefined;
        fraction?: string | undefined;
    } & Record<Exclude<keyof I, keyof ValidatorSlashEvent>, never>>(object: I): ValidatorSlashEvent;
};
export declare const ValidatorSlashEvents: {
    typeUrl: string;
    encode(message: ValidatorSlashEvents, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ValidatorSlashEvents;
    fromJSON(object: any): ValidatorSlashEvents;
    toJSON(message: ValidatorSlashEvents): unknown;
    fromPartial<I extends {
        validatorSlashEvents?: {
            validatorPeriod?: bigint | undefined;
            fraction?: string | undefined;
        }[] | undefined;
    } & {
        validatorSlashEvents?: ({
            validatorPeriod?: bigint | undefined;
            fraction?: string | undefined;
        }[] & ({
            validatorPeriod?: bigint | undefined;
            fraction?: string | undefined;
        } & {
            validatorPeriod?: bigint | undefined;
            fraction?: string | undefined;
        } & Record<Exclude<keyof I["validatorSlashEvents"][number], keyof ValidatorSlashEvent>, never>)[] & Record<Exclude<keyof I["validatorSlashEvents"], keyof {
            validatorPeriod?: bigint | undefined;
            fraction?: string | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "validatorSlashEvents">, never>>(object: I): ValidatorSlashEvents;
};
export declare const FeePool: {
    typeUrl: string;
    encode(message: FeePool, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): FeePool;
    fromJSON(object: any): FeePool;
    toJSON(message: FeePool): unknown;
    fromPartial<I extends {
        communityPool?: {
            denom?: string | undefined;
            amount?: string | undefined;
        }[] | undefined;
    } & {
        communityPool?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        }[] & ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["communityPool"][number], keyof DecCoin>, never>)[] & Record<Exclude<keyof I["communityPool"], keyof {
            denom?: string | undefined;
            amount?: string | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "communityPool">, never>>(object: I): FeePool;
};
export declare const CommunityPoolSpendProposal: {
    typeUrl: string;
    encode(message: CommunityPoolSpendProposal, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): CommunityPoolSpendProposal;
    fromJSON(object: any): CommunityPoolSpendProposal;
    toJSON(message: CommunityPoolSpendProposal): unknown;
    fromPartial<I extends {
        title?: string | undefined;
        description?: string | undefined;
        recipient?: string | undefined;
        amount?: {
            denom?: string | undefined;
            amount?: string | undefined;
        }[] | undefined;
    } & {
        title?: string | undefined;
        description?: string | undefined;
        recipient?: string | undefined;
        amount?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        }[] & ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["amount"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["amount"], keyof {
            denom?: string | undefined;
            amount?: string | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof CommunityPoolSpendProposal>, never>>(object: I): CommunityPoolSpendProposal;
};
export declare const DelegatorStartingInfo: {
    typeUrl: string;
    encode(message: DelegatorStartingInfo, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): DelegatorStartingInfo;
    fromJSON(object: any): DelegatorStartingInfo;
    toJSON(message: DelegatorStartingInfo): unknown;
    fromPartial<I extends {
        previousPeriod?: bigint | undefined;
        stake?: string | undefined;
        height?: bigint | undefined;
    } & {
        previousPeriod?: bigint | undefined;
        stake?: string | undefined;
        height?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof DelegatorStartingInfo>, never>>(object: I): DelegatorStartingInfo;
};
export declare const DelegationDelegatorReward: {
    typeUrl: string;
    encode(message: DelegationDelegatorReward, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): DelegationDelegatorReward;
    fromJSON(object: any): DelegationDelegatorReward;
    toJSON(message: DelegationDelegatorReward): unknown;
    fromPartial<I extends {
        validatorAddress?: string | undefined;
        reward?: {
            denom?: string | undefined;
            amount?: string | undefined;
        }[] | undefined;
    } & {
        validatorAddress?: string | undefined;
        reward?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        }[] & ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["reward"][number], keyof DecCoin>, never>)[] & Record<Exclude<keyof I["reward"], keyof {
            denom?: string | undefined;
            amount?: string | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof DelegationDelegatorReward>, never>>(object: I): DelegationDelegatorReward;
};
export declare const CommunityPoolSpendProposalWithDeposit: {
    typeUrl: string;
    encode(message: CommunityPoolSpendProposalWithDeposit, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): CommunityPoolSpendProposalWithDeposit;
    fromJSON(object: any): CommunityPoolSpendProposalWithDeposit;
    toJSON(message: CommunityPoolSpendProposalWithDeposit): unknown;
    fromPartial<I extends {
        title?: string | undefined;
        description?: string | undefined;
        recipient?: string | undefined;
        amount?: string | undefined;
        deposit?: string | undefined;
    } & {
        title?: string | undefined;
        description?: string | undefined;
        recipient?: string | undefined;
        amount?: string | undefined;
        deposit?: string | undefined;
    } & Record<Exclude<keyof I, keyof CommunityPoolSpendProposalWithDeposit>, never>>(object: I): CommunityPoolSpendProposalWithDeposit;
};
