package util

import "github.com/lib/pq"

const (
	PqUniqueViolationError                          = pq.ErrorCode("23505") // 'unique_violation'
	PqSchemaAndDataStatementMixingNotSupportedError = pq.ErrorCode("25007") // 'schema_and_data_statement_mixing_not_supported'
)
