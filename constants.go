package filter

// Data type constants identify the expected value type of a [Field].
const (
	DataTypeStr      = "string"   // DataTypeStr represents a string value.
	DataTypeInt      = "int"      // DataTypeInt represents an integer value.
	DataTypeBool     = "bool"     // DataTypeBool represents a boolean value.
	DataTypeDate     = "date"     // DataTypeDate represents a date (format 2006-01-02).
	DataTypeDateTime = "datetime" // DataTypeDateTime represents a datetime (format 2006-01-02T15:04:05).
	DataTypeUUID     = "uuid"     // DataTypeUUID represents a UUID value.
	DataTypeList     = "list"     // DataTypeList represents a JSON array.
)

// Operator constants define comparison operators used in [Field] filters.
const (
	OperatorEq            = "eq"      // OperatorEq — equal.
	OperatorNotEq         = "neq"     // OperatorNotEq — not equal.
	OperatorLowerThan     = "lt"      // OperatorLowerThan — less than.
	OperatorLowerThanEq   = "lte"     // OperatorLowerThanEq — less than or equal.
	OperatorGreaterThan   = "gt"      // OperatorGreaterThan — greater than.
	OperatorGreaterThanEq = "gte"     // OperatorGreaterThanEq — greater than or equal.
	OperatorBetween       = "between" // OperatorBetween — between two values.
	OperatorLike          = "like"    // OperatorLike — substring match (SQL LIKE).
	OperatorIn            = "in"      // OperatorIn — value is in a list.
)

// Sort direction constants.
const (
	SortOperatorASC  = "ASC"  // SortOperatorASC sorts in ascending order.
	SortOperatorDESC = "DESC" // SortOperatorDESC sorts in descending order.
)
