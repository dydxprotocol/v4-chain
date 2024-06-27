"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const constants_1 = require("../../src/lib/constants");
const helpers_1 = require("../../src/lib/helpers");
const constants_2 = require("../helpers/constants");
describe('helpers', () => {
    describe('convertPartialTransactionOptionsToFull', () => it.each([
        [
            'partial transactionOptions',
            {
                accountNumber: constants_2.defaultTransactionOptions.accountNumber,
                chainId: constants_2.defaultTransactionOptions.chainId,
            },
            { ...constants_2.defaultTransactionOptions, sequence: constants_1.DEFAULT_SEQUENCE },
        ],
        [
            'undefined transactionOptions',
            undefined,
            undefined,
        ],
    ])('convertPartialTransactionOptionsToFull: %s', (_name, partialTransactionOptions, expectedResult) => {
        const transactionOptions = (0, helpers_1.convertPartialTransactionOptionsToFull)(partialTransactionOptions);
        expect(expectedResult).toEqual(transactionOptions);
    }));
    describe('stripHexPrefix', () => {
        it('strips 0x prefix', () => {
            expect((0, helpers_1.stripHexPrefix)('0x123')).toEqual('123');
        });
        it('returns input if no prefix', () => {
            expect((0, helpers_1.stripHexPrefix)('10x23')).toEqual('10x23');
        });
    });
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiaGVscGVycy50ZXN0LmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vX190ZXN0c19fL2xpYi9oZWxwZXJzLnRlc3QudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7QUFDQSx1REFBMkQ7QUFDM0QsbURBQStGO0FBQy9GLG9EQUFpRTtBQUVqRSxRQUFRLENBQUMsU0FBUyxFQUFFLEdBQUcsRUFBRTtJQUN2QixRQUFRLENBQUMsd0NBQXdDLEVBQUUsR0FBRyxFQUFFLENBQUMsRUFBRSxDQUFDLElBQUksQ0FBQztRQUMvRDtZQUNFLDRCQUE0QjtZQUM1QjtnQkFDRSxhQUFhLEVBQUUscUNBQXlCLENBQUMsYUFBYTtnQkFDdEQsT0FBTyxFQUFFLHFDQUF5QixDQUFDLE9BQU87YUFDM0M7WUFDRCxFQUFFLEdBQUcscUNBQXlCLEVBQUUsUUFBUSxFQUFFLDRCQUFnQixFQUFFO1NBQzdEO1FBQ0Q7WUFDRSw4QkFBOEI7WUFDOUIsU0FBUztZQUNULFNBQVM7U0FDVjtLQUNGLENBQUMsQ0FBQyw0Q0FBNEMsRUFBRSxDQUMvQyxLQUFhLEVBQ2IseUJBQWdFLEVBQ2hFLGNBQThDLEVBQzlDLEVBQUU7UUFDRixNQUFNLGtCQUFrQixHQUE4QixJQUFBLGdEQUFzQyxFQUMxRix5QkFBeUIsQ0FDMUIsQ0FBQztRQUNGLE1BQU0sQ0FBQyxjQUFjLENBQUMsQ0FBQyxPQUFPLENBQUMsa0JBQWtCLENBQUMsQ0FBQztJQUNyRCxDQUFDLENBQUMsQ0FBQyxDQUFDO0lBRUosUUFBUSxDQUFDLGdCQUFnQixFQUFFLEdBQUcsRUFBRTtRQUM5QixFQUFFLENBQUMsa0JBQWtCLEVBQUUsR0FBRyxFQUFFO1lBQzFCLE1BQU0sQ0FBQyxJQUFBLHdCQUFjLEVBQUMsT0FBTyxDQUFDLENBQUMsQ0FBQyxPQUFPLENBQUMsS0FBSyxDQUFDLENBQUM7UUFDakQsQ0FBQyxDQUFDLENBQUM7UUFDSCxFQUFFLENBQUMsNEJBQTRCLEVBQUUsR0FBRyxFQUFFO1lBQ3BDLE1BQU0sQ0FBQyxJQUFBLHdCQUFjLEVBQUMsT0FBTyxDQUFDLENBQUMsQ0FBQyxPQUFPLENBQUMsT0FBTyxDQUFDLENBQUM7UUFDbkQsQ0FBQyxDQUFDLENBQUM7SUFDTCxDQUFDLENBQUMsQ0FBQztBQUNMLENBQUMsQ0FBQyxDQUFDIn0=