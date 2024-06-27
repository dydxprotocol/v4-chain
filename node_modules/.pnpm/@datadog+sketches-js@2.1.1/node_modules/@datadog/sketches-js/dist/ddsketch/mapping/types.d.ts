import type { IIndexMapping } from '../proto/compiled';
export interface Mapping {
    relativeAccuracy: number;
    gamma: number;
    minPossible: number;
    maxPossible: number;
    key: (value: number) => number;
    value: (key: number) => number;
    toProto(): IIndexMapping;
}
