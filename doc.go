// Package filter provides a universal filtering, pagination, and sorting
// mechanism for HTTP query parameters. It converts user-supplied query strings
// into typed filter criteria that can be applied to SQL (via goqu) and MongoDB
// queries through the adapter sub-package.
//
// # Overview
//
// The central type is [Options], which holds pagination settings (limit, page),
// sorting parameters (sortBy, desc), and a collection of [Field] filters.
// Create an Options value with [NewOptions]:
//
//	opts := filter.NewOptions(10, 1, "created_at", true)
//	opts.AddField("status", filter.OperatorEq, "active", filter.DataTypeStr)
//
// # Fields and Type Conversion
//
// Each [Field] carries a string Value that is converted to a Go type by
// [Field.GetValue] based on the Field's Type constant:
//
//   - [DataTypeStr]      → string
//   - [DataTypeInt]      → int  (or []int for [OperatorIn])
//   - [DataTypeBool]     → bool
//   - [DataTypeDate]     → time.Time or *[Between] (format "2006-01-02")
//   - [DataTypeDateTime] → time.Time or *[Between] (format "2006-01-02T15:04:05")
//   - [DataTypeUUID]     → uuid.UUID (or []uuid.UUID for [OperatorIn])
//   - [DataTypeList]     → any (JSON unmarshal)
//
// # Operators
//
// Supported comparison operators are defined as constants:
// [OperatorEq], [OperatorNotEq], [OperatorLowerThan], [OperatorLowerThanEq],
// [OperatorGreaterThan], [OperatorGreaterThanEq], [OperatorBetween],
// [OperatorLike], and [OperatorIn].
//
// Operator-datatype compatibility is validated at AddField time.
//
// # Gin Middlewares
//
// Three [gin] middlewares parse query parameters automatically:
//   - [QueryOptionsMiddlewares] stores Options in gin.Context
//   - [SingleQueryOptionsMiddlewares] passes Options to a callback
//   - [SingleQueryOptionsMiddlewaresWithDefaults] same, with custom defaults
//
// # Errors
//
// The package provides typed errors for programmatic handling:
//   - [ErrNilOptions], [ErrNilDataset], [ErrEmptyTypeTag], [ErrNotStruct] — sentinel errors
//   - [*ValidationError]  — invalid operator / data type / combination
//   - [*ConversionError]  — value conversion failure
//   - [*ParseError]       — between-range parsing failure
//
// All custom error types support errors.Is / errors.As unwrapping.
//
// # Adapters
//
// Query building is done through the adapter sub-package:
//   - adapter.GoquQuery  — builds a goqu SelectDataset
//   - adapter.MongoQueryD — builds a bson.D filter + FindOptions
//
// See the [adapter] package for details.
package filter
