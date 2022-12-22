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
