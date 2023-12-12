import { SubaccountId, SubaccountIdSDKType } from "../../subaccounts/subaccount";
import { SubaccountOpenPositionInfo, SubaccountOpenPositionInfoSDKType } from "../../clob/liquidations";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../../helpers";
export interface LiquidateSubaccountsRequest_SubaccountOpenPositionInfoEntry {
  key: number;
  value?: SubaccountOpenPositionInfo;
}
export interface LiquidateSubaccountsRequest_SubaccountOpenPositionInfoEntrySDKType {
  key: number;
  value?: SubaccountOpenPositionInfoSDKType;
}
/**
 * LiquidateSubaccountsRequest is a request message that contains a list of
 * subaccount ids that potentially need to be liquidated. The list of subaccount
 * ids should not contain duplicates. The application should re-verify these
 * subaccount ids against current state before liquidating their positions.
 */

export interface LiquidateSubaccountsRequest {
  /** The block height at which the liquidation daemon is processing. */
  blockHeight: number;
  /** The list of liquidatable subaccount ids. */

  liquidatableSubaccountIds: SubaccountId[];
  /** The list of subaccount ids with negative total net collateral. */

  negativeTncSubaccountIds: SubaccountId[];
  subaccountOpenPositionInfo?: {
    [key: number]: SubaccountOpenPositionInfo;
  };
}
/**
 * LiquidateSubaccountsRequest is a request message that contains a list of
 * subaccount ids that potentially need to be liquidated. The list of subaccount
 * ids should not contain duplicates. The application should re-verify these
 * subaccount ids against current state before liquidating their positions.
 */

export interface LiquidateSubaccountsRequestSDKType {
  /** The block height at which the liquidation daemon is processing. */
  block_height: number;
  /** The list of liquidatable subaccount ids. */

  liquidatable_subaccount_ids: SubaccountIdSDKType[];
  /** The list of subaccount ids with negative total net collateral. */

  negative_tnc_subaccount_ids: SubaccountIdSDKType[];
  subaccount_open_position_info?: {
    [key: number]: SubaccountOpenPositionInfoSDKType;
  };
}
/**
 * LiquidateSubaccountsResponse is a response message for
 * LiquidateSubaccountsRequest.
 */

export interface LiquidateSubaccountsResponse {}
/**
 * LiquidateSubaccountsResponse is a response message for
 * LiquidateSubaccountsRequest.
 */

export interface LiquidateSubaccountsResponseSDKType {}

function createBaseLiquidateSubaccountsRequest_SubaccountOpenPositionInfoEntry(): LiquidateSubaccountsRequest_SubaccountOpenPositionInfoEntry {
  return {
    key: 0,
    value: undefined
  };
}

export const LiquidateSubaccountsRequest_SubaccountOpenPositionInfoEntry = {
  encode(message: LiquidateSubaccountsRequest_SubaccountOpenPositionInfoEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== 0) {
      writer.uint32(8).uint32(message.key);
    }

    if (message.value !== undefined) {
      dydxprotocol.clob.SubaccountOpenPositionInfo.encode(message.value, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): LiquidateSubaccountsRequest_SubaccountOpenPositionInfoEntry {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLiquidateSubaccountsRequest_SubaccountOpenPositionInfoEntry();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.key = reader.uint32();
          break;

        case 2:
          message.value = dydxprotocol.clob.SubaccountOpenPositionInfo.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<LiquidateSubaccountsRequest_SubaccountOpenPositionInfoEntry>): LiquidateSubaccountsRequest_SubaccountOpenPositionInfoEntry {
    const message = createBaseLiquidateSubaccountsRequest_SubaccountOpenPositionInfoEntry();
    message.key = object.key ?? 0;
    message.value = object.value !== undefined && object.value !== null ? dydxprotocol.clob.SubaccountOpenPositionInfo.fromPartial(object.value) : undefined;
    return message;
  }

};

function createBaseLiquidateSubaccountsRequest(): LiquidateSubaccountsRequest {
  return {
    blockHeight: 0,
    liquidatableSubaccountIds: [],
    negativeTncSubaccountIds: [],
    subaccountOpenPositionInfo: {}
  };
}

export const LiquidateSubaccountsRequest = {
  encode(message: LiquidateSubaccountsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.blockHeight !== 0) {
      writer.uint32(8).uint32(message.blockHeight);
    }

    for (const v of message.liquidatableSubaccountIds) {
      SubaccountId.encode(v!, writer.uint32(18).fork()).ldelim();
    }

    for (const v of message.negativeTncSubaccountIds) {
      SubaccountId.encode(v!, writer.uint32(26).fork()).ldelim();
    }

    Object.entries(message.subaccountOpenPositionInfo).forEach(([key, value]) => {
      LiquidateSubaccountsRequest_SubaccountOpenPositionInfoEntry.encode({
        key: (key as any),
        value
      }, writer.uint32(34).fork()).ldelim();
    });
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): LiquidateSubaccountsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLiquidateSubaccountsRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.blockHeight = reader.uint32();
          break;

        case 2:
          message.liquidatableSubaccountIds.push(SubaccountId.decode(reader, reader.uint32()));
          break;

        case 3:
          message.negativeTncSubaccountIds.push(SubaccountId.decode(reader, reader.uint32()));
          break;

        case 4:
          const entry4 = LiquidateSubaccountsRequest_SubaccountOpenPositionInfoEntry.decode(reader, reader.uint32());

          if (entry4.value !== undefined) {
            message.subaccountOpenPositionInfo[entry4.key] = entry4.value;
          }

          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<LiquidateSubaccountsRequest>): LiquidateSubaccountsRequest {
    const message = createBaseLiquidateSubaccountsRequest();
    message.blockHeight = object.blockHeight ?? 0;
    message.liquidatableSubaccountIds = object.liquidatableSubaccountIds?.map(e => SubaccountId.fromPartial(e)) || [];
    message.negativeTncSubaccountIds = object.negativeTncSubaccountIds?.map(e => SubaccountId.fromPartial(e)) || [];
    message.subaccountOpenPositionInfo = Object.entries(object.subaccountOpenPositionInfo ?? {}).reduce<{
      [key: number]: SubaccountOpenPositionInfo;
    }>((acc, [key, value]) => {
      if (value !== undefined) {
        acc[Number(key)] = SubaccountOpenPositionInfo.fromPartial(value);
      }

      return acc;
    }, {});
    return message;
  }

};

function createBaseLiquidateSubaccountsResponse(): LiquidateSubaccountsResponse {
  return {};
}

export const LiquidateSubaccountsResponse = {
  encode(_: LiquidateSubaccountsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): LiquidateSubaccountsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLiquidateSubaccountsResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(_: DeepPartial<LiquidateSubaccountsResponse>): LiquidateSubaccountsResponse {
    const message = createBaseLiquidateSubaccountsResponse();
    return message;
  }

};