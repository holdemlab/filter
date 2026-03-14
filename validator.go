package filter

func validateOperator(operator string) error {
	switch operator {
	case OperatorEq, OperatorNotEq, OperatorLowerThan, OperatorLowerThanEq,
		OperatorGreaterThan, OperatorGreaterThanEq, OperatorBetween, OperatorLike, OperatorIn:
		return nil
	default:
		return &ValidationError{Kind: "operator", Values: []string{operator}}
	}
}

func validateDataType(dType string) error {
	switch dType {
	case DataTypeStr, DataTypeInt, DataTypeBool, DataTypeDate, DataTypeDateTime, DataTypeUUID, DataTypeList:
		return nil
	default:
		return &ValidationError{Kind: "data_type", Values: []string{dType}}
	}
}

// validateOperatorDataType перевіряє сумісність оператора з типом даних.
func validateOperatorDataType(operator, dType string) error {
	switch operator {
	case OperatorLike:
		if dType != DataTypeStr {
			return &ValidationError{Kind: "operator_data_type", Values: []string{operator, dType}}
		}
	case OperatorBetween:
		switch dType {
		case DataTypeDate, DataTypeDateTime, DataTypeInt:
			// ok
		default:
			return &ValidationError{Kind: "operator_data_type", Values: []string{operator, dType}}
		}
	case OperatorIn:
		switch dType {
		case DataTypeStr, DataTypeInt, DataTypeUUID, DataTypeList:
			// ok
		default:
			return &ValidationError{Kind: "operator_data_type", Values: []string{operator, dType}}
		}
	case OperatorLowerThan, OperatorLowerThanEq, OperatorGreaterThan, OperatorGreaterThanEq:
		switch dType {
		case DataTypeBool, DataTypeList:
			return &ValidationError{Kind: "operator_data_type", Values: []string{operator, dType}}
		}
	}
	return nil
}
