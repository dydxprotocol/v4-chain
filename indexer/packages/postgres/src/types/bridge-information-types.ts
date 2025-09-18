export interface BridgeInformationCreateObject {
  id?: string,
  from_address: string,
  chain_id: string,
  amount: string,
  transaction_hash?: string,
  created_at: string,
}

export enum BridgeInformationColumns {
  id = 'id',
  from_address = 'from_address',
  chain_id = 'chain_id',
  amount = 'amount',
  transaction_hash = 'transaction_hash',
  created_at = 'created_at',
}

export interface BridgeInformationQueryFilters {
  from_address?: string,
  chain_id?: string,
  transaction_hash?: string,
  has_transaction_hash?: boolean, // true for NOT NULL, false for NULL
}

export interface BridgeInformationQueryOptions {
  orderBy?: 'created_at' | 'amount',
  orderDirection?: 'ASC' | 'DESC',
  limit?: number,
  offset?: number,
}
