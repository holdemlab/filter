package filter

import (
	"fmt"
	"reflect"
	"strings"
)

// Options describes a set of filtering, pagination, and sorting criteria.
type (
	Options interface {
		// Limit returns the maximum number of results per page.
		Limit() uint
		// Page returns the 1-based page number.
		Page() int
		// SortBy returns the field name to sort by.
		SortBy() string
		// Desc reports whether the sort order is descending.
		Desc() bool
		// AddField appends a validated filter field.
		// It returns a [*ValidationError] if the operator, data type, or their
		// combination is invalid.
		AddField(name, operator, value, dType string) error
		// Fields returns all filter fields.
		Fields() []Field
		// FieldByName returns the first field with the given name and true,
		// or a zero [Field] and false if not found.
		FieldByName(name string) (Field, bool)
		// UpdateFromQueryParams populates filter fields from the exported
		// string fields of opt, using the specified struct tag as the field
		// name mapping. The "type" tag must also be present on each field.
		UpdateFromQueryParams(opt any, tag string) error
	}
	options struct {
		limit      int
		page       int
		sortBy     string
		descending bool
		fields     []Field
	}
)

// NewOptions creates a new [Options] with the given pagination and sorting
// parameters. Values of limit and page that are less than 1 are clamped to 1.
func NewOptions(limit int, page int, sortBy string, desc bool) Options {
	if limit <= 0 {
		limit = 1
	}
	if page <= 0 {
		page = 1
	}
	return &options{limit: limit, page: page, sortBy: sortBy, descending: desc}
}

func (o *options) Limit() uint {
	return uint(o.limit)
}

func (o *options) Page() int {
	return o.page
}

func (o *options) SortBy() string {
	return o.sortBy
}

func (o *options) Desc() bool {
	return o.descending
}

func (o *options) Fields() []Field {
	return o.fields
}

func (o *options) FieldByName(name string) (Field, bool) {
	for _, f := range o.fields {
		if f.Name == name {
			return f, true
		}
	}
	return Field{}, false
}

func (o *options) AddField(name, operator, value, dType string) error {
	if err := validateOperator(operator); err != nil {
		return fmt.Errorf("invalid operator: %w", err)
	}
	if err := validateDataType(dType); err != nil {
		return fmt.Errorf("invalid data type: %w", err)
	}
	if err := validateOperatorDataType(operator, dType); err != nil {
		return fmt.Errorf("incompatible operator and data type: %w", err)
	}
	o.fields = append(o.fields, Field{
		Name:     name,
		Operator: operator,
		Value:    value,
		Type:     dType,
	})
	return nil
}

func (o *options) UpdateFromQueryParams(opt any, tag string) error {
	val := reflect.ValueOf(opt)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("%w: got %T", ErrNotStruct, opt)
	}

	valType := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := valType.Field(i)
		if valTag := field.Tag.Get(tag); valTag != "" {
			if typTag := field.Tag.Get("type"); typTag != "" {
				operator := OperatorEq
				value, ok := val.Field(i).Interface().(string)
				if !ok {
					return &ConversionError{
						DataType: "string",
						Value:    fmt.Sprintf("%v", val.Field(i).Interface()),
					}
				}
				if value == "" {
					continue
				}
				if strings.Contains(value, "#") {
					split := strings.Split(value, "#")
					operator = split[0]
					value = split[1]
				}

				if err := o.AddField(valTag, operator, value, typTag); err != nil {
					return err
				}
			} else {
				return ErrEmptyTypeTag
			}

		}
	}
	return nil
}
