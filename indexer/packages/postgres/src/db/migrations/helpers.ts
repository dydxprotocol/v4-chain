// copied from this answer: https://github.com/knex/knex/issues/1699
// Generates raw SQL to updating check constraints on an enum column with minimal locking.
export function formatAlterTableEnumSql(
  tableName: string,
  columnName: string,
  enums: string[],
): string {
  const constraintName = `${tableName}_${columnName}_check`;
  return [
    `ALTER TABLE "${tableName}" DROP CONSTRAINT IF EXISTS "${constraintName}";`,
    `ALTER TABLE "${tableName}" ADD CONSTRAINT "${constraintName}" CHECK ("${columnName}" = ANY (ARRAY['${enums.join(
      '\'::text, \'',
    )}'::text])) NOT VALID;`,
  ].join('\n');
}
