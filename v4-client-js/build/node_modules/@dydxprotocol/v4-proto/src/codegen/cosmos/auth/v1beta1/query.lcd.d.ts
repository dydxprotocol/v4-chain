import { LCDClient } from "@osmonauts/lcd";
import { QueryAccountsRequest, QueryAccountsResponseSDKType, QueryAccountRequest, QueryAccountResponseSDKType, QueryAccountAddressByIDRequest, QueryAccountAddressByIDResponseSDKType, QueryParamsRequest, QueryParamsResponseSDKType, QueryModuleAccountsRequest, QueryModuleAccountsResponseSDKType, QueryModuleAccountByNameRequest, QueryModuleAccountByNameResponseSDKType, Bech32PrefixRequest, Bech32PrefixResponseSDKType, AddressBytesToStringRequest, AddressBytesToStringResponseSDKType, AddressStringToBytesRequest, AddressStringToBytesResponseSDKType, QueryAccountInfoRequest, QueryAccountInfoResponseSDKType } from "./query";
export declare class LCDQueryClient {
    req: LCDClient;
    constructor({ requestClient }: {
        requestClient: LCDClient;
    });
    accounts(params?: QueryAccountsRequest): Promise<QueryAccountsResponseSDKType>;
    account(params: QueryAccountRequest): Promise<QueryAccountResponseSDKType>;
    accountAddressByID(params: QueryAccountAddressByIDRequest): Promise<QueryAccountAddressByIDResponseSDKType>;
    params(_params?: QueryParamsRequest): Promise<QueryParamsResponseSDKType>;
    moduleAccounts(_params?: QueryModuleAccountsRequest): Promise<QueryModuleAccountsResponseSDKType>;
    moduleAccountByName(params: QueryModuleAccountByNameRequest): Promise<QueryModuleAccountByNameResponseSDKType>;
    bech32Prefix(_params?: Bech32PrefixRequest): Promise<Bech32PrefixResponseSDKType>;
    addressBytesToString(params: AddressBytesToStringRequest): Promise<AddressBytesToStringResponseSDKType>;
    addressStringToBytes(params: AddressStringToBytesRequest): Promise<AddressStringToBytesResponseSDKType>;
    accountInfo(params: QueryAccountInfoRequest): Promise<QueryAccountInfoResponseSDKType>;
}
