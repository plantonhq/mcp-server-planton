package domains

import (
	"fmt"
	"sort"
	"strings"
)

// EnumResolver maps user-supplied strings to protobuf enum values of type T.
// T must be a proto enum type whose underlying representation is int32.
type EnumResolver[T ~int32] struct {
	values     map[string]int32
	typeName   string
	excludeKey string
}

// NewEnumResolver creates a resolver for proto enum type T.
//   - values:     the proto-generated _value map (e.g. MyEnum_value)
//   - typeName:   human-readable label used in error messages (e.g. "execution status")
//   - excludeKey: the zero-value sentinel to exclude from "valid values" listings
func NewEnumResolver[T ~int32](values map[string]int32, typeName, excludeKey string) EnumResolver[T] {
	return EnumResolver[T]{values: values, typeName: typeName, excludeKey: excludeKey}
}

// Resolve maps a single string to the corresponding enum value.
func (r EnumResolver[T]) Resolve(s string) (T, error) {
	v, ok := r.values[s]
	if !ok {
		return 0, fmt.Errorf("unknown %s %q — valid values: %s",
			r.typeName, s, JoinEnumValues(r.values, r.excludeKey))
	}
	return T(v), nil
}

// ResolveSlice maps a slice of strings to the corresponding enum values.
// Returns all resolved values, or an error on the first unknown string.
func (r EnumResolver[T]) ResolveSlice(ss []string) ([]T, error) {
	out := make([]T, 0, len(ss))
	for _, s := range ss {
		v, err := r.Resolve(s)
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, nil
}

// JoinEnumValues returns a sorted, comma-separated list of the map's keys,
// excluding the specified zero-value key (e.g. "unspecified" sentinel).
func JoinEnumValues(m map[string]int32, exclude string) string {
	vals := make([]string, 0, len(m)-1)
	for k := range m {
		if k != exclude {
			vals = append(vals, k)
		}
	}
	sort.Strings(vals)
	return strings.Join(vals, ", ")
}
