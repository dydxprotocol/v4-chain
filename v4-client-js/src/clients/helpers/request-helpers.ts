export function generateQueryPath(url: string, params: {}): string {
  const definedEntries = Object.entries(params)
    .filter(([_key, value]: [string, unknown]) => value !== undefined);

  if (!definedEntries.length) {
    return url;
  }

  const paramsString = definedEntries.map(
    ([key, value]: [string, unknown]) => `${key}=${value}`,
  ).join('&');
  return `${url}?${paramsString}`;
}
