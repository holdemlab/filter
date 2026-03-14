package filter

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

// Field represents a single filter criterion with a name, comparison operator,
// string value, and target data type.
type (
	Field struct {
		Name     string // Name is the field/column name to filter on.
		Operator string // Operator is the comparison operator (e.g. OperatorEq).
		Value    string // Value is the raw string representation of the filter value.
		Type     string // Type is the data type constant (e.g. DataTypeStr).
	}
	// Between holds the start and end values of a range filter.
	Between struct {
		start any
		end   any
	}
)

// Start returns the start (lower bound) of the range.
func (b *Between) Start() any {
	return b.start
}

// End returns the end (upper bound) of the range.
func (b *Between) End() any {
	return b.end
}

// GetValue converts the string [Field.Value] into a typed Go value
// according to [Field.Type]. It returns a [*ConversionError] when the
// conversion fails and a [*ValidationError] for unknown data types.
func (f *Field) GetValue() (value any, err error) {
	switch f.Type {
	case DataTypeInt:
		if f.Value == "" {
			return 0, nil
		}
		if f.Operator == OperatorIn {
			parts := strings.Split(f.Value, ",")
			ints := make([]int, 0, len(parts))
			for _, p := range parts {
				p = strings.TrimSpace(p)
				v, err := strconv.Atoi(p)
				if err != nil {
					return nil, &ConversionError{DataType: DataTypeInt, Value: p, Err: err}
				}
				ints = append(ints, v)
			}
			return ints, nil
		}
		value, err = strconv.Atoi(f.Value)
		if err != nil {
			return nil, &ConversionError{DataType: DataTypeInt, Value: f.Value, Err: err}
		}

		return value, nil
	case DataTypeStr:
		return f.Value, nil
	case DataTypeBool:
		value, err = strconv.ParseBool(f.Value)
		if err != nil {
			return nil, &ConversionError{DataType: DataTypeBool, Value: f.Value, Err: err}
		}
		return value, nil
	case DataTypeDate:
		switch f.Operator {
		case OperatorBetween:
			return ParseBetweenOperator(formatDate, f.Value)
		case OperatorEq, OperatorNotEq, OperatorLowerThan, OperatorLowerThanEq, OperatorGreaterThan, OperatorGreaterThanEq:
			// check if the value is a timestamp format date YYYY-MM-DD
			if strings.Contains(f.Value, "-") {
				return time.Parse(formatDate, f.Value)
			}
			timestamp, err := strconv.Atoi(f.Value)
			if err != nil {
				return nil, &ConversionError{DataType: DataTypeDate, Value: f.Value, Err: err}
			}
			return time.Unix(int64(timestamp), 0), nil
		default:
			timestamp, err := strconv.Atoi(f.Value)
			if err != nil {
				return nil, &ConversionError{DataType: DataTypeDate, Value: f.Value, Err: err}
			}
			return time.Unix(int64(timestamp), 0), nil
		}
	case DataTypeDateTime:
		switch f.Operator {
		case OperatorBetween:
			return ParseBetweenOperator(formatDateTime, f.Value)
		case OperatorEq, OperatorNotEq, OperatorLowerThan, OperatorLowerThanEq, OperatorGreaterThan, OperatorGreaterThanEq:
			// check if the value is a timestamp
			if strings.Contains(f.Value, "T") {
				return time.Parse(formatDateTime, f.Value)
			}
			timestamp, err := strconv.Atoi(f.Value)
			if err != nil {
				return nil, &ConversionError{DataType: DataTypeDateTime, Value: f.Value, Err: err}
			}
			return time.Unix(int64(timestamp), 0), nil
		default:
			timestamp, err := strconv.Atoi(f.Value)
			if err != nil {
				return nil, &ConversionError{DataType: DataTypeDateTime, Value: f.Value, Err: err}
			}
			return time.Unix(int64(timestamp), 0), nil
		}
	case DataTypeUUID:
		if f.Operator == OperatorIn {
			parts := strings.Split(f.Value, ",")
			uids := make([]uuid.UUID, 0, len(parts))
			for _, p := range parts {
				p = strings.TrimSpace(p)
				uid, err := uuid.FromString(p)
				if err != nil {
					return nil, &ConversionError{DataType: DataTypeUUID, Value: p, Err: err}
				}
				uids = append(uids, uid)
			}
			return uids, nil
		}
		uid, err := uuid.FromString(f.Value)
		if err != nil {
			return nil, &ConversionError{DataType: DataTypeUUID, Value: f.Value, Err: err}
		}
		return uid, nil
	case DataTypeList:
		err = json.Unmarshal([]byte(f.Value), &value)
		if err != nil {
			return nil, &ConversionError{DataType: DataTypeList, Value: f.Value, Err: err}
		}
		return value, nil
	default:
		return nil, &ValidationError{Kind: "data_type", Values: []string{f.Type}}
	}
}
