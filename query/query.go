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
	FilterOperatorEqual       FilterOperator = "eq"
	FilterOperatorContains    FilterOperator = "contains"
	FilterOperatorGreaterThan FilterOperator = "gt"
)

func (f *FilterOperator) IsValid() bool {
	switch *f {
	case FilterOperatorEqual, FilterOperatorContains, FilterOperatorGreaterThan:
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

// eg: "i[love]brazil" -> Filter{Field: "i", Operation: FilterOperator("love"), Value: "brazil"}
func getFilter(value string) (Filter, error) {
	opStart := strings.Index(value, string(QueryParamSeparatorOperatorStart))
	if opStart == -1 {
		return Filter{}, errors.New("invalid filter")
	}

	opEnd := strings.Index(value, string(QueryParamSeparatorOperatorEnd))
	if opEnd == -1 {
		return Filter{}, errors.New("invalid filter")
	}

	o := FilterOperator(value[opStart+1 : opEnd])
	if !o.IsValid() {
		return Filter{}, errors.New("invalid filter operator")
	}

	f := value[:opStart]
	if !hasALetter(f) {
		return Filter{}, errors.New("invalid filter field")
	}

	v := value[opEnd+1:]
	if v == "" {
		return Filter{}, errors.New("invalid filter value")
	}

	return Filter{
		Field:     f,
		Operation: o,
		Value:     v,
	}, nil
}

type QueryParamSeparator string

const (
	QueryParamSeparatorArray         QueryParamSeparator = ";"
	QueryParamSeparatorMap           QueryParamSeparator = ","
	QueryParamSeparatorValue         QueryParamSeparator = ":"
	QueryParamSeparatorOperatorStart QueryParamSeparator = "["
	QueryParamSeparatorOperatorEnd   QueryParamSeparator = "]"
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

func getFilterFromQuery(queryParams map[string]string) []Filter {
	if value, ok := queryParams[string(QueryKeyFilters)]; ok {
		return getFilterFields(value)
	}

	return []Filter{}
}

func getFilterFields(value string) []Filter {
	filters := []Filter{}

	fields := splitStringBySeparator(value, QueryParamSeparatorMap)

	for _, field := range fields {
		f, err := getFilter(field)
		if err != nil {
			continue
		}

		filters = append(filters, f)
	}

	return filters
}

func GetFilterFromQuery(c *fiber.Ctx) []Filter {
	queryParams := queryParamsToMap(c)
	filter := getFilterFromQuery(queryParams)
	return filter
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
