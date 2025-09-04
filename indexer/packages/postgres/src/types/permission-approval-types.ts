export interface PermissionApprovalCreateObject {
  suborg_id: string,
  chain_id: ChainId,
  approval: string,
}

export interface PermissionApprovalFromDatabase {
  suborg_id: string,
  chain_id: ChainId,
  approval: string,
}

export enum PermissionApprovalColumns {
  suborg_id = 'suborg_id',
  chain_id = 'chain_id',
  approval = 'approval',
}

export enum ChainId {
  ARBITRUM = 'arbitrum',
  BASE = 'base',
  AVALANCHE = 'avalanche',
  OPTIMISM = 'optimism',
  ETHEREUM = 'ethereum',
}
