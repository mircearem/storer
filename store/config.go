package store

import "time"

// Error types
type dberror uint8

const (
	ERR_COL_MAKE_FAIL dberror = iota + 1 // Collection make fail
	ERR_COL_DROP_FAIL                    // Collection drop fail
	ERR_COL_NOT_FOUND                    // Collection not found
	ERR_PUT_FAIL                         // Put fail, undefined
	ERR_PUT_FAIL_CONF                    // Put fail, conflict
	ERR_GET_FAIL_UDEF                    // Get fail, undefined
	ERR_GET_FAIL_NOTF                    // Get fail, not found
	ERR_UPD_FAIL_UDEF                    // Update fail, undefined
	ERR_UPD_FAIL_NOTF                    // Update fail, not found
	ERR_DEL_FAIL_UDEF                    // Delete fail, undefined
	ERR_DEL_FAIL_NOTF                    // Delete fail, not found
)

// Storage configuration
type StoreConfig struct {
	DbName  string
	timeout time.Duration
}

// Create a new configuration struct for the store
func NewStoreConfig() StoreConfig {
	return StoreConfig{
		DbName:  "app",
		timeout: 0,
	}
}

// Change the name of the configuration
func (s StoreConfig) WithDbName(name string) StoreConfig {
	s.DbName = name
	return s
}

// Change the timeout of the database
func (s StoreConfig) WithTimeout(t time.Duration) StoreConfig {
	s.timeout = t
	return s
}

// Error struct to make error handling better
type StoreError struct {
	errorType    dberror
	errorMessage string
}

// Return the error message
func (s StoreError) Error() string {
	return s.errorMessage
}

// Return the error type
func (s StoreError) Type() dberror {
	return s.errorType
}

// Create a new error message
func NewStoreError(typ dberror, msg string) StoreError {
	return StoreError{
		errorType:    typ,
		errorMessage: msg,
	}
}
