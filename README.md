# filter

[![Go Reference](https://pkg.go.dev/badge/github.com/holdemlab/filter.svg)](https://pkg.go.dev/github.com/holdemlab/filter)

Package `filter` provides a universal filtering, pagination, and sorting mechanism for HTTP requests. It supports query building for both SQL (via goqu) and MongoDB.

## Installation

```bash
go get github.com/holdemlab/filter
```

## Package Structure

```
filter/
├── doc.go              # Package-level documentation (godoc)
├── options.go          # Options interface and implementation
├── field.go            # Field struct and value conversion
├── constants.go        # Operator and data type constants
├── errors.go           # Typed errors (sentinel + custom types)
├── validator.go        # Operator and data type validation
├── between_parser.go   # Between operator parser for dates
├── middlewares.go      # Gin middlewares for automatic query parameter parsing
└── adapter/
    ├── doc.go          # Package-level documentation (godoc)
    ├── goqu.go         # Adapter for building SQL queries (goqu)
    └── mongo.go        # Adapter for building MongoDB queries
```

## Core Components

### Options

The `Options` interface is the central component of the package. It holds pagination, sorting parameters, and a set of filter fields.

```go
type Options interface {
    Limit() uint
    Page() int
    SortBy() string
    Desc() bool
    AddField(name, operator, value, dType string) error
    Fields() []Field
    FieldByName(name string) (Field, bool)
    UpdateFromQueryParams(opt any, tag string) error
}
```

Creating an instance:

```go
opts := filter.NewOptions(10, 1, "created_at", true)
// limit=10, page=1, sort by created_at, descending=true
```

### Field

A struct describing a single filter:

```go
type Field struct {
    Name     string // Field name (column / document key)
    Operator string // Comparison operator (eq, neq, lt, gte, between, like, in, ...)
    Value    string // Value as a string
    Type     string // Data type (string, int, bool, date, datetime, uuid, list)
}
```

The `GetValue()` method converts the string `Value` into a typed value according to the `Type` field:

| Type       | Go Result Type             | Notes                                                           |
|------------|----------------------------|-----------------------------------------------------------------|
| `string`   | `string`                   |                                                                 |
| `int`      | `int`                      |                                                                 |
| `bool`     | `bool`                     |                                                                 |
| `date`     | `time.Time` / `*Between`   | Format `2006-01-02`. For `between` — two values separated by `*` |
| `datetime` | `time.Time` / `*Between`   | Format `2006-01-02T15:04:05`. For `between` — separated by `*` |
| `uuid`     | `uuid.UUID`                |                                                                 |
| `list`     | `any` (JSON unmarshal)     | JSON array, e.g. `["a","b"]`                                    |

Dates can also be passed as Unix timestamps (integers).

### Constants

#### Comparison Operators

| Constant             | Value       | Description                   |
|----------------------|-------------|-------------------------------|
| `OperatorEq`         | `eq`        | Equal                         |
| `OperatorNotEq`      | `neq`       | Not equal                     |
| `OperatorLowerThan`  | `lt`        | Less than                     |
| `OperatorLowerThanEq`| `lte`       | Less than or equal            |
| `OperatorGreaterThan`| `gt`        | Greater than                  |
| `OperatorGreaterThanEq`| `gte`     | Greater than or equal         |
| `OperatorBetween`    | `between`   | Between two values            |
| `OperatorLike`       | `like`      | Substring search (SQL LIKE)   |
| `OperatorIn`         | `in`        | Contained in a list           |

#### Data Types

`string`, `int`, `bool`, `date`, `datetime`, `uuid`, `list`

## Gin Middlewares

The package provides three Gin middlewares for automatic parsing of pagination and sorting query parameters.

### QueryOptionsMiddlewares

Stores `Options` in `gin.Context` under the key `filter_options`. You can retrieve it in a handler:

```go
router.GET("/items", filter.QueryOptionsMiddlewares(), func(c *gin.Context) {
    opts := c.MustGet(filter.OptionsContextKey).(filter.Options)
    // ...
})
```

**Query parameters:** `sort_by`, `descending`, `page`, `limit` (default limit is 10).

### SingleQueryOptionsMiddlewares

Passes `Options` directly to a callback function. Uses default values: `limit=10`, `page=1`, `sort_by=id`.

```go
router.GET("/items", filter.SingleQueryOptionsMiddlewares(func(c *gin.Context, opts filter.Options) {
    // use opts
}))
```

### SingleQueryOptionsMiddlewaresWithDefaults

Similar to `SingleQueryOptionsMiddlewares`, but allows passing custom default values via `*Params`:

```go
defaults := &filter.Params{
    Limit:      20,
    Page:       1,
    SortBy:     "created_at",
    Descending: true,
}
router.GET("/items", filter.SingleQueryOptionsMiddlewaresWithDefaults(handler, defaults))
```

## UpdateFromQueryParams

A method for automatically building filters from a struct with tags. Tags define the field name and its type:

```go
type MyFilter struct {
    Status string `form:"status" type:"string"`
    Amount string `form:"amount" type:"int"`
}

var f MyFilter
c.ShouldBind(&f)
opts.UpdateFromQueryParams(f, "form")
```

To specify an operator, use the value format `operator#value`. For example, the query parameter `?amount=gte#100` creates a filter with operator `gte` and value `100`. If `#` is absent, the default operator `eq` is used.

## Between Operator

For the `between` operator, values are separated by `*`:

```
?date=between#2024-01-01*2024-12-31
```

The `ParseBetweenOperator` parser splits the value into `start` and `end` and returns a `*Between` struct with `Start()` and `End()` methods.

## Adapters

### adapter.GoquQuery (SQL / goqu)

Accepts `context.Context`, `Options`, and `*goqu.SelectDataset`, adds WHERE conditions for each `Field`:

```go
dataset := goqu.From("users").Select("*")
dataset, err := adapter.GoquQuery(ctx, opts, dataset)
```

Operator mapping:
- `like` — `WHERE field LIKE '%value%'`
- `between` — `WHERE field BETWEEN start AND end`
- Others (`eq`, `neq`, `lt`, `lte`, `gt`, `gte`, `in`) — standard goqu `Op`

### adapter.MongoQueryD (MongoDB)

Accepts `context.Context` and `Options`, returns a `bson.D` filter and `*options.FindOptions` with pagination and sorting:

```go
bsonFilter, findOpts, err := adapter.MongoQueryD(ctx, opts)
cursor, err := collection.Find(ctx, bsonFilter, findOpts)
```

Operator mapping:
- `eq` → `$eq`, `neq` → `$ne`, `gt` → `$gt`, `gte` → `$gte`, `lt` → `$lt`, `lte` → `$lte`
- `in` → `$in`
- `between` → `$gte` + `$lte`

Pagination and sorting are configured via `FindOptions` (`SetLimit`, `SetSkip`, `SetSort`).

## Typed Errors

The package provides typed errors that support `errors.Is` and `errors.As`:

| Type | Description |
|------|-------------|
| `ErrNilOptions` | Options is nil |
| `ErrNilDataset` | Dataset is nil |
| `ErrEmptyTypeTag` | Missing `type` tag in struct |
| `ErrNotStruct` | Non-struct value provided |
| `*ValidationError` | Invalid operator, data type, or their combination |
| `*ConversionError` | Value conversion error |
| `*ParseError` | Between expression parsing error |

```go
var ve *filter.ValidationError
if errors.As(err, &ve) {
    fmt.Println(ve.Kind, ve.Values)
}
```

## Validation

When adding a field via `AddField`, the following are validated automatically:
- **Operator** — only defined constants are allowed (`eq`, `neq`, `lt`, `lte`, `gt`, `gte`, `between`, `like`, `in`)
- **Data type** — only defined constants are allowed (`string`, `int`, `bool`, `date`, `datetime`, `uuid`, `list`)
- **Compatibility** — the operator+type combination is checked (e.g., `like` only with `string`, `between` only with `date`/`datetime`/`int`)

An error is returned if the value is invalid.

## Full Usage Example

```go
// Handler
func ListUsers(c *gin.Context, opts filter.Options) {
    // Additional filters from query parameters
    type UserFilter struct {
        Name   string `form:"name"   type:"string"`
        Age    string `form:"age"    type:"int"`
        Active string `form:"active" type:"bool"`
    }

    var f UserFilter
    c.ShouldBind(&f)
    opts.UpdateFromQueryParams(f, "form")

    // SQL query via goqu
    dataset := goqu.From("users").Select("*")
    dataset, err := adapter.GoquQuery(c.Request.Context(), opts, dataset)
    if err != nil {
        // handle error
    }

    // Pagination
    offset := uint(opts.Page()) * opts.Limit()
    dataset = dataset.Limit(opts.Limit()).Offset(offset)

    // Sorting
    if opts.SortBy() != "" {
        if opts.Desc() {
            dataset = dataset.Order(goqu.C(opts.SortBy()).Desc())
        } else {
            dataset = dataset.Order(goqu.C(opts.SortBy()).Asc())
        }
    }
}

// Router
router.GET("/users", filter.SingleQueryOptionsMiddlewares(ListUsers))
```

**Request:**
```
GET /users?page=1&limit=20&sort_by=name&descending=true&name=like#John&age=gte#18
```

## Documentation

Full API documentation is available at [pkg.go.dev](https://pkg.go.dev/github.com/holdemlab/filter).
