import * as _m0 from "protobufjs/minimal";
import { DeepPartial, Long } from "../../helpers";
/** Mode controls whether canonical USDC routing is active. */

export enum Mode {
  /** MODE_UNSPECIFIED - MODE_UNSPECIFIED is invalid in configured state. */
  MODE_UNSPECIFIED = 0,

  /** MODE_DISABLED - MODE_DISABLED preserves the legacy transfer stack exactly. */
  MODE_DISABLED = 1,

  /** MODE_GRADUAL - MODE_GRADUAL enables canonical Injective routing for logical USDC. */
  MODE_GRADUAL = 2,

  /** MODE_PAUSED - MODE_PAUSED rejects canonical USDC transfers while retaining pending state. */
  MODE_PAUSED = 3,
  UNRECOGNIZED = -1,
}
/** Mode controls whether canonical USDC routing is active. */

export enum ModeSDKType {
  /** MODE_UNSPECIFIED - MODE_UNSPECIFIED is invalid in configured state. */
  MODE_UNSPECIFIED = 0,

  /** MODE_DISABLED - MODE_DISABLED preserves the legacy transfer stack exactly. */
  MODE_DISABLED = 1,

  /** MODE_GRADUAL - MODE_GRADUAL enables canonical Injective routing for logical USDC. */
  MODE_GRADUAL = 2,

  /** MODE_PAUSED - MODE_PAUSED rejects canonical USDC transfers while retaining pending state. */
  MODE_PAUSED = 3,
  UNRECOGNIZED = -1,
}
export function modeFromJSON(object: any): Mode {
  switch (object) {
    case 0:
    case "MODE_UNSPECIFIED":
      return Mode.MODE_UNSPECIFIED;

    case 1:
    case "MODE_DISABLED":
      return Mode.MODE_DISABLED;

    case 2:
    case "MODE_GRADUAL":
      return Mode.MODE_GRADUAL;

    case 3:
    case "MODE_PAUSED":
      return Mode.MODE_PAUSED;

    case -1:
    case "UNRECOGNIZED":
    default:
      return Mode.UNRECOGNIZED;
  }
}
export function modeToJSON(object: Mode): string {
  switch (object) {
    case Mode.MODE_UNSPECIFIED:
      return "MODE_UNSPECIFIED";

    case Mode.MODE_DISABLED:
      return "MODE_DISABLED";

    case Mode.MODE_GRADUAL:
      return "MODE_GRADUAL";

    case Mode.MODE_PAUSED:
      return "MODE_PAUSED";

    case Mode.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
/** Route identifies the physical settlement route for a logical USDC transfer. */

export enum Route {
  /** ROUTE_UNSPECIFIED - ROUTE_UNSPECIFIED is not a valid settlement route. */
  ROUTE_UNSPECIFIED = 0,

  /** ROUTE_INJECTIVE - ROUTE_INJECTIVE settles with Injective-origin USDC. */
  ROUTE_INJECTIVE = 1,

  /** ROUTE_NOBLE - ROUTE_NOBLE settles with Noble-origin USDC. */
  ROUTE_NOBLE = 2,
  UNRECOGNIZED = -1,
}
/** Route identifies the physical settlement route for a logical USDC transfer. */

export enum RouteSDKType {
  /** ROUTE_UNSPECIFIED - ROUTE_UNSPECIFIED is not a valid settlement route. */
  ROUTE_UNSPECIFIED = 0,

  /** ROUTE_INJECTIVE - ROUTE_INJECTIVE settles with Injective-origin USDC. */
  ROUTE_INJECTIVE = 1,

  /** ROUTE_NOBLE - ROUTE_NOBLE settles with Noble-origin USDC. */
  ROUTE_NOBLE = 2,
  UNRECOGNIZED = -1,
}
export function routeFromJSON(object: any): Route {
  switch (object) {
    case 0:
    case "ROUTE_UNSPECIFIED":
      return Route.ROUTE_UNSPECIFIED;

    case 1:
    case "ROUTE_INJECTIVE":
      return Route.ROUTE_INJECTIVE;

    case 2:
    case "ROUTE_NOBLE":
      return Route.ROUTE_NOBLE;

    case -1:
    case "UNRECOGNIZED":
    default:
      return Route.UNRECOGNIZED;
  }
}
export function routeToJSON(object: Route): string {
  switch (object) {
    case Route.ROUTE_UNSPECIFIED:
      return "ROUTE_UNSPECIFIED";

    case Route.ROUTE_INJECTIVE:
      return "ROUTE_INJECTIVE";

    case Route.ROUTE_NOBLE:
      return "ROUTE_NOBLE";

    case Route.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
/** Controls binds canonical routing to exact channels and denominations. */

export interface Controls {
  /** mode is the routing mode for canonical USDC. */
  mode: Mode;
  /**
   * logical_denom is the application USDC denom (the `uusdc` IBC hash) whose
   * backing is classified.
   */

  logicalDenom: string;
  /** noble_channel is the local IBC channel ID bound to Noble. */

  nobleChannel: string;
  /** noble_client is the IBC client the Noble channel must be bound to. */

  nobleClient: string;
  /** noble_connection is the IBC connection the Noble channel must be bound to. */

  nobleConnection: string;
  /**
   * noble_packet_denom is Noble's base denom as it appears in an inbound ICS20
   * packet (e.g. `uusdc`), not the full trace path. Outbound packets carry
   * `transfer/<noble_channel>/<noble_packet_denom>`.
   */

  noblePacketDenom: string;
  /**
   * noble_withdrawals_enabled is whether logical USDC may exit over the Noble
   * channel.
   */

  nobleWithdrawalsEnabled: boolean;
  /** injective_channel is the local IBC channel ID bound to Injective. */

  injectiveChannel: string;
  /** injective_client is the IBC client the Injective channel must be bound to. */

  injectiveClient: string;
  /**
   * injective_connection is the IBC connection the Injective channel must be
   * bound to.
   */

  injectiveConnection: string;
  /**
   * injective_packet_denom is Injective's base denom as it appears in an
   * inbound ICS20 packet (e.g. `uusdc`), not the full trace path. Outbound
   * packets carry `transfer/<injective_channel>/<injective_packet_denom>`.
   */

  injectivePacketDenom: string;
  /**
   * injective_denom is the local IBC denom of Injective-origin USDC (the hash
   * of transfer/<injective_channel>/<injective_packet_denom>).
   */

  injectiveDenom: string;
  /**
   * injective_deposits_disabled stops Injective-origin USDC received on the
   * Injective channel from being credited as logical USDC. When true, those
   * deposits are rejected with an error acknowledgement and refunded on
   * Injective. Withdrawals and in-flight settlements are unaffected.
   */

  injectiveDepositsDisabled: boolean;
  /**
   * max_pending_settlements is the cap on concurrent in-flight canonical
   * settlements, at most 100000. It bounds module state and the scans over
   * it; it is not a throttle. A withdrawal on either route needs a free slot,
   * so it is sized well above the expected in-flight volume.
   */

  maxPendingSettlements: number;
  /**
   * noble_deposits_disabled rejects Noble USDC received on the Noble channel
   * with an error acknowledgement, refunding the sender on Noble. While false,
   * those deposits mint logical USDC through ICS20 as before and are credited
   * to noble_backing. It must be true before noble_withdrawals_enabled is
   * false, or deposits would create Noble backing that cannot leave.
   */

  nobleDepositsDisabled: boolean;
}
/** Controls binds canonical routing to exact channels and denominations. */

export interface ControlsSDKType {
  /** mode is the routing mode for canonical USDC. */
  mode: ModeSDKType;
  /**
   * logical_denom is the application USDC denom (the `uusdc` IBC hash) whose
   * backing is classified.
   */

  logical_denom: string;
  /** noble_channel is the local IBC channel ID bound to Noble. */

  noble_channel: string;
  /** noble_client is the IBC client the Noble channel must be bound to. */

  noble_client: string;
  /** noble_connection is the IBC connection the Noble channel must be bound to. */

  noble_connection: string;
  /**
   * noble_packet_denom is Noble's base denom as it appears in an inbound ICS20
   * packet (e.g. `uusdc`), not the full trace path. Outbound packets carry
   * `transfer/<noble_channel>/<noble_packet_denom>`.
   */

  noble_packet_denom: string;
  /**
   * noble_withdrawals_enabled is whether logical USDC may exit over the Noble
   * channel.
   */

  noble_withdrawals_enabled: boolean;
  /** injective_channel is the local IBC channel ID bound to Injective. */

  injective_channel: string;
  /** injective_client is the IBC client the Injective channel must be bound to. */

  injective_client: string;
  /**
   * injective_connection is the IBC connection the Injective channel must be
   * bound to.
   */

  injective_connection: string;
  /**
   * injective_packet_denom is Injective's base denom as it appears in an
   * inbound ICS20 packet (e.g. `uusdc`), not the full trace path. Outbound
   * packets carry `transfer/<injective_channel>/<injective_packet_denom>`.
   */

  injective_packet_denom: string;
  /**
   * injective_denom is the local IBC denom of Injective-origin USDC (the hash
   * of transfer/<injective_channel>/<injective_packet_denom>).
   */

  injective_denom: string;
  /**
   * injective_deposits_disabled stops Injective-origin USDC received on the
   * Injective channel from being credited as logical USDC. When true, those
   * deposits are rejected with an error acknowledgement and refunded on
   * Injective. Withdrawals and in-flight settlements are unaffected.
   */

  injective_deposits_disabled: boolean;
  /**
   * max_pending_settlements is the cap on concurrent in-flight canonical
   * settlements, at most 100000. It bounds module state and the scans over
   * it; it is not a throttle. A withdrawal on either route needs a free slot,
   * so it is sized well above the expected in-flight volume.
   */

  max_pending_settlements: number;
  /**
   * noble_deposits_disabled rejects Noble USDC received on the Noble channel
   * with an error acknowledgement, refunding the sender on Noble. While false,
   * those deposits mint logical USDC through ICS20 as before and are credited
   * to noble_backing. It must be true before noble_withdrawals_enabled is
   * false, or deposits would create Noble backing that cannot leave.
   */

  noble_deposits_disabled: boolean;
}
/**
 * Ledger accounts for backing and in-flight settlement without changing the
 * logical application denomination.
 */

export interface Ledger {
  /** noble_backing is the logical supply still classified as Noble-backed. */
  nobleBacking: Uint8Array;
  /**
   * injective_backing is the Injective-origin USDC held by the module backing
   * logical supply.
   */

  injectiveBacking: Uint8Array;
  /**
   * pending_injective is the logical USDC locked on the module for in-flight
   * Injective withdrawals.
   */

  pendingInjective: Uint8Array;
  /**
   * pending_noble is the logical USDC in flight to Noble awaiting ack or
   * timeout (already removed from supply).
   */

  pendingNoble: Uint8Array;
  /**
   * refillable_capacity is logical USDC withdrawal capacity that Injective
   * withdrawals consumed and that Injective deposits may restore, capped at
   * the largest x/ratelimit baseline for the logical denom. A failed Injective
   * settlement lowers it, since x/ratelimit undoes that charge itself.
   */

  refillableCapacity: Uint8Array;
}
/**
 * Ledger accounts for backing and in-flight settlement without changing the
 * logical application denomination.
 */

export interface LedgerSDKType {
  /** noble_backing is the logical supply still classified as Noble-backed. */
  noble_backing: Uint8Array;
  /**
   * injective_backing is the Injective-origin USDC held by the module backing
   * logical supply.
   */

  injective_backing: Uint8Array;
  /**
   * pending_injective is the logical USDC locked on the module for in-flight
   * Injective withdrawals.
   */

  pending_injective: Uint8Array;
  /**
   * pending_noble is the logical USDC in flight to Noble awaiting ack or
   * timeout (already removed from supply).
   */

  pending_noble: Uint8Array;
  /**
   * refillable_capacity is logical USDC withdrawal capacity that Injective
   * withdrawals consumed and that Injective deposits may restore, capped at
   * the largest x/ratelimit baseline for the logical denom. A failed Injective
   * settlement lowers it, since x/ratelimit undoes that charge itself.
   */

  refillable_capacity: Uint8Array;
}
/**
 * PendingSettlement records enough logical and physical state to settle a
 * packet acknowledgement or timeout exactly once.
 */

export interface PendingSettlement {
  /** route is the physical settlement route of the packet. */
  route: Route;
  /** source_channel is the source channel of the IBC packet. */

  sourceChannel: string;
  /** sequence is the sequence of the IBC packet. */

  sequence: Long;
  /** sender is the dYdX sender of the MsgTransfer. */

  sender: string;
  /** logical_receiver is the counterparty receiver of the logical transfer. */

  logicalReceiver: string;
  /** physical_receiver is the counterparty receiver of the physical packet. */

  physicalReceiver: string;
  /** logical_denom is the application USDC denom of the transfer. */

  logicalDenom: string;
  /** physical_denom is the denom carried by the physical packet. */

  physicalDenom: string;
  /** amount is the transfer amount. */

  amount: Uint8Array;
}
/**
 * PendingSettlement records enough logical and physical state to settle a
 * packet acknowledgement or timeout exactly once.
 */

export interface PendingSettlementSDKType {
  /** route is the physical settlement route of the packet. */
  route: RouteSDKType;
  /** source_channel is the source channel of the IBC packet. */

  source_channel: string;
  /** sequence is the sequence of the IBC packet. */

  sequence: Long;
  /** sender is the dYdX sender of the MsgTransfer. */

  sender: string;
  /** logical_receiver is the counterparty receiver of the logical transfer. */

  logical_receiver: string;
  /** physical_receiver is the counterparty receiver of the physical packet. */

  physical_receiver: string;
  /** logical_denom is the application USDC denom of the transfer. */

  logical_denom: string;
  /** physical_denom is the denom carried by the physical packet. */

  physical_denom: string;
  /** amount is the transfer amount. */

  amount: Uint8Array;
}

function createBaseControls(): Controls {
  return {
    mode: 0,
    logicalDenom: "",
    nobleChannel: "",
    nobleClient: "",
    nobleConnection: "",
    noblePacketDenom: "",
    nobleWithdrawalsEnabled: false,
    injectiveChannel: "",
    injectiveClient: "",
    injectiveConnection: "",
    injectivePacketDenom: "",
    injectiveDenom: "",
    injectiveDepositsDisabled: false,
    maxPendingSettlements: 0,
    nobleDepositsDisabled: false
  };
}

export const Controls = {
  encode(message: Controls, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.mode !== 0) {
      writer.uint32(8).int32(message.mode);
    }

    if (message.logicalDenom !== "") {
      writer.uint32(18).string(message.logicalDenom);
    }

    if (message.nobleChannel !== "") {
      writer.uint32(26).string(message.nobleChannel);
    }

    if (message.nobleClient !== "") {
      writer.uint32(34).string(message.nobleClient);
    }

    if (message.nobleConnection !== "") {
      writer.uint32(42).string(message.nobleConnection);
    }

    if (message.noblePacketDenom !== "") {
      writer.uint32(50).string(message.noblePacketDenom);
    }

    if (message.nobleWithdrawalsEnabled === true) {
      writer.uint32(56).bool(message.nobleWithdrawalsEnabled);
    }

    if (message.injectiveChannel !== "") {
      writer.uint32(66).string(message.injectiveChannel);
    }

    if (message.injectiveClient !== "") {
      writer.uint32(74).string(message.injectiveClient);
    }

    if (message.injectiveConnection !== "") {
      writer.uint32(82).string(message.injectiveConnection);
    }

    if (message.injectivePacketDenom !== "") {
      writer.uint32(90).string(message.injectivePacketDenom);
    }

    if (message.injectiveDenom !== "") {
      writer.uint32(98).string(message.injectiveDenom);
    }

    if (message.injectiveDepositsDisabled === true) {
      writer.uint32(104).bool(message.injectiveDepositsDisabled);
    }

    if (message.maxPendingSettlements !== 0) {
      writer.uint32(112).uint32(message.maxPendingSettlements);
    }

    if (message.nobleDepositsDisabled === true) {
      writer.uint32(120).bool(message.nobleDepositsDisabled);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Controls {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseControls();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.mode = (reader.int32() as any);
          break;

        case 2:
          message.logicalDenom = reader.string();
          break;

        case 3:
          message.nobleChannel = reader.string();
          break;

        case 4:
          message.nobleClient = reader.string();
          break;

        case 5:
          message.nobleConnection = reader.string();
          break;

        case 6:
          message.noblePacketDenom = reader.string();
          break;

        case 7:
          message.nobleWithdrawalsEnabled = reader.bool();
          break;

        case 8:
          message.injectiveChannel = reader.string();
          break;

        case 9:
          message.injectiveClient = reader.string();
          break;

        case 10:
          message.injectiveConnection = reader.string();
          break;

        case 11:
          message.injectivePacketDenom = reader.string();
          break;

        case 12:
          message.injectiveDenom = reader.string();
          break;

        case 13:
          message.injectiveDepositsDisabled = reader.bool();
          break;

        case 14:
          message.maxPendingSettlements = reader.uint32();
          break;

        case 15:
          message.nobleDepositsDisabled = reader.bool();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<Controls>): Controls {
    const message = createBaseControls();
    message.mode = object.mode ?? 0;
    message.logicalDenom = object.logicalDenom ?? "";
    message.nobleChannel = object.nobleChannel ?? "";
    message.nobleClient = object.nobleClient ?? "";
    message.nobleConnection = object.nobleConnection ?? "";
    message.noblePacketDenom = object.noblePacketDenom ?? "";
    message.nobleWithdrawalsEnabled = object.nobleWithdrawalsEnabled ?? false;
    message.injectiveChannel = object.injectiveChannel ?? "";
    message.injectiveClient = object.injectiveClient ?? "";
    message.injectiveConnection = object.injectiveConnection ?? "";
    message.injectivePacketDenom = object.injectivePacketDenom ?? "";
    message.injectiveDenom = object.injectiveDenom ?? "";
    message.injectiveDepositsDisabled = object.injectiveDepositsDisabled ?? false;
    message.maxPendingSettlements = object.maxPendingSettlements ?? 0;
    message.nobleDepositsDisabled = object.nobleDepositsDisabled ?? false;
    return message;
  }

};

function createBaseLedger(): Ledger {
  return {
    nobleBacking: new Uint8Array(),
    injectiveBacking: new Uint8Array(),
    pendingInjective: new Uint8Array(),
    pendingNoble: new Uint8Array(),
    refillableCapacity: new Uint8Array()
  };
}

export const Ledger = {
  encode(message: Ledger, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.nobleBacking.length !== 0) {
      writer.uint32(10).bytes(message.nobleBacking);
    }

    if (message.injectiveBacking.length !== 0) {
      writer.uint32(18).bytes(message.injectiveBacking);
    }

    if (message.pendingInjective.length !== 0) {
      writer.uint32(26).bytes(message.pendingInjective);
    }

    if (message.pendingNoble.length !== 0) {
      writer.uint32(34).bytes(message.pendingNoble);
    }

    if (message.refillableCapacity.length !== 0) {
      writer.uint32(42).bytes(message.refillableCapacity);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Ledger {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLedger();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.nobleBacking = reader.bytes();
          break;

        case 2:
          message.injectiveBacking = reader.bytes();
          break;

        case 3:
          message.pendingInjective = reader.bytes();
          break;

        case 4:
          message.pendingNoble = reader.bytes();
          break;

        case 5:
          message.refillableCapacity = reader.bytes();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<Ledger>): Ledger {
    const message = createBaseLedger();
    message.nobleBacking = object.nobleBacking ?? new Uint8Array();
    message.injectiveBacking = object.injectiveBacking ?? new Uint8Array();
    message.pendingInjective = object.pendingInjective ?? new Uint8Array();
    message.pendingNoble = object.pendingNoble ?? new Uint8Array();
    message.refillableCapacity = object.refillableCapacity ?? new Uint8Array();
    return message;
  }

};

function createBasePendingSettlement(): PendingSettlement {
  return {
    route: 0,
    sourceChannel: "",
    sequence: Long.UZERO,
    sender: "",
    logicalReceiver: "",
    physicalReceiver: "",
    logicalDenom: "",
    physicalDenom: "",
    amount: new Uint8Array()
  };
}

export const PendingSettlement = {
  encode(message: PendingSettlement, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.route !== 0) {
      writer.uint32(8).int32(message.route);
    }

    if (message.sourceChannel !== "") {
      writer.uint32(18).string(message.sourceChannel);
    }

    if (!message.sequence.isZero()) {
      writer.uint32(24).uint64(message.sequence);
    }

    if (message.sender !== "") {
      writer.uint32(34).string(message.sender);
    }

    if (message.logicalReceiver !== "") {
      writer.uint32(42).string(message.logicalReceiver);
    }

    if (message.physicalReceiver !== "") {
      writer.uint32(50).string(message.physicalReceiver);
    }

    if (message.logicalDenom !== "") {
      writer.uint32(58).string(message.logicalDenom);
    }

    if (message.physicalDenom !== "") {
      writer.uint32(66).string(message.physicalDenom);
    }

    if (message.amount.length !== 0) {
      writer.uint32(74).bytes(message.amount);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PendingSettlement {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePendingSettlement();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.route = (reader.int32() as any);
          break;

        case 2:
          message.sourceChannel = reader.string();
          break;

        case 3:
          message.sequence = (reader.uint64() as Long);
          break;

        case 4:
          message.sender = reader.string();
          break;

        case 5:
          message.logicalReceiver = reader.string();
          break;

        case 6:
          message.physicalReceiver = reader.string();
          break;

        case 7:
          message.logicalDenom = reader.string();
          break;

        case 8:
          message.physicalDenom = reader.string();
          break;

        case 9:
          message.amount = reader.bytes();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<PendingSettlement>): PendingSettlement {
    const message = createBasePendingSettlement();
    message.route = object.route ?? 0;
    message.sourceChannel = object.sourceChannel ?? "";
    message.sequence = object.sequence !== undefined && object.sequence !== null ? Long.fromValue(object.sequence) : Long.UZERO;
    message.sender = object.sender ?? "";
    message.logicalReceiver = object.logicalReceiver ?? "";
    message.physicalReceiver = object.physicalReceiver ?? "";
    message.logicalDenom = object.logicalDenom ?? "";
    message.physicalDenom = object.physicalDenom ?? "";
    message.amount = object.amount ?? new Uint8Array();
    return message;
  }

};