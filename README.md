# filter

[![Go Reference](https://pkg.go.dev/badge/github.com/holdemlab/filter.svg)](https://pkg.go.dev/github.com/holdemlab/filter)

Пакет `filter` реалізує універсальний механізм фільтрації, пагінації та сортування для HTTP-запитів. Підтримує побудову запитів як до SQL (через goqu), так і до MongoDB.

## Встановлення

```bash
go get github.com/holdemlab/filter
```

## Структура пакету

```
filter/
├── doc.go              # Package-level documentation (godoc)
├── options.go          # Options інтерфейс та його реалізація
├── field.go            # Field структура та конвертація значень
├── constants.go        # Константи операторів та типів даних
├── errors.go           # Типізовані помилки (sentinel + custom types)
├── validator.go        # Валідація операторів і типів
├── between_parser.go   # Парсер оператора between для дат
├── middlewares.go      # Gin middlewares для автоматичного парсингу query-параметрів
└── adapter/
    ├── doc.go          # Package-level documentation (godoc)
    ├── goqu.go         # Адаптер для побудови SQL-запитів (goqu)
    └── mongo.go        # Адаптер для побудови MongoDB-запитів
```

## Основні компоненти

### Options

Інтерфейс `Options` — центральна точка пакету. Містить параметри пагінації, сортування та набір полів-фільтрів.

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

Створення:

```go
opts := filter.NewOptions(10, 1, "created_at", true)
// limit=10, page=1, сортування по created_at, descending=true
```

### Field

Структура, що описує один фільтр:

```go
type Field struct {
    Name     string // Ім'я поля (колонка / ключ документа)
    Operator string // Оператор порівняння (eq, neq, lt, gte, between, like, in, ...)
    Value    string // Значення у вигляді рядка
    Type     string // Тип даних (string, int, bool, date, datetime, uuid, list)
}
```

Метод `GetValue()` конвертує рядкове `Value` у типізоване значення відповідно до поля `Type`:

| Type       | Go-тип результату          | Примітки                                                        |
|------------|----------------------------|-----------------------------------------------------------------|
| `string`   | `string`                   |                                                                 |
| `int`      | `int`                      |                                                                 |
| `bool`     | `bool`                     |                                                                 |
| `date`     | `time.Time` / `*Between`   | Формат `2006-01-02`. Для `between` — два значення через `*`    |
| `datetime` | `time.Time` / `*Between`   | Формат `2006-01-02T15:04:05`. Для `between` — через `*`       |
| `uuid`     | `uuid.UUID`                |                                                                 |
| `list`     | `any` (JSON unmarshal)     | JSON-масив, наприклад `["a","b"]`                               |

Дати також можуть бути передані як Unix timestamp (ціле число).

### Константи

#### Оператори порівняння

| Константа            | Значення    | Опис                          |
|----------------------|-------------|-------------------------------|
| `OperatorEq`         | `eq`        | Дорівнює                      |
| `OperatorNotEq`      | `neq`       | Не дорівнює                   |
| `OperatorLowerThan`  | `lt`        | Менше                         |
| `OperatorLowerThanEq`| `lte`       | Менше або дорівнює            |
| `OperatorGreaterThan`| `gt`        | Більше                        |
| `OperatorGreaterThanEq`| `gte`     | Більше або дорівнює           |
| `OperatorBetween`    | `between`   | Між двома значеннями          |
| `OperatorLike`       | `like`      | Пошук за підрядком (SQL LIKE) |
| `OperatorIn`         | `in`        | Входить у список              |

#### Типи даних

`string`, `int`, `bool`, `date`, `datetime`, `uuid`, `list`

## Gin Middlewares

Пакет надає три gin middleware для автоматичного парсингу query-параметрів пагінації та сортування.

### QueryOptionsMiddlewares

Зберігає `Options` у `gin.Context` під ключем `filter_options`. Далі в хендлері їх можна дістати:

```go
router.GET("/items", filter.QueryOptionsMiddlewares(), func(c *gin.Context) {
    opts := c.MustGet(filter.OptionsContextKey).(filter.Options)
    // ...
})
```

**Query-параметри:** `sort_by`, `descending`, `page`, `limit` (ліміт за замовчуванням — 10).

### SingleQueryOptionsMiddlewares

Передає `Options` напряму в callback-функцію. Встановлює значення за замовчуванням: `limit=10`, `page=1`, `sort_by=id`.

```go
router.GET("/items", filter.SingleQueryOptionsMiddlewares(func(c *gin.Context, opts filter.Options) {
    // використання opts
}))
```

### SingleQueryOptionsMiddlewaresWithDefaults

Аналогічно `SingleQueryOptionsMiddlewares`, але дозволяє передати власні значення за замовчуванням через `*Params`:

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

Метод для автоматичної побудови фільтрів із структури з тегами. Теги визначають ім'я поля та його тип:

```go
type MyFilter struct {
    Status string `form:"status" type:"string"`
    Amount string `form:"amount" type:"int"`
}

var f MyFilter
c.ShouldBind(&f)
opts.UpdateFromQueryParams(f, "form")
```

Для вказання оператора використовується формат значення `operator#value`. Наприклад, query-параметр `?amount=gte#100` створить фільтр з оператором `gte` та значенням `100`. Якщо `#` відсутній — використовується оператор `eq` за замовчуванням.

## Between оператор

Для оператора `between` значення передаються через розділювач `*`:

```
?date=between#2024-01-01*2024-12-31
```

Парсер `ParseBetweenOperator` розбиває значення на `start` та `end` і повертає структуру `*Between` з методами `Start()` та `End()`.

## Адаптери

### adapter.GoquQuery (SQL / goqu)

Приймає `context.Context`, `Options` та `*goqu.SelectDataset`, додає WHERE-умови для кожного `Field`:

```go
dataset := goqu.From("users").Select("*")
dataset, err := adapter.GoquQuery(ctx, opts, dataset)
```

Маппінг операторів:
- `like` — `WHERE field LIKE '%value%'`
- `between` — `WHERE field BETWEEN start AND end`
- Решта (`eq`, `neq`, `lt`, `lte`, `gt`, `gte`, `in`) — стандартний goqu `Op`

### adapter.MongoQueryD (MongoDB)

Приймає `context.Context` та `Options`, повертає `bson.D` фільтр та `*options.FindOptions` із пагінацією та сортуванням:

```go
bsonFilter, findOpts, err := adapter.MongoQueryD(ctx, opts)
cursor, err := collection.Find(ctx, bsonFilter, findOpts)
```

Маппінг операторів:
- `eq` → `$eq`, `neq` → `$ne`, `gt` → `$gt`, `gte` → `$gte`, `lt` → `$lt`, `lte` → `$lte`
- `in` → `$in`
- `between` → `$gte` + `$lte`

Пагінація та сортування налаштовуються через `FindOptions` (`SetLimit`, `SetSkip`, `SetSort`).

## Помилки (Typed Errors)

Пакет надає типізовані помилки, які підтримують `errors.Is` та `errors.As`:

| Тип | Опис |
|------|------|
| `ErrNilOptions` | Options є nil |
| `ErrNilDataset` | Dataset є nil |
| `ErrEmptyTypeTag` | Відсутній тег `type` у struct |
| `ErrNotStruct` | Передано не-struct |
| `*ValidationError` | Невалідний оператор, тип даних або їх комбінація |
| `*ConversionError` | Помилка конвертації значення |
| `*ParseError` | Помилка парсингу between-виразу |

```go
var ve *filter.ValidationError
if errors.As(err, &ve) {
    fmt.Println(ve.Kind, ve.Values)
}
```

## Валідація

При додаванні поля через `AddField` автоматично валідуються:
- **Оператор** — дозволені тільки визначені константи (`eq`, `neq`, `lt`, `lte`, `gt`, `gte`, `between`, `like`, `in`)
- **Тип даних** — дозволені тільки визначені константи (`string`, `int`, `bool`, `date`, `datetime`, `uuid`, `list`)
- **Сумісність** — перевіряється комбінація оператор+тип (наприклад, `like` тільки з `string`, `between` тільки з `date`/`datetime`/`int`)

У разі невалідного значення повертається помилка.

## Приклад повного використання

```go
// Handler
func ListUsers(c *gin.Context, opts filter.Options) {
    // Додаткові фільтри з query-параметрів
    type UserFilter struct {
        Name   string `form:"name"   type:"string"`
        Age    string `form:"age"    type:"int"`
        Active string `form:"active" type:"bool"`
    }

    var f UserFilter
    c.ShouldBind(&f)
    opts.UpdateFromQueryParams(f, "form")

    // SQL-запит через goqu
    dataset := goqu.From("users").Select("*")
    dataset, err := adapter.GoquQuery(c.Request.Context(), opts, dataset)
    if err != nil {
        // handle error
    }

    // Пагінація
    offset := uint(opts.Page()) * opts.Limit()
    dataset = dataset.Limit(opts.Limit()).Offset(offset)

    // Сортування
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

**Запит:**
```
GET /users?page=1&limit=20&sort_by=name&descending=true&name=like#John&age=gte#18
```

## Документація

Повна API-документація доступна на [pkg.go.dev](https://pkg.go.dev/github.com/holdemlab/filter).
