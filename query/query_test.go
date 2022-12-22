package query

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOrderFromQuery(t *testing.T) {
	type args struct {
		queryParams map[string]string
	}
	tests := []struct {
		name string
		args args
		want []Order
	}{
		{
			name: "should return empty slice when query params is empty",
			args: args{queryParams: map[string]string{}},
			want: []Order{},
		},
		{
			name: "should return empty slice when query params has no order",
			args: args{queryParams: map[string]string{
				"limit":  "10",
				"offset": "0",
			}},
			want: []Order{},
		},
		{
			name: "should return empty slice when query params has empty order",
			args: args{queryParams: map[string]string{
				"limit":  "10",
				"offset": "0",
				"order":  "",
			}},
			want: []Order{},
		},
		{
			name: "should return empty slice when query params has invalid order",
			args: args{queryParams: map[string]string{
				"limit":  "10",
				"offset": "0",
				"order":  "invalid",
			}},
			want: []Order{},
		},
		{
			name: "should return empty slice when query params has invalid order",
			args: args{queryParams: map[string]string{
				"limit":  "10",
				"offset": "0",
				"order":  "invalid,invalid",
			}},
			want: []Order{},
		},
		{
			name: "should return Order slice when query params has valid order",
			args: args{queryParams: map[string]string{
				"limit":  "10",
				"offset": "0",
				"order":  "id:asc",
			}},
			want: []Order{
				{
					Field: "id",
					Asc:   true,
				},
			},
		},
		{
			name: "should return Order slice when query params has valid order",
			args: args{queryParams: map[string]string{
				"limit":  "10",
				"offset": "0",
				"order":  "id:asc,name:desc",
			}},
			want: []Order{
				{
					Field: "id",
					Asc:   true,
				},
				{
					Field: "name",
					Asc:   false,
				},
			},
		},
		{
			name: "should return Order slice when query params has valid order",
			args: args{queryParams: map[string]string{
				"limit":  "10",
				"offset": "0",
				"order":  "id:asc,name:desc,age:asc",
			}},
			want: []Order{
				{
					Field: "id",
					Asc:   true,
				},
				{
					Field: "name",
					Asc:   false,
				},
				{
					Field: "age",
					Asc:   true,
				},
			},
		},
		{
			name: "should return Order slice when query params has valid order",
			args: args{queryParams: map[string]string{
				"limit":  "10",
				"offset": "0",
				"order":  "a:.,b:*(),c:???",
			}},
			want: []Order{
				{
					Field: "a",
					Asc:   false,
				},
				{
					Field: "b",
					Asc:   false,
				},
				{
					Field: "c",
					Asc:   false,
				},
			},
		},
		{
			name: "should return Order slice with omitted values when query params has invalid order",
			args: args{queryParams: map[string]string{
				"limit":  "10",
				"offset": "0",
				"order":  "a:asc,.:desc,c:asc",
			}},
			want: []Order{
				{
					Field: "a",
					Asc:   true,
				},
				{
					Field: "c",
					Asc:   true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getOrderFromQuery(tt.args.queryParams)
			assert.True(t, reflect.DeepEqual(tt.want, got), "got: %v, want: %v", got, tt.want)
		})
	}
}

func TestHasOneLetter(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should return true when string has one letter",
			args: args{s: "a"},
			want: true,
		},
		{
			name: "should return true when string has one letter",
			args: args{s: "ab"},
			want: true,
		},
		{
			name: "should return true when string has one letter",
			args: args{s: "a1"},
			want: true,
		},
		{
			name: "should return true when string has one letter",
			args: args{s: "1a"},
			want: true,
		},
		{
			name: "should return true when string has one letter",
			args: args{s: "1a1"},
			want: true,
		},
		{
			name: "should return false when string has no letters",
			args: args{s: "1"},
			want: false,
		},
		{
			name: "should return false when string has no letters",
			args: args{s: "11"},
			want: false,
		},
		{
			name: "should return false when string has no letters",
			args: args{s: ""},
			want: false,
		},
		{
			name: "should return false when string has no letters",
			args: args{s: " "},
			want: false,
		},
		{
			name: "should return false when string has no letters",
			args: args{s: "      "},
			want: false,
		},
		{
			name: "should return false when string has no letters",
			args: args{s: "."},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasALetter(tt.args.s)
			assert.Equal(t, tt.want, got, "got: %v, want: %v", got, tt.want)
		})
	}
}

func TestFilterOperatorIsValid(t *testing.T) {
	type args struct {
		operator string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should return true when operator is valid",
			args: args{operator: "eq"},
			want: true,
		},
		{
			name: "should return true when operator is valid",
			args: args{operator: "contains"},
			want: true,
		},
		{
			name: "should return true when operator is valid",
			args: args{operator: "gt"},
			want: true,
		},
		{
			name: "should return false when operator is invalid",
			args: args{operator: "anything else"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Filter{
				Field:     "a",
				Operation: FilterOperator(tt.args.operator),
				Value:     "b",
			}
			got := f.Operation.IsValid()
			assert.Equal(t, tt.want, got, "got: %v, want: %v", got, tt.want)
		})
	}
}

func TestGetFilterFromQuery(t *testing.T) {
	type args struct {
		queryParams map[string]string
	}
	tests := []struct {
		name string
		args args
		want []Filter
	}{
		{
			name: "should return empty slice when query params has no filters",
			args: args{queryParams: map[string]string{
				"limit":  "10",
				"offset": "0",
			}},
			want: []Filter{},
		},
		{
			name: "should return empty slice when query params has invalid filters",
			args: args{queryParams: map[string]string{
				"limit":   "10",
				"offset":  "0",
				"filters": "a",
			}},
			want: []Filter{},
		},
		{
			name: "should return empty slice when query params has invalid filters",
			args: args{queryParams: map[string]string{
				"limit":   "10",
				"offset":  "0",
				"filters": "a:asc",
			}},
			want: []Filter{},
		},
		{
			name: "should return empty slice when query params has invalid filters",
			args: args{queryParams: map[string]string{
				"limit":   "10",
				"offset":  "0",
				"filters": "a:asc,b:desc",
			}},
			want: []Filter{},
		},
		{
			name: "should return empty slice when query params has invalid filters",
			args: args{queryParams: map[string]string{
				"limit":   "10",
				"offset":  "0",
				"filters": "a[eqb",
			}},
			want: []Filter{},
		},
		{
			name: "should return empty slice when query params has invalid filters",
			args: args{queryParams: map[string]string{
				"limit":   "10",
				"offset":  "0",
				"filters": "aeq]b",
			}},
			want: []Filter{},
		},
		{
			name: "should return Filter slice when query params has valid filters",
			args: args{queryParams: map[string]string{
				"limit":   "10",
				"offset":  "0",
				"filters": "a[eq]b",
			}},
			want: []Filter{
				{
					Field:     "a",
					Operation: FilterOperator("eq"),
					Value:     "b",
				},
			},
		},
		{
			name: "should return Filter slice when query params has valid filters",
			args: args{queryParams: map[string]string{
				"limit":   "10",
				"offset":  "0",
				"filters": "a[contains]b",
			}},
			want: []Filter{
				{
					Field:     "a",
					Operation: FilterOperator("contains"),
					Value:     "b",
				},
			},
		},
		{
			name: "should return Filter slice with omitted values when query params has valid and invalid filters",
			args: args{queryParams: map[string]string{
				"limit":   "10",
				"offset":  "0",
				"filters": "a[eq]b,c[]d,e[eq]f",
			}},
			want: []Filter{
				{
					Field:     "a",
					Operation: FilterOperator("eq"),
					Value:     "b",
				},
				{
					Field:     "e",
					Operation: FilterOperator("eq"),
					Value:     "f",
				},
			},
		},
		{
			name: "should return Filter slice with omitted values when query params has valid and invalid filters",
			args: args{queryParams: map[string]string{
				"limit":   "10",
				"offset":  "0",
				"filters": "a[eq]b,[eq]d,e[contains]f",
			}},
			want: []Filter{
				{
					Field:     "a",
					Operation: FilterOperator("eq"),
					Value:     "b",
				},
				{
					Field:     "e",
					Operation: FilterOperator("contains"),
					Value:     "f",
				},
			},
		},
		{
			name: "should return Filter slice with omitted values when query params has valid and invalid filters",
			args: args{queryParams: map[string]string{
				"limit":   "10",
				"offset":  "0",
				"filters": "a[eq]b,t[eq],e[contains]f",
			}},
			want: []Filter{
				{
					Field:     "a",
					Operation: FilterOperator("eq"),
					Value:     "b",
				},
				{
					Field:     "e",
					Operation: FilterOperator("contains"),
					Value:     "f",
				},
			},
		},
		{
			name: "should return Filter slice with omitted values when query params has valid and invalid filters",
			args: args{queryParams: map[string]string{
				"limit":   "10",
				"offset":  "0",
				"filters": "a[eq]b,[eq],e[contains]f",
			}},
			want: []Filter{
				{
					Field:     "a",
					Operation: FilterOperator("eq"),
					Value:     "b",
				},
				{
					Field:     "e",
					Operation: FilterOperator("contains"),
					Value:     "f",
				},
			},
		},
		{
			name: "should return Filter slice when query params has valid filters",
			args: args{queryParams: map[string]string{
				"limit":   "10",
				"offset":  "0",
				"filters": "a[eq]b;c[eq]d;e[eq]f",
			}},
			want: []Filter{
				{
					Field:     "a",
					Operation: FilterOperator("eq"),
					Value:     "b;c[eq]d;e[eq]f",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getFilterFromQuery(tt.args.queryParams)
			assert.True(t, reflect.DeepEqual(tt.want, got), "got: %v, want: %v", got, tt.want)
		})
	}
}
