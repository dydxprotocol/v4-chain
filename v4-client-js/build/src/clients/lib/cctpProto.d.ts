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
    create<I extends Exact<DeepPartial<MsgUpdateOwner>, I>>(base?: I): MsgUpdateOwner;
    fromPartial<I extends Exact<DeepPartial<MsgUpdateOwner>, I>>(object: I): MsgUpdateOwner;
};
export declare const MsgUpdateOwnerResponse: {
    encode(_: MsgUpdateOwnerResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateOwnerResponse;
    fromJSON(_: any): MsgUpdateOwnerResponse;
    toJSON(_: MsgUpdateOwnerResponse): unknown;
    create<I extends Exact<DeepPartial<MsgUpdateOwnerResponse>, I>>(base?: I): MsgUpdateOwnerResponse;
    fromPartial<I extends Exact<DeepPartial<MsgUpdateOwnerResponse>, I>>(_: I): MsgUpdateOwnerResponse;
};
export declare const MsgUpdateAttesterManager: {
    encode(message: MsgUpdateAttesterManager, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateAttesterManager;
    fromJSON(object: any): MsgUpdateAttesterManager;
    toJSON(message: MsgUpdateAttesterManager): unknown;
    create<I extends Exact<DeepPartial<MsgUpdateAttesterManager>, I>>(base?: I): MsgUpdateAttesterManager;
    fromPartial<I extends Exact<DeepPartial<MsgUpdateAttesterManager>, I>>(object: I): MsgUpdateAttesterManager;
};
export declare const MsgUpdateAttesterManagerResponse: {
    encode(_: MsgUpdateAttesterManagerResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateAttesterManagerResponse;
    fromJSON(_: any): MsgUpdateAttesterManagerResponse;
    toJSON(_: MsgUpdateAttesterManagerResponse): unknown;
    create<I extends Exact<DeepPartial<MsgUpdateAttesterManagerResponse>, I>>(base?: I): MsgUpdateAttesterManagerResponse;
    fromPartial<I extends Exact<DeepPartial<MsgUpdateAttesterManagerResponse>, I>>(_: I): MsgUpdateAttesterManagerResponse;
};
export declare const MsgUpdateTokenController: {
    encode(message: MsgUpdateTokenController, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateTokenController;
    fromJSON(object: any): MsgUpdateTokenController;
    toJSON(message: MsgUpdateTokenController): unknown;
    create<I extends Exact<DeepPartial<MsgUpdateTokenController>, I>>(base?: I): MsgUpdateTokenController;
    fromPartial<I extends Exact<DeepPartial<MsgUpdateTokenController>, I>>(object: I): MsgUpdateTokenController;
};
export declare const MsgUpdateTokenControllerResponse: {
    encode(_: MsgUpdateTokenControllerResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateTokenControllerResponse;
    fromJSON(_: any): MsgUpdateTokenControllerResponse;
    toJSON(_: MsgUpdateTokenControllerResponse): unknown;
    create<I extends Exact<DeepPartial<MsgUpdateTokenControllerResponse>, I>>(base?: I): MsgUpdateTokenControllerResponse;
    fromPartial<I extends Exact<DeepPartial<MsgUpdateTokenControllerResponse>, I>>(_: I): MsgUpdateTokenControllerResponse;
};
export declare const MsgUpdatePauser: {
    encode(message: MsgUpdatePauser, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdatePauser;
    fromJSON(object: any): MsgUpdatePauser;
    toJSON(message: MsgUpdatePauser): unknown;
    create<I extends Exact<DeepPartial<MsgUpdatePauser>, I>>(base?: I): MsgUpdatePauser;
    fromPartial<I extends Exact<DeepPartial<MsgUpdatePauser>, I>>(object: I): MsgUpdatePauser;
};
export declare const MsgUpdatePauserResponse: {
    encode(_: MsgUpdatePauserResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdatePauserResponse;
    fromJSON(_: any): MsgUpdatePauserResponse;
    toJSON(_: MsgUpdatePauserResponse): unknown;
    create<I extends Exact<DeepPartial<MsgUpdatePauserResponse>, I>>(base?: I): MsgUpdatePauserResponse;
    fromPartial<I extends Exact<DeepPartial<MsgUpdatePauserResponse>, I>>(_: I): MsgUpdatePauserResponse;
};
export declare const MsgAcceptOwner: {
    encode(message: MsgAcceptOwner, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgAcceptOwner;
    fromJSON(object: any): MsgAcceptOwner;
    toJSON(message: MsgAcceptOwner): unknown;
    create<I extends Exact<DeepPartial<MsgAcceptOwner>, I>>(base?: I): MsgAcceptOwner;
    fromPartial<I extends Exact<DeepPartial<MsgAcceptOwner>, I>>(object: I): MsgAcceptOwner;
};
export declare const MsgAcceptOwnerResponse: {
    encode(_: MsgAcceptOwnerResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgAcceptOwnerResponse;
    fromJSON(_: any): MsgAcceptOwnerResponse;
    toJSON(_: MsgAcceptOwnerResponse): unknown;
    create<I extends Exact<DeepPartial<MsgAcceptOwnerResponse>, I>>(base?: I): MsgAcceptOwnerResponse;
    fromPartial<I extends Exact<DeepPartial<MsgAcceptOwnerResponse>, I>>(_: I): MsgAcceptOwnerResponse;
};
export declare const MsgEnableAttester: {
    encode(message: MsgEnableAttester, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgEnableAttester;
    fromJSON(object: any): MsgEnableAttester;
    toJSON(message: MsgEnableAttester): unknown;
    create<I extends Exact<DeepPartial<MsgEnableAttester>, I>>(base?: I): MsgEnableAttester;
    fromPartial<I extends Exact<DeepPartial<MsgEnableAttester>, I>>(object: I): MsgEnableAttester;
};
export declare const MsgEnableAttesterResponse: {
    encode(_: MsgEnableAttesterResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgEnableAttesterResponse;
    fromJSON(_: any): MsgEnableAttesterResponse;
    toJSON(_: MsgEnableAttesterResponse): unknown;
    create<I extends Exact<DeepPartial<MsgEnableAttesterResponse>, I>>(base?: I): MsgEnableAttesterResponse;
    fromPartial<I extends Exact<DeepPartial<MsgEnableAttesterResponse>, I>>(_: I): MsgEnableAttesterResponse;
};
export declare const MsgDisableAttester: {
    encode(message: MsgDisableAttester, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgDisableAttester;
    fromJSON(object: any): MsgDisableAttester;
    toJSON(message: MsgDisableAttester): unknown;
    create<I extends Exact<DeepPartial<MsgDisableAttester>, I>>(base?: I): MsgDisableAttester;
    fromPartial<I extends Exact<DeepPartial<MsgDisableAttester>, I>>(object: I): MsgDisableAttester;
};
export declare const MsgDisableAttesterResponse: {
    encode(_: MsgDisableAttesterResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgDisableAttesterResponse;
    fromJSON(_: any): MsgDisableAttesterResponse;
    toJSON(_: MsgDisableAttesterResponse): unknown;
    create<I extends Exact<DeepPartial<MsgDisableAttesterResponse>, I>>(base?: I): MsgDisableAttesterResponse;
    fromPartial<I extends Exact<DeepPartial<MsgDisableAttesterResponse>, I>>(_: I): MsgDisableAttesterResponse;
};
export declare const MsgPauseBurningAndMinting: {
    encode(message: MsgPauseBurningAndMinting, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgPauseBurningAndMinting;
    fromJSON(object: any): MsgPauseBurningAndMinting;
    toJSON(message: MsgPauseBurningAndMinting): unknown;
    create<I extends Exact<DeepPartial<MsgPauseBurningAndMinting>, I>>(base?: I): MsgPauseBurningAndMinting;
    fromPartial<I extends Exact<DeepPartial<MsgPauseBurningAndMinting>, I>>(object: I): MsgPauseBurningAndMinting;
};
export declare const MsgPauseBurningAndMintingResponse: {
    encode(_: MsgPauseBurningAndMintingResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgPauseBurningAndMintingResponse;
    fromJSON(_: any): MsgPauseBurningAndMintingResponse;
    toJSON(_: MsgPauseBurningAndMintingResponse): unknown;
    create<I extends Exact<DeepPartial<MsgPauseBurningAndMintingResponse>, I>>(base?: I): MsgPauseBurningAndMintingResponse;
    fromPartial<I extends Exact<DeepPartial<MsgPauseBurningAndMintingResponse>, I>>(_: I): MsgPauseBurningAndMintingResponse;
};
export declare const MsgUnpauseBurningAndMinting: {
    encode(message: MsgUnpauseBurningAndMinting, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUnpauseBurningAndMinting;
    fromJSON(object: any): MsgUnpauseBurningAndMinting;
    toJSON(message: MsgUnpauseBurningAndMinting): unknown;
    create<I extends Exact<DeepPartial<MsgUnpauseBurningAndMinting>, I>>(base?: I): MsgUnpauseBurningAndMinting;
    fromPartial<I extends Exact<DeepPartial<MsgUnpauseBurningAndMinting>, I>>(object: I): MsgUnpauseBurningAndMinting;
};
export declare const MsgUnpauseBurningAndMintingResponse: {
    encode(_: MsgUnpauseBurningAndMintingResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUnpauseBurningAndMintingResponse;
    fromJSON(_: any): MsgUnpauseBurningAndMintingResponse;
    toJSON(_: MsgUnpauseBurningAndMintingResponse): unknown;
    create<I extends Exact<DeepPartial<MsgUnpauseBurningAndMintingResponse>, I>>(base?: I): MsgUnpauseBurningAndMintingResponse;
    fromPartial<I extends Exact<DeepPartial<MsgUnpauseBurningAndMintingResponse>, I>>(_: I): MsgUnpauseBurningAndMintingResponse;
};
export declare const MsgPauseSendingAndReceivingMessages: {
    encode(message: MsgPauseSendingAndReceivingMessages, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgPauseSendingAndReceivingMessages;
    fromJSON(object: any): MsgPauseSendingAndReceivingMessages;
    toJSON(message: MsgPauseSendingAndReceivingMessages): unknown;
    create<I extends Exact<DeepPartial<MsgPauseSendingAndReceivingMessages>, I>>(base?: I): MsgPauseSendingAndReceivingMessages;
    fromPartial<I extends Exact<DeepPartial<MsgPauseSendingAndReceivingMessages>, I>>(object: I): MsgPauseSendingAndReceivingMessages;
};
export declare const MsgPauseSendingAndReceivingMessagesResponse: {
    encode(_: MsgPauseSendingAndReceivingMessagesResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgPauseSendingAndReceivingMessagesResponse;
    fromJSON(_: any): MsgPauseSendingAndReceivingMessagesResponse;
    toJSON(_: MsgPauseSendingAndReceivingMessagesResponse): unknown;
    create<I extends Exact<DeepPartial<MsgPauseSendingAndReceivingMessagesResponse>, I>>(base?: I): MsgPauseSendingAndReceivingMessagesResponse;
    fromPartial<I extends Exact<DeepPartial<MsgPauseSendingAndReceivingMessagesResponse>, I>>(_: I): MsgPauseSendingAndReceivingMessagesResponse;
};
export declare const MsgUnpauseSendingAndReceivingMessages: {
    encode(message: MsgUnpauseSendingAndReceivingMessages, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUnpauseSendingAndReceivingMessages;
    fromJSON(object: any): MsgUnpauseSendingAndReceivingMessages;
    toJSON(message: MsgUnpauseSendingAndReceivingMessages): unknown;
    create<I extends Exact<DeepPartial<MsgUnpauseSendingAndReceivingMessages>, I>>(base?: I): MsgUnpauseSendingAndReceivingMessages;
    fromPartial<I extends Exact<DeepPartial<MsgUnpauseSendingAndReceivingMessages>, I>>(object: I): MsgUnpauseSendingAndReceivingMessages;
};
export declare const MsgUnpauseSendingAndReceivingMessagesResponse: {
    encode(_: MsgUnpauseSendingAndReceivingMessagesResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUnpauseSendingAndReceivingMessagesResponse;
    fromJSON(_: any): MsgUnpauseSendingAndReceivingMessagesResponse;
    toJSON(_: MsgUnpauseSendingAndReceivingMessagesResponse): unknown;
    create<I extends Exact<DeepPartial<MsgUnpauseSendingAndReceivingMessagesResponse>, I>>(base?: I): MsgUnpauseSendingAndReceivingMessagesResponse;
    fromPartial<I extends Exact<DeepPartial<MsgUnpauseSendingAndReceivingMessagesResponse>, I>>(_: I): MsgUnpauseSendingAndReceivingMessagesResponse;
};
export declare const MsgUpdateMaxMessageBodySize: {
    encode(message: MsgUpdateMaxMessageBodySize, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateMaxMessageBodySize;
    fromJSON(object: any): MsgUpdateMaxMessageBodySize;
    toJSON(message: MsgUpdateMaxMessageBodySize): unknown;
    create<I extends Exact<DeepPartial<MsgUpdateMaxMessageBodySize>, I>>(base?: I): MsgUpdateMaxMessageBodySize;
    fromPartial<I extends Exact<DeepPartial<MsgUpdateMaxMessageBodySize>, I>>(object: I): MsgUpdateMaxMessageBodySize;
};
export declare const MsgUpdateMaxMessageBodySizeResponse: {
    encode(_: MsgUpdateMaxMessageBodySizeResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateMaxMessageBodySizeResponse;
    fromJSON(_: any): MsgUpdateMaxMessageBodySizeResponse;
    toJSON(_: MsgUpdateMaxMessageBodySizeResponse): unknown;
    create<I extends Exact<DeepPartial<MsgUpdateMaxMessageBodySizeResponse>, I>>(base?: I): MsgUpdateMaxMessageBodySizeResponse;
    fromPartial<I extends Exact<DeepPartial<MsgUpdateMaxMessageBodySizeResponse>, I>>(_: I): MsgUpdateMaxMessageBodySizeResponse;
};
export declare const MsgSetMaxBurnAmountPerMessage: {
    encode(message: MsgSetMaxBurnAmountPerMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetMaxBurnAmountPerMessage;
    fromJSON(object: any): MsgSetMaxBurnAmountPerMessage;
    toJSON(message: MsgSetMaxBurnAmountPerMessage): unknown;
    create<I extends Exact<DeepPartial<MsgSetMaxBurnAmountPerMessage>, I>>(base?: I): MsgSetMaxBurnAmountPerMessage;
    fromPartial<I extends Exact<DeepPartial<MsgSetMaxBurnAmountPerMessage>, I>>(object: I): MsgSetMaxBurnAmountPerMessage;
};
export declare const MsgSetMaxBurnAmountPerMessageResponse: {
    encode(_: MsgSetMaxBurnAmountPerMessageResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetMaxBurnAmountPerMessageResponse;
    fromJSON(_: any): MsgSetMaxBurnAmountPerMessageResponse;
    toJSON(_: MsgSetMaxBurnAmountPerMessageResponse): unknown;
    create<I extends Exact<DeepPartial<MsgSetMaxBurnAmountPerMessageResponse>, I>>(base?: I): MsgSetMaxBurnAmountPerMessageResponse;
    fromPartial<I extends Exact<DeepPartial<MsgSetMaxBurnAmountPerMessageResponse>, I>>(_: I): MsgSetMaxBurnAmountPerMessageResponse;
};
export declare const MsgDepositForBurn: {
    encode(message: MsgDepositForBurn, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgDepositForBurn;
    fromJSON(object: any): MsgDepositForBurn;
    toJSON(message: MsgDepositForBurn): unknown;
    create<I extends Exact<DeepPartial<MsgDepositForBurn>, I>>(base?: I): MsgDepositForBurn;
    fromPartial<I extends Exact<DeepPartial<MsgDepositForBurn>, I>>(object: I): MsgDepositForBurn;
};
export declare const MsgDepositForBurnResponse: {
    encode(message: MsgDepositForBurnResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgDepositForBurnResponse;
    fromJSON(object: any): MsgDepositForBurnResponse;
    toJSON(message: MsgDepositForBurnResponse): unknown;
    create<I extends Exact<DeepPartial<MsgDepositForBurnResponse>, I>>(base?: I): MsgDepositForBurnResponse;
    fromPartial<I extends Exact<DeepPartial<MsgDepositForBurnResponse>, I>>(object: I): MsgDepositForBurnResponse;
};
export declare const MsgDepositForBurnWithCaller: {
    encode(message: MsgDepositForBurnWithCaller, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgDepositForBurnWithCaller;
    fromJSON(object: any): MsgDepositForBurnWithCaller;
    toJSON(message: MsgDepositForBurnWithCaller): unknown;
    create<I extends Exact<DeepPartial<MsgDepositForBurnWithCaller>, I>>(base?: I): MsgDepositForBurnWithCaller;
    fromPartial<I extends Exact<DeepPartial<MsgDepositForBurnWithCaller>, I>>(object: I): MsgDepositForBurnWithCaller;
};
export declare const MsgDepositForBurnWithCallerResponse: {
    encode(message: MsgDepositForBurnWithCallerResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgDepositForBurnWithCallerResponse;
    fromJSON(object: any): MsgDepositForBurnWithCallerResponse;
    toJSON(message: MsgDepositForBurnWithCallerResponse): unknown;
    create<I extends Exact<DeepPartial<MsgDepositForBurnWithCallerResponse>, I>>(base?: I): MsgDepositForBurnWithCallerResponse;
    fromPartial<I extends Exact<DeepPartial<MsgDepositForBurnWithCallerResponse>, I>>(object: I): MsgDepositForBurnWithCallerResponse;
};
export declare const MsgReplaceDepositForBurn: {
    encode(message: MsgReplaceDepositForBurn, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgReplaceDepositForBurn;
    fromJSON(object: any): MsgReplaceDepositForBurn;
    toJSON(message: MsgReplaceDepositForBurn): unknown;
    create<I extends Exact<DeepPartial<MsgReplaceDepositForBurn>, I>>(base?: I): MsgReplaceDepositForBurn;
    fromPartial<I extends Exact<DeepPartial<MsgReplaceDepositForBurn>, I>>(object: I): MsgReplaceDepositForBurn;
};
export declare const MsgReplaceDepositForBurnResponse: {
    encode(_: MsgReplaceDepositForBurnResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgReplaceDepositForBurnResponse;
    fromJSON(_: any): MsgReplaceDepositForBurnResponse;
    toJSON(_: MsgReplaceDepositForBurnResponse): unknown;
    create<I extends Exact<DeepPartial<MsgReplaceDepositForBurnResponse>, I>>(base?: I): MsgReplaceDepositForBurnResponse;
    fromPartial<I extends Exact<DeepPartial<MsgReplaceDepositForBurnResponse>, I>>(_: I): MsgReplaceDepositForBurnResponse;
};
export declare const MsgReceiveMessage: {
    encode(message: MsgReceiveMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgReceiveMessage;
    fromJSON(object: any): MsgReceiveMessage;
    toJSON(message: MsgReceiveMessage): unknown;
    create<I extends Exact<DeepPartial<MsgReceiveMessage>, I>>(base?: I): MsgReceiveMessage;
    fromPartial<I extends Exact<DeepPartial<MsgReceiveMessage>, I>>(object: I): MsgReceiveMessage;
};
export declare const MsgReceiveMessageResponse: {
    encode(message: MsgReceiveMessageResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgReceiveMessageResponse;
    fromJSON(object: any): MsgReceiveMessageResponse;
    toJSON(message: MsgReceiveMessageResponse): unknown;
    create<I extends Exact<DeepPartial<MsgReceiveMessageResponse>, I>>(base?: I): MsgReceiveMessageResponse;
    fromPartial<I extends Exact<DeepPartial<MsgReceiveMessageResponse>, I>>(object: I): MsgReceiveMessageResponse;
};
export declare const MsgSendMessage: {
    encode(message: MsgSendMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgSendMessage;
    fromJSON(object: any): MsgSendMessage;
    toJSON(message: MsgSendMessage): unknown;
    create<I extends Exact<DeepPartial<MsgSendMessage>, I>>(base?: I): MsgSendMessage;
    fromPartial<I extends Exact<DeepPartial<MsgSendMessage>, I>>(object: I): MsgSendMessage;
};
export declare const MsgSendMessageResponse: {
    encode(message: MsgSendMessageResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgSendMessageResponse;
    fromJSON(object: any): MsgSendMessageResponse;
    toJSON(message: MsgSendMessageResponse): unknown;
    create<I extends Exact<DeepPartial<MsgSendMessageResponse>, I>>(base?: I): MsgSendMessageResponse;
    fromPartial<I extends Exact<DeepPartial<MsgSendMessageResponse>, I>>(object: I): MsgSendMessageResponse;
};
export declare const MsgSendMessageWithCaller: {
    encode(message: MsgSendMessageWithCaller, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgSendMessageWithCaller;
    fromJSON(object: any): MsgSendMessageWithCaller;
    toJSON(message: MsgSendMessageWithCaller): unknown;
    create<I extends Exact<DeepPartial<MsgSendMessageWithCaller>, I>>(base?: I): MsgSendMessageWithCaller;
    fromPartial<I extends Exact<DeepPartial<MsgSendMessageWithCaller>, I>>(object: I): MsgSendMessageWithCaller;
};
export declare const MsgSendMessageWithCallerResponse: {
    encode(message: MsgSendMessageWithCallerResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgSendMessageWithCallerResponse;
    fromJSON(object: any): MsgSendMessageWithCallerResponse;
    toJSON(message: MsgSendMessageWithCallerResponse): unknown;
    create<I extends Exact<DeepPartial<MsgSendMessageWithCallerResponse>, I>>(base?: I): MsgSendMessageWithCallerResponse;
    fromPartial<I extends Exact<DeepPartial<MsgSendMessageWithCallerResponse>, I>>(object: I): MsgSendMessageWithCallerResponse;
};
export declare const MsgReplaceMessage: {
    encode(message: MsgReplaceMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgReplaceMessage;
    fromJSON(object: any): MsgReplaceMessage;
    toJSON(message: MsgReplaceMessage): unknown;
    create<I extends Exact<DeepPartial<MsgReplaceMessage>, I>>(base?: I): MsgReplaceMessage;
    fromPartial<I extends Exact<DeepPartial<MsgReplaceMessage>, I>>(object: I): MsgReplaceMessage;
};
export declare const MsgReplaceMessageResponse: {
    encode(_: MsgReplaceMessageResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgReplaceMessageResponse;
    fromJSON(_: any): MsgReplaceMessageResponse;
    toJSON(_: MsgReplaceMessageResponse): unknown;
    create<I extends Exact<DeepPartial<MsgReplaceMessageResponse>, I>>(base?: I): MsgReplaceMessageResponse;
    fromPartial<I extends Exact<DeepPartial<MsgReplaceMessageResponse>, I>>(_: I): MsgReplaceMessageResponse;
};
export declare const MsgUpdateSignatureThreshold: {
    encode(message: MsgUpdateSignatureThreshold, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateSignatureThreshold;
    fromJSON(object: any): MsgUpdateSignatureThreshold;
    toJSON(message: MsgUpdateSignatureThreshold): unknown;
    create<I extends Exact<DeepPartial<MsgUpdateSignatureThreshold>, I>>(base?: I): MsgUpdateSignatureThreshold;
    fromPartial<I extends Exact<DeepPartial<MsgUpdateSignatureThreshold>, I>>(object: I): MsgUpdateSignatureThreshold;
};
export declare const MsgUpdateSignatureThresholdResponse: {
    encode(_: MsgUpdateSignatureThresholdResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateSignatureThresholdResponse;
    fromJSON(_: any): MsgUpdateSignatureThresholdResponse;
    toJSON(_: MsgUpdateSignatureThresholdResponse): unknown;
    create<I extends Exact<DeepPartial<MsgUpdateSignatureThresholdResponse>, I>>(base?: I): MsgUpdateSignatureThresholdResponse;
    fromPartial<I extends Exact<DeepPartial<MsgUpdateSignatureThresholdResponse>, I>>(_: I): MsgUpdateSignatureThresholdResponse;
};
export declare const MsgLinkTokenPair: {
    encode(message: MsgLinkTokenPair, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgLinkTokenPair;
    fromJSON(object: any): MsgLinkTokenPair;
    toJSON(message: MsgLinkTokenPair): unknown;
    create<I extends Exact<DeepPartial<MsgLinkTokenPair>, I>>(base?: I): MsgLinkTokenPair;
    fromPartial<I extends Exact<DeepPartial<MsgLinkTokenPair>, I>>(object: I): MsgLinkTokenPair;
};
export declare const MsgLinkTokenPairResponse: {
    encode(_: MsgLinkTokenPairResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgLinkTokenPairResponse;
    fromJSON(_: any): MsgLinkTokenPairResponse;
    toJSON(_: MsgLinkTokenPairResponse): unknown;
    create<I extends Exact<DeepPartial<MsgLinkTokenPairResponse>, I>>(base?: I): MsgLinkTokenPairResponse;
    fromPartial<I extends Exact<DeepPartial<MsgLinkTokenPairResponse>, I>>(_: I): MsgLinkTokenPairResponse;
};
export declare const MsgUnlinkTokenPair: {
    encode(message: MsgUnlinkTokenPair, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUnlinkTokenPair;
    fromJSON(object: any): MsgUnlinkTokenPair;
    toJSON(message: MsgUnlinkTokenPair): unknown;
    create<I extends Exact<DeepPartial<MsgUnlinkTokenPair>, I>>(base?: I): MsgUnlinkTokenPair;
    fromPartial<I extends Exact<DeepPartial<MsgUnlinkTokenPair>, I>>(object: I): MsgUnlinkTokenPair;
};
export declare const MsgUnlinkTokenPairResponse: {
    encode(_: MsgUnlinkTokenPairResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUnlinkTokenPairResponse;
    fromJSON(_: any): MsgUnlinkTokenPairResponse;
    toJSON(_: MsgUnlinkTokenPairResponse): unknown;
    create<I extends Exact<DeepPartial<MsgUnlinkTokenPairResponse>, I>>(base?: I): MsgUnlinkTokenPairResponse;
    fromPartial<I extends Exact<DeepPartial<MsgUnlinkTokenPairResponse>, I>>(_: I): MsgUnlinkTokenPairResponse;
};
export declare const MsgAddRemoteTokenMessenger: {
    encode(message: MsgAddRemoteTokenMessenger, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgAddRemoteTokenMessenger;
    fromJSON(object: any): MsgAddRemoteTokenMessenger;
    toJSON(message: MsgAddRemoteTokenMessenger): unknown;
    create<I extends Exact<DeepPartial<MsgAddRemoteTokenMessenger>, I>>(base?: I): MsgAddRemoteTokenMessenger;
    fromPartial<I extends Exact<DeepPartial<MsgAddRemoteTokenMessenger>, I>>(object: I): MsgAddRemoteTokenMessenger;
};
export declare const MsgAddRemoteTokenMessengerResponse: {
    encode(_: MsgAddRemoteTokenMessengerResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgAddRemoteTokenMessengerResponse;
    fromJSON(_: any): MsgAddRemoteTokenMessengerResponse;
    toJSON(_: MsgAddRemoteTokenMessengerResponse): unknown;
    create<I extends Exact<DeepPartial<MsgAddRemoteTokenMessengerResponse>, I>>(base?: I): MsgAddRemoteTokenMessengerResponse;
    fromPartial<I extends Exact<DeepPartial<MsgAddRemoteTokenMessengerResponse>, I>>(_: I): MsgAddRemoteTokenMessengerResponse;
};
export declare const MsgRemoveRemoteTokenMessenger: {
    encode(message: MsgRemoveRemoteTokenMessenger, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgRemoveRemoteTokenMessenger;
    fromJSON(object: any): MsgRemoveRemoteTokenMessenger;
    toJSON(message: MsgRemoveRemoteTokenMessenger): unknown;
    create<I extends Exact<DeepPartial<MsgRemoveRemoteTokenMessenger>, I>>(base?: I): MsgRemoveRemoteTokenMessenger;
    fromPartial<I extends Exact<DeepPartial<MsgRemoveRemoteTokenMessenger>, I>>(object: I): MsgRemoveRemoteTokenMessenger;
};
export declare const MsgRemoveRemoteTokenMessengerResponse: {
    encode(_: MsgRemoveRemoteTokenMessengerResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgRemoveRemoteTokenMessengerResponse;
    fromJSON(_: any): MsgRemoveRemoteTokenMessengerResponse;
    toJSON(_: MsgRemoveRemoteTokenMessengerResponse): unknown;
    create<I extends Exact<DeepPartial<MsgRemoveRemoteTokenMessengerResponse>, I>>(base?: I): MsgRemoveRemoteTokenMessengerResponse;
    fromPartial<I extends Exact<DeepPartial<MsgRemoveRemoteTokenMessengerResponse>, I>>(_: I): MsgRemoveRemoteTokenMessengerResponse;
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
