declare namespace DbTypes {

  class DBError extends Error {
    name: string;
    nativeError: Error;
  }

  class CheckViolationError extends DBError {
    table: string;
    constraint: string;
  }

  class ConstraintViolationError extends DBError {}

  class DataError extends DBError {}

  class ForeignKeyViolationError extends ConstraintViolationError {
    table: string;
    constraint: string;
    schema?: string;
  }

  class NotNullViolationError extends ConstraintViolationError {
    table: string;
    column: string;
    database?: string;
    schema?: string;
  }

  class UniqueViolationError extends ConstraintViolationError {
    table: string;
    columns: string[];
    constraint: string;
    schema?: string;
  }

  function wrapError(err: Error): DBError

  export {
    wrapError,
    DBError,
    CheckViolationError,
    ConstraintViolationError,
    DataError,
    ForeignKeyViolationError,
    NotNullViolationError,
    UniqueViolationError,
  }
}

export = DbTypes
