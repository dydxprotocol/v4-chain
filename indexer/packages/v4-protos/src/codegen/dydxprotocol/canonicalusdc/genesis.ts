import { Controls, ControlsSDKType, Ledger, LedgerSDKType, PendingSettlement, PendingSettlementSDKType } from "./canonical_usdc";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** GenesisState defines canonical USDC state. */

export interface GenesisState {
  controls?: Controls;
  ledger?: Ledger;
  pendingSettlements: PendingSettlement[];
}
/** GenesisState defines canonical USDC state. */

export interface GenesisStateSDKType {
  controls?: ControlsSDKType;
  ledger?: LedgerSDKType;
  pending_settlements: PendingSettlementSDKType[];
}

function createBaseGenesisState(): GenesisState {
  return {
    controls: undefined,
    ledger: undefined,
    pendingSettlements: []
  };
}

export const GenesisState = {
  encode(message: GenesisState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.controls !== undefined) {
      Controls.encode(message.controls, writer.uint32(10).fork()).ldelim();
    }

    if (message.ledger !== undefined) {
      Ledger.encode(message.ledger, writer.uint32(18).fork()).ldelim();
    }

    for (const v of message.pendingSettlements) {
      PendingSettlement.encode(v!, writer.uint32(26).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGenesisState();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.controls = Controls.decode(reader, reader.uint32());
          break;

        case 2:
          message.ledger = Ledger.decode(reader, reader.uint32());
          break;

        case 3:
          message.pendingSettlements.push(PendingSettlement.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<GenesisState>): GenesisState {
    const message = createBaseGenesisState();
    message.controls = object.controls !== undefined && object.controls !== null ? Controls.fromPartial(object.controls) : undefined;
    message.ledger = object.ledger !== undefined && object.ledger !== null ? Ledger.fromPartial(object.ledger) : undefined;
    message.pendingSettlements = object.pendingSettlements?.map(e => PendingSettlement.fromPartial(e)) || [];
    return message;
  }

};