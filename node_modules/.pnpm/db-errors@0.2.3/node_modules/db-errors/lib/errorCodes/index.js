'use strict';

const postgres = require('./postgres');
const mysql = require('./mysql');
const mssql = require('./mssql');

module.exports = {
  postgres,
  mysql,
  mssql,
};