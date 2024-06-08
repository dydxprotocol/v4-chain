import {
  MAX_UINT_32,
} from '../../src/lib/constants';
import {
  randomInt,
  generateRandomClientId,
  clientIdFromString,
} from '../../src/lib/utils';

describe('lib/util', () => {
  describe('randomInt', () => {
    it('random values', () => {
      const maxValue: number = 9999;
      let lastValue: number = 0;
      for (let i = 0; i < 100; i++) {
        const value: number = randomInt(maxValue);

        // Within the expected bounds.
        expect(value).toBeGreaterThanOrEqual(0);
        expect(value).toBeLessThanOrEqual(maxValue);

        // No collision.
        expect(value).not.toEqual(lastValue);
        lastValue = value;
      }
    });

    it('zero', () => {
      expect(randomInt(0)).toEqual(0);
    });
  });

  describe('generateRandomClientId', () => {
    it('random values', () => {
      let lastValue: number = 0;
      for (let i = 0; i < 100; i++) {
        const value: number = generateRandomClientId();

        // Within the expected bounds.
        expect(value).toBeGreaterThanOrEqual(0);
        expect(value).toBeLessThanOrEqual(MAX_UINT_32);

        // No collision.
        expect(value).not.toEqual(lastValue);
        lastValue = value;
      }
    });
  });

  describe('clientIdFromString', () => {
    it('hard-coded', () => {
      expect(clientIdFromString('test')).toEqual(2151040146);
    });

    it('random values', () => {
      let lastValue: number = 0;
      let lastInput: number = 0;
      for (let i = 0; i < 1000; i++) {
        // Prevent input collision.
        let input: number = randomInt(MAX_UINT_32);
        if (input === lastInput) {
          input += 1;
        }

        const value: number = clientIdFromString(`${input}`);
        const valueAgain: number = clientIdFromString(`${input}`);

        // Deterministic.
        expect(value).toEqual(valueAgain);

        // Within the expected bounds.
        expect(value).toBeGreaterThanOrEqual(0);
        expect(value).toBeLessThanOrEqual(MAX_UINT_32);

        // No collision.
        expect(value).not.toEqual(lastValue);
        expect(value).not.toEqual(input);
        lastValue = value;
        lastInput = input;
      }
    });
  });
});
