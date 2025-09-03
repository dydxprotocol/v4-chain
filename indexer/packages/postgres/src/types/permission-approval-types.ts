export interface PermissionApprovalCreateObject {
  suborg_id: string,
  arbitrum_approval?: string,
  base_approval?: string,
  avalanche_approval?: string,
  optimism_approval?: string,
  ethereum_approval?: string,
}

export interface PermissionApprovalFromDatabase {
  suborg_id: string,
  arbitrum_approval?: string,
  base_approval?: string,
  avalanche_approval?: string,
  optimism_approval?: string,
  ethereum_approval?: string,
}

export enum PermissionApprovalColumns {
  suborg_id = 'suborg_id',
  arbitrum_approval = 'arbitrum_approval',
  base_approval = 'base_approval',
  avalanche_approval = 'avalanche_approval',
  optimism_approval = 'optimism_approval',
  ethereum_approval = 'ethereum_approval',
}
