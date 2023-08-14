import { RegexPattern } from '../types';

export const IntegerPattern: RegexPattern = '^\\d+$';

export const NumericPattern: RegexPattern = '^-?[0-9]+\\.?[0-9]*$';

export const NonNegativeNumericPattern: RegexPattern = '^[0-9]+\\.?[0-9]*$';
