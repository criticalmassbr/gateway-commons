package query

import (
	"sort"
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
	Field     string `json:"field"`     // the field to filter by eg: "price"
	Value     string `json:"value"`     // string representation of the value to apply the filter, eg: "10"
	Operation string `json:"operation"` // the operation to use for filtering, eg: gt (greather than)
}

type Order struct {
	Field string `json:"field"` // the field to sort by eg: "price"
	Asc   bool   `json:"asc"`   // if true, sort on ascending order, else descending
	idx   int
}

func hasOneLetter(s string) bool {
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
func getStrBeforeAndInbetweenBrackets(k string) (string, string) {
	s := strings.Index(k, "[")
	if s == -1 {
		return "", ""
	}

	e := strings.Index(k, "]")
	if e == -1 {
		return "", ""
	}

	return k[:s], k[s+1 : e]
}

func GetPaginationFromQuery(ctx *fiber.Ctx) Paginable {
	return Paginable{
		Limit:  getPositiveIntFromQueryWithFallback(ctx, "limit", 10),
		Offset: getPositiveIntFromQueryWithFallback(ctx, "offset", 0),
	}
}

func GetFilterFromQuery(c *fiber.Ctx) []Filter {
	f := []Filter{}

	c.Context().QueryArgs().VisitAll(func(key, val []byte) {
		k := string(key)
		v := string(val)

		field, op := getStrBeforeAndInbetweenBrackets(k)

		if !hasOneLetter(op) {
			return
		}

		f = append(f, Filter{Field: field, Operation: op, Value: v})
	})

	return f
}

func GetOrderFromQuery(c *fiber.Ctx) []Order {
	s := []Order{}

	c.Context().QueryArgs().VisitAll(func(key, val []byte) {
		k := string(key)
		v := string(val)

		bef, in := getStrBeforeAndInbetweenBrackets(k)

		if bef != "sort" {
			return
		}

		idx, err := strconv.Atoi(in)
		if err != nil {
			return
		}

		vSpaceIdx := strings.Index(v, " ")

		if vSpaceIdx == -1 {
			s = append(s, Order{Field: v, Asc: false, idx: idx})
			return
		}

		field := v[:vSpaceIdx]
		ascOrDescStr := v[vSpaceIdx+1:]

		if ascOrDescStr == "desc" || ascOrDescStr == "DESC" {
			s = append(s, Order{Field: field, Asc: false, idx: idx})
		} else {
			s = append(s, Order{Field: field, Asc: true, idx: idx})
		}
	})

	sort.SliceStable(s, func(i, j int) bool { return s[i].idx < s[j].idx })

	return s
}
