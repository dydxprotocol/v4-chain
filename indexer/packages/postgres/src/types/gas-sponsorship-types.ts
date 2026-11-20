export enum GasSponsorshipStatus {
  // The grant transaction was signed, and possibly broadcast, but is not yet known to be in a
  // block.
  PENDING = 'PENDING',
  // The grant transaction was included in a block and succeeded.
  CONFIRMED = 'CONFIRMED',
  // The grant transaction was rejected or failed, so no allowance was given.
  FAILED = 'FAILED',
}

export interface GasSponsorshipCreateObject {
  id?: string,
  chain_id: string,
  granter_address: string,
  grantee_address: string,
  dydx_address: string,
  spend_limit: string,
  fee_amount: string,
  status: GasSponsorshipStatus,
  transaction_hash?: string,
  expires_at: string,
  created_at: string,
}

export interface GasSponsorshipUpdateObject {
  status?: GasSponsorshipStatus,
  // Null clears the hash, for a transaction that no node accepted and so does not exist.
  transaction_hash?: string | null,
}

export enum GasSponsorshipColumns {
  id = 'id',
  chain_id = 'chain_id',
  granter_address = 'granter_address',
  grantee_address = 'grantee_address',
  dydx_address = 'dydx_address',
  spend_limit = 'spend_limit',
  fee_amount = 'fee_amount',
  status = 'status',
  transaction_hash = 'transaction_hash',
  expires_at = 'expires_at',
  created_at = 'created_at',
}

export interface GasSponsorshipQueryFilters {
  chain_id?: string,
  grantee_address?: string,
  dydx_address?: string,
  statuses?: GasSponsorshipStatus[],
  // Inclusive lower bound on created_at.
  createdOnOrAfter?: string,
}
