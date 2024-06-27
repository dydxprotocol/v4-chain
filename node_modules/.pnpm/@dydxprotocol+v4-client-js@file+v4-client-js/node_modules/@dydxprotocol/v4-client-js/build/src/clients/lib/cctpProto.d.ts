import Long from "long";
import _m0 from "protobufjs/minimal";
export declare const protobufPackage = "circle.cctp.v1";
/**
 * Copyright (c) 2023, © Circle Internet Financial, LTD.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
/** TODO add comments */
export interface MsgUpdateOwner {
    from: string;
    newOwner: string;
}
export interface MsgUpdateOwnerResponse {
}
export interface MsgUpdateAttesterManager {
    from: string;
    newAttesterManager: string;
}
export interface MsgUpdateAttesterManagerResponse {
}
export interface MsgUpdateTokenController {
    from: string;
    newTokenController: string;
}
export interface MsgUpdateTokenControllerResponse {
}
export interface MsgUpdatePauser {
    from: string;
    newPauser: string;
}
export interface MsgUpdatePauserResponse {
}
export interface MsgAcceptOwner {
    from: string;
}
export interface MsgAcceptOwnerResponse {
}
export interface MsgEnableAttester {
    from: string;
    attester: string;
}
export interface MsgEnableAttesterResponse {
}
export interface MsgDisableAttester {
    from: string;
    attester: string;
}
export interface MsgDisableAttesterResponse {
}
export interface MsgPauseBurningAndMinting {
    from: string;
}
export interface MsgPauseBurningAndMintingResponse {
}
export interface MsgUnpauseBurningAndMinting {
    from: string;
}
export interface MsgUnpauseBurningAndMintingResponse {
}
export interface MsgPauseSendingAndReceivingMessages {
    from: string;
}
export interface MsgPauseSendingAndReceivingMessagesResponse {
}
export interface MsgUnpauseSendingAndReceivingMessages {
    from: string;
}
export interface MsgUnpauseSendingAndReceivingMessagesResponse {
}
export interface MsgUpdateMaxMessageBodySize {
    from: string;
    messageSize: Long;
}
export interface MsgUpdateMaxMessageBodySizeResponse {
}
export interface MsgSetMaxBurnAmountPerMessage {
    from: string;
    localToken: string;
    amount: string;
}
export interface MsgSetMaxBurnAmountPerMessageResponse {
}
export interface MsgDepositForBurn {
    from: string;
    amount: string;
    destinationDomain: number;
    mintRecipient: Uint8Array;
    burnToken: string;
}
export interface MsgDepositForBurnResponse {
    nonce: Long;
}
export interface MsgDepositForBurnWithCaller {
    from: string;
    amount: string;
    destinationDomain: number;
    mintRecipient: Uint8Array;
    burnToken: string;
    destinationCaller: Uint8Array;
}
export interface MsgDepositForBurnWithCallerResponse {
    nonce: Long;
}
export interface MsgReplaceDepositForBurn {
    from: string;
    originalMessage: Uint8Array;
    originalAttestation: Uint8Array;
    newDestinationCaller: Uint8Array;
    newMintRecipient: Uint8Array;
}
export interface MsgReplaceDepositForBurnResponse {
}
export interface MsgReceiveMessage {
    from: string;
    message: Uint8Array;
    attestation: Uint8Array;
}
export interface MsgReceiveMessageResponse {
    success: boolean;
}
export interface MsgSendMessage {
    from: string;
    destinationDomain: number;
    recipient: Uint8Array;
    messageBody: Uint8Array;
}
export interface MsgSendMessageResponse {
    nonce: Long;
}
export interface MsgSendMessageWithCaller {
    from: string;
    destinationDomain: number;
    recipient: Uint8Array;
    messageBody: Uint8Array;
    destinationCaller: Uint8Array;
}
export interface MsgSendMessageWithCallerResponse {
    nonce: Long;
}
export interface MsgReplaceMessage {
    from: string;
    originalMessage: Uint8Array;
    originalAttestation: Uint8Array;
    newMessageBody: Uint8Array;
    newDestinationCaller: Uint8Array;
}
export interface MsgReplaceMessageResponse {
}
export interface MsgUpdateSignatureThreshold {
    from: string;
    amount: number;
}
export interface MsgUpdateSignatureThresholdResponse {
}
export interface MsgLinkTokenPair {
    from: string;
    remoteDomain: number;
    remoteToken: Uint8Array;
    localToken: string;
}
export interface MsgLinkTokenPairResponse {
}
export interface MsgUnlinkTokenPair {
    from: string;
    remoteDomain: number;
    remoteToken: Uint8Array;
    localToken: string;
}
export interface MsgUnlinkTokenPairResponse {
}
export interface MsgAddRemoteTokenMessenger {
    from: string;
    domainId: number;
    address: Uint8Array;
}
export interface MsgAddRemoteTokenMessengerResponse {
}
export interface MsgRemoveRemoteTokenMessenger {
    from: string;
    domainId: number;
}
export interface MsgRemoveRemoteTokenMessengerResponse {
}
export declare const MsgUpdateOwner: {
    encode(message: MsgUpdateOwner, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateOwner;
    fromJSON(object: any): MsgUpdateOwner;
    toJSON(message: MsgUpdateOwner): unknown;
    create<I extends {
        from?: string | undefined;
        newOwner?: string | undefined;
    } & {
        from?: string | undefined;
        newOwner?: string | undefined;
    } & { [K in Exclude<keyof I, keyof MsgUpdateOwner>]: never; }>(base?: I | undefined): MsgUpdateOwner;
    fromPartial<I_1 extends {
        from?: string | undefined;
        newOwner?: string | undefined;
    } & {
        from?: string | undefined;
        newOwner?: string | undefined;
    } & { [K_1 in Exclude<keyof I_1, keyof MsgUpdateOwner>]: never; }>(object: I_1): MsgUpdateOwner;
};
export declare const MsgUpdateOwnerResponse: {
    encode(_: MsgUpdateOwnerResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateOwnerResponse;
    fromJSON(_: any): MsgUpdateOwnerResponse;
    toJSON(_: MsgUpdateOwnerResponse): unknown;
    create<I extends {} & {} & { [K in Exclude<keyof I, never>]: never; }>(base?: I | undefined): MsgUpdateOwnerResponse;
    fromPartial<I_1 extends {} & {} & { [K_1 in Exclude<keyof I_1, never>]: never; }>(_: I_1): MsgUpdateOwnerResponse;
};
export declare const MsgUpdateAttesterManager: {
    encode(message: MsgUpdateAttesterManager, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateAttesterManager;
    fromJSON(object: any): MsgUpdateAttesterManager;
    toJSON(message: MsgUpdateAttesterManager): unknown;
    create<I extends {
        from?: string | undefined;
        newAttesterManager?: string | undefined;
    } & {
        from?: string | undefined;
        newAttesterManager?: string | undefined;
    } & { [K in Exclude<keyof I, keyof MsgUpdateAttesterManager>]: never; }>(base?: I | undefined): MsgUpdateAttesterManager;
    fromPartial<I_1 extends {
        from?: string | undefined;
        newAttesterManager?: string | undefined;
    } & {
        from?: string | undefined;
        newAttesterManager?: string | undefined;
    } & { [K_1 in Exclude<keyof I_1, keyof MsgUpdateAttesterManager>]: never; }>(object: I_1): MsgUpdateAttesterManager;
};
export declare const MsgUpdateAttesterManagerResponse: {
    encode(_: MsgUpdateAttesterManagerResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateAttesterManagerResponse;
    fromJSON(_: any): MsgUpdateAttesterManagerResponse;
    toJSON(_: MsgUpdateAttesterManagerResponse): unknown;
    create<I extends {} & {} & { [K in Exclude<keyof I, never>]: never; }>(base?: I | undefined): MsgUpdateAttesterManagerResponse;
    fromPartial<I_1 extends {} & {} & { [K_1 in Exclude<keyof I_1, never>]: never; }>(_: I_1): MsgUpdateAttesterManagerResponse;
};
export declare const MsgUpdateTokenController: {
    encode(message: MsgUpdateTokenController, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateTokenController;
    fromJSON(object: any): MsgUpdateTokenController;
    toJSON(message: MsgUpdateTokenController): unknown;
    create<I extends {
        from?: string | undefined;
        newTokenController?: string | undefined;
    } & {
        from?: string | undefined;
        newTokenController?: string | undefined;
    } & { [K in Exclude<keyof I, keyof MsgUpdateTokenController>]: never; }>(base?: I | undefined): MsgUpdateTokenController;
    fromPartial<I_1 extends {
        from?: string | undefined;
        newTokenController?: string | undefined;
    } & {
        from?: string | undefined;
        newTokenController?: string | undefined;
    } & { [K_1 in Exclude<keyof I_1, keyof MsgUpdateTokenController>]: never; }>(object: I_1): MsgUpdateTokenController;
};
export declare const MsgUpdateTokenControllerResponse: {
    encode(_: MsgUpdateTokenControllerResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateTokenControllerResponse;
    fromJSON(_: any): MsgUpdateTokenControllerResponse;
    toJSON(_: MsgUpdateTokenControllerResponse): unknown;
    create<I extends {} & {} & { [K in Exclude<keyof I, never>]: never; }>(base?: I | undefined): MsgUpdateTokenControllerResponse;
    fromPartial<I_1 extends {} & {} & { [K_1 in Exclude<keyof I_1, never>]: never; }>(_: I_1): MsgUpdateTokenControllerResponse;
};
export declare const MsgUpdatePauser: {
    encode(message: MsgUpdatePauser, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdatePauser;
    fromJSON(object: any): MsgUpdatePauser;
    toJSON(message: MsgUpdatePauser): unknown;
    create<I extends {
        from?: string | undefined;
        newPauser?: string | undefined;
    } & {
        from?: string | undefined;
        newPauser?: string | undefined;
    } & { [K in Exclude<keyof I, keyof MsgUpdatePauser>]: never; }>(base?: I | undefined): MsgUpdatePauser;
    fromPartial<I_1 extends {
        from?: string | undefined;
        newPauser?: string | undefined;
    } & {
        from?: string | undefined;
        newPauser?: string | undefined;
    } & { [K_1 in Exclude<keyof I_1, keyof MsgUpdatePauser>]: never; }>(object: I_1): MsgUpdatePauser;
};
export declare const MsgUpdatePauserResponse: {
    encode(_: MsgUpdatePauserResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdatePauserResponse;
    fromJSON(_: any): MsgUpdatePauserResponse;
    toJSON(_: MsgUpdatePauserResponse): unknown;
    create<I extends {} & {} & { [K in Exclude<keyof I, never>]: never; }>(base?: I | undefined): MsgUpdatePauserResponse;
    fromPartial<I_1 extends {} & {} & { [K_1 in Exclude<keyof I_1, never>]: never; }>(_: I_1): MsgUpdatePauserResponse;
};
export declare const MsgAcceptOwner: {
    encode(message: MsgAcceptOwner, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgAcceptOwner;
    fromJSON(object: any): MsgAcceptOwner;
    toJSON(message: MsgAcceptOwner): unknown;
    create<I extends {
        from?: string | undefined;
    } & {
        from?: string | undefined;
    } & { [K in Exclude<keyof I, "from">]: never; }>(base?: I | undefined): MsgAcceptOwner;
    fromPartial<I_1 extends {
        from?: string | undefined;
    } & {
        from?: string | undefined;
    } & { [K_1 in Exclude<keyof I_1, "from">]: never; }>(object: I_1): MsgAcceptOwner;
};
export declare const MsgAcceptOwnerResponse: {
    encode(_: MsgAcceptOwnerResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgAcceptOwnerResponse;
    fromJSON(_: any): MsgAcceptOwnerResponse;
    toJSON(_: MsgAcceptOwnerResponse): unknown;
    create<I extends {} & {} & { [K in Exclude<keyof I, never>]: never; }>(base?: I | undefined): MsgAcceptOwnerResponse;
    fromPartial<I_1 extends {} & {} & { [K_1 in Exclude<keyof I_1, never>]: never; }>(_: I_1): MsgAcceptOwnerResponse;
};
export declare const MsgEnableAttester: {
    encode(message: MsgEnableAttester, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgEnableAttester;
    fromJSON(object: any): MsgEnableAttester;
    toJSON(message: MsgEnableAttester): unknown;
    create<I extends {
        from?: string | undefined;
        attester?: string | undefined;
    } & {
        from?: string | undefined;
        attester?: string | undefined;
    } & { [K in Exclude<keyof I, keyof MsgEnableAttester>]: never; }>(base?: I | undefined): MsgEnableAttester;
    fromPartial<I_1 extends {
        from?: string | undefined;
        attester?: string | undefined;
    } & {
        from?: string | undefined;
        attester?: string | undefined;
    } & { [K_1 in Exclude<keyof I_1, keyof MsgEnableAttester>]: never; }>(object: I_1): MsgEnableAttester;
};
export declare const MsgEnableAttesterResponse: {
    encode(_: MsgEnableAttesterResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgEnableAttesterResponse;
    fromJSON(_: any): MsgEnableAttesterResponse;
    toJSON(_: MsgEnableAttesterResponse): unknown;
    create<I extends {} & {} & { [K in Exclude<keyof I, never>]: never; }>(base?: I | undefined): MsgEnableAttesterResponse;
    fromPartial<I_1 extends {} & {} & { [K_1 in Exclude<keyof I_1, never>]: never; }>(_: I_1): MsgEnableAttesterResponse;
};
export declare const MsgDisableAttester: {
    encode(message: MsgDisableAttester, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgDisableAttester;
    fromJSON(object: any): MsgDisableAttester;
    toJSON(message: MsgDisableAttester): unknown;
    create<I extends {
        from?: string | undefined;
        attester?: string | undefined;
    } & {
        from?: string | undefined;
        attester?: string | undefined;
    } & { [K in Exclude<keyof I, keyof MsgDisableAttester>]: never; }>(base?: I | undefined): MsgDisableAttester;
    fromPartial<I_1 extends {
        from?: string | undefined;
        attester?: string | undefined;
    } & {
        from?: string | undefined;
        attester?: string | undefined;
    } & { [K_1 in Exclude<keyof I_1, keyof MsgDisableAttester>]: never; }>(object: I_1): MsgDisableAttester;
};
export declare const MsgDisableAttesterResponse: {
    encode(_: MsgDisableAttesterResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgDisableAttesterResponse;
    fromJSON(_: any): MsgDisableAttesterResponse;
    toJSON(_: MsgDisableAttesterResponse): unknown;
    create<I extends {} & {} & { [K in Exclude<keyof I, never>]: never; }>(base?: I | undefined): MsgDisableAttesterResponse;
    fromPartial<I_1 extends {} & {} & { [K_1 in Exclude<keyof I_1, never>]: never; }>(_: I_1): MsgDisableAttesterResponse;
};
export declare const MsgPauseBurningAndMinting: {
    encode(message: MsgPauseBurningAndMinting, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgPauseBurningAndMinting;
    fromJSON(object: any): MsgPauseBurningAndMinting;
    toJSON(message: MsgPauseBurningAndMinting): unknown;
    create<I extends {
        from?: string | undefined;
    } & {
        from?: string | undefined;
    } & { [K in Exclude<keyof I, "from">]: never; }>(base?: I | undefined): MsgPauseBurningAndMinting;
    fromPartial<I_1 extends {
        from?: string | undefined;
    } & {
        from?: string | undefined;
    } & { [K_1 in Exclude<keyof I_1, "from">]: never; }>(object: I_1): MsgPauseBurningAndMinting;
};
export declare const MsgPauseBurningAndMintingResponse: {
    encode(_: MsgPauseBurningAndMintingResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgPauseBurningAndMintingResponse;
    fromJSON(_: any): MsgPauseBurningAndMintingResponse;
    toJSON(_: MsgPauseBurningAndMintingResponse): unknown;
    create<I extends {} & {} & { [K in Exclude<keyof I, never>]: never; }>(base?: I | undefined): MsgPauseBurningAndMintingResponse;
    fromPartial<I_1 extends {} & {} & { [K_1 in Exclude<keyof I_1, never>]: never; }>(_: I_1): MsgPauseBurningAndMintingResponse;
};
export declare const MsgUnpauseBurningAndMinting: {
    encode(message: MsgUnpauseBurningAndMinting, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUnpauseBurningAndMinting;
    fromJSON(object: any): MsgUnpauseBurningAndMinting;
    toJSON(message: MsgUnpauseBurningAndMinting): unknown;
    create<I extends {
        from?: string | undefined;
    } & {
        from?: string | undefined;
    } & { [K in Exclude<keyof I, "from">]: never; }>(base?: I | undefined): MsgUnpauseBurningAndMinting;
    fromPartial<I_1 extends {
        from?: string | undefined;
    } & {
        from?: string | undefined;
    } & { [K_1 in Exclude<keyof I_1, "from">]: never; }>(object: I_1): MsgUnpauseBurningAndMinting;
};
export declare const MsgUnpauseBurningAndMintingResponse: {
    encode(_: MsgUnpauseBurningAndMintingResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUnpauseBurningAndMintingResponse;
    fromJSON(_: any): MsgUnpauseBurningAndMintingResponse;
    toJSON(_: MsgUnpauseBurningAndMintingResponse): unknown;
    create<I extends {} & {} & { [K in Exclude<keyof I, never>]: never; }>(base?: I | undefined): MsgUnpauseBurningAndMintingResponse;
    fromPartial<I_1 extends {} & {} & { [K_1 in Exclude<keyof I_1, never>]: never; }>(_: I_1): MsgUnpauseBurningAndMintingResponse;
};
export declare const MsgPauseSendingAndReceivingMessages: {
    encode(message: MsgPauseSendingAndReceivingMessages, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgPauseSendingAndReceivingMessages;
    fromJSON(object: any): MsgPauseSendingAndReceivingMessages;
    toJSON(message: MsgPauseSendingAndReceivingMessages): unknown;
    create<I extends {
        from?: string | undefined;
    } & {
        from?: string | undefined;
    } & { [K in Exclude<keyof I, "from">]: never; }>(base?: I | undefined): MsgPauseSendingAndReceivingMessages;
    fromPartial<I_1 extends {
        from?: string | undefined;
    } & {
        from?: string | undefined;
    } & { [K_1 in Exclude<keyof I_1, "from">]: never; }>(object: I_1): MsgPauseSendingAndReceivingMessages;
};
export declare const MsgPauseSendingAndReceivingMessagesResponse: {
    encode(_: MsgPauseSendingAndReceivingMessagesResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgPauseSendingAndReceivingMessagesResponse;
    fromJSON(_: any): MsgPauseSendingAndReceivingMessagesResponse;
    toJSON(_: MsgPauseSendingAndReceivingMessagesResponse): unknown;
    create<I extends {} & {} & { [K in Exclude<keyof I, never>]: never; }>(base?: I | undefined): MsgPauseSendingAndReceivingMessagesResponse;
    fromPartial<I_1 extends {} & {} & { [K_1 in Exclude<keyof I_1, never>]: never; }>(_: I_1): MsgPauseSendingAndReceivingMessagesResponse;
};
export declare const MsgUnpauseSendingAndReceivingMessages: {
    encode(message: MsgUnpauseSendingAndReceivingMessages, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUnpauseSendingAndReceivingMessages;
    fromJSON(object: any): MsgUnpauseSendingAndReceivingMessages;
    toJSON(message: MsgUnpauseSendingAndReceivingMessages): unknown;
    create<I extends {
        from?: string | undefined;
    } & {
        from?: string | undefined;
    } & { [K in Exclude<keyof I, "from">]: never; }>(base?: I | undefined): MsgUnpauseSendingAndReceivingMessages;
    fromPartial<I_1 extends {
        from?: string | undefined;
    } & {
        from?: string | undefined;
    } & { [K_1 in Exclude<keyof I_1, "from">]: never; }>(object: I_1): MsgUnpauseSendingAndReceivingMessages;
};
export declare const MsgUnpauseSendingAndReceivingMessagesResponse: {
    encode(_: MsgUnpauseSendingAndReceivingMessagesResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUnpauseSendingAndReceivingMessagesResponse;
    fromJSON(_: any): MsgUnpauseSendingAndReceivingMessagesResponse;
    toJSON(_: MsgUnpauseSendingAndReceivingMessagesResponse): unknown;
    create<I extends {} & {} & { [K in Exclude<keyof I, never>]: never; }>(base?: I | undefined): MsgUnpauseSendingAndReceivingMessagesResponse;
    fromPartial<I_1 extends {} & {} & { [K_1 in Exclude<keyof I_1, never>]: never; }>(_: I_1): MsgUnpauseSendingAndReceivingMessagesResponse;
};
export declare const MsgUpdateMaxMessageBodySize: {
    encode(message: MsgUpdateMaxMessageBodySize, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateMaxMessageBodySize;
    fromJSON(object: any): MsgUpdateMaxMessageBodySize;
    toJSON(message: MsgUpdateMaxMessageBodySize): unknown;
    create<I extends {
        from?: string | undefined;
        messageSize?: string | number | Long.Long | undefined;
    } & {
        from?: string | undefined;
        messageSize?: string | number | (Long.Long & {
            high: number;
            low: number;
            unsigned: boolean;
            add: (addend: string | number | Long.Long) => Long.Long;
            and: (other: string | number | Long.Long) => Long.Long;
            compare: (other: string | number | Long.Long) => number;
            comp: (other: string | number | Long.Long) => number;
            divide: (divisor: string | number | Long.Long) => Long.Long;
            div: (divisor: string | number | Long.Long) => Long.Long;
            equals: (other: string | number | Long.Long) => boolean;
            eq: (other: string | number | Long.Long) => boolean;
            getHighBits: () => number;
            getHighBitsUnsigned: () => number;
            getLowBits: () => number;
            getLowBitsUnsigned: () => number;
            getNumBitsAbs: () => number;
            greaterThan: (other: string | number | Long.Long) => boolean;
            gt: (other: string | number | Long.Long) => boolean;
            greaterThanOrEqual: (other: string | number | Long.Long) => boolean;
            gte: (other: string | number | Long.Long) => boolean;
            isEven: () => boolean;
            isNegative: () => boolean;
            isOdd: () => boolean;
            isPositive: () => boolean;
            isZero: () => boolean;
            lessThan: (other: string | number | Long.Long) => boolean;
            lt: (other: string | number | Long.Long) => boolean;
            lessThanOrEqual: (other: string | number | Long.Long) => boolean;
            lte: (other: string | number | Long.Long) => boolean;
            modulo: (other: string | number | Long.Long) => Long.Long;
            mod: (other: string | number | Long.Long) => Long.Long;
            multiply: (multiplier: string | number | Long.Long) => Long.Long;
            mul: (multiplier: string | number | Long.Long) => Long.Long;
            negate: () => Long.Long;
            neg: () => Long.Long;
            not: () => Long.Long;
            notEquals: (other: string | number | Long.Long) => boolean;
            neq: (other: string | number | Long.Long) => boolean;
            or: (other: string | number | Long.Long) => Long.Long;
            shiftLeft: (numBits: number | Long.Long) => Long.Long;
            shl: (numBits: number | Long.Long) => Long.Long;
            shiftRight: (numBits: number | Long.Long) => Long.Long;
            shr: (numBits: number | Long.Long) => Long.Long;
            shiftRightUnsigned: (numBits: number | Long.Long) => Long.Long;
            shru: (numBits: number | Long.Long) => Long.Long;
            subtract: (subtrahend: string | number | Long.Long) => Long.Long;
            sub: (subtrahend: string | number | Long.Long) => Long.Long;
            toInt: () => number;
            toNumber: () => number;
            toBytes: (le?: boolean | undefined) => number[];
            toBytesLE: () => number[];
            toBytesBE: () => number[];
            toSigned: () => Long.Long;
            toString: (radix?: number | undefined) => string;
            toUnsigned: () => Long.Long;
            xor: (other: string | number | Long.Long) => Long.Long;
        } & { [K in Exclude<keyof I["messageSize"], keyof Long.Long>]: never; }) | undefined;
    } & { [K_1 in Exclude<keyof I, keyof MsgUpdateMaxMessageBodySize>]: never; }>(base?: I | undefined): MsgUpdateMaxMessageBodySize;
    fromPartial<I_1 extends {
        from?: string | undefined;
        messageSize?: string | number | Long.Long | undefined;
    } & {
        from?: string | undefined;
        messageSize?: string | number | (Long.Long & {
            high: number;
            low: number;
            unsigned: boolean;
            add: (addend: string | number | Long.Long) => Long.Long;
            and: (other: string | number | Long.Long) => Long.Long;
            compare: (other: string | number | Long.Long) => number;
            comp: (other: string | number | Long.Long) => number;
            divide: (divisor: string | number | Long.Long) => Long.Long;
            div: (divisor: string | number | Long.Long) => Long.Long;
            equals: (other: string | number | Long.Long) => boolean;
            eq: (other: string | number | Long.Long) => boolean;
            getHighBits: () => number;
            getHighBitsUnsigned: () => number;
            getLowBits: () => number;
            getLowBitsUnsigned: () => number;
            getNumBitsAbs: () => number;
            greaterThan: (other: string | number | Long.Long) => boolean;
            gt: (other: string | number | Long.Long) => boolean;
            greaterThanOrEqual: (other: string | number | Long.Long) => boolean;
            gte: (other: string | number | Long.Long) => boolean;
            isEven: () => boolean;
            isNegative: () => boolean;
            isOdd: () => boolean;
            isPositive: () => boolean;
            isZero: () => boolean;
            lessThan: (other: string | number | Long.Long) => boolean;
            lt: (other: string | number | Long.Long) => boolean;
            lessThanOrEqual: (other: string | number | Long.Long) => boolean;
            lte: (other: string | number | Long.Long) => boolean;
            modulo: (other: string | number | Long.Long) => Long.Long;
            mod: (other: string | number | Long.Long) => Long.Long;
            multiply: (multiplier: string | number | Long.Long) => Long.Long;
            mul: (multiplier: string | number | Long.Long) => Long.Long;
            negate: () => Long.Long;
            neg: () => Long.Long;
            not: () => Long.Long;
            notEquals: (other: string | number | Long.Long) => boolean;
            neq: (other: string | number | Long.Long) => boolean;
            or: (other: string | number | Long.Long) => Long.Long;
            shiftLeft: (numBits: number | Long.Long) => Long.Long;
            shl: (numBits: number | Long.Long) => Long.Long;
            shiftRight: (numBits: number | Long.Long) => Long.Long;
            shr: (numBits: number | Long.Long) => Long.Long;
            shiftRightUnsigned: (numBits: number | Long.Long) => Long.Long;
            shru: (numBits: number | Long.Long) => Long.Long;
            subtract: (subtrahend: string | number | Long.Long) => Long.Long;
            sub: (subtrahend: string | number | Long.Long) => Long.Long;
            toInt: () => number;
            toNumber: () => number;
            toBytes: (le?: boolean | undefined) => number[];
            toBytesLE: () => number[];
            toBytesBE: () => number[];
            toSigned: () => Long.Long;
            toString: (radix?: number | undefined) => string;
            toUnsigned: () => Long.Long;
            xor: (other: string | number | Long.Long) => Long.Long;
        } & { [K_2 in Exclude<keyof I_1["messageSize"], keyof Long.Long>]: never; }) | undefined;
    } & { [K_3 in Exclude<keyof I_1, keyof MsgUpdateMaxMessageBodySize>]: never; }>(object: I_1): MsgUpdateMaxMessageBodySize;
};
export declare const MsgUpdateMaxMessageBodySizeResponse: {
    encode(_: MsgUpdateMaxMessageBodySizeResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateMaxMessageBodySizeResponse;
    fromJSON(_: any): MsgUpdateMaxMessageBodySizeResponse;
    toJSON(_: MsgUpdateMaxMessageBodySizeResponse): unknown;
    create<I extends {} & {} & { [K in Exclude<keyof I, never>]: never; }>(base?: I | undefined): MsgUpdateMaxMessageBodySizeResponse;
    fromPartial<I_1 extends {} & {} & { [K_1 in Exclude<keyof I_1, never>]: never; }>(_: I_1): MsgUpdateMaxMessageBodySizeResponse;
};
export declare const MsgSetMaxBurnAmountPerMessage: {
    encode(message: MsgSetMaxBurnAmountPerMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetMaxBurnAmountPerMessage;
    fromJSON(object: any): MsgSetMaxBurnAmountPerMessage;
    toJSON(message: MsgSetMaxBurnAmountPerMessage): unknown;
    create<I extends {
        from?: string | undefined;
        localToken?: string | undefined;
        amount?: string | undefined;
    } & {
        from?: string | undefined;
        localToken?: string | undefined;
        amount?: string | undefined;
    } & { [K in Exclude<keyof I, keyof MsgSetMaxBurnAmountPerMessage>]: never; }>(base?: I | undefined): MsgSetMaxBurnAmountPerMessage;
    fromPartial<I_1 extends {
        from?: string | undefined;
        localToken?: string | undefined;
        amount?: string | undefined;
    } & {
        from?: string | undefined;
        localToken?: string | undefined;
        amount?: string | undefined;
    } & { [K_1 in Exclude<keyof I_1, keyof MsgSetMaxBurnAmountPerMessage>]: never; }>(object: I_1): MsgSetMaxBurnAmountPerMessage;
};
export declare const MsgSetMaxBurnAmountPerMessageResponse: {
    encode(_: MsgSetMaxBurnAmountPerMessageResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetMaxBurnAmountPerMessageResponse;
    fromJSON(_: any): MsgSetMaxBurnAmountPerMessageResponse;
    toJSON(_: MsgSetMaxBurnAmountPerMessageResponse): unknown;
    create<I extends {} & {} & { [K in Exclude<keyof I, never>]: never; }>(base?: I | undefined): MsgSetMaxBurnAmountPerMessageResponse;
    fromPartial<I_1 extends {} & {} & { [K_1 in Exclude<keyof I_1, never>]: never; }>(_: I_1): MsgSetMaxBurnAmountPerMessageResponse;
};
export declare const MsgDepositForBurn: {
    encode(message: MsgDepositForBurn, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgDepositForBurn;
    fromJSON(object: any): MsgDepositForBurn;
    toJSON(message: MsgDepositForBurn): unknown;
    create<I extends {
        from?: string | undefined;
        amount?: string | undefined;
        destinationDomain?: number | undefined;
        mintRecipient?: Uint8Array | undefined;
        burnToken?: string | undefined;
    } & {
        from?: string | undefined;
        amount?: string | undefined;
        destinationDomain?: number | undefined;
        mintRecipient?: Uint8Array | undefined;
        burnToken?: string | undefined;
    } & { [K in Exclude<keyof I, keyof MsgDepositForBurn>]: never; }>(base?: I | undefined): MsgDepositForBurn;
    fromPartial<I_1 extends {
        from?: string | undefined;
        amount?: string | undefined;
        destinationDomain?: number | undefined;
        mintRecipient?: Uint8Array | undefined;
        burnToken?: string | undefined;
    } & {
        from?: string | undefined;
        amount?: string | undefined;
        destinationDomain?: number | undefined;
        mintRecipient?: Uint8Array | undefined;
        burnToken?: string | undefined;
    } & { [K_1 in Exclude<keyof I_1, keyof MsgDepositForBurn>]: never; }>(object: I_1): MsgDepositForBurn;
};
export declare const MsgDepositForBurnResponse: {
    encode(message: MsgDepositForBurnResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgDepositForBurnResponse;
    fromJSON(object: any): MsgDepositForBurnResponse;
    toJSON(message: MsgDepositForBurnResponse): unknown;
    create<I extends {
        nonce?: string | number | Long.Long | undefined;
    } & {
        nonce?: string | number | (Long.Long & {
            high: number;
            low: number;
            unsigned: boolean;
            add: (addend: string | number | Long.Long) => Long.Long;
            and: (other: string | number | Long.Long) => Long.Long;
            compare: (other: string | number | Long.Long) => number;
            comp: (other: string | number | Long.Long) => number;
            divide: (divisor: string | number | Long.Long) => Long.Long;
            div: (divisor: string | number | Long.Long) => Long.Long;
            equals: (other: string | number | Long.Long) => boolean;
            eq: (other: string | number | Long.Long) => boolean;
            getHighBits: () => number;
            getHighBitsUnsigned: () => number;
            getLowBits: () => number;
            getLowBitsUnsigned: () => number;
            getNumBitsAbs: () => number;
            greaterThan: (other: string | number | Long.Long) => boolean;
            gt: (other: string | number | Long.Long) => boolean;
            greaterThanOrEqual: (other: string | number | Long.Long) => boolean;
            gte: (other: string | number | Long.Long) => boolean;
            isEven: () => boolean;
            isNegative: () => boolean;
            isOdd: () => boolean;
            isPositive: () => boolean;
            isZero: () => boolean;
            lessThan: (other: string | number | Long.Long) => boolean;
            lt: (other: string | number | Long.Long) => boolean;
            lessThanOrEqual: (other: string | number | Long.Long) => boolean;
            lte: (other: string | number | Long.Long) => boolean;
            modulo: (other: string | number | Long.Long) => Long.Long;
            mod: (other: string | number | Long.Long) => Long.Long;
            multiply: (multiplier: string | number | Long.Long) => Long.Long;
            mul: (multiplier: string | number | Long.Long) => Long.Long;
            negate: () => Long.Long;
            neg: () => Long.Long;
            not: () => Long.Long;
            notEquals: (other: string | number | Long.Long) => boolean;
            neq: (other: string | number | Long.Long) => boolean;
            or: (other: string | number | Long.Long) => Long.Long;
            shiftLeft: (numBits: number | Long.Long) => Long.Long;
            shl: (numBits: number | Long.Long) => Long.Long;
            shiftRight: (numBits: number | Long.Long) => Long.Long;
            shr: (numBits: number | Long.Long) => Long.Long;
            shiftRightUnsigned: (numBits: number | Long.Long) => Long.Long;
            shru: (numBits: number | Long.Long) => Long.Long;
            subtract: (subtrahend: string | number | Long.Long) => Long.Long;
            sub: (subtrahend: string | number | Long.Long) => Long.Long;
            toInt: () => number;
            toNumber: () => number;
            toBytes: (le?: boolean | undefined) => number[];
            toBytesLE: () => number[];
            toBytesBE: () => number[];
            toSigned: () => Long.Long;
            toString: (radix?: number | undefined) => string;
            toUnsigned: () => Long.Long;
            xor: (other: string | number | Long.Long) => Long.Long;
        } & { [K in Exclude<keyof I["nonce"], keyof Long.Long>]: never; }) | undefined;
    } & { [K_1 in Exclude<keyof I, "nonce">]: never; }>(base?: I | undefined): MsgDepositForBurnResponse;
    fromPartial<I_1 extends {
        nonce?: string | number | Long.Long | undefined;
    } & {
        nonce?: string | number | (Long.Long & {
            high: number;
            low: number;
            unsigned: boolean;
            add: (addend: string | number | Long.Long) => Long.Long;
            and: (other: string | number | Long.Long) => Long.Long;
            compare: (other: string | number | Long.Long) => number;
            comp: (other: string | number | Long.Long) => number;
            divide: (divisor: string | number | Long.Long) => Long.Long;
            div: (divisor: string | number | Long.Long) => Long.Long;
            equals: (other: string | number | Long.Long) => boolean;
            eq: (other: string | number | Long.Long) => boolean;
            getHighBits: () => number;
            getHighBitsUnsigned: () => number;
            getLowBits: () => number;
            getLowBitsUnsigned: () => number;
            getNumBitsAbs: () => number;
            greaterThan: (other: string | number | Long.Long) => boolean;
            gt: (other: string | number | Long.Long) => boolean;
            greaterThanOrEqual: (other: string | number | Long.Long) => boolean;
            gte: (other: string | number | Long.Long) => boolean;
            isEven: () => boolean;
            isNegative: () => boolean;
            isOdd: () => boolean;
            isPositive: () => boolean;
            isZero: () => boolean;
            lessThan: (other: string | number | Long.Long) => boolean;
            lt: (other: string | number | Long.Long) => boolean;
            lessThanOrEqual: (other: string | number | Long.Long) => boolean;
            lte: (other: string | number | Long.Long) => boolean;
            modulo: (other: string | number | Long.Long) => Long.Long;
            mod: (other: string | number | Long.Long) => Long.Long;
            multiply: (multiplier: string | number | Long.Long) => Long.Long;
            mul: (multiplier: string | number | Long.Long) => Long.Long;
            negate: () => Long.Long;
            neg: () => Long.Long;
            not: () => Long.Long;
            notEquals: (other: string | number | Long.Long) => boolean;
            neq: (other: string | number | Long.Long) => boolean;
            or: (other: string | number | Long.Long) => Long.Long;
            shiftLeft: (numBits: number | Long.Long) => Long.Long;
            shl: (numBits: number | Long.Long) => Long.Long;
            shiftRight: (numBits: number | Long.Long) => Long.Long;
            shr: (numBits: number | Long.Long) => Long.Long;
            shiftRightUnsigned: (numBits: number | Long.Long) => Long.Long;
            shru: (numBits: number | Long.Long) => Long.Long;
            subtract: (subtrahend: string | number | Long.Long) => Long.Long;
            sub: (subtrahend: string | number | Long.Long) => Long.Long;
            toInt: () => number;
            toNumber: () => number;
            toBytes: (le?: boolean | undefined) => number[];
            toBytesLE: () => number[];
            toBytesBE: () => number[];
            toSigned: () => Long.Long;
            toString: (radix?: number | undefined) => string;
            toUnsigned: () => Long.Long;
            xor: (other: string | number | Long.Long) => Long.Long;
        } & { [K_2 in Exclude<keyof I_1["nonce"], keyof Long.Long>]: never; }) | undefined;
    } & { [K_3 in Exclude<keyof I_1, "nonce">]: never; }>(object: I_1): MsgDepositForBurnResponse;
};
export declare const MsgDepositForBurnWithCaller: {
    encode(message: MsgDepositForBurnWithCaller, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgDepositForBurnWithCaller;
    fromJSON(object: any): MsgDepositForBurnWithCaller;
    toJSON(message: MsgDepositForBurnWithCaller): unknown;
    create<I extends {
        from?: string | undefined;
        amount?: string | undefined;
        destinationDomain?: number | undefined;
        mintRecipient?: Uint8Array | undefined;
        burnToken?: string | undefined;
        destinationCaller?: Uint8Array | undefined;
    } & {
        from?: string | undefined;
        amount?: string | undefined;
        destinationDomain?: number | undefined;
        mintRecipient?: Uint8Array | undefined;
        burnToken?: string | undefined;
        destinationCaller?: Uint8Array | undefined;
    } & { [K in Exclude<keyof I, keyof MsgDepositForBurnWithCaller>]: never; }>(base?: I | undefined): MsgDepositForBurnWithCaller;
    fromPartial<I_1 extends {
        from?: string | undefined;
        amount?: string | undefined;
        destinationDomain?: number | undefined;
        mintRecipient?: Uint8Array | undefined;
        burnToken?: string | undefined;
        destinationCaller?: Uint8Array | undefined;
    } & {
        from?: string | undefined;
        amount?: string | undefined;
        destinationDomain?: number | undefined;
        mintRecipient?: Uint8Array | undefined;
        burnToken?: string | undefined;
        destinationCaller?: Uint8Array | undefined;
    } & { [K_1 in Exclude<keyof I_1, keyof MsgDepositForBurnWithCaller>]: never; }>(object: I_1): MsgDepositForBurnWithCaller;
};
export declare const MsgDepositForBurnWithCallerResponse: {
    encode(message: MsgDepositForBurnWithCallerResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgDepositForBurnWithCallerResponse;
    fromJSON(object: any): MsgDepositForBurnWithCallerResponse;
    toJSON(message: MsgDepositForBurnWithCallerResponse): unknown;
    create<I extends {
        nonce?: string | number | Long.Long | undefined;
    } & {
        nonce?: string | number | (Long.Long & {
            high: number;
            low: number;
            unsigned: boolean;
            add: (addend: string | number | Long.Long) => Long.Long;
            and: (other: string | number | Long.Long) => Long.Long;
            compare: (other: string | number | Long.Long) => number;
            comp: (other: string | number | Long.Long) => number;
            divide: (divisor: string | number | Long.Long) => Long.Long;
            div: (divisor: string | number | Long.Long) => Long.Long;
            equals: (other: string | number | Long.Long) => boolean;
            eq: (other: string | number | Long.Long) => boolean;
            getHighBits: () => number;
            getHighBitsUnsigned: () => number;
            getLowBits: () => number;
            getLowBitsUnsigned: () => number;
            getNumBitsAbs: () => number;
            greaterThan: (other: string | number | Long.Long) => boolean;
            gt: (other: string | number | Long.Long) => boolean;
            greaterThanOrEqual: (other: string | number | Long.Long) => boolean;
            gte: (other: string | number | Long.Long) => boolean;
            isEven: () => boolean;
            isNegative: () => boolean;
            isOdd: () => boolean;
            isPositive: () => boolean;
            isZero: () => boolean;
            lessThan: (other: string | number | Long.Long) => boolean;
            lt: (other: string | number | Long.Long) => boolean;
            lessThanOrEqual: (other: string | number | Long.Long) => boolean;
            lte: (other: string | number | Long.Long) => boolean;
            modulo: (other: string | number | Long.Long) => Long.Long;
            mod: (other: string | number | Long.Long) => Long.Long;
            multiply: (multiplier: string | number | Long.Long) => Long.Long;
            mul: (multiplier: string | number | Long.Long) => Long.Long;
            negate: () => Long.Long;
            neg: () => Long.Long;
            not: () => Long.Long;
            notEquals: (other: string | number | Long.Long) => boolean;
            neq: (other: string | number | Long.Long) => boolean;
            or: (other: string | number | Long.Long) => Long.Long;
            shiftLeft: (numBits: number | Long.Long) => Long.Long;
            shl: (numBits: number | Long.Long) => Long.Long;
            shiftRight: (numBits: number | Long.Long) => Long.Long;
            shr: (numBits: number | Long.Long) => Long.Long;
            shiftRightUnsigned: (numBits: number | Long.Long) => Long.Long;
            shru: (numBits: number | Long.Long) => Long.Long;
            subtract: (subtrahend: string | number | Long.Long) => Long.Long;
            sub: (subtrahend: string | number | Long.Long) => Long.Long;
            toInt: () => number;
            toNumber: () => number;
            toBytes: (le?: boolean | undefined) => number[];
            toBytesLE: () => number[];
            toBytesBE: () => number[];
            toSigned: () => Long.Long;
            toString: (radix?: number | undefined) => string;
            toUnsigned: () => Long.Long;
            xor: (other: string | number | Long.Long) => Long.Long;
        } & { [K in Exclude<keyof I["nonce"], keyof Long.Long>]: never; }) | undefined;
    } & { [K_1 in Exclude<keyof I, "nonce">]: never; }>(base?: I | undefined): MsgDepositForBurnWithCallerResponse;
    fromPartial<I_1 extends {
        nonce?: string | number | Long.Long | undefined;
    } & {
        nonce?: string | number | (Long.Long & {
            high: number;
            low: number;
            unsigned: boolean;
            add: (addend: string | number | Long.Long) => Long.Long;
            and: (other: string | number | Long.Long) => Long.Long;
            compare: (other: string | number | Long.Long) => number;
            comp: (other: string | number | Long.Long) => number;
            divide: (divisor: string | number | Long.Long) => Long.Long;
            div: (divisor: string | number | Long.Long) => Long.Long;
            equals: (other: string | number | Long.Long) => boolean;
            eq: (other: string | number | Long.Long) => boolean;
            getHighBits: () => number;
            getHighBitsUnsigned: () => number;
            getLowBits: () => number;
            getLowBitsUnsigned: () => number;
            getNumBitsAbs: () => number;
            greaterThan: (other: string | number | Long.Long) => boolean;
            gt: (other: string | number | Long.Long) => boolean;
            greaterThanOrEqual: (other: string | number | Long.Long) => boolean;
            gte: (other: string | number | Long.Long) => boolean;
            isEven: () => boolean;
            isNegative: () => boolean;
            isOdd: () => boolean;
            isPositive: () => boolean;
            isZero: () => boolean;
            lessThan: (other: string | number | Long.Long) => boolean;
            lt: (other: string | number | Long.Long) => boolean;
            lessThanOrEqual: (other: string | number | Long.Long) => boolean;
            lte: (other: string | number | Long.Long) => boolean;
            modulo: (other: string | number | Long.Long) => Long.Long;
            mod: (other: string | number | Long.Long) => Long.Long;
            multiply: (multiplier: string | number | Long.Long) => Long.Long;
            mul: (multiplier: string | number | Long.Long) => Long.Long;
            negate: () => Long.Long;
            neg: () => Long.Long;
            not: () => Long.Long;
            notEquals: (other: string | number | Long.Long) => boolean;
            neq: (other: string | number | Long.Long) => boolean;
            or: (other: string | number | Long.Long) => Long.Long;
            shiftLeft: (numBits: number | Long.Long) => Long.Long;
            shl: (numBits: number | Long.Long) => Long.Long;
            shiftRight: (numBits: number | Long.Long) => Long.Long;
            shr: (numBits: number | Long.Long) => Long.Long;
            shiftRightUnsigned: (numBits: number | Long.Long) => Long.Long;
            shru: (numBits: number | Long.Long) => Long.Long;
            subtract: (subtrahend: string | number | Long.Long) => Long.Long;
            sub: (subtrahend: string | number | Long.Long) => Long.Long;
            toInt: () => number;
            toNumber: () => number;
            toBytes: (le?: boolean | undefined) => number[];
            toBytesLE: () => number[];
            toBytesBE: () => number[];
            toSigned: () => Long.Long;
            toString: (radix?: number | undefined) => string;
            toUnsigned: () => Long.Long;
            xor: (other: string | number | Long.Long) => Long.Long;
        } & { [K_2 in Exclude<keyof I_1["nonce"], keyof Long.Long>]: never; }) | undefined;
    } & { [K_3 in Exclude<keyof I_1, "nonce">]: never; }>(object: I_1): MsgDepositForBurnWithCallerResponse;
};
export declare const MsgReplaceDepositForBurn: {
    encode(message: MsgReplaceDepositForBurn, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgReplaceDepositForBurn;
    fromJSON(object: any): MsgReplaceDepositForBurn;
    toJSON(message: MsgReplaceDepositForBurn): unknown;
    create<I extends {
        from?: string | undefined;
        originalMessage?: Uint8Array | undefined;
        originalAttestation?: Uint8Array | undefined;
        newDestinationCaller?: Uint8Array | undefined;
        newMintRecipient?: Uint8Array | undefined;
    } & {
        from?: string | undefined;
        originalMessage?: Uint8Array | undefined;
        originalAttestation?: Uint8Array | undefined;
        newDestinationCaller?: Uint8Array | undefined;
        newMintRecipient?: Uint8Array | undefined;
    } & { [K in Exclude<keyof I, keyof MsgReplaceDepositForBurn>]: never; }>(base?: I | undefined): MsgReplaceDepositForBurn;
    fromPartial<I_1 extends {
        from?: string | undefined;
        originalMessage?: Uint8Array | undefined;
        originalAttestation?: Uint8Array | undefined;
        newDestinationCaller?: Uint8Array | undefined;
        newMintRecipient?: Uint8Array | undefined;
    } & {
        from?: string | undefined;
        originalMessage?: Uint8Array | undefined;
        originalAttestation?: Uint8Array | undefined;
        newDestinationCaller?: Uint8Array | undefined;
        newMintRecipient?: Uint8Array | undefined;
    } & { [K_1 in Exclude<keyof I_1, keyof MsgReplaceDepositForBurn>]: never; }>(object: I_1): MsgReplaceDepositForBurn;
};
export declare const MsgReplaceDepositForBurnResponse: {
    encode(_: MsgReplaceDepositForBurnResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgReplaceDepositForBurnResponse;
    fromJSON(_: any): MsgReplaceDepositForBurnResponse;
    toJSON(_: MsgReplaceDepositForBurnResponse): unknown;
    create<I extends {} & {} & { [K in Exclude<keyof I, never>]: never; }>(base?: I | undefined): MsgReplaceDepositForBurnResponse;
    fromPartial<I_1 extends {} & {} & { [K_1 in Exclude<keyof I_1, never>]: never; }>(_: I_1): MsgReplaceDepositForBurnResponse;
};
export declare const MsgReceiveMessage: {
    encode(message: MsgReceiveMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgReceiveMessage;
    fromJSON(object: any): MsgReceiveMessage;
    toJSON(message: MsgReceiveMessage): unknown;
    create<I extends {
        from?: string | undefined;
        message?: Uint8Array | undefined;
        attestation?: Uint8Array | undefined;
    } & {
        from?: string | undefined;
        message?: Uint8Array | undefined;
        attestation?: Uint8Array | undefined;
    } & { [K in Exclude<keyof I, keyof MsgReceiveMessage>]: never; }>(base?: I | undefined): MsgReceiveMessage;
    fromPartial<I_1 extends {
        from?: string | undefined;
        message?: Uint8Array | undefined;
        attestation?: Uint8Array | undefined;
    } & {
        from?: string | undefined;
        message?: Uint8Array | undefined;
        attestation?: Uint8Array | undefined;
    } & { [K_1 in Exclude<keyof I_1, keyof MsgReceiveMessage>]: never; }>(object: I_1): MsgReceiveMessage;
};
export declare const MsgReceiveMessageResponse: {
    encode(message: MsgReceiveMessageResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgReceiveMessageResponse;
    fromJSON(object: any): MsgReceiveMessageResponse;
    toJSON(message: MsgReceiveMessageResponse): unknown;
    create<I extends {
        success?: boolean | undefined;
    } & {
        success?: boolean | undefined;
    } & { [K in Exclude<keyof I, "success">]: never; }>(base?: I | undefined): MsgReceiveMessageResponse;
    fromPartial<I_1 extends {
        success?: boolean | undefined;
    } & {
        success?: boolean | undefined;
    } & { [K_1 in Exclude<keyof I_1, "success">]: never; }>(object: I_1): MsgReceiveMessageResponse;
};
export declare const MsgSendMessage: {
    encode(message: MsgSendMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgSendMessage;
    fromJSON(object: any): MsgSendMessage;
    toJSON(message: MsgSendMessage): unknown;
    create<I extends {
        from?: string | undefined;
        destinationDomain?: number | undefined;
        recipient?: Uint8Array | undefined;
        messageBody?: Uint8Array | undefined;
    } & {
        from?: string | undefined;
        destinationDomain?: number | undefined;
        recipient?: Uint8Array | undefined;
        messageBody?: Uint8Array | undefined;
    } & { [K in Exclude<keyof I, keyof MsgSendMessage>]: never; }>(base?: I | undefined): MsgSendMessage;
    fromPartial<I_1 extends {
        from?: string | undefined;
        destinationDomain?: number | undefined;
        recipient?: Uint8Array | undefined;
        messageBody?: Uint8Array | undefined;
    } & {
        from?: string | undefined;
        destinationDomain?: number | undefined;
        recipient?: Uint8Array | undefined;
        messageBody?: Uint8Array | undefined;
    } & { [K_1 in Exclude<keyof I_1, keyof MsgSendMessage>]: never; }>(object: I_1): MsgSendMessage;
};
export declare const MsgSendMessageResponse: {
    encode(message: MsgSendMessageResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgSendMessageResponse;
    fromJSON(object: any): MsgSendMessageResponse;
    toJSON(message: MsgSendMessageResponse): unknown;
    create<I extends {
        nonce?: string | number | Long.Long | undefined;
    } & {
        nonce?: string | number | (Long.Long & {
            high: number;
            low: number;
            unsigned: boolean;
            add: (addend: string | number | Long.Long) => Long.Long;
            and: (other: string | number | Long.Long) => Long.Long;
            compare: (other: string | number | Long.Long) => number;
            comp: (other: string | number | Long.Long) => number;
            divide: (divisor: string | number | Long.Long) => Long.Long;
            div: (divisor: string | number | Long.Long) => Long.Long;
            equals: (other: string | number | Long.Long) => boolean;
            eq: (other: string | number | Long.Long) => boolean;
            getHighBits: () => number;
            getHighBitsUnsigned: () => number;
            getLowBits: () => number;
            getLowBitsUnsigned: () => number;
            getNumBitsAbs: () => number;
            greaterThan: (other: string | number | Long.Long) => boolean;
            gt: (other: string | number | Long.Long) => boolean;
            greaterThanOrEqual: (other: string | number | Long.Long) => boolean;
            gte: (other: string | number | Long.Long) => boolean;
            isEven: () => boolean;
            isNegative: () => boolean;
            isOdd: () => boolean;
            isPositive: () => boolean;
            isZero: () => boolean;
            lessThan: (other: string | number | Long.Long) => boolean;
            lt: (other: string | number | Long.Long) => boolean;
            lessThanOrEqual: (other: string | number | Long.Long) => boolean;
            lte: (other: string | number | Long.Long) => boolean;
            modulo: (other: string | number | Long.Long) => Long.Long;
            mod: (other: string | number | Long.Long) => Long.Long;
            multiply: (multiplier: string | number | Long.Long) => Long.Long;
            mul: (multiplier: string | number | Long.Long) => Long.Long;
            negate: () => Long.Long;
            neg: () => Long.Long;
            not: () => Long.Long;
            notEquals: (other: string | number | Long.Long) => boolean;
            neq: (other: string | number | Long.Long) => boolean;
            or: (other: string | number | Long.Long) => Long.Long;
            shiftLeft: (numBits: number | Long.Long) => Long.Long;
            shl: (numBits: number | Long.Long) => Long.Long;
            shiftRight: (numBits: number | Long.Long) => Long.Long;
            shr: (numBits: number | Long.Long) => Long.Long;
            shiftRightUnsigned: (numBits: number | Long.Long) => Long.Long;
            shru: (numBits: number | Long.Long) => Long.Long;
            subtract: (subtrahend: string | number | Long.Long) => Long.Long;
            sub: (subtrahend: string | number | Long.Long) => Long.Long;
            toInt: () => number;
            toNumber: () => number;
            toBytes: (le?: boolean | undefined) => number[];
            toBytesLE: () => number[];
            toBytesBE: () => number[];
            toSigned: () => Long.Long;
            toString: (radix?: number | undefined) => string;
            toUnsigned: () => Long.Long;
            xor: (other: string | number | Long.Long) => Long.Long;
        } & { [K in Exclude<keyof I["nonce"], keyof Long.Long>]: never; }) | undefined;
    } & { [K_1 in Exclude<keyof I, "nonce">]: never; }>(base?: I | undefined): MsgSendMessageResponse;
    fromPartial<I_1 extends {
        nonce?: string | number | Long.Long | undefined;
    } & {
        nonce?: string | number | (Long.Long & {
            high: number;
            low: number;
            unsigned: boolean;
            add: (addend: string | number | Long.Long) => Long.Long;
            and: (other: string | number | Long.Long) => Long.Long;
            compare: (other: string | number | Long.Long) => number;
            comp: (other: string | number | Long.Long) => number;
            divide: (divisor: string | number | Long.Long) => Long.Long;
            div: (divisor: string | number | Long.Long) => Long.Long;
            equals: (other: string | number | Long.Long) => boolean;
            eq: (other: string | number | Long.Long) => boolean;
            getHighBits: () => number;
            getHighBitsUnsigned: () => number;
            getLowBits: () => number;
            getLowBitsUnsigned: () => number;
            getNumBitsAbs: () => number;
            greaterThan: (other: string | number | Long.Long) => boolean;
            gt: (other: string | number | Long.Long) => boolean;
            greaterThanOrEqual: (other: string | number | Long.Long) => boolean;
            gte: (other: string | number | Long.Long) => boolean;
            isEven: () => boolean;
            isNegative: () => boolean;
            isOdd: () => boolean;
            isPositive: () => boolean;
            isZero: () => boolean;
            lessThan: (other: string | number | Long.Long) => boolean;
            lt: (other: string | number | Long.Long) => boolean;
            lessThanOrEqual: (other: string | number | Long.Long) => boolean;
            lte: (other: string | number | Long.Long) => boolean;
            modulo: (other: string | number | Long.Long) => Long.Long;
            mod: (other: string | number | Long.Long) => Long.Long;
            multiply: (multiplier: string | number | Long.Long) => Long.Long;
            mul: (multiplier: string | number | Long.Long) => Long.Long;
            negate: () => Long.Long;
            neg: () => Long.Long;
            not: () => Long.Long;
            notEquals: (other: string | number | Long.Long) => boolean;
            neq: (other: string | number | Long.Long) => boolean;
            or: (other: string | number | Long.Long) => Long.Long;
            shiftLeft: (numBits: number | Long.Long) => Long.Long;
            shl: (numBits: number | Long.Long) => Long.Long;
            shiftRight: (numBits: number | Long.Long) => Long.Long;
            shr: (numBits: number | Long.Long) => Long.Long;
            shiftRightUnsigned: (numBits: number | Long.Long) => Long.Long;
            shru: (numBits: number | Long.Long) => Long.Long;
            subtract: (subtrahend: string | number | Long.Long) => Long.Long;
            sub: (subtrahend: string | number | Long.Long) => Long.Long;
            toInt: () => number;
            toNumber: () => number;
            toBytes: (le?: boolean | undefined) => number[];
            toBytesLE: () => number[];
            toBytesBE: () => number[];
            toSigned: () => Long.Long;
            toString: (radix?: number | undefined) => string;
            toUnsigned: () => Long.Long;
            xor: (other: string | number | Long.Long) => Long.Long;
        } & { [K_2 in Exclude<keyof I_1["nonce"], keyof Long.Long>]: never; }) | undefined;
    } & { [K_3 in Exclude<keyof I_1, "nonce">]: never; }>(object: I_1): MsgSendMessageResponse;
};
export declare const MsgSendMessageWithCaller: {
    encode(message: MsgSendMessageWithCaller, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgSendMessageWithCaller;
    fromJSON(object: any): MsgSendMessageWithCaller;
    toJSON(message: MsgSendMessageWithCaller): unknown;
    create<I extends {
        from?: string | undefined;
        destinationDomain?: number | undefined;
        recipient?: Uint8Array | undefined;
        messageBody?: Uint8Array | undefined;
        destinationCaller?: Uint8Array | undefined;
    } & {
        from?: string | undefined;
        destinationDomain?: number | undefined;
        recipient?: Uint8Array | undefined;
        messageBody?: Uint8Array | undefined;
        destinationCaller?: Uint8Array | undefined;
    } & { [K in Exclude<keyof I, keyof MsgSendMessageWithCaller>]: never; }>(base?: I | undefined): MsgSendMessageWithCaller;
    fromPartial<I_1 extends {
        from?: string | undefined;
        destinationDomain?: number | undefined;
        recipient?: Uint8Array | undefined;
        messageBody?: Uint8Array | undefined;
        destinationCaller?: Uint8Array | undefined;
    } & {
        from?: string | undefined;
        destinationDomain?: number | undefined;
        recipient?: Uint8Array | undefined;
        messageBody?: Uint8Array | undefined;
        destinationCaller?: Uint8Array | undefined;
    } & { [K_1 in Exclude<keyof I_1, keyof MsgSendMessageWithCaller>]: never; }>(object: I_1): MsgSendMessageWithCaller;
};
export declare const MsgSendMessageWithCallerResponse: {
    encode(message: MsgSendMessageWithCallerResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgSendMessageWithCallerResponse;
    fromJSON(object: any): MsgSendMessageWithCallerResponse;
    toJSON(message: MsgSendMessageWithCallerResponse): unknown;
    create<I extends {
        nonce?: string | number | Long.Long | undefined;
    } & {
        nonce?: string | number | (Long.Long & {
            high: number;
            low: number;
            unsigned: boolean;
            add: (addend: string | number | Long.Long) => Long.Long;
            and: (other: string | number | Long.Long) => Long.Long;
            compare: (other: string | number | Long.Long) => number;
            comp: (other: string | number | Long.Long) => number;
            divide: (divisor: string | number | Long.Long) => Long.Long;
            div: (divisor: string | number | Long.Long) => Long.Long;
            equals: (other: string | number | Long.Long) => boolean;
            eq: (other: string | number | Long.Long) => boolean;
            getHighBits: () => number;
            getHighBitsUnsigned: () => number;
            getLowBits: () => number;
            getLowBitsUnsigned: () => number;
            getNumBitsAbs: () => number;
            greaterThan: (other: string | number | Long.Long) => boolean;
            gt: (other: string | number | Long.Long) => boolean;
            greaterThanOrEqual: (other: string | number | Long.Long) => boolean;
            gte: (other: string | number | Long.Long) => boolean;
            isEven: () => boolean;
            isNegative: () => boolean;
            isOdd: () => boolean;
            isPositive: () => boolean;
            isZero: () => boolean;
            lessThan: (other: string | number | Long.Long) => boolean;
            lt: (other: string | number | Long.Long) => boolean;
            lessThanOrEqual: (other: string | number | Long.Long) => boolean;
            lte: (other: string | number | Long.Long) => boolean;
            modulo: (other: string | number | Long.Long) => Long.Long;
            mod: (other: string | number | Long.Long) => Long.Long;
            multiply: (multiplier: string | number | Long.Long) => Long.Long;
            mul: (multiplier: string | number | Long.Long) => Long.Long;
            negate: () => Long.Long;
            neg: () => Long.Long;
            not: () => Long.Long;
            notEquals: (other: string | number | Long.Long) => boolean;
            neq: (other: string | number | Long.Long) => boolean;
            or: (other: string | number | Long.Long) => Long.Long;
            shiftLeft: (numBits: number | Long.Long) => Long.Long;
            shl: (numBits: number | Long.Long) => Long.Long;
            shiftRight: (numBits: number | Long.Long) => Long.Long;
            shr: (numBits: number | Long.Long) => Long.Long;
            shiftRightUnsigned: (numBits: number | Long.Long) => Long.Long;
            shru: (numBits: number | Long.Long) => Long.Long;
            subtract: (subtrahend: string | number | Long.Long) => Long.Long;
            sub: (subtrahend: string | number | Long.Long) => Long.Long;
            toInt: () => number;
            toNumber: () => number;
            toBytes: (le?: boolean | undefined) => number[];
            toBytesLE: () => number[];
            toBytesBE: () => number[];
            toSigned: () => Long.Long;
            toString: (radix?: number | undefined) => string;
            toUnsigned: () => Long.Long;
            xor: (other: string | number | Long.Long) => Long.Long;
        } & { [K in Exclude<keyof I["nonce"], keyof Long.Long>]: never; }) | undefined;
    } & { [K_1 in Exclude<keyof I, "nonce">]: never; }>(base?: I | undefined): MsgSendMessageWithCallerResponse;
    fromPartial<I_1 extends {
        nonce?: string | number | Long.Long | undefined;
    } & {
        nonce?: string | number | (Long.Long & {
            high: number;
            low: number;
            unsigned: boolean;
            add: (addend: string | number | Long.Long) => Long.Long;
            and: (other: string | number | Long.Long) => Long.Long;
            compare: (other: string | number | Long.Long) => number;
            comp: (other: string | number | Long.Long) => number;
            divide: (divisor: string | number | Long.Long) => Long.Long;
            div: (divisor: string | number | Long.Long) => Long.Long;
            equals: (other: string | number | Long.Long) => boolean;
            eq: (other: string | number | Long.Long) => boolean;
            getHighBits: () => number;
            getHighBitsUnsigned: () => number;
            getLowBits: () => number;
            getLowBitsUnsigned: () => number;
            getNumBitsAbs: () => number;
            greaterThan: (other: string | number | Long.Long) => boolean;
            gt: (other: string | number | Long.Long) => boolean;
            greaterThanOrEqual: (other: string | number | Long.Long) => boolean;
            gte: (other: string | number | Long.Long) => boolean;
            isEven: () => boolean;
            isNegative: () => boolean;
            isOdd: () => boolean;
            isPositive: () => boolean;
            isZero: () => boolean;
            lessThan: (other: string | number | Long.Long) => boolean;
            lt: (other: string | number | Long.Long) => boolean;
            lessThanOrEqual: (other: string | number | Long.Long) => boolean;
            lte: (other: string | number | Long.Long) => boolean;
            modulo: (other: string | number | Long.Long) => Long.Long;
            mod: (other: string | number | Long.Long) => Long.Long;
            multiply: (multiplier: string | number | Long.Long) => Long.Long;
            mul: (multiplier: string | number | Long.Long) => Long.Long;
            negate: () => Long.Long;
            neg: () => Long.Long;
            not: () => Long.Long;
            notEquals: (other: string | number | Long.Long) => boolean;
            neq: (other: string | number | Long.Long) => boolean;
            or: (other: string | number | Long.Long) => Long.Long;
            shiftLeft: (numBits: number | Long.Long) => Long.Long;
            shl: (numBits: number | Long.Long) => Long.Long;
            shiftRight: (numBits: number | Long.Long) => Long.Long;
            shr: (numBits: number | Long.Long) => Long.Long;
            shiftRightUnsigned: (numBits: number | Long.Long) => Long.Long;
            shru: (numBits: number | Long.Long) => Long.Long;
            subtract: (subtrahend: string | number | Long.Long) => Long.Long;
            sub: (subtrahend: string | number | Long.Long) => Long.Long;
            toInt: () => number;
            toNumber: () => number;
            toBytes: (le?: boolean | undefined) => number[];
            toBytesLE: () => number[];
            toBytesBE: () => number[];
            toSigned: () => Long.Long;
            toString: (radix?: number | undefined) => string;
            toUnsigned: () => Long.Long;
            xor: (other: string | number | Long.Long) => Long.Long;
        } & { [K_2 in Exclude<keyof I_1["nonce"], keyof Long.Long>]: never; }) | undefined;
    } & { [K_3 in Exclude<keyof I_1, "nonce">]: never; }>(object: I_1): MsgSendMessageWithCallerResponse;
};
export declare const MsgReplaceMessage: {
    encode(message: MsgReplaceMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgReplaceMessage;
    fromJSON(object: any): MsgReplaceMessage;
    toJSON(message: MsgReplaceMessage): unknown;
    create<I extends {
        from?: string | undefined;
        originalMessage?: Uint8Array | undefined;
        originalAttestation?: Uint8Array | undefined;
        newMessageBody?: Uint8Array | undefined;
        newDestinationCaller?: Uint8Array | undefined;
    } & {
        from?: string | undefined;
        originalMessage?: Uint8Array | undefined;
        originalAttestation?: Uint8Array | undefined;
        newMessageBody?: Uint8Array | undefined;
        newDestinationCaller?: Uint8Array | undefined;
    } & { [K in Exclude<keyof I, keyof MsgReplaceMessage>]: never; }>(base?: I | undefined): MsgReplaceMessage;
    fromPartial<I_1 extends {
        from?: string | undefined;
        originalMessage?: Uint8Array | undefined;
        originalAttestation?: Uint8Array | undefined;
        newMessageBody?: Uint8Array | undefined;
        newDestinationCaller?: Uint8Array | undefined;
    } & {
        from?: string | undefined;
        originalMessage?: Uint8Array | undefined;
        originalAttestation?: Uint8Array | undefined;
        newMessageBody?: Uint8Array | undefined;
        newDestinationCaller?: Uint8Array | undefined;
    } & { [K_1 in Exclude<keyof I_1, keyof MsgReplaceMessage>]: never; }>(object: I_1): MsgReplaceMessage;
};
export declare const MsgReplaceMessageResponse: {
    encode(_: MsgReplaceMessageResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgReplaceMessageResponse;
    fromJSON(_: any): MsgReplaceMessageResponse;
    toJSON(_: MsgReplaceMessageResponse): unknown;
    create<I extends {} & {} & { [K in Exclude<keyof I, never>]: never; }>(base?: I | undefined): MsgReplaceMessageResponse;
    fromPartial<I_1 extends {} & {} & { [K_1 in Exclude<keyof I_1, never>]: never; }>(_: I_1): MsgReplaceMessageResponse;
};
export declare const MsgUpdateSignatureThreshold: {
    encode(message: MsgUpdateSignatureThreshold, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateSignatureThreshold;
    fromJSON(object: any): MsgUpdateSignatureThreshold;
    toJSON(message: MsgUpdateSignatureThreshold): unknown;
    create<I extends {
        from?: string | undefined;
        amount?: number | undefined;
    } & {
        from?: string | undefined;
        amount?: number | undefined;
    } & { [K in Exclude<keyof I, keyof MsgUpdateSignatureThreshold>]: never; }>(base?: I | undefined): MsgUpdateSignatureThreshold;
    fromPartial<I_1 extends {
        from?: string | undefined;
        amount?: number | undefined;
    } & {
        from?: string | undefined;
        amount?: number | undefined;
    } & { [K_1 in Exclude<keyof I_1, keyof MsgUpdateSignatureThreshold>]: never; }>(object: I_1): MsgUpdateSignatureThreshold;
};
export declare const MsgUpdateSignatureThresholdResponse: {
    encode(_: MsgUpdateSignatureThresholdResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateSignatureThresholdResponse;
    fromJSON(_: any): MsgUpdateSignatureThresholdResponse;
    toJSON(_: MsgUpdateSignatureThresholdResponse): unknown;
    create<I extends {} & {} & { [K in Exclude<keyof I, never>]: never; }>(base?: I | undefined): MsgUpdateSignatureThresholdResponse;
    fromPartial<I_1 extends {} & {} & { [K_1 in Exclude<keyof I_1, never>]: never; }>(_: I_1): MsgUpdateSignatureThresholdResponse;
};
export declare const MsgLinkTokenPair: {
    encode(message: MsgLinkTokenPair, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgLinkTokenPair;
    fromJSON(object: any): MsgLinkTokenPair;
    toJSON(message: MsgLinkTokenPair): unknown;
    create<I extends {
        from?: string | undefined;
        remoteDomain?: number | undefined;
        remoteToken?: Uint8Array | undefined;
        localToken?: string | undefined;
    } & {
        from?: string | undefined;
        remoteDomain?: number | undefined;
        remoteToken?: Uint8Array | undefined;
        localToken?: string | undefined;
    } & { [K in Exclude<keyof I, keyof MsgLinkTokenPair>]: never; }>(base?: I | undefined): MsgLinkTokenPair;
    fromPartial<I_1 extends {
        from?: string | undefined;
        remoteDomain?: number | undefined;
        remoteToken?: Uint8Array | undefined;
        localToken?: string | undefined;
    } & {
        from?: string | undefined;
        remoteDomain?: number | undefined;
        remoteToken?: Uint8Array | undefined;
        localToken?: string | undefined;
    } & { [K_1 in Exclude<keyof I_1, keyof MsgLinkTokenPair>]: never; }>(object: I_1): MsgLinkTokenPair;
};
export declare const MsgLinkTokenPairResponse: {
    encode(_: MsgLinkTokenPairResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgLinkTokenPairResponse;
    fromJSON(_: any): MsgLinkTokenPairResponse;
    toJSON(_: MsgLinkTokenPairResponse): unknown;
    create<I extends {} & {} & { [K in Exclude<keyof I, never>]: never; }>(base?: I | undefined): MsgLinkTokenPairResponse;
    fromPartial<I_1 extends {} & {} & { [K_1 in Exclude<keyof I_1, never>]: never; }>(_: I_1): MsgLinkTokenPairResponse;
};
export declare const MsgUnlinkTokenPair: {
    encode(message: MsgUnlinkTokenPair, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUnlinkTokenPair;
    fromJSON(object: any): MsgUnlinkTokenPair;
    toJSON(message: MsgUnlinkTokenPair): unknown;
    create<I extends {
        from?: string | undefined;
        remoteDomain?: number | undefined;
        remoteToken?: Uint8Array | undefined;
        localToken?: string | undefined;
    } & {
        from?: string | undefined;
        remoteDomain?: number | undefined;
        remoteToken?: Uint8Array | undefined;
        localToken?: string | undefined;
    } & { [K in Exclude<keyof I, keyof MsgUnlinkTokenPair>]: never; }>(base?: I | undefined): MsgUnlinkTokenPair;
    fromPartial<I_1 extends {
        from?: string | undefined;
        remoteDomain?: number | undefined;
        remoteToken?: Uint8Array | undefined;
        localToken?: string | undefined;
    } & {
        from?: string | undefined;
        remoteDomain?: number | undefined;
        remoteToken?: Uint8Array | undefined;
        localToken?: string | undefined;
    } & { [K_1 in Exclude<keyof I_1, keyof MsgUnlinkTokenPair>]: never; }>(object: I_1): MsgUnlinkTokenPair;
};
export declare const MsgUnlinkTokenPairResponse: {
    encode(_: MsgUnlinkTokenPairResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUnlinkTokenPairResponse;
    fromJSON(_: any): MsgUnlinkTokenPairResponse;
    toJSON(_: MsgUnlinkTokenPairResponse): unknown;
    create<I extends {} & {} & { [K in Exclude<keyof I, never>]: never; }>(base?: I | undefined): MsgUnlinkTokenPairResponse;
    fromPartial<I_1 extends {} & {} & { [K_1 in Exclude<keyof I_1, never>]: never; }>(_: I_1): MsgUnlinkTokenPairResponse;
};
export declare const MsgAddRemoteTokenMessenger: {
    encode(message: MsgAddRemoteTokenMessenger, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgAddRemoteTokenMessenger;
    fromJSON(object: any): MsgAddRemoteTokenMessenger;
    toJSON(message: MsgAddRemoteTokenMessenger): unknown;
    create<I extends {
        from?: string | undefined;
        domainId?: number | undefined;
        address?: Uint8Array | undefined;
    } & {
        from?: string | undefined;
        domainId?: number | undefined;
        address?: Uint8Array | undefined;
    } & { [K in Exclude<keyof I, keyof MsgAddRemoteTokenMessenger>]: never; }>(base?: I | undefined): MsgAddRemoteTokenMessenger;
    fromPartial<I_1 extends {
        from?: string | undefined;
        domainId?: number | undefined;
        address?: Uint8Array | undefined;
    } & {
        from?: string | undefined;
        domainId?: number | undefined;
        address?: Uint8Array | undefined;
    } & { [K_1 in Exclude<keyof I_1, keyof MsgAddRemoteTokenMessenger>]: never; }>(object: I_1): MsgAddRemoteTokenMessenger;
};
export declare const MsgAddRemoteTokenMessengerResponse: {
    encode(_: MsgAddRemoteTokenMessengerResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgAddRemoteTokenMessengerResponse;
    fromJSON(_: any): MsgAddRemoteTokenMessengerResponse;
    toJSON(_: MsgAddRemoteTokenMessengerResponse): unknown;
    create<I extends {} & {} & { [K in Exclude<keyof I, never>]: never; }>(base?: I | undefined): MsgAddRemoteTokenMessengerResponse;
    fromPartial<I_1 extends {} & {} & { [K_1 in Exclude<keyof I_1, never>]: never; }>(_: I_1): MsgAddRemoteTokenMessengerResponse;
};
export declare const MsgRemoveRemoteTokenMessenger: {
    encode(message: MsgRemoveRemoteTokenMessenger, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgRemoveRemoteTokenMessenger;
    fromJSON(object: any): MsgRemoveRemoteTokenMessenger;
    toJSON(message: MsgRemoveRemoteTokenMessenger): unknown;
    create<I extends {
        from?: string | undefined;
        domainId?: number | undefined;
    } & {
        from?: string | undefined;
        domainId?: number | undefined;
    } & { [K in Exclude<keyof I, keyof MsgRemoveRemoteTokenMessenger>]: never; }>(base?: I | undefined): MsgRemoveRemoteTokenMessenger;
    fromPartial<I_1 extends {
        from?: string | undefined;
        domainId?: number | undefined;
    } & {
        from?: string | undefined;
        domainId?: number | undefined;
    } & { [K_1 in Exclude<keyof I_1, keyof MsgRemoveRemoteTokenMessenger>]: never; }>(object: I_1): MsgRemoveRemoteTokenMessenger;
};
export declare const MsgRemoveRemoteTokenMessengerResponse: {
    encode(_: MsgRemoveRemoteTokenMessengerResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgRemoveRemoteTokenMessengerResponse;
    fromJSON(_: any): MsgRemoveRemoteTokenMessengerResponse;
    toJSON(_: MsgRemoveRemoteTokenMessengerResponse): unknown;
    create<I extends {} & {} & { [K in Exclude<keyof I, never>]: never; }>(base?: I | undefined): MsgRemoveRemoteTokenMessengerResponse;
    fromPartial<I_1 extends {} & {} & { [K_1 in Exclude<keyof I_1, never>]: never; }>(_: I_1): MsgRemoveRemoteTokenMessengerResponse;
};
/** Msg defines the Msg service. */
export interface Msg {
    AcceptOwner(request: MsgAcceptOwner): Promise<MsgAcceptOwnerResponse>;
    AddRemoteTokenMessenger(request: MsgAddRemoteTokenMessenger): Promise<MsgAddRemoteTokenMessengerResponse>;
    DepositForBurn(request: MsgDepositForBurn): Promise<MsgDepositForBurnResponse>;
    DepositForBurnWithCaller(request: MsgDepositForBurnWithCaller): Promise<MsgDepositForBurnWithCallerResponse>;
    DisableAttester(request: MsgDisableAttester): Promise<MsgDisableAttesterResponse>;
    EnableAttester(request: MsgEnableAttester): Promise<MsgEnableAttesterResponse>;
    LinkTokenPair(request: MsgLinkTokenPair): Promise<MsgLinkTokenPairResponse>;
    PauseBurningAndMinting(request: MsgPauseBurningAndMinting): Promise<MsgPauseBurningAndMintingResponse>;
    PauseSendingAndReceivingMessages(request: MsgPauseSendingAndReceivingMessages): Promise<MsgPauseSendingAndReceivingMessagesResponse>;
    ReceiveMessage(request: MsgReceiveMessage): Promise<MsgReceiveMessageResponse>;
    RemoveRemoteTokenMessenger(request: MsgRemoveRemoteTokenMessenger): Promise<MsgRemoveRemoteTokenMessengerResponse>;
    ReplaceDepositForBurn(request: MsgReplaceDepositForBurn): Promise<MsgReplaceDepositForBurnResponse>;
    ReplaceMessage(request: MsgReplaceMessage): Promise<MsgReplaceMessageResponse>;
    SendMessage(request: MsgSendMessage): Promise<MsgSendMessageResponse>;
    SendMessageWithCaller(request: MsgSendMessageWithCaller): Promise<MsgSendMessageWithCallerResponse>;
    UnlinkTokenPair(request: MsgUnlinkTokenPair): Promise<MsgUnlinkTokenPairResponse>;
    UnpauseBurningAndMinting(request: MsgUnpauseBurningAndMinting): Promise<MsgUnpauseBurningAndMintingResponse>;
    UnpauseSendingAndReceivingMessages(request: MsgUnpauseSendingAndReceivingMessages): Promise<MsgUnpauseSendingAndReceivingMessagesResponse>;
    UpdateOwner(request: MsgUpdateOwner): Promise<MsgUpdateOwnerResponse>;
    UpdateAttesterManager(request: MsgUpdateAttesterManager): Promise<MsgUpdateAttesterManagerResponse>;
    UpdateTokenController(request: MsgUpdateTokenController): Promise<MsgUpdateTokenControllerResponse>;
    UpdatePauser(request: MsgUpdatePauser): Promise<MsgUpdatePauserResponse>;
    UpdateMaxMessageBodySize(request: MsgUpdateMaxMessageBodySize): Promise<MsgUpdateMaxMessageBodySizeResponse>;
    SetMaxBurnAmountPerMessage(request: MsgSetMaxBurnAmountPerMessage): Promise<MsgSetMaxBurnAmountPerMessageResponse>;
    UpdateSignatureThreshold(request: MsgUpdateSignatureThreshold): Promise<MsgUpdateSignatureThresholdResponse>;
}
export declare const MsgServiceName = "circle.cctp.v1.Msg";
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    private readonly service;
    constructor(rpc: Rpc, opts?: {
        service?: string;
    });
    AcceptOwner(request: MsgAcceptOwner): Promise<MsgAcceptOwnerResponse>;
    AddRemoteTokenMessenger(request: MsgAddRemoteTokenMessenger): Promise<MsgAddRemoteTokenMessengerResponse>;
    DepositForBurn(request: MsgDepositForBurn): Promise<MsgDepositForBurnResponse>;
    DepositForBurnWithCaller(request: MsgDepositForBurnWithCaller): Promise<MsgDepositForBurnWithCallerResponse>;
    DisableAttester(request: MsgDisableAttester): Promise<MsgDisableAttesterResponse>;
    EnableAttester(request: MsgEnableAttester): Promise<MsgEnableAttesterResponse>;
    LinkTokenPair(request: MsgLinkTokenPair): Promise<MsgLinkTokenPairResponse>;
    PauseBurningAndMinting(request: MsgPauseBurningAndMinting): Promise<MsgPauseBurningAndMintingResponse>;
    PauseSendingAndReceivingMessages(request: MsgPauseSendingAndReceivingMessages): Promise<MsgPauseSendingAndReceivingMessagesResponse>;
    ReceiveMessage(request: MsgReceiveMessage): Promise<MsgReceiveMessageResponse>;
    RemoveRemoteTokenMessenger(request: MsgRemoveRemoteTokenMessenger): Promise<MsgRemoveRemoteTokenMessengerResponse>;
    ReplaceDepositForBurn(request: MsgReplaceDepositForBurn): Promise<MsgReplaceDepositForBurnResponse>;
    ReplaceMessage(request: MsgReplaceMessage): Promise<MsgReplaceMessageResponse>;
    SendMessage(request: MsgSendMessage): Promise<MsgSendMessageResponse>;
    SendMessageWithCaller(request: MsgSendMessageWithCaller): Promise<MsgSendMessageWithCallerResponse>;
    UnlinkTokenPair(request: MsgUnlinkTokenPair): Promise<MsgUnlinkTokenPairResponse>;
    UnpauseBurningAndMinting(request: MsgUnpauseBurningAndMinting): Promise<MsgUnpauseBurningAndMintingResponse>;
    UnpauseSendingAndReceivingMessages(request: MsgUnpauseSendingAndReceivingMessages): Promise<MsgUnpauseSendingAndReceivingMessagesResponse>;
    UpdateOwner(request: MsgUpdateOwner): Promise<MsgUpdateOwnerResponse>;
    UpdateAttesterManager(request: MsgUpdateAttesterManager): Promise<MsgUpdateAttesterManagerResponse>;
    UpdateTokenController(request: MsgUpdateTokenController): Promise<MsgUpdateTokenControllerResponse>;
    UpdatePauser(request: MsgUpdatePauser): Promise<MsgUpdatePauserResponse>;
    UpdateMaxMessageBodySize(request: MsgUpdateMaxMessageBodySize): Promise<MsgUpdateMaxMessageBodySizeResponse>;
    SetMaxBurnAmountPerMessage(request: MsgSetMaxBurnAmountPerMessage): Promise<MsgSetMaxBurnAmountPerMessageResponse>;
    UpdateSignatureThreshold(request: MsgUpdateSignatureThreshold): Promise<MsgUpdateSignatureThresholdResponse>;
}
interface Rpc {
    request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}
type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;
export type DeepPartial<T> = T extends Builtin ? T : T extends Long ? string | number | Long : T extends globalThis.Array<infer U> ? globalThis.Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P : P & {
    [K in keyof P]: Exact<P[K], I[K]>;
} & {
    [K in Exclude<keyof I, KeysOfUnion<P>>]: never;
};
export {};
