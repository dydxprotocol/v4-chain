"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.FillablePriceConfig = exports.SubaccountBlockLimits = exports.PositionBlockLimits = exports.LiquidationsConfig = void 0;
const _m0 = __importStar(require("protobufjs/minimal"));
const helpers_1 = require("../../helpers");
function createBaseLiquidationsConfig() {
    return {
        maxLiquidationFeePpm: 0,
        positionBlockLimits: undefined,
        subaccountBlockLimits: undefined,
        fillablePriceConfig: undefined
    };
}
exports.LiquidationsConfig = {
    encode(message, writer = _m0.Writer.create()) {
        if (message.maxLiquidationFeePpm !== 0) {
            writer.uint32(8).uint32(message.maxLiquidationFeePpm);
        }
        if (message.positionBlockLimits !== undefined) {
            exports.PositionBlockLimits.encode(message.positionBlockLimits, writer.uint32(18).fork()).ldelim();
        }
        if (message.subaccountBlockLimits !== undefined) {
            exports.SubaccountBlockLimits.encode(message.subaccountBlockLimits, writer.uint32(26).fork()).ldelim();
        }
        if (message.fillablePriceConfig !== undefined) {
            exports.FillablePriceConfig.encode(message.fillablePriceConfig, writer.uint32(34).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseLiquidationsConfig();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.maxLiquidationFeePpm = reader.uint32();
                    break;
                case 2:
                    message.positionBlockLimits = exports.PositionBlockLimits.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.subaccountBlockLimits = exports.SubaccountBlockLimits.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.fillablePriceConfig = exports.FillablePriceConfig.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromPartial(object) {
        var _a;
        const message = createBaseLiquidationsConfig();
        message.maxLiquidationFeePpm = (_a = object.maxLiquidationFeePpm) !== null && _a !== void 0 ? _a : 0;
        message.positionBlockLimits = object.positionBlockLimits !== undefined && object.positionBlockLimits !== null ? exports.PositionBlockLimits.fromPartial(object.positionBlockLimits) : undefined;
        message.subaccountBlockLimits = object.subaccountBlockLimits !== undefined && object.subaccountBlockLimits !== null ? exports.SubaccountBlockLimits.fromPartial(object.subaccountBlockLimits) : undefined;
        message.fillablePriceConfig = object.fillablePriceConfig !== undefined && object.fillablePriceConfig !== null ? exports.FillablePriceConfig.fromPartial(object.fillablePriceConfig) : undefined;
        return message;
    }
};
function createBasePositionBlockLimits() {
    return {
        minPositionNotionalLiquidated: helpers_1.Long.UZERO,
        maxPositionPortionLiquidatedPpm: 0
    };
}
exports.PositionBlockLimits = {
    encode(message, writer = _m0.Writer.create()) {
        if (!message.minPositionNotionalLiquidated.isZero()) {
            writer.uint32(8).uint64(message.minPositionNotionalLiquidated);
        }
        if (message.maxPositionPortionLiquidatedPpm !== 0) {
            writer.uint32(16).uint32(message.maxPositionPortionLiquidatedPpm);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBasePositionBlockLimits();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.minPositionNotionalLiquidated = reader.uint64();
                    break;
                case 2:
                    message.maxPositionPortionLiquidatedPpm = reader.uint32();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromPartial(object) {
        var _a;
        const message = createBasePositionBlockLimits();
        message.minPositionNotionalLiquidated = object.minPositionNotionalLiquidated !== undefined && object.minPositionNotionalLiquidated !== null ? helpers_1.Long.fromValue(object.minPositionNotionalLiquidated) : helpers_1.Long.UZERO;
        message.maxPositionPortionLiquidatedPpm = (_a = object.maxPositionPortionLiquidatedPpm) !== null && _a !== void 0 ? _a : 0;
        return message;
    }
};
function createBaseSubaccountBlockLimits() {
    return {
        maxNotionalLiquidated: helpers_1.Long.UZERO,
        maxQuantumsInsuranceLost: helpers_1.Long.UZERO
    };
}
exports.SubaccountBlockLimits = {
    encode(message, writer = _m0.Writer.create()) {
        if (!message.maxNotionalLiquidated.isZero()) {
            writer.uint32(8).uint64(message.maxNotionalLiquidated);
        }
        if (!message.maxQuantumsInsuranceLost.isZero()) {
            writer.uint32(16).uint64(message.maxQuantumsInsuranceLost);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSubaccountBlockLimits();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.maxNotionalLiquidated = reader.uint64();
                    break;
                case 2:
                    message.maxQuantumsInsuranceLost = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromPartial(object) {
        const message = createBaseSubaccountBlockLimits();
        message.maxNotionalLiquidated = object.maxNotionalLiquidated !== undefined && object.maxNotionalLiquidated !== null ? helpers_1.Long.fromValue(object.maxNotionalLiquidated) : helpers_1.Long.UZERO;
        message.maxQuantumsInsuranceLost = object.maxQuantumsInsuranceLost !== undefined && object.maxQuantumsInsuranceLost !== null ? helpers_1.Long.fromValue(object.maxQuantumsInsuranceLost) : helpers_1.Long.UZERO;
        return message;
    }
};
function createBaseFillablePriceConfig() {
    return {
        bankruptcyAdjustmentPpm: 0,
        spreadToMaintenanceMarginRatioPpm: 0
    };
}
exports.FillablePriceConfig = {
    encode(message, writer = _m0.Writer.create()) {
        if (message.bankruptcyAdjustmentPpm !== 0) {
            writer.uint32(8).uint32(message.bankruptcyAdjustmentPpm);
        }
        if (message.spreadToMaintenanceMarginRatioPpm !== 0) {
            writer.uint32(16).uint32(message.spreadToMaintenanceMarginRatioPpm);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseFillablePriceConfig();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.bankruptcyAdjustmentPpm = reader.uint32();
                    break;
                case 2:
                    message.spreadToMaintenanceMarginRatioPpm = reader.uint32();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromPartial(object) {
        var _a, _b;
        const message = createBaseFillablePriceConfig();
        message.bankruptcyAdjustmentPpm = (_a = object.bankruptcyAdjustmentPpm) !== null && _a !== void 0 ? _a : 0;
        message.spreadToMaintenanceMarginRatioPpm = (_b = object.spreadToMaintenanceMarginRatioPpm) !== null && _b !== void 0 ? _b : 0;
        return message;
    }
};
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoibGlxdWlkYXRpb25zX2NvbmZpZy5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uLy4uLy4uLy4uLy4uL25vZGVfbW9kdWxlcy9AZHlkeHByb3RvY29sL3Y0LXByb3RvL3NyYy9jb2RlZ2VuL2R5ZHhwcm90b2NvbC9jbG9iL2xpcXVpZGF0aW9uc19jb25maWcudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7QUFBQSx3REFBMEM7QUFDMUMsMkNBQWtEO0FBc0hsRCxTQUFTLDRCQUE0QjtJQUNuQyxPQUFPO1FBQ0wsb0JBQW9CLEVBQUUsQ0FBQztRQUN2QixtQkFBbUIsRUFBRSxTQUFTO1FBQzlCLHFCQUFxQixFQUFFLFNBQVM7UUFDaEMsbUJBQW1CLEVBQUUsU0FBUztLQUMvQixDQUFDO0FBQ0osQ0FBQztBQUVZLFFBQUEsa0JBQWtCLEdBQUc7SUFDaEMsTUFBTSxDQUFDLE9BQTJCLEVBQUUsU0FBcUIsR0FBRyxDQUFDLE1BQU0sQ0FBQyxNQUFNLEVBQUU7UUFDMUUsSUFBSSxPQUFPLENBQUMsb0JBQW9CLEtBQUssQ0FBQyxFQUFFO1lBQ3RDLE1BQU0sQ0FBQyxNQUFNLENBQUMsQ0FBQyxDQUFDLENBQUMsTUFBTSxDQUFDLE9BQU8sQ0FBQyxvQkFBb0IsQ0FBQyxDQUFDO1NBQ3ZEO1FBRUQsSUFBSSxPQUFPLENBQUMsbUJBQW1CLEtBQUssU0FBUyxFQUFFO1lBQzdDLDJCQUFtQixDQUFDLE1BQU0sQ0FBQyxPQUFPLENBQUMsbUJBQW1CLEVBQUUsTUFBTSxDQUFDLE1BQU0sQ0FBQyxFQUFFLENBQUMsQ0FBQyxJQUFJLEVBQUUsQ0FBQyxDQUFDLE1BQU0sRUFBRSxDQUFDO1NBQzVGO1FBRUQsSUFBSSxPQUFPLENBQUMscUJBQXFCLEtBQUssU0FBUyxFQUFFO1lBQy9DLDZCQUFxQixDQUFDLE1BQU0sQ0FBQyxPQUFPLENBQUMscUJBQXFCLEVBQUUsTUFBTSxDQUFDLE1BQU0sQ0FBQyxFQUFFLENBQUMsQ0FBQyxJQUFJLEVBQUUsQ0FBQyxDQUFDLE1BQU0sRUFBRSxDQUFDO1NBQ2hHO1FBRUQsSUFBSSxPQUFPLENBQUMsbUJBQW1CLEtBQUssU0FBUyxFQUFFO1lBQzdDLDJCQUFtQixDQUFDLE1BQU0sQ0FBQyxPQUFPLENBQUMsbUJBQW1CLEVBQUUsTUFBTSxDQUFDLE1BQU0sQ0FBQyxFQUFFLENBQUMsQ0FBQyxJQUFJLEVBQUUsQ0FBQyxDQUFDLE1BQU0sRUFBRSxDQUFDO1NBQzVGO1FBRUQsT0FBTyxNQUFNLENBQUM7SUFDaEIsQ0FBQztJQUVELE1BQU0sQ0FBQyxLQUE4QixFQUFFLE1BQWU7UUFDcEQsTUFBTSxNQUFNLEdBQUcsS0FBSyxZQUFZLEdBQUcsQ0FBQyxNQUFNLENBQUMsQ0FBQyxDQUFDLEtBQUssQ0FBQyxDQUFDLENBQUMsSUFBSSxHQUFHLENBQUMsTUFBTSxDQUFDLEtBQUssQ0FBQyxDQUFDO1FBQzNFLElBQUksR0FBRyxHQUFHLE1BQU0sS0FBSyxTQUFTLENBQUMsQ0FBQyxDQUFDLE1BQU0sQ0FBQyxHQUFHLENBQUMsQ0FBQyxDQUFDLE1BQU0sQ0FBQyxHQUFHLEdBQUcsTUFBTSxDQUFDO1FBQ2xFLE1BQU0sT0FBTyxHQUFHLDRCQUE0QixFQUFFLENBQUM7UUFFL0MsT0FBTyxNQUFNLENBQUMsR0FBRyxHQUFHLEdBQUcsRUFBRTtZQUN2QixNQUFNLEdBQUcsR0FBRyxNQUFNLENBQUMsTUFBTSxFQUFFLENBQUM7WUFFNUIsUUFBUSxHQUFHLEtBQUssQ0FBQyxFQUFFO2dCQUNqQixLQUFLLENBQUM7b0JBQ0osT0FBTyxDQUFDLG9CQUFvQixHQUFHLE1BQU0sQ0FBQyxNQUFNLEVBQUUsQ0FBQztvQkFDL0MsTUFBTTtnQkFFUixLQUFLLENBQUM7b0JBQ0osT0FBTyxDQUFDLG1CQUFtQixHQUFHLDJCQUFtQixDQUFDLE1BQU0sQ0FBQyxNQUFNLEVBQUUsTUFBTSxDQUFDLE1BQU0sRUFBRSxDQUFDLENBQUM7b0JBQ2xGLE1BQU07Z0JBRVIsS0FBSyxDQUFDO29CQUNKLE9BQU8sQ0FBQyxxQkFBcUIsR0FBRyw2QkFBcUIsQ0FBQyxNQUFNLENBQUMsTUFBTSxFQUFFLE1BQU0sQ0FBQyxNQUFNLEVBQUUsQ0FBQyxDQUFDO29CQUN0RixNQUFNO2dCQUVSLEtBQUssQ0FBQztvQkFDSixPQUFPLENBQUMsbUJBQW1CLEdBQUcsMkJBQW1CLENBQUMsTUFBTSxDQUFDLE1BQU0sRUFBRSxNQUFNLENBQUMsTUFBTSxFQUFFLENBQUMsQ0FBQztvQkFDbEYsTUFBTTtnQkFFUjtvQkFDRSxNQUFNLENBQUMsUUFBUSxDQUFDLEdBQUcsR0FBRyxDQUFDLENBQUMsQ0FBQztvQkFDekIsTUFBTTthQUNUO1NBQ0Y7UUFFRCxPQUFPLE9BQU8sQ0FBQztJQUNqQixDQUFDO0lBRUQsV0FBVyxDQUFDLE1BQXVDOztRQUNqRCxNQUFNLE9BQU8sR0FBRyw0QkFBNEIsRUFBRSxDQUFDO1FBQy9DLE9BQU8sQ0FBQyxvQkFBb0IsR0FBRyxNQUFBLE1BQU0sQ0FBQyxvQkFBb0IsbUNBQUksQ0FBQyxDQUFDO1FBQ2hFLE9BQU8sQ0FBQyxtQkFBbUIsR0FBRyxNQUFNLENBQUMsbUJBQW1CLEtBQUssU0FBUyxJQUFJLE1BQU0sQ0FBQyxtQkFBbUIsS0FBSyxJQUFJLENBQUMsQ0FBQyxDQUFDLDJCQUFtQixDQUFDLFdBQVcsQ0FBQyxNQUFNLENBQUMsbUJBQW1CLENBQUMsQ0FBQyxDQUFDLENBQUMsU0FBUyxDQUFDO1FBQ3hMLE9BQU8sQ0FBQyxxQkFBcUIsR0FBRyxNQUFNLENBQUMscUJBQXFCLEtBQUssU0FBUyxJQUFJLE1BQU0sQ0FBQyxxQkFBcUIsS0FBSyxJQUFJLENBQUMsQ0FBQyxDQUFDLDZCQUFxQixDQUFDLFdBQVcsQ0FBQyxNQUFNLENBQUMscUJBQXFCLENBQUMsQ0FBQyxDQUFDLENBQUMsU0FBUyxDQUFDO1FBQ2xNLE9BQU8sQ0FBQyxtQkFBbUIsR0FBRyxNQUFNLENBQUMsbUJBQW1CLEtBQUssU0FBUyxJQUFJLE1BQU0sQ0FBQyxtQkFBbUIsS0FBSyxJQUFJLENBQUMsQ0FBQyxDQUFDLDJCQUFtQixDQUFDLFdBQVcsQ0FBQyxNQUFNLENBQUMsbUJBQW1CLENBQUMsQ0FBQyxDQUFDLENBQUMsU0FBUyxDQUFDO1FBQ3hMLE9BQU8sT0FBTyxDQUFDO0lBQ2pCLENBQUM7Q0FFRixDQUFDO0FBRUYsU0FBUyw2QkFBNkI7SUFDcEMsT0FBTztRQUNMLDZCQUE2QixFQUFFLGNBQUksQ0FBQyxLQUFLO1FBQ3pDLCtCQUErQixFQUFFLENBQUM7S0FDbkMsQ0FBQztBQUNKLENBQUM7QUFFWSxRQUFBLG1CQUFtQixHQUFHO0lBQ2pDLE1BQU0sQ0FBQyxPQUE0QixFQUFFLFNBQXFCLEdBQUcsQ0FBQyxNQUFNLENBQUMsTUFBTSxFQUFFO1FBQzNFLElBQUksQ0FBQyxPQUFPLENBQUMsNkJBQTZCLENBQUMsTUFBTSxFQUFFLEVBQUU7WUFDbkQsTUFBTSxDQUFDLE1BQU0sQ0FBQyxDQUFDLENBQUMsQ0FBQyxNQUFNLENBQUMsT0FBTyxDQUFDLDZCQUE2QixDQUFDLENBQUM7U0FDaEU7UUFFRCxJQUFJLE9BQU8sQ0FBQywrQkFBK0IsS0FBSyxDQUFDLEVBQUU7WUFDakQsTUFBTSxDQUFDLE1BQU0sQ0FBQyxFQUFFLENBQUMsQ0FBQyxNQUFNLENBQUMsT0FBTyxDQUFDLCtCQUErQixDQUFDLENBQUM7U0FDbkU7UUFFRCxPQUFPLE1BQU0sQ0FBQztJQUNoQixDQUFDO0lBRUQsTUFBTSxDQUFDLEtBQThCLEVBQUUsTUFBZTtRQUNwRCxNQUFNLE1BQU0sR0FBRyxLQUFLLFlBQVksR0FBRyxDQUFDLE1BQU0sQ0FBQyxDQUFDLENBQUMsS0FBSyxDQUFDLENBQUMsQ0FBQyxJQUFJLEdBQUcsQ0FBQyxNQUFNLENBQUMsS0FBSyxDQUFDLENBQUM7UUFDM0UsSUFBSSxHQUFHLEdBQUcsTUFBTSxLQUFLLFNBQVMsQ0FBQyxDQUFDLENBQUMsTUFBTSxDQUFDLEdBQUcsQ0FBQyxDQUFDLENBQUMsTUFBTSxDQUFDLEdBQUcsR0FBRyxNQUFNLENBQUM7UUFDbEUsTUFBTSxPQUFPLEdBQUcsNkJBQTZCLEVBQUUsQ0FBQztRQUVoRCxPQUFPLE1BQU0sQ0FBQyxHQUFHLEdBQUcsR0FBRyxFQUFFO1lBQ3ZCLE1BQU0sR0FBRyxHQUFHLE1BQU0sQ0FBQyxNQUFNLEVBQUUsQ0FBQztZQUU1QixRQUFRLEdBQUcsS0FBSyxDQUFDLEVBQUU7Z0JBQ2pCLEtBQUssQ0FBQztvQkFDSixPQUFPLENBQUMsNkJBQTZCLEdBQUksTUFBTSxDQUFDLE1BQU0sRUFBVyxDQUFDO29CQUNsRSxNQUFNO2dCQUVSLEtBQUssQ0FBQztvQkFDSixPQUFPLENBQUMsK0JBQStCLEdBQUcsTUFBTSxDQUFDLE1BQU0sRUFBRSxDQUFDO29CQUMxRCxNQUFNO2dCQUVSO29CQUNFLE1BQU0sQ0FBQyxRQUFRLENBQUMsR0FBRyxHQUFHLENBQUMsQ0FBQyxDQUFDO29CQUN6QixNQUFNO2FBQ1Q7U0FDRjtRQUVELE9BQU8sT0FBTyxDQUFDO0lBQ2pCLENBQUM7SUFFRCxXQUFXLENBQUMsTUFBd0M7O1FBQ2xELE1BQU0sT0FBTyxHQUFHLDZCQUE2QixFQUFFLENBQUM7UUFDaEQsT0FBTyxDQUFDLDZCQUE2QixHQUFHLE1BQU0sQ0FBQyw2QkFBNkIsS0FBSyxTQUFTLElBQUksTUFBTSxDQUFDLDZCQUE2QixLQUFLLElBQUksQ0FBQyxDQUFDLENBQUMsY0FBSSxDQUFDLFNBQVMsQ0FBQyxNQUFNLENBQUMsNkJBQTZCLENBQUMsQ0FBQyxDQUFDLENBQUMsY0FBSSxDQUFDLEtBQUssQ0FBQztRQUNoTixPQUFPLENBQUMsK0JBQStCLEdBQUcsTUFBQSxNQUFNLENBQUMsK0JBQStCLG1DQUFJLENBQUMsQ0FBQztRQUN0RixPQUFPLE9BQU8sQ0FBQztJQUNqQixDQUFDO0NBRUYsQ0FBQztBQUVGLFNBQVMsK0JBQStCO0lBQ3RDLE9BQU87UUFDTCxxQkFBcUIsRUFBRSxjQUFJLENBQUMsS0FBSztRQUNqQyx3QkFBd0IsRUFBRSxjQUFJLENBQUMsS0FBSztLQUNyQyxDQUFDO0FBQ0osQ0FBQztBQUVZLFFBQUEscUJBQXFCLEdBQUc7SUFDbkMsTUFBTSxDQUFDLE9BQThCLEVBQUUsU0FBcUIsR0FBRyxDQUFDLE1BQU0sQ0FBQyxNQUFNLEVBQUU7UUFDN0UsSUFBSSxDQUFDLE9BQU8sQ0FBQyxxQkFBcUIsQ0FBQyxNQUFNLEVBQUUsRUFBRTtZQUMzQyxNQUFNLENBQUMsTUFBTSxDQUFDLENBQUMsQ0FBQyxDQUFDLE1BQU0sQ0FBQyxPQUFPLENBQUMscUJBQXFCLENBQUMsQ0FBQztTQUN4RDtRQUVELElBQUksQ0FBQyxPQUFPLENBQUMsd0JBQXdCLENBQUMsTUFBTSxFQUFFLEVBQUU7WUFDOUMsTUFBTSxDQUFDLE1BQU0sQ0FBQyxFQUFFLENBQUMsQ0FBQyxNQUFNLENBQUMsT0FBTyxDQUFDLHdCQUF3QixDQUFDLENBQUM7U0FDNUQ7UUFFRCxPQUFPLE1BQU0sQ0FBQztJQUNoQixDQUFDO0lBRUQsTUFBTSxDQUFDLEtBQThCLEVBQUUsTUFBZTtRQUNwRCxNQUFNLE1BQU0sR0FBRyxLQUFLLFlBQVksR0FBRyxDQUFDLE1BQU0sQ0FBQyxDQUFDLENBQUMsS0FBSyxDQUFDLENBQUMsQ0FBQyxJQUFJLEdBQUcsQ0FBQyxNQUFNLENBQUMsS0FBSyxDQUFDLENBQUM7UUFDM0UsSUFBSSxHQUFHLEdBQUcsTUFBTSxLQUFLLFNBQVMsQ0FBQyxDQUFDLENBQUMsTUFBTSxDQUFDLEdBQUcsQ0FBQyxDQUFDLENBQUMsTUFBTSxDQUFDLEdBQUcsR0FBRyxNQUFNLENBQUM7UUFDbEUsTUFBTSxPQUFPLEdBQUcsK0JBQStCLEVBQUUsQ0FBQztRQUVsRCxPQUFPLE1BQU0sQ0FBQyxHQUFHLEdBQUcsR0FBRyxFQUFFO1lBQ3ZCLE1BQU0sR0FBRyxHQUFHLE1BQU0sQ0FBQyxNQUFNLEVBQUUsQ0FBQztZQUU1QixRQUFRLEdBQUcsS0FBSyxDQUFDLEVBQUU7Z0JBQ2pCLEtBQUssQ0FBQztvQkFDSixPQUFPLENBQUMscUJBQXFCLEdBQUksTUFBTSxDQUFDLE1BQU0sRUFBVyxDQUFDO29CQUMxRCxNQUFNO2dCQUVSLEtBQUssQ0FBQztvQkFDSixPQUFPLENBQUMsd0JBQXdCLEdBQUksTUFBTSxDQUFDLE1BQU0sRUFBVyxDQUFDO29CQUM3RCxNQUFNO2dCQUVSO29CQUNFLE1BQU0sQ0FBQyxRQUFRLENBQUMsR0FBRyxHQUFHLENBQUMsQ0FBQyxDQUFDO29CQUN6QixNQUFNO2FBQ1Q7U0FDRjtRQUVELE9BQU8sT0FBTyxDQUFDO0lBQ2pCLENBQUM7SUFFRCxXQUFXLENBQUMsTUFBMEM7UUFDcEQsTUFBTSxPQUFPLEdBQUcsK0JBQStCLEVBQUUsQ0FBQztRQUNsRCxPQUFPLENBQUMscUJBQXFCLEdBQUcsTUFBTSxDQUFDLHFCQUFxQixLQUFLLFNBQVMsSUFBSSxNQUFNLENBQUMscUJBQXFCLEtBQUssSUFBSSxDQUFDLENBQUMsQ0FBQyxjQUFJLENBQUMsU0FBUyxDQUFDLE1BQU0sQ0FBQyxxQkFBcUIsQ0FBQyxDQUFDLENBQUMsQ0FBQyxjQUFJLENBQUMsS0FBSyxDQUFDO1FBQ2hMLE9BQU8sQ0FBQyx3QkFBd0IsR0FBRyxNQUFNLENBQUMsd0JBQXdCLEtBQUssU0FBUyxJQUFJLE1BQU0sQ0FBQyx3QkFBd0IsS0FBSyxJQUFJLENBQUMsQ0FBQyxDQUFDLGNBQUksQ0FBQyxTQUFTLENBQUMsTUFBTSxDQUFDLHdCQUF3QixDQUFDLENBQUMsQ0FBQyxDQUFDLGNBQUksQ0FBQyxLQUFLLENBQUM7UUFDNUwsT0FBTyxPQUFPLENBQUM7SUFDakIsQ0FBQztDQUVGLENBQUM7QUFFRixTQUFTLDZCQUE2QjtJQUNwQyxPQUFPO1FBQ0wsdUJBQXVCLEVBQUUsQ0FBQztRQUMxQixpQ0FBaUMsRUFBRSxDQUFDO0tBQ3JDLENBQUM7QUFDSixDQUFDO0FBRVksUUFBQSxtQkFBbUIsR0FBRztJQUNqQyxNQUFNLENBQUMsT0FBNEIsRUFBRSxTQUFxQixHQUFHLENBQUMsTUFBTSxDQUFDLE1BQU0sRUFBRTtRQUMzRSxJQUFJLE9BQU8sQ0FBQyx1QkFBdUIsS0FBSyxDQUFDLEVBQUU7WUFDekMsTUFBTSxDQUFDLE1BQU0sQ0FBQyxDQUFDLENBQUMsQ0FBQyxNQUFNLENBQUMsT0FBTyxDQUFDLHVCQUF1QixDQUFDLENBQUM7U0FDMUQ7UUFFRCxJQUFJLE9BQU8sQ0FBQyxpQ0FBaUMsS0FBSyxDQUFDLEVBQUU7WUFDbkQsTUFBTSxDQUFDLE1BQU0sQ0FBQyxFQUFFLENBQUMsQ0FBQyxNQUFNLENBQUMsT0FBTyxDQUFDLGlDQUFpQyxDQUFDLENBQUM7U0FDckU7UUFFRCxPQUFPLE1BQU0sQ0FBQztJQUNoQixDQUFDO0lBRUQsTUFBTSxDQUFDLEtBQThCLEVBQUUsTUFBZTtRQUNwRCxNQUFNLE1BQU0sR0FBRyxLQUFLLFlBQVksR0FBRyxDQUFDLE1BQU0sQ0FBQyxDQUFDLENBQUMsS0FBSyxDQUFDLENBQUMsQ0FBQyxJQUFJLEdBQUcsQ0FBQyxNQUFNLENBQUMsS0FBSyxDQUFDLENBQUM7UUFDM0UsSUFBSSxHQUFHLEdBQUcsTUFBTSxLQUFLLFNBQVMsQ0FBQyxDQUFDLENBQUMsTUFBTSxDQUFDLEdBQUcsQ0FBQyxDQUFDLENBQUMsTUFBTSxDQUFDLEdBQUcsR0FBRyxNQUFNLENBQUM7UUFDbEUsTUFBTSxPQUFPLEdBQUcsNkJBQTZCLEVBQUUsQ0FBQztRQUVoRCxPQUFPLE1BQU0sQ0FBQyxHQUFHLEdBQUcsR0FBRyxFQUFFO1lBQ3ZCLE1BQU0sR0FBRyxHQUFHLE1BQU0sQ0FBQyxNQUFNLEVBQUUsQ0FBQztZQUU1QixRQUFRLEdBQUcsS0FBSyxDQUFDLEVBQUU7Z0JBQ2pCLEtBQUssQ0FBQztvQkFDSixPQUFPLENBQUMsdUJBQXVCLEdBQUcsTUFBTSxDQUFDLE1BQU0sRUFBRSxDQUFDO29CQUNsRCxNQUFNO2dCQUVSLEtBQUssQ0FBQztvQkFDSixPQUFPLENBQUMsaUNBQWlDLEdBQUcsTUFBTSxDQUFDLE1BQU0sRUFBRSxDQUFDO29CQUM1RCxNQUFNO2dCQUVSO29CQUNFLE1BQU0sQ0FBQyxRQUFRLENBQUMsR0FBRyxHQUFHLENBQUMsQ0FBQyxDQUFDO29CQUN6QixNQUFNO2FBQ1Q7U0FDRjtRQUVELE9BQU8sT0FBTyxDQUFDO0lBQ2pCLENBQUM7SUFFRCxXQUFXLENBQUMsTUFBd0M7O1FBQ2xELE1BQU0sT0FBTyxHQUFHLDZCQUE2QixFQUFFLENBQUM7UUFDaEQsT0FBTyxDQUFDLHVCQUF1QixHQUFHLE1BQUEsTUFBTSxDQUFDLHVCQUF1QixtQ0FBSSxDQUFDLENBQUM7UUFDdEUsT0FBTyxDQUFDLGlDQUFpQyxHQUFHLE1BQUEsTUFBTSxDQUFDLGlDQUFpQyxtQ0FBSSxDQUFDLENBQUM7UUFDMUYsT0FBTyxPQUFPLENBQUM7SUFDakIsQ0FBQztDQUVGLENBQUMifQ==