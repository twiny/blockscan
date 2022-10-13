package api

import (
	"fmt"
	"strconv"
	"strings"
)

// parseScanQuery
// 100 200 => range
// 100 - the latest => range once reached subscribe to new
// empty => 0:latest
func (idx *Indexer) parseScanQuery(s string) (start int64, end int64, err error) {
	parts := strings.Split(s, ":")

	switch {
	case len(s) == 0:
		start = 0
		end = idx.head
		return
	case len(parts) == 1:
		start, err = strconv.ParseInt(parts[0], 10, 0)
		if err != nil {
			return
		}

		end = idx.head
		return

	case len(parts) == 2:
		start, err = strconv.ParseInt(parts[0], 10, 0)
		if err != nil {
			return
		}

		end, err = strconv.ParseInt(parts[1], 10, 0)
		if err != nil {
			return
		}

		return
	default:
		err = fmt.Errorf("range must be in format start:end")
		return
	}
}
