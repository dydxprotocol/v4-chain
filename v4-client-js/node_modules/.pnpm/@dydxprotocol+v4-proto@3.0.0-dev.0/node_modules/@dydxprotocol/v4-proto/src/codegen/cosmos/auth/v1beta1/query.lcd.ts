import { setPaginationParams } from "../../../helpers";
import { LCDClient } from "@osmonauts/lcd";
import { QueryAccountsRequest, QueryAccountsResponseSDKType, QueryAccountRequest, QueryAccountResponseSDKType, QueryAccountAddressByIDRequest, QueryAccountAddressByIDResponseSDKType, QueryParamsRequest, QueryParamsResponseSDKType, QueryModuleAccountsRequest, QueryModuleAccountsResponseSDKType, QueryModuleAccountByNameRequest, QueryModuleAccountByNameResponseSDKType, Bech32PrefixRequest, Bech32PrefixResponseSDKType, AddressBytesToStringRequest, AddressBytesToStringResponseSDKType, AddressStringToBytesRequest, AddressStringToBytesResponseSDKType, QueryAccountInfoRequest, QueryAccountInfoResponseSDKType } from "./query";
export class LCDQueryClient {
  req: LCDClient;

  constructor({
    requestClient
  }: {
    requestClient: LCDClient;
  }) {
    this.req = requestClient;
    this.accounts = this.accounts.bind(this);
    this.account = this.account.bind(this);
    this.accountAddressByID = this.accountAddressByID.bind(this);
    this.params = this.params.bind(this);
    this.moduleAccounts = this.moduleAccounts.bind(this);
    this.moduleAccountByName = this.moduleAccountByName.bind(this);
    this.bech32Prefix = this.bech32Prefix.bind(this);
    this.addressBytesToString = this.addressBytesToString.bind(this);
    this.addressStringToBytes = this.addressStringToBytes.bind(this);
    this.accountInfo = this.accountInfo.bind(this);
  }
  /* Accounts returns all the existing accounts.
  
   When called from another module, this query might consume a high amount of
   gas if the pagination field is incorrectly set.
  
   Since: cosmos-sdk 0.43 */


  async accounts(params: QueryAccountsRequest = {
    pagination: undefined
  }): Promise<QueryAccountsResponseSDKType> {
    const options: any = {
      params: {}
    };

    if (typeof params?.pagination !== "undefined") {
      setPaginationParams(options, params.pagination);
    }

    const endpoint = `cosmos/auth/v1beta1/accounts`;
    return await this.req.get<QueryAccountsResponseSDKType>(endpoint, options);
  }
  /* Account returns account details based on address. */


  async account(params: QueryAccountRequest): Promise<QueryAccountResponseSDKType> {
    const endpoint = `cosmos/auth/v1beta1/accounts/${params.address}`;
    return await this.req.get<QueryAccountResponseSDKType>(endpoint);
  }
  /* AccountAddressByID returns account address based on account number.
  
   Since: cosmos-sdk 0.46.2 */


  async accountAddressByID(params: QueryAccountAddressByIDRequest): Promise<QueryAccountAddressByIDResponseSDKType> {
    const options: any = {
      params: {}
    };

    if (typeof params?.accountId !== "undefined") {
      options.params.account_id = params.accountId;
    }

    const endpoint = `cosmos/auth/v1beta1/address_by_id/${params.id}`;
    return await this.req.get<QueryAccountAddressByIDResponseSDKType>(endpoint, options);
  }
  /* Params queries all parameters. */


  async params(_params: QueryParamsRequest = {}): Promise<QueryParamsResponseSDKType> {
    const endpoint = `cosmos/auth/v1beta1/params`;
    return await this.req.get<QueryParamsResponseSDKType>(endpoint);
  }
  /* ModuleAccounts returns all the existing module accounts.
  
   Since: cosmos-sdk 0.46 */


  async moduleAccounts(_params: QueryModuleAccountsRequest = {}): Promise<QueryModuleAccountsResponseSDKType> {
    const endpoint = `cosmos/auth/v1beta1/module_accounts`;
    return await this.req.get<QueryModuleAccountsResponseSDKType>(endpoint);
  }
  /* ModuleAccountByName returns the module account info by module name */


  async moduleAccountByName(params: QueryModuleAccountByNameRequest): Promise<QueryModuleAccountByNameResponseSDKType> {
    const endpoint = `cosmos/auth/v1beta1/module_accounts/${params.name}`;
    return await this.req.get<QueryModuleAccountByNameResponseSDKType>(endpoint);
  }
  /* Bech32Prefix queries bech32Prefix
  
   Since: cosmos-sdk 0.46 */


  async bech32Prefix(_params: Bech32PrefixRequest = {}): Promise<Bech32PrefixResponseSDKType> {
    const endpoint = `cosmos/auth/v1beta1/bech32`;
    return await this.req.get<Bech32PrefixResponseSDKType>(endpoint);
  }
  /* AddressBytesToString converts Account Address bytes to string
  
   Since: cosmos-sdk 0.46 */


  async addressBytesToString(params: AddressBytesToStringRequest): Promise<AddressBytesToStringResponseSDKType> {
    const endpoint = `cosmos/auth/v1beta1/bech32/${params.addressBytes}`;
    return await this.req.get<AddressBytesToStringResponseSDKType>(endpoint);
  }
  /* AddressStringToBytes converts Address string to bytes
  
   Since: cosmos-sdk 0.46 */


  async addressStringToBytes(params: AddressStringToBytesRequest): Promise<AddressStringToBytesResponseSDKType> {
    const endpoint = `cosmos/auth/v1beta1/bech32/${params.addressString}`;
    return await this.req.get<AddressStringToBytesResponseSDKType>(endpoint);
  }
  /* AccountInfo queries account info which is common to all account types.
  
   Since: cosmos-sdk 0.47 */


  async accountInfo(params: QueryAccountInfoRequest): Promise<QueryAccountInfoResponseSDKType> {
    const endpoint = `cosmos/auth/v1beta1/account_info/${params.address}`;
    return await this.req.get<QueryAccountInfoResponseSDKType>(endpoint);
  }

}