/* ------- TURNKEY USER TYPES ------- */

export interface TurnkeyUserCreateObject {
  suborgId: string,
  username?: string,
  email?: string,
  svmAddress: string,
  evmAddress: string,
  salt: string,
  dydxAddress?: string,
  createdAt: string,
}

export enum TurnkeyUserColumns {
  suborgId = 'suborgId',
  username = 'username',
  email = 'email',
  svmAddress = 'svmAddress',
  evmAddress = 'evmAddress',
  salt = 'salt',
  dydxAddress = 'dydxAddress',
  createdAt = 'createdAt',
}
