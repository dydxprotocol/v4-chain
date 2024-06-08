export interface AminoException {
    aminoType?: string;
    toAmino?: any;
    fromAmino?: any;
}
export interface AminoExceptions {
    [key: string]: AminoException;
}
export declare const DEFAULT_AMINO_EXCEPTIONS: AminoExceptions;
