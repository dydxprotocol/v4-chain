import seedrandom from 'seedrandom';

import config from '../config';
import adjectives from './adjectives.json';
import nouns from './nouns.json';

export function generateUsernameForSubaccount(
  subaccountId: string,
  subaccountNum: number,
  nounce: number = 0, // incremented in case of collision
): string {
  const rng = seedrandom(`${subaccountId}/${subaccountNum}/${nounce}`);
  const randomAdjective: string = adjectives[Math.floor(rng() * adjectives.length)];
  const randomNoun: string = nouns[Math.floor(rng() * nouns.length)];
  const randomNumber: string = Math.floor(rng() * 1000).toString().padStart(
    config.SUBACCOUNT_USERNAME_NUM_RANDOM_DIGITS, '0');

  const capitalizedAdjective: string = randomAdjective.charAt(
    0).toUpperCase() + randomAdjective.slice(1);
  const capitalizedNoun: string = randomNoun.charAt(0).toUpperCase() + randomNoun.slice(1);

  return `${capitalizedAdjective}${capitalizedNoun}${randomNumber}`;
}
