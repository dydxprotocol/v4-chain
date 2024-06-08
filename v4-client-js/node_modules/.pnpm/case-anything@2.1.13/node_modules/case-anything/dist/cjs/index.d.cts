/**
 * # 🐪 camelCase
 * converts a string to camelCase
 * - first lowercase then all capitalised
 * - *strips away* special characters by default
 *
 * @example
 *   camelCase('$catDog') === 'catDog'
 * @example
 *   camelCase('$catDog', { keepSpecialCharacters: true }) === '$catDog'
 */
declare function camelCase(string: string, options?: {
    keepSpecialCharacters?: boolean;
    keep?: string[];
}): string;
/**
 * # 🐫 PascalCase
 * converts a string to PascalCase (also called UpperCamelCase)
 * - all capitalised
 * - *strips away* special characters by default
 *
 * @example
 *   pascalCase('$catDog') === 'CatDog'
 * @example
 *   pascalCase('$catDog', { keepSpecialCharacters: true }) === '$CatDog'
 */
declare function pascalCase(string: string, options?: {
    keepSpecialCharacters?: boolean;
    keep?: string[];
}): string;
/**
 * # 🐫 UpperCamelCase
 * converts a string to UpperCamelCase (also called PascalCase)
 * - all capitalised
 * - *strips away* special characters by default
 *
 * @example
 *   upperCamelCase('$catDog') === 'CatDog'
 * @example
 *   upperCamelCase('$catDog', { keepSpecialCharacters: true }) === '$CatDog'
 */
declare const upperCamelCase: typeof pascalCase;
/**
 * # 🥙 kebab-case
 * converts a string to kebab-case
 * - hyphenated lowercase
 * - *strips away* special characters by default
 *
 * @example
 *   kebabCase('$catDog') === 'cat-dog'
 * @example
 *   kebabCase('$catDog', { keepSpecialCharacters: true }) === '$cat-dog'
 */
declare function kebabCase(string: string, options?: {
    keepSpecialCharacters?: boolean;
    keep?: string[];
}): string;
/**
 * # 🐍 snake_case
 * converts a string to snake_case
 * - underscored lowercase
 * - *strips away* special characters by default
 *
 * @example
 *   snakeCase('$catDog') === 'cat_dog'
 * @example
 *   snakeCase('$catDog', { keepSpecialCharacters: true }) === '$cat_dog'
 */
declare function snakeCase(string: string, options?: {
    keepSpecialCharacters?: boolean;
    keep?: string[];
}): string;
/**
 * # 📣 CONSTANT_CASE
 * converts a string to CONSTANT_CASE
 * - underscored uppercase
 * - *strips away* special characters by default
 *
 * @example
 *   constantCase('$catDog') === 'CAT_DOG'
 * @example
 *   constantCase('$catDog', { keepSpecialCharacters: true }) === '$CAT_DOG'
 */
declare function constantCase(string: string, options?: {
    keepSpecialCharacters?: boolean;
    keep?: string[];
}): string;
/**
 * # 🚂 Train-Case
 * converts strings to Train-Case
 * - hyphenated & capitalised
 * - *strips away* special characters by default
 *
 * @example
 *   trainCase('$catDog') === 'Cat-Dog'
 * @example
 *   trainCase('$catDog', { keepSpecialCharacters: true }) === '$Cat-Dog'
 */
declare function trainCase(string: string, options?: {
    keepSpecialCharacters?: boolean;
    keep?: string[];
}): string;
/**
 * # 🕊 Ada_Case
 * converts a string to Ada_Case
 * - underscored & capitalised
 * - *strips away* special characters by default
 *
 * @example
 *   adaCase('$catDog') === 'Cat_Dog'
 * @example
 *   adaCase('$catDog', { keepSpecialCharacters: true }) === '$Cat_Dog'
 */
declare function adaCase(string: string, options?: {
    keepSpecialCharacters?: boolean;
    keep?: string[];
}): string;
/**
 * # 👔 COBOL-CASE
 * converts a string to COBOL-CASE
 * - hyphenated uppercase
 * - *strips away* special characters by default
 *
 * @example
 *   cobolCase('$catDog') === 'CAT-DOG'
 * @example
 *   cobolCase('$catDog', { keepSpecialCharacters: true }) === '$CAT-DOG'
 */
declare function cobolCase(string: string, options?: {
    keepSpecialCharacters?: boolean;
    keep?: string[];
}): string;
/**
 * # 📍 Dot.notation
 * converts a string to dot.notation
 * - adds dots, does not change casing
 * - *strips away* special characters by default
 *
 * @example
 *   dotNotation('$catDog') === 'cat.Dog'
 * @example
 *   dotNotation('$catDog', { keepSpecialCharacters: true }) === '$cat.Dog'
 */
declare function dotNotation(string: string, options?: {
    keepSpecialCharacters?: boolean;
    keep?: string[];
}): string;
/**
 * # 📂 Path/case
 * converts a string to path/case
 * - adds slashes, does not change casing
 * - *keeps* special characters by default
 *
 * @example
 *   pathCase('$catDog') === '$cat/Dog'
 * @example
 *   pathCase('$catDog', { keepSpecialCharacters: false }) === 'cat/Dog'
 */
declare function pathCase(string: string, options?: {
    keepSpecialCharacters?: boolean;
    keep?: string[];
}): string;
/**
 * # 🛰 Space case
 * converts a string to space case
 * - adds spaces, does not change casing
 * - *keeps* special characters by default
 *
 * @example
 *   spaceCase('$catDog') === '$cat Dog'
 * @example
 *   spaceCase('$catDog', { keepSpecialCharacters: false }) === 'cat Dog'
 */
declare function spaceCase(string: string, options?: {
    keepSpecialCharacters?: boolean;
    keep?: string[];
}): string;
/**
 * # 🏛 Capital Case
 * converts a string to Capital Case
 * - capitalizes words and adds spaces
 * - *keeps* special characters by default
 *
 * @example
 *   capitalCase('$catDog') === '$Cat Dog'
 * @example
 *   capitalCase('$catDog', { keepSpecialCharacters: false }) === 'Cat Dog'
 *
 * ⟪ if you do not want to add spaces, use `pascalCase()` ⟫
 */
declare function capitalCase(string: string, options?: {
    keepSpecialCharacters?: boolean;
    keep?: string[];
}): string;
/**
 * # 🔡 lower case
 * converts a string to lower case
 * - makes words lowercase and adds spaces
 * - *keeps* special characters by default
 *
 * @example
 *   lowerCase('$catDog') === '$cat dog'
 * @example
 *   lowerCase('$catDog', { keepSpecialCharacters: false }) === 'cat dog'
 *
 * ⟪ if you do not want to add spaces, use the native JS `toLowerCase()` ⟫
 */
declare function lowerCase(string: string, options?: {
    keepSpecialCharacters?: boolean;
    keep?: string[];
}): string;
/**
 * # 🔠 UPPER CASE
 * converts a string to UPPER CASE
 * - makes words upper case and adds spaces
 * - *keeps* special characters by default
 *
 * @example
 *   upperCase('$catDog') === '$CAT DOG'
 * @example
 *   upperCase('$catDog', { keepSpecialCharacters: false }) === 'CAT DOG'
 *
 * ⟪ if you do not want to add spaces, use the native JS `toUpperCase()` ⟫
 */
declare function upperCase(string: string, options?: {
    keepSpecialCharacters?: boolean;
    keep?: string[];
}): string;

export { adaCase, camelCase, capitalCase, cobolCase, constantCase, dotNotation, kebabCase, lowerCase, pascalCase, pathCase, snakeCase, spaceCase, trainCase, upperCamelCase, upperCase };
