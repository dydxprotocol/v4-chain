/**
 * This function does exhaustiveness checking to ensure that you have discriminated a union so that no type remains. Use this to get the typescript compiler to help discover cases that were not considered.
 */
export declare function assertNever(value: never): never;
