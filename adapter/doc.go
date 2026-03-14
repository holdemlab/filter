// Package adapter provides query-building adapters that translate
// [filter.Options] into database-specific queries.
//
// Two adapters are available:
//
// # SQL via goqu
//
// [GoquQuery] takes an [filter.Options] and a *goqu.SelectDataset and returns
// a new dataset with WHERE clauses and ORDER BY applied:
//
//	dataset, err := adapter.GoquQuery(ctx, opts, goqu.From("users").Select("*"))
//
// # MongoDB
//
// [MongoQueryD] takes [filter.Options] and returns a bson.D filter together
// with *options.FindOptions that include limit, skip (pagination), and sort:
//
//	bsonFilter, findOpts, err := adapter.MongoQueryD(ctx, opts)
//	cursor, err := collection.Find(ctx, bsonFilter, findOpts)
//
// Both adapters accept a [context.Context] as the first argument.
package adapter
