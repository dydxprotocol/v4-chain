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
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoidXRpbC50ZXN0LmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vX190ZXN0c19fL2xpYi91dGlsLnRlc3QudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7QUFBQSx1REFFaUM7QUFDakMsK0NBSTZCO0FBRTdCLFFBQVEsQ0FBQyxVQUFVLEVBQUUsR0FBRyxFQUFFO0lBQ3hCLFFBQVEsQ0FBQyxXQUFXLEVBQUUsR0FBRyxFQUFFO1FBQ3pCLEVBQUUsQ0FBQyxlQUFlLEVBQUUsR0FBRyxFQUFFO1lBQ3ZCLE1BQU0sUUFBUSxHQUFXLElBQUksQ0FBQztZQUM5QixJQUFJLFNBQVMsR0FBVyxDQUFDLENBQUM7WUFDMUIsS0FBSyxJQUFJLENBQUMsR0FBRyxDQUFDLEVBQUUsQ0FBQyxHQUFHLEdBQUcsRUFBRSxDQUFDLEVBQUUsRUFBRTtnQkFDNUIsTUFBTSxLQUFLLEdBQVcsSUFBQSxpQkFBUyxFQUFDLFFBQVEsQ0FBQyxDQUFDO2dCQUUxQyw4QkFBOEI7Z0JBQzlCLE1BQU0sQ0FBQyxLQUFLLENBQUMsQ0FBQyxzQkFBc0IsQ0FBQyxDQUFDLENBQUMsQ0FBQztnQkFDeEMsTUFBTSxDQUFDLEtBQUssQ0FBQyxDQUFDLG1CQUFtQixDQUFDLFFBQVEsQ0FBQyxDQUFDO2dCQUU1QyxnQkFBZ0I7Z0JBQ2hCLE1BQU0sQ0FBQyxLQUFLLENBQUMsQ0FBQyxHQUFHLENBQUMsT0FBTyxDQUFDLFNBQVMsQ0FBQyxDQUFDO2dCQUNyQyxTQUFTLEdBQUcsS0FBSyxDQUFDO2FBQ25CO1FBQ0gsQ0FBQyxDQUFDLENBQUM7UUFFSCxFQUFFLENBQUMsTUFBTSxFQUFFLEdBQUcsRUFBRTtZQUNkLE1BQU0sQ0FBQyxJQUFBLGlCQUFTLEVBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQyxPQUFPLENBQUMsQ0FBQyxDQUFDLENBQUM7UUFDbEMsQ0FBQyxDQUFDLENBQUM7SUFDTCxDQUFDLENBQUMsQ0FBQztJQUVILFFBQVEsQ0FBQyx3QkFBd0IsRUFBRSxHQUFHLEVBQUU7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