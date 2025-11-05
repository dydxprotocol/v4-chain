export interface TurnkeyUserCreateObject {
  suborg_id: string,
  email?: string,
  svm_address: string,
  evm_address: string,
  smart_account_address?: string,
  salt: string,
  dydx_address?: string,
  created_at: string,
}

export enum TurnkeyUserColumns {
  suborg_id = 'suborg_id',
  email = 'email',
  svm_address = 'svm_address',
  evm_address = 'evm_address',
  smart_account_address = 'smart_account_address',
  salt = 'salt',
  dydx_address = 'dydx_address',
  created_at = 'created_at',
}
