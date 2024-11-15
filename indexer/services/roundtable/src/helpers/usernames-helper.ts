import seedrandom from 'seedrandom';

import config from '../config';
import adjectives from './adjectives.json';
import nouns from './nouns.json';

const suffixCharacters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';

export function generateUsernameForSubaccount(
  address: string,
  subaccountNum: number,
  nounce: number = 0, // incremented in case of collision
): string {
  const rng = seedrandom(`${address}/${subaccountNum}/${nounce}`);
  const randomAdjective: string = adjectives[Math.floor(rng() * adjectives.length)];
  const randomNoun: string = nouns[Math.floor(rng() * nouns.length)];
  const randomSuffix: string = Array.from(
    { length: config.SUBACCOUNT_USERNAME_SUFFIX_RANDOM_DIGITS },
    () => suffixCharacters.charAt(Math.floor(rng() * suffixCharacters.length))).join('');

  const capitalizedAdjective: string = randomAdjective.charAt(
    0).toUpperCase() + randomAdjective.slice(1);
  const capitalizedNoun: string = randomNoun.charAt(0).toUpperCase() + randomNoun.slice(1);

  return `${capitalizedAdjective}${capitalizedNoun}${randomSuffix}`;
}
