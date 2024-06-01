"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ResponsePrepareProposal = exports.ResponseApplySnapshotChunk = exports.ResponseLoadSnapshotChunk = exports.ResponseOfferSnapshot = exports.ResponseListSnapshots = exports.ResponseCommit = exports.ResponseEndBlock = exports.ResponseDeliverTx = exports.ResponseCheckTx = exports.ResponseBeginBlock = exports.ResponseQuery = exports.ResponseInitChain = exports.ResponseInfo = exports.ResponseFlush = exports.ResponseEcho = exports.ResponseException = exports.Response = exports.RequestProcessProposal = exports.RequestPrepareProposal = exports.RequestApplySnapshotChunk = exports.RequestLoadSnapshotChunk = exports.RequestOfferSnapshot = exports.RequestListSnapshots = exports.RequestCommit = exports.RequestEndBlock = exports.RequestDeliverTx = exports.RequestCheckTx = exports.RequestBeginBlock = exports.RequestQuery = exports.RequestInitChain = exports.RequestInfo = exports.RequestFlush = exports.RequestEcho = exports.Request = exports.misbehaviorTypeToJSON = exports.misbehaviorTypeFromJSON = exports.MisbehaviorType = exports.responseProcessProposal_ProposalStatusToJSON = exports.responseProcessProposal_ProposalStatusFromJSON = exports.ResponseProcessProposal_ProposalStatus = exports.responseApplySnapshotChunk_ResultToJSON = exports.responseApplySnapshotChunk_ResultFromJSON = exports.ResponseApplySnapshotChunk_Result = exports.responseOfferSnapshot_ResultToJSON = exports.responseOfferSnapshot_ResultFromJSON = exports.ResponseOfferSnapshot_Result = exports.checkTxTypeToJSON = exports.checkTxTypeFromJSON = exports.CheckTxType = exports.protobufPackage = void 0;
exports.ABCIApplicationClientImpl = exports.Snapshot = exports.Misbehavior = exports.ExtendedVoteInfo = exports.VoteInfo = exports.ValidatorUpdate = exports.Validator = exports.TxResult = exports.EventAttribute = exports.Event = exports.ExtendedCommitInfo = exports.CommitInfo = exports.ResponseProcessProposal = void 0;
/* eslint-disable */
const timestamp_1 = require("../../google/protobuf/timestamp");
const params_1 = require("../types/params");
const types_1 = require("../types/types");
const proof_1 = require("../crypto/proof");
const keys_1 = require("../crypto/keys");
const binary_1 = require("../../binary");
const helpers_1 = require("../../helpers");
exports.protobufPackage = "tendermint.abci";
var CheckTxType;
(function (CheckTxType) {
    CheckTxType[CheckTxType["NEW"] = 0] = "NEW";
    CheckTxType[CheckTxType["RECHECK"] = 1] = "RECHECK";
    CheckTxType[CheckTxType["UNRECOGNIZED"] = -1] = "UNRECOGNIZED";
})(CheckTxType || (exports.CheckTxType = CheckTxType = {}));
function checkTxTypeFromJSON(object) {
    switch (object) {
        case 0:
        case "NEW":
            return CheckTxType.NEW;
        case 1:
        case "RECHECK":
            return CheckTxType.RECHECK;
        case -1:
        case "UNRECOGNIZED":
        default:
            return CheckTxType.UNRECOGNIZED;
    }
}
exports.checkTxTypeFromJSON = checkTxTypeFromJSON;
function checkTxTypeToJSON(object) {
    switch (object) {
        case CheckTxType.NEW:
            return "NEW";
        case CheckTxType.RECHECK:
            return "RECHECK";
        case CheckTxType.UNRECOGNIZED:
        default:
            return "UNRECOGNIZED";
    }
}
exports.checkTxTypeToJSON = checkTxTypeToJSON;
var ResponseOfferSnapshot_Result;
(function (ResponseOfferSnapshot_Result) {
    /** UNKNOWN - Unknown result, abort all snapshot restoration */
    ResponseOfferSnapshot_Result[ResponseOfferSnapshot_Result["UNKNOWN"] = 0] = "UNKNOWN";
    /** ACCEPT - Snapshot accepted, apply chunks */
    ResponseOfferSnapshot_Result[ResponseOfferSnapshot_Result["ACCEPT"] = 1] = "ACCEPT";
    /** ABORT - Abort all snapshot restoration */
    ResponseOfferSnapshot_Result[ResponseOfferSnapshot_Result["ABORT"] = 2] = "ABORT";
    /** REJECT - Reject this specific snapshot, try others */
    ResponseOfferSnapshot_Result[ResponseOfferSnapshot_Result["REJECT"] = 3] = "REJECT";
    /** REJECT_FORMAT - Reject all snapshots of this format, try others */
    ResponseOfferSnapshot_Result[ResponseOfferSnapshot_Result["REJECT_FORMAT"] = 4] = "REJECT_FORMAT";
    /** REJECT_SENDER - Reject all snapshots from the sender(s), try others */
    ResponseOfferSnapshot_Result[ResponseOfferSnapshot_Result["REJECT_SENDER"] = 5] = "REJECT_SENDER";
    ResponseOfferSnapshot_Result[ResponseOfferSnapshot_Result["UNRECOGNIZED"] = -1] = "UNRECOGNIZED";
})(ResponseOfferSnapshot_Result || (exports.ResponseOfferSnapshot_Result = ResponseOfferSnapshot_Result = {}));
function responseOfferSnapshot_ResultFromJSON(object) {
    switch (object) {
        case 0:
        case "UNKNOWN":
            return ResponseOfferSnapshot_Result.UNKNOWN;
        case 1:
        case "ACCEPT":
            return ResponseOfferSnapshot_Result.ACCEPT;
        case 2:
        case "ABORT":
            return ResponseOfferSnapshot_Result.ABORT;
        case 3:
        case "REJECT":
            return ResponseOfferSnapshot_Result.REJECT;
        case 4:
        case "REJECT_FORMAT":
            return ResponseOfferSnapshot_Result.REJECT_FORMAT;
        case 5:
        case "REJECT_SENDER":
            return ResponseOfferSnapshot_Result.REJECT_SENDER;
        case -1:
        case "UNRECOGNIZED":
        default:
            return ResponseOfferSnapshot_Result.UNRECOGNIZED;
    }
}
exports.responseOfferSnapshot_ResultFromJSON = responseOfferSnapshot_ResultFromJSON;
function responseOfferSnapshot_ResultToJSON(object) {
    switch (object) {
        case ResponseOfferSnapshot_Result.UNKNOWN:
            return "UNKNOWN";
        case ResponseOfferSnapshot_Result.ACCEPT:
            return "ACCEPT";
        case ResponseOfferSnapshot_Result.ABORT:
            return "ABORT";
        case ResponseOfferSnapshot_Result.REJECT:
            return "REJECT";
        case ResponseOfferSnapshot_Result.REJECT_FORMAT:
            return "REJECT_FORMAT";
        case ResponseOfferSnapshot_Result.REJECT_SENDER:
            return "REJECT_SENDER";
        case ResponseOfferSnapshot_Result.UNRECOGNIZED:
        default:
            return "UNRECOGNIZED";
    }
}
exports.responseOfferSnapshot_ResultToJSON = responseOfferSnapshot_ResultToJSON;
var ResponseApplySnapshotChunk_Result;
(function (ResponseApplySnapshotChunk_Result) {
    /** UNKNOWN - Unknown result, abort all snapshot restoration */
    ResponseApplySnapshotChunk_Result[ResponseApplySnapshotChunk_Result["UNKNOWN"] = 0] = "UNKNOWN";
    /** ACCEPT - Chunk successfully accepted */
    ResponseApplySnapshotChunk_Result[ResponseApplySnapshotChunk_Result["ACCEPT"] = 1] = "ACCEPT";
    /** ABORT - Abort all snapshot restoration */
    ResponseApplySnapshotChunk_Result[ResponseApplySnapshotChunk_Result["ABORT"] = 2] = "ABORT";
    /** RETRY - Retry chunk (combine with refetch and reject) */
    ResponseApplySnapshotChunk_Result[ResponseApplySnapshotChunk_Result["RETRY"] = 3] = "RETRY";
    /** RETRY_SNAPSHOT - Retry snapshot (combine with refetch and reject) */
    ResponseApplySnapshotChunk_Result[ResponseApplySnapshotChunk_Result["RETRY_SNAPSHOT"] = 4] = "RETRY_SNAPSHOT";
    /** REJECT_SNAPSHOT - Reject this snapshot, try others */
    ResponseApplySnapshotChunk_Result[ResponseApplySnapshotChunk_Result["REJECT_SNAPSHOT"] = 5] = "REJECT_SNAPSHOT";
    ResponseApplySnapshotChunk_Result[ResponseApplySnapshotChunk_Result["UNRECOGNIZED"] = -1] = "UNRECOGNIZED";
})(ResponseApplySnapshotChunk_Result || (exports.ResponseApplySnapshotChunk_Result = ResponseApplySnapshotChunk_Result = {}));
function responseApplySnapshotChunk_ResultFromJSON(object) {
    switch (object) {
        case 0:
        case "UNKNOWN":
            return ResponseApplySnapshotChunk_Result.UNKNOWN;
        case 1:
        case "ACCEPT":
            return ResponseApplySnapshotChunk_Result.ACCEPT;
        case 2:
        case "ABORT":
            return ResponseApplySnapshotChunk_Result.ABORT;
        case 3:
        case "RETRY":
            return ResponseApplySnapshotChunk_Result.RETRY;
        case 4:
        case "RETRY_SNAPSHOT":
            return ResponseApplySnapshotChunk_Result.RETRY_SNAPSHOT;
        case 5:
        case "REJECT_SNAPSHOT":
            return ResponseApplySnapshotChunk_Result.REJECT_SNAPSHOT;
        case -1:
        case "UNRECOGNIZED":
        default:
            return ResponseApplySnapshotChunk_Result.UNRECOGNIZED;
    }
}
exports.responseApplySnapshotChunk_ResultFromJSON = responseApplySnapshotChunk_ResultFromJSON;
function responseApplySnapshotChunk_ResultToJSON(object) {
    switch (object) {
        case ResponseApplySnapshotChunk_Result.UNKNOWN:
            return "UNKNOWN";
        case ResponseApplySnapshotChunk_Result.ACCEPT:
            return "ACCEPT";
        case ResponseApplySnapshotChunk_Result.ABORT:
            return "ABORT";
        case ResponseApplySnapshotChunk_Result.RETRY:
            return "RETRY";
        case ResponseApplySnapshotChunk_Result.RETRY_SNAPSHOT:
            return "RETRY_SNAPSHOT";
        case ResponseApplySnapshotChunk_Result.REJECT_SNAPSHOT:
            return "REJECT_SNAPSHOT";
        case ResponseApplySnapshotChunk_Result.UNRECOGNIZED:
        default:
            return "UNRECOGNIZED";
    }
}
exports.responseApplySnapshotChunk_ResultToJSON = responseApplySnapshotChunk_ResultToJSON;
var ResponseProcessProposal_ProposalStatus;
(function (ResponseProcessProposal_ProposalStatus) {
    ResponseProcessProposal_ProposalStatus[ResponseProcessProposal_ProposalStatus["UNKNOWN"] = 0] = "UNKNOWN";
    ResponseProcessProposal_ProposalStatus[ResponseProcessProposal_ProposalStatus["ACCEPT"] = 1] = "ACCEPT";
    ResponseProcessProposal_ProposalStatus[ResponseProcessProposal_ProposalStatus["REJECT"] = 2] = "REJECT";
    ResponseProcessProposal_ProposalStatus[ResponseProcessProposal_ProposalStatus["UNRECOGNIZED"] = -1] = "UNRECOGNIZED";
})(ResponseProcessProposal_ProposalStatus || (exports.ResponseProcessProposal_ProposalStatus = ResponseProcessProposal_ProposalStatus = {}));
function responseProcessProposal_ProposalStatusFromJSON(object) {
    switch (object) {
        case 0:
        case "UNKNOWN":
            return ResponseProcessProposal_ProposalStatus.UNKNOWN;
        case 1:
        case "ACCEPT":
            return ResponseProcessProposal_ProposalStatus.ACCEPT;
        case 2:
        case "REJECT":
            return ResponseProcessProposal_ProposalStatus.REJECT;
        case -1:
        case "UNRECOGNIZED":
        default:
            return ResponseProcessProposal_ProposalStatus.UNRECOGNIZED;
    }
}
exports.responseProcessProposal_ProposalStatusFromJSON = responseProcessProposal_ProposalStatusFromJSON;
function responseProcessProposal_ProposalStatusToJSON(object) {
    switch (object) {
        case ResponseProcessProposal_ProposalStatus.UNKNOWN:
            return "UNKNOWN";
        case ResponseProcessProposal_ProposalStatus.ACCEPT:
            return "ACCEPT";
        case ResponseProcessProposal_ProposalStatus.REJECT:
            return "REJECT";
        case ResponseProcessProposal_ProposalStatus.UNRECOGNIZED:
        default:
            return "UNRECOGNIZED";
    }
}
exports.responseProcessProposal_ProposalStatusToJSON = responseProcessProposal_ProposalStatusToJSON;
var MisbehaviorType;
(function (MisbehaviorType) {
    MisbehaviorType[MisbehaviorType["UNKNOWN"] = 0] = "UNKNOWN";
    MisbehaviorType[MisbehaviorType["DUPLICATE_VOTE"] = 1] = "DUPLICATE_VOTE";
    MisbehaviorType[MisbehaviorType["LIGHT_CLIENT_ATTACK"] = 2] = "LIGHT_CLIENT_ATTACK";
    MisbehaviorType[MisbehaviorType["UNRECOGNIZED"] = -1] = "UNRECOGNIZED";
})(MisbehaviorType || (exports.MisbehaviorType = MisbehaviorType = {}));
function misbehaviorTypeFromJSON(object) {
    switch (object) {
        case 0:
        case "UNKNOWN":
            return MisbehaviorType.UNKNOWN;
        case 1:
        case "DUPLICATE_VOTE":
            return MisbehaviorType.DUPLICATE_VOTE;
        case 2:
        case "LIGHT_CLIENT_ATTACK":
            return MisbehaviorType.LIGHT_CLIENT_ATTACK;
        case -1:
        case "UNRECOGNIZED":
        default:
            return MisbehaviorType.UNRECOGNIZED;
    }
}
exports.misbehaviorTypeFromJSON = misbehaviorTypeFromJSON;
function misbehaviorTypeToJSON(object) {
    switch (object) {
        case MisbehaviorType.UNKNOWN:
            return "UNKNOWN";
        case MisbehaviorType.DUPLICATE_VOTE:
            return "DUPLICATE_VOTE";
        case MisbehaviorType.LIGHT_CLIENT_ATTACK:
            return "LIGHT_CLIENT_ATTACK";
        case MisbehaviorType.UNRECOGNIZED:
        default:
            return "UNRECOGNIZED";
    }
}
exports.misbehaviorTypeToJSON = misbehaviorTypeToJSON;
function createBaseRequest() {
    return {
        echo: undefined,
        flush: undefined,
        info: undefined,
        initChain: undefined,
        query: undefined,
        beginBlock: undefined,
        checkTx: undefined,
        deliverTx: undefined,
        endBlock: undefined,
        commit: undefined,
        listSnapshots: undefined,
        offerSnapshot: undefined,
        loadSnapshotChunk: undefined,
        applySnapshotChunk: undefined,
        prepareProposal: undefined,
        processProposal: undefined,
    };
}
exports.Request = {
    typeUrl: "/tendermint.abci.Request",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.echo !== undefined) {
            exports.RequestEcho.encode(message.echo, writer.uint32(10).fork()).ldelim();
        }
        if (message.flush !== undefined) {
            exports.RequestFlush.encode(message.flush, writer.uint32(18).fork()).ldelim();
        }
        if (message.info !== undefined) {
            exports.RequestInfo.encode(message.info, writer.uint32(26).fork()).ldelim();
        }
        if (message.initChain !== undefined) {
            exports.RequestInitChain.encode(message.initChain, writer.uint32(42).fork()).ldelim();
        }
        if (message.query !== undefined) {
            exports.RequestQuery.encode(message.query, writer.uint32(50).fork()).ldelim();
        }
        if (message.beginBlock !== undefined) {
            exports.RequestBeginBlock.encode(message.beginBlock, writer.uint32(58).fork()).ldelim();
        }
        if (message.checkTx !== undefined) {
            exports.RequestCheckTx.encode(message.checkTx, writer.uint32(66).fork()).ldelim();
        }
        if (message.deliverTx !== undefined) {
            exports.RequestDeliverTx.encode(message.deliverTx, writer.uint32(74).fork()).ldelim();
        }
        if (message.endBlock !== undefined) {
            exports.RequestEndBlock.encode(message.endBlock, writer.uint32(82).fork()).ldelim();
        }
        if (message.commit !== undefined) {
            exports.RequestCommit.encode(message.commit, writer.uint32(90).fork()).ldelim();
        }
        if (message.listSnapshots !== undefined) {
            exports.RequestListSnapshots.encode(message.listSnapshots, writer.uint32(98).fork()).ldelim();
        }
        if (message.offerSnapshot !== undefined) {
            exports.RequestOfferSnapshot.encode(message.offerSnapshot, writer.uint32(106).fork()).ldelim();
        }
        if (message.loadSnapshotChunk !== undefined) {
            exports.RequestLoadSnapshotChunk.encode(message.loadSnapshotChunk, writer.uint32(114).fork()).ldelim();
        }
        if (message.applySnapshotChunk !== undefined) {
            exports.RequestApplySnapshotChunk.encode(message.applySnapshotChunk, writer.uint32(122).fork()).ldelim();
        }
        if (message.prepareProposal !== undefined) {
            exports.RequestPrepareProposal.encode(message.prepareProposal, writer.uint32(130).fork()).ldelim();
        }
        if (message.processProposal !== undefined) {
            exports.RequestProcessProposal.encode(message.processProposal, writer.uint32(138).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.echo = exports.RequestEcho.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.flush = exports.RequestFlush.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.info = exports.RequestInfo.decode(reader, reader.uint32());
                    break;
                case 5:
                    message.initChain = exports.RequestInitChain.decode(reader, reader.uint32());
                    break;
                case 6:
                    message.query = exports.RequestQuery.decode(reader, reader.uint32());
                    break;
                case 7:
                    message.beginBlock = exports.RequestBeginBlock.decode(reader, reader.uint32());
                    break;
                case 8:
                    message.checkTx = exports.RequestCheckTx.decode(reader, reader.uint32());
                    break;
                case 9:
                    message.deliverTx = exports.RequestDeliverTx.decode(reader, reader.uint32());
                    break;
                case 10:
                    message.endBlock = exports.RequestEndBlock.decode(reader, reader.uint32());
                    break;
                case 11:
                    message.commit = exports.RequestCommit.decode(reader, reader.uint32());
                    break;
                case 12:
                    message.listSnapshots = exports.RequestListSnapshots.decode(reader, reader.uint32());
                    break;
                case 13:
                    message.offerSnapshot = exports.RequestOfferSnapshot.decode(reader, reader.uint32());
                    break;
                case 14:
                    message.loadSnapshotChunk = exports.RequestLoadSnapshotChunk.decode(reader, reader.uint32());
                    break;
                case 15:
                    message.applySnapshotChunk = exports.RequestApplySnapshotChunk.decode(reader, reader.uint32());
                    break;
                case 16:
                    message.prepareProposal = exports.RequestPrepareProposal.decode(reader, reader.uint32());
                    break;
                case 17:
                    message.processProposal = exports.RequestProcessProposal.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseRequest();
        if ((0, helpers_1.isSet)(object.echo))
            obj.echo = exports.RequestEcho.fromJSON(object.echo);
        if ((0, helpers_1.isSet)(object.flush))
            obj.flush = exports.RequestFlush.fromJSON(object.flush);
        if ((0, helpers_1.isSet)(object.info))
            obj.info = exports.RequestInfo.fromJSON(object.info);
        if ((0, helpers_1.isSet)(object.initChain))
            obj.initChain = exports.RequestInitChain.fromJSON(object.initChain);
        if ((0, helpers_1.isSet)(object.query))
            obj.query = exports.RequestQuery.fromJSON(object.query);
        if ((0, helpers_1.isSet)(object.beginBlock))
            obj.beginBlock = exports.RequestBeginBlock.fromJSON(object.beginBlock);
        if ((0, helpers_1.isSet)(object.checkTx))
            obj.checkTx = exports.RequestCheckTx.fromJSON(object.checkTx);
        if ((0, helpers_1.isSet)(object.deliverTx))
            obj.deliverTx = exports.RequestDeliverTx.fromJSON(object.deliverTx);
        if ((0, helpers_1.isSet)(object.endBlock))
            obj.endBlock = exports.RequestEndBlock.fromJSON(object.endBlock);
        if ((0, helpers_1.isSet)(object.commit))
            obj.commit = exports.RequestCommit.fromJSON(object.commit);
        if ((0, helpers_1.isSet)(object.listSnapshots))
            obj.listSnapshots = exports.RequestListSnapshots.fromJSON(object.listSnapshots);
        if ((0, helpers_1.isSet)(object.offerSnapshot))
            obj.offerSnapshot = exports.RequestOfferSnapshot.fromJSON(object.offerSnapshot);
        if ((0, helpers_1.isSet)(object.loadSnapshotChunk))
            obj.loadSnapshotChunk = exports.RequestLoadSnapshotChunk.fromJSON(object.loadSnapshotChunk);
        if ((0, helpers_1.isSet)(object.applySnapshotChunk))
            obj.applySnapshotChunk = exports.RequestApplySnapshotChunk.fromJSON(object.applySnapshotChunk);
        if ((0, helpers_1.isSet)(object.prepareProposal))
            obj.prepareProposal = exports.RequestPrepareProposal.fromJSON(object.prepareProposal);
        if ((0, helpers_1.isSet)(object.processProposal))
            obj.processProposal = exports.RequestProcessProposal.fromJSON(object.processProposal);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.echo !== undefined && (obj.echo = message.echo ? exports.RequestEcho.toJSON(message.echo) : undefined);
        message.flush !== undefined &&
            (obj.flush = message.flush ? exports.RequestFlush.toJSON(message.flush) : undefined);
        message.info !== undefined && (obj.info = message.info ? exports.RequestInfo.toJSON(message.info) : undefined);
        message.initChain !== undefined &&
            (obj.initChain = message.initChain ? exports.RequestInitChain.toJSON(message.initChain) : undefined);
        message.query !== undefined &&
            (obj.query = message.query ? exports.RequestQuery.toJSON(message.query) : undefined);
        message.beginBlock !== undefined &&
            (obj.beginBlock = message.beginBlock ? exports.RequestBeginBlock.toJSON(message.beginBlock) : undefined);
        message.checkTx !== undefined &&
            (obj.checkTx = message.checkTx ? exports.RequestCheckTx.toJSON(message.checkTx) : undefined);
        message.deliverTx !== undefined &&
            (obj.deliverTx = message.deliverTx ? exports.RequestDeliverTx.toJSON(message.deliverTx) : undefined);
        message.endBlock !== undefined &&
            (obj.endBlock = message.endBlock ? exports.RequestEndBlock.toJSON(message.endBlock) : undefined);
        message.commit !== undefined &&
            (obj.commit = message.commit ? exports.RequestCommit.toJSON(message.commit) : undefined);
        message.listSnapshots !== undefined &&
            (obj.listSnapshots = message.listSnapshots
                ? exports.RequestListSnapshots.toJSON(message.listSnapshots)
                : undefined);
        message.offerSnapshot !== undefined &&
            (obj.offerSnapshot = message.offerSnapshot
                ? exports.RequestOfferSnapshot.toJSON(message.offerSnapshot)
                : undefined);
        message.loadSnapshotChunk !== undefined &&
            (obj.loadSnapshotChunk = message.loadSnapshotChunk
                ? exports.RequestLoadSnapshotChunk.toJSON(message.loadSnapshotChunk)
                : undefined);
        message.applySnapshotChunk !== undefined &&
            (obj.applySnapshotChunk = message.applySnapshotChunk
                ? exports.RequestApplySnapshotChunk.toJSON(message.applySnapshotChunk)
                : undefined);
        message.prepareProposal !== undefined &&
            (obj.prepareProposal = message.prepareProposal
                ? exports.RequestPrepareProposal.toJSON(message.prepareProposal)
                : undefined);
        message.processProposal !== undefined &&
            (obj.processProposal = message.processProposal
                ? exports.RequestProcessProposal.toJSON(message.processProposal)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseRequest();
        if (object.echo !== undefined && object.echo !== null) {
            message.echo = exports.RequestEcho.fromPartial(object.echo);
        }
        if (object.flush !== undefined && object.flush !== null) {
            message.flush = exports.RequestFlush.fromPartial(object.flush);
        }
        if (object.info !== undefined && object.info !== null) {
            message.info = exports.RequestInfo.fromPartial(object.info);
        }
        if (object.initChain !== undefined && object.initChain !== null) {
            message.initChain = exports.RequestInitChain.fromPartial(object.initChain);
        }
        if (object.query !== undefined && object.query !== null) {
            message.query = exports.RequestQuery.fromPartial(object.query);
        }
        if (object.beginBlock !== undefined && object.beginBlock !== null) {
            message.beginBlock = exports.RequestBeginBlock.fromPartial(object.beginBlock);
        }
        if (object.checkTx !== undefined && object.checkTx !== null) {
            message.checkTx = exports.RequestCheckTx.fromPartial(object.checkTx);
        }
        if (object.deliverTx !== undefined && object.deliverTx !== null) {
            message.deliverTx = exports.RequestDeliverTx.fromPartial(object.deliverTx);
        }
        if (object.endBlock !== undefined && object.endBlock !== null) {
            message.endBlock = exports.RequestEndBlock.fromPartial(object.endBlock);
        }
        if (object.commit !== undefined && object.commit !== null) {
            message.commit = exports.RequestCommit.fromPartial(object.commit);
        }
        if (object.listSnapshots !== undefined && object.listSnapshots !== null) {
            message.listSnapshots = exports.RequestListSnapshots.fromPartial(object.listSnapshots);
        }
        if (object.offerSnapshot !== undefined && object.offerSnapshot !== null) {
            message.offerSnapshot = exports.RequestOfferSnapshot.fromPartial(object.offerSnapshot);
        }
        if (object.loadSnapshotChunk !== undefined && object.loadSnapshotChunk !== null) {
            message.loadSnapshotChunk = exports.RequestLoadSnapshotChunk.fromPartial(object.loadSnapshotChunk);
        }
        if (object.applySnapshotChunk !== undefined && object.applySnapshotChunk !== null) {
            message.applySnapshotChunk = exports.RequestApplySnapshotChunk.fromPartial(object.applySnapshotChunk);
        }
        if (object.prepareProposal !== undefined && object.prepareProposal !== null) {
            message.prepareProposal = exports.RequestPrepareProposal.fromPartial(object.prepareProposal);
        }
        if (object.processProposal !== undefined && object.processProposal !== null) {
            message.processProposal = exports.RequestProcessProposal.fromPartial(object.processProposal);
        }
        return message;
    },
};
function createBaseRequestEcho() {
    return {
        message: "",
    };
}
exports.RequestEcho = {
    typeUrl: "/tendermint.abci.RequestEcho",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.message !== "") {
            writer.uint32(10).string(message.message);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseRequestEcho();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.message = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseRequestEcho();
        if ((0, helpers_1.isSet)(object.message))
            obj.message = String(object.message);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.message !== undefined && (obj.message = message.message);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseRequestEcho();
        message.message = object.message ?? "";
        return message;
    },
};
function createBaseRequestFlush() {
    return {};
}
exports.RequestFlush = {
    typeUrl: "/tendermint.abci.RequestFlush",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseRequestFlush();
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
    fromJSON(_) {
        const obj = createBaseRequestFlush();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseRequestFlush();
        return message;
    },
};
function createBaseRequestInfo() {
    return {
        version: "",
        blockVersion: BigInt(0),
        p2pVersion: BigInt(0),
        abciVersion: "",
    };
}
exports.RequestInfo = {
    typeUrl: "/tendermint.abci.RequestInfo",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.version !== "") {
            writer.uint32(10).string(message.version);
        }
        if (message.blockVersion !== BigInt(0)) {
            writer.uint32(16).uint64(message.blockVersion);
        }
        if (message.p2pVersion !== BigInt(0)) {
            writer.uint32(24).uint64(message.p2pVersion);
        }
        if (message.abciVersion !== "") {
            writer.uint32(34).string(message.abciVersion);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseRequestInfo();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.version = reader.string();
                    break;
                case 2:
                    message.blockVersion = reader.uint64();
                    break;
                case 3:
                    message.p2pVersion = reader.uint64();
                    break;
                case 4:
                    message.abciVersion = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseRequestInfo();
        if ((0, helpers_1.isSet)(object.version))
            obj.version = String(object.version);
        if ((0, helpers_1.isSet)(object.blockVersion))
            obj.blockVersion = BigInt(object.blockVersion.toString());
        if ((0, helpers_1.isSet)(object.p2pVersion))
            obj.p2pVersion = BigInt(object.p2pVersion.toString());
        if ((0, helpers_1.isSet)(object.abciVersion))
            obj.abciVersion = String(object.abciVersion);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.version !== undefined && (obj.version = message.version);
        message.blockVersion !== undefined && (obj.blockVersion = (message.blockVersion || BigInt(0)).toString());
        message.p2pVersion !== undefined && (obj.p2pVersion = (message.p2pVersion || BigInt(0)).toString());
        message.abciVersion !== undefined && (obj.abciVersion = message.abciVersion);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseRequestInfo();
        message.version = object.version ?? "";
        if (object.blockVersion !== undefined && object.blockVersion !== null) {
            message.blockVersion = BigInt(object.blockVersion.toString());
        }
        if (object.p2pVersion !== undefined && object.p2pVersion !== null) {
            message.p2pVersion = BigInt(object.p2pVersion.toString());
        }
        message.abciVersion = object.abciVersion ?? "";
        return message;
    },
};
function createBaseRequestInitChain() {
    return {
        time: timestamp_1.Timestamp.fromPartial({}),
        chainId: "",
        consensusParams: undefined,
        validators: [],
        appStateBytes: new Uint8Array(),
        initialHeight: BigInt(0),
    };
}
exports.RequestInitChain = {
    typeUrl: "/tendermint.abci.RequestInitChain",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.time !== undefined) {
            timestamp_1.Timestamp.encode(message.time, writer.uint32(10).fork()).ldelim();
        }
        if (message.chainId !== "") {
            writer.uint32(18).string(message.chainId);
        }
        if (message.consensusParams !== undefined) {
            params_1.ConsensusParams.encode(message.consensusParams, writer.uint32(26).fork()).ldelim();
        }
        for (const v of message.validators) {
            exports.ValidatorUpdate.encode(v, writer.uint32(34).fork()).ldelim();
        }
        if (message.appStateBytes.length !== 0) {
            writer.uint32(42).bytes(message.appStateBytes);
        }
        if (message.initialHeight !== BigInt(0)) {
            writer.uint32(48).int64(message.initialHeight);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseRequestInitChain();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.time = timestamp_1.Timestamp.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.chainId = reader.string();
                    break;
                case 3:
                    message.consensusParams = params_1.ConsensusParams.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.validators.push(exports.ValidatorUpdate.decode(reader, reader.uint32()));
                    break;
                case 5:
                    message.appStateBytes = reader.bytes();
                    break;
                case 6:
                    message.initialHeight = reader.int64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseRequestInitChain();
        if ((0, helpers_1.isSet)(object.time))
            obj.time = (0, helpers_1.fromJsonTimestamp)(object.time);
        if ((0, helpers_1.isSet)(object.chainId))
            obj.chainId = String(object.chainId);
        if ((0, helpers_1.isSet)(object.consensusParams))
            obj.consensusParams = params_1.ConsensusParams.fromJSON(object.consensusParams);
        if (Array.isArray(object?.validators))
            obj.validators = object.validators.map((e) => exports.ValidatorUpdate.fromJSON(e));
        if ((0, helpers_1.isSet)(object.appStateBytes))
            obj.appStateBytes = (0, helpers_1.bytesFromBase64)(object.appStateBytes);
        if ((0, helpers_1.isSet)(object.initialHeight))
            obj.initialHeight = BigInt(object.initialHeight.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.time !== undefined && (obj.time = (0, helpers_1.fromTimestamp)(message.time).toISOString());
        message.chainId !== undefined && (obj.chainId = message.chainId);
        message.consensusParams !== undefined &&
            (obj.consensusParams = message.consensusParams
                ? params_1.ConsensusParams.toJSON(message.consensusParams)
                : undefined);
        if (message.validators) {
            obj.validators = message.validators.map((e) => (e ? exports.ValidatorUpdate.toJSON(e) : undefined));
        }
        else {
            obj.validators = [];
        }
        message.appStateBytes !== undefined &&
            (obj.appStateBytes = (0, helpers_1.base64FromBytes)(message.appStateBytes !== undefined ? message.appStateBytes : new Uint8Array()));
        message.initialHeight !== undefined &&
            (obj.initialHeight = (message.initialHeight || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseRequestInitChain();
        if (object.time !== undefined && object.time !== null) {
            message.time = timestamp_1.Timestamp.fromPartial(object.time);
        }
        message.chainId = object.chainId ?? "";
        if (object.consensusParams !== undefined && object.consensusParams !== null) {
            message.consensusParams = params_1.ConsensusParams.fromPartial(object.consensusParams);
        }
        message.validators = object.validators?.map((e) => exports.ValidatorUpdate.fromPartial(e)) || [];
        message.appStateBytes = object.appStateBytes ?? new Uint8Array();
        if (object.initialHeight !== undefined && object.initialHeight !== null) {
            message.initialHeight = BigInt(object.initialHeight.toString());
        }
        return message;
    },
};
function createBaseRequestQuery() {
    return {
        data: new Uint8Array(),
        path: "",
        height: BigInt(0),
        prove: false,
    };
}
exports.RequestQuery = {
    typeUrl: "/tendermint.abci.RequestQuery",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.data.length !== 0) {
            writer.uint32(10).bytes(message.data);
        }
        if (message.path !== "") {
            writer.uint32(18).string(message.path);
        }
        if (message.height !== BigInt(0)) {
            writer.uint32(24).int64(message.height);
        }
        if (message.prove === true) {
            writer.uint32(32).bool(message.prove);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseRequestQuery();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.data = reader.bytes();
                    break;
                case 2:
                    message.path = reader.string();
                    break;
                case 3:
                    message.height = reader.int64();
                    break;
                case 4:
                    message.prove = reader.bool();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseRequestQuery();
        if ((0, helpers_1.isSet)(object.data))
            obj.data = (0, helpers_1.bytesFromBase64)(object.data);
        if ((0, helpers_1.isSet)(object.path))
            obj.path = String(object.path);
        if ((0, helpers_1.isSet)(object.height))
            obj.height = BigInt(object.height.toString());
        if ((0, helpers_1.isSet)(object.prove))
            obj.prove = Boolean(object.prove);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.data !== undefined &&
            (obj.data = (0, helpers_1.base64FromBytes)(message.data !== undefined ? message.data : new Uint8Array()));
        message.path !== undefined && (obj.path = message.path);
        message.height !== undefined && (obj.height = (message.height || BigInt(0)).toString());
        message.prove !== undefined && (obj.prove = message.prove);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseRequestQuery();
        message.data = object.data ?? new Uint8Array();
        message.path = object.path ?? "";
        if (object.height !== undefined && object.height !== null) {
            message.height = BigInt(object.height.toString());
        }
        message.prove = object.prove ?? false;
        return message;
    },
};
function createBaseRequestBeginBlock() {
    return {
        hash: new Uint8Array(),
        header: types_1.Header.fromPartial({}),
        lastCommitInfo: exports.CommitInfo.fromPartial({}),
        byzantineValidators: [],
    };
}
exports.RequestBeginBlock = {
    typeUrl: "/tendermint.abci.RequestBeginBlock",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.hash.length !== 0) {
            writer.uint32(10).bytes(message.hash);
        }
        if (message.header !== undefined) {
            types_1.Header.encode(message.header, writer.uint32(18).fork()).ldelim();
        }
        if (message.lastCommitInfo !== undefined) {
            exports.CommitInfo.encode(message.lastCommitInfo, writer.uint32(26).fork()).ldelim();
        }
        for (const v of message.byzantineValidators) {
            exports.Misbehavior.encode(v, writer.uint32(34).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseRequestBeginBlock();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.hash = reader.bytes();
                    break;
                case 2:
                    message.header = types_1.Header.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.lastCommitInfo = exports.CommitInfo.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.byzantineValidators.push(exports.Misbehavior.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseRequestBeginBlock();
        if ((0, helpers_1.isSet)(object.hash))
            obj.hash = (0, helpers_1.bytesFromBase64)(object.hash);
        if ((0, helpers_1.isSet)(object.header))
            obj.header = types_1.Header.fromJSON(object.header);
        if ((0, helpers_1.isSet)(object.lastCommitInfo))
            obj.lastCommitInfo = exports.CommitInfo.fromJSON(object.lastCommitInfo);
        if (Array.isArray(object?.byzantineValidators))
            obj.byzantineValidators = object.byzantineValidators.map((e) => exports.Misbehavior.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.hash !== undefined &&
            (obj.hash = (0, helpers_1.base64FromBytes)(message.hash !== undefined ? message.hash : new Uint8Array()));
        message.header !== undefined && (obj.header = message.header ? types_1.Header.toJSON(message.header) : undefined);
        message.lastCommitInfo !== undefined &&
            (obj.lastCommitInfo = message.lastCommitInfo ? exports.CommitInfo.toJSON(message.lastCommitInfo) : undefined);
        if (message.byzantineValidators) {
            obj.byzantineValidators = message.byzantineValidators.map((e) => e ? exports.Misbehavior.toJSON(e) : undefined);
        }
        else {
            obj.byzantineValidators = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseRequestBeginBlock();
        message.hash = object.hash ?? new Uint8Array();
        if (object.header !== undefined && object.header !== null) {
            message.header = types_1.Header.fromPartial(object.header);
        }
        if (object.lastCommitInfo !== undefined && object.lastCommitInfo !== null) {
            message.lastCommitInfo = exports.CommitInfo.fromPartial(object.lastCommitInfo);
        }
        message.byzantineValidators = object.byzantineValidators?.map((e) => exports.Misbehavior.fromPartial(e)) || [];
        return message;
    },
};
function createBaseRequestCheckTx() {
    return {
        tx: new Uint8Array(),
        type: 0,
    };
}
exports.RequestCheckTx = {
    typeUrl: "/tendermint.abci.RequestCheckTx",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.tx.length !== 0) {
            writer.uint32(10).bytes(message.tx);
        }
        if (message.type !== 0) {
            writer.uint32(16).int32(message.type);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseRequestCheckTx();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.tx = reader.bytes();
                    break;
                case 2:
                    message.type = reader.int32();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseRequestCheckTx();
        if ((0, helpers_1.isSet)(object.tx))
            obj.tx = (0, helpers_1.bytesFromBase64)(object.tx);
        if ((0, helpers_1.isSet)(object.type))
            obj.type = checkTxTypeFromJSON(object.type);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.tx !== undefined &&
            (obj.tx = (0, helpers_1.base64FromBytes)(message.tx !== undefined ? message.tx : new Uint8Array()));
        message.type !== undefined && (obj.type = checkTxTypeToJSON(message.type));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseRequestCheckTx();
        message.tx = object.tx ?? new Uint8Array();
        message.type = object.type ?? 0;
        return message;
    },
};
function createBaseRequestDeliverTx() {
    return {
        tx: new Uint8Array(),
    };
}
exports.RequestDeliverTx = {
    typeUrl: "/tendermint.abci.RequestDeliverTx",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.tx.length !== 0) {
            writer.uint32(10).bytes(message.tx);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseRequestDeliverTx();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.tx = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseRequestDeliverTx();
        if ((0, helpers_1.isSet)(object.tx))
            obj.tx = (0, helpers_1.bytesFromBase64)(object.tx);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.tx !== undefined &&
            (obj.tx = (0, helpers_1.base64FromBytes)(message.tx !== undefined ? message.tx : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseRequestDeliverTx();
        message.tx = object.tx ?? new Uint8Array();
        return message;
    },
};
function createBaseRequestEndBlock() {
    return {
        height: BigInt(0),
    };
}
exports.RequestEndBlock = {
    typeUrl: "/tendermint.abci.RequestEndBlock",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.height !== BigInt(0)) {
            writer.uint32(8).int64(message.height);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseRequestEndBlock();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.height = reader.int64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseRequestEndBlock();
        if ((0, helpers_1.isSet)(object.height))
            obj.height = BigInt(object.height.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.height !== undefined && (obj.height = (message.height || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseRequestEndBlock();
        if (object.height !== undefined && object.height !== null) {
            message.height = BigInt(object.height.toString());
        }
        return message;
    },
};
function createBaseRequestCommit() {
    return {};
}
exports.RequestCommit = {
    typeUrl: "/tendermint.abci.RequestCommit",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseRequestCommit();
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
    fromJSON(_) {
        const obj = createBaseRequestCommit();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseRequestCommit();
        return message;
    },
};
function createBaseRequestListSnapshots() {
    return {};
}
exports.RequestListSnapshots = {
    typeUrl: "/tendermint.abci.RequestListSnapshots",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseRequestListSnapshots();
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
    fromJSON(_) {
        const obj = createBaseRequestListSnapshots();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseRequestListSnapshots();
        return message;
    },
};
function createBaseRequestOfferSnapshot() {
    return {
        snapshot: undefined,
        appHash: new Uint8Array(),
    };
}
exports.RequestOfferSnapshot = {
    typeUrl: "/tendermint.abci.RequestOfferSnapshot",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.snapshot !== undefined) {
            exports.Snapshot.encode(message.snapshot, writer.uint32(10).fork()).ldelim();
        }
        if (message.appHash.length !== 0) {
            writer.uint32(18).bytes(message.appHash);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseRequestOfferSnapshot();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.snapshot = exports.Snapshot.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.appHash = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseRequestOfferSnapshot();
        if ((0, helpers_1.isSet)(object.snapshot))
            obj.snapshot = exports.Snapshot.fromJSON(object.snapshot);
        if ((0, helpers_1.isSet)(object.appHash))
            obj.appHash = (0, helpers_1.bytesFromBase64)(object.appHash);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.snapshot !== undefined &&
            (obj.snapshot = message.snapshot ? exports.Snapshot.toJSON(message.snapshot) : undefined);
        message.appHash !== undefined &&
            (obj.appHash = (0, helpers_1.base64FromBytes)(message.appHash !== undefined ? message.appHash : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseRequestOfferSnapshot();
        if (object.snapshot !== undefined && object.snapshot !== null) {
            message.snapshot = exports.Snapshot.fromPartial(object.snapshot);
        }
        message.appHash = object.appHash ?? new Uint8Array();
        return message;
    },
};
function createBaseRequestLoadSnapshotChunk() {
    return {
        height: BigInt(0),
        format: 0,
        chunk: 0,
    };
}
exports.RequestLoadSnapshotChunk = {
    typeUrl: "/tendermint.abci.RequestLoadSnapshotChunk",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.height !== BigInt(0)) {
            writer.uint32(8).uint64(message.height);
        }
        if (message.format !== 0) {
            writer.uint32(16).uint32(message.format);
        }
        if (message.chunk !== 0) {
            writer.uint32(24).uint32(message.chunk);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseRequestLoadSnapshotChunk();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.height = reader.uint64();
                    break;
                case 2:
                    message.format = reader.uint32();
                    break;
                case 3:
                    message.chunk = reader.uint32();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseRequestLoadSnapshotChunk();
        if ((0, helpers_1.isSet)(object.height))
            obj.height = BigInt(object.height.toString());
        if ((0, helpers_1.isSet)(object.format))
            obj.format = Number(object.format);
        if ((0, helpers_1.isSet)(object.chunk))
            obj.chunk = Number(object.chunk);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.height !== undefined && (obj.height = (message.height || BigInt(0)).toString());
        message.format !== undefined && (obj.format = Math.round(message.format));
        message.chunk !== undefined && (obj.chunk = Math.round(message.chunk));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseRequestLoadSnapshotChunk();
        if (object.height !== undefined && object.height !== null) {
            message.height = BigInt(object.height.toString());
        }
        message.format = object.format ?? 0;
        message.chunk = object.chunk ?? 0;
        return message;
    },
};
function createBaseRequestApplySnapshotChunk() {
    return {
        index: 0,
        chunk: new Uint8Array(),
        sender: "",
    };
}
exports.RequestApplySnapshotChunk = {
    typeUrl: "/tendermint.abci.RequestApplySnapshotChunk",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.index !== 0) {
            writer.uint32(8).uint32(message.index);
        }
        if (message.chunk.length !== 0) {
            writer.uint32(18).bytes(message.chunk);
        }
        if (message.sender !== "") {
            writer.uint32(26).string(message.sender);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseRequestApplySnapshotChunk();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.index = reader.uint32();
                    break;
                case 2:
                    message.chunk = reader.bytes();
                    break;
                case 3:
                    message.sender = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseRequestApplySnapshotChunk();
        if ((0, helpers_1.isSet)(object.index))
            obj.index = Number(object.index);
        if ((0, helpers_1.isSet)(object.chunk))
            obj.chunk = (0, helpers_1.bytesFromBase64)(object.chunk);
        if ((0, helpers_1.isSet)(object.sender))
            obj.sender = String(object.sender);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.index !== undefined && (obj.index = Math.round(message.index));
        message.chunk !== undefined &&
            (obj.chunk = (0, helpers_1.base64FromBytes)(message.chunk !== undefined ? message.chunk : new Uint8Array()));
        message.sender !== undefined && (obj.sender = message.sender);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseRequestApplySnapshotChunk();
        message.index = object.index ?? 0;
        message.chunk = object.chunk ?? new Uint8Array();
        message.sender = object.sender ?? "";
        return message;
    },
};
function createBaseRequestPrepareProposal() {
    return {
        maxTxBytes: BigInt(0),
        txs: [],
        localLastCommit: exports.ExtendedCommitInfo.fromPartial({}),
        misbehavior: [],
        height: BigInt(0),
        time: timestamp_1.Timestamp.fromPartial({}),
        nextValidatorsHash: new Uint8Array(),
        proposerAddress: new Uint8Array(),
    };
}
exports.RequestPrepareProposal = {
    typeUrl: "/tendermint.abci.RequestPrepareProposal",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.maxTxBytes !== BigInt(0)) {
            writer.uint32(8).int64(message.maxTxBytes);
        }
        for (const v of message.txs) {
            writer.uint32(18).bytes(v);
        }
        if (message.localLastCommit !== undefined) {
            exports.ExtendedCommitInfo.encode(message.localLastCommit, writer.uint32(26).fork()).ldelim();
        }
        for (const v of message.misbehavior) {
            exports.Misbehavior.encode(v, writer.uint32(34).fork()).ldelim();
        }
        if (message.height !== BigInt(0)) {
            writer.uint32(40).int64(message.height);
        }
        if (message.time !== undefined) {
            timestamp_1.Timestamp.encode(message.time, writer.uint32(50).fork()).ldelim();
        }
        if (message.nextValidatorsHash.length !== 0) {
            writer.uint32(58).bytes(message.nextValidatorsHash);
        }
        if (message.proposerAddress.length !== 0) {
            writer.uint32(66).bytes(message.proposerAddress);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseRequestPrepareProposal();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.maxTxBytes = reader.int64();
                    break;
                case 2:
                    message.txs.push(reader.bytes());
                    break;
                case 3:
                    message.localLastCommit = exports.ExtendedCommitInfo.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.misbehavior.push(exports.Misbehavior.decode(reader, reader.uint32()));
                    break;
                case 5:
                    message.height = reader.int64();
                    break;
                case 6:
                    message.time = timestamp_1.Timestamp.decode(reader, reader.uint32());
                    break;
                case 7:
                    message.nextValidatorsHash = reader.bytes();
                    break;
                case 8:
                    message.proposerAddress = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseRequestPrepareProposal();
        if ((0, helpers_1.isSet)(object.maxTxBytes))
            obj.maxTxBytes = BigInt(object.maxTxBytes.toString());
        if (Array.isArray(object?.txs))
            obj.txs = object.txs.map((e) => (0, helpers_1.bytesFromBase64)(e));
        if ((0, helpers_1.isSet)(object.localLastCommit))
            obj.localLastCommit = exports.ExtendedCommitInfo.fromJSON(object.localLastCommit);
        if (Array.isArray(object?.misbehavior))
            obj.misbehavior = object.misbehavior.map((e) => exports.Misbehavior.fromJSON(e));
        if ((0, helpers_1.isSet)(object.height))
            obj.height = BigInt(object.height.toString());
        if ((0, helpers_1.isSet)(object.time))
            obj.time = (0, helpers_1.fromJsonTimestamp)(object.time);
        if ((0, helpers_1.isSet)(object.nextValidatorsHash))
            obj.nextValidatorsHash = (0, helpers_1.bytesFromBase64)(object.nextValidatorsHash);
        if ((0, helpers_1.isSet)(object.proposerAddress))
            obj.proposerAddress = (0, helpers_1.bytesFromBase64)(object.proposerAddress);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.maxTxBytes !== undefined && (obj.maxTxBytes = (message.maxTxBytes || BigInt(0)).toString());
        if (message.txs) {
            obj.txs = message.txs.map((e) => (0, helpers_1.base64FromBytes)(e !== undefined ? e : new Uint8Array()));
        }
        else {
            obj.txs = [];
        }
        message.localLastCommit !== undefined &&
            (obj.localLastCommit = message.localLastCommit
                ? exports.ExtendedCommitInfo.toJSON(message.localLastCommit)
                : undefined);
        if (message.misbehavior) {
            obj.misbehavior = message.misbehavior.map((e) => (e ? exports.Misbehavior.toJSON(e) : undefined));
        }
        else {
            obj.misbehavior = [];
        }
        message.height !== undefined && (obj.height = (message.height || BigInt(0)).toString());
        message.time !== undefined && (obj.time = (0, helpers_1.fromTimestamp)(message.time).toISOString());
        message.nextValidatorsHash !== undefined &&
            (obj.nextValidatorsHash = (0, helpers_1.base64FromBytes)(message.nextValidatorsHash !== undefined ? message.nextValidatorsHash : new Uint8Array()));
        message.proposerAddress !== undefined &&
            (obj.proposerAddress = (0, helpers_1.base64FromBytes)(message.proposerAddress !== undefined ? message.proposerAddress : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseRequestPrepareProposal();
        if (object.maxTxBytes !== undefined && object.maxTxBytes !== null) {
            message.maxTxBytes = BigInt(object.maxTxBytes.toString());
        }
        message.txs = object.txs?.map((e) => e) || [];
        if (object.localLastCommit !== undefined && object.localLastCommit !== null) {
            message.localLastCommit = exports.ExtendedCommitInfo.fromPartial(object.localLastCommit);
        }
        message.misbehavior = object.misbehavior?.map((e) => exports.Misbehavior.fromPartial(e)) || [];
        if (object.height !== undefined && object.height !== null) {
            message.height = BigInt(object.height.toString());
        }
        if (object.time !== undefined && object.time !== null) {
            message.time = timestamp_1.Timestamp.fromPartial(object.time);
        }
        message.nextValidatorsHash = object.nextValidatorsHash ?? new Uint8Array();
        message.proposerAddress = object.proposerAddress ?? new Uint8Array();
        return message;
    },
};
function createBaseRequestProcessProposal() {
    return {
        txs: [],
        proposedLastCommit: exports.CommitInfo.fromPartial({}),
        misbehavior: [],
        hash: new Uint8Array(),
        height: BigInt(0),
        time: timestamp_1.Timestamp.fromPartial({}),
        nextValidatorsHash: new Uint8Array(),
        proposerAddress: new Uint8Array(),
    };
}
exports.RequestProcessProposal = {
    typeUrl: "/tendermint.abci.RequestProcessProposal",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.txs) {
            writer.uint32(10).bytes(v);
        }
        if (message.proposedLastCommit !== undefined) {
            exports.CommitInfo.encode(message.proposedLastCommit, writer.uint32(18).fork()).ldelim();
        }
        for (const v of message.misbehavior) {
            exports.Misbehavior.encode(v, writer.uint32(26).fork()).ldelim();
        }
        if (message.hash.length !== 0) {
            writer.uint32(34).bytes(message.hash);
        }
        if (message.height !== BigInt(0)) {
            writer.uint32(40).int64(message.height);
        }
        if (message.time !== undefined) {
            timestamp_1.Timestamp.encode(message.time, writer.uint32(50).fork()).ldelim();
        }
        if (message.nextValidatorsHash.length !== 0) {
            writer.uint32(58).bytes(message.nextValidatorsHash);
        }
        if (message.proposerAddress.length !== 0) {
            writer.uint32(66).bytes(message.proposerAddress);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseRequestProcessProposal();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.txs.push(reader.bytes());
                    break;
                case 2:
                    message.proposedLastCommit = exports.CommitInfo.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.misbehavior.push(exports.Misbehavior.decode(reader, reader.uint32()));
                    break;
                case 4:
                    message.hash = reader.bytes();
                    break;
                case 5:
                    message.height = reader.int64();
                    break;
                case 6:
                    message.time = timestamp_1.Timestamp.decode(reader, reader.uint32());
                    break;
                case 7:
                    message.nextValidatorsHash = reader.bytes();
                    break;
                case 8:
                    message.proposerAddress = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseRequestProcessProposal();
        if (Array.isArray(object?.txs))
            obj.txs = object.txs.map((e) => (0, helpers_1.bytesFromBase64)(e));
        if ((0, helpers_1.isSet)(object.proposedLastCommit))
            obj.proposedLastCommit = exports.CommitInfo.fromJSON(object.proposedLastCommit);
        if (Array.isArray(object?.misbehavior))
            obj.misbehavior = object.misbehavior.map((e) => exports.Misbehavior.fromJSON(e));
        if ((0, helpers_1.isSet)(object.hash))
            obj.hash = (0, helpers_1.bytesFromBase64)(object.hash);
        if ((0, helpers_1.isSet)(object.height))
            obj.height = BigInt(object.height.toString());
        if ((0, helpers_1.isSet)(object.time))
            obj.time = (0, helpers_1.fromJsonTimestamp)(object.time);
        if ((0, helpers_1.isSet)(object.nextValidatorsHash))
            obj.nextValidatorsHash = (0, helpers_1.bytesFromBase64)(object.nextValidatorsHash);
        if ((0, helpers_1.isSet)(object.proposerAddress))
            obj.proposerAddress = (0, helpers_1.bytesFromBase64)(object.proposerAddress);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.txs) {
            obj.txs = message.txs.map((e) => (0, helpers_1.base64FromBytes)(e !== undefined ? e : new Uint8Array()));
        }
        else {
            obj.txs = [];
        }
        message.proposedLastCommit !== undefined &&
            (obj.proposedLastCommit = message.proposedLastCommit
                ? exports.CommitInfo.toJSON(message.proposedLastCommit)
                : undefined);
        if (message.misbehavior) {
            obj.misbehavior = message.misbehavior.map((e) => (e ? exports.Misbehavior.toJSON(e) : undefined));
        }
        else {
            obj.misbehavior = [];
        }
        message.hash !== undefined &&
            (obj.hash = (0, helpers_1.base64FromBytes)(message.hash !== undefined ? message.hash : new Uint8Array()));
        message.height !== undefined && (obj.height = (message.height || BigInt(0)).toString());
        message.time !== undefined && (obj.time = (0, helpers_1.fromTimestamp)(message.time).toISOString());
        message.nextValidatorsHash !== undefined &&
            (obj.nextValidatorsHash = (0, helpers_1.base64FromBytes)(message.nextValidatorsHash !== undefined ? message.nextValidatorsHash : new Uint8Array()));
        message.proposerAddress !== undefined &&
            (obj.proposerAddress = (0, helpers_1.base64FromBytes)(message.proposerAddress !== undefined ? message.proposerAddress : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseRequestProcessProposal();
        message.txs = object.txs?.map((e) => e) || [];
        if (object.proposedLastCommit !== undefined && object.proposedLastCommit !== null) {
            message.proposedLastCommit = exports.CommitInfo.fromPartial(object.proposedLastCommit);
        }
        message.misbehavior = object.misbehavior?.map((e) => exports.Misbehavior.fromPartial(e)) || [];
        message.hash = object.hash ?? new Uint8Array();
        if (object.height !== undefined && object.height !== null) {
            message.height = BigInt(object.height.toString());
        }
        if (object.time !== undefined && object.time !== null) {
            message.time = timestamp_1.Timestamp.fromPartial(object.time);
        }
        message.nextValidatorsHash = object.nextValidatorsHash ?? new Uint8Array();
        message.proposerAddress = object.proposerAddress ?? new Uint8Array();
        return message;
    },
};
function createBaseResponse() {
    return {
        exception: undefined,
        echo: undefined,
        flush: undefined,
        info: undefined,
        initChain: undefined,
        query: undefined,
        beginBlock: undefined,
        checkTx: undefined,
        deliverTx: undefined,
        endBlock: undefined,
        commit: undefined,
        listSnapshots: undefined,
        offerSnapshot: undefined,
        loadSnapshotChunk: undefined,
        applySnapshotChunk: undefined,
        prepareProposal: undefined,
        processProposal: undefined,
    };
}
exports.Response = {
    typeUrl: "/tendermint.abci.Response",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.exception !== undefined) {
            exports.ResponseException.encode(message.exception, writer.uint32(10).fork()).ldelim();
        }
        if (message.echo !== undefined) {
            exports.ResponseEcho.encode(message.echo, writer.uint32(18).fork()).ldelim();
        }
        if (message.flush !== undefined) {
            exports.ResponseFlush.encode(message.flush, writer.uint32(26).fork()).ldelim();
        }
        if (message.info !== undefined) {
            exports.ResponseInfo.encode(message.info, writer.uint32(34).fork()).ldelim();
        }
        if (message.initChain !== undefined) {
            exports.ResponseInitChain.encode(message.initChain, writer.uint32(50).fork()).ldelim();
        }
        if (message.query !== undefined) {
            exports.ResponseQuery.encode(message.query, writer.uint32(58).fork()).ldelim();
        }
        if (message.beginBlock !== undefined) {
            exports.ResponseBeginBlock.encode(message.beginBlock, writer.uint32(66).fork()).ldelim();
        }
        if (message.checkTx !== undefined) {
            exports.ResponseCheckTx.encode(message.checkTx, writer.uint32(74).fork()).ldelim();
        }
        if (message.deliverTx !== undefined) {
            exports.ResponseDeliverTx.encode(message.deliverTx, writer.uint32(82).fork()).ldelim();
        }
        if (message.endBlock !== undefined) {
            exports.ResponseEndBlock.encode(message.endBlock, writer.uint32(90).fork()).ldelim();
        }
        if (message.commit !== undefined) {
            exports.ResponseCommit.encode(message.commit, writer.uint32(98).fork()).ldelim();
        }
        if (message.listSnapshots !== undefined) {
            exports.ResponseListSnapshots.encode(message.listSnapshots, writer.uint32(106).fork()).ldelim();
        }
        if (message.offerSnapshot !== undefined) {
            exports.ResponseOfferSnapshot.encode(message.offerSnapshot, writer.uint32(114).fork()).ldelim();
        }
        if (message.loadSnapshotChunk !== undefined) {
            exports.ResponseLoadSnapshotChunk.encode(message.loadSnapshotChunk, writer.uint32(122).fork()).ldelim();
        }
        if (message.applySnapshotChunk !== undefined) {
            exports.ResponseApplySnapshotChunk.encode(message.applySnapshotChunk, writer.uint32(130).fork()).ldelim();
        }
        if (message.prepareProposal !== undefined) {
            exports.ResponsePrepareProposal.encode(message.prepareProposal, writer.uint32(138).fork()).ldelim();
        }
        if (message.processProposal !== undefined) {
            exports.ResponseProcessProposal.encode(message.processProposal, writer.uint32(146).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.exception = exports.ResponseException.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.echo = exports.ResponseEcho.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.flush = exports.ResponseFlush.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.info = exports.ResponseInfo.decode(reader, reader.uint32());
                    break;
                case 6:
                    message.initChain = exports.ResponseInitChain.decode(reader, reader.uint32());
                    break;
                case 7:
                    message.query = exports.ResponseQuery.decode(reader, reader.uint32());
                    break;
                case 8:
                    message.beginBlock = exports.ResponseBeginBlock.decode(reader, reader.uint32());
                    break;
                case 9:
                    message.checkTx = exports.ResponseCheckTx.decode(reader, reader.uint32());
                    break;
                case 10:
                    message.deliverTx = exports.ResponseDeliverTx.decode(reader, reader.uint32());
                    break;
                case 11:
                    message.endBlock = exports.ResponseEndBlock.decode(reader, reader.uint32());
                    break;
                case 12:
                    message.commit = exports.ResponseCommit.decode(reader, reader.uint32());
                    break;
                case 13:
                    message.listSnapshots = exports.ResponseListSnapshots.decode(reader, reader.uint32());
                    break;
                case 14:
                    message.offerSnapshot = exports.ResponseOfferSnapshot.decode(reader, reader.uint32());
                    break;
                case 15:
                    message.loadSnapshotChunk = exports.ResponseLoadSnapshotChunk.decode(reader, reader.uint32());
                    break;
                case 16:
                    message.applySnapshotChunk = exports.ResponseApplySnapshotChunk.decode(reader, reader.uint32());
                    break;
                case 17:
                    message.prepareProposal = exports.ResponsePrepareProposal.decode(reader, reader.uint32());
                    break;
                case 18:
                    message.processProposal = exports.ResponseProcessProposal.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseResponse();
        if ((0, helpers_1.isSet)(object.exception))
            obj.exception = exports.ResponseException.fromJSON(object.exception);
        if ((0, helpers_1.isSet)(object.echo))
            obj.echo = exports.ResponseEcho.fromJSON(object.echo);
        if ((0, helpers_1.isSet)(object.flush))
            obj.flush = exports.ResponseFlush.fromJSON(object.flush);
        if ((0, helpers_1.isSet)(object.info))
            obj.info = exports.ResponseInfo.fromJSON(object.info);
        if ((0, helpers_1.isSet)(object.initChain))
            obj.initChain = exports.ResponseInitChain.fromJSON(object.initChain);
        if ((0, helpers_1.isSet)(object.query))
            obj.query = exports.ResponseQuery.fromJSON(object.query);
        if ((0, helpers_1.isSet)(object.beginBlock))
            obj.beginBlock = exports.ResponseBeginBlock.fromJSON(object.beginBlock);
        if ((0, helpers_1.isSet)(object.checkTx))
            obj.checkTx = exports.ResponseCheckTx.fromJSON(object.checkTx);
        if ((0, helpers_1.isSet)(object.deliverTx))
            obj.deliverTx = exports.ResponseDeliverTx.fromJSON(object.deliverTx);
        if ((0, helpers_1.isSet)(object.endBlock))
            obj.endBlock = exports.ResponseEndBlock.fromJSON(object.endBlock);
        if ((0, helpers_1.isSet)(object.commit))
            obj.commit = exports.ResponseCommit.fromJSON(object.commit);
        if ((0, helpers_1.isSet)(object.listSnapshots))
            obj.listSnapshots = exports.ResponseListSnapshots.fromJSON(object.listSnapshots);
        if ((0, helpers_1.isSet)(object.offerSnapshot))
            obj.offerSnapshot = exports.ResponseOfferSnapshot.fromJSON(object.offerSnapshot);
        if ((0, helpers_1.isSet)(object.loadSnapshotChunk))
            obj.loadSnapshotChunk = exports.ResponseLoadSnapshotChunk.fromJSON(object.loadSnapshotChunk);
        if ((0, helpers_1.isSet)(object.applySnapshotChunk))
            obj.applySnapshotChunk = exports.ResponseApplySnapshotChunk.fromJSON(object.applySnapshotChunk);
        if ((0, helpers_1.isSet)(object.prepareProposal))
            obj.prepareProposal = exports.ResponsePrepareProposal.fromJSON(object.prepareProposal);
        if ((0, helpers_1.isSet)(object.processProposal))
            obj.processProposal = exports.ResponseProcessProposal.fromJSON(object.processProposal);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.exception !== undefined &&
            (obj.exception = message.exception ? exports.ResponseException.toJSON(message.exception) : undefined);
        message.echo !== undefined && (obj.echo = message.echo ? exports.ResponseEcho.toJSON(message.echo) : undefined);
        message.flush !== undefined &&
            (obj.flush = message.flush ? exports.ResponseFlush.toJSON(message.flush) : undefined);
        message.info !== undefined && (obj.info = message.info ? exports.ResponseInfo.toJSON(message.info) : undefined);
        message.initChain !== undefined &&
            (obj.initChain = message.initChain ? exports.ResponseInitChain.toJSON(message.initChain) : undefined);
        message.query !== undefined &&
            (obj.query = message.query ? exports.ResponseQuery.toJSON(message.query) : undefined);
        message.beginBlock !== undefined &&
            (obj.beginBlock = message.beginBlock ? exports.ResponseBeginBlock.toJSON(message.beginBlock) : undefined);
        message.checkTx !== undefined &&
            (obj.checkTx = message.checkTx ? exports.ResponseCheckTx.toJSON(message.checkTx) : undefined);
        message.deliverTx !== undefined &&
            (obj.deliverTx = message.deliverTx ? exports.ResponseDeliverTx.toJSON(message.deliverTx) : undefined);
        message.endBlock !== undefined &&
            (obj.endBlock = message.endBlock ? exports.ResponseEndBlock.toJSON(message.endBlock) : undefined);
        message.commit !== undefined &&
            (obj.commit = message.commit ? exports.ResponseCommit.toJSON(message.commit) : undefined);
        message.listSnapshots !== undefined &&
            (obj.listSnapshots = message.listSnapshots
                ? exports.ResponseListSnapshots.toJSON(message.listSnapshots)
                : undefined);
        message.offerSnapshot !== undefined &&
            (obj.offerSnapshot = message.offerSnapshot
                ? exports.ResponseOfferSnapshot.toJSON(message.offerSnapshot)
                : undefined);
        message.loadSnapshotChunk !== undefined &&
            (obj.loadSnapshotChunk = message.loadSnapshotChunk
                ? exports.ResponseLoadSnapshotChunk.toJSON(message.loadSnapshotChunk)
                : undefined);
        message.applySnapshotChunk !== undefined &&
            (obj.applySnapshotChunk = message.applySnapshotChunk
                ? exports.ResponseApplySnapshotChunk.toJSON(message.applySnapshotChunk)
                : undefined);
        message.prepareProposal !== undefined &&
            (obj.prepareProposal = message.prepareProposal
                ? exports.ResponsePrepareProposal.toJSON(message.prepareProposal)
                : undefined);
        message.processProposal !== undefined &&
            (obj.processProposal = message.processProposal
                ? exports.ResponseProcessProposal.toJSON(message.processProposal)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseResponse();
        if (object.exception !== undefined && object.exception !== null) {
            message.exception = exports.ResponseException.fromPartial(object.exception);
        }
        if (object.echo !== undefined && object.echo !== null) {
            message.echo = exports.ResponseEcho.fromPartial(object.echo);
        }
        if (object.flush !== undefined && object.flush !== null) {
            message.flush = exports.ResponseFlush.fromPartial(object.flush);
        }
        if (object.info !== undefined && object.info !== null) {
            message.info = exports.ResponseInfo.fromPartial(object.info);
        }
        if (object.initChain !== undefined && object.initChain !== null) {
            message.initChain = exports.ResponseInitChain.fromPartial(object.initChain);
        }
        if (object.query !== undefined && object.query !== null) {
            message.query = exports.ResponseQuery.fromPartial(object.query);
        }
        if (object.beginBlock !== undefined && object.beginBlock !== null) {
            message.beginBlock = exports.ResponseBeginBlock.fromPartial(object.beginBlock);
        }
        if (object.checkTx !== undefined && object.checkTx !== null) {
            message.checkTx = exports.ResponseCheckTx.fromPartial(object.checkTx);
        }
        if (object.deliverTx !== undefined && object.deliverTx !== null) {
            message.deliverTx = exports.ResponseDeliverTx.fromPartial(object.deliverTx);
        }
        if (object.endBlock !== undefined && object.endBlock !== null) {
            message.endBlock = exports.ResponseEndBlock.fromPartial(object.endBlock);
        }
        if (object.commit !== undefined && object.commit !== null) {
            message.commit = exports.ResponseCommit.fromPartial(object.commit);
        }
        if (object.listSnapshots !== undefined && object.listSnapshots !== null) {
            message.listSnapshots = exports.ResponseListSnapshots.fromPartial(object.listSnapshots);
        }
        if (object.offerSnapshot !== undefined && object.offerSnapshot !== null) {
            message.offerSnapshot = exports.ResponseOfferSnapshot.fromPartial(object.offerSnapshot);
        }
        if (object.loadSnapshotChunk !== undefined && object.loadSnapshotChunk !== null) {
            message.loadSnapshotChunk = exports.ResponseLoadSnapshotChunk.fromPartial(object.loadSnapshotChunk);
        }
        if (object.applySnapshotChunk !== undefined && object.applySnapshotChunk !== null) {
            message.applySnapshotChunk = exports.ResponseApplySnapshotChunk.fromPartial(object.applySnapshotChunk);
        }
        if (object.prepareProposal !== undefined && object.prepareProposal !== null) {
            message.prepareProposal = exports.ResponsePrepareProposal.fromPartial(object.prepareProposal);
        }
        if (object.processProposal !== undefined && object.processProposal !== null) {
            message.processProposal = exports.ResponseProcessProposal.fromPartial(object.processProposal);
        }
        return message;
    },
};
function createBaseResponseException() {
    return {
        error: "",
    };
}
exports.ResponseException = {
    typeUrl: "/tendermint.abci.ResponseException",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.error !== "") {
            writer.uint32(10).string(message.error);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseResponseException();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.error = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseResponseException();
        if ((0, helpers_1.isSet)(object.error))
            obj.error = String(object.error);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.error !== undefined && (obj.error = message.error);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseResponseException();
        message.error = object.error ?? "";
        return message;
    },
};
function createBaseResponseEcho() {
    return {
        message: "",
    };
}
exports.ResponseEcho = {
    typeUrl: "/tendermint.abci.ResponseEcho",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.message !== "") {
            writer.uint32(10).string(message.message);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseResponseEcho();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.message = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseResponseEcho();
        if ((0, helpers_1.isSet)(object.message))
            obj.message = String(object.message);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.message !== undefined && (obj.message = message.message);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseResponseEcho();
        message.message = object.message ?? "";
        return message;
    },
};
function createBaseResponseFlush() {
    return {};
}
exports.ResponseFlush = {
    typeUrl: "/tendermint.abci.ResponseFlush",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseResponseFlush();
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
    fromJSON(_) {
        const obj = createBaseResponseFlush();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseResponseFlush();
        return message;
    },
};
function createBaseResponseInfo() {
    return {
        data: "",
        version: "",
        appVersion: BigInt(0),
        lastBlockHeight: BigInt(0),
        lastBlockAppHash: new Uint8Array(),
    };
}
exports.ResponseInfo = {
    typeUrl: "/tendermint.abci.ResponseInfo",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.data !== "") {
            writer.uint32(10).string(message.data);
        }
        if (message.version !== "") {
            writer.uint32(18).string(message.version);
        }
        if (message.appVersion !== BigInt(0)) {
            writer.uint32(24).uint64(message.appVersion);
        }
        if (message.lastBlockHeight !== BigInt(0)) {
            writer.uint32(32).int64(message.lastBlockHeight);
        }
        if (message.lastBlockAppHash.length !== 0) {
            writer.uint32(42).bytes(message.lastBlockAppHash);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseResponseInfo();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.data = reader.string();
                    break;
                case 2:
                    message.version = reader.string();
                    break;
                case 3:
                    message.appVersion = reader.uint64();
                    break;
                case 4:
                    message.lastBlockHeight = reader.int64();
                    break;
                case 5:
                    message.lastBlockAppHash = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseResponseInfo();
        if ((0, helpers_1.isSet)(object.data))
            obj.data = String(object.data);
        if ((0, helpers_1.isSet)(object.version))
            obj.version = String(object.version);
        if ((0, helpers_1.isSet)(object.appVersion))
            obj.appVersion = BigInt(object.appVersion.toString());
        if ((0, helpers_1.isSet)(object.lastBlockHeight))
            obj.lastBlockHeight = BigInt(object.lastBlockHeight.toString());
        if ((0, helpers_1.isSet)(object.lastBlockAppHash))
            obj.lastBlockAppHash = (0, helpers_1.bytesFromBase64)(object.lastBlockAppHash);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.data !== undefined && (obj.data = message.data);
        message.version !== undefined && (obj.version = message.version);
        message.appVersion !== undefined && (obj.appVersion = (message.appVersion || BigInt(0)).toString());
        message.lastBlockHeight !== undefined &&
            (obj.lastBlockHeight = (message.lastBlockHeight || BigInt(0)).toString());
        message.lastBlockAppHash !== undefined &&
            (obj.lastBlockAppHash = (0, helpers_1.base64FromBytes)(message.lastBlockAppHash !== undefined ? message.lastBlockAppHash : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseResponseInfo();
        message.data = object.data ?? "";
        message.version = object.version ?? "";
        if (object.appVersion !== undefined && object.appVersion !== null) {
            message.appVersion = BigInt(object.appVersion.toString());
        }
        if (object.lastBlockHeight !== undefined && object.lastBlockHeight !== null) {
            message.lastBlockHeight = BigInt(object.lastBlockHeight.toString());
        }
        message.lastBlockAppHash = object.lastBlockAppHash ?? new Uint8Array();
        return message;
    },
};
function createBaseResponseInitChain() {
    return {
        consensusParams: undefined,
        validators: [],
        appHash: new Uint8Array(),
    };
}
exports.ResponseInitChain = {
    typeUrl: "/tendermint.abci.ResponseInitChain",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.consensusParams !== undefined) {
            params_1.ConsensusParams.encode(message.consensusParams, writer.uint32(10).fork()).ldelim();
        }
        for (const v of message.validators) {
            exports.ValidatorUpdate.encode(v, writer.uint32(18).fork()).ldelim();
        }
        if (message.appHash.length !== 0) {
            writer.uint32(26).bytes(message.appHash);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseResponseInitChain();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.consensusParams = params_1.ConsensusParams.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.validators.push(exports.ValidatorUpdate.decode(reader, reader.uint32()));
                    break;
                case 3:
                    message.appHash = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseResponseInitChain();
        if ((0, helpers_1.isSet)(object.consensusParams))
            obj.consensusParams = params_1.ConsensusParams.fromJSON(object.consensusParams);
        if (Array.isArray(object?.validators))
            obj.validators = object.validators.map((e) => exports.ValidatorUpdate.fromJSON(e));
        if ((0, helpers_1.isSet)(object.appHash))
            obj.appHash = (0, helpers_1.bytesFromBase64)(object.appHash);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.consensusParams !== undefined &&
            (obj.consensusParams = message.consensusParams
                ? params_1.ConsensusParams.toJSON(message.consensusParams)
                : undefined);
        if (message.validators) {
            obj.validators = message.validators.map((e) => (e ? exports.ValidatorUpdate.toJSON(e) : undefined));
        }
        else {
            obj.validators = [];
        }
        message.appHash !== undefined &&
            (obj.appHash = (0, helpers_1.base64FromBytes)(message.appHash !== undefined ? message.appHash : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseResponseInitChain();
        if (object.consensusParams !== undefined && object.consensusParams !== null) {
            message.consensusParams = params_1.ConsensusParams.fromPartial(object.consensusParams);
        }
        message.validators = object.validators?.map((e) => exports.ValidatorUpdate.fromPartial(e)) || [];
        message.appHash = object.appHash ?? new Uint8Array();
        return message;
    },
};
function createBaseResponseQuery() {
    return {
        code: 0,
        log: "",
        info: "",
        index: BigInt(0),
        key: new Uint8Array(),
        value: new Uint8Array(),
        proofOps: undefined,
        height: BigInt(0),
        codespace: "",
    };
}
exports.ResponseQuery = {
    typeUrl: "/tendermint.abci.ResponseQuery",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.code !== 0) {
            writer.uint32(8).uint32(message.code);
        }
        if (message.log !== "") {
            writer.uint32(26).string(message.log);
        }
        if (message.info !== "") {
            writer.uint32(34).string(message.info);
        }
        if (message.index !== BigInt(0)) {
            writer.uint32(40).int64(message.index);
        }
        if (message.key.length !== 0) {
            writer.uint32(50).bytes(message.key);
        }
        if (message.value.length !== 0) {
            writer.uint32(58).bytes(message.value);
        }
        if (message.proofOps !== undefined) {
            proof_1.ProofOps.encode(message.proofOps, writer.uint32(66).fork()).ldelim();
        }
        if (message.height !== BigInt(0)) {
            writer.uint32(72).int64(message.height);
        }
        if (message.codespace !== "") {
            writer.uint32(82).string(message.codespace);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseResponseQuery();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.code = reader.uint32();
                    break;
                case 3:
                    message.log = reader.string();
                    break;
                case 4:
                    message.info = reader.string();
                    break;
                case 5:
                    message.index = reader.int64();
                    break;
                case 6:
                    message.key = reader.bytes();
                    break;
                case 7:
                    message.value = reader.bytes();
                    break;
                case 8:
                    message.proofOps = proof_1.ProofOps.decode(reader, reader.uint32());
                    break;
                case 9:
                    message.height = reader.int64();
                    break;
                case 10:
                    message.codespace = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseResponseQuery();
        if ((0, helpers_1.isSet)(object.code))
            obj.code = Number(object.code);
        if ((0, helpers_1.isSet)(object.log))
            obj.log = String(object.log);
        if ((0, helpers_1.isSet)(object.info))
            obj.info = String(object.info);
        if ((0, helpers_1.isSet)(object.index))
            obj.index = BigInt(object.index.toString());
        if ((0, helpers_1.isSet)(object.key))
            obj.key = (0, helpers_1.bytesFromBase64)(object.key);
        if ((0, helpers_1.isSet)(object.value))
            obj.value = (0, helpers_1.bytesFromBase64)(object.value);
        if ((0, helpers_1.isSet)(object.proofOps))
            obj.proofOps = proof_1.ProofOps.fromJSON(object.proofOps);
        if ((0, helpers_1.isSet)(object.height))
            obj.height = BigInt(object.height.toString());
        if ((0, helpers_1.isSet)(object.codespace))
            obj.codespace = String(object.codespace);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.code !== undefined && (obj.code = Math.round(message.code));
        message.log !== undefined && (obj.log = message.log);
        message.info !== undefined && (obj.info = message.info);
        message.index !== undefined && (obj.index = (message.index || BigInt(0)).toString());
        message.key !== undefined &&
            (obj.key = (0, helpers_1.base64FromBytes)(message.key !== undefined ? message.key : new Uint8Array()));
        message.value !== undefined &&
            (obj.value = (0, helpers_1.base64FromBytes)(message.value !== undefined ? message.value : new Uint8Array()));
        message.proofOps !== undefined &&
            (obj.proofOps = message.proofOps ? proof_1.ProofOps.toJSON(message.proofOps) : undefined);
        message.height !== undefined && (obj.height = (message.height || BigInt(0)).toString());
        message.codespace !== undefined && (obj.codespace = message.codespace);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseResponseQuery();
        message.code = object.code ?? 0;
        message.log = object.log ?? "";
        message.info = object.info ?? "";
        if (object.index !== undefined && object.index !== null) {
            message.index = BigInt(object.index.toString());
        }
        message.key = object.key ?? new Uint8Array();
        message.value = object.value ?? new Uint8Array();
        if (object.proofOps !== undefined && object.proofOps !== null) {
            message.proofOps = proof_1.ProofOps.fromPartial(object.proofOps);
        }
        if (object.height !== undefined && object.height !== null) {
            message.height = BigInt(object.height.toString());
        }
        message.codespace = object.codespace ?? "";
        return message;
    },
};
function createBaseResponseBeginBlock() {
    return {
        events: [],
    };
}
exports.ResponseBeginBlock = {
    typeUrl: "/tendermint.abci.ResponseBeginBlock",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.events) {
            exports.Event.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseResponseBeginBlock();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.events.push(exports.Event.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseResponseBeginBlock();
        if (Array.isArray(object?.events))
            obj.events = object.events.map((e) => exports.Event.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.events) {
            obj.events = message.events.map((e) => (e ? exports.Event.toJSON(e) : undefined));
        }
        else {
            obj.events = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseResponseBeginBlock();
        message.events = object.events?.map((e) => exports.Event.fromPartial(e)) || [];
        return message;
    },
};
function createBaseResponseCheckTx() {
    return {
        code: 0,
        data: new Uint8Array(),
        log: "",
        info: "",
        gasWanted: BigInt(0),
        gasUsed: BigInt(0),
        events: [],
        codespace: "",
        sender: "",
        priority: BigInt(0),
        mempoolError: "",
    };
}
exports.ResponseCheckTx = {
    typeUrl: "/tendermint.abci.ResponseCheckTx",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.code !== 0) {
            writer.uint32(8).uint32(message.code);
        }
        if (message.data.length !== 0) {
            writer.uint32(18).bytes(message.data);
        }
        if (message.log !== "") {
            writer.uint32(26).string(message.log);
        }
        if (message.info !== "") {
            writer.uint32(34).string(message.info);
        }
        if (message.gasWanted !== BigInt(0)) {
            writer.uint32(40).int64(message.gasWanted);
        }
        if (message.gasUsed !== BigInt(0)) {
            writer.uint32(48).int64(message.gasUsed);
        }
        for (const v of message.events) {
            exports.Event.encode(v, writer.uint32(58).fork()).ldelim();
        }
        if (message.codespace !== "") {
            writer.uint32(66).string(message.codespace);
        }
        if (message.sender !== "") {
            writer.uint32(74).string(message.sender);
        }
        if (message.priority !== BigInt(0)) {
            writer.uint32(80).int64(message.priority);
        }
        if (message.mempoolError !== "") {
            writer.uint32(90).string(message.mempoolError);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseResponseCheckTx();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.code = reader.uint32();
                    break;
                case 2:
                    message.data = reader.bytes();
                    break;
                case 3:
                    message.log = reader.string();
                    break;
                case 4:
                    message.info = reader.string();
                    break;
                case 5:
                    message.gasWanted = reader.int64();
                    break;
                case 6:
                    message.gasUsed = reader.int64();
                    break;
                case 7:
                    message.events.push(exports.Event.decode(reader, reader.uint32()));
                    break;
                case 8:
                    message.codespace = reader.string();
                    break;
                case 9:
                    message.sender = reader.string();
                    break;
                case 10:
                    message.priority = reader.int64();
                    break;
                case 11:
                    message.mempoolError = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseResponseCheckTx();
        if ((0, helpers_1.isSet)(object.code))
            obj.code = Number(object.code);
        if ((0, helpers_1.isSet)(object.data))
            obj.data = (0, helpers_1.bytesFromBase64)(object.data);
        if ((0, helpers_1.isSet)(object.log))
            obj.log = String(object.log);
        if ((0, helpers_1.isSet)(object.info))
            obj.info = String(object.info);
        if ((0, helpers_1.isSet)(object.gas_wanted))
            obj.gasWanted = BigInt(object.gas_wanted.toString());
        if ((0, helpers_1.isSet)(object.gas_used))
            obj.gasUsed = BigInt(object.gas_used.toString());
        if (Array.isArray(object?.events))
            obj.events = object.events.map((e) => exports.Event.fromJSON(e));
        if ((0, helpers_1.isSet)(object.codespace))
            obj.codespace = String(object.codespace);
        if ((0, helpers_1.isSet)(object.sender))
            obj.sender = String(object.sender);
        if ((0, helpers_1.isSet)(object.priority))
            obj.priority = BigInt(object.priority.toString());
        if ((0, helpers_1.isSet)(object.mempoolError))
            obj.mempoolError = String(object.mempoolError);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.code !== undefined && (obj.code = Math.round(message.code));
        message.data !== undefined &&
            (obj.data = (0, helpers_1.base64FromBytes)(message.data !== undefined ? message.data : new Uint8Array()));
        message.log !== undefined && (obj.log = message.log);
        message.info !== undefined && (obj.info = message.info);
        message.gasWanted !== undefined && (obj.gas_wanted = (message.gasWanted || BigInt(0)).toString());
        message.gasUsed !== undefined && (obj.gas_used = (message.gasUsed || BigInt(0)).toString());
        if (message.events) {
            obj.events = message.events.map((e) => (e ? exports.Event.toJSON(e) : undefined));
        }
        else {
            obj.events = [];
        }
        message.codespace !== undefined && (obj.codespace = message.codespace);
        message.sender !== undefined && (obj.sender = message.sender);
        message.priority !== undefined && (obj.priority = (message.priority || BigInt(0)).toString());
        message.mempoolError !== undefined && (obj.mempoolError = message.mempoolError);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseResponseCheckTx();
        message.code = object.code ?? 0;
        message.data = object.data ?? new Uint8Array();
        message.log = object.log ?? "";
        message.info = object.info ?? "";
        if (object.gasWanted !== undefined && object.gasWanted !== null) {
            message.gasWanted = BigInt(object.gasWanted.toString());
        }
        if (object.gasUsed !== undefined && object.gasUsed !== null) {
            message.gasUsed = BigInt(object.gasUsed.toString());
        }
        message.events = object.events?.map((e) => exports.Event.fromPartial(e)) || [];
        message.codespace = object.codespace ?? "";
        message.sender = object.sender ?? "";
        if (object.priority !== undefined && object.priority !== null) {
            message.priority = BigInt(object.priority.toString());
        }
        message.mempoolError = object.mempoolError ?? "";
        return message;
    },
};
function createBaseResponseDeliverTx() {
    return {
        code: 0,
        data: new Uint8Array(),
        log: "",
        info: "",
        gasWanted: BigInt(0),
        gasUsed: BigInt(0),
        events: [],
        codespace: "",
    };
}
exports.ResponseDeliverTx = {
    typeUrl: "/tendermint.abci.ResponseDeliverTx",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.code !== 0) {
            writer.uint32(8).uint32(message.code);
        }
        if (message.data.length !== 0) {
            writer.uint32(18).bytes(message.data);
        }
        if (message.log !== "") {
            writer.uint32(26).string(message.log);
        }
        if (message.info !== "") {
            writer.uint32(34).string(message.info);
        }
        if (message.gasWanted !== BigInt(0)) {
            writer.uint32(40).int64(message.gasWanted);
        }
        if (message.gasUsed !== BigInt(0)) {
            writer.uint32(48).int64(message.gasUsed);
        }
        for (const v of message.events) {
            exports.Event.encode(v, writer.uint32(58).fork()).ldelim();
        }
        if (message.codespace !== "") {
            writer.uint32(66).string(message.codespace);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseResponseDeliverTx();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.code = reader.uint32();
                    break;
                case 2:
                    message.data = reader.bytes();
                    break;
                case 3:
                    message.log = reader.string();
                    break;
                case 4:
                    message.info = reader.string();
                    break;
                case 5:
                    message.gasWanted = reader.int64();
                    break;
                case 6:
                    message.gasUsed = reader.int64();
                    break;
                case 7:
                    message.events.push(exports.Event.decode(reader, reader.uint32()));
                    break;
                case 8:
                    message.codespace = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseResponseDeliverTx();
        if ((0, helpers_1.isSet)(object.code))
            obj.code = Number(object.code);
        if ((0, helpers_1.isSet)(object.data))
            obj.data = (0, helpers_1.bytesFromBase64)(object.data);
        if ((0, helpers_1.isSet)(object.log))
            obj.log = String(object.log);
        if ((0, helpers_1.isSet)(object.info))
            obj.info = String(object.info);
        if ((0, helpers_1.isSet)(object.gas_wanted))
            obj.gasWanted = BigInt(object.gas_wanted.toString());
        if ((0, helpers_1.isSet)(object.gas_used))
            obj.gasUsed = BigInt(object.gas_used.toString());
        if (Array.isArray(object?.events))
            obj.events = object.events.map((e) => exports.Event.fromJSON(e));
        if ((0, helpers_1.isSet)(object.codespace))
            obj.codespace = String(object.codespace);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.code !== undefined && (obj.code = Math.round(message.code));
        message.data !== undefined &&
            (obj.data = (0, helpers_1.base64FromBytes)(message.data !== undefined ? message.data : new Uint8Array()));
        message.log !== undefined && (obj.log = message.log);
        message.info !== undefined && (obj.info = message.info);
        message.gasWanted !== undefined && (obj.gas_wanted = (message.gasWanted || BigInt(0)).toString());
        message.gasUsed !== undefined && (obj.gas_used = (message.gasUsed || BigInt(0)).toString());
        if (message.events) {
            obj.events = message.events.map((e) => (e ? exports.Event.toJSON(e) : undefined));
        }
        else {
            obj.events = [];
        }
        message.codespace !== undefined && (obj.codespace = message.codespace);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseResponseDeliverTx();
        message.code = object.code ?? 0;
        message.data = object.data ?? new Uint8Array();
        message.log = object.log ?? "";
        message.info = object.info ?? "";
        if (object.gasWanted !== undefined && object.gasWanted !== null) {
            message.gasWanted = BigInt(object.gasWanted.toString());
        }
        if (object.gasUsed !== undefined && object.gasUsed !== null) {
            message.gasUsed = BigInt(object.gasUsed.toString());
        }
        message.events = object.events?.map((e) => exports.Event.fromPartial(e)) || [];
        message.codespace = object.codespace ?? "";
        return message;
    },
};
function createBaseResponseEndBlock() {
    return {
        validatorUpdates: [],
        consensusParamUpdates: undefined,
        events: [],
    };
}
exports.ResponseEndBlock = {
    typeUrl: "/tendermint.abci.ResponseEndBlock",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.validatorUpdates) {
            exports.ValidatorUpdate.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.consensusParamUpdates !== undefined) {
            params_1.ConsensusParams.encode(message.consensusParamUpdates, writer.uint32(18).fork()).ldelim();
        }
        for (const v of message.events) {
            exports.Event.encode(v, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseResponseEndBlock();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.validatorUpdates.push(exports.ValidatorUpdate.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.consensusParamUpdates = params_1.ConsensusParams.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.events.push(exports.Event.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseResponseEndBlock();
        if (Array.isArray(object?.validatorUpdates))
            obj.validatorUpdates = object.validatorUpdates.map((e) => exports.ValidatorUpdate.fromJSON(e));
        if ((0, helpers_1.isSet)(object.consensusParamUpdates))
            obj.consensusParamUpdates = params_1.ConsensusParams.fromJSON(object.consensusParamUpdates);
        if (Array.isArray(object?.events))
            obj.events = object.events.map((e) => exports.Event.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.validatorUpdates) {
            obj.validatorUpdates = message.validatorUpdates.map((e) => (e ? exports.ValidatorUpdate.toJSON(e) : undefined));
        }
        else {
            obj.validatorUpdates = [];
        }
        message.consensusParamUpdates !== undefined &&
            (obj.consensusParamUpdates = message.consensusParamUpdates
                ? params_1.ConsensusParams.toJSON(message.consensusParamUpdates)
                : undefined);
        if (message.events) {
            obj.events = message.events.map((e) => (e ? exports.Event.toJSON(e) : undefined));
        }
        else {
            obj.events = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseResponseEndBlock();
        message.validatorUpdates = object.validatorUpdates?.map((e) => exports.ValidatorUpdate.fromPartial(e)) || [];
        if (object.consensusParamUpdates !== undefined && object.consensusParamUpdates !== null) {
            message.consensusParamUpdates = params_1.ConsensusParams.fromPartial(object.consensusParamUpdates);
        }
        message.events = object.events?.map((e) => exports.Event.fromPartial(e)) || [];
        return message;
    },
};
function createBaseResponseCommit() {
    return {
        data: new Uint8Array(),
        retainHeight: BigInt(0),
    };
}
exports.ResponseCommit = {
    typeUrl: "/tendermint.abci.ResponseCommit",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.data.length !== 0) {
            writer.uint32(18).bytes(message.data);
        }
        if (message.retainHeight !== BigInt(0)) {
            writer.uint32(24).int64(message.retainHeight);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseResponseCommit();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 2:
                    message.data = reader.bytes();
                    break;
                case 3:
                    message.retainHeight = reader.int64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseResponseCommit();
        if ((0, helpers_1.isSet)(object.data))
            obj.data = (0, helpers_1.bytesFromBase64)(object.data);
        if ((0, helpers_1.isSet)(object.retainHeight))
            obj.retainHeight = BigInt(object.retainHeight.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.data !== undefined &&
            (obj.data = (0, helpers_1.base64FromBytes)(message.data !== undefined ? message.data : new Uint8Array()));
        message.retainHeight !== undefined && (obj.retainHeight = (message.retainHeight || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseResponseCommit();
        message.data = object.data ?? new Uint8Array();
        if (object.retainHeight !== undefined && object.retainHeight !== null) {
            message.retainHeight = BigInt(object.retainHeight.toString());
        }
        return message;
    },
};
function createBaseResponseListSnapshots() {
    return {
        snapshots: [],
    };
}
exports.ResponseListSnapshots = {
    typeUrl: "/tendermint.abci.ResponseListSnapshots",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.snapshots) {
            exports.Snapshot.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseResponseListSnapshots();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.snapshots.push(exports.Snapshot.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseResponseListSnapshots();
        if (Array.isArray(object?.snapshots))
            obj.snapshots = object.snapshots.map((e) => exports.Snapshot.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.snapshots) {
            obj.snapshots = message.snapshots.map((e) => (e ? exports.Snapshot.toJSON(e) : undefined));
        }
        else {
            obj.snapshots = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseResponseListSnapshots();
        message.snapshots = object.snapshots?.map((e) => exports.Snapshot.fromPartial(e)) || [];
        return message;
    },
};
function createBaseResponseOfferSnapshot() {
    return {
        result: 0,
    };
}
exports.ResponseOfferSnapshot = {
    typeUrl: "/tendermint.abci.ResponseOfferSnapshot",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.result !== 0) {
            writer.uint32(8).int32(message.result);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseResponseOfferSnapshot();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.result = reader.int32();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseResponseOfferSnapshot();
        if ((0, helpers_1.isSet)(object.result))
            obj.result = responseOfferSnapshot_ResultFromJSON(object.result);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.result !== undefined && (obj.result = responseOfferSnapshot_ResultToJSON(message.result));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseResponseOfferSnapshot();
        message.result = object.result ?? 0;
        return message;
    },
};
function createBaseResponseLoadSnapshotChunk() {
    return {
        chunk: new Uint8Array(),
    };
}
exports.ResponseLoadSnapshotChunk = {
    typeUrl: "/tendermint.abci.ResponseLoadSnapshotChunk",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.chunk.length !== 0) {
            writer.uint32(10).bytes(message.chunk);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseResponseLoadSnapshotChunk();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.chunk = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseResponseLoadSnapshotChunk();
        if ((0, helpers_1.isSet)(object.chunk))
            obj.chunk = (0, helpers_1.bytesFromBase64)(object.chunk);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.chunk !== undefined &&
            (obj.chunk = (0, helpers_1.base64FromBytes)(message.chunk !== undefined ? message.chunk : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseResponseLoadSnapshotChunk();
        message.chunk = object.chunk ?? new Uint8Array();
        return message;
    },
};
function createBaseResponseApplySnapshotChunk() {
    return {
        result: 0,
        refetchChunks: [],
        rejectSenders: [],
    };
}
exports.ResponseApplySnapshotChunk = {
    typeUrl: "/tendermint.abci.ResponseApplySnapshotChunk",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.result !== 0) {
            writer.uint32(8).int32(message.result);
        }
        writer.uint32(18).fork();
        for (const v of message.refetchChunks) {
            writer.uint32(v);
        }
        writer.ldelim();
        for (const v of message.rejectSenders) {
            writer.uint32(26).string(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseResponseApplySnapshotChunk();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.result = reader.int32();
                    break;
                case 2:
                    if ((tag & 7) === 2) {
                        const end2 = reader.uint32() + reader.pos;
                        while (reader.pos < end2) {
                            message.refetchChunks.push(reader.uint32());
                        }
                    }
                    else {
                        message.refetchChunks.push(reader.uint32());
                    }
                    break;
                case 3:
                    message.rejectSenders.push(reader.string());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseResponseApplySnapshotChunk();
        if ((0, helpers_1.isSet)(object.result))
            obj.result = responseApplySnapshotChunk_ResultFromJSON(object.result);
        if (Array.isArray(object?.refetchChunks))
            obj.refetchChunks = object.refetchChunks.map((e) => Number(e));
        if (Array.isArray(object?.rejectSenders))
            obj.rejectSenders = object.rejectSenders.map((e) => String(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.result !== undefined && (obj.result = responseApplySnapshotChunk_ResultToJSON(message.result));
        if (message.refetchChunks) {
            obj.refetchChunks = message.refetchChunks.map((e) => Math.round(e));
        }
        else {
            obj.refetchChunks = [];
        }
        if (message.rejectSenders) {
            obj.rejectSenders = message.rejectSenders.map((e) => e);
        }
        else {
            obj.rejectSenders = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseResponseApplySnapshotChunk();
        message.result = object.result ?? 0;
        message.refetchChunks = object.refetchChunks?.map((e) => e) || [];
        message.rejectSenders = object.rejectSenders?.map((e) => e) || [];
        return message;
    },
};
function createBaseResponsePrepareProposal() {
    return {
        txs: [],
    };
}
exports.ResponsePrepareProposal = {
    typeUrl: "/tendermint.abci.ResponsePrepareProposal",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.txs) {
            writer.uint32(10).bytes(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseResponsePrepareProposal();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.txs.push(reader.bytes());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseResponsePrepareProposal();
        if (Array.isArray(object?.txs))
            obj.txs = object.txs.map((e) => (0, helpers_1.bytesFromBase64)(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.txs) {
            obj.txs = message.txs.map((e) => (0, helpers_1.base64FromBytes)(e !== undefined ? e : new Uint8Array()));
        }
        else {
            obj.txs = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseResponsePrepareProposal();
        message.txs = object.txs?.map((e) => e) || [];
        return message;
    },
};
function createBaseResponseProcessProposal() {
    return {
        status: 0,
    };
}
exports.ResponseProcessProposal = {
    typeUrl: "/tendermint.abci.ResponseProcessProposal",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.status !== 0) {
            writer.uint32(8).int32(message.status);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseResponseProcessProposal();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.status = reader.int32();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseResponseProcessProposal();
        if ((0, helpers_1.isSet)(object.status))
            obj.status = responseProcessProposal_ProposalStatusFromJSON(object.status);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.status !== undefined &&
            (obj.status = responseProcessProposal_ProposalStatusToJSON(message.status));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseResponseProcessProposal();
        message.status = object.status ?? 0;
        return message;
    },
};
function createBaseCommitInfo() {
    return {
        round: 0,
        votes: [],
    };
}
exports.CommitInfo = {
    typeUrl: "/tendermint.abci.CommitInfo",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.round !== 0) {
            writer.uint32(8).int32(message.round);
        }
        for (const v of message.votes) {
            exports.VoteInfo.encode(v, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseCommitInfo();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.round = reader.int32();
                    break;
                case 2:
                    message.votes.push(exports.VoteInfo.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseCommitInfo();
        if ((0, helpers_1.isSet)(object.round))
            obj.round = Number(object.round);
        if (Array.isArray(object?.votes))
            obj.votes = object.votes.map((e) => exports.VoteInfo.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.round !== undefined && (obj.round = Math.round(message.round));
        if (message.votes) {
            obj.votes = message.votes.map((e) => (e ? exports.VoteInfo.toJSON(e) : undefined));
        }
        else {
            obj.votes = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseCommitInfo();
        message.round = object.round ?? 0;
        message.votes = object.votes?.map((e) => exports.VoteInfo.fromPartial(e)) || [];
        return message;
    },
};
function createBaseExtendedCommitInfo() {
    return {
        round: 0,
        votes: [],
    };
}
exports.ExtendedCommitInfo = {
    typeUrl: "/tendermint.abci.ExtendedCommitInfo",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.round !== 0) {
            writer.uint32(8).int32(message.round);
        }
        for (const v of message.votes) {
            exports.ExtendedVoteInfo.encode(v, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseExtendedCommitInfo();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.round = reader.int32();
                    break;
                case 2:
                    message.votes.push(exports.ExtendedVoteInfo.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseExtendedCommitInfo();
        if ((0, helpers_1.isSet)(object.round))
            obj.round = Number(object.round);
        if (Array.isArray(object?.votes))
            obj.votes = object.votes.map((e) => exports.ExtendedVoteInfo.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.round !== undefined && (obj.round = Math.round(message.round));
        if (message.votes) {
            obj.votes = message.votes.map((e) => (e ? exports.ExtendedVoteInfo.toJSON(e) : undefined));
        }
        else {
            obj.votes = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseExtendedCommitInfo();
        message.round = object.round ?? 0;
        message.votes = object.votes?.map((e) => exports.ExtendedVoteInfo.fromPartial(e)) || [];
        return message;
    },
};
function createBaseEvent() {
    return {
        type: "",
        attributes: [],
    };
}
exports.Event = {
    typeUrl: "/tendermint.abci.Event",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.type !== "") {
            writer.uint32(10).string(message.type);
        }
        for (const v of message.attributes) {
            exports.EventAttribute.encode(v, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseEvent();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.type = reader.string();
                    break;
                case 2:
                    message.attributes.push(exports.EventAttribute.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseEvent();
        if ((0, helpers_1.isSet)(object.type))
            obj.type = String(object.type);
        if (Array.isArray(object?.attributes))
            obj.attributes = object.attributes.map((e) => exports.EventAttribute.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.type !== undefined && (obj.type = message.type);
        if (message.attributes) {
            obj.attributes = message.attributes.map((e) => (e ? exports.EventAttribute.toJSON(e) : undefined));
        }
        else {
            obj.attributes = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseEvent();
        message.type = object.type ?? "";
        message.attributes = object.attributes?.map((e) => exports.EventAttribute.fromPartial(e)) || [];
        return message;
    },
};
function createBaseEventAttribute() {
    return {
        key: "",
        value: "",
        index: false,
    };
}
exports.EventAttribute = {
    typeUrl: "/tendermint.abci.EventAttribute",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.key !== "") {
            writer.uint32(10).string(message.key);
        }
        if (message.value !== "") {
            writer.uint32(18).string(message.value);
        }
        if (message.index === true) {
            writer.uint32(24).bool(message.index);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseEventAttribute();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.key = reader.string();
                    break;
                case 2:
                    message.value = reader.string();
                    break;
                case 3:
                    message.index = reader.bool();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseEventAttribute();
        if ((0, helpers_1.isSet)(object.key))
            obj.key = String(object.key);
        if ((0, helpers_1.isSet)(object.value))
            obj.value = String(object.value);
        if ((0, helpers_1.isSet)(object.index))
            obj.index = Boolean(object.index);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.key !== undefined && (obj.key = message.key);
        message.value !== undefined && (obj.value = message.value);
        message.index !== undefined && (obj.index = message.index);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseEventAttribute();
        message.key = object.key ?? "";
        message.value = object.value ?? "";
        message.index = object.index ?? false;
        return message;
    },
};
function createBaseTxResult() {
    return {
        height: BigInt(0),
        index: 0,
        tx: new Uint8Array(),
        result: exports.ResponseDeliverTx.fromPartial({}),
    };
}
exports.TxResult = {
    typeUrl: "/tendermint.abci.TxResult",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.height !== BigInt(0)) {
            writer.uint32(8).int64(message.height);
        }
        if (message.index !== 0) {
            writer.uint32(16).uint32(message.index);
        }
        if (message.tx.length !== 0) {
            writer.uint32(26).bytes(message.tx);
        }
        if (message.result !== undefined) {
            exports.ResponseDeliverTx.encode(message.result, writer.uint32(34).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseTxResult();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.height = reader.int64();
                    break;
                case 2:
                    message.index = reader.uint32();
                    break;
                case 3:
                    message.tx = reader.bytes();
                    break;
                case 4:
                    message.result = exports.ResponseDeliverTx.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseTxResult();
        if ((0, helpers_1.isSet)(object.height))
            obj.height = BigInt(object.height.toString());
        if ((0, helpers_1.isSet)(object.index))
            obj.index = Number(object.index);
        if ((0, helpers_1.isSet)(object.tx))
            obj.tx = (0, helpers_1.bytesFromBase64)(object.tx);
        if ((0, helpers_1.isSet)(object.result))
            obj.result = exports.ResponseDeliverTx.fromJSON(object.result);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.height !== undefined && (obj.height = (message.height || BigInt(0)).toString());
        message.index !== undefined && (obj.index = Math.round(message.index));
        message.tx !== undefined &&
            (obj.tx = (0, helpers_1.base64FromBytes)(message.tx !== undefined ? message.tx : new Uint8Array()));
        message.result !== undefined &&
            (obj.result = message.result ? exports.ResponseDeliverTx.toJSON(message.result) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseTxResult();
        if (object.height !== undefined && object.height !== null) {
            message.height = BigInt(object.height.toString());
        }
        message.index = object.index ?? 0;
        message.tx = object.tx ?? new Uint8Array();
        if (object.result !== undefined && object.result !== null) {
            message.result = exports.ResponseDeliverTx.fromPartial(object.result);
        }
        return message;
    },
};
function createBaseValidator() {
    return {
        address: new Uint8Array(),
        power: BigInt(0),
    };
}
exports.Validator = {
    typeUrl: "/tendermint.abci.Validator",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.address.length !== 0) {
            writer.uint32(10).bytes(message.address);
        }
        if (message.power !== BigInt(0)) {
            writer.uint32(24).int64(message.power);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseValidator();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.address = reader.bytes();
                    break;
                case 3:
                    message.power = reader.int64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseValidator();
        if ((0, helpers_1.isSet)(object.address))
            obj.address = (0, helpers_1.bytesFromBase64)(object.address);
        if ((0, helpers_1.isSet)(object.power))
            obj.power = BigInt(object.power.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.address !== undefined &&
            (obj.address = (0, helpers_1.base64FromBytes)(message.address !== undefined ? message.address : new Uint8Array()));
        message.power !== undefined && (obj.power = (message.power || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseValidator();
        message.address = object.address ?? new Uint8Array();
        if (object.power !== undefined && object.power !== null) {
            message.power = BigInt(object.power.toString());
        }
        return message;
    },
};
function createBaseValidatorUpdate() {
    return {
        pubKey: keys_1.PublicKey.fromPartial({}),
        power: BigInt(0),
    };
}
exports.ValidatorUpdate = {
    typeUrl: "/tendermint.abci.ValidatorUpdate",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.pubKey !== undefined) {
            keys_1.PublicKey.encode(message.pubKey, writer.uint32(10).fork()).ldelim();
        }
        if (message.power !== BigInt(0)) {
            writer.uint32(16).int64(message.power);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseValidatorUpdate();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.pubKey = keys_1.PublicKey.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.power = reader.int64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseValidatorUpdate();
        if ((0, helpers_1.isSet)(object.pubKey))
            obj.pubKey = keys_1.PublicKey.fromJSON(object.pubKey);
        if ((0, helpers_1.isSet)(object.power))
            obj.power = BigInt(object.power.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.pubKey !== undefined &&
            (obj.pubKey = message.pubKey ? keys_1.PublicKey.toJSON(message.pubKey) : undefined);
        message.power !== undefined && (obj.power = (message.power || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseValidatorUpdate();
        if (object.pubKey !== undefined && object.pubKey !== null) {
            message.pubKey = keys_1.PublicKey.fromPartial(object.pubKey);
        }
        if (object.power !== undefined && object.power !== null) {
            message.power = BigInt(object.power.toString());
        }
        return message;
    },
};
function createBaseVoteInfo() {
    return {
        validator: exports.Validator.fromPartial({}),
        signedLastBlock: false,
    };
}
exports.VoteInfo = {
    typeUrl: "/tendermint.abci.VoteInfo",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.validator !== undefined) {
            exports.Validator.encode(message.validator, writer.uint32(10).fork()).ldelim();
        }
        if (message.signedLastBlock === true) {
            writer.uint32(16).bool(message.signedLastBlock);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseVoteInfo();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.validator = exports.Validator.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.signedLastBlock = reader.bool();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseVoteInfo();
        if ((0, helpers_1.isSet)(object.validator))
            obj.validator = exports.Validator.fromJSON(object.validator);
        if ((0, helpers_1.isSet)(object.signedLastBlock))
            obj.signedLastBlock = Boolean(object.signedLastBlock);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.validator !== undefined &&
            (obj.validator = message.validator ? exports.Validator.toJSON(message.validator) : undefined);
        message.signedLastBlock !== undefined && (obj.signedLastBlock = message.signedLastBlock);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseVoteInfo();
        if (object.validator !== undefined && object.validator !== null) {
            message.validator = exports.Validator.fromPartial(object.validator);
        }
        message.signedLastBlock = object.signedLastBlock ?? false;
        return message;
    },
};
function createBaseExtendedVoteInfo() {
    return {
        validator: exports.Validator.fromPartial({}),
        signedLastBlock: false,
        voteExtension: new Uint8Array(),
    };
}
exports.ExtendedVoteInfo = {
    typeUrl: "/tendermint.abci.ExtendedVoteInfo",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.validator !== undefined) {
            exports.Validator.encode(message.validator, writer.uint32(10).fork()).ldelim();
        }
        if (message.signedLastBlock === true) {
            writer.uint32(16).bool(message.signedLastBlock);
        }
        if (message.voteExtension.length !== 0) {
            writer.uint32(26).bytes(message.voteExtension);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseExtendedVoteInfo();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.validator = exports.Validator.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.signedLastBlock = reader.bool();
                    break;
                case 3:
                    message.voteExtension = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseExtendedVoteInfo();
        if ((0, helpers_1.isSet)(object.validator))
            obj.validator = exports.Validator.fromJSON(object.validator);
        if ((0, helpers_1.isSet)(object.signedLastBlock))
            obj.signedLastBlock = Boolean(object.signedLastBlock);
        if ((0, helpers_1.isSet)(object.voteExtension))
            obj.voteExtension = (0, helpers_1.bytesFromBase64)(object.voteExtension);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.validator !== undefined &&
            (obj.validator = message.validator ? exports.Validator.toJSON(message.validator) : undefined);
        message.signedLastBlock !== undefined && (obj.signedLastBlock = message.signedLastBlock);
        message.voteExtension !== undefined &&
            (obj.voteExtension = (0, helpers_1.base64FromBytes)(message.voteExtension !== undefined ? message.voteExtension : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseExtendedVoteInfo();
        if (object.validator !== undefined && object.validator !== null) {
            message.validator = exports.Validator.fromPartial(object.validator);
        }
        message.signedLastBlock = object.signedLastBlock ?? false;
        message.voteExtension = object.voteExtension ?? new Uint8Array();
        return message;
    },
};
function createBaseMisbehavior() {
    return {
        type: 0,
        validator: exports.Validator.fromPartial({}),
        height: BigInt(0),
        time: timestamp_1.Timestamp.fromPartial({}),
        totalVotingPower: BigInt(0),
    };
}
exports.Misbehavior = {
    typeUrl: "/tendermint.abci.Misbehavior",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.type !== 0) {
            writer.uint32(8).int32(message.type);
        }
        if (message.validator !== undefined) {
            exports.Validator.encode(message.validator, writer.uint32(18).fork()).ldelim();
        }
        if (message.height !== BigInt(0)) {
            writer.uint32(24).int64(message.height);
        }
        if (message.time !== undefined) {
            timestamp_1.Timestamp.encode(message.time, writer.uint32(34).fork()).ldelim();
        }
        if (message.totalVotingPower !== BigInt(0)) {
            writer.uint32(40).int64(message.totalVotingPower);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMisbehavior();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.type = reader.int32();
                    break;
                case 2:
                    message.validator = exports.Validator.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.height = reader.int64();
                    break;
                case 4:
                    message.time = timestamp_1.Timestamp.decode(reader, reader.uint32());
                    break;
                case 5:
                    message.totalVotingPower = reader.int64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMisbehavior();
        if ((0, helpers_1.isSet)(object.type))
            obj.type = misbehaviorTypeFromJSON(object.type);
        if ((0, helpers_1.isSet)(object.validator))
            obj.validator = exports.Validator.fromJSON(object.validator);
        if ((0, helpers_1.isSet)(object.height))
            obj.height = BigInt(object.height.toString());
        if ((0, helpers_1.isSet)(object.time))
            obj.time = (0, helpers_1.fromJsonTimestamp)(object.time);
        if ((0, helpers_1.isSet)(object.totalVotingPower))
            obj.totalVotingPower = BigInt(object.totalVotingPower.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.type !== undefined && (obj.type = misbehaviorTypeToJSON(message.type));
        message.validator !== undefined &&
            (obj.validator = message.validator ? exports.Validator.toJSON(message.validator) : undefined);
        message.height !== undefined && (obj.height = (message.height || BigInt(0)).toString());
        message.time !== undefined && (obj.time = (0, helpers_1.fromTimestamp)(message.time).toISOString());
        message.totalVotingPower !== undefined &&
            (obj.totalVotingPower = (message.totalVotingPower || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMisbehavior();
        message.type = object.type ?? 0;
        if (object.validator !== undefined && object.validator !== null) {
            message.validator = exports.Validator.fromPartial(object.validator);
        }
        if (object.height !== undefined && object.height !== null) {
            message.height = BigInt(object.height.toString());
        }
        if (object.time !== undefined && object.time !== null) {
            message.time = timestamp_1.Timestamp.fromPartial(object.time);
        }
        if (object.totalVotingPower !== undefined && object.totalVotingPower !== null) {
            message.totalVotingPower = BigInt(object.totalVotingPower.toString());
        }
        return message;
    },
};
function createBaseSnapshot() {
    return {
        height: BigInt(0),
        format: 0,
        chunks: 0,
        hash: new Uint8Array(),
        metadata: new Uint8Array(),
    };
}
exports.Snapshot = {
    typeUrl: "/tendermint.abci.Snapshot",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.height !== BigInt(0)) {
            writer.uint32(8).uint64(message.height);
        }
        if (message.format !== 0) {
            writer.uint32(16).uint32(message.format);
        }
        if (message.chunks !== 0) {
            writer.uint32(24).uint32(message.chunks);
        }
        if (message.hash.length !== 0) {
            writer.uint32(34).bytes(message.hash);
        }
        if (message.metadata.length !== 0) {
            writer.uint32(42).bytes(message.metadata);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSnapshot();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.height = reader.uint64();
                    break;
                case 2:
                    message.format = reader.uint32();
                    break;
                case 3:
                    message.chunks = reader.uint32();
                    break;
                case 4:
                    message.hash = reader.bytes();
                    break;
                case 5:
                    message.metadata = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseSnapshot();
        if ((0, helpers_1.isSet)(object.height))
            obj.height = BigInt(object.height.toString());
        if ((0, helpers_1.isSet)(object.format))
            obj.format = Number(object.format);
        if ((0, helpers_1.isSet)(object.chunks))
            obj.chunks = Number(object.chunks);
        if ((0, helpers_1.isSet)(object.hash))
            obj.hash = (0, helpers_1.bytesFromBase64)(object.hash);
        if ((0, helpers_1.isSet)(object.metadata))
            obj.metadata = (0, helpers_1.bytesFromBase64)(object.metadata);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.height !== undefined && (obj.height = (message.height || BigInt(0)).toString());
        message.format !== undefined && (obj.format = Math.round(message.format));
        message.chunks !== undefined && (obj.chunks = Math.round(message.chunks));
        message.hash !== undefined &&
            (obj.hash = (0, helpers_1.base64FromBytes)(message.hash !== undefined ? message.hash : new Uint8Array()));
        message.metadata !== undefined &&
            (obj.metadata = (0, helpers_1.base64FromBytes)(message.metadata !== undefined ? message.metadata : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseSnapshot();
        if (object.height !== undefined && object.height !== null) {
            message.height = BigInt(object.height.toString());
        }
        message.format = object.format ?? 0;
        message.chunks = object.chunks ?? 0;
        message.hash = object.hash ?? new Uint8Array();
        message.metadata = object.metadata ?? new Uint8Array();
        return message;
    },
};
class ABCIApplicationClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
        this.Echo = this.Echo.bind(this);
        this.Flush = this.Flush.bind(this);
        this.Info = this.Info.bind(this);
        this.DeliverTx = this.DeliverTx.bind(this);
        this.CheckTx = this.CheckTx.bind(this);
        this.Query = this.Query.bind(this);
        this.Commit = this.Commit.bind(this);
        this.InitChain = this.InitChain.bind(this);
        this.BeginBlock = this.BeginBlock.bind(this);
        this.EndBlock = this.EndBlock.bind(this);
        this.ListSnapshots = this.ListSnapshots.bind(this);
        this.OfferSnapshot = this.OfferSnapshot.bind(this);
        this.LoadSnapshotChunk = this.LoadSnapshotChunk.bind(this);
        this.ApplySnapshotChunk = this.ApplySnapshotChunk.bind(this);
        this.PrepareProposal = this.PrepareProposal.bind(this);
        this.ProcessProposal = this.ProcessProposal.bind(this);
    }
    Echo(request) {
        const data = exports.RequestEcho.encode(request).finish();
        const promise = this.rpc.request("tendermint.abci.ABCIApplication", "Echo", data);
        return promise.then((data) => exports.ResponseEcho.decode(new binary_1.BinaryReader(data)));
    }
    Flush(request = {}) {
        const data = exports.RequestFlush.encode(request).finish();
        const promise = this.rpc.request("tendermint.abci.ABCIApplication", "Flush", data);
        return promise.then((data) => exports.ResponseFlush.decode(new binary_1.BinaryReader(data)));
    }
    Info(request) {
        const data = exports.RequestInfo.encode(request).finish();
        const promise = this.rpc.request("tendermint.abci.ABCIApplication", "Info", data);
        return promise.then((data) => exports.ResponseInfo.decode(new binary_1.BinaryReader(data)));
    }
    DeliverTx(request) {
        const data = exports.RequestDeliverTx.encode(request).finish();
        const promise = this.rpc.request("tendermint.abci.ABCIApplication", "DeliverTx", data);
        return promise.then((data) => exports.ResponseDeliverTx.decode(new binary_1.BinaryReader(data)));
    }
    CheckTx(request) {
        const data = exports.RequestCheckTx.encode(request).finish();
        const promise = this.rpc.request("tendermint.abci.ABCIApplication", "CheckTx", data);
        return promise.then((data) => exports.ResponseCheckTx.decode(new binary_1.BinaryReader(data)));
    }
    Query(request) {
        const data = exports.RequestQuery.encode(request).finish();
        const promise = this.rpc.request("tendermint.abci.ABCIApplication", "Query", data);
        return promise.then((data) => exports.ResponseQuery.decode(new binary_1.BinaryReader(data)));
    }
    Commit(request = {}) {
        const data = exports.RequestCommit.encode(request).finish();
        const promise = this.rpc.request("tendermint.abci.ABCIApplication", "Commit", data);
        return promise.then((data) => exports.ResponseCommit.decode(new binary_1.BinaryReader(data)));
    }
    InitChain(request) {
        const data = exports.RequestInitChain.encode(request).finish();
        const promise = this.rpc.request("tendermint.abci.ABCIApplication", "InitChain", data);
        return promise.then((data) => exports.ResponseInitChain.decode(new binary_1.BinaryReader(data)));
    }
    BeginBlock(request) {
        const data = exports.RequestBeginBlock.encode(request).finish();
        const promise = this.rpc.request("tendermint.abci.ABCIApplication", "BeginBlock", data);
        return promise.then((data) => exports.ResponseBeginBlock.decode(new binary_1.BinaryReader(data)));
    }
    EndBlock(request) {
        const data = exports.RequestEndBlock.encode(request).finish();
        const promise = this.rpc.request("tendermint.abci.ABCIApplication", "EndBlock", data);
        return promise.then((data) => exports.ResponseEndBlock.decode(new binary_1.BinaryReader(data)));
    }
    ListSnapshots(request = {}) {
        const data = exports.RequestListSnapshots.encode(request).finish();
        const promise = this.rpc.request("tendermint.abci.ABCIApplication", "ListSnapshots", data);
        return promise.then((data) => exports.ResponseListSnapshots.decode(new binary_1.BinaryReader(data)));
    }
    OfferSnapshot(request) {
        const data = exports.RequestOfferSnapshot.encode(request).finish();
        const promise = this.rpc.request("tendermint.abci.ABCIApplication", "OfferSnapshot", data);
        return promise.then((data) => exports.ResponseOfferSnapshot.decode(new binary_1.BinaryReader(data)));
    }
    LoadSnapshotChunk(request) {
        const data = exports.RequestLoadSnapshotChunk.encode(request).finish();
        const promise = this.rpc.request("tendermint.abci.ABCIApplication", "LoadSnapshotChunk", data);
        return promise.then((data) => exports.ResponseLoadSnapshotChunk.decode(new binary_1.BinaryReader(data)));
    }
    ApplySnapshotChunk(request) {
        const data = exports.RequestApplySnapshotChunk.encode(request).finish();
        const promise = this.rpc.request("tendermint.abci.ABCIApplication", "ApplySnapshotChunk", data);
        return promise.then((data) => exports.ResponseApplySnapshotChunk.decode(new binary_1.BinaryReader(data)));
    }
    PrepareProposal(request) {
        const data = exports.RequestPrepareProposal.encode(request).finish();
        const promise = this.rpc.request("tendermint.abci.ABCIApplication", "PrepareProposal", data);
        return promise.then((data) => exports.ResponsePrepareProposal.decode(new binary_1.BinaryReader(data)));
    }
    ProcessProposal(request) {
        const data = exports.RequestProcessProposal.encode(request).finish();
        const promise = this.rpc.request("tendermint.abci.ABCIApplication", "ProcessProposal", data);
        return promise.then((data) => exports.ResponseProcessProposal.decode(new binary_1.BinaryReader(data)));
    }
}
exports.ABCIApplicationClientImpl = ABCIApplicationClientImpl;
//# sourceMappingURL=types.js.map