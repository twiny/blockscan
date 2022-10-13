package api

import (
	"testing"
)

func TestParseScanQuery(t *testing.T) {
	idx, err := NewIndexer("config/config.yaml")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name  string
		query string
		want  []any
	}{
		{
			name:  "TestValidRange",
			query: "100:200",
			want:  []any{100, 200, nil},
		},
		{
			name:  "TestBadRange",
			query: "x:y",
			// want:  []any{0, 0, err}, // invalid range error
		},
		{
			// add test for range from id to latest
		},
	}

	for _, tc := range tests {
		start, end, err := idx.parseScanQuery(tc.query)

		s, ok := tc.want[0].(int64)
		if !ok {
			//
		}
		if start != s {
			//
		}

		//

		e, ok := tc.want[1].(int64)
		if !ok {
			//
		}
		if end != e {
			//
		}

		er, ok := tc.want[2].(error)
		if !ok {
			//
		}
		if err != er {
			//
		}
	}

}
