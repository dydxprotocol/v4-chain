import { AminoMsg, Coin } from "@cosmjs/amino";
import { AminoConverters } from "../../aminotypes";
export interface AminoMsgCreateVestingAccount extends AminoMsg {
    readonly type: "cosmos-sdk/MsgCreateVestingAccount";
    readonly value: {
        /** Bech32 account address */
        readonly from_address: string;
        /** Bech32 account address */
        readonly to_address: string;
        readonly amount: readonly Coin[];
        readonly end_time: string;
        readonly delayed: boolean;
    };
}
export declare function isAminoMsgCreateVestingAccount(msg: AminoMsg): msg is AminoMsgCreateVestingAccount;
export declare function createVestingAminoConverters(): AminoConverters;
