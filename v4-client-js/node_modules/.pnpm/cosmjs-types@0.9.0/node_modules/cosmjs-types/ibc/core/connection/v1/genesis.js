"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GenesisState = exports.protobufPackage = void 0;
/* eslint-disable */
const connection_1 = require("./connection");
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "ibc.core.connection.v1";
function createBaseGenesisState() {
    return {
        connections: [],
        clientConnectionPaths: [],
        nextConnectionSequence: BigInt(0),
        params: connection_1.Params.fromPartial({}),
    };
}
exports.GenesisState = {
    typeUrl: "/ibc.core.connection.v1.GenesisState",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.connections) {
            connection_1.IdentifiedConnection.encode(v, writer.uint32(10).fork()).ldelim();
        }
        for (const v of message.clientConnectionPaths) {
            connection_1.ConnectionPaths.encode(v, writer.uint32(18).fork()).ldelim();
        }
        if (message.nextConnectionSequence !== BigInt(0)) {
            writer.uint32(24).uint64(message.nextConnectionSequence);
        }
        if (message.params !== undefined) {
            connection_1.Params.encode(message.params, writer.uint32(34).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGenesisState();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.connections.push(connection_1.IdentifiedConnection.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.clientConnectionPaths.push(connection_1.ConnectionPaths.decode(reader, reader.uint32()));
                    break;
                case 3:
                    message.nextConnectionSequence = reader.uint64();
                    break;
                case 4:
                    message.params = connection_1.Params.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseGenesisState();
        if (Array.isArray(object?.connections))
            obj.connections = object.connections.map((e) => connection_1.IdentifiedConnection.fromJSON(e));
        if (Array.isArray(object?.clientConnectionPaths))
            obj.clientConnectionPaths = object.clientConnectionPaths.map((e) => connection_1.ConnectionPaths.fromJSON(e));
        if ((0, helpers_1.isSet)(object.nextConnectionSequence))
            obj.nextConnectionSequence = BigInt(object.nextConnectionSequence.toString());
        if ((0, helpers_1.isSet)(object.params))
            obj.params = connection_1.Params.fromJSON(object.params);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.connections) {
            obj.connections = message.connections.map((e) => (e ? connection_1.IdentifiedConnection.toJSON(e) : undefined));
        }
        else {
            obj.connections = [];
        }
        if (message.clientConnectionPaths) {
            obj.clientConnectionPaths = message.clientConnectionPaths.map((e) => e ? connection_1.ConnectionPaths.toJSON(e) : undefined);
        }
        else {
            obj.clientConnectionPaths = [];
        }
        message.nextConnectionSequence !== undefined &&
            (obj.nextConnectionSequence = (message.nextConnectionSequence || BigInt(0)).toString());
        message.params !== undefined && (obj.params = message.params ? connection_1.Params.toJSON(message.params) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGenesisState();
        message.connections = object.connections?.map((e) => connection_1.IdentifiedConnection.fromPartial(e)) || [];
        message.clientConnectionPaths =
            object.clientConnectionPaths?.map((e) => connection_1.ConnectionPaths.fromPartial(e)) || [];
        if (object.nextConnectionSequence !== undefined && object.nextConnectionSequence !== null) {
            message.nextConnectionSequence = BigInt(object.nextConnectionSequence.toString());
        }
        if (object.params !== undefined && object.params !== null) {
            message.params = connection_1.Params.fromPartial(object.params);
        }
        return message;
    },
};
//# sourceMappingURL=genesis.js.map