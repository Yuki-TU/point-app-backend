// Package entities contains generated code for schema 'point_app'.
package entities

// Code generated by xo. DO NOT EDIT.

import (
	"context"
	"database/sql"
	"fmt"
	"io"
)

var (
	// logf is used by generated code to log SQL queries.
	logf = func(string, ...interface{}) {}
	// errf is used by generated code to log SQL errors.
	errf = func(string, ...interface{}) {}
)

// logerror logs the error and returns it.
func logerror(err error) error {
	errf("ERROR: %v", err)
	return err
}

// Logf logs a message using the package logger.
func Logf(s string, v ...interface{}) {
	logf(s, v...)
}

// SetLogger sets the package logger. Valid logger types:
//
//	io.Writer
//	func(string, ...interface{}) (int, error) // fmt.Printf
//	func(string, ...interface{}) // log.Printf
func SetLogger(logger interface{}) {
	logf = convLogger(logger)
}

// Errorf logs an error message using the package error logger.
func Errorf(s string, v ...interface{}) {
	errf(s, v...)
}

// SetErrorLogger sets the package error logger. Valid logger types:
//
//	io.Writer
//	func(string, ...interface{}) (int, error) // fmt.Printf
//	func(string, ...interface{}) // log.Printf
func SetErrorLogger(logger interface{}) {
	errf = convLogger(logger)
}

// convLogger converts logger to the standard logger interface.
func convLogger(logger interface{}) func(string, ...interface{}) {
	switch z := logger.(type) {
	case io.Writer:
		return func(s string, v ...interface{}) {
			fmt.Fprintf(z, s, v...)
		}
	case func(string, ...interface{}) (int, error): // fmt.Printf
		return func(s string, v ...interface{}) {
			_, _ = z(s, v...)
		}
	case func(string, ...interface{}): // log.Printf
		return z
	}
	panic(fmt.Sprintf("unsupported logger type %T", logger))
}

// DB is the common interface for database operations that can be used with
// types from schema 'point_app'.
//
// This works with both [database/sql.DB] and [database/sql.Tx].
type DB interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

// Error is an error.
type Error string

// Error satisfies the error interface.
func (err Error) Error() string {
	return string(err)
}

// Error values.
const (
	// ErrAlreadyExists is the already exists error.
	ErrAlreadyExists Error = "already exists"
	// ErrDoesNotExist is the does not exist error.
	ErrDoesNotExist Error = "does not exist"
	// ErrMarkedForDeletion is the marked for deletion error.
	ErrMarkedForDeletion Error = "marked for deletion"
)

// ErrInsertFailed is the insert failed error.
type ErrInsertFailed struct {
	Err error
}

// Error satisfies the error interface.
func (err *ErrInsertFailed) Error() string {
	return fmt.Sprintf("insert failed: %v", err.Err)
}

// Unwrap satisfies the unwrap interface.
func (err *ErrInsertFailed) Unwrap() error {
	return err.Err
}

// ErrUpdateFailed is the update failed error.
type ErrUpdateFailed struct {
	Err error
}

// Error satisfies the error interface.
func (err *ErrUpdateFailed) Error() string {
	return fmt.Sprintf("update failed: %v", err.Err)
}

// Unwrap satisfies the unwrap interface.
func (err *ErrUpdateFailed) Unwrap() error {
	return err.Err
}

// ErrUpsertFailed is the upsert failed error.
type ErrUpsertFailed struct {
	Err error
}

// Error satisfies the error interface.
func (err *ErrUpsertFailed) Error() string {
	return fmt.Sprintf("upsert failed: %v", err.Err)
}

// Unwrap satisfies the unwrap interface.
func (err *ErrUpsertFailed) Unwrap() error {
	return err.Err
}
