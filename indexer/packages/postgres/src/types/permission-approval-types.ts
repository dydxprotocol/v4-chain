export interface PermissionApprovalCreateObject {
  suborg_id: string,
  chain_id: string,
  approval: string,
}

export interface PermissionApprovalFromDatabase {
  suborg_id: string,
  chain_id: string,
  approval: string,
}

export enum PermissionApprovalColumns {
  suborg_id = 'suborg_id',
  chain_id = 'chain_id',
  approval = 'approval',
}
