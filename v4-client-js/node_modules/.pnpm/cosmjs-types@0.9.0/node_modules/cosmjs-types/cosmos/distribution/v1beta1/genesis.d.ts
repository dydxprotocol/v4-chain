import { DecCoin } from "../../base/v1beta1/coin";
import { ValidatorAccumulatedCommission, ValidatorHistoricalRewards, ValidatorCurrentRewards, DelegatorStartingInfo, ValidatorSlashEvent, Params, FeePool } from "./distribution";
import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmos.distribution.v1beta1";
/**
 * DelegatorWithdrawInfo is the address for where distributions rewards are
 * withdrawn to by default this struct is only used at genesis to feed in
 * default withdraw addresses.
 */
export interface DelegatorWithdrawInfo {
    /** delegator_address is the address of the delegator. */
    delegatorAddress: string;
    /** withdraw_address is the address to withdraw the delegation rewards to. */
    withdrawAddress: string;
}
/** ValidatorOutstandingRewardsRecord is used for import/export via genesis json. */
export interface ValidatorOutstandingRewardsRecord {
    /** validator_address is the address of the validator. */
    validatorAddress: string;
    /** outstanding_rewards represents the outstanding rewards of a validator. */
    outstandingRewards: DecCoin[];
}
/**
 * ValidatorAccumulatedCommissionRecord is used for import / export via genesis
 * json.
 */
export interface ValidatorAccumulatedCommissionRecord {
    /** validator_address is the address of the validator. */
    validatorAddress: string;
    /** accumulated is the accumulated commission of a validator. */
    accumulated: ValidatorAccumulatedCommission;
}
/**
 * ValidatorHistoricalRewardsRecord is used for import / export via genesis
 * json.
 */
export interface ValidatorHistoricalRewardsRecord {
    /** validator_address is the address of the validator. */
    validatorAddress: string;
    /** period defines the period the historical rewards apply to. */
    period: bigint;
    /** rewards defines the historical rewards of a validator. */
    rewards: ValidatorHistoricalRewards;
}
/** ValidatorCurrentRewardsRecord is used for import / export via genesis json. */
export interface ValidatorCurrentRewardsRecord {
    /** validator_address is the address of the validator. */
    validatorAddress: string;
    /** rewards defines the current rewards of a validator. */
    rewards: ValidatorCurrentRewards;
}
/** DelegatorStartingInfoRecord used for import / export via genesis json. */
export interface DelegatorStartingInfoRecord {
    /** delegator_address is the address of the delegator. */
    delegatorAddress: string;
    /** validator_address is the address of the validator. */
    validatorAddress: string;
    /** starting_info defines the starting info of a delegator. */
    startingInfo: DelegatorStartingInfo;
}
/** ValidatorSlashEventRecord is used for import / export via genesis json. */
export interface ValidatorSlashEventRecord {
    /** validator_address is the address of the validator. */
    validatorAddress: string;
    /** height defines the block height at which the slash event occurred. */
    height: bigint;
    /** period is the period of the slash event. */
    period: bigint;
    /** validator_slash_event describes the slash event. */
    validatorSlashEvent: ValidatorSlashEvent;
}
/** GenesisState defines the distribution module's genesis state. */
export interface GenesisState {
    /** params defines all the parameters of the module. */
    params: Params;
    /** fee_pool defines the fee pool at genesis. */
    feePool: FeePool;
    /** fee_pool defines the delegator withdraw infos at genesis. */
    delegatorWithdrawInfos: DelegatorWithdrawInfo[];
    /** fee_pool defines the previous proposer at genesis. */
    previousProposer: string;
    /** fee_pool defines the outstanding rewards of all validators at genesis. */
    outstandingRewards: ValidatorOutstandingRewardsRecord[];
    /** fee_pool defines the accumulated commissions of all validators at genesis. */
    validatorAccumulatedCommissions: ValidatorAccumulatedCommissionRecord[];
    /** fee_pool defines the historical rewards of all validators at genesis. */
    validatorHistoricalRewards: ValidatorHistoricalRewardsRecord[];
    /** fee_pool defines the current rewards of all validators at genesis. */
    validatorCurrentRewards: ValidatorCurrentRewardsRecord[];
    /** fee_pool defines the delegator starting infos at genesis. */
    delegatorStartingInfos: DelegatorStartingInfoRecord[];
    /** fee_pool defines the validator slash events at genesis. */
    validatorSlashEvents: ValidatorSlashEventRecord[];
}
export declare const DelegatorWithdrawInfo: {
    typeUrl: string;
    encode(message: DelegatorWithdrawInfo, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): DelegatorWithdrawInfo;
    fromJSON(object: any): DelegatorWithdrawInfo;
    toJSON(message: DelegatorWithdrawInfo): unknown;
    fromPartial<I extends {
        delegatorAddress?: string | undefined;
        withdrawAddress?: string | undefined;
    } & {
        delegatorAddress?: string | undefined;
        withdrawAddress?: string | undefined;
    } & Record<Exclude<keyof I, keyof DelegatorWithdrawInfo>, never>>(object: I): DelegatorWithdrawInfo;
};
export declare const ValidatorOutstandingRewardsRecord: {
    typeUrl: string;
    encode(message: ValidatorOutstandingRewardsRecord, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ValidatorOutstandingRewardsRecord;
    fromJSON(object: any): ValidatorOutstandingRewardsRecord;
    toJSON(message: ValidatorOutstandingRewardsRecord): unknown;
    fromPartial<I extends {
        validatorAddress?: string | undefined;
        outstandingRewards?: {
            denom?: string | undefined;
            amount?: string | undefined;
        }[] | undefined;
    } & {
        validatorAddress?: string | undefined;
        outstandingRewards?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        }[] & ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["outstandingRewards"][number], keyof DecCoin>, never>)[] & Record<Exclude<keyof I["outstandingRewards"], keyof {
            denom?: string | undefined;
            amount?: string | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof ValidatorOutstandingRewardsRecord>, never>>(object: I): ValidatorOutstandingRewardsRecord;
};
export declare const ValidatorAccumulatedCommissionRecord: {
    typeUrl: string;
    encode(message: ValidatorAccumulatedCommissionRecord, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ValidatorAccumulatedCommissionRecord;
    fromJSON(object: any): ValidatorAccumulatedCommissionRecord;
    toJSON(message: ValidatorAccumulatedCommissionRecord): unknown;
    fromPartial<I extends {
        validatorAddress?: string | undefined;
        accumulated?: {
            commission?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
        } | undefined;
    } & {
        validatorAddress?: string | undefined;
        accumulated?: ({
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
            } & Record<Exclude<keyof I["accumulated"]["commission"][number], keyof DecCoin>, never>)[] & Record<Exclude<keyof I["accumulated"]["commission"], keyof {
                denom?: string | undefined;
                amount?: string | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["accumulated"], "commission">, never>) | undefined;
    } & Record<Exclude<keyof I, keyof ValidatorAccumulatedCommissionRecord>, never>>(object: I): ValidatorAccumulatedCommissionRecord;
};
export declare const ValidatorHistoricalRewardsRecord: {
    typeUrl: string;
    encode(message: ValidatorHistoricalRewardsRecord, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ValidatorHistoricalRewardsRecord;
    fromJSON(object: any): ValidatorHistoricalRewardsRecord;
    toJSON(message: ValidatorHistoricalRewardsRecord): unknown;
    fromPartial<I extends {
        validatorAddress?: string | undefined;
        period?: bigint | undefined;
        rewards?: {
            cumulativeRewardRatio?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            referenceCount?: number | undefined;
        } | undefined;
    } & {
        validatorAddress?: string | undefined;
        period?: bigint | undefined;
        rewards?: ({
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
            } & Record<Exclude<keyof I["rewards"]["cumulativeRewardRatio"][number], keyof DecCoin>, never>)[] & Record<Exclude<keyof I["rewards"]["cumulativeRewardRatio"], keyof {
                denom?: string | undefined;
                amount?: string | undefined;
            }[]>, never>) | undefined;
            referenceCount?: number | undefined;
        } & Record<Exclude<keyof I["rewards"], keyof ValidatorHistoricalRewards>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof ValidatorHistoricalRewardsRecord>, never>>(object: I): ValidatorHistoricalRewardsRecord;
};
export declare const ValidatorCurrentRewardsRecord: {
    typeUrl: string;
    encode(message: ValidatorCurrentRewardsRecord, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ValidatorCurrentRewardsRecord;
    fromJSON(object: any): ValidatorCurrentRewardsRecord;
    toJSON(message: ValidatorCurrentRewardsRecord): unknown;
    fromPartial<I extends {
        validatorAddress?: string | undefined;
        rewards?: {
            rewards?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            period?: bigint | undefined;
        } | undefined;
    } & {
        validatorAddress?: string | undefined;
        rewards?: ({
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
            } & Record<Exclude<keyof I["rewards"]["rewards"][number], keyof DecCoin>, never>)[] & Record<Exclude<keyof I["rewards"]["rewards"], keyof {
                denom?: string | undefined;
                amount?: string | undefined;
            }[]>, never>) | undefined;
            period?: bigint | undefined;
        } & Record<Exclude<keyof I["rewards"], keyof ValidatorCurrentRewards>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof ValidatorCurrentRewardsRecord>, never>>(object: I): ValidatorCurrentRewardsRecord;
};
export declare const DelegatorStartingInfoRecord: {
    typeUrl: string;
    encode(message: DelegatorStartingInfoRecord, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): DelegatorStartingInfoRecord;
    fromJSON(object: any): DelegatorStartingInfoRecord;
    toJSON(message: DelegatorStartingInfoRecord): unknown;
    fromPartial<I extends {
        delegatorAddress?: string | undefined;
        validatorAddress?: string | undefined;
        startingInfo?: {
            previousPeriod?: bigint | undefined;
            stake?: string | undefined;
            height?: bigint | undefined;
        } | undefined;
    } & {
        delegatorAddress?: string | undefined;
        validatorAddress?: string | undefined;
        startingInfo?: ({
            previousPeriod?: bigint | undefined;
            stake?: string | undefined;
            height?: bigint | undefined;
        } & {
            previousPeriod?: bigint | undefined;
            stake?: string | undefined;
            height?: bigint | undefined;
        } & Record<Exclude<keyof I["startingInfo"], keyof DelegatorStartingInfo>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof DelegatorStartingInfoRecord>, never>>(object: I): DelegatorStartingInfoRecord;
};
export declare const ValidatorSlashEventRecord: {
    typeUrl: string;
    encode(message: ValidatorSlashEventRecord, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ValidatorSlashEventRecord;
    fromJSON(object: any): ValidatorSlashEventRecord;
    toJSON(message: ValidatorSlashEventRecord): unknown;
    fromPartial<I extends {
        validatorAddress?: string | undefined;
        height?: bigint | undefined;
        period?: bigint | undefined;
        validatorSlashEvent?: {
            validatorPeriod?: bigint | undefined;
            fraction?: string | undefined;
        } | undefined;
    } & {
        validatorAddress?: string | undefined;
        height?: bigint | undefined;
        period?: bigint | undefined;
        validatorSlashEvent?: ({
            validatorPeriod?: bigint | undefined;
            fraction?: string | undefined;
        } & {
            validatorPeriod?: bigint | undefined;
            fraction?: string | undefined;
        } & Record<Exclude<keyof I["validatorSlashEvent"], keyof ValidatorSlashEvent>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof ValidatorSlashEventRecord>, never>>(object: I): ValidatorSlashEventRecord;
};
export declare const GenesisState: {
    typeUrl: string;
    encode(message: GenesisState, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
    fromPartial<I extends {
        params?: {
            communityTax?: string | undefined;
            baseProposerReward?: string | undefined;
            bonusProposerReward?: string | undefined;
            withdrawAddrEnabled?: boolean | undefined;
        } | undefined;
        feePool?: {
            communityPool?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
        } | undefined;
        delegatorWithdrawInfos?: {
            delegatorAddress?: string | undefined;
            withdrawAddress?: string | undefined;
        }[] | undefined;
        previousProposer?: string | undefined;
        outstandingRewards?: {
            validatorAddress?: string | undefined;
            outstandingRewards?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
        }[] | undefined;
        validatorAccumulatedCommissions?: {
            validatorAddress?: string | undefined;
            accumulated?: {
                commission?: {
                    denom?: string | undefined;
                    amount?: string | undefined;
                }[] | undefined;
            } | undefined;
        }[] | undefined;
        validatorHistoricalRewards?: {
            validatorAddress?: string | undefined;
            period?: bigint | undefined;
            rewards?: {
                cumulativeRewardRatio?: {
                    denom?: string | undefined;
                    amount?: string | undefined;
                }[] | undefined;
                referenceCount?: number | undefined;
            } | undefined;
        }[] | undefined;
        validatorCurrentRewards?: {
            validatorAddress?: string | undefined;
            rewards?: {
                rewards?: {
                    denom?: string | undefined;
                    amount?: string | undefined;
                }[] | undefined;
                period?: bigint | undefined;
            } | undefined;
        }[] | undefined;
        delegatorStartingInfos?: {
            delegatorAddress?: string | undefined;
            validatorAddress?: string | undefined;
            startingInfo?: {
                previousPeriod?: bigint | undefined;
                stake?: string | undefined;
                height?: bigint | undefined;
            } | undefined;
        }[] | undefined;
        validatorSlashEvents?: {
            validatorAddress?: string | undefined;
            height?: bigint | undefined;
            period?: bigint | undefined;
            validatorSlashEvent?: {
                validatorPeriod?: bigint | undefined;
                fraction?: string | undefined;
            } | undefined;
        }[] | undefined;
    } & {
        params?: ({
            communityTax?: string | undefined;
            baseProposerReward?: string | undefined;
            bonusProposerReward?: string | undefined;
            withdrawAddrEnabled?: boolean | undefined;
        } & {
            communityTax?: string | undefined;
            baseProposerReward?: string | undefined;
            bonusProposerReward?: string | undefined;
            withdrawAddrEnabled?: boolean | undefined;
        } & Record<Exclude<keyof I["params"], keyof Params>, never>) | undefined;
        feePool?: ({
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
            } & Record<Exclude<keyof I["feePool"]["communityPool"][number], keyof DecCoin>, never>)[] & Record<Exclude<keyof I["feePool"]["communityPool"], keyof {
                denom?: string | undefined;
                amount?: string | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["feePool"], "communityPool">, never>) | undefined;
        delegatorWithdrawInfos?: ({
            delegatorAddress?: string | undefined;
            withdrawAddress?: string | undefined;
        }[] & ({
            delegatorAddress?: string | undefined;
            withdrawAddress?: string | undefined;
        } & {
            delegatorAddress?: string | undefined;
            withdrawAddress?: string | undefined;
        } & Record<Exclude<keyof I["delegatorWithdrawInfos"][number], keyof DelegatorWithdrawInfo>, never>)[] & Record<Exclude<keyof I["delegatorWithdrawInfos"], keyof {
            delegatorAddress?: string | undefined;
            withdrawAddress?: string | undefined;
        }[]>, never>) | undefined;
        previousProposer?: string | undefined;
        outstandingRewards?: ({
            validatorAddress?: string | undefined;
            outstandingRewards?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
        }[] & ({
            validatorAddress?: string | undefined;
            outstandingRewards?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
        } & {
            validatorAddress?: string | undefined;
            outstandingRewards?: ({
                denom?: string | undefined;
                amount?: string | undefined;
            }[] & ({
                denom?: string | undefined;
                amount?: string | undefined;
            } & {
                denom?: string | undefined;
                amount?: string | undefined;
            } & Record<Exclude<keyof I["outstandingRewards"][number]["outstandingRewards"][number], keyof DecCoin>, never>)[] & Record<Exclude<keyof I["outstandingRewards"][number]["outstandingRewards"], keyof {
                denom?: string | undefined;
                amount?: string | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["outstandingRewards"][number], keyof ValidatorOutstandingRewardsRecord>, never>)[] & Record<Exclude<keyof I["outstandingRewards"], keyof {
            validatorAddress?: string | undefined;
            outstandingRewards?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
        }[]>, never>) | undefined;
        validatorAccumulatedCommissions?: ({
            validatorAddress?: string | undefined;
            accumulated?: {
                commission?: {
                    denom?: string | undefined;
                    amount?: string | undefined;
                }[] | undefined;
            } | undefined;
        }[] & ({
            validatorAddress?: string | undefined;
            accumulated?: {
                commission?: {
                    denom?: string | undefined;
                    amount?: string | undefined;
                }[] | undefined;
            } | undefined;
        } & {
            validatorAddress?: string | undefined;
            accumulated?: ({
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
                } & Record<Exclude<keyof I["validatorAccumulatedCommissions"][number]["accumulated"]["commission"][number], keyof DecCoin>, never>)[] & Record<Exclude<keyof I["validatorAccumulatedCommissions"][number]["accumulated"]["commission"], keyof {
                    denom?: string | undefined;
                    amount?: string | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["validatorAccumulatedCommissions"][number]["accumulated"], "commission">, never>) | undefined;
        } & Record<Exclude<keyof I["validatorAccumulatedCommissions"][number], keyof ValidatorAccumulatedCommissionRecord>, never>)[] & Record<Exclude<keyof I["validatorAccumulatedCommissions"], keyof {
            validatorAddress?: string | undefined;
            accumulated?: {
                commission?: {
                    denom?: string | undefined;
                    amount?: string | undefined;
                }[] | undefined;
            } | undefined;
        }[]>, never>) | undefined;
        validatorHistoricalRewards?: ({
            validatorAddress?: string | undefined;
            period?: bigint | undefined;
            rewards?: {
                cumulativeRewardRatio?: {
                    denom?: string | undefined;
                    amount?: string | undefined;
                }[] | undefined;
                referenceCount?: number | undefined;
            } | undefined;
        }[] & ({
            validatorAddress?: string | undefined;
            period?: bigint | undefined;
            rewards?: {
                cumulativeRewardRatio?: {
                    denom?: string | undefined;
                    amount?: string | undefined;
                }[] | undefined;
                referenceCount?: number | undefined;
            } | undefined;
        } & {
            validatorAddress?: string | undefined;
            period?: bigint | undefined;
            rewards?: ({
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
                } & Record<Exclude<keyof I["validatorHistoricalRewards"][number]["rewards"]["cumulativeRewardRatio"][number], keyof DecCoin>, never>)[] & Record<Exclude<keyof I["validatorHistoricalRewards"][number]["rewards"]["cumulativeRewardRatio"], keyof {
                    denom?: string | undefined;
                    amount?: string | undefined;
                }[]>, never>) | undefined;
                referenceCount?: number | undefined;
            } & Record<Exclude<keyof I["validatorHistoricalRewards"][number]["rewards"], keyof ValidatorHistoricalRewards>, never>) | undefined;
        } & Record<Exclude<keyof I["validatorHistoricalRewards"][number], keyof ValidatorHistoricalRewardsRecord>, never>)[] & Record<Exclude<keyof I["validatorHistoricalRewards"], keyof {
            validatorAddress?: string | undefined;
            period?: bigint | undefined;
            rewards?: {
                cumulativeRewardRatio?: {
                    denom?: string | undefined;
                    amount?: string | undefined;
                }[] | undefined;
                referenceCount?: number | undefined;
            } | undefined;
        }[]>, never>) | undefined;
        validatorCurrentRewards?: ({
            validatorAddress?: string | undefined;
            rewards?: {
                rewards?: {
                    denom?: string | undefined;
                    amount?: string | undefined;
                }[] | undefined;
                period?: bigint | undefined;
            } | undefined;
        }[] & ({
            validatorAddress?: string | undefined;
            rewards?: {
                rewards?: {
                    denom?: string | undefined;
                    amount?: string | undefined;
                }[] | undefined;
                period?: bigint | undefined;
            } | undefined;
        } & {
            validatorAddress?: string | undefined;
            rewards?: ({
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
                } & Record<Exclude<keyof I["validatorCurrentRewards"][number]["rewards"]["rewards"][number], keyof DecCoin>, never>)[] & Record<Exclude<keyof I["validatorCurrentRewards"][number]["rewards"]["rewards"], keyof {
                    denom?: string | undefined;
                    amount?: string | undefined;
                }[]>, never>) | undefined;
                period?: bigint | undefined;
            } & Record<Exclude<keyof I["validatorCurrentRewards"][number]["rewards"], keyof ValidatorCurrentRewards>, never>) | undefined;
        } & Record<Exclude<keyof I["validatorCurrentRewards"][number], keyof ValidatorCurrentRewardsRecord>, never>)[] & Record<Exclude<keyof I["validatorCurrentRewards"], keyof {
            validatorAddress?: string | undefined;
            rewards?: {
                rewards?: {
                    denom?: string | undefined;
                    amount?: string | undefined;
                }[] | undefined;
                period?: bigint | undefined;
            } | undefined;
        }[]>, never>) | undefined;
        delegatorStartingInfos?: ({
            delegatorAddress?: string | undefined;
            validatorAddress?: string | undefined;
            startingInfo?: {
                previousPeriod?: bigint | undefined;
                stake?: string | undefined;
                height?: bigint | undefined;
            } | undefined;
        }[] & ({
            delegatorAddress?: string | undefined;
            validatorAddress?: string | undefined;
            startingInfo?: {
                previousPeriod?: bigint | undefined;
                stake?: string | undefined;
                height?: bigint | undefined;
            } | undefined;
        } & {
            delegatorAddress?: string | undefined;
            validatorAddress?: string | undefined;
            startingInfo?: ({
                previousPeriod?: bigint | undefined;
                stake?: string | undefined;
                height?: bigint | undefined;
            } & {
                previousPeriod?: bigint | undefined;
                stake?: string | undefined;
                height?: bigint | undefined;
            } & Record<Exclude<keyof I["delegatorStartingInfos"][number]["startingInfo"], keyof DelegatorStartingInfo>, never>) | undefined;
        } & Record<Exclude<keyof I["delegatorStartingInfos"][number], keyof DelegatorStartingInfoRecord>, never>)[] & Record<Exclude<keyof I["delegatorStartingInfos"], keyof {
            delegatorAddress?: string | undefined;
            validatorAddress?: string | undefined;
            startingInfo?: {
                previousPeriod?: bigint | undefined;
                stake?: string | undefined;
                height?: bigint | undefined;
            } | undefined;
        }[]>, never>) | undefined;
        validatorSlashEvents?: ({
            validatorAddress?: string | undefined;
            height?: bigint | undefined;
            period?: bigint | undefined;
            validatorSlashEvent?: {
                validatorPeriod?: bigint | undefined;
                fraction?: string | undefined;
            } | undefined;
        }[] & ({
            validatorAddress?: string | undefined;
            height?: bigint | undefined;
            period?: bigint | undefined;
            validatorSlashEvent?: {
                validatorPeriod?: bigint | undefined;
                fraction?: string | undefined;
            } | undefined;
        } & {
            validatorAddress?: string | undefined;
            height?: bigint | undefined;
            period?: bigint | undefined;
            validatorSlashEvent?: ({
                validatorPeriod?: bigint | undefined;
                fraction?: string | undefined;
            } & {
                validatorPeriod?: bigint | undefined;
                fraction?: string | undefined;
            } & Record<Exclude<keyof I["validatorSlashEvents"][number]["validatorSlashEvent"], keyof ValidatorSlashEvent>, never>) | undefined;
        } & Record<Exclude<keyof I["validatorSlashEvents"][number], keyof ValidatorSlashEventRecord>, never>)[] & Record<Exclude<keyof I["validatorSlashEvents"], keyof {
            validatorAddress?: string | undefined;
            height?: bigint | undefined;
            period?: bigint | undefined;
            validatorSlashEvent?: {
                validatorPeriod?: bigint | undefined;
                fraction?: string | undefined;
            } | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof GenesisState>, never>>(object: I): GenesisState;
};
