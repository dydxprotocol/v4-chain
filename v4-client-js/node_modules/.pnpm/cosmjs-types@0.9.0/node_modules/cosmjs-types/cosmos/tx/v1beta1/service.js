"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ServiceClientImpl = exports.TxDecodeAminoResponse = exports.TxDecodeAminoRequest = exports.TxEncodeAminoResponse = exports.TxEncodeAminoRequest = exports.TxEncodeResponse = exports.TxEncodeRequest = exports.TxDecodeResponse = exports.TxDecodeRequest = exports.GetBlockWithTxsResponse = exports.GetBlockWithTxsRequest = exports.GetTxResponse = exports.GetTxRequest = exports.SimulateResponse = exports.SimulateRequest = exports.BroadcastTxResponse = exports.BroadcastTxRequest = exports.GetTxsEventResponse = exports.GetTxsEventRequest = exports.broadcastModeToJSON = exports.broadcastModeFromJSON = exports.BroadcastMode = exports.orderByToJSON = exports.orderByFromJSON = exports.OrderBy = exports.protobufPackage = void 0;
/* eslint-disable */
const tx_1 = require("./tx");
const pagination_1 = require("../../base/query/v1beta1/pagination");
const abci_1 = require("../../base/abci/v1beta1/abci");
const types_1 = require("../../../tendermint/types/types");
const block_1 = require("../../../tendermint/types/block");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.tx.v1beta1";
/** OrderBy defines the sorting order */
var OrderBy;
(function (OrderBy) {
    /** ORDER_BY_UNSPECIFIED - ORDER_BY_UNSPECIFIED specifies an unknown sorting order. OrderBy defaults to ASC in this case. */
    OrderBy[OrderBy["ORDER_BY_UNSPECIFIED"] = 0] = "ORDER_BY_UNSPECIFIED";
    /** ORDER_BY_ASC - ORDER_BY_ASC defines ascending order */
    OrderBy[OrderBy["ORDER_BY_ASC"] = 1] = "ORDER_BY_ASC";
    /** ORDER_BY_DESC - ORDER_BY_DESC defines descending order */
    OrderBy[OrderBy["ORDER_BY_DESC"] = 2] = "ORDER_BY_DESC";
    OrderBy[OrderBy["UNRECOGNIZED"] = -1] = "UNRECOGNIZED";
})(OrderBy || (exports.OrderBy = OrderBy = {}));
function orderByFromJSON(object) {
    switch (object) {
        case 0:
        case "ORDER_BY_UNSPECIFIED":
            return OrderBy.ORDER_BY_UNSPECIFIED;
        case 1:
        case "ORDER_BY_ASC":
            return OrderBy.ORDER_BY_ASC;
        case 2:
        case "ORDER_BY_DESC":
            return OrderBy.ORDER_BY_DESC;
        case -1:
        case "UNRECOGNIZED":
        default:
            return OrderBy.UNRECOGNIZED;
    }
}
exports.orderByFromJSON = orderByFromJSON;
function orderByToJSON(object) {
    switch (object) {
        case OrderBy.ORDER_BY_UNSPECIFIED:
            return "ORDER_BY_UNSPECIFIED";
        case OrderBy.ORDER_BY_ASC:
            return "ORDER_BY_ASC";
        case OrderBy.ORDER_BY_DESC:
            return "ORDER_BY_DESC";
        case OrderBy.UNRECOGNIZED:
        default:
            return "UNRECOGNIZED";
    }
}
exports.orderByToJSON = orderByToJSON;
/** BroadcastMode specifies the broadcast mode for the TxService.Broadcast RPC method. */
var BroadcastMode;
(function (BroadcastMode) {
    /** BROADCAST_MODE_UNSPECIFIED - zero-value for mode ordering */
    BroadcastMode[BroadcastMode["BROADCAST_MODE_UNSPECIFIED"] = 0] = "BROADCAST_MODE_UNSPECIFIED";
    /**
     * BROADCAST_MODE_BLOCK - DEPRECATED: use BROADCAST_MODE_SYNC instead,
     * BROADCAST_MODE_BLOCK is not supported by the SDK from v0.47.x onwards.
     */
    BroadcastMode[BroadcastMode["BROADCAST_MODE_BLOCK"] = 1] = "BROADCAST_MODE_BLOCK";
    /**
     * BROADCAST_MODE_SYNC - BROADCAST_MODE_SYNC defines a tx broadcasting mode where the client waits for
     * a CheckTx execution response only.
     */
    BroadcastMode[BroadcastMode["BROADCAST_MODE_SYNC"] = 2] = "BROADCAST_MODE_SYNC";
    /**
     * BROADCAST_MODE_ASYNC - BROADCAST_MODE_ASYNC defines a tx broadcasting mode where the client returns
     * immediately.
     */
    BroadcastMode[BroadcastMode["BROADCAST_MODE_ASYNC"] = 3] = "BROADCAST_MODE_ASYNC";
    BroadcastMode[BroadcastMode["UNRECOGNIZED"] = -1] = "UNRECOGNIZED";
})(BroadcastMode || (exports.BroadcastMode = BroadcastMode = {}));
function broadcastModeFromJSON(object) {
    switch (object) {
        case 0:
        case "BROADCAST_MODE_UNSPECIFIED":
            return BroadcastMode.BROADCAST_MODE_UNSPECIFIED;
        case 1:
        case "BROADCAST_MODE_BLOCK":
            return BroadcastMode.BROADCAST_MODE_BLOCK;
        case 2:
        case "BROADCAST_MODE_SYNC":
            return BroadcastMode.BROADCAST_MODE_SYNC;
        case 3:
        case "BROADCAST_MODE_ASYNC":
            return BroadcastMode.BROADCAST_MODE_ASYNC;
        case -1:
        case "UNRECOGNIZED":
        default:
            return BroadcastMode.UNRECOGNIZED;
    }
}
exports.broadcastModeFromJSON = broadcastModeFromJSON;
function broadcastModeToJSON(object) {
    switch (object) {
        case BroadcastMode.BROADCAST_MODE_UNSPECIFIED:
            return "BROADCAST_MODE_UNSPECIFIED";
        case BroadcastMode.BROADCAST_MODE_BLOCK:
            return "BROADCAST_MODE_BLOCK";
        case BroadcastMode.BROADCAST_MODE_SYNC:
            return "BROADCAST_MODE_SYNC";
        case BroadcastMode.BROADCAST_MODE_ASYNC:
            return "BROADCAST_MODE_ASYNC";
        case BroadcastMode.UNRECOGNIZED:
        default:
            return "UNRECOGNIZED";
    }
}
exports.broadcastModeToJSON = broadcastModeToJSON;
function createBaseGetTxsEventRequest() {
    return {
        events: [],
        pagination: undefined,
        orderBy: 0,
        page: BigInt(0),
        limit: BigInt(0),
    };
}
exports.GetTxsEventRequest = {
    typeUrl: "/cosmos.tx.v1beta1.GetTxsEventRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.events) {
            writer.uint32(10).string(v);
        }
        if (message.pagination !== undefined) {
            pagination_1.PageRequest.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        if (message.orderBy !== 0) {
            writer.uint32(24).int32(message.orderBy);
        }
        if (message.page !== BigInt(0)) {
            writer.uint32(32).uint64(message.page);
        }
        if (message.limit !== BigInt(0)) {
            writer.uint32(40).uint64(message.limit);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGetTxsEventRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.events.push(reader.string());
                    break;
                case 2:
                    message.pagination = pagination_1.PageRequest.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.orderBy = reader.int32();
                    break;
                case 4:
                    message.page = reader.uint64();
                    break;
                case 5:
                    message.limit = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseGetTxsEventRequest();
        if (Array.isArray(object?.events))
            obj.events = object.events.map((e) => String(e));
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageRequest.fromJSON(object.pagination);
        if ((0, helpers_1.isSet)(object.orderBy))
            obj.orderBy = orderByFromJSON(object.orderBy);
        if ((0, helpers_1.isSet)(object.page))
            obj.page = BigInt(object.page.toString());
        if ((0, helpers_1.isSet)(object.limit))
            obj.limit = BigInt(object.limit.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.events) {
            obj.events = message.events.map((e) => e);
        }
        else {
            obj.events = [];
        }
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageRequest.toJSON(message.pagination) : undefined);
        message.orderBy !== undefined && (obj.orderBy = orderByToJSON(message.orderBy));
        message.page !== undefined && (obj.page = (message.page || BigInt(0)).toString());
        message.limit !== undefined && (obj.limit = (message.limit || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGetTxsEventRequest();
        message.events = object.events?.map((e) => e) || [];
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageRequest.fromPartial(object.pagination);
        }
        message.orderBy = object.orderBy ?? 0;
        if (object.page !== undefined && object.page !== null) {
            message.page = BigInt(object.page.toString());
        }
        if (object.limit !== undefined && object.limit !== null) {
            message.limit = BigInt(object.limit.toString());
        }
        return message;
    },
};
function createBaseGetTxsEventResponse() {
    return {
        txs: [],
        txResponses: [],
        pagination: undefined,
        total: BigInt(0),
    };
}
exports.GetTxsEventResponse = {
    typeUrl: "/cosmos.tx.v1beta1.GetTxsEventResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.txs) {
            tx_1.Tx.encode(v, writer.uint32(10).fork()).ldelim();
        }
        for (const v of message.txResponses) {
            abci_1.TxResponse.encode(v, writer.uint32(18).fork()).ldelim();
        }
        if (message.pagination !== undefined) {
            pagination_1.PageResponse.encode(message.pagination, writer.uint32(26).fork()).ldelim();
        }
        if (message.total !== BigInt(0)) {
            writer.uint32(32).uint64(message.total);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGetTxsEventResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.txs.push(tx_1.Tx.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.txResponses.push(abci_1.TxResponse.decode(reader, reader.uint32()));
                    break;
                case 3:
                    message.pagination = pagination_1.PageResponse.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.total = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseGetTxsEventResponse();
        if (Array.isArray(object?.txs))
            obj.txs = object.txs.map((e) => tx_1.Tx.fromJSON(e));
        if (Array.isArray(object?.txResponses))
            obj.txResponses = object.txResponses.map((e) => abci_1.TxResponse.fromJSON(e));
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageResponse.fromJSON(object.pagination);
        if ((0, helpers_1.isSet)(object.total))
            obj.total = BigInt(object.total.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.txs) {
            obj.txs = message.txs.map((e) => (e ? tx_1.Tx.toJSON(e) : undefined));
        }
        else {
            obj.txs = [];
        }
        if (message.txResponses) {
            obj.txResponses = message.txResponses.map((e) => (e ? abci_1.TxResponse.toJSON(e) : undefined));
        }
        else {
            obj.txResponses = [];
        }
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageResponse.toJSON(message.pagination) : undefined);
        message.total !== undefined && (obj.total = (message.total || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGetTxsEventResponse();
        message.txs = object.txs?.map((e) => tx_1.Tx.fromPartial(e)) || [];
        message.txResponses = object.txResponses?.map((e) => abci_1.TxResponse.fromPartial(e)) || [];
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageResponse.fromPartial(object.pagination);
        }
        if (object.total !== undefined && object.total !== null) {
            message.total = BigInt(object.total.toString());
        }
        return message;
    },
};
function createBaseBroadcastTxRequest() {
    return {
        txBytes: new Uint8Array(),
        mode: 0,
    };
}
exports.BroadcastTxRequest = {
    typeUrl: "/cosmos.tx.v1beta1.BroadcastTxRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.txBytes.length !== 0) {
            writer.uint32(10).bytes(message.txBytes);
        }
        if (message.mode !== 0) {
            writer.uint32(16).int32(message.mode);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseBroadcastTxRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.txBytes = reader.bytes();
                    break;
                case 2:
                    message.mode = reader.int32();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseBroadcastTxRequest();
        if ((0, helpers_1.isSet)(object.txBytes))
            obj.txBytes = (0, helpers_1.bytesFromBase64)(object.txBytes);
        if ((0, helpers_1.isSet)(object.mode))
            obj.mode = broadcastModeFromJSON(object.mode);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.txBytes !== undefined &&
            (obj.txBytes = (0, helpers_1.base64FromBytes)(message.txBytes !== undefined ? message.txBytes : new Uint8Array()));
        message.mode !== undefined && (obj.mode = broadcastModeToJSON(message.mode));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseBroadcastTxRequest();
        message.txBytes = object.txBytes ?? new Uint8Array();
        message.mode = object.mode ?? 0;
        return message;
    },
};
function createBaseBroadcastTxResponse() {
    return {
        txResponse: undefined,
    };
}
exports.BroadcastTxResponse = {
    typeUrl: "/cosmos.tx.v1beta1.BroadcastTxResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.txResponse !== undefined) {
            abci_1.TxResponse.encode(message.txResponse, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseBroadcastTxResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.txResponse = abci_1.TxResponse.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseBroadcastTxResponse();
        if ((0, helpers_1.isSet)(object.txResponse))
            obj.txResponse = abci_1.TxResponse.fromJSON(object.txResponse);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.txResponse !== undefined &&
            (obj.txResponse = message.txResponse ? abci_1.TxResponse.toJSON(message.txResponse) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseBroadcastTxResponse();
        if (object.txResponse !== undefined && object.txResponse !== null) {
            message.txResponse = abci_1.TxResponse.fromPartial(object.txResponse);
        }
        return message;
    },
};
function createBaseSimulateRequest() {
    return {
        tx: undefined,
        txBytes: new Uint8Array(),
    };
}
exports.SimulateRequest = {
    typeUrl: "/cosmos.tx.v1beta1.SimulateRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.tx !== undefined) {
            tx_1.Tx.encode(message.tx, writer.uint32(10).fork()).ldelim();
        }
        if (message.txBytes.length !== 0) {
            writer.uint32(18).bytes(message.txBytes);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSimulateRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.tx = tx_1.Tx.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.txBytes = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseSimulateRequest();
        if ((0, helpers_1.isSet)(object.tx))
            obj.tx = tx_1.Tx.fromJSON(object.tx);
        if ((0, helpers_1.isSet)(object.txBytes))
            obj.txBytes = (0, helpers_1.bytesFromBase64)(object.txBytes);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.tx !== undefined && (obj.tx = message.tx ? tx_1.Tx.toJSON(message.tx) : undefined);
        message.txBytes !== undefined &&
            (obj.txBytes = (0, helpers_1.base64FromBytes)(message.txBytes !== undefined ? message.txBytes : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseSimulateRequest();
        if (object.tx !== undefined && object.tx !== null) {
            message.tx = tx_1.Tx.fromPartial(object.tx);
        }
        message.txBytes = object.txBytes ?? new Uint8Array();
        return message;
    },
};
function createBaseSimulateResponse() {
    return {
        gasInfo: undefined,
        result: undefined,
    };
}
exports.SimulateResponse = {
    typeUrl: "/cosmos.tx.v1beta1.SimulateResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.gasInfo !== undefined) {
            abci_1.GasInfo.encode(message.gasInfo, writer.uint32(10).fork()).ldelim();
        }
        if (message.result !== undefined) {
            abci_1.Result.encode(message.result, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSimulateResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.gasInfo = abci_1.GasInfo.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.result = abci_1.Result.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseSimulateResponse();
        if ((0, helpers_1.isSet)(object.gasInfo))
            obj.gasInfo = abci_1.GasInfo.fromJSON(object.gasInfo);
        if ((0, helpers_1.isSet)(object.result))
            obj.result = abci_1.Result.fromJSON(object.result);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.gasInfo !== undefined &&
            (obj.gasInfo = message.gasInfo ? abci_1.GasInfo.toJSON(message.gasInfo) : undefined);
        message.result !== undefined && (obj.result = message.result ? abci_1.Result.toJSON(message.result) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseSimulateResponse();
        if (object.gasInfo !== undefined && object.gasInfo !== null) {
            message.gasInfo = abci_1.GasInfo.fromPartial(object.gasInfo);
        }
        if (object.result !== undefined && object.result !== null) {
            message.result = abci_1.Result.fromPartial(object.result);
        }
        return message;
    },
};
function createBaseGetTxRequest() {
    return {
        hash: "",
    };
}
exports.GetTxRequest = {
    typeUrl: "/cosmos.tx.v1beta1.GetTxRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.hash !== "") {
            writer.uint32(10).string(message.hash);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGetTxRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.hash = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseGetTxRequest();
        if ((0, helpers_1.isSet)(object.hash))
            obj.hash = String(object.hash);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.hash !== undefined && (obj.hash = message.hash);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGetTxRequest();
        message.hash = object.hash ?? "";
        return message;
    },
};
function createBaseGetTxResponse() {
    return {
        tx: undefined,
        txResponse: undefined,
    };
}
exports.GetTxResponse = {
    typeUrl: "/cosmos.tx.v1beta1.GetTxResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.tx !== undefined) {
            tx_1.Tx.encode(message.tx, writer.uint32(10).fork()).ldelim();
        }
        if (message.txResponse !== undefined) {
            abci_1.TxResponse.encode(message.txResponse, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGetTxResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.tx = tx_1.Tx.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.txResponse = abci_1.TxResponse.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseGetTxResponse();
        if ((0, helpers_1.isSet)(object.tx))
            obj.tx = tx_1.Tx.fromJSON(object.tx);
        if ((0, helpers_1.isSet)(object.txResponse))
            obj.txResponse = abci_1.TxResponse.fromJSON(object.txResponse);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.tx !== undefined && (obj.tx = message.tx ? tx_1.Tx.toJSON(message.tx) : undefined);
        message.txResponse !== undefined &&
            (obj.txResponse = message.txResponse ? abci_1.TxResponse.toJSON(message.txResponse) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGetTxResponse();
        if (object.tx !== undefined && object.tx !== null) {
            message.tx = tx_1.Tx.fromPartial(object.tx);
        }
        if (object.txResponse !== undefined && object.txResponse !== null) {
            message.txResponse = abci_1.TxResponse.fromPartial(object.txResponse);
        }
        return message;
    },
};
function createBaseGetBlockWithTxsRequest() {
    return {
        height: BigInt(0),
        pagination: undefined,
    };
}
exports.GetBlockWithTxsRequest = {
    typeUrl: "/cosmos.tx.v1beta1.GetBlockWithTxsRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.height !== BigInt(0)) {
            writer.uint32(8).int64(message.height);
        }
        if (message.pagination !== undefined) {
            pagination_1.PageRequest.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGetBlockWithTxsRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.height = reader.int64();
                    break;
                case 2:
                    message.pagination = pagination_1.PageRequest.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseGetBlockWithTxsRequest();
        if ((0, helpers_1.isSet)(object.height))
            obj.height = BigInt(object.height.toString());
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageRequest.fromJSON(object.pagination);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.height !== undefined && (obj.height = (message.height || BigInt(0)).toString());
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageRequest.toJSON(message.pagination) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGetBlockWithTxsRequest();
        if (object.height !== undefined && object.height !== null) {
            message.height = BigInt(object.height.toString());
        }
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageRequest.fromPartial(object.pagination);
        }
        return message;
    },
};
function createBaseGetBlockWithTxsResponse() {
    return {
        txs: [],
        blockId: undefined,
        block: undefined,
        pagination: undefined,
    };
}
exports.GetBlockWithTxsResponse = {
    typeUrl: "/cosmos.tx.v1beta1.GetBlockWithTxsResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.txs) {
            tx_1.Tx.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.blockId !== undefined) {
            types_1.BlockID.encode(message.blockId, writer.uint32(18).fork()).ldelim();
        }
        if (message.block !== undefined) {
            block_1.Block.encode(message.block, writer.uint32(26).fork()).ldelim();
        }
        if (message.pagination !== undefined) {
            pagination_1.PageResponse.encode(message.pagination, writer.uint32(34).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGetBlockWithTxsResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.txs.push(tx_1.Tx.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.blockId = types_1.BlockID.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.block = block_1.Block.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.pagination = pagination_1.PageResponse.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseGetBlockWithTxsResponse();
        if (Array.isArray(object?.txs))
            obj.txs = object.txs.map((e) => tx_1.Tx.fromJSON(e));
        if ((0, helpers_1.isSet)(object.blockId))
            obj.blockId = types_1.BlockID.fromJSON(object.blockId);
        if ((0, helpers_1.isSet)(object.block))
            obj.block = block_1.Block.fromJSON(object.block);
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageResponse.fromJSON(object.pagination);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.txs) {
            obj.txs = message.txs.map((e) => (e ? tx_1.Tx.toJSON(e) : undefined));
        }
        else {
            obj.txs = [];
        }
        message.blockId !== undefined &&
            (obj.blockId = message.blockId ? types_1.BlockID.toJSON(message.blockId) : undefined);
        message.block !== undefined && (obj.block = message.block ? block_1.Block.toJSON(message.block) : undefined);
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageResponse.toJSON(message.pagination) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGetBlockWithTxsResponse();
        message.txs = object.txs?.map((e) => tx_1.Tx.fromPartial(e)) || [];
        if (object.blockId !== undefined && object.blockId !== null) {
            message.blockId = types_1.BlockID.fromPartial(object.blockId);
        }
        if (object.block !== undefined && object.block !== null) {
            message.block = block_1.Block.fromPartial(object.block);
        }
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageResponse.fromPartial(object.pagination);
        }
        return message;
    },
};
function createBaseTxDecodeRequest() {
    return {
        txBytes: new Uint8Array(),
    };
}
exports.TxDecodeRequest = {
    typeUrl: "/cosmos.tx.v1beta1.TxDecodeRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.txBytes.length !== 0) {
            writer.uint32(10).bytes(message.txBytes);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseTxDecodeRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.txBytes = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseTxDecodeRequest();
        if ((0, helpers_1.isSet)(object.txBytes))
            obj.txBytes = (0, helpers_1.bytesFromBase64)(object.txBytes);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.txBytes !== undefined &&
            (obj.txBytes = (0, helpers_1.base64FromBytes)(message.txBytes !== undefined ? message.txBytes : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseTxDecodeRequest();
        message.txBytes = object.txBytes ?? new Uint8Array();
        return message;
    },
};
function createBaseTxDecodeResponse() {
    return {
        tx: undefined,
    };
}
exports.TxDecodeResponse = {
    typeUrl: "/cosmos.tx.v1beta1.TxDecodeResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.tx !== undefined) {
            tx_1.Tx.encode(message.tx, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseTxDecodeResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.tx = tx_1.Tx.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseTxDecodeResponse();
        if ((0, helpers_1.isSet)(object.tx))
            obj.tx = tx_1.Tx.fromJSON(object.tx);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.tx !== undefined && (obj.tx = message.tx ? tx_1.Tx.toJSON(message.tx) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseTxDecodeResponse();
        if (object.tx !== undefined && object.tx !== null) {
            message.tx = tx_1.Tx.fromPartial(object.tx);
        }
        return message;
    },
};
function createBaseTxEncodeRequest() {
    return {
        tx: undefined,
    };
}
exports.TxEncodeRequest = {
    typeUrl: "/cosmos.tx.v1beta1.TxEncodeRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.tx !== undefined) {
            tx_1.Tx.encode(message.tx, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseTxEncodeRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.tx = tx_1.Tx.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseTxEncodeRequest();
        if ((0, helpers_1.isSet)(object.tx))
            obj.tx = tx_1.Tx.fromJSON(object.tx);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.tx !== undefined && (obj.tx = message.tx ? tx_1.Tx.toJSON(message.tx) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseTxEncodeRequest();
        if (object.tx !== undefined && object.tx !== null) {
            message.tx = tx_1.Tx.fromPartial(object.tx);
        }
        return message;
    },
};
function createBaseTxEncodeResponse() {
    return {
        txBytes: new Uint8Array(),
    };
}
exports.TxEncodeResponse = {
    typeUrl: "/cosmos.tx.v1beta1.TxEncodeResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.txBytes.length !== 0) {
            writer.uint32(10).bytes(message.txBytes);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseTxEncodeResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.txBytes = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseTxEncodeResponse();
        if ((0, helpers_1.isSet)(object.txBytes))
            obj.txBytes = (0, helpers_1.bytesFromBase64)(object.txBytes);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.txBytes !== undefined &&
            (obj.txBytes = (0, helpers_1.base64FromBytes)(message.txBytes !== undefined ? message.txBytes : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseTxEncodeResponse();
        message.txBytes = object.txBytes ?? new Uint8Array();
        return message;
    },
};
function createBaseTxEncodeAminoRequest() {
    return {
        aminoJson: "",
    };
}
exports.TxEncodeAminoRequest = {
    typeUrl: "/cosmos.tx.v1beta1.TxEncodeAminoRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.aminoJson !== "") {
            writer.uint32(10).string(message.aminoJson);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseTxEncodeAminoRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.aminoJson = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseTxEncodeAminoRequest();
        if ((0, helpers_1.isSet)(object.aminoJson))
            obj.aminoJson = String(object.aminoJson);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.aminoJson !== undefined && (obj.aminoJson = message.aminoJson);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseTxEncodeAminoRequest();
        message.aminoJson = object.aminoJson ?? "";
        return message;
    },
};
function createBaseTxEncodeAminoResponse() {
    return {
        aminoBinary: new Uint8Array(),
    };
}
exports.TxEncodeAminoResponse = {
    typeUrl: "/cosmos.tx.v1beta1.TxEncodeAminoResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.aminoBinary.length !== 0) {
            writer.uint32(10).bytes(message.aminoBinary);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseTxEncodeAminoResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.aminoBinary = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseTxEncodeAminoResponse();
        if ((0, helpers_1.isSet)(object.aminoBinary))
            obj.aminoBinary = (0, helpers_1.bytesFromBase64)(object.aminoBinary);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.aminoBinary !== undefined &&
            (obj.aminoBinary = (0, helpers_1.base64FromBytes)(message.aminoBinary !== undefined ? message.aminoBinary : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseTxEncodeAminoResponse();
        message.aminoBinary = object.aminoBinary ?? new Uint8Array();
        return message;
    },
};
function createBaseTxDecodeAminoRequest() {
    return {
        aminoBinary: new Uint8Array(),
    };
}
exports.TxDecodeAminoRequest = {
    typeUrl: "/cosmos.tx.v1beta1.TxDecodeAminoRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.aminoBinary.length !== 0) {
            writer.uint32(10).bytes(message.aminoBinary);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseTxDecodeAminoRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.aminoBinary = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseTxDecodeAminoRequest();
        if ((0, helpers_1.isSet)(object.aminoBinary))
            obj.aminoBinary = (0, helpers_1.bytesFromBase64)(object.aminoBinary);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.aminoBinary !== undefined &&
            (obj.aminoBinary = (0, helpers_1.base64FromBytes)(message.aminoBinary !== undefined ? message.aminoBinary : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseTxDecodeAminoRequest();
        message.aminoBinary = object.aminoBinary ?? new Uint8Array();
        return message;
    },
};
function createBaseTxDecodeAminoResponse() {
    return {
        aminoJson: "",
    };
}
exports.TxDecodeAminoResponse = {
    typeUrl: "/cosmos.tx.v1beta1.TxDecodeAminoResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.aminoJson !== "") {
            writer.uint32(10).string(message.aminoJson);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseTxDecodeAminoResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.aminoJson = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseTxDecodeAminoResponse();
        if ((0, helpers_1.isSet)(object.aminoJson))
            obj.aminoJson = String(object.aminoJson);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.aminoJson !== undefined && (obj.aminoJson = message.aminoJson);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseTxDecodeAminoResponse();
        message.aminoJson = object.aminoJson ?? "";
        return message;
    },
};
class ServiceClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
        this.Simulate = this.Simulate.bind(this);
        this.GetTx = this.GetTx.bind(this);
        this.BroadcastTx = this.BroadcastTx.bind(this);
        this.GetTxsEvent = this.GetTxsEvent.bind(this);
        this.GetBlockWithTxs = this.GetBlockWithTxs.bind(this);
        this.TxDecode = this.TxDecode.bind(this);
        this.TxEncode = this.TxEncode.bind(this);
        this.TxEncodeAmino = this.TxEncodeAmino.bind(this);
        this.TxDecodeAmino = this.TxDecodeAmino.bind(this);
    }
    Simulate(request) {
        const data = exports.SimulateRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.tx.v1beta1.Service", "Simulate", data);
        return promise.then((data) => exports.SimulateResponse.decode(new binary_1.BinaryReader(data)));
    }
    GetTx(request) {
        const data = exports.GetTxRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.tx.v1beta1.Service", "GetTx", data);
        return promise.then((data) => exports.GetTxResponse.decode(new binary_1.BinaryReader(data)));
    }
    BroadcastTx(request) {
        const data = exports.BroadcastTxRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.tx.v1beta1.Service", "BroadcastTx", data);
        return promise.then((data) => exports.BroadcastTxResponse.decode(new binary_1.BinaryReader(data)));
    }
    GetTxsEvent(request) {
        const data = exports.GetTxsEventRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.tx.v1beta1.Service", "GetTxsEvent", data);
        return promise.then((data) => exports.GetTxsEventResponse.decode(new binary_1.BinaryReader(data)));
    }
    GetBlockWithTxs(request) {
        const data = exports.GetBlockWithTxsRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.tx.v1beta1.Service", "GetBlockWithTxs", data);
        return promise.then((data) => exports.GetBlockWithTxsResponse.decode(new binary_1.BinaryReader(data)));
    }
    TxDecode(request) {
        const data = exports.TxDecodeRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.tx.v1beta1.Service", "TxDecode", data);
        return promise.then((data) => exports.TxDecodeResponse.decode(new binary_1.BinaryReader(data)));
    }
    TxEncode(request) {
        const data = exports.TxEncodeRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.tx.v1beta1.Service", "TxEncode", data);
        return promise.then((data) => exports.TxEncodeResponse.decode(new binary_1.BinaryReader(data)));
    }
    TxEncodeAmino(request) {
        const data = exports.TxEncodeAminoRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.tx.v1beta1.Service", "TxEncodeAmino", data);
        return promise.then((data) => exports.TxEncodeAminoResponse.decode(new binary_1.BinaryReader(data)));
    }
    TxDecodeAmino(request) {
        const data = exports.TxDecodeAminoRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.tx.v1beta1.Service", "TxDecodeAmino", data);
        return promise.then((data) => exports.TxDecodeAminoResponse.decode(new binary_1.BinaryReader(data)));
    }
}
exports.ServiceClientImpl = ServiceClientImpl;
//# sourceMappingURL=service.js.map