import 'colors';
export interface InquirerQuestion {
    _?: boolean;
    name: string;
    message: string;
    required?: boolean;
    validate?: Function;
}
export declare const required: (questions: InquirerQuestion[]) => InquirerQuestion[];
export declare const names: (questions: InquirerQuestion[]) => InquirerQuestion[];
export declare const filter: (questions: InquirerQuestion[], answers: {
    [type: string]: any;
}) => (InquirerQuestion | undefined)[];
export declare const _filter: (questions: InquirerQuestion[], answers: {
    [type: string]: any;
}) => {
    [type: string]: any;
};
export declare const prompt: (questions: InquirerQuestion[], answers: {
    [type: string]: any;
}) => Promise<any>;
