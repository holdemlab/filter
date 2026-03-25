package adapter

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/holdemlab/filter"
)

// GoquQuery builds a goqu.SelectDataset with filters, sorting and context applied.
func GoquQuery(ctx context.Context, op filter.Options, dataset *goqu.SelectDataset) (*goqu.SelectDataset, error) {
	var err error
	if op == nil {
		return nil, filter.ErrNilOptions
	}
	if dataset == nil {
		return nil, filter.ErrNilDataset
	}
	_ = ctx // reserved for future use (e.g. query execution with context)
	for _, f := range op.Fields() {
		dataset, err = goquQuery(dataset, f)
		if err != nil {
			return nil, fmt.Errorf("error in query: %w", err)
		}
	}

	if op.SortBy() != "" {
		if op.Desc() {
			dataset = dataset.Order(goqu.C(op.SortBy()).Desc())
		} else {
			dataset = dataset.Order(goqu.C(op.SortBy()).Asc())
		}
	}

	if op.Limit() > 0 {
		dataset = dataset.Limit(op.Limit())
	}

	if op.Page() > 0 && op.Limit() > 0 {
		dataset = dataset.Offset(uint(op.Page()-1) * op.Limit())
	}

	return dataset, nil
}

func goquQuery(dataset *goqu.SelectDataset, field filter.Field) (*goqu.SelectDataset, error) {
	value, err := field.GetValue()
	if err != nil {
		return nil, fmt.Errorf("error getting value for field %q: %w", field.Name, err)
	}
	switch field.Operator {
	case filter.OperatorLike:
		return dataset.Where(goqu.Ex{field.Name: goqu.Op{filter.OperatorLike: fmt.Sprintf("%%%v%%", value)}}), nil
	case filter.OperatorBetween:
		between, ok := value.(*filter.Between)
		if !ok {
			return nil, &filter.ConversionError{
				DataType: "Between",
				Value:    fmt.Sprintf("%T", value),
			}
		}
		return dataset.Where(goqu.Ex{field.Name: goqu.Op{filter.OperatorBetween: goqu.Range(between.Start(), between.End())}}), nil
	default:
		return dataset.Where(goqu.Ex{field.Name: goqu.Op{field.Operator: value}}), nil
	}
}
