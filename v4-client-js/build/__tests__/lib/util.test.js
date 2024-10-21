"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const constants_1 = require("../../src/lib/constants");
const utils_1 = require("../../src/lib/utils");
describe('lib/util', () => {
    describe('randomInt', () => {
        it('random values', () => {
            const maxValue = 9999;
            let lastValue = 0;
            for (let i = 0; i < 100; i++) {
                const value = (0, utils_1.randomInt)(maxValue);
                // Within the expected bounds.
                expect(value).toBeGreaterThanOrEqual(0);
                expect(value).toBeLessThanOrEqual(maxValue);
                // No collision.
                expect(value).not.toEqual(lastValue);
                lastValue = value;
            }
        });
        it('zero', () => {
            expect((0, utils_1.randomInt)(0)).toEqual(0);
        });
    });
    describe('generateRandomClientId', () => {
        it('random values', () => {
            let lastValue = 0;
            for (let i = 0; i < 100; i++) {
                const value = (0, utils_1.generateRandomClientId)();
                // Within the expected bounds.
                expect(value).toBeGreaterThanOrEqual(0);
                expect(value).toBeLessThanOrEqual(constants_1.MAX_UINT_32);
                // No collision.
                expect(value).not.toEqual(lastValue);
                lastValue = value;
            }
        });
    });
    describe('clientIdFromString', () => {
        it('hard-coded', () => {
            expect((0, utils_1.clientIdFromString)('test')).toEqual(2151040146);
        });
        it('random values', () => {
            let lastValue = 0;
            let lastInput = 0;
            for (let i = 0; i < 1000; i++) {
                // Prevent input collision.
                let input = (0, utils_1.randomInt)(constants_1.MAX_UINT_32);
                if (input === lastInput) {
                    input += 1;
                }
                const value = (0, utils_1.clientIdFromString)(`${input}`);
                const valueAgain = (0, utils_1.clientIdFromString)(`${input}`);
                // Deterministic.
                expect(value).toEqual(valueAgain);
                // Within the expected bounds.
                expect(value).toBeGreaterThanOrEqual(0);
                expect(value).toBeLessThanOrEqual(constants_1.MAX_UINT_32);
                // No collision.
                expect(value).not.toEqual(lastValue);
                expect(value).not.toEqual(input);
                lastValue = value;
                lastInput = input;
            }
        });
    });
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoidXRpbC50ZXN0LmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vX190ZXN0c19fL2xpYi91dGlsLnRlc3QudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7QUFBQSx1REFFaUM7QUFDakMsK0NBSTZCO0FBRTdCLFFBQVEsQ0FBQyxVQUFVLEVBQUUsR0FBRyxFQUFFO0lBQ3hCLFFBQVEsQ0FBQyxXQUFXLEVBQUUsR0FBRyxFQUFFO1FBQ3pCLEVBQUUsQ0FBQyxlQUFlLEVBQUUsR0FBRyxFQUFFO1lBQ3ZCLE1BQU0sUUFBUSxHQUFXLElBQUksQ0FBQztZQUM5QixJQUFJLFNBQVMsR0FBVyxDQUFDLENBQUM7WUFDMUIsS0FBSyxJQUFJLENBQUMsR0FBRyxDQUFDLEVBQUUsQ0FBQyxHQUFHLEdBQUcsRUFBRSxDQUFDLEVBQUUsRUFBRSxDQUFDO2dCQUM3QixNQUFNLEtBQUssR0FBVyxJQUFBLGlCQUFTLEVBQUMsUUFBUSxDQUFDLENBQUM7Z0JBRTFDLDhCQUE4QjtnQkFDOUIsTUFBTSxDQUFDLEtBQUssQ0FBQyxDQUFDLHNCQUFzQixDQUFDLENBQUMsQ0FBQyxDQUFDO2dCQUN4QyxNQUFNLENBQUMsS0FBSyxDQUFDLENBQUMsbUJBQW1CLENBQUMsUUFBUSxDQUFDLENBQUM7Z0JBRTVDLGdCQUFnQjtnQkFDaEIsTUFBTSxDQUFDLEtBQUssQ0FBQyxDQUFDLEdBQUcsQ0FBQyxPQUFPLENBQUMsU0FBUyxDQUFDLENBQUM7Z0JBQ3JDLFNBQVMsR0FBRyxLQUFLLENBQUM7WUFDcEIsQ0FBQztRQUNILENBQUMsQ0FBQyxDQUFDO1FBRUgsRUFBRSxDQUFDLE1BQU0sRUFBRSxHQUFHLEVBQUU7WUFDZCxNQUFNLENBQUMsSUFBQSxpQkFBUyxFQUFDLENBQUMsQ0FBQyxDQUFDLENBQUMsT0FBTyxDQUFDLENBQUMsQ0FBQyxDQUFDO1FBQ2xDLENBQUMsQ0FBQyxDQUFDO0lBQ0wsQ0FBQyxDQUFDLENBQUM7SUFFSCxRQUFRLENBQUMsd0JBQXdCLEVBQUUs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