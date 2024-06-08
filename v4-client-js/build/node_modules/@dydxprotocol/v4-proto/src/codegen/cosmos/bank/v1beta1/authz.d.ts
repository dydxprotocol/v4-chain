import { Coin, CoinSDKType } from "../../base/v1beta1/coin";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../../helpers";
/**
 * SendAuthorization allows the grantee to spend up to spend_limit coins from
 * the granter's account.
 *
 * Since: cosmos-sdk 0.43
 */
export interface SendAuthorization {
    spendLimit: Coin[];
    /**
     * allow_list specifies an optional list of addresses to whom the grantee can send tokens on behalf of the
     * granter. If omitted, any recipient is allowed.
     *
     * Since: cosmos-sdk 0.47
     */
    allowList: string[];
}
/**
 * SendAuthorization allows the grantee to spend up to spend_limit coins from
 * the granter's account.
 *
 * Since: cosmos-sdk 0.43
 */
export interface SendAuthorizationSDKType {
    spend_limit: CoinSDKType[];
    allow_list: string[];
}
export declare const SendAuthorization: {
    encode(message: SendAuthorization, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): SendAuthorization;
    fromPartial(object: DeepPartial<SendAuthorization>): SendAuthorization;
};
