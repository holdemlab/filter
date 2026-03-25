package adapter

import (
	"context"
	"fmt"

	"github.com/holdemlab/filter"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoQueryD builds a bson.D filter and FindOptions from filter.Options.
// The ctx parameter should be used by the caller when executing the query.
func MongoQueryD(ctx context.Context, op filter.Options) (bson.D, *options.FindOptions, error) {
	if op == nil {
		return nil, nil, filter.ErrNilOptions
	}
	_ = ctx // reserved for future use (e.g. query execution with context)
	dataset := bson.D{}
	for _, f := range op.Fields() {
		set, err := mongoQueryD(f)
		if err != nil {
			return nil, nil, fmt.Errorf("error in query: %w", err)
		}
		dataset = append(dataset, set...)
	}
	opts := options.Find()
	if op.SortBy() != "" {
		if op.Desc() {
			opts.SetSort(bson.D{bson.E{Key: op.SortBy(), Value: -1}})
		} else {
			opts.SetSort(bson.D{bson.E{Key: op.SortBy(), Value: 1}})
		}
	}
	if op.Limit() > 0 {
		opts.SetLimit(int64(op.Limit()))
	}
	if op.Page() > 0 && op.Limit() > 0 {
		opts.SetSkip(int64(op.Page()-1) * int64(op.Limit()))
	}
	return dataset, opts, nil
}

func mongoQueryD(field filter.Field) (bson.D, error) {
	value, err := field.GetValue()
	if err != nil {
		return nil, fmt.Errorf("error getting value for field %q: %w", field.Name, err)
	}
	switch field.Operator {
	case filter.OperatorBetween:
		between := value.(*filter.Between)
		return bson.D{bson.E{Key: field.Name, Value: bson.D{
			bson.E{Key: "$gte", Value: between.Start()},
			bson.E{Key: "$lte", Value: between.End()},
		}}}, nil
	case filter.OperatorGreaterThan:
		return bson.D{bson.E{Key: field.Name, Value: bson.D{bson.E{Key: "$gt", Value: value}}}}, nil
	case filter.OperatorGreaterThanEq:
		return bson.D{bson.E{Key: field.Name, Value: bson.D{bson.E{Key: "$gte", Value: value}}}}, nil
	case filter.OperatorLowerThan:
		return bson.D{bson.E{Key: field.Name, Value: bson.D{bson.E{Key: "$lt", Value: value}}}}, nil
	case filter.OperatorLowerThanEq:
		return bson.D{bson.E{Key: field.Name, Value: bson.D{bson.E{Key: "$lte", Value: value}}}}, nil
	case filter.OperatorEq:
		return bson.D{bson.E{Key: field.Name, Value: bson.D{bson.E{Key: "$eq", Value: value}}}}, nil
	case filter.OperatorNotEq:
		return bson.D{bson.E{Key: field.Name, Value: bson.D{bson.E{Key: "$ne", Value: value}}}}, nil
	case filter.OperatorIn:
		return bson.D{bson.E{Key: field.Name, Value: bson.D{bson.E{Key: "$in", Value: value}}}}, nil
	case filter.OperatorLike:
		return bson.D{bson.E{Key: field.Name, Value: bson.D{
			bson.E{Key: "$regex", Value: fmt.Sprintf("%v", value)},
			bson.E{Key: "$options", Value: "i"},
		}}}, nil
	default:
		return bson.D{bson.E{Key: field.Name, Value: value}}, nil
	}
}
