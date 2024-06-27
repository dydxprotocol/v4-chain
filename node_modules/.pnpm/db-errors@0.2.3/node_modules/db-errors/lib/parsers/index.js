const postgres = require('./postgres/DBError/parser');
const sqlite = require('./sqlite/DBError/parser');
const mysql = require('./mysql/DBError/parser');
const mssql = require('./mssql/DBError/parser');

module.exports = {
  postgres,
  sqlite,
  mysql,
  mssql,
};