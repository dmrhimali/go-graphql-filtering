package graph

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

type Date string

var ErrUnexpectedValue = errors.New("unexpected value")

// MarshalDate marshals time.Time to YYYY-MM-DD
func MarshalDate(t time.Time) graphql.Marshaler {
	if t.IsZero() {
		return graphql.Null
	}

	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, strconv.Quote(t.Format("2006-01-02")))
	})
}

// UnmarshalDate unmarshalls YYYY-MM-DD to time.Time
func UnmarshalDate(v interface{}) (time.Time, error) {
	if tmpStr, ok := v.(string); ok {
		if len(tmpStr) == 0 {
			return time.Time{}, nil
		}

		parse, err := time.Parse(time.DateOnly, tmpStr)
		if err != nil {
			return time.Time{}, fmt.Errorf("date must be in format %q: %w", time.DateOnly, err)
		}

		return parse, nil
	}

	return time.Time{}, fmt.Errorf("%w: date should be a string", ErrUnexpectedValue)
}
