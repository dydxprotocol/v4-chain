import { EncodeObject, GeneratedType } from "@cosmjs/proto-signing";
import { MsgBeginRedelegate, MsgCancelUnbondingDelegation, MsgCreateValidator, MsgDelegate, MsgEditValidator, MsgUndelegate } from "cosmjs-types/cosmos/staking/v1beta1/tx";
export declare const stakingTypes: ReadonlyArray<[string, GeneratedType]>;
export interface MsgBeginRedelegateEncodeObject extends EncodeObject {
    readonly typeUrl: "/cosmos.staking.v1beta1.MsgBeginRedelegate";
    readonly value: Partial<MsgBeginRedelegate>;
}
export declare function isMsgBeginRedelegateEncodeObject(o: EncodeObject): o is MsgBeginRedelegateEncodeObject;
export interface MsgCreateValidatorEncodeObject extends EncodeObject {
    readonly typeUrl: "/cosmos.staking.v1beta1.MsgCreateValidator";
    readonly value: Partial<MsgCreateValidator>;
}
export declare function isMsgCreateValidatorEncodeObject(o: EncodeObject): o is MsgCreateValidatorEncodeObject;
export interface MsgDelegateEncodeObject extends EncodeObject {
    readonly typeUrl: "/cosmos.staking.v1beta1.MsgDelegate";
    readonly value: Partial<MsgDelegate>;
}
export declare function isMsgDelegateEncodeObject(object: EncodeObject): object is MsgDelegateEncodeObject;
export interface MsgEditValidatorEncodeObject extends EncodeObject {
    readonly typeUrl: "/cosmos.staking.v1beta1.MsgEditValidator";
    readonly value: Partial<MsgEditValidator>;
}
export declare function isMsgEditValidatorEncodeObject(o: EncodeObject): o is MsgEditValidatorEncodeObject;
export interface MsgUndelegateEncodeObject extends EncodeObject {
    readonly typeUrl: "/cosmos.staking.v1beta1.MsgUndelegate";
    readonly value: Partial<MsgUndelegate>;
}
export declare function isMsgUndelegateEncodeObject(object: EncodeObject): object is MsgUndelegateEncodeObject;
export interface MsgCancelUnbondingDelegationEncodeObject extends EncodeObject {
    readonly typeUrl: "/cosmos.staking.v1beta1.MsgCancelUnbondingDelegation";
    readonly value: Partial<MsgCancelUnbondingDelegation>;
}
export declare function isMsgCancelUnbondingDelegationEncodeObject(object: EncodeObject): object is MsgCancelUnbondingDelegationEncodeObject;
