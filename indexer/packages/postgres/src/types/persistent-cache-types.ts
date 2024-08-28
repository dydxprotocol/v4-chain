export interface PersistentCacheCreateObject {
  key: string,
  value: string,
}

export interface PersistentCacheUpdateObject {
  key: string,
  value: string,
}

export enum PersistentCacheColumns {
  key = 'key',
  value = 'value',
}
