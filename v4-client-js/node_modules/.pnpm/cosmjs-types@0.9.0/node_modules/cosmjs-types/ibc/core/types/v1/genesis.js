"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GenesisState = exports.protobufPackage = void 0;
/* eslint-disable */
const genesis_1 = require("../../client/v1/genesis");
const genesis_2 = require("../../connection/v1/genesis");
const genesis_3 = require("../../channel/v1/genesis");
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "ibc.core.types.v1";
function createBaseGenesisState() {
    return {
        clientGenesis: genesis_1.GenesisState.fromPartial({}),
        connectionGenesis: genesis_2.GenesisState.fromPartial({}),
        channelGenesis: genesis_3.GenesisState.fromPartial({}),
    };
}
exports.GenesisState = {
    typeUrl: "/ibc.core.types.v1.GenesisState",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.clientGenesis !== undefined) {
            genesis_1.GenesisState.encode(message.clientGenesis, writer.uint32(10).fork()).ldelim();
        }
        if (message.connectionGenesis !== undefined) {
            genesis_2.GenesisState.encode(message.connectionGenesis, writer.uint32(18).fork()).ldelim();
        }
        if (message.channelGenesis !== undefined) {
            genesis_3.GenesisState.encode(message.channelGenesis, writer.uint32(26).fork()).ldelim();
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
                    message.clientGenesis = genesis_1.GenesisState.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.connectionGenesis = genesis_2.GenesisState.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.channelGenesis = genesis_3.GenesisState.decode(reader, reader.uint32());
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
        if ((0, helpers_1.isSet)(object.clientGenesis))
            obj.clientGenesis = genesis_1.GenesisState.fromJSON(object.clientGenesis);
        if ((0, helpers_1.isSet)(object.connectionGenesis))
            obj.connectionGenesis = genesis_2.GenesisState.fromJSON(object.connectionGenesis);
        if ((0, helpers_1.isSet)(object.channelGenesis))
            obj.channelGenesis = genesis_3.GenesisState.fromJSON(object.channelGenesis);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.clientGenesis !== undefined &&
            (obj.clientGenesis = message.clientGenesis ? genesis_1.GenesisState.toJSON(message.clientGenesis) : undefined);
        message.connectionGenesis !== undefined &&
            (obj.connectionGenesis = message.connectionGenesis
                ? genesis_2.GenesisState.toJSON(message.connectionGenesis)
                : undefined);
        message.channelGenesis !== undefined &&
            (obj.channelGenesis = message.channelGenesis
                ? genesis_3.GenesisState.toJSON(message.channelGenesis)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGenesisState();
        if (object.clientGenesis !== undefined && object.clientGenesis !== null) {
            message.clientGenesis = genesis_1.GenesisState.fromPartial(object.clientGenesis);
        }
        if (object.connectionGenesis !== undefined && object.connectionGenesis !== null) {
            message.connectionGenesis = genesis_2.GenesisState.fromPartial(object.connectionGenesis);
        }
        if (object.channelGenesis !== undefined && object.channelGenesis !== null) {
            message.channelGenesis = genesis_3.GenesisState.fromPartial(object.channelGenesis);
        }
        return message;
    },
};
//# sourceMappingURL=genesis.js.map