declare function cuid(): string;
declare namespace cuid {
    function fingerprint(): string
    function isCuid(value: unknown): value is string

    export { fingerprint };
    export { isCuid };
}

export = cuid;
