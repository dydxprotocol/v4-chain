"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Entry = exports.GenesisState = exports.protobufPackage = void 0;
/* eslint-disable */
const nft_1 = require("./nft");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.nft.v1beta1";
function createBaseGenesisState() {
    return {
        classes: [],
        entries: [],
    };
}
exports.GenesisState = {
    typeUrl: "/cosmos.nft.v1beta1.GenesisState",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.classes) {
            nft_1.Class.encode(v, writer.uint32(10).fork()).ldelim();
        }
        for (const v of message.entries) {
            exports.Entry.encode(v, writer.uint32(18).fork()).ldelim();
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
                    message.classes.push(nft_1.Class.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.entries.push(exports.Entry.decode(reader, reader.uint32()));
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
        if (Array.isArray(object?.classes))
            obj.classes = object.classes.map((e) => nft_1.Class.fromJSON(e));
        if (Array.isArray(object?.entries))
            obj.entries = object.entries.map((e) => exports.Entry.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.classes) {
            obj.classes = message.classes.map((e) => (e ? nft_1.Class.toJSON(e) : undefined));
        }
        else {
            obj.classes = [];
        }
        if (message.entries) {
            obj.entries = message.entries.map((e) => (e ? exports.Entry.toJSON(e) : undefined));
        }
        else {
            obj.entries = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGenesisState();
        message.classes = object.classes?.map((e) => nft_1.Class.fromPartial(e)) || [];
        message.entries = object.entries?.map((e) => exports.Entry.fromPartial(e)) || [];
        return message;
    },
};
function createBaseEntry() {
    return {
        owner: "",
        nfts: [],
    };
}
exports.Entry = {
    typeUrl: "/cosmos.nft.v1beta1.Entry",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.owner !== "") {
            writer.uint32(10).string(message.owner);
        }
        for (const v of message.nfts) {
            nft_1.NFT.encode(v, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseEntry();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.owner = reader.string();
                    break;
                case 2:
                    message.nfts.push(nft_1.NFT.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseEntry();
        if ((0, helpers_1.isSet)(object.owner))
            obj.owner = String(object.owner);
        if (Array.isArray(object?.nfts))
            obj.nfts = object.nfts.map((e) => nft_1.NFT.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.owner !== undefined && (obj.owner = message.owner);
        if (message.nfts) {
            obj.nfts = message.nfts.map((e) => (e ? nft_1.NFT.toJSON(e) : undefined));
        }
        else {
            obj.nfts = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseEntry();
        message.owner = object.owner ?? "";
        message.nfts = object.nfts?.map((e) => nft_1.NFT.fromPartial(e)) || [];
        return message;
    },
};
//# sourceMappingURL=genesis.js.map