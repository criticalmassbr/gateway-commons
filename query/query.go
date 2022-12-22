package query

import (
	"errors"
	"strconv"
	"strings"
	"unicode"

	"github.com/gofiber/fiber/v2"
)

type Paginable struct {
	Limit  int `json:"limit"`  // Maximun amount of records that should be fetched
	Offset int `json:"offset"` // Index to fetch records after
}

type Filter struct {
	Field     string         `json:"field"`     // the field to filter by eg: "price"
	Value     string         `json:"value"`     // string representation of the value to apply the filter, eg: "10", "tag1;tag2"
	Operation FilterOperator `json:"operation"` // the operation to use for filtering, eg: gt (greather than)
}

type FilterOperator string

const (
	FilterOperatorEqual    FilterOperator = "eq"
	FilterOperatorContains FilterOperator = "contains"
)

func (f *FilterOperator) IsValid() bool {
	switch *f {
	case FilterOperatorEqual, FilterOperatorContains:
		return true
	}

	return false
}

type Order struct {
	Field string `json:"field"` // the field to sort by eg: "price"
	Asc   bool   `json:"asc"`   // if true, sort on ascending order, else descending
}

type QueryKey string

const (
	QueryKeyLimit   QueryKey = "limit"
	QueryKeyOffset  QueryKey = "offset"
	QueryKeySearch  QueryKey = "search"
	QueryKeyOrder   QueryKey = "order"
	QueryKeyFilters QueryKey = "filters"
)

func hasALetter(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) {
			return true
		}
	}
	return false
}

func getPositiveIntFromQueryWithFallback(ctx *fiber.Ctx, key string, fallbackVal int) int {
	val := ctx.Query(key)
	if val == "" {
		return fallbackVal
	}

	i, err := strconv.Atoi(val)
	if err != nil {
		return fallbackVal
	}

	if i < 0 {
		return fallbackVal
	}

	return i
}

// eg: "ilove[brazil]" -> "ilove", "brazil"
func getStrBeforeAndInbetweenBrackets(k string) (string, FilterOperator) {
	s := strings.Index(k, "[")
	if s == -1 {
		return "", ""
	}

	e := strings.Index(k, "]")
	if e == -1 {
		return "", ""
	}

	return k[:s], FilterOperator(k[s+1 : e])
}

type QueryParamSeparator string

const (
	QueryParamSeparatorArray QueryParamSeparator = ";"
	QueryParamSeparatorMap   QueryParamSeparator = ","
	QueryParamSeparatorValue QueryParamSeparator = ":"
)

func splitStringBySeparator(s string, sep QueryParamSeparator) []string {
	return strings.Split(s, string(sep))
}

func GetPaginationFromQuery(ctx *fiber.Ctx) Paginable {
	return Paginable{
		Limit:  getPositiveIntFromQueryWithFallback(ctx, string(QueryKeyLimit), 10),
		Offset: getPositiveIntFromQueryWithFallback(ctx, string(QueryKeyOffset), 0),
	}
}

func GetFilterFromQuery(c *fiber.Ctx) []Filter {
	f := []Filter{}

	c.Context().QueryArgs().VisitAll(func(key, val []byte) {
		k := string(key)
		v := string(val)

		field, op := getStrBeforeAndInbetweenBrackets(k)

		if !hasALetter(string(op)) {
			return
		}

		f = append(f, Filter{Field: field, Operation: op, Value: v})
	})

	return f
}

func getOrderFromQuery(queryParams map[string]string) []Order {
	if value, ok := queryParams[string(QueryKeyOrder)]; ok {
		return getOrderFields(value)
	}

	return []Order{}
}

func getOrderFields(value string) []Order {
	order := []Order{}

	sorts := splitStringBySeparator(value, QueryParamSeparatorMap)

	for _, sort := range sorts {
		s := splitStringBySeparator(sort, QueryParamSeparatorValue)
		if len(s) != 2 {
			continue
		}

		o, err := getOrder(s[0], s[1])
		if err != nil {
			continue
		}

		order = append(order, o)
	}

	return order
}

func getOrder(field, value string) (Order, error) {
	if !hasALetter(field) {
		return Order{}, errors.New("invalid order field")
	}

	return Order{
		Field: field,
		Asc:   strings.ToLower(value) == "asc",
	}, nil
}

func GetOrderFromQuery(c *fiber.Ctx) []Order {
	queryParams := queryParamsToMap(c)
	order := getOrderFromQuery(queryParams)
	return order
}

func queryParamsToMap(c *fiber.Ctx) map[string]string {
	queryParams := make(map[string]string)

	c.Context().QueryArgs().VisitAll(func(key []byte, value []byte) {
		queryParams[string(key)] = string(value)
	})

	return queryParams
}
