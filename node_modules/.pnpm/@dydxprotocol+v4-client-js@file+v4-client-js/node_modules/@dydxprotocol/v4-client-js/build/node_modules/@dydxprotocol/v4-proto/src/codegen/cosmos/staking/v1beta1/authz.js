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
exports.StakeAuthorization_Validators = exports.StakeAuthorization = exports.authorizationTypeToJSON = exports.authorizationTypeFromJSON = exports.AuthorizationTypeSDKType = exports.AuthorizationType = void 0;
const coin_1 = require("../../base/v1beta1/coin");
const _m0 = __importStar(require("protobufjs/minimal"));
/**
 * AuthorizationType defines the type of staking module authorization type
 *
 * Since: cosmos-sdk 0.43
 */
var AuthorizationType;
(function (AuthorizationType) {
    /** AUTHORIZATION_TYPE_UNSPECIFIED - AUTHORIZATION_TYPE_UNSPECIFIED specifies an unknown authorization type */
    AuthorizationType[AuthorizationType["AUTHORIZATION_TYPE_UNSPECIFIED"] = 0] = "AUTHORIZATION_TYPE_UNSPECIFIED";
    /** AUTHORIZATION_TYPE_DELEGATE - AUTHORIZATION_TYPE_DELEGATE defines an authorization type for Msg/Delegate */
    AuthorizationType[AuthorizationType["AUTHORIZATION_TYPE_DELEGATE"] = 1] = "AUTHORIZATION_TYPE_DELEGATE";
    /** AUTHORIZATION_TYPE_UNDELEGATE - AUTHORIZATION_TYPE_UNDELEGATE defines an authorization type for Msg/Undelegate */
    AuthorizationType[AuthorizationType["AUTHORIZATION_TYPE_UNDELEGATE"] = 2] = "AUTHORIZATION_TYPE_UNDELEGATE";
    /** AUTHORIZATION_TYPE_REDELEGATE - AUTHORIZATION_TYPE_REDELEGATE defines an authorization type for Msg/BeginRedelegate */
    AuthorizationType[AuthorizationType["AUTHORIZATION_TYPE_REDELEGATE"] = 3] = "AUTHORIZATION_TYPE_REDELEGATE";
    /** AUTHORIZATION_TYPE_CANCEL_UNBONDING_DELEGATION - AUTHORIZATION_TYPE_CANCEL_UNBONDING_DELEGATION defines an authorization type for Msg/MsgCancelUnbondingDelegation */
    AuthorizationType[AuthorizationType["AUTHORIZATION_TYPE_CANCEL_UNBONDING_DELEGATION"] = 4] = "AUTHORIZATION_TYPE_CANCEL_UNBONDING_DELEGATION";
    AuthorizationType[AuthorizationType["UNRECOGNIZED"] = -1] = "UNRECOGNIZED";
})(AuthorizationType = exports.AuthorizationType || (exports.AuthorizationType = {}));
exports.AuthorizationTypeSDKType = AuthorizationType;
function authorizationTypeFromJSON(object) {
    switch (object) {
        case 0:
        case "AUTHORIZATION_TYPE_UNSPECIFIED":
            return AuthorizationType.AUTHORIZATION_TYPE_UNSPECIFIED;
        case 1:
        case "AUTHORIZATION_TYPE_DELEGATE":
            return AuthorizationType.AUTHORIZATION_TYPE_DELEGATE;
        case 2:
        case "AUTHORIZATION_TYPE_UNDELEGATE":
            return AuthorizationType.AUTHORIZATION_TYPE_UNDELEGATE;
        case 3:
        case "AUTHORIZATION_TYPE_REDELEGATE":
            return AuthorizationType.AUTHORIZATION_TYPE_REDELEGATE;
        case 4:
        case "AUTHORIZATION_TYPE_CANCEL_UNBONDING_DELEGATION":
            return AuthorizationType.AUTHORIZATION_TYPE_CANCEL_UNBONDING_DELEGATION;
        case -1:
        case "UNRECOGNIZED":
        default:
            return AuthorizationType.UNRECOGNIZED;
    }
}
exports.authorizationTypeFromJSON = authorizationTypeFromJSON;
function authorizationTypeToJSON(object) {
    switch (object) {
        case AuthorizationType.AUTHORIZATION_TYPE_UNSPECIFIED:
            return "AUTHORIZATION_TYPE_UNSPECIFIED";
        case AuthorizationType.AUTHORIZATION_TYPE_DELEGATE:
            return "AUTHORIZATION_TYPE_DELEGATE";
        case AuthorizationType.AUTHORIZATION_TYPE_UNDELEGATE:
            return "AUTHORIZATION_TYPE_UNDELEGATE";
        case AuthorizationType.AUTHORIZATION_TYPE_REDELEGATE:
            return "AUTHORIZATION_TYPE_REDELEGATE";
        case AuthorizationType.AUTHORIZATION_TYPE_CANCEL_UNBONDING_DELEGATION:
            return "AUTHORIZATION_TYPE_CANCEL_UNBONDING_DELEGATION";
        case AuthorizationType.UNRECOGNIZED:
        default:
            return "UNRECOGNIZED";
    }
}
exports.authorizationTypeToJSON = authorizationTypeToJSON;
function createBaseStakeAuthorization() {
    return {
        maxTokens: undefined,
        allowList: undefined,
        denyList: undefined,
        authorizationType: 0
    };
}
exports.StakeAuthorization = {
    encode(message, writer = _m0.Writer.create()) {
        if (message.maxTokens !== undefined) {
            coin_1.Coin.encode(message.maxTokens, writer.uint32(10).fork()).ldelim();
        }
        if (message.allowList !== undefined) {
            exports.StakeAuthorization_Validators.encode(message.allowList, writer.uint32(18).fork()).ldelim();
        }
        if (message.denyList !== undefined) {
            exports.StakeAuthorization_Validators.encode(message.denyList, writer.uint32(26).fork()).ldelim();
        }
        if (message.authorizationType !== 0) {
            writer.uint32(32).int32(message.authorizationType);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseStakeAuthorization();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.maxTokens = coin_1.Coin.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.allowList = exports.StakeAuthorization_Validators.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.denyList = exports.StakeAuthorization_Validators.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.authorizationType = reader.int32();
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
        const message = createBaseStakeAuthorization();
        message.maxTokens = object.maxTokens !== undefined && object.maxTokens !== null ? coin_1.Coin.fromPartial(object.maxTokens) : undefined;
        message.allowList = object.allowList !== undefined && object.allowList !== null ? exports.StakeAuthorization_Validators.fromPartial(object.allowList) : undefined;
        message.denyList = object.denyList !== undefined && object.denyList !== null ? exports.StakeAuthorization_Validators.fromPartial(object.denyList) : undefined;
        message.authorizationType = (_a = object.authorizationType) !== null && _a !== void 0 ? _a : 0;
        return message;
    }
};
function createBaseStakeAuthorization_Validators() {
    return {
        address: []
    };
}
exports.StakeAuthorization_Validators = {
    encode(message, writer = _m0.Writer.create()) {
        for (const v of message.address) {
            writer.uint32(10).string(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseStakeAuthorization_Validators();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.address.push(reader.string());
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
        const message = createBaseStakeAuthorization_Validators();
        message.address = ((_a = object.address) === null || _a === void 0 ? void 0 : _a.map(e => e)) || [];
        return message;
    }
};
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiYXV0aHouanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi8uLi8uLi8uLi8uLi8uLi8uLi9ub2RlX21vZHVsZXMvQGR5ZHhwcm90b2NvbC92NC1wcm90by9zcmMvY29kZWdlbi9jb3Ntb3Mvc3Rha2luZy92MWJldGExL2F1dGh6LnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7O0FBQUEsa0RBQTREO0FBQzVELHdEQUEwQztBQUUxQzs7OztHQUlHO0FBRUgsSUFBWSxpQkFnQlg7QUFoQkQsV0FBWSxpQkFBaUI7SUFDM0IsOEdBQThHO0lBQzlHLDZHQUFrQyxDQUFBO0lBRWxDLCtHQUErRztJQUMvRyx1R0FBK0IsQ0FBQTtJQUUvQixxSEFBcUg7SUFDckgsMkdBQWlDLENBQUE7SUFFakMsMEhBQTBIO0lBQzFILDJHQUFpQyxDQUFBO0lBRWpDLHlLQUF5SztJQUN6Syw2SUFBa0QsQ0FBQTtJQUNsRCwwRUFBaUIsQ0FBQTtBQUNuQixDQUFDLEVBaEJXLGlCQUFpQixHQUFqQix5QkFBaUIsS0FBakIseUJBQWlCLFFBZ0I1QjtBQUNZLFFBQUEsd0JBQXdCLEdBQUcsaUJBQWlCLENBQUM7QUFDMUQsU0FBZ0IseUJBQXlCLENBQUMsTUFBVztJQUNuRCxRQUFRLE1BQU0sRUFBRTtRQUNkLEtBQUssQ0FBQyxDQUFDO1FBQ1AsS0FBSyxnQ0FBZ0M7WUFDbkMsT0FBTyxpQkFBaUIsQ0FBQyw4QkFBOEIsQ0FBQztRQUUxRCxLQUFLLENBQUMsQ0FBQztRQUNQLEtBQUssNkJBQTZCO1lBQ2hDLE9BQU8saUJBQWlCLENBQUMsMkJBQTJCLENBQUM7UUFFdkQsS0FBSyxDQUFDLENBQUM7UUFDUCxLQUFLLCtCQUErQjtZQUNsQyxPQUFPLGlCQUFpQixDQUFDLDZCQUE2QixDQUFDO1FBRXpELEtBQUssQ0FBQyxDQUFDO1FBQ1AsS0FBSywrQkFBK0I7WUFDbEMsT0FBTyxpQkFBaUIsQ0FBQyw2QkFBNkIsQ0FBQztRQUV6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