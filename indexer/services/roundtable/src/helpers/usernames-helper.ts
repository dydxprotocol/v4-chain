import { randomInt } from 'crypto';

import config from '../config';
import adjectives from './adjectives.json';
import nouns from './nouns.json';

export function generateUsername(): string {
  const randomAdjective: string = adjectives[randomInt(0, adjectives.length)];
  const randomNoun: string = nouns[randomInt(0, nouns.length)];
  const randomNumber: string = randomInt(0, 1000).toString().padStart(
    config.SUBACCOUNT_USERNAME_NUM_RANDOM_DIGITS, '0');

  const capitalizedAdjective: string = randomAdjective.charAt(
    0).toUpperCase() + randomAdjective.slice(1);
  const capitalizedNoun: string = randomNoun.charAt(0).toUpperCase() + randomNoun.slice(1);

  return `${capitalizedAdjective}${capitalizedNoun}${randomNumber}`;
}
