"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const types_1 = require("../../src/types");
const constants_1 = require("../../src/lib/constants");
const errors_1 = require("../../src/lib/errors");
const validation_1 = require("../../src/lib/validation");
const constants_2 = require("../helpers/constants");
const long_1 = __importDefault(require("long"));
const MAX_UINT_32_PLUS_1 = constants_1.MAX_UINT_32 + 1;
const MAX_SUBACCOUNT_NUMBER_PLUS_1 = constants_1.MAX_SUBACCOUNT_NUMBER + 1;
describe('Validations', () => {
    it.each([
        ['valid', constants_2.defaultOrder, undefined],
        [
            '0 value clientId',
            { ...constants_2.defaultOrder, clientId: 0 },
            undefined,
        ],
        [
            'underflow clientId',
            { ...constants_2.defaultOrder, clientId: -1 },
            new errors_1.UserError(`clientId: ${-1} is not a valid uint32`),
        ],
        [
            'overflow clientId',
            { ...constants_2.defaultOrder, clientId: MAX_UINT_32_PLUS_1 },
            new errors_1.UserError(`clientId: ${MAX_UINT_32_PLUS_1} is not a valid uint32`),
        ],
        [
            'underflow quantums',
            { ...constants_2.defaultOrder, quantums: long_1.default.NEG_ONE },
            new errors_1.UserError(`quantums: ${-1} cannot be <= 0`),
        ],
        [
            'underflow goodTilBlock',
            { ...constants_2.defaultOrder, goodTilBlock: -1 },
            new errors_1.UserError(`goodTilBlock: ${-1} is not a valid uint32 or is 0`),
        ],
        [
            'overflow goodTilBlock',
            { ...constants_2.defaultOrder, goodTilBlock: MAX_UINT_32_PLUS_1 },
            new errors_1.UserError(`goodTilBlock: ${MAX_UINT_32_PLUS_1} is not a valid uint32 or is 0`),
        ],
        [
            '0 goodTilBlock',
            { ...constants_2.defaultOrder, goodTilBlock: 0 },
            new errors_1.UserError(`goodTilBlock: ${0} is not a valid uint32 or is 0`),
        ],
        [
            'underflow subticks',
            { ...constants_2.defaultOrder, subticks: long_1.default.NEG_ONE },
            new errors_1.UserError(`subticks: ${-1} cannot be <= 0`),
        ],
    ])('Validate order: %s', (_name, order, expectedError) => {
        const validationError = (0, validation_1.validatePlaceOrderMessage)(0, order);
        expect(validationError).toEqual(expectedError);
    });
    it.each([
        ['valid', constants_2.defaultCancelOrder, undefined],
        [
            '0 value clientId',
            { ...constants_2.defaultCancelOrder, clientId: 0 },
            undefined,
        ],
        [
            'underflow clientId',
            { ...constants_2.defaultCancelOrder, clientId: -1 },
            new errors_1.UserError(`clientId: ${-1} is not a valid uint32`),
        ],
        [
            'overflow clientId',
            { ...constants_2.defaultCancelOrder, clientId: MAX_UINT_32_PLUS_1 },
            new errors_1.UserError(`clientId: ${MAX_UINT_32_PLUS_1} is not a valid uint32`),
        ],
        [
            'underflow goodTilBlock',
            { ...constants_2.defaultCancelOrder, goodTilBlock: -1 },
            new errors_1.UserError(`goodTilBlock: ${-1} is not a valid uint32 or is 0`),
        ],
        [
            'overflow goodTilBlock',
            { ...constants_2.defaultCancelOrder, goodTilBlock: MAX_UINT_32_PLUS_1 },
            new errors_1.UserError(`goodTilBlock: ${MAX_UINT_32_PLUS_1} is not a valid uint32 or is 0`),
        ],
        [
            '0 goodTilBlock',
            { ...constants_2.defaultCancelOrder, goodTilBlock: 0 },
            new errors_1.UserError(`goodTilBlock: ${0} is not a valid uint32 or is 0`),
        ],
        [
            'contains GoodTilBlockTime',
            { ...constants_2.defaultCancelOrder, goodTilBlockTime: 10 },
            new errors_1.UserError('goodTilBlockTime is 10, but should not be set for non-stateful orders'),
        ],
        [
            'stateful order - valid',
            {
                ...constants_2.defaultCancelOrder,
                orderFlags: types_1.OrderFlags.LONG_TERM,
                goodTilBlock: undefined,
                goodTilBlockTime: 10,
            },
            undefined,
        ],
        [
            'stateful order - undefined goodTilBlockTime',
            { ...constants_2.defaultCancelOrder, orderFlags: types_1.OrderFlags.LONG_TERM },
            new errors_1.UserError(`goodTilBlockTime: ${undefined} is not a valid uint32 or is 0`),
        ],
        [
            'stateful order - zero goodTilBlockTime',
            { ...constants_2.defaultCancelOrder, orderFlags: types_1.OrderFlags.LONG_TERM, goodTilBlockTime: 0 },
            new errors_1.UserError(`goodTilBlockTime: ${0} is not a valid uint32 or is 0`),
        ],
        [
            'stateful order - underflow goodTilBlockTime',
            { ...constants_2.defaultCancelOrder, orderFlags: types_1.OrderFlags.LONG_TERM, goodTilBlockTime: -1 },
            new errors_1.UserError(`goodTilBlockTime: ${-1} is not a valid uint32 or is 0`),
        ],
        [
            'stateful order - overflow goodTilBlockTime',
            {
                ...constants_2.defaultCancelOrder,
                orderFlags: types_1.OrderFlags.LONG_TERM,
                goodTilBlockTime: MAX_UINT_32_PLUS_1,
            },
            new errors_1.UserError(`goodTilBlockTime: ${MAX_UINT_32_PLUS_1} is not a valid uint32 or is 0`),
        ],
        [
            'stateful order - has GoodTilBlock',
            {
                ...constants_2.defaultCancelOrder,
                orderFlags: types_1.OrderFlags.LONG_TERM,
                goodTilBlock: 10,
                goodTilBlockTime: 10,
            },
            new errors_1.UserError('goodTilBlock is 10, but should not be set for stateful orders'),
        ],
    ])('Validate cancel order: %s', (_name, order, expectedError) => {
        const validationError = (0, validation_1.validateCancelOrderMessage)(0, order);
        expect(validationError).toEqual(expectedError);
    });
    it.each([
        ['valid', constants_2.defaultTransfer, undefined],
        [
            'underflow senderSubaccountNumber',
            { ...constants_2.defaultTransfer, sender: { owner: constants_2.TEST_ADDRESS, number: -1 } },
            new errors_1.UserError(`senderSubaccountNumber: ${-1} cannot be < 0 or > ${constants_1.MAX_SUBACCOUNT_NUMBER}`),
        ],
        [
            'exceeds max subaccount number - senderSubaccountNumber',
            { ...constants_2.defaultTransfer, sender: { owner: constants_2.TEST_ADDRESS, number: MAX_SUBACCOUNT_NUMBER_PLUS_1 } },
            new errors_1.UserError(`senderSubaccountNumber: ${MAX_SUBACCOUNT_NUMBER_PLUS_1} cannot be < 0 or > ${constants_1.MAX_SUBACCOUNT_NUMBER}`),
        ],
        [
            '0 senderSubaccountNumber',
            { ...constants_2.defaultTransfer, sender: { owner: constants_2.TEST_ADDRESS, number: 0 } },
            undefined,
        ],
        [
            'underflow recipientSubaccountNumber',
            { ...constants_2.defaultTransfer, recipient: { owner: constants_2.TEST_ADDRESS, number: -1 } },
            new errors_1.UserError(`recipientSubaccountNumber: ${-1} cannot be < 0 or > ${constants_1.MAX_SUBACCOUNT_NUMBER}`),
        ],
        [
            'exceeds max subaccount number - recipient.subaccountNumber',
            {
                ...constants_2.defaultTransfer,
                recipient: { owner: constants_2.TEST_ADDRESS, number: MAX_SUBACCOUNT_NUMBER_PLUS_1 },
            },
            new errors_1.UserError(`recipientSubaccountNumber: ${MAX_SUBACCOUNT_NUMBER_PLUS_1} cannot be < 0 or > ${constants_1.MAX_SUBACCOUNT_NUMBER}`),
        ],
        [
            '0 recipientSubaccountNumber',
            { ...constants_2.defaultTransfer, recipient: { owner: constants_2.TEST_ADDRESS, number: 0 } },
            undefined,
        ],
        [
            'non-zero asset id',
            { ...constants_2.defaultTransfer, assetId: 1 },
            new errors_1.UserError(`asset id: ${1} not supported`),
        ],
        [
            '0 amount',
            { ...constants_2.defaultTransfer, amount: long_1.default.ZERO },
            new errors_1.UserError(`amount: ${0} cannot be <= 0`),
        ],
        [
            'too short recipientAddress',
            {
                ...constants_2.defaultTransfer,
                recipient: {
                    owner: 'dydx14063jves4u9zhm7eja5ltf3t8zspxd92qnk23',
                    number: 0,
                },
            },
            new errors_1.UserError('Error: Invalid checksum for dydx14063jves4u9zhm7eja5ltf3t8zspxd92qnk23'),
        ],
        [
            'invalid recipientAddress',
            {
                ...constants_2.defaultTransfer,
                recipient: {
                    owner: 'fakeAddress1234',
                    number: 0,
                },
            },
            new errors_1.UserError('Error: Mixed-case string fakeAddress1234'),
        ],
    ])('Validate transfer: %s', (_name, transfer, expectedError) => {
        const validationError = (0, validation_1.validateTransferMessage)(transfer);
        expect(validationError).toEqual(expectedError);
    });
    it.each([
        ['valid', 'dydx17xpfvakm2amg962yls6f84z3kell8c5leqdyt2', true],
        ['invalid: does not start with dydx1', 'dydx27xpfvakm2amg962yls6f84z3kell8c5leqdyt2', false],
        ['invalid: too short', 'dydx17xpfvakm2amg962yls6f84z3kell8c5leqdyt', false],
        ['invalid: too long', 'dydx17xpfvakm2amg962yls6f84z3kell8c5leqdyt2s', false],
    ])('Validate address: %s', (_name, address, expectedResult) => {
        const validationResult = (0, validation_1.isValidAddress)(address);
        expect(validationResult).toEqual(expectedResult);
    });
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoidmFsaWRhdGlvbi50ZXN0LmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vX190ZXN0c19fL2xpYi92YWxpZGF0aW9uLnRlc3QudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7QUFBQSwyQ0FLeUI7QUFDekIsdURBQTZFO0FBQzdFLGlEQUFpRDtBQUNqRCx5REFLa0M7QUFDbEMsb0RBRThCO0FBQzlCLGdEQUF3QjtBQUV4QixNQUFNLGtCQUFrQixHQUFXLHVCQUFXLEdBQUcsQ0FBQyxDQUFDO0FBQ25ELE1BQU0sNEJBQTRCLEdBQVcsaUNBQXFCLEdBQUcsQ0FBQyxDQUFDO0FBRXZFLFFBQVEsQ0FBQyxhQUFhLEVBQUUsR0FBRyxFQUFFO0lBQzNCLEVBQUUsQ0FBQyxJQUFJLENBQUM7UUFDTixDQUFDLE9BQU8sRUFBRSx3QkFBWSxFQUFFLFNBQVMsQ0FBQztRQUNsQztZQUNFLGtCQUFrQjtZQUNsQixFQUFFLEdBQUcsd0JBQVksRUFBRSxRQUFRLEVBQUUsQ0FBQyxFQUFFO1lBQ2hDLFNBQVM7U0FDVjtRQUNEO1lBQ0Usb0JBQW9CO1lBQ3BCLEVBQUUsR0FBRyx3QkFBWSxFQUFFLFFBQVEsRUFBRSxDQUFDLENBQUMsRUFBRTtZQUNqQyxJQUFJLGtCQUFTLENBQUMsYUFBYSxDQUFDLENBQUMsd0JBQXdCLENBQUM7U0FDdkQ7UUFDRDtZQUNFLG1CQUFtQjtZQUNuQixFQUFFLEdBQUcsd0JBQVksRUFBRSxRQUFRLEVBQUUsa0JBQWtCLEVBQUU7WUFDakQsSUFBSSxrQkFBUyxDQUFDLGFBQWEsa0JBQWtCLHdCQUF3QixDQUFDO1NBQ3ZFO1FBQ0Q7WUFDRSxvQkFBb0I7WUFDcEIsRUFBRSxHQUFHLHdCQUFZLEVBQUUsUUFBUSxFQUFFLGNBQUksQ0FBQyxPQUFPLEVBQUU7WUFDM0MsSUFBSSxrQkFBUyxDQUFDLGFBQWEsQ0FBQyxDQUFDLGlCQUFpQixDQUFDO1NBQ2hEO1FBQ0Q7WUFDRSx3QkFBd0I7WUFDeEIsRUFBRSxHQUFHLHdCQUFZ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