import { S3_LOCATION_PREFIX } from './aws';

export function castToTimestamp(column: string): string {
  return `CAST("${column}" AS timestamp) as "${column}"`;
}

export function castToDouble(column: string): string {
  return `CAST("${column}" as double) as "${column}"`;
}

export function getExternalAthenaTableCreationStatement(
  tablePrefix: string,
  rdsExportIdentifier: string,
  tableName: string,
  rawColumns: string,
): string {
  return `
      CREATE EXTERNAL TABLE IF NOT EXISTS \`${tablePrefix}_raw_${tableName}\` (
        ${rawColumns}
      )
      ROW FORMAT SERDE 'org.apache.hadoop.hive.ql.io.parquet.serde.ParquetHiveSerDe'
      STORED AS INPUTFORMAT
        'org.apache.hadoop.hive.ql.io.parquet.MapredParquetInputFormat'
      OUTPUTFORMAT
        'org.apache.hadoop.hive.ql.io.parquet.MapredParquetOutputFormat'
      LOCATION '${S3_LOCATION_PREFIX}/${rdsExportIdentifier}/dydx/public.${tableName}'
      TBLPROPERTIES ('has_encrypted_data'='false');
  `;
}

export function getAthenaTableCreationStatement(
  tablePrefix: string,
  tableName: string,
  columns: string,
): string {
  return `
    CREATE TABLE IF NOT EXISTS "${tablePrefix}_${tableName}" as select
      ${columns}
    FROM "${tablePrefix}_raw_${tableName}";
  `;
}
