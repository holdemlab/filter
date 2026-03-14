package filter

import (
	"testing"
)

func TestValidateOperator_Valid(t *testing.T) {
	valid := []string{
		OperatorEq, OperatorNotEq, OperatorLowerThan, OperatorLowerThanEq,
		OperatorGreaterThan, OperatorGreaterThanEq, OperatorBetween, OperatorLike, OperatorIn,
	}
	for _, op := range valid {
		if err := validateOperator(op); err != nil {
			t.Errorf("expected operator %q to be valid, got error: %v", op, err)
		}
	}
}

func TestValidateOperator_Invalid(t *testing.T) {
	invalid := []string{"", "unknown", "EQ", "Eq", "!=", ">=", "BETWEEN"}
	for _, op := range invalid {
		if err := validateOperator(op); err == nil {
			t.Errorf("expected operator %q to be invalid, got nil", op)
		}
	}
}

func TestValidateDataType_Valid(t *testing.T) {
	valid := []string{
		DataTypeStr, DataTypeInt, DataTypeBool, DataTypeDate,
		DataTypeDateTime, DataTypeUUID, DataTypeList,
	}
	for _, dt := range valid {
		if err := validateDataType(dt); err != nil {
			t.Errorf("expected data type %q to be valid, got error: %v", dt, err)
		}
	}
}

func TestValidateDataType_Invalid(t *testing.T) {
	invalid := []string{"", "unknown", "String", "INT", "float", "array"}
	for _, dt := range invalid {
		if err := validateDataType(dt); err == nil {
			t.Errorf("expected data type %q to be invalid, got nil", dt)
		}
	}
}

func TestValidateOperatorDataType_Compatible(t *testing.T) {
	compatible := []struct {
		operator string
		dType    string
	}{
		// eq / neq працюють з усіма типами
		{OperatorEq, DataTypeStr},
		{OperatorEq, DataTypeInt},
		{OperatorEq, DataTypeBool},
		{OperatorEq, DataTypeDate},
		{OperatorEq, DataTypeDateTime},
		{OperatorEq, DataTypeUUID},
		{OperatorEq, DataTypeList},
		{OperatorNotEq, DataTypeStr},
		{OperatorNotEq, DataTypeInt},
		// like тільки з string
		{OperatorLike, DataTypeStr},
		// between з date, datetime, int
		{OperatorBetween, DataTypeDate},
		{OperatorBetween, DataTypeDateTime},
		{OperatorBetween, DataTypeInt},
		// in з string, int, uuid, list
		{OperatorIn, DataTypeStr},
		{OperatorIn, DataTypeInt},
		{OperatorIn, DataTypeUUID},
		{OperatorIn, DataTypeList},
		// comparison з числовими / датами / рядками
		{OperatorLowerThan, DataTypeInt},
		{OperatorLowerThan, DataTypeStr},
		{OperatorLowerThan, DataTypeDate},
		{OperatorLowerThanEq, DataTypeDateTime},
		{OperatorGreaterThan, DataTypeUUID},
		{OperatorGreaterThanEq, DataTypeInt},
	}
	for _, tc := range compatible {
		if err := validateOperatorDataType(tc.operator, tc.dType); err != nil {
			t.Errorf("expected (%q, %q) to be compatible, got error: %v", tc.operator, tc.dType, err)
		}
	}
}

func TestValidateOperatorDataType_Incompatible(t *testing.T) {
	incompatible := []struct {
		operator string
		dType    string
	}{
		// like не з string
		{OperatorLike, DataTypeInt},
		{OperatorLike, DataTypeBool},
		{OperatorLike, DataTypeDate},
		{OperatorLike, DataTypeUUID},
		{OperatorLike, DataTypeList},
		// between не з string, bool, uuid, list
		{OperatorBetween, DataTypeStr},
		{OperatorBetween, DataTypeBool},
		{OperatorBetween, DataTypeUUID},
		{OperatorBetween, DataTypeList},
		// in не з bool, date, datetime
		{OperatorIn, DataTypeBool},
		{OperatorIn, DataTypeDate},
		{OperatorIn, DataTypeDateTime},
		// comparison з bool, list
		{OperatorLowerThan, DataTypeBool},
		{OperatorLowerThan, DataTypeList},
		{OperatorLowerThanEq, DataTypeBool},
		{OperatorGreaterThan, DataTypeList},
		{OperatorGreaterThanEq, DataTypeBool},
	}
	for _, tc := range incompatible {
		if err := validateOperatorDataType(tc.operator, tc.dType); err == nil {
			t.Errorf("expected (%q, %q) to be incompatible, got nil", tc.operator, tc.dType)
		}
	}
}
