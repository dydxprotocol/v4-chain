/* ------- SUBACCOUNT USERNAME TYPES ------- */

export interface SubaccountUsernamesCreateObject {
  username: string,
  subaccountId: string,
}

export enum SubaccountUsernamesColumns {
  username = 'username',
  subaccountId = 'subaccountId',
}

export interface SubaccountsWithoutUsernamesResult {
  subaccountId: string,
  address: string,
}
