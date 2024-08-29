export interface PaginationFromDatabase<T> {
  results: T[],
  total?: number,
  offset?: number,
  limit?: number,
}
