/** VaultStatus represents the status of a vault. */
export enum VaultStatus {
  /** VAULT_STATUS_UNSPECIFIED - Default value, invalid and unused. */
  VAULT_STATUS_UNSPECIFIED = 0,

  /** VAULT_STATUS_DEACTIVATED - Don’t place orders. Does not count toward global vault balances. */
  VAULT_STATUS_DEACTIVATED = 1,

  /** VAULT_STATUS_STAND_BY - Don’t place orders. Does count towards global vault balances. */
  VAULT_STATUS_STAND_BY = 2,

  /** VAULT_STATUS_QUOTING - Places orders on both sides of the book. */
  VAULT_STATUS_QUOTING = 3,

  /** VAULT_STATUS_CLOSE_ONLY - Only place orders that close the position. */
  VAULT_STATUS_CLOSE_ONLY = 4,
  UNRECOGNIZED = -1,
}
/** VaultStatus represents the status of a vault. */

export enum VaultStatusSDKType {
  /** VAULT_STATUS_UNSPECIFIED - Default value, invalid and unused. */
  VAULT_STATUS_UNSPECIFIED = 0,

  /** VAULT_STATUS_DEACTIVATED - Don’t place orders. Does not count toward global vault balances. */
  VAULT_STATUS_DEACTIVATED = 1,

  /** VAULT_STATUS_STAND_BY - Don’t place orders. Does count towards global vault balances. */
  VAULT_STATUS_STAND_BY = 2,

  /** VAULT_STATUS_QUOTING - Places orders on both sides of the book. */
  VAULT_STATUS_QUOTING = 3,

  /** VAULT_STATUS_CLOSE_ONLY - Only place orders that close the position. */
  VAULT_STATUS_CLOSE_ONLY = 4,
  UNRECOGNIZED = -1,
}
export function vaultStatusFromJSON(object: any): VaultStatus {
  switch (object) {
    case 0:
    case "VAULT_STATUS_UNSPECIFIED":
      return VaultStatus.VAULT_STATUS_UNSPECIFIED;

    case 1:
    case "VAULT_STATUS_DEACTIVATED":
      return VaultStatus.VAULT_STATUS_DEACTIVATED;

    case 2:
    case "VAULT_STATUS_STAND_BY":
      return VaultStatus.VAULT_STATUS_STAND_BY;

    case 3:
    case "VAULT_STATUS_QUOTING":
      return VaultStatus.VAULT_STATUS_QUOTING;

    case 4:
    case "VAULT_STATUS_CLOSE_ONLY":
      return VaultStatus.VAULT_STATUS_CLOSE_ONLY;

    case -1:
    case "UNRECOGNIZED":
    default:
      return VaultStatus.UNRECOGNIZED;
  }
}
export function vaultStatusToJSON(object: VaultStatus): string {
  switch (object) {
    case VaultStatus.VAULT_STATUS_UNSPECIFIED:
      return "VAULT_STATUS_UNSPECIFIED";

    case VaultStatus.VAULT_STATUS_DEACTIVATED:
      return "VAULT_STATUS_DEACTIVATED";

    case VaultStatus.VAULT_STATUS_STAND_BY:
      return "VAULT_STATUS_STAND_BY";

    case VaultStatus.VAULT_STATUS_QUOTING:
      return "VAULT_STATUS_QUOTING";

    case VaultStatus.VAULT_STATUS_CLOSE_ONLY:
      return "VAULT_STATUS_CLOSE_ONLY";

    case VaultStatus.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}